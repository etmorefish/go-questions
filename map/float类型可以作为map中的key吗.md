# float类型可以作为map中的key吗

从语法上看，是可以的。Go 语言中只要是可比较的类型都可以作为 key。除开 slice，map，functions 这几种类型，其他类型都是 OK 的。具体包括：布尔值、数字、字符串、指针、通道、接口类型、结构体、只包含上述类型的数组。这些类型的共同特征是支持 `==` 和 `!=` 操作符，`k1 == k2` 时，可认为 k1 和 k2 是同一个 key。如果是结构体，只有 hash 后的值相等以及字面值相等，才被认为是相同的 key。很多字面值相等的，hash出来的值不一定相等，比如引用。

顺便说一句，任何类型都可以作为 value，包括 map 类型。

```go
func main() {
	m := make(map[float64]int)
	m[1.4] = 1
	m[2.4] = 2
	m[math.NaN()] = 3
	m[math.NaN()] = 3

	for k, v := range m {
		fmt.Printf("[%v, %d] ", k, v)
	}

	fmt.Printf("\nk: %v, v: %d\n", math.NaN(), m[math.NaN()])
	fmt.Printf("k: %v, v: %d\n", 2.400000000001, m[2.400000000001])
	fmt.Printf("k: %v, v: %d\n", 2.4000000000000000000000001, m[2.4000000000000000000000001])

	fmt.Println(math.NaN() == math.NaN())
}


/*
[1.4, 1] [2.4, 2] [NaN, 3] [NaN, 3] 
k: NaN, v: 0
k: 2.400000000001, v: 0
k: 2.4, v: 2
false
*/
```

NAN != NAN

当用 float64 作为 key 的时候，先要将其转成 uint64 类型，再插入 key 中

```go
func main() {
	m := make(map[float64]int)
	m[2.4] = 2

    fmt.Println(math.Float64bits(2.4))
	fmt.Println(math.Float64bits(2.400000000001))
	fmt.Println(math.Float64bits(2.4000000000000000000000001))
}

/*
4612586738352862003
4612586738352864255
4612586738352862003
*/
```

`2.4` 和 `2.4000000000000000000000001` 经过 `math.Float64bits()` 函数转换后的结果是一样的。自然，二者在 map 看来，就是同一个 key 了。
