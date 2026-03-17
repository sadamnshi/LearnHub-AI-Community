<template>
  <!-- 帖子详情页面容器 -->
  <div class="post-detail-page">
    <!-- ═════════════════════════════════════════════════════════════ -->
    <!-- 加载中状态 - 显示加载动画 -->
    <!-- ═════════════════════════════════════════════════════════════ -->
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>加载帖子详情中...</p>
    </div>

    <!-- ═════════════════════════════════════════════════════════════ -->
    <!-- 错误提示 - 当获取数据失败时显示 -->
    <!-- ═════════════════════════════════════════════════════════════ -->
    <div v-else-if="error" class="error-container">
      <div class="error-icon">⚠️</div>
      <h2>出错了</h2>
      <p>{{ error }}</p>
      <button class="btn-primary" @click="goBack">⬅️ 返回列表</button>
    </div>

    <!-- ═════════════════════════════════════════════════════════════ -->
    <!-- 帖子内容 - 成功加载后显示 -->
    <!-- ═════════════════════════════════════════════════════════════ -->
    <article v-else class="post-detail">
      <!-- 工具栏 - 返回按钮和其他操作 -->
      <div class="toolbar">
        <button class="btn-back" @click="goBack" title="返回列表">
          ⬅️ 返回
        </button>
        <div class="toolbar-spacer"></div>
        <button class="btn-action" title="分享">
          📤 分享
        </button>
        <button class="btn-action like-btn" @click="toggleLike" :class="{ liked: isLiked }" title="点赞">
          {{ isLiked ? '❤️' : '🤍' }} {{ likeCount }}
        </button>
      </div>

      <!-- 帖子头部信息 -->
      <header class="post-header">
        <!-- 标题 -->
        <h1 class="post-title">{{ post.title }}</h1>

        <!-- 元信息：作者、时间、浏览次数等 -->
        <div class="post-meta">
          <!-- 作者信息 -->
          <div class="meta-item author-info">
            <img :src="post.author.avatar" :alt="post.author.username" class="avatar" />
            <span class="author-name">{{ post.author.username }}</span>
          </div>

          <!-- 发布时间 -->
          <div class="meta-item">
            <span class="label">📅</span>
            <time :datetime="post.created_at">{{ formatDate(post.created_at) }}</time>
          </div>

          <!-- 更新时间（仅在内容被更新过时显示） -->
          <div v-if="isUpdated" class="meta-item">
            <span class="label">✏️</span>
            <time :datetime="post.updated_at">{{ formatDate(post.updated_at) }}</time>
          </div>

          <!-- 浏览次数 -->
          <div class="meta-item">
            <span class="label">👁️</span>
            <span>{{ viewCount }} 浏览</span>
          </div>
        </div>
      </header>

      <!-- 分类和标签 -->
      <div class="categories-tags-section">
        <!-- 分类 - 单个分类 -->
        <div class="category">
          <span class="category-icon">{{ post.category.icon }}</span>
          <span class="category-name">{{ post.category.name }}</span>
        </div>

        <!-- 标签 - 多个标签 -->
        <div v-if="post.tags && post.tags.length > 0" class="tags">
          <span v-for="tag in post.tags" :key="tag" class="tag">
            #{{ tag }}
          </span>
        </div>
      </div>

      <!-- 分隔线 -->
      <hr class="divider" />

      <!-- 文章主体内容 -->
      <section class="post-content">
        <!-- 
          使用 v-html 渲染 Markdown 转换后的 HTML
          注意：确保内容来自可信源，避免 XSS 攻击
          在生产环境中应该使用专门的 Markdown 解析库
          例如：marked.js、markdown-it 等
        -->
        <div class="content-body" v-html="renderedContent"></div>
      </section>

      <!-- 分隔线 -->
      <hr class="divider" />

      <!-- 帖子底部 - 点赞和分享 -->
      <footer class="post-footer">
        <div class="footer-actions">
          <button 
            class="action-btn like" 
            @click="toggleLike"
            :class="{ liked: isLiked }"
          >
            {{ isLiked ? '❤️' : '🤍' }} 
            {{ isLiked ? '已点赞' : '点赞' }}
            <span class="count">({{ likeCount }})</span>
          </button>
          <button class="action-btn share">
            📤 分享
          </button>
          <button class="action-btn comment">
            💬 评论
            <span class="count">(0)</span>
          </button>
        </div>
      </footer>
    </article>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchPostDetail } from '@/api/post'

// ═════════════════════════════════════════════════════════════════════════
// 路由和页面导航
// ═════════════════════════════════════════════════════════════════════════
const route = useRoute()
const router = useRouter()

