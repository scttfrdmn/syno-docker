# syno-docker: Simple Container Deployment for Synology DSM 7.2+

## Project Overview

A **single Go binary** that makes deploying containers to Synology DSM 7.2+ as simple as deploying to any Docker host. Distributed via Homebrew for zero-friction installation.

```bash
# Install
brew install syno-docker

# One-time setup 
syno-docker init 192.168.1.100

# Deploy anything
syno-docker run nginx:latest -p 8080:80 -v /volume1/web:/usr/share/nginx/html
syno-docker deploy docker-compose.yml
```

## Why Go is Perfect for This

### Technical Advantages
- **Single static binary** - perfect for Homebrew distribution
- **Excellent Docker SDK** - official `github.com/docker/docker` client
- **Cross-platform** - works on macOS, Linux, Windows out of the box
- **SSH libraries** - robust SSH client support with `golang.org/x/crypto/ssh`
- **YAML/JSON handling** - standard library support for configs
- **Goroutines** - concurrent health checks and deployments

### Distribution Benefits
```bash
# Users get a single command installation
brew install syno-docker

# Behind the scenes: static binary, no dependencies
# Works immediately without Python, Node, or other runtimes
```

## Simplified Architecture (DSM 7.2+ Only)

### Known Constants (No Detection Needed)
- **Service name:** `pkg-ContainerManager-dockerd`
- **Config path:** `/var/packages/ContainerManager/etc/dockerd.json`
- **Docker binary:** `/usr/local/bin/docker` (PATH issue is known)
- **Socket:** `/var/run/docker.sock`
- **Service restart:** `sudo systemctl restart pkg-ContainerManager-dockerd`

### Core Structure
```go
// main.go - CLI interface using cobra
// pkg/synology/ - DSM 7.2+ specific logic
// pkg/deploy/ - Docker deployment logic  
// pkg/config/ - Configuration management
```

## Minimal Viable Product

### 1. Core Commands
```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "syno-docker",
    Short: "Deploy containers to Synology DSM 7.2+",
}

// cmd/init.go
var initCmd = &cobra.Command{
    Use:   "init [host]",
    Short: "Setup connection to Synology NAS",
    Run:   initCommand,
}

// cmd/run.go  
var runCmd = &cobra.Command{
    Use:   "run [image]",
    Short: "Deploy a single container",
    Run:   runCommand,
}

// cmd/deploy.go
var deployCmd = &cobra.Command{
    Use:   "deploy [compose-file]", 
    Short: "Deploy from docker-compose.yml",
    Run:   deployCommand,
}
```

### 2. Connection Manager
```go
// pkg/synology/connection.go
type Connection struct {
    Host     string
    User     string
    SSHKey   string
    client   *ssh.Client
    docker   *client.Client
}

func (c *Connection) Connect() error {
    // 1. SSH connection with key auth
    // 2. Setup Docker client over SSH
    // 3. Handle PATH issue: use full path /usr/local/bin/docker
    return nil
}

func (c *Connection) ExecuteDocker(args []string) error {
    // Always use: /usr/local/bin/docker [args]
    cmd := "/usr/local/bin/docker " + strings.Join(args, " ")
    return c.runSSHCommand(cmd)
}
```

### 3. Simple Configuration
```go
// pkg/config/config.go
type Config struct {
    Host    string `yaml:"host"`
    User    string `yaml:"user"`
    SSHKey  string `yaml:"ssh_key_path"`
    Default struct {
        VolumePath string `yaml:"volume_path"` // /volume1/docker
    } `yaml:"default"`
}

// Stored in ~/.syno-docker/config.yaml
```

### 4. Container Deployment
```go
// pkg/deploy/container.go
type ContainerDeploy struct {
    Image   string
    Ports   []string  // "8080:80"
    Volumes []string  // "/volume1/data:/app/data"
    Env     []string  // "KEY=value"
    Name    string
}

func (cd *ContainerDeploy) Deploy(conn *synology.Connection) error {
    // 1. docker pull [image]
    // 2. docker run with all flags
    // 3. Return container ID
}
```

