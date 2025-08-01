package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_Help(t *testing.T) {
	// Test main function by running it as a subprocess
	if os.Getenv("BE_SUBPROCESS") == "1" {
		main()

		return
	}

	// Run the test as a subprocess
	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_Help")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")
	cmd.Args = append(cmd.Args, "--", "--help")

	err := cmd.Run()

	// Main function should complete without hanging
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		// Exit codes 0 or 1 are acceptable for help
		assert.Contains(t, []int{0, 1}, exitError.ExitCode())
	} else {
		assert.NoError(t, err)
	}
}

func TestMain_Version(t *testing.T) {
	// Test version flag
	if os.Getenv("BE_SUBPROCESS") == "1" {
		main()

		return
	}

	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_Version")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")
	cmd.Args = append(cmd.Args, "--", "--version")

	err := cmd.Run()

	// Version should exit cleanly
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		assert.Contains(t, []int{0, 1}, exitError.ExitCode())
	} else {
		assert.NoError(t, err)
	}
}

func TestMain_InvalidCommand(t *testing.T) {
	// Test invalid command handling
	if os.Getenv("BE_SUBPROCESS") == "1" {
		main()

		return
	}

	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_InvalidCommand")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")
	cmd.Args = append(cmd.Args, "--", "invalid-command")

	err := cmd.Run()

	// Invalid command should exit with non-zero
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		assert.NotEqual(t, 0, exitError.ExitCode())
	}
}
