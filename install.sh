#!/usr/bin/env bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Installation directory (defaults to /usr/local/bin, can be overridden)
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="runslit"
REPO_URL="https://github.com/rxj18/runslit.git"
TMP_DIR=$(mktemp -d)

# Cleanup function
cleanup() {
    if [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}
trap cleanup EXIT

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║       Runslit Installation Script      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed. Please install Go first.${NC}"
    echo -e "${YELLOW}→ Visit: https://golang.org/doc/install${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed ($(go version))${NC}"

# Check if fzf is installed
if ! command -v fzf &> /dev/null; then
    echo -e "${YELLOW}⚠️  fzf is not installed (required for 'runslit test' command)${NC}"
    echo -e "${YELLOW}→ Attempting to install fzf...${NC}"
    
    # Detect OS and install fzf
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            echo -e "${YELLOW}→ Installing fzf via Homebrew...${NC}"
            brew install fzf
        else
            echo -e "${YELLOW}→ Homebrew not found. Installing fzf manually...${NC}"
            git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
            ~/.fzf/install --bin
            export PATH="$PATH:$HOME/.fzf/bin"
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if command -v apt-get &> /dev/null; then
            echo -e "${YELLOW}→ Installing fzf via apt...${NC}"
            sudo apt-get update && sudo apt-get install -y fzf
        elif command -v yum &> /dev/null; then
            echo -e "${YELLOW}→ Installing fzf via yum...${NC}"
            sudo yum install -y fzf
        elif command -v pacman &> /dev/null; then
            echo -e "${YELLOW}→ Installing fzf via pacman...${NC}"
            sudo pacman -S --noconfirm fzf
        else
            echo -e "${YELLOW}→ Installing fzf manually...${NC}"
            git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
            ~/.fzf/install --bin
            export PATH="$PATH:$HOME/.fzf/bin"
        fi
    else
        echo -e "${YELLOW}→ Unsupported OS. Installing fzf manually...${NC}"
        git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
        ~/.fzf/install --bin
        export PATH="$PATH:$HOME/.fzf/bin"
    fi
    
    # Verify fzf installation
    if command -v fzf &> /dev/null; then
        echo -e "${GREEN}✓ fzf installed successfully${NC}"
    else
        echo -e "${YELLOW}⚠️  fzf installation may require manual setup${NC}"
        echo -e "${YELLOW}→ Visit: https://github.com/junegunn/fzf#installation${NC}"
        echo -e "${YELLOW}→ runslit will work, but 'runslit test' requires fzf${NC}"
    fi
else
    echo -e "${GREEN}✓ fzf is installed${NC}"
fi

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo -e "${RED}✗ Git is not installed. Please install Git first.${NC}"
    exit 1
fi

# Clone or use existing repository
if [ -f "go.mod" ] && grep -q "runslit" go.mod 2>/dev/null; then
    # We're already in the runslit directory
    echo -e "${GREEN}✓ Running from runslit repository${NC}"
    BUILD_DIR="."
else
    # Download the repository
    echo -e "${YELLOW}→ Cloning runslit repository...${NC}"
    if git clone --depth 1 "$REPO_URL" "$TMP_DIR/runslit"; then
        echo -e "${GREEN}✓ Repository cloned${NC}"
        BUILD_DIR="$TMP_DIR/runslit"
    else
        echo -e "${RED}✗ Failed to clone repository${NC}"
        exit 1
    fi
fi

# Build the binary
echo -e "${YELLOW}→ Building runslit...${NC}"
cd "$BUILD_DIR"
if go build -o "$TMP_DIR/$BINARY_NAME" .; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Check if we need sudo for installation
if [ -w "$INSTALL_DIR" ]; then
    SUDO=""
else
    SUDO="sudo"
    echo -e "${YELLOW}→ Requesting sudo access to install to $INSTALL_DIR${NC}"
fi

# Install the binary
echo -e "${YELLOW}→ Installing to $INSTALL_DIR/$BINARY_NAME...${NC}"
if $SUDO cp "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"; then
    echo -e "${GREEN}✓ Installed successfully${NC}"
else
    echo -e "${RED}✗ Installation failed${NC}"
    exit 1
fi

# Make it executable
$SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Verify installation
if command -v runslit &> /dev/null; then
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║     Installation Complete! 🎉         ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Installed at:${NC} $INSTALL_DIR/$BINARY_NAME"
    echo -e "${BLUE}Version:${NC} $(runslit help | head -n 1)"
    echo ""
    echo -e "${YELLOW}Get started:${NC}"
    echo "  runslit config    # Configure kube-manifests path"
    echo "  runslit init      # Initialize SLIT environment"
    echo "  runslit help      # Show all commands"
else
    echo ""
    echo -e "${YELLOW}⚠️  Installation complete, but 'runslit' is not in PATH${NC}"
    echo -e "${YELLOW}→ Add $INSTALL_DIR to your PATH or run:${NC}"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
fi
