package events

import (
	"fmt"
	"strings"
	"sync"

	"xoney/exchange"
	"xoney/internal"
)

type Event interface {
	Occur(connector exchange.Connector) error
}

type OpenOrder struct {
	order exchange.Order
}

func (o *OpenOrder) Occur(connector exchange.Connector) error {
	return connector.PlaceOrder(o.order)
}

func NewOpenOrder(order exchange.Order) *OpenOrder {
	return &OpenOrder{order: order}
}

type CancelOrder struct {
	id exchange.OrderID
}

func (o *CancelOrder) Occur(connector exchange.Connector) error {
	return connector.CancelOrder(o.id)
}

func NewCancelOrder(id exchange.OrderID) *CancelOrder {
	return &CancelOrder{id: id}
}

type EditOrder struct {
	cancelID exchange.OrderID
	order    exchange.Order
}

func (e *EditOrder) Occur(connector exchange.Connector) error {
	if err := connector.CancelOrder(e.cancelID); err != nil {
		return fmt.Errorf("error canceling order: %w", err)
	}

	if err := connector.PlaceOrder(e.order); err != nil {
		return fmt.Errorf("error placing order: %w", err)
	}

	return nil
}

func NewEditOrder(cancelID exchange.OrderID, newOrder exchange.Order) *EditOrder {
	return &EditOrder{
		cancelID: cancelID,
		order:    newOrder,
	}
}

type Sequential struct {
	actions []Event
}

func (s *Sequential) Occur(connector exchange.Connector) error {
	for _, action := range s.actions {
		if err := action.Occur(connector); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sequential) Add(actions ...Event) {
	s.actions = append(s.actions, actions...)
}

func (s *Sequential) Events() []Event {
	return s.actions
}

func NewSequential(actions ...Event) *Sequential {
	return &Sequential{actions: actions}
}

type Parallel struct {
	actions []Event
}

func (p *Parallel) Occur(connector exchange.Connector) error {
	var wg sync.WaitGroup
	errorsChan := make(chan string, len(p.actions))

	for _, action := range p.actions {
		wg.Add(1)
		go func(act Event) {
			defer wg.Done()
			if err := act.Occur(connector); err != nil {
				errorsChan <- err.Error()
			}
		}(action)
	}

	wg.Wait()
	close(errorsChan)

	var errorsList []string
	for err := range errorsChan {
		errorsList = append(errorsList, err)
	}

	if len(errorsList) > 0 {
		return fmt.Errorf("errors occurred in parallel execution: %s", strings.Join(errorsList, "; "))
	}
	return nil
}

func (p *Parallel) Add(actions ...Event) {
	p.actions = internal.Append(p.actions, actions...)
}

func NewParallel(actions ...Event) *Parallel {
	return &Parallel{actions: actions}
}
