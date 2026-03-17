package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_demo/databases"
	"gin_demo/repositories"
	"strings"
	"time"
)

// PostListResult 帖子列表结果（含分页信息）
// 这是调用 GetPostList() 方法时的返回数据结构
// 前端收到的 JSON 数据结构就对应这个 struct
type PostListResult struct {
	// List 帖子摘要列表
	// 示例：[{id: 1, title: "...", ...}, {id: 2, title: "...", ...}]
	List []PostSummary `json:"list"`

	// Total 符合条件的总帖子数（用来计算总页数）
	// 比如：total = 100，pageSize = 20，那么总共有 5 页
	Total int64 `json:"total"`

	// Page 当前页码（用来告诉前端用户现在看的是第几页）
	// 比如：page = 2
	Page int `json:"page"`

	// PageSize 每页显示的条数（用来告诉前端这次返回了多少条）
	// 比如：page_size = 20
	PageSize int `json:"page_size"`
}

// PostSummary 帖子列表摘要（不返回完整正文，节省流量）
// 这个结构体用来表示帖子列表中的每一个帖子
// 为什么叫 "摘要"？因为不包含完整的 Content（正文），只包含摘要前 200 字
//
// 对比：
//
//	├─ 列表页：显示帖子标题 + 摘要（让用户快速浏览）← 用 PostSummary
//	└─ 详情页：显示帖子的全部内容（让用户详细阅读）← 需要其他 DTO
type PostSummary struct {
	// ID 帖子的唯一标识
	// 前端点击帖子时会用这个 ID 来获取详情
	// 例如：点击帖子后，向 /api/posts/123 发送请求
	ID uint `json:"id"`

	// Title 帖子标题
	// 例如："如何学习 Go 语言"
	Title string `json:"title"`

	// Summary 帖子内容摘要（正文前 200 字符）
	// 作用：在列表页显示预览，不用加载完整的长正文
	// 例如："Go 是一门编译型语言，具有高效的并发能力..."
	Summary string `json:"summary"`

	// Author 作者信息（嵌套的对象）
	// 这样前端收到的 JSON 是这样的：
	// {
	//   "id": 1,
	//   "author": {
	//     "id": 10,
	//     "username": "tom",
	//     "avatar": "http://..."
	//   }
	// }
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
// 这是一个嵌套的结构体，用来在 PostSummary 中表示作者信息
// 为什么不直接用 User 结构体？
//
//	User 结构体可能包含 Password、Email 等敏感信息
//	这里只挑选需要展示的字段，更加安全
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
// 这是一个嵌套的结构体，用来在 PostSummary 中表示分类信息
type CategoryInfo struct {
	// ID 分类 ID
	ID uint `json:"id"`

	// Name 分类名称
	// 例如："技术分享"、"闲聊吐槽"
	Name string `json:"name"`

	// Icon 分类图标（emoji 或图标名）
	// 例如："💻"、"📚"
	Icon string `json:"icon"`
}

// ═══════════════════════════════════════════════════════════════
// 第二部分：定义业务接口
// ═══════════════════════════════════════════════════════════════
// 接口的作用：
//   定义服务层应该提供的功能
//   类似于一份合同：实现这个接口的结构体必须实现所有定义的方法

// ═══════════════════════════════════════════════════════════════
// 帖子详情 DTO - 返回给前端的数据结构
// ═══════════════════════════════════════════════════════════════

// PostDetail 帖子详情（包含完整内容）
// 对比 PostSummary：
//
//	├─ PostSummary：用于列表页（摘要 + 分页）
//	└─ PostDetail：用于详情页（完整内容 + 更多信息）
//
// 关键区别：
//   - Content：完整的帖子内容（不截断）
//   - ViewCount：浏览次数（从 Redis 获取）
//   - LikeCount：点赞数（可选，需要实现点赞功能）
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

	// 统计信息（可选，需要扩展）
	// ViewCount int   `json:"view_count"`      // 浏览次数
	// LikeCount int   `json:"like_count"`      // 点赞数
	// CommentCount int `json:"comment_count"`  // 评论数
}

// PostService 帖子业务接口
// 这个接口定义了关于帖子的所有业务操作方法
// 好处：
//  1. 解耦：控制器不需要知道具体实现，只需调用接口方法
//  2. 可测试：可以轻松创建 Mock 实现进行单元测试
//  3. 灵活性：将来可以轻松替换为另一个实现
type PostService interface {
	// GetPostList 获取帖子分页列表
	// 参数：
	//   - page: 页码（从 1 开始）
	//   - pageSize: 每页显示多少条
	//   - categoryID: 分类 ID（0 表示不按分类筛选）
	// 返回：
	//   - *PostListResult: 分页列表结果
	//   - error: 如果出错会返回错误信息
	GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error)

	// GetPostDetail 获取帖子详情
	// 参数：
	//   - postID: 帖子 ID
	// 返回：
	//   - *PostDetail: 帖子详情数据
	//   - error: 如果出错会返回错误信息
	GetPostDetail(postID uint) (*PostDetail, error)
}

