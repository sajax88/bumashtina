package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncomeForm_Validate(t *testing.T) {
	defaultConfig := GetDefaultTaxesConfig()
	validForm := IncomeForm{
		Month:            1,
		Year:             2026,
		MonthIncomeCents: 100000,
		TaxedIncomeCents: 70000,
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
		{"Zero income but taxed income non-zero", func() IncomeForm { f := validForm; f.MonthIncomeCents = 0; f.TaxedIncomeCents = 70000; return f }(), false},
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

func TestGetMonthlyAlignmentResult(t *testing.T) {
	config := TaxesConfig{
		MinInsuranceIncomeCents:      62020,
		MaxInsuranceIncomeCents:      230000,
		PensionPercentagePartOne:     14.8,
		PensionPercentagePartTwo:     5.0,
		HealthInsurancePercentage:    8.0,
		PregnancyInsurancePercentage: 3.5,
	}

	tests := []struct {
		name                         string
		form                         IncomeForm
		averageMonthlyTaxedIncome    int64
		expectedFinalInsuranceIncome int64
		expectedRecalculatedSocSec   int64
	}{
		{
			name: "Average income within min/max range",
			form: IncomeForm{
				Month:                         1,
				MonthIncomeCents:              150000,
				SocialSecurityReallyPaidCents: 20000,
				TaxesConfig:                   config,
			},
			averageMonthlyTaxedIncome:    100000,
			expectedFinalInsuranceIncome: 100000,
			// 14 800 (PensionOne) + 5000 (PensionTwo) + 8000 (Health) = 27 800
			expectedRecalculatedSocSec: 27800,
		},
		{
			name: "Average income below minimum gets clamped up",
			form: IncomeForm{
				Month:                         2,
				MonthIncomeCents:              30000,
				SocialSecurityReallyPaidCents: 10000,
				TaxesConfig:                   config,
			},
			averageMonthlyTaxedIncome:    50000,
			expectedFinalInsuranceIncome: 62020, // clamped to MinInsuranceIncomeCents
			// 9179 (PensionOne) + 3101 (PensionTwo) + 4962 (Health) = 17 242
			expectedRecalculatedSocSec: 17242,
		},
		{
			name: "Average income above maximum gets clamped down",
			form: IncomeForm{
				Month:                         3,
				MonthIncomeCents:              500000,
				SocialSecurityReallyPaidCents: 60000,
				TaxesConfig:                   config,
			},
			averageMonthlyTaxedIncome:    300000,
			expectedFinalInsuranceIncome: 230000, // clamped to MaxInsuranceIncomeCents
			// 34 040 (PensionOne) + 11 500 (PensionTwo) + 18 400 (Health) = 63 940
			expectedRecalculatedSocSec: 63940,
		},
		{
			name: "Pregnancy insurance enabled adds extra to PensionPartOne",
			form: IncomeForm{
				Month:                         4,
				MonthIncomeCents:              150000,
				SocialSecurityReallyPaidCents: 20000,
				TaxesConfig:                   config,
				Settings: Settings{
					IsPregnancyInsuranceEnabled: true,
				},
			},
			averageMonthlyTaxedIncome:    100000,
			expectedFinalInsuranceIncome: 100000,
			// PensionOne = 14800 + 3500 (pregnancy) = 18300; PensionTwo = 5000; Health = 8000 => 31300
			expectedRecalculatedSocSec: 31300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMonthlyAlignmentResult(tt.form, tt.averageMonthlyTaxedIncome)

			if result.Month != tt.form.Month {
				t.Errorf("Month = %d; want %d", result.Month, tt.form.Month)
			}
			if result.GrossIncomeCents != tt.form.MonthIncomeCents {
				t.Errorf("GrossIncomeCents = %d; want %d", result.GrossIncomeCents, tt.form.MonthIncomeCents)
			}
			if result.AverageTaxedIncomeCents != tt.averageMonthlyTaxedIncome {
				t.Errorf("AverageTaxedIncomeCents = %d; want %d", result.AverageTaxedIncomeCents, tt.averageMonthlyTaxedIncome)
			}
			if result.FinalSocSecIncomeCents != tt.expectedFinalInsuranceIncome {
				t.Errorf("FinalSocSecIncomeCents = %d; want %d", result.FinalSocSecIncomeCents, tt.expectedFinalInsuranceIncome)
			}
			if result.PaidSocialSecurityCents != tt.form.SocialSecurityReallyPaidCents {
				t.Errorf("PaidSocialSecurityCents = %d; want %d", result.PaidSocialSecurityCents, tt.form.SocialSecurityReallyPaidCents)
			}
			if result.RecalculatedSocialSecurityCents != tt.expectedRecalculatedSocSec {
				t.Errorf("RecalculatedSocialSecurityCents = %d; want %d", result.RecalculatedSocialSecurityCents, tt.expectedRecalculatedSocSec)
			}
		})
	}
}

func TestGetYearlyAlignmentResult(t *testing.T) {
	baseConfig := TaxesConfig{
		MinInsuranceIncomeCents:   55066,
		MaxInsuranceIncomeCents:   211164,
		ExpensesPercentage:        25.0,
		TaxPercentage:             10.0,
		HealthInsurancePercentage: 8.0,
		PensionPercentagePartOne:  14.8,
		PensionPercentagePartTwo:  5.0,
	}

	t.Run("No active months", func(t *testing.T) {
		var forms []IncomeForm

		result := GetYearlyAlignmentResult(forms)
		expected := YearlyAlignmentResult{
			IsCalculated:                    false,
			YearlyGrossIncomeCents:          0,
			TaxedYearlyIncomeCents:          0,
			TaxesReallyPaidCents:            0,
			SocialSecurityReallyPaidCents:   0,
			RecalculatedSocialSecurityCents: 0,
			RecalculatedTaxCents:            0,
			SocialSecurityDiffCents:         0,
			TaxesDiffCents:                  0,
			ExpensesPercentage:              0,
			Months:                          []MonthlyAlignmentResult(nil),
		}

		assert.Equal(t, result, expected)
	})

	t.Run("Exact calculation, income between limits, no alignment needed", func(t *testing.T) {
		config := baseConfig
		var forms = []IncomeForm{
			{
				Month:                         9,
				MonthIncomeCents:              100000,
				ExpensesCents:                 25000,
				SocialSecurityReallyPaidCents: 20850,
				TaxesReallyPaidCents:          5416, // 3rd quarter
				TaxesConfig:                   config,
			},
			{
				Month:                         12,
				MonthIncomeCents:              200000,
				ExpensesCents:                 50000,
				SocialSecurityReallyPaidCents: 41700,
				TaxesReallyPaidCents:          0,
				TaxesConfig:                   config,
			},
		}

		result := GetYearlyAlignmentResult(forms)

		expected := YearlyAlignmentResult{
			IsCalculated:                    true,
			YearlyGrossIncomeCents:          300000,
			TaxedYearlyIncomeCents:          225000, // (100000 - 25% expenses) + (200000 - 25% expenses)
			TaxesReallyPaidCents:            5416,
			SocialSecurityReallyPaidCents:   62550, // 20850 + 41700
			RecalculatedSocialSecurityCents: 62550, // 31275 * 2
			RecalculatedTaxCents:            16246, // Recalculated based on taxed income minus recalculated insurance
			SocialSecurityDiffCents:         0,     // 62550 - 62550
			TaxesDiffCents:                  10830, // Left to pay: 16246 - 5416
			ExpensesPercentage:              25.0,
			Months: []MonthlyAlignmentResult{
				{
					Month:                           9,
					GrossIncomeCents:                100000,
					AverageTaxedIncomeCents:         112500, // (75000 + 150000) / 2
					FinalSocSecIncomeCents:          112500, // Clamped between min/max
					PaidSocialSecurityCents:         20850,
					RecalculatedSocialSecurityCents: 31275, // 14.8% + 5% + 8% of 112500
				},
				{
					Month:                           12,
					GrossIncomeCents:                200000,
					AverageTaxedIncomeCents:         112500,
					FinalSocSecIncomeCents:          112500,
					PaidSocialSecurityCents:         41700,
					RecalculatedSocialSecurityCents: 31275, // 14.8% + 5% + 8% of 112500
				},
			},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("Minimal insurance income, big alignment", func(t *testing.T) {
		config := baseConfig
		var forms = []IncomeForm{
			{
				Month:                         9,
				MonthIncomeCents:              200000,
				ExpensesCents:                 50000,
				SocialSecurityReallyPaidCents: 15309, // On the minimal income
				TaxesReallyPaidCents:          13469,
				TaxesConfig:                   config,
			},
			{
				Month:                         11,
				MonthIncomeCents:              200000,
				ExpensesCents:                 50000,
				SocialSecurityReallyPaidCents: 15309,
				TaxesReallyPaidCents:          0,
				TaxesConfig:                   config,
			},
			{
				Month:                         12,
				MonthIncomeCents:              200000,
				ExpensesCents:                 50000,
				SocialSecurityReallyPaidCents: 15309,
				TaxesReallyPaidCents:          0,
				TaxesConfig:                   config,
			},
		}

		result := GetYearlyAlignmentResult(forms)

		expected := YearlyAlignmentResult{
			IsCalculated:                    true,
			YearlyGrossIncomeCents:          600000,
			TaxedYearlyIncomeCents:          450000,
			TaxesReallyPaidCents:            13469,
			SocialSecurityReallyPaidCents:   45927,
			RecalculatedSocialSecurityCents: 125100,
			RecalculatedTaxCents:            32490,
			SocialSecurityDiffCents:         79173,
			TaxesDiffCents:                  19021,
			ExpensesPercentage:              25.0,
			Months: []MonthlyAlignmentResult{
				{
					Month:                           9,
					GrossIncomeCents:                200000,
					AverageTaxedIncomeCents:         150000,
					FinalSocSecIncomeCents:          150000,
					PaidSocialSecurityCents:         15309,
					RecalculatedSocialSecurityCents: 41700,
				},
				{
					Month:                           11,
					GrossIncomeCents:                200000,
					AverageTaxedIncomeCents:         150000,
					FinalSocSecIncomeCents:          150000,
					PaidSocialSecurityCents:         15309,
					RecalculatedSocialSecurityCents: 41700,
				},
				{
					Month:                           12,
					GrossIncomeCents:                200000,
					AverageTaxedIncomeCents:         150000,
					FinalSocSecIncomeCents:          150000,
					PaidSocialSecurityCents:         15309,
					RecalculatedSocialSecurityCents: 41700,
				},
			},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("Paid too much, overpayment left", func(t *testing.T) {
		config := baseConfig
		var forms = []IncomeForm{
			{
				Month:                         9,
				MonthIncomeCents:              100000,
				ExpensesCents:                 25000,
				SocialSecurityReallyPaidCents: 58703, // Max insurance income
				TaxesReallyPaidCents:          1630,  // 3rd quarter
				TaxesConfig:                   config,
			},
			{
				Month:                         12,
				MonthIncomeCents:              200000,
				ExpensesCents:                 50000,
				SocialSecurityReallyPaidCents: 58703, // Max insurance income
				TaxesReallyPaidCents:          0,
				TaxesConfig:                   config,
			},
		}

		result := GetYearlyAlignmentResult(forms)

		expected := YearlyAlignmentResult{
			IsCalculated:                    true,
			YearlyGrossIncomeCents:          300000,
			TaxedYearlyIncomeCents:          225000, // (100000 - 25% expenses) + (200000 - 25% expenses)
			TaxesReallyPaidCents:            1630,
			SocialSecurityReallyPaidCents:   117406, // 58703 * 2
			RecalculatedSocialSecurityCents: 62550,  // 31275 * 2
			RecalculatedTaxCents:            16246,  // Recalculated based on taxed income minus recalculated insurance
			SocialSecurityDiffCents:         54856,  // 62550 - 117406, overpayment should be positive
			TaxesDiffCents:                  14616,  // Left to pay: 16246 - 1630
			ExpensesPercentage:              25.0,
			Months: []MonthlyAlignmentResult{
				{
					Month:                           9,
					GrossIncomeCents:                100000,
					AverageTaxedIncomeCents:         112500, // (75000 + 150000) / 2
					FinalSocSecIncomeCents:          112500, // Clamped between min/max
					PaidSocialSecurityCents:         58703,
					RecalculatedSocialSecurityCents: 31275, // 14.8% + 5% + 8% of 112500
				},
				{
					Month:                           12,
					GrossIncomeCents:                200000,
					AverageTaxedIncomeCents:         112500,
					FinalSocSecIncomeCents:          112500,
					PaidSocialSecurityCents:         58703,
					RecalculatedSocialSecurityCents: 31275, // 14.8% + 5% + 8% of 112500
				},
			},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("Max insurance threshold raised", func(t *testing.T) {
		config1 := baseConfig
		config2 := baseConfig
		config2.MaxInsuranceIncomeCents = 230000

		var forms = []IncomeForm{
			{
				Month:                         9,
				MonthIncomeCents:              400000,
				ExpensesCents:                 100000,
				SocialSecurityReallyPaidCents: 58703,   // For nax insurance income, old
				TaxesReallyPaidCents:          24130,   // 3rd quarter
				TaxesConfig:                   config1, // Here the threshold is 211164, we stop at it and
			},
			{
				Month:                         12,
				MonthIncomeCents:              400000,
				ExpensesCents:                 100000,
				SocialSecurityReallyPaidCents: 63940, // For max insurance income, new
				TaxesReallyPaidCents:          0,
				TaxesConfig:                   config2, // Here the threshold is 230000, we now pay more for social security
			},
		}

		result := GetYearlyAlignmentResult(forms)

		expected := YearlyAlignmentResult{
			IsCalculated:                    true,
			YearlyGrossIncomeCents:          800000,
			TaxedYearlyIncomeCents:          600000, // (400000 - 25% expenses) + (400000 - 25% expenses)
			TaxesReallyPaidCents:            24130,
			SocialSecurityReallyPaidCents:   122643, // 58703 + 63940
			RecalculatedSocialSecurityCents: 122643, // The same, as we went for max threshold
			RecalculatedTaxCents:            47736,  // Left for month 12
			SocialSecurityDiffCents:         0,
			TaxesDiffCents:                  23606, // Left to pay, a bit lower than for month 9 because we paid more social security
			ExpensesPercentage:              25.0,
			Months: []MonthlyAlignmentResult{
				{
					Month:                           9,
					GrossIncomeCents:                400000,
					AverageTaxedIncomeCents:         300000,
					FinalSocSecIncomeCents:          211164, // Old max income limit
					PaidSocialSecurityCents:         58703,
					RecalculatedSocialSecurityCents: 58703, // Stays on max threshold for that month
				},
				{
					Month:                           12,
					GrossIncomeCents:                400000,
					AverageTaxedIncomeCents:         300000,
					FinalSocSecIncomeCents:          230000, // New max income limit
					PaidSocialSecurityCents:         63940,
					RecalculatedSocialSecurityCents: 63940, // Stays on max threshold for that month
				},
			},
		}

		assert.Equal(t, result, expected)
	})
}
