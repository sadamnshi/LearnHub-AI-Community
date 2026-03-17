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
