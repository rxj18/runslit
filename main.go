package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rxj18/runslit/internal"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		internal.ShowHelp()
		return
	}

	cmd := strings.ToLower(strings.TrimSpace(args[0]))

	switch cmd {
	case "config":
		internal.Configure()
	case "status":
		internal.ShowStatus()
	case "sync":
		internal.Sync()
	case "delete", "destroy":
		internal.Delete()
	case "test":
		internal.RunTest()
	case "help":
		internal.ShowHelp()
	default:
		fmt.Println("Unknown command:", args[0])
		internal.ShowHelp()
		os.Exit(1)
	}
}
