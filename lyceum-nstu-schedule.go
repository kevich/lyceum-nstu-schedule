package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
	"kevich/lyceum-nstu-schedule/internal"
	"kevich/lyceum-nstu-schedule/internal/api"
	"kevich/lyceum-nstu-schedule/tools"
	"net/http"
	"os"
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

func fetchDataAndReturnAliceResponse(class string, date string) string {
	apiClient := api.ScheduleAPI{Client: &http.Client{}, BaseURL: api.BaseUrl}
	return fetchDataAndReturnAliceResponseWithClient(class, date, &apiClient)
}

func fetchDataAndReturnAliceResponseWithClient(class string, date string, apiClient *api.ScheduleAPI) string {
	jsonData, err := apiClient.ApiGetData()
	tools.CheckError(err, "failed getting json %v")
	var input domain.ScheduleDataJSON
	err = json.Unmarshal(jsonData, &input)
	tools.CheckError(err, "failed parsing json %v")
	reformatted, err := internal.ReformatSchedule(input)
	tools.CheckError(err, "could not reformat schedule")
	my := reformatted[class][date]
	return internal.FormatDayToAliceResponse(date, my)
}

type Config struct {
	Class string
	Date  string
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("schedule", flag.ContinueOnError)
	class := fs.String("class", "7а", "schedule for a class")
	date := fs.String("date", "", "schedule for a date")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return &Config{
		Class: *class,
		Date:  *date,
	}, nil
}

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}
	fmt.Println(fetchDataAndReturnAliceResponse(config.Class, config.Date))
}
