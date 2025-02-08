package internal

import (
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
	"kevich/lyceum-nstu-schedule/tools"
	"strconv"
	"time"
)

func isCancelledSubject[T any](s T) bool {
	if s, ok := any(s).(string); ok {
		return s == "F"
	}
	return false
}

func ReformatSchedule(jsonData domain.ScheduleDataJSON) domain.Schedule {
	schedule := make(domain.Schedule)
	nowDateString := time.Now().Format("2006-01-02")
	todaysDay, err := time.Parse("2006-01-02", nowDateString)
	tools.CheckError(err, "Could not get current date")
	todayWeekday := time.Duration(todaysDay.Weekday()) - 1
	firstDayOfWeek := todaysDay.Add(-24 * todayWeekday * time.Hour)
	for classKey, class := range jsonData.CLASS_SCHEDULE.School {
		className := jsonData.CLASSES[classKey]
		schedule[className] = make(domain.ClassSchedule)
		for dayNumber := 1; dayNumber < 7; dayNumber++ {
			dayOfWeek := firstDayOfWeek.Add(24 * time.Duration(dayNumber-1) * time.Hour)
			dayOfWeekString := dayOfWeek.Format("02.01.2006")
			schedule[className][dayOfWeekString] = make(domain.DaySchedule)
			for lessonNumber := 1; lessonNumber <= jsonData.LESSONSINDAY; lessonNumber++ {
				dayLesson := fmt.Sprintf("%d%.2d", dayNumber, lessonNumber)
				jsonClassDaySchedule, found := class[dayLesson]
				if found != false {
					for i := range jsonClassDaySchedule.Subject {
						replacement := jsonData.CLASS_EXCHANGE[classKey][dayOfWeekString][strconv.Itoa(lessonNumber)]
						classDaySchedule := jsonClassDaySchedule
						subject := jsonData.SUBJECTS[classDaySchedule.Subject[i]]
						teacher := jsonData.TEACHERS[classDaySchedule.Teacher[i]]
						roomNumber := jsonData.ROOMS[classDaySchedule.Room[i]]
						isCancelled := false
						if isCancelledSubject(replacement.Subject) {
							isCancelled = true
						} else {
							if replacementSubject, ok := replacement.Subject.([]string); ok {
								subject = jsonData.SUBJECTS[replacementSubject[i]]
							}
							if replacement.Teacher != nil {
								teacher = jsonData.TEACHERS[replacement.Teacher[i]]
							}
							if replacement.Room != nil {
								roomNumber = jsonData.TEACHERS[replacement.Room[i]]
							}
						}

						if i == 0 {
							schedule[className][dayOfWeekString][lessonNumber] = make([]domain.Lesson, len(classDaySchedule.Subject))
						}
						schedule[className][dayOfWeekString][lessonNumber][i] = domain.Lesson{
							Number:      lessonNumber,
							TimeRange:   jsonData.LESSON_TIMES[strconv.Itoa(lessonNumber)],
							Name:        subject,
							Teacher:     teacher,
							RoomNumber:  roomNumber,
							IsCancelled: isCancelled,
						}
					}

				}
			}
		}
	}

	return schedule
}
