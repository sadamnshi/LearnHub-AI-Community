<template>
  <div id="app">
    <header >

    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
import { clearExpiredToken } from '@/utils/token'

export default {
  name: 'App',
  created() {
    // 应用启动时检查并清除过期的 Token
    clearExpiredToken()
  }
}
</script>

<style>
/* Global Glassmorphism Flowing Style */
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* Fluid glowing background */
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 20%, #a1c4fd 50%, #c2e9fb 80%, #8fd3f4 100%);
  background-size: 400% 400%;
  animation: gradientFlow 15s ease infinite;
  min-height: 100vh;
  color: #2c3e50;
}

@keyframes gradientFlow {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

:root {
  --glass-bg: rgba(255, 255, 255, 0.25);
  --glass-border: rgba(255, 255, 255, 0.4);
  --glass-blur: blur(16px);
  --glass-shadow: 0 8px 32px rgba(31, 38, 135, 0.15);
}

/* Force specific containers to be transparent so the body background shows through */
#app,
#app .page,
#app .home-page,
#app .ai-chat-container,
#app .create-post-container,
#app .post-detail-page,
#app .container,
#app .profile-container,
#app .main-content,
#app main.main-content {
  background: transparent !important;
}

/* Apply Glassmorphism to core UI elements */
#app .card,
#app .post-card,
#app .navbar,
#app .chat-window,
#app .chat-history,
#app .chat-input-area,
#app .message-item,
#app .message-content,
#app .empty-state,
#app .form-wrapper,
#app .sidebar-card,
#app .modal,
#app .info-card,
#app .profile-section,
#app .post-detail,
#app .category,
#app .tag,
#app .post-footer,
#app .action-btn {
  background: var(--glass-bg) !important;
  backdrop-filter: var(--glass-blur) !important;
  -webkit-backdrop-filter: var(--glass-blur) !important;
  border: 1px solid var(--glass-border) !important;
  box-shadow: var(--glass-shadow) !important;
  border-radius: 16px !important;
  color: #1e293b !important;
}

/* Slightly tweak specific elements for better UI */
#app .navbar {
  border-radius: 0 0 16px 16px !important;
  margin-bottom: 20px;
}

/* Input & Textarea Glass */
#app input,
#app textarea,
#app select {
  background: rgba(255, 255, 255, 0.3) !important;
  border: 1px solid rgba(255, 255, 255, 0.5) !important;
  border-radius: 12px !important;
  color: #0f172a !important;
  backdrop-filter: blur(8px) !important;
  -webkit-backdrop-filter: blur(8px) !important;
  padding: 10px 14px;
}

#app input:focus,
#app textarea:focus,
#app select:focus {
  outline: none !important;
  background: rgba(255, 255, 255, 0.5) !important;
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.6) !important;
}

#app input::placeholder,
#app textarea::placeholder {
  color: rgba(15, 23, 42, 0.5) !important;
}

/* Button Glassmorphism */
#app button,
#app .btn-primary,
#app .btn-secondary,
#app .btn-outline {
  background: rgba(255, 255, 255, 0.4) !important;
  backdrop-filter: blur(8px) !important;
  -webkit-backdrop-filter: blur(8px) !important;
  border: 1px solid rgba(255, 255, 255, 0.6) !important;
  color: #1e293b !important;
  font-weight: 600 !important;
  border-radius: 12px !important;
  transition: all 0.3s ease !important;
}

#app button:hover,
#app .btn-primary:hover,
#app .btn-secondary:hover,
#app .btn-outline:hover {
  background: rgba(255, 255, 255, 0.6) !important;
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1) !important;
}

#app button:disabled {
  background: rgba(255, 255, 255, 0.2) !important;
  color: rgba(30, 41, 59, 0.5) !important;
  cursor: not-allowed;
  transform: none !important;
}

/* Text glow to ensure readability */
#app h1, #app h2, #app h3, #app .navbar-brand {
  text-shadow: 0 2px 4px rgba(255, 255, 255, 0.6) !important;
  color: #1a202c !important;
}

/* Code & Pre Blocks */
#app .content-body code,
#app .content-body pre,
#app .content-body blockquote {
  background: rgba(255, 255, 255, 0.3) !important;
  backdrop-filter: blur(8px) !important;
  -webkit-backdrop-filter: blur(8px) !important;
  border: 1px solid rgba(255, 255, 255, 0.5) !important;
  border-left: 4px solid rgba(255, 255, 255, 0.8) !important; /* For blockquotes */
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
  color: #1e293b !important;
}

#app .content-body pre {
  border-left: 1px solid rgba(255, 255, 255, 0.5) !important;
}

/* Active Nav Links */
#app .nav-link.active {
  background: rgba(255, 255, 255, 0.6) !important;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05) !important;
}
</style>
