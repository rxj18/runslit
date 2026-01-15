package internal

import "fmt"

func ShowHelp() {
	printBanner()

	fmt.Println("Usage: runslit <command>")
	fmt.Println()

	fmt.Printf("%sSetup Commands:%s\n", Blue, Reset)
	fmt.Println("  config     Set/update kube-manifests path")
	fmt.Println()

	fmt.Printf("%sSLIT Commands:%s\n", Blue, Reset)
	fmt.Println("  init       Initialize SLIT environment")
	fmt.Println("  sync       Deploy/sync the SLIT helmfile")
	fmt.Println("  delete     Destroy the SLIT deployment")
	fmt.Println("  status     Show current configuration")
	fmt.Println("  test       Select and run tests from ./slit directory")
	fmt.Println()

	fmt.Printf("%sOther:%s\n", Blue, Reset)
	fmt.Println("  help       Show this help message")
	fmt.Println()

	fmt.Printf("%sExamples:%s\n", Blue, Reset)
	fmt.Println("  runslit config     # Set kube-manifests path")
	fmt.Println("  runslit init       # Initialize SLIT environment")
	fmt.Println("  runslit sync       # Deploy your environment")
	fmt.Println("  runslit test       # Select and run a test")
	fmt.Println("  runslit delete     # Destroy your environment")
}
