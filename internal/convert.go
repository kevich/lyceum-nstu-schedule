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

func getReplacementSubject(replacement interface{}) []string {
	value := reflect.ValueOf(replacement)

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil
	}

	var result []string
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func getLesson(jsonData domain.ScheduleDataJSON, lessonNumber int, timeRange [2]string, subject []string, teacher []string, roomNumber []string, isCancelled bool, isReplaced bool) []domain.Lesson {
	lessons := make([]domain.Lesson, len(subject))

	for i := range subject {
		lessons[i] = domain.Lesson{
			Number:      lessonNumber,
			TimeRange:   timeRange,
			Name:        jsonData.SUBJECTS[subject[i]],
			Teacher:     jsonData.TEACHERS[teacher[i]],
			RoomNumber:  jsonData.ROOMS[roomNumber[i]],
			IsCancelled: isCancelled,
			IsReplaced:  isReplaced,
		}
	}

	return lessons
}

func ReformatSchedule(jsonData domain.ScheduleDataJSON) (domain.Schedule, error) {
	schedule := make(domain.Schedule)
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
	replacement := jsonData.CLASS_EXCHANGE[classKey][dayOfWeekString][strconv.Itoa(lessonNumber)]
	var (
		subject     []string
		teacher     []string
		room        []string
		isCancelled bool
		isReplaced  bool
	)

	switch {
	case replacement.Subject != nil:
		if isCancelledSubject(replacement.Subject) {
			subject = jsonClassDaySchedule.Subject
			teacher = jsonClassDaySchedule.Teacher
			room = jsonClassDaySchedule.Room
			isCancelled = true
		} else {
			subject = getReplacementSubject(replacement.Subject)
			teacher = replacement.Teacher
			room = replacement.Room
			isReplaced = true
		}

	case found:
		subject = jsonClassDaySchedule.Subject
		teacher = jsonClassDaySchedule.Teacher
		room = jsonClassDaySchedule.Room

	default:
		return
	}

	schedule[className][dayOfWeekString][lessonNumber] = getLesson(
		jsonData,
		lessonNumber,
		jsonData.LESSON_TIMES[strconv.Itoa(lessonNumber)],
		subject,
		teacher,
		room,
		isCancelled,
		isReplaced,
	)
}
