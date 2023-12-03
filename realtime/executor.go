package realtime

import (
	"context"
	"fmt"
	"sync"

	"xoney/common/data"
	evnt "xoney/events"
	conn "xoney/exchange"
	st "xoney/strategy"
)

type Runner interface {
	Run(ctx context.Context, system st.Tradable) error
}

type Executor struct {
	connector conn.Connector
	system st.Tradable
}

func (e *Executor) Run(ctx context.Context, system st.Tradable) error {
	if err := e.setup(); err != nil {
		return fmt.Errorf("error during setup: %w", err)
	}

	if err := e.execute(ctx); err != nil {
		return fmt.Errorf("error during executing: %w", err)
	}

	return nil
}

func (e *Executor) setup() error {
	charts := e.getCharts()
	return e.system.Start(charts)
}
func (e *Executor) execute(ctx context.Context) error {
	candleFlow := e.listenCandles(ctx)

	for candle := range candleFlow {
		events, err := e.system.Next(candle)
		if err != nil {
			return err
		}

		if err := e.processEvents(events); err != nil {
			return err
		}
	}

	return nil
}
func (e *Executor) getCharts() data.ChartContainer {
	panic("TODO: implement")
}
func (e *Executor) listenCandles(ctx context.Context) chan data.InstrumentCandle {
	panic("TODO: implement")
}
func (e *Executor) processEvents(events []evnt.Event) error {
	errors := make(chan error, len(events))

	var wg sync.WaitGroup
	wg.Add(len(events))

	for _, ev := range events {
		go func (event evnt.Event) {
			defer wg.Done()
			errors <- event.Occur(e.connector)
		}(ev)
	}

	wg.Wait()

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
