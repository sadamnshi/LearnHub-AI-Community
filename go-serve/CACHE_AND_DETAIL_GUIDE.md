# 🚀 缓存模式 + 详情页完整实现指南

## 📚 核心概念讲解

### 1. 旁路缓存模式（Cache-Aside Pattern）

#### 什么是旁路缓存？

```
┌─────────────┐
│   前端请求   │
└──────┬──────┘
       │
       ▼
┌──────────────────────┐
│  1. 查询缓存          │
│  Redis 中有吗？       │
└──┬────────────────┬──┘
   │ 有             │ 没有
   ▼                ▼
┌────────┐    ┌─────────────────┐
│返回缓存 │    │ 查询数据库       │
└────┬───┘    └────────┬────────┘
     │                 │
     │                 ▼
     │         ┌──────────────────┐
     │         │ 2. 存入缓存      │
     │         │ 数据写到 Redis   │
     │         └────────┬─────────┘
     │                  │
     └──────────┬───────┘
                │
                ▼
         ┌────────────────┐
         │ 返回给前端     │
         └────────────────┘
```

**核心思想：**
1. 先查询缓存，有就直接返回（快！）
2. 缓存没有，查询数据库，然后存入缓存（为下次做准备）
3. 返回数据给前端

**优点：**
- ✅ 解耦：缓存和数据库分离
- ✅ 灵活：可以独立管理缓存策略
- ✅ 简单：实现起来最容易

**缺点：**
- ❌ 缓存穿透：如果查询的数据不存在，会频繁查询数据库
- ❌ 缓存雪崩：所有缓存同时过期，导致大量请求打到数据库

---

### 2. 旁路缓存 vs 其他模式对比

```
┌──────────────────┬────────────────┬──────────────┬──────────────┐
│      模式名      │   实现复杂度   │    性能      │    安全性    │
├──────────────────┼────────────────┼──────────────┼──────────────┤
│ Cache-Aside      │  ⭐⭐（简单） │ ⭐⭐⭐    │ ⭐⭐       │
│ （旁路缓存）     │                │  （中等）    │              │
├──────────────────┼────────────────┼──────────────┼──────────────┤
│ Write-Through    │  ⭐⭐⭐       │ ⭐⭐⭐⭐  │ ⭐⭐⭐⭐  │
│ （写穿）         │  （复杂）      │  （高）      │  （高）      │
├──────────────────┼────────────────┼──────────────┼──────────────┤
│ Write-Behind     │  ⭐⭐⭐⭐     │ ⭐⭐⭐⭐⭐ │ ⭐⭐       │
│ （写回）         │  （很复杂）    │  （最高）    │  （低）      │
└──────────────────┴────────────────┴──────────────┴──────────────┘

适用场景：
- Cache-Aside：适合读多写少的场景（论坛帖子）✅
- Write-Through：适合对一致性要求高的场景（订单、支付）
- Write-Behind：适合高并发写入场景（日志系统）
```

---

## 🛠️ 实现思路详解

### 第一部分：Controller 层 - 获取帖子详情

#### 设计思路

```
HTTP 请求
    │
    ├─ GET /api/posts/:id
    │
    ▼
PostController.GetPostDetail()
    │
    ├─ 1️⃣ 从 URL 提取帖子 ID
    │     c.Param("id") → "123"
    │
    ├─ 2️⃣ 参数验证
    │     strconv.ParseUint() → uint(123)
    │
    ├─ 3️⃣ 调用服务层
    │     ctl.service.GetPostDetail(uint(123))
    │
    ├─ 4️⃣ 错误处理
    │     err != nil → 返回 404
    │
    └─ 5️⃣ 返回 JSON 响应
          c.JSON(200, gin.H{...})
```

#### 代码实现步骤

**Step 1: 在 Controller 中添加 GetPostDetail 方法**

