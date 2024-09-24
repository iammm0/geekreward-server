package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// BountyController 结构体
type BountyController struct {
	bountyService    services.BountyService
	milestoneService services.MilestoneService
}

// NewBountyController 创建新的 BountyController 实例
func NewBountyController(bountyService services.BountyService, milestoneService services.MilestoneService) *BountyController {
	return &BountyController{bountyService: bountyService, milestoneService: milestoneService}
}

// CreateBounty 创建悬赏令处理函数
func (ctl *BountyController) CreateBounty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input dtos.BountyDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	bounty, err := ctl.bountyService.CreateBounty(input, userID.(uint))
	if err != nil {
		log.Printf("Error creating bounty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bounty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty created successfully", "bounty": bounty})
}

// GetBounties 获取所有悬赏令处理函数
func (ctl *BountyController) GetBounties(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
		return
	}

	bounties, err := ctl.bountyService.GetBounties(limit, offset)
	if err != nil {
		log.Printf("Error fetching bounties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bounties"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

// GetBounty 获取指定悬赏令的处理函数
func (ctl *BountyController) GetBounty(c *gin.Context) {
	idParam := c.Param("bounty_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	bounty, err := ctl.bountyService.GetBounty(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bounty not found"})
		return
	}

	c.JSON(http.StatusOK, bounty)
}

// UpdateBounty 更新悬赏令处理函数
func (ctl *BountyController) UpdateBounty(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("bounty_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input dtos.BountyDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	bounty, err := ctl.bountyService.UpdateBounty(uint(id), input)
	if err != nil {
		log.Printf("Error updating bounty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bounty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty updated successfully", "bounty": bounty})
}

// DeleteBounty 删除悬赏令处理函数
func (ctl *BountyController) DeleteBounty(c *gin.Context) {
	idParam := c.Param("bounty_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = ctl.bountyService.DeleteBounty(uint(id))
	if err != nil {
		log.Printf("Error deleting bounty: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Bounty not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty deleted successfully"})
}

// LikeBounty 点赞悬赏令处理函数
func (ctl *BountyController) LikeBounty(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	if err := ctl.bountyService.LikeBounty(userID, uint(bountyID)); err != nil {
		log.Printf("Error liking bounty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like bounty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty liked successfully"})
}

// CommentOnBounty 评论悬赏令处理函数
func (ctl *BountyController) CommentOnBounty(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	comment, err := ctl.bountyService.CommentOnBounty(userID, uint(bountyID), input.Content)
	if err != nil {
		log.Printf("Error commenting on bounty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to comment on bounty"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// RateBounty 评分悬赏令处理函数
// RateBounty 处理用户评分悬赏令的请求
func (ctl *BountyController) RateBounty(c *gin.Context) {
	userID := c.MustGet("user_id").(uint) // 从上下文中获取用户ID
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	var input struct {
		Score float64 `json:"score" binding:"required,min=1,max=5"` // 假设评分在1到5之间
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// 调用服务层方法来更新或创建评分记录
	if err := ctl.bountyService.RateBounty(userID, uint(bountyID), input.Score); err != nil {
		log.Printf("Error rating bounty for user %d and bounty %d: %v", userID, bountyID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate bounty"})
		return
	}

	// 返回简单的成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Bounty rated successfully"})
}

// GetBountyComments 获取悬赏令评论处理函数
func (ctl *BountyController) GetBountyComments(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	comments, err := ctl.bountyService.GetComments(uint(bountyID))
	if err != nil {
		log.Printf("Error fetching bounty comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bounty comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetUserBountyInteraction 获取用户悬赏令互动信息处理函数
func (ctl *BountyController) GetUserBountyInteraction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	interaction, err := ctl.bountyService.GetUserBountyInteraction(userID, uint(bountyID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果是记录未找到错误，返回404状态码
			c.JSON(http.StatusNotFound, gin.H{"error": "No interaction found for this bounty"})
		} else {
			// 否则返回500内部服务器错误
			log.Printf("Error fetching user bounty interaction for user %d and bounty %d: %v", userID, bountyID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interaction"})
		}
		return
	}

	c.JSON(http.StatusOK, interaction)
}

// GetMilestonesByBountyID 获取指定悬赏令的里程碑
func (ctl *BountyController) GetMilestonesByBountyID(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.Atoi(bountyIDStr)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	milestones, err := ctl.milestoneService.GetMilestonesByBountyID(uint(bountyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch milestones"})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

// GetBountiesByUser 获取指定用户发布的悬赏令
func (ctl *BountyController) GetBountiesByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bounties, err := ctl.bountyService.GetBountiesByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bounties"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

// GetReceivedBounties 获取指定用户接收的悬赏令
func (ctl *BountyController) GetReceivedBounties(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bounties, err := ctl.bountyService.GetReceivedBounties(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch received bounties"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

func (ctl *BountyController) UnlikeBounty(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil || bountyID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	if err := ctl.bountyService.UnlikeBounty(userID, uint(bountyID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}
