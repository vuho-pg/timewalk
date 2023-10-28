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
		YearField:      Field[int](From(0)),
		MonthField:     Field[time.Month](From(time.January)),
		DayField:       Field[int](From(1)),
		DayOfWeekField: Field[time.Weekday](From(time.Sunday)),
		HourField:      Field[int](From(0)),
		MinuteField:    Field[int](From(0)),
		SecondField:    Field[int](From(0)),
		Duration:       0 * time.Second,
		StartAt:        startAt,
		Location:       time.Local.String(),
		Loc:            time.Local,
	}
}

func (s *Schedule) SetLoc(loc *time.Location) *Schedule {
	s.Loc = loc
	s.Location = loc.String()
	return s
}

func (s *Schedule) SetLocString(loc string) *Schedule {
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

func (s *Schedule) Nearest(t time.Time) *time.Time {
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
	nY = s.YearField.Nearest(y)
	if nY == nil {
		return nil
	}
	over = over || *nY < t.Year()
month:
	if over {
		nM = s.MonthField.Nearest(time.December)
	} else {
		nM = s.MonthField.Nearest(m)
	}
	if nM == nil {
		y--
		goto year
	}
	over = over || *nM < t.Month()
day:
	if over {
		nD = s.DayField.Nearest(31)
	} else {
		nD = s.DayField.Nearest(d)
	}
	if nD == nil {
		m--
		goto month
	}
	over = over || *nD < t.Day()
hour:
	if over {
		nH = s.HourField.Nearest(23)
	} else {
		nH = s.HourField.Nearest(h)
	}
	if nH == nil {
		d--
		goto day
	}
	over = over || *nH < t.Hour()
minute:
	if over {
		nMin = s.MinuteField.Nearest(59)
	} else {
		nMin = s.MinuteField.Nearest(minute)
	}
	if nMin == nil {
		h--
		goto hour
	}
	over = over || *nMin < t.Minute()
	// second
	if over {
		nSec = s.SecondField.Nearest(59)
	} else {
		nSec = s.SecondField.Nearest(sec)

	}
	if nSec == nil {
		minute--
		goto minute
	}
	return Ptr(time.Date(*nY, *nM, *nD, *nH, *nMin, *nSec, 0, s.Loc))
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
