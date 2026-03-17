<template>
  <section class="page">
    <div class="card">
      <h2>用户登录</h2>
      <p v-if="error" class="alert error">错误：{{ error }}</p>
      <p v-if="message" class="alert success">{{ message }}</p>
      <p v-if="username" class="welcome">欢迎：{{ username }}</p>

      <form class="form" @submit.prevent="handleLogin">
        <label for="username">账号</label>
        <input
          id="username"
          v-model.trim="form.username"
          type="text"
          name="username"
          autocomplete="username"
          required
        >

        <label for="password">密码</label>
        <input
          id="password"
          v-model.trim="form.password"
          type="password"
          name="password"
          autocomplete="current-password"
          required
        >

        <button type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <router-link class="link" to="/register">没有账号？去注册</router-link>
    </div>
  </section>
</template>

<script>
import { login } from '@/api/auth'

export default {
  name: 'LoginView',
  data() {
    return {
      form: {
        username: '',
        password: ''
      },
      loading: false,
      error: '',
      message: '',
      username: ''
    }
  },
  methods: {
    async handleLogin() {
      this.error = ''
      this.message = ''
      this.username = ''
      this.loading = true

      try {
        // 使用封装的 API 方法
        const response = await login(this.form)
        
        console.log('Login response:', response)
        
        // 后端返回格式：{ code: 0, msg: "...", data: { token, user } }
        const { token, user } = response.data || {}
        
        // 保存 token 到 localStorage
        if (token) {
          localStorage.setItem('auth_token', token)
        }
        // 同时保存用户名，供首页导航栏显示
        if (user?.username) {
          localStorage.setItem('username', user.username)
        }
        
        console.log('Login success:', { token, user })
        
        // 触发自定义事件，通知其他组件登录状态已更新
        // 使用 setTimeout 确保事件在下一个事件循环中被触发，以便 HomeView 能捕获
        setTimeout(() => {
          console.log('Dispatching login-success event')
          window.dispatchEvent(new CustomEvent('login-success', {
            detail: { token, user }
          }))
        }, 0)
        
        this.message = '登录成功！'
        this.username = user?.username || this.form.username
        
        // 登录成功后跳转到首页
        // 使用 await 确保路由跳转完成后再关闭加载状态
        await this.$router.push('/')
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
  background: #2563eb;
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

.welcome {
  margin-bottom: 8px;
  font-size: 14px;
  color: #0f172a;
}

.link {
  display: inline-block;
  margin-top: 16px;
  color: #2563eb;
  text-decoration: none;
}
</style>
