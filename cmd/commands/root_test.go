package commands_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peterovchinnikov/flow-test-go/cmd/commands"
)

func TestNewGlobalState(t *testing.T) {
	state := commands.NewGlobalState()
	assert.NotNil(t, state)
}

func TestGlobalState_Structure(t *testing.T) {
	// Test the global state structure and basic operations
	state := commands.NewGlobalState()
	assert.NotNil(t, state)

	// Test that we can create multiple instances without issues
	state2 := commands.NewGlobalState()
	assert.NotNil(t, state2)

	// They should be different instances
	assert.NotSame(t, state, state2)
}

func TestExecute_Help(t *testing.T) {
	// Test the Execute function by running it as a subprocess
	// This is the recommended way to test commands that call os.Exit
	if os.Getenv("BE_SUBPROCESS") == "1" {
		state := commands.NewGlobalState()
		// Simulate help command
		os.Args = []string{"flow-test-go", "--help"}

		commands.Execute(state)

		return
	}

	// Run the test as a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_Help")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Help command should exit with code 0
	if exitError, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, 0, exitError.ExitCode())
	} else {
		assert.NoError(t, err)
	}

	// Should contain help text
	output := stdout.String() + stderr.String()
	assert.Contains(t, output, "flow-test-go")
}

func TestExecute_InvalidCommand(t *testing.T) {
	// Test error handling by running invalid command as subprocess
	if os.Getenv("BE_SUBPROCESS") == "1" {
		state := commands.NewGlobalState()
		// Simulate invalid command
		os.Args = []string{"flow-test-go", "invalid-command"}

		commands.Execute(state)

		return
	}

	// Run the test as a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_InvalidCommand")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Invalid command should exit with non-zero code
	if exitError, ok := err.(*exec.ExitError); ok {
		assert.NotEqual(t, 0, exitError.ExitCode())
	} else {
		// If no error, check that stderr contains error message
		assert.Contains(t, stderr.String(), "Error")
	}
}

func TestExecute_Version(t *testing.T) {
	// Test version command
	if os.Getenv("BE_SUBPROCESS") == "1" {
		state := commands.NewGlobalState()
		os.Args = []string{"flow-test-go", "--version"}

		commands.Execute(state)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_Version")
	cmd.Env = append(os.Environ(), "BE_SUBPROCESS=1")

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Version command should exit with code 0
	if exitError, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, 0, exitError.ExitCode())
	} else {
		assert.NoError(t, err)
	}

	// Should contain version
	output := stdout.String() + stderr.String()
	assert.True(t, strings.Contains(output, "1.0.0") || strings.Contains(output, "version"))
}
