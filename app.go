package main

import (
	"context"
	"fmt"
	"log"
	"os"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx   context.Context
	cache *lru.Cache[string, []byte]
}

func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
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

	c, err := LoadConfig(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	c.User = u
	err = SaveConfig(a, c)
	if err != nil {
		ShowErrorDialog(a.ctx, "Грешка при запазване на файла", err.Error())
		return ""
	}

	return "Успешно запазено"
}

func (a *App) SaveSettingsConfig(s Settings) string {
	c, err := LoadConfig(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	c.Settings = s
	err = SaveConfig(a, c)
	if err != nil {
		ShowErrorDialog(a.ctx, "Грешка при запазване на файла", err.Error())
		return ""
	}

	return "Успешно запазено"
}

// TODO: same as in Users config
func (a *App) SaveTaxesConfig(t TaxesConfig) string {
	isValid, errorMessage := t.Validate()
	if !isValid {
		ShowWarningDialog(a.ctx, "Грешна конфигурация", errorMessage)
		return ""
	}

	c, err := LoadConfig(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	c.TaxesConfig = t
	err = SaveConfig(a, c)
	if err != nil {
		ShowErrorDialog(a.ctx, "Грешка при запазване на файла", err.Error())
		return ""
	}

	return "Успешно запазено"
}

func (a *App) LoadUserConfig() UserConfig {
	c, err := LoadConfig(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return UserConfig{}
	}
	return c.User
}

func (a *App) LoadSettingsConfig() Settings {
	c, err := LoadConfig(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return Settings{}
	}
	return c.Settings
}

func (a *App) LoadTaxesConfig() TaxesConfig {
	c, err := LoadConfig(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return TaxesConfig{}
	}
	return c.TaxesConfig
}

func (a *App) LoadTaxesConfigLabels() []string {
	return GetLabelsForTaxesConfig()
}

func (a *App) LoadAllIncomeData() []IncomeForm {
	rows, err := GetIncomeData(a)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		log.Fatal(err)
	}
	return rows
}

func (a *App) LoadIncomeDataForMonth(month int, year int) IncomeForm {
	row, err := GetDataFromFileForMonth(a, month, year)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
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
	existingForm, err := GetDataFromFileForMonth(a, int(f.Month), int(f.Year))
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return err.Error()
	}

	if existingForm.Month > 0 {
		ShowWarningDialog(a.ctx, "", fmt.Sprintf("Данните за %d/%d вече съществуват, изтрийте ги, за да въведете нови", f.Month, f.Year))
		return ""
	}

	f.TaxesConfig = a.LoadTaxesConfig()
	f.Settings = a.LoadSettingsConfig()

	isValid, errorMessage := f.Validate()
	if !isValid {
		ShowWarningDialog(a.ctx, "", errorMessage)
		return ""
	}

	// Calculate approximate taxes and social security, save them together with the form
	CalculateSocialSecurity(&f)
	CalculateTaxForMonth(&f)

	err = AddDataToFile(a, f)
	if err != nil {
		return err.Error()
	}

	return "Успешно запазено"
}

func (a *App) UpdateForm(f IncomeForm) string {
	existingForm, err := GetDataFromFileForMonth(a, int(f.Month), int(f.Year))
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	if existingForm.Month > 0 {
		// Allow updating only some values
		existingForm.TaxesReallyPaidCents = f.TaxesReallyPaidCents
		existingForm.SocialSecurityReallyPaidParts = f.SocialSecurityReallyPaidParts
		existingForm.SocialSecurityReallyPaidCents = f.SocialSecurityReallyPaidParts.PensionPartOneCents + f.SocialSecurityReallyPaidParts.PensionPartTwoCents + f.SocialSecurityReallyPaidParts.HealthInsuranceCents
	} else {
		ShowErrorDialog(a.ctx, "", fmt.Sprintf("Данните за месец %d не са намерени", existingForm.Month))
		return ""
	}

	// Remove previous record and save the new one
	err = DeleteDataFromFile(a, int(f.Month), int(f.Year))
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	err = AddDataToFile(a, existingForm)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
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

	err := DeleteDataFromFile(a, month, year)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}
	return "Успешно изтрито"
}

