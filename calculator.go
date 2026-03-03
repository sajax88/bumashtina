package main

type IncomeForm struct {
	Month            int16
	Year             int16
	MonthIncomeCents int64
	TaxedIncomeCents int64
	DayEnd           int16
	DayStart         int16
	WorkDaysTotal    int16

	TaxesToPayCents          int64
	SocialSecurityToPayCents int64

	TaxesConfig TaxesConfig
}

func CalculateSocialSecurity(f IncomeForm, settings Settings) int32 {
	// За да определите осигуровките си
	// следва да умножите осигурителен доход по 0.278 или 0.313
	// в зависимост от това дали имате включено майчинство в настройките.

	insurancePercent := f.TaxesConfig.PensionPercentage + f.TaxesConfig.HealthInsurancePercentage
	if settings.IsPregnancyInsuranceEnabled {
		insurancePercent += f.TaxesConfig.PregnancyInsurancePercentage
	}

	res := float32(f.TaxedIncomeCents) * insurancePercent / 100

	return int32(res)
}

func CalculateTaxForMonth(f IncomeForm, insurance int32) int32 {
	// За да определите авансовия си данък
	// следва да извадите от доход  признатите разходи 25%
	// и платените осигуровки

	expenses := int32(float32(f.MonthIncomeCents) * f.TaxesConfig.ExpensesPercentage / 100)

	return int32(f.MonthIncomeCents) - expenses - insurance
}

func CalculateAdvanceTaxForThreeMonths(f IncomeForm, insurance int32) int32 {
	// За да определите авансовия си данък за тримесечието
	// следва да извадите от доход за 3 месеца признатите разходи 25%
	// и платените осигуровки за трите месеца. Данък 10%
	return 0 // TODO
}
