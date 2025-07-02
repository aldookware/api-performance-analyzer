
# 🚀 API Performance Analyzer

**Detect and fix API performance bottlenecks in Go REST APIs before they hit production**

[![GitHub Marketplace](https://img.shields.io/badge/GitHub-Marketplace-blue?logo=github)](https://github.com/marketplace/actions/api-performance-analyzer)
[![Docker](https://img.shields.io/badge/Docker-Available-blue?logo=docker)](https://hub.docker.com/r/aldookware/api-performance-analyzer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Transform your Go APIs into high-performance, production-ready services with automated performance analysis and actionable recommendations.

## 🎯 What Does It Detect?

### ⚡ Performance Issues
- **N+1 Database Queries** - Detect GORM queries that could be optimized with preloading
- **Missing Database Indexes** - Identify tables that need indexes for faster queries  
- **Large Response Payloads** - Flag APIs returning excessive data
- **Memory Leaks** - Spot in-memory storage that doesn't scale
- **Inefficient Algorithms** - Find O(n²) loops and unoptimized code

### 🔒 Security Vulnerabilities  
- **SQL Injection Risks** - Raw SQL queries without parameterization
- **Missing CORS Protection** - APIs vulnerable to cross-origin attacks
- **Hardcoded Secrets** - Credentials and API keys in source code
- **Input Validation Gaps** - Endpoints missing proper validation

### 📊 Code Quality Issues
- **Missing Error Handling** - Unhandled error conditions
- **Poor API Design** - REST API anti-patterns
- **Lack of Rate Limiting** - APIs without protection against abuse
- **Missing Logging/Monitoring** - Observability gaps

## � Quick Start

### GitHub Actions (Recommended)

Add this to your `.github/workflows/api-performance.yml`:

```yaml
name: API Performance Check
on: [push, pull_request]

jobs:
  performance:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Analyze API Performance
      uses: aldookware/api-performance-analyzer@v1
      with:
        code-path: '.'
        output-format: 'github'
        severity-threshold: 'medium'
        fail-on-issues: 'true'
```

### CLI Tool

```bash
# Download and run locally
go install github.com/aldookware/api-performance-analyzer/cmd/analyzer@latest

# Analyze your codebase
analyzer --path=. --format=markdown --verbose

# Generate JSON report
analyzer --path=./src --format=json --threshold=high
```

### Docker

```bash
# Pull and run
docker run --rm -v "$(pwd):/workspace" \
  aldookware/api-performance-analyzer:latest \
  /workspace markdown medium false

# Build from source
docker build -t api-analyzer .
docker run --rm -v "$(pwd):/code" api-analyzer /code json high true
```

## 📊 Output Formats

### GitHub Actions Annotations
Perfect for CI/CD - creates inline comments on your PRs:
```bash
::error file=main.go,line=45::N+1 Query detected - use Preload() instead
::warning file=handlers.go,line=12::Missing CORS middleware
```

### Markdown Report
Beautiful, shareable reports:
```markdown
# 🚀 API Performance Analysis Report

## 🟢 Overall Performance Score: 83/100

| File | Score | Issues |
|------|-------|--------|
| main.go | A+ (95/100) | 2 |
| handlers.go | B (78/100) | 5 |
```

### JSON Output
For integration with other tools:
```json
{
  "overall_score": 83,
  "files": [
    {
      "path": "main.go",
      "performance_score": 95,
      "grade": "A+",
      "issues": [...]
    }
  ]
}
```

### SARIF Format
Compatible with GitHub Security tab and other security tools.

## ⚙️ Configuration

### Input Parameters

| Parameter | Description | Default | Options |
|-----------|-------------|---------|---------|
| `code-path` | Path to analyze | `.` | Any valid path |
| `output-format` | Output format | `markdown` | `markdown`, `json`, `github`, `sarif` |
| `severity-threshold` | Minimum severity | `medium` | `low`, `medium`, `high`, `critical` |
| `fail-on-issues` | Fail CI on issues | `false` | `true`, `false` |

### Example Configurations

**Strict Mode** (fails on any high+ issues):
```yaml
- uses: aldookware/api-performance-analyzer@v1
  with:
    severity-threshold: 'high'
    fail-on-issues: 'true'
```

**Development Mode** (warnings only):
```yaml
- uses: aldookware/api-performance-analyzer@v1
  with:
    severity-threshold: 'low'
    fail-on-issues: 'false'
    output-format: 'markdown'
```

**Security Focus** (JSON output for processing):
```yaml
- uses: aldookware/api-performance-analyzer@v1
  with:
    output-format: 'json'
    severity-threshold: 'medium'
```

## 📈 Sample Analysis Results

### Before Optimization
```
🚨 Performance Score: 45/100 (F)

Critical Issues:
- N+1 Query in getUserOrders() - 💥 Each user triggers separate query
- Missing index on users.email - 🐌 Full table scan on every login
- 50MB response in /api/export - 📡 Massive payload kills mobile users

High Priority:
- SQL injection in searchUsers() - 🛡️ Use parameterized queries
- No rate limiting - 🚫 API abuse vulnerability
```

### After Applying Fixes
```
✅ Performance Score: 92/100 (A+)

Improvements Applied:
✓ Used Preload() for user orders - 🚀 40x faster queries
✓ Added database indexes - ⚡ Login time: 2s → 0.1s  
✓ Implemented pagination - 📱 Mobile-friendly responses
✓ Fixed SQL injection - 🛡️ Secure parameterized queries
✓ Added rate limiting - 🚫 Protected against abuse
```

## 🛠️ Advanced Usage

### Custom Rules
Create `.api-analyzer.yml` in your repo:
```yaml
rules:
  performance:
    n_plus_one: error
    missing_indexes: warning
    large_payloads: error
  security:
    sql_injection: error
    missing_cors: warning
    hardcoded_secrets: error

ignore_files:
  - "test/**"
  - "mock/**"
  - "vendor/**"

custom_patterns:
  - pattern: "SELECT.*WHERE.*="
    message: "Consider using prepared statements"
    severity: high
```

### Framework Support
Currently supports:
- ✅ **Gin** - Full support
- ✅ **GORM** - Database optimization
- ✅ **Standard Library** - Basic patterns
- 🚧 **Echo** - Coming soon
- 🚧 **Fiber** - Coming soon

### Integration Examples

**With SonarQube:**
```bash
analyzer --format=sarif --output=results.sarif
sonar-scanner -Dsonar.go.coverage.reportPaths=results.sarif
```

**With Slack Notifications:**
```yaml
- name: Notify on Performance Issues
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: failure
    text: "🚨 API Performance issues detected!"
```

## 🏗️ Local Development

### Building from Source
```bash
# Clone repository
git clone https://github.com/aldookware/api-performance-analyzer
cd api-performance-analyzer

# Build CLI
go build -o analyzer ./cmd/analyzer

# Build Docker image
docker build -t api-analyzer .

# Run tests
go test ./...
```

### Project Structure
```
api-performance-analyzer/
├── cmd/analyzer/          # CLI entrypoint
├── internal/analysis/     # Core analysis engine
├── .github/workflows/     # Sample workflows
├── static/               # Web UI (optional)
├── action.yml           # GitHub Action definition
├── Dockerfile           # Container definition
└── entrypoint.sh        # Docker entrypoint
```

### Testing Your Changes
```bash
# Test locally
./analyzer --path=./examples --format=markdown --verbose

# Test Docker build
docker build -t test-analyzer .
docker run --rm -v "$(pwd):/workspace" test-analyzer /workspace json medium false

# Test GitHub Action (requires act)
act -P ubuntu-latest=nektos/act-environments-ubuntu:18.04
```

## 📊 Performance Benchmarks

| Codebase Size | Analysis Time | Memory Usage |
|---------------|---------------|--------------|
| Small (< 1K LOC) | 0.5s | 10MB |
| Medium (1K-10K LOC) | 2-5s | 25MB |
| Large (10K-100K LOC) | 10-30s | 50MB |
| Enterprise (100K+ LOC) | 1-3min | 100MB |

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

### Quick Start for Contributors
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Ensure tests pass: `go test ./...`
5. Submit a pull request

### Reporting Issues
- 🐛 [Bug Report](https://github.com/aldookware/api-performance-analyzer/issues/new?template=bug_report.md)
- ✨ [Feature Request](https://github.com/aldookware/api-performance-analyzer/issues/new?template=feature_request.md)
- 📚 [Documentation Issue](https://github.com/aldookware/api-performance-analyzer/issues/new?template=docs.md)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Go AST parsing powered by the standard library
- Static analysis inspired by tools like golangci-lint
- Performance patterns from Google's Go style guide
- Security checks based on OWASP guidelines

---

**Ready to optimize your Go APIs? [Get started now!](https://github.com/marketplace/actions/api-performance-analyzer) 🚀**

<div align="center">
  <strong>Made with ❤️ for the Go community</strong>
</div>
