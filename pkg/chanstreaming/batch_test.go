package chanstreaming_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
	"time"
)

func TestBatch(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
			time.Sleep(time.Millisecond * 10) // Simulate input delay
		}
	}()

	batched := ch.Batch[int](3, 50*time.Millisecond)(source)
	var results [][]int
	for batch := range batched {
		results = append(results, batch)
	}

	expectedBatches := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	}

	if len(results) != len(expectedBatches) {
		t.Errorf("Expected %d batches, got %d", len(expectedBatches), len(results))
	}
	for i, batch := range results {
		for j, val := range batch {
			if val != expectedBatches[i][j] {
				t.Errorf("Batch %d, Expected %d, got %d", i, expectedBatches[i][j], val)
			}
		}
	}
}

func TestBatchWithDecreasingFrequency(t *testing.T) {
	source := make(chan int, 1) // Large buffer to prevent blocking
	flushInterval := 200 * time.Millisecond
	batchSize := 100

	// Start batched processing
	batched := ch.Batch[int](batchSize, flushInterval)(source)
	buffered := ch.Buffered[[]int](1000)(batched)
	logged := ch.Apply(func(x []int) { t.Log("Batch received:", len(x)) })(buffered)

	// Start measuring time
	start := time.Now()
	// Source emits data with increasing intervals between elements
	go func() {
		defer close(source)
		for i := 0; i < 400; i++ {
			sleepTime := i / 20
			if sleepTime > 0 {
				time.Sleep(time.Millisecond * time.Duration(sleepTime))
			}
			source <- i
		}
	}()

	results := ch.ToSlice(logged)
	minBatchSize := batchSize
	// **Validate batch size decreasing over time**
	for i := 0; i < len(results); i++ {
		currentBatchSize := len(results[i])
		if currentBatchSize > batchSize {
			t.Error("Batch size exceeds limit", len(results[i]))
		}
		// allow +2 to account for batch including events happening at the timeframe edge
		if currentBatchSize > minBatchSize+2 {
			t.Error("Batch size not decreasing", currentBatchSize, " >> ", minBatchSize)
		}
		minBatchSize = min(minBatchSize, currentBatchSize)
	}

	t.Log("Batch test passed:", len(results), "batches emitted in", time.Since(start))
}

// Generates 8M events and sends them into a channel at controlled speed
func generateHighThroughputSource(ratePerSecond int, size int) <-chan int {
	maxSleepTime := 100 * time.Millisecond
	itemsPerMaxSleepTime := ratePerSecond * (int(time.Second) / int(maxSleepTime))
	out := make(chan int, min(10000, ratePerSecond)) // Buffered with 10k for smooth flow
	go func() {
		defer close(out)
		nextWrite := time.Now()
		for i := 0; i < size; i++ {
			out <- i
			if i%itemsPerMaxSleepTime == 0 {
				nextWrite = nextWrite.Add(maxSleepTime)
				sleepDuration := time.Until(nextWrite)
				if sleepDuration > 0 {
					time.Sleep(sleepDuration)
				}
			}
		}
	}()
	return out
}

// Test Batch Processing at High Throughput, piping 8M events through a noop batch processor and measuring throughput.
// Should process 8M events in 20000 batches in around 4 seconds.
func TestBatchWithHighThroughput(t *testing.T) {
	totalSize := 8000000
	batchSize := 20000
	source := generateHighThroughputSource(batchSize, totalSize)
	flushInterval := 34 * time.Millisecond

	// Start batched processing
	batched := ch.Batch[int](batchSize, flushInterval)(source)

	// Track processing time
	start := time.Now()
	totalEvents := 0
	batchCount := 0

	// Process batches
	for batch := range batched {
		totalEvents += len(batch)
		batchCount++
	}

	elapsed := time.Since(start)

	// Validate total count
	if totalEvents != totalSize {
		t.Error("Expected", totalSize, "events, but got", totalEvents)
	}
	t.Log("Processed", totalEvents, "events in", batchCount, "batches over", elapsed)
}
