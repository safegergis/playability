package ai

import (
	"context"
	"os"
	"playability/types"
	"strings"
	"testing"
)

func TestNewClaudeAIService(t *testing.T) {
	t.Run("creates service with valid API key", func(t *testing.T) {
		service := NewClaudeAIService("test-api-key")
		if service == nil {
			t.Error("NewClaudeAIService() returned nil")
		}
		if service.client == nil {
			t.Error("NewClaudeAIService() client is nil")
		}
	})

	t.Run("creates service with empty API key", func(t *testing.T) {
		service := NewClaudeAIService("")
		if service == nil {
			t.Error("NewClaudeAIService() returned nil")
		}
		// Service should still be created but will fail on API calls
	})
}

func TestClaudeAIService_Moderation(t *testing.T) {
	t.Run("nil report returns error", func(t *testing.T) {
		service := NewClaudeAIService("test-key")
		ctx := context.Background()

		_, err := service.Moderation(ctx, nil)

		if err == nil {
			t.Error("Moderation() should return error when report is nil")
		}
		if !strings.Contains(err.Error(), "report cannot be nil") {
			t.Errorf("Moderation() error = %v, want 'report cannot be nil'", err)
		}
	})

	t.Run("empty report text returns error", func(t *testing.T) {
		service := NewClaudeAIService("test-key")
		ctx := context.Background()

		report := &types.ReportRow{
			ID:     1,
			Report: "", // Empty report
		}

		_, err := service.Moderation(ctx, report)

		if err == nil {
			t.Error("Moderation() should return error when report text is empty")
		}
		if !strings.Contains(err.Error(), "report text cannot be empty") {
			t.Errorf("Moderation() error = %v, want 'report text cannot be empty'", err)
		}
	})

	t.Run("invalid API key returns error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping API call test in short mode")
		}

		service := NewClaudeAIService("invalid-key")
		ctx := context.Background()

		report := &types.ReportRow{
			ID:     1,
			Report: "This is a test report",
		}

		_, err := service.Moderation(ctx, report)

		if err == nil {
			t.Error("Moderation() should return error with invalid API key")
		}
	})
}

func TestClaudeAIService_SummarizeReports(t *testing.T) {
	t.Run("nil reports returns error", func(t *testing.T) {
		service := NewClaudeAIService("test-key")
		ctx := context.Background()

		_, err := service.SummarizeReports(ctx, nil)

		if err == nil {
			t.Error("SummarizeReports() should return error when reports is nil")
		}
		if !strings.Contains(err.Error(), "reports cannot be nil or empty") {
			t.Errorf("SummarizeReports() error = %v, want 'reports cannot be nil or empty'", err)
		}
	})

	t.Run("empty reports slice returns error", func(t *testing.T) {
		service := NewClaudeAIService("test-key")
		ctx := context.Background()

		reports := []types.ReportRow{}

		_, err := service.SummarizeReports(ctx, reports)

		if err == nil {
			t.Error("SummarizeReports() should return error when reports is empty")
		}
		if !strings.Contains(err.Error(), "reports cannot be nil or empty") {
			t.Errorf("SummarizeReports() error = %v, want 'reports cannot be nil or empty'", err)
		}
	})

	t.Run("invalid API key returns error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping API call test in short mode")
		}

		service := NewClaudeAIService("invalid-key")
		ctx := context.Background()

		reports := []types.ReportRow{
			{
				ID:       1,
				Platform: types.PC,
				Report:   "Great accessibility features",
			},
		}

		_, err := service.SummarizeReports(ctx, reports)

		if err == nil {
			t.Error("SummarizeReports() should return error with invalid API key")
		}
	})
}

// Integration tests - only run if CLAUDE_API_KEY is set
func TestClaudeAIService_ModerationIntegration(t *testing.T) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test - CLAUDE_API_KEY not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("safe content integration", func(t *testing.T) {
		service := NewClaudeAIService(apiKey)
		ctx := context.Background()

		report := &types.ReportRow{
			ID:     1,
			Report: "This game has excellent closed captions and color blind modes. Great accessibility features!",
		}

		result, err := service.Moderation(ctx, report)

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

func TestClaudeAIService_SummarizeReportsIntegration(t *testing.T) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test - CLAUDE_API_KEY not set")
	}

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("multiple reports integration", func(t *testing.T) {
		service := NewClaudeAIService(apiKey)
		ctx := context.Background()

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

		summary, err := service.SummarizeReports(ctx, reports)

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
		summaryLower := strings.ToLower(summary)
		if !strings.Contains(summaryLower, "pc") && !strings.Contains(summaryLower, "playstation") {
			t.Error("SummarizeReports() summary should mention the platforms from reports")
		}
	})
}
