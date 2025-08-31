package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"shorturl/internal/config"
	"shorturl/internal/models"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to create or update the database schema.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrations() error {
	cfg := GetConfig()

	// Initialize database connections
	config.InitDatabaseWithConfig(cfg)

	log.Println("Running database migrations...")

	// Auto migrate the schema
	err := config.DB.AutoMigrate(&models.URL{}, &models.AuthToken{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
