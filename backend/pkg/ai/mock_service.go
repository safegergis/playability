package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"playability/types"
)

// MockAIService is a mock implementation of AIService for testing
type MockAIService struct {
	// ModerationError if set, Moderation will return this error
	ModerationError error
	// ModerationResponse if set, Moderation will return this response
	ModerationResponse *types.ModerationResponse
	// SummarizeError if set, SummarizeReports will return this error
	SummarizeError error
	// SummarizeResponse if set, SummarizeReports will return this response
	SummarizeResponse string
	// CallCount tracks how many times methods were called
	ModerationCallCount   int
	SummarizeCallCount    int
	LastModerationReport  *types.ReportRow
	LastSummarizeReports  []types.ReportRow
}

// Moderation is a mock implementation that returns pre-configured responses
func (m *MockAIService) Moderation(ctx context.Context, report *types.ReportRow) ([]byte, error) {
	m.ModerationCallCount++
	m.LastModerationReport = report

	if m.ModerationError != nil {
		return nil, m.ModerationError
	}

	if m.ModerationResponse != nil {
		return json.Marshal(m.ModerationResponse)
	}

	// Default: no violation
	return json.Marshal(types.ModerationResponse{
		Violation:  false,
		Categories: []string{},
	})
}

// SummarizeReports is a mock implementation that returns pre-configured responses
func (m *MockAIService) SummarizeReports(ctx context.Context, reports []types.ReportRow) (string, error) {
	m.SummarizeCallCount++
	m.LastSummarizeReports = reports

	if m.SummarizeError != nil {
		return "", m.SummarizeError
	}

	if m.SummarizeResponse != "" {
		return m.SummarizeResponse, nil
	}

	// Default: generic summary
	return fmt.Sprintf("Mock summary for %d reports", len(reports)), nil
}
