#!/bin/bash

# Declare the name of the image and the container
IMAGE_NAME="ascii-image-1"
CONTAINER_NAME="ascii-container-1"
PORT_MAPPING="8080:8080"

# Build the Docker image
echo "Building the Docker image..."
docker build -f Dockerfile -t $IMAGE_NAME .

# List Docker images
echo "Listing Docker images..."
docker images

# Run the Docker container
echo "Running the Docker container..."
docker container run -p $PORT_MAPPING --detach --name $CONTAINER_NAME $IMAGE_NAME

# List all containers
echo "Listing all Docker containers..."
docker ps -a

# Execute a command in the running container
echo "Executing a command in the running container..."
docker exec -it $CONTAINER_NAME /bin/bash
