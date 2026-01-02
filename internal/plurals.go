package internal

import (
	"kevich/lyceum-nstu-schedule/tools"
	"time"

	"github.com/goodsign/monday"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LessonsFormatter() *message.Printer {
	message.Set(language.Russian, "%d уроков", plural.Selectf(
		1, "%d",
		"=1", "один урок",
		"=2", "два урока",
		"=3", "три урока",
		"=4", "четыре урока",
		"other", "%[1]d уроков",
	))
	return message.NewPrinter(language.Russian)
}

func FormatDate(dateString string) string {
	date, err := time.Parse("02.01.2006", dateString)
	tools.CheckError(err, "Could not get date")
	return monday.Format(date, "2 January", monday.LocaleRuRU)
}
