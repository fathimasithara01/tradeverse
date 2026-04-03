package bootstrap

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
)

// Controllers holds all admin controllers
type Controllers struct {
	Auth             *controllers.AuthController
	User             *controllers.UserController
	Role             *controllers.RoleController
	Dashboard        *controllers.DashboardController
	Permission       *controllers.PermissionController
	Activity         *controllers.ActivityController
	AdminWallet      *controllers.AdminWalletController
	Subscription     *controllers.SubscriptionController
	Signal           *controllers.SignalController
	Transaction      *controllers.TransactionController
	Commission       *controllers.CommissionController
	WebConfiguration *controllers.WebConfigurationController
}

func InitControllers(svc *Services) *Controllers {
	return &Controllers{
		Auth:             controllers.NewAuthController(svc.User),
		User:             controllers.NewUserController(svc.User),
		Role:             controllers.NewRoleController(svc.Role),
		Dashboard:        controllers.NewDashboardController(svc.Dashboard, svc.MarketData),
		Permission:       controllers.NewPermissionController(svc.Permission, svc.Role),
		Activity:         controllers.NewActivityController(svc.Activity),
		AdminWallet:      controllers.NewAdminWalletController(svc.AdminWallet),
		Subscription:     controllers.NewSubscriptionController(svc.Subscription, svc.SubscriptionPlan),
		Signal:           controllers.NewSignalController(svc.LiveSignal),
		Transaction:      controllers.NewTransactionController(svc.Transaction),
		Commission:       controllers.NewCommissionController(svc.Commission),
		WebConfiguration: controllers.NewWebConfigurationController(svc.WebConfiguration),
	}
}
