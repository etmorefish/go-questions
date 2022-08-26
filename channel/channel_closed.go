package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	ch <- 18
	close(ch)
	x, ok := <-ch
	if ok {
		fmt.Println("received: ", x)
	} else {
		fmt.Println("received: ", x, ok)
	}

	// ch <- 90
	xx, ok := <-ch
	if !ok {
		fmt.Println("channel closed, data invalid.", xx)
	}
	// close(ch)

}
