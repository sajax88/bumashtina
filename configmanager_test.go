package main

import (
	"testing"
)

func TestIsDigitsOnly(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"0", true},
		{"", false},
		{"abc", false},
		{"12a3", false},
		{" ", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isDigitsOnly(tt.input); got != tt.expected {
				t.Errorf("isDigitsOnly(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUserConfigIsPopulated(t *testing.T) {
	tests := []struct {
		name     string
		u        UserConfig
		expected bool
	}{
		{
			"Fully Populated",
			UserConfig{FirstName: "John", LastName: "Doe", Egn: "1234567890", Bulstat: "123456789"},
			true,
		},
		{
			"Missing FirstName",
			UserConfig{LastName: "Doe", Egn: "1234567890", Bulstat: "123456789"},
			false,
		},
		{
			"Missing LastName",
			UserConfig{FirstName: "John", Egn: "1234567890", Bulstat: "123456789"},
			false,
		},
		{
			"Missing Egn",
			UserConfig{FirstName: "John", LastName: "Doe", Bulstat: "123456789"},
			false,
		},
		{
			"Missing Bulstat",
			UserConfig{FirstName: "John", LastName: "Doe", Egn: "1234567890"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.isPopulated(); got != tt.expected {
				t.Errorf("isPopulated() = %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestUserConfigValidate(t *testing.T) {
	tests := []struct {
		name       string
		u          UserConfig
		wantOk     bool
		wantErrMsg string
	}{
		{
			"Valid Config",
			UserConfig{FirstName: "Тодор", LastName: "Колев", Egn: "1234567890", Bulstat: "123456789"},
			true,
			"",
		},
		{
			"Not Populated",
			UserConfig{FirstName: "A"},
			false,
			"Попълнете всички полета",
		},
		{
			"Invalid EGN (too short)",
			UserConfig{FirstName: "A", LastName: "B", Egn: "123", Bulstat: "123456789"},
			false,
			"Невалиден ЕГН",
		},
		{
			"Invalid EGN (chars)",
			UserConfig{FirstName: "A", LastName: "B", Egn: "123456789a", Bulstat: "123456789"},
			false,
			"Невалиден ЕГН",
		},
		{
			"Invalid Bulstat (too long)",
			UserConfig{FirstName: "A", LastName: "B", Egn: "1234567890", Bulstat: "12345678901"},
			false,
			"Невалиден Булстат",
		},
		{
			"Invalid Bulstat (chars)",
			UserConfig{FirstName: "A", LastName: "B", Egn: "1234567890", Bulstat: "12345678a"},
			false,
			"Невалиден Булстат",
		},
		{
			"Invalid Phone",
			UserConfig{FirstName: "A", LastName: "B", Egn: "1234567890", Bulstat: "123456789", Phone: "abc"},
			false,
			"Невалиден телефон",
		},
		{
			"Valid Email",
			UserConfig{FirstName: "A", LastName: "B", Egn: "1234567890", Bulstat: "123456789", Email: "test@example.com"},
			true,
			"",
		},
		{
			"Invalid Email",
			UserConfig{FirstName: "A", LastName: "B", Egn: "1234567890", Bulstat: "123456789", Email: "invalid-email"},
			false,
			"Невалиден email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotMsg := tt.u.Validate()
			if gotOk != tt.wantOk {
				t.Errorf("Validate() ok = %v; want %v", gotOk, tt.wantOk)
			}
			if gotMsg != tt.wantErrMsg {
				t.Errorf("Validate() msg = %q; want %q", gotMsg, tt.wantErrMsg)
			}
		})
	}
}

func TestTaxesConfigValidate(t *testing.T) {
	validTaxes := GetDefaultTaxesConfig()

	tests := []struct {
		name       string
		t          TaxesConfig
		wantOk     bool
		wantErrMsg string
	}{
		{"Valid", validTaxes, true, ""},
		{
			"Min Income <= 0",
			func() TaxesConfig { v := validTaxes; v.MinInsuranceIncomeCents = 0; return v }(),
			false,
			"Невалиден минимален осигурителен доход",
		},
		{
			"Max Income <= 0",
			func() TaxesConfig { v := validTaxes; v.MaxInsuranceIncomeCents = 0; return v }(),
			false,
			"Невалиден максимален осигурителен доход",
		},
		{
			"Min >= Max",
			func() TaxesConfig {
				v := validTaxes
				v.MinInsuranceIncomeCents = 100
				v.MaxInsuranceIncomeCents = 100
				return v
			}(),
			false,
			"Минималният осигурителен доход трябва да е по-малък от максималния",
		},
		{
			"Expenses <= 0",
			func() TaxesConfig { v := validTaxes; v.ExpensesPercentage = 0; return v }(),
			false,
			"Невалиден процент признати разходи",
		},
		{
			"Tax <= 0",
			func() TaxesConfig { v := validTaxes; v.TaxPercentage = 0; return v }(),
			false,
			"Невалидна данъчна ставка",
		},
		{
			"Health Insurance <= 0",
			func() TaxesConfig { v := validTaxes; v.HealthInsurancePercentage = 0; return v }(),
			false,
			"Невалиден процент за Здравно осигуряване",
		},
		{
			"Pregnancy Insurance <= 0",
			func() TaxesConfig { v := validTaxes; v.PregnancyInsurancePercentage = 0; return v }(),
			false,
			"Невалиден процент за фонд Общо заболяване и майчинство",
		},
		{
			"Pension part 1 <= 0",
			func() TaxesConfig { v := validTaxes; v.PensionPercentagePartOne = 0; return v }(),
			false,
			"Невалиден процент за ДОО",
		},
		{
			"Pension part 2 <= 0",
			func() TaxesConfig { v := validTaxes; v.PensionPercentagePartTwo = 0; return v }(),
			false,
			"Невалиден процент за ДЗПО",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotMsg := tt.t.Validate()
			if gotOk != tt.wantOk {
				t.Errorf("Validate() ok = %v; want %v", gotOk, tt.wantOk)
			}
			if gotMsg != tt.wantErrMsg {
				t.Errorf("Validate() msg = %q; want %q", gotMsg, tt.wantErrMsg)
			}
		})
	}
}

func TestGetDefaultTaxesConfig(t *testing.T) {
	cfg := GetDefaultTaxesConfig()
	if cfg.MinInsuranceIncomeCents == 0 {
		t.Error("Expected non-zero MinInsuranceIncomeCents")
	}
	if cfg.MaxInsuranceIncomeCents == 0 {
		t.Error("Expected non-zero MaxInsuranceIncomeCents")
	}
	if cfg.ExpensesPercentage == 0 {
		t.Error("Expected non-zero ExpensesPercentage")
	}
	if cfg.TaxPercentage == 0 {
		t.Error("Expected non-zero TaxPercentage")
	}
	if cfg.HealthInsurancePercentage == 0 {
		t.Error("Expected non-zero HealthInsurancePercentage")
	}
	if cfg.PregnancyInsurancePercentage == 0 {
		t.Error("Expected non-zero PregnancyInsurancePercentage")
	}
	if cfg.PensionPercentagePartOne == 0 {
		t.Error("Expected non-zero PensionPercentagePartOne")
	}
	if cfg.PensionPercentagePartTwo == 0 {
		t.Error("Expected non-zero PensionPercentagePartTwo")
	}
}

func TestGetLabelsForTaxesConfig(t *testing.T) {
	labels := GetLabelsForTaxesConfig()
	if len(labels) != 8 {
		t.Errorf("Expected 8 labels, got %d", len(labels))
	}
}

func TestLoadConfig(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("Empty Config", func(t *testing.T) {
		cfg, err := LoadConfig(app)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.TaxesConfig != GetDefaultTaxesConfig() {
			t.Error("Expected default taxes config for empty storage")
		}
	})

	t.Run("Valid Saved Config", func(t *testing.T) {
		expected := Config{
			User:        UserConfig{FirstName: "John"},
			TaxesConfig: GetDefaultTaxesConfig(),
		}
		err := SaveConfig(app, expected)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		cfg, err := LoadConfig(app)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.User.FirstName != "John" {
			t.Errorf("Expected FirstName John, got %s", cfg.User.FirstName)
		}
	})

	t.Run("Invalid Taxes in Config", func(t *testing.T) {
		invalidCfg := Config{
			TaxesConfig: TaxesConfig{MinInsuranceIncomeCents: -1},
		}
		// SaveConfig does not perform validation, so we can save invalid config here
		err := SaveConfig(app, invalidCfg)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		cfg, err := LoadConfig(app)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.TaxesConfig != GetDefaultTaxesConfig() {
			t.Error("Expected default taxes config when saved config is invalid")
		}
	})
}

func TestSaveConfig(t *testing.T) {
	app, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := Config{
		User: UserConfig{FirstName: "Jane"},
	}

	err := SaveConfig(app, cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify by loading
	loaded, _ := LoadConfig(app)
	if loaded.User.FirstName != "Jane" {
		t.Errorf("Expected Jane, got %s", loaded.User.FirstName)
	}
}
