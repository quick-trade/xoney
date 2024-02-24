package events

import (
	"fmt"
	"sync"

	"github.com/quick-trade/xoney/errors"
	"github.com/quick-trade/xoney/exchange"
	"github.com/quick-trade/xoney/internal"
)

// Event is an interface designed for interaction with the exchange.
// Any strategy's response to new information is encapsulated as an Event.
// It defines the behavior for events to interact with the environment through
// the execution of their Occur method, which applies the event's effects to the
// given exchange.Connector.
type Event interface {
	// Occur executes the event's effects on the provided exchange.Connector.
	// It returns an error if the event could not be successfully executed.
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

// Sequential represents a collection of events that are meant to be executed
// in sequence. Each event in the actions slice is executed in order, and
// if any event returns an error, the sequential execution is stopped and the error
// is returned. This is used to ensure that a series of dependent events occur
// in a specific order without interruption.
type Sequential struct {
	actions []Event
}

// Occur executes each event in the Sequential actions slice in order. If an event
// fails, it returns the error and cancelling processing the remaining events.
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

// NewSequential creates a new Sequential instance with the provided actions.
// It returns a pointer to the Sequential struct which can then be used to
// execute the actions in the order they were added. If any of the actions
// return an error during execution, the subsequent actions are not executed.
//
// actions: A variadic number of Event interface implementations that represent
// the events to be added to the Sequential object.
//
// Returns: A pointer to the newly created Sequential object.
func NewSequential(actions ...Event) *Sequential {
	return &Sequential{actions: actions}
}

// Parallel represents a collection of events that are meant to be executed
// concurrently. Unlike Sequential, which executes events in a strict order,
// Parallel initiates all its contained events simultaneously, and they
// may complete in any order. This is useful for events that are independent
// from one another and can be run in parallel to improve efficiency.
type Parallel struct {
	actions []Event
}

// Occur concurrently executes all events in the Parallel structure using
// goroutines, waiting for all to complete with a WaitGroup.
//
// Errors from event executions are sent to a buffered errors channel, which
// is used to aggregate a ParallelExecutionError if any events fail. The
// channel's buffer size equals the number of actions to avoid blocking.
//
// Returns a ParallelExecutionError containing all individual errors if any
// event fails, otherwise nil.
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
		return errors.NewParallelExecutionError(errorsList)
	}
	return nil
}

func (p *Parallel) Add(actions ...Event) {
	p.actions = internal.Append(p.actions, actions...)
}

// NewParallel creates a new Parallel object that can execute events concurrently.
// Events added to the Parallel object will be initiated simultaneously, and may
// complete in any order.
//
// actions: A variadic number of Event interface implementations that represent
// the events to be added to the Parallel object.
//
// Returns: A pointer to the newly created Parallel object.
func NewParallel(actions ...Event) *Parallel {
	return &Parallel{actions: actions}
}
