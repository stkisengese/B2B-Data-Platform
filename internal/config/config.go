package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	DataSources DataSourcesConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Path string
}

type DataSourcesConfig struct {
	// OpenCorporates OpenCorporatesConfig
	CompaniesHouse CompaniesHouseConfig
}

// type OpenCorporatesConfig struct {
// 	APIKey   string
// 	Enabled  bool
// 	RateLimit int
// }

type CompaniesHouseConfig struct {
	APIKey    string
	Enabled   bool
	RateLimit int
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.path", "b2b.db")
	// viper.SetDefault("datasources.opencorporates.enabled", false)
	// viper.SetDefault("datasources.opencorporates.ratelimit", 5)
	viper.SetDefault("datasources.companieshouse.enabled", false)
	viper.SetDefault("datasources.companieshouse.ratelimit", 5)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
