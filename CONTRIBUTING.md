# Contributing to tv

Thank you for your interest in contributing to tv! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.19 or later
- Git
- (Optional) golangci-lint for linting
- (Optional) goreleaser for releases

### Getting Started

1. Fork and clone the repository:
```bash
git clone https://github.com/YOUR_USERNAME/tv.git
cd tv
```

2. Install dependencies:
```bash
make deps
```

3. Build the project:
```bash
make build
```

4. Run tests:
```bash
make test
```

## Development Workflow

### Making Changes

1. Create a new branch for your feature or bugfix:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and ensure tests pass:
```bash
make check  # Runs tests and linters
```

3. Commit your changes with clear commit messages:
```bash
git commit -m "feat: add new feature"
```

We follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions or changes
- `chore:` Maintenance tasks
- `refactor:` Code refactoring

4. Push to your fork and create a pull request:
```bash
git push origin feature/your-feature-name
```

## Available Make Commands

```bash
make build      # Build the binary
make test       # Run tests
make lint       # Run linters
make check      # Run tests and linters
make install    # Install to /usr/local/bin
make clean      # Remove build artifacts
make snapshot   # Create local test release
make help       # Show all commands
```

## Testing

- Write tests for new features
- Ensure all tests pass before submitting PR
- Aim for good test coverage
- Test files should be named `*_test.go`

Run tests with:
```bash
go test -v ./...
# or
make test
```

## Code Style

- Follow standard Go formatting (run `gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and manageable

Run linter:
```bash
make lint
```

## Pull Request Process

1. Update README.md with details of changes if needed
2. Update tests and ensure they pass
3. Ensure the PR description clearly describes the problem and solution
4. Link any relevant issues
5. Wait for review and address any feedback

## Release Process (Maintainers Only)

### Creating a Release

1. Ensure all tests pass and code is merged to master
2. Tag the release:
```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

3. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create GitHub release
   - Publish to package managers (Homebrew, Snap, Scoop)
   - Update checksums

### Manual Release (if needed)

```bash
# Requires GITHUB_TOKEN environment variable
export GITHUB_TOKEN=your_token
make release
```

### Testing Release Locally

```bash
make snapshot
# Binaries will be in ./dist/
```

## Package Distribution

The project automatically publishes to:

- **GitHub Releases**: Binary downloads
- **Homebrew**: `brew install codechenx/tv/tv`
- **Snap**: `snap install codechenx-tv`
- **Scoop (Windows)**: `scoop install tv`
- **Go**: `go install github.com/codechenx/tv@latest`
- **Debian/Ubuntu**: `.deb` packages
- **CentOS/Fedora**: `.rpm` packages
- **Alpine**: `.apk` packages

Configuration is in `.goreleaser.yml`.

## Getting Help

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Be respectful and constructive in discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
