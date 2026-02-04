package scheduler

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/berkaycubuk/subtrack/internal/services"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron     *cron.Cron
	subSvc   *services.SubscriptionService
	stopChan chan struct{}
}

func NewScheduler(subSvc *services.SubscriptionService) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		subSvc:   subSvc,
		stopChan: make(chan struct{}),
	}
}

func (s *Scheduler) Start() error {
	log.Println("Starting SubTrack scheduler...")

	if _, err := s.cron.AddFunc("0 0 9,21 * * *", s.runCheck); err != nil {
		return err
	}

	s.cron.Start()
	log.Println("Scheduler started (runs at 9:00 AM and 9:00 PM)")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Println("Received shutdown signal...")
	case <-s.stopChan:
		log.Println("Received stop signal...")
	}

	s.Stop()
	return nil
}

func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) runCheck() {
	log.Println("Running subscription check...")

	if err := s.subSvc.UpdatePastDuePayments(); err != nil {
		log.Printf("Error updating past due payments: %v", err)
	}

	subs, err := s.subSvc.CheckUpcomingPayments()
	if err != nil {
		log.Printf("Error checking upcoming payments: %v", err)
		return
	}

	if len(subs) == 0 {
		log.Println("No upcoming payments to notify")
		return
	}

	log.Printf("Found %d subscriptions with upcoming payments", len(subs))

	if err := s.subSvc.SendNotifications(subs); err != nil {
		log.Printf("Error sending notifications: %v", err)
	}
}

func (s *Scheduler) StopGracefully() {
	close(s.stopChan)
}
