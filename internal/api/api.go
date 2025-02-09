package api

import (
	"fmt"
	"io"
	"kevich/lyceum-nstu-schedule/tools"
	"net/http"
	"regexp"
	"strings"
)

const BaseUrl = "https://lyceum.nstu.ru/"

type ScheduleAPI struct {
	Client  *http.Client
	BaseURL string
}

func (internal *ScheduleAPI) ApiGetData() ([]byte, error) {
	var jsonText string
	response, err := internal.Client.Get(internal.BaseURL + "/rasp/schedule.html")
	tools.CheckError(err, "Could not get schedule")
	if response.Body == nil {
		return nil, fmt.Errorf("could not get schedule body")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		tools.CheckError(err, "Could not get schedule body")
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	tools.CheckError(err, "Could not read schedule")
	re := regexp.MustCompile(`src="(nika_data.+js)`)
	matches := re.FindSubmatch(body)
	if matches == nil {
		return nil, fmt.Errorf("could not find schedule file")
	}
	file := string(matches[1])

	response, err = internal.Client.Get(internal.BaseURL + "/rasp/" + file)
	tools.CheckError(err, "Could not get schedule")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		tools.CheckError(err, "Could not get schedule body")
	}(response.Body)
	body, err = io.ReadAll(response.Body)
	tools.CheckError(err, "Could not read schedule")
	re = regexp.MustCompile(`(?s)NIKA=(.+);`)
	matches = re.FindSubmatch(body)
	if matches == nil {
		return nil, fmt.Errorf("could not find schedule data")
	}
	jsonText = string(matches[1])
	fmt.Println(jsonText)

	return []byte(strings.Trim(jsonText, "\n ")), err
}
