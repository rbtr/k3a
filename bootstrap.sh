#!/bin/bash
set -xe
K3A=$(which k3a || echo "/usr/local/bin/k3a")
K3A_CLUSTER=$(whoami)-hyperscale
K3A_SUBSCRIPTION=$AZURE_SUBSCRIPTION
go build -o $K3A ./cmd/k3a
k3a cluster create --cluster $K3A_CLUSTER --subscription $K3A_SUBSCRIPTION --region eastus2
k3a pool create --k8s-version 1.34.0 --cluster $K3A_CLUSTER --name k3a-controlplane --instance-count 1 --subscription $K3A_SUBSCRIPTION --role control-plane --sku Standard_D96s_v5 --ssh-key ~/.ssh/id_ed25519.pub 
k3a kubeconfig --cluster $K3A_CLUSTER > ~/.kube/config
k3a pool create --cluster $K3A_CLUSTER --name etcd --size Standard_E96as_v6 --count 1 --ssh-key ~/.ssh/id_ed25519.pub


k3a pool create --k8s-version 1.34.0 --region eastus2 --cluster $K3A_CLUSTER --name k3a-controlplane --instance-count 1 --subscription $K3A_SUBSCRIPTION --role control-plane --sku Standard_D96s_v5 --ssh-key ~/.ssh/id_ed25519.pub --etcd-addr 10.1.0.4

helm upgrade --install -n kube-system cilium cilium/cilium --version 1.19.0-pre.0 --values cilium-values.yaml

k3a pool create --cluster $K3A_CLUSTER --name k3a-worker-1 --role worker --sku Standard_D2_v3 --instance-count 10 --region eastus2 --subscription $K3A_SUBSCRIPTION --ssh-key ~/.ssh/id_ed25519.pub

---
export AZURE_SUBSCRIPTION=37deca37-c375-4a14-b90a-043849bd2bf1
export K3A_CLUSTER=evanbaker-hyperscale
kubectl delete node $(kubectl get nodes | grep NotReady | awk '{print $1;}')
