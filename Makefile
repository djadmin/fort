.PHONY: build test clean install lint

BINARY  := fort
VERSION := 0.1.0
LDFLAGS := -ldflags="-X main.version=$(VERSION)"
GO      := mise exec -- go

build:
	$(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/fort

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY)

install:
	GOBIN="$$HOME/.local/bin" $(GO) install $(LDFLAGS) ./cmd/fort

lint:
	$(GO) vet ./...

.DEFAULT_GOAL := build
