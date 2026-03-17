<template>
  <section class="page">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-content">
        <h1>用户中心</h1>
        <button @click="handleLogout" class="btn-logout-header">退出登录</button>
      </div>
    </header>

    <div class="container">
      <!-- 左侧：用户信息卡片 -->
      <div class="main-content">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>

        <!-- 错误提示 -->
        <div v-if="error && !loading" class="alert alert-error">
          <span class="alert-icon">⚠️</span>
          <div>
            <p class="alert-title">加载失败</p>
            <p class="alert-message">{{ error }}</p>
          </div>
        </div>

        <!-- 用户信息展示 -->
        <div v-if="profile && !loading" class="profile-section">
          <!-- 头像和基本信息 -->
          <div class="profile-header">
            <div class="avatar-container">
              <img
                v-if="profile.avatar"
                :src="profile.avatar"
                alt="头像"
                class="avatar"
              />
              <div v-else class="avatar-placeholder">
                {{ profile.username ? profile.username[0].toUpperCase() : '?' }}
              </div>
            </div>

            <div class="user-meta">
              <h2 class="username">{{ profile.username }}</h2>
              <div class="role-badge" :class="'badge-' + profile.role">
                {{ roleLabel(profile.role) }}
              </div>
              <p class="user-id">用户 ID: {{ profile.id }}</p>
            </div>
          </div>

          <!-- 详细信息卡片 -->
          <div class="info-card">
            <div class="info-group">
              <h3 class="group-title">📧 联系方式</h3>
              <div class="info-item">
                <span class="item-label">邮箱</span>
                <span class="item-value">{{ profile.email || '暂未填写' }}</span>
              </div>
            </div>

            <div class="info-group">
              <h3 class="group-title">👤 个人信息</h3>
              <div class="info-item">
                <span class="item-label">个人简介</span>
                <p class="item-value bio">
                  {{ profile.bio || '这个人很懒，什么都没写～' }}
                </p>
              </div>
            </div>

            <div class="info-group">
              <h3 class="group-title">📅 账号信息</h3>
              <div class="info-item">
                <span class="item-label">注册时间</span>
                <span class="item-value">{{ formatDate(profile.created_at) }}</span>
              </div>
              <div class="info-item">
                <span class="item-label">最后更新</span>
                <span class="item-value">{{ formatDate(profile.updated_at) }}</span>
              </div>
            </div>
          </div>

          <!-- 操作按钮组 -->
          <div class="action-buttons">
            <button @click="loadProfile" :disabled="loading" class="btn btn-primary">
              🔄 刷新信息
            </button>
            <button @click="showChangePassword = true" class="btn btn-secondary">
              🔐 修改密码
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧边栏：快速操作 -->
      <aside class="sidebar">
        <div class="sidebar-card">
          <h3>⚡ 快速操作</h3>
          <nav class="quick-menu">
            <a href="#" class="menu-item" @click.prevent="loadProfile">
              <span class="menu-icon">🔄</span>
              <span>刷新信息</span>
            </a>
            <a href="#" class="menu-item" @click.prevent="showChangePassword = true">
              <span class="menu-icon">🔐</span>
              <span>修改密码</span>
            </a>
            <a href="/" class="menu-item">
              <span class="menu-icon">🏠</span>
              <span>返回首页</span>
            </a>
          </nav>
        </div>

        <div class="sidebar-card">
          <h3>✨ 账号状态</h3>
          <div class="status-item">
            <span class="status-label">账号状态</span>
            <span class="status-value active">✓ 正常</span>
          </div>
          <div class="status-item">
            <span class="status-label">邮箱验证</span>
            <span v-if="profile" class="status-value" :class="profile.email ? 'verified' : 'unverified'">
              {{ profile.email ? '✓ 已验证' : '○ 未验证' }}
            </span>
            <span v-else class="status-value unverified">○ 加载中...</span>
          </div>
        </div>
      </aside>
    </div>

    <!-- 修改密码模态框 -->
    <div v-if="showChangePassword" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>修改密码</h3>
          <button @click="closeModal" class="modal-close">✕</button>
        </div>

        <div class="modal-body">
          <!-- 错误提示 -->
          <div v-if="passwordError" class="alert alert-error modal-alert">
            <span class="alert-icon">⚠️</span>
            <p>{{ passwordError }}</p>
          </div>

          <!-- 成功提示 -->
          <div v-if="passwordSuccess" class="alert alert-success modal-alert">
            <span class="alert-icon">✓</span>
            <p>{{ passwordSuccess }}</p>
          </div>

          <!-- 密码输入表单 -->
          <form @submit.prevent="submitChangePassword">
            <div class="form-group">
              <label for="oldPassword" class="form-label">旧密码 *</label>
              <input
                id="oldPassword"
                v-model="passwordForm.oldPassword"
                type="password"
                class="form-input"
                placeholder="请输入旧密码"
                required
                :disabled="changingPassword"
              />
            </div>

            <div class="form-group">
              <label for="newPassword" class="form-label">新密码 *</label>
              <input
                id="newPassword"
                v-model="passwordForm.newPassword"
                type="password"
                class="form-input"
                placeholder="请输入新密码（至少 6 位）"
                required
                minlength="6"
                :disabled="changingPassword"
              />
              <small class="form-hint">密码长度至少 6 位，建议包含大小写字母和数字</small>
            </div>

            <div class="form-group">
              <label for="confirmPassword" class="form-label">确认新密码 *</label>
              <input
                id="confirmPassword"
                v-model="passwordForm.confirmPassword"
                type="password"
                class="form-input"
                placeholder="请再次输入新密码"
                required
                :disabled="changingPassword"
              />
            </div>

            <!-- 密码匹配提示 -->
            <div v-if="passwordForm.newPassword && passwordForm.confirmPassword" class="password-check">
              <span
                v-if="passwordForm.newPassword === passwordForm.confirmPassword"
                class="check-success"
              >
                ✓ 密码匹配
              </span>
              <span v-else class="check-error">✕ 密码不匹配</span>
            </div>
          </form>
        </div>

        <div class="modal-footer">
          <button
            @click="closeModal"
            :disabled="changingPassword"
            class="btn btn-secondary"
          >
            取消
          </button>
          <button
            @click="submitChangePassword"
            :disabled="changingPassword || !isPasswordFormValid"
            class="btn btn-primary"
          >
            {{ changingPassword ? '处理中...' : '确认修改' }}
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import { getUserProfile, logout, updatePassword } from '@/api/auth'

