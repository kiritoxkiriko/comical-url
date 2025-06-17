package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	config2 "shorturl/internal/config"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitDatabase() {
	// Default initialization for backward compatibility
	cfg := &config2.Config{
		Database: config2.DatabaseConfig{
			Type:     "mysql",
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			Database: "shorturl",
		},
		Redis: config2.RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
	}
	InitDatabaseWithConfig(cfg)
}

func InitDatabaseWithConfig(cfg *config2.Config) {
	// Database connection
	var err error
	var dialector gorm.Dialector
	
	switch cfg.Database.Type {
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.Database.DSN)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.DSN)
	default: // mysql
		dialector = mysql.Open(cfg.Database.DSN)
	}
	
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", cfg.Database.Type, err)
	}

	// Redis connection
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err = Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	fmt.Println("Database connections established successfully")
}