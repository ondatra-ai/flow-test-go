package main

import (
    "fmt"
    "log"
    
    "github.com/spf13/cobra"
    "github.com/tmc/langgraphgo"
    "github.com/revrost/go-openrouter"
)

var rootCmd = &cobra.Command{
    Use:   "flow-test-go",
    Short: "A CLI tool for orchestrating AI agents",
    Long:  `flow-test-go is a CLI tool that orchestrates AI agents using LangGraph and OpenRouter.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("✅ flow-test-go CLI is working!")
        fmt.Println("   - Cobra: imported successfully")
        fmt.Println("   - LangGraph: imported successfully")
        fmt.Println("   - OpenRouter: imported successfully")
        
        // Basic type check to ensure imports work
        _ = &langgraphgo.Graph{}
        _ = &openrouter.Client{}
    },
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}
