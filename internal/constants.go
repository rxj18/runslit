package internal

import (
	"os"
	"path/filepath"
)

const (
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Yellow = "\033[1;33m"
	Blue   = "\033[0;34m"
	Cyan   = "\033[0;36m"
	Reset  = "\033[0m"
)

var (
	ConfigFile       = filepath.Join(os.Getenv("HOME"), ".runslit.config")
	SlitHelmfileName = "helmfile-slit.yaml.gotmpl"
)
