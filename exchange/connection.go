package exchange

import (
	"xoney/common"
	"xoney/common/data"
	"xoney/errors"
	"xoney/internal/structures"
)

type OrderHeap struct {
	heap structures.Heap[Order]
}

func (o OrderHeap) IndexByID(id uint) (int, error) {
	for index, order := range o.heap.Members {
		if order.ID() == id {
			return index, nil
		}
	}

	return -1, errors.NewNoLimitOrderError(id)
}

func (o *OrderHeap) RemoveByID(id uint) error {
	index, err := o.IndexByID(id)
	if err != nil {
		return err
	}

	return o.heap.RemoveAt(index)
}

type Connector interface {
	PlaceOrder(order Order) error
	CancelOrder(id uint) error
	Transfer(quantity float64, currency data.Currency, target data.Exchange) error
	Portfolio() *common.Portfolio
}

type Simulator struct {
	prices      map[data.Currency]float64
	portfolio   common.Portfolio
	limitOrders OrderHeap
}

func (s *Simulator) CancelOrder(id uint) error {
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

func (s *Simulator) Portfolio() *common.Portfolio {
	return &s.portfolio
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

func NewSimulator() *Simulator {
	panic("TODO: implement")
}
