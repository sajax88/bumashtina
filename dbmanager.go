package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

// We don't have a real database, so we just store the data in the file for now.

func getDataPath() (string, error) {
	homeDir, dirErr := os.UserConfigDir()

	if dirErr != nil {
		log.Fatal(dirErr)
	}

	dataPath := filepath.Join(homeDir, "bumashtina", "data", "data.json")
	err := os.MkdirAll(filepath.Dir(dataPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return dataPath, nil
}

func AddDataToFile(f IncomeForm) error {
	dataPath, err := getDataPath()
	if err != nil {
		return err
	}

	allData, err := GetAllDataFromFile()
	if err != nil {
		return err
	}

	allData = append(allData, f)
	allData = sortRows(allData)

	data, err := json.Marshal(allData)
	if err != nil {
		return err
	}

	_, err = SaveToFile(dataPath, data)

	return err
}

func GetDataFromFileForMonth(month int, year int) (IncomeForm, error) {
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

func GetDataFromFileForYear(year int) ([]IncomeForm, error) {
	row, err := GetAllDataFromFile()
	if err != nil {
		log.Fatal(err)
		return []IncomeForm{}, err
	}

	var rows []IncomeForm
	for _, f := range row {
		if f.Year == int16(year) {
			rows = append(rows, f)
		}
	}

	return rows, nil
}

func GetDataFromFileForQuarter(quarter int, year int, result *CalculatedTax) ([]IncomeForm, error) {
	row, err := GetAllDataFromFile()
	if err != nil {
		log.Fatal(err)
		return []IncomeForm{}, err
	}

	monthStart := (quarter-1)*3 + 1
	monthEnd := monthStart + 2

	var rows []IncomeForm
	for _, f := range row {
		if (f.Month >= int16(monthStart) && f.Month <= int16(monthEnd)) && f.Year == int16(year) {
			rows = append(rows, f)
		}
	}

	result.MonthStart = monthStart
	result.MonthEnd = monthEnd

	return rows, nil
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
		return err
	}

	_, err = SaveToFile(dataPath, data)
	return err
}

func GetAllDataFromFile() ([]IncomeForm, error) {
	// TODO: cache in app once it was read? Refresh the cache when written!
	// Caching: https://v3.wails.io/guides/performance/#caching

	dataPath, err := getDataPath()
	if err != nil {
		log.Fatal(err)
	}

	savedData, err := os.ReadFile(dataPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []IncomeForm{}, nil
		}
		log.Fatal(err)
	}

	if len(savedData) > 0 {
		var rows []IncomeForm
		err := json.Unmarshal(savedData, &rows)
		if err != nil {
			return nil, err
		}
		rows = sortRows(rows)

		return rows, nil
	}

	return []IncomeForm{}, nil
}

// Sort the rows by month and year
func sortRows(rows []IncomeForm) []IncomeForm {
	compareByMonthAndYear := func(a, b IncomeForm) int {
		return -cmp.Compare(
			fmt.Sprintf("%d%02d", a.Year, a.Month),
			fmt.Sprintf("%d%02d", b.Year, b.Month),
		)
	}

	slices.SortFunc(rows, compareByMonthAndYear)

	return rows
}
