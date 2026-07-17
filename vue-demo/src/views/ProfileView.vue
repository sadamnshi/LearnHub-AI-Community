<template>
  <section class="profile-page">
    <!-- Header -->
    <header class="header">
      <div class="header-inner">
        <h1>用户中心</h1>
        <button @click="handleLogout" class="btn-logout">退出登录</button>
      </div>
    </header>

    <div class="container">
      <!-- Main -->
      <div class="main-content">
        <!-- Loading -->
        <div v-if="loading" class="loading-container glass-card">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>

        <!-- Error -->
        <div v-if="error && !loading" class="alert-error glass-card">
          <span>⚠️</span>
          <div><strong>加载失败</strong><p>{{ error }}</p></div>
        </div>

        <!-- Profile -->
        <div v-if="profile && !loading" class="profile-card glass-card">
          <!-- Profile Header -->
          <div class="profile-header">
            <div class="avatar-container">
              <img v-if="profile.avatar" :src="profile.avatar" alt="头像" class="avatar" />
              <div v-else class="avatar-placeholder">
                {{ profile.username ? profile.username[0].toUpperCase() : '?' }}
              </div>
            </div>
            <div class="user-meta">
              <h2 class="username">{{ profile.username }}</h2>
              <div class="role-badge">{{ roleLabel(profile.role) }}</div>
              <p class="user-id">ID: {{ profile.id }}</p>
            </div>
          </div>

          <!-- Info -->
          <div class="info-section">
            <div class="info-group">
              <h3>联系方式</h3>
              <div class="info-item">
                <span class="info-label">邮箱</span>
                <span class="info-value">{{ profile.email || '暂未填写' }}</span>
              </div>
            </div>
            <div class="info-group">
              <h3>个人简介</h3>
              <p class="bio">{{ profile.bio || '这个人很懒，什么都没写～' }}</p>
            </div>
            <div class="info-group">
              <h3>账号信息</h3>
              <div class="info-item">
                <span class="info-label">注册时间</span>
                <span class="info-value">{{ formatDate(profile.created_at) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">最后更新</span>
                <span class="info-value">{{ formatDate(profile.updated_at) }}</span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="action-bar">
            <button @click="loadProfile" :disabled="loading" class="btn-glass">刷新信息</button>
            <button @click="showChangePassword = true" class="btn-accent">修改密码</button>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <aside class="sidebar">
        <div class="sidebar-card glass-card">
          <h3>快速操作</h3>
          <nav class="quick-menu">
            <a href="#" class="menu-item" @click.prevent="loadProfile">
              <span class="menu-dot"></span>刷新信息
            </a>
            <a href="#" class="menu-item" @click.prevent="showChangePassword = true">
              <span class="menu-dot"></span>修改密码
            </a>
            <router-link to="/" class="menu-item">
              <span class="menu-dot"></span>返回首页
            </router-link>
          </nav>
        </div>

        <div class="sidebar-card glass-card">
          <h3>账号状态</h3>
          <div class="status-row">
            <span>账号状态</span>
            <span class="status-active">
              <span class="pulse-dot"></span> 正常
            </span>
          </div>
          <div class="status-row">
            <span>邮箱验证</span>
            <span v-if="profile" :class="profile.email ? 'status-verified' : 'status-unverified'">
              {{ profile.email ? '已验证' : '未验证' }}
            </span>
          </div>
        </div>
      </aside>
    </div>

    <!-- Password Modal -->
    <transition name="modal-fade">
      <div v-if="showChangePassword" class="modal-overlay" @click="closeModal">
        <div class="modal glass-card" @click.stop>
          <div class="modal-header">
            <h3>修改密码</h3>
            <button @click="closeModal" class="modal-close">✕</button>
          </div>
          <div class="modal-body">
            <div v-if="passwordError" class="alert-error modal-alert">{{ passwordError }}</div>
            <div v-if="passwordSuccess" class="alert-success modal-alert">{{ passwordSuccess }}</div>

            <form @submit.prevent="submitChangePassword">
              <div class="form-field">
                <label>旧密码 *</label>
                <input v-model="passwordForm.oldPassword" type="password" placeholder="请输入旧密码" required :disabled="changingPassword" />
              </div>
              <div class="form-field">
                <label>新密码 *</label>
                <input v-model="passwordForm.newPassword" type="password" placeholder="至少 6 位" required minlength="6" :disabled="changingPassword" />
                <small class="form-hint">建议包含大小写字母和数字</small>
              </div>
              <div class="form-field">
                <label>确认新密码 *</label>
                <input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入" required :disabled="changingPassword" />
              </div>
              <div v-if="passwordForm.newPassword && passwordForm.confirmPassword" class="password-match" :class="{ match: passwordForm.newPassword === passwordForm.confirmPassword }">
                {{ passwordForm.newPassword === passwordForm.confirmPassword ? '✓ 密码匹配' : '✕ 密码不匹配' }}
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button @click="closeModal" :disabled="changingPassword" class="btn-glass">取消</button>
            <button @click="submitChangePassword" :disabled="changingPassword || !isPasswordFormValid" class="btn-accent">
              {{ changingPassword ? '处理中...' : '确认修改' }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </section>
</template>

<script>
import { getUserProfile, logout, updatePassword } from '@/api/auth'

export default {
  name: 'ProfileView',
  data() {
    return {
      profile: null,
      loading: false,
      error: '',
      showChangePassword: false,
      changingPassword: false,
      passwordError: '',
      passwordSuccess: '',
      passwordForm: { oldPassword: '', newPassword: '', confirmPassword: '' }
    }
  },
  computed: {
    isPasswordFormValid() {
      return (
        this.passwordForm.oldPassword.length > 0 &&
        this.passwordForm.newPassword.length >= 6 &&
        this.passwordForm.confirmPassword.length >= 6 &&
        this.passwordForm.newPassword === this.passwordForm.confirmPassword
      )
    }
  },
  mounted() { this.loadProfile() },
  methods: {
    async loadProfile() {
      this.loading = true
      this.error = ''
      try {
        const response = await getUserProfile()
        this.profile = response.data?.data || response.data
        if (!this.profile) this.error = '用户信息为空'
      } catch (error) {
        this.error = error.message || '加载用户信息失败'
        if (error.message.includes('401') || error.message.includes('token')) {
          setTimeout(() => { this.$router.push('/login') }, 2000)
        }
      } finally {
        this.loading = false
      }
    },
    closeModal() {
      this.showChangePassword = false
      setTimeout(() => {
        this.passwordForm = { oldPassword: '', newPassword: '', confirmPassword: '' }
        this.passwordError = ''
        this.passwordSuccess = ''
      }, 300)
    },
    async submitChangePassword() {
      if (!this.isPasswordFormValid) return
      this.changingPassword = true
      this.passwordError = ''
      this.passwordSuccess = ''
      try {
        await updatePassword({
          old_password: this.passwordForm.oldPassword,
          new_password: this.passwordForm.newPassword
        })
        this.passwordSuccess = '密码修改成功！'
        setTimeout(() => { this.closeModal() }, 2000)
      } catch (error) {
        if (error.message.includes('旧密码')) {
          this.passwordError = '旧密码错误，请重新输入'
        } else if (error.message.includes('相同')) {
          this.passwordError = '新密码不能与旧密码相同'
        } else {
          this.passwordError = error.message || '密码修改失败'
        }
      } finally {
        this.changingPassword = false
      }
    },
    handleLogout() {
      if (confirm('确定要退出登录吗？')) {
        logout()
        this.$router.push('/login')
      }
    },
    roleLabel(role) {
      return { admin: '管理员', moderator: '版主', user: '普通用户' }[role] || role || '未知'
    },
    formatDate(dateString) {
      if (!dateString) return '—'
      try {
        return new Date(dateString).toLocaleString('zh-CN', {
          year: 'numeric', month: '2-digit', day: '2-digit',
          hour: '2-digit', minute: '2-digit'
        })
      } catch { return dateString }
    }
  }
}
</script>

<style scoped>
/* ═══════════════════════════════════════════════════════════════════════
   Header
   ═══════════════════════════════════════════════════════════════════════ */

.header {
  background: rgba(12, 10, 29, 0.6);
  backdrop-filter: blur(24px) saturate(1.5);
  -webkit-backdrop-filter: blur(24px) saturate(1.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-md) var(--space-lg);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.btn-logout {
  padding: 8px 16px;
  background: rgba(248, 113, 113, 0.15);
  border: 1px solid rgba(248, 113, 113, 0.2);
  color: var(--danger);
  border-radius: var(--glass-radius-xs);
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-logout:hover {
  background: rgba(248, 113, 113, 0.25);
  transform: translateY(-1px);
}

/* ═══════════════════════════════════════════════════════════════════════
   Layout
   ═══════════════════════════════════════════════════════════════════════ */

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: var(--space-lg);
}

.glass-card {
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid var(--glass-border);
  border-radius: var(--glass-radius);
  box-shadow: var(--glass-shadow);
}

/* ═══════════════════════════════════════════════════════════════════════
   Loading / Error
   ═══════════════════════════════════════════════════════════════════════ */

.loading-container {
  padding: var(--space-2xl);
  text-align: center;
  color: var(--text-tertiary);
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid rgba(255, 255, 255, 0.08);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto var(--space-md);
}

@keyframes spin { to { transform: rotate(360deg); } }

.alert-error {
  padding: var(--space-md);
  background: var(--danger-glow);
  border: 1px solid rgba(248, 113, 113, 0.2);
  border-radius: var(--glass-radius-sm);
  color: var(--danger);
  font-size: 14px;
  display: flex;
  gap: var(--space-sm);
}

.alert-error strong { display: block; margin-bottom: 4px; }
.alert-error p { margin: 0; font-size: 13px; }

/* ═══════════════════════════════════════════════════════════════════════
   Profile Card
   ═══════════════════════════════════════════════════════════════════════ */

.profile-card {
  overflow: hidden;
  animation: fadeInUp 0.5s var(--ease-out-expo);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}

.profile-header {
  padding: var(--space-xl);
  display: flex;
  gap: var(--space-lg);
  align-items: flex-start;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.03);
}

.avatar-container { flex-shrink: 0; }

.avatar, .avatar-placeholder {
  width: 90px;
  height: 90px;
  border-radius: 50%;
  border: 3px solid rgba(255, 255, 255, 0.12);
  transition: all var(--duration-slow) var(--ease-spring);
}

.avatar {
  object-fit: cover;
}

.avatar-placeholder {
  background: var(--accent-gradient-soft);
  backdrop-filter: blur(8px);
  color: var(--text-primary);
  font-size: 36px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-container:hover .avatar,
.avatar-container:hover .avatar-placeholder {
  transform: scale(1.08) rotate(3deg);
  box-shadow: 0 12px 32px var(--accent-glow);
}

.user-meta { flex: 1; }

.username {
  font-size: 26px;
  font-weight: 800;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.3px;
}

.role-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: rgba(129, 140, 248, 0.12);
  color: var(--accent-primary);
  border: 1px solid rgba(129, 140, 248, 0.2);
  margin-bottom: 8px;
}

.user-id {
  margin: 0;
  font-size: 13px;
  color: var(--text-tertiary);
}

/* Info Section */
.info-section {
  padding: var(--space-lg) var(--space-xl);
}

.info-group {
  margin-bottom: var(--space-lg);
}

.info-group:last-child { margin-bottom: 0; }

.info-group h3 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-tertiary);
  margin: 0 0 var(--space-md) 0;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.info-item {
  display: flex;
  gap: var(--space-md);
  padding: 10px 14px;
  border-radius: var(--glass-radius-xs);
  margin-bottom: 6px;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.info-item:hover {
  background: rgba(255, 255, 255, 0.04);
  transform: translateX(4px);
}

.info-label {
  width: 90px;
  flex-shrink: 0;
  color: var(--text-tertiary);
  font-size: 13px;
  font-weight: 500;
}

.info-value {
  color: var(--text-secondary);
  font-size: 13px;
}

.bio {
  color: var(--text-secondary);
  font-style: italic;
  line-height: 1.7;
  font-size: 14px;
  padding: 10px 14px;
  margin: 0;
}

/* Actions */
.action-bar {
  display: flex;
  gap: var(--space-sm);
  padding: var(--space-lg) var(--space-xl);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
}

.btn-glass {
  padding: 10px 20px;
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  color: var(--text-primary);
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-glass:hover:not(:disabled) {
  background: var(--glass-bg-hover);
  transform: translateY(-1px);
}

.btn-glass:disabled { opacity: 0.5; cursor: not-allowed; }

.btn-accent {
  padding: 10px 20px;
  background: var(--accent-gradient);
  color: var(--text-on-accent);
  border: none;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-spring);
  box-shadow: 0 4px 12px var(--accent-glow);
}

.btn-accent:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px var(--accent-glow);
}

.btn-accent:disabled { opacity: 0.5; cursor: not-allowed; }

/* ═══════════════════════════════════════════════════════════════════════
   Sidebar
   ═══════════════════════════════════════════════════════════════════════ */

.sidebar {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  animation: fadeInUp 0.6s var(--ease-out-expo) 100ms both;
}

.sidebar-card {
  padding: var(--space-lg);
}

.sidebar-card h3 {
  margin: 0 0 var(--space-md) 0;
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
}

.quick-menu {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  color: var(--text-secondary);
  text-decoration: none;
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-primary);
  padding-left: 18px;
}

.menu-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent-primary);
  opacity: 0.5;
  transition: opacity var(--duration-fast);
}

