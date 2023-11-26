#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 <protocol>"
    echo "Example: $0 http"
    exit 1
}

# Check if a service name is provided
if [ -z "$1" ]; then
    usage
fi

PROTOCOL_NAME="$1"

echo "Running Mean Application Response Time Test..."
sleep 1

echo "16 bytes payload size"
k6 run -e Q=1 -e RES=16 mean_app_response_"$PROTOCOL_NAME".js
sleep 2

echo "256 bytes payload size"
k6 run -e Q=1 -e RES=256 mean_app_response_"$PROTOCOL_NAME".js
sleep 2

echo "4 kb payload size"
k6 run -e Q=1 -e RES=4k mean_app_response_"$PROTOCOL_NAME".js
sleep 2

echo "64 kb payload size"
k6 run -e Q=1 -e RES=64k mean_app_response_"$PROTOCOL_NAME".js
sleep 2

echo "1 mb payload size"
k6 run -e Q=1 -e RES=1m mean_app_response_"$PROTOCOL_NAME".js
sleep 2

echo "Test Finished"