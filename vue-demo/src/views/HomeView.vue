<template>
  <div class="home-page">
    <!-- Glass Navbar -->
    <header class="navbar">
      <div class="navbar-inner">
        <router-link to="/" class="navbar-brand">
          <span class="brand-icon"> </span>
          <span class="brand-text">LearnHub</span>
        </router-link>
        <nav class="navbar-nav">
          <router-link to="/" class="nav-link active">首页</router-link>
          <template v-if="isLoggedIn">
            <router-link to="/profile" class="nav-link">{{ username }}</router-link>
            <button @click="handleLogout" class="btn-ghost">退出</button>
          </template>
          <template v-else>
            <router-link to="/login" class="nav-link btn-accent">登录</router-link>
            <router-link to="/register" class="nav-link btn-glass-outline">注册</router-link>
          </template>
        </nav>
      </div>
    </header>

    <!-- Main Content -->
    <main class="main-content">
      <div class="page-header">
        <div>
          <h1 class="page-title">最新帖子</h1>
          <p class="page-subtitle">探索社区最新动态与知识分享</p>
        </div>
        <div class="header-actions" v-if="isLoggedIn">
          <router-link to="/ai-chat" class="btn-glass">
            <span class="btn-icon">✨</span> AI 助手
          </router-link>
          <router-link to="/post/create" class="btn-accent">
            <span class="btn-icon">✏️</span> 发布帖子
          </router-link>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- Error -->
      <div v-if="error" class="error-state glass-card">
        <span class="error-icon">⚠️</span>
        {{ error }}
        <button class="btn-glass-sm" @click="loadPosts">重试</button>
      </div>

      <!-- Post List -->
      <div v-if="!loading && !error" class="post-list">
        <div v-if="posts.length === 0" class="empty-state glass-card">
          <div class="empty-icon"> </div>
          <p>暂无帖子</p>
          <router-link to="/post/create" class="btn-accent-sm">来发第一篇吧</router-link>
        </div>

        <article
          v-for="(post, index) in posts"
          :key="post.id"
          class="post-card glass-card"
          :style="{ animationDelay: `${index * 60}ms` }"
          @click="goToPost(post.id)"
        >
          <span v-if="post.is_pinned" class="badge-pinned">📌 置顶</span>

          <div class="post-card-body">
            <h2 class="post-title">{{ post.title }}</h2>
            <p class="post-summary">{{ post.summary }}</p>
            <div class="post-tags" v-if="post.tags && post.tags.length">
              <span v-for="tag in post.tags" :key="tag.id" class="tag"># {{ tag.name }}</span>
            </div>
          </div>

          <footer class="post-card-footer">
            <div class="post-author">
              <img :src="post.author.avatar || defaultAvatar" :alt="post.author.username" class="avatar" />
              <span>{{ post.author.username }}</span>
            </div>
            <span v-if="post.category && post.category.name" class="category-badge">
              {{ post.category.icon }} {{ post.category.name }}
            </span>
            <div class="post-stats">
              <span>👁 {{ post.view_count }}</span>
              <span>👍 {{ post.like_count }}</span>
              <span>💬 {{ post.comment_count }}</span>
            </div>
            <span class="post-time">{{ formatTime(post.created_at) }}</span>
          </footer>
        </article>
      </div>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="pagination">
        <button class="btn-glass-sm" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <span class="page-info">第 {{ page }} 页 / 共 {{ totalPages }} 页（{{ total }} 篇）</span>
        <button class="btn-glass-sm" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
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
      posts: [],
      total: 0,
      page: 1,
      pageSize: 5,
      loading: false,
      error: '',
      defaultAvatar: 'https://api.dicebear.com/7.x/thumbs/svg?seed=default',
      authToken: '',
      currentUsername: ''
    }
  },
  computed: {
    isLoggedIn() { return !!this.authToken },
    username() { return this.currentUsername || '我' },
    totalPages() { return Math.ceil(this.total / this.pageSize) }
  },
  mounted() {
    this.authToken = localStorage.getItem('auth_token') || ''
    this.currentUsername = localStorage.getItem('username') || ''
    this._onStorageChange = (e) => this.onStorageChange(e)
    this._onLoginSuccess = (e) => this.onLoginSuccess(e)
    window.addEventListener('storage', this._onStorageChange)
    window.addEventListener('login-success', this._onLoginSuccess)
    this.loadPosts()
  },
  beforeUnmount() {
    window.removeEventListener('storage', this._onStorageChange)
    window.removeEventListener('login-success', this._onLoginSuccess)
  },
  methods: {
    async loadPosts() {
      this.loading = true
      this.error = ''
      try {
        const res = await getPostList({ page: this.page, page_size: this.pageSize })
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
    changePage(newPage) {
      this.page = newPage
      this.loadPosts()
      window.scrollTo({ top: 0, behavior: 'smooth' })
    },
    goToPost(id) { this.$router.push(`/posts/${id}`) },
    handleLogout() {
      logout()
      this.authToken = ''
      this.currentUsername = ''
      this.$router.go(0)
    },
    onStorageChange(e) {
      if (e.key === 'auth_token') this.authToken = e.newValue || ''
      if (e.key === 'username') this.currentUsername = e.newValue || ''
    },
    onLoginSuccess(e) {
      const { token, user } = e.detail
      this.authToken = token || ''
      this.currentUsername = user?.username || ''
    },
    formatTime(isoString) {
      if (!isoString) return ''
      const date = new Date(isoString)
      const now = new Date()
      const diff = Math.floor((now - date) / 1000)
      if (diff < 60) return '刚刚'
      if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`
      if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`
      if (diff < 86400 * 7) return `${Math.floor(diff / 86400)} 天前`
      return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
    }
  }
}
</script>

<style scoped>
/* ═══════════════════════════════════════════════════════════════════════
   Navbar
   ═══════════════════════════════════════════════════════════════════════ */

.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(12, 10, 29, 0.6);
  backdrop-filter: blur(24px) saturate(1.5);
  -webkit-backdrop-filter: blur(24px) saturate(1.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.navbar-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-lg);
  height: 64px;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: var(--text-primary);
  font-weight: 700;
  font-size: 20px;
  letter-spacing: -0.3px;
}