```go
// 在 PostController 中添加这个方法

func (ctl *PostController) GetPostDetail(c *gin.Context) {
    // 第 1 步：从 URL 参数中提取帖子 ID
    // GET /api/posts/123
    // c.Param("id") 会提取路径中的 :id 参数
    // 例如：URL 是 /api/posts/123，那么 c.Param("id") 返回 "123"
    postIDStr := c.Param("id")
    
    // 第 2 步：转换字符串为数字
    // strconv.ParseUint(string, 进制, 位数)
    // 参数说明：
    //   - postIDStr: "123"
    //   - 10: 十进制
    //   - 64: 转换为 uint64（Go 的标准做法）
    // 返回：(uint64, error)
    postID, err := strconv.ParseUint(postIDStr, 10, 64)
    
    // 如果转换失败（用户输入非数字），返回 400 Bad Request
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 400,
            "msg":  "无效的帖子 ID",
        })
        return
    }
    
    // 第 3 步：调用服务层获取详情
    // 服务层会：
    //   1. 从缓存查询（旁路缓存）
    //   2. 查询数据库
    //   3. 存入缓存
    //   4. 返回结果
    detail, err := ctl.service.GetPostDetail(uint(postID))
    
    // 第 4 步：错误处理
    if err != nil {
        // 如果服务层返回错误，说明帖子不存在或其他问题
        c.JSON(http.StatusNotFound, gin.H{
            "code": 404,
            "msg":  "帖子不存在或已被删除",
        })
        return
    }
    
    // 第 5 步：返回成功响应
    // 统一的响应格式
    c.JSON(http.StatusOK, gin.H{
        "code": 0,
        "msg":  "success",
        "data": detail,
    })
}
```

---

### 第二部分：Service 层 - 旁路缓存实现

#### 缓存设计

```go
// Redis 中的 Key 设计
// 格式：post:{id}:detail
// 示例：
//   - post:1:detail    → 第 1 篇帖子的详情
//   - post:123:detail  → 第 123 篇帖子的详情

// 缓存过期时间设计
// - 短期缓存（1 小时）：用于列表页（更新频繁）
// - 中期缓存（6 小时）：用于详情页（更新较少）
// - 长期缓存（24 小时）：用于静态内容（几乎不更新）
```

#### 旁路缓存实现步骤

```go
// 伪代码流程

func (s *postService) GetPostDetail(postID uint) (*PostDetail, error) {
    
    // ════════════════════════════════════════════════════════════
    // 步骤 1：尝试从缓存获取
    // ════════════════════════════════════════════════════════════
    cacheKey := fmt.Sprintf("post:%d:detail", postID)
    cachedData, err := s.cache.Get(cacheKey)  // Redis GET
    
    // 如果缓存存在，直接返回（快！）
    if err == nil && cachedData != "" {
        // 反序列化 JSON 字符串为 PostDetail 对象
        var detail PostDetail
        json.Unmarshal([]byte(cachedData), &detail)
        return &detail, nil  // 缓存命中 ✓
    }
    
    // ════════════════════════════════════════════════════════════
    // 步骤 2：缓存未命中，查询数据库
    // ════════════════════════════════════════════════════════════
    post, err := s.repo.FindByID(postID)  // 数据库查询
    if err != nil {
        return nil, err  // 帖子不存在
    }
    
    // 业务逻辑检查...
    if post.Status != "published" {
        return nil, fmt.Errorf("帖子不存在或已被删除")
    }
    
    // 转换为 DTO...
    detail := buildPostDetail(post)  // 这个函数在下面定义
    
    // ════════════════════════════════════════════════════════════
    // 步骤 3：将结果存入缓存
    // ════════════════════════════════════════════════════════════
    // 序列化为 JSON 字符串
    data, _ := json.Marshal(detail)
    
    // 存入 Redis，6 小时后过期
    s.cache.Set(cacheKey, string(data), 6*time.Hour)
    
    // ════════════════════════════════════════════════════════════
    // 步骤 4：返回结果
    // ════════════════════════════════════════════════════════════
    return detail, nil
}
```

---

### 第三部分：Vue 前端 - 详情页面

#### 页面设计

