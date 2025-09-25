package cron

import (
	"log"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"                    // For admin subscription cron
	customerService "github.com/fathimasithara01/tradeverse/internal/customer/service" // For customer trader subscription cron
	"github.com/robfig/cron/v3"
)

// SubscriptionCronJob encapsulates the admin subscription service for cron tasks
type SubscriptionCronJob struct {
	SubscriptionService service.ISubscriptionService
}

// NewSubscriptionCronJob creates a new SubscriptionCronJob instance
func NewSubscriptionCronJob(subService service.ISubscriptionService) *SubscriptionCronJob {
	return &SubscriptionCronJob{
		SubscriptionService: subService,
	}
}

// Run for SubscriptionCronJob
func (j *SubscriptionCronJob) Run() {
	err := j.SubscriptionService.DeactivateExpiredSubscriptions()
	if err != nil {
		log.Printf("CRON JOB FAILED (DeactivateExpiredSubscriptions): %v", err)
	}
}

// TraderSubscriptionCronJob encapsulates the customer trader subscription service for cron tasks
type TraderSubscriptionCronJob struct {
	CustomerService customerService.CustomerService
}

// NewTraderSubscriptionCronJob creates a new TraderSubscriptionCronJob instance
func NewTraderSubscriptionCronJob(custService customerService.CustomerService) *TraderSubscriptionCronJob {
	return &TraderSubscriptionCronJob{
		CustomerService: custService,
	}
}

// Run for TraderSubscriptionCronJob
func (j *TraderSubscriptionCronJob) Run() {
	err := j.CustomerService.DeactivateExpiredTraderSubscriptions()
	if err != nil {
		log.Printf("CRON JOB FAILED (DeactivateExpiredTraderSubscriptions): %v", err)
	}
}

// StartCronJobs initializes and starts all cron jobs
func StartCronJobs(adminSubService service.ISubscriptionService, custSubService customerService.CustomerService) {
	c := cron.New() // Create a new cron scheduler

	// Schedule the Admin Subscription Deactivation job
	// --- TEMPORARY CRON JOB TESTING OVERRIDE ---
	// This example runs every minute for testing.
	// For production, change this back to a sensible schedule like "0 2 * * *" (daily at 2 AM)
	_, err := c.AddJob("* * * * *", NewSubscriptionCronJob(adminSubService)) // Runs every minute
	if err != nil {
		log.Fatalf("Error scheduling Admin Subscription Deactivation cron job: %v", err)
	}
	// --- END TEMPORARY OVERRIDE ---
	log.Println("Admin Subscription Deactivation cron job scheduled (running every minute for testing)...")

	// Schedule the Trader Subscription Deactivation job
	// --- TEMPORARY CRON JOB TESTING OVERRIDE ---
	// This example runs every minute for testing.
	// For production, change this back to a sensible schedule like "0 2 * * *" (daily at 2 AM)
	_, err = c.AddJob("* * * * *", NewTraderSubscriptionCronJob(custSubService)) // Runs every minute
	if err != nil {
		log.Fatalf("Error scheduling Trader Subscription Deactivation cron job: %v", err)
	}
	// --- END TEMPORARY OVERRIDE ---
	log.Println("Trader Subscription Deactivation cron job scheduled (running every minute for testing)...")

	log.Println("All cron jobs started...")
	c.Start() // Start the cron scheduler in a non-blocking way
}
