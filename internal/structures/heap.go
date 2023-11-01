package structures

import (
	"xoney/errors"
	"xoney/internal"
)

type Equaler[T any] interface {
	IsEqual(other *T) bool
}

type Heap[T Equaler[T]] struct {
	Members []T
}

func (h *Heap[T]) Len() int { return len(h.Members) }
func (h *Heap[T]) Contains(v *T) bool {
	for i := range h.Members {
		if h.Members[i].IsEqual(v) {
			return true
		}
	}

	return false
}
func (h *Heap[T]) Add(v T) { h.Members = internal.Append(h.Members, v) }

func (h *Heap[T]) RemoveAt(index int) error {
	last := h.Len() - 1
	if last < index {
		return errors.NewOutOfIndexError(index)
	}

	h.Members[index], h.Members[last] = h.Members[last], h.Members[index]
	h.Members = h.Members[:last]

	return nil
}

func (h *Heap[T]) Index(v *T) (int, error) {
	for i := range h.Members {
		if h.Members[i].IsEqual(v) {
			return i, nil
		}
	}

	return 0, errors.ValueNotFoundError{}
}

func (h *Heap[T]) Remove(v *T) error {
	idx, err := h.Index(v)
	if err != nil {
		return err
	}

	return h.RemoveAt(idx)
}

func NewHeap[T Equaler[T]](members ...T) *Heap[T] {
	return &Heap[T]{Members: members}
}
