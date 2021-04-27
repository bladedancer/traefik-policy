#!/bin/sh

echo ======================
echo === Create Cluster ===
echo ======================

k3d cluster delete traefikgw
k3d cluster create --no-lb --update-default-kubeconfig=false --wait traefikgw --volume `pwd`/config:/tmp/config
export KUBECONFIG=$(k3d kubeconfig write traefikgw)
kubectl cluster-info

docker pull traefik:v2.4

k3d image import --cluster traefikgw traefik:v2.4

echo ==========================
echo === Installing Traefik ===
echo ==========================
kubectl create configmap traefik --from-file=./config/traefik/traefik.yaml
kubectl apply -f deploy/traefik

echo "RUN:"
echo "export KUBECONFIG=$(k3d kubeconfig write traefikgw)"
