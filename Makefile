.PHONY: build

# define variable APP
APP=prom-proxy

build:
	@echo "Building..."
	@go build -o bin/$(APP) cmd/$(APP)/main.go
