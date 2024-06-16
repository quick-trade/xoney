package toolkit

import (
	"fmt"
	"math"
	"sync"

	"github.com/quick-trade/xoney/common"
	"github.com/quick-trade/xoney/common/data"
	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/events"
	"github.com/quick-trade/xoney/exchange"
	"github.com/quick-trade/xoney/internal"
	st "github.com/quick-trade/xoney/strategy"
)

type (
	BaseWeight       float64
	PortfolioWeights map[data.Currency]BaseWeight
)

func NewPortfolioWeights(distribution map[data.Currency]BaseWeight, epsilon float64) (*PortfolioWeights, error) {
	weights := PortfolioWeights(distribution)
	if err := weights.isValid(epsilon); err != nil {
		return nil, err
	}
	return &weights, nil
}

func (f PortfolioWeights) isValid(epsilon float64) error {
	sumWeights := 0.0

	for _, weight := range f {
		sumWeights += math.Abs(float64(weight))
	}

	if math.Abs(sumWeights-1) > epsilon {
		return errors.NewInvalidWeightsError(sumWeights)
	}

	return nil
}

func (pw PortfolioWeights) Synchronize(
	current common.BaseDistribution,
	prices map[data.Currency]float64,
	mainCurrency data.Currency,
) (target common.BaseDistribution, err error) {
	totalQuote := 0.0
	totalQuoteWeight := 0.0
	missingCurrencyErr := errors.NewMissingCurrencyError(internal.DefaultCapacity)
	success := true

	// Calculate the total value of quote distribution
	for currency, amount := range current {
		price, ok := prices[currency]
		if !ok && currency == mainCurrency {
			price, ok = 1, true
		}
		if !ok {
			missingCurrencyErr.Add(currency.String())
			success = false
		} else {
			totalQuote += amount * price
			totalQuoteWeight += float64(pw[currency]) * price
		}
	}

	if !success {
		return nil, missingCurrencyErr
	}

	// Calculate the target base distribution based on weights
	target = make(common.BaseDistribution)

	for currency, weight := range pw {
		target[currency] = float64(weight) * totalQuote / totalQuoteWeight
	}

	if !success {
		return nil, missingCurrencyErr
	}

	return target, nil
}

type CapitalAllocator interface {
	Start(charts data.ChartContainer) error
	Next(candle data.InstrumentCandle) (PortfolioWeights, error)
	MinDurations() st.Durations
}

type RebalancePortfolio struct {
	weights             PortfolioWeights
	currentDistribution common.BaseDistribution
	lastPrices          map[data.Currency]float64
	mainCurrency        data.Currency
}

func (r *RebalancePortfolio) Occur(connector exchange.Connector) error {
	target, err := r.getTargetDistribution(connector)
	if err != nil {
		return fmt.Errorf("failed to get target assets distribution: %w", err)
	}

	rebalanceEvents, err := r.rebalance(target)
	if err != nil {
		return fmt.Errorf("failed to initialize asset rebalance: %w", err)
	}
	err = rebalanceEvents.Occur(connector)

	if err != nil {
		return fmt.Errorf("failed to rebalance assets: %w", err)
	}

	return nil
}

func (r *RebalancePortfolio) rebalance(target common.BaseDistribution) (events.Event, error) {
	difference := r.calculateDifference(target)
	sellDifferences, buyDifferences := r.sortDifference(difference)

	sellEvents, err := r.newOrders(sellDifferences)
	if err != nil {
		return nil, fmt.Errorf("failed to create sell orders: %w", err)
	}

	buyEvents, err := r.newOrders(buyDifferences)
	if err != nil {
		return nil, fmt.Errorf("failed to create buy orders: %w", err)
	}

	return events.NewSequential(sellEvents, buyEvents), nil
}

