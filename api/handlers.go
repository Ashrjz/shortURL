package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type shortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type UpdateURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	userID, err := createUser(req.Username, hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	token, err := generateToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, passwordHash, err := getUserByUsername(req.Username)
	if err != nil || !checkPassword(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func shortenURL(c *gin.Context) {
	var req shortenRequest

	// Validate and bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. 'url' field is required and must be a valid URL",
		})
		return
	}

	// Additional validation: ensure URL has scheme
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "URL must start with http:// or https://",
		})
		return
	}

	url, err := createURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create short URL",
		})
		return
	}

	c.JSON(http.StatusCreated, url)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Server is running",
	})
}

func getShortURL(c *gin.Context) {
	shortCode := c.Param("code")

	url, err := getURLByShortCode(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve URL",
		})
		return
	}

	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	c.JSON(http.StatusOK, url)
}

func updateShortURL(c *gin.Context) {
	shortCode := c.Param("code")
	var req UpdateURLRequest

	// Validate and bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. 'url' field is required and must be a valid URL",
		})
		return
	}

	// Additional validation: ensure URL has scheme
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "URL must start with http:// or https://",
		})
		return
	}

	// Update the URL
	url, err := updateURL(shortCode, req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update URL",
		})
		return
	}

	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	deleteCachedURL(shortCode)
	c.JSON(http.StatusOK, url)
}

func deleteShortURL(c *gin.Context) {
	shortCode := c.Param("code")

	deleted, err := deleteURL(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete URL",
		})
		return
	}

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	deleteCachedURL(shortCode)
	c.Status(http.StatusNoContent)
}

func redirectURL(c *gin.Context) {
	shortCode := c.Param("code")

	if cachedURL, err := getCachedURL(shortCode); err == nil {
		recordAccess(shortCode)
		c.Redirect(http.StatusFound, cachedURL)
		return
	}

	url, err := getURLByShortCode(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve URL",
		})
		return
	}

	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	setCachedURL(shortCode, url.URL)
	recordAccess(shortCode)
	c.Redirect(http.StatusFound, url.URL)
}

func getURLStatsHandler(c *gin.Context) {
	shortCode := c.Param("code")

	stats, err := getURLStats(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve stats",
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
