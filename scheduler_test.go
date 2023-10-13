package timewalk

import (
	"fmt"
	"testing"
)

func TestNewSchedule(t *testing.T) {
	s := NewSchedule().
		Year(NewField(NewUnit().From(2023).To(2025))).
		Month(NewField(NewUnit().From(1).To(4), NewUnit().At(10))).
		Day(NewField(NewUnit().From(1).To(31))).
		Hour(NewField(NewUnit().From(0).To(23)))
	fmt.Println(s)
}
