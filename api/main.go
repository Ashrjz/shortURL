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
	r.POST("/register", register)
	r.POST("/login", login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(authMiddleware())
	{
		protected.POST("/shorten", shortenURL)
		protected.PUT("/shorten/:code", updateShortURL)
		protected.DELETE("/shorten/:code", deleteShortURL)
	}

	// Public routes (no auth needed)
	r.GET("/shorten/:code", getShortURL)
	r.GET("/shorten/:code/stats", getURLStatsHandler)
	r.GET("/:code", redirectURL)

	// Start server
	r.Run(":8080")
}
