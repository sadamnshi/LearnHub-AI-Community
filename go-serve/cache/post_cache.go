package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_demo/databases"
	"gin_demo/repositories"
	"gin_demo/services"
	"log"
	"strings"
	"time"
)

const (
	KeyFormat  = "post:%d:detail"      // %d 是占位符，表示帖子 ID
	Expiration = 10 * time.Minute      // 缓存过期时间
	TimeFormat = "2006-01-02 15:04:05" // 时间格式化字符串
)

type PostCache interface {
	// GetPostDetail 从缓存获取帖子详情（缓存未命中则查询数据库并缓存）
	GetPostDetail(postID uint) (*services.PostDetail, error)

	// InvalidatePostDetail 删除帖子详情缓存（编辑后需要清空）
	InvalidatePostDetail(postID uint) error
}

type postCacheImpl struct {
	repo repositories.PostRepository
}

// NewPostCache 创建缓存实例
func NewPostCache(repo repositories.PostRepository) PostCache {
	return &postCacheImpl{
		repo: repo,
	}
}

// 🎯 改进建议 #2：参数验证
func (c *postCacheImpl) GetPostDetail(postID uint) (*services.PostDetail, error) {
	// 参数验证：postID 必须是正整数
	if postID == 0 {
		return nil, fmt.Errorf("invalid post id: %d", postID)
	}
	//创建缓存键
	cacheKey := fmt.Sprintf(KeyFormat, postID)
	//尝试从redis中查找缓存
	detail, err := c.getFromCache(cacheKey)
	if err == nil && detail != nil {
		// 🎉 缓存命中，直接返回（~10ms）
		return detail, nil
	}

	post, err := c.repo.FindByID(postID)
	if err != nil {
		log.Printf("Failed to find post %d from database: %v", postID, err)
		return nil, fmt.Errorf("failed to get post detail: %w", err)
	}
	//检查帖子是否是已经发布的
	if post.Status != "published" {
		return nil, fmt.Errorf("post %d not found or has been deleted", postID)
	}
	// 处理标签（从字符串转换为数组）
	var tags []string
	if post.Tags != "" {
		tags = strings.Split(post.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	//将从数据库中查找到帖子信息进行dto转化
	detail = &services.PostDetail{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		CreatedAt: post.CreatedAt.Format(TimeFormat),
		UpdatedAt: post.UpdatedAt.Format(TimeFormat),
		Author: services.AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		},
		Category: services.CategoryInfo{
			ID:   post.Category.ID,
			Name: post.Category.Name,
			Icon: post.Category.Icon,
		},
		Tags: tags,
	}

	//异步缓存，不阻塞主流程
	go func() {
		err := databases.RDB.Set(context.Background(), cacheKey, detail, Expiration).Err()
		if err != nil {
			// 只记录日志，不影响返回结果
			log.Printf("Failed to cache post detail %d: %v", postID, err)
		}
	}()

	return detail, nil
}

// ═══════════════════════════════════════════════════════════════════════════════════
// InvalidatePostDetail 删除缓存
// ═══════════════════════════════════════════════════════════════════════════════════
// 使用场景：
//   - 用户编辑帖子后清空缓存
//   - 帖子被删除时清空缓存
//   - 点赞/评论时更新缓存
func (c *postCacheImpl) InvalidatePostDetail(postID uint) error {
	if postID == 0 {
		return fmt.Errorf("invalid post id: %d", postID)
	}

	cacheKey := fmt.Sprintf(KeyFormat, postID)
	err := databases.RDB.Del(context.Background(), cacheKey).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache for post %d: %v", postID, err)
		// 注意：缓存失效失败不应该导致业务异常，只记录日志
		// 最坏的情况就是返回旧数据，等到过期时间后会自动清空
		return nil
	}
	return nil
}

func (c *postCacheImpl) getFromCache(cacheKey string) (*services.PostDetail, error) {
	cacheData, err := databases.RDB.Get(context.Background(), cacheKey).Result()
	if err != nil {
		return nil, err
	}
	if cacheData == "" {
		return nil, nil
	}
	//反序列化
	var detail services.PostDetail
	err = json.Unmarshal([]byte(cacheData), &detail)
	if err != nil {
		log.Printf("Failed to unmarshal cache data for key %s: %v", cacheKey, err)
		_ = databases.RDB.Del(context.Background(), cacheKey).Err()
		return nil, nil
	}
	return &detail, nil
}
