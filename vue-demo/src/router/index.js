import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfileView from '../views/ProfileView.vue'
import HomeView from '../views/HomeView.vue'
import PostDetailView from '../views/PostDetailView.vue'

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
    component: { template: '<div style="padding: 40px; text-align: center;"><h2>发布帖子功能开发中...</h2><p><router-link to="/">返回首页</router-link></p></div>' },
    meta: { requiresAuth: true } // 标记需要认证
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫：检查是否需要登录
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('auth_token')
  
  if (to.meta.requiresAuth && !token) {
    // 需要认证但没有 token，跳转到登录页
    next('/login')
  } else {
    next()
  }
})

export default router

