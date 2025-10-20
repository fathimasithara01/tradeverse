package bootstrap

import (
	"log"

	"github.com/fathimasithara01/tradeverse/config"
	adminRepo "github.com/fathimasithara01/tradeverse/internal/admin/repository"
	adminService "github.com/fathimasithara01/tradeverse/internal/admin/service"

	"github.com/fathimasithara01/tradeverse/internal/trader/controllers"
	"github.com/fathimasithara01/tradeverse/internal/trader/cron"
	"github.com/fathimasithara01/tradeverse/internal/trader/repository"
	"github.com/fathimasithara01/tradeverse/internal/trader/router"
	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/migrations"
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
	userService := adminService.NewUserService(userRepo, roleRepo, cfg.JWT.Secret)

	authController := controllers.NewAuthController(userService)

	tradeSignlRepo := repository.NewSignalRepository(db)
	profileRepo := repository.NewTraderProfileRepository(db)
	walletrepo := repository.NewGormWalletRepository(db)
	subRepo := repository.NewSubscriberRepository(db)
	liveRepo := repository.NewLiveTradeRepository(db)
	traderSubsRepo := repository.NewTraderSubscriptionRepository(db)

	subService := service.NewSubscriberService(subRepo)
	liveService := service.NewLiveTradeService(liveRepo)
	profileService := service.NewTraderProfileService(profileRepo)
	walletService := service.NewWalletService(walletrepo)
	tradeSignlService := service.NewSignalService(tradeSignlRepo)
	traderSubsService := service.NewTraderSubscriptionService(traderSubsRepo, db)

	subController := controllers.NewSubscriberController(subService)
	liveController := controllers.NewLiveTradeController(liveService)
	profileController := controllers.NewTraderProfileController(profileService)
	walletController := controllers.NewWalletController(walletService)
	tradeSignlController := controllers.NewSignalController(tradeSignlService)
	traderSubsController := controllers.NewTraderSubscriptionController(traderSubsService)

	marketDataRepo := repository.NewMarketDataRepository(db)
	marketDataService := service.NewMarketDataService(marketDataRepo)
	marketDataHandler := controllers.NewMarketDataHandler(marketDataService)

	r := router.SetupRouter(cfg, authController, profileController, walletController, subController, liveController, tradeSignlController, marketDataHandler, traderSubsController)

	cron.StartSignalCronJobs(service.NewSignalService(repository.NewSignalRepository(db)))

	return &App{
		engine: r,
		port:   cfg.Server.TraderPort,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Customer API server starting on http://localhost:%s", a.port)
	return a.engine.Run(":" + a.port)
}
