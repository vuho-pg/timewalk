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
	StartAt        *time.Time           `json:"start_at"`
	Location       string               `json:"location"`
	Loc            *time.Location       `json:"-"`
}

func NewSchedule() *Schedule {
	return &Schedule{
		YearField:      Field[int](From(0)),
		MonthField:     Field[time.Month](From(time.January)),
		DayField:       Field[int](From(1)),
		DayOfWeekField: Field[time.Weekday](From(time.Sunday)),
		HourField:      Field[int](From(0)),
		MinuteField:    Field[int](From(0)),
		SecondField:    Field[int](From(0)),
		Duration:       0 * time.Second,
		StartAt:        nil,
		Location:       time.Local.String(),
		Loc:            time.Local,
	}
}

func (s *Schedule) WithStartAt(t time.Time) *Schedule {
	s.StartAt = ptr(t)
	return s
}

func (s *Schedule) WithLoc(loc *time.Location) *Schedule {
	s.Loc = loc
	s.Location = loc.String()
	return s
}

func (s *Schedule) WithLocString(loc string) *Schedule {
	s.Location = loc
	tz, err := time.LoadLocation(loc)
	if err != nil {
		tz = time.Local
	}
	s.Loc = tz
	return s
}

func (s *Schedule) WithDuration(dur time.Duration) *Schedule {
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

func (s *Schedule) Previous(t time.Time) *time.Time {
	t = t.In(s.Loc)
	y := t.Year()
	m := t.Month()
	d := t.Day()
	h := t.Hour()
	minute := t.Minute()
	sec := t.Second()
	var (
		nY   *int
		nM   *time.Month
		nD   *int
		nH   *int
		nMin *int
		nSec *int
	)
	over := false
year:
	nY = s.YearField.Previous(y)
	if nY == nil {
		return nil
	}
	over = over || *nY < t.Year()
month:
	if over {
		nM = s.MonthField.Previous(time.December)
	} else {
		nM = s.MonthField.Previous(m)
	}
	if nM == nil {
		y--
		goto year
	}
	over = over || *nM < t.Month()
day:
	if over {
		nD = s.DayField.Previous(maxDay(*nY, *nM))
	} else {
		nD = s.DayField.Previous(d)
	}
	if nD == nil {
		m--
		goto month
	}
	over = over || *nD < t.Day()
hour:
	if over {
		nH = s.HourField.Previous(23)
	} else {
		nH = s.HourField.Previous(h)
	}
	if nH == nil {
		d--
		goto day
	}
	over = over || *nH < t.Hour()
minute:
	if over {
		nMin = s.MinuteField.Previous(59)
	} else {
		nMin = s.MinuteField.Previous(minute)
	}
	if nMin == nil {
		h--
		goto hour
	}
	over = over || *nMin < t.Minute()
	// second
	if over {
		nSec = s.SecondField.Previous(59)
	} else {
		nSec = s.SecondField.Previous(sec)

	}
	if nSec == nil {
		minute--
		goto minute
	}
	return ptr(time.Date(*nY, *nM, *nD, *nH, *nMin, *nSec, 0, s.Loc))
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
	if s.StartAt != nil {
		b.WriteString(", start from ")
		b.WriteString(s.StartAt.In(s.Loc).Format(time.RFC850))
	}
	if s.Duration != 0 {
		b.WriteString(" with ")
		b.WriteString(s.Duration.String())
		b.WriteString(" duration")
	}

	return b.String()
}
