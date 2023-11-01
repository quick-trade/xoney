package internal

func Append[T any](slice []T, elems ...T) []T {
	var newSlice []T
	if len(slice) == cap(slice) {
		newCapacity := cap(slice) * 10
		newSlice = make([]T, len(slice), newCapacity)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = append(slice, elems...)
	return slice
}
