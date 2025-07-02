package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aldookware/api-performance-analyzer/internal/analysis"
)

type Config struct {
	CodePath          string
	OutputFormat      string
	SeverityThreshold string
	FailOnIssues      bool
	Verbose           bool
}

func main() {
	config := parseFlags()

	if config.Verbose {
		fmt.Printf("🚀 API Performance Analyzer\n")
		fmt.Printf("Analyzing code in: %s\n", config.CodePath)
		fmt.Printf("Output format: %s\n", config.OutputFormat)
		fmt.Printf("Severity threshold: %s\n", config.SeverityThreshold)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	results, err := analyzeCodebase(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Analysis failed: %v\n", err)
		os.Exit(1)
	}

	// Output results
	switch config.OutputFormat {
	case "json":
		outputJSON(results)
	case "sarif":
		outputSARIF(results)
	case "github":
		outputGitHub(results)
	default:
		outputMarkdown(results)
	}

	// Set GitHub Action outputs if running in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		setGitHubActionOutputs(results)
	}

	// Check if we should fail the build
	if config.FailOnIssues && hasHighSeverityIssues(results, config.SeverityThreshold) {
		fmt.Fprintf(os.Stderr, "❌ Performance issues found above threshold (%s)\n", config.SeverityThreshold)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Printf("✅ Analysis complete! Found %d total issues across %d files\n",
			countTotalIssues(results), len(results))
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.CodePath, "path", ".", "Path to analyze")
	flag.StringVar(&config.OutputFormat, "format", "markdown", "Output format (markdown, json, sarif, github)")
	flag.StringVar(&config.SeverityThreshold, "threshold", "medium", "Severity threshold (low, medium, high, critical)")
	flag.BoolVar(&config.FailOnIssues, "fail-on-issues", false, "Fail on issues above threshold")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")
	flag.Parse()

	// Override with environment variables if present (for GitHub Actions)
	if path := os.Getenv("INPUT_CODE_PATH"); path != "" {
		config.CodePath = path
	}
	if format := os.Getenv("INPUT_OUTPUT_FORMAT"); format != "" {
		config.OutputFormat = format
	}
	if threshold := os.Getenv("INPUT_SEVERITY_THRESHOLD"); threshold != "" {
		config.SeverityThreshold = threshold
	}
	if failStr := os.Getenv("INPUT_FAIL_ON_ISSUES"); failStr == "true" {
		config.FailOnIssues = true
	}

	return config
}

func analyzeCodebase(config Config) ([]analysis.FileAnalysis, error) {
	var results []analysis.FileAnalysis

	err := filepath.WalkDir(config.CodePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only analyze Go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip vendor, .git, and test files for faster analysis
		if strings.Contains(path, "vendor/") ||
			strings.Contains(path, ".git/") ||
			strings.Contains(path, "_test.go") {
			return nil
		}

		if config.Verbose {
			fmt.Printf("📁 Analyzing: %s\n", path)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		analysisResult := analysis.AnalyzeCode(string(content), "go", path)

		results = append(results, analysis.FileAnalysis{
			FilePath: path,
			Analysis: analysisResult,
		})

		return nil
	})

	return results, err
}

