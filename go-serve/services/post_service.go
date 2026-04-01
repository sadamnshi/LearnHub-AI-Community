package services

import (
	"fmt"
	"gin_demo/models"
	"gin_demo/repositories"
	"strings"
)

// 为什么在 services 包中定义？
//   - services 包依赖这个接口
//   - cache 包实现这个接口
//   - 这样避免循环导入（services 不需要导入 cache）

type PostDetailCacher interface {
	// GetPostDetail 从缓存获取帖子详情
	GetPostDetail(postID uint) (*PostDetail, error)
}

type PostService interface {
	// GetPostList 获取帖子分页列表
	GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error)

	// GetPostDetail 获取帖子详情
	GetPostDetail(postID uint) (*PostDetail, error)

	// CreatePost 创建新帖子
	CreatePost(authorID uint, req *CreatePostRequest) (*PostDetail, error)
}

// PostListResult 帖子列表结果（含分页信息）
type PostListResult struct {
	// List 帖子摘要列表
	List []PostSummary `json:"list"`

	// Total 符合条件的总帖子数（用来计算总页数）
	Total int64 `json:"total"`

	// Page 当前页码（用来告诉前端用户现在看的是第几页）

	Page int `json:"page"`

	// PageSize 每页显示的条数（用来告诉前端这次返回了多少条）

	PageSize int `json:"page_size"`
}

// PostSummary用来保存列表页展示的帖子摘要信息，不返回整个帖子的信息
type PostSummary struct {
	// ID 帖子的唯一标识
	// 前端点击帖子时会用这个 ID 来获取详情
	ID uint `json:"id"`

	// Title 帖子标题
	Title string `json:"title"`

	// Summary 帖子内容摘要（正文前 200 字符）
	// 作用：在列表页显示预览，不用加载完整的长正文
	Summary string `json:"summary"`

	// Author 作者信息（嵌套的对象）

	Author AuthorInfo `json:"author"`

	// Category 分类信息（嵌套的对象）
	Category CategoryInfo `json:"category"`

	// Status 帖子状态
	// 可能的值："draft"（草稿）、"published"（已发布）、"banned"（已禁）
	Status string `json:"status"`

	// CreatedAt 创建时间（格式化后的字符串）
	// 例如："2026-03-15 10:30:45"
	CreatedAt string `json:"created_at"`

	// Tags 帖子标签数组（已解析）
	// 示例：["golang", "并发", "最佳实践"]
	// 数据来源：从数据库的 Tags 字段（逗号分隔字符串）解析而来
	Tags []string `json:"tags"`
}

// AuthorInfo 作者摘要信息（不暴露密码等敏感字段）

type AuthorInfo struct {
	// ID 作者的用户 ID
	ID uint `json:"id"`

	// Username 作者的用户名
	// 例如："tom"
	Username string `json:"username"`

	// Avatar 作者的头像 URL
	// 例如："https://example.com/avatars/tom.jpg"
	Avatar string `json:"avatar"`
}

// CategoryInfo 分类摘要

type CategoryInfo struct {
	// ID 分类 ID
	ID uint `json:"id"`

	// Name 分类名称

	Name string `json:"name"`

	// Icon 分类图标（emoji 或图标名）

	Icon string `json:"icon"`
}

type PostDetail struct {
	// 基础信息
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"` // ← 完整内容，不摘要
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"` // ← 新增：更新时间

	// 作者信息
	Author AuthorInfo `json:"author"`

	// 分类信息
	Category CategoryInfo `json:"category"`

	// 标签信息 - 字符串数组（已解析）
	// 示例：["golang", "并发", "最佳实践"]
	Tags []string `json:"tags"`
}

// 创建帖子模型
type CreatePostRequest struct {
	// Title 帖子标题
	// binding:"required" - 表示必填项，如果不提供会自动验证失败
	// binding:"max=200" - 表示最多 200 个字符
	Title string `json:"title" binding:"required,max=200"`

	// Content 帖子内容（支持 Markdown）
	// binding:"required" - 表示必填项
	// binding:"min=10" - 表示至少 10 个字符（防止垃圾内容）
	Content string `json:"content" binding:"required,min=10"`

	// CategoryID 帖子分类 ID
	// 可选字段（binding 中没有 required）
	// 如果用户不提供，默认为 0（不分类）
	CategoryID uint `json:"category_id"`

	// Tags 帖子标签（逗号分隔的字符串）
	// 示例："golang,数据库,缓存"
	// 可选字段
	Tags string `json:"tags" binding:"max=200"`
}

