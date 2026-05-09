# Runslit

A CLI utility for managing SLIT (Service Level Integration Testing) environments with Helmfile.

## Features

- 🚀 Easy SLIT environment initialization
- 🔧 Helmfile-based deployment management
- 🧪 Interactive test runner with fzf
- 📦 Zero external dependencies (except Go)
- 🎯 Devstack label management

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://gist.githubusercontent.com/rxj18/9225f1112813c7f8b50c13026ddc3664/raw/install.sh | bash
```

> **Note:** 
> - The installation script will automatically clone the repository, build, and install
> - `fzf` will be automatically installed if it's not already present on your system
> - Temporary files are cleaned up automatically

### Manual Installation

#### Option 1: Using the install script (from repository)

```bash
git clone https://github.com/rxj18/runslit.git
cd runslit
./install.sh
```

#### Option 2: Using Make

```bash
git clone https://github.com/rxj18/runslit.git
cd runslit
make install
```

#### Option 3: Using Go

```bash
go install github.com/rxj18/runslit@latest
```

#### Option 4: Build from source

```bash
git clone https://github.com/rxj18/runslit.git
cd runslit
go build -o runslit .
mkdir -p ~/.local/bin
mv runslit ~/.local/bin/
export PATH="$PATH:$HOME/.local/bin"  # Add to your shell profile
```

### Custom Installation Directory

The default installation directory is `~/.local/bin`. To install to a different location:

```bash
# Install to /usr/local/bin
INSTALL_DIR=/usr/local/bin ./install.sh

# Or with Make
make install INSTALL_DIR=/usr/local/bin
```

## Prerequisites

- Go 1.21+ (for building)
- Helmfile (for deployment)
- kubectl (configured with cluster access)
- fzf (for interactive test selection) - **automatically installed by install script**

## Quick Start

### 1. Configure kube-manifests path

```bash
runslit config
```

### 2. Initialize SLIT environment

```bash
runslit init
```

Select your environment (stage/slit) and provide a devstack label.

### 3. Deploy your environment

```bash
runslit sync
```

### 4. Run tests

```bash
runslit test
```

### 5. Check status

```bash
runslit status
```

### 6. Destroy environment

```bash
runslit delete
```

## Commands

| Command | Description |
|---------|-------------|
| `config` | Set/update kube-manifests path |
| `init` | Initialize SLIT environment |
| `sync` | Deploy/sync the SLIT helmfile |
| `test` | Select and run tests interactively |
| `status` | Show current configuration |
| `delete` | Destroy the SLIT deployment |
| `help` | Show help message |

## Configuration

Configuration is stored in `~/.runslit.config` as JSON:

```json
{
  "kube_manifests_path": "/path/to/kube-manifests",
  "runslit_install_dir": "/Users/you/.runslit",
  "slit_env": "slit",
  "devstack_label": "your-label"
}
```

## Test Runner

The `runslit test` command provides an interactive test runner:

1. Scans `./slit` directory for test files
2. Extracts test suites and test cases using Go AST
3. Presents tests in fzf for selection
4. Runs selected test with `DEVSTACK_LABEL` environment variable

**Example:**
```
TestNetBankingPayment -> TestVerifySuccess | ./slit/netbanking
```

Runs: `DEVSTACK_LABEL=your-label go test -v -run TestNetBankingPayment/TestVerifySuccess ./slit/netbanking`

## Uninstallation

```bash
# Using Make
make uninstall

# Or manually
rm ~/.local/bin/runslit
rm ~/.runslit.config
rm -rf ~/.runslit  # Optional: remove installation directory
```

## Development

### Build

```bash
make build
```

### Run locally

```bash
go run main.go <command>
```

### Clean

```bash
make clean
```

## Troubleshooting

### Command not found after installation

Make sure `~/.local/bin` is in your PATH:

```bash
echo $PATH | grep "$HOME/.local/bin"
```

If not, add it to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
export PATH="$PATH:$HOME/.local/bin"
```

Then reload your shell:

```bash
source ~/.bashrc  # or ~/.zshrc
```

### Permission denied

The install script may need sudo access. Run:

```bash
sudo ./install.sh
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
