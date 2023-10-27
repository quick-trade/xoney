package events

import "xoney/trade"

type Event interface {
	HandleTrades(trades *trade.TradeHeap)
}
