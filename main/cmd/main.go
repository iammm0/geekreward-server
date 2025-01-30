package main

import (
	controllers2 "GeekReward/main/inernal/app/controllers"
	repositories2 "GeekReward/main/inernal/app/repositories"
	"GeekReward/main/inernal/app/routes"
	services2 "GeekReward/main/inernal/app/services"
	"GeekReward/main/inernal/app/validators"
	"GeekReward/main/pkg/database"
	"GeekReward/main/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

func main() {
	// 设置 Gin 为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化日志记录器
	logger.InitLogger()

	// 初始化 Viper 以读取 config.yaml
	viper.SetConfigName("config") // 不包括扩展名
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // 从当前目录查找配置文件

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config.yaml file: %s", err)
	}

	// 通过 Viper 获取配置信息
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.name")

	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetInt("redis.port")

	// 使用日志记录器记录数据库和Redis连接信息
	logger.InfoLogger.WithFields(logrus.Fields{
		"Database Name": dbName,
		"Database Host": dbHost,
		"Database Port": dbPort,
		"Database User": dbUser,
	}).Info("Connecting to database")

	logger.InfoLogger.WithFields(logrus.Fields{
		"Redis Host": redisHost,
		"Redis Port": redisPort,
	}).Info("Connecting to Redis")

	// 使用 Viper 配置数据库连接
	database.ConnectDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)

	// 初始化仓库
	userRepo := repositories2.NewUserRepository(database.DB)
	bountyRepo := repositories2.NewBountyRepository(database.DB)
	geekRepo := repositories2.NewGeekRepository(database.DB)
	milestoneRepo := repositories2.NewMilestoneRepository(database.DB)
	notificationRepo := repositories2.NewNotificationRepository(database.DB)
	applicationRepo := repositories2.NewApplicationRepository(database.DB)
	invitationRepo := repositories2.NewInvitationRepository(database.DB)

	// 初始化服务
	authService := services2.NewAuthService(userRepo)
	bountyService := services2.NewBountyService(userRepo, bountyRepo, applicationRepo, notificationRepo, milestoneRepo)
	geekService := services2.NewGeekService(geekRepo, invitationRepo)
	userService := services2.NewUserService(userRepo)
	milestoneService := services2.NewMilestoneService(milestoneRepo, bountyRepo)
	notificationService := services2.NewNotificationService(notificationRepo)
	applicationService := services2.NewApplicationService(applicationRepo, bountyRepo)
	invitationService := services2.NewInvitationService(invitationRepo, userRepo)

	// 初始化控制器
	authController := controllers2.NewAuthController(authService, notificationService)
	bountyController := controllers2.NewBountyController(bountyService, milestoneService, notificationService)
	geekController := controllers2.NewGeekController(geekService, notificationService)
	userController := controllers2.NewUserController(userService, notificationService)
	notificationController := controllers2.NewNotificationController(notificationService)
	applicationController := controllers2.NewApplicationController(applicationService, bountyService, notificationService)
	milestoneController := controllers2.NewMilestoneController(milestoneService, notificationService)
	invitationController := controllers2.NewInvitationController(invitationService, notificationService)
	attachmentController := controllers2.NewAttachmentController(notificationService)

	// 初始化验证器
	validatorInstance, err := utils.NewValidator()
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to initialize validator: %v", err)
	}

	// 设置路由
	r := routes.SetupRouter(
		authController,
		bountyController,
		geekController,
		userController,
		notificationController,
		applicationController,
		milestoneController,
		invitationController,
		attachmentController,
	)

	// 传递给需要的组件或通过中间件设置到上下文中
	r.Use(func(c *gin.Context) {
		c.Set("validator", validatorInstance)
		c.Next()
	})

	// 获取端口号，默认为8080
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	// 启动Gin服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
