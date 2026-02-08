package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SaveConfig() string {
	c := Config{"Joe", "Test", "Smith"} // TODO
	err := SaveConfigToFile(c)
	if err != nil {
		return err.Error()
	}
	return "Success" // TODO: BG localization
}

func (a *App) LoadConfig() Config {
	c, err := LoadConfigFromFile()
	if err != nil {
		return Config{}
	}
	return c
}

func (a *App) GenerateDeclarationOne(v string) string {
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
	return fmt.Sprintf("Success %s, %d bytes writen to %s", v, len(content), path1)
}
