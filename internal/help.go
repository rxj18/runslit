package internal

import "fmt"

func ShowHelp() {
	fmt.Println("Usage: runslit <command>")
	fmt.Println()

	fmt.Printf("%sCommands:%s\n", Blue, Reset)
	fmt.Println("  config    Configure runslit (path, label, images)")
	fmt.Println("  sync      Deploy selected releases")
	fmt.Println("  delete    Destroy selected releases")
	fmt.Println("  status    Show current configuration")
	fmt.Println("  test      Select and run a test")
	fmt.Println("  help      Show this help message")
	fmt.Println()

	fmt.Printf("%sExamples:%s\n", Blue, Reset)
	fmt.Println("  runslit config    # Set up or update any field")
	fmt.Println("  runslit sync      # Deploy releases")
	fmt.Println("  runslit delete    # Destroy releases")
	fmt.Println("  runslit test      # Pick and run a test")
}
