// Package main provides the entry point for the flow-test-go CLI tool.
package main

import (
	"github.com/ondatra-ai/flow-test-go/cmd/commands"
)

func main() {
	// Create application state
	state := commands.NewGlobalState()

	commands.Execute(state)
}