// ═════════════════════════════════════════════════════════════════════════
// 响应式数据
// ═════════════════════════════════════════════════════════════════════════

// post 帖子详情数据
// 初始值为 null，加载完成后会被赋值
// 结构：{ id, title, content, author, category, tags, created_at, updated_at }
const post = ref(null)

// loading 加载状态
// true: 正在加载，显示加载动画
// false: 加载完成，显示内容
const loading = ref(true)

// error 错误信息
// null: 没有错误
// string: 错误提示文本
const error = ref(null)

// isLiked 当前用户是否已点赞
// 这个字段在实际应用中应该从后端返回
const isLiked = ref(false)

// likeCount 点赞数
// 初始值为 0，应该从 post 数据中读取
const likeCount = ref(0)

// viewCount 浏览次数
// 初始值为 0，应该从 post 数据中读取
const viewCount = ref(0)

// ═════════════════════════════════════════════════════════════════════════
// 计算属性
// ═════════════════════════════════════════════════════════════════════════

// isUpdated 是否被更新过
// 比较 created_at 和 updated_at，如果不同则说明被更新过
const isUpdated = computed(() => {
  if (!post.value) return false
  return post.value.created_at !== post.value.updated_at
})

// renderedContent 渲染后的 HTML 内容
// 将 Markdown 文本转换为 HTML
// 在实际应用中，应该使用专门的 Markdown 解析库
// 这里为了简单起见，只做基本的行尾换行处理
const renderedContent = computed(() => {
  if (!post.value) return ''
  
  let html = post.value.content
  
  // 简单的 Markdown 处理：
  // 1. 替换换行符为 <br>（为了保留行尾）
  html = html.replace(/\n/g, '<br>')
  
  // 实际应用中应该这样做：
  // import { marked } from 'marked'
  // return marked(post.value.content)
  
  return html
})

// ═════════════════════════════════════════════════════════════════════════
// 方法
// ═════════════════════════════════════════════════════════════════════════

// goBack 返回列表页面
const goBack = () => {
  // 方式 1: 使用 router.back() - 返回历史记录前一个页面
  router.back()
  
  // 方式 2: 返回到固定的首页
  // router.push('/')
}

// toggleLike 切换点赞状态
const toggleLike = async () => {
  // 切换点赞状态
  isLiked.value = !isLiked.value
  
  // 更新点赞数
  if (isLiked.value) {
    likeCount.value++
  } else {
    likeCount.value--
  }
  
  // 在实际应用中，这里应该调用后端 API：
  // await likePost(post.value.id)
}

// formatDate 格式化日期
// 参数：ISO 8601 格式的日期字符串
// 返回：更易读的格式
const formatDate = (dateStr) => {
  // dateStr 格式：2026-03-15 10:30:45
  // 直接返回，因为后端已经格式化了
  return dateStr
  
  // 如果需要进一步处理，可以这样：
  // const date = new Date(dateStr)
  // return date.toLocaleDateString('zh-CN')
}

// loadPostDetail 加载帖子详情
// 这是核心方法，从后端获取数据
const loadPostDetail = async () => {
  try {
    // 重置错误信息
    error.value = null
    loading.value = true

    // 从 URL 参数中获取帖子 ID
    // route.params.id 对应 URL 中的 :id 参数
    // 例如：/posts/123 中，id = "123"
    const postId = route.params.id

    // 调用 API 获取帖子详情
    // fetchPostDetail() 是在 api/post.js 中定义的方法
    // 它会发送 GET 请求到 /api/posts/{postId}
    const data = await fetchPostDetail(postId)

    // 保存数据到响应式变量
    post.value = data

    // 初始化点赞和浏览计数
    // 在实际应用中，这些值应该从 post 数据中读取
    likeCount.value = data.like_count || 0
    viewCount.value = data.view_count || 0

  } catch (err) {
    // 捕获并保存错误信息
    error.value = err.message || '获取帖子详情失败，请重试'
    console.error('加载帖子详情出错：', err)

  } finally {
    // 无论成功还是失败，都关闭加载动画
    loading.value = false
  }
}

// ═════════════════════════════════════════════════════════════════════════
// 生命周期钩子
// ═════════════════════════════════════════════════════════════════════════

// onMounted 组件挂载时
// 此时 DOM 已经就绪，可以进行网络请求、修改 document 等操作
onMounted(() => {
  // 加载帖子详情
  loadPostDetail()
})
</script>

<style scoped>
/* ═════════════════════════════════════════════════════════════════════════ */
/* 页面容器 */
/* ═════════════════════════════════════════════════════════════════════════ */

