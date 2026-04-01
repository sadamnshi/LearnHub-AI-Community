package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func RespondError(message string) gin.H {
	return gin.H{"error": message}
}

// ErrorResponse 返回统一的错误响应格式
func ErrorResponse(c *gin.Context, statusCode int, errorCode string, message string) {
	c.JSON(statusCode, gin.H{
		"code":  -1,
		"error": errorCode,
		"msg":   message,
	})
}
