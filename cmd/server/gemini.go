package main

import (
	"context"

	"google.golang.org/genai"
)

func IssueAiPrompt(prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
