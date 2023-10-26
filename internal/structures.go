package internal

import (
	"xoney/pkg/common"
	"xoney/pkg/common/events/trade"
	"xoney/pkg/errors"
)

type Equaler[T any] interface {
	IsEqual(other T) bool
}

type heap[T Equaler[T]] struct {
	members []T
}

func (h *heap[T]) Len() int { return len(h.members) }
func (h *heap[T]) Contains(v T) bool {
	for i := range h.members {
		if v.IsEqual(h.members[i]) {
			return true
		}
	}
	return false
}
func (h *heap[T]) Add(v T) { h.members = append(h.members, v) }

func (h *heap[T]) RemoveAt(i int) error {
	last := h.Len() - 1
	if last < i {
		return errors.NewOutOfIndexError(i)
	}

	h.members[i], h.members[last] = h.members[last], h.members[i]
	h.members = h.members[:last]
	return nil
}
func (h *heap[T]) Index(v T) (int, error) {
	for i := range h.members {
		if h.members[i].IsEqual(v) {
			return i, nil
		}
	}
	return 0, errors.ValueNotFoundError{}
}
func (h *heap[T]) Remove(v T) error {
	idx, err := h.Index(v)
	if err != nil {
		return err
	}
	h.RemoveAt(idx)
	return nil
}
func newHeap[T Equaler[T]](members ...T) *heap[T] {
	return &heap[T]{members: members}
}

type TradeHeap struct {
	heap[trade.Trade]
}

func (h *TradeHeap) Update(candle common.Candle) {
	for i := range h.members {
		(&h.members[i]).Update(candle)
	}
}
