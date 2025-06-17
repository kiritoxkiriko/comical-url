package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"shorturl/config"
	"shorturl/handlers"
	"shorturl/middleware"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server to handle URL shortening requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServer() error {
	cfg := GetConfig()

	// Initialize database connections
	config.InitDatabaseWithConfig(cfg)

	// Set Gin mode
	if cfg.App.Name != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Initialize handlers
	urlHandler := handlers.NewURLHandler()
	authHandler := handlers.NewAuthHandler()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"app":       cfg.App.Name,
		})
	})

	// Auth routes (no auth required for token creation initially)
	auth := r.Group("/api/auth")
	{
		auth.POST("/tokens", authHandler.CreateToken)
		auth.DELETE("/tokens/:token", middleware.TokenAuth(), authHandler.RevokeToken)
		auth.GET("/tokens", middleware.TokenAuth(), authHandler.ListTokens)
	}

	// URL routes
	api := r.Group("/api")
	if cfg.App.RequireAuth {
		api.Use(middleware.TokenAuth()) // Mandatory auth for all API routes
	} else {
		api.Use(middleware.OptionalTokenAuth()) // Optional auth for all API routes
	}
	{
		api.POST("/shorten", urlHandler.CreateURL)
		api.GET("/info/:key", urlHandler.GetURLInfo)
		api.DELETE("/urls/:key", urlHandler.RevokeURL)
		api.POST("/auto-revoke", urlHandler.AutoRevoke)
	}

	// Direct redirect route (no /api prefix)
	r.GET("/:key", urlHandler.RedirectURL)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting %s on %s", cfg.App.Name, addr)

	if err := r.Run(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
