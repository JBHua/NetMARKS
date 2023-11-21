#!/bin/bash

# Step 1: Build the Docker image and tag it as "netmarks-scheduler:latest"
docker build -t netmarks-scheduler .

# Step 2: Load local image into minikube
minikube image load netmarks-scheduler:latest

# Step 3: Delete previous scheduler deployment
kubectl delete -n kube-system deployment netmark-scheduler

# Step 4: Deploy the newer version
kubectl apply -f ./infra/deployment.yaml
