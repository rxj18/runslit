package internal

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Config represents the runslit configuration
type Config struct {
	KubeManifestsPath string `json:"kube_manifests_path,omitempty"`
	DevstackLabel     string `json:"devstack_label,omitempty"`
	NBPlusImage       string `json:"nbplus_image,omitempty"`
	MockGWImage       string `json:"mockgw_image,omitempty"`
	TTL               string `json:"ttl,omitempty"`
	LastTest          string `json:"last_test,omitempty"`
}

func (c *Config) ttl() string {
	if c.TTL != "" {
		return c.TTL
	}
	return "12h"
}

var (
	cachedConfig  *Config
	cachedModTime time.Time
)

func loadConfig() (*Config, error) {
	// Get file info to check modification time
	fileInfo, err := os.Stat(ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}

	modTime := fileInfo.ModTime()

	// Check if we have a cached version and if the file hasn't been modified
	if cachedConfig != nil && modTime.Equal(cachedModTime) {
		return cachedConfig, nil
	}

	// File has been modified or not cached, read it
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}

	// Parse JSON config
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Update cache
	cachedConfig = cfg
	cachedModTime = modTime

	return cfg, nil
}

func (cfg *Config) saveConfigToFile() error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(ConfigFile, data, 0644); err != nil {
		return err
	}

	// Invalidate cache after writing
	cachedConfig = nil
	cachedModTime = time.Time{}

	return nil
}
