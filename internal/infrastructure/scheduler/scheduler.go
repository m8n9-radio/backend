package scheduler

import (
	"context"
	"log"
	"sync/atomic"

	"hub/internal/delivery/http/service"

	"github.com/robfig/cron/v3"
)

type (
	Scheduler interface {
		Start()
	}

	scheduler struct {
		cron            *cron.Cron
		listenerService service.ListenerService
		isRunning       atomic.Bool
	}
)

func NewScheduler(listenerService service.ListenerService) Scheduler {
	return &scheduler{
		cron:            cron.New(cron.WithSeconds()),
		listenerService: listenerService,
	}
}

func (s *scheduler) Start() {
	// Track listeners every 10 seconds
	_, err := s.cron.AddFunc("*/10 * * * * *", func() {
		// Skip if previous job is still running
		if !s.isRunning.CompareAndSwap(false, true) {
			log.Println("Skipping listener tracking - previous job still running")
			return
		}
		defer s.isRunning.Store(false)

		ctx := context.Background()
		if err := s.listenerService.TrackCurrentListeners(ctx); err != nil {
			log.Printf("Failed to track listeners: %v", err)
		}
	})

	if err != nil {
		log.Printf("Failed to add cron job: %v", err)
		return
	}

	s.cron.Start()
	log.Println("Scheduler started - tracking listeners every 10 seconds")
}
