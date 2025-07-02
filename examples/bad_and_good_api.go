package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Order represents an order in the system
type Order struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	Amount int  `json:"amount"`
	User   User `json:"user"`
}

var db *gorm.DB

// BAD: This function has multiple performance and security issues
func getUserOrders(c *gin.Context) {
	userID := c.Param("user_id")

	// Security issue: SQL injection vulnerability
	query := fmt.Sprintf("SELECT * FROM orders WHERE user_id = %s", userID)
	rows, err := db.Raw(query).Rows()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var orders []Order

	// Performance issue: N+1 query problem
	for rows.Next() {
		var order Order
		if err := db.ScanRows(rows, &order); err != nil {
			continue
		}

		// Each iteration triggers a separate query - N+1 problem!
		db.Model(&order).Related(&order.User)
		orders = append(orders, order)
	}

	// Performance issue: Large payload - returning everything
	c.JSON(200, orders)
}

// GOOD: Optimized version
func getUserOrdersOptimized(c *gin.Context) {
	userIDStr := c.Param("user_id")

	// Security: Proper input validation and parameterized query
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var orders []Order

	// Performance: Use Preload to solve N+1 problem
	result := db.Preload("User").Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Performance: Pagination to avoid large payloads
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var paginatedOrders []Order
	db.Preload("User").Where("user_id = ?", userID).
		Limit(limit).Offset(offset).Find(&paginatedOrders)

	c.JSON(200, gin.H{
		"orders": paginatedOrders,
		"page":   page,
		"limit":  limit,
	})
}

// BAD: Missing middleware and security
func setupBadRouter() *gin.Engine {
	r := gin.Default()

	// Missing CORS middleware
	// Missing rate limiting
	// Missing authentication middleware

	r.GET("/users/:user_id/orders", getUserOrders)

	return r
}

// GOOD: Proper middleware setup
func setupGoodRouter() *gin.Engine {
	r := gin.Default()

	// Add essential middleware
	r.Use(corsMiddleware())
	r.Use(rateLimitMiddleware())
	r.Use(authMiddleware())
	r.Use(loggingMiddleware())

	// API versioning
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users/:user_id/orders", getUserOrdersOptimized)
	}

	return r
}

// Placeholder middleware functions (would be implemented)
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rate limiting logic would go here
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authentication logic would go here
		c.Next()
	}
}

func loggingMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

func main() {
	// BAD: Hardcoded secrets
	password := "mysecretpassword123"
	fmt.Println("Database password:", password)

	// Set up router
	r := setupGoodRouter()

	// BAD: Hardcoded port
	r.Run(":8080")
}
