# 📖 帖子详情功能完整讲解

## 🎯 核心概念

你的问题涉及如何实现"获取帖子详情"功能。这是一个与"获取帖子列表"完全不同的业务逻辑。

### 💡 关键区别

```
获取列表 (GetPostList)          获取详情 (GetPostDetail)
├─ 返回多条数据                ├─ 返回单条数据
├─ 返回摘要（节省流量）       ├─ 返回完整内容
├─ 需要分页                   ├─ 不需要分页
├─ 关注性能（快速加载）       ├─ 关注准确性和完整性
├─ 可以缓存较短时间            ├─ 可以缓存较长时间
└─ URL: /api/posts/home?page=1 └─ URL: /api/posts/:id
```

---

## 📊 完整的业务流程

```
┌─────────────────────────────────────────────────────────────┐
│ 用户在浏览器打开帖子详情页                                    │
│ 例如：http://localhost:3000/posts/123                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────▼────────────────┐
        │ 前端 Vue 组件生命周期          │
        │ mounted() 或 onMounted()       │
        │ 从 URL 提取帖子 ID: 123        │
        └──────────────┬─────────────────┘
                       │
        ┌──────────────▼────────────────────────────┐
        │ 发送 HTTP 请求给后端                      │
        │ GET /api/posts/123                       │
        │ Content-Type: application/json           │
        └──────────────┬───────────────────────────┘
                       │
        ┌──────────────▼─────────────────────────────────────┐
        │ 后端 PostController.GetPostDetail(c *gin.Context)  │
        │ 1. 从 URL 参数提取帖子 ID                          │
        │ 2. 调用 service.GetPostDetail(postID)              │
        └──────────────┬────────────────────────────────────┘
                       │
        ┌──────────────▼──────────────────────────────────┐
        │ 服务层 PostService.GetPostDetail()              │
        │ 1. 调用仓储层查询数据库                         │
        │ 2. 检查帖子是否存在                            │
        │ 3. 检查帖子是否已发布                          │
        │ 4. 转换数据为 PostDetail DTO                   │
        │ 5. 处理标签字符串为数组                        │
        │ 6. 格式化时间                                  │
        │ 7. 增加浏览次数（可选）                        │
        └──────────────┬──────────────────────────────────┘
                       │
        ┌──────────────▼──────────────────────────────┐
        │ 仓储层 PostRepository.FindByID()             │
        │ 1. 构建 SQL 查询                            │
        │ 2. WHERE id = 123                          │
        │ 3. Preload 作者和分类                       │
        │ 4. 返回 Post 模型                          │
        └──────────────┬───────────────────────────────┘
                       │
        ┌──────────────▼───────────────────────┐
        │ PostgreSQL 数据库                    │
        │ 执行查询，返回结果                    │
        └──────────────┬────────────────────────┘
                       │
        ┌──────────────▼──────────────────────────────────────┐
        │ 数据回传链路                                        │
        │ Repository → Service → Controller → HTTP 响应      │
        └──────────────┬─────────────────────────────────────┘
                       │
        ┌──────────────▼──────────────────────────────┐
        │ HTTP 200 + JSON 响应                        │
        │ {                                          │
        │   "code": 0,                              │
        │   "msg": "success",                       │
        │   "data": {                               │
        │     "id": 123,                            │
        │     "title": "...",                       │
        │     "content": "完整内容...",               │
        │     "author": {...},                      │
        │     "category": {...},                    │
        │     "tags": [...],                        │
        │     "created_at": "...",                  │
        │     "updated_at": "..."                   │
        │   }                                       │
        │ }                                         │
        └──────────────┬───────────────────────────┘
                       │
        ┌──────────────▼──────────────┐
        │ 前端收到响应数据            │
        │ 1. 解析 JSON                │
        │ 2. 用 Markdown 渲染 content │
        │ 3. 显示作者和分类信息       │
        │ 4. 显示标签                 │
        │ 5. 更新页面标题            │
        └──────────────────────────────┘
```

