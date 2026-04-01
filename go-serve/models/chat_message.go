package models

import (
	"time"
)

// ChatMessage 表示用户和 AI 的聊天消息
type ChatMessage struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      string     `json:"role"` // "user" 或 "assistant"
	Content   string     `gorm:"type:text" json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ChatRequest 请求体结构
type ChatRequest struct {
	Message string `json:"message" binding:"required,min=1,max=2000"`
}

// ChatResponse 响应体结构
type ChatResponse struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

// ChatHistoryResponse 聊天历史响应
type ChatHistoryResponse struct {
	ID      uint   `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// AliYunResponse 阿里云千问 API 响应
type AliYunResponse struct {
	Output struct {
		Text string `json:"text"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}
