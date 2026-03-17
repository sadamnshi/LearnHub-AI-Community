package controllers

import (
	"gin_demo/util"
	"net/http"

	"github.com/gin-gonic/gin"

	"gin_demo/models"
	"gin_demo/services"
)

// UserController 控制层：处理 HTTP 输入/输出，不写业务规则。

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// HandleRegisterJSON 提供 JSON 接口注册，便于前后端分离或移动端调用。
func (ctl *UserController) HandleRegisterJSON(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, util.RespondError("获取用户信息失败"))
		return
	}
	if err := ctl.service.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, util.RespondError("注册失败"))
		return
	}
	util.Success(c, gin.H{"data": "register success"})
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
		c.JSON(http.StatusUnauthorized, util.RespondError("用户名或密码错误"))
		return
	}
	util.Success(c, gin.H{"token": token, "user": u})
}

// Profile 受保护的用户资料接口：从数据库查询完整用户信息返回。
// JWT 中间件已将 userID 注入上下文。
func (ctl *UserController) ShowProfile(c *gin.Context) {
	userID := uint(c.GetInt64("userID")) //从jwt中拿到id

	user, err := ctl.service.GetProfile(userID) //回到服务层查找用户
	if err != nil {
		c.JSON(http.StatusNotFound, util.RespondError("用户未找到")) //如果没有找到用户，返回404
		return
	}

	util.Success(c, gin.H{"data": user})
}

func (ctl *UserController) UpdatePassword(c *gin.Context) {
	var payload struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, util.RespondError("获取新旧密码失败"))
		return
	}

	// Token 已由 authMiddleware 验证通过，并将 userID 存入了上下文
	// 这里直接从上下文取出，不需要再手动解析 Token
	userID := uint(c.GetInt64("userID"))

	// 调用 Service 处理：验证旧密码、检查新旧密码不同、哈希新密码、写库
	if err := ctl.service.UpdatePassword(userID, payload.OldPassword, payload.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, util.RespondError("密码修改失败"))
		return
	}

	util.Success(c, gin.H{"data": "密码修改成功"})
}
