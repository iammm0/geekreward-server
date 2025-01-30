package controllers

import (
	dtos2 "GeekReward/main/inernal/app/models/dtos"
	"GeekReward/main/inernal/app/models/tables"
	services2 "GeekReward/main/inernal/app/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// BountyController 结构体
type BountyController struct {
	bountyService       services2.BountyService
	milestoneService    services2.MilestoneService
	notificationService services2.NotificationService
}

// NewBountyController 创建新的 BountyController 实例
func NewBountyController(
	bountyService services2.BountyService,
	milestoneService services2.MilestoneService,
	notificationService services2.NotificationService,
) *BountyController {
	return &BountyController{
		bountyService:       bountyService,
		milestoneService:    milestoneService,
		notificationService: notificationService,
	}
}

// CreateBounty 创建悬赏令处理函数
func (ctl *BountyController) CreateBounty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的ID类型"})
		return
	}

	var input dtos2.BountyDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	bounty, err := ctl.bountyService.CreateBounty(input, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建悬赏令失败"})
		return
	}

	// 发送通知给发布者（自己）
	err = ctl.notificationService.CreateNotification(&tables.Notification{
		UserID:      bounty.UserID,
		Type:        "BountyCreated",
		Title:       "悬赏令已创建",
		Description: fmt.Sprintf("您已成功创建悬赏令 '%s'。", bounty.Title),
		Metadata: map[string]interface{}{
			"bounty_id":    bounty.ID.String(),
			"bounty_title": bounty.Title,
		},
		IsRead: false,
	})
	if err != nil {
		// 记录错误但不影响主流程
		fmt.Println("发送通知失败:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "悬赏令创建成功", "bounty": bounty})
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
// GetBounties GET /bounties?status=Settled&publisher_id=...&receiver_id=...&limit=10&offset=0
func (ctl *BountyController) GetBounties(c *gin.Context) {
	// 1. 解析分页
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分页限制（limit）"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分页限制（offset）"})
		return
	}

	// 2. 解析可选status
	statusStr := c.Query("status")
	var status *tables.BountyStatus
	if statusStr != "" {
		s := tables.BountyStatus(statusStr)
		status = &s
	}

	// 3. 解析可选 publisher_id
	publisherIDStr := c.Query("publisher_id")
	var publisherID *uuid.UUID
	if publisherIDStr != "" {
		pubID, err := uuid.Parse(publisherIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publisher_id"})
			return
		}
		publisherID = &pubID
	}

	// 4. 解析可选 receiver_id
	receiverIDStr := c.Query("receiver_id")
	var receiverID *uuid.UUID
	if receiverIDStr != "" {
		recID, err := uuid.Parse(receiverIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiver_id"})
			return
		}
		receiverID = &recID
	}

	// 5. 构造过滤器
	filters := dtos2.BountyFilter{
		Status:      status,
		PublisherID: publisherID,
		ReceiverID:  receiverID,
		Limit:       limit,
		Offset:      offset,
	}

	// 6. 调用服务层
	bounties, err := ctl.bountyService.FindBounties(filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	var input dtos2.BountyDTO
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

	c.JSON(http.StatusOK, gin.H{"message": "悬赏令删除成功"})
}

// LikeBounty 点赞悬赏令处理函数
func (ctl *BountyController) LikeBounty(c *gin.Context) {
	currentUserIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	currentUserID, ok := currentUserIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的UserID类型"})
		return
	}

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	if err := ctl.bountyService.LikeBounty(currentUserID, bountyID); err != nil {
		log.Printf("点赞悬赏令错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞悬赏令失败"})
		return
	}

	// 2. 获取悬赏令信息, 找到发布者
	bounty, err := ctl.bountyService.GetBounty(bountyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询悬赏令失败"})
		return
	}

	// 3. 通过 NotificationService 的便捷方法创建通知
	err = ctl.notificationService.CreateBountyLikeNotification(
		currentUserID, // actor
		bounty.UserID, // publisher
		bounty.ID,     // bountyID
		bounty.Title,  // bountyTitle
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "悬赏令喜欢成功"})
}

