package ai

import (
	"os"
	"playability/types"
	"strings"
	"testing"
)

// Note: These tests require CLAUDE_API_KEY to be set and will make real API calls
// For true unit testing, the AI client should be dependency-injected
// Current implementation uses environment variables and direct API calls

func TestModeration(t *testing.T) {
	t.Run("missing API key", func(t *testing.T) {
		// Temporarily unset API key
		originalKey := os.Getenv("CLAUDE_API_KEY")
		os.Unsetenv("CLAUDE_API_KEY")
		defer os.Setenv("CLAUDE_API_KEY", originalKey)

		report := &types.ReportRow{
			ID:     1,
			Report: "This is a test report",
		}

		_, err := Moderation(report)

		if err == nil {
			t.Error("Moderation() should return error when API key is missing")
		}
		if !strings.Contains(err.Error(), "CLAUDE_API_KEY environment variable is not set") {
			t.Errorf("Moderation() error = %v, want 'CLAUDE_API_KEY environment variable is not set'", err)
		}
	})

	t.Run("nil report", func(t *testing.T) {
		originalKey := os.Getenv("CLAUDE_API_KEY")
		os.Setenv("CLAUDE_API_KEY", "test-key")
		defer func() {
			if originalKey == "" {
				os.Unsetenv("CLAUDE_API_KEY")
			} else {
				os.Setenv("CLAUDE_API_KEY", originalKey)
			}
		}()

		_, err := Moderation(nil)

		if err == nil {
			t.Error("Moderation() should return error when report is nil")
		}
		if !strings.Contains(err.Error(), "report cannot be nil") {
			t.Errorf("Moderation() error = %v, want 'report cannot be nil'", err)
		}
	})

	t.Run("empty report text", func(t *testing.T) {
		originalKey := os.Getenv("CLAUDE_API_KEY")
		os.Setenv("CLAUDE_API_KEY", "test-key")
		defer func() {
			if originalKey == "" {
				os.Unsetenv("CLAUDE_API_KEY")
			} else {
				os.Setenv("CLAUDE_API_KEY", originalKey)
			}
		}()

		report := &types.ReportRow{
			ID:     1,
			Report: "", // Empty report
		}

		_, err := Moderation(report)

		if err == nil {
			t.Error("Moderation() should return error when report text is empty")
		}
		if !strings.Contains(err.Error(), "report text cannot be empty") {
			t.Errorf("Moderation() error = %v, want 'report text cannot be empty'", err)
		}
	})

	// The following tests would require mocking the Anthropic API client
	// This requires refactoring the Moderation function to accept a client interface

	t.Run("API call with valid report - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")

		// Example of what the test would look like with proper DI:
		// mockClient := &MockAnthropicClient{
		//     Response: &anthropic.MessagesResponse{
		//         Content: []anthropic.MessageContent{
		//             {Text: `{"violation": false, "categories": []}`},
		//         },
		//     },
		// }
		//
		// report := &types.ReportRow{
		//     ID:     1,
		//     Report: "This game has excellent accessibility features",
		// }
		//
		// result, err := ModerationWithClient(report, mockClient)
		// ...
	})

	t.Run("violation detection - hate speech - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("violation detection - harassment - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("safe content passes moderation - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("API error handling - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("empty API response - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("JSON parsing error - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})
}

// Integration test - only runs if CLAUDE_API_KEY is set
func TestModerationIntegration(t *testing.T) {
	if os.Getenv("CLAUDE_API_KEY") == "" {
		t.Skip("Skipping integration test - CLAUDE_API_KEY not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("safe content integration", func(t *testing.T) {
		report := &types.ReportRow{
			ID:     1,
			Report: "This game has excellent closed captions and color blind modes. Great accessibility features!",
		}

		result, err := Moderation(report)

		if err != nil {
			t.Errorf("Moderation() unexpected error = %v", err)
		}

		if result == nil || len(result) == 0 {
			t.Error("Moderation() returned empty result")
		}

		// Result should be valid JSON with violation: false
		if !strings.Contains(string(result), "violation") {
			t.Error("Moderation() result should contain 'violation' field")
		}
	})
}
