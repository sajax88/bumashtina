package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const MONEY_DIVIDER = 100

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

type Config struct {
	User        UserConfig
	Settings    Settings
	TaxesConfig TaxesConfig
}

func (c UserConfig) IsValid() bool {
	// TODO: only cyrillic, number of digits for EGN and Bulstat, etc.
	notEmpty := c.FirstName != "" && c.LastName != "" && c.Egn != "" && c.Bulstat != ""
	numbersOnly := isDigitsOnly(c.Egn) && isDigitsOnly(c.Bulstat)
	return notEmpty && numbersOnly
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
	s := Settings{IsPregnancyInsuranceEnabled: true} // Default
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
		json.Unmarshal(savedConfig, &c)
	}

	// TODO: if no Taxes Config?

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

// Money is stored with Divider = 100 to avoid floating point issues.
// E.g. 550,66 euro is stored as 55066.
type TaxesConfig struct {
	MinInsuranceIncomeCents      int32
	MaxInsuranceIncomeCents      int32
	ExpensesPercentage           float32
	TaxPercentage                float32
	PensionPercentage            float32
	HealthInsurancePercentage    float32
	PregnancyInsurancePercentage float32
	PensionPercentagePartOne     float32
	PensionPercentagePartTwo     float32
}

func GetDefaultTaxesConfig() TaxesConfig {
	return TaxesConfig{
		MinInsuranceIncomeCents:      55066,  // Euro cents
		MaxInsuranceIncomeCents:      211164, // Euro cents
		ExpensesPercentage:           25.0,   // Признати разходи
		TaxPercentage:                10.0,   // Данък върху дохода
		PensionPercentage:            19.8,   // За фонд „Пенсии“
		HealthInsurancePercentage:    8.0,    // За фонд „Здравно осигуряване“
		PregnancyInsurancePercentage: 3.5,    // За фонд „Майчинство“ - see Settings.IsPregnancyInsuranceEnabled
		PensionPercentagePartOne:     14.8,   // Използва се в Декларация 1, 6
		PensionPercentagePartTwo:     5.0,    // Използва се в Декларация 1, 6
	}
}

func GetLabelsForTaxesConfig() []string {
	return []string{
		"Минимален осигурителен доход, EUR",
		"Максимален осигурителен доход, EUR",
		"Процент признати разходи",
		"Данъчна ставка",
		"Процент за фонд „Пенсии“",
		"Процент за фонд „Здравно осигуряване“",
		"Процент за фонд „Майчинство“",
		"Процент за фонд „Пенсии“ – за родените след 31 декември 1959 г.",
		"Процент за ДЗПО – Универсален пенсионен фонд за родените след 31 декември 1959 г.",
	}
}
