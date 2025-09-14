# SynoDeploy

[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/synodeploy)](https://goreportcard.com/report/github.com/scttfrdmn/synodeploy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/badge/release-v0.1.0-blue.svg)](https://github.com/scttfrdmn/synodeploy/releases/tag/v0.1.0)

**SynoDeploy** is a CLI tool that simplifies Docker container deployment to Synology NAS devices running DSM 7.2+. It handles SSH connection management, Docker client setup, and path resolution issues specific to Synology Container Manager.

## Features

- 🚀 **One-command deployment** - Deploy containers as easily as `synodeploy run nginx`
- 🔐 **SSH key authentication** - Secure connection using your existing SSH keys
- 📦 **docker-compose support** - Deploy complex multi-container applications
- 🎯 **DSM 7.2+ optimized** - Built specifically for Container Manager
- 🔧 **PATH resolution** - Automatically handles Docker binary location issues
- 📂 **Volume path helpers** - Smart handling of Synology volume paths
- 🔄 **Container lifecycle** - Deploy, list, and remove containers easily
- ⚡ **Single binary** - No dependencies, just download and use

## Quick Start

### Installation

```bash
# Install via Homebrew (recommended)
brew install scttfrdmn/tap/synodeploy

# Or download binary from releases
curl -L https://github.com/scttfrdmn/synodeploy/releases/latest/download/synodeploy-$(uname -s)-$(uname -m) -o synodeploy
chmod +x synodeploy
sudo mv synodeploy /usr/local/bin/
```

### Setup

```bash
# One-time setup - connect to your Synology NAS
synodeploy init 192.168.1.100
```

This will:
- Test SSH connection to your NAS
- Verify Container Manager is running
- Save connection details to `~/.synodeploy/config.yaml`

### Deploy Your First Container

```bash
# Deploy Nginx web server
synodeploy run nginx:latest \
  --name web-server \
  --port 8080:80 \
  --volume /volume1/web:/usr/share/nginx/html

# Deploy from docker-compose.yml
synodeploy deploy ./docker-compose.yml

# List running containers
synodeploy ps

# Remove container
synodeploy rm web-server
```

## Commands

### `synodeploy init <host>`

Setup connection to your Synology NAS.

```bash
synodeploy init 192.168.1.100 \
  --user admin \
  --port 22 \
  --key ~/.ssh/id_rsa \
  --volume-path /volume1/docker
```

### `synodeploy run <image>`

Deploy a single container.

```bash
synodeploy run postgres:13 \
  --name database \
  --port 5432:5432 \
  --volume /volume1/postgres:/var/lib/postgresql/data \
  --env POSTGRES_PASSWORD=secretpassword \
  --restart unless-stopped
```

**Options:**
- `--name` - Container name (auto-generated if not specified)
- `--port` - Port mappings (format: `host:container`)
- `--volume` - Volume mappings (format: `host:container`)
- `--env` - Environment variables (format: `KEY=value`)
- `--restart` - Restart policy (`no`, `always`, `unless-stopped`, `on-failure`)
- `--network` - Network mode (default: `bridge`)
- `--user` - User to run container as (format: `uid:gid`)
- `--workdir` - Working directory inside container
- `--command` - Command to run in container

### `synodeploy deploy <compose-file>`

Deploy from docker-compose.yml file.

```bash
synodeploy deploy ./docker-compose.yml \
  --project my-app \
  --env-file .env
```

**Supported compose features:**
- Multi-service deployments
- Port mappings
- Volume mounts
- Environment variables
- Environment variable substitution
- Restart policies
- Networks (basic support)
- Dependencies (deployment order only)

### `synodeploy ps`

List containers.

```bash
# Show running containers
synodeploy ps

# Show all containers (including stopped)
synodeploy ps --all
```

### `synodeploy rm <container>`

Remove container.

```bash
# Remove stopped container
synodeploy rm web-server

# Force remove running container
synodeploy rm web-server --force
```

## Configuration

SynoDeploy stores configuration in `~/.synodeploy/config.yaml`:

```yaml
host: 192.168.1.100
port: 22
user: admin
ssh_key_path: /home/user/.ssh/id_rsa
defaults:
  volume_path: /volume1/docker
  network: bridge
```

## Volume Path Handling

SynoDeploy automatically handles Synology volume paths:

```bash
# These are equivalent:
synodeploy run nginx -v /volume1/web:/usr/share/nginx/html
synodeploy run nginx -v ./web:/usr/share/nginx/html  # Expands to /volume1/docker/web
synodeploy run nginx -v web:/usr/share/nginx/html    # Expands to /volume1/docker/web
```

## Requirements

### Synology NAS
- DSM 7.2 or later
- Container Manager installed and running
- SSH access enabled
- User with administrator privileges

### Local Machine
- SSH key pair configured
- Network access to your NAS

## Troubleshooting

### Connection Issues

```bash
# Test SSH connection manually
ssh admin@192.168.1.100

# Check if Container Manager is running
ssh admin@192.168.1.100 'systemctl status pkg-ContainerManager-dockerd'
```

### Docker Command Not Found

This is automatically handled by SynoDeploy, but if you see this error, it means:
- Container Manager is not installed/running
- Your user doesn't have the right permissions
- There's a PATH issue (SynoDeploy handles this automatically)

### Permission Denied

```bash
# Ensure your user is in the docker group
ssh admin@192.168.1.100 'sudo synogroup --member docker admin'
```

### Port Already in Use

```bash
# Check what's using the port
ssh admin@192.168.1.100 'netstat -tlnp | grep :8080'
```

## Development

### Building from Source

```bash
git clone https://github.com/scttfrdmn/synodeploy.git
cd synodeploy
make build
```

### Running Tests

```bash
make test              # Run tests
make quality-check     # Run all quality checks
make coverage         # Generate coverage report
```

### Quality Checks

SynoDeploy maintains Go Report Card A+ grade with:

- `gofmt` - Code formatting
- `go vet` - Static analysis
- `golangci-lint` - Comprehensive linting
- `staticcheck` - Advanced static analysis
- `ineffassign` - Ineffectual assignment detection
- `misspell` - Spelling mistakes
- `gocyclo` - Cyclomatic complexity

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run quality checks (`make quality-check`)
5. Run tests (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Docker SDK](https://github.com/docker/docker) - Docker client library
- [SSH package](https://golang.org/x/crypto/ssh) - SSH client implementation
- Synology Community - For documenting DSM 7.2+ changes

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/scttfrdmn/synodeploy/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/synodeploy/discussions)
- 📧 Email: support@synodeploy.com

---

**Made with ❤️ for the Synology community**