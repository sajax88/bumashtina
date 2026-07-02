package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	lru "github.com/hashicorp/golang-lru/v2"
)

const ConfigKey = "CONFIG"
const DataKey = "DATA"

const ConfigFile = "config.json"
const DataFile = "data.json"

func LoadData(cache *lru.Cache[string, []byte], fileName string, cacheKey string) ([]byte, error) {
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

func SaveData(cache *lru.Cache[string, []byte], data []byte, fileName string, cacheKey string) error {
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
		return "", err
	}

	return filePath, nil
}

func exportIntoZip(zipFilePath string) error {
	configFilePath, err := getFilePath(ConfigFile)
	if err != nil {
		return err
	}
	dataFilePath, err := getFilePath(DataFile)
	if err != nil {
		return err
	}

	archive, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	files := []string{configFilePath, dataFilePath}

	for _, fPath := range files {
		reader, err := os.Open(fPath)
		if err != nil {
			return err
		}
		defer reader.Close()

		fileName := filepath.Base(fPath)
		writer, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}
		if _, err := io.Copy(writer, reader); err != nil {
			return err
		}
		reader.Close()
	}

	return nil
}
