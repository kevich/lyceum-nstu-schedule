package internal

import (
	"fmt"
	"io"
	"kevich/lyceum-nstu-schedule/tools"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const baseUrl = "https://lyceum.nstu.ru/"

// const localFile = ""

const localFile = "./demo.json"

func ApiGetData() []byte {
	var jsonText string
	if localFile != "" {
		jsonTextBytes, err := os.ReadFile(localFile)
		tools.CheckError(err, "Error reading local file")
		jsonText = string(jsonTextBytes)
	} else {
		response, err := http.Get(baseUrl + "/rasp/schedule.html")
		tools.CheckError(err, "Could not get schedule")
		if response.Body == nil {
			fmt.Println("Could not get schedule body")
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
			fmt.Println("Could not find schedule file")
		}
		file := string(matches[1])

		response, err = http.Get(baseUrl + "rasp/" + file)
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
			fmt.Println("Could not find schedule data")
		}
		jsonText = string(matches[1])
		fmt.Println(jsonText)
		os.WriteFile("./demo.json", matches[1], 0666)
	}

	return []byte(strings.Trim(jsonText, "\n "))
}
