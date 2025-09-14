# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-09-14

### Added
- **Core CLI Framework**: Complete command-line interface with `init`, `run`, `deploy`, `ps`, `rm` commands
- **SSH Connection Management**: Secure SSH key authentication with automatic Docker binary path resolution
- **Container Deployment**: Single container deployment with `synodeploy run` supporting ports, volumes, environment variables
- **Docker Compose Support**: Multi-container deployment with `synodeploy deploy` for compose files
- **Container Management**: List containers with `synodeploy ps` and remove with `synodeploy rm`
- **Configuration System**: Persistent configuration in `~/.synodeploy/config.yaml` with validation
- **DSM 7.2+ Optimization**: Built specifically for Container Manager with known constants and paths
- **Volume Path Helpers**: Smart volume path handling with Synology volume validation
- **Environment Variable Expansion**: Support for ${VAR} substitution in compose files
- **Input Validation**: Comprehensive validation for Docker images, container names, ports, volumes
- **Quality Assurance**: Go Report Card A+ compliance with automated quality checks
- **Cross-Platform Builds**: Support for macOS and Linux (AMD64, ARM64)
- **Professional Documentation**: Complete README, installation guide, and usage documentation

### Security
- SSH key-only authentication (no password storage)
- Input validation and sanitization for all user data
- Path traversal prevention for volume mappings
- Secure defaults for container deployment
- Sensitive information scanning in git hooks

### Developer Experience
- Comprehensive unit tests with high coverage
- Git pre-commit and pre-push hooks with quality enforcement
- Makefile with complete build automation
- Goreleaser configuration for automated releases
- GitHub Actions workflows for CI/CD
- MIT license with proper copyright attribution

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