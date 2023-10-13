package timewalk

import (
	"fmt"
	"strings"
)

type Field struct {
	Schedulers []*Unit `json:"schedulers"`
	Name       string  `json:"name"`
}

func NewField(unit ...*Unit) *Field {
	return &Field{
		Schedulers: unit,
	}
}

func (f *Field) Named(name string) *Field {
	f.Name = name
	for _, u := range f.Schedulers {
		u.Named(name)
	}
	return f
}

func (f *Field) String() string {
	b := strings.Builder{}
	b.WriteString(f.Schedulers[0].String())
	for _, v := range f.Schedulers[1:] {
		b.WriteString(" and ")
		b.WriteString(v.String())
	}
	return b.String()
}

type Unit struct {
	Name      string `json:"name"`
	IsRange   bool   `json:"is_range"`
	IsStep    bool   `json:"is_step"`
	IsValue   bool   `json:"is_value"`
	Value     *int   `json:"value"`
	ValueFrom int    `json:"value_from"`
	ValueTo   *int   `json:"value_to"`
	ValueStep *int   `json:"value_step"`
}

func NewUnit() *Unit {
	return &Unit{Name: "unit"}
}

func (s *Unit) Named(name string) *Unit {
	s.Name = name
	return s
}

func (s *Unit) At(value int) *Unit {
	s.IsValue = true
	s.Value = &value
	return s
}

func (s *Unit) Step(step int) *Unit {
	s.IsStep = true
	s.ValueStep = &step
	return s
}

func (s *Unit) From(value int) *Unit {
	s.IsRange = true
	s.ValueFrom = value
	return s
}

func (s *Unit) To(value int) *Unit {
	s.IsRange = true
	s.ValueTo = &value
	return s
}

func (s *Unit) String() string {
	b := strings.Builder{}
	if s.IsStep && s.ValueStep != nil {
		b.WriteString("every ")
		b.WriteString(fmt.Sprint(*s.ValueStep))
		b.WriteString(" ")
		b.WriteString(s.Name)
	}
	if s.IsRange {
		if s.IsStep {
			b.WriteString(" ")
		}
		b.WriteString("from ")
		b.WriteString(s.Name)
		b.WriteString(" no ")
		b.WriteString(fmt.Sprint(s.ValueFrom))
		if s.ValueTo != nil {
			b.WriteString(" through ")
			b.WriteString(s.Name)
			b.WriteString(" no ")
			b.WriteString(fmt.Sprint(*s.ValueTo))
		}
	}
	if s.IsValue && s.Value != nil {
		b.WriteString("at ")
		b.WriteString(s.Name)
		b.WriteString(" no ")
		b.WriteString(fmt.Sprint(*s.Value))
	}
	return b.String()
}
