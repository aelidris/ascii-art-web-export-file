#!/bin/bash

# Declare the name of the image and the container
CONTAINER_NAME="ascii-container-1"
IMAGE_NAME="ascii-image-1"

# Stop the container
echo "Stop the container"
docker stop $CONTAINER_NAME

# Listing the runing containers
echo "Listing the runing container ..."
docker ps 

# Remove the container
echo "Remove the container"
docker rm $CONTAINER_NAME

# listing all containers
echo "listing all containers"
docker ps -a

# Remove the image of container
echo "Remove the image of container"
docker image rm $IMAGE_NAME

# Listing all images
echo "Listing all images"
docker images

