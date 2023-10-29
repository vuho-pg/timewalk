package timewalk

import "strings"

type TField[T TimeUnit] []*TUnit[T]

func Field[T TimeUnit](unit ...*TUnit[T]) TField[T] {
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

func (f TField[T]) Previous(data T) *T {
	var res *T
	for _, u := range f {
		now := u.Previous(data)
		if now != nil {
			if res == nil {
				res = now
				continue
			}
			res = ptr(max(*res, *now))
		}
	}
	return res
}
