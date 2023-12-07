package realtime

import (
	"context"
	"fmt"
	"sync"

	"xoney/common/data"
	ev "xoney/events"
	conn "xoney/exchange"
	"xoney/internal"
	st "xoney/strategy"
)

type DataSupplier interface {
	GetCharts(instruments st.Durations) (data.ChartContainer, error)
	StreamCandles(ctx context.Context, instruments []data.Instrument) <-chan data.InstrumentCandle
}

type Runner interface {
	Run(ctx context.Context, system st.Tradable) error
}

type Executor struct {
	connector conn.Connector
	supplier  DataSupplier
	system    st.Tradable
}

func NewExecutor(connector conn.Connector, supplier DataSupplier) *Executor {
	return &Executor{
		connector: connector,
		supplier:  supplier,
		system:    nil,
	}
}

func (e *Executor) Run(ctx context.Context, system st.Tradable) error {
	if err := e.setup(); err != nil {
		return fmt.Errorf("error during setup: %w", err)
	}

	if err := e.execute(ctx); err != nil {
		return fmt.Errorf("error during executing: %w", err)
	}

	if err := e.stop(); err != nil {
		return fmt.Errorf("error during stopping: %w", err)
	}

	return nil
}

func (e *Executor) setup() error {
	charts, err := e.getCharts()
	if err != nil {
		return fmt.Errorf("error fetching charts: %w", err)
	}

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

func (e *Executor) getCharts() (data.ChartContainer, error) {
	return e.supplier.GetCharts(e.system.MinDurations())
}

func (e *Executor) listenCandles(ctx context.Context) <-chan data.InstrumentCandle {
	durations := e.system.MinDurations()
	instruments := internal.MapKeys(durations)

	return e.supplier.StreamCandles(ctx, instruments)
}

func (e *Executor) processEvents(events []ev.Event) error {
	errors := make(chan error, len(events))

	var wg sync.WaitGroup

	wg.Add(len(events))

	for _, event := range events {
		go func(event ev.Event) {
			defer wg.Done()
			errors <- event.Occur(e.connector)
		}(event)
	}

	wg.Wait()

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) stop() error {
	firstErr := e.connector.CancelAllOrders()

	if err := e.connector.SellAll(); err != nil {
		if firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
