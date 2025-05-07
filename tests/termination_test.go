package chanstreamingtests_test

import (
	"context"
	"errors"
	"testing"
	"time"

	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"github.com/stretchr/testify/assert"
)

func TestMuted(t *testing.T) {
	source := make(chan ch.Result[int], 10)
	go func() {
		defer close(source)
		source <- ch.Result[int]{Data: 1}
		source <- ch.Result[int]{Error: errors.New("error")}
	}()
	muted := ch.Muted[int](source)

	result := ch.ToSlice(muted)
	assert.Equal(t, 1, len(result))
}

// CtxTestCase is a struct that holds context and the expected output, representing a test case.
type CtxTestCase struct {
	Ctx    context.Context
	Output []int
}

func TestWithContext(t *testing.T) {
	theSlice := []int{1, 2, 3, 4, 5}

	deadline := time.Now().Add(250 * time.Millisecond)
	deadlineCtx, cancelFunc := context.WithDeadline(context.Background(), deadline)
	defer cancelFunc()

	ctxTests := []CtxTestCase{
		{deadlineCtx, theSlice[:2]},
		{context.Background(), theSlice},
		{context.TODO(), theSlice},
	}

	doneInvoked := false
	doTest := func(ctxTest CtxTestCase) {
		source := ch.FromSlice(theSlice)
		throttled := ch.Throttle[int](100 * time.Millisecond)(source)
		withCtx := ch.WithContext[int](ctxTest.Ctx)(throttled)
		whenDone := ch.WhenDone[ch.Result[int]](func() { doneInvoked = true })(withCtx)
		muted := ch.Muted[int](whenDone)

		result := ch.ToSlice(muted)
		assert.Equal(t, ctxTest.Output, result)
		assert.Equal(t, true, doneInvoked)
	}

	for _, ctxTest := range ctxTests {
		doTest(ctxTest)
		ctxTest.Ctx.Deadline()
	}
}

func TestToContext(t *testing.T) {
	source := ch.FromSlice([]int{1, 2, 3, 4, 5})
	throttled := ch.Throttle[int](100 * time.Millisecond)(source)
	theCtx := ch.ToContext(throttled)

	for {
		select {
		case <-theCtx.Done():
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
