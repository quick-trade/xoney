package exchange

import (
	"fmt"
	"math"

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

func newOrderHeap(capacity int) OrderHeap {
	return OrderHeap{
		heap: structures.Heap[Order]{
			Members: make([]Order, 0, capacity),
		},
	}
}

type SymbolPrice struct {
	Symbol data.Symbol
	Price  float64
}
func NewSymbolPrice(symbol data.Symbol, price float64) *SymbolPrice {
	return &SymbolPrice{
		Symbol: symbol,
		Price:  price,
	}
}
type Connector interface {
	PlaceOrder(order Order) error
	CancelOrder(id OrderID) error
	CancelAllOrders() error
	Transfer(quantity float64, currency data.Currency, target data.Exchange) error
	Portfolio() common.Portfolio
	SellAll() error
	GetPrices(symbols []data.Symbol) <-chan SymbolPrice
}

type Simulator interface {
	Connector
	Cleanup() error
	Total() (float64, error)
	UpdatePrice(candle data.InstrumentCandle) error
}

type MarginSimulator struct {
	prices         map[data.Symbol]float64
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

	s.PlaceLimitOrder(order)

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

func (s *MarginSimulator) PlaceLimitOrder(order Order) {
	s.limitOrders.heap.Add(order)
}

func (s *MarginSimulator) updateLimits(symbol data.Symbol, high, low float64) error {
	var err error

	s.limitOrders.heap.Filter(func(order *Order) bool {
		if order.symbol == symbol && order.CrossesPrice(high, low) {
			execErr := s.executeMarketOrder(*order)
			if execErr != nil {
				err = execErr
			}

			return false
		}

		return true
	})

	return err
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
	quote := symbol.Quote()

	if quote == s.portfolio.MainCurrency() {
		s.prices[symbol] = candle.Close
	}

	return s.updateLimits(symbol, candle.High, candle.Low)
}

func (s *MarginSimulator) CancelAllOrders() error {
	s.limitOrders.heap.Members = make([]Order, 0, internal.DefaultCapacity)

	return nil
}

func (s *MarginSimulator) Total() (float64, error) {
	return s.portfolio.Total(s.prices)
}

func (s *MarginSimulator) Portfolio() common.Portfolio {
	return s.portfolio.Copy()
}

func (s *MarginSimulator) SellAll() error {
	balance := s.portfolio.Assets()

	var firstErr error

	for symbol, price := range s.prices {
		amount := balance[symbol.Base()]

		err := s.PlaceOrder(
			*NewOrder(symbol, Market, orderSideFromBalance(amount), price, math.Abs(amount)),
		)
		if firstErr == nil && err != nil {
			firstErr = fmt.Errorf("error during placing selling order: %w", err)
		}
	}

	return firstErr
}
func orderSideFromBalance(balance float64) OrderSide {
	if balance > 0 {
		return Sell
	}
	return Buy
}

func (s *MarginSimulator) GetPrices(symbols []data.Symbol) <-chan SymbolPrice {
	prices := make(chan SymbolPrice, len(symbols))

	defer close(prices)

	for _, symbol := range symbols {
		prices <- SymbolPrice{
			Symbol: symbol,
			Price:  s.prices[symbol],
		}
	}

	return prices
}

func (s *MarginSimulator) Cleanup() error {
	err := s.CancelAllOrders()
	if err != nil {
		return fmt.Errorf("order cleanup failed: %w", err)
	}
	return nil
}

func NewMarginSimulator(portfolio common.Portfolio) MarginSimulator {
	return MarginSimulator{
		prices:         make(map[data.Symbol]float64, internal.DefaultCapacity),
		portfolio:      portfolio,
		startPortfolio: portfolio.Copy(),
		limitOrders:    newOrderHeap(internal.DefaultCapacity),
	}
}

type SpotSimulator struct{ MarginSimulator }

func NewSpotSimulator(portfolio common.Portfolio) SpotSimulator {
	return SpotSimulator{
		MarginSimulator: NewMarginSimulator(portfolio),
	}
}

func (s *SpotSimulator) PlaceOrder(order Order) error {
	if err := s.validOrder(order); err != nil {
		return fmt.Errorf("error validating order: %w", err)
	}

	if order.orderType == Market {
		return s.executeMarketOrder(order)
	}

	s.PlaceLimitOrder(order)

	return nil
}

func (s *SpotSimulator) validOrder(order Order) error {
	baseQuantity := order.amount
	quoteQuantity := baseQuantity * order.price

	symbol := order.symbol
	quote := symbol.Quote()
	base := symbol.Base()

	if order.side == Buy {
		if quoteQuantity > s.portfolio.Balance(quote) {
			return errors.NewNotEnoughFundsError(quote.String(), quoteQuantity)
		}
	} else {
		if baseQuantity > s.portfolio.Balance(base) {
			return errors.NewNotEnoughFundsError(base.String(), baseQuantity)
		}
	}

	return nil
}
