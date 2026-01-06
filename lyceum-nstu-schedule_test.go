package main

import (
	"fmt"
	"kevich/lyceum-nstu-schedule/internal/api"
	"kevich/lyceum-nstu-schedule/tools"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const scheduleHtmlTemplate = `<!DOCTYPE html>
<html><head>
<script type="text/javascript" src="nika_data_test.js"></script>
</head><body></body></html>`

const scheduleJSONTemplate = `{
	"LESSONSINDAY": 12,
	"DAY_NAMES": ["Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"],
	"TEACHERS": {"1": "Иванов И.И.", "2": "Петров П.П."},
	"SUBJECTS": {"1": "Математика", "2": "Русский язык", "3": "Физика"},
	"CLASSES": {"1": "7а", "2": "8б"},
	"ROOMS": {"1": "101", "2": "202", "3": "303"},
	"CLASSGROUPS": {},
	"LESSON_TIMES": {
		"1": ["08:30", "09:15"],
		"2": ["09:25", "10:10"],
		"3": ["10:20", "11:05"]
	},
	"CLASS_SCHEDULE": {
		"1": {
			"1": {
				"101": {"s": ["1"], "t": ["1"], "g": [], "r": ["1"]},
				"102": {"s": ["2"], "t": ["2"], "g": [], "r": ["2"]},
				"103": {"s": ["3"], "t": ["1"], "g": [], "r": ["3"]}
			}
		}
	},
	"PERIODS": {
		"1": {"b": "05.01.2026", "e": "12.01.2026"}
	},
	"CLASS_EXCHANGE": {}
}`

func createTestServer(scheduleJSON string) *httptest.Server {
	rawJS := `// nika_data.js
var NIKA=%s;`
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/rasp/schedule.html", func(rw http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(rw, scheduleHtmlTemplate)
		tools.CheckError(err, "failed returning response")
	})
	serverMux.HandleFunc("/rasp/nika_data_test.js", func(rw http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(rw, fmt.Sprintf(rawJS, scheduleJSON))
		tools.CheckError(err, "failed returning response")
	})
	return httptest.NewServer(serverMux)
}

