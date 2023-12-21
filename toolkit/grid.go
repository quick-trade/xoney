package toolkit

import (
	"time"
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
	st "xoney/strategy"
)

type LevelID uint64

type GridLevel struct {
	price float64
	amount float64
	id LevelID
}
func NewGridLevel(price, amount float64, id LevelID) *GridLevel {
	return &GridLevel{
		price: price,
		amount: amount,
		id: id,
	}
}

type grid struct {
	levels   []GridLevel
	executed LevelID
	orders   map[LevelID]exchange.Order
}

func (g *grid) setLevels(levels []GridLevel) {
	added, modified, canceled := g.checkNewLevels(levels)
}
func (g *grid) checkNewLevels(levels []GridLevel) (
	added    []GridLevel,
	modified []GridLevel,
	canceled []GridLevel,
) {
	panic("TODO: implement")
}
func (g *grid) updateOrders(candle data.Candle) ([]events.Event, error) {
	panic("TODO: implement")
}


type GridGenerator interface {
	Instrument() data.Instrument
	MinDuration() time.Duration
	Start(chart data.Chart) error
	Next(candle data.Candle) ([]GridLevel, error)
}


type GridBot struct {
	grid grid
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

	g.grid.setLevels(levels)

	return g.grid.updateOrders(candle.Candle)
}

func (g *GridBot) Start(charts data.ChartContainer) error {
	return g.strategy.Start(charts[g.strategy.Instrument()])
}
