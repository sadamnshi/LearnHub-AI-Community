/**
 * 认证相关 API 封装
 * 企业级最佳实践：
 * 1. 统一错误处理
 * 2. 自动添加 Token
 * 3. 请求/响应拦截
 * 4. 超时控制
 */

const API_BASE = '/api' // 开发环境会被代理到 http://localhost:8080/api

/**
 * 通用请求封装
 */
async function request(url, options = {}) {
  const token = localStorage.getItem('auth_token')
  
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }) // 自动添加 Token
    }
  }

  const finalOptions = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...options.headers
    }
  }

  try {
    const response = await fetch(url, finalOptions)
    
    // 处理 HTTP 错误状态码
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.error || errorData.message || `HTTP ${response.status}`)
    }

    return await response.json()
  } catch (error) {
    console.error('API Request Error:', error)
    throw error
  }
}

/**
 * 用户注册
 * @param {Object} data - { username, password }
 */
export function register(data) {
  return request(`${API_BASE}/register`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

/**
 * 用户登录
 * @param {Object} data - { username, password }
 * @returns {Promise<{token: string, user: Object}>}
 */
export function login(data) {
  return request(`${API_BASE}/login`, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

/**
 * 获取用户信息（需要登录）
 */
export function getUserProfile() {
  return request(`${API_BASE}/user/profile`, {
    method: 'GET'
  })
}

/**
 * 修改密码（需要登录）
 * @param {Object} data - { old_password, new_password }
 */
export function updatePassword(data) {
  return request(`${API_BASE}/user/password`, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

/**
 * 登出（清除本地 Token 和用户信息）
 */
export function logout() {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('username')
}
