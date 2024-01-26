// toolkit_test.go
package toolkit_test

import (
	"math"
	"testing"
	"xoney/common"
	"xoney/common/data"
	"xoney/toolkit"
)

const epsilon = 0.001

// closeEnough возвращает true, если a и b близки с учетом погрешности epsilon.
func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// equalFloatSlices сравнивает два среза чисел с учетом погрешности epsilon.
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
	weights := make(toolkit.PortfolioWeights)

	for currency, amount := range portfolio {
		weights[currency] = toolkit.BaseWeight(amount / totalCapital)
	}

	dist, err := toolkit.NewPortfolioWeights(weights)
	if err != nil {
		panic(err)
	}
	return *dist
}

func TestSynchronize(t *testing.T) {
	currentBasePortfolio := common.BaseDistribution{
		data.NewCurrency("USD", "EXCHANGE"): 1. / 100.,
		data.NewCurrency("BTC", "EXCHANGE"): 1. / 500.,
		data.NewCurrency("ETH", "EXCHANGE"): 1. / 10.,
	}

	targetWeights := CurrentWeights(currentBasePortfolio)

	currentPrices := map[data.Currency]float64{
		data.NewCurrency("USD", "EXCHANGE"): 100,
		data.NewCurrency("BTC", "EXCHANGE"): 500,
		data.NewCurrency("ETH", "EXCHANGE"): 10,
	}

	expectedResult := currentBasePortfolio

	result, err := targetWeights.Synchronize(currentBasePortfolio, currentPrices)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !equalFloatMaps(result, expectedResult) {
		t.Errorf("Expected %v, got %v", expectedResult, result)
	}
}
