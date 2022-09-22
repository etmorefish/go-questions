package main

import (
	"fmt"
	"sort"
)

func main() {
	m := map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}
	fmt.Println("first range  ")
	for i, v := range m {
		fmt.Printf("m[%v]=%v ", i, v)
	}
	fmt.Println("\nsecond range  ")
	for i, v := range m {
		fmt.Printf("m[%v]=%v ", i, v)
	}

	// 实现有序遍历
	// 把k单独取出来放到切片
	var sl []int
	for k := range m {
		sl = append(sl, k)
	}
	sort.Ints(sl)
	for _, k := range sl {
		fmt.Println(k, m[k])
	}
}
