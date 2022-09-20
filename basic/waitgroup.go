package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println(" before down")
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			wg.Wait()
			fmt.Println(i, "wait down")
		}(i)
	}
	go func() {
		time.Sleep(time.Second * 3)
	}()
	select {}
	// wg.Wait()

}
