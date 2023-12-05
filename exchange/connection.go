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

func (o OrderHeap) IndexByID(id uint64) (int, error) {
	for index, order := range o.heap.Members {
		if order.ID() == id {
			return index, nil
		}
	}

	return -1, errors.NewNoLimitOrderError(id)
}

func (o *OrderHeap) RemoveByID(id uint64) error {
	index, err := o.IndexByID(id)
	if err != nil {
		return err
	}

	return o.heap.RemoveAt(index)
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
	CancelOrder(id uint64) error
	CancelAllOrders() error
	Transfer(quantity float64, currency data.Currency, target data.Exchange) error
	Portfolio() common.Portfolio
	SellAll() error
}

type Simulator struct {
	prices      map[data.Currency]float64
	portfolio   common.Portfolio
	startPortfolio common.Portfolio
	limitOrders OrderHeap
}

func (s *Simulator) CancelOrder(id uint64) error {
	return s.limitOrders.RemoveByID(id)
}

func (s *Simulator) PlaceOrder(order Order) error {
	if order.orderType == Market {
		return s.executeMarketOrder(order)
	}

	s.executeLimitOrder(order)

	return nil
}

func (s *Simulator) executeMarketOrder(order Order) error {
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

func (s *Simulator) executeBuyOrder(base, quote data.Currency, baseQuantity, quoteQuantity float64) error {
	if quoteQuantity > s.portfolio.Balance(quote) {
		return errors.NewNotEnoughFundsError(quote.String(), quoteQuantity)
	}

	s.portfolio.Increase(base, baseQuantity)
	s.portfolio.Decrease(quote, quoteQuantity)

	return nil
}

func (s *Simulator) executeSellOrder(base, quote data.Currency, baseQuantity, quoteQuantity float64) error {
	if baseQuantity > s.portfolio.Balance(base) {
		return errors.NewNotEnoughFundsError(quote.String(), quoteQuantity)
	}

	s.portfolio.Decrease(base, baseQuantity)
	s.portfolio.Increase(quote, quoteQuantity)

	return nil
}

func (s *Simulator) executeLimitOrder(order Order) {
	s.limitOrders.heap.Add(order)
}

func (s *Simulator) updateLimits(high, low float64) error {
	for i, order := range s.limitOrders.heap.Members {
		if crossesPrice(order, high, low) {
			s.limitOrders.heap.RemoveAt(i)
			return s.executeMarketOrder(order)
		}
	}
	return nil
}

func (s *Simulator) Transfer(quantity float64, currency data.Currency, target data.Exchange) error {
	if s.portfolio.Balance(currency) < quantity {
		return errors.NewNotEnoughFundsError(currency.String(), quantity)
	}

	s.portfolio.Decrease(currency, quantity)

	currency.Exchange = target
	s.portfolio.Increase(currency, quantity)

	return nil
}

func (s *Simulator) UpdatePrice(candle data.InstrumentCandle) error {
	symbol := candle.Symbol()
	base := symbol.Base()
	quote := symbol.Quote()

	if quote == s.portfolio.MainCurrency() {
		s.prices[base] = candle.Close
	}

	return s.updateLimits(candle.High, candle.Low)
}

func (s *Simulator) CancelAllOrders() error {
	clear(s.limitOrders.heap.Members)
	return nil
}

func (s *Simulator) Total() (float64, error) {
	return s.portfolio.Total(s.prices)
}

func (s *Simulator) Portfolio() common.Portfolio {
	return s.portfolio
}

func (s *Simulator) SellAll() error {
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
func (s *Simulator) Cleanup() {
	s.CancelAllOrders()
	s.portfolio = s.startPortfolio
}

func NewSimulator(portfolio common.Portfolio) Simulator {
	return Simulator{
		prices:      make(map[data.Currency]float64, internal.DefaultCapacity),
		portfolio:   portfolio,
		startPortfolio: portfolio,
		limitOrders: newOrderHeap(internal.DefaultCapacity),
	}
}
