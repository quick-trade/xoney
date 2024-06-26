package exchange_test

import (
	goErrors "errors"
	"testing"
	"time"

	"github.com/quick-trade/xoney/common"
	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/exchange"
)

func usd() data.Currency {
	return data.NewCurrency("USD", "BINANCE")
}

func btc() data.Currency {
	return data.NewCurrency("BTC", "BINANCE")
}

func timeStart() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

func btcUSD() data.Symbol {
	return *data.NewSymbol("BTC", "USD", "BINANCE")
}

func timeframe() data.TimeFrame {
	timeframe, _ := data.NewTimeFrame(time.Hour, "1h")

	return *timeframe
}

func instrument() data.Instrument {
	return data.NewInstrument(btcUSD(), timeframe())
}

func portfolioUSD() common.Portfolio {
	portfolio := common.NewPortfolio(usd())
	portfolio.Set(usd(), 5000)

	return portfolio
}

func marginSimulator() exchange.MarginSimulator {
	return exchange.NewMarginSimulator(portfolioUSD(), 0)
}

func spotSimulator() exchange.SpotSimulator {
	return exchange.NewSpotSimulator(portfolioUSD(), 0)
}

func TestMarginSimulator_PlaceMarketOrder_Buy(t *testing.T) {
	simulator := marginSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 0.1
	order, _ := exchange.NewOrder(symbol, exchange.Market, exchange.Buy, price, amount)

	err := simulator.PlaceOrder(*order)
	if err != nil {
		t.Errorf("Error placing market buy order: %v", err)
	}

	expectedBalance := 5000.0 - amount*price
	if simulator.Portfolio().Balance(usd()) != expectedBalance {
		t.Errorf("Expected balance after market buy order: %v, got: %v", expectedBalance, simulator.Portfolio().Balance(usd()))
	}
}

func TestMarginSimulator_PlaceMarketOrder_Sell(t *testing.T) {
	simulator := marginSimulator()
	symbol := btcUSD()
	btc := symbol.Base()
	price := 50000.0
	amount := 0.1
	order, _ := exchange.NewOrder(symbol, exchange.Market, exchange.Sell, price, amount)

	err := simulator.PlaceOrder(*order)
	if err != nil {
		t.Errorf("Error placing market sell order: %v", err)
	}

	expectedBalance := 5000.0 + amount*price
	expectedBTC := -amount
	if simulator.Portfolio().Balance(usd()) != expectedBalance {
		t.Errorf("Expected balance after market sell order: %v, got: %v", expectedBalance, simulator.Portfolio().Balance(usd()))
	}
	if simulator.Portfolio().Balance(btc) != expectedBTC {
		t.Errorf("Expected balance after market sell order: %v, got: %v", expectedBalance, simulator.Portfolio().Balance(usd()))
	}
}

func TestMarginSimulator_ExecuteLimitOrder_ImmediateExecution(t *testing.T) {
	simulator := marginSimulator()
	symbol := data.NewSymbol("BTC", "USD", "BINANCE")
	price := 50000.0
	amount := 0.1

	// Place a limit order that should be executed immediately
	limitOrder, _ := exchange.NewOrder(*symbol, exchange.Limit, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*limitOrder)
	if err != nil {
		t.Errorf("Error placing limit order: %v", err)
	}

	// Update price to cross the limit order immediately
	x := 100.0
	candle := data.NewCandle(price-1000, price+200, price-1200, price-x, 0, timeStart())
	iCandle := data.NewInstrumentCandle(*candle, instrument())
	err = simulator.UpdatePrice(*iCandle)

	if err != nil {
		t.Errorf("Error updating price: %v", err)
	}

	// Check if the limit order was executed
	expectedBalance := 5000.0 - amount*price

	balance := simulator.Portfolio().Balance(usd())
	if balance != expectedBalance {
		t.Errorf("Expected balance after immediate limit order execution: %v, got: %v", expectedBalance, balance)
	}

	balance = simulator.Portfolio().Balance(btc())
	if balance != amount {
		t.Errorf("Expected balance after immediate limit order: %v, got: %v", amount, balance)
	}
}

