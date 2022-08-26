# 如何实现两种get

Go 语言中读取 map 有两种语法：带 comma 和 不带 comma。当要查询的 key 不在 map 里，带 comma 的用法会返回一个 bool 型变量提示 key 是否在 map 中；而不带 comma 的语句则会返回一个 key 对应 value 类型的零值。如果 value 是 int 型就会返回 0，如果 value 是 string 类型，就会返回空字符串。

```go
package main

import "fmt"

func main() {
	ageMap := make(map[string]int)
	ageMap["qcrao"] = 18

    // 不带 comma 用法
	age1 := ageMap["stefno"]
	fmt.Println(age1)

    // 带 comma 用法
	age2, ok := ageMap["stefno"]
	fmt.Println(age2, ok)
}

// 0
// 0 false
```

两种语法对应到底层两个不同的函数。

```go
// src/runtime/map.go
func mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer
func mapaccess2(t *maptype, h *hmap, key unsafe.Pointer) (unsafe.Pointer, bool)
```