```
┌─────────────────────────────────────────┐
│              帖子详情页                  │
├─────────────────────────────────────────┤
│ ⬅️ 返回  |  分享  📤  |  点赞 ❤️         │
├─────────────────────────────────────────┤
│ 标题：Go 语言并发编程最佳实践              │
│ 作者：Tom  📅 2026-03-15  👁️ 浏览数      │
│ 分类：💻 技术分享                        │
│ 标签：[golang] [并发] [goroutine]       │
├─────────────────────────────────────────┤
│ 完整内容（Markdown 渲染）                │
│                                        │
│ 在 Go 中，并发是通过 goroutine 和      │
│ channel 实现的。                       │
│                                        │
│ Goroutine 是一个轻量级的线程，由 Go   │
│ 运行时管理。                          │
│                                        │
│ ... 完整内容 ...                      │
│                                        │
├─────────────────────────────────────────┤
│ 💬 评论区（可选）                       │
├─────────────────────────────────────────┤
│ 相关推荐：                              │
│  - [Go] 如何优化并发程序性能            │
│  - [Go] channel 深度解析               │
└─────────────────────────────────────────┘
```

#### Vue 3 Composition API 实现

```javascript
// views/PostDetailView.vue

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchPostDetail } from '@/api/post'

const route = useRoute()
const router = useRouter()

// 数据状态
const post = ref(null)           // 帖子详情数据
const loading = ref(true)        // 加载状态
const error = ref(null)          // 错误信息

// 生命周期：组件挂载时加载数据
onMounted(async () => {
  try {
    loading.value = true
    const postId = route.params.id
    
    // 调用 API 获取详情
    post.value = await fetchPostDetail(postId)
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})

// 返回列表
const goBack = () => {
  router.back()  // 或 router.push('/')
}

// 点赞功能（可选）
const likePost = () => {
  // TODO: 实现点赞逻辑
}
</script>

<template>
  <!-- 加载中 -->
  <div v-if="loading" class="loading">加载中...</div>
  
  <!-- 错误提示 -->
  <div v-else-if="error" class="error">{{ error }}</div>
  
  <!-- 详情内容 -->
  <div v-else class="post-detail">
    <!-- 工具栏 -->
    <div class="toolbar">
      <button @click="goBack">⬅️ 返回</button>
      <button @click="likePost">❤️ 点赞</button>
    </div>
    
    <!-- 标题和基本信息 -->
    <header class="post-header">
      <h1>{{ post.title }}</h1>
      <div class="meta">
        <span>作者：{{ post.author.username }}</span>
        <span>📅 {{ post.created_at }}</span>
        <span>👁️ {{ post.view_count }} 浏览</span>
      </div>
    </header>
    
    <!-- 分类和标签 -->
    <div class="categories-tags">
      <span class="category">{{ post.category.icon }} {{ post.category.name }}</span>
      <span v-for="tag in post.tags" :key="tag" class="tag">
        {{ tag }}
      </span>
    </div>
    
    <!-- 文章内容 -->
    <article class="post-content">
      <div v-html="markdownToHtml(post.content)"></div>
    </article>
  </div>
</template>

<style scoped>
.loading, .error {
  text-align: center;
  padding: 20px;
}

.post-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.post-header {
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
  margin-bottom: 15px;
}

.post-header h1 {
  margin: 0 0 10px 0;
}

.meta {
  display: flex;
  gap: 20px;
  color: #666;
  font-size: 14px;
}

.categories-tags {
  margin-bottom: 20px;
}

.category {
  display: inline-block;
  background: #f0f0f0;
  padding: 5px 10px;
  border-radius: 3px;
  margin-right: 10px;
}

.tag {
  display: inline-block;
  background: #e8f4f8;
  padding: 5px 10px;
  border-radius: 3px;
  margin-right: 5px;
}

.post-content {
  line-height: 1.8;
  font-size: 16px;
}
</style>
```

---

### 第四部分：API 层 - 前端请求

