#!/bin/bash
# Function to display usage information
usage() {
    echo "Usage: $0 <protocol> <service>"
    echo "Example: $0 http fish"
    exit 1
}

# Check if a service name is provided
if [ -z "$1" ]; then
    usage
fi


PROTOCOL="$1"
SVC_NAME="$2"


echo "Running Mean Application Response Time Test..."
sleep 1

echo "16 bytes payload size"
k6 run -e Q=1 -e RES=16 -e SVC_NAME="$SVC_NAME" mean_app_response_"$PROTOCOL".js
sleep 2

echo "256 bytes payload size"
k6 run -e Q=1 -e RES=256 -e SVC_NAME="$SVC_NAME" mean_app_response_"$PROTOCOL".js
sleep 2

echo "4 kb payload size"
k6 run -e Q=1 -e RES=4k -e SVC_NAME="$SVC_NAME" mean_app_response_"$PROTOCOL".js
sleep 2

echo "64 kb payload size"
k6 run -e Q=1 -e RES=64k -e SVC_NAME="$SVC_NAME" mean_app_response_"$PROTOCOL".js
sleep 2

echo "1 mb payload size"
k6 run -e Q=1 -e RES=1m -e SVC_NAME="$SVC_NAME" mean_app_response_"$PROTOCOL".js

echo "Test Finished"