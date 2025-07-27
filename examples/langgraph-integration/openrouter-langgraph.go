package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/tmc/langgraphgo/graph"
    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/schema"
    openrouter "github.com/revrost/go-openrouter"
)

// OpenRouterLLM adapts OpenRouter client to LangChain LLM interface
type OpenRouterLLM struct {
    client *openrouter.Client
    model  string
}

// NewOpenRouterLLM creates a new OpenRouter LLM adapter
func NewOpenRouterLLM(apiKey string, model string) *OpenRouterLLM {
    client := openrouter.NewClient(
        apiKey,
        openrouter.WithXTitle("Flow Test Go"),
        openrouter.WithHTTPReferer("https://github.com/peterovchinnikov/flow-test-go"),
    )
    
    return &OpenRouterLLM{
        client: client,
        model:  model,
    }
}

// GenerateContent implements the LangChain LLM interface
func (llm *OpenRouterLLM) GenerateContent(
    ctx context.Context,
    messages []llms.MessageContent,
    options ...llms.CallOption,
) (*llms.ContentResponse, error) {
    // Convert LangChain messages to OpenRouter format
    orMessages := make([]openrouter.ChatCompletionMessage, len(messages))
    for i, msg := range messages {
        content := ""
        for _, part := range msg.Parts {
            if text, ok := part.(llms.TextContent); ok {
                content += text.Text
            }
        }
        
        orMessages[i] = openrouter.ChatCompletionMessage{
            Role:    string(msg.Role),
            Content: openrouter.Content{Text: content},
        }
    }
    
    // Create request
    req := openrouter.ChatCompletionRequest{
        Model:    llm.model,
        Messages: orMessages,
    }
    
    // Execute request
    resp, err := llm.client.CreateChatCompletion(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Convert response back to LangChain format
    choices := make([]*llms.ContentChoice, len(resp.Choices))
    for i, choice := range resp.Choices {
        choices[i] = &llms.ContentChoice{
            Content: llms.TextParts(
                schema.ChatMessageTypeAI,
                choice.Message.Content,
            ),
        }
    }
    
    return &llms.ContentResponse{
        Choices: choices,
    }, nil
}

// BuildFlowWithOpenRouter shows how to build a LangGraph flow using OpenRouter
func BuildFlowWithOpenRouter() (*graph.CompiledGraph, error) {
    // Create LangGraph
    g := graph.NewMessageGraph()
    
    // Create OpenRouter LLM
    llm := NewOpenRouterLLM(
        os.Getenv("OPENROUTER_API_KEY"),
        openrouter.GPT4Turbo, // or any model like openrouter.DeepseekV3
    )
    
    // Add a code analysis node
    g.AddNode("analyze", func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
        // Add analysis prompt
        messages := append(state, llms.TextParts(
            schema.ChatMessageTypeHuman,
            "Analyze this code for potential issues, focusing on security and performance.",
        ))
        
        // Use OpenRouter for analysis
        response, err := llm.GenerateContent(ctx, messages)
        if err != nil {
            return nil, err
        }
        
        // Return updated state
        return append(state, response.Choices[0].Content), nil
    })
    
    // Add a decision node
    g.AddNode("check-issues", func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
        // Get last message (analysis result)
        lastMsg := state[len(state)-1]
        
        // Simple check for critical issues
        content := ""
        for _, part := range lastMsg.Parts {
            if text, ok := part.(llms.TextContent); ok {
                content += text.Text
            }
        }
        
        if containsCriticalIssue(content) {
            return append(state, llms.TextParts(
                schema.ChatMessageTypeSystem,
                "CRITICAL_ISSUES_FOUND",
            )), nil
        }
        
        return append(state, llms.TextParts(
            schema.ChatMessageTypeSystem,
            "NO_CRITICAL_ISSUES",
        )), nil
    })
    
    // Add conditional routing
    g.AddConditionalEdge("check-issues", func(ctx context.Context, state []llms.MessageContent) string {
        lastMsg := state[len(state)-1]
        for _, part := range lastMsg.Parts {
            if text, ok := part.(llms.TextContent); ok {
                if text.Text == "CRITICAL_ISSUES_FOUND" {
                    return "create-urgent-issue"
                }
            }
        }
        return "approve"
    })
    
    // Add action nodes
    g.AddNode("create-urgent-issue", func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
        // Use a more powerful model for creating detailed issue
        urgentLLM := NewOpenRouterLLM(
            os.Getenv("OPENROUTER_API_KEY"),
            openrouter.Claude35Opus, // Use Claude for detailed writing
        )
        
        messages := append(state, llms.TextParts(
            schema.ChatMessageTypeHuman,
            "Create a detailed GitHub issue for the critical problems found. Include reproduction steps and suggested fixes.",
        ))
        
        response, err := urgentLLM.GenerateContent(ctx, messages)
        if err != nil {
            return nil, err
        }
        
        return append(state, response.Choices[0].Content), nil
    })
    
    g.AddNode("approve", func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
        return append(state, llms.TextParts(
            schema.ChatMessageTypeAI,
            "Code review completed. No critical issues found. ✅",
        )), nil
    })
    
    // Set up flow
    g.SetEntryPoint("analyze")
    g.AddEdge("analyze", "check-issues")
    g.AddEdge("create-urgent-issue", graph.END)
    g.AddEdge("approve", graph.END)
    
    return g.Compile()
}

// Example with dynamic model selection
func selectModelForTask(taskType string) string {
    modelMap := map[string]string{
        "code-analysis":  openrouter.DeepseekV3,        // Best for code
        "creative":       openrouter.Claude35Sonnet,    // Best for writing
        "quick-response": openrouter.GPT35Turbo,        // Fast and cheap
        "complex":        openrouter.GPT4O,             // Most capable
        "vision":         openrouter.GPT4Vision,        // For images
    }
    
    if model, ok := modelMap[taskType]; ok {
        return model
    }
    return openrouter.GPT4Turbo // Default
}

func containsCriticalIssue(content string) bool {
    // Simple check - in real implementation would be more sophisticated
    criticalKeywords := []string{
        "security vulnerability",
        "sql injection",
        "memory leak",
        "race condition",
        "buffer overflow",
    }
    
    for _, keyword := range criticalKeywords {
        if contains(content, keyword) {
            return true
        }
    }
    return false
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && s[0:len(substr)] == substr
}

func main() {
    // Build and run the flow
    flow, err := BuildFlowWithOpenRouter()
    if err != nil {
        panic(err)
    }
    
    // Execute with some code to analyze
    ctx := context.Background()
    result, err := flow.Invoke(ctx, []llms.MessageContent{
        llms.TextParts(
            schema.ChatMessageTypeHuman,
            `Please analyze this Go code:
            
            func getUserData(id string) (User, error) {
                query := "SELECT * FROM users WHERE id = '" + id + "'"
                return db.Query(query)
            }`,
        ),
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Analysis complete:", result)
}
