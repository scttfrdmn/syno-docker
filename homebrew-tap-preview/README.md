# Scott Friedman's Homebrew Tap

Official Homebrew tap for tools developed by Scott Friedman.

## Available Formulae

### [syno-docker](https://github.com/scttfrdmn/syno-docker)
Deploy containers to Synology DSM 7.2+ with ease.

```bash
brew install scttfrdmn/tap/syno-docker
```

**Features:**
- 🚀 One-command container deployment
- 🔐 SSH key & ssh-agent support
- 👤 Administrator user compatibility
- 📦 docker-compose support
- 🎯 DSM 7.2+ optimized
- 🧪 Integration tested on real hardware

## Installation

```bash
# Add the tap
brew tap scttfrdmn/tap

# Install any formula
brew install <formula-name>

# Example: Install syno-docker
brew install syno-docker
```

## Usage

### syno-docker
```bash
# Setup connection to your Synology NAS
syno-docker init your-nas.local --user your-username

# Deploy a container
syno-docker run nginx:latest --port 8080:80

# Deploy from docker-compose
syno-docker deploy docker-compose.yml

# List containers
syno-docker ps

# Remove container
syno-docker rm container-name
```

## Requirements

### For syno-docker
- **Synology NAS**: DSM 7.2+ with Container Manager
- **SSH Access**: Enabled with key-based authentication or ssh-agent
- **Network Access**: Local network connectivity to your NAS

## Support

- 📖 [syno-docker Documentation](https://github.com/scttfrdmn/syno-docker/tree/main/docs)
- 🐛 [Issue Tracker](https://github.com/scttfrdmn/syno-docker/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/syno-docker/discussions)

## Contributing

If you'd like to contribute a formula to this tap:

1. Fork this repository
2. Create a new formula in the `Formula/` directory
3. Test the formula locally
4. Submit a pull request

### Formula Guidelines
- Follow [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- Include a proper `test do` block
- Use semantic versioning
- Include comprehensive documentation

## About

This tap is maintained by [Scott Friedman](https://github.com/scttfrdmn) and contains formulae for various development tools and utilities.

**License:** Individual formulae may have different licenses. Check each formula for specific license information.

**Automated Updates:** This tap is automatically updated via GitHub Actions and GoReleaser when new releases are published.