/**
 * AI 聊天 API 模块
 * 负责与后端 AI 聊天接口的通信
 */

/**
 * 通用 HTTP 请求函数
 * @param {string} method - HTTP 方法 (GET, POST, PUT, DELETE 等)
 * @param {string} url - API 路径
 * @param {Object} data - 请求数据（JSON 格式）
 * @returns {Promise} 响应数据
 */
async function request(method, url, data = null) {
  const token = localStorage.getItem('auth_token')

  const options = {
    method: method,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    },
  }

  // 只有 POST、PUT、PATCH 需要 body
  if (data && ['POST', 'PUT', 'PATCH'].includes(method)) {
    options.body = JSON.stringify(data)
  }

  try {
    const response = await fetch(url, options)

    // 尝试解析 JSON 响应
    const responseData = await response.json().catch(() => null)

    if (!response.ok) {
      // 提取错误信息
      const errorMsg = responseData?.msg || responseData?.error || responseData?.message || `HTTP ${response.status}`
      throw new Error(errorMsg)
    }

    return responseData
  } catch (error) {
    console.error(`❌ API request failed: ${method} ${url}`, error)
    throw error
  }
}

/**
 * 发送消息给 AI
 * @param {string} message - 用户消息内容
 * @returns {Promise} AI 回复
 */
export async function sendMessage(message) {
  if (!message || message.trim().length === 0) {
    throw new Error('消息不能为空')
  }

  if (message.length > 2000) {
    throw new Error('消息长度不能超过 2000 个字符')
  }

  return request('POST', '/api/chat/send', {
    message: message.trim(),
  })
}

/**
 * 获取聊天历史记录
 * @param {number} limit - 返回的最多消息数（默认 20，最多 100）
 * @returns {Promise} 聊天历史数组
 */
export async function getChatHistory(limit = 20) {
  const validLimit = Math.min(Math.max(1, limit), 100)
  return request('GET', `/api/chat/history?limit=${validLimit}`)
}

/**
 * 清空聊天历史
 * @returns {Promise}
 */
export async function clearChatHistory() {
  return request('DELETE', '/api/chat/history')
}

// 导出所有聊天相关的 API
export default {
  sendMessage,
  getChatHistory,
  clearChatHistory,
}
