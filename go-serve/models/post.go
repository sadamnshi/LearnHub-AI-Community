package models

import (
	"gorm.io/gorm"
)

// Post 帖子模型
type Post struct {
	gorm.Model
	Title      string `gorm:"size:200;not null"`
	Content    string `gorm:"type:text;not null"` // Markdown 格式
	AuthorID   uint   `gorm:"not null;index"`
	Author     User   `gorm:"foreignKey:AuthorID"`
	CategoryID uint   `gorm:"index"`
	Category   Category `gorm:"foreignKey:CategoryID"`
	Tags       string `gorm:"size:200"` // 逗号分隔的标签字符串，如 "golang,database,缓存"
	Status     string `gorm:"default:'published'"` // draft | published | banned
}

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;size:64"`
	Description string `gorm:"size:255"`
	Icon        string `gorm:"size:128"` // emoji 或图标名
}
