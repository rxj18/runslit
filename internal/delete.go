package internal

import (
	"os"
	"path/filepath"
)

func Delete() {
	printBanner()

	kubePath, err := checkKubePath()
	if err != nil {
		fatal(err.Error())
	}

	helmfilePath := getHelmfilePath(kubePath)
	helmfileDir := filepath.Dir(helmfilePath)

	if _, err := os.Stat(helmfilePath); err != nil {
		fatal("SLIT helmfile not found")
	}

	info("Destroying SLIT environment...")

	err = runCommand(
		helmfileDir,
		"helmfile",
		"-f", SlitHelmfileName,
		"destroy",
	)
	if err != nil {
		fatal("helmfile destroy failed")
	}

	success("SLIT environment destroyed")
}
