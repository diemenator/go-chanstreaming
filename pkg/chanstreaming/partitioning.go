package chanstreaming

// Merge takes multiple input channels and combines them into a single output channel.
// It launches `len(sources)` worker goroutines to consume from the sources concurrently
// until all of them have been closed.
func Merge[T any](sources []<-chan T) <-chan T {
	out := make(chan T, 1)
	done := make(chan struct{}, 1)
	srcCount := len(sources)
	go func() {
		for _, source := range sources {
			go func(src <-chan T) {
				defer func() {
					done <- struct{}{}
				}()
				for v := range src {
					out <- v
				}
			}(source) // Launch a goroutine for each source, passing source explicitly
		}
	}()

	go func() {
		for i := 0; i < srcCount; i++ {
			<-done
		}
		close(done)
		close(out)
	}()

	return out
}

// Partition splits an input channel into `maxPartitions` separate output channels
// based on a `partitioner` function that returns a partition number.
func Partition[T any](maxPartitions int, partitioner func(T) int) func(in <-chan T) []<-chan T {
	return func(in <-chan T) []<-chan T {
		partitions := make([]chan T, maxPartitions)
		// Create partition channels with size 1
		for i := range partitions {
			partitions[i] = make(chan T, 1)
		}

		// launch source consumer
		go func() {
			// Close all partitions when the input channel is closed
			defer func() {
				for _, ch := range partitions {
					close(ch)
				}
			}()

			for item := range in {
				idx := partitioner(item)
				if idx < 0 {
					idx = -idx // Convert negative to positive
				}

				idx = idx % maxPartitions

				partitions[idx] <- item
			}
		}()

		// Convert to read-only channels
		readOnlyOuts := make([]<-chan T, maxPartitions)
		for i, ch := range partitions {
			readOnlyOuts[i] = ch
		}

		return readOnlyOuts
	}
}

func MergeSlice[T any](sources []<-chan T) <-chan T {
	return Merge(sources)
}

func Merged[T any](sources ...<-chan T) <-chan T {
	return MergeSlice(sources)
}
