package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gems/pkg/message"
	"io"
	"net/http"
	"strconv"
)

type openaiModel struct {
	baseUrl     string
	apiKey      string
	chatModel   string
	embedModel  string
	temperature string
}

type openaiChatParams struct {
	Model       string            `json:"model"`
	Messages    []message.Message `json:"messages"`
	Temperature float64           `json:"temperature"`
	Stream      bool              `json:"stream"`
}

type Choice struct {
	Message message.Message `json:"message"`
}

type openaiChatResponse struct {
	Choices []Choice `json:"choices"`
}

func (m *openaiModel) Chat(messages []message.Message) (string, error) {
	fmt.Println("Ich bin OpenAI 🌻 Ich weiß alles 🌻")
	chatEndpoint := fmt.Sprintf("%s/chat/completions", m.baseUrl)

	temperature, err := strconv.ParseFloat(m.temperature, 64)
	if err != nil {
		return "", fmt.Errorf("can't convert temperature from string to float64: %s", err)
	}
	jsonParams, err := json.Marshal(&openaiChatParams{
		Model:       m.chatModel,
		Messages:    messages,
		Temperature: temperature,
		Stream:      false,
	})
	if err != nil {
		return "", fmt.Errorf("cannot marshal the json data: %s", err)
	}

	// make a new request
	req, err := http.NewRequest("POST", chatEndpoint, bytes.NewReader(jsonParams))
	if err != nil {
		return "", fmt.Errorf("failed when making a request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))

	// create a client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot send the request with the client: %s", err)
	}

	byteResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read the response body %s", err)
	}
	structResp := &openaiChatResponse{}
	err = json.Unmarshal(byteResp, structResp)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal the response %s", err)
	}
	generatedContent := structResp.Choices[0].Message.Content

	return generatedContent, nil
}

func (m *openaiModel) Embed(text string) ([]float32, error) {
	return make([]float32, 0), nil
}
