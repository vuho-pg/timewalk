package timewalk

import "time"

type Schedulers []*Schedule

func (s Schedulers) InProgress(t time.Time) bool {
	for _, v := range s {
		v := v
		if !v.Enable {
			continue
		}
		if v.InProgress(t) {
			return true
		}
	}
	return false
}
