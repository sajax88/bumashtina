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

	// Actually paid (vs. calculated)
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
	Notes              string
	TaxesConfig        TaxesConfig
}

func (f IncomeForm) Validate() (bool, string) {
	if f.WorkDaysTotal < 0 || f.WorkDaysTotal > 31 {
		return false, "Невалиден брой на общите работни дни"
	}

	if f.WorkDaysSickLeave > f.WorkDaysTotal {
		return false, "Дните в болничен не могат да надвишават общите работни дни"
	}

	if f.DayStart > 0 && f.DayEnd > 0 && f.DayStart >= f.DayEnd {
		return false, "Началният ден трябва да бъде преди крайния ден"
	}

	if f.DayStart < 0 || f.DayStart > 31 {
		return false, "Невалиден начален ден"
	}

	if f.DayEnd < 0 || f.DayEnd > 31 {
		return false, "Невалиден краен ден"
	}

	if f.MonthIncomeCents < 0 {
		return false, "Невалиден месечен доход"
	}

	if f.Year < MinYear {
		return false, fmt.Sprintf("Дати преди %d година не се поддържат", MinYear)
	}

	if f.MonthIncomeCents > 0 &&
		(f.TaxedIncomeCents < f.TaxesConfig.MinInsuranceIncomeCents || f.TaxedIncomeCents > f.TaxesConfig.MaxInsuranceIncomeCents) {
		return false, fmt.Sprintf(
			"Осигурителният доход трябва да бъде между %.2f и %.2f EUR",
			float64(f.TaxesConfig.MinInsuranceIncomeCents)/MoneyDivider,
			float64(f.TaxesConfig.MaxInsuranceIncomeCents)/MoneyDivider,
		)
	}

	return true, ""
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

func calculateInsuranceFromPercentage(taxedIncomeCents int64, insurancePercent float64) int64 {
	insuranceCents := float64(taxedIncomeCents) * insurancePercent / MoneyDivider
	return int64(math.Round(insuranceCents))
}

func CalculateTaxForMonth(f IncomeForm) int64 {
	expenses := math.Round(float64(f.MonthIncomeCents) * f.TaxesConfig.ExpensesPercentage / MoneyDivider)
	var insurance float64
	if f.SocialSecurityReallyPaidCents > 0 {
		insurance = float64(f.SocialSecurityReallyPaidCents)
	} else {
		insurance = float64(f.SocialSecurityToPayCents)
	}

	return int64(math.Round(math.Round(float64(f.MonthIncomeCents)-expenses-insurance) * f.TaxesConfig.TaxPercentage / MoneyDivider))
}

func CalculateIncomeForThreeMonths(forms []IncomeForm) (int64, error) {
	if len(forms) > 3 {
		return 0, fmt.Errorf("очаквах данни за 3 месеца, получих %d", len(forms))
	}
	var incomeTotalCents int64
	for _, f := range forms {
		incomeTotalCents += f.MonthIncomeCents
	}
	return incomeTotalCents, nil
}

func CalculateAdvanceTaxForThreeMonths(forms []IncomeForm, result *CalculatedTax) error {
	if len(forms) > 3 {
		return fmt.Errorf("въведете данни за 3 месеца, в момента има %d записа", len(forms))
	}

	// За да определите авансовия си данък за тримесечието
	// следва да извадите от доход за 3 месеца признатите разходи 25%
	// и платените осигуровки за трите месеца. Данък 10%
	var incomeTotalCents int64
	var taxPercent float64
	var expensesPercent float64
	var paidInsuranceCents int64
	for _, f := range forms {
		taxPercent = f.TaxesConfig.TaxPercentage
		expensesPercent = f.TaxesConfig.ExpensesPercentage

		incomeTotalCents += f.MonthIncomeCents

		if f.SocialSecurityReallyPaidCents > 0 {
			paidInsuranceCents += f.SocialSecurityReallyPaidCents
		} else {
			paidInsuranceCents += f.SocialSecurityToPayCents
			result.Notes = "Въз основа на изчислените осигуровки"
		}
	}

	// (total income for 3 months - expenses - insurances) * taxPercentage
	incomeWithDeductions := float64(incomeTotalCents) - float64(incomeTotalCents)*expensesPercent/MoneyDivider - float64(paidInsuranceCents)
	advanceTax := incomeWithDeductions * taxPercent / MoneyDivider

	result.PaidInsuranceCents = paidInsuranceCents
	result.ExpensesCents = int64(math.Round(float64(incomeTotalCents) * expensesPercent / MoneyDivider))

	result.TaxCents = int64(math.Round(advanceTax))

	return nil
}
