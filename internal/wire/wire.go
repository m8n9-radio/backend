//go:build wireinject
// +build wireinject

package wire

//go:generate go run github.com/google/wire/cmd/wire

import (
	"hub/internal/config"
	"hub/internal/database"
	"hub/internal/delivery/http/handler"
	"hub/internal/delivery/http/repository"
	"hub/internal/delivery/http/server"
	"hub/internal/delivery/http/service"
	"hub/internal/infrastructure/icecast"
	"hub/internal/infrastructure/scheduler"
	"hub/internal/logger"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Application holds all dependencies for the serve command
type Application struct {
	Config    config.Config
	Logger    *logger.Logger
	Database  database.Database
	Server    server.Server
	Scheduler scheduler.Scheduler
}

// MigrateApp holds dependencies for migrate commands
type MigrateApp struct {
	Config config.Config
	Logger *logger.Logger
	DSN    string
}

// ProvideConfig creates a new Config instance
func ProvideConfig() config.Config {
	return config.NewConfig()
}

// ProvideLogger creates a new Logger instance
func ProvideLogger(cfg config.Config) *logger.Logger {
	return logger.NewLogger(cfg.LogLevel())
}

// ProvideDSN extracts the database connection string from config
func ProvideDSN(cfg config.Config) string {
	dsn, _, _ := cfg.DatabaseConnection()
	return dsn
}

// ProvideDatabase creates a new Database instance
func ProvideDatabase(cfg config.Config, log *logger.Logger) database.Database {
	dsn, minConns, maxConns := cfg.DatabaseConnection()
	return database.NewDatabase(dsn, minConns, maxConns, log)
}

// ProvidePool extracts the pgxpool.Pool from Database
func ProvidePool(db database.Database) *pgxpool.Pool {
	return db.Pool()
}

// ProvideTrackRepository creates a new TrackRepository
func ProvideTrackRepository(pool *pgxpool.Pool) repository.TrackRepository {
	return repository.NewTrackRepository(pool)
}

// ProvideTrackService creates a new TrackService
func ProvideTrackService(repo repository.TrackRepository) service.TrackService {
	return service.NewTrackService(repo)
}

// ProvideTrackHandler creates a new TrackHandler
func ProvideTrackHandler(svc service.TrackService) handler.TrackHandler {
	return handler.NewTrackHandler(svc)
}

// ProvideReactionRepository creates a new ReactionRepository
func ProvideReactionRepository(pool *pgxpool.Pool) repository.ReactionRepository {
	return repository.NewReactionRepository(pool)
}

// ProvideReactionService creates a new ReactionService
func ProvideReactionService(reactionRepo repository.ReactionRepository, trackRepo repository.TrackRepository) service.ReactionService {
	return service.NewReactionService(reactionRepo, trackRepo)
}

// ProvideReactionHandler creates a new ReactionHandler
func ProvideReactionHandler(svc service.ReactionService) handler.ReactionHandler {
	return handler.NewReactionHandler(svc)
}

// ProvideIcecastClient creates a new IcecastClient
func ProvideIcecastClient(cfg config.Config) (icecast.Client, error) {
	return icecast.NewClient(cfg)
}

// ProvideRadioService creates a new RadioService
func ProvideRadioService(icecastClient icecast.Client) service.RadioService {
	return service.NewRadioService(icecastClient)
}

// ProvideRadioHandler creates a new RadioHandler
func ProvideRadioHandler(svc service.RadioService) handler.RadioHandler {
	return handler.NewRadioHandler(svc)
}

// ProvideListenerRepository creates a new ListenerRepository
func ProvideListenerRepository(pool *pgxpool.Pool) repository.ListenerRepository {
	return repository.NewListenerRepository(pool)
}

// ProvideStatisticsRepository creates a new StatisticsRepository
func ProvideStatisticsRepository(pool *pgxpool.Pool) repository.StatisticsRepository {
	return repository.NewStatisticsRepository(pool)
}

// ProvideStatisticsService creates a new StatisticsService
func ProvideStatisticsService(repo repository.StatisticsRepository) service.StatisticsService {
	return service.NewStatisticsService(repo)
}

// ProvideStatisticsHandler creates a new StatisticsHandler
func ProvideStatisticsHandler(svc service.StatisticsService) handler.StatisticsHandler {
	return handler.NewStatisticsHandler(svc)
}

// ProvideListenerService creates a new ListenerService
func ProvideListenerService(
	icecastClient icecast.Client,
	listenerRepo repository.ListenerRepository,
	trackRepo repository.TrackRepository,
) service.ListenerService {
	return service.NewListenerService(icecastClient, listenerRepo, trackRepo)
}

// ProvideServer creates a new Server instance
func ProvideServer(log *logger.Logger, pool *pgxpool.Pool, trackHandler handler.TrackHandler, reactionHandler handler.ReactionHandler, radioHandler handler.RadioHandler, statisticsHandler handler.StatisticsHandler) server.Server {
	return server.NewServer(log, pool, trackHandler, reactionHandler, radioHandler, statisticsHandler)
}

// ProvideScheduler creates a new Scheduler instance
func ProvideScheduler(listenerService service.ListenerService) scheduler.Scheduler {
	return scheduler.NewScheduler(listenerService)
}

// ProvideApplication creates the Application struct
func ProvideApplication(
	cfg config.Config,
	log *logger.Logger,
	db database.Database,
	srv server.Server,
	sched scheduler.Scheduler,
) *Application {
	return &Application{
		Config:    cfg,
		Logger:    log,
		Database:  db,
		Server:    srv,
		Scheduler: sched,
	}
}

// ProvideMigrateApp creates the MigrateApp struct
func ProvideMigrateApp(cfg config.Config, log *logger.Logger, dsn string) *MigrateApp {
	return &MigrateApp{
		Config: cfg,
		Logger: log,
		DSN:    dsn,
	}
}

// ProviderSet is the main set of providers
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideDSN,
	ProvideDatabase,
	ProvidePool,
	ProvideTrackRepository,
	ProvideTrackService,
	ProvideTrackHandler,
	ProvideReactionRepository,
	ProvideReactionService,
	ProvideReactionHandler,
	ProvideIcecastClient,
	ProvideRadioService,
	ProvideRadioHandler,
	ProvideListenerRepository,
	ProvideListenerService,
	ProvideStatisticsRepository,
	ProvideStatisticsService,
	ProvideStatisticsHandler,
	ProvideServer,
	ProvideScheduler,
	ProvideApplication,
)

// MigrateProviderSet is the set of providers for migrate commands
var MigrateProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideDSN,
	ProvideMigrateApp,
)

// InitializeApp creates a fully wired Application
func InitializeApp() (*Application, func(), error) {
	wire.Build(ProviderSet)
	return nil, nil, nil
}

// InitializeMigrateApp creates dependencies for migrate commands
func InitializeMigrateApp() (*MigrateApp, error) {
	wire.Build(MigrateProviderSet)
	return nil, nil
}