.brand-icon {
  font-size: 24px;
}

.brand-text {
  background: var(--accent-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.navbar-nav {
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav-link {
  text-decoration: none;
  color: var(--text-secondary);
  padding: 8px 14px;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 500;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.nav-link:hover,
.nav-link.active {
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.08);
}

/* ═══════════════════════════════════════════════════════════════════════
   Buttons
   ═══════════════════════════════════════════════════════════════════════ */

.btn-accent {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  color: var(--text-on-accent);
  background: var(--accent-gradient);
  padding: 9px 18px;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all var(--duration-normal) var(--ease-spring);
  box-shadow: 0 4px 16px var(--accent-glow);
}

.btn-accent:hover {
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 8px 24px var(--accent-glow);
}

.btn-accent-sm {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  color: var(--text-on-accent);
  background: var(--accent-gradient);
  padding: 7px 14px;
  border-radius: var(--glass-radius-xs);
  font-size: 13px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all var(--duration-normal) var(--ease-spring);
}

.btn-glass {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  color: var(--text-primary);
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  padding: 9px 18px;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-glass:hover {
  background: var(--glass-bg-hover);
  border-color: var(--glass-border-bright);
  transform: translateY(-1px);
}

.btn-glass-sm {
  display: inline-flex;
  align-items: center;
  color: var(--text-primary);
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  padding: 7px 14px;
  border-radius: var(--glass-radius-xs);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--duration-normal) var(--ease-out-expo);
  font-family: inherit;
}

.btn-glass-sm:hover:not(:disabled) {
  background: var(--glass-bg-hover);
  border-color: var(--glass-border-bright);
  transform: translateY(-1px);
}

.btn-glass-sm:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.btn-glass-outline {
  text-decoration: none;
  color: var(--text-secondary);
  border: 1px solid var(--glass-border);
  padding: 7px 16px;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 500;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-glass-outline:hover {
  color: var(--text-primary);
  border-color: var(--glass-border-bright);
  background: rgba(255, 255, 255, 0.05);
}

.btn-ghost {
  background: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--text-secondary);
  padding: 7px 14px;
  border-radius: var(--glass-radius-xs);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-ghost:hover {
  border-color: var(--danger);
  color: var(--danger);
  background: var(--danger-glow);
}

.btn-icon {
  font-size: 15px;
}

/* ═══════════════════════════════════════════════════════════════════════
   Main Content
   ═══════════════════════════════════════════════════════════════════════ */

.main-content {
  max-width: 900px;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
}

.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: var(--space-xl);
}

.page-title {
  font-size: 32px;
  font-weight: 800;
  color: var(--text-primary);
  letter-spacing: -0.5px;
  line-height: 1.2;
}

.page-subtitle {
  color: var(--text-tertiary);
  font-size: 15px;
  margin-top: 6px;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-shrink: 0;
}

/* ═══════════════════════════════════════════════════════════════════════
   Glass Card Base
   ═══════════════════════════════════════════════════════════════════════ */

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid var(--glass-border);
  border-radius: var(--glass-radius);
  box-shadow: var(--glass-shadow);
}

/* ═══════════════════════════════════════════════════════════════════════
   States
   ═══════════════════════════════════════════════════════════════════════ */

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-tertiary);
  padding: 80px 0;
  font-size: 15px;
}

