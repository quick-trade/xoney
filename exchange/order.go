package exchange

import (
	"xoney/common/data"
)

type OrderType string

const (
	Market OrderType = "market"
	Limit  OrderType = "limit"
)

type OrderSide string

const (
	Buy  OrderSide = "buy"
	Sell OrderSide = "sell"
)

type Order struct {
	symbol data.Symbol
	orderType  OrderType
	side   OrderSide
	id     uint
	price  float64
	amount float64
}

func (o Order) Symbol() data.Symbol { return o.symbol }
func (o Order) Type() OrderType     { return o.orderType }
func (o Order) Side() OrderSide     { return o.side }
func (o Order) ID() uint            { return o.id }
func (o Order) Price() float64      { return o.price }
func (o Order) Amount() float64     { return o.amount }

func (o Order) IsEqual(other *Order) bool {
	return o.id == other.id
}
