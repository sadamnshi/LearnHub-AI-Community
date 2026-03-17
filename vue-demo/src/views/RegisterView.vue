<template>
  <section class="page">
    <div class="card">
      <h2>用户注册</h2>
      <p v-if="error" class="alert error">错误：{{ error }}</p>
      <p v-if="message" class="alert success">{{ message }}</p>

      <form class="form" @submit.prevent="handleRegister">
        <label for="username">用户名</label>
        <input
          id="username"
          v-model.trim="form.username"
          type="text"
          name="username"
          required
        >

        <label for="password">密码</label>
        <input
          id="password"
          v-model.trim="form.password"
          type="password"
          name="password"
          required
        >

        <button type="submit" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>

      <router-link class="link" to="/login">已有账号？去登录</router-link>
    </div>
  </section>
</template>

<script>
import { register } from '@/api/auth'

export default {
  name: 'RegisterView',
  data() {
    return {
      form: {
        username: '',
        password: ''
      },
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
        // 使用封装的 API 方法
        const data = await register(this.form)
        this.message = data.message || '注册成功，即将跳转到登录页...'
        
        // 2秒后自动跳转到登录页
        setTimeout(() => {
          this.$router.push('/login')
        }, 2000)
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
.page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fb;
  padding: 32px 16px;
}

.card {
  width: min(420px, 100%);
  background: #fff;
  padding: 32px;
  border-radius: 16px;
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.12);
  text-align: left;
}

h2 {
  margin-bottom: 20px;
  color: #1e293b;
}

.form {
  display: grid;
  gap: 12px;
}

label {
  font-size: 14px;
  color: #475569;
}

input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d0d7e2;
  border-radius: 10px;
  font-size: 14px;
}

button {
  margin-top: 8px;
  padding: 12px;
  background: #16a34a;
  color: #fff;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
}

button:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

.alert {
  margin-bottom: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 14px;
}

.error {
  background: #fee2e2;
  color: #b91c1c;
}

.success {
  background: #dcfce7;
  color: #15803d;
}

.link {
  display: inline-block;
  margin-top: 16px;
  color: #2563eb;
  text-decoration: none;
}
</style>
