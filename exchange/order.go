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
	symbol    data.Symbol
	orderType OrderType
	side      OrderSide
	id        uint64
	price     float64
	amount    float64
}

func (o Order) Symbol() data.Symbol { return o.symbol }
func (o Order) Type() OrderType     { return o.orderType }
func (o Order) Side() OrderSide     { return o.side }
func (o Order) ID() uint64          { return o.id }
func (o Order) Price() float64      { return o.price }
func (o Order) Amount() float64     { return o.amount }

func (o Order) IsEqual(other *Order) bool {
	return o.id == other.id
}

func NewOrder(symbol data.Symbol, orderType OrderType, side OrderSide, price, amount float64) *Order {
	return &Order{
		symbol:    symbol,
		orderType: orderType,
		side:      side,
		price:     price,
		amount:    amount,
		id:        0, // BUGFIX: make exchange id or random
	}
}

func crossesPrice(order Order, high, low float64) bool {
	if order.side == Buy {
		return low < order.price
	}
	return high > order.price
}
