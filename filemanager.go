package main

import (
	"os"
	"path/filepath"
)

func SaveDeclaration(dir string, content []byte) (string, error) {
	path1, err := filepath.Abs(dir + "/dat1.txt") // TODO: name
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path1, content, 0644)
	if err != nil {
		return "", err
	}

	return path1, nil
}
