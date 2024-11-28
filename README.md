<div align="center">
  <img src="assets/logo.png" width="340" height="340">

# Xoney
</div>

A simple, fast, and powerful library for algorithmic trading in Go, **with no dependencies**.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Strategies](#strategies)
  - [Bollinger Bands Strategy](#bollinger-bands-strategy)
  - [Grid Trading Bot](#grid-trading-bot)
  - [Custom Strategy](#custom-strategy)
- [Backtesting](#backtesting)
- [Portfolio Management](#portfolio-management)
- [Exchange Simulation](#exchange-simulation)

## Features
- Zero external dependencies
- Event-driven architecture
- Built-in strategy implementations
- Backtesting engine with market and limit orders support
- Portfolio management tools
- Exchange simulation for testing
- Support for both spot and margin trading
- Customizable commission rates
- Performance metrics calculation (Sharpe ratio, CARA utility)

## Installation

```bash
go get github.com/quick-trade/xoney
```

## Quick Start

```go
package main

import (
    "github.com/quick-trade/xoney/backtest"
    "github.com/quick-trade/xoney/common"
    "github.com/quick-trade/xoney/exchange"
)

func main() {
    // Create a portfolio with initial USD balance
    currency := data.NewCurrency("USD", "BINANCE")
    portfolio := common.NewPortfolio(currency)
    portfolio.Set(currency, 10000)

    // Initialize exchange simulator with 0.1% commission
    simulator := exchange.NewMarginSimulator(portfolio, 0.001)
    
    // Create backtester
    tester := backtest.NewBacktester(&simulator)
    
    // Run your strategy
    equity, err := tester.Backtest(charts, strategy)
    if err != nil {
        panic(err)
    }
}
```

## Strategies

### Bollinger Bands Strategy

The Bollinger Bands strategy is a popular technical analysis tool that uses standard deviations to determine overbought or oversold conditions.

```go
package main

import (
    "github.com/quick-trade/xoney/strategy"
    "github.com/quick-trade/xoney/common/data"
)

func main() {
    // Create BTC/USD instrument with 15-minute timeframe
    btcUsd := data.NewSymbol("BTC", "USD", "BINANCE")
    timeframe, _ := data.NewTimeFrame(time.Minute*15, "15m")
    instrument := data.NewInstrument(*btcUsd, *timeframe)
    
    // Initialize Bollinger Bands strategy
    // Parameters: period=30 days, deviation=2
    strategy := NewBBStrategy(60*24*30, 2, instrument)
    
    // Run backtest
    equity, err := tester.Backtest(charts, &strategy)
    if err != nil {
        panic(err)
    }
}
```

### Grid Trading Bot

Grid trading is a strategy that places multiple orders at regular intervals above and below a set price, aiming to profit from market oscillations.

```go
package main

import (
    "github.com/quick-trade/xoney/toolkit"
    "github.com/quick-trade/xoney/common/data"
)

func main() {
    // Create auto grid generator
    // Parameters: lookback=100 candles, levels=10, deviation=1.5, total_amount=0.5 (in base currency)
    generator := NewAutoGrid(100, 10, 1.5, 0.5)
    
    // Initialize grid bot
    bot := toolkit.NewGridBot(generator, instrument)
    
    // Run backtest
    equity, err := tester.Backtest(charts, bot)
    if err != nil {
        panic(err)
    }
}
```

### Custom Strategy

You can create your own strategy by implementing the `Tradable` interface:

```go
type Tradable interface {
    Start(charts data.ChartContainer) error
    Next(candle data.InstrumentCandle) (events.Event, error)
    MinDurations() Durations
}
```

Example of a simple moving average crossover strategy:

```go
type MACrossStrategy struct {
    instrument data.Instrument
    shortPeriod int
    longPeriod int
    shortMA []float64
    longMA []float64
}

func (m *MACrossStrategy) Next(candle data.InstrumentCandle) (events.Event, error) {
    // Strategy logic here
    return event, nil
}
```

## Backtesting

Xoney provides a powerful backtesting engine that supports:
- Market and limit orders
- Custom commission rates
- Portfolio tracking
- Performance metrics

```go
func runBacktest() {
    // Initialize simulator with 0.1% commission
    simulator := exchange.NewMarginSimulator(portfolio, 0.001)
    tester := backtest.NewBacktester(&simulator)
    
    // Run backtest
    equity, err := tester.Backtest(charts, strategy)
    if err != nil {
        panic(err)
    }
    
    // Access performance metrics
    history := equity.Deposit()
    balanceHistory := equity.PortfolioHistory()
}
```

## Portfolio Management

Xoney includes tools for portfolio management and rebalancing. All weights and orders are specified in base currency:

```go
// Create a rebalancing strategy with weights in base currency
weights := map[data.Currency]float64{
    btc:  0.4,  // 0.4 BTC
    eth:  0.3,  // 0.3 ETH
    usdt: 0.3,  // 0.3 USDT
}

rebalancer := toolkit.NewRebalancePortfolio(weights)
```

## Exchange Simulation

Test your strategies with the built-in exchange simulator:

```go
// Create simulator with 0.1% commission
simulator := exchange.NewMarginSimulator(portfolio, 0.001)

// Place market order
order, _ := exchange.NewOrder(
    symbol,
    exchange.Market,
    exchange.Buy,
    price,
    amount,
)
err := simulator.PlaceOrder(order)
```

The simulator supports:
- Market and limit orders
- Spot and margin trading
- Custom commission rates
- Portfolio tracking
- Price feeds

For more examples and detailed documentation, visit our [documentation](https://github.com/quick-trade/xoney/docs).

## Toolkit

**Warning**: The toolkit package is currently in experimental stage and should be used with extreme caution. The functionality may be unstable and contain bugs.

### Grid Trading Bot

Grid trading is a strategy that places multiple orders at regular intervals above and below a set price, aiming to profit from market oscillations.

```go
package main

import (
    "github.com/quick-trade/xoney/toolkit"
    "github.com/quick-trade/xoney/common/data"
)

func main() {
    // Create auto grid generator
    // Parameters: lookback=100 candles, levels=10, deviation=1.5, total_amount=0.5 (in base currency)
    generator := NewAutoGrid(100, 10, 1.5, 0.5)
    
    // Initialize grid bot
    bot := toolkit.NewGridBot(generator, instrument)
    
    // Run backtest
    equity, err := tester.Backtest(charts, bot)
    if err != nil {
        panic(err)
    }
}
```

### Portfolio Management

Xoney includes tools for portfolio management and rebalancing. All weights and orders are specified in base currency:

```go
// Create a rebalancing strategy with weights in base currency
weights := map[data.Currency]float64{
    btc:  0.4,  // 0.4 BTC
    eth:  0.3,  // 0.3 ETH
    usdt: 0.3,  // 0.3 USDT
}

rebalancer := toolkit.NewRebalancePortfolio(weights)
```

### Custom Events

You can create custom events by implementing the Event interface:

```go
type CustomEvent struct {
    timestamp time.Time
    data      interface{}
}

func (e *CustomEvent) Occur(connector exchange.Connector) error {
    // Implement your event logic here
    return nil
}

// Create sequential events (executed one after another)
sequential := events.NewSequential(event1, event2, event3)

// Create parallel events (executed simultaneously)
parallel := events.NewParallel(event1, event2, event3)
```

Example of a custom order event:

```go
type OrderEvent struct {
    order exchange.Order
}

func (e *OrderEvent) Occur(connector exchange.Connector) error {
    return connector.PlaceOrder(e.order)
}

// Create and use the event
order, _ := exchange.NewOrder(
    symbol,
    exchange.Market,
    exchange.Buy,
    price,  // price in base currency
    amount,  // amount in base currency
)
event := &OrderEvent{order: order}
```
