package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"kevich/lyceum-nstu-schedule/domain"
	"kevich/lyceum-nstu-schedule/tools"
	"os"
	"testing"
)

func TestHandler(T *testing.T) {
	jsonTextBytes, err := os.ReadFile("./test/data/requests/empty_first_request.json")
	tools.CheckError(err, "Error reading local file")
	var event domain.Event
	err = json.Unmarshal(jsonTextBytes, &event)
	tools.CheckError(err, "failed parsing json %v")
	var response *domain.Response
	response, err = Handler("", event)
	if response != nil {
		assert.Equal(T, "Привет, я могу рассказать расписание инженерного лицея НГТУ. Расписание какого класса и в какой день вас интересует?", response.Response.Text)
	} else {
		T.Errorf("Response should be returned")
	}
}
