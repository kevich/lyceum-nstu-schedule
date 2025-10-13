package domain

type Lesson struct {
	Number      int
	TimeRange   [2]string
	Name        string
	Teacher     string
	RoomNumber  string
	IsCancelled bool
	IsReplaced  bool
}
type DaySchedule map[int][]Lesson
type ClassSchedule map[string]DaySchedule

type Schedule map[string]ClassSchedule
