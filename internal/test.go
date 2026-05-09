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

	slitDir := "./slit"
	if _, err := os.Stat(slitDir); os.IsNotExist(err) {
		fatal("./slit directory not found in current directory")
	}

	absSlitDir, err := filepath.Abs(slitDir)
	if err != nil {
		fatal(err.Error())
	}

	cfg, err := loadConfig()
	if err != nil {
		fatal(err.Error())
	}
	if cfg.DevstackLabel == "" {
		fatal("devstack label not configured — run 'runslit config'")
	}

	info("Scanning for tests...")

	testCases, err := extractTestCases(slitDir)
	if err != nil {
		fatal(fmt.Sprintf("failed to extract tests: %s", err.Error()))
	}
	if len(testCases) == 0 {
		fatal("no tests found in ./slit directory")
	}

	fmt.Printf("%sFound %d test case(s)%s\n\n", Green, len(testCases), Reset)

	selected, err := selectTestWithFzf(testCases, cfg.LastTest)
	if err != nil {
		fatal(fmt.Sprintf("test selection failed: %s", err.Error()))
	}
	if selected == "" {
		fmt.Println("No test selected.")
		return
	}

	testRunner, testMethod, relPath := parseSelectedTest(selected)

	moduleRoot := filepath.Dir(absSlitDir)
	testPath := "./slit"
	if relPath != "" {
		testPath = "./slit/" + relPath
	}

	fmt.Println()
	info(fmt.Sprintf("Running: %s/%s in %s", testRunner, testMethod, testPath))
	info(fmt.Sprintf("DEVSTACK_LABEL: %s", cfg.DevstackLabel))
	fmt.Println()

	// Save last-run test before executing so it persists even on failure.
	cfg.LastTest = selected
	_ = cfg.saveConfigToFile()

	err = runCommandWithEnv(
		moduleRoot,
		[]string{"DEVSTACK_LABEL=" + cfg.DevstackLabel},
		"go", "test", "-v", "-run", testRunner+"/"+testMethod, testPath,
	)
	if err != nil {
		fmt.Println()
		fatal("test execution failed")
	}

	fmt.Println()
	success("Test completed successfully!")
}

// extractTestCases walks ./slit, finds all testify suite runners in every
// directory, and returns one entry per (runner, method) pair.
func extractTestCases(slitDir string) ([]string, error) {
	// Group test files by directory.
	dirFiles := make(map[string][]string)
	err := filepath.Walk(slitDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			dir := filepath.Dir(path)
			dirFiles[dir] = append(dirFiles[dir], path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var testCases []string

	for dir, files := range dirFiles {
		fset := token.NewFileSet()
		var parsed []*ast.File
		for _, f := range files {
			af, err := parser.ParseFile(fset, f, nil, 0)
			if err != nil {
				continue
			}
			parsed = append(parsed, af)
		}

		// Collect ALL suite runners in this directory (not just the first).
		var suites []suiteInfo
		for _, af := range parsed {
			suites = append(suites, findSuiteInfos(af)...)
		}
		if len(suites) == 0 {
			continue
		}

		relPath, _ := filepath.Rel(slitDir, dir)
		if relPath == "." {
			relPath = ""
		}

		for _, si := range suites {
			var methods []string
			for _, af := range parsed {
				methods = append(methods, extractSuiteTestMethods(af, si.SuiteTypeName)...)
			}
			for _, m := range methods {
				testCases = append(testCases, fmt.Sprintf("%s -> %s | %s", si.TestRunnerName, m, relPath))
			}
		}
	}

	return testCases, nil
}

type suiteInfo struct {
	TestRunnerName string // e.g. TestFpxPayment
	SuiteTypeName  string // e.g. FpxPaymentSlitTestSuite
}

// findSuiteInfos returns ALL testify suite runners in a file.
// A runner is a top-level Test* function that calls suite.Run(t, &SomeType{}).
func findSuiteInfos(file *ast.File) []suiteInfo {
	var result []suiteInfo

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Run" || len(call.Args) < 2 {
				return true
			}
			unary, ok := call.Args[1].(*ast.UnaryExpr)
			if !ok {
				return true
			}
			comp, ok := unary.X.(*ast.CompositeLit)
			if !ok {
				return true
			}
			ident, ok := comp.Type.(*ast.Ident)
			if !ok {
				return true
			}
			result = append(result, suiteInfo{
				TestRunnerName: fn.Name.Name,
				SuiteTypeName:  ident.Name,
			})
			return false
		})
	}

	return result
}

// extractSuiteTestMethods returns all Test* methods on the named suite type.
func extractSuiteTestMethods(file *ast.File, suiteTypeName string) []string {
	var methods []string
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}
		if !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		var recv string
		switch t := fn.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			if id, ok := t.X.(*ast.Ident); ok {
				recv = id.Name
			}
		case *ast.Ident:
			recv = t.Name
		}
		if recv == suiteTypeName {
			methods = append(methods, fn.Name.Name)
		}
	}
	return methods
}

// parseSelectedTest splits "TestRunner -> TestMethod | relPath".
func parseSelectedTest(s string) (runner, method, relPath string) {
	parts := strings.SplitN(s, " | ", 2)
	if len(parts) != 2 {
		return s, "", ""
	}
	relPath = strings.TrimSpace(parts[1])
	info := strings.TrimSpace(parts[0])
	if before, after, ok := strings.Cut(info, " -> "); ok {
		return strings.TrimSpace(before), strings.TrimSpace(after), relPath
	}
	return info, "", relPath
}

// selectTestWithFzf opens fzf for test selection.
// If lastTest is set, it is moved to the top of the list so the fzf cursor
// starts there — no query text in the search box.
func selectTestWithFzf(testCases []string, lastTest string) (string, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", fmt.Errorf("fzf not found in PATH — please install fzf")
	}

	// Move last-run test to top so the cursor lands on it without --query.
	if lastTest != "" {
		filtered := make([]string, 0, len(testCases))
		for _, tc := range testCases {
			if tc != lastTest {
				filtered = append(filtered, tc)
			}
		}
		testCases = append([]string{lastTest}, filtered...)
	}

	cmd := exec.Command("fzf",
		"--height", "80%",
		"--border",
		"--prompt", "Search tests: ",
		"--layout", "reverse",
		"--no-info",
	)
	cmd.Stdin = strings.NewReader(strings.Join(testCases, "\n"))
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

