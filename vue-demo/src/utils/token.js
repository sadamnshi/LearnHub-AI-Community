/**
 * Token 工具函数
 * 用于处理 JWT Token 的验证、过期检查等操作
 */

/**
 * 检查 Token 是否过期
 * @param {string} token - JWT Token 字符串
 * @returns {boolean} true = 已过期，false = 未过期
 * 
 * @example
 * const token = localStorage.getItem('auth_token')
 * if (isTokenExpired(token)) {
 *   console.log('Token 已过期，需要重新登录')
 * }
 */
export function isTokenExpired(token) {
  if (!token) {
    return true  // 无 Token，认为已过期
  }

  try {
    // JWT 由三部分组成，用 "." 分隔
    // 格式：header.payload.signature
    const parts = token.split('.')
    
    // 验证 Token 格式
    if (parts.length !== 3) {
      console.warn('Token 格式错误：应包含 3 部分')
      return true
    }

    // 解码 payload（第二部分）
    // payload 是 Base64 编码，需要 atob 解码
    const payload = JSON.parse(atob(parts[1]))

    // 检查是否包含 exp 字段
    if (!payload.exp) {
      console.warn('Token 中缺少 exp 字段')
      return true
    }

    // exp 是过期时间，单位为秒（Unix 时间戳）
    // 需要转换为毫秒来与 Date.now() 比较
    const expiresAt = payload.exp * 1000
    const now = Date.now()

    // 如果当前时间 >= 过期时间，则已过期
    const expired = now >= expiresAt
    
    return expired
  } catch (error) {
    console.error('Token 解析错误:', error)
    // 解析失败，认为已过期
    return true
  }
}

/**
 * 检查 Token 是否有效（存在且未过期）
 * @param {string} token - JWT Token 字符串
 * @returns {boolean} true = 有效，false = 无效或已过期
 * 
 * @example
 * if (!isTokenValid(token)) {
 *   this.$router.push('/login')
 * }
 */
export function isTokenValid(token) {
  return !!token && !isTokenExpired(token)
}

/**
 * 获取 Token 剩余时间（秒）
 * @param {string} token - JWT Token 字符串
 * @returns {number} 剩余秒数（0 = 已过期）
 * 
 * @example
 * const remaining = getTokenTimeRemaining(token)
 * console.log(`Token 剩余时间：${remaining} 秒`)
 */
export function getTokenTimeRemaining(token) {
  if (!token) {
    return 0
  }

  try {
    const parts = token.split('.')
    if (parts.length !== 3) return 0

    const payload = JSON.parse(atob(parts[1]))
    if (!payload.exp) return 0

    const expiresAt = payload.exp * 1000
    const remaining = Math.floor((expiresAt - Date.now()) / 1000)
    
    // 返回正数或 0
    return Math.max(0, remaining)
  } catch (error) {
    console.error('获取 Token 剩余时间失败:', error)
    return 0
  }
}

/**
 * 获取 Token 过期时间
 * @param {string} token - JWT Token 字符串
 * @returns {Date|null} 过期时间的 Date 对象，解析失败返回 null
 * 
 * @example
 * const expireTime = getTokenExpireTime(token)
 * console.log(`Token 过期时间：${expireTime}`)
 */
export function getTokenExpireTime(token) {
  if (!token) {
    return null
  }

  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    const payload = JSON.parse(atob(parts[1]))
    if (!payload.exp) return null

    // exp 单位为秒，转换为毫秒后创建 Date 对象
    return new Date(payload.exp * 1000)
  } catch (error) {
    console.error('获取 Token 过期时间失败:', error)
    return null
  }
}

/**
 * 格式化剩余时间
 * @param {number} seconds - 剩余秒数
 * @returns {string} 格式化后的字符串，如 "1小时30分钟"
 * 
 * @example
 * const remaining = getTokenTimeRemaining(token)
 * console.log(formatTimeRemaining(remaining))
 * // 输出："1小时30分钟"
 */
export function formatTimeRemaining(seconds) {
  if (seconds <= 0) {
    return '已过期'
  }

  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  const parts = []
  if (hours > 0) parts.push(`${hours}小时`)
  if (minutes > 0) parts.push(`${minutes}分钟`)
  if (secs > 0 || parts.length === 0) parts.push(`${secs}秒`)

  return parts.join(' ')
}

/**
 * 解析 Token 获取用户信息
 * @param {string} token - JWT Token 字符串
 * @returns {Object|null} Token 中的 payload 对象
 * 
 * @example
 * const payload = decodeToken(token)
 * console.log(payload.sub)      // 用户 ID
 * console.log(payload.name)     // 用户名
 * console.log(payload.exp)      // 过期时间
 */
export function decodeToken(token) {
  if (!token) {
    return null
  }

  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    const payload = JSON.parse(atob(parts[1]))
    return payload
  } catch (error) {
    console.error('Token 解析失败:', error)
    return null
  }
}

/**
 * 获取 Token 中的用户 ID
 * @param {string} token - JWT Token 字符串
 * @returns {number|null} 用户 ID
 */
export function getUserIdFromToken(token) {
  const payload = decodeToken(token)
  return payload?.sub || null
}

/**
 * 获取 Token 中的用户名
 * @param {string} token - JWT Token 字符串
 * @returns {string|null} 用户名
 */
export function getUsernameFromToken(token) {
  const payload = decodeToken(token)
  return payload?.name || null
}

/**
 * 清除过期的 Token
 * 在 App 初始化时调用此函数
 * 
 * @example
 * clearExpiredToken()  // 如果 Token 已过期则清除
 */
export function clearExpiredToken() {
  const token = localStorage.getItem('auth_token')
  
  if (token && isTokenExpired(token)) {
    console.warn('检测到 Token 已过期，自动清除')
    localStorage.removeItem('auth_token')
    localStorage.removeItem('username')
    return true  // 已清除
  }
  
  return false  // 未清除
}

/**
 * 检查 Token 是否即将过期（在指定时间内）
 * @param {string} token - JWT Token 字符串
 * @param {number} minutesBefore - 多少分钟后过期则认为即将过期（默认 30 分钟）
 * @returns {boolean} true = 即将过期，false = 还有足够时间
 * 
 * @example
 * if (isTokenExpiringSoon(token, 15)) {
 *   console.log('Token 将在 15 分钟后过期，建议续约')
 * }
 */
export function isTokenExpiringSoon(token, minutesBefore = 30) {
  if (!token || isTokenExpired(token)) {
    return false
  }

  const remaining = getTokenTimeRemaining(token)
  const secondsBefore = minutesBefore * 60

  return remaining <= secondsBefore && remaining > 0
}
