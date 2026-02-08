package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	FirstName  string
	MiddleName string
	LastName   string
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
	var c Config

	configPath, err := getConfigPath()
	if err != nil {
		return c, err
	}

	savedConfig, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		// There is a config file but we couldn't read it
		log.Fatal(err)
		return c, err
	}

	if savedConfig != nil {
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
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, config, 0600)
}
