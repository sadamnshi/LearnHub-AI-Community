package router

import (
	"errors"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"gin_demo/controllers"
	"gin_demo/repositories"
	"gin_demo/services"
)

// authMiddleware 简单 JWT 校验中间件：
// 1. 从 Authorization: Bearer <token> 或 Cookie auth_token 读取 token
// 2. 解析与验证签名及过期
// 3. 将关键信息（userID, username）放入上下文
// 企业级增强（未实现但注释说明）：
// - Token 黑名单（登出/强制失效）
// - Redis 存储用户会话信息
// - 角色/权限校验（RBAC/ABAC）
// - 刷新 Token 流程（Refresh Token 双 Token 模式）
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 token: 先尝试 Header, 再尝试 Cookie
		authHeader := c.GetHeader("Authorization")
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			if cookie, err := c.Cookie("auth_token"); err == nil {
				tokenString = cookie
			}
		}
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}
		// 解析 Token（与 services 中使用的同一秘钥）
		parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// 签名算法校验避免被替换为 none
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("CHANGE_ME_TO_ENV_SECRET"), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
			return
		}
		// 读取关键字段
		if sub, ok := claims["sub"].(float64); ok {
			c.Set("userID", int64(sub))
		}
		if name, ok := claims["name"].(string); ok {
			c.Set("username", name)
		}
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 中间件必须在所有路由注册之前 Use，否则对已注册路由不生效
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 依赖注入：repository -> service -> controller
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	postRepo := repositories.NewPostRepository()
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	// JSON API 路由（前后端分离使用）
	api := r.Group("/api")
	{
		// 公开接口（无需登录）
		api.POST("/register", userController.HandleRegisterJSON)
		api.POST("/login", userController.HandleLoginJSON)
		api.GET("/posts", postController.GetPosts) // 帖子列表

		// 受保护路由：需要 JWT
		apiAuth := api.Group("/user").Use(authMiddleware())
		{
			apiAuth.GET("/profile", userController.Profile)
		}
	}

	return r
}
