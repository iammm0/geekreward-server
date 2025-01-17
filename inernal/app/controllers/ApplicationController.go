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

	if err := ctl.applicationService.CreateApplication(bountyID, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "申请失败"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applications)
}
