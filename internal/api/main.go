package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Role string

const (
	pollinationsURL = "https://text.pollinations.ai/openai"
	modelsListURL   = "https://text.pollinations.ai/models"
	chanBufferSize  = 0
)

const systemPrompt = `You are a very smart and helpful assistant. 
You are able to answer any question and provide information on a wide range of topics. 
You are also able to generate text in a variety of styles and formats, including poetry, prose, and technical writing. 
You are always polite and respectful, and you strive to provide the best possible answers to your users' questions.

If the user asks for an image, use the pollinations url using this template using url encoded description (translate in english if needed),
encode the description for URL (space become %20, etc), set width and height to something reasonable or use the requirements 
from the user (adapt to square of 2 values, you may use landscape, square and portrait resolution), and use a random seed value 
to force generation of new image, the seed should be a valid unsigned integer:

![Image description](https://image.pollinations.ai/prompt/{description}?nologo=true&private=true&enhance=true&width={width}&height={height}&seed={seed})

You answer using markdown. Use the language of the user.
`

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
	Role    Role     `json:"role"`
	Choices []Choice `json:"choices"`
	Id      string   `json:"id"`
}

type Choice struct {
	FinishReason string `json:"finish_reason"`
	Delta        Delta  `json:"delta"`
}

type Delta struct {
	Content string `json:"content"`
}

// Ask sends a request to the OpenAI API and returns a channel to receive the response chunks and the updated message history.
// TODO: find the seed in the prompts and manage a real uint rand value, because the LLM always want to provide 12345 :(
func Ask(prompt string, history []*Message, model string) (chan *OpenAIChunk, []*Message) {
	if len(history) == 0 {
		history = []*Message{
			{
				Role:    System,
				Content: systemPrompt,
			},
		}
	}

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

	fmt.Println("Asking:", prompt, history)
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
		if content != "" {
			chunk.Role = Assistant
			stream <- chunk
		}
	}

	return nil
}

// GetModels fetches the list of available models from the OpenAI API.
func GetModels() []ModelDefinition {
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
