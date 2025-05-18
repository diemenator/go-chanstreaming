package chanstreamingtests_test

import (
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestThrottle(t *testing.T) {
	source := make(chan int)
	go func() {
		defer close(source)
		for i := range 5 {
			source <- i
		}
	}()

	throttle := ch.Throttle[int](10 * time.Millisecond)
	out := throttle(source)

	start := time.Now()
	result := ch.ToSlice(out)
	elapsed := time.Since(start)

	assert.Equal(t, result, []int{0, 1, 2, 3, 4})
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond) // 5 elements with 10 milliseconds delay before emitting each of it
}

func TestJitter(t *testing.T) {
	source := make(chan int)
	go func() {
		defer close(source)
		for i := range 5 {
			source <- i
		}
	}()

	throttle := ch.Jitter[int](20 * time.Millisecond)
	out := throttle(source)

	// TODO: it's better to test the histogram of delays
	start := time.Now()
	result := ch.ToSlice(out)
	elapsed := time.Since(start)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
	assert.LessOrEqual(t, elapsed, 100*time.Millisecond)
}
