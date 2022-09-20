package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	// 1. write
	m.Store("xxml", 18)
	m.Store("phil", 16)
	// 2. read
	age, _ := m.Load("xxml")
	fmt.Println("age: ", age)
	// 3. range
	m.Range(func(key, value interface{}) bool {
		name := key.(string)
		age := value.(int)
		fmt.Println(name, age)
		return true
	})
	// 4. delete
	m.Delete("xxml")
	age, ok := m.Load("xxml")
	fmt.Println(age, ok)
	// 5. read or write
	m.LoadOrStore("apple", 10)
	age, _ = m.Load("phil")
	fmt.Println("age: ", age)

}
