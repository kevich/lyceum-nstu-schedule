package internal

import (
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
	"kevich/lyceum-nstu-schedule/tools"
	"reflect"
	"strconv"
	"time"
)

func isCancelledSubject[T any](s T) bool {
	if s, ok := any(s).(string); ok {
		return s == "F"
	}
	return false
}

func ReformatSchedule(jsonData domain.ScheduleDataJSON) (domain.Schedule, error) {
	schedule := make(domain.Schedule)
	// nowDateString := time.Now().Format("2006-01-02")
	// todayDay, err := time.Parse("2006-01-02", nowDateString)
	// tools.CheckError(err, "Could not get current date")
	// todayWeekday := time.Duration(todayDay.Weekday()) - 1
	// firstDayOfWeek := todayDay.Add(-24 * todayWeekday * time.Hour)
	periods := reflect.ValueOf(jsonData.PERIODS).MapKeys()
	if len(periods) == 0 {
		return nil, fmt.Errorf("could not find schedule data")
	}
	for _, period := range periods {
		for classKey, class := range jsonData.CLASS_SCHEDULE[period.String()] {
			startDate, err := time.Parse("02.01.2006", jsonData.PERIODS[period.String()].StartDate)
			endDate, err := time.Parse("02.01.2006", jsonData.PERIODS[period.String()].EndDate)
			tools.CheckError(err, "Could not get start date")
			reformatWeek(jsonData, classKey, schedule, startDate, endDate, class)
		}
	}

	return schedule, nil
}

func reformatWeek(jsonData domain.ScheduleDataJSON, classKey string, schedule domain.Schedule, startDate time.Time, endDate time.Time, class domain.ClassScheduleJSON) {
	className := jsonData.CLASSES[classKey]
	schedule[className] = make(domain.ClassSchedule)
	currentDate := startDate
	for currentDate.Before(endDate) {
		reformatDay(jsonData, classKey, schedule, currentDate, class, className)
		currentDate = currentDate.AddDate(0, 0, 1)
	}
}

func reformatDay(jsonData domain.ScheduleDataJSON, classKey string, schedule domain.Schedule, currentDay time.Time, class domain.ClassScheduleJSON, className string) {
	dayNumber := int(currentDay.Weekday())
	dayOfWeekString := currentDay.Format("02.01.2006")
	schedule[className][dayOfWeekString] = make(domain.DaySchedule)
	for lessonNumber := 1; lessonNumber <= jsonData.LESSONSINDAY; lessonNumber++ {
		reformatLesson(jsonData, classKey, schedule, dayNumber, lessonNumber, class, dayOfWeekString, className)
	}
}

func reformatLesson(jsonData domain.ScheduleDataJSON, classKey string, schedule domain.Schedule, dayNumber int, lessonNumber int, class domain.ClassScheduleJSON, dayOfWeekString string, className string) {
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
