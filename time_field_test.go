package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTField_PreviousInPool(t *testing.T) {
	fields := Field[int](
		At(10),
		From(15).To(30),
		At(40),
	)
	assert.Equal(t, -1, fields.PreviousInPool(9, []int{}))
	assert.Equal(t, 10, fields.PreviousInPool(13, []int{10}))
	assert.Equal(t, 16, fields.PreviousInPool(20, []int{15, 16}))
	assert.Equal(t, 30, fields.PreviousInPool(35, []int{30, 35, 40}))
	assert.Equal(t, 40, fields.PreviousInPool(45, []int{40}))
}

func TestTField_NextInPool(t *testing.T) {
	fields := Field[int](
		At(10),
		From(15).To(30),
		At(40),
	)
	assert.Equal(t, -1, fields.NextInPool(9, []int{}))
	assert.Equal(t, 10, fields.NextInPool(3, []int{10}))
	assert.Equal(t, 15, fields.NextInPool(10, []int{15, 16}))
	assert.Equal(t, 40, fields.NextInPool(35, []int{30, 35, 40}))
	assert.Equal(t, -1, fields.NextInPool(45, []int{40}))
}

func TestTField_Previous(t *testing.T) {
	fields := Field[int](
		At(10),
		From(15).To(30),
		At(40),
	)
	assert.Equal(t, -1, fields.Previous(9))
	assert.Equal(t, 10, fields.Previous(13))
	assert.Equal(t, 20, fields.Previous(20))
	assert.Equal(t, 30, fields.Previous(35))
	assert.Equal(t, 40, fields.Previous(45))
}

func TestTField_Next(t *testing.T) {
	fields := Field[int](
		At(10),
		From(15).To(30),
		At(40),
	)
	assert.Equal(t, 10, fields.Next(9))
	assert.Equal(t, 15, fields.Next(13))
	assert.Equal(t, 20, fields.Next(20))
	assert.Equal(t, 40, fields.Next(35))
	assert.Equal(t, -1, fields.Next(45))

}

func TestField_String(t *testing.T) {
	fields := Field[int](
		At(10),
		From(15).To(30),
		At(40),
	)

	assert.Equal(t, "at 10th unit and from 15th unit through 30th unit and at 40th unit", fields.String("unit"))
}
