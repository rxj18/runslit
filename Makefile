.PHONY: build install uninstall clean test help

BINARY_NAME=runslit
INSTALL_DIR=$(HOME)/.local/bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✓ Build complete: ./$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed: $(INSTALL_DIR)/$(BINARY_NAME)"
	@if ! echo $$PATH | grep -q "$(INSTALL_DIR)"; then \
		echo "⚠️  $(INSTALL_DIR) is not in your PATH"; \
		echo "→ Add to your shell profile: export PATH=\"\$$PATH:$(INSTALL_DIR)\""; \
	fi

uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
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
