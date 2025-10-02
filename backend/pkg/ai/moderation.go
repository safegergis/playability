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

func Moderation(report *types.ReportRow) ([]byte, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Printf("[Moderation] CLAUDE_API_KEY environment variable is not set")
		return nil, fmt.Errorf("CLAUDE_API_KEY environment variable is not set")
	}

	if report == nil {
		return nil, fmt.Errorf("report cannot be nil")
	}

	log.Printf("[Moderation] Starting moderation for report ID: %d", report.ID)

	client := anthropic.NewClient(apiKey)

	temperature := float32(0.0)
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

	reportText := report.Report
	if reportText == "" {
		log.Printf("[Moderation] Report text is empty for report ID: %d", report.ID)
		return nil, fmt.Errorf("report text cannot be empty")
	}

	log.Printf("[Moderation] Analyzing report text (length: %d characters)", len(reportText))

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
    {{
    "violation": <Boolean field denoting whether the message should be moderated>,
    "categories": [Comma-separated list of violated categories],
    "explanation": [Optional. Only include if there is a violation.]
    }}`, reportText, unsafeCategoryStr)
	log.Printf("[Moderation] Sending request to Claude API")

	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
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
		} else {
			log.Printf("[Moderation] Unexpected error calling Claude API for report ID %d: %v", report.ID, err)
			return nil, fmt.Errorf("error calling Claude API: %w", err)
		}
	}

	if len(resp.Content) == 0 {
		log.Printf("[Moderation] Claude API returned empty content for report ID %d", report.ID)
		return nil, fmt.Errorf("Claude API returned empty response")
	}

	responseText := resp.Content[0].GetText()
	log.Printf("[Moderation] Claude API response for report ID %d (length: %d): %s", report.ID, len(responseText), responseText)

	return []byte(responseText), nil
}
