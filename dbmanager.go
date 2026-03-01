package main

import (
	"encoding/json"
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

	dataPath := filepath.Join(homeDir, "bumashtina", "data", "data.json")
	err := os.MkdirAll(filepath.Dir(dataPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return dataPath, nil
}

func SaveDataToFile(f IncomeForm) error {
	dataPath, err := getDataPath()
	if err != nil {
		return err
	}

	// TODO: load all data, insert or update, save new data

	// TODO
	data, err := json.Marshal(f)
	log.Print("Saving data:", string(data))
	if err != nil {
		log.Print(err)
		return err
	}

	return os.WriteFile(dataPath, data, 0600)
}
