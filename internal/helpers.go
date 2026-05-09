package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func checkKubePath() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	return validateKubePath(cfg.KubeManifestsPath)
}

func validateKubePath(kubePath string) (string, error) {
	if kubePath == "" {
		return "", fmt.Errorf("kube-manifests path not configured — run 'runslit config'")
	}
	fi, err := os.Stat(kubePath)
	if err != nil || !fi.IsDir() {
		return "", fmt.Errorf("kube-manifests path does not exist: %s", kubePath)
	}
	return kubePath, nil
}

func chartPath(kubePath, relChart string) string {
	return filepath.Join(kubePath, relChart)
}

var devstackLabelRe = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

func validateDevstackLabel(label string) error {
	if label == "" {
		return fmt.Errorf("DEVSTACK_LABEL cannot be empty")
	}
	if !devstackLabelRe.MatchString(label) {
		return fmt.Errorf("DEVSTACK_LABEL can only contain alphanumeric characters and hyphens")
	}
	return nil
}
