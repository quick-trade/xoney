package events

import (
	"xoney/exchange"
)

type Event interface {
	Occur(connector *exchange.Connector)
}
