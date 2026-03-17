.PHONY: build run test lint clean docker

BINARY_NAME=presence
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

test:
	go test -race -cover ./...

test-coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME) coverage.out coverage.html

docker:
	docker build -t presence:$(VERSION) .

docker-run:
	docker compose up -d
