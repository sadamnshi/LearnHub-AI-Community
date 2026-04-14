<template>
  <div class="create-post-container">
    <div class="form-wrapper">
      <h1>发布新帖子</h1>
      
      <form @submit.prevent="handleSubmit">
        <!-- 标题输入框 -->
        <div class="form-group">
          <label for="title">帖子标题 *</label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            placeholder="请输入帖子标题（最多200字符）"
            maxlength="200"
            required
          />
          <span class="char-count">{{ form.title.length }}/200</span>
        </div>

        <!-- 内容输入框 -->
        <div class="form-group">
          <label for="content">帖子内容 *</label>
          <div class="editor-wrapper">
            <textarea
              id="content"
              v-model="form.content"
              placeholder="请输入帖子内容（支持Markdown，至少10字符）"
              rows="10"
              minlength="10"
              required
            ></textarea>
            <span class="char-count">{{ form.content.length }} 字符</span>
          </div>
        </div>

        <!-- 分类选择 -->
        <div class="form-group">
          <label for="category">帖子分类</label>
          <select
            id="category"
            v-model.number="form.category_id"
          >
            <option value="0">不选择分类</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.icon }} {{ category.name }}
            </option>
          </select>
        </div>

        <!-- 标签输入框 -->
        <div class="form-group">
          <label for="tags">帖子标签</label>
          <input
            id="tags"
            v-model="form.tags"
            type="text"
            placeholder="用逗号分隔多个标签，如：golang,数据库,缓存（最多200字符）"
            maxlength="200"
          />
          <div class="tags-preview" v-if="form.tags">
            <span v-for="tag in parsedTags" :key="tag" class="tag">
              {{ tag }}
            </span>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="error-message">
          ⚠️ {{ error }}
        </div>

        <!-- 提交按钮 -->
        <div class="button-group">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading">发布中...</span>
            <span v-else>发布帖子</span>
          </button>
          <router-link to="/" class="btn btn-secondary">取消</router-link>
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
      form: {
        title: '',
        content: '',
        category_id: 0,
        tags: ''
      },
      categories: [],
      loading: false,
      error: '',
      categoriesLoading: true
    }
  },
  computed: {
    // 解析标签数组用于预览
    parsedTags() {
      return this.form.tags
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag.length > 0)
    }
  },
  mounted() {
    // 组件挂载时加载分类列表
    this.loadCategories()
  },
  methods: {
    // 加载分类列表
    async loadCategories() {
      try {
        this.categoriesLoading = true
        this.categories = await getCategories()
        console.log('✅ 分类加载成功:', this.categories)
      } catch (err) {
        console.error('❌ 分类加载失败:', err)
        this.error = '加载分类失败，请刷新页面重试'
        // 即使加载失败，也提供一个空分类列表，让用户可以继续发布
        this.categories = []
      } finally {
        this.categoriesLoading = false
      }
    },
    async handleSubmit() {
      // 验证标题
      if (!this.form.title.trim()) {
        this.error = '请输入帖子标题'
        return
      }

      if (this.form.title.length > 200) {
        this.error = '标题不能超过200字符'
        return
      }

      // 验证内容
      if (!this.form.content.trim()) {
        this.error = '请输入帖子内容'
        return
      }

      if (this.form.content.length < 10) {
        this.error = '内容至少需要10个字符'
        return
      }

      // 验证标签
      if (this.form.tags.length > 200) {
        this.error = '标签不能超过200字符'
        return
      }

      this.error = ''
      this.loading = true

      try {
        // 调用后端 API 创建帖子
        const response = await createPost({
          title: this.form.title.trim(),
          content: this.form.content.trim(),
          category_id: this.form.category_id,
          tags: this.form.tags.trim()
        })

        // 创建成功，跳转到帖子详情页
        this.$router.push({
          name: 'post-detail',
          params: { id: response.id }
        })

        // 可选：显示成功提示
        alert('帖子发布成功！')
      } catch (err) {
        console.error('发布帖子失败:', err)
        this.error = err.message || '发布帖子失败，请重试'
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.create-post-container {
  min-height: 100vh;
  background: #f5f7fb;
  padding: 40px 20px;
}

.form-wrapper {
  max-width: 800px;
  margin: 0 auto;
  background: #ffffff;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

h1 {
  margin: 0 0 32px 0;
  font-size: 28px;
  color: #0f172a;
  text-align: center;
}

.form-group {
  margin-bottom: 24px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #334155;
}

input[type="text"],
input[type="email"],
textarea,
select {
  width: 100%;
  padding: 10px 12px;
  font-size: 14px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-family: inherit;
  transition: border-color 0.3s ease;
}

input[type="text"]:focus,
input[type="email"]:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

textarea {
  resize: vertical;
  min-height: 200px;
  font-family: 'Courier New', monospace;
}

.editor-wrapper {
  position: relative;
}

.char-count {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
  text-align: right;
}

.tags-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.tag {
  display: inline-block;
  padding: 4px 12px;
  background: #e0f2fe;
  color: #0369a1;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.error-message {
  padding: 12px 16px;
  margin-bottom: 24px;
  background: #fee2e2;
  color: #991b1b;
  border-radius: 6px;
  border-left: 4px solid #dc2626;
  font-size: 14px;
}

.button-group {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 32px;
}

.btn {
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  display: inline-block;
  text-align: center;
}

.btn-primary {
  background: #3b82f6;
  color: #ffffff;
  min-width: 120px;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.btn-primary:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  opacity: 0.6;
}

.btn-secondary {
  background: #e2e8f0;
  color: #334155;
}

.btn-secondary:hover {
  background: #cbd5e1;
  transform: translateY(-2px);
}

/* 响应式设计 */
@media (max-width: 640px) {
  .form-wrapper {
    padding: 24px;
  }

  h1 {
    font-size: 22px;
    margin-bottom: 24px;
  }

  .form-group {
    margin-bottom: 18px;
  }

  .button-group {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>
