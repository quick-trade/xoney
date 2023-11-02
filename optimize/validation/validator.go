package validation

import (
	"fmt"
	"sync"
	"xoney/common/data"

	st "xoney/strategy"
)

type Validator struct {
	charts  data.ChartContainer
	sampler *Sampler
}

func (v *Validator) Validate(system st.Optimizable) (chan EquityResult, error) {
	samples, err := (*v.sampler).Samples(v.charts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve samples: %w", err)
	}

	nSamples := len(samples)
	equities := make(chan EquityResult, nSamples)

	// Running validation process
	var wg sync.WaitGroup

	wg.Add(nSamples)

	validate := func(sPair SamplePair) {
		defer wg.Done()
		equities <- sPair.test(system)
	}
	for _, sp := range samples {
		go validate(sp)
	}

	// Waiting until all data to be written
	go func() {
		wg.Wait()
		close(equities)
	}()

	return equities, nil
}
