package controllers

import (
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// ApplicationController 结构体
type ApplicationController struct {
	applicationService  services.ApplicationService
	bountyService       services.BountyService
	notificationService services.NotificationService
}

// NewApplicationController 创建新的 ApplicationController 实例
func NewApplicationController(
	applicationService services.ApplicationService,
	bountyService services.BountyService,
	notificationService services.NotificationService,
) *ApplicationController {
	return &ApplicationController{
		applicationService:  applicationService,
		bountyService:       bountyService,
		notificationService: notificationService,
	}
}

// CreateApplication 创建新的悬赏任务申请
func (ctl *ApplicationController) CreateApplication(c *gin.Context) {
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

	// 解析可选的 note
	var input struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		// note不是必填，所以如果解析失败也许只是json不对
		// 你可做更友好的处理
	}

	// 获取 bounty, 判断是否是自己发的
	bounty, err := ctl.bountyService.GetBounty(bountyID)
	if err != nil || bounty == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "悬赏令未找到"})
		return
	}
	if bounty.UserID == uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能申请自己发布的悬赏令"})
		return
	}

	// 检查用户是否已经申请过该悬赏令
	hasApplied, err := ctl.applicationService.HasUserApplied(bountyID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查申请记录失败"})
		return
	}

	if hasApplied {
		c.JSON(http.StatusBadRequest, gin.H{"error": "你已经申请过此悬赏令"})
		return
	}

	// 创建申请
	if err := ctl.applicationService.CreateApplication(bountyID, uid, input.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "申请成功"})
}

// GetApplications 获取某个悬赏任务的所有申请
func (ctl *ApplicationController) GetApplications(c *gin.Context) {
	// 从路由字段获取 bounty_id 的字段值
	// 声明并赋值 bountyIDStr 字段
	bountyIDStr := c.Param("bounty_id")
	// 使用 uuid.Parse() 类型转化函数
	bountyID, err := uuid.Parse(bountyIDStr)
	// 该类型转化函数会返回两个可能的值 ，如若 err 不为空 那么就返回对应的状态码与错误信息
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 悬赏令ID"})
		return
	}

	// 这里要不要验证 "是否当前用户就是发布者"? 视业务需求
	// 例如先 getBounty(bountyID), 如果 bounty.UserID != currentUser => 403

	// 调用控制器对象所绑定的 applicationService 模块的方法的引用
	applications, err := ctl.applicationService.GetApplications(bountyID)
	// 同样
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取悬赏令的申请信息失败"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// ApproveApplication 批准申请
func (ctl *ApplicationController) ApproveApplication(c *gin.Context) {
	applicationIDStr := c.Param("application_id")
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的申请ID"})
		return
	}

	// 可能需要添加发布者身份验证，确保只有悬赏令的发布者才能批准申请
	// 例如，通过JWT解析当前用户ID，并检查是否为悬赏令的发布者

	if err := ctl.applicationService.ApproveApplication(applicationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批准申请失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批准申请成功"})
}

// RejectApplication 拒绝申请
func (ctl *ApplicationController) RejectApplication(c *gin.Context) {
	applicationIDStr := c.Param("application_id")
	applicationID, err := uuid.Parse(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的申请ID"})
		return
	}

	// 可能需要添加发布者身份验证，确保只有悬赏令的发布者才能拒绝申请

	if err := ctl.applicationService.RejectApplication(applicationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "拒绝申请失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "拒绝申请成功"})
}

// GetPublicApplications 获取公开的申请信息
func (ctl *ApplicationController) GetPublicApplications(c *gin.Context) {
	bountyIDStr := c.Param("bounty_id")
	bountyID, err := uuid.Parse(bountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的悬赏令ID"})
		return
	}

	applications, err := ctl.applicationService.GetPublicApplications(bountyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公开申请失败"})
		return
	}

	c.JSON(http.StatusOK, applications)
}
