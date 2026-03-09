package models

import (
	"gorm.io/gorm"
)

// Post 帖子模型
type Post struct {
	gorm.Model
	Title        string   `gorm:"size:200;not null"`
	Content      string   `gorm:"type:text;not null"` // Markdown 格式
	AuthorID     uint     `gorm:"not null;index"`
	Author       User     `gorm:"foreignKey:AuthorID"`
	CategoryID   uint     `gorm:"index"`
	Category     Category `gorm:"foreignKey:CategoryID"`
	Tags         []Tag    `gorm:"many2many:post_tags;"`
	ViewCount    int      `gorm:"default:0"`
	LikeCount    int      `gorm:"default:0"`
	CommentCount int      `gorm:"default:0"`
	IsPinned     bool     `gorm:"default:false"`       // 置顶
	Status       string   `gorm:"default:'published'"` // draft | published | banned
}

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;size:64"`
	Description string `gorm:"size:255"`
	Icon        string `gorm:"size:128"` // emoji 或图标名
	PostCount   int    `gorm:"default:0"`
}

// Tag 标签模型
type Tag struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex;size:64"`
	PostCount int    `gorm:"default:0"`
}
