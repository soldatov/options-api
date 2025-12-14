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
  -e CONFIG_FILE=/root/options.json \
  options-api
```

## Volume Mounting

For development or persistent configuration:

```bash
docker run -d --name options-api \
  -p 8080:8080 \
  -v $(pwd)/templates:/root/templates \
  -v $(pwd)/options.json:/root/options.json \
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

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  options-api:
    image: your-dockerhub-username/options-api:latest
    ports:
      - "8080:8080"
    restart: unless-stopped
    volumes:
      - ./options.json:/root/options.json
    environment:
      - CONFIG_FILE=/root/options.json
```

Deploy with:

```bash
docker-compose up -d
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