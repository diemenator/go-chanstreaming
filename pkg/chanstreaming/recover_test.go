package chanstreaming_test

import (
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestMapUnorderedSafe(t *testing.T) {
	// will emit numbers and panic on even ones
	theSlice := []int{5, 4, 3, 2, 1}
	out := ch.FromSlice(theSlice)
	mapped := ch.MapUnorderedSafe[int, int](func(x int) int {
		time.Sleep(100 * time.Millisecond * time.Duration(x))
		if x%2 == 0 {
			panic(x)
		}
		return x
	}, 5)(out)
	muted := ch.Muted(mapped)

	result := ch.ToSet(muted)
	expected := map[int]struct{}{
		1: {},
		3: {},
		5: {},
	}
	assert.Equal(t, expected, result)
}

func TestMapSafe(t *testing.T) {
	theSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	out := ch.FromSlice(theSlice)
	mapped := ch.MapSafe[int, int](func(x int) int {
		time.Sleep(100 * time.Millisecond * time.Duration(x))
		if x%2 == 0 {
			panic(x)
		}
		return x
	}, 10)(out)
	muted := ch.Muted(mapped)

	result := ch.ToSlice(muted)
	expected := []int{1, 3, 5, 7, 9}
	assert.Equal(t, expected, result)
}
