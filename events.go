package main

import (
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var filesToSend = []string{}

func (a *App) setupEvents() {
	runtime.OnFileDrop(a.ctx, a.onFileDrop)
}

// onFileDrop is called when files are dropped on the window. We check if
// the model can maanage the files and if so, we add them to the list of files to send.
// TODO: limit the number of files to send
// TODO: mangage audio files
func (a *App) onFileDrop(w, y int, files []string) {
	log.Println("Dropped files:", files)
	if !currentModel.Vision {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: "Sorry, the current model cannot read input files",
			Title:   "Error",
		})
	} else {
		a.addFiles(files)
	}
}

func (a *App) addFiles(files []string) {
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
