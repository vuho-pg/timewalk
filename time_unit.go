package timewalk

import (
	"fmt"
	"strings"
	"time"
)

type UnitType int

const (
	TValue UnitType = 1 << iota
	TRange
	TStep
)

func (u UnitType) Is(t UnitType) bool {
	return u&t != 0
}

type TimeUnit interface {
	int | time.Weekday | time.Month
}

type Unit[T TimeUnit] struct {
	Type      UnitType `json:"type"`
	Value     *T       `json:"value,omitempty"`
	ValueFrom *T       `json:"value_from,omitempty"`
	ValueTo   *T       `json:"value_to,omitempty"`
	ValueStep *T       `json:"value_step,omitempty"`
}

func (u *Unit[T]) At(value T) *Unit[T] {
	u.Type |= TValue
	u.Value = &value
	return u
}

func At[T TimeUnit](value T) *Unit[T] {
	return ptr(Unit[T]{}).At(value)
}

func From[T TimeUnit](value T) *Unit[T] {
	return ptr(Unit[T]{}).From(value)
}

func Every[T TimeUnit](step T) *Unit[T] {
	return ptr(Unit[T]{}).Every(step)
}

func (u *Unit[T]) Every(step T) *Unit[T] {
	u.Type |= TStep
	u.ValueStep = &step
	return u
}

func (u *Unit[T]) From(value T) *Unit[T] {
	u.Type |= TRange
	u.ValueFrom = &value
	return u
}

func (u *Unit[T]) To(value T) *Unit[T] {
	u.Type |= TRange
	u.ValueTo = &value
	return u
}

func (u *Unit[T]) String(unitName string) string {
	b := strings.Builder{}
	if u.Type.Is(TStep) && u.ValueStep != nil {
		b.WriteString("every ")
		b.WriteString(fmt.Sprint(*u.ValueStep))
		b.WriteString(" ")
		b.WriteString(unitName)
	}
	if u.Type.Is(TRange) {
		if u.Type.Is(TStep) {
			b.WriteString(" ")
		}
		b.WriteString("from ")
		b.WriteString(ordinalSuffix(*u.ValueFrom, unitName))
		if u.ValueTo != nil {
			b.WriteString(" through ")
			b.WriteString(ordinalSuffix(*u.ValueTo, unitName))
		}
	}
	if u.Type.Is(TValue) && u.Value != nil {
		b.WriteString("at ")
		b.WriteString(ordinalSuffix(*u.Value, unitName))
	}
	return b.String()
}

func (u *Unit[T]) Previous(data T) *T {
	// [] range
	// [/] range step
	// * value
	// / step
	// o data

	if u.Type.Is(TValue) {
		// * o
		if *u.Value <= data {
			return u.Value
		}
	}
	if u.Type.Is(TRange) {

		// o [ ]
		if data < *u.ValueFrom {
			return nil
		}

		// [ o ]
		if u.ValueTo == nil || *u.ValueTo > data {
			if u.Type.Is(TStep) {
				dist := data - *u.ValueFrom
				stepCnt := dist / (*u.ValueStep)
				return ptr(*u.ValueFrom + stepCnt*(*u.ValueStep))
			}

			return &data
		}
		// [] o
		if u.ValueTo != nil && *u.ValueTo <= data {
			if u.Type.Is(TStep) {
				mod := (*u.ValueTo - *u.ValueFrom) % (*u.ValueStep)
				return ptr(*u.ValueTo - mod)
			}
			return u.ValueTo
		}
	}

	if u.Type.Is(TStep) {
		if data < 0 {
			return nil
		}
		stepCnt := data / (*u.ValueStep)
		return ptr(stepCnt * (*u.ValueStep))
	}

	return nil
}
