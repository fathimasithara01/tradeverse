package cron

import (
	"context"
	"log"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/robfig/cron/v3"
)

func StartSignalCronJobs(signalService service.ISignalService) {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		log.Println("Cron: Updating Pending Signals Current Price...")
		if err := signalService.UpdatePendingSignalsCurrentPrice(context.Background()); err != nil {
			log.Printf("Error updating pending signals: %v", err)
		}
	})

	c.AddFunc("@every 30s", func() {
		log.Println("Cron: Checking Active Signals Status...")
		if err := signalService.UpdateActiveSignalStatuses(context.Background()); err != nil {
			log.Printf("Error updating active signal statuses: %v", err)
		}
	})

	c.Start()
	log.Println("Signal cron jobs started.")
}
