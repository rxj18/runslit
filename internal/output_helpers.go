package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func info(msg string) {
	fmt.Fprintln(os.Stdout, Yellow+"→ "+msg+Reset)
}

func success(msg string) {
	fmt.Fprintln(os.Stdout, Green+"✓ "+msg+Reset)
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, Red+"✗ "+msg+Reset)
	os.Exit(1)
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runCommandWithEnv(dir string, env []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), env...)
	return cmd.Run()
}
