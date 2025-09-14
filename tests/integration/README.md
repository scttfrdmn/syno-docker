# SynoDeploy Integration Test Suite

This comprehensive integration test suite validates SynoDeploy functionality against a real Synology NAS running DSM 7.2+ with Container Manager.

## 🎯 **Test Coverage**

### Core Functionality Tests
- **Basic Deployment**: Single container deployment with various configurations
- **Compose Deployment**: Multi-container applications using docker-compose.yml
- **Lifecycle Management**: Container start, stop, restart, and removal operations
- **Volume Mapping**: Host path mounting and named volume management
- **Network Connectivity**: Inter-container communication and external access
- **Error Handling**: Invalid configurations and failure scenarios

### Advanced Scenarios
- **Performance Tests**: Resource usage and deployment speed benchmarks
- **Stress Tests**: High load and concurrent deployment scenarios
- **Edge Cases**: Unusual configurations and boundary conditions
- **Security Tests**: Permission validation and access control

## 🛠️ **Prerequisites**

### Synology NAS Requirements
- **DSM Version**: 7.2 or later
- **Container Manager**: Installed and running
- **SSH Service**: Enabled (Control Panel → Terminal & SNMP)
- **User Account**: Admin privileges with docker group membership
- **Storage**: At least 2GB free space on `/volume1/`
- **Network**: Accessible from test machine

### Test Environment Setup
```bash
# Required tools
- Go 1.21+
- SSH client with key-based authentication
- Network connectivity to NAS
- Docker client (for comparison/validation)

# Optional tools
- curl/wget (for HTTP endpoint testing)
- PostgreSQL client (for database connectivity tests)
```

## 📋 **Setup Instructions**

### 1. Prepare Your NAS

**Enable SSH Access:**
```bash
# On your NAS via SSH or DSM Terminal
sudo synogroup --add docker your-username
sudo systemctl enable sshd
sudo systemctl start sshd
```

**Create Test Directory:**
```bash
sudo mkdir -p /volume1/synodeploy-test
sudo chown your-username:users /volume1/synodeploy-test
sudo chmod 755 /volume1/synodeploy-test
```

### 2. Configure SSH Keys
```bash
# Generate SSH key if needed
ssh-keygen -t rsa -b 4096 -C "synodeploy-test"

# Copy public key to NAS
ssh-copy-id your-username@your-nas-ip

# Test SSH connection
ssh your-username@your-nas-ip "docker version"
```

### 3. Configure Test Environment
```bash
# Copy and customize test configuration
cp tests/integration/config/test_config.yaml.example test_config.yaml
# Edit test_config.yaml with your NAS details

# Set environment variables (alternative to config file)
export NAS_HOST="192.168.1.100"
export NAS_USER="admin"
export NAS_SSH_KEY="~/.ssh/id_rsa"
```

## 🚀 **Running Tests**

### Basic Test Run
```bash
# Run all integration tests
go test -v -integration \
    -nas-host=192.168.1.100 \
    -nas-user=admin \
    -nas-key=~/.ssh/id_rsa \
    ./tests/integration/...

# Run specific test suite
go test -v -integration \
    -nas-host=192.168.1.100 \
    -run TestBasicDeployment \
    ./tests/integration/
```

### Advanced Test Options
```bash
# Run with custom timeout
go test -v -integration -timeout=10m \
    -nas-host=192.168.1.100 \
    ./tests/integration/...

# Run performance benchmarks
go test -v -integration -bench=. \
    -nas-host=192.168.1.100 \
    ./tests/integration/...

# Skip cleanup for debugging
go test -v -integration \
    -nas-host=192.168.1.100 \
    -cleanup=false \
    ./tests/integration/...

# Parallel execution (experimental)
go test -v -integration -parallel=4 \
    -nas-host=192.168.1.100 \
    ./tests/integration/...
```

### Makefile Targets
```bash
# Run integration tests (requires NAS_HOST environment variable)
make integration-test

# Run with full cleanup and reporting
make integration-test-full

# Generate integration test report
make integration-test-report
```

## 📊 **Test Scenarios**

### 1. Basic Deployment Tests
- Single container deployment (nginx, postgres, redis)
- Port mapping validation
- Volume mounting verification
- Environment variable passing
- Restart policy enforcement

### 2. Compose Deployment Tests
- Multi-service application deployment
- Service dependency management
- Named volume creation and mounting
- Custom network configuration
- Environment variable substitution

### 3. Lifecycle Management Tests
- Container start/stop operations
- Graceful shutdown handling
- Container removal with cleanup
- Force removal scenarios
- Container recreation and updates

