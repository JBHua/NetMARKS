#!/bin/bash

# Step 0: Remove old Docker image. Yes, manually remove it
docker image rm --no-prune netmarks-scheduler

# Step 1: Build the Docker image and tag it as "netmarks-scheduler:latest"
#docker build --build-arg CACHEBUST=$(date +%s) -t netmarks-scheduler .
docker build --no-cache -t netmarks-scheduler:"$1" .

# Step 2: Load local image into minikube
minikube image load netmarks-scheduler:"$1"

# Step 3: Delete previous scheduler deployment
kubectl delete -n kube-system deployment netmark-scheduler

# Step 4: Deploy the newer version
kubectl apply -f ./infra/deployment.yaml
