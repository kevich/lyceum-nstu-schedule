package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"kevich/lyceum-nstu-schedule/domain"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	jsonTextBytes, err := os.ReadFile("./test/data/requests/empty_first_request.json")
	assert.NoError(t, err, "Error reading local file")
	var event domain.Event
	err = json.Unmarshal(jsonTextBytes, &event)
	assert.NoError(t, err, "failed parsing json %v")
	var response *domain.Response
	response, err = Handler("", event)
	if response != nil {
		assert.Equal(t, "Привет, я могу рассказать расписание инженерного лицея НГТУ. Расписание какого класса и в какой день вас интересует?", response.Response.Text)
	} else {
		t.Errorf("Response should be returned")
	}
}
