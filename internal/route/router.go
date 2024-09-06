package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/P0H1ng/Cardinal/conf"
	"github.com/P0H1ng/Cardinal/frontend"
	"github.com/P0H1ng/Cardinal/internal/asteroid"
	"github.com/P0H1ng/Cardinal/internal/auth"
	"github.com/P0H1ng/Cardinal/internal/auth/manager"
	"github.com/P0H1ng/Cardinal/internal/auth/team"
	"github.com/P0H1ng/Cardinal/internal/bulletin"
	"github.com/P0H1ng/Cardinal/internal/container"
	"github.com/P0H1ng/Cardinal/internal/dynamic_config"
	"github.com/P0H1ng/Cardinal/internal/game"
	"github.com/P0H1ng/Cardinal/internal/healthy"
	"github.com/P0H1ng/Cardinal/internal/livelog"
	"github.com/P0H1ng/Cardinal/internal/locales"
	"github.com/P0H1ng/Cardinal/internal/logger"
	"github.com/P0H1ng/Cardinal/internal/misc/webhook"
	"github.com/P0H1ng/Cardinal/internal/timer"
	"github.com/P0H1ng/Cardinal/internal/upload"
	"github.com/P0H1ng/Cardinal/internal/utils"
)

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	// CORS Header
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{"Authorization", "Content-type", "User-Agent"},
		AllowOrigins: []string{"*"},
	}))

	api := r.Group("/api")
	api.Use(locales.Middleware()) // i18n
	// Sentry
	// if conf.Get().Sentry {
	// 	api.Use(sentrygin.New(sentrygin.Options{
	// 		Repanic: true,
	// 	}))
	// }

	// Frontend
	if !conf.Get().SeparateFrontend {
		r.Use(static.Serve("/", frontend.FS()))
	}

	// Cardinal basic info
	api.Any("/", func(c *gin.Context) {
		c.JSON(utils.MakeSuccessJSON("Cardinal"))
	})

	api.GET("/base", func(c *gin.Context) {
		c.JSON(utils.MakeSuccessJSON(gin.H{
			"Title":    dynamic_config.Get(utils.TITLE_CONF),
			"Language": dynamic_config.Get(utils.DEFAULT_LANGUAGE),
		}))
	})

	api.GET("/time", __(timer.GetTime))

	// Static files
	api.Static("/uploads", "./uploads")

	// Team login
	api.POST("/login", __(team.TeamLogin))

	// Team logout
	api.GET("/logout", __(team.TeamLogout))

	// Live log
	api.GET("/livelog", livelog.GlobalStreamHandler)

	// Submit flag
	api.POST("/flag", __(game.SubmitFlag))

	// Asteroid websocket
	api.GET("/asteroid", func(c *gin.Context) {
		asteroid.ServeWebSocket(c)
	})

	// For team
	teamRouter := api.Group("/team")
	teamRouter.Use(auth.TeamAuthRequired())
	{
		teamRouter.GET("/info", __(team.GetTeamInfo))
		teamRouter.GET("/gameboxes", __(game.GetSelfGameBoxes))
		teamRouter.GET("/gameboxes/all", __(game.GetOthersGameBox))
		teamRouter.GET("/rank", func(c *gin.Context) {
			c.JSON(utils.MakeSuccessJSON(gin.H{"Title": game.GetRankListTitle(), "Rank": game.GetRankList()}))
		})
		teamRouter.GET("/bulletins", __(bulletin.GetAllBulletins))
		teamRouter.POST("/start", __(asteroid.Start))
	}

	// Manager login
	api.POST("/manager/login", __(manager.ManagerLogin))

	// Manager logout
	api.GET("/manager/logout", __(manager.ManagerLogout))

	// For manager
	check := api.Group("/manager").Use(auth.AdminAuthRequired())
	managerRouter := api.Group("/manager").Use(auth.AdminAuthRequired(), auth.ManagerRequired())
	{
		// Challenge
		managerRouter.GET("/challenges", __(game.GetAllChallenges))
		managerRouter.POST("/challenge", __(game.NewChallenge))
		managerRouter.PUT("/challenge", __(game.EditChallenge))
		managerRouter.DELETE("/challenge", __(game.DeleteChallenge))
		managerRouter.POST("/challenge/visible", __(game.SetVisible))

		// GameBox
		managerRouter.GET("/gameboxes", __(game.GetGameBoxes))
		managerRouter.POST("/gameboxes", __(game.NewGameBoxes))
		managerRouter.PUT("/gamebox", __(game.EditGameBox))
		managerRouter.GET("/gameboxes/sshTest", __(game.TestAllSSH))
		managerRouter.POST("/gameboxes/sshTest", __(game.TestSSH))
		managerRouter.GET("/gameboxes/refreshFlag", func(c *gin.Context) {
			game.RefreshFlag()
			// TODO: i18n
			c.JSON(utils.MakeSuccessJSON("更新 Flag 操作已執行，請在系統狀態查看是否有錯誤訊息"))
		})
		managerRouter.GET("/gameboxes/reset", __(game.ResetAllGameBoxes))

		// Team
		managerRouter.GET("/teams", __(team.GetAllTeams))
		managerRouter.POST("/teams", __(team.NewTeams))
		managerRouter.PUT("/team", __(team.EditTeam))
		managerRouter.DELETE("/team", __(team.DeleteTeam))
		managerRouter.POST("/team/resetPassword", __(team.ResetTeamPassword))

		// Manager
		managerRouter.GET("/managers", __(manager.GetAllManager))
		managerRouter.POST("/manager", __(manager.NewManager))
		managerRouter.GET("/manager/token", __(manager.RefreshManagerToken))
		managerRouter.GET("/manager/changePassword", __(manager.ChangeManagerPassword))
		managerRouter.DELETE("/manager", __(manager.DeleteManager))

		// Flag
		managerRouter.GET("/flags", __(game.GetFlags))
		managerRouter.POST("/flag/generate", __(game.GenerateFlag))
		managerRouter.GET("/flag/export", __(game.ExportFlag))

		// Asteroid Unity3D
		managerRouter.GET("/asteroid/status", __(asteroid.GetAsteroidStatus))
		managerRouter.POST("/asteroid/attack", __(asteroid.Attack))
		managerRouter.POST("/asteroid/rank", __(asteroid.Rank))
		managerRouter.POST("/asteroid/status", __(asteroid.Status))
		managerRouter.POST("/asteroid/round", __(asteroid.Round))
		managerRouter.POST("/asteroid/easterEgg", __(asteroid.EasterEgg))
		managerRouter.POST("/asteroid/time", __(asteroid.Time))
		managerRouter.POST("/asteroid/clear", __(asteroid.Clear))
		managerRouter.POST("/asteroid/clearAll", __(asteroid.ClearAll))

		// Check
		check.POST("/checkDown", __(game.CheckDown))

		// Bulletin
		managerRouter.GET("/bulletins", __(bulletin.GetAllBulletins))
		managerRouter.POST("/bulletin", __(bulletin.NewBulletin))
		managerRouter.PUT("/bulletin", __(bulletin.EditBulletin))
		managerRouter.DELETE("/bulletin", __(bulletin.DeleteBulletin))

		// File
		managerRouter.POST("/uploadPicture", __(upload.UploadPicture))
		managerRouter.GET("/dir", __(upload.GetDir))

		// Docker
		managerRouter.POST("/docker/findImage", __(container.GetImageData))

		// Log
		managerRouter.GET("/logs", __(logger.GetLogs))
		managerRouter.GET("/rank", func(c *gin.Context) {
			c.JSON(utils.MakeSuccessJSON(gin.H{"Title": game.GetRankListTitle(), "Rank": game.GetManagerRankList()}))
		})
		managerRouter.GET("/panel", __(healthy.Panel))

		// WebHook
		managerRouter.GET("/webhooks", __(webhook.GetWebHook))
		managerRouter.POST("/webhook", __(webhook.NewWebHook))
		managerRouter.PUT("/webhook", __(webhook.EditWebHook))
		managerRouter.DELETE("/webhook", __(webhook.DeleteWebHook))

		// Config
		managerRouter.GET("/configs", __(dynamic_config.GetAllConfig))
		managerRouter.GET("/config", __(dynamic_config.GetConfig))
		managerRouter.PUT("/config", __(dynamic_config.SetConfig))
	}

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(utils.MakeErrJSON(404, 40400,
			locales.I18n.T(c.GetString("lang"), "general.not_found"),
		))
	})

	// 405
	r.NoMethod(func(c *gin.Context) {
		c.JSON(utils.MakeErrJSON(405, 40500,
			locales.I18n.T(c.GetString("lang"), "general.method_not_allow"),
		))
	})

	return r
}
