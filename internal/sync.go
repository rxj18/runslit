package internal

import "fmt"

const (
	releaseNBPlus = "payments-nbplus"
	releaseMockGW = "mock-go"
)

func Sync() {

	kubePath, err := checkKubePath()
	if err != nil {
		fatal(err.Error())
	}

	cfg, err := loadConfig()
	if err != nil {
		fatal(err.Error())
	}

	if cfg.DevstackLabel == "" {
		fatal("not configured — run 'runslit config' first")
	}

	releases, err := selectReleases("Select releases to deploy")
	if err != nil {
		fatal(err.Error())
	}
	if len(releases) == 0 {
		fmt.Println("No releases selected.")
		return
	}

	label := cfg.DevstackLabel
	ttl := cfg.ttl()

	for _, rel := range releases {
		switch rel {
		case releaseNBPlus:
			img := cfg.NBPlusImage
			if img == "" {
				fatal("payments-nbplus image not set — run 'runslit config'")
			}
			info(fmt.Sprintf("Deploying %s-%s", rel, label))
			err = runCommand(".",
				"helm", "upgrade", "--install",
				rel+"-"+label,
				chartPath(kubePath, NBPlusChartPath),
				"--namespace", NBPlusNamespace,
				"--set", "devstack_label="+label,
				"--set", "payments_nbplus_live_app_env=slit",
				"--set", "payments_nbplus_test_app_env=slit",
				"--set", "image="+img,
				"--set", "ttl="+ttl,
				"--set", "create_pg_ledger_acknowledgment_worker=false",
				"--set", "create_outbox_relay_worker=false",
				"--set", "create_sqs_recon_worker=false",
				"--set", "ephemeral_db=false",
				"--wait",
				"--timeout", "200s",
			)
			if err != nil {
				fatal("helm upgrade failed for " + rel)
			}
			success(rel + " deployed")

		case releaseMockGW:
			img := cfg.MockGWImage
			if img == "" {
				fatal("mock-go image not set — run 'runslit config'")
			}
			info(fmt.Sprintf("Deploying %s-%s", rel, label))
			err = runCommand(".",
				"helm", "upgrade", "--install",
				rel+"-"+label,
				chartPath(kubePath, MockGWChartPath),
				"--namespace", MockGWNamespace,
				"--set", "devstack_label="+label,
				"--set", "image="+img,
				"--set", "ttl="+ttl,
				"--wait",
				"--timeout", "200s",
			)
			if err != nil {
				fatal("helm upgrade failed for " + rel)
			}
			success(rel + " deployed")
		}
	}

	fmt.Println()
	success("Done")
}
