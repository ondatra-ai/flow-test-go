package commands_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peterovchinnikov/flow-test-go/cmd/commands"
)

func TestNewGlobalState(t *testing.T) {
	state := commands.NewGlobalState()
	assert.NotNil(t, state)
}

func TestGlobalState_Structure(t *testing.T) {
	state1 := commands.NewGlobalState()
	state2 := commands.NewGlobalState()

	// Each state should be independent
	assert.NotSame(t, state1, state2)
}

func TestCommandIntegration(t *testing.T) {
	state := commands.NewGlobalState()

	// Test that we can create commands with the state
	cmd := commands.CreateListCommand(state)
	require.NotNil(t, cmd)

	// Test integration with state
	assert.Equal(t, "list", cmd.Use)
	assert.Contains(t, cmd.Short, "List")
}

func TestListCommand_Help(t *testing.T) {
	state := commands.NewGlobalState()
	cmd := commands.CreateListCommand(state)

	// Test help output
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	// Set help flag and execute
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()

	// Help should not error and should contain expected content
	require.NoError(t, err)

	helpOutput := output.String()
	assert.Contains(t, helpOutput, "list")
	assert.Contains(t, helpOutput, "List")
}

func TestListCommand_Structure(t *testing.T) {
	state := commands.NewGlobalState()
	cmd := commands.CreateListCommand(state)

	// Test command structure
	assert.Equal(t, "list", cmd.Use)
	assert.False(t, cmd.Hidden)
	assert.NotEmpty(t, cmd.Short)
	assert.False(t, cmd.HasSubCommands())
}

func TestListCommand_Flags(t *testing.T) {
	state := commands.NewGlobalState()
	cmd := commands.CreateListCommand(state)

	// Test that command has flag functionality
	flags := cmd.Flags()
	assert.NotNil(t, flags)

	// Test basic command functionality
	assert.NotEmpty(t, cmd.Use)
	assert.NotEmpty(t, cmd.Short)
}
