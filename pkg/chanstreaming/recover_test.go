package chanstreaming_test

import (
	"github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"maps"
	"slices"
	"testing"
	"time"
)

func TestMapUnorderedSafe(t *testing.T) {
	// will emit numbers and panic on even ones
	theSlice := []int{5, 4, 3, 2, 1}
	ch := chanstreaming.FromSlice(theSlice)
	mapped := chanstreaming.MapUnorderedSafe[int, int](func(x int) int {
		time.Sleep(100 * time.Millisecond * time.Duration(x))
		if x%2 == 0 {
			panic(x)
		}
		return x
	}, 5)(ch)
	muted := chanstreaming.Muted(mapped)
	unordered := chanstreaming.ToSet(muted)
	expected := map[int]struct{}{
		1: {},
		3: {},
		5: {},
	}
	if !maps.Equal(unordered, expected) {
		t.Errorf("Expected %v, got %v", expected, unordered)
	}
}

func TestMapSafe(t *testing.T) {
	theSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch := chanstreaming.FromSlice(theSlice)
	mapped := chanstreaming.MapSafe[int, int](func(x int) int {
		time.Sleep(100 * time.Millisecond * time.Duration(x))
		if x%2 == 0 {
			panic(x)
		}
		return x
	}, 10)(ch)

	muted := chanstreaming.Muted(mapped)
	sliced := chanstreaming.ToSlice(muted)
	expected := []int{1, 3, 5, 7, 9}
	if !slices.Equal(sliced, expected) {
		t.Errorf("Expected %v, got %v", expected, sliced)
	}
}
