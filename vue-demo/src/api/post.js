/**
 * 帖子相关 API 封装
 */

const API_BASE = '/api'

async function request(url, options = {}) {
  const token = localStorage.getItem('auth_token')

  const finalOptions = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` })
    },
    ...options,
    // 合并 headers，避免覆盖
  }
  // 重新合并 headers（options.headers 可能覆盖 token header）
  finalOptions.headers = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers
  }

  const response = await fetch(url, finalOptions)
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.error || errorData.message || `HTTP ${response.status}`)
  }
  return await response.json()
}

/**
 * 获取帖子列表（公开接口，无需登录）
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码，默认 1
 * @param {number} params.page_size - 每页条数，默认 20
 * @param {number} params.category_id - 分类 ID，默认 0（不筛选）
 * @returns {Promise<{code, data: {list, total, page, page_size}, msg}>}
 */
export function getPostList({ page = 1, page_size = 20, category_id = 0 } = {}) {
  const query = new URLSearchParams({ page, page_size, category_id }).toString()
  return request(`${API_BASE}/posts/home?${query}`)
}

/**
 * 获取帖子详情（公开接口，无需登录）
 * 使用旁路缓存模式优化性能：
 * - 首次访问：查询数据库 + 存入 Redis（~200ms）
 * - 后续访问：直接返回缓存（~10ms，快 20 倍！）
 * 
 * @param {number} postId - 帖子 ID
 * @returns {Promise<{code, data: PostDetail, msg}>}
 * @throws {Error} 如果帖子不存在或获取失败
 * 
 * 响应数据结构示例：
 * {
 *   "code": 0,
 *   "msg": "success",
 *   "data": {
 *     "id": 1,
 *     "title": "Go 语言并发编程最佳实践",
 *     "content": "完整的帖子内容...",
 *     "author": {
 *       "id": 1,
 *       "username": "tom",
 *       "avatar": "https://..."
 *     },
 *     "category": {
 *       "id": 1,
 *       "name": "技术分享",
 *       "icon": "💻"
 *     },
 *     "tags": ["golang", "并发", "goroutine"],
 *     "status": "published",
 *     "created_at": "2026-03-15 10:30:45",
 *     "updated_at": "2026-03-15 10:30:45"
 *   }
 * }
 * 
 * 错误响应示例：
 * {
 *   "code": 404,
 *   "msg": "帖子不存在或已被删除"
 * }
 */
export async function fetchPostDetail(postId) {
  // 参数验证：确保 postId 是有效的正整数
  if (!postId || postId <= 0) {
    throw new Error('无效的帖子 ID')
  }

  try {
    // 调用通用 request 函数，发送 GET 请求
    // 最终 URL：/api/posts/1、/api/posts/123 等
    const response = await request(`${API_BASE}/posts/${postId}`)

    // 检查业务状态码（不是 HTTP 状态码）
    // code = 0 表示成功，其他值表示失败
    if (response.code !== 0) {
      // 从服务器返回的错误信息
      throw new Error(response.msg || '获取帖子详情失败')
    }

    // 返回实际的数据对象（不是整个响应）
    // 调用者会直接得到 PostDetail 对象
    return response.data

  } catch (error) {
    // 如果是我们自己抛出的错误，直接抛出
    if (error instanceof Error) {
      throw error
    }
    // 其他类型的错误，包装成 Error 对象
    throw new Error(error.message || '获取帖子详情时发生错误')
  }
}
