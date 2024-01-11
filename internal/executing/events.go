package executing

import (
	"fmt"

	"xoney/events"
	"xoney/exchange"
)

func ProcessEvent(connector exchange.Connector, event events.Event) error {
	if event == nil {
		return nil
	}

	if err := event.Occur(connector); err != nil {
		return fmt.Errorf("failed to process event: %w", err)
	}

	return nil
}
