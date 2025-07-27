// Package main provides the entry point for the flow-test-go CLI tool.
package main

import (
	"github.com/peterovchinnikov/flow-test-go/cmd/flow-test-go/commands"
)

func main() {
	commands.Execute()
}
