package exchange_test

import (
	"testing"

	"xoney/common/data"
	"xoney/exchange"
)

func eth() data.Symbol {
	return *data.NewSymbol("ETH", "USD", "BINANCE")
}

func orderBTC() *exchange.Order {
	order, _ := exchange.NewOrder(
		btcUSD(),
		exchange.Market,
		exchange.Buy,
		50000.0,
		1.0,
	)
	return order
}

func orderBTCLimit() *exchange.Order {
	order, _ := exchange.NewOrder(
		btcUSD(),
		exchange.Limit,
		exchange.Buy,
		50000.0,
		1.0,
	)
	return order
}

func orderBTC4k() *exchange.Order {
	order, _ := exchange.NewOrder(
		btcUSD(),
		exchange.Market,
		exchange.Buy,
		4000.0,
		1.0,
	)
	return order
}

func orderETH() *exchange.Order {
	order, _ := exchange.NewOrder(
		eth(),
		exchange.Limit,
		exchange.Sell,
		200.0,
		5.0,
	)
	return order
}

func orderETHBuy() *exchange.Order {
	order, _ := exchange.NewOrder(
		eth(),
		exchange.Limit,
		exchange.Buy,
		200.0,
		5.0,
	)
	return order
}

func orderETHlikeBTC() *exchange.Order {
	order, _ := exchange.NewOrder(
		eth(),
		exchange.Market,
		exchange.Buy,
		50000.0,
		1.0,
	)
	return order
}

func orderETH1() *exchange.Order {
	order, _ := exchange.NewOrder(
		eth(),
		exchange.Limit,
		exchange.Sell,
		200.0,
		1.0,
	)
	return order
}

func TestOrderMethods(t *testing.T) {
	order1 := orderBTC()
	order2 := orderETH()

	if result := order1.Symbol(); result != btcUSD() {
		t.Errorf("Expected Symbol: BTCUSD, got: %s", result)
	}

	if result := order2.Type(); result != exchange.Limit {
		t.Errorf("Expected Type: Limit, got: %s", result)
	}

	if result := order1.Side(); result != exchange.Buy {
		t.Errorf("Expected Side: Buy, got: %s", result)
	}

	if result := order1.Price(); result != 50000.0 {
		t.Errorf("Expected Price: 50000.0, got: %f", result)
	}

	if result := order2.Amount(); result != 5.0 {
		t.Errorf("Expected Amount: 5.0, got: %f", result)
	}

	if !order1.IsEqual(order1) {
		t.Errorf("Expected order1 to be equal to itself")
	}

	if !order1.CrossesPrice(55000.0, 49000.0) {
		t.Errorf("Expected CrossesPrice to be true for a Buy order")
	}

	if order2.CrossesPrice(199.0, 170.0) {
		t.Errorf("Expected CrossesPrice to be false for a Sell order")
	}

	if order1.ID() == order2.ID() {
		t.Errorf("Expected different IDs for different orders: %v and %v", order1.ID(), order2.ID())
	}
}

func TestIsEqual(t *testing.T) {
	btc1 := orderBTC()
	btc2 := orderBTCLimit()
	eth1 := orderETH()
	eth2 := orderETHBuy()
	ethLikeBTC := orderETHlikeBTC()
	btc4k := orderBTC4k()
	eth1eth := orderETH1()

	if btc1.IsEqual(btc2) {
		t.Error("Expected .IsEqual to be false for limit and market orders")
	}

	if eth1.IsEqual(eth2) {
		t.Error("Expected .IsEqual to be false for buy and sell orders")
	}

	if btc1.IsEqual(ethLikeBTC) {
		t.Error("Expected .IsEqual to be false for BTC and ETH orders")
	}

	if btc1.IsEqual(btc4k) {
		t.Error("Expected .IsEqual to be false for orders with different prices")
	}

	if eth1.IsEqual(eth1eth) {
		t.Error("Expected .IsEqual to be false for orders with different amount")
	}
}
