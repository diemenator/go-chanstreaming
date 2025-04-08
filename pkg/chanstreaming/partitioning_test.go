package chanstreaming_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
	"time"
)

func TestMergeTwoSources(t *testing.T) {
	source1 := make(chan int, 5)
	source2 := make(chan int, 5)

	// Fill source1 with values (fast producer)
	go func() {
		defer close(source1)
		for i := 1; i <= 5; i++ {
			source1 <- i
			time.Sleep(time.Millisecond * 10) // Simulate slight delay
		}
	}()

	// Fill source2 with values (slower producer)
	go func() {
		defer close(source2)
		for i := 100; i <= 104; i++ {
			source2 <- i
			time.Sleep(time.Millisecond * 20) // Simulate longer delay
		}
	}()

	// Merge both sources
	out := ch.Merge([]<-chan int{source1, source2})

	// Collect results
	results := []int{}
	for v := range out {
		results = append(results, v)
	}

	// Validate results (order not guaranteed)
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	// Ensure all expected values exist in the output (order is not fixed)
	expectedValues := map[int]bool{
		1: true, 2: true, 3: true, 4: true, 5: true,
		100: true, 101: true, 102: true, 103: true, 104: true,
	}
	for _, v := range results {
		if !expectedValues[v] {
			t.Error("Unexpected value received:", v)
		}
	}

	t.Log("Merge test passed with results:", results)
}

func TestPartition(t *testing.T) {
	source := make(chan int, 10)

	// Populate source channel
	go func() {
		defer close(source)
		for i := 0; i < 10; i++ {
			source <- i
		}
	}()

	// Partition into 2 streams (even & odd numbers)
	partitions := ch.Partition(2, func(i int) int {
		return i % 2
	})(source)

	// Make the consumer slow to ensure backpressure
	time.Sleep(time.Millisecond * 500)

	// Collect results
	evenResults := []int{}
	oddResults := []int{}

	done := make(chan struct{})
	go func() {
		for v := range partitions[0] {
			evenResults = append(evenResults, v)
		}
		done <- struct{}{}
	}()
	go func() {
		for v := range partitions[1] {
			oddResults = append(oddResults, v)
		}
		done <- struct{}{}
	}()

	<-done
	<-done

	expectedEvens := []int{0, 2, 4, 6, 8}
	expectedOdds := []int{1, 3, 5, 7, 9}

	if len(evenResults) != len(expectedEvens) {
		t.Error("Expected evens", expectedEvens, "got", evenResults)
	}
	if len(oddResults) != len(expectedOdds) {
		t.Error("Expected odds", expectedOdds, "got", oddResults)
	}
}