.spinner {
  width: 22px;
  height: 22px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: var(--space-lg);
  color: var(--danger);
  font-size: 14px;
}

.error-icon {
  font-size: 18px;
}

.empty-state {
  text-align: center;
  padding: var(--space-2xl);
  color: var(--text-tertiary);
}

.empty-icon {
  font-size: 56px;
  margin-bottom: var(--space-md);
  opacity: 0.6;
}

.empty-state p {
  margin-bottom: var(--space-md);
  font-size: 16px;
}

/* ═══════════════════════════════════════════════════════════════════════
   Post Cards
   ═══════════════════════════════════════════════════════════════════════ */

.post-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.post-card {
  padding: var(--space-lg);
  cursor: pointer;
  position: relative;
  transition: all var(--duration-slow) var(--ease-out-expo);
  animation: fadeInUp var(--duration-slow) var(--ease-out-expo) both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.post-card:hover {
  background: var(--glass-bg-hover);
  border-color: var(--glass-border-bright);
  transform: translateY(-3px);
  box-shadow: var(--glass-shadow-hover), 0 0 0 1px rgba(129, 140, 248, 0.1);
}

.badge-pinned {
  position: absolute;
  top: var(--space-md);
  right: var(--space-md);
  font-size: 12px;
  color: var(--warning);
  background: rgba(251, 191, 36, 0.12);
  padding: 4px 10px;
  border-radius: 8px;
  border: 1px solid rgba(251, 191, 36, 0.2);
  font-weight: 500;
}

.post-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px;
  line-height: 1.4;
  letter-spacing: -0.2px;
}

.post-summary {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.7;
  margin: 0 0 14px;
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
  margin-bottom: 14px;
}

.tag {
  font-size: 12px;
  color: var(--accent-primary);
  background: rgba(129, 140, 248, 0.1);
  border: 1px solid rgba(129, 140, 248, 0.15);
  padding: 3px 10px;
  border-radius: 20px;
  font-weight: 500;
}

.post-card-footer {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  font-size: 13px;
  color: var(--text-tertiary);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  padding-top: 12px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-weight: 500;
}

.avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.category-badge {
  background: rgba(255, 255, 255, 0.06);
  padding: 3px 10px;
  border-radius: 8px;
  font-size: 12px;
  color: var(--text-secondary);
}

.post-stats {
  display: flex;
  gap: 12px;
  margin-left: auto;
}

.post-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* ═══════════════════════════════════════════════════════════════════════
   Pagination
   ═══════════════════════════════════════════════════════════════════════ */

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-md);
  margin-top: var(--space-xl);
  padding-bottom: var(--space-2xl);
}

.page-info {
  color: var(--text-tertiary);
  font-size: 14px;
}

/* ═══════════════════════════════════════════════════════════════════════
   Responsive
   ═══════════════════════════════════════════════════════════════════════ */

@media (max-width: 768px) {
  .navbar-inner { padding: 0 var(--space-md); }
  .main-content { padding: var(--space-lg) var(--space-md); }
  .page-header { flex-direction: column; align-items: flex-start; gap: var(--space-md); }
  .page-title { font-size: 26px; }
  .post-card { padding: var(--space-md); }
  .post-card-footer { gap: 10px; font-size: 12px; }
  .post-stats { margin-left: 0; }
  .header-actions { width: 100%; }
}

@media (max-width: 480px) {
  .navbar-nav { gap: 4px; }
  .nav-link { padding: 6px 10px; font-size: 13px; }
  .brand-text { font-size: 17px; }
  .post-title { font-size: 16px; }
}
</style>
