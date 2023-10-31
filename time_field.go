package timewalk

import "strings"

type TField[T TimeUnit] []*Unit[T]

func Field[T TimeUnit](unit ...*Unit[T]) TField[T] {
	return TField[T](unit)
}

func (f TField[T]) String(unitName string) string {
	b := strings.Builder{}
	b.WriteString(f[0].String(unitName))
	for _, v := range f[1:] {
		b.WriteString(" and ")
		b.WriteString(v.String(unitName))
	}
	return b.String()
}

// PreviousInPool returns the previous value in the pool
func (f TField[T]) PreviousInPool(data T, pool []T) T {
	if len(pool) == 0 {
		return -1
	}
	res := T(-1)
	for _, v := range pool {
		if v <= data && f.Match(v) {
			res = max(res, v)
		}
	}
	return res
}

func (f TField[T]) Match(data T) bool {
	for _, u := range f {
		if u.Match(data) {
			return true
		}
	}
	return false
}

func (f TField[T]) Previous(data T) T {
	var res = T(-1)
	for _, u := range f {
		now := u.Previous(data)
		if now != -1 {
			res = max(res, now)
		}
	}
	return res
}
