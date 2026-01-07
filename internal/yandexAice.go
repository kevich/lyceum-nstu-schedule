package internal

import (
	"context"
	"fmt"
	"kevich/lyceum-nstu-schedule/domain"
	"time"
)

type ScheduleFetcher func(class string, date string) string

type AliceHandler struct {
	Fetcher     ScheduleFetcher
	TimeNowFunc func() time.Time // For testing, defaults to time.Now
}

func NewAliceHandler(fetcher ScheduleFetcher) *AliceHandler {
	return &AliceHandler{
		Fetcher:     fetcher,
		TimeNowFunc: time.Now,
	}
}

const GreetingText = "Привет, я могу рассказать расписание инженерного лицея НГТУ. Расписание какого класса и в какой день вас интересует?"

func (h *AliceHandler) Handle(ctx context.Context, event domain.Event) (*domain.Response, error) {
	text := GreetingText
	if event.Session.MessageID != 0 {
		text = "Простите, я не смогла распознать команду. Попробуйте еще раз."
	}

	// Check if we have a class_and_date intent
	intent := event.Request.NLU.Intents.ClassAndDate
	if intent.Slots.ClassNumber.Value != nil && intent.Slots.ClassCharacter.Value != nil {
		className := h.extractClassName(intent.Slots)
		date := h.resolveDate(event.Meta.Timezone, intent.Slots.Date)

		// If no date slot, try to resolve from day_of_week slot
		if date == "" {
			date = h.resolveDayOfWeek(event.Meta.Timezone, intent.Slots.DayOfWeek)
		}

		if className != "" && date != "" && h.Fetcher != nil {
			text = h.Fetcher(className, date)
		}
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

func (h *AliceHandler) extractClassName(slots domain.SlotsClassAndDate) string {
	classNum, ok := slots.ClassNumber.Value.(float64)
	if !ok {
		return ""
	}
	classChar, ok := slots.ClassCharacter.Value.(string)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d%s", int(classNum), classChar)
}

var russianWeekdays = map[string]time.Weekday{
	"понедельник": time.Monday,
	"вторник":     time.Tuesday,
	"среда":       time.Wednesday,
	"среду":       time.Wednesday,
	"четверг":     time.Thursday,
	"пятница":     time.Friday,
	"пятницу":     time.Friday,
	"суббота":     time.Saturday,
	"субботу":     time.Saturday,
	"воскресенье": time.Sunday,
}

func (h *AliceHandler) resolveDayOfWeek(timezone string, dayOfWeekSlot domain.Slot) string {
	if dayOfWeekSlot.Value == nil {
		return ""
	}

	dayName, ok := dayOfWeekSlot.Value.(string)
	if !ok {
		return ""
	}

	targetWeekday, ok := russianWeekdays[dayName]
	if !ok {
		return ""
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	now := h.TimeNowFunc().In(loc)
	currentWeekday := now.Weekday()

	// Calculate days until the target weekday
	daysUntil := int(targetWeekday) - int(currentWeekday)
	if daysUntil <= 0 {
		daysUntil += 7 // Move to next week if today or past
	}

	targetDate := now.AddDate(0, 0, daysUntil)
	return targetDate.Format("02.01.2006")
}

func (h *AliceHandler) resolveDate(timezone string, dateSlot domain.Slot) string {
	if dateSlot.Value == nil {
		return ""
	}

	valueMap, ok := dateSlot.Value.(map[string]interface{})
	if !ok {
		return ""
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	now := h.TimeNowFunc().In(loc)

	if dayIsRelative, ok := valueMap["day_is_relative"].(bool); ok && dayIsRelative {
		if dayOffset, ok := valueMap["day"].(float64); ok {
			targetDate := now.AddDate(0, 0, int(dayOffset))
			return targetDate.Format("02.01.2006")
		}
	}

	day, dayOk := valueMap["day"].(float64)
	month, monthOk := valueMap["month"].(float64)
	year, yearOk := valueMap["year"].(float64)

	if dayOk && monthOk {
		if !yearOk {
			year = float64(now.Year())
		}
		return fmt.Sprintf("%02d.%02d.%d", int(day), int(month), int(year))
	}

	return ""
}

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
