package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	initDB()
	initRedis()
	defer closeDB()

	// Create router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// CORS middleware - allow frontend dev server
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.GET("/health", healthCheck)
	r.POST("/register", register)
	r.POST("/login", login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(authMiddleware())
	{
		protected.POST("/shorten", shortenURL)
		protected.GET("/shorten", listURLs)
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
