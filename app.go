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

func (a *App) LoadThisMonthActions() string {
	return "" // TODO

	// TODO
	//- Декларация 1 за дължими осигуровки (по Булстат) – всеки месец до 25-то число на следващия месец;
	//- Декларация 6 за дължими осигурителни вноски (по ЕГН) – до 30.04 на следващата календарна година;
	//- Декларация по чл. 55 от ЗДДФЛ (по ЕГН) – за първите три тримесечия, до края на месеца, следващ
	//тримесечието;
	// - Декларация по чл. 50 от ЗДДФЛ (по ЕГН) – до 30.04 на следващата календарна година;
}

func (a *App) SaveIncomeForm(f IncomeForm) string {
	// TODO: validation

	existingForm, err := GetDataFromFileForMonth(int(f.Month), int(f.Year))
	if err != nil {
		return err.Error()
	}

	if existingForm.Month > 0 {
		return fmt.Sprintf("Данните за %d/%d вече съществуват, изтрийте ги за да въведете нови", f.Month, f.Year)
	}

	f.TaxesConfig = a.LoadTaxesConfig()
	f.Settings = a.LoadSettingsConfig()

	// Calculate approximate taxes and social security, save them together with the form
	CalculateSocialSecurity(&f)
	f.TaxesToPayCents = CalculateTaxForMonth(f)

	err = SaveDataToFile(f)
	if err != nil {
		return err.Error()
	}

	return "Успешно запазено"
}

func (a *App) UpdateForm(f IncomeForm) string {
	// TODO: validation

	existingForm, err := GetDataFromFileForMonth(int(f.Month), int(f.Year))
	if err != nil {
		return err.Error()
	}

	if existingForm.Month > 0 {
		// Allow to update only some values
		existingForm.SocialSecurityReallyPaidCents = f.SocialSecurityReallyPaidCents
		existingForm.TaxesReallyPaidCents = f.TaxesReallyPaidCents
	} else {
		return fmt.Sprintf("Данните за месец %d не са намерени", existingForm.Month) // TODO: return err object?
	}

	// Remove previous record and save the new one
	err = DeleteDataFromFile(int(f.Month), int(f.Year))
	if err != nil {
		return err.Error()
	}

	err = SaveDataToFile(existingForm)
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
	var res string

	// TODO: check that personal data is entered!

	incomeForm, err := GetDataFromFileForMonth(month, year)
	if err != nil {
		return err.Error()
	}

	personalData := a.LoadUserConfig()
	settings := a.LoadSettingsConfig()

	content, err := MakeDeclarationOne(incomeForm, personalData, settings)
	if err != nil {
		return err.Error()
	}

	res, err = a.SaveDeclarationToFile(content, fmt.Sprintf("Declaration_1_%d_%02d", year, month))
	if err != nil {
		return err.Error()
	}

	return res
}

func (a *App) PreviewDeclarationSix(year int) string {

	// TODO: check that personal data is entered!

	// TODO: calculate insurances to pay
	// TODO: display them to the user for correction if needed, then submit and actually make the declaration?

	// TODO: calc for total year - 14.8 % PensionPercentagePartOne OR + 3.5% = 18.3% for PregnancyInsurancePercentage
	// TODO: calc for total year - 5% PensionPercentagePartTwo
	// TODO: calc for total year - 8% HealthInsurancePercentage

	return "" // TODO
}

// TODO: pass sums []float64
func (a *App) GenerateDeclarationSix(year int) string {
	var res string

	// TODO: check that personal data is entered!

	personalData := a.LoadUserConfig()

	sums := []float64{1500.50, 1200.34, 1200.00}
	if len(sums) != 3 {
		return "Неверни данни"
	}

	content, err := MakeDeclarationSix(year, personalData, sums)
	if err != nil {
		return err.Error()
	}

	res, err = a.SaveDeclarationToFile(content, fmt.Sprintf("Declaration_6_%d", year))
	if err != nil {
		return err.Error()
	}

	return res
}

func (a *App) SaveDeclarationToFile(content []byte, filename string) (string, error) {
	var err error
	var options runtime.SaveDialogOptions

	options.DefaultFilename = fmt.Sprintf("%s.txt", filename)
	options.DefaultDirectory, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}
	saveFilePath, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return "", err
	}

	saveFilePath, err = SaveDeclaration(saveFilePath, content)
	if err != nil {
		return "", err
	}

	// TODO a.log LogPrint(ctx, err.Error()), MessageDialog

	return fmt.Sprintf("Успешно, файлът беше записан в %s", saveFilePath), nil
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
	result.TaxCents, err = CalculateAdvanceTaxForThreeMonths(rows, &result)
	if err != nil {

		log.Fatal(err)
	}

	return result
}

func (a *App) SavePaidTaxForQuarter(quarter int, year int, amount float32) string {
	// TODO: object with status instead of string, in all cases?
	// Find the last month in the quarter and save the paid tax there
	lastMonth := (quarter-1)*3 + 3
	row, err := GetDataFromFileForMonth(lastMonth, year)
	if err != nil {
		return err.Error()
	}
	if row.Month != int16(lastMonth) {
		return "Въведете данните за трите месеца за да запазите платения данък за тримесечието"
	}
	row.TaxesReallyPaidCents = int64(amount * MONEY_DIVIDER)
	return a.UpdateForm(row)
}
