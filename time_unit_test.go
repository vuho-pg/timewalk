package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeUnit_Nearest(t *testing.T) {

	// value
	v := TUnit[int]{
		Type:  TValue,
		Value: Ptr(10),
	}
	// * o
	assert.Equal(t, Ptr(10), v.Nearest(11))
	// o *
	assert.Nil(t, v.Nearest(9))
	// o|*
	assert.Equal(t, Ptr(10), v.Nearest(10))

	// range
	r := TUnit[int]{
		Type:      TRange,
		ValueFrom: Ptr(10),
		ValueTo:   Ptr(20),
	}
	// [] o
	assert.Equal(t, Ptr(20), r.Nearest(21))
	// o []
	assert.Nil(t, r.Nearest(9))
	// [ o ]
	assert.Equal(t, Ptr(15), r.Nearest(15))

	// range step
	rs := TUnit[int]{
		Type:      TRange | TStep,
		ValueFrom: Ptr(10),
		ValueTo:   Ptr(20),
		ValueStep: Ptr(3),
	}
	// o [/]
	assert.Nil(t, rs.Nearest(9))
	// [ o / ]
	assert.Equal(t, Ptr(10), rs.Nearest(12))
	assert.Equal(t, Ptr(13), rs.Nearest(15))
	// [/] o
	assert.Equal(t, Ptr(19), rs.Nearest(21))

	// step
	s := TUnit[int]{
		Type:      TStep,
		ValueStep: Ptr(3),
	}
	// o /
	assert.Nil(t, s.Nearest(-1))
	// / o /
	assert.Equal(t, Ptr(3), s.Nearest(4))
	assert.Equal(t, Ptr(6), s.Nearest(7))
	assert.Equal(t, Ptr(300), s.Nearest(300))

}
