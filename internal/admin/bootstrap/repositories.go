package bootstrap

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/repository"
	customerRepo "github.com/fathimasithara01/tradeverse/internal/customer/repository/customerrepo"

	"gorm.io/gorm"
)

type Repositories struct {
	User             repository.IUserRepository
	Role             repository.IRoleRepository
	Dashboard        repository.IDashboardRepository
	Permission       repository.IPermissionRepository
	Activity         repository.IActivityRepository
	SubscriptionPlan repository.ISubscriptionPlanRepository
	Subscription     repository.ISubscriptionRepository
	AdminWallet      repository.IAdminWalletRepository
	Signal           repository.ISignalRepository
	Transaction      repository.ITransactionRepository
	Commission       repository.ICommissionRepository
	WebConfig        repository.IWebConfigurationRepository

	CustomerSubscription *customerRepo.CustomerSubscriptionRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:                 repository.NewUserRepository(db),
		Role:                 repository.NewRoleRepository(db),
		Dashboard:            repository.NewDashboardRepository(db),
		Permission:           repository.NewPermissionRepository(db),
		Activity:             repository.NewActivityRepository(db),
		SubscriptionPlan:     repository.NewSubscriptionPlanRepository(db),
		Subscription:         repository.NewSubscriptionRepository(db),
		AdminWallet:          repository.NewAdminWalletRepository(db),
		Signal:               repository.NewSignalRepository(db),
		Transaction:          repository.NewTransactionRepository(db),
		Commission:           repository.NewCommissionRepository(db),
		WebConfig:            repository.NewWebConfigurationRepository(db),
		CustomerSubscription: customerRepo.NewCustomerSubscriptionRepository(db), // ← initialize

	}
}
