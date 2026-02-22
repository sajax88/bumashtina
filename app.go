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
	a.config.User = c
	err := SaveConfigToFile(a.config)
	if err != nil {
		return err.Error()
	}

	return "Success" // TODO: BG localization
}

func (a *App) SaveSettingsConfig(c Settings) string {
	// TODO: validation
	a.config.Settings = c
	err := SaveConfigToFile(a.config)
	if err != nil {
		return err.Error()
	}

	return "Success" // TODO: BG localization
}

func (a *App) LoadUserConfig() UserConfig {
	if a.config != (Config{}) { // TODO
		return a.config.User
	}

	c, err := LoadConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}
	a.config = c
	return c.User
}

func (a *App) LoadSettingsConfig() Settings {
	if a.config != (Config{}) { // TODO
		return a.config.Settings
	}

	c, err := LoadConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}
	a.config = c
	return c.Settings
}

func (a *App) LoadTaxesConfig() TaxesConfig {
	return GetTaxesConfig()
}

func (a *App) GenerateDeclarationOne() string {
	content, err := MakeDeclarationOne()
	if err != nil {
		return err.Error()
	}

	var options runtime.OpenDialogOptions // TODO: options default
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
	return fmt.Sprintf("Success, %d bytes writen to %s", len(content), path1)
}

func (a *App) GenerateDeclarationSix() string {
	// TODO
	return "TODO"
}

func (a *App) GetRemindersText() string {
	// TODO

	// Декларации за данъци и осигуровки:
	// Декларация 1 за дължими осигуровки – всеки месец от 25-о число на следващия месец;
	// Декларация 6 за дължими осигурителни вноски – до 30.04 на следващата календарна година;
	// Декларация по чл. 55 от ЗДДФЛ – до края на месеца следващ тримесечието, за което декларацията се подава (само за първите три тримесечия);
	// Декларация по чл. 50 от ЗДДФЛ – до 30.04 на годината следваща тази, за която се подава декларацията.
	// При ДДС регистрация – ежемесечни ДДС декларации до 15-о число на следващия месец.

	return "TODO"
}
