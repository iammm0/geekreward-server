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
	milestoneService    services.MilestoneService
	notificationService services.NotificationService
}

func NewMilestoneController(
	milestoneService services.MilestoneService,
	notificationService services.NotificationService,
) *MilestoneController {
	return &MilestoneController{
		milestoneService:    milestoneService,
		notificationService: notificationService,
	}
}

// CreateMilestone 创建新的里程碑
func (ctl *MilestoneController) CreateMilestone(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	var input dtos.MilestoneDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的传输模型"})
		return
	}

	milestone, err := ctl.milestoneService.CreateMilestone(bountyID, input)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的里程碑ID"})
		return
	}

	var input dtos.MilestoneUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的传输模型"})
		return
	}

	err = ctl.milestoneService.UpdateMilestone(milestoneID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "里程碑更新成功"})
}

// DeleteMilestone 删除里程碑
func (ctl *MilestoneController) DeleteMilestone(c *gin.Context) {
	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的里程碑ID"})
		return
	}

	err = ctl.milestoneService.DeleteMilestone(milestoneID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "里程碑删除成功"})
}

// GetMilestonesByBountyID 获取指定悬赏令的所有里程碑
func (ctl *MilestoneController) GetMilestonesByBountyID(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	// 从Query或Param获取bounty_id
	// bountyIDStr := c.Query("bounty_id") // 或者使用路由: c.Param("bounty_id")

	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	milestoneIDStr := c.Param("milestone_id")
	milestoneID, err := uuid.Parse(milestoneIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的里程碑ID"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	var input dtos.MilestoneUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的传输模型"})
		return
	}

	err = ctl.milestoneService.UpdateMilestoneByReceiver(bountyID, milestoneID, userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "里程碑更新成功"})
}
