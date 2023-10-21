package timewalk

import (
	"strings"
	"time"
)

type Schedule struct {
	YearField      TField[int]          `json:"year"`
	MonthField     TField[time.Month]   `json:"month"`       //1-12
	DayField       TField[int]          `json:"day"`         //1-31
	DayOfWeekField TField[time.Weekday] `json:"day_of_week"` //0-6
	HourField      TField[int]          `json:"hour"`        //0-23
	MinuteField    TField[int]          `json:"minute"`      //0-59
	SecondField    TField[int]          `json:"second"`      //0-59
	Duration       time.Duration        `json:"duration"`
	StartAt        time.Time            `json:"start_at"`
	Location       string               `json:"location"`
	Loc            *time.Location       `json:"-"`
}

func NewSchedule(startAt time.Time) *Schedule {
	return &Schedule{
		YearField:      Field[int](Unit[int]().At(0)),
		MonthField:     Field[time.Month](Unit[time.Month]().At(0)),
		DayField:       Field[int](Unit[int]().At(0)),
		DayOfWeekField: Field[time.Weekday](Unit[time.Weekday]().At(0)),
		HourField:      Field[int](Unit[int]().At(0)),
		MinuteField:    Field[int](Unit[int]().At(0)),
		SecondField:    Field[int](Unit[int]().At(0)),
		Duration:       0,
		StartAt:        startAt,
		Location:       time.Local.String(),
		Loc:            time.Local,
	}
}

func (s *Schedule) SetLocation(loc string) *Schedule {
	s.Location = loc
	tz, err := time.LoadLocation(loc)
	if err != nil {
		tz = time.Local
	}
	s.Loc = tz
	return s
}

func (s *Schedule) SetDuration(dur time.Duration) *Schedule {
	s.Duration = dur
	return s
}

func (s *Schedule) Year(units ...*TUnit[int]) *Schedule {
	s.YearField = units
	return s
}

func (s *Schedule) Month(units ...*TUnit[time.Month]) *Schedule {
	s.MonthField = units
	return s
}

func (s *Schedule) Day(units ...*TUnit[int]) *Schedule {
	s.DayField = units
	return s
}

func (s *Schedule) DayOfWeek(units ...*TUnit[time.Weekday]) *Schedule {
	s.DayOfWeekField = units
	return s
}

func (s *Schedule) Hour(field ...*TUnit[int]) *Schedule {
	s.HourField = field
	return s
}

func (s *Schedule) Minute(field ...*TUnit[int]) *Schedule {
	s.MinuteField = field
	return s
}

func (s *Schedule) Second(field ...*TUnit[int]) *Schedule {
	s.SecondField = field
	return s
}

func (s *Schedule) NearestBefore(t time.Time) time.Time {
	t = t.In(s.Loc)
	y := t.Year()
	m := t.Month()
	d := t.Day()
	h := t.Hour()
	min := t.Minute()
	sec := t.Second()
	rem := 0
	return t
}

func (s *Schedule) String() string {
	b := strings.Builder{}
	pre := false
	if s.YearField != nil {
		pre = true
		b.WriteString(s.YearField.String("year"))
	}
	if s.MonthField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.MonthField.String("month"))
	}
	if s.DayField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.DayField.String("day"))
	}
	if s.DayOfWeekField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.DayOfWeekField.String(""))
	}
	if s.HourField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.HourField.String("hour"))
	}
	if s.MinuteField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.MinuteField.String("minute"))
	}
	if s.SecondField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.SecondField.String("second"))
	}
	b.WriteString(" start at ")
	b.WriteString(s.StartAt.In(s.Loc).Format(time.RFC850))
	b.WriteString(" with ")
	b.WriteString(s.Duration.String())
	b.WriteString(" duration")

	return b.String()
}
