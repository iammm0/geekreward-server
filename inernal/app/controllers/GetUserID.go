package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserID 从上下文中获取并断言 user_id 为 uuid.UUID 类型
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("Unauthorized: user_id not found in context")
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("Internal Server Error: user_id is not of type uuid.UUID")
	}

	return userID, nil
}
