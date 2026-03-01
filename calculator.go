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

func CalculateSocialSecurity() {
	// TODO
	// За да определите осигуровките си
	// следва да умножите осигурителен доход (макс по дефолт) по 0.278 или 0.313
	// в зависимост от това дали имате включено майчинство в настройките.
}

func CalculateAdvanceTaxForThreeMonths() {
	// TODO
	//За да определите авансовия си данък за тримесечието
	//следва да извадите от доход за 3 месеца признатите разходи 25%
	//и платените осигуровки за трите месеца. Данък 10%
}