export default {
  name: 'ProfileView',
  data() {
    return {
      // 个人资料
      profile: null,
      loading: false,
      error: '',

      // 修改密码相关
      showChangePassword: false,
      changingPassword: false,
      passwordError: '',
      passwordSuccess: '',
      passwordForm: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    }
  },
  computed: {
    // 密码表单验证
    isPasswordFormValid() {
      return (
        this.passwordForm.oldPassword.length > 0 &&
        this.passwordForm.newPassword.length >= 6 &&
        this.passwordForm.confirmPassword.length >= 6 &&
        this.passwordForm.newPassword === this.passwordForm.confirmPassword
      )
    }
  },
  mounted() {
    this.loadProfile()
  },
  methods: {
    /**
     * 加载用户资料
     */
    async loadProfile() {
      this.loading = true
      this.error = ''
      try {
        const response = await getUserProfile()
        // 后端返回 { code: 0, msg: "...", data: { data: user } } 格式
        console.log('getUserProfile response:', response)
        this.profile = response.data?.data || response.data
        
        // 确保 profile 不为 null
        if (!this.profile) {
          this.error = '用户信息为空'
        }
      } catch (error) {
        this.error = error.message || '加载用户信息失败'
        console.error('Error loading profile:', error)
        // Token 过期，2 秒后跳转到登录页
        if (error.message.includes('401') || error.message.includes('token')) {
          setTimeout(() => {
            this.$router.push('/login')
          }, 2000)
        }
      } finally {
        this.loading = false
      }
    },

    /**
     * 关闭修改密码模态框
     */
    closeModal() {
      this.showChangePassword = false
      // 清空表单和提示
      setTimeout(() => {
        this.passwordForm = {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        }
        this.passwordError = ''
        this.passwordSuccess = ''
      }, 300)
    },

    /**
     * 提交修改密码表单
     */
    async submitChangePassword() {
      // 验证表单
      if (!this.isPasswordFormValid) {
        this.passwordError = '请检查密码输入'
        return
      }

      this.changingPassword = true
      this.passwordError = ''
      this.passwordSuccess = ''

      try {
        const response = await updatePassword({
          old_password: this.passwordForm.oldPassword,
          new_password: this.passwordForm.newPassword
        })

        // 修改成功
        this.passwordSuccess = '密码修改成功！'
        
        // 2 秒后关闭模态框
        setTimeout(() => {
          this.closeModal()
        }, 2000)
      } catch (error) {
        // 密码修改失败
        if (error.message.includes('旧密码')) {
          this.passwordError = '旧密码错误，请重新输入'
        } else if (error.message.includes('相同')) {
          this.passwordError = '新密码不能与旧密码相同'
        } else {
          this.passwordError = error.message || '密码修改失败，请稍后重试'
        }
      } finally {
        this.changingPassword = false
      }
    },

    /**
     * 退出登录
     */
    handleLogout() {
      if (confirm('确定要退出登录吗？')) {
        logout()
        this.$router.push('/login')
      }
    },

    /**
     * 将后端 role 字段转换为中文显示
     */
    roleLabel(role) {
      const map = {
        admin: '管理员',
        moderator: '版主',
        user: '普通用户'
      }
      return map[role] || role || '未知'
    },

    /**
     * 格式化日期时间
     */
    formatDate(dateString) {
      if (!dateString) return '—'
      try {
        const date = new Date(dateString)
        return date.toLocaleString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch {
        return dateString
      }
    }
  }
}
</script>

<style scoped>
/* ============================================================================
   全局样式
   ============================================================================ */

* {
  box-sizing: border-box;
}

.page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding-top: 0;
}

