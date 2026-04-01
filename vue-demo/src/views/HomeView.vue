<template>
  <div class="home-page">
    <!-- 顶部导航栏 -->
   <header class="navbar">
      <div class="navbar-brand">📚 LearnHub</div>
      <nav class="navbar-nav">
        <router-link to="/" class="nav-link active">首页</router-link>
         <!-- 已登录显示"个人中心"，未登录显示"登录" -->
        <template v-if="isLoggedIn">
          <router-link to="/profile" class="nav-link">👤 {{ username }}</router-link>
          <button @click="handleLogout" class="btn-logout">退出</button>
        </template>
        <template v-else>
          <router-link to="/login" class="nav-link btn-primary">登录</router-link>
          <router-link to="/register" class="nav-link btn-outline">注册</router-link>
        </template>
      </nav>
    </header>

    <!-- 主体内容 -->
    <main class="main-content">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1>最新帖子</h1>
        <div class="header-actions" v-if="isLoggedIn">
          <router-link to="/ai-chat" class="btn-secondary">
            🤖 AI助手
          </router-link>
          <router-link to="/post/create" class="btn-primary">
            ✏️ 发布帖子
          </router-link>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- 错误状态 -->
      <div v-if="error" class="error-state">
        ⚠️ {{ error }}
        <button @click="loadPosts">重试</button>
      </div>

      <!-- 帖子列表 -->
      <div v-if="!loading && !error" class="post-list">
        <!-- 空状态 -->
        <div v-if="posts.length === 0" class="empty-state">
          <p>暂无帖子，<router-link to="/post/create">来发第一篇吧</router-link></p>
        </div>

        <!-- 帖子卡片 -->
        <article
          v-for="post in posts"
          :key="post.id"
          class="post-card"
          @click="goToPost(post.id)"
        >
          <!-- 置顶标记 -->
          <span v-if="post.is_pinned" class="badge-pinned">📌 置顶</span>

          <div class="post-card-body">
            <h2 class="post-title">{{ post.title }}</h2>
            <p class="post-summary">{{ post.summary }}</p>

            <!-- 标签 -->
            <div class="post-tags" v-if="post.tags && post.tags.length">
              <span
                v-for="tag in post.tags"
                :key="tag.id"
                class="tag"
              ># {{ tag.name }}</span>
            </div>
          </div>

          <footer class="post-card-footer">
            <!-- 作者信息 -->
            <div class="post-author">
              <img
                :src="post.author.avatar || defaultAvatar"
                :alt="post.author.username"
                class="avatar"
              />
              <span>{{ post.author.username }}</span>
            </div>

            <!-- 分类 -->
            <span v-if="post.category && post.category.name" class="category-badge">
              {{ post.category.icon }} {{ post.category.name }}
            </span>

            <!-- 统计数字 -->
            <div class="post-stats">
              <span>👁 {{ post.view_count }}</span>
              <span>👍 {{ post.like_count }}</span>
              <span>💬 {{ post.comment_count }}</span>
            </div>

            <!-- 发布时间 -->
            <span class="post-time">{{ formatTime(post.created_at) }}</span>
          </footer>
        </article>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="pagination">
        <button :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <span>第 {{ page }} 页 / 共 {{ totalPages }} 页（{{ total }} 篇）</span>
        <button :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
      </div>
    </main>
  </div>
</template>

<script>
import { getPostList } from '@/api/post'
import { logout } from '@/api/auth'

