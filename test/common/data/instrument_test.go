package data_test

import (
	"testing"
	"time"

	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/errors"
)

func TestSymbolString(t *testing.T) {
	symbol := data.NewSymbol("BTC", "USD", "BINANCE")

	// Expected string representation
	expectedString := "BINANCE:BTC/USD"

	// Compare the actual and expected string representations
	if result := symbol.String(); result != expectedString {
		t.Errorf("Expected: %s, got: %s", expectedString, result)
	}
}

func TestSymbolExchange(t *testing.T) {
	symbol := data.NewSymbol("BTC", "USD", "BINANCE")

	// Expected exchange
	expectedExchange := data.Exchange("BINANCE")

	// Compare the actual and expected exchange
	if result := symbol.Exchange(); result != expectedExchange {
		t.Errorf("Expected: %s, got: %s", expectedExchange, result)
	}
}

func TestTimeFrameIncorrect(t *testing.T) {
	duration := time.Duration(-time.Hour)
	_, err := data.NewTimeFrame(duration, "incorrect tf")

	expected := "invalid duration: -1h0m0s."

	if err.Error() != expected {
		t.Errorf("Expected IncorrectDurationError, got: %s", err.Error())
	}
	if _, ok := err.(errors.IncorrectDurationError); !ok {
		t.Errorf("Expected IncorrectDurationError, got: %v", err)
	}
}

func Daily() data.TimeFrame {
	tf, err := data.NewTimeFrame(time.Hour*24, "1D")
	if err != nil {
		panic("Error creating TimeFrame")
	}

	return *tf
}

func TestEquityTimeframe(t *testing.T) {
	// Create a new Equity instance with a specific timeframe
	equity := data.NewEquity(Daily(), 100)

	// Expected timeframe
	expectedTimeframe := Daily()

	// Compare the actual and expected timeframe
	if result := equity.Timeframe(); result != expectedTimeframe {
		t.Errorf("Expected: %s, got: %s", expectedTimeframe.Name, result.Name)
	}
}

func TestEquityDeposit(t *testing.T) {
	// Create a new Equity instance
	equity := data.NewEquity(Daily(), 100)

	// Deposit some values
	values := []float64{1000.0, 1500.0, 1200.0}
	for _, value := range values {
		equity.AddValue(value, time.Now())
	}

	// Expected deposit history
	expectedDeposit := values

	// Compare the actual and expected deposit history
	if result := equity.Deposit(); !isEqual(result, expectedDeposit) {
		t.Errorf("Expected: %v, got: %v", expectedDeposit, result)
	}
}

func TestEquityPortfolioHistory(t *testing.T) {
	equity := data.NewEquity(Daily(), 100)

	usd := data.NewCurrency("USD", "BINANCE")
	btc := data.NewCurrency("BTC", "BINANCE")

	// Add portfolio values
	portfolio1 := map[data.Currency]float64{usd: 1000.0, btc: 0.5}
	equity.AddPortfolio(portfolio1)

	portfolio2 := map[data.Currency]float64{usd: 1200.0, btc: 0.7}
	equity.AddPortfolio(portfolio2)

	// Expected portfolio history
	expectedPortfolio := map[data.Currency][]float64{
		usd: {1000.0, 1200.0},
		btc: {0.5, 0.7},
	}

	// Compare the actual and expected portfolio history
	if result := equity.PortfolioHistory(); !isEqualMap(result, expectedPortfolio) {
		t.Errorf("Expected: %v, got: %v", expectedPortfolio, result)
	}
}

// Helper function to check if two slices are equal.
func isEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Helper function to check if two maps of slices are equal.
func isEqualMap(a, b map[data.Currency][]float64) bool {
	if len(a) != len(b) {
		return false
	}

	for key := range a {
		if !isEqual(a[key], b[key]) {
			return false
		}
	}

	return true
}

func TestEquityNow(t *testing.T) {
	// Create a new Equity instance
	equity := data.NewEquity(Daily(), 100)

	// Add some values to the equity
	values := []float64{1000.0, 1500.0, 1200.0}
	for _, value := range values {
		equity.AddValue(value, time.Now())
	}

	// Expected value from the Now method
	expectedNow := values[len(values)-1]

	// Compare the actual and expected Now value
	if result := equity.Now(); result != expectedNow {
		t.Errorf("Expected: %f, got: %f", expectedNow, result)
	}
}
