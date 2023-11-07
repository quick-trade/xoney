package internal

const DefaultCapacity = 100
const CapacityMultiplier = 10

func Append[T any](slice []T, elems ...T) []T {
	var newSlice []T

	if len(slice) == cap(slice) {
		newCapacity := cap(slice) * CapacityMultiplier
		newSlice = make([]T, len(slice), newCapacity)
		copy(newSlice, slice)
		slice = newSlice
	}

	slice = append(slice, elems...)

	return slice
}
