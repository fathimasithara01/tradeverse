package bootstrap

import (
	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	customerService "github.com/fathimasithara01/tradeverse/internal/customer/service"

	"gorm.io/gorm"
)

type Services struct {
	User                 service.IUserService
	Role                 service.IRoleService
	Dashboard            service.IDashboardService
	Permission           service.IPermissionService
	Activity             service.IActivityService
	SubscriptionPlan     service.ISubscriptionPlanService
	Subscription         service.ISubscriptionService
	AdminWallet          service.IAdminWalletService
	LiveSignal           service.ILiveSignalService
	Transaction          service.ITransactionService
	MarketData           service.IMarketDataService
	Commission           service.ICommissionService
	WebConfiguration     service.IWebConfigurationService
	CustomerSubscription *customerService.CustomerSubscriptionService
}

func InitServices(repos *Repositories, db *gorm.DB, cfg *config.Config) *Services {
	adminWalletService := service.NewAdminWalletService(repos.AdminWallet, db)

	customerSubService := customerService.NewCustomerSubscriptionService(
		repos.CustomerSubscription,
		repos.SubscriptionPlan,
		adminWalletService,
		repos.User,
		db,
	)

	return &Services{
		User:                 service.NewUserService(repos.User, repos.Role, cfg.JWT.Secret),
		Role:                 service.NewRoleService(repos.Role, repos.Permission, repos.User),
		Dashboard:            service.NewDashboardService(repos.Dashboard),
		Permission:           service.NewPermissionService(repos.Permission),
		Activity:             service.NewActivityService(repos.Activity),
		SubscriptionPlan:     service.NewSubscriptionPlanService(repos.SubscriptionPlan),
		AdminWallet:          adminWalletService,
		Subscription:         service.NewSubscriptionService(repos.Subscription, repos.SubscriptionPlan, repos.User, adminWalletService, db),
		LiveSignal:           service.NewLiveSignalService(repos.Signal),
		Transaction:          service.NewTransactionService(repos.Transaction),
		MarketData:           service.NewMarketDataService(),
		Commission:           service.NewCommissionService(repos.Commission, db),
		WebConfiguration:     service.NewWebConfigurationService(repos.WebConfig),
		CustomerSubscription: customerSubService,
	}
}
