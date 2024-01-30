package structures

import (
	"xoney/errors"
	"xoney/internal"
)

type Equaler[T any] interface {
	IsEqual(other *T) bool
}

// Heap is a utility structure for efficient management of a collection of objects.
// It provides methods for adding, removing, and checking for elements.
// Removal from any position is more efficient than a slice, operating in O(1) time complexity,
// but does not maintain the order of the elements.
type Heap[T Equaler[T]] struct {
	Members []T // Members contains the collection of heap elements.
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int { return len(h.Members) }

// Contains checks if the heap contains the element v.
// It compares elements using the Equaler interface.
func (h *Heap[T]) Contains(v *T) bool {
	for i := range h.Members {
		if h.Members[i].IsEqual(v) {
			return true
		}
	}
	return false
}

// Add inserts a new element v into the heap.
func (h *Heap[T]) Add(v T) { h.Members = internal.Append(h.Members, v) }

// RemoveAt removes the element at the specified index in the heap.
// It swaps the element to remove with the last element and then truncates the slice.
// This method operates in O(1) time complexity.
func (h *Heap[T]) RemoveAt(index int) error {
	last := h.Len() - 1
	if last < index {
		return errors.NewOutOfIndexError(index)
	}
	h.Members[index], h.Members[last] = h.Members[last], h.Members[index]
	h.Members = h.Members[:last]
	return nil
}

// Index finds the index of the element v in the heap.
// It uses the Equaler interface to compare elements.
func (h *Heap[T]) Index(v *T) (int, error) {
	for i := range h.Members {
		if h.Members[i].IsEqual(v) {
			return i, nil
		}
	}
	return 0, errors.ValueNotFoundError{}
}

// Remove looks for the element v and removes it from the heap.
// It finds the index of the element and then uses RemoveAt to remove it.
func (h *Heap[T]) Remove(v *T) error {
	idx, err := h.Index(v)
	if err != nil {
		return err
	}
	return h.RemoveAt(idx)
}

// Filter iterates over the elements of the heap, removing any elements
// for which the keep function returns false.
// This operation can potentially reorder the elements in the heap.
func (h *Heap[T]) Filter(keep func(*T) bool) {
	for i := h.Len() - 1; i >= 0; i-- {
		member := &h.Members[i]
		if !keep(member) {
			h.RemoveAt(i)
		}
	}
}

// NewHeap creates a new Heap with the specified initial capacity.
func NewHeap[T Equaler[T]](capacity int) *Heap[T] {
	return &Heap[T]{Members: make([]T, 0, capacity)}
}
