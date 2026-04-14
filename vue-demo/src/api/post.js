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
    // 后端返回的是 gin.H 格式，包含 error/msg 字段
    const errorMsg = errorData.error || errorData.msg || errorData.message || `HTTP ${response.status}`
    throw new Error(errorMsg)
  }
  
  // 检查响应是否是 JSON 格式
  const data = await response.json().catch(() => null)
  if (!data) {
    throw new Error('服务器响应无效')
  }
  
  return data
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

/**
 * 创建新帖子（需要登录，受保护接口）
 * 
 * @param {Object} postData - 帖子数据
 * @param {string} postData.title - 帖子标题（必填，最多200字符）
 * @param {string} postData.content - 帖子内容（必填，支持Markdown，至少10字符）
 * @param {number} postData.category_id - 分类ID（可选，默认0）
 * @param {string} postData.tags - 逗号分隔的标签（可选，最多200字符）
 * @returns {Promise<PostDetail>} 创建后的帖子详情
 * @throws {Error} 如果创建失败或未登录
 * 
 * 调用示例：
 * const newPost = await createPost({
 *   title: '我的第一篇帖子',
 *   content: '这是帖子的完整内容...',
 *   category_id: 1,
 *   tags: 'javascript,vue,前端'
 * })
 * 
 * 成功响应示例：
 * {
 *   "id": 123,
 *   "title": "我的第一篇帖子",
 *   "content": "这是帖子的完整内容...",
 *   "author": {
 *     "id": 1,
 *     "username": "alice",
 *     "avatar": "https://..."
 *   },
 *   "category": {
 *     "id": 1,
 *     "name": "技术分享",
 *     "icon": "💻"
 *   },
 *   "tags": ["javascript", "vue", "前端"],
 *   "status": "published",
 *   "created_at": "2026-03-19 15:30:45",
 *   "updated_at": "2026-03-19 15:30:45"
 * }
 * 
 * 错误响应示例：
 * {
 *   "code": 401,
 *   "msg": "未授权，请先登录"
 * }
 */
export async function createPost(postData) {
  // 参数验证
  if (!postData.title || !postData.title.trim()) {
    throw new Error('帖子标题不能为空')
  }

  if (!postData.content || !postData.content.trim()) {
    throw new Error('帖子内容不能为空')
  }

  if (postData.content.trim().length < 10) {
    throw new Error('帖子内容至少需要10个字符')
  }

  const token = localStorage.getItem('auth_token')
  if (!token) {
    throw new Error('未登录，请先登录后再发布帖子')
  }

  // 调试信息：检查 Token 格式
  if (!token.startsWith('eyJ')) {
    console.warn('⚠️ Token 格式可能不正确，预期以 "eyJ" 开头')
  }
  if (token.length < 50) {
    console.warn('⚠️ Token 长度过短，可能不是有效的 JWT')
  }

  try {
    // 发送 POST 请求到创建帖子接口
    const response = await request(`${API_BASE}/posts/create`, {
      method: 'POST',
      body: JSON.stringify({
        title: postData.title.trim(),
        content: postData.content.trim(),
        category_id: postData.category_id || 0,
        tags: (postData.tags || '').trim()
      })
    })

    // 检查业务状态码
    if (response.code !== 0) {
      throw new Error(response.msg || '创建帖子失败')
    }

    // 返回创建后的帖子详情
    return response.data

  } catch (error) {
    // 更详细的错误日志
    console.error('❌ 创建帖子错误详情:', {
      message: error.message,
      stack: error.stack,
      token: token ? `${token.substring(0, 50)}...` : 'null'
    })
    
    if (error instanceof Error) {
      throw error
    }
    throw new Error(error.message || '创建帖子时发生错误')
  }
}

/**
 * 获取所有帖子分类（公开接口，无需登录）
 * 
 * @returns {Promise<CategoryInfo[]>} 分类列表
 * @throws {Error} 如果获取失败
 * 
 * 响应数据结构示例：
 * [
 *   {
 *     "id": 1,
 *     "name": "技术分享",
 *     "icon": "💻"
 *   },
 *   {
 *     "id": 2,
 *     "name": "生活杂谈",
 *     "icon": "😊"
 *   }
 * ]
 */
export async function getCategories() {
  try {
    // 调用通用 request 函数，发送 GET 请求
    const response = await request(`${API_BASE}/posts/categories`)

    // 检查业务状态码
    if (response.code !== 0) {
      throw new Error(response.msg || '获取分类列表失败')
    }

    // 返回分类列表（data 直接是数组）
    return response.data || []

  } catch (error) {
    console.error('❌ 获取分类列表错误:', error.message)
    
    if (error instanceof Error) {
      throw error
    }
    throw new Error(error.message || '获取分类列表时发生错误')
  }
}

