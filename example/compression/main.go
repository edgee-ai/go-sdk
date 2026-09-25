// Example: Token compression with Edgee Gateway SDK
//
// This example demonstrates how to:
//  1. Turn tool-result trimming on for a single request
//  2. Access compression metrics from the response
//
// Tool-result trimming shortens the output of tool calls (here a long `ls -la`
// listing) before it reaches the model. The per-request toggles
// (ToolResultTrimming, ToolSurfaceReduction, OutputBrevity) override the API key
// settings for this request only; leave one nil to keep the key's setting.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edgee-ai/go-sdk/edgee"
)

// lsOutput builds a long directory listing, the kind of tool output coding agents send back.
func lsOutput() string {
	lines := []string{"total 800"}
	for i := 0; i < 200; i++ {
		lines = append(lines, fmt.Sprintf("-rw-r--r--  1 user  staff  %d Jan  1 12:00 src/components/module_%03d.tsx", 1000+i, i))
	}
	return strings.Join(lines, "\n")
}

func main() {
	// Create client with API key from environment variable
	client, err := edgee.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	output := lsOutput()

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("Edgee Token Compression Example")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Example: Large tool result with tool-result trimming turned on")
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Tool output length: %d characters\n", len(output))
	fmt.Println()

	callID := "call_1"
	description := "Run a shell command and return its output."
	input := edgee.InputObject{
		Messages: []edgee.Message{
			{Role: "user", Content: "How many files are in src/components?"},
			{Role: "assistant", ToolCalls: []edgee.ToolCall{{
				ID:       callID,
				Type:     "function",
				Function: edgee.FunctionCall{Name: "Bash", Arguments: `{"command":"ls -la src/components"}`},
			}}},
			{Role: "tool", ToolCallID: &callID, Content: output},
		},
		Tools: []edgee.Tool{{
			Type: "function",
			Function: edgee.FunctionDefinition{
				Name:        "Bash",
				Description: &description,
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{"command": map[string]any{"type": "string"}},
					"required":   []string{"command"},
				},
			},
		}},
		ToolResultTrimming: edgee.BoolPtr(true),
	}

	response, err := client.Send("anthropic/claude-haiku-4-5", input)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Response: %s\n", response.Text())
	fmt.Println()

	// Display usage information
	if response.Usage != nil {
		fmt.Println("Token Usage:")
		fmt.Printf("  Prompt tokens:     %d\n", response.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", response.Usage.CompletionTokens)
		fmt.Printf("  Total tokens:      %d\n", response.Usage.TotalTokens)
		fmt.Println()
	}

	// Display compression information
	if response.Compression != nil {
		fmt.Println("Compression Metrics:")
		fmt.Printf("  Saved tokens:  %d\n", response.Compression.SavedTokens)
		fmt.Printf("  Reduction:     %.1f%%\n", response.Compression.Reduction)
		fmt.Printf("  Cost savings:  $%.3f\n", float64(response.Compression.CostSavings)/1000000)
		fmt.Printf("  Time:          %d ms\n", response.Compression.TimeMs)
		if response.Compression.Reduction > 0 {
			originalTokens := int(float64(response.Compression.SavedTokens) * 100 / response.Compression.Reduction)
			tokensAfter := originalTokens - response.Compression.SavedTokens
			fmt.Println()
			fmt.Println("  💡 Without compression, this request would have used")
			fmt.Printf("     %d input tokens.\n", originalTokens)
			fmt.Printf("     With compression, only %d tokens were processed!\n", tokensAfter)
		}
	} else {
		fmt.Println("No compression data available in response.")
		fmt.Println("Note: Compression data is only returned when trimming actually shortened")
		fmt.Println("      a tool result.")
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))

	os.Exit(0)
}
