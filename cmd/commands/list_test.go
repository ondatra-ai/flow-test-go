package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peterovchinnikov/flow-test-go/cmd/commands"
)

func TestCreateListCommand(t *testing.T) {
	state := commands.NewGlobalState()
	cmd := commands.CreateListCommand(state)

	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List available flows", cmd.Short)
	assert.Contains(t, cmd.Long, "List all available flows")
}

func TestListFlows_CommandStructure(t *testing.T) {
	state := commands.NewGlobalState()
	cmd := commands.CreateListCommand(state)

	// Test command structure
	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List available flows", cmd.Short)
	assert.Contains(t, cmd.Long, "List all available flows")
	assert.NotNil(t, cmd.RunE)
}

func TestListFlows_CommandProperties(t *testing.T) {
	state := commands.NewGlobalState()
	cmd := commands.CreateListCommand(state)

	// Test that the command has the expected properties
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List available flows", cmd.Short)
	assert.Contains(t, cmd.Long, "List all available flows")
	assert.Contains(t, cmd.Long, "flow-test-go list")
	assert.NotNil(t, cmd.RunE)
}
