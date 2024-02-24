package exchange

import (
	"fmt"
	"math"

	"github.com/quick-trade/xoney/common"
	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/internal"
	"github.com/quick-trade/xoney/internal/structures"
)

// OrderHeap is a utility structure for efficient management of a collection of orders.
// It leverages the performance characteristics of structures.Heap[T] to provide
// efficient operations such as finding an order by ID and removing orders.
type OrderHeap struct {
	heap structures.Heap[Order]
}

// IndexByID searches for an order with the given id and returns its index in the heap.
// If the order is not found, it returns an error.
func (o OrderHeap) IndexByID(id OrderID) (int, error) {
	for index, order := range o.heap.Members {
		if order.ID() == id {
			return index, nil
		}
	}
	return -1, errors.NewNoLimitOrderError(uint64(id))
}

// RemoveByID removes the order with the given id from the heap.
// It first finds the index of the order and then removes it using the RemoveAt method.
// If the order is not found, it returns an error.
func (o *OrderHeap) RemoveByID(id OrderID) error {
	index, err := o.IndexByID(id)
	if err != nil {
		return err
	}
	return o.heap.RemoveAt(index)
}

// newOrderHeap creates a new OrderHeap with the specified initial capacity.
// This allows for preallocation of memory to improve performance.
func newOrderHeap(capacity int) OrderHeap {
	return OrderHeap{
		heap: structures.Heap[Order]{
			Members: make([]Order, 0, capacity),
		},
	}
}

type SymbolPrice struct {
	Symbol data.Symbol
	Price  float64
}

func NewSymbolPrice(symbol data.Symbol, price float64) *SymbolPrice {
	return &SymbolPrice{
		Symbol: symbol,
		Price:  price,
	}
}

// Connector is a critical interface for interacting with the exchange.
// It provides all the necessary methods to retrieve real-time information from
// the exchange and to send various types of requests, such as placing orders
// or viewing the portfolio balance.
type Connector interface {
	PlaceOrder(order Order) error                                                  // Places a new order on the exchange.
	CancelOrder(id OrderID) error                                                  // Cancels an existing order using its ID.
	CancelAllOrders() error                                                        // Cancels all existing orders.
	Transfer(quantity float64, currency data.Currency, target data.Exchange) error // Transfers a quantity of currency to a target exchange.
	Portfolio() common.Portfolio                                                   // Retrieves the current state of the portfolio.
	SellAll() error                                                                // Executes the sale of all assets in the portfolio.
	GetPrices(symbols []data.Symbol) (<-chan SymbolPrice, <-chan error)            // Retrieves real-time prices for the specified symbols.
}

// Simulator is a key component for backtesting any trading system.
// It is responsible for computing the total balance (profitability) and
// simulating real trades. Prices are updated by the Backtester structure
// based on the historical data provided (ChartContainer).
type Simulator interface {
	Connector
	Cleanup() error                                 // Typically used to reset the simulation to its initial state.
	Total() (float64, error)                        // Calculates the total balance of the portfolio.
	UpdatePrice(candle data.InstrumentCandle) error // Updates the price based on a new candle data.
}

// MarginSimulator is a structure used for testing trading strategies with
// margin trading capabilities. It allows for the simulation of leveraged
// and short positions.
type MarginSimulator struct {
	prices         map[data.Currency]float64 // Current simulated prices for each currency.
	portfolio      common.Portfolio          // The trading portfolio including current holdings.
	startPortfolio common.Portfolio          // The portfolio at the start of the simulation to compare against.
	limitOrders    OrderHeap                 // Heap of limit orders to manage order execution.
	commission     float64                   // Commission fees for executing trades within the simulator.
}

func (s *MarginSimulator) CancelOrder(id OrderID) error {
	return s.limitOrders.RemoveByID(id)
}

func (s *MarginSimulator) PlaceOrder(order Order) error {
	if order.orderType == Market {
		return s.executeMarketOrder(order)
	}

	s.PlaceLimitOrder(order)

	return nil
}

func (s *MarginSimulator) executeMarketOrder(order Order) error {
	baseQuantity := order.amount
	quoteQuantity := baseQuantity * order.price

	symbol := order.symbol
	quote := symbol.Quote()
	base := symbol.Base()

	commission := s.commission * quoteQuantity
	s.portfolio.Decrease(quote, commission)

	if order.side == Buy {
		return s.executeBuyOrder(base, quote, baseQuantity, quoteQuantity)
	}

	return s.executeSellOrder(base, quote, baseQuantity, quoteQuantity)
}

func (s *MarginSimulator) executeBuyOrder(base, quote data.Currency, baseQuantity, quoteQuantity float64) error {
	s.portfolio.Increase(base, baseQuantity)
	s.portfolio.Decrease(quote, quoteQuantity)

	return nil
}

func (s *MarginSimulator) executeSellOrder(base, quote data.Currency, baseQuantity, quoteQuantity float64) error {
	s.portfolio.Decrease(base, baseQuantity)
	s.portfolio.Increase(quote, quoteQuantity)

	return nil
}

