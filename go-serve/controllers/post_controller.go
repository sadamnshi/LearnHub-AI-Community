package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gin_demo/services"
)

// PostController 帖子控制器
type PostController struct {
	service services.PostService
}

func NewPostController(service services.PostService) *PostController {
	return &PostController{service: service}
}

// GetPosts 获取帖子列表（无需登录，公开接口）
//
// GET /api/posts?page=1&page_size=20&category_id=0
//
// 查询参数：
//   - page: 页码，默认 1
//   - page_size: 每页条数，默认 20，最大 50
//   - category_id: 分类 ID，默认 0（不筛选）
func (ctl *PostController) GetPosts(c *gin.Context) {
	// 解析查询参数，提供默认值
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 64)

	result, err := ctl.service.GetPostList(page, pageSize, uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子列表失败"})
		return
	}

	// 统一响应格式：{ "code": 0, "data": {...}, "msg": "success" }
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": result,
	})
}
