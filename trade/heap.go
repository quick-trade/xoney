package trade

import (
	"xoney/common/data"
	st "xoney/internal/structures"
)

type TradeHeap struct {
	st.Heap[Trade]
}

func (h *TradeHeap) Update(candle data.Candle) {
	for i := range h.Members {
		(&h.Members[i]).Update(candle)
	}
}

func NewTradeHeap(trades ...Trade) *TradeHeap {
	return &TradeHeap{Heap: *st.NewHeap[Trade](trades...)}
}
