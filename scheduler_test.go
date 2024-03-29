package timewalk

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSchedule_String(t *testing.T) {
	now := time.Now()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	assert.NoError(t, err)
	s := Scheduler().StartAt(&now).EndAt(ptr(now.Add(time.Hour))).WithLoc(loc).WithDuration(time.Hour).
		Year(At(2023)).
		Month(At(time.December)).
		Week(From(1).To(5)).
		Day(At(22)).
		DayOfWeek(At(time.Tuesday)).
		Hour(At(3)).Minute(At(11)).Second(At(0))
	str := fmt.Sprintf("at 2023rd year, at December, from 1st week through 5th week, at 22nd day, at Tuesday, at 3rd hour, at 11th minute, at 0th second, start from %v, end at %v with 1h0m0s duration", now.In(loc).Format(time.RFC850), now.In(loc).Add(time.Hour).Format(time.RFC850))
	assert.Equal(t, str, s.String())
}

func TestSchedule_JSON(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	assert.NoError(t, err)
	s := Scheduler().StartAt(ptr(time.Now())).EndAt(ptr(time.Now().Add(time.Hour))).WithLoc(loc).WithDuration(time.Hour).
		Year(At(2023), From(2025).To(2030).Every(2)).
		Month(At(time.January), From(time.March).To(time.December).Every(2)).
		Day(At(1), From(2).To(31).Every(3)).
		DayOfWeek(At(time.Tuesday)).
		Hour(At(3)).Minute(At(11)).Second(At(0))

	sJSON, err := json.MarshalIndent(s, "", "\t")
	assert.NoError(t, err)
	fromJSON, err := ScheduleFromJSON(string(sJSON))
	assert.NoError(t, err)
	assert.Equal(t, s.String(), fromJSON.String())
	fmt.Println(string(sJSON))
}

func TestSchedule_StartAt(t *testing.T) {
	now := time.Now()
	s := Scheduler().StartAt(&now)
	assert.Equal(t, now.Unix(), s.Start)
	assert.NotNil(t, s.StartTime)
	assert.Equal(t, now, *s.StartTime)
	s.StartAt(nil)
	assert.Equal(t, int64(0), s.Start)
	assert.Nil(t, s.StartTime)
}

func TestSchedule_EndAt(t *testing.T) {
	now := time.Now()
	s := Scheduler().EndAt(&now)
	assert.Equal(t, now.Unix(), s.End)
	assert.NotNil(t, s.EndTime)
	assert.Equal(t, now, *s.EndTime)
	s.EndAt(nil)
	assert.Equal(t, int64(0), s.End)
	assert.Nil(t, s.EndTime)
}

func TestSchedule_WithLoc(t *testing.T) {
	loc := time.UTC
	s := Scheduler().WithLoc(loc)
	assert.Equal(t, loc, s.Loc)
	assert.Equal(t, loc.String(), s.Location)
}

func TestSchedule_WithLocString(t *testing.T) {
	loc := time.UTC.String()
	s := Scheduler().WithLocString(loc)
	assert.Equal(t, loc, s.Location)
	assert.Equal(t, loc, s.Loc.String())
}

func TestSchedule_WithLocStringDefault(t *testing.T) {
	loc := "invalid"
	s := Scheduler().WithLocString(loc)
	assert.Equal(t, time.Local.String(), s.Location)
	assert.Equal(t, time.Local.String(), s.Loc.String())
}

func TestSchedule_WithDuration(t *testing.T) {
	dur := time.Hour
	s := Scheduler().WithDuration(dur)
	assert.Equal(t, dur, s.Duration)
}

func TestSchedule_Next_Any(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	s := Scheduler()
	assert.Equal(t, &now, s.Next(now))
}

func TestSchedule_Previous_Any(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	s := Scheduler()
	assert.Equal(t, &now, s.Previous(now))
}

func TestSchedule_Previous_Nil(t *testing.T) {
	s := Scheduler().Year(From(2025))
	assert.Nil(t, s.Previous(time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)))
}

func TestSchedule_Next_Nil(t *testing.T) {
	s := Scheduler().Year(From(0).To(2025))
	assert.Nil(t, s.Next(time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)))
}

func TestSchedule_Next_NextYear(t *testing.T) {
	s := Scheduler().Year(From(2023).To(2025))
	assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local), *s.Next(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)))
}

