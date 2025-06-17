package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shorturl/config"
	"shorturl/models"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type CreateTokenRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateTokenResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

func (h *AuthHandler) CreateToken(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := uuid.New().String()
	authToken := &models.AuthToken{
		Token:    token,
		Name:     req.Name,
		IsActive: true,
	}

	if err := config.DB.Create(authToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	response := CreateTokenResponse{
		Token: token,
		Name:  req.Name,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) RevokeToken(c *gin.Context) {
	token := c.Param("token")

	result := config.DB.Model(&models.AuthToken{}).Where("token = ?", token).Update("is_active", false)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}

func (h *AuthHandler) ListTokens(c *gin.Context) {
	var tokens []models.AuthToken
	if err := config.DB.Where("is_active = ?", true).Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}
