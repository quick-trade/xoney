package events

import "xoney/trade"

type VolumeDistributor interface {
	TradeVolume(trade trade.Trade) // TODO: add system state
}
