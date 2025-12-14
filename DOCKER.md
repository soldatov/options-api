# Docker Deployment Guide

This guide explains how to build and deploy the Options API application using Docker.

## Prerequisites

1. **Docker Desktop** installed and running
2. **Docker Hub** account
3. **Make** command available (usually pre-installed on Linux/macOS)

## Quick Start

### 1. Configure Docker Hub Username

Edit the `Makefile` and update the `DOCKERHUB_USERNAME` variable:

```makefile
DOCKERHUB_USERNAME = your-actual-dockerhub-username
```

### 2. Build and Push to Docker Hub

```bash
# Login to Docker Hub (one-time setup)
make login

# Build and push the image
make push
```

The image will be tagged as: `your-dockerhub-username/options-api:latest`

## Available Commands

### Development Commands

```bash
# Show all available commands
make help

# Build the Docker image locally
make build

# Run the container locally
make run

# Run with volume mount for template development
make dev-run

# View container logs
make logs

# Get interactive shell in running container
make shell

# Stop and remove container
make stop
```

### Deployment Commands

```bash
# Login to Docker Hub
make login

# Build and push to Docker Hub
make push

# Full deployment cycle (clean, login, push)
make deploy

# Clean up local images and containers
make clean
```

## Docker Configuration

### Multi-stage Build

The `Dockerfile` uses a multi-stage build process:

1. **Builder Stage**: Compiles the Go application
2. **Runtime Stage**: Creates a minimal Alpine Linux image with just the binary

### Image Features

- **Size**: Optimized using Alpine Linux (~15MB)
- **Security**: No shell access in final image
- **Port**: Exposes port 8080
- **Templates**: Includes HTML templates in the image

## Running the Container

### Using Docker Command

```bash
# Build and run
docker build -t options-api .
docker run -d --name options-api -p 8080:8080 options-api
```

### Using Make Commands

```bash
# Build and run
make build
make run
```

### Accessing the Application

After starting the container, access the application at:
- **Local**: http://localhost:8080
- **Remote**: http://your-server-ip:8080

## Environment Variables

The application can be configured using environment variables:

```bash
docker run -d --name options-api \
  -p 8080:8080 \
  -e CONFIG_FILE=/app/options.json \
  options-api
```

## Volume Mounting

For development or persistent configuration:

```bash
docker run -d --name options-api \
  -p 8080:8080 \
  -v $(pwd)/templates:/app/templates \
  -v $(pwd)/options.json:/app/options.json \
  options-api
```

## Docker Hub Integration

### Automated Tagging

The Makefile supports different tags:

```bash
# Build with specific tag
make build TAG=v1.0.0

# Push with specific tag
make push TAG=v1.0.0
```

### Repository Structure

- **Image Name**: `options-api`
- **Repository**: `your-dockerhub-username/options-api`
- **Default Tag**: `latest`

## Troubleshooting

### Docker Daemon Not Running

```bash
# On macOS
open /Applications/Docker.app

# On Linux
sudo systemctl start docker
sudo systemctl enable docker
```

### Permission Issues

```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
# Logout and login again
```

### Port Already in Use

```bash
# Find what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port
docker run -d --name options-api -p 8081:8080 options-api
```

## Production Deployment

### Using Docker Compose

The project includes a pre-configured `docker-compose.yml` file for easy deployment:

```yaml
version: '3.8'

services:
  options-api:
    build: .
    container_name: options-api-compose
    ports:
      - "8080:8080"
    volumes:
      - ./options.json:/app/options.json
      - ./templates:/app/templates
    environment:
      - CONFIG_FILE=/app/options.json
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - options-api-network

networks:
  options-api-network:
    driver: bridge
```

#### Quick Start with Docker Compose

1. **Ensure Docker daemon is running**
2. **Navigate to project directory** (where docker-compose.yml is located)
3. **Start the service**:

```bash
docker-compose up --build -d
```

4. **Access the application** at http://localhost:8080
5. **Stop the service**:

```bash
docker-compose down
```

#### Docker Compose Features

- **Volume Mounting**: The `options.json` file is mounted from the same directory as docker-compose.yml
- **Template Hot-reloading**: Templates directory is mounted for development
- **Health Checks**: Automatic health monitoring with curl checks
- **Auto-restart**: Container restarts automatically if it fails
- **Network Isolation**: Uses dedicated Docker network
- **Environment Configuration**: Config file path controlled by `CONFIG_FILE` environment variable

#### Custom Configuration

Place your `options.json` file in the same directory as `docker-compose.yml`:

```bash
# Create custom configuration
cat > options.json << 'EOF'
{
  "fields": [
    {"name": "customField1", "value": "defaultValue1"},
    {"name": "customField2", "value": 42}
  ]
}
EOF

# Start with custom config
docker-compose up --build -d
```

#### Development Mode

For development with live template reloading:

```bash
# Run in development mode
docker-compose up --build

# View logs
docker-compose logs -f

# Get shell access
docker-compose exec options-api /bin/sh
```

#### Production Deployment

```bash
# Deploy to production
docker-compose -f docker-compose.yml up -d

# Update the service
docker-compose pull
docker-compose up -d

# Scale the service
docker-compose up -d --scale options-api=3
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: options-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: options-api
  template:
    metadata:
      labels:
        app: options-api
    spec:
      containers:
      - name: options-api
        image: your-dockerhub-username/options-api:latest
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: options-api-service
spec:
  selector:
    app: options-api
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```