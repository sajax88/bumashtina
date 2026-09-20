package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtractWorkDays(t *testing.T) {
	tests := []struct {
		name             string
		expectedWorkDays int16
		form             IncomeForm
		want             int16
	}{
		{
			name:             "No DayStart or DayEnd - no change",
			expectedWorkDays: 20,
			form:             IncomeForm{Year: 2026, Month: 1, DayStart: 0, DayEnd: 0},
			want:             20,
		},
		{
			name:             "DayStart = 1 has no effect (job starts on the 1st)",
			expectedWorkDays: 20,
			form:             IncomeForm{Year: 2026, Month: 1, DayStart: 1, DayEnd: 0},
			want:             20,
		},
		{
			name:             "DayEnd = 1 means no work days",
			expectedWorkDays: 0,
			form:             IncomeForm{Year: 2026, Month: 1, DayStart: 0, DayEnd: 1},
			want:             0,
		},
		{
			name:             "DayStart only - subtracts weekdays before start",
			expectedWorkDays: 20,
			// January 2026 starts on Thursday. Days 1-5: Thu, Fri, Sat, Sun, Mon.
			// Weekdays before day 6 (excluding weekends): Thu, Fri, Mon = 3 days.
			form: IncomeForm{Year: 2026, Month: 1, DayStart: 6, DayEnd: 0},
			want: 17,
		},
		{
			name:             "DayEnd only - subtracts weekdays after end",
			expectedWorkDays: 20,
			// January 2026: days 26-31 are Mon, Tue, Wed, Thu, Fri, Sat.
			// Weekdays after day 25 (excluding weekends): Mon, Tue, Wed, Thu, Fri = 5 days.
			form: IncomeForm{Year: 2026, Month: 1, DayStart: 0, DayEnd: 25},
			want: 15,
		},
		{
			name:             "DayStart and DayEnd combined",
			expectedWorkDays: 20,
			form:             IncomeForm{Year: 2026, Month: 1, DayStart: 6, DayEnd: 25},
			want:             12,
		},
		{
			name:             "DayEnd equals last day of month - no subtraction",
			expectedWorkDays: 20,
			form:             IncomeForm{Year: 2026, Month: 1, DayStart: 0, DayEnd: 31},
			want:             20,
		},
		{
			name:             "DayStart in a shorter month (February)",
			expectedWorkDays: 20,
			// February 2026 starts on Sunday. Day1=Sun, 2=Mon, 3=Tue, 4=Wed, 5=Thu.
			// Weekdays before day 5 (excluding weekends): Mon, Tue, Wed = 3 days.
			form: IncomeForm{Year: 2026, Month: 2, DayStart: 5, DayEnd: 0},
			want: 17,
		},
		{
			name:             "Result is clamped to 0 when subtractions exceed expected work days",
			expectedWorkDays: 2,
			form:             IncomeForm{Year: 2026, Month: 1, DayStart: 6, DayEnd: 0},
			want:             0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subtractWorkDays(tt.expectedWorkDays, tt.form)
			assert.Equal(t, tt.want, got)
		})
	}
}
