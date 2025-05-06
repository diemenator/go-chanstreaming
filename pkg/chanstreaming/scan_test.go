package chanstreaming_test

import (
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()
	scanned := ch.Scan[int, int](func(acc, el int) int { return acc + el }, 0)(source)

	result := ch.ToSlice(scanned)
	expected := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	assert.Equal(t, expected, result)
}

func TestFold(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()
	folded := ch.Fold[int, int](func(acc, el int) int { return acc + el }, 0)(source)
	result := <-folded

	assert.Equal(t, 55, result)
}

func TestWithSlidingWindowTimed(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()

	slidingWindow := ch.WithSlidingWindowTimed[int](time.Millisecond * 100)(source)
	hardCopies := ch.Mapped[[]int, []int](func(x []int) []int {
		result := make([]int, len(x))
		copy(result, x)
		return result
	})(slidingWindow)

	results := ch.ToSlice(hardCopies)
	if len(results) != 10 { // including empty starting window
		t.Error("Expected 10 results, got", len(results))
	}

	for i, r := range results {
		if len(r) <= 12 && len(r) >= 9 && i > 10 {
			t.Error("Expected 9 to 12 results in window", i, "got", len(r))
		} else if len(r) <= i+1 && i > 10 {
			t.Error("Expected up to", i+1, "results in window", i, "got", len(r))
		}
	}
}

func TestWithSlidingWindowCount(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()

	slidingWindow := ch.WithSlidingWindowCount[int](5)(source)
	hardCopies := ch.Mapped[[]int, []int](func(x []int) []int {
		result := make([]int, len(x))
		copy(result, x)
		return result
	})(slidingWindow)
	results := ch.ToSlice(hardCopies)
	if len(results) != 10 {
		t.Error("Expected 10 results, got", len(results))
	}
	for i, r := range results {
		if (len(r) >= 6 || len(r) <= 4) && i > 5 {
			t.Error("Expected between 4 and 6 results in window", i, "got", len(r))
		} else if len(r) > i+1 && i <= 5 {
			t.Error("Expected up to", i+1, "results in window", i, "got", len(r))
		}
	}
}
