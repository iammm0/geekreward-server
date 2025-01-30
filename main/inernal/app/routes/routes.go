package routes

import (
	controllers2 "GeekReward/main/inernal/app/controllers"
	"GeekReward/main/inernal/app/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter(
	authController *controllers2.AuthController,
	bountyController *controllers2.BountyController,
	geekController *controllers2.GeekController,
	userController *controllers2.UserController,
	notificationController *controllers2.NotificationController,
	applicationController *controllers2.ApplicationController,
	milestoneController *controllers2.MilestoneController,
	invitationController *controllers2.InvitationController,
	attachmentController *controllers2.AttachmentController,
) *gin.Engine {
	// 创建Gin路由引擎实例
	r := gin.Default()

	// 配置CORS中间件，允许来自前端的跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 修改为你的前端URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// 设置静态文件路由
	// r.Static("/static", "./frontend/dist/static")
	r.Static("/uploads", "./uploads")

	// 注册路由和控制器
	api := r.Group("")
	// api.Use(middlewares.JWTAuthMiddleware())
	{
		// 处理文件上传
		api.POST("/attachment", attachmentController.UploadAttachment)

		// 用户认证相关路由
		api.POST("/register", authController.Register) // 用户注册
		api.POST("/login", authController.Login)       // 用户登录

		// 悬赏令相关路由
		// GET /bounties?status=Settled&publisher_id=...&receiver_id=...&limit=10&offset=0
		api.GET("/bounties", bountyController.GetBounties)                                                                       // 获取所有悬赏令
		api.POST("/bounties", middlewares.JWTAuthMiddleware(), bountyController.CreateBounty)                                    // 创建悬赏令（需JWT认证）
		api.GET("/bounties/:bounty_id", bountyController.GetBounty)                                                              // 获取指定悬赏令
		api.GET("/bounties/:bounty_id/comments", bountyController.GetBountyComments)                                             // 获取指定悬赏令的评论
		api.PUT("/bounties/:bounty_id", middlewares.JWTAuthMiddleware(), bountyController.UpdateBounty)                          // 更新悬赏令（需JWT认证）
		api.DELETE("/bounties/:bounty_id", middlewares.JWTAuthMiddleware(), bountyController.DeleteBounty)                       // 删除悬赏令（需JWT认证）
		api.POST("/bounties/:bounty_id/like", middlewares.JWTAuthMiddleware(), bountyController.LikeBounty)                      // 点赞悬赏令（需JWT认证）
		api.DELETE("/bounties/:bounty_id/unlike", middlewares.JWTAuthMiddleware(), bountyController.UnlikeBounty)                // 取消点赞悬赏令（需JWT认证）
		api.POST("/bounties/:bounty_id/comment", middlewares.JWTAuthMiddleware(), bountyController.CommentOnBounty)              // 评论悬赏令（需JWT认证）
		api.POST("/bounties/:bounty_id/rate", middlewares.JWTAuthMiddleware(), bountyController.RateBounty)                      // 评分悬赏令（需JWT认证）
		api.GET("/bounties/:bounty_id/interaction", middlewares.JWTAuthMiddleware(), bountyController.GetUserBountyInteraction)  // 获取用户的悬赏令互动信息（需JWT认证）
		api.POST("/bounties/:bounty_id/settle-accounts", middlewares.JWTAuthMiddleware(), bountyController.SettleBountyAccounts) // 结算悬赏令（需JWT认证）
		// 发布方取消
		api.POST("/bounties/:bounty_id/cancel-settlement/publisher", middlewares.JWTAuthMiddleware(), bountyController.CancelSettlementByPublisher)
		// 接收方取消
		api.POST("/bounties/:bounty_id/cancel-settlement/receiver", middlewares.JWTAuthMiddleware(), bountyController.CancelSettlementByReceiver)

		// 极客相关路由
		api.GET("/geeks", geekController.GetTopGeeks)                                                                        // 获取极客排行榜概览信息
		api.GET("/geeks/:id", geekController.GetGeekByID)                                                                    // 获取指定ID的极客公开信息
		api.POST("/geeks/:id/invitation", middlewares.JWTAuthMiddleware(), geekController.SendInvitation)                    // 向特定极客发出组队邀请（需JWT认证）
		api.POST("/geeks/:id", middlewares.JWTAuthMiddleware(), geekController.ExpressAffection)                             // 向特定极客或团队表达好感（需JWT认证）
		api.PUT("/invitation/:invitation_id/accept", middlewares.JWTAuthMiddleware(), invitationController.AcceptInvitation) // 接受组队邀请（需JWT认证）
		api.PUT("/invitation/:invitation_id/reject", middlewares.JWTAuthMiddleware(), invitationController.RejectInvitation) // 拒绝组队邀请（需JWT认证）
		api.POST("/geeks/:id/express-affection", middlewares.JWTAuthMiddleware(), geekController.ExpressAffection)

		// 用户信息相关路由
		api.GET("/user/profile", middlewares.JWTAuthMiddleware(), userController.GetUserInfo)                     // 获取用户信息（需JWT认证）
		api.PUT("/user/profile", middlewares.JWTAuthMiddleware(), userController.UpdateUserInfo)                  // 更新用户信息（需JWT认证）
		api.GET("/user/bounties", middlewares.JWTAuthMiddleware(), bountyController.GetBountiesByUser)            // 获取用户发布的悬赏令（需JWT认证）
		api.GET("/user/received-bounties", middlewares.JWTAuthMiddleware(), bountyController.GetReceivedBounties) // 获取用户接收的悬赏令（需JWT认证）

		// 通知相关路由
		api.GET("/notifications", middlewares.JWTAuthMiddleware(), notificationController.GetUserNotifications)            // 获取用户的所有通知（需JWT认证）
		api.PUT("/notifications/:id/read", middlewares.JWTAuthMiddleware(), notificationController.MarkNotificationAsRead) // 标记通知为已读（需JWT认证）
		api.DELETE("/notifications/:id", middlewares.JWTAuthMiddleware(), notificationController.DeleteNotification)       // 删除通知（需JWT认证）
		// 下面这个多为测试/管理员使用
		api.POST("/notifications", middlewares.JWTAuthMiddleware(), notificationController.CreateNotification)

		// 悬赏令申请相关路由
		api.POST("/applications/:bounty_id", middlewares.JWTAuthMiddleware(), applicationController.CreateApplication)              // 向特定悬赏任务发出悬赏令申请（需JWT认证）
		api.GET("/applications/:bounty_id/private", middlewares.JWTAuthMiddleware(), applicationController.GetApplications)         // 发布者获取某个悬赏令的所有申请（需JWT认证）
		api.GET("applications/:bounty_id/public", middlewares.JWTAuthMiddleware(), applicationController.GetPublicApplications)     // 公开的申请信息
		api.PUT("/applications/:application_id/approve", middlewares.JWTAuthMiddleware(), applicationController.ApproveApplication) // 批准悬赏令申请（需JWT认证）
		api.PUT("/applications/:application_id/reject", middlewares.JWTAuthMiddleware(), applicationController.RejectApplication)   // 拒绝悬赏令申请（需JWT认证）

		// 里程碑相关路由
		api.GET("/bounties/:bounty_id/milestones", milestoneController.GetMilestonesByBountyID)                                                           // 获取指定悬赏令的里程碑
		api.POST("/bounties/:bounty_id/milestones", middlewares.JWTAuthMiddleware(), milestoneController.CreateMilestone)                                 // 悬赏令发布者公布里程碑
		api.PUT("/bounties/:bounty_id/milestones/:milestone_id/promulgator", middlewares.JWTAuthMiddleware(), milestoneController.UpdateMilestone)        // 悬赏令发布者更新里程碑（需JWT认证）
		api.PUT("/bounties/:bounty_id/milestones/:milestone_id/receiver", middlewares.JWTAuthMiddleware(), milestoneController.UpdateMilestoneByReceiver) // 悬赏零接收者更新里程碑（需JWT认证）
		api.DELETE("/bounties/:bounty_id/milestones/:milestone_id", middlewares.JWTAuthMiddleware(), milestoneController.DeleteMilestone)                 // 删除里程碑（需JWT认证）

		// 新增的悬赏令状态相关路由
		api.POST("/bounties/:bounty_id/confirm-milestones", middlewares.JWTAuthMiddleware(), bountyController.ConfirmMilestones) // 接收者确认提交所有里程碑
		api.POST("/bounties/:bounty_id/verify-milestones", middlewares.JWTAuthMiddleware(), bountyController.VerifyMilestones)   // 发布者审核并确认所有里程碑
		api.POST("/bounties/:bounty_id/settle", middlewares.JWTAuthMiddleware(), bountyController.ApplySettlement)               // 接收者申请悬赏令清算
	}

	// 处理未匹配的路由，返回前端应用的 index.html（如果使用SPA）

	// 处理未匹配的路由，返回前端应用的 index.html
	// r.NoRoute(func(c *gin.Context) {
	// 	// 如果是请求的静态文件，返回对应文件
	// 	if strings.HasPrefix(c.Request.URL.Path, "/static/") {
	// 		c.File("./frontend/dist" + c.Request.URL.Path)
	// 	} else {
	// 		// 否则返回 SPA 的 index.html
	// 		c.File("./frontend/dist/index.html")
	// 	}
	// })

	return r
}
