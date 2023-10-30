package timewalk

import (
	"strings"
	"time"
)

var (
	AnyYear      = Field[int](From(0))
	AnyMonth     = Field[time.Month](From(time.January))
	AnyDay       = Field[int](From(1))
	AnyDayOfWeek = Field[time.Weekday](From(time.Sunday))
	AnyHour      = Field[int](From(0))
	AnyMinute    = Field[int](From(0))
	AnySecond    = Field[int](From(0))
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
	Start          int64                `json:"start"`
	Location       string               `json:"location"`
	StartTime      *time.Time           `json:"-"`
	Loc            *time.Location       `json:"-"`
}

func Scheduler() *Schedule {
	return &Schedule{
		Location: time.Local.String(),
		Loc:      time.Local,
	}
}

func (s *Schedule) StartAt(t time.Time) *Schedule {
	s.StartTime = &t
	s.Start = t.Unix()
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
	s.Location = tz.String()
	return s
}

func (s *Schedule) WithDuration(dur time.Duration) *Schedule {
	s.Duration = dur
	return s
}

func (s *Schedule) Year(units ...*Unit[int]) *Schedule {
	s.YearField = units
	return s
}

func (s *Schedule) Month(units ...*Unit[time.Month]) *Schedule {
	s.MonthField = units
	return s
}

func (s *Schedule) Day(units ...*Unit[int]) *Schedule {
	s.DayField = units
	return s
}

func (s *Schedule) DayOfWeek(units ...*Unit[time.Weekday]) *Schedule {
	s.DayOfWeekField = units
	return s
}

func (s *Schedule) Hour(field ...*Unit[int]) *Schedule {
	s.HourField = field
	return s
}

func (s *Schedule) Minute(field ...*Unit[int]) *Schedule {
	s.MinuteField = field
	return s
}

func (s *Schedule) Second(field ...*Unit[int]) *Schedule {
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
	yearField := s.YearField
	if len(yearField) == 0 {
		yearField = AnyYear
	}
	nY = yearField.Previous(y)
	if nY == nil {
		return nil
	}
	over = over || *nY < t.Year()
month:
	monthField := s.MonthField
	if len(monthField) == 0 {
		monthField = AnyMonth
	}
	if over {
		nM = monthField.Previous(time.December)
	} else {
		nM = monthField.Previous(m)
	}
	if nM == nil {
		y--
		goto year
	}
	over = over || *nM < t.Month()
day:
	dayField := s.DayField
	if len(dayField) == 0 {
		dayField = AnyDay
	}
	if over {
		nD = dayField.Previous(maxDay(*nY, *nM))
	} else {
		nD = dayField.Previous(d)
	}
	if nD == nil {
		m--
		goto month
	}
	over = over || *nD < t.Day()
hour:
	hourField := s.HourField
	if len(hourField) == 0 {
		hourField = AnyHour
	}
	if over {
		nH = hourField.Previous(23)
	} else {
		nH = hourField.Previous(h)
	}
	if nH == nil {
		d--
		goto day
	}
	over = over || *nH < t.Hour()
minute:
	minField := s.MinuteField
	if len(minField) == 0 {
		minField = AnyMinute
	}
	if over {
		nMin = minField.Previous(59)
	} else {
		nMin = minField.Previous(minute)
	}
	if nMin == nil {
		h--
		goto hour
	}
	over = over || *nMin < t.Minute()
	// second
	secField := s.SecondField
	if len(secField) == 0 {
		secField = AnySecond
	}
	if over {
		nSec = secField.Previous(59)
	} else {
		nSec = secField.Previous(sec)

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
	if s.StartTime != nil {
		b.WriteString(", start from ")
		b.WriteString(s.StartTime.In(s.Loc).Format(time.RFC850))
	}
	if s.Duration != 0 {
		b.WriteString(" with ")
		b.WriteString(s.Duration.String())
		b.WriteString(" duration")
	}

	return b.String()
}