export default {
  name: 'HomeView',
  data() {
    return {
      posts: [],       // 当前页帖子列表
      total: 0,        // 总条数
      page: 1,         // 当前页码
      pageSize: 5,    // 每页条数
      loading: false,
      error: '',
      defaultAvatar: 'https://api.dicebear.com/7.x/thumbs/svg?seed=default',
      // 使用响应式数据存储登录状态，而不仅依赖 localStorage
      authToken: '',
      currentUsername: ''
    }
  },
  computed: {
    // 判断是否已登录
    isLoggedIn() {
      return !!this.authToken
    },
    // 获取用户名
    username() {
      return this.currentUsername || '我'
    },
    // 总页数
    totalPages() {
      return Math.ceil(this.total / this.pageSize)
    }
  },
  mounted() {
    // 初始化登录状态：从 localStorage 读取
    this.authToken = localStorage.getItem('auth_token') || ''
    this.currentUsername = localStorage.getItem('username') || ''
    
    // 输出日志用于调试
    console.log('HomeView mounted:', { authToken: this.authToken, currentUsername: this.currentUsername })
    
    // 监听 storage 变化（来自其他标签页或窗口）
    // 使用箭头函数确保 this 指向组件实例
    this._onStorageChange = (event) => this.onStorageChange(event)
    this._onLoginSuccess = (event) => this.onLoginSuccess(event)
    
    window.addEventListener('storage', this._onStorageChange)
    window.addEventListener('login-success', this._onLoginSuccess)
    
    this.loadPosts()
  },
  beforeUnmount() {
    // 组件卸载时移除事件监听
    window.removeEventListener('storage', this._onStorageChange)
    window.removeEventListener('login-success', this._onLoginSuccess)
  },
  methods: {
    // 加载帖子列表
    async loadPosts() {
      this.loading = true
      this.error = ''
      try {
        const res = await getPostList({ page: this.page, page_size: this.pageSize })
        // 后端统一响应格式：{ code: 0, data: { list, total, page, page_size } }
        if (res.code === 0) {
          this.posts = res.data.list || []
          this.total = res.data.total || 0
        } else {
          this.error = res.msg || '加载失败'
        }
      } catch (e) {
        this.error = e.message || '网络错误，请稍后重试'
      } finally {
        this.loading = false
      }
    },

    // 切换页码
    changePage(newPage) {
      this.page = newPage
      this.loadPosts()
      // 滚动到顶部
      window.scrollTo({ top: 0, behavior: 'smooth' })
    },

    // 跳转帖子详情
    goToPost(id) {
      this.$router.push(`/posts/${id}`)
    },

    // 退出登录
    handleLogout() {
      logout()
      // 清空本地状态
      this.authToken = ''
      this.currentUsername = ''
      // 退出后刷新页面（首页是公开的，刷新即可）
      this.$router.go(0)
    },

    // 监听 localStorage 变化
    onStorageChange(event) {
      console.log('Storage changed:', event.key, event.newValue)
      if (event.key === 'auth_token') {
        this.authToken = event.newValue || ''
      }
      if (event.key === 'username') {
        this.currentUsername = event.newValue || ''
      }
    },

    // 监听登录成功事件
    onLoginSuccess(event) {
      console.log('Login success event received:', event.detail)
      const { token, user } = event.detail
      this.authToken = token || ''
      this.currentUsername = user?.username || ''
    },

    // 格式化时间：显示为"X 分钟前"或日期
    formatTime(isoString) {
      if (!isoString) return ''
      const date = new Date(isoString)
      const now = new Date()
      const diff = Math.floor((now - date) / 1000) // 秒差

      if (diff < 60) return '刚刚'
      if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`
      if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`
      if (diff < 86400 * 7) return `${Math.floor(diff / 86400)} 天前`

      // 超过 7 天显示完整日期
      return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
    }
  }
}
</script>

<style scoped>
/* ========== 导航栏 ========== */
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.navbar-brand {
  font-size: 20px;
  font-weight: 700;
  color: #1a73e8;
}

.navbar-nav {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-link {
  text-decoration: none;
  color: #555;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background 0.2s;
  font-size: 14px;
}

.nav-link:hover,
.nav-link.active {
  background: #f0f4ff;
  color: #1a73e8;
}

.btn-primary {
  background: #1a73e8;
  color: #fff !important;
  padding: 6px 16px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  border: none;
}

.btn-primary:hover {
  background: #1558b0;
}

.btn-outline {
  border: 1px solid #1a73e8;
  color: #1a73e8 !important;
  padding: 5px 15px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
}

.btn-logout {
  background: none;
  border: 1px solid #ddd;
  color: #888;
  padding: 5px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}

.btn-logout:hover {
  border-color: #f56565;
  color: #f56565;
}

/* ========== 主体内容 ========== */
.main-content {
  max-width: 860px;
  margin: 32px auto;
  padding: 0 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 22px;
  color: #222;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.btn-secondary {
  background: #9c27b0;
  color: #fff !important;
  padding: 8px 16px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  border: none;
  transition: background 0.2s;
}

.btn-secondary:hover {
  background: #7b1fa2;
}

/* ========== 状态 ========== */
.loading-state {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #888;
  padding: 40px 0;
  justify-content: center;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #ddd;
  border-top-color: #1a73e8;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state {
  text-align: center;
  padding: 32px;
  color: #e53e3e;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #aaa;
  font-size: 15px;
}

/* ========== 帖子卡片 ========== */
.post-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.post-card {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 10px;
  padding: 20px;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.15s;
  position: relative;
}

.post-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.badge-pinned {
  position: absolute;
  top: 14px;
  right: 14px;
  font-size: 12px;
  color: #e67e22;
  background: #fff8f0;
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid #f0d0a0;
}

.post-title {
  font-size: 17px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px;
  line-height: 1.4;
}

.post-summary {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0 0 12px;
  /* 最多显示 3 行 */
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  font-size: 12px;
  color: #1a73e8;
  background: #e8f0fe;
  padding: 2px 8px;
  border-radius: 12px;
}

.post-card-footer {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  font-size: 13px;
  color: #888;
  border-top: 1px solid #f0f0f0;
  padding-top: 12px;
  margin-top: 4px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #555;
  font-weight: 500;
}

.avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.category-badge {
  background: #f5f5f5;
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
}

.post-stats {
  display: flex;
  gap: 10px;
  margin-left: auto;
}

.post-time {
  font-size: 12px;
  color: #bbb;
}

/* ========== 分页 ========== */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
  padding-bottom: 40px;
  color: #555;
  font-size: 14px;
}

.pagination button {
  padding: 6px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  border-color: #1a73e8;
  color: #1a73e8;
}

.pagination button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
