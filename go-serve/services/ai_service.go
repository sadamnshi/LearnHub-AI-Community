package services

import (
	"context"
	"fmt"
	"gin_demo/models"
	"gin_demo/repositories"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	EnvAliyunAPIKey  = "ALIYUN_API_KEY"
	EnvAliyunBaseURL = "ALIYUN_BASE_URL"
	EnvOpenAIModel   = "OpenAIModel"

	// maxConcurrentAPICalls 最多允许同时进行多少个 AI API 调用。
	// 原理：AI API 提供商（如阿里云）都有并发限制，超出会被限速（429 Too Many Requests）。
	// 使用带缓冲的 channel 作为"令牌桶"：想调用 API 必须先拿到令牌，用完后归还。
	// 同时只有 maxConcurrentAPICalls 个令牌，所以最多 maxConcurrentAPICalls 个并发。
	maxConcurrentAPICalls = 5

	// asyncSaveMaxRetries 异步落库的最大重试次数
	asyncSaveMaxRetries = 3

	// asyncSaveBaseDelay 重试的基础延迟（指数退避: 500ms, 1000ms, 2000ms）
	asyncSaveBaseDelay = 500 * time.Millisecond
)

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("关键环境变量缺失: %s", key)
	}
	return value
}

// GetOpenAIAPIKey 获取 OpenAI API Key（延迟初始化，避免启动时加载）
func GetOpenAIAPIKey() string {
	return mustGetEnv(EnvAliyunAPIKey)
}

// GetOpenAIAPIBaseURL 获取 OpenAI API 基础地址（延迟初始化，避免启动时加载）
func GetOpenAIAPIBaseURL() string {
	return mustGetEnv(EnvAliyunBaseURL)
}

// GetOpenAIModel 获取 OpenAI 模型名称（延迟初始化，避免启动时加载）
func GetOpenAIModel() string {
	return mustGetEnv(EnvOpenAIModel)
}

// AIService 定义 AI 聊天服务接口
type AIService interface {
	// SendMessage 发送消息并获取 AI 回复（支持多轮对话上下文）
	SendMessage(userID uint, message string) (string, error)

	// GetChatHistory 获取聊天历史
	GetChatHistory(userID uint, limit int) ([]models.ChatHistoryResponse, error)

	// ClearChatHistory 清空聊天历史
	ClearChatHistory(userID uint) error
}

// aiServiceImpl AI 服务实现
//
// 新增字段 semaphore：带缓冲的 channel，用作"并发信号量"（Semaphore）。
//
// Semaphore 核心原理：
//
//	把一个容量为 N 的 buffered channel 当成 N 个"令牌"。
//	  - 获取令牌（acquire）：向 channel 发送一个空结构体 -> semaphore <- struct{}{}
//	                         如果 channel 已满（N 个令牌全被占用），当前 goroutine 会阻塞等待
//	  - 释放令牌（release）：从 channel 接收一个空结构体 <- semaphore
//	                         腾出一个槽，让等待中的 goroutine 得以继续
//
//	这样就实现了：同一时刻最多 N 个 goroutine 能进入受保护的代码区（AI API 调用）。
type aiServiceImpl struct {
	chatRepo  repositories.ChatRepository
	client    *openai.Client
	semaphore chan struct{} // 并发信号量：限制同时进行的 AI API 调用数
}

// NewAIService 创建 AI 服务实例
func NewAIService(chatRepo repositories.ChatRepository) AIService {
	config := openai.DefaultConfig(GetOpenAIAPIKey())
	config.BaseURL = GetOpenAIAPIBaseURL()
	client := openai.NewClientWithConfig(config)

	return &aiServiceImpl{
		chatRepo: chatRepo,
		client:   client,
		// 创建容量为 maxConcurrentAPICalls 的 buffered channel 作为信号量
		// 初始时 channel 为空（没有值），可以接受 maxConcurrentAPICalls 次 send 操作
		semaphore: make(chan struct{}, maxConcurrentAPICalls),
	}
}

// historyResult 用于在 goroutine 与主流程之间传递历史记录查询结果
type historyResult struct {
	messages []models.ChatMessage
	err      error
}

