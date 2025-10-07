package ai

import (
	"context"
	"playability/types"
)

// AIService defines the interface for AI-powered content moderation and summarization
type AIService interface {
	// Moderation analyzes a report for policy violations and returns a JSON response
	Moderation(ctx context.Context, report *types.ReportRow) ([]byte, error)

	// SummarizeReports generates an accessibility summary from multiple user reports
	SummarizeReports(ctx context.Context, reports []types.ReportRow) (string, error)
}
