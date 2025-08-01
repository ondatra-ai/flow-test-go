package main

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain_Help(t *testing.T) {
	// Test main function by running it as a subprocess
	if os.Getenv("BE_SUBPROCESS") == "1" {
		os.Args = []string{"flow-test-go", "--help"}

		main()

		return
	}

	// Run the test as a subprocess
	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_Help")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")

	err := cmd.Run()

	// Help should exit cleanly
	require.NoError(t, err, "Help command should not error")
}

func TestMain_Version(t *testing.T) {
	// Test version flag
	if os.Getenv("BE_SUBPROCESS") == "1" {
		os.Args = []string{"flow-test-go", "--version"}

		main()

		return
	}

	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_Version")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")

	err := cmd.Run()

	// Version should exit cleanly
	require.NoError(t, err, "Version command should not error")
}

func TestMain_InvalidCommand(t *testing.T) {
	// Test invalid command handling
	if os.Getenv("BE_SUBPROCESS") == "1" {
		os.Args = []string{"flow-test-go", "invalid-command"}

		main()

		return
	}

	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_InvalidCommand")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")

	err := cmd.Run()

	// Invalid command MUST exit with non-zero code
	require.Error(t, err, "Invalid command should always fail")

	var exitError *exec.ExitError
	require.ErrorAs(t, err, &exitError, "Should be an exit error")
	assert.NotEqual(t, 0, exitError.ExitCode(), "Invalid command should exit with non-zero code")
}
