package cardinal_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/P0H1ng/Cardinal/conf"
	"github.com/P0H1ng/Cardinal/internal/asteroid"
	"github.com/P0H1ng/Cardinal/internal/bootstrap"
	"github.com/P0H1ng/Cardinal/internal/db"
	"github.com/P0H1ng/Cardinal/internal/dynamic_config"
	"github.com/P0H1ng/Cardinal/internal/game"
	"github.com/P0H1ng/Cardinal/internal/livelog"
	"github.com/P0H1ng/Cardinal/internal/misc/webhook"
	"github.com/P0H1ng/Cardinal/internal/route"
	"github.com/P0H1ng/Cardinal/internal/store"
	"github.com/P0H1ng/Cardinal/internal/timer"
	"github.com/P0H1ng/Cardinal/internal/utils"
	log "unknwon.dev/clog/v2"
)

var managerToken = utils.GenerateToken()

var checkToken string

var team = make([]struct {
	Name      string `json:"Name"`
	Password  string `json:"Password"`
	Token     string `json:"token"`
	AccessKey string `json:"access_key"`
}, 0)

var router *gin.Engine

func TestMain(m *testing.M) {
	prepare()
	log.Trace("Cardinal Test ready...")
	m.Run()

	os.Exit(0)
}

func prepare() {
	_ = log.NewConsole(100)
	
	log.Trace("Prepare for Cardinal test environment...")

	gin.SetMode(gin.ReleaseMode)

	conf.Init()

	// Init MySQL database.
	db.InitMySQL()

	// Test manager account e99:qwe1qwe2qwe3
	db.MySQL.Create(&db.Manager{
		Name:     "e99",
		Password: utils.AddSalt("qwe1qwe2qwe3"),
		Token:    managerToken,
		IsCheck:  false,
	})

	// Refresh the dynamic config from the database.
	dynamic_config.Init()

	bootstrap.GameToTimerBridge()
	timer.Init()

	asteroid.Init(game.AsteroidGreetData)

	// Cache
	store.Init()
	webhook.RefreshWebHookStore()

	// Live log
	livelog.Init()

	// Web router.
	router = route.Init()
}
