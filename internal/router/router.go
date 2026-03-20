package router

import (
	"baokaobao/internal/handler"
	"baokaobao/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerGroup struct {
	Auth     *handler.AuthHandler
	User     *handler.UserHandler
	Question *handler.QuestionHandler
	Quiz     *handler.QuizHandler
	Score    *handler.ScoreHandler
	Admin    *handler.AdminHandler
}

func SetupRouter(group *HandlerGroup) *gin.Engine {
	if middleware.GetJWT() == nil {
		middleware.InitJWT()
	}

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
			auth.POST("/login_by_wechat", group.Auth.LoginByWechat)
			auth.POST("/decrypt_phone", middleware.MiniProgramAuth(), group.Auth.DecryptPhone)
			auth.POST("/logout", middleware.MiniProgramAuth(), group.Auth.Logout)
		}

		user := api.Group("/user")
		user.Use(middleware.MiniProgramAuth())
		{
			user.GET("/profile", group.User.GetProfile)
			user.PUT("/profile", group.User.UpdateProfile)
			user.POST("/avatar", group.User.UploadAvatar)
		}

		question := api.Group("/questions")
		question.Use(middleware.MiniProgramAuth())
		{
			question.GET("", group.Question.List)
			question.GET("/:id", group.Question.Get)
			question.GET("/random", group.Question.Random)
		}

		quiz := api.Group("/quiz")
		quiz.Use(middleware.MiniProgramAuth())
		{
			quiz.POST("/submit", group.Quiz.Submit)
			quiz.GET("/history", group.Quiz.History)
			quiz.GET("/wrong_questions", group.Quiz.WrongQuestions)
			quiz.POST("/add_wrong/:qid", group.Quiz.AddWrong)
		}

		score := api.Group("/score")
		score.Use(middleware.MiniProgramAuth())
		{
			score.GET("/my", group.Score.MyScore)
			score.GET("/ranking", group.Score.Ranking)
			score.GET("/stats", group.Score.Stats)
		}
	}

	admin := r.Group("/admin/api/v1")
	{
		admin.POST("/login", group.Admin.Login)
		admin.POST("/logout", middleware.AdminAuth(), group.Admin.Logout)

		adminProtected := admin.Group("")
		adminProtected.Use(middleware.AdminAuth())
		{
			adminProtected.GET("/dashboard", group.Admin.Dashboard)
			adminProtected.GET("/users", group.Admin.ListUsers)
			adminProtected.GET("/users/:id", group.Admin.GetUser)
			adminProtected.PUT("/users/:id/status", group.Admin.UpdateUserStatus)

			adminProtected.GET("/question_banks", group.Admin.ListQuestionBanks)
			adminProtected.POST("/question_banks", group.Admin.CreateQuestionBank)
			adminProtected.PUT("/question_banks/:id", group.Admin.UpdateQuestionBank)
			adminProtected.DELETE("/question_banks/:id", group.Admin.DeleteQuestionBank)

			adminProtected.GET("/questions", group.Admin.ListQuestions)
			adminProtected.POST("/questions", group.Admin.CreateQuestion)
			adminProtected.PUT("/questions/:id", group.Admin.UpdateQuestion)
			adminProtected.DELETE("/questions/:id", group.Admin.DeleteQuestion)
			adminProtected.POST("/questions/import", group.Admin.ImportQuestions)

			adminProtected.GET("/stats/overview", group.Admin.StatsOverview)
			adminProtected.GET("/stats/users", group.Admin.StatsUsers)
			adminProtected.GET("/stats/questions", group.Admin.StatsQuestions)
		}
	}

	zap.S().Info("Router initialized")
	return r
}
