package main

import (
	"PolAIn/internal/api"
	"context"
	"log"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

	stream, history := api.Ask(prompt, a.history, currentModel.Name)
	a.history = history
	buffer := ""
	for chunk := range stream {
		runtime.EventsEmit(a.ctx, "chunk", chunk)
		buffer += chunk.Choices[0].Delta.Content
	}
	a.history = append(history, &api.Message{
		Role:    api.Assistant,
		Content: buffer,
	})
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
