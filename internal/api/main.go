package api

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Role string

const (
	pollinationsURL = "https://text.pollinations.ai/openai"
	modelsListURL   = "https://text.pollinations.ai/models"
	chanBufferSize  = 0
)

//go:embed prompts/unity.txt
var unityPrompt string

//go:embed prompts/default.txt
var defaultPrompt string

var modelList map[string]ModelDefinition

const (
	Assistant Role = "assistant"
	User      Role = "user"
	System    Role = "system"
)

var sseHeaders = map[string]string{
	"Content-Type": "application/json",
	"Accept":       "text/event-stream",
}

type ModelDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Uncensorded bool   `json:"uncensored,omitempty"`
	Reasoning   bool   `json:"reasoning,omitempty"`
}

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Stream   bool       `json:"stream"`
	Messages []*Message `json:"messages"`
	Model    string     `json:"model"`
	Private  bool       `json:"private"`
}

type OpenAIChunk struct {
	Role     Role     `json:"role"`
	Choices  []Choice `json:"choices"`
	Thinking bool     `json:"thinking"`
	Id       string   `json:"id"`
}

type Choice struct {
	FinishReason string `json:"finish_reason"`
	Delta        Delta  `json:"delta"`
}

type Delta struct {
	Content string `json:"content"`
}

func init() {
	// keep the model list in memory
	modelList = make(map[string]ModelDefinition)
	models := GetModels()
	for _, model := range models {
		modelList[model.Name] = model
	}
}

// Ask sends a request to the OpenAI API and returns a channel to receive the response chunks and the updated message history.
// TODO: find the seed in the prompts and manage a real uint rand value, because the LLM always want to provide 12345 :(
func Ask(prompt string, history []*Message, model string) (chan *OpenAIChunk, []*Message) {
	history = fixSystemPrompt(history, model)

	history = append(history, &Message{
		Role:    User,
		Content: prompt,
	})

	chunk := make(chan *OpenAIChunk, chanBufferSize)

	go CallAPI(&OpenAIRequest{
		Stream:   true,
		Private:  true,
		Messages: history,
		Model:    model,
	}, chunk)

	return chunk, history
}

// CallAPI sends a request to the OpenAI API and streams the response to the provided channel.
func CallAPI(r *OpenAIRequest, stream chan *OpenAIChunk) error {
	defer close(stream)

	client := &http.Client{}
	data, err := json.Marshal(r)
	if err != nil {
		log.Println("Error marshalling request:", err)
		return err
	}
	dataReader := bytes.NewReader(data)

	req, err := http.NewRequest(
		http.MethodPost,
		pollinationsURL,
		dataReader,
	)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	for k, v := range sseHeaders {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	model := GetModel(r.Model)
	// let's go!
	hadThought := false
	thinking := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 6 || line[:6] != "data: " {
			continue
		}

		// drop "data: " prefix
		data := line[6:]
		chunk := &OpenAIChunk{}
		err := json.Unmarshal([]byte(data), chunk)
		if err != nil {
			log.Println("Error unmarshalling chunk:", err)
			continue
		}

		// get the content
		if len(chunk.Choices) == 0 {
			continue
		}
		choice := chunk.Choices[0]
		if choice.FinishReason != "" {
			return nil
		}
		content := choice.Delta.Content

		// case of reasoning
		if model.Reasoning && !hadThought {
			if strings.Contains(content, "<think>") {
				thinking = true
			} else if strings.Contains(content, "</think>") {
				thinking = false
				hadThought = true
			}
			chunk.Thinking = thinking
			// clean the content
			content = strings.ReplaceAll(content, "<think>", "")
			content = strings.ReplaceAll(content, "</think>", "")
		}

		if content != "" {
			chunk.Role = Assistant
			stream <- chunk
		}
	}

	return nil
}

func GetModel(name string) ModelDefinition {
	if model, ok := modelList[name]; ok {
		return model
	}
	return ModelDefinition{}
}

// GetModels fetches the list of available models from the OpenAI API.
func GetModels() []ModelDefinition {
	if len(modelList) != 0 {
		return func() []ModelDefinition {
			models := make([]ModelDefinition, 0, len(modelList))
			for _, model := range modelList {
				models = append(models, model)
			}
			return models
		}()
	}

	resp, err := http.Get(modelsListURL)
	if err != nil {
		log.Println("Error fetching models:", err)
		return nil
	}

	response := []ModelDefinition{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error decoding models:", err)
		return nil
	}
	return response
}

// fixSystemPrompt checks if the first message in the history is a system prompt.
func fixSystemPrompt(history []*Message, model string) []*Message {
	var systemPrompt string
	switch model {
	case "unity":
		systemPrompt = unityPrompt
	default:
		systemPrompt = defaultPrompt
	}

	if len(history) == 0 {
		// no history, add the system prompt
		history = []*Message{
			{
				Role:    System,
				Content: systemPrompt,
			},
		}
	} else if history[0].Role != System {
		// the first message is not a system prompt, add it
		history = append([]*Message{
			{
				Role:    System,
				Content: systemPrompt,
			},
		}, history...)
	}
	return history
}
