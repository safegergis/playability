package ai

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"playability/types"

	"github.com/liushuangls/go-anthropic/v2"
)

// SummarizeReports generates an AI-powered accessibility summary from multiple user reports
func SummarizeReports(reports []types.ReportRow) (string, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Printf("[AI] CLAUDE_API_KEY environment variable is not set")
		return "", fmt.Errorf("CLAUDE_API_KEY environment variable is not set")
	}

	if reports == nil || len(reports) == 0 {
		return "", fmt.Errorf("reports cannot be nil or empty")
	}

	client := anthropic.NewClient(apiKey)

	// Format reports for the prompt
	formattedReports := ""
	for i, report := range reports {
		formattedReports += fmt.Sprintf("Review %d: %s; Platform: %s\n", i+1, report.Report, report.Platform.String())
	}

	summarizationPrompt := fmt.Sprintf(`Analyze the game accessibility reviews provided below and create a summary paragraph to help other players make informed decisions.
Reviews:
<reviews>
%s
</reviews>

Analysis Requirements:
- Identify which platforms are mentioned (PC, PlayStation, Xbox, Switch, Mobile, etc.)
- Note features consistently praised across multiple reviews
- Note features consistently reported as missing or inadequate
- Track mentions of new versions/updates that added accessibility features
- Identify platform disparities (features on some platforms but not others)

Output ONLY a single paragraph that:
- Helps players understand what accessibility features are available
- Clearly indicates which platforms have which features
- Highlights any significant gaps or limitations players should know about
- Mentions positive features that work well for different accessibility needs
- Uses clear, player-friendly language that's easy to understand
- Provides actionable information to help players choose the right platform or version`, formattedReports)

	log.Printf("[AI] Sending summarization request to Claude API")

	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Dot5HaikuLatest,
		Messages: []anthropic.Message{
			{
				Role: "user",
				Content: []anthropic.MessageContent{
					{
						Type: "text",
						Text: &summarizationPrompt,
					},
				},
			},
		},
		MaxTokens: 1024,
	})

	if err != nil {
		var aiErr *anthropic.APIError
		if errors.As(err, &aiErr) {
			log.Printf("[AI] Claude API error - Type: %s, Message: %s", aiErr.Type, aiErr.Message)
			return "", fmt.Errorf("Claude API error (type: %s): %s", aiErr.Type, aiErr.Message)
		}
		log.Printf("[AI] Unexpected error calling Claude API: %v", err)
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	if len(resp.Content) == 0 {
		log.Printf("[AI] Claude API returned empty content")
		return "", fmt.Errorf("Claude API returned empty response")
	}

	summary := resp.Content[0].GetText()
	log.Printf("[AI] Claude API response (length: %d)", len(summary))

	return summary, nil
}