func TestMarginSimulator_ExecuteLimitOrder_DelayedExecution(t *testing.T) {
	simulator := marginSimulator()
	symbol := data.NewSymbol("BTC", "USD", "BINANCE")
	price := 50000.0
	amount := 0.1

	// Place a limit order that should be executed after two updates
	limitOrder, _ := exchange.NewOrder(*symbol, exchange.Limit, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*limitOrder)
	if err != nil {
		t.Errorf("Error placing limit order: %v", err)
	}

	// Update price to cross the limit order after two updates
	candle1 := data.NewCandle(price+1, price+124, price-50, price-1000, 0, timeStart())
	iCandle1 := data.NewInstrumentCandle(*candle1, instrument())
	err = simulator.UpdatePrice(*iCandle1)

	if err != nil {
		t.Errorf("Error updating price: %v", err)
	}

	candle2 := data.NewCandle(price-1000, price-600, price-1350, price-900, 0, timeStart().Add(time.Hour))
	iCandle2 := data.NewInstrumentCandle(*candle2, instrument())
	err = simulator.UpdatePrice(*iCandle2)

	if err != nil {
		t.Errorf("Error updating price: %v", err)
	}

	// Check if the limit order was executed
	expectedBalance := 5000.0 - amount*price

	balance := simulator.Portfolio().Balance(usd())
	if balance != expectedBalance {
		t.Errorf("Expected balance after delayed limit order execution: %v, got: %v", expectedBalance, balance)
	}

	balance = simulator.Portfolio().Balance(btc())
	if balance != amount {
		t.Errorf("Expected balance after immediate limit order: %v, got: %v", amount, balance)
	}
}

func TestMarginSimulator_ExecuteMultipleLimitOrders(t *testing.T) {
	simulator := marginSimulator()
	symbol := data.NewSymbol("BTC", "USD", "BINANCE")

	price1 := 50000.0
	price2 := 51000.0

	amount1 := 0.1
	amount2 := 0.2

	// Place two limit orders
	limitOrder1, _ := exchange.NewOrder(*symbol, exchange.Limit, exchange.Buy, price1, amount1)
	err := simulator.PlaceOrder(*limitOrder1)
	if err != nil {
		t.Errorf("Error placing limit order 1: %v", err)
	}

	limitOrder2, _ := exchange.NewOrder(*symbol, exchange.Limit, exchange.Buy, price2, amount2)
	err = simulator.PlaceOrder(*limitOrder2)

	if err != nil {
		t.Errorf("Error placing limit order 2: %v", err)
	}

	candle := data.NewCandle(price1+1200, price1+4000, price1+1001, price1+1111, 0, timeStart())
	iCandle := data.NewInstrumentCandle(*candle, instrument())
	err = simulator.UpdatePrice(*iCandle)

	if err != nil {
		t.Errorf("Error updating price: %v", err)
	}

	assets := simulator.Portfolio().Assets()

	if assets[btc()] != 0 {
		t.Errorf("Unexpected order execution, balance: %fBTC", assets[btc()])
	}
	if assets[usd()] != 5000 {
		t.Errorf("Unexpected order execution, balance: %fUSD", assets[usd()])
	}

	candle = data.NewCandle(price2+200, price2+3000, price2, price1+111, 0, timeStart())
	iCandle = data.NewInstrumentCandle(*candle, instrument())
	err = simulator.UpdatePrice(*iCandle)

	if err != nil {
		t.Errorf("Error updating price: %v", err)
	}

	assets = simulator.Portfolio().Assets()
	if assets[btc()] != amount2 {
		t.Errorf("Incorrect order execution, balance: %fBTC, expected: %fBTC", assets[btc()], amount2)
	}
}

func TestMarginSimulator_CancelOrder_NonExistingOrder(t *testing.T) {
	simulator := marginSimulator()

	// Try to cancel a non-existing order
	nonExistingOrderID := exchange.OrderID(123)
	err := simulator.CancelOrder(nonExistingOrderID)

	// Check if the error is of the expected type
	expectedError := errors.NewNoLimitOrderError(123)
	if !goErrors.Is(err, expectedError) {
		t.Errorf("Expected NoLimitOrderError, got: %v", err)
	}
	if err.Error() != "there is no such limit order with ID: 123." {
		t.Errorf("Incorrect text of error: %v", err)
	}
}

