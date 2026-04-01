import { createRouter, createWebHistory } from 'vue-router'
import { isTokenValid, isTokenExpired, clearExpiredToken } from '@/utils/token'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfileView from '../views/ProfileView.vue'
import HomeView from '../views/HomeView.vue'
import PostDetailView from '../views/PostDetailView.vue'
import CreatePostView from '../views/CreatePostView.vue'
import AIChatView from '../views/AIChatView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView  // 首页（帖子列表，无需登录）
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView
  },
  {
    path: '/profile',
    name: 'profile',
    component: ProfileView,
    meta: { requiresAuth: true } // 标记需要认证
  },
  {
    path: '/posts/:id',
    name: 'post-detail',
    component: PostDetailView,
    meta: { 
      title: '帖子详情',
      requiresAuth: false // 不需要登录，公开访问
    }
  },
  {
    path: '/post/create',
    name: 'post-create',
    component: CreatePostView,
    meta: { requiresAuth: true } // 标记需要认证
  },
  {
    path: '/ai-chat',
    name: 'ai-chat',
    component: AIChatView,
    meta: { requiresAuth: true, title: 'AI 聊天' } // 标记需要认证
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫：检查是否需要登录
router.beforeEach((to, from, next) => {
  // 首先检查并清除过期的 Token
  const token = localStorage.getItem('auth_token')
  if (token && isTokenExpired(token)) {
    console.warn('⏰ 检测到 Token 已过期，自动清除')
    localStorage.removeItem('auth_token')
    localStorage.removeItem('username')
  }

  // 需要认证的路由
  if (to.meta.requiresAuth && !isTokenValid(localStorage.getItem('auth_token'))) {
    // 需要认证但没有有效 token，跳转到登录页
    console.log('🔒 需要登录，重定向到登录页')
    next('/login')
  } else {
    next()
  }
})

export default router

