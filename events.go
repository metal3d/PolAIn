package main

import (
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var filesToSend = []string{}

func (a *App) setupEvents() {
	runtime.OnFileDrop(a.ctx, a.onFileDrop)
}

func (a *App) onFileDrop(w, y int, files []string) {
	log.Println("Dropped files:", files)
	if !currentModel.Vision {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Message: "Sorry, the current model cannot read input files",
			Title:   "Error",
		})
	} else {
		a.AddFiles(files)
	}
}
