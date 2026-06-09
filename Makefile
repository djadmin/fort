.PHONY: build test clean install lint build-landing serve-landing tag gen check-gen

BINARY  := fort
# Derived from git so it never goes stale. Clean tag -> "0.3.0"; dev build -> "0.3.0-19-gabc123-dirty".
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//')
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

# Regenerate committed artifacts (currently the sample report). Deterministic:
# same input -> same output, so the result is safe to commit. Run after changing
# a check, the report template, or the sample data.
gen:
	$(GO) run ./cmd/sample-report

# Fail if the committed sample report is out of date. CI runs this to guarantee
# landing/sample-report.html always matches the checks and template.
check-gen: gen
	@git diff --exit-code landing/sample-report.html \
		|| { echo "landing/sample-report.html is stale — run 'make gen' and commit the result"; exit 1; }

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
