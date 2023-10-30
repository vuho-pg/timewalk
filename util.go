package timewalk

import (
	"fmt"
	"time"
)

func maxDay(year int, month time.Month) int {
	switch month {

	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if year%4 == 0 && year%100 != 0 || year%400 == 0 {
			return 29
		}
		return 28
	default:
		return 31
	}
}

func ordinalSuffix[T TimeUnit](x T, suffix string) string {
	switch any(x).(type) {
	case int:
		switch x % 100 {
		case 11, 12, 13:
			return fmt.Sprint(x, "th", " ", suffix)
		}
		switch x % 10 {
		case 1:
			return fmt.Sprint(x, "st", " ", suffix)
		case 2:
			return fmt.Sprint(x, "nd", " ", suffix)
		case 3:
			return fmt.Sprint(x, "rd", " ", suffix)
		default:
			return fmt.Sprint(x, "th", " ", suffix)
		}
	default:
		return fmt.Sprint(x)
	}

}

func ptr[T any](value T) *T {
	return &value
}

func max[T TimeUnit](values ...T) T {
	now := values[0]
	for _, v := range values[1:] {
		if v > now {
			now = v
		}
	}
	return now
}
