APP_NAME := gwm
VERSION := $(shell cat VERSION)

.PHONY: help build run docker-build docker-run clean
all: help

# Default target
help:
	@echo "Available targets:"
	@echo "  build           Build the binary"
	@echo "  run             Run the binary"
	@echo "  deps            Update our dependencies"
	@echo "  clean           Clean up the binary and Docker image"
	@echo "  docker-build    Build the Docker image"
	@echo "  docker-run      Run the Docker container"
	@echo ""
	@echo "Use 'make <target>' to run a specific target."

# Build the binary
build:
	go build -o ./bin/$(APP_NAME) ./src

# Run the binary
run:
	./bin/$(APP_NAME)

# Update our dependencies
deps:
	go mod download
	go mod tidy

# Clean up the binary and Docker image
clean:
	rm -f ./bin/$(APP_NAME)
	docker rmi -f $(APP_NAME):$(VERSION)

# Build the Docker image
docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

# Run the Docker container
docker-run:
	docker run --rm -it --env DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix $(APP_NAME):$(VERSION)

