// toolkit_test.go
package toolkit_test

import (
	"math"
	"testing"

	"github.com/quick-trade/xoney/common"
	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/exchange"
	"github.com/quick-trade/xoney/internal"
	"github.com/quick-trade/xoney/toolkit"
)

const epsilon = 0.001

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func equalFloatMaps(expected, actual map[data.Currency]float64) bool {
	if len(expected) != len(actual) {
		return false
	}

	for currency, expectedValue := range expected {
		actualValue, ok := actual[currency]
		if !ok || !closeEnough(expectedValue, actualValue) {
			return false
		}
	}

	return true
}

// TotalCapital calculates the absolute sum of a portfolio.
func TotalCapital(portfolio common.BaseDistribution) float64 {
	total := 0.0

	for _, amount := range portfolio {
		total += math.Abs(amount)
	}

	return total
}

// CurrentWeights calculates the weights of each currency in the portfolio.
func CurrentWeights(portfolio common.BaseDistribution) toolkit.PortfolioWeights {
	totalCapital := TotalCapital(portfolio)
	weights := make(map[data.Currency]toolkit.BaseWeight)

	for currency, amount := range portfolio {
		weights[currency] = toolkit.BaseWeight(amount / totalCapital)
	}

	dist, err := toolkit.NewPortfolioWeights(weights)
	if err != nil {
		panic(err)
	}
	return *dist
}

func btc() data.Currency {
	return data.NewCurrency("BTC", "EXCHANGE")
}

func eth() data.Currency {
	return data.NewCurrency("ETH", "EXCHANGE")
}

func usd() data.Currency {
	return data.NewCurrency("USD", "EXCHANGE")
}

func TestSynchronize(t *testing.T) {
	currentBasePortfolio := common.BaseDistribution{
		usd(): 1. / 100.,
		btc(): 1. / 500.,
		eth(): 1. / 10.,
	}

	targetWeights := CurrentWeights(currentBasePortfolio)

	currentPrices := map[data.Currency]float64{
		usd(): 100,
		btc(): 500,
		eth(): 10,
	}

	expectedResult := currentBasePortfolio

	result, err := targetWeights.Synchronize(currentBasePortfolio, currentPrices, usd())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !equalFloatMaps(result, expectedResult) {
		t.Errorf("Expected %v, got %v", expectedResult, result)
	}
}

type MockConnector struct {
	orders    []*exchange.Order
	portfolio common.Portfolio
	prices    map[data.Symbol]float64
}

func (mc *MockConnector) PlaceOrder(order exchange.Order) error {
	mc.orders = internal.Append(mc.orders, &order)

	return nil
}

func (mc *MockConnector) CheckOrder(amount float64, symbol data.Symbol, side exchange.OrderSide, orderType exchange.OrderType) (bool, *exchange.Order) {
	for _, order := range mc.orders {
		if order.Symbol() == symbol {
			if order.Side() != side {
				continue
			}
			if order.Type() != orderType {
				continue
			}
			if order.Amount() != amount {
				continue
			}

			return true, order
		}
	}

	return false, nil
}

func (m *MockConnector) CancelOrder(id exchange.OrderID) error {
	return nil
}

func (m *MockConnector) CancelAllOrders() error {
	return nil
}

func (m *MockConnector) Transfer(quantity float64, currency data.Currency, target data.Exchange) error {
	return nil
}

func (m *MockConnector) Portfolio() common.Portfolio {
	return m.portfolio
}

func (m *MockConnector) SellAll() error {
	return nil
}

func (m *MockConnector) GetPrices(symbols []data.Symbol) (<-chan exchange.SymbolPrice, <-chan error) {
	priceChan := make(chan exchange.SymbolPrice)
	errChan := make(chan error, 1) // Buffer of 1 for non-blocking send on error

	go func() {
		defer close(priceChan)
		defer close(errChan)
		for _, symbol := range symbols {
			price, ok := m.prices[symbol]
			if !ok {
				errChan <- errors.NewNoPriceError(symbol.String())
				return
			}
			priceChan <- *exchange.NewSymbolPrice(symbol, price)
		}
	}()

	return priceChan, errChan
}

