package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func initpool() {
	pool = &sync.Pool{New: func() interface{} {
		fmt.Println("Creating a new person...")
		return new(Person)
	},
	}
}
func main() {
	initpool()
	p := pool.Get().(*Person)
	fmt.Println("首次从poll中获取：", p)
	p.Name = "first"
	fmt.Println(p.Name)

	pool.Put(p)

	fmt.Println(pool.Get().(*Person))
	fmt.Println(pool.Get().(*Person).Name)
	fmt.Println("--------------------------------")
	fmt.Println(pool.Get().(*Person))

}
