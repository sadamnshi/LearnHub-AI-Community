<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import { clearExpiredToken } from '@/utils/token'

export default {
  name: 'App',
  created() {
    clearExpiredToken()
  }
}
</script>

<style>
/* ═══════════════════════════════════════════════════════════════════════
   Premium Aurora Glass — Global Foundation
   ═══════════════════════════════════════════════════════════════════════ */

:root {
  /* Glass tokens */
  --glass-bg: rgba(255, 255, 255, 0.08);
  --glass-bg-hover: rgba(255, 255, 255, 0.12);
  --glass-bg-active: rgba(255, 255, 255, 0.16);
  --glass-border: rgba(255, 255, 255, 0.15);
  --glass-border-bright: rgba(255, 255, 255, 0.25);
  --glass-blur: blur(20px);
  --glass-shadow: 0 8px 32px rgba(0, 0, 0, 0.25);
  --glass-shadow-hover: 0 16px 48px rgba(0, 0, 0, 0.3);
  --glass-radius: 20px;
  --glass-radius-sm: 14px;
  --glass-radius-xs: 10px;

  /* Accent palette */
  --accent-primary: #818cf8;
  --accent-secondary: #a78bfa;
  --accent-glow: rgba(129, 140, 248, 0.4);
  --accent-gradient: linear-gradient(135deg, #818cf8 0%, #a78bfa 50%, #c084fc 100%);
  --accent-gradient-soft: linear-gradient(135deg, rgba(129, 140, 248, 0.3) 0%, rgba(167, 139, 250, 0.3) 100%);

  /* Text */
  --text-primary: #f1f5f9;
  --text-secondary: rgba(241, 245, 249, 0.7);
  --text-tertiary: rgba(241, 245, 249, 0.45);
  --text-on-accent: #ffffff;

  /* Semantic */
  --danger: #f87171;
  --danger-glow: rgba(248, 113, 113, 0.3);
  --success: #34d399;
  --success-glow: rgba(52, 211, 153, 0.3);
  --warning: #fbbf24;

  /* Spacing */
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;

  /* Transitions */
  --ease-out-expo: cubic-bezier(0.16, 1, 0.3, 1);
  --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);
  --duration-fast: 150ms;
  --duration-normal: 250ms;
  --duration-slow: 400ms;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  scroll-behavior: smooth;
}

body {
  font-family: 'Outfit', 'Noto Sans SC', -apple-system, BlinkMacSystemFont, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--text-primary);
  min-height: 100vh;
  overflow-x: hidden;

  /* Deep aurora gradient background */
  background: #0c0a1d;
  background-image:
    radial-gradient(ellipse 80% 60% at 10% 20%, rgba(99, 102, 241, 0.25) 0%, transparent 60%),
    radial-gradient(ellipse 60% 50% at 85% 15%, rgba(168, 85, 247, 0.2) 0%, transparent 55%),
    radial-gradient(ellipse 70% 55% at 50% 80%, rgba(59, 130, 246, 0.18) 0%, transparent 60%),
    radial-gradient(ellipse 50% 40% at 75% 60%, rgba(236, 72, 153, 0.12) 0%, transparent 50%);
  background-attachment: fixed;
}

/* Animated floating orbs */
body::before,
body::after {
  content: '';
  position: fixed;
  border-radius: 50%;
  filter: blur(80px);
  z-index: 0;
  pointer-events: none;
  animation: floatOrb 20s ease-in-out infinite;
}

body::before {
  width: 500px;
  height: 500px;
  background: rgba(99, 102, 241, 0.15);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

body::after {
  width: 400px;
  height: 400px;
  background: rgba(168, 85, 247, 0.12);
  bottom: -80px;
  right: -80px;
  animation-delay: -10s;
}

@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(60px, -40px) scale(1.1); }
  50% { transform: translate(-30px, 50px) scale(0.95); }
  75% { transform: translate(40px, 30px) scale(1.05); }
}

#app {
  position: relative;
  z-index: 1;
  min-height: 100vh;
}

/* ═══════════════════════════════════════════════════════════════════════
   Scrollbar
   ═══════════════════════════════════════════════════════════════════════ */

::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.25);
}

/* ═══════════════════════════════════════════════════════════════════════
   Selection
   ═══════════════════════════════════════════════════════════════════════ */

::selection {
  background: rgba(129, 140, 248, 0.4);
  color: #fff;
}
</style>
