.PHONY: all build-api build-cli build-runtimes build-ui install deploy clean

all: build-api build-cli build-runtimes build-ui

# Build API server
build-api:
	cd api && make build

# Build CLI
build-cli:
	cd cli && make build

# Build runtime images
build-runtimes:
	docker build -t kube-serverless-nodejs:latest runtimes/nodejs
	docker build -t kube-serverless-python:latest runtimes/python
	docker build -t kube-serverless-go:latest runtimes/go

# Build UI
build-ui:
	docker build -t kube-serverless-ui:latest ui

# Build API Docker image
docker-api:
	cd api && make docker-build

# Install platform to Kubernetes
install:
	cd k8s && ./install.sh

# Deploy example functions
examples:
	cd cli && ./ksls deploy -f ../examples/nodejs-hello.yaml
	cd cli && ./ksls deploy -f ../examples/python-data-processor.yaml

# Clean build artifacts
clean:
	cd api && make clean
	cd cli && make clean

# Run tests
test:
	cd api && make test
	cd cli && make test

# Development setup
dev:
	kubectl port-forward -n kube-serverless svc/kube-serverless-api 8080:80 &
	cd ui && npm start
