# Installation Guide

This guide covers various ways to install SynoDeploy on your system.

## Homebrew (Recommended)

The easiest way to install SynoDeploy is via Homebrew:

```bash
# Add the tap
brew tap scttfrdmn/tap

# Install synodeploy
brew install synodeploy

# Verify installation
synodeploy version
```

### Updating via Homebrew

```bash
brew upgrade synodeploy
```

## Direct Download

Download the latest binary for your platform from the [releases page](https://github.com/scttfrdmn/synodeploy/releases):

### macOS

```bash
# Intel Mac
curl -L https://github.com/scttfrdmn/synodeploy/releases/latest/download/synodeploy-darwin-amd64 -o synodeploy

# Apple Silicon Mac
curl -L https://github.com/scttfrdmn/synodeploy/releases/latest/download/synodeploy-darwin-arm64 -o synodeploy

# Make executable and install
chmod +x synodeploy
sudo mv synodeploy /usr/local/bin/
```

### Linux

```bash
# Intel/AMD 64-bit
curl -L https://github.com/scttfrdmn/synodeploy/releases/latest/download/synodeploy-linux-amd64 -o synodeploy

# ARM 64-bit
curl -L https://github.com/scttfrdmn/synodeploy/releases/latest/download/synodeploy-linux-arm64 -o synodeploy

# Make executable and install
chmod +x synodeploy
sudo mv synodeploy /usr/local/bin/
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
git clone https://github.com/scttfrdmn/synodeploy.git
cd synodeploy

# Build using Make
make build

# Or build using Go directly
go build -o bin/synodeploy main.go

# Install to system
sudo cp bin/synodeploy /usr/local/bin/
```

### Development Build

For development purposes, you can create a symlink:

```bash
make dev-install
```

This creates a symlink in `/usr/local/bin/synodeploy` pointing to your development build.

## Verification

After installation, verify SynoDeploy is working correctly:

```bash
# Check version
synodeploy version

# View help
synodeploy --help

# Check if all dependencies are available
synodeploy init --help
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
ls -la /usr/local/bin/synodeploy

# Fix permissions if needed
sudo chmod +x /usr/local/bin/synodeploy
```

### Homebrew Issues

If Homebrew installation fails:

```bash
# Update Homebrew
brew update

# Check for issues
brew doctor

# Try installing again
brew install scttfrdmn/tap/synodeploy
```

## Uninstallation

### Homebrew

```bash
brew uninstall synodeploy
brew untap scttfrdmn/tap
```

### Manual Installation

```bash
# Remove binary
sudo rm /usr/local/bin/synodeploy

# Remove configuration (optional)
rm -rf ~/.synodeploy
```

## Next Steps

After installation, proceed to the [Usage Guide](usage.md) to set up your first connection to your Synology NAS.