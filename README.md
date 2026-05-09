# runslit

A CLI tool for managing SLIT (Service Level Integration Testing) environments. Deploys `payments-nbplus` and `mock-go` releases directly via `helm`, and provides an interactive test runner for the `./slit` directory.

## Prerequisites

- Go 1.21+ (for building from source)
- [helm](https://helm.sh/docs/intro/install/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) configured with cluster access
- [fzf](https://github.com/junegunn/fzf) (for `runslit test`)

## Installation

### Quick install

```bash
curl -fsSL https://raw.githubusercontent.com/rxj18/runslit/main/install.sh | bash
```

The script checks for `helm` and `fzf`, clones the repo, builds the binary, and installs it to `~/.local/bin`.

### Custom install directory

```bash
INSTALL_DIR=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/rxj18/runslit/main/install.sh | bash
```

### From source

```bash
git clone https://github.com/rxj18/runslit.git
cd runslit
make install
```

### Via Go

```bash
go install github.com/rxj18/runslit@latest
```

## Usage

### 1. Configure

```bash
runslit config
```

Interactive menu to set:
- `kube-manifests path` — local path to your kube-manifests repo
- `devstack label` — your personal label (e.g. `rituraj`, `pr-123`)
- `payments-nbplus image SHA`
- `mock-go image SHA`

Configuration is saved to `~/.runslit.config`.

### 2. Deploy

```bash
runslit sync
```

Select which releases to deploy (payments-nbplus, mock-go, or both) and runslit runs `helm upgrade --install` for each.

### 3. Run tests

```bash
runslit test
```

Scans `./slit` for testify suite test cases, opens fzf for selection, and runs the chosen test with `DEVSTACK_LABEL` set. The last-run test is pre-selected on next invocation.

### 4. Check status

```bash
runslit status
```

### 5. Destroy

```bash
runslit delete
```

Select which releases to uninstall.

## Commands

| Command | Description |
|---------|-------------|
| `config` | Configure runslit interactively |
| `sync` | Deploy selected releases via helm |
| `delete` | Destroy selected releases via helm |
| `status` | Show current configuration |
| `test` | Select and run a test from `./slit` |
| `help` | Show help |

## Configuration file

```json
{
  "kube_manifests_path": "/path/to/kube-manifests",
  "devstack_label": "your-label",
  "nbplus_image": "abc123...",
  "mockgw_image": "def456..."
}
```

## Uninstall

```bash
make uninstall

# or manually
rm ~/.local/bin/runslit
rm ~/.runslit.config
```

## Development

```bash
make build       # build binary
make install     # build and install
make clean       # remove build artifacts
go run main.go   # run locally
```
