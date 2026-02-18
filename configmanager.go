package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

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
	notEmpty := c.FirstName != "" && c.LastName != "" && c.Egn != "" && c.Bulstat != "" && c.YearOfBirth != ""
	numbersOnly := isDigitsOnly(c.Egn) && isDigitsOnly(c.Bulstat) && isDigitsOnly(c.YearOfBirth)
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

	configPath := filepath.Join(configDir, "bumashtina.conf")

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

func GetHardcodedConfig() {
	// •    19,8% за фонд „Пенсии“ за родените преди 1 януари 1960 г., както и лицата по чл. 4б от КСО;
	// •    14,8% за фонд „Пенсии“ за родените след 31 декември 1959 г.; - DEPENDING ON SETTINGS
	// •    3,5% за фонд „Общо заболяване и майчинство“* - OPTIONAL, IN SETTINGS
	// •    5% за ДЗПО – Универсален пенсионен фонд за родените след 31 декември 1959 г.; - DEPENDING ON SETTINGS
	// •    8% за здравно осигуряване.
	// Общо 0.278

	// Мин.осиг.доход: 550,66 евро
	// Макс.осиг.доход: 2111.64 евро
	// признатите разходи 25%
	// Данък 10%
}