.menu-item:hover .menu-dot { opacity: 1; }

.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  font-size: 13px;
  color: var(--text-secondary);
}

.status-row:last-child { border-bottom: none; }

.status-active {
  color: var(--success);
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success);
  box-shadow: 0 0 8px var(--success-glow);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.4); }
  70% { box-shadow: 0 0 0 8px rgba(52, 211, 153, 0); }
  100% { box-shadow: 0 0 0 0 rgba(52, 211, 153, 0); }
}

.status-verified { color: var(--success); font-weight: 600; }
.status-unverified { color: var(--warning); font-weight: 600; }

/* ═══════════════════════════════════════════════════════════════════════
   Modal
   ═══════════════════════════════════════════════════════════════════════ */

.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-md);
}

.modal {
  max-width: 420px;
  width: 100%;
  background: rgba(20, 18, 40, 0.85) !important;
  backdrop-filter: blur(32px) saturate(1.5) !important;
  -webkit-backdrop-filter: blur(32px) saturate(1.5) !important;
  border: 1px solid var(--glass-border-bright) !important;
  animation: modalIn 0.35s var(--ease-spring);
}

@keyframes modalIn {
  from { transform: translateY(20px) scale(0.96); opacity: 0; }
  to { transform: translateY(0) scale(1); opacity: 1; }
}

