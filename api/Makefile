BINARY_NAME=kube-serverless-api
DOCKER_IMAGE=kube-serverless-api:latest

.PHONY: all build run test clean docker-build docker-push

all: build

build:
	go build -o ${BINARY_NAME} .

run:
	go run .

test:
	go test -v ./...

clean:
	go clean
	rm -f ${BINARY_NAME}

docker-build:
	docker build -t ${DOCKER_IMAGE} .

docker-push:
	docker push ${DOCKER_IMAGE}

install:
	go mod download
