package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type MilestoneController struct {
	milestoneService services.MilestoneService
}

func NewMilestoneController(milestoneService services.MilestoneService) *MilestoneController {
	return &MilestoneController{milestoneService: milestoneService}
}

// GetMilestonesByBountyID 获取指定悬赏令的里程碑
func (ctl *MilestoneController) GetMilestonesByBountyID(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.Atoi(bountyIDStr)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	milestones, err := ctl.milestoneService.GetMilestonesByBountyID(uint(bountyID))
	if err != nil {
		log.Printf("Error fetching milestones: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch milestones"})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

// CreateMilestone 创建新的里程碑
func (ctl *MilestoneController) CreateMilestone(c *gin.Context) {
	var input dtos.MilestoneDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	if err := ctl.milestoneService.CreateMilestone(&input); err != nil {
		log.Printf("Error creating milestone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone created successfully"})
}

// UpdateMilestone 更新里程碑
func (ctl *MilestoneController) UpdateMilestone(c *gin.Context) {
	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := strconv.Atoi(milestoneIDStr)
	if err != nil || milestoneID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Milestone ID"})
		return
	}

	var input dtos.MilestoneDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	if err := ctl.milestoneService.UpdateMilestone(uint(milestoneID), &input); err != nil {
		log.Printf("Error updating milestone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update milestone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated successfully"})
}

// DeleteMilestone 删除里程碑
func (ctl *MilestoneController) DeleteMilestone(c *gin.Context) {
	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := strconv.Atoi(milestoneIDStr)
	if err != nil || milestoneID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Milestone ID"})
		return
	}

	if err := ctl.milestoneService.DeleteMilestone(uint(milestoneID)); err != nil {
		log.Printf("Error deleting milestone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete milestone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted successfully"})
}
