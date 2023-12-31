package internal_test

import (
	"testing"

	"xoney/errors"
	"xoney/internal/structures"
)

type MockEqualer int

func (m MockEqualer) IsEqual(other *MockEqualer) bool {
	return m == *other
}

func TestHeapAddContains(t *testing.T) {
	h := structures.NewHeap[MockEqualer](5)
	elem := MockEqualer(42)

	h.Add(elem)

	if !h.Contains(&elem) {
		t.Errorf("Expected heap to contain %v, but it didn't", elem)
	}
}

func TestHeapRemove(t *testing.T) {
	h := structures.NewHeap[MockEqualer](5)
	elem := MockEqualer(42)

	h.Add(elem)
	err := h.Remove(&elem)

	if err != nil || h.Contains(&elem) {
		t.Errorf("Expected heap to remove %v, but it didn't", elem)
	}
}

func TestHeapFilter(t *testing.T) {
	h := structures.NewHeap[MockEqualer](5)
	elem1, elem2, elem3 := MockEqualer(99), MockEqualer(24), MockEqualer(15)

	h.Add(elem1)
	h.Add(elem2)
	h.Add(elem3)
	h.Add(elem2)

	h.Filter(func(e *MockEqualer) bool {
		return *e == 24
	})

	if h.Contains(&elem1) || !h.Contains(&elem2) || h.Contains(&elem3) {
		t.Errorf("Expected heap to only contain %v after filtering, but it didn't", elem1)
	}
}

func TestHeapIndex(t *testing.T) {
	h := structures.NewHeap[MockEqualer](5)
	elem := MockEqualer(42)

	h.Add(elem)

	index, err := h.Index(&elem)

	if err != nil || index != 0 {
		t.Errorf("Expected index of %v to be 0, but got %d", elem, index)
	}
}

func TestHeapRemoveAt(t *testing.T) {
	h := structures.NewHeap[MockEqualer](5)
	elem := MockEqualer(42)

	h.Add(elem)
	err := h.RemoveAt(0)

	if err != nil || h.Contains(&elem) {
		t.Errorf("Expected heap to remove %v at index 0, but it didn't", elem)
	}
}

func TestHeapLen(t *testing.T) {
	h := structures.NewHeap[MockEqualer](5)
	elem1, elem2 := MockEqualer(42), MockEqualer(24)

	h.Add(elem1)
	h.Add(elem2)

	length := h.Len()

	if length != 2 {
		t.Errorf("Expected heap length to be 2, but got %d", length)
	}
}

func TestHeapIndexError(t *testing.T) {
	// Create a Heap with some elements
	myHeap := structures.NewHeap[MockEqualer](5)
	myHeap.Add(1)
	myHeap.Add(2)
	myHeap.Add(3)

	// Attempt to remove an element at an out-of-bounds index
	index := 5
	err := myHeap.RemoveAt(index)

	// Check if the error is of type errors.OutOfIndexError
	if _, ok := err.(errors.OutOfIndexError); !ok {
		t.Errorf("Expected OutOfIndexError, got: %v", err)
	}
}

func TestHeapNotFoundError(t *testing.T) {
	// Create a Heap with some elements
	myHeap := structures.NewHeap[MockEqualer](5)
	myHeap.Add(1)
	myHeap.Add(2)
	myHeap.Add(3)

	// Attempt to get the index of a value not present in the heap
	var valueNotInHeap MockEqualer = 42
	_, err := myHeap.Index(&valueNotInHeap)

	// Check if the error is of type ValueNotFoundError
	if _, ok := err.(errors.ValueNotFoundError); !ok {
		t.Errorf("Expected ValueNotFoundError, got: %v", err)
	}
}

func TestHeapRemoveError(t *testing.T) {
	// Create a Heap with some elements
	myHeap := structures.NewHeap[MockEqualer](5)
	myHeap.Add(1)
	myHeap.Add(2)
	myHeap.Add(3)

	// Attempt to remove a value not present in the heap
	valueNotInHeap := MockEqualer(42)
	err := myHeap.Remove(&valueNotInHeap)

	// Check if the error is of type ValueNotFoundError
	if _, ok := err.(errors.ValueNotFoundError); !ok {
		t.Errorf("Expected ValueNotFoundError, got: %v", err)
	}
}
