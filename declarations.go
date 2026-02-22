package main

import (
	"bytes"
	"io"
	"slices"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func MakeDeclarationOne() ([]byte, error) {
	endSymbol := ""
	fiedls := []string{
		"2",       // TODO
		"2026",    // TODO
		"булстат", // TODO
		"егн",     // TODO
		"0",
		"ИМЯ",      // TODO
		"ИНИЦИАЛЫ", // TODO
		"12",       // Self-employed
		"00",
		"00",
		"00",
		"00",
		"00",
		"00",
		"00",
		"00",
		"00", // TODO: check start/end days
		"00",
		"2000", // TODO
		"20",   // TODO
		"00",
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
		"2111.64", // TODO
		"0.00",
		"14.80", // TODO
		"0.00",
		"8.00", // TODO
		"0.00",
		"0.00",
		"0.00",
		"0.00",
		"5.00", // TODO
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

	result := joinFields(fiedls, ",") + "\r\n" + endSymbol

	result, err := toWindows1251(result)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func joinFields(fields []string, separator string) string {
	result := ""
	withQuotes := []int{2, 3, 5, 6, 18, 28, 30}
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
	// TODO
	return "", nil
	// decoder := charmap.Windows1251.NewDecoder()
	// decodedReader := transform.NewReader(bytes, decoder)
	// decodedString, err := io.ReadAll(decodedReader)
	// return string(decodedString), err
}
