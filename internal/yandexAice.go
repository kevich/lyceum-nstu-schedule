package internal

import (
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
)

func FormatDayToAliceResponse(day string, daySchedule domain.DaySchedule) string {
	dayString := ""
	lessonFormatter := LessonsFormatter()
	xLessons := lessonFormatter.Sprintf("%d уроков", len(daySchedule))
	dayString += fmt.Sprintf("%s будет %s\n", FormatDate(day), xLessons)
	dayString += fmt.Sprintf("Уроки начинаются в %s\n", daySchedule[0][0].TimeRange[0])

	for _, lessons := range daySchedule {
		for index, lesson := range lessons {
			name := lesson.Name
			if name == "" {
				name = "нет урока"
			}
			if lesson.IsCancelled {
				name += " - отменен"
			}
			if len(lessons) > 1 {
				dayString += fmt.Sprintf("%d группа - %s\n", index+1, name)
			} else {
				dayString += fmt.Sprintf("%s\n", name)
			}

		}
	}

	dayString += fmt.Sprintf("Уроки закончатся в %s\n", daySchedule[len(daySchedule)-1][0].TimeRange[1])
	return dayString
}
