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
	filePath, err := prepareFilePath(fileName)
	if err != nil {
		return nil, err
	}
	savedData, err := ReadFromFile(filePath)
	if err != nil || savedData == nil {
		return nil, err
	}

	dataCopy := make([]byte, len(savedData))
	copy(dataCopy, savedData)
	cache.Add(cacheKey, dataCopy)

	return savedData, nil
}

func SaveData(cache *lru.Cache[string, []byte], data []byte, fileName string, cacheKey string) error {
	filePath, err := prepareFilePath(fileName)
	if err != nil {
		return err
	}

	err = SaveToFile(filePath, data)
	if err != nil {
		return err
	}

	// Update cache
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	cache.Add(cacheKey, dataCopy)

	return nil
}

func prepareFilePath(filename string) (string, error) {
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
	configFilePath, err := prepareFilePath(ConfigFile)
	if err != nil {
		return err
	}
	dataFilePath, err := prepareFilePath(DataFile)
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

		fileName := filepath.Base(fPath)
		writer, err := zipWriter.Create(fileName)
		if err != nil {
			reader.Close()
			return err
		}
		if _, err := io.Copy(writer, reader); err != nil {
			reader.Close()
			return err
		}

		reader.Close()
	}

	return nil
}

func importFromZip(cache *lru.Cache[string, []byte], zipFilePath string) error {
	archive, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		var filePath string
		if f.Name == ConfigFile {
			filePath, err = prepareFilePath(ConfigFile)
		} else if f.Name == DataFile {
			filePath, err = prepareFilePath(DataFile)
		} else {
			// Skip whatever else we have there
			continue
		}
		if err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			dstFile.Close()
			fileInArchive.Close()
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	cache.Purge()

	return nil
}