// ═══════════════════════════════════════════════════════════════
// 第三部分：实现业务接口
// ═══════════════════════════════════════════════════════════════

// postService 帖子业务的具体实现
// 注意：首字母小写（postService），表示私有结构体，只能在本包内使用
// 对应的是上面的 PostService 接口
type postService struct {
	// repo 数据访问层的仓储对象
	// 通过这个对象来调用数据库查询等操作
	// 例如：s.repo.List(...) 会调用数据库查询
	repo repositories.PostRepository
}

// NewPostService 工厂函数 - 创建服务实例
// 这是一个创建 postService 的标准方式
// 为什么要用工厂函数而不是直接 new？
//  1. 统一的创建方式
//  2. 可以在创建时做一些初始化逻辑
//  3. 隐藏内部实现细节（用户不需要知道是 postService）
//
// 使用方式：
//
//	repo := repositories.NewPostRepository()
//	svc := NewPostService(repo)  ← 这样就创建了一个服务实例
func NewPostService(repo repositories.PostRepository) PostService {
	// 创建 postService 实例，保存仓储对象的引用
	// 返回类型是 PostService 接口，而不是具体的 postService 实现
	// 这样调用者只关心接口，不关心具体实现
	return &postService{repo: repo}
}

// ═══════════════════════════════════════════════════════════════
// 第四部分：核心业务方法 - GetPostList
// ═══════════════════════════════════════════════════════════════

