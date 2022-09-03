package main

import "fmt"

func main() {
	x := foo()
	fmt.Println(*x)
}

func foo() *int {
	t := 3
	return &t
}

// go build -gcflags '-m -l' main.go
