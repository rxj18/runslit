package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunTest() {
	printBanner()

	// Check if ./slit directory exists
	slitDir := "/Users/rituraj.paul/razorpay/nbplus/slit"
	if _, err := os.Stat(slitDir); os.IsNotExist(err) {
		fatal("./slit directory not found in current directory")
	}

	// Get absolute path of slit directory
	absSlitDir, err := filepath.Abs(slitDir)
	if err != nil {
		fatal(fmt.Sprintf("failed to get absolute path of slit directory: %s", err.Error()))
	}

	info("Scanning for tests in ./slit directory...")

	// Extract all test suites and test cases
	testCases, err := extractTestCases(slitDir)
	if err != nil {
		fatal(fmt.Sprintf("failed to extract tests: %s", err.Error()))
	}

	if len(testCases) == 0 {
		fatal("no tests found in ./slit directory")
	}

	fmt.Printf("%sFound %d test case(s)%s\n\n", Green, len(testCases), Reset)

	// Let user select test using fzf
	selectedTest, err := selectTestWithFzf(testCases)
	if err != nil {
		fatal(fmt.Sprintf("test selection failed: %s", err.Error()))
	}

	if selectedTest == "" {
		fmt.Println("No test selected")
		return
	}

	// Parse the selection to build test run command
	testRunner, testMethod, relPath := parseSelectedTest(selectedTest)

	// Get the module root directory (parent of slit directory)
	moduleRoot := filepath.Dir(absSlitDir)

	// Build test path relative to module root
	// Must start with ./ to indicate filesystem path, not package import path
	testPath := "./slit"
	if relPath != "" {
		testPath = "./slit/" + relPath
	}

	// Load config to get DEVSTACK_LABEL
	cfg, err := loadConfig()
	if err != nil {
		fatal(fmt.Sprintf("failed to load config: %s", err.Error()))
	}

	if cfg.DevstackLabel == "" {
		fatal("DEVSTACK_LABEL not configured. Run 'runslit init' first.")
	}

	// Print what we're running
	fmt.Println()
	info(fmt.Sprintf("Running: %s/%s in %s", testRunner, testMethod, testPath))
	info(fmt.Sprintf("DEVSTACK_LABEL: %s", cfg.DevstackLabel))
	info(fmt.Sprintf("Working directory: %s", moduleRoot))
	fmt.Println()

	// Run the test with DEVSTACK_LABEL environment variable from module root
	err = runCommandWithEnv(moduleRoot, []string{fmt.Sprintf("DEVSTACK_LABEL=%s", cfg.DevstackLabel)}, "go", "test", "-v", "-run", fmt.Sprintf("%s/%s", testRunner, testMethod), testPath)
	if err != nil {
		fmt.Println()
		fatal("test execution failed")
	}

	fmt.Println()
	success("Test completed successfully!")
}

