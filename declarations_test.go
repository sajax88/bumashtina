package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeDeclarationOne(t *testing.T) {
	tests := []struct {
		name     string
		form     IncomeForm
		user     UserConfig
		settings Settings
		wantErr  bool
		errMsg   string
		expected []byte
	}{
		{
			name: "Valid Input",
			form: IncomeForm{
				Month:             7,
				Year:              2026,
				TaxedIncomeCents:  206320,
				DayStart:          0,
				DayEnd:            0,
				WorkDaysTotal:     22,
				WorkDaysSickLeave: 2,
				TaxesConfig:       GetDefaultTaxesConfig(),
			},
			user: UserConfig{
				FirstName:  "Иван",
				MiddleName: "Петров",
				LastName:   "Иванов",
				Egn:        "1234567890",
				Bulstat:    "987654321",
			},
			settings: Settings{
				IsPregnancyInsuranceEnabled: false,
			},
			wantErr:  false,
			expected: readFileOrPanic("testdata/declaration_one_valid.txt"),
		},
		{
			name: "Negative Work Days Real",
			form: IncomeForm{
				Month:             7,
				Year:              2026,
				DayStart:          1,
				DayEnd:            31,
				WorkDaysTotal:     10,
				WorkDaysSickLeave: 15,
			},
			user:     UserConfig{},
			settings: Settings{},
			wantErr:  true,
			errMsg:   "Невалидна стойност на работни дни: -5",
		},
		{
			name: "Work Days Real Over Max Limit",
			form: IncomeForm{
				Month:             7,
				Year:              2026,
				DayStart:          1,
				DayEnd:            31,
				WorkDaysTotal:     40,
				WorkDaysSickLeave: 5,
			},
			user:     UserConfig{},
			settings: Settings{},
			wantErr:  true,
			errMsg:   "Невалидна стойност на работни дни: 35",
		},
		{
			name: "Zero Taxed Income",
			form: IncomeForm{
				Month:             7,
				Year:              2026,
				TaxedIncomeCents:  0,
				DayStart:          1,
				DayEnd:            31,
				WorkDaysTotal:     22,
				WorkDaysSickLeave: 0,
				TaxesConfig: TaxesConfig{
					PensionPercentagePartOne:     5.0,
					PensionPercentagePartTwo:     2.5,
					PregnancyInsurancePercentage: 3.0,
					HealthInsurancePercentage:    8.0,
				},
			},
			user: UserConfig{},
			settings: Settings{
				IsPregnancyInsuranceEnabled: true,
			},
			wantErr: true,
			errMsg:  "Нулев осигурителен доход, декларацията не се подава",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MakeDeclarationOne(tt.form, tt.user, tt.settings)

			if (err != nil) != tt.wantErr {
				t.Errorf("MakeDeclarationOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("MakeDeclarationOne() error = %v, expected error message = %v", err.Error(), tt.errMsg)
			}

			if !tt.wantErr && !bytes.Equal(result, tt.expected) {
				t.Errorf(
					"MakeDeclarationOne() result is not equal to expected. Diff from expected: %v",
					cmp.Diff(tt.expected, result),
				)
			}
		})
	}
}

func readFileOrPanic(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}
