package structures

import "xoney/errors"

type Equaler[T any] interface {
	IsEqual(other T) bool
}

type Heap[T Equaler[T]] struct {
	Members []T
}

func (h *Heap[T]) Len() int { return len(h.Members) }
func (h *Heap[T]) Contains(v T) bool {
	for i := range h.Members {
		if v.IsEqual(h.Members[i]) {
			return true
		}
	}
	return false
}
func (h *Heap[T]) Add(v T) { h.Members = append(h.Members, v) }

func (h *Heap[T]) RemoveAt(i int) error {
	last := h.Len() - 1
	if last < i {
		return errors.NewOutOfIndexError(i)
	}

	h.Members[i], h.Members[last] = h.Members[last], h.Members[i]
	h.Members = h.Members[:last]
	return nil
}
func (h *Heap[T]) Index(v T) (int, error) {
	for i := range h.Members {
		if h.Members[i].IsEqual(v) {
			return i, nil
		}
	}
	return 0, errors.ValueNotFoundError{}
}
func (h *Heap[T]) Remove(v T) error {
	idx, err := h.Index(v)
	if err != nil {
		return err
	}
	h.RemoveAt(idx)
	return nil
}