func TestSchedule_Previous_PrevYear(t *testing.T) {
	s := Scheduler().Year(From(2023).To(2025)).Month(From(time.February).To(time.March))
	assert.Equal(t, time.Date(2023, 3, 31, 23, 59, 59, 0, time.Local), *s.Previous(time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)))
}

func TestSchedule_Previous_PrevMonth(t *testing.T) {
	s := Scheduler().Year(At(2023)).Month(At(time.January), At(time.February)).Day(From(2))
	assert.Equal(t, time.Date(2023, 1, 31, 23, 59, 59, 0, time.Local), *s.Previous(time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local)))
}

func TestSchedule_Previous_PrevDay(t *testing.T) {
	s := Scheduler().Year(At(2023)).Month(At(time.January)).Day(At(1), At(2)).Hour(From(2))
	assert.Equal(t, time.Date(2023, 1, 1, 23, 59, 59, 0, time.Local), *s.Previous(time.Date(2023, 1, 2, 1, 0, 0, 0, time.Local)))
}

func TestSchedule_Previous_PrevHour(t *testing.T) {
	s := Scheduler().Year(At(2023)).Month(At(time.January)).Day(At(1)).Hour(At(1), At(2)).Minute(From(2))
	assert.Equal(t, time.Date(2023, 1, 1, 1, 59, 59, 0, time.Local), *s.Previous(time.Date(2023, 1, 1, 2, 0, 0, 0, time.Local)))
}

func TestSchedule_Previous_PrevMinute(t *testing.T) {
	s := Scheduler().Year(At(2023)).Month(At(time.January)).Day(At(1)).Hour(At(1)).Minute(At(1), At(2)).Second(From(2))
	assert.Equal(t, time.Date(2023, 1, 1, 1, 1, 59, 0, time.Local), *s.Previous(time.Date(2023, 1, 1, 1, 2, 0, 0, time.Local)))
}

func TestSchedule(t *testing.T) {
	// every Tue and Thu, 10:30
	//s := Scheduler().DayOfWeek(At(time.Tuesday), At(time.Thursday)).Hour(At(10)).Minute(At(30)).Second(At(0))
	//now := time.Date(2023, 10, 31, 11, 0, 0, 0, time.Local)
	//assert.Equal(t, ptr(time.Date(2023, 10, 31, 10, 30, 0, 0, time.Local)), s.Previous(now))
	//now = time.Date(2023, 10, 31, 0, 0, 0, 0, time.Local)
	//assert.Equal(t, ptr(time.Date(2023, 10, 26, 10, 30, 0, 0, time.Local)), s.Previous(now))
	//// every tue and 2 week
	//s = Scheduler().Week(From(1).Every(2)).DayOfWeek(At(time.Tuesday)).Hour(At(12)).Minute(At(30)).Second(At(0))
	//assert.Equal(t, ptr(time.Date(2023, 10, 31, 12, 30, 0, 0, time.Local)), s.Previous(time.Date(2023, 11, 1, 0, 0, 0, 0, time.Local)))
	//assert.Equal(t, ptr(time.Date(2023, 10, 17, 12, 30, 0, 0, time.Local)), s.Previous(time.Date(2023, 10, 31, 0, 0, 0, 0, time.Local)))
	//// every 2 weak
	s := Scheduler().Week(From(1).Every(2)).Hour(At(12)).Minute(At(30)).Second(At(0))
	assert.Equal(t, ptr(time.Date(2023, 10, 7, 12, 30, 0, 0, time.Local)), s.Previous(time.Date(2023, 10, 13, 0, 0, 0, 0, time.Local)))
	assert.Equal(t, ptr(time.Date(2023, 10, 6, 12, 30, 0, 0, time.Local)), s.Previous(time.Date(2023, 10, 7, 0, 0, 0, 0, time.Local)))
	assert.Equal(t, ptr(time.Date(2023, 10, 31, 12, 30, 0, 0, time.Local)), s.Previous(time.Date(2023, 11, 1, 0, 0, 0, 0, time.Local)))
	assert.Equal(t, ptr(time.Date(2023, 11, 1, 12, 30, 0, 0, time.Local)), s.Previous(time.Date(2023, 11, 2, 0, 0, 0, 0, time.Local)))
}
