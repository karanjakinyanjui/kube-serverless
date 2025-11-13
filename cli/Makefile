BINARY_NAME=ksls
INSTALL_PATH=/usr/local/bin

.PHONY: all build install clean

all: build

build:
	go build -o ${BINARY_NAME} .

install: build
	cp ${BINARY_NAME} ${INSTALL_PATH}/

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test -v ./...
