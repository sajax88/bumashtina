package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	lru "github.com/hashicorp/golang-lru/v2"
)

const ConfigKey = "CONFIG"
const DataKey = "DATA"

const ConfigFile = "config.json"
const DataFile = "data.json"

// TODO: do we need this?
// StorageConfig can be used for saving and loading config and data
type StorageConfig struct {
	CacheKey string
	FileName string
	Obj      interface{}
}

// TODO: replace all SaveToFile and ReadFromFile with this

// TODO: reuse the same methods for saving and loading data,

func LoadConfig(a *App) (Config, error) {
	var c Config

	data, err := loadData(a.cache, ConfigFile, ConfigKey)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	// Set default taxes config
	if c.TaxesConfig.TaxPercentage == 0 {
		c.TaxesConfig = GetDefaultTaxesConfig()
	}

	return c, err
}

func SaveConfig(a *App, c Config) error {
	config, err := json.Marshal(c) // TODO: all marshalling and unmarshalling - in file service
	if err != nil {
		return err
	}

	return saveData(
		a.cache,
		config,
		ConfigFile,
		ConfigKey,
	)
}

func loadData(cache *lru.Cache[string, []byte], fileName string, cacheKey string) ([]byte, error) {
	// First try to get from cache
	cachedData, inCache := cache.Get(cacheKey)
	if inCache && len(cachedData) > 0 {
		return cachedData, nil
	}

	// Try to get from file
	filePath, err := getFilePath(fileName)
	if err != nil {
		return nil, err
	}
	savedData, err := ReadFromFile(filePath)
	if err != nil || len(savedData) == 0 {
		return nil, err
	}

	if inCache {
		cache.Remove(cacheKey)
	}
	cache.Add(cacheKey, savedData)

	return savedData, nil
}

func saveData(cache *lru.Cache[string, []byte], data []byte, fileName string, cacheKey string) error {
	filePath, err := getFilePath(fileName)
	if err != nil {
		return err
	}

	_, err = SaveToFile(filePath, data)
	if err != nil {
		return err
	}

	// Update cache
	cache.Remove(cacheKey)
	cache.Add(cacheKey, data)

	return nil
}

func getFilePath(filename string) (string, error) {
	d, dirErr := os.UserConfigDir()
	if dirErr != nil {
		return "", dirErr
	}

	filePath := filepath.Join(d, "bumashtina", "data", filename)
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", err // TODO: remove almost all log.Fatal!
	}

	return filePath, nil
}
