package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// `OpenAIClient` struct holds an API key for interacting with the OpenAI API.
type OpenAIClient struct {
	APIKey string
}

// `Message` struct represents a JSON-serializable message with fields for role and content.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// `Choice` struct defines a choice with an index, a message (using `Message` struct), optional log probabilities,
// and a finish reason, suitable for JSON serialization.
type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	Logprobs     []string `json:"logprobs,omitempty"`
	FinishReason string   `json:"finish_reason"`
}

// `CompletionResponse` struct encapsulates a list of `Choice` structs,
// intended for JSON serialization as "choices".
type CompletionResponse struct {
	Choices []Choice `json:"choices"`
}

// NewOpenAIClient creates a new instance of the OpenAI client.
// apiKey: The API key for accessing the OpenAI API.
// Returns: A pointer to the newly created OpenAIClient instance.
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

// genereateCompletionRequest generates a JSON request for the OpenAI completion API.
// It takes a prompt text and a ChatSystem pointer as parameters.
// The function constructs a JSON object with the specified model and messages,
// including the system message and the user's prompt.
// It returns a byte slice containing the JSON request and an error if any.
func (c *OpenAIClient) genereateCompletionRequest(promptText string, ChatSystem *ChatSystem) ([]byte, error) {
	// Define the messages array with the system message and user's prompt
	messages := []Message{
		{Role: "system", Content: ChatSystem.SystemMessage},
		{Role: "user", Content: promptText}, // Use the text variable here
	}

	// Define the request JSON structure
	requestJson := struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}{
		Model:    "gpt-3.5-turbo-16k",
		Messages: messages,
	}

	// Marshal the request JSON structure to a byte slice with indentation
	requestJSONString, err := json.MarshalIndent(requestJson, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return []byte("Failed to generate the JSON"), nil
	}

	// Return the byte slice containing the JSON request and nil error
	return requestJSONString, nil
}
