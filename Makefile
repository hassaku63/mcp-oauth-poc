SHELL := /bin/sh

# Config
ISSUER ?= http://localhost:8080
AUTHORIZE_ENDPOINT ?= $(ISSUER)/oauth2/authorize
TOKEN_ENDPOINT ?= $(ISSUER)/oauth2/token

BIN_DIR := bin
AS_BIN := $(BIN_DIR)/as-poc
CLI_BIN := $(BIN_DIR)/auth-cli

.PHONY: all build build-as build-auth-cli run-as auth-new auth-url auth-complete fmt clean help

all: build

build: build-as build-auth-cli

build-as:
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -o $(AS_BIN) ./cmd/as-poc
	@echo "Built $(AS_BIN)"

build-auth-cli:
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -o $(CLI_BIN) ./cmd/auth-cli
	@echo "Built $(CLI_BIN)"

run-as:
	GO111MODULE=on go run ./cmd/as-poc

# Convenience wrappers (edit args as needed)
auth-new:
	@echo "Example: $(CLI_BIN) new --redirect-uri http://127.0.0.1:53219/callback --resource $(ISSUER) --scope 'mcp.read mcp.write' --client-id mcp-cli-12345" && \
	true

auth-url:
	@echo "Example: $(CLI_BIN) url --session-id <id> --authorize-endpoint $(AUTHORIZE_ENDPOINT) --open" && \
	true

auth-complete:
	@echo "Example: $(CLI_BIN) complete --session-id <id> --token-endpoint $(TOKEN_ENDPOINT) --client-id mcp-cli-12345 --callback-url '<full_url>'" && \
	true

fmt:
	@gofmt -s -w .

clean:
	rm -rf $(BIN_DIR) .auth-cli-sessions

help:
	@echo "Targets:"
	@echo "  build            Build AS and CLI binaries"
	@echo "  run-as           Run Authorization Server PoC on :8080"
	@echo "  auth-new         Print example command to create a new auth session"
	@echo "  auth-url         Print example command to generate the authorize URL"
	@echo "  auth-complete    Print example command to complete the flow (token exchange)"
	@echo "  fmt              gofmt -s -w ."
	@echo "  clean            Remove build artifacts and CLI sessions"

