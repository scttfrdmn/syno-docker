# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of SynoDeploy CLI tool
- SSH connection management with key authentication
- Docker client setup over SSH tunnels
- Container deployment with `synodeploy run` command
- docker-compose support with `synodeploy deploy` command
- Container management with `synodeploy ps` and `synodeploy rm` commands
- Configuration management in `~/.synodeploy/config.yaml`
- Comprehensive unit tests with >90% coverage
- Go Report Card A+ quality checks (gofmt, govet, golint, staticcheck, etc.)
- Git pre-commit and pre-push hooks with quality enforcement
- Makefile with build automation and quality checks
- Cross-platform build support (macOS, Linux)
- DSM 7.2+ Container Manager support
- PATH resolution for Docker binary (`/usr/local/bin/docker`)
- Volume path validation and Synology-specific path handling
- Environment variable expansion in compose files
- Port mapping validation and conflict detection
- Restart policy support (no, always, unless-stopped, on-failure)
- Network mode configuration
- User and working directory specification
- Command override support
- Comprehensive error handling with actionable messages

### Security
- SSH key-only authentication (no password storage)
- Input validation for all user-provided data
- Path traversal prevention for volume mappings
- Sensitive information scanning in git hooks
- Secure defaults for container deployment

## [0.1.0] - 2025-01-XX

### Added
- Initial release of SynoDeploy
- Basic container deployment functionality
- SSH connection management
- Configuration system
- Command-line interface with cobra framework

---

## Release Checklist

When releasing a new version:

- [ ] Update version in this file
- [ ] Update version in `.goreleaser.yaml`
- [ ] Update version in documentation
- [ ] Run `make release` to ensure all checks pass
- [ ] Create git tag: `git tag -a v<version> -m "Release v<version>"`
- [ ] Push tag: `git push origin v<version>`
- [ ] GitHub Actions will automatically create release and update Homebrew formula