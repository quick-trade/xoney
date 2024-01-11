package toolkit

import (
	"testing"
	"time"

	"xoney/common/data"
	"xoney/common"
	"xoney/events"
	"xoney/exchange"
)


func btcUSD() data.Symbol {
	return *data.NewSymbol("BTC", "USD", "BINANCE")
}

func timeframe() data.TimeFrame {
	timeframe, _ := data.NewTimeFrame(time.Hour, "1h")

	return *timeframe
}

func btcUSD1h() data.Instrument {
	return data.NewInstrument(btcUSD(), timeframe())
}

func startTime() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestNewGridLevel(t *testing.T) {
	price := 100.0
	amount := 0.5

	level := NewGridLevel(price, amount)

	if level.price != price {
		t.Errorf("Expected price: %v, got: %v", price, level.price)
	}

	if level.amount != amount {
		t.Errorf("Expected amount: %v, got: %v", amount, level.amount)
	}

	if level.id == 0 {
		t.Errorf("Expected non-zero level ID, got: %v", level.id)
	}
}
type MockGridGenerator struct {
	counter int
}

func (g *MockGridGenerator) MinDuration(timeframe data.TimeFrame) time.Duration {
	return time.Hour
}

func (g *MockGridGenerator) Start(chart data.Chart) error {
	g.counter = 0

	return nil
}

func (g *MockGridGenerator) Next(candle data.Candle) ([]GridLevel, error) {
	g.counter++

	if g.counter == 5 {
		// Return a new grid with price 105 and amount 0.1 on the fifth candle
		return []GridLevel{*NewGridLevel(105, 0.1)}, nil
	} else if g.counter == 1 {
		return []GridLevel{*NewGridLevel(100, 0.5)}, nil
	}
	return nil, nil
}

type mockConnector struct {
	PlaceOrderFunc      func(order exchange.Order) error
	CancelOrderFunc     func(id exchange.OrderID) error
	CancelAllOrdersFunc func() error
	TransferFunc        func(quantity float64, currency data.Currency, target data.Exchange) error
	PortfolioFunc       func() common.Portfolio
	SellAllFunc         func() error
}

func (m *mockConnector) PlaceOrder(order exchange.Order) error {
	if m.PlaceOrderFunc != nil {
		return m.PlaceOrderFunc(order)
	}
	return nil
}

func (m *mockConnector) CancelOrder(id exchange.OrderID) error {
	if m.CancelOrderFunc != nil {
		return m.CancelOrderFunc(id)
	}
	return nil
}

func (m *mockConnector) CancelAllOrders() error {
	if m.CancelAllOrdersFunc != nil {
		return m.CancelAllOrdersFunc()
	}
	return nil
}

func (m *mockConnector) Transfer(quantity float64, currency data.Currency, target data.Exchange) error {
	if m.TransferFunc != nil {
		return m.TransferFunc(quantity, currency, target)
	}
	return nil
}

func (m *mockConnector) Portfolio() common.Portfolio {
	if m.PortfolioFunc != nil {
		return m.PortfolioFunc()
	}
	return common.Portfolio{}
}

func (m *mockConnector) SellAll() error {
	if m.SellAllFunc != nil {
		return m.SellAllFunc()
	}
	return nil
}


