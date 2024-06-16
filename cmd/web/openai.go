package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenAIClient represents the OpenAI API client configuration.
type OpenAIClient struct {
	APIKey string
}

// CompletionResponse represents the response structure from OpenAI's completion endpoint.

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	Logprobs     []string `json:"logprobs,omitempty"`
	FinishReason string   `json:"finish_reason"`
}

type CompletionResponse struct {
	Choices []Choice `json:"choices"`
}

// NewOpenAIClient creates a new instance of the OpenAI client.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: apiKey,
	}
}

// CompleteText sends a prompt to OpenAI's completion endpoint and returns the generated text.
func (c *OpenAIClient) GetCompletionResponse(promptText string, ChatSystem *ChatSystem) ([]byte, error) {

	url := "https://api.openai.com/v1/chat/completions"

	requestJSONString, err := c.genereateCompletionRequest(promptText, ChatSystem)

	if err != nil {
		return []byte(""), err // TODO: Handle error properly
	}

	convertRequest := string(requestJSONString)

	// Create HTTP client
	client := &http.Client{}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(convertRequest))
	if err != nil {
		return []byte(""), err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}
	//defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	// Parse JSON response
	var completionResponse CompletionResponse
	if err := json.Unmarshal(body, &completionResponse); err != nil {
		return []byte(""), err
	}

	if len(completionResponse.Choices) > 0 {
		return []byte(completionResponse.Choices[0].Message.Content), nil
	}

	return []byte(""), fmt.Errorf("no completion response received")
}

func (c *OpenAIClient) genereateCompletionRequest(promptText string, ChatSystem *ChatSystem) ([]byte, error) {
	messages := []Message{
		{Role: "system", Content: ChatSystem.SystemMessage},
		{Role: "user", Content: promptText}, // Use the text variable here
	}

	requestJson := struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}{
		Model:    "gpt-3.5-turbo-16k",
		Messages: messages,
	}

	requestJSONString, err := json.MarshalIndent(requestJson, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return []byte("Failed to generate the JSON"), nil
	}

	return requestJSONString, nil

}