// SendMessage 发送消息给 AI 并获取回复（支持多轮对话）
//
// ⚡ 进阶三阶段并发流水线（v2）：
//
//	在原有三阶段流水线的基础上，新增两项 Goroutine 进阶技术：
//
//	┌──────────────────────────────────────────────────────────────────────┐
//	│ 阶段一（WaitGroup 并行）                                             │
//	│   goroutine-1: 保存用户消息(DB写, ~20ms) ──┐                        │
//	│   goroutine-2: 获取历史记录(DB读, ~30ms) ──┘ 并行，耗时≈30ms        │
//	├──────────────────────────────────────────────────────────────────────┤
//	│ 阶段二（Semaphore 限流 + 串行 AI API 调用）          ← 新增 ⭐       │
//	│   先获取"并发令牌"（若令牌已被占满则等待），防止 API 限速             │
//	│   → 调用 AI API(~1000ms) → 释放令牌                                 │
//	├──────────────────────────────────────────────────────────────────────┤
//	│ 阶段三（异步重试落库，fire-and-forget goroutine）    ← 升级 ⭐       │
//	│   goroutine-3 内部加入指数退避重试（最多3次）                        │
//	│   第1次失败 → 等500ms → 重试                                         │
//	│   第2次失败 → 等1000ms → 重试                                        │
//	│   第3次失败 → 记录日志，放弃（可接受的降级）                         │
//	└──────────────────────────────────────────────────────────────────────┘
func (s *aiServiceImpl) SendMessage(userID uint, message string) (string, error) {
	// ══════════════════════════════════════════════════════════════════════
	// 阶段一：并行执行「保存用户消息」和「获取历史记录」（WaitGroup）
	// ══════════════════════════════════════════════════════════════════════
	var (
		saveErr error
		history historyResult
		wg      sync.WaitGroup
	)

	wg.Add(2)

	// goroutine-1：异步保存用户消息到数据库（DB 写操作）
	go func() {
		defer wg.Done()
		userMsg := &models.ChatMessage{
			UserID:  userID,
			Role:    "user",
			Content: message,
		}
		saveErr = s.chatRepo.SaveMessage(userMsg)
	}()

	// goroutine-2：异步获取历史消息（DB 读操作，用于构建 AI 上下文）
	go func() {
		defer wg.Done()
		history.messages, history.err = s.chatRepo.GetRecentMessagesForContext(userID, 20)
	}()

	wg.Wait()

	// "保存用户消息"失败 → 关键错误，中止本次请求
	if saveErr != nil {
		log.Printf("Failed to save user message: %v", saveErr)
		return "", fmt.Errorf("failed to save message: %w", saveErr)
	}

	// "获取历史记录"失败 → 非关键错误，继续执行（AI 丢失上下文但仍能回答）
	if history.err != nil {
		log.Printf("Failed to get chat history: %v", history.err)
	}

	// 由于阶段一两个操作并行发起，历史查询可能不含当前用户消息，手动追加
	currentMsg := models.ChatMessage{UserID: userID, Role: "user", Content: message}
	historyWithCurrent := append(history.messages, currentMsg)

	// ══════════════════════════════════════════════════════════════════════
	// 阶段二：Semaphore 限流 → 调用 AI API
	// ══════════════════════════════════════════════════════════════════════
	// 【新增】在调用 AI API 之前，先向信号量 channel 发送一个空结构体（"获取令牌"）。
	// 如果当前已有 maxConcurrentAPICalls 个 goroutine 在调用 AI API，
	// 这一步会阻塞，直到某个 goroutine 退出并"归还令牌"。
	// 这就像排队取号：号码（令牌）就这么多，拿到号才能进服务区。
	log.Printf("⏳ 等待 AI API 并发令牌... (当前信号量使用数: %d/%d)", len(s.semaphore), maxConcurrentAPICalls)
	s.semaphore <- struct{}{} // 获取令牌（acquire），channel 满时此行阻塞

	// 用 defer 确保函数退出时一定"归还令牌"，即使 AI API 调用 panic 也不会死锁
	defer func() { <-s.semaphore }() // 释放令牌（release），从 channel 取出一个值腾出槽位

	log.Printf("✅ 已获取 AI API 并发令牌，开始调用 (当前信号量使用数: %d/%d)", len(s.semaphore), maxConcurrentAPICalls)

	aiReply, err := s.callOpenAIAPIWithContext(historyWithCurrent)
	if err != nil {
		log.Printf("Failed to call OpenAI API: %v", err)
		return "", fmt.Errorf("failed to get AI response: %w", err)
	}

	// ══════════════════════════════════════════════════════════════════════
	// 阶段三：异步落库 + 指数退避重试（Exponential Backoff Retry）
	// ══════════════════════════════════════════════════════════════════════
	// 【升级】原来是简单的 fire-and-forget，失败只打日志。
	// 现在加入指数退避重试：失败后等待一段时间再重试，等待时间指数增长。
	//
	// 指数退避（Exponential Backoff）的原理：
	//   第 1 次重试：等 500ms  (baseDelay * 2^0)
	//   第 2 次重试：等 1000ms (baseDelay * 2^1)
	//   第 3 次重试：等 2000ms (baseDelay * 2^2)
	//   → 这样做避免了在数据库刚恢复时被大量重试请求瞬间淹没（避免"惊群效应"）
	//
	// 为什么用 goroutine？→ 用户不需要等待这步，AI 回复已经拿到可以立即返回
	go func() {
		aiMsg := &models.ChatMessage{
			UserID:  userID,
			Role:    "assistant",
			Content: aiReply,
		}
		s.saveWithRetry(aiMsg)
	}()

	return aiReply, nil
}

