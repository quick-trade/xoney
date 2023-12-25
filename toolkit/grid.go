package toolkit

import (
	"math"
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
	spent float64
}

func (g *grid) SetLevels(levels []GridLevel, candle data.Candle) []events.Event {
	if g.checkNewGrid(levels) {
		return g.setNewGrid(levels, candle)
	}
	return []events.Event{}
}
func (g *grid) setNewGrid(levels []GridLevel, candle data.Candle) []events.Event {
	cancelEvents := g.cancelLevelsOrders(g.levels)
	sellAll := g.undoExecuted(candle)

	g.levels = levels

	if sellAll != nil {
		cancelEvents = internal.Append(cancelEvents, sellAll)
	}

	return cancelEvents
}
func (g *grid) cancelLevelsOrders(canceled []GridLevel) []events.Event {
	orderEvents := make([]events.Event, 0, len(g.levels))

	var cancelOrder events.Event
	for _, level := range canceled {
		order, ok := g.orders[level.id]

		if ok {
			cancelOrder = events.NewCloseOrder(order.ID())
			orderEvents = internal.Append(orderEvents, cancelOrder)

			delete(g.orders, level.id)
		}
	}
	return orderEvents
}

func (g *grid) checkNewGrid(levels []GridLevel) bool {
	if len(g.levels) != len(levels) { return false }

	// The map is needed to quickly find ID's
	gridLevelsIDs := make(map[LevelID]struct{})
	for _, level := range g.levels {
		gridLevelsIDs[level.id] = struct{}{}
	}

	for _, level := range levels {
		if !internal.Contains(gridLevelsIDs, level.id) {
			return false
		}
	}

	return true
}

func (g *grid) UpdateOrders(candle data.Candle) []events.Event {
	orderEvents := make([]events.Event, 0, len(g.levels))

	for _, level := range g.levels {
		g.processIfExecuted(level, candle)

		if level.id != g.executed {
			orderEvents = internal.Append(orderEvents, g.editOrder(level, candle.Close))
		}
	}

	return orderEvents
}
func (g *grid) processIfExecuted(level GridLevel, candle data.Candle) {
	order, ok := g.orders[level.id]
	if ok && order.CrossesPrice(candle.High, candle.Low) {
		g.executed = level.id
		g.registerOrderExpenses(order)

		delete(g.orders, level.id)
	}
}
func (g *grid) registerOrderExpenses(order exchange.Order) {
	if order.Side() == exchange.Buy {
		g.spent += order.Amount()
	} else {
		g.spent -= order.Amount()
	}
}
func (g *grid) undoExecuted(candle data.Candle) events.Event {
	price := candle.Close
	amount := math.Abs(g.spent)

	if amount == 0 { return nil }

	var side exchange.OrderSide
	if g.spent > 0 {
		side = exchange.Sell
	} else {
		side = exchange.Buy
	}

	return events.NewOpenOrder(
		*exchange.NewOrder(
			g.symbol,
			exchange.Market,
			side,
			price,
			amount,
		),
	)
}

func (g *grid) editOrder(level GridLevel, currPrice float64) events.Event {
	newOrder := orderByLevel(level, currPrice, g.symbol)

	if order, ok := g.orders[level.id]; ok {
		return events.NewEditOrder(order.ID(), newOrder)
	}

	return events.NewOpenOrder(newOrder)
}

func newGrid() *grid {
	levels := make([]GridLevel, 0)
	orders := make(map[LevelID]exchange.Order, internal.DefaultCapacity)

	return &grid{
		symbol:   data.Symbol{},
		levels:   levels,
		executed: 0,
		orders:   orders,
	}
}

type GridGenerator interface {
	MinDuration(timeframe data.TimeFrame) time.Duration
	Start(chart data.Chart) error
	Next(candle data.Candle) ([]GridLevel, error) // new levels can be nil
}

type GridBot struct {
	grid     grid
	strategy GridGenerator
	instrument data.Instrument
}

func (g *GridBot) MinDurations() st.Durations {
	return st.Durations{
		g.instrument: g.strategy.MinDuration(g.instrument.Timeframe()),
	}
}

func (g *GridBot) Next(candle data.InstrumentCandle) ([]events.Event, error) {
	levels, err := g.strategy.Next(candle.Candle)
	if err != nil {
		return nil, err
	}

	var levelEvents []events.Event
	if levels != nil {
		levelEvents = g.grid.SetLevels(levels, candle.Candle)
	}

	updateEvents := g.grid.UpdateOrders(candle.Candle)


	return internal.Append(levelEvents, updateEvents...), nil
}

func (g *GridBot) Start(charts data.ChartContainer) error {
	return g.strategy.Start(charts[g.instrument])
}

func NewGridBot(strategy GridGenerator, instrument data.Instrument) *GridBot {
	return &GridBot{
		grid:     *newGrid(),
		strategy: strategy,
		instrument: instrument,
	}
}