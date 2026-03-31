#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
CLUSTER_NAME="local"

# Create the cluster if it doesn't already exist
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo "Cluster '${CLUSTER_NAME}' already exists, skipping creation."
else
  kind create cluster --name "${CLUSTER_NAME}" --config "${SCRIPT_DIR}/kind.yaml"
fi

KUBE_CONTEXT="kind-${CLUSTER_NAME}"

# Install the flux-operator via Helm
helm upgrade --install flux-operator \
  oci://ghcr.io/controlplaneio-fluxcd/charts/flux-operator \
  --namespace flux-system \
  --create-namespace \
  --wait \
  --kube-context "${KUBE_CONTEXT}"

# Apply the FluxInstance to install Flux controllers
kubectl apply --context "${KUBE_CONTEXT}" -f "${SCRIPT_DIR}/flux-instance.yaml"

# Build and load app images into the cluster
for variant in server client; do
  echo "Building ${variant} image..."
  docker build --build-arg VARIANT="${variant}" -t "${variant}:latest" "${REPO_DIR}/app"
  echo "Loading ${variant} image into kind cluster..."
  kind load docker-image "${variant}:latest" --name "${CLUSTER_NAME}"
done
