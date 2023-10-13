package timewalk

import "strings"

type Schedule struct {
	YearField      *Field `json:"year"`
	MonthField     *Field `json:"month"`
	DayField       *Field `json:"day"`
	DayOfWeekField *Field `json:"day_of_week"`
	HourField      *Field `json:"hour"`
}

func NewSchedule() *Schedule {
	return &Schedule{}
}
func (s *Schedule) Year(field *Field) *Schedule {
	s.YearField = field
	if field.Name == "" {
		field.Named("year")
	}
	return s
}

func (s *Schedule) Month(field *Field) *Schedule {
	s.MonthField = field
	if field.Name == "" {
		field.Named("month")
	}
	return s
}

func (s *Schedule) Day(field *Field) *Schedule {
	s.DayField = field
	if field.Name == "" {
		field.Named("day")
	}
	return s
}

func (s *Schedule) DayOfWeek(field *Field) *Schedule {
	s.DayOfWeekField = field
	if field.Name == "" {
		field.Named("day of week")
	}
	return s
}

func (s *Schedule) Hour(field *Field) *Schedule {
	s.HourField = field
	if field.Name == "" {
		field.Named("hour")
	}
	return s
}

func (s *Schedule) String() string {
	b := strings.Builder{}
	pre := false
	if s.YearField != nil {
		pre = true
		b.WriteString(s.YearField.String())
	}
	if s.MonthField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.MonthField.String())
	}
	if s.DayField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.DayField.String())
	}
	if s.DayOfWeekField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.DayOfWeekField.String())
	}
	if s.HourField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.HourField.String())
	}

	return b.String()
}