.post-detail-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 加载状态 */
/* ═════════════════════════════════════════════════════════════════════════ */

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  gap: 20px;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-container p {
  color: #666;
  font-size: 18px;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 错误状态 */
/* ═════════════════════════════════════════════════════════════════════════ */

.error-container {
  max-width: 500px;
  margin: 100px auto;
  padding: 40px;
  background: white;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.error-icon {
  font-size: 60px;
  margin-bottom: 20px;
}

.error-container h2 {
  color: #e74c3c;
  margin-bottom: 10px;
}

.error-container p {
  color: #666;
  margin-bottom: 30px;
  line-height: 1.6;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 帖子容器 */
/* ═════════════════════════════════════════════════════════════════════════ */

.post-detail {
  max-width: 900px;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 工具栏 */
/* ═════════════════════════════════════════════════════════════════════════ */

.toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 30px;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.toolbar-spacer {
  flex: 1;
}

.btn-back,
.btn-action {
  padding: 8px 16px;
  border: none;
  background: none;
  color: #666;
  cursor: pointer;
  font-size: 14px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.btn-back:hover,
.btn-action:hover {
  background: #e9ecef;
  color: #333;
}

.like-btn {
  color: #e74c3c;
}

.like-btn.liked {
  color: #e74c3c;
  background: #ffe6e6;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 帖子头部 */
/* ═════════════════════════════════════════════════════════════════════════ */

.post-header {
  padding: 40px 30px 20px;
}

.post-title {
  font-size: 36px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 20px 0;
  line-height: 1.4;
}

.post-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 20px;
  color: #666;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e9ecef;
  border: 2px solid #f0f0f0;
}

.author-name {
  color: #3498db;
  font-weight: 600;
}

.label {
  font-size: 16px;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 分类和标签 */
/* ═════════════════════════════════════════════════════════════════════════ */

.categories-tags-section {
  padding: 0 30px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.category {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f0f7ff;
  border: 1px solid #b3d9ff;
  border-radius: 20px;
  color: #1890ff;
  font-size: 13px;
  font-weight: 600;
}

.category-icon {
  font-size: 16px;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  display: inline-block;
  padding: 6px 12px;
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 16px;
  font-size: 13px;
  color: #555;
  transition: all 0.3s ease;
  cursor: pointer;
}

.tag:hover {
  background: #e8e8e8;
  border-color: #999;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 分隔线 */
/* ═════════════════════════════════════════════════════════════════════════ */

.divider {
  margin: 20px 0;
  border: none;
  border-top: 1px solid #e9ecef;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 文章内容 */
/* ═════════════════════════════════════════════════════════════════════════ */

.post-content {
  padding: 30px;
  font-size: 16px;
  line-height: 1.8;
  color: #333;
}

.content-body {
  word-break: break-word;
  overflow-wrap: break-word;
}

.content-body p {
  margin: 16px 0;
}

.content-body h1,
.content-body h2,
.content-body h3 {
  margin: 24px 0 16px 0;
  color: #1a1a1a;
  font-weight: 700;
}

.content-body h1 { font-size: 28px; }
.content-body h2 { font-size: 24px; }
.content-body h3 { font-size: 20px; }

.content-body code {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  color: #e74c3c;
}

.content-body pre {
  background: #2d2d2d;
  color: #f8f8f2;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 16px 0;
}

.content-body blockquote {
  border-left: 4px solid #ddd;
  margin: 16px 0;
  padding: 12px 16px;
  background: #f9f9f9;
  color: #666;
  font-style: italic;
}

.content-body ul,
.content-body ol {
  margin: 16px 0;
  padding-left: 30px;
}

.content-body li {
  margin: 8px 0;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 底部操作栏 */
/* ═════════════════════════════════════════════════════════════════════════ */

.post-footer {
  padding: 30px;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
}

.footer-actions {
  display: flex;
  gap: 16px;
}

.action-btn {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.action-btn:hover {
  border-color: #999;
  background: #f9f9f9;
}

.action-btn.like.liked {
  background: #ffe6e6;
  border-color: #e74c3c;
  color: #e74c3c;
}

.action-btn .count {
  color: #999;
  font-weight: normal;
}

/* ═════════════════════════════════════════════════════════════════════════ */
/* 响应式设计 */
/* ═════════════════════════════════════════════════════════════════════════ */

@media (max-width: 768px) {
  .post-detail-page {
    padding: 10px;
  }

  .post-title {
    font-size: 24px;
  }

  .post-header,
  .post-content,
  .post-footer,
  .categories-tags-section {
    padding: 20px;
  }

  .post-meta {
    gap: 12px;
    font-size: 12px;
  }

  .footer-actions {
    flex-direction: column;
  }

  .toolbar {
    padding: 12px 20px;
  }
}
</style>
