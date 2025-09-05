package main

import (
	"context"
	"encoding/json"
	"kevich/lyceum-nstu-schedule/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	jsonTextBytes, err := os.ReadFile("./test/data/requests/empty_first_request.json")
	assert.NoError(t, err, "Error reading local file")
	var event domain.Event
	err = json.Unmarshal(jsonTextBytes, &event)
	assert.NoError(t, err, "failed parsing json %v")
	var response *domain.Response
	response, err = Handler(context.Background(), event)
	if response != nil {
		assert.Equal(t, "Привет, я могу рассказать расписание инженерного лицея НГТУ. Расписание какого класса и в какой день вас интересует?", response.Response.Text)
	} else {
		t.Errorf("Response should be returned")
	}
}
