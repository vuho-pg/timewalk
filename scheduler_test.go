package timewalk

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNewSchedule(t *testing.T) {
	s := NewSchedule(time.Now()).
		SetLocString(time.Local.String()).
		SetDuration(time.Hour*24).
		Year(Unit[int]().From(2023).To(2025)).
		Month(Unit[time.Month]().From(1).To(4), Unit[time.Month]().At(10)).
		Day(Unit[int]().From(1).To(31)).
		DayOfWeek(Unit[time.Weekday]().From(time.Monday).To(time.Friday)).
		Hour(Unit[int]().At(10)).Minute(Unit[int]().At(0)).Second(Unit[int]().At(0))
	fmt.Println(s)
	j, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(j))
}

func TestSchedule_Nearest(t *testing.T) {
	// daily 10:30 -> 11:00, 11:15-11:30
	s := NewSchedule(time.Unix(0, 0)).SetLocString(time.Local.String()).
		SetDuration(30 * time.Minute).
		Year(From(0)).
		Month(From(time.January)).
		Day(From(1)).
		DayOfWeek(From(time.Sunday)).
		Hour(At(10)).Minute(At(30))
	fmt.Println(s.Nearest(time.Date(2020, 11, 1, 10, 29, 0, 0, time.Local)))
	fmt.Println(s.Nearest(time.Date(2020, 11, 3, 0, 0, 0, 0, time.Local)))

}
