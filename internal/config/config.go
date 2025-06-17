package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`     // mysql, postgres, sqlite
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	DSN      string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AppConfig struct {
	Name            string `mapstructure:"name"`
	DefaultExpire   string `mapstructure:"default_expire"`     // e.g., "30d", "1y"
	KeyLength       int    `mapstructure:"key_length"`
	CacheDuration   string `mapstructure:"cache_duration"`     // e.g., "7d", "168h"
	RequireAuth     bool   `mapstructure:"require_auth"`       // mandatory token auth
	AllowCustomKeys bool   `mapstructure:"allow_custom_keys"`  // allow custom short keys
	MaxURLLength    int    `mapstructure:"max_url_length"`     // max URL length
}

var GlobalConfig *Config

func LoadConfig(configFile string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/shorturl")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// Set defaults
	setDefaults()

	// Enable environment variable support
	viper.SetEnvPrefix("SHORTURL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Build DSN if not provided
	if config.Database.DSN == "" {
		switch config.Database.Type {
		case "postgres", "postgresql":
			config.Database.DSN = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				config.Database.Host,
				config.Database.Port,
				config.Database.User,
				config.Database.Password,
				config.Database.Database,
			)
		case "sqlite":
			config.Database.DSN = config.Database.Database
		default: // mysql
			config.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.Database.User,
				config.Database.Password,
				config.Database.Host,
				config.Database.Port,
				config.Database.Database,
			)
		}
	}

	GlobalConfig = &config
	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	// Database defaults
	viper.SetDefault("database.type", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.database", "shorturl")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// App defaults
	viper.SetDefault("app.name", "Short URL Service")
	viper.SetDefault("app.default_expire", "30d")
	viper.SetDefault("app.key_length", 6)
	viper.SetDefault("app.cache_duration", "7d")
	viper.SetDefault("app.require_auth", false)
	viper.SetDefault("app.allow_custom_keys", true)
	viper.SetDefault("app.max_url_length", 2048)
}
