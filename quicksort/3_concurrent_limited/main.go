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

var MAXGOROUTINES = 10000

// concurrentqsLimited uses go routines to perform quick sort.
func concurrentqsLimited(a []int, done chan int, workers chan int) {
	// report to caller that we're finished
	if done != nil {
		defer func() { done <- 1 }()
	}

	if len(a) < 2 {
		return
	}
	// since we may use the doneChannel synchronously
	// we need to buffer it so the synchronous code will
	// continue executing and not block waiting for a read
	doneChannel := make(chan int, 1)

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

	select {
	case <-workers:
		// if we have spare workers, use a goroutine
		// for parallelization
		go concurrentqsLimited(a[:left], doneChannel, workers)
	default:
		// if no spare workers, sort synchronously
		concurrentqsLimited(a[:left], nil, workers)
		// calling this here as opposed to using the defer
		doneChannel <- 1
	}

	// use the existing goroutine to sort above the pivot
	concurrentqsLimited(a[left+1:], nil, workers)
	// if we used a goroutine we'll need to wait for
	// the async signal on this channel, if not there
	// will already be a value in the channel and it shouldn't block
	<-doneChannel
	return
}

// ConcurrentQSLimited ...
func ConcurrentQSLimited(s []int) {
	if len(s) <= 1 {
		return
	}
	workers := make(chan int, MAXGOROUTINES-1)
	for i := 0; i < (MAXGOROUTINES - 1); i++ {
		workers <- 1
	}
	concurrentqsLimited(s, nil, workers)
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

	randomSlice = generateRandomSlice(rand.Intn(20))
	fmt.Println("randomSlice : ", randomSlice)

	ConcurrentQSLimited(randomSlice)
	fmt.Println("sortedSlice : ", randomSlice)

}
