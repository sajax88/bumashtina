package main

import (
	"os"
	"path/filepath"
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
)

// setupTestEnv sets up a temporary directory and environment variables for testing.
// It returns the initialized app and the cleanup function
func setupTestEnv(t *testing.T) (*App, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "bumashtina_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %s", err.Error())
	}
	// User configs and data are stored in user config dir, temporarily replace it with temp dir
	oldXdg := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cache, err := lru.New[string, []byte](10)
	if err != nil {
		panic(err.Error())
	}
	app := &App{cache: cache}

	cleanup := func() {
		os.RemoveAll(tmpDir)
		os.Setenv("XDG_CONFIG_HOME", oldXdg)
	}

	return app, cleanup
}

func TestAddDataToEmptyFile(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	data := IncomeForm{Month: 1, Year: 2026, MonthIncomeCents: 10000, TaxedIncomeCents: 5000}

	err := AddDataToFile(app, data)
	if err != nil {
		t.Fatalf("AddDataToFile failed: %v", err)
	}

	rows, err := GetIncomeData(app)
	if err != nil {
		t.Fatalf("GetIncomeData failed: %v", err)
	}

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
	if rows[0].MonthIncomeCents != 10000 {
		t.Errorf("Expected MonthIncomeCents = 10000, got %d", rows[0].MonthIncomeCents)
	}
	if rows[0].TaxedIncomeCents != 5000 {
		t.Errorf("Expected TaxedIncomeCents = 5000, got %d", rows[0].TaxedIncomeCents)
	}

	// Try to add the same data again
	err = AddDataToFile(app, data)
	if err == nil {
		t.Error("Expected error for duplicate data, got nil")
	}
}

func TestGetDataFromFileForMonth(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	_ = AddDataToFile(app, IncomeForm{Month: 1, Year: 2026, MonthIncomeCents: 10000})
	_ = AddDataToFile(app, IncomeForm{Month: 2, Year: 2026, MonthIncomeCents: 15000})

	t.Run("ExistingMonth", func(t *testing.T) {
		f, err := GetDataFromFileForMonth(app, 2, 2026)
		if err != nil {
			t.Errorf("GetDataFromFileForMonth failed: %v", err)
		}
		if f.MonthIncomeCents != 15000 {
			t.Errorf("Expected MonthIncomeCents = 10000, got %d", f.MonthIncomeCents)
		}
	})

	t.Run("NonExistingMonth", func(t *testing.T) {
		f, err := GetDataFromFileForMonth(app, 12, 2026)
		if err != nil {
			t.Errorf("GetDataFromFileForMonth failed: %v", err)
		}
		if f.Month != 0 {
			t.Errorf("Expected empty result, got month %d", f.Month)
		}
	})
}

func TestGetDataFromFileForYear(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	_ = AddDataToFile(app, IncomeForm{Month: 1, Year: 2026})
	_ = AddDataToFile(app, IncomeForm{Month: 2, Year: 2026})
	_ = AddDataToFile(app, IncomeForm{Month: 1, Year: 2027})

	rows, err := GetDataFromFileForYear(app, 2026)
	if err != nil {
		t.Fatalf("GetDataFromFileForYear failed: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("Expected 2 rows for 2026, got %d", len(rows))
	}
}

func TestGetDataFromFileForQuarter(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	_ = AddDataToFile(app, IncomeForm{Month: 1, Year: 2026})
	_ = AddDataToFile(app, IncomeForm{Month: 2, Year: 2026})
	_ = AddDataToFile(app, IncomeForm{Month: 3, Year: 2026})
	_ = AddDataToFile(app, IncomeForm{Month: 4, Year: 2026})

	var result CalculatedTax
	rows, err := GetDataFromFileForQuarter(app, 1, 2026, &result)
	if err != nil {
		t.Fatalf("GetDataFromFileForQuarter failed: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows for Q1, got %d", len(rows))
	}
	if result.MonthStart != 1 || result.MonthEnd != 3 {
		t.Errorf("Expected months 1-3, got %d-%d", result.MonthStart, result.MonthEnd)
	}
}

func TestDeleteDataFromFile(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	_ = AddDataToFile(app, IncomeForm{Month: 1, Year: 2026})
	_ = AddDataToFile(app, IncomeForm{Month: 2, Year: 2026})

	err := DeleteDataFromFile(app, 1, 2026)
	if err != nil {
		t.Fatalf("DeleteDataFromFile failed: %v", err)
	}

	rows, err := GetIncomeData(app)
	if err != nil {
		t.Fatalf("GetIncomeData failed: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("Expected 1 row after deletion, got %d", len(rows))
	}
	if rows[0].Month != 2 {
		t.Errorf("Expected month 2 to remain, got %d", rows[0].Month)
	}
}

func TestSortRows(t *testing.T) {
	rows := []IncomeForm{
		{Month: 1, Year: 2026},
		{Month: 3, Year: 2026},
		{Month: 2, Year: 2026},
		{Month: 1, Year: 2027},
		{Month: 10, Year: 2026},
	}
	sorted := sortRows(rows)

	if len(sorted) != 4 {
		t.Fatalf("Expected 4 rows, got %d", len(sorted))
	}

	// sortRows sorts descending (Year, Month)
	expectedOrder := []struct {
		year  int16
		month int16
	}{
		{2027, 1},
		{2026, 10},
		{2026, 3},
		{2026, 2},
		{2026, 1},
	}

	for i, exp := range expectedOrder {
		if sorted[i].Year != exp.year || sorted[i].Month != exp.month {
			t.Errorf("At index %d: expected %d-%02d, got %d-%02d", i, exp.year, exp.month, sorted[i].Year, sorted[i].Month)
		}
	}
}

func TestGetIncomeDataError(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	// Manually corrupt the data file
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	filePath := filepath.Join(xdgConfig, "bumashtina", "data", DataFile)

	// Create the directory first because prepareFilePath does it when Save/Load is called
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		t.Fatalf("Failed to create data dir: %v", err)
	}

	err = os.WriteFile(filePath, []byte("invalid json"), 0600)
	if err != nil {
		t.Fatalf("Failed to write corrupted file: %v", err)
	}
	app.cache.Purge()

	_, err = GetIncomeData(app)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}
