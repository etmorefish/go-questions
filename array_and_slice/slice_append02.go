package main

import "fmt"

func main() {
	s := []int{1, 2}
	s = append(s, 4, 5, 6)
	fmt.Printf("len=%d. cap=%d\n", len(s), cap(s))
}

// len=5. cap=6