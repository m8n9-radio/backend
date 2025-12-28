package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config interface {
		Port() int
		LogLevel() string
		DatabaseConnection() (string, int, int)
		RedisConnection() (string, string)
		IcecastConnection() (string, string, string, string)
		SchedulerEnabled() bool
	}
	config struct {
		port     int
		logLevel string

		dbHost string
		dbPort int
		dbName string
		dbUser string
		dbPass string

		dbMaxConns int
		dbMinConns int

		icecastHost     string
		icecastPort     int
		icecastUser     string
		icecastPassword string
		icecastMount    string

		redis_host     string
		redis_port     int
		redis_user     string
		redis_password string
		redis_db       int
		redis_prefix   string

		scheduler bool
	}
)

func NewConfig() Config {
	_ = godotenv.Load()
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "dev")
	viper.SetDefault("DB_USERNAME", "dev")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_MIN_CONNS", "1")
	viper.SetDefault("DB_MAX_CONNS", "4")

	viper.SetDefault("ICECAST_HOST", "127.0.0.1")
	viper.SetDefault("ICECAST_PORT", "8000")
	viper.SetDefault("ICECAST_USER", "admin")
	viper.SetDefault("ICECAST_PASSWORD", "changeme")
	viper.SetDefault("ICECAST_MOUNT", "/mp3")

	viper.SetDefault("REDIS_HOST", "127.0.0.1")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_USERNAME", "admin")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", "1")
	viper.SetDefault("REDIS_PREFIX", "be___")

	viper.SetDefault("SCHEDULER_ENABLED", "false")

	return &config{
		port:     viper.GetInt("PORT"),
		logLevel: viper.GetString("LOG_LEVEL"),

		dbHost: viper.GetString("DB_HOST"),
		dbPort: viper.GetInt("DB_PORT"),
		dbName: viper.GetString("DB_NAME"),
		dbUser: viper.GetString("DB_USERNAME"),
		dbPass: viper.GetString("DB_PASSWORD"),

		dbMaxConns: viper.GetInt("DB_MAX_CONNS"),
		dbMinConns: viper.GetInt("DB_MIN_CONNS"),

		icecastHost:     viper.GetString("ICECAST_HOST"),
		icecastPort:     viper.GetInt("ICECAST_PORT"),
		icecastUser:     viper.GetString("ICECAST_USER"),
		icecastPassword: viper.GetString("ICECAST_PASSWORD"),
		icecastMount:    viper.GetString("ICECAST_MOUNT"),

		redis_host:     viper.GetString("REDIS_HOST"),
		redis_port:     viper.GetInt("REDIS_PORT"),
		redis_user:     viper.GetString("REDIS_USERNAME"),
		redis_password: viper.GetString("REDIS_PASSWORD"),
		redis_db:       viper.GetInt("REDIS_DB"),
		redis_prefix:   viper.GetString("REDIS_PREFIX"),

		scheduler: viper.GetBool("SCHEDULER_ENABLED"),
	}
}

func (c *config) Port() int {
	return c.port
}

func (c *config) LogLevel() string {
	return c.logLevel
}

func (c *config) DatabaseConnection() (string, int, int) {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.dbUser, c.dbPass, c.dbHost, c.dbPort, c.dbName,
	), c.dbMinConns, c.dbMaxConns
}

func (c *config) IcecastConnection() (string, string, string, string) {
	return fmt.Sprintf("http://%s:%d", c.icecastHost, c.icecastPort),
		c.icecastUser,
		c.icecastPassword,
		c.icecastMount
}

func (c *config) RedisConnection() (string, string) {
	return fmt.Sprintf(
		"redis://%s:%s@%s:%d/%d", c.redis_user, c.redis_password, c.redis_host, c.redis_port, c.redis_db,
	), c.redis_prefix
}

func (c *config) SchedulerEnabled() bool {
	return c.scheduler
}
