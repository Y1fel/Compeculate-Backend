package main

import (
	"log"

	"WechatGo/config"
	"WechatGo/handlers"
	"WechatGo/middleware"
	"WechatGo/services"
	"WechatGo/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建数据库服务
	dbService, err := services.NewDatabaseService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建JWT工具
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret)

	// 创建处理器
	authHandler := handlers.NewAuthHandler(dbService, jwtUtil)
	userHandler := handlers.NewUserHandler(dbService)
	quizHandler := handlers.NewQuizHandler(dbService)
	rankingHandler := handlers.NewRankingHandler(dbService)

	// 创建Gin实例
	r := gin.New()

	// 使用中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关接口
		user := api.Group("/user")
		{
			// 用户登录（不需要认证）
			user.POST("/login", authHandler.Login)

			// 需要认证的用户接口
			userAuth := user.Group("")
			userAuth.Use(middleware.AuthMiddleware(jwtUtil))
			{
				userAuth.GET("/stats", userHandler.GetUserStats)
				userAuth.GET("/profile", userHandler.GetUserProfile)
			}
		}

		// 答题相关接口
		quiz := api.Group("/quiz")
		quiz.Use(middleware.AuthMiddleware(jwtUtil))
		{
			quiz.POST("/submit", quizHandler.SubmitQuizResult)
		}

		// 排行榜相关接口
		ranking := api.Group("/ranking")
		ranking.Use(middleware.AuthMiddleware(jwtUtil))
		{
			ranking.GET("/list", rankingHandler.GetRankingList)
		}
	}

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	err = r.Run(":" + cfg.Server.Port)
}
