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

func (a *App) SaveUserConfig(u UserConfig) string {
	isValid, errorMessage := u.Validate()
	if !isValid {
		ShowWarningDialog(a.ctx, "Грешна конфигурация", errorMessage)
		return ""
	}

	a.config.User = u
	err := SaveConfigToFile(a.config)
	if err != nil {
		ShowErrorDialog(a.ctx, "Грешка при запазване на файла", err.Error())
		return ""
	}

	return "Успешно запазено"
}

func (a *App) SaveSettingsConfig(c Settings) string {
	a.config.Settings = c
	err := SaveConfigToFile(a.config)
	if err != nil {
		ShowErrorDialog(a.ctx, "Грешка при запазване на файла", err.Error())
	}

	return "Успешно запазено"
}

func (a *App) SaveTaxesConfig(t TaxesConfig) string {
	isValid, errorMessage := t.Validate()
	if !isValid {
		ShowWarningDialog(a.ctx, "Грешна конфигурация", errorMessage)
		return ""
	}

	a.config.TaxesConfig = t
	err := SaveConfigToFile(a.config)
	if err != nil {
		ShowErrorDialog(a.ctx, "Грешка при запазване на файла", err.Error())
		return ""
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

func (a *App) LoadAlerts() string {
	personalData := a.LoadUserConfig()
	if !personalData.isPopulated() {
		return "Попълнете личните си данни за да генерирате декларации. Проверете \"Общо заболяване и майчинство\" в Настройките."
	}
	return ""
}

func (a *App) SaveIncomeForm(f IncomeForm) string {

	existingForm, err := GetDataFromFileForMonth(int(f.Month), int(f.Year))
	if err != nil {
		return err.Error()
	}

	if existingForm.Month > 0 {
		ShowWarningDialog(a.ctx, "", fmt.Sprintf("Данните за %d/%d вече съществуват, изтрийте ги за да въведете нови", f.Month, f.Year))
		return ""
	}

	f.TaxesConfig = a.LoadTaxesConfig()
	f.Settings = a.LoadSettingsConfig()

	if f.WorkDaysSickLeave > f.WorkDaysTotal {
		ShowWarningDialog(a.ctx, "", "Дните в болничен не могат да надвишават общите работни дни")
		return ""
	}

	if f.DayStart > 0 && f.DayEnd > 0 && f.DayStart >= f.DayEnd {
		ShowWarningDialog(a.ctx, "", "Началният ден трябва да бъде преди крайния ден")
		return ""
	}

	if f.MonthIncomeCents == 0 || f.TaxedIncomeCents == 0 || f.Year == 0 {
		ShowWarningDialog(a.ctx, "", "Моля, попълнете всички задължителни полета")
		return ""
	}

	if f.TaxedIncomeCents < f.TaxesConfig.MinInsuranceIncomeCents || f.TaxedIncomeCents > f.TaxesConfig.MaxInsuranceIncomeCents {
		ShowWarningDialog(
			a.ctx,
			"",
			fmt.Sprintf(
				"Осигурителният доход трябва да бъде между %.2f и %.2f EUR",
				float64(f.TaxesConfig.MinInsuranceIncomeCents)/MoneyDivider,
				float64(f.TaxesConfig.MaxInsuranceIncomeCents)/MoneyDivider,
			),
		)
	}

	// TODO: validation

	// Calculate approximate taxes and social security, save them together with the form
	CalculateSocialSecurity(&f)
	f.TaxesToPayCents = CalculateTaxForMonth(f)

	err = SaveDataToFile(f)
	if err != nil {
		return err.Error()
	}

	// TODO: link to entered data
	return "Успешно запазено"
}

func (a *App) UpdateForm(f IncomeForm) string {
	// TODO: validation

	existingForm, err := GetDataFromFileForMonth(int(f.Month), int(f.Year))
	if err != nil {
		return err.Error()
	}

	if existingForm.Month > 0 {
		// Allow updating only some values
		existingForm.TaxesReallyPaidCents = f.TaxesReallyPaidCents
		existingForm.SocialSecurityReallyPaidParts = f.SocialSecurityReallyPaidParts
		existingForm.SocialSecurityReallyPaidCents = f.SocialSecurityReallyPaidParts.PensionPartOneCents + f.SocialSecurityReallyPaidParts.PensionPartTwoCents + f.SocialSecurityReallyPaidParts.HealthInsuranceCents
	} else {
		// TODO: return err object? or open error dialog?
		return fmt.Sprintf("Данните за месец %d не са намерени", existingForm.Month)
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
	question := fmt.Sprintf(
		"Сигурни ли сте, че искате да изтриете данните за %d/%d? Изтритите данни не се възстановяват!",
		month,
		year,
	)
	answer := ShowQuestionDialog(a.ctx, "Потвърдете действието", question, "")

	if answer == "No" {
		return ""
	}

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

	res += "\nСлед плащането в НАП попълнете платени осигуровки в Въведени данни за Декларацията 6"

	return res
}

func (a *App) PreviewDeclarationSix(year int) SocialSecurityParts {

	// TODO: check that personal data is entered!

	sums := SocialSecurityParts{}
	rows, err := GetDataFromFileForYear(year)
	if err != nil {
		// Error is logged in GetDataFromFileForYear
		return sums
	}

	for _, f := range rows {
		paidSums := f.SocialSecurityReallyPaidParts
		sums.PensionPartOneCents += paidSums.PensionPartOneCents
		sums.PensionPartTwoCents += paidSums.PensionPartTwoCents
		sums.HealthInsuranceCents += paidSums.HealthInsuranceCents
	}

	return sums
}

func (a *App) GenerateDeclarationSix(year int, sums SocialSecurityParts) string {
	var res string

	// TODO: check that personal data is entered!
	personalData := a.LoadUserConfig()

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

	// TODO: MessageDialog

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
	row.TaxesReallyPaidCents = int64(amount * MoneyDivider)
	return a.UpdateForm(row)
}