func (a *App) GenerateDeclarationOne(month int, year int) string {
	var res string
	incomeForm, err := GetDataFromFileForMonth(a, month, year)
	if err != nil {
		return err.Error()
	}

	personalData := a.LoadUserConfig()
	settings := a.LoadSettingsConfig()

	if !personalData.isPopulated() {
		ShowWarningDialog(a.ctx, "", "Попълнете личните си данни за да генерирате декларацията")
		return ""
	}

	content, err := MakeDeclarationOne(incomeForm, personalData, settings)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	res, err = a.SaveDeclarationToFile(content, fmt.Sprintf("Declaration_1_%d_%02d", year, month))
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}
	if res == "" {
		return res
	}

	res += "\nСлед плащането в НАП попълнете платени осигуровки в Въведени данни за Декларацията 6"

	return res
}

func (a *App) PreviewDeclarationSix(year int) SocialSecurityParts {
	sums := SocialSecurityParts{}
	rows, err := GetDataFromFileForYear(a, year)
	if err != nil {
		// Error is logged in GetDataFromFileForYear
		return sums
	}

	personalData := a.LoadUserConfig()
	if !personalData.isPopulated() {
		ShowWarningDialog(a.ctx, "", "Попълнете личните си данни за да генерирате декларацията")
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

	personalData := a.LoadUserConfig()
	if !personalData.isPopulated() {
		ShowWarningDialog(a.ctx, "", "Попълнете личните си данни за да генерирате декларацията")
		return ""
	}

	if year < MinYear {
		ShowErrorDialog(a.ctx, "", "Невалидна година")
		return ""
	}

	content, err := MakeDeclarationSix(year, personalData, sums)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}

	res, err = a.SaveDeclarationToFile(content, fmt.Sprintf("Declaration_6_%d", year))
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return ""
	}
	if res == "" {
		return res
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
	if saveFilePath == "" {
		// Can happen when user cancels the dialog
		return "", nil
	}

	saveFilePath, err = SaveToFile(saveFilePath, content)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Успешно, файлът беше записан в %s", saveFilePath), nil
}

func (a *App) CalculateTaxForQuarter(quarter int, year int) CalculatedTax {
	result := CalculatedTax{
		TotalIncomeCents: 0,
		TaxCents:         0,
		Quarter:          quarter,
		Year:             year,
	}

	if quarter > 4 || quarter < 1 {
		ShowErrorDialog(a.ctx, "", "Невалидно тримесечие")
		return result
	}

	if year < MinYear {
		ShowErrorDialog(a.ctx, "", "Невалидна година")
		return result
	}

	rows, err := GetDataFromFileForQuarter(a, quarter, year, &result)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return result
	}

	if len(rows) == 0 {
		return result
	}

	result.TotalIncomeCents, err = CalculateIncomeForThreeMonths(rows)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return result
	}

	err = CalculateAdvanceTaxForThreeMonths(rows, &result)
	if err != nil {
		ShowErrorDialog(a.ctx, "", err.Error())
		return result
	}

	result.TaxesConfig = rows[0].TaxesConfig
	return result
}

func (a *App) SavePaidTaxForQuarter(quarter int, year int, amount float32) string {
	// Find the last month in the quarter and save the paid tax there
	lastMonth := (quarter-1)*3 + 3
	row, err := GetDataFromFileForMonth(a, lastMonth, year)
	if err != nil {
		ShowWarningDialog(a.ctx, "", err.Error())
		return ""
	}
	if row.Month != int16(lastMonth) {
		ShowWarningDialog(a.ctx, "", "Въведете данните за трите месеца, за да запазите платения данък за тримесечието. Ако сте прекъснали дейност, въведете нулев доход и нулев осигурителен доход.")
		return ""
	}
	row.TaxesReallyPaidCents = int64(amount * MoneyDivider)
	return a.UpdateForm(row)
}
