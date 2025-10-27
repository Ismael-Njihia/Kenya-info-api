# Contributing to Kenya Info API

Thank you for your interest in contributing to the Kenya Info API! This document provides guidelines and instructions for contributing.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

1. **Fork the Repository**
   - Click the "Fork" button at the top right of the repository page

2. **Clone Your Fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/Kenya-info-api.git
   cd Kenya-info-api
   ```

3. **Add Upstream Remote**
   ```bash
   git remote add upstream https://github.com/Ismael-Njihia/Kenya-info-api.git
   ```

## Development Setup

1. **Install Prerequisites**
   - Go 1.21 or higher
   - MongoDB (or MongoDB Atlas account)
   - Docker (optional)

2. **Install Dependencies**
   ```bash
   make install
   ```

3. **Set Up Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Generate Swagger Documentation**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   make swagger
   ```

5. **Run Tests**
   ```bash
   make test
   ```

6. **Run the Application**
   ```bash
   make run
   ```

## How to Contribute

### Reporting Bugs

1. **Check Existing Issues**
   - Search existing issues to avoid duplicates

2. **Create a New Issue**
   - Use a clear and descriptive title
   - Describe the steps to reproduce the bug
   - Include error messages and logs
   - Specify your environment (OS, Go version, etc.)

### Suggesting Features

1. **Check Existing Issues**
   - See if the feature has already been suggested

2. **Create a Feature Request**
   - Use a clear and descriptive title
   - Explain the problem the feature would solve
   - Describe the proposed solution
   - Consider alternative solutions

### Code Contributions

1. **Find an Issue**
   - Look for issues labeled `good first issue` or `help wanted`
   - Comment on the issue to let others know you're working on it

2. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Your Changes**
   - Write clean, maintainable code
   - Follow the coding standards
   - Add tests for new functionality
   - Update documentation as needed

4. **Test Your Changes**
   ```bash
   make test
   make build
   ```

5. **Commit Your Changes**
   - Follow the commit guidelines below
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

6. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**
   - Go to the original repository
   - Click "New Pull Request"
   - Select your fork and branch
   - Fill in the PR template

## Coding Standards

### Go Style Guide

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `golangci-lint` before committing

### Code Organization

```
internal/
├── handlers/     # HTTP request handlers
├── middleware/   # HTTP middleware
├── models/       # Data models
├── repository/   # Database operations
└── logger/       # Logging utilities
```

### Naming Conventions

- **Files**: Use lowercase with underscores (e.g., `county_handler.go`)
- **Packages**: Use lowercase, single word (e.g., `handlers`, `models`)
- **Functions**: Use camelCase (e.g., `GetCountyByID`)
- **Exported**: Start with uppercase (e.g., `County`, `CreateCounty`)
- **Unexported**: Start with lowercase (e.g., `parseRequest`)

### Error Handling

```go
if err != nil {
    logger.Log.Error("Failed to process request", zap.Error(err))
    return fmt.Errorf("process request: %w", err)
}
```

### Logging

Use structured logging with Zap:

```go
logger.Log.Info("Processing request",
    zap.String("id", id),
    zap.String("action", "create"),
)
```

### Testing

- Write unit tests for all new code
- Aim for at least 80% code coverage
- Use table-driven tests when appropriate
- Mock external dependencies

Example:
```go
func TestCountyHandler_GetByID(t *testing.T) {
    // Setup
    // Test
    // Assert
}
```

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples
```
feat(counties): add search functionality
fix(database): resolve connection timeout issue
docs(readme): update installation instructions
test(handlers): add unit tests for ward handler
```

## Pull Request Process

1. **Update Documentation**
   - Update README.md if needed
   - Add/update API documentation
   - Include code comments for complex logic

2. **Run Tests**
   ```bash
   make test
   make lint
   ```

3. **Update Swagger Docs**
   ```bash
   make swagger
   ```

4. **Create Pull Request**
   - Use a clear title following commit guidelines
   - Fill in the PR template completely
   - Link related issues

5. **Code Review**
   - Address reviewer feedback
   - Keep discussions professional and constructive
   - Make requested changes in new commits

6. **Merge**
   - PRs require at least one approval
   - Maintainer will merge once approved

## Development Workflow

### Syncing Your Fork

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

### Working on Multiple Features

```bash
# Create feature branch from main
git checkout main
git pull upstream main
git checkout -b feature/new-feature

# Work on feature
# ...

# Create another feature branch
git checkout main
git checkout -b feature/another-feature
```

## Testing Guidelines

### Unit Tests

- Test individual functions and methods
- Mock external dependencies
- Use testify for assertions

### Integration Tests

- Test complete request/response cycles
- Use test database
- Clean up test data

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/handlers/...
```

## Documentation

### Code Comments

- Add comments for exported functions
- Explain complex logic
- Use godoc format

Example:
```go
// GetCountyByID retrieves a county by its unique identifier.
// It returns an error if the county is not found or if there's a database error.
func (h *CountyHandler) GetCountyByID(c *gin.Context) {
    // Implementation
}
```

### API Documentation

- Use Swagger annotations for all endpoints
- Include request/response examples
- Document error cases

## Need Help?

- Ask questions in GitHub Issues
- Tag issues with `question` label
- Reach out to maintainers

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project README

Thank you for contributing to Kenya Info API! 🎉
