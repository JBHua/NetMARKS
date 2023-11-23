#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 <service-name>"
    echo "Example: $0 fish"
    exit 1
}

# Check if a service name is provided
if [ -z "$1" ]; then
    usage
fi

# First argument is the service name
SERVICE_NAME="netmarks-$1"

# Create a temporary file
TMP_FILE=$(mktemp)

# Run minikube command in the background and redirect output to temporary file
minikube service $SERVICE_NAME > $TMP_FILE &
#MINIKUBE_PID=$!

## Wait for a few seconds and then send SIGINT to the minikube command
#{ sleep 3; kill -INT $MINIKUBE_PID; } &
#
sleep 3

# Wait for the minikube command to terminate


# Run minikube command to get service address and extract the URL
# Extract the service address from the temporary file
SERVICE_ADDRESS=$(grep -m 1 'http://' $TMP_FILE)

# Check if the service address was found
if [ -z "$SERVICE_ADDRESS" ]; then
    echo "Service address not found for $SERVICE_NAME"
    exit 1
fi

# Print the extracted service address
echo "Service address: $SERVICE_ADDRESS"

# Clean up - remove the temporary file
rm $TMP_FILE

# Run the k6 command with the extracted service address
k6 run -e SERVICE_URL="$SERVICE_ADDRESS" test.js

wait $MINIKUBE_PID