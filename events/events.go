package events

import (
	"fmt"
	"xoney/exchange"
)

type Event interface {
	Occur(connector exchange.Connector) error
}

type OpenOrder struct {
	order exchange.Order
}

func (o *OpenOrder) Occur(connector exchange.Connector) error {
	return connector.PlaceOrder(o.order)
}

func NewOpenOrder(order exchange.Order) *OpenOrder {
	return &OpenOrder{order: order}
}

type CancelOrder struct {
	id exchange.OrderID
}

func (o *CancelOrder) Occur(connector exchange.Connector) error {
	return connector.CancelOrder(o.id)
}

func NewCloseOrder(id exchange.OrderID) *CancelOrder {
	return &CancelOrder{id: id}
}

type EditOrder struct {
	cancelID exchange.OrderID
	order    exchange.Order
}

func (e *EditOrder) Occur(connector exchange.Connector) error {
	if err := connector.CancelOrder(e.cancelID); err != nil {
		return fmt.Errorf("error canceling order: %w", err)
	}

	if err := connector.PlaceOrder(e.order); err != nil {
		return fmt.Errorf("error placing order: %w", err)
	}

	return nil
}

func NewEditOrder(cancelID exchange.OrderID, newOrder exchange.Order) *EditOrder {
	return &EditOrder{
		cancelID: cancelID,
		order:    newOrder,
	}
}

type Sequential struct {
	actions []Event
}
func (s *Sequential) Occur(connector exchange.Connector) error {
	for _, action := range s.actions {
		if err := action.Occur(connector); err != nil {
			return err
		}
	}
	return nil
}
func (s *Sequential) Add(actions ...Event) {
	s.actions = append(s.actions, actions...)
}
func (s *Sequential) Events() []Event {
	return s.actions
}
func NewSequential(actions ...Event) *Sequential {
	return &Sequential{actions: actions}
}
