package main

import "fmt"

func sum(s []int, ch chan int) {
	result := 0
	for _, e := range s {
		result += e
	}
	ch <- result
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}
	fmt.Println("Initial slice : ", s)

	ch := make(chan int)
	go sum(s[:len(s)/2], ch)
	left := <-ch
	fmt.Println("First half sum : ", left)

	go sum(s[len(s)/2:], ch)
	right := <-ch
	fmt.Println("Second half sum : ", right)

	final := left + right
	fmt.Println("Final sum : ", final)
}
