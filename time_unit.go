package timewalk

import (
	"fmt"
	"strings"
	"time"
)

type TimeUnit interface {
	int | time.Weekday | time.Month
}

type TUnit[T TimeUnit] struct {
	Type      UnitType `json:"type"`
	Value     *T       `json:"value,omitempty"`
	ValueFrom *T       `json:"value_from,omitempty"`
	ValueTo   *T       `json:"value_to,omitempty"`
	ValueStep *T       `json:"value_step,omitempty"`
}

func Unit[T TimeUnit]() *TUnit[T] {
	return &TUnit[T]{}
}

func (u *TUnit[T]) At(value T) *TUnit[T] {
	u.Type |= TValue
	u.Value = &value
	return u
}

func At[T TimeUnit](value T) *TUnit[T] {
	return Unit[T]().At(value)
}

func From[T TimeUnit](value T) *TUnit[T] {
	return Unit[T]().From(value)
}

func Every[T TimeUnit](step T) *TUnit[T] {
	return Unit[T]().Every(step)
}

func (u *TUnit[T]) Every(step T) *TUnit[T] {
	u.Type |= TStep
	u.ValueStep = &step
	return u
}

func (u *TUnit[T]) From(value T) *TUnit[T] {
	u.Type |= TRange
	u.ValueFrom = &value
	return u
}

func (u *TUnit[T]) To(value T) *TUnit[T] {
	u.Type |= TRange
	u.ValueTo = &value
	return u
}

func (u *TUnit[T]) String(unitName string) string {
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

func (u *TUnit[T]) Previous(data T) *T {
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
