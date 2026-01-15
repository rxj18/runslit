package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Initiate() {
	kubePath, err := checkKubePath()
	if err != nil {
		fatal(err.Error())
	}

	printBanner()
	fmt.Printf("%s→ Using kube-manifests at: %s%s\n\n", Yellow, kubePath, Reset)

	reader := bufio.NewReader(os.Stdin)

	env, err := selectEnvironment(reader)
	if err != nil {
		fatal(err.Error())
	}

	fmt.Printf("%s→ Selected environment: %s%s\n\n", Yellow, env, Reset)

	label := prompt(reader, "Enter your DEVSTACK_LABEL (e.g., your-name, pr-123): ")
	if err := validateDevstackLabel(label); err != nil {
		fatal(err.Error())
	}

	helmfileDir := filepath.Join(kubePath, "helmfile")
	valuesFile := filepath.Join(
		helmfileDir,
		"charts",
		"payments-nbplus",
		"values.yaml",
	)

	slitHelmfile := getHelmfilePath(kubePath)

	info(fmt.Sprintf("Initializing %s environment...", strings.ToUpper(env)))

	if err := updateValuesFile(valuesFile, env); err != nil {
		fatal("failed to update values.yaml")
	}

	if err := createSlitHelmfile(slitHelmfile, env, label); err != nil {
		fatal("failed to create SLIT helmfile")
	}

	if err := saveSlitConfig(env, label); err != nil {
		fatal("failed to save SLIT config")
	}

	fmt.Println()
	success(fmt.Sprintf("%s environment initialized successfully!", strings.ToUpper(env)))
	fmt.Println()
	fmt.Printf("%sNext steps:%s\n", Blue, Reset)
	fmt.Println("  • Run 'runslit sync' to deploy your SLIT environment")
	fmt.Println("  • Run 'runslit delete' to destroy your SLIT environment")
	fmt.Println("  • Run 'runslit status' to check current configuration")
}

// helper func ...
func checkKubePath() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}

	return validateKubePath(cfg.KubeManifestsPath)
}

func validateKubePath(kubePath string) (string, error) {
	if kubePath == "" {
		return "", fmt.Errorf("kube-manifests path not configured")
	}

	fileInfo, err := os.Stat(kubePath)
	if err != nil || !fileInfo.IsDir() {
		return "", fmt.Errorf("kube-manifests path does not exist: %s", kubePath)
	}

	return kubePath, nil
}

func getHelmfilePath(kubePath string) string {
	return filepath.Join(kubePath, "helmfile", SlitHelmfileName)
}

func prompt(reader *bufio.Reader, msg string) string {
	fmt.Print(msg)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func selectEnvironment(reader *bufio.Reader) (string, error) {
	fmt.Printf("%sSelect environment:%s\n", Blue, Reset)
	fmt.Println("  1) stage")
	fmt.Println("  2) slit")
	fmt.Println()

	choice := prompt(reader, "Enter choice [1-2]: ")

	switch choice {
	case "1":
		return "stage", nil
	case "2":
		return "slit", nil
	default:
		return "", fmt.Errorf("invalid choice")
	}
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

func updateValuesFile(valuesFile, env string) error {
	data, err := os.ReadFile(valuesFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	foundLive := false
	foundTest := false

	for i, line := range lines {
		if strings.HasPrefix(line, "payments_nbplus_live_app_env:") {
			lines[i] = "payments_nbplus_live_app_env: " + env
			foundLive = true
		}
		if strings.HasPrefix(line, "payments_nbplus_test_app_env:") {
			lines[i] = "payments_nbplus_test_app_env: " + env
			foundTest = true
		}
		// Early exit once both keys are found
		if foundLive && foundTest {
			break
		}
	}

	return os.WriteFile(valuesFile, []byte(strings.Join(lines, "\n")), 0644)
}

func createSlitHelmfile(path, env, label string) error {
	content := fmt.Sprintf(`# Auto-generated SLIT Helmfile
# Environment: %s
# Devstack Label: %s
# Generated on: %s

helmDefaults:
  cleanupOnFail: false
  wait: true
  recreatePods: true
  createNamespace: false
  timeout: 200
  historyMax: 1

environments:
  default:
    values:
      - devstack_label: %s
      - ttl: 12h
      - secret: {{ randAlphaNum 12 | lower }}
---

releases:
  - name: payments-nbplus-{{ .Values.devstack_label }}
    namespace: payments-nbplus
    chart: ./charts/payments-nbplus
    values:
      - image: ce29e52abe4de1ffc1108bd70f69e8244d70ffd1
      - devstack_label: {{ .Values.devstack_label }}
      - ttl: {{ .Values.ttl }}
      - secret: '{{ .Values.secret }}'
      - create_pg_ledger_acknowledgment_worker: false
      - create_outbox_relay_worker: false
      - create_sqs_recon_worker: false
      - ephemeral_db: false
      - db_password: {{ randAlphaNum 12 | lower }}

  - name: mock-go-{{ .Values.devstack_label }}
    namespace: perf
    chart: ./charts/mock-gateway
    values:
      - image: 0679b428295e966ec7a3ce661577b034eda701f1
      - devstack_label: {{ .Values.devstack_label }}
      - ttl: {{ .Values.ttl }}
`,
		env,
		label,
		time.Now().Format(time.RFC1123),
		label,
	)

	return os.WriteFile(path, []byte(content), 0644)
}

func saveSlitConfig(env, label string) error {
	// Load existing config to preserve other fields
	cfg, err := loadConfig()
	if err != nil {
		cfg = &Config{}
	}

	// Update only the SLIT-specific fields
	cfg.SlitEnv = env
	cfg.DevstackLabel = label

	return cfg.saveConfigToFile()
}
