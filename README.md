# syno-docker

[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/syno-docker)](https://goreportcard.com/report/github.com/scttfrdmn/syno-docker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/badge/release-v0.2.1-blue.svg)](https://github.com/scttfrdmn/syno-docker/releases/tag/v0.2.1)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)
[![Integration Tests](https://img.shields.io/badge/integration%20tests-passing-brightgreen.svg)](#integration-tests)

**syno-docker** is a comprehensive Docker management CLI tool for Synology NAS devices running DSM 7.2+. It provides the full Docker workflow - from image management to container lifecycle to system maintenance - all optimized for Synology Container Manager.

**✅ Verified Working** on real Synology hardware with comprehensive integration testing.

## Features

### Core Deployment
- 🚀 **One-command deployment** - Deploy containers as easily as `syno-docker run nginx`
- 📦 **Docker Compose support** - Deploy complex multi-container applications
- 🔧 **PATH resolution** - Automatically handles Docker binary location issues
- 📂 **Volume path helpers** - Smart handling of Synology volume paths

### Container Management
- 🔄 **Complete lifecycle** - Start, stop, restart, remove containers
- 📋 **Container inspection** - Detailed container information and logs
- 🖥️ **Interactive execution** - Run commands inside containers (`exec`)
- 📊 **Resource monitoring** - Real-time container statistics

### Image & System Management
- 🏗️ **Image operations** - Pull, list, remove images with advanced filtering
- 📦 **Volume management** - Create, list, inspect, and clean up volumes
- 🌐 **Network management** - Create, list, inspect networks; connect/disconnect containers
- 🧹 **System maintenance** - Disk usage, system info, and cleanup tools
- 📤 **Import/Export** - Backup and restore containers

### Infrastructure
- 🔐 **SSH key & ssh-agent support** - Works with both SSH key files and ssh-agent
- 👤 **Administrator user support** - Compatible with both `admin` and custom admin users
- 🎯 **DSM 7.2+ optimized** - Built specifically for Container Manager
- ⚡ **Single binary** - No dependencies, just download and use
- 🧪 **Integration tested** - Verified on real Synology hardware

## Quick Start

### Installation

```bash
# Install via Homebrew (recommended)
brew install scttfrdmn/tap/syno-docker

# Or download binary from releases
curl -L https://github.com/scttfrdmn/syno-docker/releases/latest/download/syno-docker-$(uname -s)-$(uname -m) -o syno-docker
chmod +x syno-docker
sudo mv syno-docker /usr/local/bin/
```

### Setup

```bash
# One-time setup - connect to your Synology NAS
syno-docker init 192.168.1.100

# Or with custom admin username (if not using 'admin')
syno-docker init your-nas.local --user your-username

# For ssh-agent users (automatically detected)
syno-docker init your-nas.local --user your-username
```

This will:
- Test SSH connection to your NAS (supports both SSH keys and ssh-agent)
- Verify Container Manager is running
- Test Docker command execution
- Save connection details to `~/.syno-docker/config.yaml`

### Deploy Your First Container

```bash
# Deploy Nginx web server
syno-docker run nginx:latest \
  --name web-server \
  --port 8080:80 \
  --volume /volume1/web:/usr/share/nginx/html

# Deploy from docker-compose.yml
syno-docker deploy ./docker-compose.yml

# List running containers and monitor resources
syno-docker ps
syno-docker stats

# Get container logs and execute commands
syno-docker logs web-server --follow
syno-docker exec web-server /bin/bash

# Remove container when done
syno-docker rm web-server
```

## Commands Overview

syno-docker provides **20+ commands** covering the complete Docker workflow:

### **Container Lifecycle**
- `syno-docker run` - Deploy single containers with full configuration options
- `syno-docker ps` - List containers (running/all) with detailed status
- `syno-docker start/stop/restart` - Control container state
- `syno-docker rm` - Remove containers (with force option)

### **Container Operations**
- `syno-docker logs` - View container logs (follow, tail, timestamps)
- `syno-docker exec` - Execute commands inside containers (interactive/non-interactive)
- `syno-docker stats` - Real-time resource usage statistics
- `syno-docker inspect` - Detailed container/image/volume information

### **Image Management**
- `syno-docker pull` - Pull images from registries (platform-specific, all tags)
- `syno-docker images` - List images (all, dangling, with digests)
- `syno-docker rmi` - Remove images (force, preserve parents)
- `syno-docker import/export` - Backup and restore containers

### **Volume Management**
- `syno-docker volume ls` - List volumes with driver information
- `syno-docker volume create` - Create volumes with custom drivers/labels
- `syno-docker volume rm` - Remove volumes (with force)
- `syno-docker volume inspect` - Detailed volume information
- `syno-docker volume prune` - Clean unused volumes

### **Network Management**
- `syno-docker network ls` - List networks with filtering options
- `syno-docker network create` - Create custom networks with CIDR, gateways
- `syno-docker network rm` - Remove networks
- `syno-docker network inspect` - Detailed network information
- `syno-docker network connect/disconnect` - Attach/detach containers
- `syno-docker network prune` - Clean unused networks

### **System Operations**
- `syno-docker system df` - Show Docker disk usage
- `syno-docker system info` - Display Docker system information
- `syno-docker system prune` - Clean unused containers, images, networks

### **Multi-Container Applications**
- `syno-docker deploy` - Deploy from docker-compose.yml files
- `syno-docker init` - Setup connection to Synology NAS

### **Key Command Examples**

```bash
# Container lifecycle
syno-docker run nginx:latest --name web --port 80:80 --restart unless-stopped
syno-docker logs web --follow --timestamps
syno-docker exec -it web /bin/bash
syno-docker restart web
syno-docker stop web && syno-docker rm web

# Image management
syno-docker pull postgres:13 --platform linux/arm64
syno-docker images --dangling
syno-docker rmi old-image --force

# Volume operations
syno-docker volume create my-data --driver local
syno-docker volume ls --quiet
syno-docker volume inspect my-data
syno-docker volume rm my-data --force

# Network operations
syno-docker network create my-app-net --driver bridge --subnet 172.20.0.0/16
syno-docker network ls --filter driver=bridge
syno-docker network connect my-app-net web-server --alias web
syno-docker network disconnect my-app-net web-server

# System maintenance
syno-docker system df --verbose
syno-docker system prune --all --volumes --force
syno-docker stats --all --no-stream

# Force remove running container
syno-docker rm web-server --force
```

## Configuration

syno-docker stores configuration in `~/.syno-docker/config.yaml`:

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

syno-docker automatically handles Synology volume paths:

```bash
# These are equivalent:
syno-docker run nginx -v /volume1/web:/usr/share/nginx/html
syno-docker run nginx -v ./web:/usr/share/nginx/html  # Expands to /volume1/docker/web
syno-docker run nginx -v web:/usr/share/nginx/html    # Expands to /volume1/docker/web
```

## Requirements

### Synology NAS
- DSM 7.2 or later
- Container Manager installed and running
- SSH access enabled (Control Panel → Terminal & SNMP)
- User with administrator privileges and docker group membership

### Local Machine
- SSH key pair configured OR ssh-agent running
- Network access to your NAS
- Go 1.21+ (for building from source)

## Troubleshooting

### Connection Issues

```bash
# Test SSH connection manually
ssh admin@192.168.1.100

# Check if Container Manager is running
ssh admin@192.168.1.100 'systemctl status pkg-ContainerManager-dockerd'
```

### Docker Command Not Found

This is automatically handled by syno-docker, but if you see this error, it means:
- Container Manager is not installed/running
- Your user doesn't have the right permissions
- There's a PATH issue (syno-docker handles this automatically)

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
git clone https://github.com/scttfrdmn/syno-docker.git
cd syno-docker
make build
```

### Running Tests

```bash
make test              # Run unit tests
make quality-check     # Run all quality checks
make coverage         # Generate coverage report
```

### Integration Tests

syno-docker includes comprehensive integration tests that validate functionality against real Synology hardware:

```bash
# Test connection to your NAS
go test -v -run TestConnectionToChubChub ./tests/integration/

# Test Docker commands via SSH
go test -v -run TestDirectDockerCommands ./tests/integration/

# Full end-to-end testing
go test -v -run TestSynoDockerEndToEnd ./tests/integration/

# All integration tests
go test -v ./tests/integration/
```

**Integration test coverage:**
- ✅ SSH connectivity and authentication (ssh-agent + key file)
- ✅ Docker command execution over SSH
- ✅ Container deployment, lifecycle, and removal
- ✅ HTTP endpoint validation for deployed services
- ✅ Volume mounting and file system access
- ✅ Error handling for invalid configurations
- ✅ Container Manager service status validation

### Quality Checks

syno-docker maintains Go Report Card A+ grade with:

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
- 🐛 [Issue Tracker](https://github.com/scttfrdmn/syno-docker/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/syno-docker/discussions)
- 📧 Email: support@syno-docker.com

---

**Made with ❤️ for the Synology community**