type postService struct {
	repo  repositories.PostRepository
	cache PostDetailCacher // ← 使用 PostDetailCacher 接口
}

// 架构原理：
//   - postService 依赖 repo 和 cache
//   - cache 内部也持有 repo（但是同一个对象）
//   - 这不是"两个 repo"，而是"同一个 repo 被两个地方引用"
func NewPostService(repo repositories.PostRepository, cache PostDetailCacher) PostService {
	return &postService{
		repo:  repo,
		cache: cache,
	}
}

func (s *postService) GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error) {

	// 检查页码是否有效
	if page < 1 {
		// 如果页码小于 1，纠正为 1（显示第一页）
		page = 1
	}

	// 检查每页条数是否有效
	if pageSize < 1 || pageSize > 20 {
		pageSize = 5 // 默认每页 5 条，最多 50 条
	}

	posts, total, err := s.repo.List(page, pageSize, categoryID)

	// 如果数据库查询出错，直接返回错误
	// 比如：网络问题、数据库离线、SQL 语法错误等
	if err != nil {
		return nil, err
	}

	// 将模型转换为 DTO，并进行业务逻辑处理
	// make([]PostSummary, 0, len(posts))
	// 这行创建一个空切片，用来存储转换后的数据
	// 参数说明：
	//   - 第一个参数 0：初始长度为 0（开始是空的）
	//   - 第二个参数 len(posts)：预分配容量（避免后续追加时频繁重新分配内存）
	// 为什么要预分配容量？性能优化！
	list := make([]PostSummary, 0, len(posts))

	// 遍历从数据库查询出来的每一个帖子
	for _, p := range posts {

		// 将正文赋值给 summary 变量（暂时都是完整正文）
		summary := p.Content

		// 将字符串转换为 rune 切片
		// 为什么要转换？
		//   - Go 中的 string 是按字节存储的
		//   - 中文等多字节字符按字节截取会乱码
		//   - 转换为 rune 切片后，按字符计算（不管几个字节）
		// 示例：
		//   "Hello" 有 5 个字节，也有 5 个字符
		//   "你好" 有 6 个字节，但只有 2 个字符
		//   所以必须用 rune 来正确计算中文字符数
		runes := []rune(summary)

		// 检查是否超过 200 字
		if len(runes) > 200 {
			// 如果超过，截取前 200 字，后面加上 "..."
			summary = string(runes[:200]) + "..."
		}

		// 声明一个空的标签切片
		var tags []string

		// 只有当标签不为空时才处理
		// p.Tags 可能是 ""（空字符串）或 "golang,database"
		if p.Tags != "" {
			// strings.Split() 将字符串按 "," 分割成切片
			// 示例：
			//   strings.Split("golang,database,缓存", ",")
			//   返回：["golang", "database", "缓存"]
			tags = strings.Split(p.Tags, ",")

			// 清理每个标签的首尾空白符
			// 为什么需要清理？
			//   - 有时候标签字符串可能是 "golang, database, 缓存"（逗号后面有空格）
			//   - Split 后会变成 ["golang", " database", " 缓存"]（多了空格）
			//   - 前端显示会不整洁
			for i, tag := range tags {
				// strings.TrimSpace() 删除字符串首尾的空格、tab、换行等
				// 示例：
				//   strings.TrimSpace(" golang ")
				//   返回："golang"
				tags[i] = strings.TrimSpace(tag)
			}
		}

		list = append(list, PostSummary{
			ID:        p.ID,
			Title:     p.Title,
			Summary:   summary,
			Status:    p.Status,
			CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
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
			Tags: tags,
		})
	}

	// ═══════════════════════════════════════════════════════════════
	// 步骤 5：返回最终结果
	// ═══════════════════════════════════════════════════════════════
	// 构建返回给控制器的结果对象
	return &PostListResult{
		List:     list,     // 转换后的帖子摘要列表
		Total:    total,    // 数据库返回的总帖子数
		Page:     page,     // 当前页码
		PageSize: pageSize, // 每页显示的条数
	}, nil // nil 表示没有错误
}

//  4. 增加浏览次数（可选，需要 Redis）
//  5. 返回详情给前端
//
// 参数：
//   - postID: 要查询的帖子 ID
//
// 返回：
//   - *PostDetail: 帖子详情数据
//   - error: 如果出错返回错误（比如帖子不存在）
//
// 与 GetPostList 的区别：

