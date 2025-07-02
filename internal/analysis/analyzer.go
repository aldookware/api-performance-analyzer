package analysis

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"time"
)

// CodeAnalysis represents the complete analysis results
type CodeAnalysis struct {
	SecurityIssues    []SecurityIssue    `json:"security_issues"`
	PerformanceHints  []PerformanceHint  `json:"performance_hints"`
	BestPractices     []BestPractice     `json:"best_practices"`
	AIRecommendations []AIRecommendation `json:"ai_recommendations"`
	AnalysisTime      time.Time          `json:"analysis_time"`
	CodeComplexity    int                `json:"code_complexity"`
	PerformanceScore  int                `json:"performance_score"`
	PerformanceGrade  string             `json:"performance_grade"`
	FilePath          string             `json:"file_path,omitempty"`
}

type SecurityIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	LineNumber  int    `json:"line_number"`
	Suggestion  string `json:"suggestion"`
}

type PerformanceHint struct {
	Issue       string `json:"issue"`
	Impact      string `json:"impact"`
	Solution    string `json:"solution"`
	CodeExample string `json:"code_example"`
	LineNumber  int    `json:"line_number"`
	Severity    string `json:"severity"`
}

type BestPractice struct {
	Category    string `json:"category"`
	Current     string `json:"current"`
	Recommended string `json:"recommended"`
	Reasoning   string `json:"reasoning"`
}

type AIRecommendation struct {
	Type           string  `json:"type"`
	Confidence     float64 `json:"confidence"`
	Recommendation string  `json:"recommendation"`
	AutoFixCode    string  `json:"auto_fix_code"`
}

// FileAnalysis represents analysis results for a single file
type FileAnalysis struct {
	FilePath string       `json:"file_path"`
	Analysis CodeAnalysis `json:"analysis"`
}

// AnalyzeCode performs comprehensive analysis on Go code
func AnalyzeCode(code, codeType, filePath string) CodeAnalysis {
	analysis := CodeAnalysis{
		AnalysisTime:      time.Now(),
		SecurityIssues:    []SecurityIssue{},
		PerformanceHints:  []PerformanceHint{},
		BestPractices:     []BestPractice{},
		AIRecommendations: []AIRecommendation{},
		FilePath:          filePath,
	}

	// Parse Go code
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, code, parser.ParseComments)
	if err != nil {
		// If parsing fails, return basic analysis
		analysis.SecurityIssues = append(analysis.SecurityIssues, SecurityIssue{
			Type:        "syntax_error",
			Description: "Code contains syntax errors",
			Severity:    "high",
			Suggestion:  "Fix syntax errors before analysis: " + err.Error(),
		})
		return analysis
	}

	// Analyze for common issues
	analysis.SecurityIssues = detectSecurityIssues(node, fset, code)
	analysis.PerformanceHints = detectPerformanceIssues(node, fset, code)
	analysis.BestPractices = suggestBestPractices(node, fset, code)
	analysis.AIRecommendations = generateAIRecommendations(code, codeType)
	analysis.CodeComplexity = calculateComplexity(node)

	// Calculate performance score
	analysis.PerformanceScore, analysis.PerformanceGrade = calculatePerformanceScore(analysis.PerformanceHints)

	return analysis
}

func detectSecurityIssues(node *ast.File, fset *token.FileSet, code string) []SecurityIssue {
	var issues []SecurityIssue

	// Check for missing CORS
	if !strings.Contains(code, "cors") && !strings.Contains(code, "CORS") {
		issues = append(issues, SecurityIssue{
			Type:        "missing_cors",
			Description: "No CORS middleware detected - this can cause browser security issues",
			Severity:    "medium",
			Suggestion:  "Add CORS middleware: router.Use(cors.Default())",
		})
	}

	// Check for missing error handling
	if strings.Contains(code, "BindJSON") && !strings.Contains(code, "StatusBadRequest") {
		issues = append(issues, SecurityIssue{
			Type:        "insufficient_error_handling",
			Description: "JSON binding without proper error responses",
			Severity:    "medium",
			Suggestion:  "Return proper HTTP error codes for invalid JSON",
		})
	}

	// Check for hardcoded secrets
	if strings.Contains(code, "password") || strings.Contains(code, "secret") || strings.Contains(code, "token") {
		issues = append(issues, SecurityIssue{
			Type:        "potential_hardcoded_secrets",
			Description: "Potential hardcoded secrets detected",
			Severity:    "high",
			Suggestion:  "Use environment variables for sensitive data",
		})
	}

	// Check for SQL injection risks
	if strings.Contains(code, "SELECT") || strings.Contains(code, "INSERT") || strings.Contains(code, "UPDATE") || strings.Contains(code, "DELETE") {
		if !strings.Contains(code, "Prepare") && !strings.Contains(code, "?") {
			issues = append(issues, SecurityIssue{
				Type:        "sql_injection_risk",
				Description: "Raw SQL queries detected - potential injection risk",
				Severity:    "high",
				Suggestion:  "Use parameterized queries or ORM with proper escaping",
			})
		}
	}

	return issues
}

