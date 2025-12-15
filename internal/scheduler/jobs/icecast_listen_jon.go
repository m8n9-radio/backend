package jobs

import (
	"context"
	"hub/internal/logger"
	"hub/internal/scheduler"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	scheduler.RegisterJob(
		"icecast-listeners",
		"*/3 * * * * *",
		NewIcecastListenJob,
	)
}

type IcecastListenJob struct {
	logger *logger.Logger
	pool   *pgxpool.Pool
}

func NewIcecastListenJob(logger *logger.Logger, pool *pgxpool.Pool) scheduler.Job {
	return &IcecastListenJob{
		logger: logger,
		pool:   pool,
	}
}

func (j *IcecastListenJob) Execute(ctx context.Context) error {
	start := time.Now()

	//
	//code here
	//

	j.logger.Infof("finis in %d", time.Since(start))
	return nil
}
