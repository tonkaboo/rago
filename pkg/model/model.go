package model

import (
	"encoding/json"
	"fmt"
	"gems/pkg/message"
	"io"
	"log"
	"os"
)

type Model interface {
	Chat([]message.Message) (string, error)
	Embed(string) ([]float32, error)
}

type modelConfig struct {
	ApiType     string `json:"apiType"`
	ChatModel   string `json:"chatModel"`
	EmbedModel  string `json:"embedModel"`
	BaseUrl     string `json:"baseUrl"`
	Port        string `json:"port"`
	ApiKey      string `json:"apiKey"`
	Temperature string `json:"temperature"`
}

// get the llm configuration from the specified file
func getConfig() *modelConfig {
	// open the specified file
	configFile, err := os.Open("../../pkg/model/modelConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	// read from the json file into a struct
	fileBytes, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config modelConfig
	json.Unmarshal(fileBytes, &config)
	return &config
}

func NewModel() (Model, error) {
	var config *modelConfig = getConfig()
	var model Model

	// return the correct model according to the ApiType
	if config.ApiType == "openai" {
		model = &openaiModel{
			baseUrl:     config.BaseUrl,
			apiKey:      config.ApiKey,
			chatModel:   config.ChatModel,
			embedModel:  config.EmbedModel,
			temperature: config.Temperature,
		}
		return model, nil
	} else if config.ApiType == "ollama" {
		model = &ollamaModel{
			port:        config.Port,
			chatModel:   config.ChatModel,
			embedModel:  config.EmbedModel,
			temperature: config.Temperature,
		}
		return model, nil
	} else {
		// if the specified type isn't a known type, return error
		return nil, fmt.Errorf("unsupported API type: %s", config.ApiType)
	}
}
