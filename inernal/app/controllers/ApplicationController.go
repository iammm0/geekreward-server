package controllers

import (
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// ApplicationController 结构体
type ApplicationController struct {
	applicationService services.ApplicationService
}

// NewApplicationController 创建新的 ApplicationController 实例
func NewApplicationController(applicationService services.ApplicationService) *ApplicationController {
	return &ApplicationController{applicationService: applicationService}
}

// CreateApplication 创建新的悬赏任务申请
func (ctl *ApplicationController) CreateApplication(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var input struct {
		BountyID uuid.UUID `json:"bounty_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	if err := ctl.applicationService.CreateApplication(input.BountyID, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application created successfully"})
}

// GetApplications 获取某个悬赏任务的所有申请
func (ctl *ApplicationController) GetApplications(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	applications, err := ctl.applicationService.GetApplications(bountyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// ApproveApplication 批准申请
func (ctl *ApplicationController) ApproveApplication(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Application ID"})
		return
	}

	if err := ctl.applicationService.ApproveApplication(applicationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application approved successfully"})
}

// RejectApplication 拒绝申请
func (ctl *ApplicationController) RejectApplication(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Application ID"})
		return
	}

	if err := ctl.applicationService.RejectApplication(applicationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application rejected successfully"})
}

// GetPublicApplications 获取公开的申请信息
func (ctl *ApplicationController) GetPublicApplications(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	applications, err := ctl.applicationService.GetPublicApplications(bountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applications)
}
