package controllers

import (
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type InvitationController struct {
	invitationService   services.InvitationService
	notificationService services.NotificationService
}

func NewInvitationController(
	invitationService services.InvitationService,
	notificationService services.NotificationService,
) *InvitationController {
	return &InvitationController{
		invitationService:   invitationService,
		notificationService: notificationService,
	}
}

// AcceptInvitation 接受组队邀请
func (ctl *InvitationController) AcceptInvitation(c *gin.Context) {
	invitationIDStr := c.Param("invitation_id")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邀请ID"})
		return
	}

	// 获取当前登录用户的ID从上下文
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	if err := ctl.invitationService.AcceptInvitation(invitationID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邀请被接受"})
}

// RejectInvitation 拒绝组队邀请
func (ctl *InvitationController) RejectInvitation(c *gin.Context) {
	invitationIDStr := c.Param("invitation_id")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邀请ID"})
		return
	}

	// 获取当前登录用户的ID从上下文
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	if err := ctl.invitationService.RejectInvitation(invitationID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邀请被拒绝"})
}
