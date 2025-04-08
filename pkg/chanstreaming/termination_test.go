package chanstreaming_test

import (
	"context"
	"errors"
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
	"time"
)

func TestMuted(t *testing.T) {
	source := make(chan ch.Result[int], 10)
	go func() {
		defer close(source)
		source <- ch.Result[int]{Data: 1}
		source <- ch.Result[int]{Error: errors.New("error")}
	}()
	muted := ch.Muted[int](source)
	results := ch.ToSlice(muted)
	if len(results) != 1 {
		t.Error("Expected 1 results, got", len(results))
	}
}

type CtxTest struct {
	Ctx    context.Context
	Output []int
}

func TestWithContext(t *testing.T) {
	theSlice := []int{1, 2, 3, 4, 5}

	deadline := time.Now().Add(250 * time.Millisecond)
	deadlineCtx, cancelFunc := context.WithDeadline(context.Background(), deadline)
	defer cancelFunc()

	ctxTests := []CtxTest{
		{deadlineCtx, theSlice[:2]},
		{context.Background(), theSlice},
		{context.TODO(), theSlice},
	}

	doneInvoked := false
	doTest := func(ctxTest CtxTest) {
		source := ch.FromSlice(theSlice)
		throttled := ch.Throttle[int](100 * time.Millisecond)(source)
		withCtx := ch.WithContext[int](ctxTest.Ctx)(throttled)
		whenDone := ch.WhenDone[ch.Result[int]](func() { doneInvoked = true })(withCtx)
		muted := ch.Muted[int](whenDone)
		results := ch.ToSlice(muted)
		if len(results) != len(ctxTest.Output) {
			t.Error("Expected", len(ctxTest.Output), "results, got", len(results))
		}
		for i, v := range results {
			if v != ctxTest.Output[i] {
				t.Error("Expected", ctxTest.Output[i], "got", v)
			}
		}
		if !doneInvoked {
			t.Error("Expected doneInvoked to be true")
		}
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
