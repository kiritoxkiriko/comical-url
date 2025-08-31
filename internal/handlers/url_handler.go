package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"shorturl/internal/services"
)

type URLHandler struct {
	urlService *services.URLService
}

func NewURLHandler() *URLHandler {
	return &URLHandler{
		urlService: services.NewURLService(),
	}
}

type CreateURLRequest struct {
	LongURL   string `json:"long_url" binding:"required"`
	CustomKey string `json:"custom_key,omitempty"`
	Passkey   string `json:"passkey,omitempty"`
	ExpiresIn string `json:"expires_in,omitempty"` // e.g., "10s", "1h", "7d", "1y"
}

type CreateURLResponse struct {
	ShortKey  string     `json:"short_key"`
	ShortURL  string     `json:"short_url"`
	LongURL   string     `json:"long_url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (h *URLHandler) CreateURL(c *gin.Context) {
	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.urlService.CreateShortURL(req.LongURL, req.CustomKey, req.Passkey, req.ExpiresIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseURL := c.Request.Host
	if c.Request.TLS == nil {
		baseURL = "http://" + baseURL
	} else {
		baseURL = "https://" + baseURL
	}

	response := CreateURLResponse{
		ShortKey:  url.ShortKey,
		ShortURL:  baseURL + "/" + url.ShortKey,
		LongURL:   url.LongURL,
		ExpiresAt: url.ExpiresAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortKey := c.Param("key")
	passkey := c.Query("passkey")

	longURL, err := h.urlService.GetLongURL(shortKey, passkey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func (h *URLHandler) GetURLInfo(c *gin.Context) {
	shortKey := c.Param("key")
	passkey := c.Query("passkey")

	longURL, err := h.urlService.GetLongURL(shortKey, passkey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_key": shortKey,
		"long_url":  longURL,
	})
}

func (h *URLHandler) RevokeURL(c *gin.Context) {
	shortKey := c.Param("key")

	err := h.urlService.RevokeURL(shortKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL revoked successfully"})
}

func (h *URLHandler) AutoRevoke(c *gin.Context) {
	err := h.urlService.AutoRevokeExpiredURLs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto-revoke completed"})
}
