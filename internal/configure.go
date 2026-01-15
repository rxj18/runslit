package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Configure() {
	printBanner()

	cfg, _ := loadConfig()
	currentPath := cfg.KubeManifestsPath

	fmt.Printf("%sConfigure kube-manifests path%s\n\n", Blue, Reset)

	if currentPath != "" {
		fmt.Printf("Current path: %s%s%s\n\n", Cyan, currentPath, Reset)
	}

	fmt.Print("Enter the path to your kube-manifests repository: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("→ No changes made")
		return
	}

	absPath, err := expandPath(input)
	if err != nil {
		fatal("invalid path")
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil || !fileInfo.IsDir() {
		fatal(fmt.Sprintf("directory does not exist: %s", absPath))
	}

	helmfileDir := filepath.Join(absPath, "helmfile")
	if _, err := os.Stat(helmfileDir); err != nil {
		fmt.Printf("%sThis doesn't look like a kube-manifests repo (no helmfile directory)%s\n", Red, Reset)
		fmt.Print("Continue anyway? [y/N]: ")

		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))
		if confirm != "y" {
			os.Exit(1)
		}
	}

	if err := saveConfig(absPath); err != nil {
		fatal("failed to save config")
	}

	fmt.Println()
	fmt.Printf("%s✓ kube-manifests path set to: %s%s\n", Green, absPath, Reset)
}

// helper func ...
func expandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, strings.TrimPrefix(p, "~"))
	}

	return filepath.Abs(p)
}

func saveConfig(kubePath string) error {
	// Load existing config to preserve other fields
	cfg, err := loadConfig()
	if err != nil {
		cfg = &Config{}
	}

	// Update only the fields we're setting
	cfg.KubeManifestsPath = kubePath
	cfg.RunslitInstallDir = filepath.Join(os.Getenv("HOME"), ".runslit")

	return cfg.saveConfigToFile()
}
