package timewalk

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_maxDay(t *testing.T) {
	// 31
	months := []time.Month{
		time.January, time.March, time.May, time.July, time.August, time.October, time.December,
	}
	for _, month := range months {
		assert.Equal(t, 31, maxDay(2020, month))
	}
	// 30
	months = []time.Month{
		time.April, time.June, time.September, time.November,
	}
	for _, month := range months {
		assert.Equal(t, 30, maxDay(2020, month))
	}
	// leap
	assert.Equal(t, 29, maxDay(2020, time.February))
	// non-leap
	assert.Equal(t, 28, maxDay(2019, time.February))
}

func Test_ordinalSuffix(t *testing.T) {
	suffix := "suffix"
	//st
	assert.Equal(t, "1st "+suffix, ordinalSuffix(1, suffix))
	assert.Equal(t, "11th "+suffix, ordinalSuffix(11, suffix))
	assert.Equal(t, "21st "+suffix, ordinalSuffix(21, suffix))
	//nd
	assert.Equal(t, "2nd "+suffix, ordinalSuffix(2, suffix))
	assert.Equal(t, "12th "+suffix, ordinalSuffix(12, suffix))
	assert.Equal(t, "22nd "+suffix, ordinalSuffix(22, suffix))
	//rd
	assert.Equal(t, "3rd "+suffix, ordinalSuffix(3, suffix))
	assert.Equal(t, "13th "+suffix, ordinalSuffix(13, suffix))
	assert.Equal(t, "23rd "+suffix, ordinalSuffix(23, suffix))
	// th
	assert.Equal(t, "4th "+suffix, ordinalSuffix(4, suffix))
	assert.Equal(t, "111th "+suffix, ordinalSuffix(111, suffix))
	assert.Equal(t, "101st "+suffix, ordinalSuffix(101, suffix))

	// month
	assert.Equal(t, "February", ordinalSuffix(time.February, suffix))
}

func Test_ptr(t *testing.T) {
	value := 10
	assert.Equal(t, &value, ptr(value))
}
