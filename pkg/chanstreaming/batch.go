package chanstreaming

import "time"

// NumberType is a type constraint for numeric types used for weighted batching
type NumberType interface {
	~float32 | ~float64 | ~int | ~int32 | ~int64
}

// weightedBatchElement is a helper struct to represent elements with their size or tick signals
type weightedBatchElement[T any, N NumberType] struct {
	// element is the actual data with pre-computed size
	element     T
	elementSize N
	// tick is a signal to flush the buffer
	tick bool
}

// BatchWeighted reshapes source channel into channel of batches.
// Each input element size is computed with provided `sizeFn`.
// A batch is accumulated until maxSize or maxCount or maxInterval limit is reached.
func BatchWeighted[T any, N NumberType](
	sizeFn func(element T) N,
	maxSize N,
	maxCount int,
	maxInterval time.Duration,
) func(in <-chan T) <-chan []T {
	tickSize := max(maxInterval/1000, time.Millisecond*time.Duration(5))
	return func(in <-chan T) <-chan []T {
		// launch consumer that reshapes source stream into data + tick signals
		done := make(chan struct{})
		dataAndTickChannel := make(chan weightedBatchElement[T, N], 1)
		go func() {
			defer func() {
				close(done)
				// signal completion of original source - used to shut down the ticker

			}()
			for sourceItem := range in {
				dataAndTickChannel <- weightedBatchElement[T, N]{
					element:     sourceItem,
					elementSize: sizeFn(sourceItem),
				}
			}
		}()

		go func() {
			defer close(dataAndTickChannel)
			ticker := time.Tick(tickSize)
			for {
				select {
				case <-done:
					// read completion of original source
					return
				case _, ok := <-ticker:
					// read a timer tick
					if !ok {
						return
					}
					dataAndTickChannel <- weightedBatchElement[T, N]{tick: true}
				}
			}
		}()

		// launch consumer for data items interleaved with timer events,
		// write batches on buffer overflow or time.Since(lastFlush) > maxInterval on timer tick
		out := make(chan []T, 1)
		buffer := make([]T, 0, maxCount)
		bufferSize := N(0)
		lastFlush := time.Now()
		go func() {
			defer close(out)
			for signal := range dataAndTickChannel {
				timeToFlush := false

				if signal.tick && time.Since(lastFlush) > maxInterval {
					timeToFlush = true
				}

				if !signal.tick {
					buffer = append(buffer, signal.element)
					bufferSize = bufferSize + sizeFn(signal.element)
					if (len(buffer) >= maxCount) || (bufferSize >= maxSize) {
						timeToFlush = true
					}
				}

				if timeToFlush {
					lastFlush = time.Now()
					out <- buffer
					buffer = make([]T, 0, maxCount)
					bufferSize = N(0)
				}
			}

			// tick the remaining buffer
			if len(buffer) > 0 {
				out <- buffer
			}
		}()

		return out
	}
}

// Batch reshapes the source channel into channel of batches based on elapsed interval and accumulated sizes.
func Batch[T any](maxElements int, maxInterval time.Duration) func(in <-chan T) <-chan []T {
	weightFunc := func(element T) int {
		return 1
	}
	return BatchWeighted[T, int](weightFunc, maxElements, maxElements, maxInterval)
}
