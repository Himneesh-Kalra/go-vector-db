# Makefile for GO-VECTOR-DB

CLI_DIR=cli
SERVER_DIR=server
BIN_DIR=bin
CLI_BIN=$(BIN_DIR)/vecdb-cli
SERVER_BIN=$(BIN_DIR)/vecdb

# Config files
ETC_DIR=/etc/vecdb
SERVER_CONFIG=server/server_config.json
CLI_CONFIG=cli/cli_config.json

.PHONY: all build clean install test bench fmt check-go

all: build


check-go:
	@command -v go >/dev/null 2>&1 || { \
		echo "❌ Go is not installed."; \
		echo "👉 Please install Go from https://golang.org/dl/ before proceeding."; \
		exit 1; \
	}

build: check-go
	@echo "🔨 Building server..."
	cd $(SERVER_DIR) && go build -o ../$(SERVER_BIN)
	@echo "🔨 Building CLI..."
	cd $(CLI_DIR) && go build -o ../$(CLI_BIN)

clean:
	@echo "🧹 Cleaning binaries..."
	rm -rf $(BIN_DIR)

install: build
	@echo "📂 Installing binaries to /usr/local/bin..."
	sudo cp $(CLI_BIN) /usr/local/bin/vecdb-cli
	sudo cp $(SERVER_BIN) /usr/local/bin/vecdb

	@echo "📁 Ensuring config directory exists at $(ETC_DIR)..."
	sudo mkdir -p $(ETC_DIR)

	@echo "📄 Copying server_config.json to $(ETC_DIR)..."
	sudo cp $(SERVER_CONFIG) $(ETC_DIR)/server_config.json

	@echo "📄 Copying cli_config.json to $(ETC_DIR)..."
	sudo cp $(CLI_CONFIG) $(ETC_DIR)/cli_config.json

	@echo "✅ Installation complete! You can now use vecdb and vecdb-cli globally."

# Run tests
test:
	go test -v ./...

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Format all code
fmt:
	go fmt ./...
