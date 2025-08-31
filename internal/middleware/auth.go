package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"shorturl/internal/config"
	"shorturl/internal/models"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		var authToken models.AuthToken
		if err := config.DB.Where("token = ? AND is_active = ?", token, true).First(&authToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("auth_token", authToken)
		c.Next()
	}
}

func OptionalTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				token := tokenParts[1]
				var authToken models.AuthToken
				if err := config.DB.Where("token = ? AND is_active = ?", token, true).First(&authToken).Error; err == nil {
					c.Set("auth_token", authToken)
				}
			}
		}
		c.Next()
	}
}