func TestMarginSimulator_CancelOrder_ExistingOrder(t *testing.T) {
	simulator := marginSimulator()
	symbol := data.NewSymbol("BTC", "USD", "BINANCE")
	price := 50000.0
	amount := 0.1

	// Place a limit order to have an existing order
	limitOrder, _ := exchange.NewOrder(*symbol, exchange.Limit, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*limitOrder)
	if err != nil {
		t.Fatalf("Error placing limit order: %v", err)
	}

	// Get the ID of the placed order
	orderID := limitOrder.ID()

	// Cancel the existing order
	err = simulator.CancelOrder(orderID)

	// Check if there is no error
	if err != nil {
		t.Errorf("Error cancelling existing order: %v", err)
	}

	// Check if the order is removed from the portfolio
	if simulator.CancelOrder(orderID) == nil {
		t.Error("Expected 0 open orders after cancellation")
	}
}

func TestMarginSimulator_Transfer_SuccessfulTransfer(t *testing.T) {
	simulator := marginSimulator()
	initialBalanceUSD := simulator.Portfolio().Balance(usd())
	initialBalanceBTC := simulator.Portfolio().Balance(btc())

	// Transfer 1000 USD to another exchange
	transferAmount := 1000.0
	err := simulator.Transfer(transferAmount, usd(), data.Exchange("OtherExchange"))
	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error during successful transfer: %v", err)
	}

	// Check if the balance of USD decreased by the transfer amount
	expectedBalanceUSD := initialBalanceUSD - transferAmount
	if simulator.Portfolio().Balance(usd()) != expectedBalanceUSD {
		t.Errorf("Expected USD balance after successful transfer: %v, got: %v", expectedBalanceUSD, simulator.Portfolio().Balance(usd()))
	}

	// Check if the balance of BTC remained unchanged
	if simulator.Portfolio().Balance(btc()) != initialBalanceBTC {
		t.Errorf("Expected BTC balance to remain unchanged after successful transfer: %v, got: %v", initialBalanceBTC, simulator.Portfolio().Balance(btc()))
	}
}

func TestMarginSimulator_Transfer_InsufficientFunds(t *testing.T) {
	simulator := marginSimulator()
	initialBalanceUSD := simulator.Portfolio().Balance(usd())
	transferAmount := initialBalanceUSD + 100.0

	// Attempt to transfer an amount exceeding the available funds
	err := simulator.Transfer(transferAmount, usd(), data.Exchange("OtherExchange"))

	// Check if the error is of the expected type
	expectedError := errors.NewNotEnoughFundsError(usd().String(), transferAmount)
	if !goErrors.Is(err, expectedError) {
		t.Errorf("Expected NotEnoughFundsError, got: %v", err)
	}

	// Check if the balances remain unchanged
	if simulator.Portfolio().Balance(usd()) != initialBalanceUSD {
		t.Errorf("Expected USD balance to remain unchanged after insufficient funds transfer: %v, got: %v", initialBalanceUSD, simulator.Portfolio().Balance(usd()))
	}
}

func TestMarginSimulator_Transfer_TargetExchangeUpdate(t *testing.T) {
	simulator := marginSimulator()
	initialBalanceUSD := simulator.Portfolio().Balance(usd())
	initialBalanceBTC := simulator.Portfolio().Balance(btc())

	// Transfer 1000 USD to another exchange
	transferAmount := 1000.0
	targetExchange := data.Exchange("OtherExchange")
	err := simulator.Transfer(transferAmount, usd(), targetExchange)
	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error during successful transfer: %v", err)
	}

	// Check if the balance of USD decreased by the transfer amount
	expectedBalanceUSD := initialBalanceUSD - transferAmount
	if simulator.Portfolio().Balance(usd()) != expectedBalanceUSD {
		t.Errorf("Expected USD balance after successful transfer: %v, got: %v", expectedBalanceUSD, simulator.Portfolio().Balance(usd()))
	}

	// Check if the balance of BTC remained unchanged
	if simulator.Portfolio().Balance(btc()) != initialBalanceBTC {
		t.Errorf("Expected BTC balance to remain unchanged after successful transfer: %v, got: %v", initialBalanceBTC, simulator.Portfolio().Balance(btc()))
	}

	// Check if the currency exchange was updated
	if simulator.Portfolio().Balance(data.NewCurrency("USD", targetExchange)) != transferAmount {
		t.Errorf("Expected updated balance in the target exchange after successful transfer: %v, got: %v", transferAmount, simulator.Portfolio().Balance(data.NewCurrency("USD", targetExchange)))
	}
}

