package routes

import (
	"backend/handlers"
	"backend/middleware"
	"backend/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize services
	verificationService := services.NewVerificationService()
	userStatsService := services.NewUserStatsService()
	userService := services.NewUserService(verificationService, userStatsService)
	questionService := services.NewQuestionService()
	questionStatsService := services.NewQuestionStatsService()
	quizService := services.NewQuizService(questionService, userStatsService, questionStatsService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, verificationService)
	userHandler := handlers.NewUserHandler(userService)
	userStatsHandler := handlers.NewUserStatsHandler(userStatsService)
	questionHandler := handlers.NewQuestionHandler(questionService)
	questionStatsHandler := handlers.NewQuestionStatsHandler(questionStatsService)
	quizHandler := handlers.NewQuizHandler(quizService)

	// Authentication routes
	authRoutes := r.Group("/auth")
	{
		// 无需认证的接口
		authRoutes.POST("/register-request", authHandler.RegisterRequest)
		authRoutes.POST("/complete-registration", authHandler.CompleteRegistration)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("login-admin", authHandler.LoginAdmin)
		authRoutes.POST("/send-verification", authHandler.SendVerificationCode)
		authRoutes.POST("/verify-code", authHandler.VerifyCode)
		authRoutes.POST("/update-password", authHandler.UpdatePassword)

		// 需要认证的接口
		authRoutes.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
	}

	// User profile routes
	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/profile", userHandler.GetProfile)
		userRoutes.POST("/update", userHandler.UpdateProfile)
		userRoutes.GET("/presign-upload", userHandler.PresignUploadURL)
	}

	// User profile routes
	userStatsRoutes := r.Group("/user-stats")
	userStatsRoutes.Use(middleware.AuthMiddleware())
	{
		userStatsRoutes.GET("/", userStatsHandler.GetUserStats)
	}

	// Question routes
	questionRoutes := r.Group("/question")
	questionRoutes.Use(middleware.AuthMiddleware()) //需要认证（通常只有管理员可以创建题目）
	{
		// 生成题目并加入题库
		questionRoutes.POST("/generate", questionHandler.GenerateQuestion)
		// 获取单个题目
		questionRoutes.GET("/question/:id", questionHandler.GetQuestionById)
		// 获取题目列表 支持分页和按类别过滤
		questionRoutes.GET("/questions", questionHandler.GetQuestionList)
		// 批量删除题目
		questionRoutes.POST("/delete", questionHandler.BatchDeleteQuestions)
	}

	// Question stats routes
	questionStatsRoutes := r.Group("/question-stats")
	questionStatsRoutes.Use(middleware.AuthMiddleware())
	{
		questionStatsRoutes.GET("/dimension-distribution", questionStatsHandler.GetDimensionDistribution)
		questionStatsRoutes.GET("/dimension-accuracy", questionStatsHandler.GetDimensionAccuracy)
	}

	// Quiz routes
	quizRoutes := r.Group("/quiz")
	// 使用JWT认证中间件保护Quiz相关的路由
	quizRoutes.Use(middleware.AuthMiddleware())
	{
		quizRoutes.POST("/new", quizHandler.CreateQuiz)
		quizRoutes.GET("/:id", quizHandler.GetQuiz)
		quizRoutes.POST("/submit", quizHandler.SubmitQuiz)
		quizRoutes.GET("/history", quizHandler.GetUserQuizHistory)
	}

	return r
}
