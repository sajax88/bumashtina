package main

import (
	"log"
	"os"
	"path/filepath"
)

// We don't have a real database, so we just store the data in the file for now.

// TODO: save income and paid social security and calculated taxes for each month

func getDataPath() (string, error) {
	homeDir, dirErr := os.UserHomeDir()

	if dirErr != nil {
		log.Fatal(dirErr)
		return "", dirErr
	}

	dataPath := filepath.Join(homeDir, "bumashtina", "data.json")
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return dataPath, nil
}

func SaveDataToFile() error {
	dataPath, err := getDataPath()
	if err != nil {
		return err
	}

	// TODO
	data := []byte{} // TODO
	// config, err := json.Marshal(c)
	// log.Print("Saving data:", c, string(config))
	// if err != nil {
	// 	log.Print(err)
	// 	return err
	// }

	return os.WriteFile(dataPath, data, 0600)
}
