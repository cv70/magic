package infra

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"magic/utils"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 读取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "missing authorization header",
			})
			c.Abort()
			return
		}

		// 提取 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 将 userID 注入 context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
