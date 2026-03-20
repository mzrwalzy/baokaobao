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
	jwtSDK := jwt.NewJWT(config.GlobalConfig.JWT.Secret, config.GlobalConfig.JWT.ExpireHours)
	wechatSDK := wechat.NewWechatSDK(config.GlobalConfig.Wechat.AppID, config.GlobalConfig.Wechat.Secret)

	svc := service.NewService(repo, jwtSDK, wechatSDK)
	h := handler.NewHandler(svc, wechatSDK)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login_by_wechat", h.LoginByWechat)
			auth.POST("/decrypt_phone", middleware.MiniProgramAuth(), h.DecryptPhone)
			auth.POST("/logout", middleware.MiniProgramAuth(), h.Logout)
		}

		user := api.Group("/user")
		user.Use(middleware.MiniProgramAuth())
		{
			user.GET("/profile", h.GetProfile)
			user.PUT("/profile", h.UpdateProfile)
			user.POST("/avatar", h.UploadAvatar)
		}

		question := api.Group("/questions")
		question.Use(middleware.MiniProgramAuth())
		{
			question.GET("", h.ListQuestions)
			question.GET("/:id", h.GetQuestion)
			question.GET("/random", h.RandomQuestions)
		}

		quiz := api.Group("/quiz")
		quiz.Use(middleware.MiniProgramAuth())
		{
			quiz.POST("/submit", h.SubmitAnswer)
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
			adminProtected.PUT("/users/:id/status", h.UpdateUserStatus)

			adminProtected.GET("/question_banks", h.ListQuestionBanks)
			adminProtected.POST("/question_banks", h.CreateQuestionBank)
			adminProtected.PUT("/question_banks/:id", h.UpdateQuestionBank)
			adminProtected.DELETE("/question_banks/:id", h.DeleteQuestionBank)

			adminProtected.GET("/questions", h.ListAllQuestions)
			adminProtected.POST("/questions", h.CreateQuestion)
			adminProtected.PUT("/questions/:id", h.UpdateQuestion)
			adminProtected.DELETE("/questions/:id", h.DeleteQuestion)
			adminProtected.POST("/questions/import", h.ImportQuestions)

			adminProtected.GET("/stats/overview", h.GetStatsOverview)
			adminProtected.GET("/stats/users", h.GetUserStats)
			adminProtected.GET("/stats/questions", h.GetQuestionStats)
		}
	}

	zap.S().Info("Router initialized")
	return r
}
