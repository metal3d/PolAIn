package main

import (
	"PolAIn/internal/api"
	"context"
	"log"
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
	a.setupEvents()
}

// shutdown is called when the app shuts down
func (a *App) shutdown(ctx context.Context) {
	log.Println("Shutting down...")
}
