package timewalk

import (
	"fmt"
	"strings"
	"time"
)

type TimeUnit interface {
	int | TimeOfDay | time.Weekday | time.Month
}

type TUnit[T TimeUnit] struct {
	unit[T]
}

func Unit[T TimeUnit]() *TUnit[T] {
	return &TUnit[T]{}
}

func (s *TUnit[T]) At(value T) *TUnit[T] {
	s.IsValue = true
	s.Value = &value
	return s
}

func (s *TUnit[T]) Step(step T) *TUnit[T] {
	s.IsStep = true
	s.ValueStep = &step
	return s
}

func (s *TUnit[T]) From(value T) *TUnit[T] {
	s.IsRange = true
	s.ValueFrom = &value
	return s
}

func (s *TUnit[T]) To(value T) *TUnit[T] {
	s.IsRange = true
	s.ValueTo = &value
	return s
}

func (s *TUnit[T]) String(unitName string) string {
	b := strings.Builder{}
	if s.IsStep && s.ValueStep != nil {
		b.WriteString("every ")
		b.WriteString(fmt.Sprint(*s.ValueStep))
		b.WriteString(" ")
		b.WriteString(unitName)
	}
	if s.IsRange {
		if s.IsStep {
			b.WriteString(" ")
		}
		b.WriteString("from ")
		b.WriteString(wrapUnit(*s.ValueFrom, unitName))
		if s.ValueTo != nil {
			b.WriteString(" through ")
			b.WriteString(wrapUnit(*s.ValueTo, unitName))
		}
	}
	if s.IsValue && s.Value != nil {
		b.WriteString("at ")
		b.WriteString(wrapUnit(*s.Value, unitName))
	}
	return b.String()
}