func (s *postService) GetPostDetail(postID uint) (*PostDetail, error) {

	return s.cache.GetPostDetail(postID)
}

// 参数：
//   - authorID: 创建者的用户 ID（从 JWT Token 中提取）
//   - req: 前端提交的创建帖子请求（已通过 Gin 的自动验证）
//
// 返回：
//   - *PostDetail: 创建后的帖子详情
//   - error: 如果出错返回错误（比如数据库插入失败）
func (s *postService) CreatePost(authorID uint, req *CreatePostRequest) (*PostDetail, error) {

	// ═════════════════════════════════════════════════════════════════════════════════
	// 步骤 1️⃣：参数验证（额外的业务规则检查）
	// ═════════════════════════════════════════════════════════════════════════════════
	// 注意：
	//   - 基础的格式校验（比如 required、max、min）已经由 Gin 的 binding 完成
	//   - 这里做额外的业务逻辑检查

	// 检查标题是否为空或只包含空格
	if strings.TrimSpace(req.Title) == "" {
		return nil, fmt.Errorf("标题不能为空")
	}

	// 检查内容是否为空或只包含空格
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("内容不能为空")
	}

	// ═════════════════════════════════════════════════════════════════════════════════
	// 步骤 2️⃣：构建 Post 数据库模型
	// ═════════════════════════════════════════════════════════════════════════════════
	// 创建一个新的 Post 结构体，填入前端提交的数据和默认值
	post := &models.Post{
		// 从请求中直接获取的字段
		Title:      strings.TrimSpace(req.Title), // 去掉首尾空格
		Content:    req.Content,
		AuthorID:   authorID,
		CategoryID: req.CategoryID,
		Tags:       strings.TrimSpace(req.Tags),
		Status:     "published", // 默认设置为已发布（可改为 "draft" 草稿）
	}

	// ═════════════════════════════════════════════════════════════════════════════════
	// 步骤 3️⃣：调用仓储层保存到数据库
	// ═════════════════════════════════════════════════════════════════════════════════
	// s.repo.Create(post) 会：
	//   1. 执行 INSERT SQL 将数据插入数据库
	//   2. GORM 自动生成 ID、CreatedAt、UpdatedAt
	//   3. 返回修改后的 post 对象（包含新生成的 ID）
	// 可能的错误：
	//   - 数据库连接失败
	//   - 唯一索引冲突（如果 Title 设置了唯一索引）
	//   - 外键约束失败（如 AuthorID 或 CategoryID 不存在）
	createdPost, err := s.repo.Create(post)
	if err != nil {
		// 返回错误，由控制器决定如何响应
		return nil, fmt.Errorf("创建帖子失败: %w", err)
	}

	// ═════════════════════════════════════════════════════════════════════════════════
	// 步骤 4️⃣：处理标签（从字符串转换为切片）
	// ═════════════════════════════════════════════════════════════════════════════════
	// 数据库中存储的标签是逗号分隔的字符串，需要转换为数组返回给前端
	var tags []string
	if createdPost.Tags != "" {
		tags = strings.Split(createdPost.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// ═════════════════════════════════════════════════════════════════════════════════
	// 步骤 5️⃣：构建 PostDetail 响应对象
	// ═════════════════════════════════════════════════════════════════════════════════
	// 将创建后的 Post 模型转换为 PostDetail DTO
	// 注意：我们还需要获取作者和分类的完整信息
	// 但通常后续的 Preload 会在 FindByID 中进行
	// 这里为了演示，我们直接使用 createdPost 中的关联信息
	detail := &PostDetail{
		ID:        createdPost.ID,
		Title:     createdPost.Title,
		Content:   createdPost.Content,
		Status:    createdPost.Status,
		CreatedAt: createdPost.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: createdPost.UpdatedAt.Format("2006-01-02 15:04:05"),
		// 这里作者和分类可能不完整（取决于 Create 后是否自动 Preload）
		// 实际应用中可能需要额外查询
		Author: AuthorInfo{
			ID:       createdPost.Author.ID,
			Username: createdPost.Author.Username,
			Avatar:   createdPost.Author.Avatar,
		},
		Category: CategoryInfo{
			ID:   createdPost.Category.ID,
			Name: createdPost.Category.Name,
			Icon: createdPost.Category.Icon,
		},
		Tags: tags,
	}

	// ═════════════════════════════════════════════════════════════════════════════════
	// 步骤 6️⃣：返回结果
	// ═════════════════════════════════════════════════════════════════════════════════
	return detail, nil
}
