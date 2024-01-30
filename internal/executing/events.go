package executing

import (
	"fmt"
	"xoney/events"
	"xoney/exchange"
)

// ProcessEvent takes an exchange.Connector and an events.Event as parameters and processes the
// given event. This function is extracted to avoid repetition and ensure consistent execution
// across tests and live trading environments. It handles the occurrence of the event
// within the exchange connector context and returns an error if the operation fails.
func ProcessEvent(connector exchange.Connector, event events.Event) error {
	if event == nil {
		return nil
	}

	if err := event.Occur(connector); err != nil {
		return fmt.Errorf("failed to process event: %w", err)
	}

	return nil
}
