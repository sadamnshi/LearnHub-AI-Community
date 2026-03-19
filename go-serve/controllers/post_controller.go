package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gin_demo/services"
)

// PostController 帖子控制器
type PostController struct {
	service services.PostService
}

func NewPostController(service services.PostService) *PostController {
	return &PostController{service: service}
}

// GetPosts 获取帖子列表（无需登录，公开接口）
//
// GET /api/posts?page=1&page_size=20&category_id=0
//
// 查询参数：
//   - page: 页码，默认 1
//   - page_size: 每页条数，默认 20，最大 50
//   - category_id: 分类 ID，默认 0（不筛选）
func (ctl *PostController) GetPosts(c *gin.Context) {
	// 解析查询参数，提供默认值
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 64)

	result, err := ctl.service.GetPostList(page, pageSize, uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子列表失败"})
		return
	}

	// 统一响应格式：{ "code": 0, "data": {...}, "msg": "success" }
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": result,
	})
}

// ═══════════════════════════════════════════════════════════════════════════
// 获取帖子详情
// ═══════════════════════════════════════════════════════════════════════════

// GetPostDetail 获取单个帖子的详细信息（无需登录，公开接口）
//
// 📌 路由：GET /api/posts/:id
//
// 🔗 URL 示例：
//   - GET /api/posts/1     → 获取 ID 为 1 的帖子详情
//   - GET /api/posts/123   → 获取 ID 为 123 的帖子详情
//
// ⚙️ 工作流程：
//  1. 从 URL 路径参数提取帖子 ID
//  2. 参数验证（必须是有效的正整数）
//  3. 调用服务层获取详情
//  4. 错误处理（帖子不存在返回 404）
//  5. 返回 JSON 响应
//
// 📤 成功响应（HTTP 200）：
//
//	{
//	  "code": 0,
//	  "msg": "success",
//	  "data": {
//	    "id": 1,
//	    "title": "Go 语言最佳实践",
//	    "content": "完整的帖子内容...",
//	    "author": {...},
//	    "category": {...},
//	    "tags": ["golang", "后端"],
//	    "created_at": "2026-03-15 10:30:45",
//	    "updated_at": "2026-03-15 10:30:45"
//	  }
//	}
//
// ❌ 失败响应（HTTP 404）：
//
//	{
//	  "code": 404,
//	  "msg": "帖子不存在或已被删除"
//	}
//
// ⚡ 性能优化：使用旁路缓存模式
//   - 首次访问：查询数据库 + 存入 Redis（~200ms）
//   - 后续访问：直接返回缓存（~10ms，快 20 倍！）
func (ctl *PostController) GetPostDetail(c *gin.Context) {

	// ═════════════════════════════════════════════════════════════════════════
	// 第 1 步：从 URL 路径参数提取帖子 ID
	// ═════════════════════════════════════════════════════════════════════════
	// 在路由中，:id 是一个动态参数
	// 例如：
	//   路由定义：GET /api/posts/:id
	//   用户请求：GET /api/posts/123
	//   c.Param("id") 返回："123"（字符串形式）
	//
	// 为什么是字符串？因为 URL 中所有的数据都是字符串，需要手动转换
	postIDStr := c.Param("id")

	// ═════════════════════════════════════════════════════════════════════════
	// 第 2 步：参数验证 - 转换字符串为数字
	// ═════════════════════════════════════════════════════════════════════════
	// strconv.ParseUint() 的参数说明：
	//   - postIDStr: 要转换的字符串，例如 "123"
	//   - 10: 十进制数字（不是二进制或十六进制）
	//   - 64: 转换为 uint64（64 位无符号整数）
	//
	// 返回值：(uint64, error)
	//   - uint64：转换成功的数字，例如 123
	//   - error：转换失败的错误（比如用户输入了 "abc"）
	//
	// 可能的错误场景：
	//   1. 用户输入非数字：/api/posts/abc → ParseUint 返回错误
	//   2. 用户输入负数：/api/posts/-123 → ParseUint 返回错误（无符号数不能是负数）
	//   3. 用户输入浮点数：/api/posts/1.5 → ParseUint 返回错误
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	// 如果转换失败，返回 400 Bad Request（请求参数格式错误）
	if err != nil {
		// 这里返回 400 而不是 404，是因为 400 表示"你的输入格式错了"
		// 404 表示"数据库中没有这个资源"
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的帖子 ID，必须是正整数",
		})
		return
	}

	// ═════════════════════════════════════════════════════════════════════════
	// 第 3 步：调用服务层获取帖子详情
	// ═════════════════════════════════════════════════════════════════════════
	// ctl.service.GetPostDetail() 这个方法会：
	//
	//   1️⃣ 查询 Redis 缓存
	//      └─ 如果缓存存在，直接返回（快！~10ms）
	//
	//   2️⃣ 缓存未命中，查询数据库
	//      └─ 执行 SQL: SELECT * FROM posts WHERE id = ?
	//
	//   3️⃣ 检查帖子是否已发布
	//      └─ 只有 status = "published" 的帖子才能被访问
	//
	//   4️⃣ 转换数据库模型为 DTO
	//      └─ Post (数据库模型) → PostDetail (API 响应)
	//      └─ 包括：标签解析、时间格式化、提取作者和分类信息
	//
	//   5️⃣ 将结果存入 Redis 缓存
	//      └─ Key: post:{id}:detail
	//      └─ 过期时间：6 小时
	//
	//   6️⃣ 返回 PostDetail 结果
	detail, err := ctl.service.GetPostDetail(uint(postID))

	// ═════════════════════════════════════════════════════════════════════════
	// 第 4 步：错误处理
	// ═════════════════════════════════════════════════════════════════════════
	// 可能出现的错误：
	//   1. 帖子不存在：数据库中没有 id = postID 的记录
	//   2. 帖子是草稿：status != "published"（不公开）
	//   3. 数据库连接错误：网络问题、数据库离线
	//   4. Redis 错误：缓存服务不可用（不影响功能，会回退到直接查询数据库）
	if err != nil {
		// 统一返回 404（用户不需要知道具体原因）
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "帖子不存在或已被删除",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": detail,
	})
}

