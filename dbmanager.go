package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"slices"
)

// We don't have a real database, so we just store the data in the file for now.

func AddDataToFile(a *App, f IncomeForm) error {
	allData, err := GetIncomeData(a)
	if err != nil {
		return err
	}

	allData = append(allData, f)

	err = SaveIncomeData(a, allData)

	return err
}

func GetDataFromFileForMonth(a *App, month int, year int) (IncomeForm, error) {
	rows, err := GetIncomeData(a)
	if err != nil {
		log.Fatal(err)
		return IncomeForm{}, err
	}

	for _, f := range rows {
		if f.Month == int16(month) && f.Year == int16(year) {
			return f, nil
		}
	}

	return IncomeForm{}, nil
}

func GetDataFromFileForYear(a *App, year int) ([]IncomeForm, error) {
	row, err := GetIncomeData(a)
	if err != nil {
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

func GetDataFromFileForQuarter(a *App, quarter int, year int, result *CalculatedTax) ([]IncomeForm, error) {
	row, err := GetIncomeData(a)
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

func DeleteDataFromFile(a *App, month int, year int) error {
	allData, err := GetIncomeData(a)
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
	err = SaveIncomeData(a, allData)

	return err
}

func GetIncomeData(a *App) ([]IncomeForm, error) {
	var rows []IncomeForm

	data, err := LoadData(a.cache, DataFile, DataKey)
	if err != nil {
		return rows, err
	}

	err = json.Unmarshal(data, &rows)
	if err != nil {
		return rows, err
	}

	rows = sortRows(rows)

	return rows, nil
}

func SaveIncomeData(a *App, rows []IncomeForm) error {
	rows = sortRows(rows)
	incomeData, err := json.Marshal(rows)
	if err != nil {
		return err
	}

	return SaveData(
		a.cache,
		incomeData,
		DataFile,
		DataKey,
	)
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
