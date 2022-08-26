# 如何比较两个map相等

map 深度相等的条件：

```
1、都为 nil
2、非空、长度相等，指向同一个 map 实体对象
3、相应的 key 指向的 value “深度”相等
```

```go
package main

import "fmt"

func main() {
	var m map[string]int
	var n map[string]int

	fmt.Println(m == nil)
	fmt.Println(n == nil)

	// 不能通过编译
	fmt.Println(m == n)
	// ./map_equal.go:13:14: invalid operation: m == n (map can only be compared to nil)
}
```

因此只能是遍历map 的每个元素，比较元素是否都是深度相等。
