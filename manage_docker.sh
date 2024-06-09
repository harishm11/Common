#!/bin/bash

IMAGE_NAME="api-gateway"
CONTAINER_NAME="api-gateway-container"
PORT="3000"

function build_image() {
    echo "Building Docker image..."
    docker build -t $IMAGE_NAME .
}

function run_container() {
    echo "Running Docker container..."
    docker run -d -p $PORT:$PORT --name $CONTAINER_NAME $IMAGE_NAME
}

function stop_container() {
    echo "Stopping Docker container..."
    docker stop $CONTAINER_NAME
}

function remove_container() {
    echo "Removing Docker container..."
    docker rm $CONTAINER_NAME
}

function show_usage() {
    echo "Usage: $0 {build|run|stop|remove|restart}"
}

function restart_container() {
    stop_container
    remove_container
    run_container
}

case "$1" in
    build)
        build_image
        ;;
    run)
        run_container
        ;;
    stop)
        stop_container
        ;;
    remove)
        remove_container
        ;;
    restart)
        restart_container
        ;;
    *)
        show_usage
        exit 1
        ;;
esac
