package controllers

import (
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ApplicationController struct {
	applicationService services.ApplicationService
}

// NewApplicationController 创建新的 ApplicationController 实例
func NewApplicationController(applicationService services.ApplicationService) *ApplicationController {
	return &ApplicationController{applicationService: applicationService}
}

// CreateApplication 创建新的悬赏任务申请
func (ctl *ApplicationController) CreateApplication(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		BountyID uint `json:"bounty_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	if err := ctl.applicationService.CreateApplication(input.BountyID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application created successfully"})
}

// GetApplications 获取某个悬赏任务的所有申请
func (ctl *ApplicationController) GetApplications(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	applications, err := ctl.applicationService.GetApplications(uint(bountyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// ApproveApplication 批准申请
func (ctl *ApplicationController) ApproveApplication(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Application ID"})
		return
	}

	if err := ctl.applicationService.ApproveApplication(uint(applicationID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application approved successfully"})
}

// RejectApplication 拒绝申请
func (ctl *ApplicationController) RejectApplication(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Application ID"})
		return
	}

	if err := ctl.applicationService.RejectApplication(uint(applicationID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application rejected successfully"})
}
