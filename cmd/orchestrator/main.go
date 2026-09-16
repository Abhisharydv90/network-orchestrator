package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LatencyResult struct {
	URL          string `json:"url"`
	LatencyMS    int64  `json:"latency_ms"`
	SuggestedWaitMS int64 `json:"suggested_wait_ms"`
	Status       string `json:"status"`
}

func measureLatency(url string) LatencyResult {
	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	
	// Make a GET request (HEAD is faster, but GET is more realistic for test scenarios)
	resp, err := client.Get(url)
	if err != nil {
		return LatencyResult{
			URL:    url,
			Status: "error",
		}
	}
	defer resp.Body.Close()

	// Calculate latency in milliseconds
	latency := time.Since(start).Milliseconds()

	// The core logic: Math.max(500, latency * 1.5)
	// If the network is fast, wait at least 500ms. If slow, wait 1.5x the latency.
	suggestedWait := latency * 3 / 2
	if suggestedWait < 500 {
		suggestedWait = 500
	}

	return LatencyResult{
		URL:             url,
		LatencyMS:       latency,
		SuggestedWaitMS: suggestedWait,
		Status:          "success",
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: orchestrator measure <target_url>")
		fmt.Println("Example: orchestrator measure http://localhost:3000")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "measure":
		if len(os.Args) < 3 {
			fmt.Println("❌ Usage: orchestrator measure <target_url>")
			os.Exit(1)
		}
		targetURL := os.Args[2]
		
		fmt.Printf("📡 Measuring network latency to: %s\n\n", targetURL)
		
		result := measureLatency(targetURL)
		
		if result.Status == "error" {
			fmt.Printf("❌ Failed to connect to %s\n", targetURL)
			os.Exit(1)
		}

		fmt.Printf("📊 Latency:        %d ms\n", result.LatencyMS)
		fmt.Printf("⏱️  Suggested Wait: %d ms\n", result.SuggestedWaitMS)
		fmt.Printf("✅ Status:         %s\n", result.Status)
		
		// Output as JSON for CI/CD pipelines to consume
		jsonOutput, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println("\n📄 JSON Output (for CI/CD):")
		fmt.Println(string(jsonOutput))

	default:
		fmt.Println("❌ Unknown command. Use 'measure'.")
		os.Exit(1)
	}
}