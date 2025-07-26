package routes

import (
	"online-exam-system/controllers"
	"online-exam-system/middleware"
	"online-exam-system/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// API版本前缀
	api := r.Group("/api/v1")

	// 公开路由（无需认证）
	public := api.Group("/")
	{
		// 认证相关
		auth := public.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register) // 注册功能（可能需要管理员权限）
		}
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.TenantMiddleware()) // 添加租户中间件
	{
		// 用户相关
		user := protected.Group("/user")
		{
			user.GET("/profile", controllers.GetProfile)
			user.PUT("/profile", controllers.UpdateProfile)
			user.PUT("/password", controllers.ChangePassword)
		}

		// 科目相关（所有角色都可以查看）
		protected.GET("/subjects", controllers.GetSubjects)
		protected.GET("/subjects/all", controllers.GetAllSubjects)
		protected.GET("/subjects/:id", controllers.GetSubject)
		protected.GET("/subjects/:id/stats", controllers.GetSubjectStats)

		// 题目相关
		protected.GET("/questions", controllers.GetQuestions)
		protected.GET("/questions/stats", controllers.GetQuestionStats)
		protected.GET("/questions/:id", controllers.GetQuestion)
		protected.GET("/questions/:id/analyze", controllers.AnalyzeQuestion) // AI分析题目

		// 试卷相关
		paper := protected.Group("/papers")
		{
			paper.GET("/", controllers.GetPapers)
			paper.GET("/:id", controllers.GetPaper)
		}

		// 考试相关
		exam := protected.Group("/exams")
		{
			exam.GET("/", controllers.GetExams)
			exam.GET("/student", controllers.GetStudentExams) // 学生考试列表
			exam.GET("/:id", controllers.GetExam)
			exam.POST("/:id/start", controllers.StartExam) // 学生开始考试
			exam.GET("/:id/result", controllers.GetExamResult) // 考试结果
			exam.GET("/:id/analysis", controllers.GetExamAnalysis) // 考试分析
		}

		// 答题相关（学生专用）
		answer := protected.Group("/answers")
		answer.Use(middleware.RoleMiddleware(models.RoleStudent))
		{
			answer.POST("/exam/:exam_id", controllers.SubmitAnswer) // 提交单个答案
			answer.POST("/exam/:exam_id/submit", controllers.SubmitExam) // 提交整份试卷
			answer.GET("/exam/:exam_id", controllers.GetStudentAnswers) // 获取学生答案
		}

		// AI问答相关
		ai := protected.Group("/ai")
		{
			ai.POST("/chat", controllers.ChatWithAI)
			ai.GET("/history", controllers.GetAIChatHistory)
			ai.DELETE("/history", controllers.ClearAIChatHistory)
		}

		// 练习相关（学生专用）
		practice := protected.Group("/practice")
		practice.Use(middleware.RoleMiddleware(models.RoleStudent))
		{
			practice.GET("/recommendations", controllers.GetPracticeRecommendations) // 获取推荐练习
			practice.POST("/start", controllers.StartPractice)                      // 开始练习
			practice.POST("/:practice_id/answer", controllers.SubmitPracticeAnswer)  // 提交答案
			practice.POST("/:practice_id/complete", controllers.CompletePractice)    // 完成练习
			practice.GET("/history", controllers.GetPracticeHistory)                // 练习历史
			practice.GET("/stats", controllers.GetPracticeStats)                    // 练习统计
			practice.GET("/wrong-questions", controllers.GetWrongQuestions)         // 获取错题列表
			practice.POST("/review/start", controllers.StartWrongQuestionReview)    // 开始错题复习
		}

		// 统计相关
		stats := protected.Group("/stats")
		{
			stats.GET("/student", middleware.RoleMiddleware(models.RoleStudent), controllers.GetStudentStats)
			stats.GET("/teacher", middleware.RoleMiddleware(models.RoleTeacher), controllers.GetTeacherStats)
		}
	}

	// 管理员专用路由
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.TenantMiddleware()) // 添加租户中间件
	admin.Use(middleware.RoleMiddleware(models.RoleAdmin))
	{
		// 用户管理
		users := admin.Group("/users")
		{
			users.GET("/", controllers.GetUsers)
			users.GET("/:id", controllers.GetUser)
			users.POST("/", controllers.CreateUser)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
			users.PUT("/:id/password", controllers.ResetPassword)
			users.POST("/import", controllers.BatchImportUsers)
		}

		// 科目管理
		subjects := admin.Group("/subjects")
		{
			subjects.POST("/", controllers.CreateSubject)
			subjects.PUT("/:id", controllers.UpdateSubject)
			subjects.DELETE("/:id", controllers.DeleteSubject)
		}

		// 仪表板统计
		admin.GET("/dashboard", controllers.GetDashboardStats)
	}

	// 教师专用路由
	teacher := api.Group("/teacher")
	teacher.Use(middleware.AuthMiddleware())
	teacher.Use(middleware.TenantMiddleware()) // 添加租户中间件
	teacher.Use(middleware.RoleMiddleware(models.RoleTeacher, models.RoleAdmin)) // 管理员也可以访问
	{
		// 题目管理
		questions := teacher.Group("/questions")
		{
			questions.POST("/", controllers.CreateQuestion)
			questions.PUT("/:id", controllers.UpdateQuestion)
			questions.DELETE("/:id", controllers.DeleteQuestion)
			questions.POST("/import", controllers.BatchImportQuestions)
		}

		// 试卷管理
		papers := teacher.Group("/papers")
		{
			papers.POST("/", controllers.CreatePaper)
			papers.POST("/auto", controllers.AutoCreatePaper) // 自动组卷
			papers.PUT("/:id", controllers.UpdatePaper)
			papers.DELETE("/:id", controllers.DeletePaper)
		}

		// 考试管理
		exams := teacher.Group("/exams")
		{
			exams.POST("/", controllers.CreateExam)
			exams.PUT("/:id", controllers.UpdateExam)
			exams.DELETE("/:id", controllers.DeleteExam)
		}
	}

	// 健康检查
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "在线考试系统API运行正常",
		})
	})
}