//go:build wireinject
// +build wireinject

package wire

//go:generate go run github.com/google/wire/cmd/wire

import (
	"hub/internal/application/listener"
	"hub/internal/application/radio"
	appreaction "hub/internal/application/reaction"
	appshared "hub/internal/application/shared"
	"hub/internal/application/statistics"
	apptrack "hub/internal/application/track"
	"hub/internal/config"
	"hub/internal/database"
	domainreaction "hub/internal/domain/reaction"
	"hub/internal/domain/track"
	"hub/internal/infrastructure/cache"
	"hub/internal/infrastructure/events"
	"hub/internal/infrastructure/icecast"
	"hub/internal/infrastructure/metrics"
	"hub/internal/infrastructure/persistence/postgres"
	"hub/internal/infrastructure/scheduler"
	"hub/internal/interfaces/http/handler"
	"hub/internal/interfaces/http/server"
	"hub/internal/logger"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Application holds all dependencies for the serve command
type Application struct {
	Config    config.Config
	Logger    *logger.Logger
	Database  database.Database
	Server    *server.Server
	Scheduler scheduler.Scheduler
}

// MigrateApp holds dependencies for migrate commands
type MigrateApp struct {
	Config config.Config
	Logger *logger.Logger
	DSN    string
}

func ProvideConfig() config.Config {
	return config.NewConfig()
}

func ProvideLogger(cfg config.Config) *logger.Logger {
	return logger.NewLogger(cfg.LogLevel())
}

func ProvideDSN(cfg config.Config) string {
	dsn, _, _ := cfg.DatabaseConnection()
	return dsn
}

func ProvideDatabase(cfg config.Config, log *logger.Logger) database.Database {
	dsn, minConns, maxConns := cfg.DatabaseConnection()
	return database.NewDatabase(dsn, minConns, maxConns, log)
}

func ProvidePool(db database.Database) *pgxpool.Pool {
	return db.Pool()
}

func ProvideEventPublisher() appshared.EventPublisher {
	return events.NewInMemoryPublisher()
}

func ProvideCache(cfg config.Config, log *logger.Logger) (cache.Cache, error) {
	return cache.NewCache(cfg, log)
}

func ProvideRedisClient(c cache.Cache) *redis.Client {
	return c.Client()
}

func ProvideMetrics() *metrics.Metrics {
	return metrics.NewMetrics()
}

func ProvideUnitOfWork(pool *pgxpool.Pool) *postgres.UnitOfWork {
	return postgres.NewUnitOfWork(pool)
}

func ProvideTrackRepository(pool *pgxpool.Pool) *postgres.TrackRepository {
	return postgres.NewTrackRepository(pool)
}

func ProvideTrackDomainRepository(repo *postgres.TrackRepository) track.Repository {
	return repo
}

func ProvideReactionRepository(pool *pgxpool.Pool) domainreaction.Repository {
	return postgres.NewReactionRepository(pool)
}

func ProvideListenerRepository(pool *pgxpool.Pool) *postgres.ListenerRepository {
	return postgres.NewListenerRepository(pool)
}

func ProvideStatisticsRepository(pool *pgxpool.Pool) *postgres.StatisticsRepository {
	return postgres.NewStatisticsRepository(pool)
}

func ProvideListenerAdapter(repo *postgres.ListenerRepository) *postgres.ListenerAdapter {
	return postgres.NewListenerAdapter(repo)
}

func ProvideTrackListenerAdapter(repo *postgres.TrackRepository) *postgres.TrackListenerAdapter {
	return postgres.NewTrackListenerAdapter(repo)
}

func ProvideUpsertTrackHandler(repo track.Repository, pub appshared.EventPublisher) *apptrack.UpsertTrackHandler {
	return apptrack.NewUpsertTrackHandler(repo, pub)
}

func ProvideGetTrackHandler(repo track.Repository) *apptrack.GetTrackHandler {
	return apptrack.NewGetTrackHandler(repo)
}

func ProvideAddReactionHandler(rr domainreaction.Repository, tr track.Repository, pub appshared.EventPublisher) *appreaction.AddReactionHandler {
	return appreaction.NewAddReactionHandler(rr, tr, pub)
}

func ProvideCheckReactionHandler(rr domainreaction.Repository) *appreaction.CheckReactionHandler {
	return appreaction.NewCheckReactionHandler(rr)
}

func ProvideIcecastClient(cfg config.Config) (icecast.Client, error) {
	return icecast.NewClient(cfg)
}

func ProvideRadioService(ic icecast.Client) radio.Service {
	return radio.NewService(ic)
}

func ProvideStatisticsService(repo *postgres.StatisticsRepository) statistics.Service {
	return statistics.NewService(repo)
}

func ProvideListenerService(ic icecast.Client, la *postgres.ListenerAdapter, ta *postgres.TrackListenerAdapter, log *logger.Logger) listener.Service {
	return listener.NewService(ic, la, ta, log)
}

func ProvideTrackHandler(uh *apptrack.UpsertTrackHandler, gh *apptrack.GetTrackHandler) *handler.TrackHandler {
	return handler.NewTrackHandler(uh, gh)
}

func ProvideReactionHandler(ah *appreaction.AddReactionHandler, ch *appreaction.CheckReactionHandler) *handler.ReactionHandler {
	return handler.NewReactionHandler(ah, ch)
}

func ProvideRadioHandler(svc radio.Service) *handler.RadioHandler {
	return handler.NewRadioHandler(svc)
}

func ProvideStatisticsHandler(svc statistics.Service) *handler.StatisticsHandler {
	return handler.NewStatisticsHandler(svc)
}

func ProvideHealthHandler(pool *pgxpool.Pool, redisClient *redis.Client) *handler.HealthHandler {
	return handler.NewHealthHandler(pool, redisClient)
}

func ProvideRouter(th *handler.TrackHandler, rh *handler.ReactionHandler, rah *handler.RadioHandler, sh *handler.StatisticsHandler, hh *handler.HealthHandler) *server.Router {
	return server.NewRouter(th, rh, rah, sh, hh)
}

func ProvideServer(router *server.Router, log *logger.Logger) *server.Server {
	return server.NewServer(router, log)
}

func ProvideScheduler(ls listener.Service, log *logger.Logger) scheduler.Scheduler {
	return scheduler.NewScheduler(ls, log)
}

func ProvideApplication(cfg config.Config, log *logger.Logger, db database.Database, srv *server.Server, sched scheduler.Scheduler) *Application {
	return &Application{Config: cfg, Logger: log, Database: db, Server: srv, Scheduler: sched}
}

func ProvideMigrateApp(cfg config.Config, log *logger.Logger, dsn string) *MigrateApp {
	return &MigrateApp{Config: cfg, Logger: log, DSN: dsn}
}

var ProviderSet = wire.NewSet(
	ProvideConfig, ProvideLogger, ProvideDSN, ProvideDatabase, ProvidePool, ProvideEventPublisher,
	ProvideCache, ProvideRedisClient, ProvideMetrics, ProvideUnitOfWork,
	ProvideTrackRepository, ProvideTrackDomainRepository, ProvideReactionRepository,
	ProvideListenerRepository, ProvideStatisticsRepository,
	ProvideListenerAdapter, ProvideTrackListenerAdapter,
	ProvideUpsertTrackHandler, ProvideGetTrackHandler, ProvideAddReactionHandler, ProvideCheckReactionHandler,
	ProvideIcecastClient, ProvideRadioService, ProvideStatisticsService, ProvideListenerService,
	ProvideTrackHandler, ProvideReactionHandler, ProvideRadioHandler, ProvideStatisticsHandler, ProvideHealthHandler,
	ProvideRouter, ProvideServer, ProvideScheduler, ProvideApplication,
)

var MigrateProviderSet = wire.NewSet(ProvideConfig, ProvideLogger, ProvideDSN, ProvideMigrateApp)

func InitializeApp() (*Application, func(), error) {
	wire.Build(ProviderSet)
	return nil, nil, nil
}

func InitializeMigrateApp() (*MigrateApp, error) {
	wire.Build(MigrateProviderSet)
	return nil, nil
}