---

## 🔧 实现思路详解

### 1️⃣ **数据模型对比**

```go
// 列表页使用的 DTO
type PostSummary struct {
    ID        uint
    Title     string
    Summary   string       // ← 摘要（200 字）
    Content   (不包含)     // ← 省略完整内容
    Tags      []string
    Author    AuthorInfo
    Category  CategoryInfo
    CreatedAt string
}

// 详情页使用的 DTO
type PostDetail struct {
    ID        uint
    Title     string
    Content   string       // ← 完整内容
    Summary   (不包含)     // ← 不需要摘要
    Tags      []string
    Author    AuthorInfo
    Category  CategoryInfo
    CreatedAt string
    UpdatedAt string       // ← 新增：更新时间
}
```

**为什么要分开？**
- 前端列表页只需要标题和摘要，减少网络流量
- 前端详情页需要完整内容，不能截断
- 这样设计更灵活，易于扩展

### 2️⃣ **关键业务逻辑**

#### A. 存在性检查

```
用户请求 ID 为 999 的帖子
    ↓
数据库中没有这个帖子
    ↓
仓储层返回错误：record not found
    ↓
服务层捕获错误
    ↓
返回 nil, error
    ↓
控制器处理错误
    ↓
返回 404 Not Found
```

#### B. 状态检查

```
用户请求 ID 为 5 的帖子
    ↓
帖子存在，但状态是 "draft"（草稿）
    ↓
普通用户不应该看到草稿
    ↓
服务层检查状态 != "published"
    ↓
返回错误
    ↓
控制器返回 404（隐藏草稿的存在）
```

#### C. 标签转换

```
数据库中的数据：
    Tags = "golang,并发编程,goroutine"

服务层处理：
    1. 检查是否为空
    2. Split(",") → ["golang", "并发编程", "goroutine"]
    3. TrimSpace 清理每个标签
    4. 返回干净的数组

前端收到：
    {
        "tags": ["golang", "并发编程", "goroutine"]
    }
```

#### D. 时间格式化

```
数据库中的数据：
    CreatedAt = time.Time{} (Go 时间对象)
    UpdatedAt = time.Time{}

服务层处理：
    1. post.CreatedAt.Format("2006-01-02 15:04:05")
    2. post.UpdatedAt.Format("2006-01-02 15:04:05")

前端收到：
    {
        "created_at": "2026-03-15 10:30:45",
        "updated_at": "2026-03-15 10:30:45"
    }
```

---

## 💻 代码实现详解

### 步骤 1: 定义 DTO 结构体

```go
type PostDetail struct {
    ID        uint         // 帖子 ID
    Title     string       // 标题
    Content   string       // 完整内容 ← 关键！
    Status    string       // 状态
    CreatedAt string       // 创建时间
    UpdatedAt string       // 更新时间 ← 新增
    Author    AuthorInfo   // 作者
    Category  CategoryInfo // 分类
    Tags      []string     // 标签数组
}
```

### 步骤 2: 在 Interface 中声明方法

```go
type PostService interface {
    GetPostList(page, pageSize int, categoryID uint) (*PostListResult, error)
    GetPostDetail(postID uint) (*PostDetail, error)  // ← 新增
}
```

### 步骤 3: 实现具体方法

```go
func (s *postService) GetPostDetail(postID uint) (*PostDetail, error) {
    // 1. 查询数据库
    post, err := s.repo.FindByID(postID)
    if err != nil {
        return nil, err  // 帖子不存在
    }

    // 2. 检查状态
    if post.Status != "published" {
        return nil, fmt.Errorf("帖子不存在或已被删除")
    }

    // 3. 处理标签
    var tags []string
    if post.Tags != "" {
        tags = strings.Split(post.Tags, ",")
        for i, tag := range tags {
            tags[i] = strings.TrimSpace(tag)
        }
    }

    // 4. 构建返回对象
    return &PostDetail{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,  // 完整内容
        Status:    post.Status,
        Tags:      tags,
        CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
        Author:    AuthorInfo{...},
        Category:  CategoryInfo{...},
    }, nil
}
```

