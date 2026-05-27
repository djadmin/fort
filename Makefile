.PHONY: build test clean install lint build-landing serve-landing tag

BINARY  := fort
VERSION := 0.1.1
LDFLAGS := -ldflags="-X main.version=$(VERSION)"
GO      := mise exec -- go

build:
	$(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/fort

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY) fort-landing

install:
	GOBIN="$$HOME/.local/bin" $(GO) install $(LDFLAGS) ./cmd/fort

lint:
	$(GO) vet ./...

build-landing:
	$(GO) build -o fort-landing ./cmd/landing

serve-landing: build-landing
	./fort-landing

# Tag and push to trigger GoReleaser release workflow
tag:
	@read -p "Version (e.g. v0.1.0): " v && \
	git tag $$v && \
	git push origin $$v

.DEFAULT_GOAL := build
