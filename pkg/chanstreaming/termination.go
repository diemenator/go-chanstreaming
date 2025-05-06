package chanstreaming

import (
	"context"
)

// WithContext transforms the source channel to a channel that completes on cancellation of the context.Context
func WithContext[T any](ctx context.Context) func(in <-chan T) <-chan Result[T] {
	return func(in <-chan T) <-chan Result[T] {
		output := make(chan Result[T], 1)
		go func() {
			defer close(output)
			for {
				select {
				case <-ctx.Done():
					// ctx.Err() != nil when ctx.Done() is closed
					output <- Result[T]{Error: ctx.Err()}
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					output <- Result[T]{Data: data}
				}
			}
		}()

		return output
	}
}

// ViaKillSwitch implements graceful termination on a signal received or channel closing from the `killSwitch` channel.
func ViaKillSwitch[T any, K any](killSwitch <-chan K) func(in <-chan T) <-chan T {
	return func(in <-chan T) <-chan T {
		output := make(chan T, 1)
		go func() {
			defer close(output)
			for {
				select {
				case <-killSwitch:
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					output <- data
				}
			}
		}()

		return output
	}
}

// WhenDone creates a new channel that tracks source termination and invokes the `callback` when the source is done.
func WhenDone[T any](callback func()) func(in <-chan T) <-chan T {
	return func(in <-chan T) <-chan T {
		out := make(chan T, 1)
		go func() {
			defer callback()
			defer close(out)
			for x := range in {
				out <- x
			}
		}()

		return out
	}
}

// ToContext transforms the source channel to a context.Context. The context cancels on source channel closing.
func ToContext[T any](in <-chan T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-in
		cancel()
	}()

	return ctx
}