### 步骤 4: 在 Controller 中调用

```go
func (ctl *PostController) GetPostDetail(c *gin.Context) {
    // 从 URL 参数提取帖子 ID
    // 例如：GET /api/posts/123
    postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    // 调用服务层
    detail, err := ctl.service.GetPostDetail(uint(postID))
    
    if err != nil {
        // 如果出错，返回 404
        c.JSON(http.StatusNotFound, gin.H{
            "code": 404,
            "msg":  "帖子不存在",
        })
        return
    }

    // 返回成功响应
    c.JSON(http.StatusOK, gin.H{
        "code": 0,
        "msg":  "success",
        "data": detail,
    })
}
```

### 步骤 5: 在 Router 中注册路由

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    // ... 其他路由 ...

    postController := controllers.NewPostController(postService)
    
    apiPost := r.Group("/api/posts")
    {
        apiPost.GET("/home", postController.GetPosts)           // 列表
        apiPost.GET("/:id", postController.GetPostDetail)       // 详情 ← 新增
    }

    return r
}
```

---

## 📱 前端调用示例

### Vue 3 Composition API

```javascript
// views/PostDetailView.vue
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

export default {
    setup() {
        const route = useRoute()
        const post = ref(null)
        const loading = ref(true)
        const error = ref(null)

        // 组件加载时获取详情
        onMounted(async () => {
            try {
                const postId = route.params.id
                const response = await fetch(
                    `http://localhost:8080/api/posts/${postId}`
                )
                const json = await response.json()
                
                if (json.code === 0) {
                    post.value = json.data
                } else {
                    error.value = "帖子不存在"
                }
            } catch (e) {
                error.value = e.message
            } finally {
                loading.value = false
            }
        })

        return {
            post,
            loading,
            error
        }
    }
}

// 模板
<template>
    <div v-if="loading">加载中...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else>
        <h1>{{ post.title }}</h1>
        <p>作者：{{ post.author.username }}</p>
        <p>分类：{{ post.category.name }}</p>
        <p>标签：{{ post.tags.join(', ') }}</p>
        <div v-html="markdownToHtml(post.content)"></div>
    </div>
</template>
```

---

## 🎨 JSON 响应格式

### 成功响应

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "title": "Go 语言并发编程最佳实践",
    "content": "在 Go 中，并发是通过 goroutine 和 channel 实现的。\n\nGoroutine 是一个轻量级的线程，由 Go 运行时管理。\n\n相比操作系统线程，goroutine 的开销极小，可以轻松创建数百万个 goroutine。\n\n...",
    "author": {
      "id": 1,
      "username": "tom",
      "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Tom"
    },
    "category": {
      "id": 1,
      "name": "技术分享",
      "icon": "💻"
    },
    "tags": ["golang", "并发", "goroutine", "channel"],
    "status": "published",
    "created_at": "2026-03-15 10:30:45",
    "updated_at": "2026-03-15 10:30:45"
  }
}
```

### 失败响应（帖子不存在）

```json
{
  "code": 404,
  "msg": "帖子不存在",
  "error": "record not found"
}
```

### 失败响应（服务器错误）

```json
{
  "code": 500,
  "msg": "获取帖子详情失败",
  "error": "database connection error"
}
```

---

## ⚙️ 扩展功能建议

### 1. 权限检查（草稿访问）

```go
// 检查当前用户是否是作者
func (s *postService) GetPostDetail(postID uint, userID *uint) (*PostDetail, error) {
    post, err := s.repo.FindByID(postID)
    if err != nil {
        return nil, err
    }

    // 只有已发布的帖子，任何人都能看
    // 草稿，只有作者能看
    if post.Status == "draft" && (userID == nil || *userID != post.AuthorID) {
        return nil, fmt.Errorf("无权限访问")
    }

    // ... 其他逻辑 ...
}
```

