# Variables
IMAGE_NAME = options-api
DOCKERHUB_USERNAME = soldatovdev
TAG = latest
PLATFORMS = linux/amd64,linux/arm64

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build        - Build the Docker image for current platform"
	@echo "  buildx       - Build multi-architecture Docker image"
	@echo "  push         - Build and push multi-arch image to Docker Hub"
	@echo "  push-single  - Build and push single-arch image to Docker Hub"
	@echo "  run          - Run the Docker container locally"
	@echo "  stop         - Stop the running container"
	@echo "  clean        - Remove Docker image and container"
	@echo "  login        - Login to Docker Hub"
	@echo ""
	@echo "Before using 'push' command, make sure to:"
	@echo "1. Update DOCKERHUB_USERNAME in this Makefile"
	@echo "2. Run 'make login' to authenticate with Docker Hub"

# Build Docker image for current platform
.PHONY: build
build:
	@echo "Building Docker image for current platform..."
	docker build -t $(IMAGE_NAME):$(TAG) .

# Build multi-architecture Docker image
.PHONY: buildx
buildx:
	@echo "Setting up Docker buildx..."
	docker buildx create --name multiarch --driver docker-container --use || docker buildx use multiarch
	docker buildx inspect --bootstrap
	@echo "Building multi-architecture Docker image..."
	docker buildx build --platform $(PLATFORMS) -t $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) .

# Push to Docker Hub (multi-architecture)
.PHONY: push
push:
	@echo "Setting up Docker buildx..."
	docker buildx create --name multiarch --driver docker-container --use || docker buildx use multiarch
	docker buildx inspect --bootstrap
	@echo "Building and pushing multi-architecture image to Docker Hub..."
	docker buildx build --platform $(PLATFORMS) -t $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) --push .
	@echo "Multi-arch image pushed successfully: $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)"

# Push to Docker Hub (single architecture)
.PHONY: push-single
push-single: build
	@echo "Pushing single-architecture image to Docker Hub..."
	docker tag $(IMAGE_NAME):$(TAG) $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)
	docker push $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)
	@echo "Single-arch image pushed successfully: $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)"

# Run container locally
.PHONY: run
run:
	@echo "Running container..."
	docker run -d --name $(IMAGE_NAME) -p 8080:8080 $(IMAGE_NAME):$(TAG)
	@echo "Container is running at http://localhost:8080"

# Stop container
.PHONY: stop
stop:
	@echo "Stopping container..."
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true

# Clean up
.PHONY: clean
clean: stop
	@echo "Removing Docker images..."
	docker rmi $(IMAGE_NAME):$(TAG) || true
	docker rmi $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) || true
	docker buildx rm multiarch || true

# Login to Docker Hub
.PHONY: login
login:
	@echo "Logging in to Docker Hub..."
	docker login

# Development run with volume mount
.PHONY: dev-run
dev-run:
	@echo "Running container with volume mount for development..."
	docker run -d --name $(IMAGE_NAME)-dev -p 8080:8080 -v $(PWD)/templates:/app/templates $(IMAGE_NAME):$(TAG)
	@echo "Dev container is running at http://localhost:8080"

# View logs
.PHONY: logs
logs:
	docker logs -f $(IMAGE_NAME)

# Interactive shell in container
.PHONY: shell
shell:
	docker exec -it $(IMAGE_NAME) /bin/sh

# Test multi-architecture image
.PHONY: test-multi
test-multi:
	@echo "Testing multi-architecture image..."
	@echo "Testing AMD64 platform..."
	docker run --rm --platform linux/amd64 $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) /bin/sh -c "uname -m && echo 'AMD64 test passed'"
	@echo "Testing ARM64 platform..."
	docker run --rm --platform linux/arm64 $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) /bin/sh -c "uname -m && echo 'ARM64 test passed'"

# Full build and deploy cycle
.PHONY: deploy
deploy: clean login push
	@echo "Deployment complete!"

# Full build and deploy cycle with testing
.PHONY: deploy-test
deploy-test: clean login push test-multi
	@echo "Deployment with testing complete!"