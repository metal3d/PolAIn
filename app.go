package main

import (
	"PolAIn/internal/api"
	"context"
	"log"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Rendered struct represents the rendered response for the view.
type Rendered struct {
	// Chunk is the chunk received from the OpenAI API.
	Chunk *api.OpenAIChunk `json:"chunk"`
	// Html is the cummulated HTML from the beginning of the reponse.
	Html string `json:"html"`
	// ThinkingHTML is the cummulated HTML from the beginning of the reponse of the reasoning model.
	ThinkingHTML string `json:"thinkingHtml"`
}

// App struct
type App struct {
	ctx     context.Context
	history []*api.Message
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Ask sends a prompt to the OpenAI API and returns the response.
func (a *App) Ask(prompt string) error {
	runtime.EventsEmit(a.ctx, "ask-start", prompt)
	defer runtime.EventsEmit(a.ctx, "ask-done", prompt)

	// call the AI API
	stream, history := api.Ask(prompt, a.history, currentModel.Name)
	a.history = history

	// on chunk received, fix the markdown, create HTML and emit the event
	var buffer, html, thinkingBuffer, thinkingHtml string
	for chunk := range stream {
		if chunk.Thinking {
			thinkingBuffer += chunk.Choices[0].Delta.Content
			thinkingBuffer = fixKatex(thinkingBuffer)
			thinkingHtml = string(MDtoHTML(thinkingBuffer))
		} else {
			buffer += chunk.Choices[0].Delta.Content
			buffer = fixKatex(buffer)
			html = string(MDtoHTML(buffer))
		}

		runtime.EventsEmit(a.ctx, "chunk", Rendered{
			Chunk:        chunk,
			Html:         string(html),
			ThinkingHTML: string(thinkingHtml),
		})
	}
	a.history = append(history, &api.Message{
		Role:    api.Assistant,
		Content: buffer,
	})
	log.Println("Final buffer", buffer)
	return nil
}

// NewConversation creates a new conversation, it removes the history and send an event.
func (a *App) NewConversation() error {
	md, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   a.Translate("menu.conversation.new"),
		Message: a.Translate("conversation.new.confirm"),
	})
	if err != nil {
		log.Println("Error showing dialog:", err)
		return err
	}
	// probably not needed, but I'm not sure about the text
	if slices.Contains([]string{"no", "cancel", "cancelled"}, strings.ToLower(md)) {
		return nil
	}
	a.history = []*api.Message{}
	runtime.EventsEmit(a.ctx, "new-conversation", a.history)
	return nil
}

// GetSelectedModel returns the selected model.
func (a *App) GetSelectedModel() *ModelPresentation {
	return currentModel
}
