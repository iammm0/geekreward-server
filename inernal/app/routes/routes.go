package routes

import (
	"GeekReward/inernal/app/controllers"
	"GeekReward/inernal/app/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

// SetupRouter 配置路由
func SetupRouter(
	authController *controllers.AuthController,
	bountyController *controllers.BountyController,
	geekController *controllers.GeekController,
	userController *controllers.UserController,
	notificationController *controllers.NotificationController,
	applicationController *controllers.ApplicationController,
) *gin.Engine {
	// 创建Gin路由引擎实例
	r := gin.Default()

	// 配置CORS中间件，允许来自前端的跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"}, // 修改为你的前端URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 设置静态文件路由
	r.Static("/static", "./frontend/dist/static")

	// 注册路由和控制器
	api := r.Group("/api")
	{
		// 用户认证相关路由
		api.POST("/register", authController.Register) // 用户注册
		api.POST("/login", authController.Login)       // 用户登录

		// 悬赏令及里程碑相关路由
		api.GET("/bounties", bountyController.GetBounties)                                                 // 获取所有悬赏令
		api.POST("/bounties", middlewares.JWTAuthMiddleware(), bountyController.CreateBounty)              // 创建悬赏令（需JWT认证）
		api.GET("/bounties/:bounty_id", bountyController.GetBounty)                                        // 获取指定悬赏令
		api.GET("/bounties/:bounty_id/comments", bountyController.GetBountyComments)                       // 获取指定悬赏令的评论
		api.PUT("/bounties/:bounty_id", middlewares.JWTAuthMiddleware(), bountyController.UpdateBounty)    // 更新悬赏令（需JWT认证）
		api.DELETE("/bounties/:bounty_id", middlewares.JWTAuthMiddleware(), bountyController.DeleteBounty) // 删除悬赏令（需JWT认证）
		api.GET("/bounties/:bounty_id/milestones", bountyController.GetMilestonesByBountyID)               // 获取指定悬赏令的里程碑
		api.POST("/bounties/:bounty_id/like", middlewares.JWTAuthMiddleware(), bountyController.LikeBounty)
		api.DELETE("/bounties/:bounty_id/unlike", middlewares.JWTAuthMiddleware(), bountyController.UnlikeBounty)               // 点赞悬赏令（需JWT认证）
		api.POST("/bounties/:bounty_id/comment", middlewares.JWTAuthMiddleware(), bountyController.CommentOnBounty)             // 评论悬赏令（需JWT认证）
		api.POST("/bounties/:bounty_id/rate", middlewares.JWTAuthMiddleware(), bountyController.RateBounty)                     // 评分悬赏令（需JWT认证）
		api.GET("/bounties/:bounty_id/interaction", middlewares.JWTAuthMiddleware(), bountyController.GetUserBountyInteraction) // 获取用户的悬赏令互动信息（需JWT认证）

		// 极客相关路由
		api.GET("/geeks", geekController.GetTopGeeks)     // 获取极客排行榜
		api.GET("/geeks/:id", geekController.GetGeekByID) // 获取指定ID的极客信息

		// 用户信息相关路由
		api.GET("/user/profile", middlewares.JWTAuthMiddleware(), userController.GetUserInfo)                     // 获取用户信息（需JWT认证）
		api.PUT("/user/profile", middlewares.JWTAuthMiddleware(), userController.UpdateUserInfo)                  // 更新用户信息（需JWT认证）
		api.GET("/user/bounties", middlewares.JWTAuthMiddleware(), bountyController.GetBountiesByUser)            // 获取用户发布的悬赏令（需JWT认证）
		api.GET("/user/received-bounties", middlewares.JWTAuthMiddleware(), bountyController.GetReceivedBounties) // 获取用户接收的悬赏令（需JWT认证）

		// 通知相关路由
		api.GET("/notifications", middlewares.JWTAuthMiddleware(), notificationController.GetUserNotifications)            // 获取用户的所有通知（需JWT认证）
		api.PUT("/notifications/:id/read", middlewares.JWTAuthMiddleware(), notificationController.MarkNotificationAsRead) // 标记通知为已读（需JWT认证）
		api.DELETE("/notifications/:id", middlewares.JWTAuthMiddleware(), notificationController.DeleteNotification)       // 删除通知（需JWT认证）

		// 悬赏令申请相关路由
		api.POST("/applications", middlewares.JWTAuthMiddleware(), applicationController.CreateApplication)             // 创建悬赏令申请（需JWT认证）
		api.GET("/applications/:bounty_id", middlewares.JWTAuthMiddleware(), applicationController.GetApplications)     // 获取某个悬赏令的所有申请（需JWT认证）
		api.PUT("/applications/:id/approve", middlewares.JWTAuthMiddleware(), applicationController.ApproveApplication) // 批准悬赏令申请（需JWT认证）
		api.PUT("/applications/:id/reject", middlewares.JWTAuthMiddleware(), applicationController.RejectApplication)   // 拒绝悬赏令申请（需JWT认证）
	}

	// 处理未匹配的路由，返回前端应用的 index.html
	r.NoRoute(func(c *gin.Context) {
		// 如果是请求的静态文件，返回对应文件
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.File("./frontend/dist" + c.Request.URL.Path)
		} else {
			// 否则返回 SPA 的 index.html
			c.File("./frontend/dist/index.html")
		}
	})

	return r
}