```javascript
// api/post.js

import request from '@/utils/request'

/**
 * 获取帖子详情
 * @param {number} postId - 帖子 ID
 * @returns {Promise<PostDetail>}
 */
export const fetchPostDetail = async (postId) => {
  const response = await request.get(`/api/posts/${postId}`)
  
  if (response.code !== 0) {
    throw new Error(response.msg || '获取帖子详情失败')
  }
  
  return response.data
}

/**
 * 获取帖子列表
 * @param {object} params - 查询参数
 * @returns {Promise<PostListResult>}
 */
export const fetchPostList = async (params) => {
  const response = await request.get('/api/posts/home', { params })
  
  if (response.code !== 0) {
    throw new Error(response.msg || '获取帖子列表失败')
  }
  
  return response.data
}
```

---

### 第五部分：Router 配置

```javascript
// router/index.js

import PostDetailView from '@/views/PostDetailView.vue'

const routes = [
  // ... 其他路由 ...
  
  {
    path: '/posts/:id',
    name: 'PostDetail',
    component: PostDetailView,
    // 可选：添加路由元信息
    meta: {
      title: '帖子详情',
      requiresAuth: false  // 不需要登录
    }
  }
]
```

---

## 💾 Redis 缓存操作详解

### Key 的设计原则

```
// ❌ 不好的设计
post:detail:1        // 不清楚是什么数据
post_1_detail        // 使用下划线，容易混淆
post-detail-1        // 顺序不清晰

// ✅ 好的设计
post:{id}:detail     // 明确的类型和内容
user:{id}:profile    // 用户资料
category:{id}:info   // 分类信息
post:{id}:comments   // 帖子评论列表
```

### 缓存操作代码

```go
package cache

import (
    "github.com/redis/go-redis/v9"
    "time"
)

type CacheClient struct {
    client *redis.Client
}

// Get 从 Redis 获取数据
func (c *CacheClient) Get(key string) (string, error) {
    return c.client.Get(context.Background(), key).Result()
}

// Set 存入 Redis，指定过期时间
func (c *CacheClient) Set(key string, value string, expiration time.Duration) error {
    return c.client.Set(context.Background(), key, value, expiration).Err()
}

// Delete 删除缓存
func (c *CacheClient) Delete(key string) error {
    return c.client.Del(context.Background(), key).Err()
}

// Exists 检查 Key 是否存在
func (c *CacheClient) Exists(key string) (bool, error) {
    result, err := c.client.Exists(context.Background(), key).Result()
    return result > 0, err
}
```

---

## 🎯 完整的数据流向

```
用户点击帖子链接
  │
  ▼ (前端)
router.push(`/posts/${postId}`)
  │
  ├─ URL 变化：/posts/123
  ├─ PostDetailView 组件挂载
  └─ onMounted() 触发
       │
       ▼ fetchPostDetail(123)
       HTTP GET /api/posts/123
       │
       ▼ (后端 - Controller)
       PostController.GetPostDetail()
       ├─ c.Param("id") → "123"
       ├─ ParseUint("123") → uint(123)
       └─ service.GetPostDetail(123)
            │
            ▼ (后端 - Service)
            GetPostDetail(123)
            ├─ Step 1: 查询缓存
            │  └─ Redis GET post:123:detail
            │     ├─ 命中 → 返回缓存数据
            │     └─ 未命中 → 继续
            │
            ├─ Step 2: 查询数据库
            │  └─ SELECT * FROM posts WHERE id = 123
            │     ├─ 数据库查询 + Preload Author, Category
            │     └─ 返回 Post 模型
            │
            ├─ Step 3: 业务逻辑检查
            │  └─ Status == "published"?
            │
            ├─ Step 4: 转换为 DTO
            │  └─ Model → PostDetail
            │
            ├─ Step 5: 存入缓存
            │  └─ Redis SET post:123:detail "{...json}" EX 21600
            │     (21600 秒 = 6 小时)
            │
            └─ Step 6: 返回结果
               └─ &PostDetail{...}
            │
            ▼
       PostController 收到结果
       ├─ 检查错误
       └─ c.JSON(200, gin.H{
            "code": 0,
            "msg": "success",
            "data": detail
          })
       │
       ▼ HTTP 200 + JSON 响应
       │
       ▼ (前端)
       fetchPostDetail() 返回数据
       ├─ post.value = response.data
       ├─ loading.value = false
       └─ 模板渲染
            │
            ▼
       用户看到完整的帖子详情页面
```

