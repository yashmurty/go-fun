package main

import (
	"math/rand"
	"testing"
)

func TestConcurrentQS(t *testing.T) {
	sortSize := 1000000
	unsorted := make([]int, 0, sortSize)
	unsorted = rand.Perm(sortSize)
	ConcurrentQS(unsorted)
	for i := 0; i < sortSize; i++ {
		if unsorted[i] != i {
			t.Errorf("expecting sorted slice")
			return
		}
	}
}

func BenchmarkSimpleQS1000000(b *testing.B) {
	sortSize := 1000000
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		SimpleQS(unsorted)
	}
}

func BenchmarkConcurrentQS1000000(b *testing.B) {
	sortSize := 1000000
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		ConcurrentQS(unsorted)
	}
}
