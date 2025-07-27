package main

// Example of how our JSON configuration would be converted to LangGraph

import (
    "context"
    "fmt"
    
    "github.com/tmc/langgraphgo/graph"
    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/schema"
)

// This shows how we would build a LangGraph from our JSON config
func BuildGraphFromConfig(config *FlowConfig) (*graph.CompiledGraph, error) {
    // Create a new message graph
    g := graph.NewMessageGraph()
    
    // Add nodes for each step in our config
    for stepID, step := range config.Steps {
        switch step.Type {
        case "prompt":
            g.AddNode(stepID, createPromptNode(step))
        case "condition":
            // For conditions, we add conditional edges instead of nodes
            for _, condition := range step.Conditions {
                g.AddConditionalEdge(
                    getPreviousStep(stepID, config),
                    createConditionFunc(condition.Expression),
                    condition.Next,
                )
            }
        case "end":
            g.AddNode(stepID, graph.END)
        }
    }
    
    // Add regular edges based on "next" field
    for stepID, step := range config.Steps {
        if step.Next != "" && step.Type != "condition" {
            g.AddEdge(stepID, step.Next)
        }
    }
    
    // Set entry point
    g.SetEntryPoint(getFirstStep(config))
    
    // Compile the graph
    return g.Compile()
}

// Create a prompt node that integrates with our MCP and AI providers
func createPromptNode(step Step) graph.NodeFunc {
    return func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
        // 1. Get MCP tools if specified
        tools := []Tool{}
        if step.MCPServer != "" {
            mcpTools, err := mcpManager.GetTools(step.MCPServer, step.Tools)
            if err != nil {
                return nil, fmt.Errorf("failed to get MCP tools: %w", err)
            }
            tools = append(tools, mcpTools...)
        }
        
        // 2. Build prompt with variable substitution
        prompt := substituteVariables(step.Prompt.Template, extractVariables(state))
        
        // 3. Get AI provider (or use default)
        provider := aiManager.GetProvider(step.Model)
        
        // 4. Execute prompt with tools
        response, err := provider.GenerateContent(ctx, 
            append(state, llms.TextParts(schema.ChatMessageTypeHuman, prompt)),
            llms.WithTools(convertTools(tools)),
        )
        if err != nil {
            return nil, err
        }
        
        // 5. Append response to state
        return append(state, 
            llms.TextParts(schema.ChatMessageTypeAI, response.Choices[0].Content),
        ), nil
    }
}

// Create a condition function for conditional edges
func createConditionFunc(expression string) graph.ConditionalFunc {
    return func(ctx context.Context, state []llms.MessageContent) string {
        // Evaluate the expression against the last message
        result := evaluateExpression(expression, state)
        if result {
            return "true"
        }
        return "false"
    }
}

// Example usage
func main() {
    // Load flow configuration
    config, err := LoadFlowConfig("examples/langgraph-integration/example-flow.json")
    if err != nil {
        panic(err)
    }
    
    // Build LangGraph from config
    compiled, err := BuildGraphFromConfig(config)
    if err != nil {
        panic(err)
    }
    
    // Execute the flow
    ctx := context.Background()
    result, err := compiled.Invoke(ctx, []llms.MessageContent{
        llms.TextParts(schema.ChatMessageTypeHuman, "Review the latest pull request"),
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Flow completed:", result)
}
