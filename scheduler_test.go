package timewalk

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNewSchedule(t *testing.T) {
	s := NewSchedule(time.Now()).
		SetLocation(time.Local.String()).
		SetDuration(time.Hour*24).
		Year(Unit[int]().From(2023).To(2025)).
		Month(Unit[time.Month]().From(1).To(4), Unit[time.Month]().At(10)).
		Day(Unit[int]().From(1).To(31)).
		DayOfWeek(Unit[time.Weekday]().From(time.Monday).To(time.Friday)).
		//Hour(Unit[TimeOfDay]().From(TimeOfDay{Hour: 0, Minute: 0, Second: 0}).To(TimeOfDay{Hour: 23, Minute: 59, Second: 59}))
		Hour(Unit[int]().At(10)).Minute(Unit[int]().At(0)).Second(Unit[int]().At(0))
	fmt.Println(s)
	j, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(j))
}
