package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edgee-cloud/go-sdk/edgee"
)

func main() {
	// Create client with API key from environment variable
	client, err := edgee.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("Edgee Token Compression Example")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Example: Request with compression enabled
	fmt.Println("Example: Request with compression enabled")
	fmt.Println(strings.Repeat("-", 70))

	// Create input object with compression settings
	input := edgee.InputObject{
		Messages: []edgee.Message{
			{Role: "user", Content: "Explain quantum computing in simple terms."},
		},
	}

	// Set compression parameters
	enableCompression := true
	compressionRate := 0.5
	input.EnableCompression = &enableCompression
	input.CompressionRate = &compressionRate

	response, err := client.Send("gpt-4o", input)
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
		fmt.Printf("  Input tokens:  %d\n", response.Compression.InputTokens)
		fmt.Printf("  Saved tokens:  %d\n", response.Compression.SavedTokens)
		fmt.Printf("  Compression rate: %.2f%%\n", response.Compression.Rate*100)
		fmt.Printf("  Token savings: %d tokens saved!\n", response.Compression.SavedTokens)
	} else {
		fmt.Println("No compression data available in response.")
		fmt.Println("Note: Compression data is only returned when compression is enabled")
		fmt.Println("      and supported by your API key configuration.")
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))

	os.Exit(0)
}