func TestGridBot_Next(t *testing.T) {
	generator := &MockGridGenerator{counter: 0}

	// Create a GridBot
	bot := NewGridBot(generator, btcUSD1h())

	// Create a fake candle
	candle := data.NewCandle(90, 110, 85, 105, 0, startTime())

	// Call the Next method
	event, err := bot.Next(*data.NewInstrumentCandle(*candle, btcUSD1h()))
	if err != nil {
		t.Fatalf("Unexpected error in Next: %v", err)
	}
	if event == nil {
		t.Fatalf("Expected an event, got nil")
	}
	Events := event.(*events.Sequential).Events()

	// Check if events were generated
	if len(Events) != 1 {
		t.Errorf("Expected 1 event, got %v", len(Events))
	}

	// Check if the event is of type OpenOrder
	openOrder, ok := Events[0].(*events.OpenOrder)

	if !ok {
		t.Errorf("Expected OpenOrder event, got %v", Events[0])
	} else {
		// Create a mock connector
		mockConnector := mockConnector{
			PlaceOrderFunc: func(order exchange.Order) error {
				// Check the order attributes
				if order.Amount() != 0.5 {
					t.Errorf("Expected order amount to be 0.5, got %v", order.Amount())
				}
				if order.Price() != 100 {
					t.Errorf("Expected order price to be 100, got %v", order.Price())
				}
				return nil
			},
		}
		openOrder.Occur(&mockConnector)
	}

	for i := 1; i < 4; i++ {
		candle.TimeClose = candle.TimeClose.Add(timeframe().Duration)
		event, err = bot.Next(*data.NewInstrumentCandle(*candle, btcUSD1h()))
		Events = event.(*events.Sequential).Events()
		if err != nil {
			t.Fatalf("Unexpected error in Next: %v", err)
		}

		// Check if events were generated
		if len(Events) != 0 {
			t.Errorf("Expected 0 events, got %v", len(Events))
		}
	}

	candle.TimeClose = candle.TimeClose.Add(timeframe().Duration)

	event, err = bot.Next(*data.NewInstrumentCandle(*candle, btcUSD1h()))
	Events = event.(*events.Sequential).Events()
	if err != nil {
		t.Fatalf("Unexpected error in Next: %v", err)
	}

	if len(Events) != 2 {
		t.Errorf("Expected 2 events, got %v", len(Events))
	}

	closeOrder, ok := Events[0].(*events.OpenOrder)
	if !ok {
		t.Errorf("Expected OpenOrder event, got %v", Events[0])
	} else {
		mockConnector := &mockConnector{
			PlaceOrderFunc: func(order exchange.Order) error {
				if order.Amount() != 0.5 {
					t.Errorf("Expected order amount to be 0.5, got %v", order.Amount())
				}
				if order.Price() != candle.Close {
					t.Errorf("Expected order price to be 100, got %v", order.Price())
				}
				if order.Type() != exchange.Market {
					t.Errorf("Expected order type to be Market, got %v", order.Type())
				}
				if order.Side() != exchange.Sell {
					t.Errorf("Expected order side to be Sell, got %v", order.Side())
				}
				return nil
			},
		}
		closeOrder.Occur(mockConnector)
	}

	openOrder, ok = Events[1].(*events.OpenOrder)
	if !ok {
		t.Errorf("Expected OpenOrder event, got %v", Events[0])
	}

	mockConnector := &mockConnector{
		PlaceOrderFunc: func(order exchange.Order) error {
			if order.Amount() != 0.1 {
				t.Errorf("Expected order amount to be 0.5, got %v", order.Amount())
			}
			if order.Price() != 105 {
				t.Errorf("Expected order price to be 100, got %v", order.Price())
			}
			if order.Type() != exchange.Limit {
				t.Errorf("Expected order type to be Market, got %v", order.Type())
			}
			if order.Side() != exchange.Sell {
				t.Errorf("Expected order side to be Sell, got %v", order.Side())
			}
			return nil
		},
	}
	openOrder.Occur(mockConnector)
}

func TestGridBot_Duration(t *testing.T) {
	generator := &MockGridGenerator{counter: 0}

	// Create a GridBot
	bot := NewGridBot(generator, btcUSD1h())

	duration := bot.MinDurations()[btcUSD1h()]
	if duration != time.Hour {
		t.Errorf("Expected duration %v, got %v", time.Hour, duration)
	}
}


// Test GridBot MinDurations method
func TestGridBot_MinDurations(t *testing.T) {
	generator := &MockGridGenerator{counter: 0}

	// Create a GridBot
	bot := NewGridBot(generator, btcUSD1h())

	// Call the MinDurations method
	durations := bot.MinDurations()

	// Check if the duration for btcUSD1h is correct
	expectedDuration := time.Hour
	if durations[btcUSD1h()] != expectedDuration {
		t.Errorf("Expected duration %v for btcUSD1h, got %v", expectedDuration, durations[btcUSD1h()])
	}
}
