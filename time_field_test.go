package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTField_Previous(t *testing.T) {
	fields := Field[int](
		At(10),
		From(15).To(30),
		At(40),
	)
	assert.Nil(t, fields.Previous(9))
	assert.Equal(t, ptr(10), fields.Previous(13))
	assert.Equal(t, ptr(20), fields.Previous(20))
	assert.Equal(t, ptr(30), fields.Previous(35))
	assert.Equal(t, ptr(40), fields.Previous(45))
}
