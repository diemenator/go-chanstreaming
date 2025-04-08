package chanstreaming_test

import (
	"github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
)

func TestUnfoldResult(t *testing.T) {
	ch := chanstreaming.UnfoldSafe[int, int](func(state int) (int, int, bool) {
		if state <= 10 {
			return state + 1, state, true
		} else {
			return state, state, false
		}
	}, 1)
	muted := chanstreaming.Muted(ch)
	sliced := chanstreaming.ToSlice(muted)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if len(sliced) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(sliced))

	}
	for i, val := range sliced {
		if val != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], val)
		}
	}
}

func TestUnfoldPanicking(t *testing.T) {
	ch := chanstreaming.UnfoldSafe[int, int](func(state int) (int, int, bool) {
		if state <= 10 {
			return state + 1, state, true
		} else {
			panic(state)
		}
	}, 1)
	muted := chanstreaming.Muted(ch)
	sliced := chanstreaming.ToSlice(muted)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if len(sliced) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(sliced))
	}
	for i, val := range sliced {
		if val != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], val)
		}
	}
}
