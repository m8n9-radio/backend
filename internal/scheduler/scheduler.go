package scheduler

import (
	"context"
	"hub/internal/config"
	"hub/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Job interface {
	Execute(ctx context.Context) error
}

type JobFactory func(logger *logger.Logger, pool *pgxpool.Pool) Job

type jobDefinition struct {
	name     string
	schedule string
	factory  JobFactory
}

var registeredJobs []jobDefinition

func RegisterJob(name, schedule string, factory JobFactory) {
	registeredJobs = append(registeredJobs, jobDefinition{
		name:     name,
		schedule: schedule,
		factory:  factory,
	})
}

type (
	Scheduler interface {
		Start(ctx context.Context) error
	}

	scheduler struct {
		cron   *cron.Cron
		logger *logger.Logger
		config config.Config
		pool   *pgxpool.Pool
		ctx    context.Context
	}
)

func NewScheduler(cfg config.Config, logger *logger.Logger, pool *pgxpool.Pool) Scheduler {
	return &scheduler{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger),
			)),
		logger: logger,
		config: cfg,
		pool:   pool,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	s.ctx = ctx

	if !s.config.SchedulerEnabled() {
		s.logger.Info("scheduler is disabled")
		<-ctx.Done()
		return ctx.Err()
	}

	s.registerAllJobs()

	s.cron.Start()
	s.logger.Info("scheduler started")

	<-ctx.Done()

	stopCtx := s.cron.Stop()
	<-stopCtx.Done()

	s.logger.Info("scheduler stopped gracefully")
	return ctx.Err()
}

func (s *scheduler) registerAllJobs() {
	if len(registeredJobs) == 0 {
		s.logger.Warn("no jobs registered - add job files in internal/scheduler/jobs/")
		return
	}

	for _, jobDef := range registeredJobs {
		job := jobDef.factory(s.logger, s.pool)
		name := jobDef.name
		schedule := jobDef.schedule

		_, err := s.cron.AddFunc(schedule, func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Errorf("[%s] panic: %v", name, r)
				}
			}()

			s.logger.Infof("[%s] executing", name)

			if err := job.Execute(s.ctx); err != nil {
				s.logger.WithError(err).Errorf("[%s] failed", name)
			}
		})

		if err != nil {
			s.logger.WithError(err).Errorf("failed to register job: %s", name)
		} else {
			s.logger.Infof("registered job: %s (schedule: %s)", name, schedule)
		}
	}
}
