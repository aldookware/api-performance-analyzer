package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/your-username/api-performance-analyzer/internal/analysis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AnalyzeRequest struct {
	Code string `json:"code" binding:"required"`
	Type string `json:"type"`
}

// analyzeCode analyzes Go REST API code for patterns and issues
func analyzeCode(c *gin.Context) {
	var request AnalyzeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if strings.TrimSpace(request.Code) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code cannot be empty"})
		return
	}

	// Use the new analysis package
	result := analysis.AnalyzeCode(request.Code, request.Type, "input.go")
	c.JSON(http.StatusOK, result)
}

// getHealth returns API health status
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "API Performance Analyzer",
		"version":   "1.0.0",
	})
}

// getStats returns usage statistics
func getStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_analyses":      0, // In production, track this in database
		"supported_languages": []string{"go", "gin", "echo"},
		"features": []string{
			"N+1 Query Detection",
			"Missing Index Analysis",
			"Large Payload Detection",
			"Caching Opportunities",
			"Security Analysis",
			"Performance Scoring",
		},
	})
}

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Serve static files (for the web interface)
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", getHealth)
		api.GET("/stats", getStats)
		api.POST("/analyze", analyzeCode)
	}

	fmt.Println("🚀 API Performance Analyzer starting...")
	fmt.Println("📊 Dashboard: http://localhost:8080")
	fmt.Println("🔗 API: http://localhost:8080/api/v1")
	fmt.Println("❤️  Health: http://localhost:8080/api/v1/health")

	router.Run(":8080")
}
