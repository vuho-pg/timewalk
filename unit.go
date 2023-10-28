package timewalk

type UnitType int

const (
	TValue UnitType = 1 << iota
	TRange
	TStep
)

func (u UnitType) Is(t UnitType) bool {
	return u&t != 0
}
