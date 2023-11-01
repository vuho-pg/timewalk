package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_String(t *testing.T) {
	unitName := "unit"
	// step
	s := Every(3)
	assert.Equal(t, "every 3 "+unitName, s.String(unitName))

	// range
	// with limit
	r := From(0).To(10)
	assert.Equal(t, "from 0th "+unitName+" through 10th "+unitName, r.String(unitName))
	// without limit
	r = From(0)
	assert.Equal(t, "from 0th "+unitName, r.String(unitName))
	// with step and limit
	r = From(0).To(10).Every(3)
	assert.Equal(t, "every 3 "+unitName+" from 0th "+unitName+" through 10th "+unitName, r.String(unitName))
	// with step and without limit
	r = From(0).Every(3)
	assert.Equal(t, "every 3 "+unitName+" from 0th "+unitName, r.String(unitName))

	// value
	v := At(10)
	assert.Equal(t, "at 10th "+unitName, v.String(unitName))

}

func TestUnit_Match(t *testing.T) {
	// value
	v := At(10)
	assert.True(t, v.Match(10))
	assert.False(t, v.Match(9))
	// range
	r := From(10).To(20)
	assert.True(t, r.Match(10))
	assert.True(t, r.Match(15))
	assert.True(t, r.Match(20))
	assert.False(t, r.Match(9))
	assert.False(t, r.Match(21))
	// range step
	r = r.Every(2)
	assert.True(t, r.Match(10))
	assert.False(t, r.Match(11))
	assert.True(t, r.Match(12))
	assert.False(t, r.Match(13))
	// step
	s := Every(2)
	assert.True(t, s.Match(2))
	assert.False(t, s.Match(3))
	// unknown
	u := Unit[int]{
		Type: TUnknown,
	}
	assert.False(t, u.Match(0))
}

func TestUnit_Previous(t *testing.T) {

	// value
	v := Unit[int]{
		Type:  TValue,
		Value: ptr(10),
	}
	// * o
	assert.Equal(t, 10, v.Previous(11))
	// o *
	assert.Equal(t, -1, v.Previous(9))
	// o|*
	assert.Equal(t, 10, v.Previous(10))

	// range
	r := Unit[int]{
		Type:      TRange,
		ValueFrom: ptr(10),
		ValueTo:   ptr(20),
	}
	// [] o
	assert.Equal(t, 20, r.Previous(21))
	// o []
	assert.Equal(t, -1, r.Previous(9))
	// [ o ]
	assert.Equal(t, 15, r.Previous(15))

	// range step
	rs := Unit[int]{
		Type:      TRange | TStep,
		ValueFrom: ptr(10),
		ValueTo:   ptr(20),
		ValueStep: ptr(3),
	}
	// o [/]
	assert.Equal(t, -1, rs.Previous(9))
	// [ o / ]
	assert.Equal(t, 10, rs.Previous(12))
	assert.Equal(t, 13, rs.Previous(15))
	// [/] o
	assert.Equal(t, 19, rs.Previous(21))

	// step
	s := Unit[int]{
		Type:      TStep,
		ValueStep: ptr(3),
	}
	// o /
	assert.Equal(t, -1, s.Previous(-1))
	// / o /
	assert.Equal(t, 3, s.Previous(4))
	assert.Equal(t, 6, s.Previous(7))
	assert.Equal(t, 300, s.Previous(300))

	// unknown
	u := Unit[int]{
		Type: TUnknown,
	}
	assert.Equal(t, -1, u.Previous(0))
}

func TestUnit_Next(t *testing.T) {

	// value
	v := Unit[int]{
		Type:  TValue,
		Value: ptr(10),
	}
	// * o
	assert.Equal(t, 10, v.Next(9))
	// o *
	assert.Equal(t, -1, v.Next(11))
	// o|*
	assert.Equal(t, 10, v.Next(10))

	// range
	r := Unit[int]{
		Type:      TRange,
		ValueFrom: ptr(10),
		ValueTo:   ptr(20),
	}
	// [] o
	assert.Equal(t, -1, r.Next(21))
	// o []
	assert.Equal(t, 10, r.Next(9))
	// [ o ]
	assert.Equal(t, 15, r.Next(15))

	// range step
	rs := Unit[int]{
		Type:      TRange | TStep,
		ValueFrom: ptr(10),
		ValueTo:   ptr(20),
		ValueStep: ptr(3),
	}
	// o [/]
	assert.Equal(t, 10, rs.Next(9))
	// [ o / ]
	assert.Equal(t, 13, rs.Next(11))
	assert.Equal(t, 13, rs.Next(13))
	// [/] o
	assert.Equal(t, -1, rs.Next(21))
	assert.Equal(t, -1, rs.Next(20))

	// step
	s := Unit[int]{
		Type:      TStep,
		ValueStep: ptr(3),
	}
	// o /
	assert.Equal(t, 3, s.Next(1))
	// / o /
	assert.Equal(t, -1, s.Next(-1))
	assert.Equal(t, 6, s.Next(4))
	assert.Equal(t, 6, s.Next(6))
	assert.Equal(t, 300, s.Next(300))
	// unknown
	u := Unit[int]{
		Type: TUnknown,
	}
	assert.Equal(t, -1, u.Next(0))
}
