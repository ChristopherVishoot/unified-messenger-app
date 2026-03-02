#!/bin/bash

# Cross-compile for Linux amd64 (for Docker/Kubernetes)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Build Docker image
docker build -t golang-service:v4 .

# Load into kind cluster (uncomment if using kind)
kind load docker-image golang-service:v4 --name hello

