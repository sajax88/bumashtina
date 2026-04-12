package main

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

func CalculateAdvanceTaxForThreeMonths(f IncomeForm, insurance int32) int32 {
	// За да определите авансовия си данък за тримесечието
	// следва да извадите от доход за 3 месеца признатите разходи 25%
	// и платените осигуровки за трите месеца. Данък 10%
	return 0 // TODO
}
