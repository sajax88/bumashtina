package main

import (
	"fmt"
	"math"
)

type IncomeForm struct {
	// Initial user input
	Month             int16
	Year              int16
	MonthIncomeCents  int64
	TaxedIncomeCents  int64
	DayEnd            int16
	DayStart          int16
	WorkDaysTotal     int16
	WorkDaysReal      int16
	WorkDaysSickLeave int16

	// Calculated
	TaxesToPayCents          int64
	SocialSecurityToPayCents int64
	SocialSecurityToPayParts SocialSecurityParts

	// Actually paid (vs calculated)
	// Taxes - only for the last month in the quarter
	TaxesReallyPaidCents          int64
	SocialSecurityReallyPaidCents int64
	SocialSecurityReallyPaidParts SocialSecurityParts // Used in Declaration 6 if entered, otherwise use calculated

	// Config at the time of saving
	TaxesConfig TaxesConfig
	Settings    Settings
}

type SocialSecurityParts struct {
	PensionPartOneCents  int64
	PensionPartTwoCents  int64
	HealthInsuranceCents int64
}

type CalculatedTax struct {
	TotalIncomeCents   int64
	TaxCents           int64
	ExpensesCents      int64
	PaidInsuranceCents int64
	Quarter            int
	Year               int
	MonthStart         int
	MonthEnd           int
}

func CalculateSocialSecurity(f *IncomeForm) {
	// ДОО
	pensionPartOne := calculateInsuranceFromPercentage(f.TaxedIncomeCents, f.TaxesConfig.PensionPercentagePartOne)
	if f.Settings.IsPregnancyInsuranceEnabled {
		// Общо заболяване и майчинство
		pensionPartOne += calculateInsuranceFromPercentage(f.TaxedIncomeCents, f.TaxesConfig.PregnancyInsurancePercentage)
	}

	// ДЗПО
	pensionPartTwo := calculateInsuranceFromPercentage(f.TaxedIncomeCents, f.TaxesConfig.PensionPercentagePartTwo)
	// НЗОК
	healthInsurance := calculateInsuranceFromPercentage(f.TaxedIncomeCents, f.TaxesConfig.HealthInsurancePercentage)

	f.SocialSecurityToPayParts = SocialSecurityParts{
		PensionPartOneCents:  pensionPartOne,
		PensionPartTwoCents:  pensionPartTwo,
		HealthInsuranceCents: healthInsurance,
	}

	f.SocialSecurityToPayCents = pensionPartOne + pensionPartTwo + healthInsurance
}

func calculateInsuranceFromPercentage(taxedIncomeCents int64, insurancePercent float32) int64 {
	insuranceCents := float32(taxedIncomeCents) * insurancePercent / 100
	return int64(math.Round(float64(insuranceCents)))
}

func CalculateTaxForMonth(f IncomeForm) int64 {
	expenses := math.Round(float64(f.MonthIncomeCents) * float64(f.TaxesConfig.ExpensesPercentage) / 100)
	var insurance float64
	if f.SocialSecurityReallyPaidCents > 0 {
		insurance = float64(f.SocialSecurityReallyPaidCents)
	} else {
		insurance = float64(f.SocialSecurityToPayCents)
	}

	return int64(math.Round(math.Round(float64(f.MonthIncomeCents)-expenses-insurance) * float64(f.TaxesConfig.TaxPercentage) / 100))
}

func CalculateIncomeForThreeMonths(forms []IncomeForm) (int64, error) {
	if len(forms) > 3 {
		return 0, fmt.Errorf("Очаквах данни за 3 месеца, получих %d", len(forms))
	}
	var incomeTotalCents int64
	for _, f := range forms {
		incomeTotalCents += f.MonthIncomeCents
	}
	return incomeTotalCents, nil
}

func CalculateAdvanceTaxForThreeMonths(forms []IncomeForm, result *CalculatedTax) (int64, error) {
	if len(forms) > 3 {
		return 0, fmt.Errorf("Очаквах данни за 3 месеца, получих %d", len(forms))
	}

	// За да определите авансовия си данък за тримесечието
	// следва да извадите от доход за 3 месеца признатите разходи 25%
	// и платените осигуровки за трите месеца. Данък 10%
	var incomeTotalCents int64
	var taxPercent float32
	var expensesPercent float32
	var paidInsuranceCents int64
	for _, f := range forms {
		taxPercent = f.TaxesConfig.TaxPercentage
		expensesPercent = f.TaxesConfig.ExpensesPercentage

		incomeTotalCents += f.MonthIncomeCents
		paidInsuranceCents += f.SocialSecurityReallyPaidCents // TODO: use calculated insurance if no really paid
	}
	// (total income for 3 months - expenses - insurances) * taxPercentage
	incomeWithDeductions := float32(incomeTotalCents) - float32(incomeTotalCents)*expensesPercent/100 - float32(paidInsuranceCents)
	advanceTax := incomeWithDeductions * taxPercent / 100

	result.PaidInsuranceCents = paidInsuranceCents
	result.ExpensesCents = int64(math.Round(float64(float32(incomeTotalCents) * expensesPercent / 100)))

	return int64(math.Round(float64(advanceTax))), nil
}