.modal-fade-enter-active { transition: opacity 0.3s; }
.modal-fade-leave-active { transition: opacity 0.2s; }
.modal-fade-enter-from,
.modal-fade-leave-to { opacity: 0; }

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-lg);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.modal-close {
  background: none;
  border: none;
  color: var(--text-tertiary);
  font-size: 20px;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all var(--duration-fast);
}

.modal-close:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-primary);
}

.modal-body { padding: var(--space-lg); }

.modal-alert { margin-bottom: var(--space-md); }

.modal-footer {
  display: flex;
  gap: var(--space-sm);
  padding: var(--space-lg);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.modal-footer .btn-glass,
.modal-footer .btn-accent { flex: 1; justify-content: center; }

/* Form fields in modal */
.form-field {
  margin-bottom: var(--space-md);
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
  padding: 11px 14px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--glass-radius-xs);
  font-size: 14px;
  color: var(--text-primary);
  font-family: inherit;
  outline: none;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.form-field input::placeholder { color: var(--text-tertiary); }

.form-field input:focus {
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.form-field input:disabled { opacity: 0.5; cursor: not-allowed; }

.form-hint {
  display: block;
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 6px;
}

.password-match {
  text-align: center;
  padding: 8px;
  border-radius: var(--glass-radius-xs);
  font-size: 13px;
  font-weight: 600;
  color: var(--danger);
  background: var(--danger-glow);
  margin-bottom: var(--space-md);
}

.password-match.match {
  color: var(--success);
  background: var(--success-glow);
}

.alert-success {
  padding: 12px 16px;
  background: var(--success-glow);
  border: 1px solid rgba(52, 211, 153, 0.2);
  border-radius: var(--glass-radius-xs);
  color: var(--success);
  font-size: 14px;
  font-weight: 500;
}

/* ═══════════════════════════════════════════════════════════════════════
   Responsive
   ═══════════════════════════════════════════════════════════════════════ */

@media (max-width: 768px) {
  .container {
    grid-template-columns: 1fr;
    padding: var(--space-md);
  }
  .sidebar {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--space-sm);
  }
  .profile-header { flex-direction: column; align-items: center; text-align: center; }
}

@media (max-width: 480px) {
  .header-inner { padding: var(--space-md); flex-direction: column; gap: var(--space-sm); }
  .sidebar { grid-template-columns: 1fr; }
  .action-bar { flex-direction: column; }
  .avatar, .avatar-placeholder { width: 72px; height: 72px; font-size: 28px; }
  .username { font-size: 22px; }
}
</style>
