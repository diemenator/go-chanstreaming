package chanstreamingtests_test

import (
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
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

	result := ch.ToSlice(out)
	expected := []int{1, 2, 3, 4, 5, 100, 101, 102, 103, 104}
	assert.ElementsMatch(t, expected, result)
}

func TestPartition(t *testing.T) {
	source := make(chan int, 10)

	// Populate source channel
	go func() {
		defer close(source)
		for i := range 10 {
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

	evenExpected := []int{0, 2, 4, 6, 8}
	oddExpected := []int{1, 3, 5, 7, 9}

	assert.Equal(t, evenExpected, evenResults)
	assert.Equal(t, oddExpected, oddResults)
}
