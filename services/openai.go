package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"io"

	"github.com/davidalvarez305/yd_cocktails/constants"
	"github.com/davidalvarez305/yd_cocktails/types"
)

func GetOpenAICompletionsResponse(prompt string, maxTokens int) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       "gpt-4o-mini",
		"messages":    []map[string]interface{}{{"role": "user", "content": prompt}},
		"max_tokens":  maxTokens,
		"temperature": 0.7,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+constants.OpenAIApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed: %s", body)
	}

	var openAiCompletionsResponse types.OpenAIChatCompletionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAiCompletionsResponse); err != nil {
		return "", err
	}

	if len(openAiCompletionsResponse.Choices) > 0 {
		return openAiCompletionsResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}
