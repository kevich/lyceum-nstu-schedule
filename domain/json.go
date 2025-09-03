package domain

type LessonJSON struct {
	Subject []string `json:"s"`
	Teacher []string `json:"t"`
	Group   []string `json:"g"`
	Room    []string `json:"r"`
}

type PeriodsJSON struct {
	StartDate string `json:"b"`
	EndDate   string `json:"e"`
}

type SchoolWrappedJSON[T any] struct {
	Schools map[string]T `json:"-"`
}

type ReplacementLessonJSON struct {
	Subject interface{} `json:"s"`
	Teacher []string    `json:"t"`
	Group   []string    `json:"g"`
	Room    []string    `json:"r"`
}

type ClassScheduleJSON map[string]LessonJSON
type ReplacementClassScheduleJSON map[string]ReplacementLessonJSON
type ReplacementClassDayJSON map[string]ReplacementClassScheduleJSON

type ScheduleDataJSON struct {
	LESSONSINDAY   int                                     `json:"LESSONSINDAY"`
	DAY_NAMES      []string                                `json:"DAY_NAMES"`
	TEACHERS       map[string]string                       `json:"TEACHERS"`
	SUBJECTS       map[string]string                       `json:"SUBJECTS"`
	CLASSES        map[string]string                       `json:"CLASSES"`
	ROOMS          map[string]string                       `json:"ROOMS"`
	CLASSGROUPS    map[string]string                       `json:"CLASSGROUPS"`
	LESSON_TIMES   map[string][2]string                    `json:"LESSON_TIMES"`
	CLASS_SCHEDULE map[string]map[string]ClassScheduleJSON `json:"CLASS_SCHEDULE"`
	PERIODS        map[string]PeriodsJSON                  `json:"PERIODS"`
	CLASS_EXCHANGE map[string]ReplacementClassDayJSON      `json:"CLASS_EXCHANGE"`
}
