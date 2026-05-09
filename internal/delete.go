package internal

import "fmt"

func Delete() {

	cfg, err := loadConfig()
	if err != nil {
		fatal(err.Error())
	}

	if cfg.DevstackLabel == "" {
		fatal("not configured — run 'runslit config' first")
	}

	releases, err := selectReleases("Select releases to destroy")
	if err != nil {
		fatal(err.Error())
	}
	if len(releases) == 0 {
		fmt.Println("No releases selected.")
		return
	}

	label := cfg.DevstackLabel

	for _, rel := range releases {
		info(fmt.Sprintf("Uninstalling %s-%s", rel, label))
		ns := NBPlusNamespace
		if rel == releaseMockGW {
			ns = MockGWNamespace
		}
		err = runCommand(".", "helm", "uninstall", rel+"-"+label, "--namespace", ns)
		if err != nil {
			fatal("helm uninstall failed for " + rel)
		}
		success(rel + " uninstalled")
	}

	fmt.Println()
	success("Done")
}
