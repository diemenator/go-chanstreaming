package chanstreamingtests_test

import (
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	elementsSeen := 0
	theElementCallback := func(el int) {
		elementsSeen++
	}
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()
	applied := ch.Apply(theElementCallback)(source)
	for range applied {

	}

	assert.Equal(t, 10, elementsSeen)
}

func TestMap(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := range 10 {
			source <- i
		}
	}()

	out := ch.Map(func(i int) int {
		return i * 2
	}, 5)(source)

	result := ch.ToSlice(out)
	expected := []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}
	assert.Equal(t, expected, result)
}

func TestMapUnordered(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := range 10 {
			source <- i
		}
	}()

	out := ch.MapUnordered(func(i int) int {
		time.Sleep(time.Millisecond * 100 * time.Duration(10-i))
		return i * 2
	}, 10)(source)

	result := ch.ToSlice(out)
	expected := []int{18, 16, 14, 12, 10, 8, 6, 4, 2, 0}
	assert.Equal(t, expected, result)
}
