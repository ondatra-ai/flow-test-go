package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peterovchinnikov/flow-test-go/cmd/commands"
)

func TestMain_StateCreation(t *testing.T) {
	// Test that main creates proper state
	// We can't test main() directly due to os.Exit, but we can test the components
	state := commands.NewGlobalState()
	assert.NotNil(t, state)
}

func TestMain_ComponentIntegration(t *testing.T) {
	// Test that main can create the components it needs
	state := commands.NewGlobalState()

	// Test that we can create the list command that would be used by main
	cmd := commands.CreateListCommand(state)
	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
}

func TestMain_StateConsistency(t *testing.T) {
	// Test that multiple state instances work correctly
	state1 := commands.NewGlobalState()
	state2 := commands.NewGlobalState()

	// States should be independent
	assert.NotSame(t, state1, state2)

	// Both should be able to create commands
	cmd1 := commands.CreateListCommand(state1)
	cmd2 := commands.CreateListCommand(state2)

	assert.NotNil(t, cmd1)
	assert.NotNil(t, cmd2)
	assert.NotSame(t, cmd1, cmd2)
}
