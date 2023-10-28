package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTUnit_String(t *testing.T) {
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

func TestTUnit_NearestFrom(t *testing.T) {

	// value
	v := TUnit[int]{
		Type:  TValue,
		Value: ptr(10),
	}
	// * o
	assert.Equal(t, ptr(10), v.NearestBefore(11))
	// o *
	assert.Nil(t, v.NearestBefore(9))
	// o|*
	assert.Equal(t, ptr(10), v.NearestBefore(10))

	// range
	r := TUnit[int]{
		Type:      TRange,
		ValueFrom: ptr(10),
		ValueTo:   ptr(20),
	}
	// [] o
	assert.Equal(t, ptr(20), r.NearestBefore(21))
	// o []
	assert.Nil(t, r.NearestBefore(9))
	// [ o ]
	assert.Equal(t, ptr(15), r.NearestBefore(15))

	// range step
	rs := TUnit[int]{
		Type:      TRange | TStep,
		ValueFrom: ptr(10),
		ValueTo:   ptr(20),
		ValueStep: ptr(3),
	}
	// o [/]
	assert.Nil(t, rs.NearestBefore(9))
	// [ o / ]
	assert.Equal(t, ptr(10), rs.NearestBefore(12))
	assert.Equal(t, ptr(13), rs.NearestBefore(15))
	// [/] o
	assert.Equal(t, ptr(19), rs.NearestBefore(21))

	// step
	s := TUnit[int]{
		Type:      TStep,
		ValueStep: ptr(3),
	}
	// o /
	assert.Nil(t, s.NearestBefore(-1))
	// / o /
	assert.Equal(t, ptr(3), s.NearestBefore(4))
	assert.Equal(t, ptr(6), s.NearestBefore(7))
	assert.Equal(t, ptr(300), s.NearestBefore(300))

}