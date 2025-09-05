package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
	"kevich/lyceum-nstu-schedule/internal"
	"kevich/lyceum-nstu-schedule/internal/api"
	"kevich/lyceum-nstu-schedule/tools"
	"net/http"
)

func Handler(ctx context.Context, event domain.Event) (*domain.Response, error) {

	text := "Привет, я могу рассказать расписание инженерного лицея НГТУ. Расписание какого класса и в какой день вас интересует?"
	if event.Request.OriginalUtterance != "" {
		text = event.Request.OriginalUtterance
	}

	return &domain.Response{
		Version: event.Version,
		Session: domain.ResponseSession{
			SessionID: event.Session.SessionID,
			MessageID: event.Session.MessageID,
			UserID:    event.Session.UserID,
		},
		Response: domain.ResponsePayload{
			Text:       text,
			EndSession: false,
		},
	}, nil
}

func main() {
	apiClient := api.ScheduleAPI{Client: &http.Client{}, BaseURL: api.BaseUrl}
	jsonData, err := apiClient.ApiGetData()
	tools.CheckError(err, "failed getting json %v")
	var input domain.ScheduleDataJSON
	err = json.Unmarshal(jsonData, &input)
	tools.CheckError(err, "failed parsing json %v")
	reformatted, err := internal.ReformatSchedule(input)
	tools.CheckError(err, "could not reformat schedule")
	my := reformatted["6а"]["04.02.2025"]
	fmt.Println(reformatted)
	fmt.Println(my)
}
