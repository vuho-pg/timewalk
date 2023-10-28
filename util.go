package timewalk

import (
	"cmp"
	"fmt"
)

func wrapUnit[T TimeUnit](x T, suffix string) string {
	anyX := any(x)
	// int
	intX, ok := anyX.(int)
	if ok {
		switch intX % 10 {
		case 1:
			return fmt.Sprint(x, "st", " ", suffix)
		case 2:
			return fmt.Sprint(x, "nd", " ", suffix)
		case 3:
			return fmt.Sprint(x, "rd", " ", suffix)
		default:
			return fmt.Sprint(x, "th", " ", suffix)
		}
	}
	// time.Weekday
	// time.Month
	return fmt.Sprint(x)
}

func Max[T cmp.Ordered](arr ...T) T {
	vMax := arr[0]
	for _, x := range arr {
		if x > vMax {
			vMax = x
		}
	}
	return vMax
}

func Min[T cmp.Ordered](arr ...T) T {
	vMin := arr[0]
	for _, x := range arr {
		if x < vMin {
			vMin = x
		}
	}
	return vMin
}

func Ptr[T any](value T) *T {
	return &value
}