## Project Structure (Minimal)

```
syno-docker/
├── main.go                   # Entry point
├── cmd/
│   ├── root.go              # Root command setup
│   ├── init.go              # syno-docker init
│   ├── run.go               # syno-docker run  
│   └── deploy.go            # syno-docker deploy
├── pkg/
│   ├── synology/
│   │   ├── connection.go    # SSH + Docker client
│   │   └── paths.go         # DSM 7.2+ specific paths
│   ├── deploy/
│   │   ├── container.go     # Single container deployment
│   │   └── compose.go       # docker-compose deployment
│   └── config/
│       └── config.go        # Configuration management
├── go.mod
├── go.sum
├── Makefile                 # Build automation
└── README.md
```

## Key Dependencies
```go
// go.mod
module github.com/username/syno-docker

require (
    github.com/spf13/cobra v1.8.0          // CLI framework
    github.com/docker/docker v24.0.0       // Docker client
    golang.org/x/crypto v0.17.0            // SSH client
    gopkg.in/yaml.v3 v3.0.1               // YAML parsing
)
```

## Homebrew Distribution

### 1. Goreleaser Configuration
```yaml
# .goreleaser.yaml
builds:
  - env: [CGO_ENABLED=0]
    goos: [linux, darwin]
    goarch: [amd64, arm64]

brews:
  - name: syno-docker
    homepage: https://github.com/username/syno-docker
    description: "Deploy containers to Synology DSM 7.2+"
    repository:
      owner: username
      name: homebrew-syno-docker
```

### 2. Release Process
```bash
# Automated via GitHub Actions
git tag v1.0.0
git push --tags
# Goreleaser builds binaries and updates Homebrew formula
```

### 3. User Experience
```bash
# Installation
brew install username/syno-docker/syno-docker

# Usage
syno-docker init 192.168.1.100
syno-docker run portainer/portainer-ce:latest -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock
```

## Implementation Strategy

### Phase 1: Core Functionality (2-4 weeks)
- [ ] SSH connection with key authentication
- [ ] Docker client over SSH with PATH fixes
- [ ] Basic `syno-docker run` command
- [ ] Configuration management
- [ ] Homebrew formula

### Phase 2: Enhanced Features (2-3 weeks)  
- [ ] docker-compose support
- [ ] Container management (stop, start, remove)
- [ ] Volume path validation and creation
- [ ] Better error handling and user feedback

### Phase 3: Polish (1-2 weeks)
- [ ] Comprehensive documentation
- [ ] Integration tests with real Synology hardware
- [ ] Performance optimizations
- [ ] Community feedback integration

## Success Metrics

- **Installation friction:** `brew install syno-docker` → working deployment in under 5 minutes
- **Binary size:** < 20MB static binary
- **Startup time:** < 100ms for most commands
- **Memory usage:** < 10MB during operation
- **Error clarity:** Every error includes next steps

## Why This Approach Works

### 1. Leverages Go's Strengths
- Static binaries eliminate dependency hell
- Cross-compilation handles multiple platforms
- Excellent standard library for networking/SSH
- Strong Docker ecosystem integration

### 2. Homebrew is Perfect Distribution
- Trusted by developers (brew is everywhere)
- Automatic updates via `brew upgrade`
- No manual PATH configuration needed
- Version management built-in

### 3. DSM 7.2+ Focus Simplifies Everything
- No version detection needed
- Known paths and service names
- Consistent Container Manager behavior
- Modern systemd service management

### 4. Follows Unix Philosophy
- Do one thing well (deploy containers)
- Compose with existing tools (docker-compose, ssh)
- Simple, predictable interface
- No unnecessary features

This approach gives you a professional, maintainable tool that users will actually want to use. The Go ecosystem provides excellent libraries for all the required functionality, and Homebrew distribution makes adoption frictionless.

Ready to start building?