func (s *MarginSimulator) PlaceLimitOrder(order Order) {
	s.limitOrders.heap.Add(order)
}

func (s *MarginSimulator) updateLimits(symbol data.Symbol, high, low float64) error {
	var err error

	s.limitOrders.heap.Filter(func(order *Order) bool {
		if order.symbol == symbol && order.CrossesPrice(high, low) {
			execErr := s.executeMarketOrder(*order)
			if execErr != nil {
				err = execErr
			}

			return false
		}

		return true
	})

	return err
}

func (s *MarginSimulator) Transfer(quantity float64, currency data.Currency, target data.Exchange) error {
	if s.portfolio.Balance(currency) < quantity {
		return errors.NewNotEnoughFundsError(currency.String(), quantity)
	}

	s.portfolio.Decrease(currency, quantity)

	currency.Exchange = target
	s.portfolio.Increase(currency, quantity)

	return nil
}

func (s *MarginSimulator) UpdatePrice(candle data.InstrumentCandle) error {
	symbol := candle.Symbol()
	base := symbol.Base()
	quote := symbol.Quote()

	if quote == s.portfolio.MainCurrency() {
		s.prices[base] = candle.Close
	}

	return s.updateLimits(symbol, candle.High, candle.Low)
}

func (s *MarginSimulator) CancelAllOrders() error {
	s.limitOrders.heap.Members = make([]Order, 0, internal.DefaultCapacity)

	return nil
}

func (s *MarginSimulator) Total() (float64, error) {
	return s.portfolio.Total(s.prices)
}

func (s *MarginSimulator) Portfolio() common.Portfolio {
	return s.portfolio.Copy()
}

func (s *MarginSimulator) SellAll() error {
	balance := s.portfolio.Assets()

	var firstErr error

	for currency, price := range s.prices {
		amount := balance[currency]

		if amount == 0 {
			continue
		}

		symbol := data.NewSymbolFromCurrencies(currency, s.portfolio.MainCurrency())

		order, _ := NewOrder(*symbol, Market, orderSideFromBalance(amount), price, math.Abs(amount))
		err := s.PlaceOrder(*order)

		if firstErr == nil && err != nil {
			firstErr = fmt.Errorf("error during placing selling order: %w", err)
		}
	}

	return firstErr
}

func (s *MarginSimulator) GetPrices(symbols []data.Symbol) (<-chan SymbolPrice, <-chan error) {
	prices := make(chan SymbolPrice, len(symbols))
	defer close(prices)

	err := make(chan error)
	defer close(err)

	for _, symbol := range symbols {
		if symbol.Quote() != s.portfolio.MainCurrency() {
			err <- fmt.Errorf("cannot get prices for non-main quote currency, got: %v", symbol.String())

			return prices, err
		}

		prices <- *NewSymbolPrice(symbol, s.prices[symbol.Base()])
	}

	err <- nil

	return prices, err
}

func (s *MarginSimulator) Cleanup() error {
	err := s.CancelAllOrders()
	if err != nil {
		return fmt.Errorf("order cleanup failed: %w", err)
	}
	return nil
}

func NewMarginSimulator(portfolio common.Portfolio, commission float64) MarginSimulator {
	return MarginSimulator{
		prices:         make(common.BaseDistribution, internal.DefaultCapacity),
		portfolio:      portfolio,
		startPortfolio: portfolio.Copy(),
		limitOrders:    newOrderHeap(internal.DefaultCapacity),
		commission:     commission,
	}
}

func orderSideFromBalance(balance float64) OrderSide {
	if balance > 0 {
		return Sell
	}
	return Buy
}

// SpotSimulator represents a trading simulator for spot markets.
// It only supports long positions and does not allow the use of leverage.
type SpotSimulator struct{ MarginSimulator }

func NewSpotSimulator(portfolio common.Portfolio, commission float64) SpotSimulator {
	return SpotSimulator{
		MarginSimulator: NewMarginSimulator(portfolio, commission),
	}
}

func (s *SpotSimulator) PlaceOrder(order Order) error {
	if err := s.validOrder(order); err != nil {
		return fmt.Errorf("error validating order: %w", err)
	}

	if order.orderType == Market {
		return s.executeMarketOrder(order)
	}

	s.PlaceLimitOrder(order)

	return nil
}

func (s *SpotSimulator) validOrder(order Order) error {
	baseQuantity := order.amount
	quoteQuantity := baseQuantity * order.price

	symbol := order.symbol
	quote := symbol.Quote()
	base := symbol.Base()

	if order.side == Buy {
		if quoteQuantity > s.portfolio.Balance(quote) {
			return errors.NewNotEnoughFundsError(quote.String(), quoteQuantity)
		}
	} else {
		if baseQuantity > s.portfolio.Balance(base) {
			return errors.NewNotEnoughFundsError(base.String(), baseQuantity)
		}
	}

	return nil
}