func (m *MockConnector) SetPrice(symbol data.Symbol, price float64) {
	m.prices[symbol] = price
}

// NewMockConnector creates a new instance of MockConnector with initialized fields.
func NewMockConnector(portfolio common.Portfolio) *MockConnector {
	return &MockConnector{
		orders:    make([]*exchange.Order, 0),
		portfolio: portfolio,
		prices:    make(map[data.Symbol]float64, internal.DefaultCapacity),
	}
}

// mapSubtract takes two maps of data.Symbol to float64 and returns a new map with the subtraction result.
func mapSubtract[T comparable](a, b map[T]float64) map[T]float64 {
	result := make(map[T]float64)
	for k, v := range a {
		if bv, ok := b[k]; ok {
			result[k] = v - bv
		} else {
			result[k] = v
		}
	}
	for k, bv := range b {
		if _, ok := a[k]; !ok {
			result[k] = -bv
		}
	}
	return result
}

// convertToPortfolio converts a BaseDistribution to a common.Portfolio using a base asset.
func convertToPortfolio(distribution common.BaseDistribution, mainCurrency data.Currency) common.Portfolio {
	portfolio := common.NewPortfolio(mainCurrency)

	for currency, amount := range distribution {
		portfolio.Set(currency, amount)
	}

	return portfolio
}

func TestRebalanceMarketOrders(t *testing.T) {
	// Hardcoded initial distribution and desired weights
	initialDistribution := common.BaseDistribution{
		usd(): 1000,
		btc(): 2,
		eth(): 20,
	}
	initialPortfolio := convertToPortfolio(initialDistribution, usd())
	// Initialize the mock connector
	mockConnector := NewMockConnector(initialPortfolio)

	desiredWeights := map[data.Currency]toolkit.BaseWeight{
		usd(): 0.5,
		btc(): 0.3,
		eth(): 0.2,
	}

	// Hardcoded current prices
	currentPrices := map[data.Currency]float64{
		usd(): 1,
		btc(): 50000,
		eth(): 4000,
	}
	for currency, price := range currentPrices {
		mockConnector.SetPrice(*data.NewSymbolFromCurrencies(currency, usd()), price)
	}

	// Create new PortfolioWeights with desired weights
	portfolioWeights, err := toolkit.NewPortfolioWeights(desiredWeights)
	if err != nil {
		t.Fatalf("Error creating portfolio weights: %v", err)
	}

	// Calculate target distribution based on PortfolioWeights.Synchronize
	targetDistribution, err := portfolioWeights.Synchronize(initialDistribution, currentPrices, usd())
	if err != nil {
		t.Fatalf("Error synchronizing portfolio: %v", err)
	}

	// Calculate the difference between initial and target distributions
	rebalance := toolkit.NewRebalancePortfolio(*portfolioWeights)
	difference := mapSubtract(targetDistribution, initialDistribution)

	// Generate orders based on the difference and check them with MockConnector
	err = rebalance.Occur(mockConnector)
	if err != nil {
		t.Fatalf("Error occurred during rebalancing: %v", err)
	}

	for currency, expectedDiff := range difference {
		if currency == usd() {
			continue
		}
		expectedSymbol := data.NewSymbolFromCurrencies(currency, usd())
		expectedSide := exchange.Sell
		if expectedDiff > 0 {
			expectedSide = exchange.Buy
		}

		// Check if the order exists in MockConnector
		ok, order := mockConnector.CheckOrder(math.Abs(expectedDiff), *expectedSymbol, expectedSide, exchange.Market)
		if !ok || order == nil {
			t.Errorf("Expected: Currency: %s, Expected Side: %s, Expected Volume: %v, Actual Order: %v", currency.String(), expectedSide, math.Abs(expectedDiff), order)
		}
	}
}
