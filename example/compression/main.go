package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edgee-ai/go-sdk/edgee"
)

// Large context document to demonstrate input compression
// IMPORTANT: Only USER messages are compressed. System messages are not compressed.
const LARGE_CONTEXT = `
The History and Impact of Artificial Intelligence

Artificial intelligence (AI) has evolved from a theoretical concept to a 
transformative technology that influences nearly every aspect of modern life. 
The field began in earnest in the 1950s when pioneers like Alan Turing and 
John McCarthy laid the groundwork for machine intelligence.

Early developments focused on symbolic reasoning and expert systems. These 
rule-based approaches dominated the field through the 1970s and 1980s, with 
systems like MYCIN demonstrating practical applications in medical diagnosis. 
However, these early systems were limited by their inability to learn from data 
and adapt to new situations.

The resurgence of neural networks in the 1980s and 1990s, particularly with 
backpropagation algorithms, opened new possibilities. Yet it wasn't until the 
2010s, with the advent of deep learning and the availability of massive datasets 
and computational power, that AI truly began to revolutionize industries.

Modern AI applications span numerous domains:
- Natural language processing enables machines to understand and generate human language
- Computer vision allows machines to interpret visual information from the world
- Robotics combines AI with mechanical systems for autonomous operation
- Healthcare uses AI for diagnosis, drug discovery, and personalized treatment
- Finance leverages AI for fraud detection, algorithmic trading, and risk assessment
- Transportation is being transformed by autonomous vehicles and traffic optimization

The development of large language models like GPT, BERT, and others has 
particularly accelerated progress in natural language understanding and generation. 
These models, trained on vast amounts of text data, can perform a wide range of 
language tasks with remarkable proficiency.

Despite remarkable progress, significant challenges remain. Issues of bias, 
interpretability, safety, and ethical considerations continue to be areas of 
active research and debate. The AI community is working to ensure that these 
powerful technologies are developed and deployed responsibly, with consideration 
for their societal impact.

Looking forward, AI is expected to continue advancing rapidly, with potential 
breakthroughs in areas like artificial general intelligence, quantum machine 
learning, and brain-computer interfaces. The integration of AI into daily life 
will likely deepen, raising important questions about human-AI collaboration, 
workforce transformation, and the future of human cognition itself.
`

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

	// Example: Request with compression enabled and large input
	fmt.Println("Example: Large user message with compression enabled")
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Input context length: %d characters\n", len(LARGE_CONTEXT))
	fmt.Println()

	// NOTE: Only USER messages are compressed
	// Put the large context in the user message to demonstrate compression
	userMessage := fmt.Sprintf(`Here is some context about AI:

%s

Based on this context, summarize the key milestones in AI development in 3 bullet points.`, LARGE_CONTEXT)

	// Create input object with compression settings
	compressionRate := 0.5
	semanticThreshold := 60
	input := edgee.InputObject{
		Messages: []edgee.Message{
			{Role: "user", Content: userMessage},
		},
		CompressionModel: "agentic",
		CompressionConfiguration: &edgee.CompressionConfiguration{
			Rate:                          &compressionRate,
			SemanticPreservationThreshold: &semanticThreshold,
		},
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
		fmt.Println("Note: Compression data is only returned when compression is enabled")
		fmt.Println("      and supported by your API key configuration.")
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))

	os.Exit(0)
}
