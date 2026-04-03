package bootstrap

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/admin/cron"
	"gorm.io/gorm"
)

func InitCron(s *Services, db *gorm.DB) {
	cron.StartCronJobs(
		s.Subscription,
		s.CustomerSubscription,
		s.LiveSignal,
		db,
	)
	log.Println("[Bootstrap] Cron jobs initialized")
}
