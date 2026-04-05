package main

import (
	"os"
)

func SaveDeclaration(filepath string, content []byte) (string, error) {
	err := os.WriteFile(filepath, content, 0644)
	if err != nil {
		return "", err
	}

	return filepath, nil
}