func detectPerformanceIssues(node *ast.File, fset *token.FileSet, code string) []PerformanceHint {
	var hints []PerformanceHint

	// 🔥 NEW: Detect N+1 Query Patterns - THE MONEY MAKER
	hints = append(hints, detectN1QueryPatterns(node, fset, code)...)

	// 🔥 NEW: Detect Missing Database Indexes
	hints = append(hints, detectMissingIndexes(node, fset, code)...)

	// 🔥 NEW: Detect Large Payload Issues
	hints = append(hints, detectLargePayloadIssues(node, fset, code)...)

	// 🔥 NEW: Detect Missing Caching Opportunities
	hints = append(hints, detectMissingCaching(node, fset, code)...)

	// Check for in-memory data storage
	if strings.Contains(code, "var albums = []") || strings.Contains(code, "albums = append(") {
		hints = append(hints, PerformanceHint{
			Issue:    "In-memory data storage",
			Impact:   "🟡 MEDIUM: Data lost on restart, not scalable, memory leaks",
			Solution: "Implement database persistence (PostgreSQL, MongoDB, etc.)",
			Severity: "medium",
			CodeExample: `// Use database instead
import "database/sql"
db, err := sql.Open("postgres", connectionString)
// Store albums in database table`,
		})
	}

	// Check for missing connection pooling
	if strings.Contains(code, "sql.Open") && !strings.Contains(code, "SetMaxOpenConns") {
		hints = append(hints, PerformanceHint{
			Issue:    "No database connection pooling",
			Impact:   "🟠 HIGH: Resource exhaustion, poor scalability",
			Solution: "Configure connection pool settings",
			Severity: "high",
			CodeExample: `db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)`,
		})
	}

	return hints
}

// Detect N+1 Query Patterns - THE CRITICAL PERFORMANCE ISSUE
func detectN1QueryPatterns(node *ast.File, fset *token.FileSet, code string) []PerformanceHint {
	var hints []PerformanceHint

	// Pattern 1: Loop with database calls
	ast.Inspect(node, func(n ast.Node) bool {
		if forStmt, ok := n.(*ast.ForStmt); ok {
			// Check if there's a database call inside the loop
			hasDBCall := false
			ast.Inspect(forStmt, func(inner ast.Node) bool {
				if callExpr, ok := inner.(*ast.CallExpr); ok {
					if isDBCall(callExpr) {
						hasDBCall = true
						return false
					}
				}
				return true
			})

			if hasDBCall {
				pos := fset.Position(forStmt.Pos())
				hints = append(hints, PerformanceHint{
					Issue:      "Potential N+1 Query Pattern",
					Impact:     "🔴 CRITICAL: Could execute hundreds of database queries instead of one",
					Solution:   "Use JOIN queries or eager loading to fetch related data in one query",
					Severity:   "critical",
					LineNumber: pos.Line,
					CodeExample: `// ❌ Bad: N+1 queries
for _, user := range users {
    posts := db.Where("user_id = ?", user.ID).Find(&posts)
}

// ✅ Good: Single query with JOIN
db.Preload("Posts").Find(&users)`,
				})
			}
		}
		return true
	})

	// Pattern 2: Range loop with database calls
	ast.Inspect(node, func(n ast.Node) bool {
		if rangeStmt, ok := n.(*ast.RangeStmt); ok {
			hasDBCall := false
			ast.Inspect(rangeStmt, func(inner ast.Node) bool {
				if callExpr, ok := inner.(*ast.CallExpr); ok {
					if isDBCall(callExpr) {
						hasDBCall = true
						return false
					}
				}
				return true
			})

			if hasDBCall {
				pos := fset.Position(rangeStmt.Pos())
				hints = append(hints, PerformanceHint{
					Issue:      "N+1 Query in Range Loop",
					Impact:     "🔴 CRITICAL: Database call inside range loop creates N+1 queries",
					Solution:   "Extract database calls outside the loop or use batch queries",
					Severity:   "critical",
					LineNumber: pos.Line,
					CodeExample: `// ❌ Bad: Query in range loop
for _, order := range orders {
    db.Model(&order).Related(&order.Items)
}

// ✅ Good: Preload related data
db.Preload("Items").Find(&orders)`,
				})
			}
		}
		return true
	})

	// Pattern 3: GORM Related() calls in loops
	if strings.Contains(code, "for") && strings.Contains(code, ".Related(") {
		hints = append(hints, PerformanceHint{
			Issue:    "GORM N+1 Query with Related()",
			Impact:   "🔴 CRITICAL: Each Related() call executes a separate query",
			Solution: "Use Preload() to fetch related data efficiently",
			Severity: "critical",
			CodeExample: `// ❌ Bad: N+1 with Related()
for _, order := range orders {
    db.Model(&order).Related(&order.Items)
}

// ✅ Good: Preload related data
db.Preload("Items").Find(&orders)`,
		})
	}

	return hints
}

