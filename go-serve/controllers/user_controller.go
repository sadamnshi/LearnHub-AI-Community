package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin_demo/models"
	"gin_demo/services"
)

// UserController 控制层：处理 HTTP 输入/输出，不写业务规则。
// 主要职责：
// 1. 参数绑定与校验返回错误（ShouldBind）
// 2. 调用 service 完成业务动作
// 3. 格式化响应：HTML / JSON
// 4. 不直接操作数据库（保持分层）
// 5. 不泄漏内部错误细节（返回业务语义）
type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// ShowLoginPage 展示登录页面
func (ctl *UserController) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// ShowRegisterPage 展示注册页面（示例，如果需要单独页面）
func (ctl *UserController) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// HandleRegister 处理表单注册示例（HTML 表单提交）。
// 企业级注意：
// 1. 限制注册频率（防刷）——使用滑块验证码/限流策略（未演示）。
// 2. 邮箱/手机验证——发送验证码后再激活（未演示）。
// 3. 密码强度校验——正则 + 长度 + 黑名单（示例简化）。
func (ctl *UserController) HandleRegister(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	if err := ctl.service.Register(&user); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{"msg": "注册成功，请登录"})
}

// HandleRegisterJSON 提供 JSON 接口注册，便于前后端分离或移动端调用。
func (ctl *UserController) HandleRegisterJSON(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctl.service.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}

// HandleLogin 处理表单登录并签发 JWT（写入 Cookie）。
// 安全注意：
// 1. Cookie 设置 HttpOnly 避免被 JS 读取。
// 2. Secure 应该在 HTTPS 场景开启（演示环境可能是 http）。
// 3. SameSite=Lax/Strict 减少 CSRF 风险。
// 4. 也可以签发在 Authorization Header 中（JSON 接口）。
func (ctl *UserController) HandleLogin(c *gin.Context) {
	var form models.User
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	token, u, err := ctl.service.Authenticate(form.Username, form.Password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": err.Error()})
		return
	}
	// 写入 JWT 到 Cookie（或返回给前端由前端存储）。
	c.SetCookie("auth_token", token, 7200, "/", "", false, true) // 2h 过期
	c.HTML(http.StatusOK, "login.html", gin.H{"msg": "登录成功", "username": u.Username})
}

// HandleLoginJSON JSON 登录接口：返回 JWT 在响应体中。
func (ctl *UserController) HandleLoginJSON(c *gin.Context) {
	var payload struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, u, err := ctl.service.Authenticate(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": u.ID, "username": u.Username}})
}

// Profile 受保护的用户资料接口：从数据库查询完整用户信息返回。
// JWT 中间件已将 userID 注入上下文。
func (ctl *UserController) Profile(c *gin.Context) {
	userID := uint(c.GetInt64("userID"))

	user, err := ctl.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"bio":      user.Bio,
		"role":     user.Role,
	})
}
