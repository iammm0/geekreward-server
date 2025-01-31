package main

import (
	"GeekReward/inernal/app/controllers"
	"GeekReward/inernal/app/repositories"
	"GeekReward/inernal/app/routes"
	"GeekReward/inernal/app/services"
	utils "GeekReward/inernal/app/validators"
	"GeekReward/pkg/database"
	"GeekReward/pkg/logger"
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
	userRepo := repositories.NewUserRepository(database.DB)
	bountyRepo := repositories.NewBountyRepository(database.DB)
	geekRepo := repositories.NewGeekRepository(database.DB)
	milestoneRepo := repositories.NewMilestoneRepository(database.DB)
	notificationRepo := repositories.NewNotificationRepository(database.DB)
	applicationRepo := repositories.NewApplicationRepository(database.DB)
	invitationRepo := repositories.NewInvitationRepository(database.DB)

	// 初始化服务
	authService := services.NewAuthService(userRepo)
	bountyService := services.NewBountyService(userRepo, bountyRepo, applicationRepo, notificationRepo, milestoneRepo)
	geekService := services.NewGeekService(geekRepo, invitationRepo)
	userService := services.NewUserService(userRepo)
	milestoneService := services.NewMilestoneService(milestoneRepo, bountyRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	applicationService := services.NewApplicationService(applicationRepo, bountyRepo)
	invitationService := services.NewInvitationService(invitationRepo, userRepo)

	// 初始化控制器
	authController := controllers.NewAuthController(authService, notificationService)
	bountyController := controllers.NewBountyController(bountyService, milestoneService, notificationService)
	geekController := controllers.NewGeekController(geekService, notificationService)
	userController := controllers.NewUserController(userService, notificationService)
	notificationController := controllers.NewNotificationController(notificationService)
	applicationController := controllers.NewApplicationController(applicationService, bountyService, notificationService)
	milestoneController := controllers.NewMilestoneController(milestoneService, notificationService)
	invitationController := controllers.NewInvitationController(invitationService, notificationService)
	attachmentController := controllers.NewAttachmentController(notificationService)

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
