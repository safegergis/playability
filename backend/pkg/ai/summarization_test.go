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

func TestSummarizeReports(t *testing.T) {
	t.Run("missing API key", func(t *testing.T) {
		// Temporarily unset API key
		originalKey := os.Getenv("CLAUDE_API_KEY")
		os.Unsetenv("CLAUDE_API_KEY")
		defer os.Setenv("CLAUDE_API_KEY", originalKey)

		reports := []types.ReportRow{
			{
				ID:       1,
				Platform: types.PC,
				Report:   "Great accessibility features",
			},
		}

		_, err := SummarizeReports(reports)

		if err == nil {
			t.Error("SummarizeReports() should return error when API key is missing")
		}
		if !strings.Contains(err.Error(), "CLAUDE_API_KEY environment variable is not set") {
			t.Errorf("SummarizeReports() error = %v, want 'CLAUDE_API_KEY environment variable is not set'", err)
		}
	})

	t.Run("nil reports", func(t *testing.T) {
		originalKey := os.Getenv("CLAUDE_API_KEY")
		os.Setenv("CLAUDE_API_KEY", "test-key")
		defer func() {
			if originalKey == "" {
				os.Unsetenv("CLAUDE_API_KEY")
			} else {
				os.Setenv("CLAUDE_API_KEY", originalKey)
			}
		}()

		_, err := SummarizeReports(nil)

		if err == nil {
			t.Error("SummarizeReports() should return error when reports is nil")
		}
		if !strings.Contains(err.Error(), "reports cannot be nil or empty") {
			t.Errorf("SummarizeReports() error = %v, want 'reports cannot be nil or empty'", err)
		}
	})

	t.Run("empty reports slice", func(t *testing.T) {
		originalKey := os.Getenv("CLAUDE_API_KEY")
		os.Setenv("CLAUDE_API_KEY", "test-key")
		defer func() {
			if originalKey == "" {
				os.Unsetenv("CLAUDE_API_KEY")
			} else {
				os.Setenv("CLAUDE_API_KEY", originalKey)
			}
		}()

		reports := []types.ReportRow{}

		_, err := SummarizeReports(reports)

		if err == nil {
			t.Error("SummarizeReports() should return error when reports is empty")
		}
		if !strings.Contains(err.Error(), "reports cannot be nil or empty") {
			t.Errorf("SummarizeReports() error = %v, want 'reports cannot be nil or empty'", err)
		}
	})

	// The following tests would require mocking the Anthropic API client
	// This requires refactoring the SummarizeReports function to accept a client interface

	t.Run("successful summarization with multiple reports - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")

		// Example of what the test would look like with proper DI:
		// mockClient := &MockAnthropicClient{
		//     Response: &anthropic.MessagesResponse{
		//         Content: []anthropic.MessageContent{
		//             {Text: "This game offers comprehensive accessibility features across all platforms..."},
		//         },
		//     },
		// }
		//
		// reports := []types.ReportRow{
		//     {ID: 1, Platform: types.PC, Report: "Excellent closed captions"},
		//     {ID: 2, Platform: types.Playstation, Report: "Great color blind modes"},
		// }
		//
		// summary, err := SummarizeReportsWithClient(reports, mockClient)
		// if err != nil {
		//     t.Errorf("SummarizeReports() unexpected error = %v", err)
		// }
		// if summary == "" {
		//     t.Error("SummarizeReports() returned empty summary")
		// }
	})

	t.Run("platform tracking in summary - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("API error handling - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("empty API response - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})

	t.Run("single report summarization - requires API mocking", func(t *testing.T) {
		t.Skip("Requires Anthropic API client mocking - client should be dependency injected")
	})
}

// Integration test - only runs if CLAUDE_API_KEY is set
func TestSummarizeReportsIntegration(t *testing.T) {
	if os.Getenv("CLAUDE_API_KEY") == "" {
		t.Skip("Skipping integration test - CLAUDE_API_KEY not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("multiple reports integration", func(t *testing.T) {
		reports := []types.ReportRow{
			{
				ID:       1,
				Platform: types.PC,
				Report:   "The PC version has excellent closed captions and color blind modes. All text is readable.",
			},
			{
				ID:       2,
				Platform: types.Playstation,
				Report:   "PlayStation version has good controller remapping but captions are sometimes hard to read.",
			},
			{
				ID:       3,
				Platform: types.PC,
				Report:   "Full controller support on PC works great with accessibility options.",
			},
		}

		summary, err := SummarizeReports(reports)

		if err != nil {
			t.Errorf("SummarizeReports() unexpected error = %v", err)
		}

		if summary == "" {
			t.Error("SummarizeReports() returned empty summary")
		}

		// Summary should be a single paragraph
		if strings.Count(summary, "\n\n") > 0 {
			t.Error("SummarizeReports() should return a single paragraph without double newlines")
		}

		// Summary should mention platforms
		if !strings.Contains(strings.ToLower(summary), "pc") && !strings.Contains(strings.ToLower(summary), "playstation") {
			t.Error("SummarizeReports() summary should mention the platforms from reports")
		}
	})
}
