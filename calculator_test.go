package main

import (
	"testing"
)

func TestIncomeForm_Validate(t *testing.T) {
	defaultConfig := GetDefaultTaxesConfig()
	validForm := IncomeForm{
		Month:            1,
		Year:             2026,
		MonthIncomeCents: 100000,
		TaxedIncomeCents: 60000,
		WorkDaysTotal:    22,
		TaxesConfig:      defaultConfig,
	}

	tests := []struct {
		name    string
		form    IncomeForm
		isValid bool
	}{
		{"Valid form", validForm, true},
		{"Invalid WorkDaysTotal negative", func() IncomeForm { f := validForm; f.WorkDaysTotal = -1; return f }(), false},
		{"Invalid WorkDaysTotal too high", func() IncomeForm { f := validForm; f.WorkDaysTotal = 32; return f }(), false},
		{"Sick leave more than total days", func() IncomeForm { f := validForm; f.WorkDaysSickLeave = 23; return f }(), false},
		{"DayStart >= DayEnd", func() IncomeForm { f := validForm; f.DayStart = 10; f.DayEnd = 5; return f }(), false},
		{"Invalid DayStart", func() IncomeForm { f := validForm; f.DayStart = -1; return f }(), false},
		{"Invalid DayEnd", func() IncomeForm { f := validForm; f.DayEnd = 32; return f }(), false},
		{"Negative MonthIncomeCents", func() IncomeForm { f := validForm; f.MonthIncomeCents = -1; return f }(), false},
		{"Invalid Month 0", func() IncomeForm { f := validForm; f.Month = 0; return f }(), false},
		{"Invalid Month 13", func() IncomeForm { f := validForm; f.Month = 13; return f }(), false},
		{"Invalid Year", func() IncomeForm { f := validForm; f.Year = 2025; return f }(), false},
		{"TaxedIncome < MinInsurance", func() IncomeForm { f := validForm; f.TaxedIncomeCents = 50000; return f }(), false},
		{"TaxedIncome > MaxInsurance", func() IncomeForm { f := validForm; f.TaxedIncomeCents = 300000; return f }(), false},
		{"Zero income but taxed income non-zero", func() IncomeForm { f := validForm; f.MonthIncomeCents = 0; f.TaxedIncomeCents = 100; return f }(), false},
		{"Zero income and zero taxed income", func() IncomeForm { f := validForm; f.MonthIncomeCents = 0; f.TaxedIncomeCents = 0; return f }(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := tt.form.Validate()
			if got != tt.isValid {
				t.Errorf("Validate() = %v, %v; want %v", got, msg, tt.isValid)
			}
		})
	}
}

func TestCalculateSocialSecurity(t *testing.T) {
	config := TaxesConfig{
		PensionPercentagePartOne:     14.8,
		PensionPercentagePartTwo:     5.0,
		HealthInsurancePercentage:    8.0,
		PregnancyInsurancePercentage: 3.5,
	}

	tests := []struct {
		name               string
		taxedIncome        int64
		pregnancyEnabled   bool
		expectedPensionOne int64
		expectedPensionTwo int64
		expectedHealth     int64
		expectedTotal      int64
	}{
		{
			name:               "Pregnancy disabled",
			taxedIncome:        100000,
			pregnancyEnabled:   false,
			expectedPensionOne: 14800,
			expectedPensionTwo: 5000,
			expectedHealth:     8000,
			expectedTotal:      27800,
		},
		{
			name:               "Pregnancy enabled",
			taxedIncome:        100000,
			pregnancyEnabled:   true,
			expectedPensionOne: 14800 + 3500,
			expectedPensionTwo: 5000,
			expectedHealth:     8000,
			expectedTotal:      31300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &IncomeForm{
				TaxedIncomeCents: tt.taxedIncome,
				TaxesConfig:      config,
				Settings: Settings{
					IsPregnancyInsuranceEnabled: tt.pregnancyEnabled,
				},
			}
			CalculateSocialSecurity(f)

			if f.SocialSecurityToPayParts.PensionPartOneCents != tt.expectedPensionOne {
				t.Errorf("PensionPartOneCents = %d; want %d", f.SocialSecurityToPayParts.PensionPartOneCents, tt.expectedPensionOne)
			}
			if f.SocialSecurityToPayParts.PensionPartTwoCents != tt.expectedPensionTwo {
				t.Errorf("PensionPartTwoCents = %d; want %d", f.SocialSecurityToPayParts.PensionPartTwoCents, tt.expectedPensionTwo)
			}
			if f.SocialSecurityToPayParts.HealthInsuranceCents != tt.expectedHealth {
				t.Errorf("HealthInsuranceCents = %d; want %d", f.SocialSecurityToPayParts.HealthInsuranceCents, tt.expectedHealth)
			}
			if f.SocialSecurityToPayCents != tt.expectedTotal {
				t.Errorf("SocialSecurityToPayCents = %d; want %d", f.SocialSecurityToPayCents, tt.expectedTotal)
			}
		})
	}
}

func TestCalculateInsuranceFromPercentage(t *testing.T) {
	tests := []struct {
		income   int64
		percent  float64
		expected int64
	}{
		{100100, 5.0, 5005},
		{100000, 14.8, 14800},
	}

	for _, tt := range tests {
		got := calculateInsuranceFromPercentage(tt.income, tt.percent)
		if got != tt.expected {
			t.Errorf("calculateInsuranceFromPercentage(%d, %f) = %d; want %d", tt.income, tt.percent, got, tt.expected)
		}
	}
}

