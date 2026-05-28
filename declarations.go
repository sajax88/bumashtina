package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func firstSymbol(s string) string {
	r := []rune(s)
	if len(r) > 1 {
		return string(r[:1])
	}
	return s
}

func MakeDeclarationOne(f IncomeForm, u UserConfig) ([]byte, error) {
	// TODO: VALIDATE days

	f.WorkDaysReal = f.WorkDaysTotal - f.WorkDaysSickLeave

	initials := strings.ToUpper(firstSymbol(u.FirstName) + firstSymbol(u.MiddleName))
	endSymbol := ""

	fields := []string{
		fmt.Sprintf("%d", f.Month),
		fmt.Sprintf("%d", f.Year),
		u.Bulstat,
		u.Egn,
		"0",
		strings.ToUpper(u.LastName),
		initials,
		"12", // Самоосигуряващо се лице
		fmt.Sprintf("%02d", f.DayStart),
		fmt.Sprintf("%02d", f.DayEnd),
		"00",
		"00",
		"00",
		"00",
		"00",
		"00",
		"00",
		"00",
		fmt.Sprintf("%02d00", f.WorkDaysTotal),
		fmt.Sprintf("%02d", f.WorkDaysReal),
		fmt.Sprintf("%02d", f.WorkDaysSickLeave),
		"00",
		"00",
		"00",
		"00",
		"000",
		"000",
		"0",
		"",
		"0",
		"0000",
		"0.00",
		"0.00",
		fmt.Sprintf("%.2f", (float32(f.TaxedIncomeCents) / float32(MONEY_DIVIDER))),
		"0.00",
		fmt.Sprintf("%.2f", f.TaxesConfig.PensionPercentagePartOne),
		"0.00",
		fmt.Sprintf("%.2f", f.TaxesConfig.HealthInsurancePercentage),
		"0.00",
		"0.00",
		"0.00",
		"0.00",
		fmt.Sprintf("%.2f", f.TaxesConfig.PensionPercentagePartTwo),
		"0.00",
		"0.00",
		"0.00",
		"0.00",
		"0.00",
		"0.00",
		"0.00",
		"0",
		"0",
		"NRAD12008",
	}

	withQuotes := []int{2, 3, 5, 6, 18, 28, 30}
	result := joinFields(fields, ",", withQuotes) + "\r\n" + endSymbol

	result, err := toWindows1251(result)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func MakeDeclarationSix(year int, u UserConfig) ([]byte, error) {
	// TODO: VALIDATE year

	endSymbol := ""
	// TODO:  + u.MiddleName + " " ?
	fullName := strings.ToUpper(u.FirstName + " " + u.LastName)

	fields := []string{
		u.Egn,
		fullName,
		"0",
		strings.ToUpper(u.Email),
		u.Phone,
		u.Phone,
		"13",
		fmt.Sprintf("%d", year),
		"",
		"0",
		"",
		"1500.50", // TODO - which field is it, re-check the fields in NAP program
		"",
		"1200.34", // TODO - which field is it, re-check the fields in NAP program
		"",
		"1200.00", // TODO - which field is it, re-check the fields in NAP program
	}
	fields = append(fields, make([]string, 42)...)
	fields = append(fields, "NRAD62007")

	withQuotes := []int{0, 1, 3, 4, 5}
	result := joinFields(fields, ",", withQuotes) + "\r\n" + endSymbol

	result, err := toWindows1251(result)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func joinFields(fields []string, separator string, withQuotes []int) string {
	result := ""
	for i, field := range fields {
		if slices.Contains(withQuotes, i) || i == len(fields)-1 {
			field = `"` + field + `"`
		}
		result += field
		if i < len(fields)-1 {
			result += separator
		}
	}
	return result
}

func toWindows1251(utf8String string) (string, error) {
	reader := transform.NewReader(bytes.NewBufferString(utf8String), charmap.Windows1251.NewEncoder())
	encodedBytes, err := io.ReadAll(reader)
	return string(encodedBytes), err
}

func fromWindows1251(bytes []byte) (string, error) {
	// TODO: do we need this?
	return "", nil
	// decoder := charmap.Windows1251.NewDecoder()
	// decodedReader := transform.NewReader(bytes, decoder)
	// decodedString, err := io.ReadAll(decodedReader)
	// return string(decodedString), err
}
