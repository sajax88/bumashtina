package main

import "fmt"

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

	TaxesToPayCents          int64
	SocialSecurityToPayCents int64

	TaxesConfig TaxesConfig
	Settings    Settings
}

type CalculatedTax struct {
	TotalIncomeCents int64
	TaxCents         int64
	Quarter          int
	Year             int
}

func CalculateSocialSecurity(f IncomeForm, settings Settings) int32 {
	insurancePercent := f.TaxesConfig.PensionPercentage + f.TaxesConfig.HealthInsurancePercentage
	if settings.IsPregnancyInsuranceEnabled {
		insurancePercent += f.TaxesConfig.PregnancyInsurancePercentage
	}

	res := float32(f.TaxedIncomeCents) * insurancePercent / 100

	return int32(res)
}

func CalculateTaxForMonth(f IncomeForm, insurance int32) int32 {
	expenses := int32(float32(f.MonthIncomeCents) * f.TaxesConfig.ExpensesPercentage / 100)

	return (int32(f.MonthIncomeCents) - expenses - insurance) * int32(f.TaxesConfig.TaxPercentage) / 100
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

// TODO: pass actually paid insurance if present, otherwise calculated
func CalculateAdvanceTaxForThreeMonths(forms []IncomeForm, paidInsuranceCents int32) (float32, error) {
	if len(forms) > 3 {
		return 0, fmt.Errorf("Очаквах данни за 3 месеца, получих %d", len(forms))
	}

	// За да определите авансовия си данък за тримесечието
	// следва да извадите от доход за 3 месеца признатите разходи 25%
	// и платените осигуровки за трите месеца. Данък 10%
	var incomeTotalCents int64
	var taxPercent float32
	var expensesPercent float32
	for _, f := range forms {
		incomeTotalCents += f.MonthIncomeCents
		taxPercent = f.TaxesConfig.TaxPercentage
		expensesPercent = f.TaxesConfig.ExpensesPercentage
	}
	// (total income for 3 months - expenses - insurances) * taxPercentage
	incomeWithDeductions := float32(incomeTotalCents) - float32(incomeTotalCents)*expensesPercent - float32(paidInsuranceCents)
	advanceTax := incomeWithDeductions * taxPercent / 100 / MONEY_DIVIDER

	return advanceTax, nil
}
