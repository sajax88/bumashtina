package main

import (
	"errors"
	"os"
)

// ReadFromFile returns empty data if the file does not exist
func ReadFromFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

func SaveToFile(filePath string, content []byte) error {
	err := os.WriteFile(filePath, content, 0600)
	if err != nil {
		return err
	}

	return nil
}
