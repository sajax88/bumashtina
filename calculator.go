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

	IsMonthSkipped bool // Дейността е прекъсната, този месец не се брои при изравняване на вноските

	// Calculated
	TaxesToPayCents          int64
	ExpensesCents            int64
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
		return false, "Невалидна стойност на работни дни"
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

	if f.IsMonthSkipped && (f.MonthIncomeCents != 0 || f.TaxedIncomeCents != 0) {
		return false, "Ако не сте упражнявали дейност, въведете нулев доход и осигурителен доход"
	}

	if f.MonthIncomeCents == 0 && !f.IsMonthSkipped && f.TaxedIncomeCents > f.TaxesConfig.MinInsuranceIncomeCents {
		return false, fmt.Sprintf("Ако нямате приходи, осигурителният доход не може да е над минималения")
	}

	if f.Month < 1 || f.Month > 12 {
		return false, fmt.Sprintf("Невалиден месец %d", f.Month)
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
	insuranceCents := float64(taxedIncomeCents) * insurancePercent / 100
	return int64(math.Round(insuranceCents))
}

func CalculateSocialSecurityAndTaxForMonth(f *IncomeForm) {
	CalculateSocialSecurity(f)

	expenses := math.Round(float64(f.MonthIncomeCents) * f.TaxesConfig.ExpensesPercentage / 100)
	f.ExpensesCents = int64(expenses)

	var insurance float64
	if f.SocialSecurityReallyPaidCents > 0 {
		insurance = float64(f.SocialSecurityReallyPaidCents)
	} else {
		insurance = float64(f.SocialSecurityToPayCents)
	}

	taxedIncome := math.Round(float64(f.MonthIncomeCents) - expenses - insurance)
	f.TaxesToPayCents = int64(math.Round(taxedIncome * f.TaxesConfig.TaxPercentage / 100))
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
		// We suppose that the taxes are the same for all months of the quarter!
		if taxPercent == 0 {
			taxPercent = f.TaxesConfig.TaxPercentage
			expensesPercent = f.TaxesConfig.ExpensesPercentage
		}

		incomeTotalCents += f.MonthIncomeCents

		if f.SocialSecurityReallyPaidCents > 0 {
			paidInsuranceCents += f.SocialSecurityReallyPaidCents
		} else {
			paidInsuranceCents += f.SocialSecurityToPayCents
			result.Notes = "Въз основа на изчислените осигуровки"
		}
	}

	// (total income for 3 months - expenses - insurances) * taxPercentage
	incomeWithDeductions := float64(incomeTotalCents) - float64(incomeTotalCents)*expensesPercent/100 - float64(paidInsuranceCents)
	advanceTax := incomeWithDeductions * taxPercent / 100

	result.PaidInsuranceCents = paidInsuranceCents
	result.ExpensesCents = int64(math.Round(float64(incomeTotalCents) * expensesPercent / 100))

	taxCents := int64(math.Round(advanceTax))

	if taxCents > 0 {
		result.TaxCents = taxCents
	} else {
		result.TaxCents = 0
	}

	return nil
}

type MonthlyAlignmentResult struct {
	AverageTaxedIncomeCents   int64
	FinalInsuranceIncomeCents int64 // Осигурителен доход (= облагаем между лимитите)
}

type YearlyAlignmentResult struct {
	IsCalculated           bool
	YearlyGrossIncomeCents int64
	TaxedYearlyIncomeCents int64
	TaxesReallyPaidCents   int64
	Months                 []MonthlyAlignmentResult
}

// TODO: unit tests
func GetYearlyAlignmentResult(forms []IncomeForm) YearlyAlignmentResult {
	var result YearlyAlignmentResult

	// First, count the total yearly income
	activeMonths := 0
	for _, f := range forms {
		if f.IsMonthSkipped {
			continue
		}
		activeMonths++

		result.YearlyGrossIncomeCents += f.MonthIncomeCents
		result.TaxedYearlyIncomeCents += f.MonthIncomeCents - f.ExpensesCents
		result.TaxesReallyPaidCents += f.TaxesReallyPaidCents
	}

	averageMonthlyTaxedIncome := result.TaxedYearlyIncomeCents / int64(activeMonths)

	// Now recalculate the average monthly taxed income and adjust it between the min and max taxed income
	for _, f := range forms {
		if f.IsMonthSkipped {
			continue
		}

		monthlyResult := getMonthlyAlignmentResult(f, averageMonthlyTaxedIncome)
		result.Months = append(result.Months, monthlyResult)

		//Пример: Годишен облагаем доход от 18 000 €,
		//	дейност през всичките 12 месеца → 18 000 ÷ 12 = 1 500 € окончателен месечен осигурителен доход.
		//		Тъй като доходът се разпределя по равно, получената стойност е еднаква за всеки активен месец.
		//		Под 550,66 € → вдига се до 550,66 €
		//	Над 2 111,64 € → намалява се до 2 111,64 €
		//

		//	След като знаеш окончателния си осигурителен доход,
		//		изчисляваш годишните осигуровки върху него. НАП сравнява тази сума с осигуровките,
		//		които вече си платил авансово:
		//	Платил си по-малко от дължимото → доплащаш разликата до 30.04
		//	Платил си повече от дължимото → надвнесеното се приспада от бъдещи задължения или ти се възстановява
	}

	// TODO: build income form, CalculateSocialSecurity(f *IncomeForm)

	result.IsCalculated = true
	return result
}

func getMonthlyAlignmentResult(f IncomeForm, averageMonthlyTaxedIncome int64) MonthlyAlignmentResult {
	minTaxedIncome := f.TaxesConfig.MinInsuranceIncomeCents
	maxTaxedIncome := f.TaxesConfig.MaxInsuranceIncomeCents

	insuranceIncome := averageMonthlyTaxedIncome
	if insuranceIncome < minTaxedIncome {
		insuranceIncome = minTaxedIncome
	}
	if insuranceIncome > maxTaxedIncome {
		insuranceIncome = maxTaxedIncome
	}

	monthlyResult := MonthlyAlignmentResult{
		AverageTaxedIncomeCents:   averageMonthlyTaxedIncome,
		FinalInsuranceIncomeCents: insuranceIncome,
		// TODO: count insurance
	}

	return monthlyResult
}
