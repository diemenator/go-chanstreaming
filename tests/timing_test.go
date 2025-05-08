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
	n := 100
	go func() {
		defer close(source)
		for i := range n {
			source <- i
		}
	}()

	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// r = rand.New(rand.NewSource(1))
	throttle := ch.Jitter[int](10 * time.Millisecond)
	out := throttle(source)

	count := []int{0, 0, 0}

	for i := range n {
		start := time.Now()
		v := <-out
		elapsed := time.Since(start)
		assert.Equal(t, i, v)
		if elapsed < 1*time.Millisecond {
			count[0]++
		} else if elapsed < 6*time.Millisecond {
			count[1]++
		} else {
			count[2]++
		}
	}

	assert.GreaterOrEqual(t, count[0], n/2)
	assert.GreaterOrEqual(t, count[1], n/3)
	assert.Equal(t, 0, count[2])
}
