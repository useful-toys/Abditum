.PHONY: build test lint vet clean

BINARY=abditum
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY) ./cmd/abditum

test:
	CGO_ENABLED=0 go test ./... -race -count=1 -v

lint:
	golangci-lint run

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