func TestMarginSimulator_SellAll_SuccessfulSell(t *testing.T) {
	simulator := marginSimulator()
	initialBalanceUSD := simulator.Portfolio().Balance(usd())
	initialBalanceBTC := simulator.Portfolio().Balance(btc())

	// Set prices for BTC and ETH
	candleBTC := data.NewCandle(50000.0, 51000.0, 49000.0, 50500.0, 0, timeStart())
	simulator.UpdatePrice(*data.NewInstrumentCandle(*candleBTC, instrument()))

	candleETH := data.NewCandle(3000.0, 3100.0, 2900.0, 3050.0, 0, timeStart())
	simulator.UpdatePrice(*data.NewInstrumentCandle(*candleETH, instrument()))

	// Execute SellAll
	err := simulator.SellAll()
	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error during successful SellAll: %v", err)
	}

	// Check if the balance of USD increased by selling BTC
	expectedBalanceUSD := initialBalanceUSD + initialBalanceBTC*50500.0
	if simulator.Portfolio().Balance(usd()) != expectedBalanceUSD {
		t.Errorf("Expected USD balance after successful SellAll: %v, got: %v", expectedBalanceUSD, simulator.Portfolio().Balance(usd()))
	}

	// Check if the balance of BTC became zero
	if simulator.Portfolio().Balance(btc()) != 0 {
		t.Errorf("Expected BTC balance to be zero after successful SellAll, got: %v", simulator.Portfolio().Balance(btc()))
	}

	// Check if the balance of ETH remains unchanged
	if simulator.Portfolio().Balance(data.NewCurrency("ETH", "BINANCE")) != 0 {
		t.Errorf("Expected ETH balance to remain unchanged after successful SellAll, got: %v", simulator.Portfolio().Balance(data.NewCurrency("ETH", "BINANCE")))
	}
}

func TestMarginSimulator_SellAll_ErrorPlacingOrder(t *testing.T) {
	simulator := marginSimulator()
	initialBalanceUSD := simulator.Portfolio().Balance(usd())
	initialBalanceBTC := simulator.Portfolio().Balance(btc())

	// Set prices for BTC and ETH
	candleBTC := data.NewCandle(50000.0, 51000.0, 49000.0, 50500.0, 0, timeStart())
	simulator.UpdatePrice(*data.NewInstrumentCandle(*candleBTC, instrument()))

	// Execute SellAll
	err := simulator.SellAll()
	if err != nil {
		t.Errorf("Unexpected error, got: %v", err)
	}

	// Check if the balances remain unchanged
	if simulator.Portfolio().Balance(usd()) != initialBalanceUSD {
		t.Errorf("Expected USD balance to remain unchanged after error placing order: %v, got: %v", initialBalanceUSD, simulator.Portfolio().Balance(usd()))
	}
	if simulator.Portfolio().Balance(btc()) != initialBalanceBTC {
		t.Errorf("Expected BTC balance to remain unchanged after error placing order: %v, got: %v", initialBalanceBTC, simulator.Portfolio().Balance(btc()))
	}
}

