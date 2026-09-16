package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type LatencyResult struct {
	URL             string `json:"url"`
	LatencyMS       int64  `json:"latency_ms"`
	SuggestedWaitMS int64  `json:"suggested_wait_ms"`
	Status          string `json:"status"`
}

type RetryDecision struct {
	Error       string `json:"error"`
	ShouldRetry bool   `json:"should_retry"`
	Reason      string `json:"reason"`
	MaxRetries  int    `json:"max_retries"`
}

func measureLatency(url string) LatencyResult {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return LatencyResult{URL: url, Status: "error"}
	}
	defer resp.Body.Close()

	latency := time.Since(start).Milliseconds()
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

// isTransientNetworkError checks if the error is a network issue that should be retried
func isTransientNetworkError(errStr string) bool {
	networkErrors := []string{"socket hang up", "econnreset", "etimedout", "network error", "econnrefused", "timeout"}
	lowerErr := strings.ToLower(errStr)
	for _, ne := range networkErrors {
		if strings.Contains(lowerErr, ne) {
			return true
		}
	}
	return false
}

func evaluateRetry(errMsg string, maxRetries int) RetryDecision {
	if isTransientNetworkError(errMsg) {
		return RetryDecision{
			Error:       errMsg,
			ShouldRetry: true,
			Reason:      "Transient network error detected. Safe to retry.",
			MaxRetries:  maxRetries,
		}
	}
	return RetryDecision{
		Error:       errMsg,
		ShouldRetry: false,
		Reason:      "Real test failure detected. Do NOT retry. Fix the code.",
		MaxRetries:  maxRetries,
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
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
		jsonOutput, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println("\n📄 JSON Output (for CI/CD):")
		fmt.Println(string(jsonOutput))

	case "retry":
		if len(os.Args) < 3 {
			fmt.Println("❌ Usage: orchestrator retry <error_message>")
			os.Exit(1)
		}
		errMsg := os.Args[2]
		decision := evaluateRetry(errMsg, 3)
		fmt.Printf("🧠 Evaluating error: %s\n\n", errMsg)
		fmt.Printf("🔍 Should Retry: %v\n", decision.ShouldRetry)
		fmt.Printf("📝 Reason:       %s\n", decision.Reason)
		fmt.Printf("🔄 Max Retries:  %d\n", decision.MaxRetries)
		jsonOutput, _ := json.MarshalIndent(decision, "", "  ")
		fmt.Println("\n📄 JSON Output (for CI/CD):")
		fmt.Println(string(jsonOutput))

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Network-Aware Test Orchestrator")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  orchestrator measure <target_url>        - Measure latency and calculate dynamic wait time")
	fmt.Println("  orchestrator retry <error_message>       - Evaluate if an error is transient and should be retried")
}