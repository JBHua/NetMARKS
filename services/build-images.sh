#!/bin/bash

# Usage function to display help for the script
usage() {
  echo "Usage: $0 [service_name|all] [version]"
  echo "service_name: the name of the service to build (must be a directory under services/)"
  echo "all: build all services"
  echo "version: the version tag to apply to the Docker images. E.g. 0.1.12"
  exit 1
}

# Check if the correct number of arguments was provided
if [ $# -ne 2 ]; then
  usage
fi

SERVICE_NAME=$1
VERSION=$2

# Function to build a Docker image for a service
build_image() {
  local service=$1
  local version=$2
  echo "Building Docker image for service: $service, version: $version"
  docker build -t "$service:$version" -f "./$service/Dockerfile" .. || { echo "Docker build failed for $service"; exit 1; }
}

# Function to push local image to minikube
push_image() {
  local service=$1
  local version=$2
  echo "Pushing image for service: $service, version: $version to minikube"
  minikube image load "$service":"$version"
}

# Build all services if the first argument is 'all'
if [ "$SERVICE_NAME" = "all" ]; then
  # Find all directories in the current directory (assumed to be services/)
  for dir in */ ; do
    # Remove the trailing slash to get the service name
    service="${dir%/}"
    build_image "$service" "$VERSION"
    push_image "$service" "$VERSION"
  done
else
  # Build the specified service if it is a directory
  if [ -d "$SERVICE_NAME" ]; then
    build_image "$SERVICE_NAME" "$VERSION"
    push_image "$SERVICE_NAME" "$VERSION"
  else
    echo "Error: Service directory '$SERVICE_NAME' does not exist."
    exit 1
  fi
fi

echo "Docker build script finished."
