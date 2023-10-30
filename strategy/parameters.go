package strategy

type Parameter interface {
	Name() string
}
type IntParameter struct {
	name string
	min  int
	max  int
}

func (i *IntParameter) Name() string { return i.name }
func NewIntParameter(name string, min int, max int) *IntParameter {
	return &IntParameter{name: name, min: min, max: max}
}

type FloatParameter struct {
	name string
	min  float64
	max  float64
}

func (f *FloatParameter) Name() string { return f.name }
func (f *FloatParameter) Min() float64 { return f.min }
func NewFloatParameter(name string, min float64, max float64) *FloatParameter {
	return &FloatParameter{name: name, min: min, max: max}
}

type CategoricalParameter[T any] struct {
	name       string
	Categories []T
}

func (c *CategoricalParameter[T]) Name() string {
	return c.name
}
