package exchange

import (
	"xoney/common/data"
	"xoney/internal"
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

type OrderID uint64

type Order struct {
	symbol     data.Symbol
	orderType  OrderType
	side       OrderSide
	internalID OrderID
	price      float64
	amount     float64
}

func (o Order) Symbol() data.Symbol { return o.symbol }
func (o Order) Type() OrderType     { return o.orderType }
func (o Order) Side() OrderSide     { return o.side }
func (o Order) ID() OrderID         { return o.internalID }
func (o Order) Price() float64      { return o.price }
func (o Order) Amount() float64     { return o.amount }

func (o Order) IsEqual(other *Order) bool {
	return o.internalID == other.internalID
}

func (o Order) CrossesPrice(high, low float64) bool {
	if o.side == Buy {
		return low < o.price
	}

	return high > o.price
}

func NewOrder(symbol data.Symbol, orderType OrderType, side OrderSide, price, amount float64) *Order {
	return &Order{
		symbol:     symbol,
		orderType:  orderType,
		side:       side,
		price:      price,
		amount:     amount,
		internalID: OrderID(internal.RandomUint64()),
	}
}