/* ============================================================================
   顶部导航栏
   ============================================================================ */

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 16px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.btn-logout-header {
  padding: 8px 16px;
  background: #dc2626;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
}

.btn-logout-header:hover {
  background: #b91c1c;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3);
}

/* ============================================================================
   主容器和布局
   ============================================================================ */

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 24px;
}

@media (max-width: 768px) {
  .container {
    grid-template-columns: 1fr;
    padding: 16px 12px;
    gap: 16px;
  }
}

/* ============================================================================
   加载状态
   ============================================================================ */

.loading-container {
  background: #fff;
  border-radius: 12px;
  padding: 48px 24px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ============================================================================
   提示信息
   ============================================================================ */

.alert {
  border-radius: 10px;
  padding: 16px;
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.alert-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.alert-title {
  font-weight: 600;
  margin: 0 0 4px 0;
  font-size: 14px;
}

.alert-message {
  margin: 0;
  font-size: 13px;
}

.alert-error {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fecaca;
}

.alert-success {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.modal-alert {
  margin-bottom: 16px;
}

/* ============================================================================
   主要内容区域
   ============================================================================ */

.main-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ============================================================================
   用户资料卡片
   ============================================================================ */

.profile-section {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

/* 资料头部 */
.profile-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 32px 24px;
  display: flex;
  gap: 24px;
  align-items: flex-start;
  color: #fff;
}

@media (max-width: 480px) {
  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 24px 16px;
  }
}

/* 头像 */
.avatar-container {
  flex-shrink: 0;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.avatar-placeholder {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  font-size: 40px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4px solid rgba(255, 255, 255, 0.3);
}

/* 用户元数据 */
.user-meta {
  flex: 1;
}

.username {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.role-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 12px;
  background: rgba(255, 255, 255, 0.2);
}

.badge-admin {
  background: rgba(255, 255, 255, 0.25) !important;
}

.badge-moderator {
  background: rgba(255, 255, 255, 0.25) !important;
}

.badge-user {
  background: rgba(255, 255, 255, 0.25) !important;
}

.user-id {
  margin: 0;
  font-size: 13px;
  opacity: 0.9;
}

/* 信息卡片 */
.info-card {
  padding: 24px;
  border-top: 1px solid #e2e8f0;
}

.info-group {
  margin-bottom: 24px;
}

.info-group:last-child {
  margin-bottom: 0;
}

.group-title {
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  margin: 0 0 12px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 0;
}

.info-item:last-child {
  margin-bottom: 0;
}

.item-label {
  width: 100px;
  flex-shrink: 0;
  color: #64748b;
  font-weight: 500;
  font-size: 13px;
}

.item-value {
  flex: 1;
  color: #1e293b;
  word-break: break-word;
  font-size: 13px;
}

.bio {
  color: #475569;
  font-style: italic;
  line-height: 1.6;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 12px;
  padding: 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

@media (max-width: 480px) {
  .action-buttons {
    flex-direction: column;
  }
}

.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: #e2e8f0;
  color: #475569;
}

.btn-secondary:hover:not(:disabled) {
  background: #cbd5e1;
  transform: translateY(-2px);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* ============================================================================
   右侧边栏
   ============================================================================ */

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.sidebar-card h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #1e293b;
  font-weight: 700;
}

/* 快速菜单 */
.quick-menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  color: #475569;
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 14px;
}

.menu-item:hover {
  background: #f1f5f9;
  color: #667eea;
}

.menu-icon {
  font-size: 16px;
}

/* 账号状态 */
.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #e2e8f0;
  font-size: 13px;
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  color: #64748b;
  font-weight: 500;
}

.status-value {
  font-weight: 600;
}

.status-value.active {
  color: #16a34a;
}

.status-value.verified {
  color: #16a34a;
}

.status-value.unverified {
  color: #ea580c;
}

/* ============================================================================
   模态框
   ============================================================================ */

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
  max-width: 400px;
  width: 100%;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #1e293b;
  font-weight: 700;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s;
}

