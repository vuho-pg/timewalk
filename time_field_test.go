package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTField_NearestBefore(t *testing.T) {
	fields := Field[int](
		Unit[int]().At(10),
		Unit[int]().From(15).To(30),
		Unit[int]().At(40),
	)
	assert.Nil(t, fields.NearestBefore(9))
	assert.Equal(t, ptr(10), fields.NearestBefore(13))
	assert.Equal(t, ptr(20), fields.NearestBefore(20))
	assert.Equal(t, ptr(30), fields.NearestBefore(35))
	assert.Equal(t, ptr(40), fields.NearestBefore(45))
}
