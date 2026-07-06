#!/bin/bash

IMAGE="forum:latest"
CONTAINER="forum-app"

docker build -t $IMAGE .

docker run -d -p 8080:8080 --name $CONTAINER $IMAGE

echo "list of images:"
docker images

echo "list of containers:"
docker ps