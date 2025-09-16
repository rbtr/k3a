#!/bin/bash
set -xe
echo $K3A_CLUSTER
echo $K3A_SUBSCRIPTION
go build -o k3a ./cmd/k3a && echo "Build successful"
# Create cluster infrastructure with integrated PostgreSQL Flexible Server
./k3a pool delete --cluster $K3A_CLUSTER --name k3a-worker-0 
./k3a pool delete --cluster $K3A_CLUSTER --name k3a-controlplane 
./k3a cluster delete --cluster $K3A_CLUSTER --subscription $K3A_SUBSCRIPTION
