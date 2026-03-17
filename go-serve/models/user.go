package models

import (
	"gorm.io/gorm"
)

// User 领域模型：对应 users 表。
// 企业级改进说明：
// 1. Username 设置唯一索引，防止重复注册。
// 2. Password 字段在入库前会被服务层使用 bcrypt 进行哈希，不存明文。
// 3. 生产环境通常还会添加 Email / Mobile / Status / LastLoginAt 等字段，这里保持最小示例。
// 4. 若需要软删除能力，gorm.Model 已包含 DeletedAt。
// 5. 通过 form / json tag 使其能同时绑定表单与 JSON 请求，binding 标签用于基础的必填校验。
// 注意：Password 在请求阶段接收的是明文，入库前会替换成哈希值。
type User struct {
	gorm.Model
	Username string  `gorm:"uniqueIndex;size:64" form:"username" json:"username" binding:"required"`
	Password string  `gorm:"size:255" form:"password" json:"password" binding:"required"` // 保存 bcrypt 哈希
	Email    *string `gorm:"uniqueIndex;size:128"`                                        // 新增
	Avatar   string  `gorm:"size:512"`                                                    // 头像 URL
	Bio      string  `gorm:"size:500"`                                                    // 个人简介
	Role     string  `gorm:"default:'user'"`                                              // user | admin | moderator
}
