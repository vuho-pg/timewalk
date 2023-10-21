package timewalk

type unit[T comparable] struct {
	IsRange   bool `json:"is_range,omitempty"`
	IsStep    bool `json:"is_step,omitempty"`
	IsValue   bool `json:"is_value,omitempty"`
	Value     *T   `json:"value,omitempty"`
	ValueFrom *T   `json:"value_from,omitempty"`
	ValueTo   *T   `json:"value_to,omitempty"`
	ValueStep *T   `json:"value_step,omitempty"`
}
