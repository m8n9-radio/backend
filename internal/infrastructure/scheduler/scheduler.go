package scheduler

import (
	"context"
	"sync/atomic"

	"hub/internal/application/listener"
	"hub/internal/logger"

	"github.com/robfig/cron/v3"
)

type (
	Scheduler interface {
		Start()
		Stop(ctx context.Context) error
	}

	scheduler struct {
		cron            *cron.Cron
		listenerService listener.Service
		logger          *logger.Logger
		isRunning       atomic.Bool
		isStarted       atomic.Bool
	}
)

func NewScheduler(listenerService listener.Service, log *logger.Logger) Scheduler {
	return &scheduler{
		cron:            cron.New(cron.WithSeconds()),
		listenerService: listenerService,
		logger:          log,
	}
}

func (s *scheduler) Start() {
	if !s.isStarted.CompareAndSwap(false, true) {
		s.logger.Warn("Scheduler already started - preventing duplicate jobs")
		return
	}

	_, err := s.cron.AddFunc("*/3 * * * * *", func() {
		if !s.isRunning.CompareAndSwap(false, true) {
			s.logger.Debug("Skipping listener tracking - previous job still running")
			return
		}
		defer s.isRunning.Store(false)

		ctx := context.Background()
		if err := s.listenerService.TrackCurrentListeners(ctx); err != nil {
			s.logger.Errorf("Failed to track listeners: %v", err)
		}
	})

	if err != nil {
		s.logger.Errorf("Failed to add cron job: %v", err)
		s.isStarted.Store(false)
		return
	}

	s.cron.Start()
	s.logger.Info("Scheduler started - tracking listeners every 3 seconds")
}

func (s *scheduler) Stop(ctx context.Context) error {
	if !s.isStarted.Load() {
		s.logger.Warn("Scheduler is not running")
		return nil
	}

	s.logger.Info("Stopping scheduler...")

	stopCtx := s.cron.Stop()

	select {
	case <-stopCtx.Done():
		s.isStarted.Store(false)
		s.logger.Info("Scheduler stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("Scheduler stop timed out")
		return ctx.Err()
	}
}
