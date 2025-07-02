# Contributing to API Performance Analyzer

Thank you for your interest in contributing to API Performance Analyzer! 🎉

## 🚀 Getting Started

### Prerequisites
- Go 1.23 or later
- Docker (for containerized testing)
- Git

### Local Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/aldookware/api-performance-analyzer.git
   cd api-performance-analyzer
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Build and test**
   ```bash
   # Build the CLI tool
   go build -o analyzer ./cmd/analyzer
   
   # Run tests
   go test ./...
   
   # Test locally
   ./analyzer --path=. --format=markdown --verbose
   ```

4. **Build Docker image**
   ```bash
   docker build -t api-analyzer-dev .
   ```

## 🛠️ Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow Go conventions and best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   # Run unit tests
   go test ./...
   
   # Test CLI functionality
   ./analyzer --path=./examples --format=json
   
   # Test Docker build
   docker build -t test-analyzer .
   ```

4. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: add new detection for xyz pattern"
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request**
   - Use a clear, descriptive title
   - Include a detailed description of changes
   - Reference any related issues

### Code Structure

```
api-performance-analyzer/
├── cmd/analyzer/          # CLI entrypoint
│   └── main.go           # Command-line interface
├── internal/analysis/     # Core analysis logic
│   └── analyzer.go       # Main analysis engine
├── .github/workflows/     # CI/CD and sample workflows
├── static/               # Web UI assets
├── action.yml           # GitHub Action definition
├── Dockerfile           # Container definition
└── entrypoint.sh        # Docker entrypoint script
```

## 📝 Coding Guidelines

### Go Code Standards
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions
- Handle errors explicitly
- Use meaningful variable names

### Adding New Detection Rules

1. **Define the pattern** in `internal/analysis/analyzer.go`
2. **Add detection logic** to the appropriate function
3. **Include severity and suggestions**
4. **Add test cases**
5. **Update documentation**

Example:
```go
// Add to analyzeSecurityIssues function
if strings.Contains(content, "password=") {
    issues = append(issues, SecurityIssue{
        Type:        "hardcoded_password",
        Description: "Hardcoded password detected",
        Severity:    "critical",
        LineNumber:  lineNumber,
        Suggestion:  "Use environment variables or secure vault",
    })
}
```

### Testing Guidelines

- Write unit tests for new functions
- Include edge cases and error scenarios
- Use table-driven tests when appropriate
- Test with realistic Go code examples

Example test:
```go
func TestDetectHardcodedPasswords(t *testing.T) {
    tests := []struct {
        name     string
        code     string
        expected bool
    }{
        {
            name:     "detects hardcoded password",
            code:     `password := "secret123"`,
            expected: true,
        },
        {
            name:     "ignores password variable",
            code:     `password := os.Getenv("PASSWORD")`,
            expected: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := detectHardcodedPasswords(tt.code)
            if (len(result) > 0) != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, len(result) > 0)
            }
        })
    }
}
```

## 🎯 Areas for Contribution

### High Priority
- **New Detection Rules** - Add more performance/security patterns
- **Framework Support** - Echo, Fiber, Chi support
- **Output Formats** - SARIF, XML, custom formats
- **Performance Improvements** - Faster analysis for large codebases

### Medium Priority
- **Better Error Messages** - More helpful suggestions
- **Configuration Options** - Custom rules, ignore patterns
- **Documentation** - More examples and tutorials
- **Test Coverage** - Improve test coverage

### Low Priority
- **Web UI Enhancements** - Better dashboard interface
- **Integration Examples** - More CI/CD examples
- **Localization** - Multi-language support

## 🐛 Reporting Issues

### Bug Reports
When reporting bugs, please include:
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Sample code that triggers the issue

### Feature Requests
For new features, please describe:
- Use case and motivation
- Proposed solution
- Alternatives considered
- Implementation details (if known)

## 📋 Pull Request Checklist

Before submitting a PR, ensure:

- [ ] Code follows Go conventions
- [ ] Tests pass (`go test ./...`)
- [ ] New functionality includes tests
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] No breaking changes (or clearly documented)
- [ ] Performance impact is considered

## 🎉 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes for significant contributions
- GitHub contributors page

## 📞 Getting Help

- **Questions**: Open a GitHub Discussion
- **Issues**: Create a GitHub Issue
- **Ideas**: Join our community discussions

Thank you for helping make API Performance Analyzer better! 🚀