func (r *RebalancePortfolio) newOrders(differences common.BaseDistribution) (events.Event, error) {
	Events := make([]events.Event, 0, len(differences))

	for currency, amount := range differences {
		if amount == 0 {
			continue
		}

		var side exchange.OrderSide

		// Assuming r.mainCurrency is the quote currency
		if currency == r.mainCurrency {
			continue
		}
		symbol := data.NewSymbolFromCurrencies(currency, r.mainCurrency)
		if amount < 0 {
			side = exchange.Sell
		} else {
			side = exchange.Buy
		}

		price, priceExists := r.lastPrices[currency]
		if !priceExists {
			return nil, errors.NewNoPriceError(currency.String())
		}

		order, err := exchange.NewOrder(*symbol, exchange.Market, side, price, math.Abs(amount))
		if err != nil {
			return nil, fmt.Errorf("failed to create a new order for %s: %w", currency, err)
		}

		openOrder := events.NewOpenOrder(*order)
		Events = internal.Append(Events, events.Event(openOrder))
	}

	return events.NewParallel(Events...), nil
}

func (r *RebalancePortfolio) sortDifference(difference common.BaseDistribution) (sellDifferences, buyDifferences common.BaseDistribution) {
	sellDifferences = make(common.BaseDistribution, len(difference))
	buyDifferences = make(common.BaseDistribution, len(difference))

	for currency, amount := range difference {
		if amount < 0 {
			sellDifferences[currency] = amount
		} else if amount != 0 {
			buyDifferences[currency] = amount
		}
	}

	return sellDifferences, buyDifferences
}

func (r *RebalancePortfolio) calculateDifference(target common.BaseDistribution) common.BaseDistribution {
	difference := make(common.BaseDistribution)

	for currency, targetVolume := range target {
		currentVolume, exists := r.currentDistribution[currency]
		if exists {
			difference[currency] = targetVolume - currentVolume
		} else {
			difference[currency] = targetVolume
		}
	}

	return difference
}

func (r *RebalancePortfolio) getTargetDistribution(connector exchange.Connector) (common.BaseDistribution, error) {
	portfolio := connector.Portfolio()

	r.currentDistribution = portfolio.Assets()

	r.mainCurrency = portfolio.MainCurrency()
	symbols := make([]data.Symbol, 0, len(r.weights))

	for currency := range r.weights {
		if currency == r.mainCurrency {
			continue
		}
		symbol := data.NewSymbolFromCurrencies(currency, r.mainCurrency)
		symbols = append(symbols, *symbol)
	}

	// Get the portfolio and prices asynchronously
	var currentDistribution common.BaseDistribution

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		currentDistribution = connector.Portfolio().Assets()
	}()

	prices, err := r.getPrices(symbols, connector)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}

	wg.Wait()

	// Convert prices into a map where keys are Currencies and values are float64
	r.lastPrices = make(map[data.Currency]float64)
	for _, symbolPrice := range prices {
		r.lastPrices[symbolPrice.Symbol.Base()] = symbolPrice.Price
	}

	target, err := r.weights.Synchronize(currentDistribution, r.lastPrices, r.mainCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to synchronize portfolio with weights: %w", err)
	}

	return target, nil
}

func (r *RebalancePortfolio) getPrices(
	symbols []data.Symbol,
	connector exchange.Connector,
) ([]exchange.SymbolPrice, error) {
	pricesChan, err := connector.GetPrices(symbols)
	prices := make([]exchange.SymbolPrice, 0, len(symbols))

	for {
		select {
		case price, ok := <-pricesChan:
			if ok {
				prices = internal.Append(prices, price)
			}
		case e, ok := <-err:
			if ok && e != nil {
				return nil, e
			}
		}

		if len(prices) == len(symbols) {
			break
		}
	}
	return prices, nil
}

func NewRebalancePortfolio(weights PortfolioWeights) *RebalancePortfolio {
	return &RebalancePortfolio{weights: weights, currentDistribution: nil}
}

type CapitalAllocationBot struct {
	allocator CapitalAllocator
}

func NewCapitalAllocationBot(allocator CapitalAllocator) *CapitalAllocationBot {
	return &CapitalAllocationBot{allocator: allocator}
}

func (c *CapitalAllocationBot) MinDurations() st.Durations {
	return c.allocator.MinDurations()
}

func (c *CapitalAllocationBot) Start(charts data.ChartContainer) error {
	return c.allocator.Start(charts)
}

func (c *CapitalAllocationBot) Next(candle data.InstrumentCandle) ([]events.Event, error) {
	weights, err := c.allocator.Next(candle)
	if err != nil {
		return nil, err
	}

	return []events.Event{NewRebalancePortfolio(weights)}, nil
}
