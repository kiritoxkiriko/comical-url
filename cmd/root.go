package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"shorturl/internal/config"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "shorturl",
	Short: "A URL shortening service",
	Long: `A fast and reliable URL shortening service built with Go, Gin, and GORM.

This application provides URL shortening capabilities with optional authentication,
custom keys, expiration dates, and passkey protection.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().StringP("host", "H", "", "server host")
	rootCmd.PersistentFlags().IntP("port", "p", 0, "server port")
	rootCmd.PersistentFlags().String("db-host", "", "database host")
	rootCmd.PersistentFlags().String("db-user", "", "database user")
	rootCmd.PersistentFlags().String("db-password", "", "database password")
	rootCmd.PersistentFlags().String("db-name", "", "database name")
	rootCmd.PersistentFlags().String("redis-host", "", "redis host")
	rootCmd.PersistentFlags().String("redis-password", "", "redis password")
}

func GetConfig() *config.Config {
	return cfg
}