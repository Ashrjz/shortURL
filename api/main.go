package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	initDB()
	defer closeDB()

	// Create router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Routes
	r.GET("/health", healthCheck)
	r.POST("/shorten", shortenURL)
	r.GET("/shorten/:code", getShortURL)
	r.PUT("/shorten/:code", updateShortURL)
	r.DELETE("/shorten/:code", deleteShortURL)

	// Start server
	r.Run(":8080")
}
