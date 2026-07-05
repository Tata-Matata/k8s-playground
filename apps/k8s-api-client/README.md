# k8s-api-client

Go web application that calls the Kubernetes API to retrieve certain data and renders it as HTML on different endpoints like `/pods`.

## Features

- Uses in-cluster service account credentials when running in Kubernetes.
- Falls back to `KUBECONFIG`, then `$HOME/.kube/config`, for local development.
- Renders a simple HTML page listing pods, namespace, phase, IP, node, and age.
- Includes starter RBAC and Deployment manifests for lab scenarios.

## Run locally

```bash
GOTOOLCHAIN=local go mod tidy
GOTOOLCHAIN=local go run .
```

Then open `http://localhost:8080/pods?namespace=default`.

To use a non-default kubeconfig file locally:

```bash
KUBECONFIG=/path/to/kubeconfig GOTOOLCHAIN=local go run .
```

## Build container

```bash
docker build -t docker.io/your-dockerhub-user/k8s-api-client:latest .
docker push docker.io/your-dockerhub-user/k8s-api-client:latest
```

## Deploy to Kubernetes

Update the image in `deploy/deployment.yaml`, then apply:

```bash
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/deployment.yaml
kubectl port-forward deployment/k8s-api-client 8080:8080
```

