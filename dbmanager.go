package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"slices"
)

// We don't have a real database, so we just store the data in the file for now.

// TODO: save income and paid social security and calculated taxes for each month

func getDataPath() (string, error) {
	homeDir, dirErr := os.UserConfigDir()

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

	// TODO: insert or update
	allData, err := GetAllDataFromFile()
	if err != nil {
		return err
	}

	// TODO: sort?
	allData = append(allData, f)

	// TODO
	data, err := json.Marshal(allData)

	if err != nil {
		log.Print(err)
		return err
	}

	return os.WriteFile(dataPath, data, 0600)
}

func GetDataFromFile(month int, year int) (IncomeForm, error) {
	row, err := GetAllDataFromFile()
	if err != nil {
		log.Fatal(err)
		return IncomeForm{}, err
	}

	for _, f := range row {
		if f.Month == int16(month) && f.Year == int16(year) {
			return f, nil
		}
	}

	return IncomeForm{}, nil
}

func DeleteDataFromFile(month int, year int) error {
	dataPath, err := getDataPath()
	if err != nil {
		return err
	}

	allData, err := GetAllDataFromFile()
	if err != nil {
		log.Fatal(err)
		return err
	}

	indexToDelete := -1
	for i, f := range allData {
		if f.Month == int16(month) && f.Year == int16(year) {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return nil
	}

	allData = slices.Delete(allData, indexToDelete, indexToDelete+1)
	data, err := json.Marshal(allData)

	if err != nil {
		log.Print(err)
		return err
	}

	return os.WriteFile(dataPath, data, 0600)
}

func GetAllDataFromFile() ([]IncomeForm, error) {
	dataPath, err := getDataPath()
	if err != nil {
		log.Fatal(err)
		return []IncomeForm{}, err
	}

	savedData, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatal(err)
		return []IncomeForm{}, err
	}

	if len(savedData) > 0 {
		var rows []IncomeForm
		json.Unmarshal(savedData, &rows)
		return rows, nil
	}

	return []IncomeForm{}, nil
}
