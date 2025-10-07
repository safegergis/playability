package ai

import (
	"context"
	"errors"
	"fmt"
	"log"
	"playability/types"

	"github.com/liushuangls/go-anthropic/v2"
)

// ClaudeAIService implements AIService using Anthropic's Claude API
type ClaudeAIService struct {
	client *anthropic.Client
}

// NewClaudeAIService creates a new ClaudeAIService with the given API key
func NewClaudeAIService(apiKey string) *ClaudeAIService {
	if apiKey == "" {
		log.Printf("[ClaudeAIService] Warning: API key is empty - AI calls will fail")
	}
	client := anthropic.NewClient(apiKey)
	return &ClaudeAIService{client: client}
}

// Moderation analyzes a report for policy violations and returns a JSON response
func (c *ClaudeAIService) Moderation(ctx context.Context, report *types.ReportRow) ([]byte, error) {
	if report == nil {
		return nil, fmt.Errorf("report cannot be nil")
	}

	reportText := report.Report
	if reportText == "" {
		log.Printf("[Moderation] Report text is empty for report ID: %d", report.ID)
		return nil, fmt.Errorf("report text cannot be empty")
	}

	log.Printf("[Moderation] Starting moderation for report ID: %d", report.ID)
	log.Printf("[Moderation] Analyzing report text (length: %d characters)", len(reportText))

	// Define unsafe categories
	unsafeCategories := map[string]string{
		"Hate":           "Content that expresses, incites, or promotes hate based on race, gender, ethnicity, religion, nationality, sexual orientation, disability status, or caste.",
		"Harassment":     "Content that targets individuals or groups with the intent to intimidate, humiliate, or abuse.",
		"Violence":       "Content that promotes or glorifies violence or celebrates the suffering or humiliation of others.",
		"Self-harm":      "Content that promotes, encourages, or depicts acts of self-harm, such as suicide, cutting, and eating disorders.",
		"Sexual":         "Content meant to arouse sexual excitement, such as the description of sexual activity, or that promotes sexual services (excluding sex education and wellness).",
		"Minors":         "Sexual or suggestive content involving minors; content that exploits or sexualizes minors in any way.",
		"Hate_speech":    "Content that promotes hate or violence against individuals or groups based on protected characteristics.",
		"Extremism":      "Content that promotes or encourages extremist ideologies or activities.",
		"Graphic":        "Content that depicts death, violence, or serious physical injury in graphic detail.",
		"Illegal":        "Content that promotes or encourages illegal activities.",
		"Misinformation": "Content that spreads false or misleading information that could cause harm.",
	}

	// Construct unsafe category string
	var unsafeCategoryStr string
	for category, definition := range unsafeCategories {
		unsafeCategoryStr += category + ": " + definition + "\n"
	}

	// Construct the prompt for Claude
	assessmentPrompt := fmt.Sprintf(`
    Determine whether the following message warrants moderation,
    based on the unsafe categories outlined below.

    Message:
    <message>%s</message>

    Unsafe Categories:
    <categories>
    %s
    </categories>

    Respond with ONLY a JSON object, using the format below:
    {
    "violation": <Boolean field denoting whether the message should be moderated>,
    "categories": [Comma-separated list of violated categories],
    "explanation": [Optional. Only include if there is a violation.]
    }`, reportText, unsafeCategoryStr)

	log.Printf("[Moderation] Sending request to Claude API")

	temperature := float32(0.0)
	resp, err := c.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Haiku20240307,
		Messages: []anthropic.Message{
			{
				Role: "user",
				Content: []anthropic.MessageContent{
					{
						Type: "text",
						Text: &assessmentPrompt,
					},
				},
			},
		},
		MaxTokens:   200,
		Temperature: &temperature,
	})

	if err != nil {
		var aiErr *anthropic.APIError
		if errors.As(err, &aiErr) {
			log.Printf("[Moderation] Claude API error for report ID %d - Type: %s, Message: %s", report.ID, aiErr.Type, aiErr.Message)
			return nil, fmt.Errorf("Claude API error (type: %s): %s", aiErr.Type, aiErr.Message)
		}
		log.Printf("[Moderation] Unexpected error calling Claude API for report ID %d: %v", report.ID, err)
		return nil, fmt.Errorf("error calling Claude API: %w", err)
	}

	if len(resp.Content) == 0 {
		log.Printf("[Moderation] Claude API returned empty content for report ID %d", report.ID)
		return nil, fmt.Errorf("Claude API returned empty response")
	}

	responseText := resp.Content[0].GetText()
	log.Printf("[Moderation] Claude API response for report ID %d (length: %d): %s", report.ID, len(responseText), responseText)

	return []byte(responseText), nil
}

// SummarizeReports generates an AI-powered accessibility summary from multiple user reports
func (c *ClaudeAIService) SummarizeReports(ctx context.Context, reports []types.ReportRow) (string, error) {
	if reports == nil || len(reports) == 0 {
		return "", fmt.Errorf("reports cannot be nil or empty")
	}

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

	resp, err := c.client.CreateMessages(ctx, anthropic.MessagesRequest{
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
