package repositories

import (
	"gin_demo/databases"
	"gin_demo/models"
)

type ChatRepository interface {
	SaveMessage(message *models.ChatMessage) error

	GetMessagesByUserID(userID uint, limit int) ([]models.ChatMessage, error)

	// GetRecentMessagesForContext 获取用户最近的 N 条消息作为对话上下文
	// 用于构建多轮对话时的历史记录上下文
	GetRecentMessagesForContext(userID uint, contextLimit int) ([]models.ChatMessage, error)

	// DeleteMessagesByUserID 删除用户的所有聊天记录
	DeleteMessagesByUserID(userID uint) error
}

type chatRepository struct{}

// NewChatRepository 创建聊天消息仓库实例
func NewChatRepository() ChatRepository {
	return &chatRepository{}
}

// SaveMessage 保存单条聊天消息
func (r *chatRepository) SaveMessage(message *models.ChatMessage) error {
	return databases.DB.Create(message).Error
}

// GetMessagesByUserID 获取用户的聊天历史
func (r *chatRepository) GetMessagesByUserID(userID uint, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage

	err := databases.DB.
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

// GetRecentMessagesForContext 获取用户最近的 N 条消息作为对话上下文
// 按时间升序返回，便于构建完整的对话历史
func (r *chatRepository) GetRecentMessagesForContext(userID uint, contextLimit int) ([]models.ChatMessage, error) {
	if contextLimit <= 0 {
		contextLimit = 10 // 默认获取最近 10 条消息（5 轮对话）
	}
	if contextLimit > 50 {
		contextLimit = 50 // 最多 50 条消息，避免 Token 过多
	}

	var messages []models.ChatMessage

	// 先按 created_at 降序获取最新的 contextLimit 条记录
	// 然后再按 created_at 升序排序，保证对话顺序正确
	err := databases.DB.
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(contextLimit).
		Order("created_at ASC").
		Find(&messages).Error

	return messages, err
}

// DeleteMessagesByUserID 删除用户的所有消息
func (r *chatRepository) DeleteMessagesByUserID(userID uint) error {
	// 使用逻辑删除（软删除）
	return databases.DB.Where("user_id = ?", userID).Delete(&models.ChatMessage{}).Error
}
