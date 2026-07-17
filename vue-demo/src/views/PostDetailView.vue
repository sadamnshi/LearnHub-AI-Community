<template>
  <div class="post-detail-page">
    <!-- Loading -->
    <div v-if="loading" class="loading-container">
      <div class="spinner-lg"></div>
      <p>加载帖子详情中...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="error-container glass-card">
      <div class="error-emoji">⚠️</div>
      <h2>出错了</h2>
      <p>{{ error }}</p>
      <button class="btn-glass" @click="goBack">← 返回列表</button>
    </div>

    <!-- Post -->
    <article v-else class="post-detail glass-card">
      <!-- Toolbar -->
      <div class="toolbar">
        <button class="toolbar-btn" @click="goBack" title="返回列表">← 返回</button>
        <div class="toolbar-spacer"></div>
        <button class="toolbar-btn" title="分享">↗ 分享</button>
        <button class="toolbar-btn like-btn" @click="toggleLike" :class="{ liked: isLiked }" title="点赞">
          {{ isLiked ? '♥' : '♡' }} {{ likeCount }}
        </button>
      </div>

      <!-- Header -->
      <header class="post-header">
        <h1 class="post-title">{{ post.title }}</h1>
        <div class="post-meta">
          <div class="meta-item author-info">
            <img :src="post.author.avatar" :alt="post.author.username" class="avatar" />
            <span class="author-name">{{ post.author.username }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">📅</span>
            <time :datetime="post.created_at">{{ formatDate(post.created_at) }}</time>
          </div>
          <div v-if="isUpdated" class="meta-item">
            <span class="meta-label">✏️</span>
            <time :datetime="post.updated_at">{{ formatDate(post.updated_at) }}</time>
          </div>
          <div class="meta-item">
            <span class="meta-label">👁</span>
            <span>{{ viewCount }} 浏览</span>
          </div>
        </div>
      </header>

      <!-- Category & Tags -->
      <div class="categories-tags">
        <div class="category-chip">
          <span>{{ post.category.icon }}</span>
          <span>{{ post.category.name }}</span>
        </div>
        <div v-if="post.tags && post.tags.length > 0" class="tags">
          <span v-for="tag in post.tags" :key="tag" class="tag-chip">#{{ tag }}</span>
        </div>
      </div>

      <hr class="divider" />

      <!-- Content -->
      <section class="post-content">
        <div class="content-body" v-html="renderedContent"></div>
      </section>

      <hr class="divider" />

      <!-- Footer Actions -->
      <footer class="post-footer">
        <div class="footer-actions">
          <button class="action-btn" @click="toggleLike" :class="{ liked: isLiked }">
            {{ isLiked ? '♥' : '♡' }}
            {{ isLiked ? '已点赞' : '点赞' }}
            <span class="action-count">({{ likeCount }})</span>
          </button>
          <button class="action-btn">
            ↗ 分享
          </button>
          <button class="action-btn">
            ✎ 评论
            <span class="action-count">(0)</span>
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

const route = useRoute()
const router = useRouter()

const post = ref(null)
const loading = ref(true)
const error = ref(null)
const isLiked = ref(false)
const likeCount = ref(0)
const viewCount = ref(0)

const isUpdated = computed(() => {
  if (!post.value) return false
  return post.value.created_at !== post.value.updated_at
})

const renderedContent = computed(() => {
  if (!post.value) return ''
  return post.value.content.replace(/\n/g, '<br>')
})

const goBack = () => { router.back() }

const toggleLike = async () => {
  isLiked.value = !isLiked.value
  likeCount.value += isLiked.value ? 1 : -1
}

const formatDate = (dateStr) => dateStr

const loadPostDetail = async () => {
  try {
    error.value = null
    loading.value = true
    const data = await fetchPostDetail(route.params.id)
    post.value = data
    likeCount.value = data.like_count || 0
    viewCount.value = data.view_count || 0
  } catch (err) {
    error.value = err.message || '获取帖子详情失败，请重试'
  } finally {
    loading.value = false
  }
}

onMounted(() => { loadPostDetail() })
</script>

<style scoped>
.post-detail-page {
  min-height: 100vh;
  padding: var(--space-lg);
}

/* Glass card base */
.glass-card {
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid var(--glass-border);
  border-radius: var(--glass-radius);
  box-shadow: var(--glass-shadow);
}

/* Loading */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: var(--space-lg);
  color: var(--text-tertiary);
  font-size: 15px;
}

