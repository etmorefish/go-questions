package main

func main() {
	// 无缓冲通道
	ch1 := make(chan int)
	// 有缓冲通道
	ch2 := make(chan int, 10)
}