func TestMarginSimulator_SellAll_NoOrdersPlaced(t *testing.T) {
	simulator := marginSimulator()
	initialBalanceUSD := simulator.Portfolio().Balance(usd())
	initialBalanceBTC := simulator.Portfolio().Balance(btc())

	// Execute SellAll when there are no prices set
	err := simulator.SellAll()
	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error during SellAll with no prices set: %v", err)
	}

	// Check if the balances remain unchanged
	if simulator.Portfolio().Balance(usd()) != initialBalanceUSD {
		t.Errorf("Expected USD balance to remain unchanged after SellAll with no prices set: %v, got: %v", initialBalanceUSD, simulator.Portfolio().Balance(usd()))
	}
	if simulator.Portfolio().Balance(btc()) != initialBalanceBTC {
		t.Errorf("Expected BTC balance to remain unchanged after SellAll with no prices set: %v, got: %v", initialBalanceBTC, simulator.Portfolio().Balance(btc()))
	}
}

func TestSpotSimulator_ExecuteMarketOrder_Buy(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 0.1

	// Execute a market buy order
	marketOrder, _ := exchange.NewOrder(symbol, exchange.Market, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*marketOrder)
	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error during market buy order execution: %v", err)
	}

	// Check if the balance of USD decreased by the order amount * price
	expectedBalanceUSD := 5000.0 - amount*price
	if simulator.Portfolio().Balance(usd()) != expectedBalanceUSD {
		t.Errorf("Expected USD balance after market buy order: %v, got: %v", expectedBalanceUSD, simulator.Portfolio().Balance(usd()))
	}

	// Check if the balance of BTC increased by the order amount
	expectedBalanceBTC := amount
	if simulator.Portfolio().Balance(btc()) != expectedBalanceBTC {
		t.Errorf("Expected BTC balance after market buy order: %v, got: %v", expectedBalanceBTC, simulator.Portfolio().Balance(btc()))
	}
}

func TestSpotSimulator_PlaceOrder_MarketBuy_Success(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 0.1
	order, _ := exchange.NewOrder(symbol, exchange.Market, exchange.Buy, price, amount)

	err := simulator.PlaceOrder(*order)
	if err != nil {
		t.Errorf("Unexpected error during successful market buy order placement: %v", err)
	}

	// Check if the balance of USD decreased by the quote quantity
	expectedBalanceUSD := 5000.0 - amount*price
	if simulator.Portfolio().Balance(usd()) != expectedBalanceUSD {
		t.Errorf("Expected USD balance after successful market buy order placement: %v, got: %v", expectedBalanceUSD, simulator.Portfolio().Balance(usd()))
	}

	// Check if the balance of BTC increased by the base quantity
	expectedBalanceBTC := amount
	if simulator.Portfolio().Balance(btc()) != expectedBalanceBTC {
		t.Errorf("Expected BTC balance after successful market buy order placement: %v, got: %v", expectedBalanceBTC, simulator.Portfolio().Balance(btc()))
	}
}

func TestSpotSimulator_ExecuteMarketOrder_Buy_InsufficientFunds(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	amount := 5.0
	price := 50000.0

	// Execute market buy order with insufficient funds
	order, _ := exchange.NewOrder(symbol, exchange.Market, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*order)
	if err == nil || !goErrors.Is(err, errors.NewNotEnoughFundsError(usd().String(), amount*price)) {
		t.Errorf("Expected NotEnoughFundsError, got: %v", err)
	}

	// Check if the balances remain unchanged
	if simulator.Portfolio().Balance(usd()) != 5000.0 {
		t.Errorf("Expected USD balance to remain unchanged after insufficient funds: %v, got: %v", 5000.0, simulator.Portfolio().Balance(usd()))
	}
	if simulator.Portfolio().Balance(btc()) != 0.0 {
		t.Errorf("Expected BTC balance to remain unchanged after insufficient funds: %v, got: %v", 0.0, simulator.Portfolio().Balance(btc()))
	}
}

func TestSpotSimulator_ExecuteLimitOrder_Sell_InsufficientFunds(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 0.1

	// Execute limit sell order
	order, _ := exchange.NewOrder(symbol, exchange.Limit, exchange.Sell, price, amount)
	err := simulator.PlaceOrder(*order)
	if !goErrors.Is(err, errors.NewNotEnoughFundsError(btc().String(), amount)) {
		t.Errorf("Expected NotEnoughFundsError, got: %v", err)
	}

	// Check if the balance of BTC increased by the order amount
	expectedBalanceBTC := amount
	if simulator.Portfolio().Balance(btc()) != 0 {
		t.Errorf("Expected BTC balance after limit sell order: %v, got: %v", expectedBalanceBTC, simulator.Portfolio().Balance(btc()))
	}

	// Check if the balance of USD remains unchanged
	if simulator.Portfolio().Balance(usd()) != 5000.0 {
		t.Errorf("Expected USD balance to remain unchanged after limit sell order: %v, got: %v", 5000.0, simulator.Portfolio().Balance(usd()))
	}
}