.spinner-lg {
  width: 44px;
  height: 44px;
  border: 3px solid rgba(255, 255, 255, 0.08);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* Error */
.error-container {
  max-width: 500px;
  margin: 80px auto;
  padding: var(--space-2xl);
  text-align: center;
  animation: fadeInUp 0.5s var(--ease-out-expo);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.error-emoji {
  font-size: 56px;
  margin-bottom: var(--space-md);
}

.error-container h2 {
  color: var(--danger);
  margin-bottom: var(--space-sm);
  font-size: 22px;
}

.error-container p {
  color: var(--text-secondary);
  margin-bottom: var(--space-xl);
  line-height: 1.6;
}

/* Post Detail */
.post-detail {
  max-width: 860px;
  margin: 0 auto;
  overflow: hidden;
  animation: fadeInUp 0.5s var(--ease-out-expo);
}

/* Toolbar */
.toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md) var(--space-xl);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.toolbar-spacer { flex: 1; }

.toolbar-btn {
  padding: 8px 14px;
  background: none;
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--text-secondary);
  border-radius: var(--glass-radius-xs);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.toolbar-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-primary);
  border-color: var(--glass-border);
}

.like-btn.liked {
  color: var(--danger);
  border-color: rgba(248, 113, 113, 0.25);
  background: var(--danger-glow);
}

/* Header */
.post-header {
  padding: var(--space-2xl) var(--space-xl) var(--space-lg);
}

.post-title {
  font-size: 34px;
  font-weight: 800;
  color: var(--text-primary);
  margin: 0 0 var(--space-lg) 0;
  line-height: 1.3;
  letter-spacing: -0.5px;
}

.post-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--space-lg);
  color: var(--text-secondary);
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.author-info {
  gap: 10px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.author-name {
  color: var(--accent-primary);
  font-weight: 600;
}

.meta-label {
  font-size: 15px;
}

/* Category & Tags */
.categories-tags {
  padding: 0 var(--space-xl);
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--space-sm);
}

.category-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: rgba(129, 140, 248, 0.1);
  border: 1px solid rgba(129, 140, 248, 0.2);
  border-radius: 20px;
  color: var(--accent-primary);
  font-size: 13px;
  font-weight: 600;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag-chip {
  padding: 5px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.tag-chip:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--glass-border);
  color: var(--text-primary);
}

/* Divider */
.divider {
  margin: var(--space-lg) var(--space-xl);
  border: none;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

/* Content */
.post-content {
  padding: var(--space-lg) var(--space-xl);
  font-size: 16px;
  line-height: 1.9;
  color: var(--text-secondary);
}

.content-body {
  word-break: break-word;
}

.content-body :deep(p) { margin: 16px 0; }
.content-body :deep(h1),
.content-body :deep(h2),
.content-body :deep(h3) {
  margin: 28px 0 14px 0;
  color: var(--text-primary);
  font-weight: 700;
}

.content-body :deep(h1) { font-size: 26px; }
.content-body :deep(h2) { font-size: 22px; }
.content-body :deep(h3) { font-size: 18px; }

.content-body :deep(code) {
  background: rgba(255, 255, 255, 0.06);
  padding: 2px 8px;
  border-radius: 6px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 0.9em;
  color: var(--accent-primary);
}

.content-body :deep(pre) {
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  padding: var(--space-md);
  border-radius: var(--glass-radius-sm);
  overflow-x: auto;
  margin: 16px 0;
}

.content-body :deep(pre code) {
  background: none;
  padding: 0;
  color: var(--text-secondary);
}

.content-body :deep(blockquote) {
  border-left: 3px solid var(--accent-primary);
  margin: 16px 0;
  padding: 12px 20px;
  background: rgba(129, 140, 248, 0.05);
  border-radius: 0 var(--glass-radius-xs) var(--glass-radius-xs) 0;
  color: var(--text-secondary);
  font-style: italic;
}

.content-body :deep(ul),
.content-body :deep(ol) {
  margin: 16px 0;
  padding-left: 28px;
}

.content-body :deep(li) { margin: 8px 0; }

/* Footer */
.post-footer {
  padding: var(--space-xl);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.footer-actions {
  display: flex;
  gap: var(--space-sm);
}

.action-btn {
  flex: 1;
  padding: 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--glass-radius-sm);
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--glass-border);
  color: var(--text-primary);
  transform: translateY(-1px);
}

.action-btn.liked {
  background: var(--danger-glow);
  border-color: rgba(248, 113, 113, 0.2);
  color: var(--danger);
}

.action-count {
  color: var(--text-tertiary);
  font-weight: 400;
}

.btn-glass {
  display: inline-flex;
  align-items: center;
  color: var(--text-primary);
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  padding: 10px 20px;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-glass:hover {
  background: var(--glass-bg-hover);
  transform: translateY(-1px);
}

/* Responsive */
@media (max-width: 768px) {
  .post-detail-page { padding: var(--space-md); }
  .post-title { font-size: 24px; }
  .post-header, .post-content, .post-footer, .categories-tags { padding: var(--space-lg) var(--space-md); }
  .toolbar, .divider { padding-left: var(--space-md); padding-right: var(--space-md); }
  .divider { margin-left: var(--space-md); margin-right: var(--space-md); }
  .post-meta { gap: var(--space-md); font-size: 13px; }
  .footer-actions { flex-direction: column; }
}
</style>
