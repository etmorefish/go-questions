package main

import "fmt"

func IsClosed(ch <-chan int) bool {
	select {
	case <-ch:
		fmt.Println("1")
		return true
	default:
		fmt.Println("2")
	}
	fmt.Println("3")

	return false
}

func main() {
	c := make(chan int)
	fmt.Println(IsClosed(c)) // false
	close(c)
	fmt.Println(IsClosed(c)) // true
}
