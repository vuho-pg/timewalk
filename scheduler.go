package timewalk

import (
	"encoding/json"
	"strings"
	"sync"
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
	once           sync.Once
	YearField      TField[int]          `json:"year"`
	MonthField     TField[time.Month]   `json:"month"`       //1-12
	DayField       TField[int]          `json:"day"`         //1-31
	DayOfWeekField TField[time.Weekday] `json:"day_of_week"` //0-6
	HourField      TField[int]          `json:"hour"`        //0-23
	MinuteField    TField[int]          `json:"minute"`      //0-59
	SecondField    TField[int]          `json:"second"`      //0-59
	Duration       time.Duration        `json:"duration"`
	Start          int64                `json:"start"`
	End            int64                `json:"end"`
	Location       string               `json:"location"`
	StartTime      *time.Time           `json:"-"`
	EndTime        *time.Time           `json:"-"`
	Loc            *time.Location       `json:"-"`
}

func Scheduler() *Schedule {
	s := &Schedule{
		Location: time.Local.String(),
	}
	s.once.Do(s.correct)
	return s
}

func (s *Schedule) correct() {
	if s.Start != 0 {
		s.StartTime = ptr(time.Unix(s.Start, 0))
	}
	if s.End != 0 {
		s.EndTime = ptr(time.Unix(s.End, 0))
	}
	s.WithLocString(s.Location)
}

func ScheduleFromJSON(data string) (*Schedule, error) {
	var s *Schedule
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return nil, err
	}
	s.once.Do(s.correct)
	return s, nil
}

func (s *Schedule) StartAt(t *time.Time) *Schedule {
	s.StartTime = t
	if t != nil {
		s.Start = t.Unix()
	} else {
		s.Start = 0
	}
	return s
}

func (s *Schedule) EndAt(t *time.Time) *Schedule {
	s.EndTime = t
	if t != nil {
		s.End = t.Unix()
	} else {
		s.End = 0
	}
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
	s.once.Do(s.correct)
	t = t.In(s.Loc)
	y, m, d, h, minute, sec := t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()
	var (
		nY   *int
		nM   *time.Month
		nD   *int
		nH   *int
		nMin *int
		nSec *int
	)
	oY, oM, oD, oH, oMin := false, false, false, false, false
year:
	yearField := s.YearField
	if len(yearField) == 0 {
		yearField = AnyYear
	}
	nY = yearField.Previous(y)
	if nY == nil {
		return nil
	}
	oY = *nY < t.Year()

month:
	monthField := s.MonthField
	if len(monthField) == 0 {
		monthField = AnyMonth
	}
	if oY {
		nM = monthField.Previous(time.December)
	} else {
		nM = monthField.Previous(m)
	}
	if nM == nil {
		y--
		goto year
	}
	oM = oY || *nM < t.Month()
	maxDayOfMonth := maxDay(*nY, *nM)
	dayPool := make([]int, 0)
	validateWD := false
	if len(s.DayOfWeekField) > 0 {
		validateWD = true
		dayOfWeekField := s.DayOfWeekField
		if len(dayOfWeekField) == 0 {
			dayOfWeekField = AnyDayOfWeek
		}
		firstDayOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, s.Loc)
		wd := firstDayOfMonth.Weekday()
		for i := 1; i <= maxDayOfMonth; i++ {
			if dayOfWeekField.Match(wd) {
				dayPool = append(dayPool, i)
			}
			wd = (wd + 1) % 7
		}
	}
day:
	dayField := s.DayField
	if len(dayField) == 0 {
		dayField = AnyDay
	}
	if oM {
		if validateWD {
			nD = dayField.PreviousInPool(maxDayOfMonth, dayPool)
		} else {
			nD = dayField.Previous(maxDayOfMonth)
		}
	} else {
		if validateWD {
			nD = dayField.PreviousInPool(d, dayPool)
		} else {
			nD = dayField.Previous(d)
		}
	}
	if nD == nil {
		m--
		goto month
	}
	oD = oM || *nD < t.Day()

hour:
	hourField := s.HourField
	if len(hourField) == 0 {
		hourField = AnyHour
	}
	if oD {
		nH = hourField.Previous(23)
	} else {
		nH = hourField.Previous(h)
	}
	if nH == nil {
		d--
		goto day
	}
	oH = oD || *nH < t.Hour()
minute:
	minField := s.MinuteField
	if len(minField) == 0 {
		minField = AnyMinute
	}
	if oH {
		nMin = minField.Previous(59)
	} else {
		nMin = minField.Previous(minute)
	}
	if nMin == nil {
		h--
		goto hour
	}
	oMin = oH || *nMin < t.Minute()
	// second
	secField := s.SecondField
	if len(secField) == 0 {
		secField = AnySecond
	}
	if oMin {
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
	s.once.Do(s.correct)
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
	if s.EndTime != nil {
		b.WriteString(", end at ")
		b.WriteString(s.EndTime.In(s.Loc).Format(time.RFC850))
	}
	if s.Duration != 0 {
		b.WriteString(" with ")
		b.WriteString(s.Duration.String())
		b.WriteString(" duration")
	}

	return b.String()
}
