package main

import (
	"os"
	"runtime/trace"
)

var cache = map[interface{}]interface{}{}

// 预期能被快速释放的内存因被根对象引用而没有得到迅速释放
func keepalloc() {
	for i := 0; i < 10000; i++ {
		m := make([]byte, 1<<10)
		cache[i] = m
	}
}

// goroutine泄漏
func keepalloc2() {
	for i := 0; i < 10000; i++ {
		go func() {
			select {}
		}()
	}
}

var ch = make(chan struct{})

func keepalloc3() {
	for i := 0; i < 10000; i++ {
		go func() {
			// 没有接受方，goroutine会一直阻塞
			ch <- struct{}{}
		}()
	}
}
func main() {
	f, _ := os.Create("trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	keepalloc()
	keepalloc2()
}

// go tool trace trace.out
