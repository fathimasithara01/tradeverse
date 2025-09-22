package bootstrap

import (
	"log"

	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	adminService "github.com/fathimasithara01/tradeverse/internal/admin/service"

	"github.com/fathimasithara01/tradeverse/internal/trader/controllers"
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/internal/trader/router"
	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/migrations"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	port   string
}

func InitializeApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := migrations.ConnectDB(*cfg)
	if err != nil {
		return nil, err
	}

	userRepo := adminRepo.NewUserRepository(db)
	roleRepo := adminRepo.NewRoleRepository(db)
	userService := adminService.NewUserService(userRepo, roleRepo, cfg.JWTSecret)

	authController := controllers.NewAuthController(userService)

	tradeRepo := repository.NewTradeRepository(db)
	subRepo := repository.NewSubscriberRepository(db)
	liveRepo := repository.NewLiveTradeRepository(db)

	tradeService := service.NewTradeService(tradeRepo)
	subService := service.NewSubscriberService(subRepo)
	liveService := service.NewLiveTradeService(liveRepo)

	tradeController := controllers.NewTradeController(tradeService)
	subController := controllers.NewSubscriberController(subService)
	liveController := controllers.NewLiveTradeController(liveService)

	r := router.SetupRouter(cfg, authController, tradeController, subController, liveController)

	return &App{
		engine: r,
		port:   cfg.TraderPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Customer API server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
