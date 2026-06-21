package main

import (
	"encoding/json"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
)

const MoneyDivider = 100

type UserConfig struct {
	FirstName  string
	MiddleName string
	LastName   string
	Egn        string
	Bulstat    string
	Phone      string
	Email      string
}

type Settings struct {
	IsPregnancyInsuranceEnabled bool
}

// TaxesConfig Money is stored with Divider = 100 to avoid floating point issues.
// E.g., 550,66 euro is stored as "55066".
type TaxesConfig struct {
	MinInsuranceIncomeCents      int64
	MaxInsuranceIncomeCents      int64
	ExpensesPercentage           float32 // TODO: to float64?
	TaxPercentage                float32
	HealthInsurancePercentage    float32
	PregnancyInsurancePercentage float32
	PensionPercentagePartOne     float32
	PensionPercentagePartTwo     float32
}

type Config struct {
	User        UserConfig
	Settings    Settings
	TaxesConfig TaxesConfig
}

func (u UserConfig) isPopulated() bool {
	if u.FirstName == "" || u.LastName == "" || u.Egn == "" || u.Bulstat == "" {
		return false
	}
	return true
}

func (u UserConfig) Validate() (bool, string) {
	if !u.isPopulated() {
		return false, "Попълнете всички полета"
	}

	if !isDigitsOnly(u.Egn) || len(u.Egn) != 10 {
		return false, "Невалиден ЕГН"
	}

	if !isDigitsOnly(u.Bulstat) || len(u.Bulstat) < 9 || len(u.Bulstat) > 10 {
		return false, "Невалиден Булстат"
	}

	if len(u.Phone) > 0 && !isDigitsOnly(u.Phone) {
		return false, "Невалиден телефон"
	}

	if len(u.Email) > 0 {
		_, err := mail.ParseAddress(u.Email)
		if err != nil {
			return false, "Невалиден email"
		}
	}

	return true, ""
}

func isDigitsOnly(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func getConfigPath() (string, error) {
	configDir, dirErr := os.UserConfigDir()

	if dirErr != nil {
		log.Fatal(dirErr)
		return "", dirErr
	}

	configPath := filepath.Join(configDir, "bumashtina", "data", "config.json")
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return configPath, nil
}

func LoadConfigFromFile() (Config, error) {
	s := Settings{IsPregnancyInsuranceEnabled: false} // Default
	t := GetDefaultTaxesConfig()
	c := Config{User: UserConfig{}, Settings: s, TaxesConfig: t}

	configPath, err := getConfigPath()
	if err != nil {
		return c, err
	}

	savedConfig, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		// There is a config file but we couldn't read it
		return c, err
	}

	if len(savedConfig) > 0 {
		err = json.Unmarshal(savedConfig, &c)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func SaveConfigToFile(c Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	config, err := json.Marshal(c)
	log.Print("Saving config:", c, string(config))
	if err != nil {
		log.Print(err)
		return err
	}

	return os.WriteFile(configPath, config, 0600)
}

func GetDefaultTaxesConfig() TaxesConfig {
	return TaxesConfig{
		MinInsuranceIncomeCents:      55066,  // Euro cents
		MaxInsuranceIncomeCents:      211164, // Euro cents
		ExpensesPercentage:           25.0,   // Признати разходи
		TaxPercentage:                10.0,   // Данък върху дохода
		HealthInsurancePercentage:    8.0,    // За фонд „Здравно осигуряване“
		PregnancyInsurancePercentage: 3.5,    // За фонд „Общо заболяване и майчинство“ - see Settings.IsPregnancyInsuranceEnabled
		PensionPercentagePartOne:     14.8,   // Фонд "Пенсии", ДОО, използва се в Декларация 1, 6
		PensionPercentagePartTwo:     5.0,    // Фонд "Пенсии", ДЗПО, използва се в Декларация 1, 6
	}
}

func GetLabelsForTaxesConfig() []string {
	return []string{
		"Минимален осигурителен доход, EUR",
		"Максимален осигурителен доход, EUR",
		"Процент признати разходи",
		"Данъчна ставка",
		"Процент за фонд „Здравно осигуряване“",
		"Процент за фонд „Общо заболяване и майчинство“",
		"Процент за фонд „Пенсии“ ДОО – за родените след 31 декември 1959 г.",
		"Процент за ДЗПО – Универсален пенсионен фонд за родените след 31 декември 1959 г.",
	}
}