---

## 🧪 测试场景

### 场景 1：首次访问（缓存未命中）

```
时间流程：
├─ T0: 用户访问 /posts/1
├─ T1: Controller 处理请求
├─ T2: Service 查询缓存 → 未命中
├─ T3: Service 查询数据库 → 成功
├─ T4: Service 存入缓存
├─ T5: 返回响应给前端（~200ms）
└─ T6: 用户看到详情页面

Redis 状态：
├─ Before: post:1:detail 不存在
└─ After: post:1:detail → JSON 数据（6小时过期）
```

### 场景 2：第二次访问（缓存命中）

```
时间流程：
├─ T0: 用户再次访问 /posts/1
├─ T1: Controller 处理请求
├─ T2: Service 查询缓存 → 命中！ ✓
├─ T3: 直接反序列化并返回
└─ T4: 返回响应给前端（~10ms，快 20 倍！）

性能对比：
├─ 首次访问：200ms（数据库查询）
├─ 第二次访问：10ms（缓存命中）
└─ 性能提升：20 倍！
```

### 场景 3：帖子不存在

```
HTTP GET /api/posts/999

时间流程：
├─ T1: Service 查询缓存 → 未命中
├─ T2: Service 查询数据库 → 记录不存在
├─ T3: 返回错误

返回响应：
HTTP 404
{
  "code": 404,
  "msg": "帖子不存在或已被删除"
}
```

---

## ⚡ 性能对比

```
┌──────────────────┬─────────────────┬──────────────────┬─────────┐
│ 场景             │ 缓存前响应时间   │ 缓存后响应时间   │ 提升   │
├──────────────────┼─────────────────┼──────────────────┼─────────┤
│ 首次访问（未命中）│ 200ms（DB查询） │ 200ms + 缓存存储 │ ✓      │
│                  │                 │                  │        │
│ 后续访问（命中）  │ N/A             │ 10ms（缓存查询） │ 20倍 ⚡ │
│                  │                 │                  │        │
│ 1000 个并发请求  │ 1000×200ms      │ ~20ms（缓存）    │ 100倍 ⚡⚡│
│ （热点数据）      │ = 200 秒        │                  │        │
└──────────────────┴─────────────────┴──────────────────┴─────────┘
```

---

## 📋 总结核心要点

### Controller 层
- ✅ 从 URL 参数提取 ID: `c.Param("id")`
- ✅ 参数验证: `strconv.ParseUint()`
- ✅ 调用服务层：`ctl.service.GetPostDetail()`
- ✅ 错误处理：返回 404
- ✅ 响应格式：统一的 `gin.H{}`

### Service 层（旁路缓存）
- ✅ Step 1: 查询缓存 `cache.Get(key)`
- ✅ Step 2: 缓存未命中，查询数据库
- ✅ Step 3: 业务逻辑检查
- ✅ Step 4: 转换为 DTO
- ✅ Step 5: 存入缓存 `cache.Set(key, value, duration)`
- ✅ Step 6: 返回结果

### Vue 前端
- ✅ 路由参数：`:id`
- ✅ 生命周期：`onMounted()` 获取数据
- ✅ 状态管理：`ref()`
- ✅ 条件渲染：`v-if/v-else`
- ✅ 数据绑定：`{{ post.title }}`

### Redis 缓存
- ✅ Key 设计：`post:{id}:detail`
- ✅ 过期时间：6 小时
- ✅ 序列化：`JSON.Marshal/Unmarshal`
- ✅ 操作：`Get/Set/Delete`

现在让我为你实现这些代码！
