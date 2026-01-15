package internal

import (
	"fmt"
	"os"
)

func ShowStatus() {
	printBanner()

	fmt.Printf("%sGlobal Configuration:%s\n", Blue, Reset)

	cfg, err := loadConfig()
	if err != nil {
		fatal("failed to load global config")
	}

	kubePath := cfg.KubeManifestsPath
	if kubePath == "" {
		fmt.Printf("%s✗ kube-manifests path not configured%s\n\n", Red, Reset)
		fmt.Println("Run 'runslit config' to set it.")
		return
	}

	fmt.Printf("  kube-manifests: %s\n", kubePath)
	fmt.Printf("  Config file:    %s\n", ConfigFile)

	// Validate kube path (reuses already loaded config)
	if _, err := validateKubePath(kubePath); err != nil {
		fmt.Println()
		fmt.Printf("%s✗ %s%s\n", Red, err.Error(), Reset)
		return
	}

	fmt.Println()
	fmt.Printf("%sSLIT Configuration:%s\n", Blue, Reset)

	slitEnv := cfg.SlitEnv
	devstackLabel := cfg.DevstackLabel

	if slitEnv == "" || devstackLabel == "" {
		fmt.Printf("%s→ No SLIT environment initialized.%s\n", Yellow, Reset)
		fmt.Println("Run 'runslit init' to initialize.")
		return
	}

	fmt.Printf("  Environment:    %s\n", slitEnv)
	fmt.Printf("  Devstack Label: %s\n", devstackLabel)
	fmt.Println()

	helmfilePath := getHelmfilePath(kubePath)
	if _, err := os.Stat(helmfilePath); err == nil {
		fmt.Printf("%s✓ SLIT helmfile exists%s\n", Green, Reset)
	} else {
		fmt.Printf("%s→ SLIT helmfile not found%s\n", Yellow, Reset)
	}
}
