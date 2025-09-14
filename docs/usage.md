# Usage Guide

This guide covers how to use syno-docker to deploy and manage containers on your Synology NAS.

## Initial Setup

### 1. Prepare Your Synology NAS

Before using syno-docker, ensure your NAS is properly configured:

1. **Install Container Manager**
   - Open Package Center on your NAS
   - Search for and install "Container Manager"
   - Start the Container Manager service

2. **Enable SSH Access**
   - Go to Control Panel → Terminal & SNMP
   - Enable SSH service
   - Note the port number (default: 22)

3. **Setup SSH Key Authentication** (Recommended)
   ```bash
   # Generate SSH key if you don't have one
   ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

   # Copy your public key to the NAS
   ssh-copy-id admin@your-nas-ip
   ```

### 2. Initialize syno-docker

```bash
# Basic initialization (uses ssh-agent if available)
syno-docker init 192.168.1.100

# With custom admin username
syno-docker init your-nas.local --user your-username

# With all custom settings
syno-docker init your-nas.local \
  --user your-username \
  --port 22 \
  --key ~/.ssh/id_rsa \
  --volume-path /volume1/docker

# Examples for different setups
syno-docker init chubchub.local --user scttfrdmn    # Custom admin user
syno-docker init 192.168.1.100 --user admin        # Standard admin user
```

This command will:
- Test SSH connection to your NAS
- Verify Docker/Container Manager is accessible
- Save configuration to `~/.syno-docker/config.yaml`

## Container Deployment

### Single Container Deployment

Deploy a simple container:

```bash
syno-docker run nginx:latest
```

Deploy with custom configuration:

```bash
syno-docker run nginx:latest \
  --name web-server \
  --port 8080:80 \
  --port 8443:443 \
  --volume /volume1/web:/usr/share/nginx/html \
  --volume /volume1/config/nginx:/etc/nginx/conf.d \
  --env NGINX_HOST=example.com \
  --restart unless-stopped
```

### Advanced Container Options

```bash
syno-docker run postgres:13 \
  --name database \
  --port 5432:5432 \
  --volume db-data:/var/lib/postgresql/data \
  --env POSTGRES_DB=myapp \
  --env POSTGRES_USER=appuser \
  --env POSTGRES_PASSWORD=secretpassword \
  --restart unless-stopped \
  --user 1000:1000 \
  --network bridge \
  --workdir /var/lib/postgresql
```

### Docker Compose Deployment

Deploy a multi-container application:

```bash
# Deploy from current directory
syno-docker deploy docker-compose.yml

# With custom project name
syno-docker deploy docker-compose.yml --project my-awesome-app

# With environment file
syno-docker deploy docker-compose.yml --env-file .env.production
```

Example `docker-compose.yml`:

```yaml
version: '3.8'
services:
  web:
    image: nginx:latest
    ports:
      - "8080:80"
    volumes:
      - ./html:/usr/share/nginx/html
      - ./nginx.conf:/etc/nginx/nginx.conf
    environment:
      - NGINX_HOST=${HOST:-localhost}
    restart: unless-stopped
    depends_on:
      - api

  api:
    image: node:16-alpine
    ports:
      - "3000:3000"
    volumes:
      - ./app:/usr/src/app
      - /volume1/uploads:/usr/src/app/uploads
    environment:
      - NODE_ENV=production
      - DATABASE_URL=${DATABASE_URL}
    working_dir: /usr/src/app
    command: ["npm", "start"]
    restart: unless-stopped

  db:
    image: postgres:13
    volumes:
      - db_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    restart: unless-stopped

volumes:
  db_data:
```

## Container Management

### List Containers

```bash
# Show running containers
syno-docker ps

# Show all containers (including stopped)
syno-docker ps --all
```

Output example:
```
CONTAINER ID  NAME        IMAGE         STATUS                 PORTS
a1b2c3d4e5f6  web-server  nginx:latest  Up 2 hours            0.0.0.0:8080->80/tcp
f6e5d4c3b2a1  database    postgres:13   Up 2 hours            0.0.0.0:5432->5432/tcp
```

### Remove Containers

```bash
# Remove stopped container
syno-docker rm web-server

# Force remove running container
syno-docker rm web-server --force

# Remove multiple containers
syno-docker rm web-server database api-server --force
```

## Volume Path Handling

syno-docker provides convenient volume path handling:

### Absolute Paths
```bash
# Full path (works as expected)
syno-docker run nginx -v /volume1/web:/usr/share/nginx/html
```

### Relative Paths
```bash
# Relative to default volume path
syno-docker run nginx -v web:/usr/share/nginx/html
# Expands to: /volume1/docker/web:/usr/share/nginx/html

# Current directory relative
syno-docker run nginx -v ./web:/usr/share/nginx/html
# Expands to: /volume1/docker/web:/usr/share/nginx/html
```

### Named Volumes
```bash
# Docker named volumes (managed by Docker)
syno-docker run postgres -v postgres-data:/var/lib/postgresql/data
```

## Environment Variables

### Single Container
```bash
syno-docker run app:latest \
  --env NODE_ENV=production \
  --env DATABASE_URL=postgres://user:pass@db:5432/myapp \
  --env API_KEY=your-secret-key
```

### Docker Compose with .env File

Create `.env` file:
```bash
NODE_ENV=production
DATABASE_URL=postgres://user:pass@localhost:5432/myapp
API_KEY=your-secret-key
POSTGRES_DB=myapp
POSTGRES_USER=appuser
POSTGRES_PASSWORD=secretpassword
```

Deploy with environment:
```bash
syno-docker deploy docker-compose.yml --env-file .env
```

Use in `docker-compose.yml`:
```yaml
services:
  app:
    image: myapp:latest
    environment:
      - NODE_ENV=${NODE_ENV}
      - DATABASE_URL=${DATABASE_URL}
      - API_KEY=${API_KEY}
```

