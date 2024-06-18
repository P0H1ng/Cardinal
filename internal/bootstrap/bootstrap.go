package bootstrap

import (
	log "unknwon.dev/clog/v2"

	"github.com/P0H1ng/Cardinal/conf"
	"github.com/P0H1ng/Cardinal/internal/asteroid"
	"github.com/P0H1ng/Cardinal/internal/db"
	"github.com/P0H1ng/Cardinal/internal/dynamic_config"
	"github.com/P0H1ng/Cardinal/internal/game"
	"github.com/P0H1ng/Cardinal/internal/install"
	"github.com/P0H1ng/Cardinal/internal/livelog"
	"github.com/P0H1ng/Cardinal/internal/misc"
	"github.com/P0H1ng/Cardinal/internal/misc/webhook"
	"github.com/P0H1ng/Cardinal/internal/route"
	"github.com/P0H1ng/Cardinal/internal/store"
	"github.com/P0H1ng/Cardinal/internal/timer"
)

func init() {
	// Init log
	_ = log.NewConsole(100)
}

// LinkStart starts the Cardinal.
func LinkStart() {
	// Install
	install.Init()

	// Config
	conf.Init()

	// Check version
	misc.CheckVersion()

	// Sentry
	misc.Sentry()

	// Init MySQL database.
	db.InitMySQL()

	// Check manager
	install.InitManager()

	// Refresh the dynamic config from the database.
	dynamic_config.Init()

	// Check if the database need update.
	misc.CheckDatabaseVersion()

	// Game timer.
	GameToTimerBridge()
	timer.Init()

	// Cache
	store.Init()
	webhook.RefreshWebHookStore()

	// Unity3D Asteroid
	asteroid.Init(game.AsteroidGreetData)

	// Live log
	livelog.Init()

	// Web router.
	router := route.Init()

	log.Fatal("Failed to start web server: %v", router.Run(conf.Get().Port))
}
