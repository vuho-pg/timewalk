package timewalk

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSchedule(t *testing.T) {
	s := NewSchedule().
		WithLocString(time.Local.String()).
		WithDuration(time.Hour*24).
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
	var s *Schedule
	//start := time.Unix(0, 0)
	// year
	// every 2 year
	s = NewSchedule().Year(Every(2))
	fmt.Println(s)
	assert.Equal(t, ptr(time.Date(2020, 12, 31, 23, 59, 59, 0, time.Local)), s.Nearest(time.Date(2021, 12, 14, 1, 3, 0, 0, time.Local)))
	assert.Equal(t, ptr(time.Date(2020, 4, 1, 9, 0, 0, 0, time.Local)), s.Nearest(time.Date(2020, 4, 1, 9, 0, 0, 0, time.Local)))

	// from 2020 to 2025
	s = NewSchedule().Year(From(2020).To(2025))
	fmt.Println(s)
	// from 2020 to 2025 every 2 year
	// at 2020

	// month

	// day

	// hour

	// min

	// sec

}
