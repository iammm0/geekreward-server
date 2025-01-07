package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// MilestoneController 结构体
type MilestoneController struct {
	milestoneService services.MilestoneService
}

func NewMilestoneController(milestoneService services.MilestoneService) *MilestoneController {
	return &MilestoneController{milestoneService: milestoneService}
}

// CreateMilestone 创建新的里程碑
func (ctl *MilestoneController) CreateMilestone(c *gin.Context) {
	var input dtos.MilestoneDTO

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	milestone, err := ctl.milestoneService.CreateMilestone(input, bountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}

// UpdateMilestone 更新里程碑
func (ctl *MilestoneController) UpdateMilestone(c *gin.Context) {
	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	var input dtos.MilestoneUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = ctl.milestoneService.UpdateMilestone(milestoneID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated successfully"})
}

// DeleteMilestone 删除里程碑
func (ctl *MilestoneController) DeleteMilestone(c *gin.Context) {
	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	err = ctl.milestoneService.DeleteMilestone(milestoneID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted successfully"})
}

// GetMilestonesByBountyID 获取指定悬赏令的所有里程碑
func (ctl *MilestoneController) GetMilestonesByBountyID(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	milestones, err := ctl.milestoneService.GetMilestonesByBountyID(bountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

// UpdateMilestoneByReceiver 悬赏零接收者更新里程碑
func (ctl *MilestoneController) UpdateMilestoneByReceiver(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
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

	var input dtos.MilestoneUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = ctl.milestoneService.UpdateMilestoneByReceiver(bountyID, milestoneID, userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated successfully"})
}