func TestCalculateSocialSecurityAndTaxForMonth(t *testing.T) {
	config := TaxesConfig{
		PensionPercentagePartOne:  10.0,
		PensionPercentagePartTwo:  5.0,
		HealthInsurancePercentage: 5.0,
		ExpensesPercentage:        25.0,
		TaxPercentage:             10.0,
	}

	t.Run("Calculated insurance", func(t *testing.T) {
		f := &IncomeForm{
			MonthIncomeCents: 100000,
			TaxedIncomeCents: 100000,
			TaxesConfig:      config,
		}
		CalculateSocialSecurityAndTaxForMonth(f)
		// Insurance = 10% + 5% + 5% = 20% of 100000 = 20000
		// Expenses = 25% of 100000 = 25000
		// Taxable = 100000 - 25000 - 20000 = 55000
		// Tax = 10% of 55000 = 5500
		if f.SocialSecurityToPayCents != 20000 {
			t.Errorf("SocialSecurityToPayCents = %d; want 20000", f.SocialSecurityToPayCents)
		}
		if f.ExpensesCents != 25000 {
			t.Errorf("ExpensesCents = %d; want 25000", f.ExpensesCents)
		}
		if f.TaxesToPayCents != 5500 {
			t.Errorf("TaxesToPayCents = %d; want 5500", f.TaxesToPayCents)
		}
	})

	t.Run("Really paid insurance", func(t *testing.T) {
		f := &IncomeForm{
			MonthIncomeCents:              100000,
			TaxedIncomeCents:              100000,
			SocialSecurityReallyPaidCents: 30000,
			TaxesConfig:                   config,
		}
		CalculateSocialSecurityAndTaxForMonth(f)
		// Insurance used for tax = 30000
		// Expenses = 25000
		// Taxable = 100000 - 25000 - 30000 = 45000
		// Tax = 10% of 45000 = 4500
		if f.TaxesToPayCents != 4500 {
			t.Errorf("TaxesToPayCents = %d; want 4500", f.TaxesToPayCents)
		}
	})
}

func TestCalculateIncomeForThreeMonths(t *testing.T) {
	forms := []IncomeForm{
		{MonthIncomeCents: 10000},
		{MonthIncomeCents: 20000},
		{MonthIncomeCents: 30000},
	}

	result, err := CalculateIncomeForThreeMonths(forms)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 60000 {
		t.Errorf("Got %d; want 60000", result)
	}

	_, err = CalculateIncomeForThreeMonths(append(forms, IncomeForm{}))
	if err == nil {
		t.Error("Expected error for 4 months, got nil")
	}
}

func TestCalculateAdvanceTaxForThreeMonths(t *testing.T) {
	config := TaxesConfig{
		ExpensesPercentage: 25.0,
		TaxPercentage:      10.0,
	}
	forms := []IncomeForm{
		{MonthIncomeCents: 100000, SocialSecurityToPayCents: 20000, TaxesConfig: config},
		{MonthIncomeCents: 100000, SocialSecurityToPayCents: 20000, TaxesConfig: config},
		{MonthIncomeCents: 100000, SocialSecurityReallyPaidCents: 25000, TaxesConfig: config},
	}

	var result CalculatedTax
	err := CalculateAdvanceTaxForThreeMonths(forms, &result)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Total Income = 300000
	// Paid Insurance = 20000 + 20000 + 25000 = 65000
	// Expenses = 25% of 300000 = 75000
	// Taxable = 300000 - 75000 - 65000 = 160000
	// Tax = 10% of 160000 = 16000

	if result.PaidInsuranceCents != 65000 {
		t.Errorf("PaidInsuranceCents = %d; want 65000", result.PaidInsuranceCents)
	}
	if result.ExpensesCents != 75000 {
		t.Errorf("ExpensesCents = %d; want 75000", result.ExpensesCents)
	}
	if result.TaxCents != 16000 {
		t.Errorf("TaxCents = %d; want 16000", result.TaxCents)
	}

	err = CalculateAdvanceTaxForThreeMonths(make([]IncomeForm, 4), &result)
	if err == nil {
		t.Error("Expected error for 4 months, got nil")
	}
}

func TestCalculateTaxesForThreeMonthsLowIncome(t *testing.T) {
	config := TaxesConfig{
		PensionPercentagePartOne:  10.0,
		PensionPercentagePartTwo:  5.0,
		HealthInsurancePercentage: 5.0,
		ExpensesPercentage:        25.0,
		TaxPercentage:             10.0,
	}
	forms := []IncomeForm{
		{MonthIncomeCents: 10000, SocialSecurityToPayCents: 60000, TaxesConfig: config},
		{MonthIncomeCents: 10000, SocialSecurityToPayCents: 60000, TaxesConfig: config},
		{MonthIncomeCents: 10000, SocialSecurityToPayCents: 60000, TaxesConfig: config},
	}

	// Total Income = 30000
	// Paid Insurance = 60000 * 3 = 180000
	// Expenses = 25% of 30000 = 7500
	// Taxable = 30000 - 7500 - 60000 = -37500
	// Tax must be 0 because taxable income is negative

	var result CalculatedTax
	err := CalculateAdvanceTaxForThreeMonths(forms, &result)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.TaxCents != 0 {
		t.Errorf("TaxCents = %d; want 0", result.TaxCents)
	}
}
