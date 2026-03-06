# Edgee Go SDK

Lightweight, type-safe Go SDK for the [Edgee AI Gateway](https://www.edgee.ai).

[![Go Reference](https://pkg.go.dev/badge/github.com/edgee-ai/go-sdk.svg)](https://pkg.go.dev/github.com/edgee-ai/go-sdk)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## Installation

```bash
go get github.com/edgee-ai/go-sdk/edgee
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/edgee-ai/go-sdk/edgee"
)

func main() {
    client, err := edgee.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

<<<<<<< HEAD
    response, err := client.Send("gpt-5.2", "What is the capital of France?")
=======
    response, err := client.Send("anthropic/claude-haiku-4-5", "What is the capital of France?")
>>>>>>> d846ba2 (feat: update compression response to new API format)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Text())
    // "The capital of France is Paris."
}
```

## Send Method

The `Send()` method makes non-streaming chat completion requests:

```go
response, err := client.Send("anthropic/claude-haiku-4-5", "Hello, world!")
if err != nil {
    log.Fatal(err)
}

// Access response
fmt.Println(response.Text())         // Text content
fmt.Println(response.FinishReason()) // Finish reason
fmt.Println(response.ToolCalls())    // Tool calls (if any)

// Access usage and compression info
if response.Usage != nil {
    fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
}

if response.Compression != nil {
    fmt.Printf("Saved tokens: %d\n", response.Compression.SavedTokens)
    fmt.Printf("Reduction: %.1f%%\n", response.Compression.Reduction)
    fmt.Printf("Cost savings: $%.3f\n", float64(response.Compression.CostSavings)/1000000)
    fmt.Printf("Time: %d ms\n", response.Compression.TimeMs)
}
```

## Stream Method

The `Stream()` method enables real-time streaming responses:

```go
chunkChan, errChan := client.Stream("anthropic/claude-haiku-4-5", "Tell me a story")

for {
    select {
    case chunk, ok := <-chunkChan:
        if !ok {
            return
        }
        if text := chunk.Text(); text != "" {
            fmt.Print(text)
        }
        
        if reason := chunk.FinishReason(); reason != "" {
            fmt.Printf("\nFinished: %s\n", reason)
        }
    case err := <-errChan:
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

## Features

- ✅ **Type-safe** - Strong typing with Go structs and interfaces
- ✅ **OpenAI-compatible** - Works with any model supported by Edgee
- ✅ **Streaming** - Real-time response streaming with channels
- ✅ **Tool calling** - Full support for function calling
- ✅ **Flexible input** - Accept strings, InputObject, or maps
- ✅ **Compression info** - Access token compression metrics in responses
- ✅ **Minimal dependencies** - Uses only standard library and essential packages

## Documentation

For complete documentation, examples, and API reference, visit:

**👉 [Official Go SDK Documentation](https://www.edgee.ai/docs/sdk/go)**

The documentation includes:
- [Configuration guide](https://www.edgee.ai/docs/sdk/go/configuration) - Multiple ways to configure the SDK
- [Send method](https://www.edgee.ai/docs/sdk/go/send) - Complete guide to non-streaming requests
- [Stream method](https://www.edgee.ai/docs/sdk/go/stream) - Streaming responses guide
- [Tools](https://www.edgee.ai/docs/sdk/go/tools) - Function calling guide

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
