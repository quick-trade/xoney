package events

import (
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
