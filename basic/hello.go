package main

import "fmt"

func main() {
	fmt.Println("hello world!")
}

// go build -gcflags "-N -l" -o hello hello.go
// gdb hello -> info files
