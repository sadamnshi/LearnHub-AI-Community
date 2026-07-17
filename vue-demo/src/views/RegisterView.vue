<template>
  <section class="auth-page">
    <div class="auth-card glass-card">
      <div class="auth-header">
        <div class="auth-logo">✨</div>
        <h2>创建账号</h2>
        <p class="auth-subtitle">加入 LearnHub 社区</p>
      </div>

      <p v-if="error" class="alert alert-error">{{ error }}</p>
      <p v-if="message" class="alert alert-success">{{ message }}</p>

      <form class="auth-form" @submit.prevent="handleRegister">
        <div class="form-field">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model.trim="form.username"
            type="text"
            required
            placeholder="请输入用户名"
          >
        </div>
        <div class="form-field">
          <label for="password">密码</label>
          <input
            id="password"
            v-model.trim="form.password"
            type="password"
            required
            placeholder="请输入密码"
          >
        </div>
        <button type="submit" class="btn-submit" :disabled="loading">
          <span v-if="loading" class="btn-loading"></span>
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>

      <router-link class="auth-link" to="/login">已有账号？去登录</router-link>
    </div>
  </section>
</template>

<script>
import { register } from '@/api/auth'

export default {
  name: 'RegisterView',
  data() {
    return {
      form: { username: '', password: '' },
      loading: false,
      error: '',
      message: ''
    }
  },
  methods: {
    async handleRegister() {
      this.error = ''
      this.message = ''
      this.loading = true
      try {
        const data = await register(this.form)
        this.message = data.message || '注册成功，即将跳转到登录页...'
        setTimeout(() => { this.$router.push('/login') }, 2000)
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
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

.auth-card {
  width: min(420px, 100%);
  padding: var(--space-2xl) var(--space-xl);
  animation: fadeInUp 0.5s var(--ease-out-expo);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.auth-header {
  text-align: center;
  margin-bottom: var(--space-xl);
}

.auth-logo {
  font-size: 48px;
  margin-bottom: var(--space-md);
}

.auth-header h2 {
  font-size: 28px;
  font-weight: 800;
  color: var(--text-primary);
  margin-bottom: 6px;
  letter-spacing: -0.5px;
}

.auth-subtitle {
  color: var(--text-tertiary);
  font-size: 14px;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.form-field label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-field input {
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

.form-field input::placeholder {
  color: var(--text-tertiary);
}

.form-field input:focus {
  border-color: var(--accent-primary);
  background: rgba(129, 140, 248, 0.06);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.btn-submit {
  margin-top: var(--space-sm);
  padding: 13px;
  background: linear-gradient(135deg, #34d399 0%, #10b981 100%);
  color: #fff;
  border: none;
  border-radius: var(--glass-radius-sm);
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-spring);
  box-shadow: 0 4px 16px var(--success-glow);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px) scale(1.01);
  box-shadow: 0 8px 24px var(--success-glow);
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
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

.alert {
  padding: 12px 16px;
  border-radius: var(--glass-radius-xs);
  font-size: 13px;
  margin-bottom: var(--space-md);
  font-weight: 500;
}

.alert-error {
  background: var(--danger-glow);
  color: var(--danger);
  border: 1px solid rgba(248, 113, 113, 0.2);
}

.alert-success {
  background: var(--success-glow);
  color: var(--success);
  border: 1px solid rgba(52, 211, 153, 0.2);
}

.auth-link {
  display: block;
  text-align: center;
  margin-top: var(--space-lg);
  color: var(--accent-primary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color var(--duration-fast);
}

.auth-link:hover {
  color: var(--accent-secondary);
}
</style>
