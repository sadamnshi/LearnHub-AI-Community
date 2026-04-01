package repositories

import (
	"gin_demo/databases"
	"gin_demo/models"
)

// ChatRepository 定义聊天消息的数据操作接口
type ChatRepository interface {
	// SaveMessage 保存聊天消息（用户或 AI 的消息）
	SaveMessage(message *models.ChatMessage) error

	// GetMessagesByUserID 获取用户的聊天历史记录
	// 参数：
	//   - userID: 用户 ID
	//   - limit: 返回的最多消息数量
	// 返回按时间升序排列的消息列表
	GetMessagesByUserID(userID uint, limit int) ([]models.ChatMessage, error)

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

// DeleteMessagesByUserID 删除用户的所有消息
func (r *chatRepository) DeleteMessagesByUserID(userID uint) error {
	// 使用逻辑删除（软删除）
	return databases.DB.Where("user_id = ?", userID).Delete(&models.ChatMessage{}).Error
}
