package router

import (
	"baokaobao/internal/config"
	"baokaobao/internal/handler"
	"baokaobao/internal/middleware"
	"baokaobao/internal/pkg/jwt"
	"baokaobao/internal/pkg/wechat"
	"baokaobao/internal/repository"
	"baokaobao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRouterWithDB(db *gorm.DB) *gin.Engine {
	middleware.InitJWT()

	repo := repository.NewRepository(db)
	middleware.BlacklistChecker = func(token string) bool {
		blacklisted, _ := repo.IsTokenBlacklisted(token)
		return blacklisted
	}

	jwtSDK := jwt.NewJWT(config.GlobalConfig.JWT.Secret, config.GlobalConfig.JWT.ExpireHours)
	wechatSDK := wechat.NewWechatSDK(config.GlobalConfig.Wechat.AppID, config.GlobalConfig.Wechat.Secret)

	svc := service.NewService(repo, jwtSDK, wechatSDK)
	h := handler.NewHandler(svc, wechatSDK)

	switch config.GlobalConfig.App.Mode {
	case gin.DebugMode, gin.TestMode, gin.ReleaseMode:
		gin.SetMode(config.GlobalConfig.App.Mode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Static("/uploads", "./uploads")

	r.GET("/health", func(c *gin.Context) {
		dbStatus := "ok"
		if sqlDB, err := db.DB(); err != nil || sqlDB.Ping() != nil {
			dbStatus = "error"
		}
		c.JSON(200, gin.H{"status": "ok", "db": dbStatus})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimit())
		{
			auth.POST("/login_by_wechat", h.LoginByWechat)
			auth.POST("/decrypt_phone", middleware.MiniProgramAuth(), h.DecryptPhone)
			auth.POST("/logout", middleware.MiniProgramAuth(), h.Logout)
		}

		user := api.Group("/user")
		user.Use(middleware.MiniProgramAuth(), middleware.RateLimit())
		{
			user.GET("/profile", h.GetProfile)
			user.PUT("/profile", h.UpdateProfile)
			user.POST("/avatar", h.UploadAvatar)
		}

		questionBank := api.Group("/question_banks")
		{
			questionBank.GET("", h.ListPublicQuestionBanks)
			questionBank.GET("/:id/access", middleware.MiniProgramAuth(), h.GetQuestionBankAccessStatus)
			questionBank.GET("/:id", h.GetPublicQuestionBankDetail)
		}

		question := api.Group("/questions")
		question.Use(middleware.MiniProgramAuth())
		{
			question.GET("", h.ListQuestions)
			question.GET("/:id", h.GetQuestion)
			question.GET("/random", h.RandomQuestions)
		}

		quiz := api.Group("/quiz")
		quiz.Use(middleware.MiniProgramAuth(), middleware.RateLimit())
		{
			quiz.POST("/submit", h.SubmitAnswer)
			quiz.POST("/submit_exam", h.SubmitExam)
			quiz.GET("/history", h.GetQuizHistory)
			quiz.GET("/wrong_questions", h.GetWrongQuestions)
			quiz.POST("/add_wrong/:qid", h.AddToWrongQuestions)
		}

		score := api.Group("/score")
		score.Use(middleware.MiniProgramAuth())
		{
			score.GET("/my", h.GetMyScore)
			score.GET("/ranking", h.GetRanking)
			score.GET("/stats", h.GetStats)
		}

		feedback := api.Group("/feedbacks")
		feedback.Use(middleware.MiniProgramAuth())
		{
			feedback.POST("", h.CreateFeedback)
		}

		exam := api.Group("/exam_records")
		exam.Use(middleware.MiniProgramAuth())
		{
			exam.GET("", h.GetExamRecords)
		}

		my := api.Group("/my")
		my.Use(middleware.MiniProgramAuth())
		{
			my.GET("/banks", h.GetMyPurchasedBanks)
		}
	}

	admin := r.Group("/admin/api/v1")
	{
		admin.POST("/login", h.AdminLogin)
		admin.POST("/logout", middleware.AdminAuth(), h.AdminLogout)
		admin.GET("/question_template", h.DownloadQuestionTemplate) // Public - no auth needed

		adminProtected := admin.Group("")
		adminProtected.Use(middleware.AdminAuth())
		{
			adminProtected.GET("/dashboard", h.GetDashboard)
			adminProtected.GET("/users", h.ListAllUsers)
			adminProtected.GET("/users/:id", h.GetUserDetail)
			adminProtected.GET("/question_banks", h.ListQuestionBanks)
			adminProtected.GET("/questions", h.ListAllQuestions)
		adminProtected.GET("/stats/overview", h.GetStatsOverview)
		adminProtected.GET("/stats/users", h.GetUserStats)
		adminProtected.GET("/stats/questions", h.GetQuestionStats)
		adminProtected.GET("/admin_users", h.ListAdminUsers)
		adminProtected.GET("/admin_users/:id", h.GetAdminUserDetail)
		adminProtected.GET("/audit_logs", h.ListAuditLogs)
		adminProtected.GET("/bank_stats", h.GetBankStatsList)
		adminProtected.GET("/bank_stats/:id", h.GetBankStats)
		adminProtected.GET("/feedbacks", h.ListFeedbacks)

		adminWrite := adminProtected.Group("")
			adminWrite.Use(middleware.RateLimit())
			{
			adminWrite.PUT("/users/:id/status", h.UpdateUserStatus)
			adminWrite.POST("/question_banks", h.CreateQuestionBank)
			adminWrite.PUT("/question_banks/:id", h.UpdateQuestionBank)
			adminWrite.DELETE("/question_banks/:id", h.DeleteQuestionBank)
			adminWrite.POST("/questions", h.CreateQuestion)
			adminWrite.PUT("/questions/:id", h.UpdateQuestion)
			adminWrite.DELETE("/questions/:id", h.DeleteQuestion)
			adminWrite.POST("/questions/import", h.ImportQuestions)
			adminWrite.POST("/admin_users", h.CreateAdminUser)
			adminWrite.PUT("/admin_users/:id", h.UpdateAdminUser)
			adminWrite.PUT("/admin_users/:id/reset_password", h.ResetAdminPassword)
			adminWrite.DELETE("/admin_users/:id", h.DeleteAdminUser)
			adminWrite.GET("/admin_users/:id/banks", h.ListAdminBankPermissions)
			adminWrite.POST("/admin_users/:id/banks", h.GrantAdminBankAccess)
			adminWrite.DELETE("/admin_users/:id/banks/:bankId", h.RevokeAdminBankAccess)
			adminWrite.PUT("/feedbacks/:id/status", h.UpdateFeedbackStatus)
			}
		}
	}

	zap.S().Info("Router initialized")
	return r
}
