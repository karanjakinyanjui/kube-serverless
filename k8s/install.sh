#!/bin/bash
set -e

echo "Installing Kube-Serverless Platform..."

# Create namespace
echo "Creating namespace..."
kubectl apply -f namespace.yaml

# Install CRD
echo "Installing Custom Resource Definition..."
kubectl apply -f crd.yaml

# Wait for CRD to be established
echo "Waiting for CRD to be ready..."
kubectl wait --for condition=established --timeout=60s crd/functions.serverless.kube.io

# Create RBAC
echo "Creating RBAC resources..."
kubectl apply -f rbac.yaml

# Create ConfigMap
echo "Creating ConfigMap..."
kubectl apply -f configmap.yaml

# Deploy monitoring
echo "Deploying monitoring stack..."
kubectl apply -f monitoring.yaml

# Deploy API server
echo "Deploying API server..."
kubectl apply -f api-deployment.yaml

# Wait for API to be ready
echo "Waiting for API server to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/kube-serverless-api -n kube-serverless

echo ""
echo "Installation complete!"
echo ""
echo "To access the API server:"
echo "  kubectl get svc kube-serverless-api -n kube-serverless"
echo ""
echo "To view logs:"
echo "  kubectl logs -f deployment/kube-serverless-api -n kube-serverless"