// Detect queries that would benefit from database indexes
func detectMissingIndexes(node *ast.File, fset *token.FileSet, code string) []PerformanceHint {
	var hints []PerformanceHint

	// Look for WHERE clauses on potentially unindexed columns
	wherePatterns := []string{
		"WHERE name",
		"WHERE email",
		"WHERE status",
		"WHERE created_at",
		"LIKE '%",
		".Where(\"name",
		".Where(\"email",
		".Where(\"status",
	}

	for _, pattern := range wherePatterns {
		if strings.Contains(code, pattern) {
			hints = append(hints, PerformanceHint{
				Issue:    "Potential Missing Database Index",
				Impact:   "🟠 HIGH: Query will scan entire table without proper index",
				Solution: "Add database index on frequently queried columns",
				Severity: "high",
				CodeExample: `-- Add these indexes to your migration:
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- For GORM, add to your struct:
type User struct {
    Email string \` + "`gorm:\"index\"`" + `
}`,
			})
			break // Only add one hint for missing indexes
		}
	}

	return hints
}

// Detect large payload serialization issues
func detectLargePayloadIssues(node *ast.File, fset *token.FileSet, code string) []PerformanceHint {
	var hints []PerformanceHint

	// Look for JSON responses without pagination
	if (strings.Contains(code, "c.JSON") || strings.Contains(code, "json.Marshal")) &&
		(strings.Contains(code, ".Find(&") || strings.Contains(code, "SELECT *")) &&
		!strings.Contains(code, "Limit") && !strings.Contains(code, "LIMIT") {
		hints = append(hints, PerformanceHint{
			Issue:    "Large Dataset Response Without Pagination",
			Impact:   "🟠 HIGH: Returning large datasets will cause slow responses and high memory usage",
			Solution: "Implement pagination for large data responses",
			Severity: "high",
			CodeExample: `// ❌ Bad: Return all records
var users []User
db.Find(&users)
c.JSON(200, users)

// ✅ Good: Paginated response
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
offset := (page - 1) * limit

var users []User
var total int64
db.Model(&User{}).Count(&total)
db.Offset(offset).Limit(limit).Find(&users)

c.JSON(200, gin.H{
    "data": users, 
    "page": page,
    "total": total,
    "pages": (total + int64(limit) - 1) / int64(limit),
})`,
		})
	}

	return hints
}

// Detect expensive operations that should be cached
func detectMissingCaching(node *ast.File, fset *token.FileSet, code string) []PerformanceHint {
	var hints []PerformanceHint

	expensiveOperations := []string{
		"calculateReport",
		"generateStats",
		"getAggregated",
		"COUNT(",
		"SUM(",
		"AVG(",
		"GROUP BY",
		"ORDER BY",
	}

	for _, op := range expensiveOperations {
		if strings.Contains(code, op) && !strings.Contains(code, "cache") && !strings.Contains(code, "Cache") {
			hints = append(hints, PerformanceHint{
				Issue:    "Expensive Operation Without Caching",
				Impact:   "🟡 MEDIUM: Repeated expensive calculations will slow down your API",
				Solution: "Cache expensive operation results with Redis or in-memory cache",
				Severity: "medium",
				CodeExample: `// Add caching for expensive operations
func getExpensiveReport(c *gin.Context) {
    cacheKey := "report_" + c.Query("type")
    
    // Check cache first
    if cached, exists := cache.Get(cacheKey); exists {
        c.JSON(200, cached)
        return
    }
    
    // Calculate only if not cached
    result := calculateExpensiveReport()
    cache.Set(cacheKey, result, 5*time.Minute)
    c.JSON(200, result)
}

// Or use middleware caching
router.GET("/reports", cache.CachePage(store, time.Minute*5, getReports))`,
			})
			break
		}
	}

	return hints
}

