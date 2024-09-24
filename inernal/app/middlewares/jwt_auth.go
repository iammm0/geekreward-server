package middlewares

import (
	"GeekReward/pkg/utils"
	"github.com/gin-gonic/gin" // 导入Gin框架
	"net/http"                 // 导入http包
	"strings"                  // 导入strings包
)

// JWTAuthMiddleware 验证JWT的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token
		tokenString := c.GetHeader("Authorization")

		// 检查 token 的格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}

		// 提取Token字符串
		tokenString = tokenString[len("Bearer "):]

		// 使用utils包中的ValidateJWT函数来验证Token
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户ID存入上下文供后续使用
		c.Set("user_id", userID)

		// 继续处理请求
		c.Next()
	}
}
