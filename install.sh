#!/usr/bin/env bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="runslit"
REPO_URL="https://github.com/rxj18/runslit.git"
TMP_DIR=$(mktemp -d)

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

echo -e "${BLUE}runslit installer${NC}"
echo ""

# Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed. Visit: https://golang.org/doc/install${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go $(go version | awk '{print $3}')${NC}"

# helm (required for sync/delete)
if ! command -v helm &> /dev/null; then
    echo -e "${RED}✗ helm is not installed (required for runslit sync/delete)${NC}"
    echo -e "${YELLOW}→ Visit: https://helm.sh/docs/intro/install/${NC}"
    exit 1
fi
echo -e "${GREEN}✓ helm $(helm version --short 2>/dev/null | head -1)${NC}"

# fzf (required for runslit test)
if ! command -v fzf &> /dev/null; then
    echo -e "${YELLOW}⚠  fzf not found (required for runslit test) — attempting install...${NC}"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew install fzf
        else
            git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install --bin
            export PATH="$PATH:$HOME/.fzf/bin"
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y fzf
        elif command -v yum &> /dev/null; then
            sudo yum install -y fzf
        elif command -v pacman &> /dev/null; then
            sudo pacman -S --noconfirm fzf
        else
            git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install --bin
            export PATH="$PATH:$HOME/.fzf/bin"
        fi
    else
        git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install --bin
        export PATH="$PATH:$HOME/.fzf/bin"
    fi

    if command -v fzf &> /dev/null; then
        echo -e "${GREEN}✓ fzf installed${NC}"
    else
        echo -e "${YELLOW}⚠  fzf install failed — runslit test will not work until fzf is in PATH${NC}"
        echo -e "${YELLOW}→ Visit: https://github.com/junegunn/fzf#installation${NC}"
    fi
else
    echo -e "${GREEN}✓ fzf${NC}"
fi

# git
if ! command -v git &> /dev/null; then
    echo -e "${RED}✗ Git is not installed.${NC}"
    exit 1
fi

# Clone or build in place
if [ -f "go.mod" ] && grep -q "runslit" go.mod 2>/dev/null; then
    echo -e "${YELLOW}→ Building from current directory...${NC}"
    BUILD_DIR="."
else
    echo -e "${YELLOW}→ Cloning repository...${NC}"
    if ! git clone --depth 1 "$REPO_URL" "$TMP_DIR/runslit"; then
        echo -e "${RED}✗ Failed to clone repository${NC}"
        exit 1
    fi
    BUILD_DIR="$TMP_DIR/runslit"
fi

BUILT_BINARY="$TMP_DIR/${BINARY_NAME}-bin"

echo -e "${YELLOW}→ Building...${NC}"
if ! go build -C "$BUILD_DIR" -o "$BUILT_BINARY" .; then
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

mkdir -p "$INSTALL_DIR"
cp "$BUILT_BINARY" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo ""
echo -e "${GREEN}✓ Installed to $INSTALL_DIR/$BINARY_NAME${NC}"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠  $INSTALL_DIR is not in your PATH${NC}"
    echo -e "${YELLOW}→ Add to your shell profile:${NC}"
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    export PATH="$PATH:$INSTALL_DIR"
fi

echo ""
echo -e "${BLUE}Get started:${NC}"
echo "  runslit config    # Configure kube-manifests path, label, and images"
echo "  runslit help      # Show all commands"
