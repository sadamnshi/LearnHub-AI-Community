<template>
  <div class="ai-chat-container">
    <div class="chat-window">
      <!-- Chat History -->
      <div class="chat-history" ref="chatHistoryRef">
        <template v-if="messages.length === 0">
          <div class="empty-state">
            <div class="empty-orb">
              <span class="empty-icon">✨</span>
            </div>
            <p class="empty-title">开始与 AI 对话</p>
            <p class="empty-desc">输入你的问题，获取智能帮助与建议</p>
          </div>
        </template>

        <template v-else>
          <div
            v-for="(msg, index) in messages"
            :key="index"
            :class="['message-item', msg.role]"
          >
            <div class="message-avatar">
              <span v-if="msg.role === 'user'" class="avatar-user">你</span>
              <span v-else class="avatar-ai">AI</span>
            </div>
            <div class="message-bubble">
              <div class="message-role">{{ msg.role === 'user' ? '你' : 'AI 助手' }}</div>
              <div class="message-text">{{ msg.content }}</div>
              <div class="message-time">{{ formatTime(msg.time) }}</div>
            </div>
          </div>

          <!-- Loading -->
          <div v-if="loading" class="message-item assistant">
            <div class="message-avatar"><span class="avatar-ai">AI</span></div>
            <div class="message-bubble">
              <div class="message-role">AI 助手</div>
              <div class="loading-dots">
                <span></span><span></span><span></span>
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- Input Area -->
      <div class="chat-input-area">
        <div class="input-wrapper">
          <textarea
            v-model="inputMessage"
            class="message-input"
            placeholder="输入你的问题... (Shift+Enter 换行)"
            :disabled="loading"
            @keydown.enter.prevent="handleEnter"
          ></textarea>
          <div class="input-footer">
            <div class="char-count">{{ inputMessage.length }} / 2000</div>
            <div class="button-group">
              <button v-if="messages.length > 0" class="btn-glass-sm" @click="clearHistory" :disabled="loading">
                清空历史
              </button>
              <button class="btn-accent-sm" @click="sendMessage" :disabled="!inputMessage.trim() || loading">
                {{ loading ? '发送中...' : '发送' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Toast -->
    <transition name="toast">
      <div v-if="errorMessage" class="error-toast">
        <span>{{ errorMessage }}</span>
        <button class="toast-close" @click="errorMessage = ''">✕</button>
      </div>
    </transition>
  </div>
</template>

<script>
import { sendMessage, getChatHistory, clearChatHistory } from '@/api/ai'

export default {
  name: 'AIChatView',
  data() {
    return {
      messages: [],
      inputMessage: '',
      loading: false,
      errorMessage: '',
      username: ''
    }
  },
  computed: {
    isLoggedIn() { return !!localStorage.getItem('auth_token') }
  },
  methods: {
    async sendMessage() {
      if (!this.inputMessage.trim()) { this.showError('请输入消息'); return }
      const message = this.inputMessage.trim()
      this.inputMessage = ''
      try {
        this.loading = true
        this.messages.push({ role: 'user', content: message, time: new Date() })
        this.$nextTick(() => this.scrollToBottom())
        const response = await sendMessage(message)
        if (response?.data?.message) {
          this.messages.push({ role: 'assistant', content: response.data.message, time: new Date() })
          this.$nextTick(() => this.scrollToBottom())
        }
      } catch (error) {
        this.showError(`发送失败: ${error.message}`)
        if (this.messages[this.messages.length - 1]?.role === 'user') this.messages.pop()
      } finally {
        this.loading = false
      }
    },
    handleEnter(event) {
      if (event.shiftKey) this.inputMessage += '\n'
      else this.sendMessage()
    },
    async loadChatHistory() {
      try {
        const response = await getChatHistory(50)
        if (response?.data && Array.isArray(response.data)) {
          this.messages = response.data.map(msg => ({
            role: msg.role, content: msg.content, time: new Date(msg.time)
          }))
          this.$nextTick(() => this.scrollToBottom())
        }
      } catch { /* silent */ }
    },
    async clearHistory() {
      if (!confirm('确定要清空所有聊天记录吗？')) return
      try {
        this.loading = true
        await clearChatHistory()
        this.messages = []
        this.showError('聊天历史已清空')
      } catch (error) {
        this.showError(`清空失败: ${error.message}`)
      } finally {
        this.loading = false
      }
    },
    scrollToBottom() {
      const el = this.$refs.chatHistoryRef
      if (el) el.scrollTop = el.scrollHeight
    },
    showError(message) {
      this.errorMessage = message
      setTimeout(() => { this.errorMessage = '' }, 4000)
    },
    formatTime(time) {
      if (!time) return ''
      const date = new Date(time)
      const diffMins = Math.floor((new Date() - date) / 60000)
      if (diffMins < 1) return '刚刚'
      if (diffMins < 60) return `${diffMins} 分钟前`
      if (diffMins < 1440) return `${Math.floor(diffMins / 60)} 小时前`
      return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
    }
  },
  mounted() {
    if (!this.isLoggedIn) { this.$router.push('/login'); return }
    this.username = localStorage.getItem('username') || '用户'
    this.loadChatHistory()
  }
}
</script>

<style scoped>
.ai-chat-container {
  display: flex;
  height: 100vh;
  padding: var(--space-md);
}

.chat-window {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 880px;
  height: 100%;
  margin: 0 auto;
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid var(--glass-border);
  border-radius: var(--glass-radius);
  box-shadow: var(--glass-shadow);
  overflow: hidden;
}

/* ═══════════════════════════════════════════════════════════════════════
   Chat History
   ═══════════════════════════════════════════════════════════════════════ */

.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-lg);
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  animation: fadeIn 0.6s var(--ease-out-expo);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-orb {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: var(--accent-gradient-soft);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(129, 140, 248, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--space-lg);
  box-shadow: 0 0 60px var(--accent-glow);
  animation: float 4s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-icon { font-size: 40px; }

.empty-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.empty-desc {
  font-size: 14px;
  color: var(--text-tertiary);
}

/* Messages */
.message-item {
  display: flex;
  margin-bottom: var(--space-md);
  animation: msgIn 0.35s var(--ease-out-expo);
}

@keyframes msgIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.message-item.user {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
  margin: 0 10px;
}

.avatar-user, .avatar-ai {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  font-size: 13px;
  font-weight: 700;
}

.avatar-user {
  background: var(--accent-gradient);
  color: #fff;
}

.avatar-ai {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid var(--glass-border);
  color: var(--accent-primary);
}

.message-bubble {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: var(--glass-radius-sm);
  position: relative;
}

.message-item.assistant .message-bubble {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px var(--glass-radius-sm) var(--glass-radius-sm) var(--glass-radius-sm);
}

.message-item.user .message-bubble {
  background: rgba(129, 140, 248, 0.15);
  border: 1px solid rgba(129, 140, 248, 0.2);
  border-radius: var(--glass-radius-sm) 4px var(--glass-radius-sm) var(--glass-radius-sm);
}

.message-role {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-tertiary);
  margin-bottom: 4px;
}

.message-text {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-primary);
  word-break: break-word;
  white-space: pre-wrap;
}

