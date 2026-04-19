package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type shortenRequest struct {
	URL string `json:"url" binding:"required,url"`
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
