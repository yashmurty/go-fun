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

func BenchmarkSimpleQS(b *testing.B) {
	sortSize := 1000000
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		SimpleQS(unsorted)
	}
}

func BenchmarkConcurrentQS(b *testing.B) {
	sortSize := 1000000
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		ConcurrentQS(unsorted)
	}
}

func BenchmarkConcurrentQSLimited1000000(b *testing.B) {
	sortSize := 1000000
	MAXGOROUTINES = 1000000
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		ConcurrentQSLimited(unsorted)
	}
}

func BenchmarkConcurrentQSLimited1000(b *testing.B) {

	sortSize := 1000000
	MAXGOROUTINES = 1000
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		ConcurrentQSLimited(unsorted)
	}
}

func BenchmarkConcurrentQSLimited10(b *testing.B) {
	sortSize := 1000000
	MAXGOROUTINES = 10
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		ConcurrentQSLimited(unsorted)
	}
}

func BenchmarkConcurrentQSLimited2(b *testing.B) {
	sortSize := 1000000
	MAXGOROUTINES = 2
	unsorted := make([]int, 0, sortSize)
	for i := 0; i < b.N; i++ {
		unsorted = rand.Perm(sortSize)
		ConcurrentQSLimited(unsorted)
	}
}