### 4. Volume and Storage Tests
- Host path mounting (`/volume1/`, `/volume2/`)
- Named volume creation and management
- Volume permission validation
- Data persistence verification
- Storage cleanup operations

### 5. Network Connectivity Tests
- Container-to-container communication
- External network access validation
- Port binding and forwarding
- Custom bridge network creation
- DNS resolution within containers

### 6. Error Handling Tests
- Invalid image names
- Non-existent volume paths
- Port conflicts
- Resource exhaustion scenarios
- Network connectivity failures

## 🔧 **Test Configuration**

### Environment Variables
```bash
# Required
export NAS_HOST="192.168.1.100"          # NAS IP address
export NAS_USER="admin"                   # SSH username
export NAS_SSH_KEY="~/.ssh/id_rsa"        # SSH private key path

# Optional
export NAS_PORT="22"                      # SSH port (default: 22)
export TEST_VOLUME_PATH="/volume1/test"   # Test volume directory
export TEST_TIMEOUT="5m"                  # Test timeout
export TEST_PARALLEL="false"             # Parallel execution
export TEST_CLEANUP="true"               # Cleanup after tests
```

### Configuration File (test_config.yaml)
```yaml
nas:
  host: "192.168.1.100"
  user: "admin"
  ssh_key_path: "~/.ssh/id_rsa"

test:
  volume_path: "/volume1/synodeploy-test"
  cleanup: true
  timeout: "5m"

scenarios:
  basic_deployment: true
  compose_deployment: true
  performance: false  # Enable for performance testing
```

## 📈 **Test Output and Reporting**

### Console Output
```bash
=== RUN   TestIntegration
=== RUN   TestIntegration/BasicDeployment
    Deploying container: test-nginx-a8b9c2d3
    Container deployed successfully: c4f1e2d9a8b7
    Testing HTTP connectivity to http://192.168.1.100:8080
    ✓ HTTP endpoint accessible and returning expected content
=== RUN   TestIntegration/ComposeDeployment
    Deploying compose project: test-stack-x7y9z1
    Services deployed: web, api, db, cache
    ✓ All services started successfully
    ✓ Inter-service connectivity verified
--- PASS: TestIntegration (45.67s)
```

### Test Reports
```bash
# Generate detailed test report
go test -v -integration -json ./tests/integration/... > test_report.json

# Generate coverage report
go test -integration -coverprofile=integration_coverage.out ./tests/integration/...
go tool cover -html=integration_coverage.out -o integration_coverage.html
```

## 🧹 **Cleanup and Troubleshooting**

### Manual Cleanup
```bash
# Remove all test containers
ssh your-nas-ip "docker ps -a | grep 'test-' | awk '{print \$1}' | xargs docker rm -f"

# Remove test volumes
ssh your-nas-ip "docker volume ls | grep 'test-' | awk '{print \$2}' | xargs docker volume rm"

# Clean test directory
ssh your-nas-ip "rm -rf /volume1/synodeploy-test/*"
```

### Common Issues

**SSH Connection Failed:**
```bash
# Verify SSH access
ssh -v your-username@your-nas-ip

# Check SSH key permissions
chmod 600 ~/.ssh/id_rsa
chmod 644 ~/.ssh/id_rsa.pub
```

**Docker Permission Denied:**
```bash
# Verify docker group membership
ssh your-nas-ip "groups"
ssh your-nas-ip "docker ps"
```

**Volume Mount Errors:**
```bash
# Check volume path permissions
ssh your-nas-ip "ls -la /volume1/synodeploy-test"
ssh your-nas-ip "mkdir -p /volume1/synodeploy-test && chmod 755 /volume1/synodeploy-test"
```

## 🎯 **Best Practices**

### Test Development
1. **Idempotent Tests**: Each test should clean up after itself
2. **Unique Names**: Use random suffixes for test resources
3. **Timeout Handling**: Set appropriate timeouts for operations
4. **Error Validation**: Test both success and failure scenarios
5. **Resource Limits**: Be mindful of NAS resource constraints

### CI/CD Integration
1. **Environment Variables**: Use env vars for CI configuration
2. **Test Reports**: Generate JSON/XML reports for CI systems
3. **Parallel Execution**: Use `-parallel` flag judiciously
4. **Cleanup Verification**: Ensure all resources are cleaned up
5. **Test Isolation**: Each test run should be independent

This integration test suite provides comprehensive validation of SynoDeploy functionality against real Synology hardware, ensuring production readiness and reliability.