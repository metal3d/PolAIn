package main

import (
	"PolAIn/internal/api"
	"fmt"
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

// GetSelectedModel returns the selected model.
func (a *App) GetSelectedModel() *ModelPresentation {
	return currentModel
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

	if len(strings.TrimSpace(buffer)) == 0 {
		return fmt.Errorf("%s", a.Translate("model.empty.response"))
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

// SelectFiles is called when the user press image or audio button.
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
	a.addFiles([]string{filename})
}

// RemoveFile is called when the user press delete button on the image or audio file.
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
