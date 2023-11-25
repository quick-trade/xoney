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
	return -1, errors.ValueNotFoundError{}
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
	portfolio common.Portfolio
	limitOrders OrderHeap
}

func (s *Simulator) CancelOrder(id uint) error {
	return s.limitOrders.RemoveByID(id)
}

func (s *Simulator) PlaceOrder(order Order) error {
	if order.type_ == Market {
		return s.executeMarketOrder(order)
	}
	return s.executeLimitOrder(order)
}
func (s *Simulator) executeMarketOrder(order Order) error {
	panic("TODO: Implement")
}
func (s *Simulator) executeLimitOrder(order Order) error {
	panic("TODO: Implement")
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

func NewSimulator() *Simulator {
	return nil
}
