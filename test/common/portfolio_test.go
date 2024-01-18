package common_test

import (
	"testing"
	"strings"

	"xoney/common"
	"xoney/common/data"
	"xoney/errors"
)

func usd() data.Currency {
	return data.NewCurrency("USD", "EXCHANGE")
}

func usdPair(currency data.Currency) data.Symbol {
	return *data.NewSymbolFromCurrencies(currency, usd())
}

func portfolio() common.Portfolio {
	return common.NewPortfolio(usd())
}

func TestSet(t *testing.T) {
	USD := usd()
	p := portfolio()
	p.Set(USD, 100)

	if p.Balance(USD) != 100 {
		t.Error("Portfolio.Set() should impact portfolio.Balance")
	}
}

func TestCopy(t *testing.T) {
	USD := usd()
	p := portfolio()
	p.Set(USD, 100)

	pc := p.Copy()
	p.Set(USD, 50)

	if pc.Balance(USD) != 100 {
		t.Error("Portfolio.Copy() is not working properly")
	}
}

func TestPortfolioTotal(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Set some assets in the portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 0.5)
	portfolio.Set(data.NewCurrency("ETH", "EXCHANGE"), 2.0)

	// Create a map of prices for assets
	btc := data.NewCurrency("BTC", "EXCHANGE")
	eth := data.NewCurrency("ETH", "EXCHANGE")
	prices := map[data.Currency]float64{
		btc: 50000.0,
		eth: 2000.0,
	}

	// Calculate the total value of the portfolio
	result, err := portfolio.Total(prices)
	expectedTotal := 29000.0

	// Compare the actual and expected total value
	if result != expectedTotal {
		t.Errorf("Expected Total: %f, got: %f", expectedTotal, result)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Without necessary data
	delete(prices, btc)

	_, err = portfolio.Total(prices)

	// Compare the actual and expected error
	expected := errors.NewMissingCurrencyError(1)
	expected.Add(btc.String())

	if err.Error() != expected.Error() {
		t.Errorf("Expected %v, got: %v", expected, err)
	}
	if strings.Contains(err.Error(), ", ") {
		t.Errorf("Unexpected ', ' in error, got: %v", err)
	}

	// Without both currencies
	delete(prices, eth)

	_, err = portfolio.Total(prices)

	if !strings.Contains(err.Error(), ", ") {
		t.Errorf("Expected ', ' in error, got: %v", err)
	}
}

func TestPortfolioBalance(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Set some assets in the portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 0.5)
	portfolio.Set(data.NewCurrency("ETH", "EXCHANGE"), 2.0)

	// Check the balance of a specific currency
	expectedBalance := 2.0
	if result := portfolio.Balance(data.NewCurrency("ETH", "EXCHANGE")); result != expectedBalance {
		t.Errorf("Expected Balance: %f, got: %f", expectedBalance, result)
	}

	// Check the balance of a currency not in the portfolio
	expectedBalance = 0.0
	if result := portfolio.Balance(data.NewCurrency("LTC", "EXCHANGE")); result != expectedBalance {
		t.Errorf("Expected Balance: %f, got: %f", expectedBalance, result)
	}
}

func TestPortfolioAssets(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Set some assets in the portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 0.5)
	portfolio.Set(data.NewCurrency("ETH", "EXCHANGE"), 2.0)

	// Get the assets map from the portfolio
	expectedAssets := map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
		data.NewCurrency("ETH", "EXCHANGE"): 2.0,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}
}

func assetsMapEquals(a, b map[data.Currency]float64) bool {
	if len(a) != len(b) {
		return false
	}

	for currency, quantityA := range a {
		quantityB, ok := b[currency]
		if !ok || quantityA != quantityB {
			return false
		}
	}

	return true
}

func TestPortfolioSet(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Set an asset in the portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 0.5)

	// Check if the asset is set correctly
	expectedAssets := map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}

	// Set another asset in the portfolio
	portfolio.Set(data.NewCurrency("ETH", "EXCHANGE"), 2.0)

	// Check if both assets are set correctly
	expectedAssets = map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
		data.NewCurrency("ETH", "EXCHANGE"): 2.0,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}
}

func TestPortfolioIncrease(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Increase an asset in the portfolio
	portfolio.Increase(data.NewCurrency("BTC", "EXCHANGE"), 0.5)

	// Check if the asset is increased correctly
	expectedAssets := map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}

	// Increase another asset in the portfolio
	portfolio.Increase(data.NewCurrency("ETH", "EXCHANGE"), 2.0)

	// Check if both assets are increased correctly
	expectedAssets = map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
		data.NewCurrency("ETH", "EXCHANGE"): 2.0,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}
}

func TestPortfolioDecrease(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Set some assets in the portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 1.0)
	portfolio.Set(data.NewCurrency("ETH", "EXCHANGE"), 3.0)

	// Decrease an asset in the portfolio
	portfolio.Decrease(data.NewCurrency("BTC", "EXCHANGE"), 0.5)

	// Check if the asset is decreased correctly
	expectedAssets := map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
		data.NewCurrency("ETH", "EXCHANGE"): 3.0,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}

	// Decrease another asset in the portfolio
	portfolio.Decrease(data.NewCurrency("ETH", "EXCHANGE"), 1.0)

	// Check if both assets are decreased correctly
	expectedAssets = map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 0.5,
		data.NewCurrency("ETH", "EXCHANGE"): 2.0,
	}
	if result := portfolio.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}
}

func TestPortfolioMainCurrency(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Check if the main currency is as expected
	expectedMainCurrency := mainCurrency
	if result := portfolio.MainCurrency(); result != expectedMainCurrency {
		t.Errorf("Expected Main Currency: %v, got: %v", expectedMainCurrency, result)
	}
}

func TestPortfolioCopy(t *testing.T) {
	// Create a new Portfolio with a main currency
	mainCurrency := data.NewCurrency("USD", "EXCHANGE")
	portfolio := common.NewPortfolio(mainCurrency)

	// Set some assets in the portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 1.0)
	portfolio.Set(data.NewCurrency("ETH", "EXCHANGE"), 3.0)

	// Create a copy of the portfolio
	portfolioCopy := portfolio.Copy()

	// Check if the copy has the same assets
	expectedAssets := map[data.Currency]float64{
		data.NewCurrency("BTC", "EXCHANGE"): 1.0,
		data.NewCurrency("ETH", "EXCHANGE"): 3.0,
	}
	if result := portfolioCopy.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}

	// Modify the original portfolio
	portfolio.Set(data.NewCurrency("BTC", "EXCHANGE"), 2.0)

	// Check if the copy remains unchanged
	if result := portfolioCopy.Assets(); !assetsMapEquals(result, expectedAssets) {
		t.Errorf("Expected Assets: %v, got: %v", expectedAssets, result)
	}
}
