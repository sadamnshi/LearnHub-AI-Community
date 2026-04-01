package services

import (
	"context"
	"fmt"
	"gin_demo/models"
	"gin_demo/repositories"
	"log"
	"os"
	"time"

	// 使用官方 OpenAI SDK，兼容 DashScope 等 OpenAI-Compatible 服务
	"github.com/sashabaranov/go-openai"
)

// 大模型配置
const (
	OpenAIModel   = "qwen3.5-plus"
	
)

// OpenAIAPIKey 从环境变量获取 API Key
var OpenAIAPIKey = func() string {
	if key := os.Getenv("ALIYUN_API_KEY"); key != "" {
		return key
	}
	// 开发环境默认值（示例密钥，实际使用请替换为真实 Key）
	return "sk-93ce5b4b6eb24050890898cd5849c84b"
}()

// OpenAIAPIBaseURL 从环境变量获取 API 基础地址，未设置则使用默认地址
var OpenAIAPIBaseURL = func() string {
	if baseURL := os.Getenv("ALIYUN_BASE_URL"); baseURL != "" {
		return baseURL
	}
	return "https://dashscope.aliyuncs.com/compatible-mode/v1"
}()

// AIService 定义 AI 聊天服务接口
type AIService interface {
	// SendMessage 发送消息并获取 AI 回复
	SendMessage(userID uint, message string) (string, error)

	// GetChatHistory 获取聊天历史
	GetChatHistory(userID uint, limit int) ([]models.ChatHistoryResponse, error)

	// ClearChatHistory 清空聊天历史
	ClearChatHistory(userID uint) error
}

type aiServiceImpl struct {
	chatRepo repositories.ChatRepository
	client   *openai.Client
}

// NewAIService 创建 AI 服务实例
func NewAIService(chatRepo repositories.ChatRepository) AIService {
	// 使用 OpenAI SDK 创建客户端（兼容 DashScope 等 OpenAI-Compatible 服务）
	config := openai.DefaultConfig(OpenAIAPIKey)
	config.BaseURL = OpenAIAPIBaseURL
	client := openai.NewClientWithConfig(config)

	return &aiServiceImpl{
		chatRepo: chatRepo,
		client:   client,
	}
}

// SendMessage 发送消息给 AI 并获取回复
func (s *aiServiceImpl) SendMessage(userID uint, message string) (string, error) {
	// 📝 1. 保存用户消息到数据库
	userMsg := &models.ChatMessage{
		UserID:  userID,
		Role:    "user",
		Content: message,
	}

	if err := s.chatRepo.SaveMessage(userMsg); err != nil {
		log.Printf("Failed to save user message: %v", err)
		return "", fmt.Errorf("failed to save message: %w", err)
	}

	// 🤖 2. 调用 OpenAI API 获取回复
	aiReply, err := s.callOpenAIAPI(message)
	if err != nil {
		log.Printf("Failed to call OpenAI API: %v", err)
		return "", fmt.Errorf("failed to get AI response: %w", err)
	}

	// 💾 3. 保存 AI 回复到数据库
	aiMsg := &models.ChatMessage{
		UserID:  userID,
		Role:    "assistant",
		Content: aiReply,
	}

	if err := s.chatRepo.SaveMessage(aiMsg); err != nil {
		log.Printf("Failed to save AI message: %v", err)
		// 即使保存失败也返回 AI 回复，因为 API 调用已经成功
	}

	return aiReply, nil
}

// callOpenAIAPI 调用 OpenAI API
// 使用 sashabaranov/go-openai SDK，简化代码并提高稳定性
func (s *aiServiceImpl) callOpenAIAPI(message string) (string, error) {
	// ⏱️ 创建带超时的上下文（30秒）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 🤖 调用 OpenAI Chat Completion API
	// 这是 go-openai SDK 的标准用法，自动处理：
	// - 请求头（Content-Type, Authorization）
	// - 请求体编码和响应解析
	// - 错误处理
	response, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		// 可选参数
		MaxTokens:   2000, // 限制回复长度
		Temperature: 0.7,  // 控制创意程度（0-2之间，0表示确定性，2表示随机性强）
	})

	if err != nil {
		// SDK 会返回详细的错误信息
		log.Printf("❌ OpenAI API error: %v", err)
		return "", fmt.Errorf("openai api error: %w", err)
	}

	// ✅ 提取回复内容
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI API")
	}

	reply := response.Choices[0].Message.Content

	// 📊 记录 token 使用情况
	log.Printf("✅ OpenAI API success - Input tokens: %d, Output tokens: %d, Total: %d",
		response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens)

	return reply, nil
}

// GetChatHistory 获取聊天历史记录
func (s *aiServiceImpl) GetChatHistory(userID uint, limit int) ([]models.ChatHistoryResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20 // 默认返回 20 条
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
