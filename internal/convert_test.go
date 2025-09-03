package internal

import (
	"encoding/json"
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestReformatSchedule(t *testing.T) {
	tests := []struct {
		name             string
		inpusJsonFile    string
		expectedJsonFile string
	}{
		{"test", "./../test/data/schedule_case1_in.json", "./../test/data/schedule_case1_out.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inJsonTextBytes, err := os.ReadFile(tt.inpusJsonFile)
			assert.NoError(t, err, "Error reading local file")
			outJsonTextBytes, err := os.ReadFile(tt.expectedJsonFile)
			assert.NoError(t, err, "Error reading local file")
			var input domain.ScheduleDataJSON
			err = json.Unmarshal(inJsonTextBytes, &input)
			var output domain.Schedule
			err = json.Unmarshal(outJsonTextBytes, &output)
			assert.NoError(t, err, "failed parsing json %v")
			reformattedSchedule, err := ReformatSchedule(input)
			reformattedScheduleJSON, err := json.Marshal(reformattedSchedule)
			if diff := cmp.Diff(output, reformattedSchedule); diff != "" {
				fmt.Println(string(reformattedScheduleJSON))
				t.Errorf("ReformatSchedule() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
