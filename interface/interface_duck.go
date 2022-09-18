package main

import "fmt"

type IGreeting interface {
	sayHello()
}

func sayHello(i IGreeting) {
	i.sayHello()
}

type Go struct{}

func (g Go) sayHello() {
	fmt.Println("Hi, i'm Go")
}

type Python struct{}

func (p Python) sayHello() {
	fmt.Println("Hi, i'm Python")
}

func main() {
	golang := Go{}
	python := Python{}

	sayHello(golang)
	sayHello(python)

}
