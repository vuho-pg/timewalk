package timewalk

import "time"

type Time struct {
	Year      int
	Month     time.Month
	Week      int
	Day       int
	DayOfWeek time.Weekday
	Hour      int
	Minute    int
	Second    int
	Loc       *time.Location
}

func EmptyTime(loc *time.Location) *Time {
	return &Time{
		Year:      -1,
		Month:     -1,
		Week:      -1,
		Day:       -1,
		DayOfWeek: -1,
		Hour:      -1,
		Minute:    -1,
		Second:    -1,
		Loc:       loc,
	}
}

func T(t time.Time) *Time {
	return &Time{
		Year:      t.Year(),
		Month:     t.Month(),
		Week:      (t.Day()-1)/7 + 1,
		Day:       t.Day(),
		DayOfWeek: t.Weekday(),
		Hour:      t.Hour(),
		Minute:    t.Minute(),
		Second:    t.Second(),
		Loc:       t.Location(),
	}
}

func (t *Time) ToTime() time.Time {
	return time.Date(t.Year, t.Month, t.Day, t.Hour, t.Minute, t.Second, 0, t.Loc)
}
