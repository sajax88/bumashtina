package main

import (
	"errors"
	"os"
)

// ReadFromFile returns empty data if the file does not exist
func ReadFromFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		}
		return []byte{}, err
	}

	return data, nil
}

func SaveToFile(filepath string, content []byte) (string, error) {
	err := os.WriteFile(filepath, content, 0644)
	if err != nil {
		return "", err
	}

	return filepath, nil
}