// extractTestCases scans the slit directory and extracts test suites and their test cases using AST
func extractTestCases(slitDir string) ([]string, error) {
	var testCases []string

	// Group test files by directory
	dirFiles := make(map[string][]string)

	err := filepath.Walk(slitDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only collect Go test files
		if info.IsDir() || !strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		dir := filepath.Dir(path)
		dirFiles[dir] = append(dirFiles[dir], path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process each directory
	for dir, files := range dirFiles {
		fset := token.NewFileSet()

		// Parse all files in the directory
		var parsedFiles []*ast.File
		for _, file := range files {
			parsed, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				continue
			}
			parsedFiles = append(parsedFiles, parsed)
		}

		// Find suite info from any file in the directory
		var suiteInfo *suiteInfo
		for _, file := range parsedFiles {
			suiteInfo = findTestSuiteInfo(file)
			if suiteInfo != nil {
				break
			}
		}

		if suiteInfo == nil {
			continue
		}

		// Extract test methods from ALL files in the directory
		var testMethods []string
		for _, file := range parsedFiles {
			methods := extractSuiteTestMethods(file, suiteInfo.SuiteTypeName)
			testMethods = append(testMethods, methods...)
		}

		// Get relative path for test command
		relPath, _ := filepath.Rel(slitDir, dir)
		if relPath == "." {
			relPath = ""
		}

		// Build test case entries
		for _, method := range testMethods {
			// Store just the relative path from slit directory (e.g., "netbanking", "fpx")
			// Format: TestRunner -> TestMethod | relPath
			testCases = append(testCases, fmt.Sprintf("%s -> %s | %s", suiteInfo.TestRunnerName, method, relPath))
		}
	}

	return testCases, nil
}

type suiteInfo struct {
	TestRunnerName string // e.g., TestFpxPayment
	SuiteTypeName  string // e.g., FpxPaymentSlitTestSuite
}

// findTestSuiteInfo finds the main test function that runs the suite
func findTestSuiteInfo(file *ast.File) *suiteInfo {
	var info *suiteInfo

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(funcDecl.Name.Name, "Test") {
			return true
		}

		// Look for suite.Run call
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if it's suite.Run
			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok || selExpr.Sel.Name != "Run" {
				return true
			}

			// Get the suite type from the second argument
			if len(callExpr.Args) >= 2 {
				if unaryExpr, ok := callExpr.Args[1].(*ast.UnaryExpr); ok {
					if compLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
						if ident, ok := compLit.Type.(*ast.Ident); ok {
							info = &suiteInfo{
								TestRunnerName: funcDecl.Name.Name,
								SuiteTypeName:  ident.Name,
							}
						}
					}
				}
			}

			return false
		})

		if info != nil {
			return false
		}

		return true
	})

	return info
}

// extractSuiteTestMethods extracts all test methods for a given suite type
func extractSuiteTestMethods(file *ast.File, suiteTypeName string) []string {
	var testMethods []string

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if it's a method on the suite
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			return true
		}

		// Get the receiver type
		var recvType string
		switch t := funcDecl.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				recvType = ident.Name
			}
		case *ast.Ident:
			recvType = t.Name
		}

		// Check if it's a method on our suite and starts with Test
		if recvType == suiteTypeName && strings.HasPrefix(funcDecl.Name.Name, "Test") {
			testMethods = append(testMethods, funcDecl.Name.Name)
		}

		return true
	})

	return testMethods
}

// parseSelectedTest parses the selected test string into components
func parseSelectedTest(selectedTest string) (testRunner, testMethod, relPath string) {
	// Parse format: TestRunner -> TestMethod | relPath
	parts := strings.Split(selectedTest, " | ")
	if len(parts) != 2 {
		return selectedTest, "", ""
	}

	relPath = strings.TrimSpace(parts[1])
	testInfo := strings.TrimSpace(parts[0])

	// Parse: TestRunner -> TestMethod
	if strings.Contains(testInfo, " -> ") {
		infoParts := strings.Split(testInfo, " -> ")
		testRunner = strings.TrimSpace(infoParts[0])
		testMethod = strings.TrimSpace(infoParts[1])
		return testRunner, testMethod, relPath
	}

	// Simple test without hierarchy
	return testInfo, "", relPath
}

// runCommandWithEnv executes a command with custom environment variables
func runCommandWithEnv(dir string, env []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Inherit existing environment and add custom ones
	cmd.Env = append(os.Environ(), env...)

	return cmd.Run()
}

// selectTestWithFzf uses fzf for interactive test selection
func selectTestWithFzf(testCases []string) (string, error) {
	// Check if fzf is available
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", fmt.Errorf("fzf not found in PATH. Please install fzf: https://github.com/junegunn/fzf")
	}

	// Prepare fzf command with better options for hierarchical display
	cmd := exec.Command("fzf",
		"--height", "50%",
		"--border",
		"--prompt", "Select test: ",
		"--preview-window", "hidden",
		"--layout", "reverse",
	)

	// Feed test cases to fzf via stdin
	cmd.Stdin = strings.NewReader(strings.Join(testCases, "\n"))
	cmd.Stderr = os.Stderr

	// Capture fzf output
	output, err := cmd.Output()
	if err != nil {
		// User cancelled (Ctrl+C or ESC)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