func (ctl *PostController) CreatePost(c *gin.Context) {
	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 1️⃣：从上下文获取登录用户的 ID
	// ═════════════════════════════════════════════════════════════════════════════
	// AuthMiddleware 中间件会将用户 ID 存入 c.Set("userID", ...)
	// 这里通过 c.GetInt64("userID") 取出，并转换为 uint
	userIDInterface, exists := c.Get("userID")
	if !exists {
		// 如果上下文中没有 userID，说明中间件验证失败（不应该发生）
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权，请先登录",
		})
		return
	}

	// 类型断言：userIDInterface 是 int64 类型
	// int64 → uint：将有符号整数转为无符号整数
	userID := uint(userIDInterface.(int64))

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 2️⃣：解析请求 JSON 数据
	// ═════════════════════════════════════════════════════════════════════════════
	// c.ShouldBindJSON() 会：
	//   1. 读取请求体
	//   2. 反序列化 JSON 数据到结构体
	//   3. 自动验证 binding 标签（required、max、min 等）
	// 返回值：
	//   - nil：绑定和验证成功
	//   - error：绑定或验证失败（会自动返回 400 错误给前端）
	var req services.CreatePostRequest

	// 关于 ShouldBindJSON 的错误处理：
	// 如果 ShouldBindJSON 返回错误（比如 JSON 格式错误、验证失败），
	// Gin 会自动返回 400 错误给前端，但不会中断执行
	// 所以我们需要手动检查并返回
	if err := c.ShouldBindJSON(&req); err != nil {
		// 验证失败，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数格式错误或验证失败",
			"error": err.Error(),
		})
		return
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 3️⃣：调用服务层创建帖子
	// ═════════════════════════════════════════════════════════════════════════════
	// ctl.service.CreatePost() 会：
	//   1. 进行额外的业务逻辑验证
	//   2. 构建 Post 数据库模型
	//   3. 保存到数据库（自动生成 ID、时间戳）
	//   4. 返回 PostDetail 给前端
	detail, err := ctl.service.CreatePost(userID, &req)

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 4️⃣：错误处理
	// ═════════════════════════════════════════════════════════════════════════════
	if err != nil {
		// 创建失败，返回 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "创建帖子失败",
			"error": err.Error(),
		})
		return
	}

	// ═════════════════════════════════════════════════════════════════════════════
	// 步骤 5️⃣：返回成功响应
	// ═════════════════════════════════════════════════════════════════════════════
	// 返回创建后的帖子详情，前端可以用帖子 ID 跳转到详情页
	// 示例响应：
	// {
	//   "code": 0,
	//   "msg": "success",
	//   "data": {
	//     "id": 123,
	//     "title": "...",
	//     "content": "...",
	//     ...
	//   }
	// }
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg":  "success",
		"data": detail,
	})
}
