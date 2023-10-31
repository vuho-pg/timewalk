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
	AnyWeek      = Field[int](From(0))
	AnyDay       = Field[int](From(1))
	AnyDayOfWeek = Field[time.Weekday](From(time.Sunday))
	AnyHour      = Field[int](From(0))
	AnyMinute    = Field[int](From(0))
	AnySecond    = Field[int](From(0))
)

type Schedule struct {
	once           sync.Once
	YearField      TField[int]          `json:"year"`
	MonthField     TField[time.Month]   `json:"month"` //1-12
	WeekField      TField[int]          `json:"week"`
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

func (s *Schedule) Week(units ...*Unit[int]) *Schedule {
	s.WeekField = units
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
	now := T(t)
	upper := T(t)
	res := EmptyTime(s.Loc)

	oY, oM, oW, oD, oH, oMin := false, false, false, false, false, false
	upper.Year = now.Year
year:
	yearField := s.YearField
	if len(yearField) == 0 {
		yearField = AnyYear
	}
	res.Year = yearField.Previous(upper.Year)
	if res.Year == -1 {
		return nil
	}
	oY = res.Year < t.Year()
	upper.Month = now.Month
	if oY {
		upper.Month = time.December
	}

month:
	monthField := s.MonthField
	if len(monthField) == 0 {
		monthField = AnyMonth
	}
	res.Month = monthField.Previous(upper.Month)
	if res.Month == -1 {
		upper.Year--
		goto year
	}
	oM = oY || res.Month < t.Month()
	upper.Week = now.Week
	if oM {
		upper.Week = 5
	}
	maxDayOfMonth := maxDay(res.Year, res.Month)
	dayPool := make([]int, 0)
	poolDayValidate := false
	if len(s.DayOfWeekField) > 0 {
		poolDayValidate = true
		dayOfWeekField := s.DayOfWeekField
		if len(dayOfWeekField) == 0 {
			dayOfWeekField = AnyDayOfWeek
		}
		firstDayOfMonth := time.Date(res.Year, res.Month, 1, 0, 0, 0, 0, s.Loc)
		wd := firstDayOfMonth.Weekday()
		for i := 1; i <= maxDayOfMonth; i++ {
			if dayOfWeekField.Match(wd) {
				dayPool = append(dayPool, i)
			}
			wd = (wd + 1) % 7
		}
	}
	upper.Day = now.Day
	if oM {
		upper.Day = maxDayOfMonth
	}
week:
	weekField := s.WeekField
	wDayPoolValidate := false
	wDayPool := make([]int, 0)
	if len(weekField) != 0 {
		wDayPoolValidate = true
		res.Week = weekField.Previous(upper.Week)
		if res.Week == -1 {
			upper.Month--
			goto month
		}
		// handle day pool

		// 1 2 3 4 5 6 7
		// 8 9 10 11 12 13 14
		// 15 16 17 18 19 20 21
		// 22 23 24 25 26 27 28
		// 29 30 31
		start := (res.Week-1)*7 + 1
		end := min(res.Week*7, maxDayOfMonth)
		for i := start; i <= end; i++ {
			wDayPool = append(wDayPool, i)
		}
		if poolDayValidate {
			wDayPool = intersect(dayPool, wDayPool)
		}

	}
	oW = oM || (res.Week != -1 && res.Week < upper.Week)

day:
	dayField := s.DayField
	if len(dayField) == 0 {
		dayField = AnyDay
	}
	if wDayPoolValidate {
		res.Day = dayField.PreviousInPool(upper.Day, wDayPool)
	} else if poolDayValidate {
		res.Day = dayField.PreviousInPool(upper.Day, dayPool)
	} else {
		res.Day = dayField.Previous(upper.Day)
	}
	if res.Day == -1 {
		if len(weekField) == 0 {
			upper.Month--
			goto month
		}
		upper.Week--
		goto week
	}
	oD = oW || res.Day < t.Day()
	upper.Hour = now.Hour
	if oD {
		upper.Hour = 23
	}

hour:
	hourField := s.HourField
	if len(hourField) == 0 {
		hourField = AnyHour
	}
	res.Hour = hourField.Previous(upper.Hour)
	if res.Hour == -1 {
		upper.Day--
		goto day
	}
	oH = oD || res.Hour < t.Hour()
	upper.Minute = now.Minute
	if oH {
		upper.Minute = 59
	}
minute:
	minField := s.MinuteField
	if len(minField) == 0 {
		minField = AnyMinute
	}
	res.Minute = minField.Previous(upper.Minute)
	if res.Minute == -1 {
		upper.Hour--
		goto hour
	}
	oMin = oH || res.Minute < t.Minute()
	upper.Second = now.Second
	if oMin {
		upper.Second = 59
	}
	// second
	secField := s.SecondField
	if len(secField) == 0 {
		secField = AnySecond
	}
	res.Second = secField.Previous(upper.Second)
	if res.Second == -1 {
		upper.Minute--
		goto minute
	}
	return ptr(res.ToTime())
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
	if s.WeekField != nil {
		if pre {
			b.WriteString(", ")
		}
		pre = true
		b.WriteString(s.WeekField.String("week"))
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
