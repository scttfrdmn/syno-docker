# Installation Guide

This guide covers various ways to install syno-docker on your system.

## Homebrew (Recommended)

The easiest way to install syno-docker is via Homebrew:

```bash
# Add the tap
brew tap scttfrdmn/syno-docker

# Install syno-docker
brew install syno-docker

# Verify installation
syno-docker version
```

### Updating via Homebrew

```bash
brew upgrade syno-docker
```

## Direct Download

Download the latest binary for your platform from the [releases page](https://github.com/scttfrdmn/syno-docker/releases):

### macOS

```bash
# Intel Mac
curl -L https://github.com/scttfrdmn/syno-docker/releases/latest/download/syno-docker-darwin-amd64 -o syno-docker

# Apple Silicon Mac
curl -L https://github.com/scttfrdmn/syno-docker/releases/latest/download/syno-docker-darwin-arm64 -o syno-docker

# Make executable and install
chmod +x syno-docker
sudo mv syno-docker /usr/local/bin/
```

### Linux

```bash
# Intel/AMD 64-bit
curl -L https://github.com/scttfrdmn/syno-docker/releases/latest/download/syno-docker-linux-amd64 -o syno-docker

# ARM 64-bit
curl -L https://github.com/scttfrdmn/syno-docker/releases/latest/download/syno-docker-linux-arm64 -o syno-docker

# Make executable and install
chmod +x syno-docker
sudo mv syno-docker /usr/local/bin/
```

## Build from Source

If you prefer to build from source:

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional)

### Building

```bash
# Clone the repository
git clone https://github.com/scttfrdmn/syno-docker.git
cd syno-docker

# Build using Make
make build

# Or build using Go directly
go build -o bin/syno-docker main.go

# Install to system
sudo cp bin/syno-docker /usr/local/bin/
```

### Development Build

For development purposes, you can create a symlink:

```bash
make dev-install
```

This creates a symlink in `/usr/local/bin/syno-docker` pointing to your development build.

## SSH Setup

syno-docker requires SSH access to your Synology NAS. It supports both ssh-agent and SSH key files.

### SSH Agent (Recommended)

If you use ssh-agent (most modern development environments):

```bash
# Verify ssh-agent is running and has keys
ssh-add -l

# Test connection to your NAS
ssh your-username@your-nas-hostname

# syno-docker will automatically use ssh-agent
syno-docker init your-nas-hostname --user your-username
```

### SSH Key Files

If you prefer using SSH key files directly:

```bash
# Generate SSH key if needed
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

# Copy public key to NAS
ssh-copy-id your-username@your-nas-hostname

# Specify key path in syno-docker
syno-docker init your-nas-hostname --user your-username --key ~/.ssh/id_rsa
```

## Verification

After installation, verify syno-docker is working correctly:

```bash
# Check version
syno-docker --version

# View help
syno-docker --help

# Initialize connection to your NAS
syno-docker init your-nas-hostname --user your-username
```

## System Requirements

### Local Machine

- **Operating System**: macOS 10.15+, Linux (any modern distribution)
- **Architecture**: x86_64 (Intel/AMD) or arm64 (Apple Silicon, ARM)
- **Network**: Access to your Synology NAS
- **SSH**: SSH client (usually pre-installed)

### Synology NAS

- **DSM Version**: 7.2 or later
- **Package**: Container Manager installed and running
- **SSH Access**: SSH service enabled
- **User Account**: Admin privileges required

## Troubleshooting Installation

### Command Not Found

If you get "command not found" after installation:

```bash
# Check if /usr/local/bin is in your PATH
echo $PATH

# Add to PATH if missing (add to your shell profile)
export PATH="/usr/local/bin:$PATH"
```

### Permission Denied

If you get permission errors:

```bash
# Check file permissions
ls -la /usr/local/bin/syno-docker

# Fix permissions if needed
sudo chmod +x /usr/local/bin/syno-docker
```

### Homebrew Issues

If Homebrew installation fails:

```bash
# Update Homebrew
brew update

# Check for issues
brew doctor

# Try installing again
brew install scttfrdmn/syno-docker/syno-docker
```

## Uninstallation

### Homebrew

```bash
brew uninstall syno-docker
brew untap scttfrdmn/tap
```

### Manual Installation

```bash
# Remove binary
sudo rm /usr/local/bin/syno-docker

# Remove configuration (optional)
rm -rf ~/.syno-docker
```

## Next Steps

After installation, proceed to the [Usage Guide](usage.md) to set up your first connection to your Synology NAS.