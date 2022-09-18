package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// 1. 比较单个值和多个值
	fmt.Println("Go runs on")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Println(os)
	}

	// 2. 每个分支设置比较条件
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning")
	case t.Hour() < 17:
		fmt.Println("Good afternoon")
	default:
		fmt.Println("Good evening")
	}

	// 3. 使用fallthrough关键字
	t = time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning")
		fallthrough
	case t.Hour() < 17:
		fmt.Println("Good afternoon")
		fallthrough
	default:
		fmt.Println("Good evening")
	}

}
