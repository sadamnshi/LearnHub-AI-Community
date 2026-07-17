<template>
  <div class="create-post-page">
    <div class="form-wrapper glass-card">
      <h1 class="form-title">发布新帖子</h1>

      <form @submit.prevent="handleSubmit">
        <!-- Title -->
        <div class="form-group">
          <label for="title">帖子标题 *</label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            placeholder="请输入帖子标题"
            maxlength="200"
            required
          />
          <span class="char-count">{{ form.title.length }} / 200</span>
        </div>

        <!-- Content -->
        <div class="form-group">
          <label for="content">帖子内容 *</label>
          <div class="editor-wrapper">
            <textarea
              id="content"
              v-model="form.content"
              placeholder="请输入帖子内容（支持 Markdown，至少 10 字符）"
              rows="10"
              minlength="10"
              required
            ></textarea>
            <span class="char-count">{{ form.content.length }} 字符</span>
          </div>
        </div>

        <!-- Category -->
        <div class="form-group">
          <label for="category">帖子分类</label>
          <select id="category" v-model.number="form.category_id">
            <option value="0">不选择分类</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.icon }} {{ category.name }}
            </option>
          </select>
        </div>

        <!-- Tags -->
        <div class="form-group">
          <label for="tags">帖子标签</label>
          <input
            id="tags"
            v-model="form.tags"
            type="text"
            placeholder="用逗号分隔，如：golang, 数据库, 缓存"
            maxlength="200"
          />
          <div class="tags-preview" v-if="form.tags">
            <span v-for="tag in parsedTags" :key="tag" class="tag-chip">{{ tag }}</span>
          </div>
        </div>

        <!-- Error -->
        <div v-if="error" class="alert-error">
          ⚠️ {{ error }}
        </div>

        <!-- Actions -->
        <div class="button-group">
          <button type="submit" class="btn-accent" :disabled="loading">
            <span v-if="loading" class="btn-loading"></span>
            {{ loading ? '发布中...' : '发布帖子' }}
          </button>
          <router-link to="/" class="btn-glass">取消</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { createPost, getCategories } from '../api/post'

export default {
  name: 'CreatePostView',
  data() {
    return {
      form: { title: '', content: '', category_id: 0, tags: '' },
      categories: [],
      loading: false,
      error: ''
    }
  },
  computed: {
    parsedTags() {
      return this.form.tags.split(',').map(t => t.trim()).filter(t => t.length > 0)
    }
  },
  mounted() {
    this.loadCategories()
  },
  methods: {
    async loadCategories() {
      try {
        this.categories = await getCategories()
      } catch {
        this.categories = []
      }
    },
    async handleSubmit() {
      if (!this.form.title.trim()) { this.error = '请输入帖子标题'; return }
      if (this.form.title.length > 200) { this.error = '标题不能超过200字符'; return }
      if (!this.form.content.trim()) { this.error = '请输入帖子内容'; return }
      if (this.form.content.length < 10) { this.error = '内容至少需要10个字符'; return }
      if (this.form.tags.length > 200) { this.error = '标签不能超过200字符'; return }

      this.error = ''
      this.loading = true
      try {
        const response = await createPost({
          title: this.form.title.trim(),
          content: this.form.content.trim(),
          category_id: this.form.category_id,
          tags: this.form.tags.trim()
        })
        this.$router.push({ name: 'post-detail', params: { id: response.id } })
      } catch (err) {
        this.error = err.message || '发布帖子失败，请重试'
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.create-post-page {
  min-height: 100vh;
  padding: var(--space-xl) var(--space-md);
}

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid var(--glass-border);
  border-radius: var(--glass-radius);
  box-shadow: var(--glass-shadow);
}

.form-wrapper {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-2xl) var(--space-xl);
  animation: fadeInUp 0.5s var(--ease-out-expo);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.form-title {
  margin: 0 0 var(--space-xl) 0;
  font-size: 30px;
  font-weight: 800;
  color: var(--text-primary);
  text-align: center;
  letter-spacing: -0.5px;
}

.form-group {
  margin-bottom: var(--space-lg);
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--glass-radius-sm);
  font-size: 15px;
  color: var(--text-primary);
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
  outline: none;
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: var(--text-tertiary);
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  border-color: var(--accent-primary);
  background: rgba(129, 140, 248, 0.06);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.form-group select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath d='M6 8L1 3h10z' fill='rgba(241,245,249,0.4)'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  padding-right: 36px;
}

.form-group select option {
  background: #1e1b3a;
  color: var(--text-primary);
}

.form-group textarea {
  resize: vertical;
  min-height: 200px;
  font-family: 'JetBrains Mono', 'Fira Code', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.7;
}

.editor-wrapper {
  position: relative;
}

.char-count {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: var(--text-tertiary);
  text-align: right;
}

.tags-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 10px;
}

.tag-chip {
  display: inline-block;
  padding: 4px 12px;
  background: rgba(129, 140, 248, 0.1);
  border: 1px solid rgba(129, 140, 248, 0.2);
  border-radius: 20px;
  font-size: 12px;
  color: var(--accent-primary);
  font-weight: 500;
}

.alert-error {
  padding: 14px 18px;
  margin-bottom: var(--space-lg);
  background: var(--danger-glow);
  border: 1px solid rgba(248, 113, 113, 0.2);
  border-left: 3px solid var(--danger);
  border-radius: var(--glass-radius-xs);
  color: var(--danger);
  font-size: 14px;
}

.button-group {
  display: flex;
  gap: var(--space-md);
  justify-content: center;
  margin-top: var(--space-xl);
}

.btn-accent {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 32px;
  background: var(--accent-gradient);
  color: var(--text-on-accent);
  border: none;
  border-radius: var(--glass-radius-sm);
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-spring);
  box-shadow: 0 4px 16px var(--accent-glow);
  min-width: 140px;
}

.btn-accent:hover:not(:disabled) {
  transform: translateY(-2px) scale(1.01);
  box-shadow: 0 8px 24px var(--accent-glow);
}

.btn-accent:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-loading {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.btn-glass {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 28px;
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  color: var(--text-primary);
  border-radius: var(--glass-radius-sm);
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-glass:hover {
  background: var(--glass-bg-hover);
  border-color: var(--glass-border-bright);
  transform: translateY(-1px);
}

/* Responsive */
@media (max-width: 640px) {
  .form-wrapper { padding: var(--space-lg) var(--space-md); }
  .form-title { font-size: 24px; }
  .button-group { flex-direction: column; }
  .btn-accent, .btn-glass { width: 100%; }
}
</style>
