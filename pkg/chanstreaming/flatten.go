package chanstreaming

func FlatMap[T any, R any](f func(v T) <-chan R) func(in <-chan T) <-chan R {
	return func(in <-chan T) <-chan R {
		out := make(chan R, 1)
		go func() {
			defer close(out)
			for v := range in {
				theChannel := f(v)
				for x := range theChannel {
					out <- x
				}
			}
		}()
		return out
	}
}

func FlatMapSlice[T any, R any](f func(v T) []R) func(in <-chan T) <-chan R {
	return func(in <-chan T) <-chan R {
		out := make(chan R, 1)
		go func() {
			defer close(out)
			for v := range in {
				theSlice := f(v)
				for _, x := range theSlice {
					out <- x
				}
			}
		}()
		return out
	}
}

func Identity[T any](v T) T {
	return v
}

func Concat[T any](channels ...<-chan T) <-chan T {
	channelOfChannels := FromSlice(channels)
	return FlatMap(Identity[<-chan T])(channelOfChannels)
}
