package main

import (
	"context"
	"fmt"
	"log"

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

func (a *App) SaveUserConfig(c UserConfig) string {
	if !c.IsValid() {
		return "Invalid config" // TODO: BG localization
	}
	a.config.user = c
	err := SaveConfigToFile(a.config)
	if err != nil {
		return err.Error()
	}

	return "Success" // TODO: BG localization
}

func (a *App) LoadUserConfig() UserConfig {
	if a.config != (Config{}) { // TODO
		return a.config.user
	}

	c, err := LoadConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}
	a.config = c
	return c.user
}

func (a *App) GenerateDeclarationOne() string {
	content := []byte("test,test,test") // TODO

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
