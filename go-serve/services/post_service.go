package services

import (
	"gin_demo/repositories"
)

// PostListResult 帖子列表结果（含分页信息）
type PostListResult struct {
	List     []PostSummary `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// PostSummary 帖子列表摘要（不返回完整正文，节省流量）
type PostSummary struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Summary      string       `json:"summary"` // 正文前 200 字
	Author       AuthorInfo   `json:"author"`
	Category     CategoryInfo `json:"category"`
	Tags         []TagInfo    `json:"tags"`
	ViewCount    int          `json:"view_count"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	IsPinned     bool         `json:"is_pinned"`
	CreatedAt    string       `json:"created_at"`
}

// AuthorInfo 作者摘要信息（不暴露密码等敏感字段）
type AuthorInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// CategoryInfo 分类摘要
type CategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// TagInfo 标签摘要
type TagInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// PostService 帖子业务接口
type PostService interface {
	GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error)
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo: repo}
}

// GetPostList 获取帖子分页列表
func (s *postService) GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error) {
	// 参数边界保护
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20 // 默认每页 20 条，最多 50 条
	}

	posts, total, err := s.repo.List(page, pageSize, categoryID)
	if err != nil {
		return nil, err
	}

	// 将 model 转换为 DTO（数据传输对象），避免直接暴露数据库结构
	list := make([]PostSummary, 0, len(posts))
	for _, p := range posts {
		// 生成摘要：取正文前 200 字
		summary := p.Content
		runes := []rune(summary)
		if len(runes) > 200 {
			summary = string(runes[:200]) + "..."
		}

		// 构建标签列表
		tags := make([]TagInfo, 0, len(p.Tags))
		for _, t := range p.Tags {
			tags = append(tags, TagInfo{ID: t.ID, Name: t.Name})
		}

		list = append(list, PostSummary{
			ID:      p.ID,
			Title:   p.Title,
			Summary: summary,
			Author: AuthorInfo{
				ID:       p.Author.ID,
				Username: p.Author.Username,
				Avatar:   p.Author.Avatar,
			},
			Category: CategoryInfo{
				ID:   p.Category.ID,
				Name: p.Category.Name,
				Icon: p.Category.Icon,
			},
			Tags:         tags,
			ViewCount:    p.ViewCount,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
			IsPinned:     p.IsPinned,
			CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &PostListResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
