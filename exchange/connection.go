package exchange

import (
	"fmt"
	"xoney/common"
	"xoney/common/data"
	"xoney/errors"
	"xoney/internal"
	"xoney/internal/structures"
)

type OrderHeap struct {
	heap structures.Heap[Order]
}

func (o OrderHeap) IndexByID(id OrderID) (int, error) {
	for index, order := range o.heap.Members {
		if order.ID() == id {
			return index, nil
		}
	}

	return -1, errors.NewNoLimitOrderError(uint64(id))
}

func (o *OrderHeap) RemoveByID(id OrderID) error {
	index, err := o.IndexByID(id)
	if err != nil {
		return err
	}

	return o.heap.RemoveAt(index)
}

func (o *OrderHeap) OrderByID(id OrderID) (Order, error) {
	index, err := o.IndexByID(id)
	if err != nil {
		return Order{}, err
	}

	return o.heap.Members[index], nil
}

func newOrderHeap(capacity int) OrderHeap {
	return OrderHeap{
		heap: structures.Heap[Order]{
			Members: make([]Order, 0, capacity),
		},
	}
}

type Connector interface {
	PlaceOrder(order Order) error
	CancelOrder(id OrderID) error
	CancelAllOrders() error
	Transfer(quantity float64, currency data.Currency, target data.Exchange) error
	Portfolio() common.Portfolio
	SellAll() error
}

type Simulator interface {
	Connector
	Cleanup() error
	Total() (float64, error)
	UpdatePrice(candle data.InstrumentCandle) error
}

type MarginSimulator struct {
	prices         map[data.Currency]float64
	portfolio      common.Portfolio
	startPortfolio common.Portfolio
	limitOrders    OrderHeap
}

func (s *MarginSimulator) CancelOrder(id OrderID) error {
	return s.limitOrders.RemoveByID(id)
}

func (s *MarginSimulator) PlaceOrder(order Order) error {
	if order.orderType == Market {
		return s.executeMarketOrder(order)
	}

	s.executeLimitOrder(order)

	return nil
}

func (s *MarginSimulator) executeMarketOrder(order Order) error {
	baseQuantity := order.amount
	quoteQuantity := baseQuantity * order.price

	symbol := order.symbol
	quote := symbol.Quote()
	base := symbol.Base()

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

func (s *MarginSimulator) executeLimitOrder(order Order) {
	s.limitOrders.heap.Add(order)
}

func (s *MarginSimulator) updateLimits(high, low float64) error {
	for i, order := range s.limitOrders.heap.Members {
		if order.CrossesPrice(high, low) {
			s.limitOrders.heap.RemoveAt(i)
			// Removing an element in a loop by index in this case is safe
			// because at the first operation we exit the loop,
			// without causing errors/collisions

			return s.executeMarketOrder(order)
		}
	}

	return nil
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

	return s.updateLimits(candle.High, candle.Low)
}

func (s *MarginSimulator) CancelAllOrders() error {
	clear(s.limitOrders.heap.Members)

	return nil
}

func (s *MarginSimulator) Total() (float64, error) {
	return s.portfolio.Total(s.prices)
}

func (s *MarginSimulator) Portfolio() common.Portfolio {
	return s.portfolio
}

func (s *MarginSimulator) SellAll() error {
	mainAsset := s.portfolio.MainCurrency().Asset
	balance := s.portfolio.Assets()

	var firstErr error

	for currency, price := range s.prices {
		pair := data.NewSymbol(currency.Asset, mainAsset, currency.Exchange)

		amount := balance[currency]

		err := s.PlaceOrder(*NewOrder(*pair, Market, Sell, price, amount))
		if firstErr == nil {
			firstErr = fmt.Errorf("error during placing selling order: %w", err)
		}
	}

	return firstErr
}

func (s *MarginSimulator) Cleanup() error {
	err := s.CancelAllOrders()
	s.portfolio = s.startPortfolio

	return fmt.Errorf("order cleanup failed: %w", err)
}

func NewMarginSimulator(portfolio common.Portfolio) MarginSimulator {
	return MarginSimulator{
		prices:         make(map[data.Currency]float64, internal.DefaultCapacity),
		portfolio:      portfolio,
		startPortfolio: portfolio,
		limitOrders:    newOrderHeap(internal.DefaultCapacity),
	}
}

type SpotSimulator struct{ MarginSimulator }

func (s *SpotSimulator) PlaceOrder(order Order) error {
	if order.orderType == Market {
		return s.executeMarketOrder(order)
	}

	s.executeLimitOrder(order)

	return nil
}

func (s *SpotSimulator) executeMarketOrder(order Order) error {
	baseQuantity := order.amount
	quoteQuantity := baseQuantity * order.price

	symbol := order.symbol
	quote := symbol.Quote()
	base := symbol.Base()

	if order.side == Buy {
		return s.executeBuyOrder(base, quote, baseQuantity, quoteQuantity)
	}

	return s.executeSellOrder(base, quote, baseQuantity, quoteQuantity)
}

func (s *SpotSimulator) executeBuyOrder(base, quote data.Currency, baseQuantity, quoteQuantity float64) error {
	if quoteQuantity > s.portfolio.Balance(quote) {
		return errors.NewNotEnoughFundsError(quote.String(), quoteQuantity)
	}

	return s.MarginSimulator.executeBuyOrder(base, quote, baseQuantity, quoteQuantity)
}

func (s *SpotSimulator) executeSellOrder(base, quote data.Currency, baseQuantity, quoteQuantity float64) error {
	if baseQuantity > s.portfolio.Balance(base) {
		return errors.NewNotEnoughFundsError(quote.String(), quoteQuantity)
	}

	return s.MarginSimulator.executeSellOrder(base, quote, baseQuantity, quoteQuantity)
}
