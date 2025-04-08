package chanstreaming_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"slices"
	"testing"
)

func TestFlatMap(t *testing.T) {
	theSlice := []int{1, 2, 3, 4, 5}
	// Create a channel of integers from the slice
	channel := ch.FromSlice[int](theSlice)
	channel = ch.FlatMap[int, int](func(x int) <-chan int {
		return ch.FromSlice[int]([]int{x, x})
	})(channel)

	result := ch.ToSlice(channel)
	expected := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFlatMapSlice(t *testing.T) {
	theSlice := []int{1, 2, 3, 4, 5}
	// Create a channel of integers from the slice
	channel := ch.FromSlice[int](theSlice)
	channel = ch.FlatMapSlice[int, int](func(x int) []int {
		return []int{x, x}
	})(channel)
	result := ch.ToSlice(channel)
	expected := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