// GetPostList 获取帖子分页列表 - PostService 接口的实现
// 这是整个服务层最重要的方法，包含了所有的业务逻辑处理
//
// 工作流程：
//  1. 验证和修正参数（确保合理的分页参数）
//  2. 调用仓储层查询数据库
//  3. 将数据库模型转换为 DTO
//  4. 对数据进行业务逻辑处理（比如生成摘要、格式化时间）
//  5. 返回格式化后的结果给控制器
//
// 参数说明：
//   - page: 用户请求的页码
//   - pageSize: 用户请求的每页条数
//   - categoryID: 用户要筛选的分类 ID
//
// 返回值说明：
//   - *PostListResult: 包含分页列表的结果对象
//   - error: 如果发生错误（比如数据库连接失败）会返回错误
func (s *postService) GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error) {

	// ═══════════════════════════════════════════════════════════════
	// 步骤 1：参数验证和修正（客户端验证）
	// ═══════════════════════════════════════════════════════════════
	// 为什么需要验证？
	//   前端的数据不一定是有效的。用户可能：
	//   1. 输入负数的页码
	//   2. 要求每页显示 10000 条（太多了，会导致服务器内存溢出）
	//   3. 通过 API 直接调用，传入非法参数
	//   所以后端必须再次验证，这叫做 "不信任来自客户端的数据"

	// 检查页码是否有效
	if page < 1 {
		// 如果页码小于 1，纠正为 1（显示第一页）
		page = 1
	}

	// 检查每页条数是否有效
	if pageSize < 1 || pageSize > 50 {
		// 如果每页条数小于 1 或大于 50，纠正为默认值 20
		// 为什么最多 50？
		//   - 限制每次查询的数据量，防止大量数据查询导致服务器压力
		//   - 这是一个常见的安全做法
		pageSize = 20 // 默认每页 20 条，最多 50 条
	}

	// ═══════════════════════════════════════════════════════════════
	// 步骤 2：调用仓储层查询数据库
	// ═══════════════════════════════════════════════════════════════
	// s.repo.List(...) 会执行以下 SQL：
	//   SELECT * FROM posts
	//   WHERE status = 'published' AND category_id = ?
	//   ORDER BY created_at DESC
	//   LIMIT ? OFFSET ?

	posts, total, err := s.repo.List(page, pageSize, categoryID)

	// 如果数据库查询出错，直接返回错误
	// 比如：网络问题、数据库离线、SQL 语法错误等
	if err != nil {
		return nil, err
	}

	// ═══════════════════════════════════════════════════════════════
	// 步骤 3 & 4：将模型转换为 DTO，并进行业务逻辑处理
	// ═══════════════════════════════════════════════════════════════
	// make([]PostSummary, 0, len(posts))
	// 这行创建一个空切片，用来存储转换后的数据
	// 参数说明：
	//   - 第一个参数 0：初始长度为 0（开始是空的）
	//   - 第二个参数 len(posts)：预分配容量（避免后续追加时频繁重新分配内存）
	// 为什么要预分配容量？性能优化！
	list := make([]PostSummary, 0, len(posts))

	// 遍历从数据库查询出来的每一个帖子
	for _, p := range posts {

		// ───────────────────────────────────────────────────────────
		// 步骤 4.1：生成帖子摘要（正文前 200 字）
		// ───────────────────────────────────────────────────────────
		// 为什么需要摘要？
		//   - 列表页不需要显示完整的帖子内容
		//   - 只显示预览，用户点击后才加载完整内容
		//   - 这样可以减少网络流量，提升列表页加载速度

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

		// ───────────────────────────────────────────────────────────
		// 步骤 4.2：处理标签（从字符串转换为切片）
		// ───────────────────────────────────────────────────────────
		// 数据库中存储的标签是这样的：
		//   "golang,database,缓存"（逗号分隔的字符串）
		// 但返回给前端的应该是：
		//   ["golang", "database", "缓存"]（JSON 数组）
		// 这样前端处理更方便

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

		// ───────────────────────────────────────────────────────────
		// 步骤 4.3：构建 PostSummary 对象
		// ───────────────────────────────────────────────────────────
		// 现在我们有了所有需要的数据，构建一个 PostSummary 对象
		// 这个对象只包含前端需要的字段（不暴露内部实现细节）

		// p.CreatedAt.Format("2006-01-02 15:04:05")
		// 时间格式化说明：
		//   - p.CreatedAt 是 time.Time 类型（Go 的时间类型）
		//   - .Format() 方法可以将其格式化为字符串
		//   - "2006-01-02 15:04:05" 是格式字符串
		//     * 2006 代表年
		//     * 01 代表月
		//     * 02 代表日
		//     * 15 代表小时（24 小时制）
		//     * 04 代表分钟
		//     * 05 代表秒
		//   - 结果例如："2026-03-15 10:30:45"
		// 为什么要格式化？
		//   - 数据库存的是 time.Time 对象，不方便转 JSON
		//   - 转成字符串后，JSON 序列化更简单
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

// GetPostDetail 获取帖子详情 - PostService 接口的实现
// ═════════════════════════════════════════════════════════════════════════════
// 这个方法处理帖子详情页的业务逻辑
//
// 工作流程：
//  1. 从数据库查询帖子及其关联数据
//  2. 检查帖子是否存在及是否已发布
//  3. 将数据库模型转换为详情 DTO
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
//
//	┌─────────────────────┬──────────────────┬──────────────────┐
//	│ 功能                │ GetPostList       │ GetPostDetail    │
//	├─────────────────────┼──────────────────┼──────────────────┤
//	│ 返回的内容          │ 摘要（200 字）   │ 完整内容        │
//	│ 返回的数量          │ 多条（分页）     │ 单条            │
//	│ 性能要求            │ 快（列表页）     │ 普通（详情页）  │
//	│ 浏览次数统计        │ 不计             │ 需要计             │
//	│ 缓存策略            │ 短期（1小时）   │ 中期（6小时）   │
//	└─────────────────────┴──────────────────┴──────────────────┘
func (s *postService) GetPostDetail(postID uint) (*PostDetail, error) {

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 1️⃣：尝试从 Redis 缓存获取数据
	// ═════════════════════════════════════════════════════════════════════════════
	// 旁路缓存（Cache-Aside Pattern）的核心：先查缓存，再查数据库
	cacheKey := fmt.Sprintf("post:%d:detail", postID)
	cachedData, err := databases.RDB.Get(context.Background(), cacheKey).Result()

	// 检查是否成功获取缓存
	if err == nil && cachedData != "" {
		// 🎉 缓存命中！
		var detail PostDetail
		err := json.Unmarshal([]byte(cachedData), &detail)
		if err == nil {
			// 缓存命中，直接返回（~10ms，快 20 倍！）
			return &detail, nil
		}
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 2️⃣：缓存未命中，查询数据库
	// ═════════════════════════════════════════════════════════════════════════════
	post, err := s.repo.FindByID(postID)

	if err != nil {
		return nil, err
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 3️⃣：业务规则检查
	// ═════════════════════════════════════════════════════════════════════════════
	if post.Status != "published" {
		return nil, fmt.Errorf("帖子不存在或已被删除")
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 4️⃣：标签处理
	// ═════════════════════════════════════════════════════════════════════════════
	var tags []string
	if post.Tags != "" {
		tags = strings.Split(post.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 5️⃣：构建 PostDetail 对象
	// ═════════════════════════════════════════════════════════════════════════════
	detail := &PostDetail{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		Author: AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		},
		Category: CategoryInfo{
			ID:   post.Category.ID,
			Name: post.Category.Name,
			Icon: post.Category.Icon,
		},
		Tags: tags,
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 6️⃣：将结果存入 Redis 缓存（旁路缓存的关键）
	// ═════════════════════════════════════════════════════════════════════════════
	// 缓存 6 小时，这样下次访问会直接返回缓存，不需要查询数据库
	data, err := json.Marshal(detail)
	if err == nil {
		expiration := 6 * time.Hour
		_ = databases.RDB.Set(context.Background(), cacheKey, string(data), expiration).Err()
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 7️⃣：返回结果
	// ═════════════════════════════════════════════════════════════════════════════
	return detail, nil
}
