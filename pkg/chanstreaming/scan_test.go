package chanstreaming_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
	"time"
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
	results := []int{}
	for v := range scanned {
		results = append(results, v)
	}
	if len(results) != 10 {
		t.Error("Expected 10 results, got", len(results))
	}
	sum := 0
	for i, v := range results {
		sum += i + 1
		if v != sum {
			t.Error("Expected", sum, "got", v)
		}
	}
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
	if result != 55 {
		t.Error("Expected 55, got", result)
	}
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
