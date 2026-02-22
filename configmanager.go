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
	FirstName   string
	MiddleName  string
	LastName    string
	Egn         string
	Bulstat     string
	YearOfBirth string
}

type Settings struct {
	IsPregnancyInsuranceEnabled bool
}

type Config struct {
	User     UserConfig
	Settings Settings
}

func (c UserConfig) IsValid() bool {
	notEmpty := c.FirstName != "" && c.LastName != "" && c.Egn != "" && c.Bulstat != ""
	numbersOnly := isDigitsOnly(c.Egn) && isDigitsOnly(c.Bulstat) && (c.YearOfBirth == "" || isDigitsOnly(c.YearOfBirth))
	return notEmpty && numbersOnly
}

func isDigitsOnly(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func getConfigPath() (string, error) {
	configDir, dirErr := os.UserHomeDir()

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
	s := Settings{IsPregnancyInsuranceEnabled: true} // TODO: default settings
	c := Config{User: UserConfig{}, Settings: s}

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

// Money is stored with Divider (100) to avoid floating point issues.
// E.g. 550,66 euro is stored as 55066.
type TaxesConfig struct {
	MinInsuranceIncomeCents      int32
	MaxInsuranceIncomeCents      int32
	ExpensesPercentage           float32
	TaxPercentage                float32
	PensionPercentage            float32
	HealthInsurancePercentage    float32
	PregnancyInsurancePercentage float32
}

func GetTaxesConfig() TaxesConfig {
	return TaxesConfig{
		MinInsuranceIncomeCents:      55066,  // 550,66 euro
		MaxInsuranceIncomeCents:      211164, // 2111,64 euro
		ExpensesPercentage:           25.0,   // 25% признати разходи
		TaxPercentage:                10.0,   // 10% данък върху дохода
		PensionPercentage:            19.8,   // 19,8% за фонд „Пенсии“
		HealthInsurancePercentage:    8.0,    // 8% за фонд „Здравно осигуряване“
		PregnancyInsurancePercentage: 3.5,    // 3,5% за фонд „Майчинство“ - see Settings.IsPregnancyInsuranceEnabled
	}
}
