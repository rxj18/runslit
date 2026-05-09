package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

type menuItem struct {
	label  string
	getter func(*Config) string
	setter func(*Config, string) error
}

var configMenu = []menuItem{
	{
		label:  "kube-manifests path",
		getter: func(c *Config) string { return c.KubeManifestsPath },
		setter: setKubePath,
	},
	{
		label:  "devstack label",
		getter: func(c *Config) string { return c.DevstackLabel },
		setter: func(c *Config, v string) error {
			if err := validateDevstackLabel(v); err != nil {
				return err
			}
			c.DevstackLabel = v
			return nil
		},
	},
	{
		label:  "payments-nbplus image SHA",
		getter: func(c *Config) string { return c.NBPlusImage },
		setter: func(c *Config, v string) error { c.NBPlusImage = v; return nil },
	},
	{
		label:  "mock-go image SHA",
		getter: func(c *Config) string { return c.MockGWImage },
		setter: func(c *Config, v string) error { c.MockGWImage = v; return nil },
	},
	{
		label:  "TTL (e.g. 12h, 24h)",
		getter: func(c *Config) string { return c.ttl() },
		setter: func(c *Config, v string) error { c.TTL = v; return nil },
	},
}

func Configure() {

	for {
		cfg, err := loadConfig()
		if err != nil {
			fatal("failed to load config")
		}

		// Build select options showing current values inline.
		opts := make([]huh.Option[int], len(configMenu))
		for i, item := range configMenu {
			val := item.getter(cfg)
			if val != "" {
				opts[i] = huh.NewOption(fmt.Sprintf("%-28s %s", item.label, val), i)
			} else {
				opts[i] = huh.NewOption(fmt.Sprintf("%-28s (not set)", item.label), i)
			}
		}
		choice := 0
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Title("Configure runslit").
					Description("Select a field to update").
					Options(opts...).
					Value(&choice),
			),
		).WithKeyMap(huhKeyMap()).Run()

		if err == huh.ErrUserAborted {
			fmt.Println()
			return
		}
		if err != nil {
			fatal(err.Error())
		}

		item := configMenu[choice]
		current := item.getter(cfg)
		newVal := current

		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(item.label).
					Placeholder(current).
					Value(&newVal),
			),
		).WithKeyMap(huhKeyMap()).Run()

		if err == huh.ErrUserAborted || newVal == current {
			continue
		}
		if err != nil {
			fatal(err.Error())
		}

		if err := item.setter(cfg, newVal); err != nil {
			fmt.Printf("%s✗ %s%s\n\n", Red, err.Error(), Reset)
			continue
		}

		if err := cfg.saveConfigToFile(); err != nil {
			fmt.Printf("%s✗ failed to save%s\n\n", Red, Reset)
			continue
		}

		success("saved")
		fmt.Println()
	}
}

func setKubePath(cfg *Config, val string) error {
	abs, err := expandPath(val)
	if err != nil {
		return fmt.Errorf("invalid path")
	}
	fi, err := os.Stat(abs)
	if err != nil || !fi.IsDir() {
		return fmt.Errorf("directory does not exist: %s", abs)
	}
	if _, err := os.Stat(filepath.Join(abs, "helmfile")); err != nil {
		fmt.Printf("%s↳ no helmfile/charts directory found, continuing anyway%s\n", Yellow, Reset)
	}
	cfg.KubeManifestsPath = abs
	return nil
}

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
