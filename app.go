package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
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

func (a *App) GenerateDeclarationOne(v string) string {
	content := []byte("test,test,test")

	var options runtime.OpenDialogOptions
	dir, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return err.Error()
	}

	// TODO: OpenDirectoryDialog(ctx context.Context, dialogOptions OpenDialogOptions) (string, error)
	// TODO: SaveFileDialog?
	path1, err := filepath.Abs(dir + "/dat1.txt")
	if err != nil {
		// TODO a.log LogPrint(ctx, err.Error())
		return err.Error()
	}

	err = os.WriteFile(path1, content, 0644)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Success %s, %d bytes writen to %s", v, len(content), path1)
}
