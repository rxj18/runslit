package internal

import (
	"fmt"
	"os"
)

var (
	out = os.Stdout
	err = os.Stderr
)

func info(msg string) {
	fmt.Fprintln(out, Yellow+"→ "+msg+Reset)
}

func success(msg string) {
	fmt.Fprintln(out, Green+"✓ "+msg+Reset)
}

func fatal(msg string) {
	fmt.Fprintln(err, Red+"✗ "+msg+Reset)
	os.Exit(1)
}
