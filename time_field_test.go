package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTField_Nearest(t *testing.T) {
	fields := Field[int](
		Unit[int]().At(10),
		Unit[int]().From(15).To(30),
		Unit[int]().At(40),
	)
	assert.Nil(t, fields.Nearest(9))
	assert.Equal(t, Ptr(10), fields.Nearest(13))
	assert.Equal(t, Ptr(20), fields.Nearest(20))
	assert.Equal(t, Ptr(30), fields.Nearest(35))
	assert.Equal(t, Ptr(40), fields.Nearest(45))
}
