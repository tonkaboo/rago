package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gems/pkg/message"
	"io"
	"log"
	"net/http"
)

type ollamaModel struct {
	port        string
	chatModel   string
	embedModel  string
	temperature string
}

type ollamaChatParams struct {
	Model    string            `json:"model"`
	Messages []message.Message `json:"messages"`
	Stream   bool              `json:"stream"`
}

type ollamaChatResponse struct {
	Message message.Message `json:"message"`
}

type ollamaEmbedParams struct {
	Model string `json:"model"`
	Text  string `json:"prompt"`
}

type ollamaEmbedResp struct {
	Embedding []float32 `json:"embedding"`
}

func (m *ollamaModel) Chat(messages []message.Message) (string, error) {
	// fmt.Println("Ich bin Ollama 🦙 Ich weiß fast alles 🦙")
	chatEndpoint := fmt.Sprintf("http://127.0.0.1:%s/api/chat", m.port)

	// prepare the params for the request
	jsonParams, err := json.Marshal(&ollamaChatParams{
		Model:    m.chatModel,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode the request params: %s", err)
	}
	reqBody := bytes.NewReader(jsonParams) // convert the json bytes to something readable (ie a reader)

	// send the request
	resp, err := http.Post(chatEndpoint, "application/json", reqBody)
	if err != nil {
		return "", fmt.Errorf("can't get a response: %s", err)
	}
	defer resp.Body.Close()

	// read from the response to a buffer
	var bufResp bytes.Buffer
	if _, err = bufResp.ReadFrom(resp.Body); err != nil {
		return "", fmt.Errorf("can't read from the response: %s", err)
	}

	// unmarshal the json bytes into a struct
	structResp := &ollamaChatResponse{}
	err = json.Unmarshal(bufResp.Bytes(), structResp)
	if err != nil {
		log.Println(err)
	}
	generatedContent := structResp.Message.Content

	return generatedContent, nil
}

func (m *ollamaModel) Embed(text string) ([]float32, error) {
	embedEndpoint := fmt.Sprintf("http://127.0.0.1:%s/api/embeddings", m.port)

	jsonParams, err := json.Marshal(&ollamaEmbedParams{
		Model: m.embedModel,
		Text:  text,
	})
	if err != nil {
		return nil, fmt.Errorf("an error occured when encoding the request params: %s", err)
	}
	reqBody := bytes.NewReader(jsonParams)

	resp, err := http.Post(embedEndpoint, "application/json", reqBody)
	if err != nil {
		return nil, fmt.Errorf("cannot get a response: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	byteResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read the response body: %s", err)
	}
	structResp := &ollamaEmbedResp{}
	if err := json.Unmarshal(byteResp, structResp); err != nil {
		return nil, fmt.Errorf("cannot parse the response: %s", err)
	}
	if structResp.Embedding == nil {
		return nil, fmt.Errorf("received nil embedding")
	}
	embedding := structResp.Embedding

	return embedding, nil
}
