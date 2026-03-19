package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// authMiddleware 简单 JWT 校验中间件：
// 1. 从 Authorizatio	n: Bearer <token> 或 Cookie auth_token 读取 token
// 2. 解析与验证签名及过期
// 3. 将关键信息（userID, username）放入上下文
// 企业级增强（未实现但注释说明）：
// - Token 黑名单（登出/强制失效）
// - Redis 存储用户会话信息
// - 角色/权限校验（RBAC/ABAC）
// - 刷新 Token 流程（Refresh Token 双 Token 模式）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ── 第一步：从请求头取出 Token 字符串 ──────────────────────────────
		//
		// 前端登录成功后，服务端会返回一个 Token（长得像 "eyJhbGci..."）
		// 前端每次请求受保护接口时，需要在请求头里带上：
		//   Authorization: Bearer eyJhbGci...
		//
		// 这里取出整个 Authorization 头的值，例如 "Bearer eyJhbGci..."
		authHeader := c.GetHeader("Authorization")

		// 检查 Authorization 头是否为空
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header required"})
			return
		}

		// 检查是否以 "Bearer " 开头
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header format"})
			return
		}

		// 去掉前面的 "Bearer " 前缀（7个字符），只保留纯 Token 字符串
		tokenString := authHeader[7:]

		// ── 第二步：用密钥解析并验证 Token ────────────────────────────────
		//
		// Token 由三部分组成，用 "." 分隔：头部.载荷.签名
		// 例如：eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOjEsIm5hbWUiOiJhZG1pbiJ9.xxxxxx
		//   头部(header)  ：声明算法类型，如 HS256
		//   载荷(payload) ：存放用户信息，如 userID、username、过期时间
		//   签名(signature)：用服务器密钥对前两部分签名，防止被篡改
		//
		// jwt.Parse 会同时做两件事：
		//   1. 解码载荷，取出里面的用户信息
		//   2. 用密钥重新计算签名，与 Token 里的签名比对，不一致则说明 Token 被篡改
		parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// 返回密钥，必须与 user_service.go 里 GenerateToken 用的密钥完全一致
			// 只有密钥相同，签名验证才能通过
			return []byte("CHANGE_ME_TO_ENV_SECRET"), nil
		})

		// 解析失败（Token 格式错误、签名不匹配、已过期等），拒绝请求
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// ── 第三步：从 Token 载荷中取出用户信息 ──────────────────────────
		//
		// 登录时 user_service.go 里的 GenerateToken 把以下字段写进了 Token：
		//   "sub"  → 用户ID（数字）
		//   "name" → 用户名（字符串）
		//   "exp"  → 过期时间
		//
		// 现在把它们读出来，存入本次请求的上下文，供后续 Controller 直接使用
		claims := parsed.Claims.(jwt.MapClaims)

		// "sub" 存的是用户ID，JWT 里数字统一用 float64 存储，需转成 int64
		userID := int64(claims["sub"].(float64))
		username := claims["name"].(string)

		// 把用户信息存入上下文，Controller 里可以用 c.GetInt64("userID") 取出
		c.Set("userID", userID)
		c.Set("username", username)

		// ── 第四步：验证通过，放行请求 ────────────────────────────────────
		c.Next()
	}
}
