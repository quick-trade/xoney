package events_test

import (
	"errors"
	"fmt"
	"testing"

	"xoney/common"
	"xoney/common/data"
	"xoney/events"
	"xoney/exchange"
)


func usd() data.Currency {
	return data.NewCurrency("USD", "NYSE")
}


func btcUSD() data.Symbol {
	return *data.NewSymbol("BTC", "USD", "BINANCE")
}

type MockConnector struct {
	CancelOrderID   exchange.OrderID
	PlaceOrderCalled int
	PlacedOrder     *exchange.Order
	CancelOrderError error
	PlaceOrderError  error
}

func (m *MockConnector) PlaceOrder(order exchange.Order) error {
	m.PlaceOrderCalled++
	m.PlacedOrder = &order
	return m.PlaceOrderError
}

func (m *MockConnector) CancelOrder(id exchange.OrderID) error {
	m.CancelOrderID = id
	return m.CancelOrderError
}

func (m *MockConnector) CancelAllOrders() error {
	return nil
}

func (m *MockConnector) Transfer(quantity float64, currency data.Currency, target data.Exchange) error {
	return nil
}

func (m *MockConnector) Portfolio() common.Portfolio {
	return common.NewPortfolio(usd())
}

func (m *MockConnector) SellAll() error {
	return nil
}

func (m *MockConnector) GetPrices(symbols []data.Symbol) <-chan exchange.SymbolPrice {
	return nil
}

func TestCancelOrder_Occur(t *testing.T) {
	orderID := exchange.OrderID(123)
	cancelOrder := events.NewCancelOrder(orderID)

	mockConnector := &MockConnector{}
	err := cancelOrder.Occur(mockConnector)

	if err != nil {
		t.Errorf("Unexpected error during CancelOrder.Occur: %v", err)
	}

	if mockConnector.CancelOrderID != orderID {
		t.Errorf("Expected CancelOrderID: %v, got: %v", orderID, mockConnector.CancelOrderID)
	}
}

func TestEditOrder_Occur(t *testing.T) {
	cancelID := exchange.OrderID(456)
	newOrder := exchange.NewOrder(btcUSD(), exchange.Market, exchange.Buy, 50000.0, 0.1)
	editOrder := events.NewEditOrder(cancelID, *newOrder)

	mockConnector := &MockConnector{}
	err := editOrder.Occur(mockConnector)

	if err != nil {
		t.Errorf("Unexpected error during EditOrder.Occur: %v", err)
	}

	if mockConnector.CancelOrderID != cancelID {
		t.Errorf("Expected CancelOrderID: %v, got: %v", cancelID, mockConnector.CancelOrderID)
	}

	if mockConnector.PlaceOrderCalled == 0 {
		t.Errorf("Expected PlaceOrder to be called, but it wasn't")
	}

	if !newOrder.IsEqual(mockConnector.PlacedOrder) {
		t.Errorf("Expected PlacedOrder to be equal to newOrder, but they are not equal")
	}
}

func TestEditOrder_Error(t *testing.T) {
	cancelID := exchange.OrderID(456)
	newOrder := exchange.NewOrder(btcUSD(), exchange.Market, exchange.Buy, 50000.0, 0.1)
	editOrder := events.NewEditOrder(cancelID, *newOrder)

	mockConnector := &MockConnector{
		CancelOrderError: errors.New("mock cancel order error"),
		PlaceOrderError:  errors.New("mock place order error"),
	}

	err := editOrder.Occur(mockConnector)

	if err == nil {
		t.Error("Expected error during EditOrder.Occur, but got nil")
	}

	expectedError := fmt.Sprintf("error canceling order: %v", mockConnector.CancelOrderError)
	if err.Error() != expectedError {
		t.Errorf("Expected error message: %v, got: %v", expectedError, err.Error())
	}
}

func TestEditOrder_Occur_PlaceOrderError(t *testing.T) {
	cancelID := exchange.OrderID(456)
	newOrder := exchange.NewOrder(btcUSD(), exchange.Market, exchange.Buy, 50000.0, 0.1)
	editOrder := events.NewEditOrder(cancelID, *newOrder)

	mockConnector := &MockConnector{
		CancelOrderError: nil,
		PlaceOrderError:  errors.New("mock place order error"),
	}

	err := editOrder.Occur(mockConnector)

	if err == nil {
		t.Error("Expected error during EditOrder.Occur, but got nil")
	}

	expectedError := fmt.Sprintf("error placing order: %v", mockConnector.PlaceOrderError)
	if err.Error() != expectedError {
		t.Errorf("Expected error message: %v, got: %v", expectedError, err.Error())
	}
}
