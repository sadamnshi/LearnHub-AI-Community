<template>
  <div class="ai-chat-container">
    <!-- 聊天窗口 -->
    <div class="chat-window">
      <!-- 聊天记录 -->
      <div class="chat-history" ref="chatHistoryRef">
        <template v-if="messages.length === 0">
          <div class="empty-state">
            <div class="empty-icon">🤖</div>
            <p class="empty-text">开始聊天吧！</p>
            <p class="empty-desc">与 AI 助手进行对话，获取帮助和建议</p>
          </div>
        </template>

        <template v-else>
          <div
            v-for="(msg, index) in messages"
            :key="index"
            :class="['message-item', msg.role]"
          >
            <div class="message-avatar">
              <span v-if="msg.role === 'user'">👤</span>
              <span v-else>🤖</span>
            </div>
            <div class="message-content">
              <div class="message-role">
                {{ msg.role === 'user' ? '你' : 'AI 助手' }}
              </div>
              <div class="message-text">{{ msg.content }}</div>
              <div class="message-time">{{ formatTime(msg.time) }}</div>
            </div>
          </div>

          <!-- 加载指示器 -->
          <div v-if="loading" class="message-item assistant">
            <div class="message-avatar">🤖</div>
            <div class="message-content">
              <div class="message-role">AI 助手</div>
              <div class="loading-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- 聊天输入区 -->
      <div class="chat-input-area">
        <div class="input-wrapper">
          <textarea
            v-model="inputMessage"
            class="message-input"
            placeholder="输入你的问题...（支持 Shift+Enter 换行，Enter 发送）"
            :disabled="loading"
            @keydown.enter.prevent="handleEnter"
          ></textarea>

          <div class="input-footer">
            <div class="char-count">
              {{ inputMessage.length }} / 2000
            </div>

            <div class="button-group">
              <button
                v-if="messages.length > 0"
                class="btn btn-secondary"
                @click="clearHistory"
                :disabled="loading"
              >
                🗑️ 清空历史
              </button>

              <button
                class="btn btn-primary"
                @click="sendMessage"
                :disabled="!inputMessage.trim() || loading"
              >
                {{ loading ? '发送中...' : '发送' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误提示 -->
    <transition name="fade">
      <div v-if="errorMessage" class="error-notification">
        <span class="error-icon">⚠️</span>
        <span class="error-text">{{ errorMessage }}</span>
        <button class="close-btn" @click="errorMessage = ''">×</button>
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
      username: '',
    }
  },
  computed: {
    isLoggedIn() {
      return !!localStorage.getItem('auth_token')
    },
  },
  methods: {
    /**
     * 发送消息
     */
    async sendMessage() {
      if (!this.inputMessage.trim()) {
        this.showError('请输入消息')
        return
      }

      const message = this.inputMessage.trim()
      this.inputMessage = ''

      try {
        this.loading = true

        // 立即显示用户消息
        this.messages.push({
          role: 'user',
          content: message,
          time: new Date(),
        })

        // 滚动到底部
        this.$nextTick(() => {
          this.scrollToBottom()
        })

        // 调用 AI API
        const response = await sendMessage(message)

        if (response?.data?.message) {
          this.messages.push({
            role: 'assistant',
            content: response.data.message,
            time: new Date(),
          })

          // 再次滚动到底部
          this.$nextTick(() => {
            this.scrollToBottom()
          })
        }
      } catch (error) {
        console.error('❌ 发送消息失败:', error)
        this.showError(`发送失败: ${error.message}`)

        // 移除失败的用户消息
        if (this.messages[this.messages.length - 1].role === 'user') {
          this.messages.pop()
        }
      } finally {
        this.loading = false
      }
    },

    /**
     * 处理 Enter 键
     */
    handleEnter(event) {
      if (event.shiftKey) {
        // Shift+Enter 换行
        this.inputMessage += '\n'
      } else {
        // Enter 发送
        this.sendMessage()
      }
    },

    /**
     * 加载聊天历史
     */
    async loadChatHistory() {
      try {
        const response = await getChatHistory(50)

        if (response?.data && Array.isArray(response.data)) {
          this.messages = response.data.map((msg) => ({
            role: msg.role,
            content: msg.content,
            time: new Date(msg.time),
          }))

          // 滚动到底部
          this.$nextTick(() => {
            this.scrollToBottom()
          })
        }
      } catch (error) {
        console.error('❌ 加载聊天历史失败:', error)
        // 不显示错误，静默处理
      }
    },

    /**
     * 清空聊天历史
     */
    async clearHistory() {
      if (!confirm('确定要清空所有聊天记录吗？')) {
        return
      }

      try {
        this.loading = true
        await clearChatHistory()
        this.messages = []
        this.showError('聊天历史已清空')
      } catch (error) {
        console.error('❌ 清空聊天历史失败:', error)
        this.showError(`清空失败: ${error.message}`)
      } finally {
        this.loading = false
      }
    },

    /**
     * 滚动到聊天窗口底部
     */
    scrollToBottom() {
      const chatWindow = this.$refs.chatHistoryRef
      if (chatWindow) {
        chatWindow.scrollTop = chatWindow.scrollHeight
      }
    },

    /**
     * 显示错误信息
     */
    showError(message) {
      this.errorMessage = message
      setTimeout(() => {
        this.errorMessage = ''
      }, 4000)
    },

    /**
     * 格式化时间
     */
    formatTime(time) {
      if (!time) return ''

      const date = new Date(time)
      const now = new Date()
      const diffMs = now - date
      const diffMins = Math.floor(diffMs / 60000)

      if (diffMins < 1) {
        return '刚刚'
      } else if (diffMins < 60) {
        return `${diffMins} 分钟前`
      } else if (diffMins < 1440) {
        return `${Math.floor(diffMins / 60)} 小时前`
      } else {
        return date.toLocaleDateString('zh-CN', {
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
        })
      }
    },
  },
  mounted() {
    // 检查登录状态
    if (!this.isLoggedIn) {
      this.$router.push('/login')
      return
    }

    // 获取用户名
    this.username = localStorage.getItem('username') || '用户'

    // 加载聊天历史
    this.loadChatHistory()
  },
}
</script>