.message-time {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 6px;
}

/* Loading Dots */
.loading-dots {
  display: flex;
  gap: 5px;
  padding: 4px 0;
}

.loading-dots span {
  width: 7px;
  height: 7px;
  background: var(--accent-primary);
  border-radius: 50%;
  animation: bounce 1.4s infinite;
}

.loading-dots span:nth-child(2) { animation-delay: 0.16s; }
.loading-dots span:nth-child(3) { animation-delay: 0.32s; }

@keyframes bounce {
  0%, 80%, 100% { transform: translateY(0); opacity: 0.4; }
  40% { transform: translateY(-8px); opacity: 1; }
}

/* ═══════════════════════════════════════════════════════════════════════
   Input Area
   ═══════════════════════════════════════════════════════════════════════ */

.chat-input-area {
  padding: var(--space-md) var(--space-lg);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.message-input {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--glass-radius-sm);
  font-size: 14px;
  font-family: inherit;
  color: var(--text-primary);
  resize: none;
  max-height: 120px;
  min-height: 56px;
  outline: none;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.message-input::placeholder { color: var(--text-tertiary); }

.message-input:focus {
  border-color: var(--accent-primary);
  background: rgba(129, 140, 248, 0.04);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.message-input:disabled { opacity: 0.5; }

.input-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.char-count {
  font-size: 12px;
  color: var(--text-tertiary);
}

.button-group {
  display: flex;
  gap: 8px;
}

.btn-glass-sm {
  padding: 8px 16px;
  background: var(--glass-bg);
  backdrop-filter: blur(12px);
  border: 1px solid var(--glass-border);
  color: var(--text-secondary);
  border-radius: var(--glass-radius-xs);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-out-expo);
}

.btn-glass-sm:hover:not(:disabled) {
  background: var(--glass-bg-hover);
  color: var(--text-primary);
}

.btn-glass-sm:disabled { opacity: 0.4; cursor: not-allowed; }

.btn-accent-sm {
  padding: 8px 20px;
  background: var(--accent-gradient);
  color: var(--text-on-accent);
  border: none;
  border-radius: var(--glass-radius-xs);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all var(--duration-normal) var(--ease-spring);
  box-shadow: 0 4px 12px var(--accent-glow);
}

.btn-accent-sm:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px var(--accent-glow);
}

.btn-accent-sm:disabled { opacity: 0.4; cursor: not-allowed; }

/* ═══════════════════════════════════════════════════════════════════════
   Error Toast
   ═══════════════════════════════════════════════════════════════════════ */

.error-toast {
  position: fixed;
  bottom: var(--space-lg);
  right: var(--space-lg);
  padding: 12px 18px;
  background: rgba(248, 113, 113, 0.15);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(248, 113, 113, 0.25);
  border-radius: var(--glass-radius-sm);
  color: var(--danger);
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  z-index: 100;
}

.toast-close {
  background: none;
  border: none;
  color: var(--danger);
  font-size: 16px;
  cursor: pointer;
  padding: 0;
  margin-left: 6px;
  opacity: 0.7;
}

.toast-close:hover { opacity: 1; }

.toast-enter-active { transition: all 0.35s var(--ease-spring); }
.toast-leave-active { transition: all 0.2s ease-in; }
.toast-enter-from { opacity: 0; transform: translateY(20px) scale(0.95); }
.toast-leave-to { opacity: 0; transform: translateY(10px); }

/* ═══════════════════════════════════════════════════════════════════════
   Responsive
   ═══════════════════════════════════════════════════════════════════════ */

@media (max-width: 768px) {
  .ai-chat-container { padding: 0; }
  .chat-window { border-radius: 0; border-left: none; border-right: none; }
  .message-bubble { max-width: 82%; }
  .button-group { width: 100%; }
  .btn-glass-sm, .btn-accent-sm { flex: 1; }
  .input-footer { flex-direction: column; align-items: flex-start; gap: 10px; }
}

@media (max-width: 480px) {
  .chat-history { padding: var(--space-md); }
  .chat-input-area { padding: var(--space-md); }
  .avatar-user, .avatar-ai { width: 30px; height: 30px; font-size: 11px; }
  .message-bubble { padding: 10px 12px; font-size: 13px; }
}
</style>
