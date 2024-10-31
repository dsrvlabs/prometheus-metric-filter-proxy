.PHONY: build

# define variable APP
APP=prom-proxy

all: lint test build

lint:
	@golint ./...

test:
	@go test -v ./...

build:
	@echo "Building..."
	@go build -o bin/$(APP) cmd/$(APP)/main.go
