package chanstreaming_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
)

func TestFromSlice(t *testing.T) {
	source := ch.FromSlice([]int{1, 2, 3, 4, 5})
	results := []int{}
	for v := range source {
		results = append(results, v)
	}
	if len(results) != 5 {
		t.Error("Expected 5 results, got", len(results))
	}
	for i, v := range results {
		if v != i+1 {
			t.Error("Expected", i+1, "got", v)
		}
	}
}

func TestCollectWhile(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()
	out, tailChannel := ch.CollectWhile[int](func(i int) bool { return i < 5 })(source)
	if len(out) != 4 {
		t.Error("Expected 4 results, got", len(out))
	}
	for i, v := range out {
		if v != i+1 {
			t.Error("Expected", i+1, "got", v)
		}
	}
	tailSlice := ch.ToSlice(tailChannel)
	if len(tailSlice) != 6 {
		t.Error("Expected 6 results in tail, got", len(tailSlice))
	}
}
