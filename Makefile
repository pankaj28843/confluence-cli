.PHONY: build install uninstall test race coverage lint setup clean release

BINARY := confluence
INSTALL_DIR := $(HOME)/.local/bin

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commit=$(COMMIT)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/confluence/

install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed $(BINARY) to $(INSTALL_DIR)/$(BINARY)"

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "Removed $(INSTALL_DIR)/$(BINARY)"

test:
	go test ./... -count=1 -timeout 60s

race:
	go test ./... -race -count=1 -timeout 60s

coverage:
	go test ./internal/... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out | grep total
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' > COVERAGE.txt

lint:
	gofmt -w .
	go vet ./...
	go mod tidy

setup:
	cp scripts/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed"

release:
	@mkdir -p build
	FULL_LDFLAGS="-s -w $(LDFLAGS)"; \
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$$FULL_LDFLAGS" -o build/$(BINARY)-darwin-arm64      ./cmd/confluence/ && \
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$$FULL_LDFLAGS" -o build/$(BINARY)-darwin-amd64      ./cmd/confluence/ && \
	GOOS=linux   GOARCH=amd64 go build -ldflags "$$FULL_LDFLAGS" -o build/$(BINARY)-linux-amd64       ./cmd/confluence/ && \
	GOOS=windows GOARCH=amd64 go build -ldflags "$$FULL_LDFLAGS" -o build/$(BINARY)-windows-amd64.exe ./cmd/confluence/
	@echo "Built:" && ls -lh build/

clean:
	rm -f $(BINARY) coverage.out
	rm -rf build/
