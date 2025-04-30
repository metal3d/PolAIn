package main

import (
	"PolAIn/internal/api"
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var filesToSend = []string{}

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
	runtime.OnFileDrop(ctx, func(x, y int, files []string) {
		log.Println("Dropped files:", files)
		if !currentModel.Vision {
			runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: "Sorry, the current model cannot read input files",
				Title:   "Error",
			})
		} else {
			a.AddFiles(files)
		}
	})
}

func (a *App) AddFiles(files []string) {
	for _, f := range files {
		content, err := encodeFile(f)
		if err != nil {
			log.Println("Error encoding image:", err)
			continue
		}
		filesToSend = append(filesToSend, content)
		runtime.EventsEmit(a.ctx, "register-files", content)
	}
}

func (a *App) SelectFiles(filetype string) {
	filters := []runtime.FileFilter{}
	if filetype == "image" {
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Images",
			Pattern:     "*.jpeg;*.jpg;*.png;*.gif;*.webp",
		})
	}
	filename, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: filters,
	})
	if filename == "" || err != nil {
		return
	}
	a.AddFiles([]string{filename})
}

func (a *App) RemoveFile(pos int) bool {
	response, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "Delete Image",
		Message: "Are you sure you want to delete this image?",
	})
	if err != nil {
		log.Println("Error showing dialog:", err)
		return false
	}
	if response == "yes" {
		filesToSend = slices.Delete(filesToSend, pos, pos+1)
	}
	return true
}

// Ask sends a prompt to the OpenAI API and returns the response.
func (a *App) Ask(prompt string) error {
	runtime.EventsEmit(a.ctx, "ask-start", prompt)
	defer runtime.EventsEmit(a.ctx, "ask-done", prompt)
	defer func() {
		filesToSend = []string{}
	}()

	message := api.MessageContent{
		Type: "text",
		Text: &prompt,
	}

	toSend := []api.MessageContent{message}

	// append the files to send
	if currentModel.Vision {
		for _, f := range filesToSend {
			toSend = append(toSend, api.MessageContent{
				Type: "image_url",
				ImageURL: &map[string]string{
					"url": f,
				},
			},
			)
		}
	}

	// call the AI API
	stream, history := api.Ask(toSend, a.history, currentModel.Name)
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
		Content: []api.MessageContent{{Type: "text", Text: &buffer}},
	})
	log.Println("Final buffer", buffer)
	if len(strings.TrimSpace(buffer)) == 0 {
		return fmt.Errorf("model.empty.response")
	}
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
	filesToSend = []string{}
	runtime.EventsEmit(a.ctx, "new-conversation", a.history)
	return nil
}

// GetSelectedModel returns the selected model.
func (a *App) GetSelectedModel() *ModelPresentation {
	return currentModel
}