func TestFetchDataAndReturnAliceResponse(t *testing.T) {
	tests := []struct {
		name           string
		class          string
		date           string
		scheduleJSON   string
		expectedOutput string
	}{
		{
			name:         "returns schedule for 7а on Monday",
			class:        "7а",
			date:         "05.01.2026",
			scheduleJSON: scheduleJSONTemplate,
			expectedOutput: `5 января будет три урока
Уроки начинаются в 08:30
Математика
Русский язык
Физика
Уроки закончатся в 11:05
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createTestServer(tt.scheduleJSON)
			defer server.Close()

			apiClient := &api.ScheduleAPI{Client: server.Client(), BaseURL: server.URL}
			result := fetchDataAndReturnAliceResponseWithClient(tt.class, tt.date, apiClient)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func TestFetchDataAndReturnAliceResponseWithReplacement(t *testing.T) {
	scheduleWithReplacement := `{
		"LESSONSINDAY": 12,
		"DAY_NAMES": ["Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"],
		"TEACHERS": {"1": "Иванов И.И.", "2": "Петров П.П."},
		"SUBJECTS": {"1": "Математика", "2": "Русский язык", "3": "Физика"},
		"CLASSES": {"1": "7а"},
		"ROOMS": {"1": "101", "2": "202"},
		"CLASSGROUPS": {},
		"LESSON_TIMES": {
			"1": ["08:30", "09:15"],
			"2": ["09:25", "10:10"]
		},
		"CLASS_SCHEDULE": {
			"1": {
				"1": {
					"101": {"s": ["1"], "t": ["1"], "g": [], "r": ["1"]},
					"102": {"s": ["2"], "t": ["2"], "g": [], "r": ["2"]}
				}
			}
		},
		"PERIODS": {
			"1": {"b": "05.01.2026", "e": "12.01.2026"}
		},
		"CLASS_EXCHANGE": {
			"1": {
				"05.01.2026": {
					"1": {"s": ["3"], "t": ["2"], "g": [], "r": ["2"]}
				}
			}
		}
	}`

	server := createTestServer(scheduleWithReplacement)
	defer server.Close()

	apiClient := &api.ScheduleAPI{Client: server.Client(), BaseURL: server.URL}
	result := fetchDataAndReturnAliceResponseWithClient("7а", "05.01.2026", apiClient)

	assert.Contains(t, result, "Физика", "Should contain replaced lesson (Физика instead of Математика)")
	assert.Contains(t, result, "Русский язык", "Should contain regular lesson")
	assert.Contains(t, result, "два урока", "Should have 2 lessons")
}

func TestFetchDataAndReturnAliceResponseWithCancellation(t *testing.T) {
	scheduleWithCancellation := `{
		"LESSONSINDAY": 12,
		"DAY_NAMES": ["Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"],
		"TEACHERS": {"1": "Иванов И.И.", "2": "Петров П.П."},
		"SUBJECTS": {"1": "Математика", "2": "Русский язык"},
		"CLASSES": {"1": "7а"},
		"ROOMS": {"1": "101", "2": "202"},
		"CLASSGROUPS": {},
		"LESSON_TIMES": {
			"1": ["08:30", "09:15"],
			"2": ["09:25", "10:10"]
		},
		"CLASS_SCHEDULE": {
			"1": {
				"1": {
					"101": {"s": ["1"], "t": ["1"], "g": [], "r": ["1"]},
					"102": {"s": ["2"], "t": ["2"], "g": [], "r": ["2"]}
				}
			}
		},
		"PERIODS": {
			"1": {"b": "05.01.2026", "e": "12.01.2026"}
		},
		"CLASS_EXCHANGE": {
			"1": {
				"05.01.2026": {
					"1": {"s": "F", "t": [], "g": [], "r": []}
				}
			}
		}
	}`

	server := createTestServer(scheduleWithCancellation)
	defer server.Close()

	apiClient := &api.ScheduleAPI{Client: server.Client(), BaseURL: server.URL}
	result := fetchDataAndReturnAliceResponseWithClient("7а", "05.01.2026", apiClient)

	assert.Contains(t, result, "отменен", "Should indicate cancelled lesson")
	assert.Contains(t, result, "Математика", "Should show original subject name")
}

func TestFetchDataAndReturnAliceResponseMultipleClasses(t *testing.T) {
	scheduleMultipleClasses := `{
		"LESSONSINDAY": 12,
		"DAY_NAMES": ["Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"],
		"TEACHERS": {"1": "Иванов И.И.", "2": "Петров П.П."},
		"SUBJECTS": {"1": "Математика", "2": "Физика"},
		"CLASSES": {"1": "7а", "2": "8б"},
		"ROOMS": {"1": "101", "2": "202"},
		"CLASSGROUPS": {},
		"LESSON_TIMES": {
			"1": ["08:30", "09:15"]
		},
		"CLASS_SCHEDULE": {
			"1": {
				"1": {
					"101": {"s": ["1"], "t": ["1"], "g": [], "r": ["1"]}
				},
				"2": {
					"101": {"s": ["2"], "t": ["2"], "g": [], "r": ["2"]}
				}
			}
		},
		"PERIODS": {
			"1": {"b": "05.01.2026", "e": "12.01.2026"}
		},
		"CLASS_EXCHANGE": {}
	}`

	server := createTestServer(scheduleMultipleClasses)
	defer server.Close()

	apiClient := &api.ScheduleAPI{Client: server.Client(), BaseURL: server.URL}

	result7a := fetchDataAndReturnAliceResponseWithClient("7а", "05.01.2026", apiClient)
	assert.Contains(t, result7a, "Математика", "7а should have Математика")

	result8b := fetchDataAndReturnAliceResponseWithClient("8б", "05.01.2026", apiClient)
	assert.Contains(t, result8b, "Физика", "8б should have Физика")
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name:           "default values when no flags provided",
			args:           []string{},
			expectedConfig: &Config{Class: "7а", Date: ""},
			expectError:    false,
		},
		{
			name:           "custom class flag",
			args:           []string{"-class", "8б"},
			expectedConfig: &Config{Class: "8б", Date: ""},
			expectError:    false,
		},
		{
			name:           "custom date flag",
			args:           []string{"-date", "05.01.2026"},
			expectedConfig: &Config{Class: "7а", Date: "05.01.2026"},
			expectError:    false,
		},
		{
			name:           "both class and date flags",
			args:           []string{"-class", "9в", "-date", "10.02.2026"},
			expectedConfig: &Config{Class: "9в", Date: "10.02.2026"},
			expectError:    false,
		},
		{
			name:           "flags with equals sign syntax",
			args:           []string{"-class=10г", "-date=15.03.2026"},
			expectedConfig: &Config{Class: "10г", Date: "15.03.2026"},
			expectError:    false,
		},
		{
			name:           "unknown flag returns error",
			args:           []string{"-unknown", "value"},
			expectedConfig: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseFlags(tt.args)

			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid flags")
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, config)
			}
		})
	}
}