func TestSpotSimulator_PlaceLimitOrder_Sell_InsufficientFunds(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 0.1

	// Execute limit sell order with insufficient funds
	order, _ := exchange.NewOrder(symbol, exchange.Limit, exchange.Sell, price, amount)
	err := simulator.PlaceOrder(*order)
	if err == nil || !goErrors.Is(err, errors.NewNotEnoughFundsError(btc().String(), amount)) {
		t.Errorf("Expected NotEnoughFundsError, got: %v", err)
	}

	// Check if the balances remain unchanged
	if simulator.Portfolio().Balance(btc()) != 0.0 {
		t.Errorf("Expected BTC balance to remain unchanged after insufficient funds: %v, got: %v", 0.0, simulator.Portfolio().Balance(btc()))
	}

	if simulator.Portfolio().Balance(usd()) != 5000.0 {
		t.Errorf("Expected USD balance to remain unchanged after insufficient funds: %v, got: %v", 5000.0, simulator.Portfolio().Balance(usd()))
	}
}

func TestSpotSimulator_PlaceLimitOrder_Buy_Success(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 0.1

	// Execute limit buy order
	order, _ := exchange.NewOrder(symbol, exchange.Limit, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*order)
	if err != nil {
		t.Errorf("Unexpected error during successful limit buy order placement: %v", err)
	}

	// Set prices for BTC
	candleBTC := data.NewCandle(price, price, price, price, 0, timeStart())
	simulator.UpdatePrice(*data.NewInstrumentCandle(*candleBTC, instrument()))

	// Check if the balance of USD decreased by the quote quantity
	expectedBalanceUSD := 5000.0 - amount*price
	if simulator.Portfolio().Balance(usd()) != expectedBalanceUSD {
		t.Errorf("Expected USD balance after successful limit buy order placement: %v, got: %v", expectedBalanceUSD, simulator.Portfolio().Balance(usd()))
	}

	// Check if the balance of BTC increased by the base quantity
	expectedBalanceBTC := amount
	if simulator.Portfolio().Balance(btc()) != expectedBalanceBTC {
		t.Errorf("Expected BTC balance after successful limit buy order placement: %v, got: %v", expectedBalanceBTC, simulator.Portfolio().Balance(btc()))
	}
}

func TestSpotSimulator_PlaceLimitOrder_Buy_InsufficientFunds(t *testing.T) {
	simulator := spotSimulator()
	symbol := btcUSD()
	price := 50000.0
	amount := 5.0

	// Set prices for BTC
	candleBTC := data.NewCandle(price, price, price, price, 0, timeStart())
	simulator.UpdatePrice(*data.NewInstrumentCandle(*candleBTC, instrument()))

	// Execute limit buy order with insufficient funds
	order, _ := exchange.NewOrder(symbol, exchange.Limit, exchange.Buy, price, amount)
	err := simulator.PlaceOrder(*order)
	if err == nil || !goErrors.Is(err, errors.NewNotEnoughFundsError(usd().String(), amount*price)) {
		t.Errorf("Expected NotEnoughFundsError, got: %v", err)
	}

	// Check if the balances remain unchanged
	if simulator.Portfolio().Balance(usd()) != 5000.0 {
		t.Errorf("Expected USD balance to remain unchanged after insufficient funds: %v, got: %v", 5000.0, simulator.Portfolio().Balance(usd()))
	}

	if simulator.Portfolio().Balance(btc()) != 0.0 {
		t.Errorf("Expected BTC balance to remain unchanged after insufficient funds: %v, got: %v", 0.0, simulator.Portfolio().Balance(btc()))
	}
}
