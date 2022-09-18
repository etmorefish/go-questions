package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

/*
func main() {
	var i interface{} = new(Student)
	s := i.(Student)

	fmt.Println(s)
	// panic: interface conversion: interface {} is *main.Student, not main.Student
}
*/

// 这里直接发生了 panic，线上代码可能并不适合这样做，可以采用“安全断言”的语法：
func main() {
	var i interface{} = new(Student)
	s, ok := i.(Student)
	if ok {
		fmt.Println(s)
	}
}
