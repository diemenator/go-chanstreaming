package chanstreaming

// FromSlice creates a readonly channel from a slice
func FromSlice[T any](data []T) <-chan T {
	output := make(chan T, 1)
	go func() {
		defer close(output)
		for _, v := range data {
			output <- v
		}
	}()
	return output
}

// Empty creates an empty readonly channel from a slice
func Empty[T any]() <-chan T {
	return FromSlice[T]([]T{})
}

// ToSlice collects all elements from the channel and returns them as a slice
func ToSlice[T any](in <-chan T) []T {
	collected := make([]T, 0)
	for x := range in {
		collected = append(collected, x)
	}
	return collected
}

// ToSet collects all elements from the channel and returns them as a map with the elements as keys
func ToSet[T comparable](in <-chan T) map[T]struct{} {
	collected := make(map[T]struct{}, 0)
	for x := range in {
		collected[x] = struct{}{}
	}
	return collected
}

// CollectWhile collects elements from the source channel while the `predicate` function is true,
// returning the slice of collected elements and the `tail` channel that starts with the first element that failed the `predicate` check
func CollectWhile[T any](predicate func(T) bool) func(in <-chan T) ([]T, <-chan T) {
	return func(in <-chan T) ([]T, <-chan T) {
		out := make(chan T, 1)
		collected := make([]T, 0)
		for x := range in {
			if predicate(x) {
				collected = append(collected, x)
			} else {
				out <- x
				break
			}
		}

		go func() {
			defer close(out)
			for x := range in {
				out <- x
			}
		}()
		return collected, out
	}
}
