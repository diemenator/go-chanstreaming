package chanstreaming_test

import (
	"testing"

	"github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestUnfoldResult(t *testing.T) {
	out := ch.UnfoldSafe[int, int](func(state int) (int, int, bool) {
		if state <= 10 {
			return state + 1, state, true
		} else {
			return state, state, false
		}
	}, 1)
	muted := ch.Muted(out)

	result := ch.ToSlice(muted)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, expected, result)
}

func TestUnfoldPanicking(t *testing.T) {
	out := ch.UnfoldSafe[int, int](func(state int) (int, int, bool) {
		if state <= 10 {
			return state + 1, state, true
		} else {
			panic(state)
		}
	}, 1)
	muted := chanstreaming.Muted(out)

	result := chanstreaming.ToSlice(muted)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, expected, result)
}
