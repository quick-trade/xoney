package validation

import (
	"xoney/common/data"
)

type WFSampler struct{}

func (w *WFSampler) Samples(data data.ChartContainer) ([]SamplePair, error) {
	panic("TODO")
}
