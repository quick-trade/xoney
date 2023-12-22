package toolkit

import (
	"time"
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
	"xoney/internal"

	st "xoney/strategy"
)

type LevelID uint64

type GridLevel struct {
	price  float64
	amount float64
	id     LevelID
}

func NewGridLevel(price, amount float64) *GridLevel {
	return &GridLevel{
		price:  price,
		amount: amount,
		id:     LevelID(internal.RandomUint64()),
	}
}

func orderByLevel(level GridLevel, currentPrice float64, symbol data.Symbol) exchange.Order {
	var side exchange.OrderSide

	if level.price < currentPrice {
		side = exchange.Buy
	} else {
		side = exchange.Sell
	}

	return *exchange.NewOrder(
		symbol,
		exchange.Limit,
		side,
		level.price,
		level.amount,
	)
}

type grid struct {
	symbol   data.Symbol
	levels   []GridLevel
	executed LevelID
	orders   map[LevelID]exchange.Order
}

func (g *grid) setLevels(levels []GridLevel) []events.Event {
	orderEvents := make([]events.Event, 0, len(g.levels))

	// Modified and added levels are processing in g.updateOrders()
	canceled := g.checkCanceledLevels(levels)

	for _, level := range canceled {
		order := g.orders[level.id]

		var cancelOrder events.Event = events.NewCloseOrder(order.ID())
		orderEvents = internal.Append(orderEvents, cancelOrder)

		delete(g.orders, level.id)
	}

	g.levels = levels

	return orderEvents
}

func (g *grid) checkCanceledLevels(levels []GridLevel) []GridLevel {
	// The map is needed to quickly find keys
	paramLevels := make(map[LevelID]struct{})
	for _, level := range levels {
		paramLevels[level.id] = struct{}{}
	}

	var canceledLevels []GridLevel

	for _, level := range g.levels {
		if !internal.Contains(paramLevels, level.id) {
			canceledLevels = internal.Append(canceledLevels, level)
		}
	}

	return canceledLevels
}

func (g *grid) updateOrders(candle data.Candle) []events.Event {
	orderEvents := make([]events.Event, 0, len(g.levels))

	for _, level := range g.levels {
		order, ok := g.orders[level.id]
		if ok && order.CrossesPrice(candle.High, candle.Low) {
			g.executed = level.id
			delete(g.orders, level.id)
		}

		if level.id == g.executed {
			continue
		}

		orderEvents = internal.Append(orderEvents, g.editOrder(level, candle.Close))
	}

	return orderEvents
}

func (g *grid) editOrder(level GridLevel, currPrice float64) events.Event {
	newOrder := orderByLevel(level, currPrice, g.symbol)

	if order, ok := g.orders[level.id]; ok {
		return events.NewEditOrder(order.ID(), newOrder)
	}

	return events.NewOpenOrder(newOrder)
}

func newGrid(symbol data.Symbol) *grid {
	levels := make([]GridLevel, 0)
	orders := make(map[LevelID]exchange.Order, internal.DefaultCapacity)

	return &grid{
		symbol:   symbol,
		levels:   levels,
		executed: 0,
		orders:   orders,
	}
}

type GridGenerator interface {
	Instrument() data.Instrument
	MinDuration() time.Duration
	Start(chart data.Chart) error
	Next(candle data.Candle) ([]GridLevel, error)
}

type GridBot struct {
	grid     grid
	strategy GridGenerator
}

func (g *GridBot) MinDurations() st.Durations {
	return st.Durations{
		g.strategy.Instrument(): g.strategy.MinDuration(),
	}
}

func (g *GridBot) Next(candle data.InstrumentCandle) ([]events.Event, error) {
	levels, err := g.strategy.Next(candle.Candle)
	if err != nil {
		return nil, err
	}

	levelEvents := g.grid.setLevels(levels)
	updateEvents := g.grid.updateOrders(candle.Candle)

	return internal.Append(levelEvents, updateEvents...), nil
}

func (g *GridBot) Start(charts data.ChartContainer) error {
	return g.strategy.Start(charts[g.strategy.Instrument()])
}

func NewGridBot(strategy GridGenerator) *GridBot {
	instrument := strategy.Instrument()
	symbol := instrument.Symbol()

	return &GridBot{
		grid:     *newGrid(symbol),
		strategy: strategy,
	}
}
