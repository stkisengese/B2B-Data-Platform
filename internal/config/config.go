package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the entire application configuration
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Collector   CollectorConfig   `mapstructure:"collector"`
	DataSources DataSourcesConfig `mapstructure:"datasources"`
	Logging     LoggingConfig     `mapstructure:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// CollectorConfig holds configuration for the collector service
type CollectorConfig struct {
	WorkerCount     int           `mapstructure:"worker_count"`
	QueueSize       int           `mapstructure:"queue_size"`
	RetryAttempts   int           `mapstructure:"retry_attempts"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	JobRetention    time.Duration `mapstructure:"job_retention"`
}

// DataSourcesConfig holds configuration for external data sources
type DataSourcesConfig struct {
	OpenCorporates OpenCorporatesConfig `mapstructure:"opencorporates"`
	CompaniesHouse CompaniesHouseConfig `mapstructure:"companieshouse"`
}

// OpenCorporatesConfig holds configuration for the OpenCorporates data source
type OpenCorporatesConfig struct {
	APIKey    string `mapstructure:"api_key"`
	Enabled   bool   `mapstructure:"enabled"`
	RateLimit int    `mapstructure:"rate_limit"`
}

// CompaniesHouseConfig holds configuration for the Companies House data source
type CompaniesHouseConfig struct {
	APIKey    string `mapstructure:"api_key"`
	Enabled   bool   `mapstructure:"enabled"`
	RateLimit int    `mapstructure:"rate_limit"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.path", "b2b.db")

	// Collector defaults
	viper.SetDefault("collector.worker_count", 5)
	viper.SetDefault("collector.queue_size", 100)
	viper.SetDefault("collector.retry_attempts", 3)
	viper.SetDefault("collector.retry_delay", "5s")
	viper.SetDefault("collector.shutdown_timeout", "30s")
	viper.SetDefault("collector.cleanup_interval", "1h")
	viper.SetDefault("collector.job_retention", "24h")

	// DataSource defaults
	viper.SetDefault("datasources.opencorporates.enabled", false)
	viper.SetDefault("datasources.opencorporates.rate_limit", 5)
	viper.SetDefault("datasources.companieshouse.enabled", false)
	viper.SetDefault("datasources.companieshouse.rate_limit", 10)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