<style scoped>
.ai-chat-container {
  display: flex;
  height: 100vh;
  background: var(--background);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue',
    Arial, sans-serif;
}

.chat-window {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 900px;
  height: 100%;
  margin: 0 auto;
  background: var(--card);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--border);
}

/* 聊天历史区 */
.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: var(--background);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--muted-foreground);
}

.empty-icon {
  font-size: 80px;
  margin-bottom: 20px;
  opacity: 0.7;
}

.empty-text {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--foreground);
}

.empty-desc {
  font-size: 14px;
  color: var(--muted-foreground);
}

/* 消息项 */
.message-item {
  display: flex;
  margin-bottom: 16px;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-item.user {
  flex-direction: row-reverse;
}

.message-avatar {
  font-size: 32px;
  margin: 0 12px;
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  background: var(--card);
  border-radius: var(--radius);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid var(--border);
}

.message-item.user .message-content {
  background: var(--primary);
  color: var(--primary-foreground);
  border-radius: var(--radius) 4px var(--radius) var(--radius);
  border: none;
}

.message-item.assistant .message-content {
  background: var(--secondary);
  border-radius: 4px var(--radius) var(--radius) var(--radius);
}

.message-role {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 4px;
  opacity: 0.7;
}

.message-text {
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
  white-space: pre-wrap;
}

.message-time {
  font-size: 11px;
  margin-top: 6px;
  opacity: 0.6;
}

/* 加载动画 */
.loading-dots {
  display: flex;
  gap: 4px;
}

.loading-dots span {
  width: 8px;
  height: 8px;
  background: var(--primary);
  border-radius: 50%;
  animation: bounce 1.4s infinite;
}

.loading-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.loading-dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes bounce {
  0%,
  80%,
  100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-8px);
  }
}

/* 输入区 */
.chat-input-area {
  padding: 20px;
  background: var(--card);
  border-top: 1px solid var(--border);
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-input {
  padding: 12px;
  border: 2px solid var(--border);
  border-radius: var(--radius);
  font-size: 14px;
  font-family: inherit;
  resize: none;
  max-height: 120px;
  min-height: 60px;
  transition: border-color 0.2s;
  background: var(--input-background);
  color: var(--foreground);
}

.message-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.message-input:disabled {
  background: var(--accent);
  color: var(--muted-foreground);
}

/* 输入底部 */
.input-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.char-count {
  font-size: 12px;
  color: #999;
}

/* 按钮组 */
.button-group {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: var(--radius);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--primary);
  color: var(--primary-foreground);
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary);
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.btn-primary:disabled {
  background: var(--muted);
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--secondary);
  color: var(--secondary-foreground);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--secondary);
  opacity: 0.9;
  transform: translateY(-2px);
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 错误通知 */
.error-notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 12px 16px;
  background: var(--destructive);
  color: var(--destructive-foreground);
  border-radius: var(--radius);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 8px;
  animation: slideInUp 0.3s ease-out;
}

.error-icon {
  font-size: 18px;
}

.error-text {
  font-size: 14px;
}

.close-btn {
  background: none;
  border: none;
  color: var(--destructive-foreground);
  font-size: 20px;
  cursor: pointer;
  padding: 0;
  margin-left: 8px;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .message-content {
    max-width: 85%;
  }

  .ai-chat-container {
    background: white;
  }

  .chat-window {
    border-radius: 0;
    box-shadow: none;
  }

  .button-group {
    width: 100%;
    flex-direction: column;
  }

  .btn {
    flex: 1;
  }

  .input-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .chat-history {
    padding: 12px;
  }

  .chat-input-area {
    padding: 12px;
  }

  .message-avatar {
    font-size: 24px;
    margin: 0 8px;
  }

  .message-content {
    max-width: 90%;
    padding: 10px 12px;
    font-size: 13px;
  }
}
</style> 
