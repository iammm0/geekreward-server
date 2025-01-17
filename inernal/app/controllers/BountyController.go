package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var input dtos.BountyDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	bounty, err := ctl.bountyService.CreateBounty(input, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bounty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty created successfully", "bounty": bounty})
}

// GetBounties 获取所有悬赏令处理函数
// 此处控制器的控制逻辑可以通过悬赏令的状态进行延展
// const (
//
//	悬赏令已创建并处于发布状态  BountyStatusCreated             BountyStatus = "Created"
//	悬赏令的里程碑被确认  BountyStatusMilestonesConfirmed BountyStatus = "MilestonesConfirmed"
//	悬赏令里程碑确认  BountyStatusMilestonesVerified  BountyStatus = "MilestonesVerified"
//	悬赏令被接收中  BountyStatusSettling            BountyStatus = "Settling"
//	悬赏令已经被解决  BountyStatusSettled             BountyStatus = "Settled"
//	悬赏令被取消   BountyStatusCancelled           BountyStatus = "Cancelled"
//
// )
// 需要实现通过提取查询参数来获取相应状态的悬赏令
func (ctl *BountyController) GetBounties(c *gin.Context) {
	// bountyIDStr := c.Query("bounty_id")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的限制整数"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的偏移"})
		return
	}

	bounties, err := ctl.bountyService.GetBounties(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取悬赏令失败"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

// GetBounty 获取指定悬赏令的处理函数
func (ctl *BountyController) GetBounty(c *gin.Context) {
	idParam := c.Param("bounty_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	bounty, err := ctl.bountyService.GetBounty(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到悬赏令"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取悬赏令失败"})
		}
		return
	}

	c.JSON(http.StatusOK, bounty)
}

// UpdateBounty 更新悬赏令处理函数
func (ctl *BountyController) UpdateBounty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	_, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	idParam := c.Param("bounty_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	var input dtos.BountyDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的传输模型",
			"details": err.Error(),
		})
		return
	}

	bounty, err := ctl.bountyService.UpdateBounty(id, input)
	if err != nil {
		log.Printf("更新悬赏令时发生错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新悬赏令失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "悬赏令更新成功", "bounty": bounty})
}

// DeleteBounty 删除悬赏令处理函数
func (ctl *BountyController) DeleteBounty(c *gin.Context) {
	idParam := c.Param("bounty_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	err = ctl.bountyService.DeleteBounty(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "悬赏令未找到"})
		} else {
			log.Printf("删除悬赏令时发生错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除悬赏令失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty deleted successfully"})
}

// LikeBounty 点赞悬赏令处理函数
func (ctl *BountyController) LikeBounty(c *gin.Context) {
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

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	if err := ctl.bountyService.LikeBounty(uid, bountyID); err != nil {
		log.Printf("Error liking bounty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like bounty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty liked successfully"})
}

// CommentOnBounty 评论悬赏令处理函数
func (ctl *BountyController) CommentOnBounty(c *gin.Context) {
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

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	comment, err := ctl.bountyService.CommentOnBounty(uid, bountyID, input.Content)
	if err != nil {
		log.Printf("Error commenting on bounty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to comment on bounty"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// RateBounty 评分悬赏令处理函数
func (ctl *BountyController) RateBounty(c *gin.Context) {
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

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	var input struct {
		Score float64 `json:"score" binding:"required,min=1,max=5"` // 假设评分在1到5之间
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	// 调用服务层方法来更新或创建评分记录
	if err := ctl.bountyService.RateBounty(uid, bountyID, input.Score); err != nil {
		log.Printf("Error rating bounty for user %s and bounty %s: %v", uid, bountyID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate bounty"})
		return
	}

	// 返回简单的成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Bounty rated successfully"})
}

// GetBountyComments 获取悬赏令评论处理函数
func (ctl *BountyController) GetBountyComments(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	milestones, err := ctl.milestoneService.GetMilestonesByBountyID(bountyID)
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

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	bounties, err := ctl.bountyService.GetBountiesByUserID(uid)
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

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	bounties, err := ctl.bountyService.GetReceivedBounties(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch received bounties"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

// UnlikeBounty 取消点赞悬赏令处理函数
func (ctl *BountyController) UnlikeBounty(c *gin.Context) {
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

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	if err := ctl.bountyService.UnlikeBounty(uid, bountyID); err != nil {
		// 如果服务层返回具体的错误消息，直接返回该消息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}

// GetUserBountyInteraction 获取用户在指定悬赏令上的互动信息
func (ctl *BountyController) GetUserBountyInteraction(c *gin.Context) {
	// 从上下文中获取当前用户的 user_id
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

	// 获取 bounty_id 路径参数并解析为 uuid.UUID
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bounty ID"})
		return
	}

	// 调用服务层方法获取互动信息
	interaction, err := ctl.bountyService.GetUserBountyInteraction(uid, bountyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Interaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get interaction"})
		}
		return
	}

	c.JSON(http.StatusOK, interaction)
}

// SettleBountyAccounts 结算悬赏令账户
func (ctl *BountyController) SettleBountyAccounts(c *gin.Context) {
	// 获取悬赏令ID从URL参数
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	// 调用服务层进行结算
	err = ctl.bountyService.SettleBountyAccounts(bountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bounty accounts settled successfully"})
}

// ConfirmMilestones 接收者确认提交所有里程碑
func (ctl *BountyController) ConfirmMilestones(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
		return
	}

	// 从上下文中获取当前用户的 user_id
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

	// 调用服务层方法确认里程碑
	err = ctl.bountyService.ConfirmMilestones(bountyID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "当前悬赏令的所有里程碑均被成功确认"})
}

// VerifyMilestones 发布者审核并确认所有里程碑
func (ctl *BountyController) VerifyMilestones(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bounty ID"})
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

	err = ctl.bountyService.VerifyMilestones(bountyID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All milestones verified by publisher"})
}

// ApplySettlement 接收者申请悬赏令清算
func (ctl *BountyController) ApplySettlement(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
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

	err = ctl.bountyService.ApplySettlement(bountyID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settlement request applied successfully"})
}