// Helper function to detect database calls
func isDBCall(callExpr *ast.CallExpr) bool {
	if selector, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		dbMethods := []string{
			"Find", "First", "Last", "Take", "Where", "Select", "Order", "Limit", "Offset",
			"Create", "Save", "Update", "UpdateColumn", "UpdateColumns", "Updates",
			"Delete", "Unscoped", "Raw", "Exec", "Scan", "Rows", "Row",
			"Count", "Group", "Having", "Joins", "Preload", "Related", "Association",
		}
		for _, method := range dbMethods {
			if selector.Sel.Name == method {
				return true
			}
		}
	}
	return false
}

func suggestBestPractices(node *ast.File, fset *token.FileSet, code string) []BestPractice {
	var practices []BestPractice

	practices = append(practices, BestPractice{
		Category:    "Error Handling",
		Current:     "Basic error handling",
		Recommended: "Structured error responses with proper HTTP status codes",
		Reasoning:   "Better client experience, easier debugging, API consistency",
	})

	if !strings.Contains(code, "middleware") {
		practices = append(practices, BestPractice{
			Category:    "Request Validation",
			Current:     "No input validation middleware",
			Recommended: "Add request validation and sanitization",
			Reasoning:   "Prevent invalid data, improve security, reduce debugging time",
		})
	}

	if !strings.Contains(code, "Logger") {
		practices = append(practices, BestPractice{
			Category:    "Logging",
			Current:     "No structured logging",
			Recommended: "Add request/response logging with correlation IDs",
			Reasoning:   "Essential for debugging production issues and monitoring",
		})
	}

	if strings.Contains(code, "Run(\"localhost:8080\")") {
		practices = append(practices, BestPractice{
			Category:    "Configuration",
			Current:     "Hardcoded port and host",
			Recommended: "Use environment variables for configuration",
			Reasoning:   "Enables different environments (dev, staging, prod)",
		})
	}

	return practices
}

func generateAIRecommendations(code, codeType string) []AIRecommendation {
	var recommendations []AIRecommendation

	// Architecture recommendations
	if strings.Contains(code, "func get") && strings.Contains(code, "func post") {
		recommendations = append(recommendations, AIRecommendation{
			Type:           "architecture",
			Confidence:     0.85,
			Recommendation: "Consider implementing repository pattern for cleaner data access",
			AutoFixCode: `type AlbumRepository interface {
    GetAll(ctx context.Context) ([]Album, error)
    GetByID(ctx context.Context, id string) (*Album, error)
    Create(ctx context.Context, album Album) error
    Update(ctx context.Context, album Album) error
    Delete(ctx context.Context, id string) error
}`,
		})
	}

	// Middleware recommendations
	if !strings.Contains(code, "Recovery()") {
		recommendations = append(recommendations, AIRecommendation{
			Type:           "middleware",
			Confidence:     0.92,
			Recommendation: "Add essential middleware for production readiness",
			AutoFixCode: `router.Use(gin.Logger())
router.Use(gin.Recovery())
router.Use(cors.Default())
router.Use(rateLimitMiddleware())`,
		})
	}

	// API versioning
	if !strings.Contains(code, "/v1/") && !strings.Contains(code, "/api/v") {
		recommendations = append(recommendations, AIRecommendation{
			Type:           "api_design",
			Confidence:     0.78,
			Recommendation: "Add API versioning for future compatibility",
			AutoFixCode: `v1 := router.Group("/api/v1")
{
    v1.GET("/albums", getAlbums)
    v1.POST("/albums", postAlbums)
    v1.GET("/albums/:id", getAlbumByID)
}`,
		})
	}

	// Testing recommendations
	if !strings.Contains(code, "_test.go") {
		recommendations = append(recommendations, AIRecommendation{
			Type:           "testing",
			Confidence:     0.88,
			Recommendation: "Add unit tests for better code reliability",
			AutoFixCode: `func TestGetAlbums(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/albums", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}`,
		})
	}

	return recommendations
}

func calculateComplexity(node *ast.File) int {
	complexity := 0

	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
			complexity++
		}
		return true
	})

	return complexity
}

// Calculate performance score based on detected issues
func calculatePerformanceScore(hints []PerformanceHint) (int, string) {
	score := 100

	for _, hint := range hints {
		switch hint.Severity {
		case "critical":
			score -= 25
		case "high":
			score -= 15
		case "medium":
			score -= 10
		case "low":
			score -= 5
		}
	}

	if score < 0 {
		score = 0
	}

	var grade string
	switch {
	case score >= 95:
		grade = "A+"
	case score >= 90:
		grade = "A"
	case score >= 80:
		grade = "B"
	case score >= 70:
		grade = "C"
	case score >= 60:
		grade = "D"
	default:
		grade = "F"
	}

	return score, grade
}
