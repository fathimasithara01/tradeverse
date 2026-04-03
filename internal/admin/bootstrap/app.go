package bootstrap

import (
	"context"
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/database"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	Port   string
}

func InitializeApp(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.ConnectDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := database.RunMigrations(db); err != nil {
		return nil, err
	}

	repos := InitRepositories(db)
	services := InitServices(repos, db, cfg)
	r := InitRouter(services, cfg, db)
	SetupTemplatesAndStatic(r)
	InitCron(services, db)

	return &App{
		engine: r,
		Port:   cfg.Server.AdminPort,
	}, nil
}
func (a *App) Engine() *gin.Engine {
	return a.engine
}

func (a *App) Run() error {
	log.Printf("Server starting on http://localhost:%s", a.Port)
	return a.engine.Run(":" + a.Port)
}
