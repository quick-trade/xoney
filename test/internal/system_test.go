package internal_test

import (
	"reflect"
	"sort"
	"testing"

	"xoney/internal"
)

func Map() map[string]int {
	return map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
}

func TestMapContains(t *testing.T) {
	m := Map()
	items := []string{"a", "b", "c"}

	for _, i := range items {
		if !internal.Contains(m, i) {
			t.Errorf("map contains %v, but .Contains() have a bug", i)
		}
	}
}

func TestMapKeys(t *testing.T) {
	emptyMap := make(map[int]string)
	emptyResult := internal.MapKeys(emptyMap)
	if len(emptyResult) != 0 {
		t.Errorf("Expected an empty slice, got: %v", emptyResult)
	}

	nonEmptyMap := map[int]string{1: "one", 2: "two", 3: "three"}
	expectedResult := []int{1, 2, 3}
	result := internal.MapKeys(nonEmptyMap)

	sort.Ints(result)
	sort.Ints(expectedResult)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected: %v, got: %v", expectedResult, result)
	}

	mixedTypeMap := map[any]string{"a": "apple", 2: "banana", 3.14: "cherry"}
	expectedMixedType := []any{"a", 2, 3.14}
	resultMixedType := internal.MapKeys(mixedTypeMap)

	success := true

	if len(resultMixedType) != len(expectedMixedType) {
		success = false
	} else {
		for _, vReal := range resultMixedType {
			found := false
			for _, vExp := range expectedMixedType {
				if vExp == vReal {
					found = true
				}
			}

			if !found {
				success = false
			}
		}
	}

	if !success {
		t.Errorf("Expected: %v, got: %v", expectedMixedType, resultMixedType)
	}
}