// CommentOnBounty 评论悬赏令处理函数
func (ctl *BountyController) CommentOnBounty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的传输模型",
			"details": err.Error(),
		})
		return
	}

	comment, err := ctl.bountyService.PostComment(uid, bountyID, input.Content)
	if err != nil {
		log.Printf("在悬赏令下评论出错: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "在悬赏令下评论失败"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// RateBounty 评分悬赏令处理函数
func (ctl *BountyController) RateBounty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	var input struct {
		Score float64 `json:"score" binding:"required,min=1,max=5"` // 假设评分在1到5之间
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的传输模型",
			"details": err.Error(),
		})
		return
	}

	// 调用服务层方法来更新或创建评分记录
	if err := ctl.bountyService.RateBounty(uid, bountyID, input.Score); err != nil {
		log.Printf("Error rating bounty for user %s and bounty %s: %v", uid, bountyID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评分悬赏令失败"})
		return
	}

	// 返回简单的成功响应
	c.JSON(http.StatusOK, gin.H{"message": "评分悬赏令成功"})
}

// GetBountyComments 获取悬赏令评论处理函数
func (ctl *BountyController) GetBountyComments(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	comments, err := ctl.bountyService.GetCommentsByBountyID(bountyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	// 直接把comments返回,
	// 因为comment里预加载了User => comment.User.Username, etc.
	c.JSON(http.StatusOK, comments)
}

// GetBountiesByUser 获取指定用户发布的悬赏令
func (ctl *BountyController) GetBountiesByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	bounties, err := ctl.bountyService.GetBountiesByUserID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取悬赏令失败"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

// GetReceivedBounties 获取指定用户接收的悬赏令
func (ctl *BountyController) GetReceivedBounties(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	bounties, err := ctl.bountyService.GetReceivedBounties(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户已接受悬赏令失败"})
		return
	}

	c.JSON(http.StatusOK, bounties)
}

// UnlikeBounty 取消点赞悬赏令处理函数
func (ctl *BountyController) UnlikeBounty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	// 获取 bounty_id 路径参数并解析为 uuid.UUID
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	// 调用服务层方法获取互动信息
	interaction, err := ctl.bountyService.GetUserBountyInteraction(uid, bountyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "交互信息未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取交互信息失败"})
		}
		return
	}

	c.JSON(http.StatusOK, interaction)
}

// CancelSettlementByPublisher 发布方取消处于Settling状态的悬赏令
// POST /bounties/:bounty_id/cancel-settlement/publisher
func (ctl *BountyController) CancelSettlementByPublisher(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := ctl.bountyService.CancelSettlementByPublisher(bountyID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布方取消清算成功 (Cancelled with penalty)"})
}

// CancelSettlementByReceiver 接收方取消处于Settling状态的悬赏令
// POST /bounties/:bounty_id/cancel-settlement/receiver
func (ctl *BountyController) CancelSettlementByReceiver(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := ctl.bountyService.CancelSettlementByReceiver(bountyID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "接收方取消清算成功 (Cancelled with penalty)"})
}

// SettleBountyAccounts 结算悬赏令账户
func (ctl *BountyController) SettleBountyAccounts(c *gin.Context) {
	// 获取悬赏令ID从URL参数
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	// 调用服务层进行结算
	err = ctl.bountyService.SettleBountyAccounts(bountyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "悬赏令结算成功"})
}

// ConfirmMilestones 接收者确认提交所有里程碑
func (ctl *BountyController) ConfirmMilestones(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	// 从上下文中获取当前用户的 user_id
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
		return
	}

	err = ctl.bountyService.VerifyMilestones(bountyID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "所有的里程碑均被发布者所确认"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
		return
	}

	err = ctl.bountyService.ApplySettlement(bountyID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "悬赏令清算申请成功"})
}
