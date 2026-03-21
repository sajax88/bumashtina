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

	return "Успешно запазено"
}

func (a *App) SaveSettingsConfig(c Settings) string {
	// TODO: validation
	a.config.Settings = c
	err := SaveConfigToFile(a.config)
	if err != nil {
		return err.Error()
	}

	return "Успешно запазено"
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

func (a *App) LoadTaxesConfigLabels() []string {
	return GetLabelsForTaxesConfig()
}

func (a *App) LoadAllIncomeData() []IncomeForm {
	rows, err := GetAllDataFromFile()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func (a *App) SaveIncomeForm(f IncomeForm) string {
	// TODO: validation
	// TODO: calculate taxes and social security, save them together with the form
	// TaxesToPayCents          int64
	// SocialSecurityToPayCents int64
	// TODO: possibility to save really paid taxes and social security
	f.TaxesConfig = GetTaxesConfig()
	err := SaveDataToFile(f)
	if err != nil {
		return err.Error()
	}

	return "Успешно запазено"
}

func (a *App) DeleteData(month int, year int) string {
	err := DeleteDataFromFile(month, year)
	if err != nil {
		return err.Error()
	}
	return "Успешно изтрито"
}

func (a *App) GenerateDeclarationOne(month int, year int) string {
	incomeForm, err := GetDataFromFile(month, year)
	if err != nil {
		return err.Error()
	}

	personalData := a.LoadUserConfig()

	content, err := MakeDeclarationOne(incomeForm, personalData)
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