func outputMarkdown(results []analysis.FileAnalysis) {
	fmt.Printf("# 🚀 API Performance Analysis Report\n\n")

	totalIssues := 0
	criticalIssues := 0
	highIssues := 0
	securityIssues := 0
	performanceIssues := 0

	// Summary stats
	for _, result := range results {
		totalIssues += len(result.Analysis.PerformanceHints) + len(result.Analysis.SecurityIssues)
		performanceIssues += len(result.Analysis.PerformanceHints)
		securityIssues += len(result.Analysis.SecurityIssues)

		for _, hint := range result.Analysis.PerformanceHints {
			if hint.Severity == "critical" {
				criticalIssues++
			} else if hint.Severity == "high" {
				highIssues++
			}
		}
	}

	// Overall Performance Score
	if len(results) > 0 {
		avgScore := 0
		for _, result := range results {
			avgScore += result.Analysis.PerformanceScore
		}
		avgScore = avgScore / len(results)

		scoreEmoji := "🔴"
		if avgScore >= 80 {
			scoreEmoji = "🟢"
		} else if avgScore >= 60 {
			scoreEmoji = "🟡"
		}

		fmt.Printf("## %s Overall Performance Score: %d/100\n\n", scoreEmoji, avgScore)
	}

	// Summary table
	fmt.Printf("| Metric | Count |\n")
	fmt.Printf("|--------|-------|\n")
	fmt.Printf("| 📁 Files Analyzed | %d |\n", len(results))
	fmt.Printf("| 🚨 Critical Issues | %d |\n", criticalIssues)
	fmt.Printf("| ⚠️ High Priority Issues | %d |\n", highIssues)
	fmt.Printf("| ⚡ Performance Issues | %d |\n", performanceIssues)
	fmt.Printf("| 🔒 Security Issues | %d |\n", securityIssues)
	fmt.Printf("| 📊 Total Issues | %d |\n\n", totalIssues)

	// File-by-file analysis
	for _, result := range results {
		if len(result.Analysis.PerformanceHints) > 0 || len(result.Analysis.SecurityIssues) > 0 {
			fmt.Printf("## 📄 %s\n\n", result.FilePath)
			fmt.Printf("**Performance Score:** %d/100 (%s)\n\n",
				result.Analysis.PerformanceScore, result.Analysis.PerformanceGrade)

			// Performance Issues
			if len(result.Analysis.PerformanceHints) > 0 {
				fmt.Printf("### ⚡ Performance Issues\n\n")
				for _, hint := range result.Analysis.PerformanceHints {
					severityEmoji := getSeverityEmoji(hint.Severity)
					fmt.Printf("#### %s %s\n", severityEmoji, hint.Issue)
					fmt.Printf("**Impact:** %s\n\n", hint.Impact)
					fmt.Printf("**Solution:** %s\n\n", hint.Solution)
					if hint.LineNumber > 0 {
						fmt.Printf("**Line:** %d\n\n", hint.LineNumber)
					}
					if hint.CodeExample != "" {
						fmt.Printf("```go\n%s\n```\n\n", hint.CodeExample)
					}
				}
			}

			// Security Issues
			if len(result.Analysis.SecurityIssues) > 0 {
				fmt.Printf("### 🔒 Security Issues\n\n")
				for _, issue := range result.Analysis.SecurityIssues {
					severityEmoji := getSeverityEmoji(issue.Severity)
					fmt.Printf("#### %s %s\n", severityEmoji, strings.ReplaceAll(issue.Type, "_", " "))
					fmt.Printf("**Description:** %s\n\n", issue.Description)
					fmt.Printf("**Suggestion:** %s\n\n", issue.Suggestion)
					if issue.LineNumber > 0 {
						fmt.Printf("**Line:** %d\n\n", issue.LineNumber)
					}
				}
			}

			fmt.Printf("---\n\n")
		}
	}

	if totalIssues == 0 {
		fmt.Printf("## 🎉 Excellent!\n\nNo performance or security issues detected. Your API is well-optimized!\n\n")
	}

	fmt.Printf("---\n")
	fmt.Printf("*Generated by [API Performance Analyzer](https://github.com/marketplace/actions/api-performance-analyzer)*\n")
}

func getSeverityEmoji(severity string) string {
	switch severity {
	case "critical":
		return "🔴"
	case "high":
		return "🟠"
	case "medium":
		return "🟡"
	case "low":
		return "🔵"
	default:
		return "⚪"
	}
}

func outputJSON(results []analysis.FileAnalysis) {
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(jsonData))
}

func outputSARIF(results []analysis.FileAnalysis) {
	// SARIF format for GitHub Security tab integration
	sarif := map[string]interface{}{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":           "API Performance Analyzer",
						"version":        "1.0.0",
						"informationUri": "https://github.com/marketplace/actions/api-performance-analyzer",
					},
				},
				"results": convertToSARIFResults(results),
			},
		},
	}

	jsonData, _ := json.MarshalIndent(sarif, "", "  ")
	fmt.Println(string(jsonData))
}

