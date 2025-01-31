package middlewares

import (
	"GeekReward/pkg/utils"
	"github.com/gin-gonic/gin" // 导入Gin框架
	"github.com/google/uuid"
	"net/http" // 导入http包
	"strings"  // 导入strings包
)

// JWTAuthMiddleware 验证JWT的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// 提取Token字符串
		tokenString := parts[1]

		// 使用utils包中的ValidateJWT函数来验证Token
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 确保userID为uuid.UUID类型（如果ValidateJWT返回的是uuid.UUID，无需断言）
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// 将user_id设置到上下文中供后续使用
		c.Set("user_id", userID)

		// 继续处理请求
		c.Next()
	}
}
