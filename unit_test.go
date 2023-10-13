package timewalk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_At(t *testing.T) {
	u := NewUnit().Named("day").At(1)
	assert.Equal(t, "at day no 1", u.String())
}

func TestUnit_From(t *testing.T) {
	u := NewUnit().Named("day").From(1)
	assert.Equal(t, "from day no 1", u.String())
}

func TestUnit_To(t *testing.T) {
	u := NewUnit().Named("day").To(10)
	assert.Equal(t, "from day no 0 to day no 10", u.String())
}

func TestUnit_FromTo(t *testing.T) {
	u := NewUnit().Named("day").From(3).To(10)
	assert.Equal(t, "from day no 3 to day no 10", u.String())
}

func TestUnit_Step(t *testing.T) {
	u := NewUnit().Named("day").Step(2)
	assert.Equal(t, "every 2 day", u.String())
	fmt.Println(u)
}

func TestUnit_StepRange(t *testing.T) {
	u := NewUnit().Named("day").Step(2).From(1).To(10)
	assert.Equal(t, "every 2 day from day no 1 to day no 10", u.String())
}
