package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	config Config
}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SaveConfig(c Config) string {
	err := SaveConfigToFile(c)
	if err != nil {
		return err.Error()
	}
	a.config = c
	return "Success" // TODO: BG localization
}

func (a *App) LoadConfig() Config {
	if a.config.LastName != "" { // TODO
		return a.config
	}

	c, err := LoadConfigFromFile()
	if err != nil {
		return Config{}
	}
	a.config = c
	return c
}

func (a *App) GenerateDeclarationOne() string {
	content := []byte("test,test,test")

	var options runtime.OpenDialogOptions // TODO: options
	dir, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return err.Error()
	}

	path1, err := SaveDeclaration(dir, content)
	if err != nil {
		return err.Error()
	}

	// TODO a.log LogPrint(ctx, err.Error()), MessageDialog

	// TODO: update message
	return fmt.Sprintf("Success %s, %d bytes writen to %s", len(content), path1)
}
