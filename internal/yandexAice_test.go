package internal

import (
	"context"
	"encoding/json"
	"kevich/lyceum-nstu-schedule/domain"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name             string
		inputJsonFile    string
		expectedResponse string
		mockFetcher      ScheduleFetcher
		mockTime         time.Time
	}{
		{
			name:             "empty first request returns greeting",
			inputJsonFile:    "../test/data/requests/empty_first_request.json",
			expectedResponse: GreetingText,
			mockFetcher:      nil,
		},
		{
			name:          "class_and_date intent returns schedule",
			inputJsonFile: "../test/data/requests/second_request_6a_tomorrow.json",
			expectedResponse: `7 января будет 5 уроков
Уроки начинаются в 08:30
Математика
Физика
Русский язык
История
Литература
Уроки закончатся в 13:05
`,
			mockFetcher: func(class string, date string) string {
				assert.Equal(t, "6а", class)
				assert.Equal(t, "07.01.2026", date) // tomorrow from 06.01.2026
				return `7 января будет 5 уроков
Уроки начинаются в 08:30
Математика
Физика
Русский язык
История
Литература
Уроки закончатся в 13:05
`
			},
			mockTime: time.Date(2026, 1, 6, 10, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonTextBytes, err := os.ReadFile(test.inputJsonFile)
			assert.NoError(t, err, "Error reading local file")
			var event domain.Event
			err = json.Unmarshal(jsonTextBytes, &event)
			assert.NoError(t, err, "failed parsing json %v")

			handler := NewAliceHandler(test.mockFetcher)
			if !test.mockTime.IsZero() {
				handler.TimeNowFunc = func() time.Time { return test.mockTime }
			}

			response, err := handler.Handle(context.Background(), event)
			assert.NoError(t, err)
			if response != nil {
				assert.Equal(t, test.expectedResponse, response.Response.Text)
			} else {
				t.Errorf("Response should be returned")
			}
		})
	}
}

func TestFormatDayToAliceResponse(t *testing.T) {
	tests := []struct {
		name        string
		day         string
		daySchedule domain.DaySchedule
		expected    string
	}{
		{
			name: "single lesson",
			day:  "15.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Математика"},
				},
			},
			expected: "15 января будет один урок\n" +
				"Уроки начинаются в 08:30\n" +
				"Математика\n" +
				"Уроки закончатся в 09:15\n",
		},
		{
			name: "multiple lessons",
			day:  "16.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Математика"},
				},
				{
					{Number: 2, TimeRange: [2]string{"09:25", "10:10"}, Name: "Физика"},
				},
				{
					{Number: 3, TimeRange: [2]string{"10:20", "11:05"}, Name: "Русский язык"},
				},
			},
			expected: "16 января будет три урока\n" +
				"Уроки начинаются в 08:30\n" +
				"Математика\n" +
				"Физика\n" +
				"Русский язык\n" +
				"Уроки закончатся в 11:05\n",
		},
		{
			name: "lesson with cancelled flag",
			day:  "17.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Математика"},
				},
				{
					{Number: 2, TimeRange: [2]string{"09:25", "10:10"}, Name: "Физика", IsCancelled: true},
				},
			},
			expected: "17 января будет два урока\n" +
				"Уроки начинаются в 08:30\n" +
				"Математика\n" +
				"Физика - отменен\n" +
				"Уроки закончатся в 10:10\n",
		},
		{
			name: "lesson with empty name",
			day:  "18.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Математика"},
				},
				{
					{Number: 2, TimeRange: [2]string{"09:25", "10:10"}, Name: ""},
				},
			},
			expected: "18 января будет два урока\n" +
				"Уроки начинаются в 08:30\n" +
				"Математика\n" +
				"нет урока\n" +
				"Уроки закончатся в 10:10\n",
		},
		{
			name: "lesson with multiple groups",
			day:  "19.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Английский язык"},
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Немецкий язык"},
				},
				{
					{Number: 2, TimeRange: [2]string{"09:25", "10:10"}, Name: "Информатика"},
				},
			},
			expected: "19 января будет два урока\n" +
				"Уроки начинаются в 08:30\n" +
				"1 группа - Английский язык\n" +
				"2 группа - Немецкий язык\n" +
				"Информатика\n" +
				"Уроки закончатся в 10:10\n",
		},
		{
			name: "five lessons",
			day:  "20.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Математика"},
				},
				{
					{Number: 2, TimeRange: [2]string{"09:25", "10:10"}, Name: "Физика"},
				},
				{
					{Number: 3, TimeRange: [2]string{"10:20", "11:05"}, Name: "Русский язык"},
				},
				{
					{Number: 4, TimeRange: [2]string{"11:25", "12:10"}, Name: "История"},
				},
				{
					{Number: 5, TimeRange: [2]string{"12:20", "13:05"}, Name: "Литература"},
				},
			},
			expected: "20 января будет 5 уроков\n" +
				"Уроки начинаются в 08:30\n" +
				"Математика\n" +
				"Физика\n" +
				"Русский язык\n" +
				"История\n" +
				"Литература\n" +
				"Уроки закончатся в 13:05\n",
		},
		{
			name: "empty lesson name with cancelled flag",
			day:  "21.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "", IsCancelled: true},
				},
			},
			expected: "21 января будет один урок\n" +
				"Уроки начинаются в 08:30\n" +
				"нет урока - отменен\n" +
				"Уроки закончатся в 09:15\n",
		},
		{
			name: "multiple groups with cancelled lesson",
			day:  "22.01.2024",
			daySchedule: domain.DaySchedule{
				{
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Английский язык", IsCancelled: true},
					{Number: 1, TimeRange: [2]string{"08:30", "09:15"}, Name: "Немецкий язык"},
				},
			},
			expected: "22 января будет один урок\n" +
				"Уроки начинаются в 08:30\n" +
				"1 группа - Английский язык - отменен\n" +
				"2 группа - Немецкий язык\n" +
				"Уроки закончатся в 09:15\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDayToAliceResponse(tt.day, tt.daySchedule)
			assert.Equal(t, tt.expected, result, "FormatDayToAliceResponse() mismatch")
		})
	}
}
