package internal

import "fmt"

func ShowStatus() {

	fmt.Printf("%sConfiguration:%s\n", Blue, Reset)

	cfg, err := loadConfig()
	if err != nil {
		fatal("failed to load config")
	}

	kubePath := cfg.KubeManifestsPath
	if kubePath == "" {
		fmt.Printf("%s  ✗ not configured%s\n\n", Red, Reset)
		fmt.Println("  Run 'runslit config' to get started.")
		return
	}

	fmt.Printf("  kube-manifests: %s\n", kubePath)
	fmt.Printf("  config file:    %s\n", ConfigFile)

	if _, err := validateKubePath(kubePath); err != nil {
		fmt.Printf("%s  ✗ %s%s\n", Red, err.Error(), Reset)
		return
	}

	fmt.Println()
	fmt.Printf("%sSLIT:%s\n", Blue, Reset)

	if cfg.DevstackLabel == "" {
		fmt.Printf("%s  → not configured — run 'runslit config'%s\n", Yellow, Reset)
		return
	}

	fmt.Printf("  devstack label: %s\n", cfg.DevstackLabel)

	fmt.Println()
	fmt.Printf("%sImages:%s\n", Blue, Reset)
	printImage("payments-nbplus", cfg.NBPlusImage)
	printImage("mock-go        ", cfg.MockGWImage)
	fmt.Printf("  TTL:             %s\n", cfg.ttl())
}

func printImage(name, image string) {
	if image != "" {
		fmt.Printf("  %s: %s\n", name, image)
	} else {
		fmt.Printf("  %s: %s(not set)%s\n", name, Yellow, Reset)
	}
}
