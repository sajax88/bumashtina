package main

import "fmt"

func ValidateWorkDays(f IncomeForm) (bool, string) {
	expectedWorkDays := GetWorkDaysNumber(f.Year, f.Month)

	if f.DayStart > 0 {
		// TODO Subtract only work days
	}

	if f.DayEnd > 0 {
		// TODO Subtract only work days
	}

	if expectedWorkDays != f.WorkDaysTotal {
		return false, fmt.Sprintf("Изчислени работни дни са %d, вие въведохте %d.", expectedWorkDays, f.WorkDaysTotal) // TODO: CHECK MESSAGE
	}

	return true, ""
}

func GetWorkDaysNumber(year int16, month int16) int16 {
	if year < MinYear {
		return 0
	}

	if month < 1 || month > 12 {
		return 0
	}

	return getWorkDaysFromCalendar(year, month)
}

func getWorkDaysFromCalendar(year int16, month int16) int16 {
	calendar := map[int16]map[int16]int16{
		2026: {
			1:  20,
			2:  20,
			3:  21,
			4:  20,
			5:  18,
			6:  22,
			7:  23,
			8:  21,
			9:  20,
			10: 22,
			11: 21,
			12: 20,
		},
		2027: {
			1:  20,
			2:  20,
			3:  22,
			4:  21,
			5:  17,
			6:  22,
			7:  22,
			8:  22,
			9:  20,
			10: 21,
			11: 22,
			12: 20,
		},
	}

	y, ok := calendar[year]
	if ok {
		return y[month]
	}

	// Decided not to fall back to the naive calculation, better to make the user check and enter the days
	return 0
}
