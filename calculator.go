package main

import (
	"fmt"
	"math"
)

type IncomeForm struct {
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
	// TODO
	TaxesToPayCents          int64
	SocialSecurityToPayCents int64

	// Actually paid
	// TODO
	TaxesReallyPaidCents          int64
	SocialSecurityReallyPaidCents int64

	TaxesConfig TaxesConfig
	Settings    Settings
}

type CalculatedTax struct {
	TotalIncomeCents int64
	TaxCents         int64
	Quarter          int
	Year             int
	MonthStart       int
	MonthEnd         int
}

func CalculateSocialSecurity(f IncomeForm) int64 {
	insurancePercent := f.TaxesConfig.PensionPercentage + f.TaxesConfig.HealthInsurancePercentage
	if f.Settings.IsPregnancyInsuranceEnabled {
		insurancePercent += f.TaxesConfig.PregnancyInsurancePercentage
	}

	insuranceCents := float32(f.TaxedIncomeCents) * insurancePercent / 100

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

func CalculateAdvanceTaxForThreeMonths(forms []IncomeForm) (int64, error) {
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

	return int64(math.Round(float64(advanceTax))), nil
}
