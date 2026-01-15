.PHONY: build install uninstall clean test help

BINARY_NAME=runslit
INSTALL_DIR=/usr/local/bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✓ Build complete: ./$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed: $(INSTALL_DIR)/$(BINARY_NAME)"

uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled"

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@echo "✓ Clean complete"

test:
	@echo "Running tests..."
	@go test ./...

help:
	@echo "Runslit Makefile Commands:"
	@echo "  make build      - Build the binary"
	@echo "  make install    - Build and install to $(INSTALL_DIR)"
	@echo "  make uninstall  - Remove from $(INSTALL_DIR)"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