func outputGitHub(results []analysis.FileAnalysis) {
	// GitHub-specific output format for annotations
	for _, result := range results {
		for _, hint := range result.Analysis.PerformanceHints {
			level := "warning"
			if hint.Severity == "critical" || hint.Severity == "high" {
				level = "error"
			}

			line := hint.LineNumber
			if line == 0 {
				line = 1
			}

			fmt.Printf("::%s file=%s,line=%d::%s: %s\n",
				level, result.FilePath, line, hint.Issue, hint.Impact)
		}

		for _, issue := range result.Analysis.SecurityIssues {
			level := "warning"
			if issue.Severity == "high" {
				level = "error"
			}

			line := issue.LineNumber
			if line == 0 {
				line = 1
			}

			fmt.Printf("::%s file=%s,line=%d::%s: %s\n",
				level, result.FilePath, line, issue.Type, issue.Description)
		}
	}
}

func convertToSARIFResults(results []analysis.FileAnalysis) []map[string]interface{} {
	var sarifResults []map[string]interface{}

	for _, result := range results {
		// Performance issues
		for _, hint := range result.Analysis.PerformanceHints {
			sarifResult := map[string]interface{}{
				"ruleId":  "performance/" + strings.ReplaceAll(hint.Issue, " ", "_"),
				"level":   mapSeverityToSARIF(hint.Severity),
				"message": map[string]string{"text": hint.Impact},
				"locations": []map[string]interface{}{
					{
						"physicalLocation": map[string]interface{}{
							"artifactLocation": map[string]string{"uri": result.FilePath},
							"region": map[string]int{
								"startLine": maxInt(hint.LineNumber, 1),
							},
						},
					},
				},
			}
			sarifResults = append(sarifResults, sarifResult)
		}

		// Security issues
		for _, issue := range result.Analysis.SecurityIssues {
			sarifResult := map[string]interface{}{
				"ruleId":  "security/" + issue.Type,
				"level":   mapSeverityToSARIF(issue.Severity),
				"message": map[string]string{"text": issue.Description},
				"locations": []map[string]interface{}{
					{
						"physicalLocation": map[string]interface{}{
							"artifactLocation": map[string]string{"uri": result.FilePath},
							"region": map[string]int{
								"startLine": maxInt(issue.LineNumber, 1),
							},
						},
					},
				},
			}
			sarifResults = append(sarifResults, sarifResult)
		}
	}

	return sarifResults
}

func mapSeverityToSARIF(severity string) string {
	switch severity {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}

func setGitHubActionOutputs(results []analysis.FileAnalysis) {
	totalIssues := countTotalIssues(results)
	avgScore := calculateAverageScore(results)

	// Set GitHub Action outputs using the new format
	fmt.Printf("::set-output name=issues-found::%d\n", totalIssues)
	fmt.Printf("::set-output name=performance-score::%d\n", avgScore)
	fmt.Printf("::set-output name=files-analyzed::%d\n", len(results))

	// Output JSON results
	jsonData, _ := json.Marshal(results)
	fmt.Printf("::set-output name=analysis-results::%s\n", string(jsonData))
}

func hasHighSeverityIssues(results []analysis.FileAnalysis, threshold string) bool {
	for _, result := range results {
		for _, hint := range result.Analysis.PerformanceHints {
			if shouldFailOnSeverity(hint.Severity, threshold) {
				return true
			}
		}
		for _, issue := range result.Analysis.SecurityIssues {
			if shouldFailOnSeverity(issue.Severity, threshold) {
				return true
			}
		}
	}
	return false
}

func shouldFailOnSeverity(issueSeverity, threshold string) bool {
	severityLevels := map[string]int{
		"low":      1,
		"medium":   2,
		"high":     3,
		"critical": 4,
	}

	return severityLevels[issueSeverity] >= severityLevels[threshold]
}

func countTotalIssues(results []analysis.FileAnalysis) int {
	total := 0
	for _, result := range results {
		total += len(result.Analysis.SecurityIssues)
		total += len(result.Analysis.PerformanceHints)
	}
	return total
}

func calculateAverageScore(results []analysis.FileAnalysis) int {
	if len(results) == 0 {
		return 100
	}

	total := 0
	for _, result := range results {
		total += result.Analysis.PerformanceScore
	}
	return total / len(results)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
