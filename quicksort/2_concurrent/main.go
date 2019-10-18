package main

import (
	"fmt"
	"math/rand"
	"time"
)

// SimpleQS is a simple implementation of quick sort without using any go routines.
func SimpleQS(a []int) []int {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1
	pivot := rand.Int() % len(a)

	// Swap pivot with right most element
	a[right], a[pivot] = a[pivot], a[right]

	for i := range a {
		if a[i] < a[right] {
			// Swap i'th position with left position
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	// Swap right position (previous pivot) with left position ((left -1) is smaller than right)
	a[right], a[left] = a[left], a[right]

	SimpleQS(a[:left])
	SimpleQS(a[left+1:])

	return a
}

// concurrentqs uses go routines to perform quick sort.
func concurrentqs(a []int, done chan int) {

	if len(a) < 2 {
		done <- 1
		return
	}
	left, right := 0, len(a)-1
	pivot := rand.Int() % len(a)

	// Swap pivot with right most element
	a[right], a[pivot] = a[pivot], a[right]

	for i := range a {
		if a[i] < a[right] {
			// Swap i'th position with left position
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	// Swap right position (previous pivot) with left position ((left -1) is smaller than right)
	a[right], a[left] = a[left], a[right]

	childChan := make(chan int)
	go concurrentqs(a[:left], childChan)
	go concurrentqs(a[left+1:], childChan)
	<-childChan
	<-childChan

	done <- 1
	return
}

// ConcurrentQS ...
func ConcurrentQS(s []int) {
	d := make(chan int)
	go concurrentqs(s, d)
	<-d
	return
}

func generateRandomSlice(length int) []int {
	slice := make([]int, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		slice[i] = rand.Intn(500)
	}
	return slice
}

func main() {
	rand.Seed(time.Now().UnixNano())

	randomSlice := generateRandomSlice(rand.Intn(20))
	fmt.Println("randomSlice : ", randomSlice)

	ConcurrentQS(randomSlice)
	fmt.Println("sortedSlice : ", randomSlice)

}
