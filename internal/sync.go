package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Sync() {
	printBanner()

	kubePath, err := checkKubePath()
	if err != nil {
		fatal(err.Error())
	}

	cfg, err := loadConfig()
	if err != nil {
		fatal(err.Error())
	}

	helmfilePath := getHelmfilePath(kubePath)
	helmfileDir := filepath.Dir(helmfilePath)

	if _, err := os.Stat(helmfilePath); err != nil {
		fatal("SLIT helmfile not found — run 'runslit init'")
	}

	info(fmt.Sprintf(
		"Deploying SLIT (%s) with label '%s'",
		cfg.SlitEnv,
		cfg.DevstackLabel,
	))

	err = runCommand(
		helmfileDir,
		"helmfile",
		"-f", SlitHelmfileName,
		"sync",
	)
	if err != nil {
		fatal("helmfile apply failed")
	}

	success("SLIT environment deployed successfully")
}

// helper func ...
func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
