package chanstreaming_test

import (
	ch "github.com/diemenator/go-chanstreaming/pkg/chanstreaming"
	"testing"
	"time"
)

func TestApply(t *testing.T) {
	elementsSeen := 0
	theElementCallback := func(el int) {
		elementsSeen++
	}
	source := make(chan int, 10)
	go func() {
		defer close(source)
		for i := 1; i <= 10; i++ {
			source <- i
		}
	}()
	applied := ch.Apply[int](theElementCallback)(source)
	for range applied {

	}

	if elementsSeen != 10 {
		t.Error("Expected 10 elements seen, got", elementsSeen)
	}
}

func TestMap(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			source <- i
		}
		close(source)
	}()

	out := ch.Map(func(i int) int {
		return i * 2
	}, 5)(source)
	var results []int
	for i := range out {
		results = append(results, i)
	}
	if len(results) != 10 {
		t.Error("Expected 10 results, got", len(results))
	}
	for i, v := range results {
		if v != i*2 {
			t.Errorf("Expected %d, got %d", i*2, v)
		}
	}
}

func TestMapUnordered(t *testing.T) {
	source := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			source <- i
		}
		close(source)
	}()

	out := ch.MapUnordered(func(i int) int {
		time.Sleep(time.Millisecond * 100 * time.Duration(10-i))
		return i * 2
	}, 10)(source)
	var results []int
	for i := range out {
		results = append(results, i)
	}
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
	for i, v := range results {
		if v != (9-i)*2 {
			t.Errorf("Expected %d, got %d", i*2, v)
		}
	}
}