.modal-close:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

.modal-footer .btn {
  flex: 1;
}

/* ============================================================================
   表单
   ============================================================================ */

.form-group {
  margin-bottom: 16px;
}

.form-group:last-of-type {
  margin-bottom: 8px;
}

.form-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 6px;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s;
  font-family: inherit;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input:disabled {
  background: #f1f5f9;
  color: #94a3b8;
  cursor: not-allowed;
}

.form-hint {
  display: block;
  font-size: 12px;
  color: #64748b;
  margin-top: 6px;
}

/* 密码检查 */
.password-check {
  font-size: 13px;
  font-weight: 600;
  margin: 12px 0 16px 0;
  padding: 8px 12px;
  border-radius: 6px;
  text-align: center;
}

.check-success {
  color: #16a34a;
  background: #dcfce7;
  display: block;
}

.check-error {
  color: #dc2626;
  background: #fee2e2;
  display: block;
}

/* ============================================================================
   响应式设计
   ============================================================================ */

@media (max-width: 1024px) {
  .container {
    grid-template-columns: 1fr;
  }

  .sidebar {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
  }
}

@media (max-width: 768px) {
  .header h1 {
    font-size: 24px;
  }

  .profile-header {
    padding: 20px 16px;
    gap: 16px;
  }

  .avatar,
  .avatar-placeholder {
    width: 80px;
    height: 80px;
    font-size: 32px;
  }

  .username {
    font-size: 24px;
  }

  .info-card {
    padding: 16px;
  }

  .action-buttons {
    padding: 16px;
  }

  .sidebar {
    grid-template-columns: 1fr;
  }

  .modal {
    max-width: 90%;
  }
}

@media (max-width: 480px) {
  .page {
    padding-top: 0;
  }

  .header-content {
    padding: 0 16px;
    flex-direction: column;
    gap: 12px;
    text-align: center;
  }

  .header h1 {
    font-size: 20px;
  }

  .btn-logout-header {
    width: 100%;
  }

  .container {
    padding: 16px;
    gap: 12px;
  }

  .profile-header {
    padding: 16px;
  }

  .avatar,
  .avatar-placeholder {
    width: 70px;
    height: 70px;
    font-size: 28px;
  }

  .username {
    font-size: 20px;
  }

  .item-label {
    width: 70px;
  }

  .modal {
    max-width: calc(100% - 32px);
  }

  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 16px;
  }
}
</style>
