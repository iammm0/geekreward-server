package controllers

import (
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type InvitationController struct {
	invitationService services.InvitationService
}

func NewInvitationController(invitationService services.InvitationService) *InvitationController {
	return &InvitationController{invitationService: invitationService}
}

// AcceptInvitation 接受组队邀请
func (ctl *InvitationController) AcceptInvitation(c *gin.Context) {
	invitationIDStr := c.Param("invitation_id")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := ctl.invitationService.AcceptInvitation(invitationID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted"})
}

// RejectInvitation 拒绝组队邀请
func (ctl *InvitationController) RejectInvitation(c *gin.Context) {
	invitationIDStr := c.Param("invitation_id")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := ctl.invitationService.RejectInvitation(invitationID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation rejected"})
}
