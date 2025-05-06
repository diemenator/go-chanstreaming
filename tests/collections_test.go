package chanstreamingtests_test

import (
	"testing"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	source := ch.FromSlice([]int{1, 2, 3, 4, 5})

	result := ch.ToSlice(source)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestCollectWhile(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()
	result, tailChannel := ch.CollectWhile[int](func(i int) bool { return i < 5 })(source)

	tailResult := ch.ToSlice(tailChannel)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
	assert.Equal(t, []int{5, 6, 7, 8, 9, 10}, tailResult)
}

func TestEmpty(t *testing.T) {
	source := ch.Empty[int]()

	result := ch.ToSlice(source)
	assert.Empty(t, result)
}
