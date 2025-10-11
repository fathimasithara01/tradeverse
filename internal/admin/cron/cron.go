package cron

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"                    // For admin subscription cron
	customerService "github.com/fathimasithara01/tradeverse/internal/customer/service" // For customer trader subscription cron
	"github.com/robfig/cron/v3"
)

type SubscriptionCronJob struct {
	SubscriptionService service.ISubscriptionService
}

func NewSubscriptionCronJob(subService service.ISubscriptionService) *SubscriptionCronJob {
	return &SubscriptionCronJob{
		SubscriptionService: subService,
	}
}

func (j *SubscriptionCronJob) Run() {
	err := j.SubscriptionService.DeactivateExpiredSubscriptions()
	if err != nil {
		log.Printf("CRON JOB FAILED (DeactivateExpiredSubscriptions): %v", err)
	}
}

type TraderSubscriptionCronJob struct {
	CustomerService customerService.AdminSubscriptionService
}

func NewTraderSubscriptionCronJob(custService customerService.AdminSubscriptionService) *TraderSubscriptionCronJob {
	return &TraderSubscriptionCronJob{
		CustomerService: custService,
	}
}

func (j *TraderSubscriptionCronJob) Run() {
	err := j.CustomerService.DeactivateExpiredTraderSubscriptions()
	if err != nil {
		log.Printf("CRON JOB FAILED (DeactivateExpiredTraderSubscriptions): %v", err)
	}
}

func StartCronJob(adminSubService service.ISubscriptionService, custSubService customerService.AdminSubscriptionService) {
	c := cron.New() // Create a new cron scheduler

	_, err := c.AddJob("* * * * *", NewSubscriptionCronJob(adminSubService)) // Runs every minute
	if err != nil {
		log.Fatalf("Error scheduling Admin Subscription Deactivation cron job: %v", err)
	}
	log.Println("Admin Subscription Deactivation cron job scheduled (running every minute for testing)...")

	_, err = c.AddJob("* * * * *", NewTraderSubscriptionCronJob(custSubService)) // Runs every minute
	if err != nil {
		log.Fatalf("Error scheduling Trader Subscription Deactivation cron job: %v", err)
	}
	log.Println("Trader Subscription Deactivation cron job scheduled (running every minute for testing)...")

	log.Println("All cron jobs started...")
	c.Start() // Start the cron scheduler in a non-blocking way
}