### 2. 浏览次数统计（使用 Redis）

```go
func (s *postService) GetPostDetail(postID uint) (*PostDetail, error) {
    // ... 前面的逻辑 ...

    // 异步增加浏览次数
    go func() {
        viewKey := fmt.Sprintf("post:%d:views", postID)
        databases.RedisClient.Incr(viewKey)
        
        // 定期同步到数据库
        // SQL: UPDATE posts SET view_count = view_count + 1 WHERE id = ?
    }()

    return detail, nil
}
```

### 3. 缓存优化（Redis 缓存详情）

```go
func (s *postService) GetPostDetail(postID uint) (*PostDetail, error) {
    // 1. 先从 Redis 查询
    cacheKey := fmt.Sprintf("post:%d:detail", postID)
    cached := databases.RedisClient.Get(cacheKey).Val()
    if cached != "" {
        // 反序列化 JSON
        var detail PostDetail
        json.Unmarshal([]byte(cached), &detail)
        return &detail, nil
    }

    // 2. 如果缓存不存在，查询数据库
    detail, err := s.GetPostDetailFromDB(postID)
    if err != nil {
        return nil, err
    }

    // 3. 存入缓存（6 小时过期）
    data, _ := json.Marshal(detail)
    databases.RedisClient.Set(cacheKey, string(data), 6*time.Hour)

    return detail, nil
}
```

### 4. 相关帖子推荐

```go
type PostDetailWithRelated struct {
    *PostDetail
    RelatedPosts []PostSummary `json:"related_posts"` // 同分类的其他帖子
}

func (s *postService) GetPostDetailWithRelated(postID uint) (*PostDetailWithRelated, error) {
    detail, err := s.GetPostDetail(postID)
    if err != nil {
        return nil, err
    }

    // 查询同分类的其他帖子
    related, _ := s.repo.ListByCategory(detail.Category.ID, 5, postID)
    
    return &PostDetailWithRelated{
        PostDetail:   detail,
        RelatedPosts: related,
    }, nil
}
```

---

## 🧪 测试用例

### 测试 1：获取存在的已发布帖子

```bash
curl "http://localhost:8080/api/posts/1"
```

**预期：**
- 状态码 200
- 返回完整的帖子内容

### 测试 2：获取不存在的帖子

```bash
curl "http://localhost:8080/api/posts/999"
```

**预期：**
- 状态码 404
- 返回错误信息

### 测试 3：获取草稿帖子（无权限）

```bash
curl "http://localhost:8080/api/posts/5"  # 假设 5 是草稿
```

**预期：**
- 状态码 404
- 隐藏帖子存在的事实

### 测试 4：验证返回的 JSON 格式

```bash
curl "http://localhost:8080/api/posts/1" | jq '.data.tags'
```

**预期：**
- 返回字符串数组 `["golang", "并发", ...]`

---

## 📊 性能对比

| 指标 | GetPostList | GetPostDetail |
|------|------------|---------------|
| 返回的数据量 | 小（摘要） | 大（完整） |
| 缓存时间 | 短（1 小时） | 长（6 小时） |
| 数据库查询 | 2 次 | 2 次 |
| 网络流量 | 低（列表）| 中（详情） |
| 响应时间 | 快（< 50ms） | 快（< 100ms） |

---

## 总结

**帖子详情功能的核心要点：**

1. ✅ **单条查询** - 返回一篇帖子，不分页
2. ✅ **完整内容** - 返回 Content，不摘要
3. ✅ **完整信息** - 包括 UpdatedAt，额外元数据
4. ✅ **状态检查** - 只显示已发布的帖子
5. ✅ **错误处理** - 不存在的帖子返回 404
6. ✅ **数据转换** - 标签解析、时间格式化
7. ✅ **可扩展性** - 易于添加权限、缓存、统计等功能

现在你有了完整的思路和实现代码，可以测试这个功能了！
