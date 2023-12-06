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

type CloseOrder struct {
	id uint64
}
func (o *CloseOrder) Occur(connector exchange.Connector) error {
	return connector.CancelOrder(o.id)
}
func NewCloseOrder(id uint64) *CloseOrder {
	return &CloseOrder{id: id}
}
