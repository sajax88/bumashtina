package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
		return "Грешна конфигурация"
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

func (a *App) SaveTaxesConfig(t TaxesConfig) string {
	// TODO: validation
	a.config.TaxesConfig = t
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
	if a.config != (Config{}) { // TODO
		return a.config.TaxesConfig
	}

	c, err := LoadConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}
	a.config = c
	return c.TaxesConfig
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

func (a *App) LoadIncomeDataForMonth(month int, year int) IncomeForm {
	row, err := GetDataFromFileForMonth(month, year)
	if err != nil {
		log.Fatal(err)
	}
	return row
}

func (a *App) SaveIncomeForm(f IncomeForm) string {
	// TODO: validation

	// TODO: calculate approximate taxes and social security, save them together with the form?
	// TaxesToPayCents          int64
	// SocialSecurityToPayCents int64
	// TODO: update and save actually paid?

	// TODO: check if we're trying to override, ask to delete the previos value first

	f.TaxesConfig = a.LoadTaxesConfig()
	f.Settings = a.LoadSettingsConfig()
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
	incomeForm, err := GetDataFromFileForMonth(month, year)
	if err != nil {
		return err.Error()
	}

	personalData := a.LoadUserConfig()

	content, err := MakeDeclarationOne(incomeForm, personalData)
	if err != nil {
		return err.Error()
	}

	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("Declaration1_%d_%02d.txt", year, month)
	options.DefaultDirectory, err = os.UserHomeDir()
	if err != nil {
		return err.Error()
	}
	saveFilePath, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return err.Error()
	}

	saveFilePath, err = SaveDeclaration(saveFilePath, content)
	if err != nil {
		return err.Error()
	}

	// TODO a.log LogPrint(ctx, err.Error()), MessageDialog

	return fmt.Sprintf("Успешно, файлът беше записан в %s", saveFilePath)
}

func (a *App) CalculateTaxForQuarter(quarter int, year int) CalculatedTax {
	// TODO: validate

	result := CalculatedTax{
		TotalIncomeCents: 0,
		TaxCents:         0,
		Quarter:          quarter,
		Year:             year,
	}

	rows, err := GetDataFromFileForQuarter(quarter, year, &result)
	if err != nil {

		log.Fatal(err)
	}

	if len(rows) == 0 {
		return result
	}

	result.TotalIncomeCents, err = CalculateIncomeForThreeMonths(rows)
	if err != nil {

		log.Fatal(err)
	}

	// TODO: if really paid insurance was not entered, add notification?
	result.TaxCents, err = CalculateAdvanceTaxForThreeMonths(rows)
	if err != nil {

		log.Fatal(err)
	}

	return result
}

func (a *App) GenerateDeclarationSix() string {
	// TODO
	return "TODO"
}
