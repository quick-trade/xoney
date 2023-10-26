package events

import (
	"xoney/internal"
)

type Event interface {
	HandleTrades(trades *internal.TradeHeap)
}
