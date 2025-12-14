# Variables
IMAGE_NAME = options-api
DOCKERHUB_USERNAME = soldatovanatoly
TAG = latest

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build        - Build the Docker image"
	@echo "  push         - Build and push the Docker image to Docker Hub"
	@echo "  run          - Run the Docker container locally"
	@echo "  stop         - Stop the running container"
	@echo "  clean        - Remove Docker image and container"
	@echo "  login        - Login to Docker Hub"
	@echo ""
	@echo "Before using 'push' command, make sure to:"
	@echo "1. Update DOCKERHUB_USERNAME in this Makefile"
	@echo "2. Run 'make login' to authenticate with Docker Hub"

# Build Docker image
.PHONY: build
build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):$(TAG) .

# Push to Docker Hub
.PHONY: push
push: build
	@echo "Pushing image to Docker Hub..."
	docker tag $(IMAGE_NAME):$(TAG) $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)
	docker push $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)
	@echo "Image pushed successfully: $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG)"

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
	@echo "Removing Docker image..."
	docker rmi $(IMAGE_NAME):$(TAG) || true
	docker rmi $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(TAG) || true

# Login to Docker Hub
.PHONY: login
login:
	@echo "Logging in to Docker Hub..."
	docker login

# Development run with volume mount
.PHONY: dev-run
dev-run:
	@echo "Running container with volume mount for development..."
	docker run -d --name $(IMAGE_NAME)-dev -p 8080:8080 -v $(PWD)/templates:/root/templates $(IMAGE_NAME):$(TAG)
	@echo "Dev container is running at http://localhost:8080"

# View logs
.PHONY: logs
logs:
	docker logs -f $(IMAGE_NAME)

# Interactive shell in container
.PHONY: shell
shell:
	docker exec -it $(IMAGE_NAME) /bin/sh

# Full build and deploy cycle
.PHONY: deploy
deploy: clean login push
	@echo "Deployment complete!"