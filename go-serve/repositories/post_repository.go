package repositories

import (
	"gin_demo/databases"
	"gin_demo/models"
)

// PostRepository 帖子数据访问接口
type PostRepository interface {
	// List 分页查询帖子列表，可按分类筛选
	// page 从 1 开始，pageSize 每页条数，categoryID=0 表示不筛选
	List(page, pageSize int, categoryID uint) ([]models.Post, int64, error)
	// FindByID 查询帖子详情（预加载作者、分类、标签）
	FindByID(id uint) (*models.Post, error)
}

type postRepository struct{}

func NewPostRepository() PostRepository {
	return &postRepository{}
}

// List 分页查询帖子列表
func (r *postRepository) List(page, pageSize int, categoryID uint) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := databases.DB.Model(&models.Post{}).Where("status = ?", "published")

	// 按分类筛选（categoryID=0 时不过滤）
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 先统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载关联数据，按创建时间倒序（最新在前）
	offset := (page - 1) * pageSize
	err := query.
		Preload("Author").   // 预加载作者信息，避免 N+1 查询
		Preload("Category"). // 预加载分类信息
		Preload("Tags").     // 预加载标签
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error

	return posts, total, err
}

// FindByID 查询帖子详情
func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := databases.DB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
