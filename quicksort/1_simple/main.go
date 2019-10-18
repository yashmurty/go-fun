package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateRandomSlice(length int) []int {
	slice := make([]int, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		slice[i] = rand.Intn(500)
	}
	return slice
}

func quicksort(a []int) []int {
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

	quicksort(a[:left])
	quicksort(a[left+1:])

	return a
}

func main() {
	rand.Seed(time.Now().UnixNano())

	randomSlice := generateRandomSlice(rand.Intn(20))
	fmt.Println("randomSlice : ", randomSlice)

	sortedSlice := quicksort(randomSlice)
	fmt.Println("sortedSlice : ", sortedSlice)
}
