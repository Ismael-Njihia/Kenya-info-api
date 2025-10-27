# Contributing to Kenya Info API

Thank you for your interest in contributing to the Kenya Info API! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior vs. actual behavior
- Your environment (OS, Go version, etc.)
- Any relevant logs or screenshots

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:
- A clear description of the feature
- The problem it solves
- Possible implementation approach
- Any alternative solutions you've considered

### Pull Requests

1. **Fork the repository** and create your branch from `main`
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clean, readable code
   - Follow Go best practices and conventions
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   make test
   make lint
   ```

4. **Commit your changes**
   - Use clear, descriptive commit messages
   - Follow the conventional commits format:
     - `feat:` for new features
     - `fix:` for bug fixes
     - `docs:` for documentation changes
     - `test:` for test changes
     - `refactor:` for code refactoring
     - `chore:` for maintenance tasks

5. **Push to your fork** and submit a pull request
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Wait for review**
   - Address any feedback from maintainers
   - Make requested changes
   - Keep your PR up to date with the main branch

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ismael-Njihia/Kenya-info-api.git
   cd Kenya-info-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your MongoDB connection string
   ```

4. **Run the application**
   ```bash
   make run
   ```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for code formatting
- Run `golangci-lint` before submitting PR
- Write clear, self-documenting code
- Add comments for complex logic

### Project Structure

- `cmd/api/` - Application entry point
- `internal/` - Private application code
  - `config/` - Configuration management
  - `database/` - Database connection
  - `handlers/` - HTTP handlers
  - `middleware/` - HTTP middleware
  - `models/` - Data models
  - `services/` - Business logic
- `pkg/` - Public packages
- `docs/` - Documentation

### Testing

- Write tests for all new features
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Example test:
```go
func TestCountyService(t *testing.T) {
    // Arrange
    // Act
    // Assert
}
```

### Documentation

- Update README.md for significant changes
- Add godoc comments for public functions
- Update Swagger annotations for API changes
- Include examples where helpful

## Database Changes

When modifying data models:
1. Update the model in `internal/models/`
2. Update related services in `internal/services/`
3. Update handlers in `internal/handlers/`
4. Update tests
5. Document changes in PR description

## API Changes

When modifying API endpoints:
1. Update handler functions
2. Update Swagger annotations
3. Update tests
4. Document changes in PR description
5. Consider backward compatibility

## Questions?

If you have questions, feel free to:
- Open an issue
- Reach out to maintainers
- Check existing issues and discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to Kenya Info API! 🇰🇪
