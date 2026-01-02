package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessonsFormatter(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{"one lesson", 1, "один урок"},
		{"two lessons", 2, "два урока"},
		{"three lessons", 3, "три урока"},
		{"four lessons", 4, "четыре урока"},
		{"five lessons", 5, "5 уроков"},
		{"ten lessons", 10, "10 уроков"},
		{"twenty lessons", 20, "20 уроков"},
		{"zero lessons", 0, "0 уроков"},
		{"twenty one lessons", 21, "21 уроков"},
		{"hundred lessons", 100, "100 уроков"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LessonsFormatter()
			result := formatter.Sprintf("%d уроков", tt.count)
			assert.Equal(t, tt.expected, result, "LessonsFormatter() mismatch")
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name       string
		dateString string
		expected   string
	}{
		{"January date", "15.01.2024", "15 января"},
		{"February date", "02.02.2024", "2 февраля"},
		{"March date", "10.03.2024", "10 марта"},
		{"April date", "01.04.2024", "1 апреля"},
		{"May date", "25.05.2024", "25 мая"},
		{"June date", "30.06.2024", "30 июня"},
		{"July date", "14.07.2024", "14 июля"},
		{"August date", "08.08.2024", "8 августа"},
		{"September date", "01.09.2024", "1 сентября"},
		{"October date", "31.10.2024", "31 октября"},
		{"November date", "15.11.2024", "15 ноября"},
		{"December date", "25.12.2024", "25 декабря"},
		{"single digit day", "05.01.2024", "5 января"},
		{"double digit day", "15.01.2024", "15 января"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDate(tt.dateString)
			assert.Equal(t, tt.expected, result, "FormatDate() mismatch for date %s", tt.dateString)
		})
	}
}
