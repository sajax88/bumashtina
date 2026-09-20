package main

import (
	"fmt"
	"time"
)

func ValidateWorkDays(f IncomeForm) (bool, string) {
	expectedWorkDays := GetWorkDaysNumber(f.Year, f.Month)

	if f.DayStart > 0 || f.DayEnd > 0 {
		expectedWorkDays = subtractWorkDays(expectedWorkDays, f)
	}

	if expectedWorkDays != f.WorkDaysTotal {
		return false, fmt.Sprintf("Изчислените работни дни са %d, а вие въведохте %d.", expectedWorkDays, f.WorkDaysTotal)
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

/**
* If the user enters a start or/and end date,
* we need to subtract the work days from the start/end of the month
* from the expected number of work days
 */
func subtractWorkDays(expectedWorkDays int16, f IncomeForm) int16 {
	dayStart := int(f.DayStart)
	dayEnd := int(f.DayEnd)

	dt := time.Date(int(f.Year), time.Month(f.Month), 1, 0, 0, 0, 0, time.UTC)
	for dt.Month() == time.Month(f.Month) {
		// TODO v1.2: check official holidays if exist for this year, unit test cases for this
		if dt.Weekday() != time.Saturday && dt.Weekday() != time.Sunday {
			if dayStart > 1 && dt.Day() < dayStart {
				expectedWorkDays--
			}
			if dayEnd > 0 && dt.Day() > dayEnd {
				expectedWorkDays--
			}
		}
		dt = dt.AddDate(0, 0, 1)
	}
	if expectedWorkDays < 0 {
		return 0
	}
	return expectedWorkDays
}
