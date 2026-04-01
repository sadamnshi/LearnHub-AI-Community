package controllers

import (
	"gin_demo/middleware"
	"gin_demo/models"
	"gin_demo/services"
	"gin_demo/util"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiService services.AIService
}

// NewAIController 创建 AI 控制器
func NewAIController(aiService services.AIService) *AIController {
	return &AIController{
		aiService: aiService,
	}
}

// SendMessage 处理用户消息请求
// POST /api/chat/send
func (c *AIController) SendMessage(ctx *gin.Context) {
	// 👤 获取用户身份
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		util.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized", "请先登录")
		return
	}

	// 📝 解析请求体
	var req models.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(ctx, http.StatusBadRequest, "invalid_request", "消息内容不能为空，且长度在 1-2000 之间")
		return
	}

	// 🤖 调用 AI 服务获取回复
	reply, err := c.aiService.SendMessage(userID, req.Message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		util.ErrorResponse(ctx, http.StatusInternalServerError, "error", "获取 AI 回复失败")
		return
	}

	// ✅ 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": models.ChatResponse{
			Message: reply,
			Role:    "assistant",
		},
	})
}

// GetHistory 获取聊天历史记录
// GET /api/chat/history?limit=20
func (c *AIController) GetHistory(ctx *gin.Context) {
	// 👤 获取用户身份
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		util.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized", "请先登录")
		return
	}

	// 📊 解析分页参数
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	// 📜 获取聊天历史
	history, err := c.aiService.GetChatHistory(userID, limit)
	if err != nil {
		log.Printf("Failed to get chat history: %v", err)
		util.ErrorResponse(ctx, http.StatusInternalServerError, "error", "获取聊天历史失败")
		return
	}

	// ✅ 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": history,
	})
}

// ClearHistory 清空聊天历史记录
// DELETE /api/chat/history
func (c *AIController) ClearHistory(ctx *gin.Context) {
	// 👤 获取用户身份
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		util.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized", "请先登录")
		return
	}

	// 🗑️ 清空聊天历史
	err = c.aiService.ClearChatHistory(userID)
	if err != nil {
		log.Printf("Failed to clear chat history: %v", err)
		util.ErrorResponse(ctx, http.StatusInternalServerError, "error", "清空聊天历史失败")
		return
	}

	// ✅ 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": nil,
	})
}