// saveWithRetry 带指数退避重试的消息落库
//
// 这是阶段三的核心实现，在后台 goroutine 中安全运行。
//
// 指数退避时间序列（baseDelay = 500ms）：
//
//	attempt=1: sleep(500ms  * 2^0) = 500ms
//	attempt=2: sleep(500ms  * 2^1) = 1000ms
//	attempt=3: sleep(500ms  * 2^2) = 2000ms（最后一次，失败后放弃）
func (s *aiServiceImpl) saveWithRetry(msg *models.ChatMessage) {
	var lastErr error

	for attempt := 1; attempt <= asyncSaveMaxRetries; attempt++ {
		err := s.chatRepo.SaveMessage(msg)
		if err == nil {
			// 落库成功
			if attempt > 1 {
				log.Printf("✅ 异步落库重试成功 (第 %d 次尝试, role=%s, userID=%d)", attempt, msg.Role, msg.UserID)
			}
			return
		}

		lastErr = err
		log.Printf("⚠️  异步落库失败 (第 %d/%d 次, role=%s, userID=%d): %v",
			attempt, asyncSaveMaxRetries, msg.Role, msg.UserID, err)

		// 不是最后一次尝试，则等待后重试
		if attempt < asyncSaveMaxRetries {
			// 指数退避：等待时间 = baseDelay * 2^(attempt-1)
			// attempt=1 → 500ms, attempt=2 → 1000ms
			backoff := asyncSaveBaseDelay * (1 << uint(attempt-1))
			log.Printf("⏳ 等待 %v 后进行第 %d 次重试...", backoff, attempt+1)
			time.Sleep(backoff)
		}
	}

	// 全部重试耗尽，记录告警日志（实际生产中可以接入告警系统）
	log.Printf("❌ 异步落库彻底失败，已重试 %d 次放弃 (role=%s, userID=%d): %v",
		asyncSaveMaxRetries, msg.Role, msg.UserID, lastErr)
}

// callOpenAIAPIWithContext 调用 OpenAI API，并传入聊天历史作为上下文
func (s *aiServiceImpl) callOpenAIAPIWithContext(historyMessages []models.ChatMessage) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `你是一个有用的AI助手。请根据对话历史提供连贯、自然的回复。
如果用户的问题与之前的对话相关，请参考之前的内容来理解上下文。`,
		},
	}

	for _, msg := range historyMessages {
		var role string
		if msg.Role == "user" {
			role = openai.ChatMessageRoleUser
		} else if msg.Role == "assistant" {
			role = openai.ChatMessageRoleAssistant
		} else {
			continue
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	log.Printf("📤 Sending %d messages to OpenAI API (including %d history messages)",
		len(messages), len(historyMessages))

	response, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       GetOpenAIModel(),
		Messages:    messages,
		MaxTokens:   2000,
		Temperature: 0.7,
	})
	if err != nil {
		log.Printf("❌ OpenAI API error: %v", err)
		return "", fmt.Errorf("openai api error: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI API")
	}

	reply := response.Choices[0].Message.Content

	log.Printf("✅ OpenAI API success - Input tokens: %d, Output tokens: %d, Total: %d, Context size: %d messages",
		response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens, len(historyMessages))

	return reply, nil
}

// GetChatHistory 获取聊天历史记录
func (s *aiServiceImpl) GetChatHistory(userID uint, limit int) ([]models.ChatHistoryResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	messages, err := s.chatRepo.GetMessagesByUserID(userID, limit)
	if err != nil {
		log.Printf("Failed to get chat history: %v", err)
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	var history []models.ChatHistoryResponse
	for _, msg := range messages {
		history = append(history, models.ChatHistoryResponse{
			ID:      msg.ID,
			Role:    msg.Role,
			Content: msg.Content,
			Time:    msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return history, nil
}

// ClearChatHistory 清空用户的聊天历史
func (s *aiServiceImpl) ClearChatHistory(userID uint) error {
	err := s.chatRepo.DeleteMessagesByUserID(userID)
	if err != nil {
		log.Printf("Failed to clear chat history for user %d: %v", userID, err)
		return fmt.Errorf("failed to clear chat history: %w", err)
	}

	log.Printf("✅ Chat history cleared for user %d", userID)
	return nil
}