## Network Configuration

### Default Bridge Network
```bash
# Uses default bridge network
syno-docker run nginx --port 8080:80
```

### Custom Network Mode
```bash
# Host networking (use NAS IP directly)
syno-docker run nginx --network host

# No networking
syno-docker run batch-job --network none
```

## Port Mapping

### Basic Port Mapping
```bash
# Map port 8080 on NAS to port 80 in container
syno-docker run nginx --port 8080:80
```

### Multiple Ports
```bash
syno-docker run nginx \
  --port 8080:80 \
  --port 8443:443 \
  --port 9090:9000
```

### Dynamic Ports
```bash
# Let Docker assign random ports
syno-docker run nginx --port :80 --port :443
```

## Restart Policies

```bash
# Never restart
syno-docker run nginx --restart no

# Always restart
syno-docker run nginx --restart always

# Restart unless manually stopped (recommended)
syno-docker run nginx --restart unless-stopped

# Restart on failure only
syno-docker run nginx --restart on-failure
```

## User and Permissions

### Running as Specific User
```bash
# Run as user ID 1000, group ID 1000
syno-docker run nginx --user 1000:1000

# Run as current user (useful for file permissions)
syno-docker run nginx --user $(id -u):$(id -g)
```

### File Permissions
When mounting volumes, ensure proper permissions:

```bash
# On your NAS, create directory with proper permissions
ssh admin@your-nas-ip 'sudo mkdir -p /volume1/app-data && sudo chown 1000:1000 /volume1/app-data'

# Deploy container with matching user
syno-docker run myapp --user 1000:1000 --volume /volume1/app-data:/data
```

## Common Use Cases

### Web Server
```bash
syno-docker run nginx:latest \
  --name web-server \
  --port 80:80 \
  --port 443:443 \
  --volume /volume1/web:/usr/share/nginx/html \
  --volume /volume1/certs:/etc/nginx/certs \
  --restart unless-stopped
```

### Database Server
```bash
syno-docker run postgres:13 \
  --name database \
  --port 5432:5432 \
  --volume /volume1/postgres:/var/lib/postgresql/data \
  --env POSTGRES_DB=myapp \
  --env POSTGRES_USER=appuser \
  --env POSTGRES_PASSWORD=secretpassword \
  --restart unless-stopped
```

### Media Server
```bash
syno-docker run plex/plex-media-server:latest \
  --name plex \
  --port 32400:32400 \
  --volume /volume1/media:/data \
  --volume /volume1/plex-config:/config \
  --env PLEX_CLAIM=your-claim-token \
  --network host \
  --restart unless-stopped
```

### Development Environment
```bash
syno-docker run node:16-alpine \
  --name dev-env \
  --port 3000:3000 \
  --volume /volume1/projects/myapp:/workspace \
  --workdir /workspace \
  --command npm run dev \
  --env NODE_ENV=development
```

## Troubleshooting

### Connection Issues

**SSH Authentication:**
```bash
# Test SSH connection manually
ssh your-username@your-nas-ip

# If using ssh-agent, verify keys are loaded
ssh-add -l

# If using key files, test specific key
ssh -i ~/.ssh/id_rsa your-username@your-nas-ip
```

**Container Manager:**
```bash
# Verify Container Manager is running
ssh your-username@your-nas-ip 'systemctl status pkg-ContainerManager-dockerd'

# Test Docker access
ssh your-username@your-nas-ip '/usr/local/bin/docker version'
```

### Container Won't Start
```bash
# Check container logs via SSH
ssh admin@your-nas-ip 'docker logs container-name'

# Check container status
syno-docker ps --all
```

### Port Already in Use
```bash
# Find what's using the port
ssh admin@your-nas-ip 'netstat -tlnp | grep :8080'

# Use different port
syno-docker run nginx --port 8081:80
```

### Permission Issues
```bash
# Ensure user is in docker group
ssh admin@your-nas-ip 'groups'

# Add user to docker group if needed
ssh admin@your-nas-ip 'sudo synogroup --member docker $(whoami)'
```

### Volume Mount Issues
```bash
# Check if path exists
ssh admin@your-nas-ip 'ls -la /volume1/your-path'

# Create path if needed
ssh admin@your-nas-ip 'mkdir -p /volume1/your-path'

# Fix permissions
ssh admin@your-nas-ip 'sudo chown -R 1000:1000 /volume1/your-path'
```

## Configuration Management

### View Current Configuration
```bash
cat ~/.syno-docker/config.yaml
```

### Update Configuration
```bash
# Re-run init to update settings
syno-docker init 192.168.1.100 --user newuser

# Or edit the file directly
vim ~/.syno-docker/config.yaml
```

### Multiple NAS Devices
Currently, syno-docker supports one NAS configuration at a time. To manage multiple devices:

```bash
# Switch to different NAS
syno-docker init nas1.local
# ... deploy containers ...

# Switch to another NAS
syno-docker init nas2.local
# ... deploy containers ...
```

## Best Practices

1. **Use SSH Keys**: Always use SSH key authentication instead of passwords
2. **Volume Paths**: Use absolute paths starting with `/volume1/` for clarity
3. **Restart Policies**: Use `unless-stopped` for production containers
4. **Resource Limits**: Monitor resource usage on your NAS
5. **Backups**: Backup your configuration and important volumes
6. **Updates**: Regularly update container images and syno-docker itself
7. **Security**: Don't expose unnecessary ports to the internet
8. **Monitoring**: Monitor container logs and resource usage

## Next Steps

- Explore the [examples](examples/) directory for more complex deployments
- Check out the [FAQ](faq.md) for common questions
- Join the community discussions for tips and tricks