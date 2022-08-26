# map的实现原理

两个关键点：map 是由 `key-value` 对组成的；`key` 只会出现一次。

和 map 相关的操作主要是：

1. 增加一个 k-v 对 —— Add or insert；
2. 删除一个 k-v 对 —— Remove or delete；
3. 修改某个 k 对应的 v —— Reassign；
4. 查询某个 k 对应的 v —— Lookup；

简单说就是最基本的 `增删查改`。

map 的设计也被称为 “The dictionary problem”，它的任务是设计一种数据结构用来维护一个集合的数据，并且可以同时对集合进行增删查改的操作。最主要的数据结构有两种：`哈希查找表（Hash table）`、`搜索树（Search tree）`。

哈希查找表用一个哈希函数将 key 分配到不同的桶（bucket，也就是数组的不同 index）。这样，开销主要在哈希函数的计算以及数组的常数访问时间。在很多场景下，哈希查找表的性能很高。

哈希查找表一般会存在“碰撞”的问题，就是说不同的 key 被哈希到了同一个 bucket。一般有两种应对方法：`链表法`和`开放地址法`。`链表法`将一个 bucket 实现成一个链表，落在同一个 bucket 中的 key 都会插入这个链表。`开放地址法`则是碰撞发生后，通过一定的规律，在数组的后面挑选“空位”，用来放置新的 key。

搜索树法一般采用自平衡搜索树，包括：AVL 树，红黑树。

自平衡搜索树法的最差搜索效率是 O(logN)，而哈希查找表最差是 O(N)。当然，哈希查找表的平均查找效率是 O(1)，如果哈希函数设计的很好，最坏的情况基本不会出现。还有一点，遍历自平衡搜索树，返回的 key 序列，一般会按照从小到大的顺序；而哈希查找表则是乱序的。

# map 的底层如何实现 
前面说了 map 实现的几种方案，Go 语言采用的是哈希查找表，并且使用链表解决哈希冲突。

## map 内存模型 
在源码中，表示 map 的结构体是 hmap，它是 hashmap 的“缩写”：

```go
type hmap struct {
	count     int    // 当前 map 中的元素个数，用于 len() 操作。
	flags     uint8  // 用于记录 map 当前状态，如是否正在执行写操作，后面会具体介绍。
	B         uint8  // buckets 的对数 log_2, 该值用于控制常规桶数组的长度。
	noverflow uint16 // overflow 的 bucket 近似数 
	hash0     uint32 // 计算 key 的哈希的时候会传入哈希函数 hash seed

	buckets    unsafe.Pointer // 指向当前桶数组地址的指针，数组长度 2^B ,如果元素个数为0，就为 nil。
	oldbuckets unsafe.Pointer // 用于扩容过程中，保存扩容旧的桶数组的指针，仅在扩容阶段非 nil。
	nevacuate  uintptr        // 用于在扩容过程中记录迁移进度，小于这个数的桶索引是已经迁移完成的。
  
	extra *mapextra // 保存了 map 的一些可选数据。
}

// 保存了 map 的一些可选数据，不是每个 map 都有该数据，目前主要保存了溢出桶相关数据。
// bmap 是桶的结构，下面会再单独介绍。
type mapextra struct {
  // 下面 2 个是比较特别的指针数组，和 GC 相关，仅在 bmap 不包含指针时候使用，溢出桶部分会详细介绍。
	overflow    *[]*bmap // 存储 buckets 用到的溢出桶指针集合。
	oldoverflow *[]*bmap // 存储 oldbuckets 用到的溢出桶指针集合。
  
	nextOverflow *bmap // 存储下一个可用的溢出桶指针。
}

// 编译时确定的 map 相关元数据。
type maptype struct {
	typ    _type  // map 自身的类型元数据。
	key    *_type // map key 的类型元数据。
	elem   *_type // map value 的类型元数据。
	bucket *_type // bucket 的类型元数据。

	hasher     func(unsafe.Pointer, uintptr) uintptr // key 的 hash 函数。
	keysize    uint8  // key 结构的大小。
	elemsize   uint8  // value 结构的大小。
	bucketsize uint16 // bucket 结构的大小。
	flags      uint32 // 用于标识 key/value 的行为，如是否被转为指针了。
}
```

#### 桶的数量

上面源码的注释其实已经提到了，桶（不包含溢出桶）的数量由 `hmap.B` 控制，具体值为 `2^B`，计算逻辑如下：

```go
// 传入的参数是目标 map 的 B 值。
func bucketShift(b uint8) uintptr {
	// 与操作主要是为了防止溢出：goarch.PtrSize 是个常量，表示该平台上一个指针的字节数，
	// 所以 goarch.PtrSize * 8 就是一个指针的位数，-1 后就是个二进制全 1 的值，
	// 也是该平台位操作能移动的最大位移。
	return uintptr(1) << (b & (goarch.PtrSize*8 - 1))
} 
```

#### hash 算法

hash 算法是在编译时候确定并赋值给 `maptype.hasher` 的，不同类型会由编译器直接决定一个合适的算法，覆盖类型全面的算法可以参考：`src/runtime/alg.go:typehash`。

另外，我们再仔细看下 `hash` 函数的签名，除了元素的指针以外，还需要一个 `uintptr`，这个参数其实是个 `hash seed`。

```go
func(unsafe.Pointer, uintptr) uintptr // hash 函数签名。
复制代码
```

map 中使用的 `hash seed` 保存在 `hmap.hash0`，这个字段是个随机数，会在 map 初始化及每次元素清零时候重新生成。因此每个 map 会有自己独立的哈希结果，也就有不一样的元素分布。这样以预防利用哈希碰撞做的 DoS 攻击，假设哈希算法稳定不变，容易被攻击者利用计时攻击等方式刻意制造出哈希碰撞，用于 DoS 攻击；

#### flags 的位含义

`flags` 用于记录 map 的状态，目前主要记录了 map 是否正在 `写入`、`迭代`、`扩容` 等，这里暂且留个印象，后续不理解再回来参考。

```go
const (
	// flags
	iterator     = 1 // 表明有迭代器正在使用 buckets。
	oldIterator  = 2 // 表明有迭代器正在使用 oldbuckets。
	hashWriting  = 4 // 表明有 goroutine 正在做写操作。
	sameSizeGrow = 8 // 表明当前正在进行等量扩容。
)
```

### 桶的结构

```go
const (
	// bucketCnt 定义了每个桶中存储的元素个数，就是常量 8 个。
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits

	// 定义了 map 中 Key 和 Value 能存放的大小上限，单位字节，超过的话会为被转为存储元素的指针。
	maxKeySize  = 128
	maxElemSize = 128
)

// 桶的数据结构。
type bmap struct {
	// 槽位状态数组，槽位无数据时用于快速判断槽位状态，
	// 槽位有数据时则存储了 hash(key) 的最高字节。
	tophash [bucketCnt]uint8
```

根据常量的定义，每个桶中能存放 `8 (bucketCnt)` 个元素。可以看到 `bmap` 的结构定义中仅有一个 `tophash` 字段，这是一个长度为 8 的 uint8 数组，保存了桶内 8 个槽位的信息，用于快速判断槽位状态。

- 槽位无数据时：对应的 `tophash` 保存了几个特殊的状态值，用于快速判断槽位状态；
- 槽位有数据时：对应的 `tophash` 则存储了 `hash(key)` 的最高字节；

下面列出了无数据时 tophash 的状态常量，后面阅读源码时候再回来参考，现在有个印象就可以了。

```go
const (
	// tophashs
	emptyRest      = 0 // 代表对应的 key/value 为空，且本元素之后的位置都是空的，无需继续寻找，是槽位的最始状态。
	emptyOne       = 1 // 代表对应的 key/value 为空，但曾经有过值，是个碎片空位。
	evacuatedX     = 2 // 元素被迁移到新桶中的 x（高位）桶了。
	evacuatedY     = 3 // 元素被迁移到新桶中的 y（低位）桶了。
	evacuatedEmpty = 4 // 表明该位置为空，且走过迁移逻辑了。
	minTopHash     = 5 // tophash 的最小值，计算出来如果小于这个数，会直接加上这个数，主要是用于避免和前面这几个特殊值冲突。
)
```



【引申1】slice 和 map 分别作为函数参数时有什么区别？

注意，这个函数返回的结果：`*hmap`，它是一个指针，而我们之前讲过的 `makeslice` 函数返回的是 `Slice` 结构体：

```go
// runtime/slice.go
type slice struct {
    array unsafe.Pointer // 元素指针
    len   int // 长度 
    cap   int // 容量
}
```

结构体内部包含底层的数据指针。

makemap 和 makeslice 的区别，带来一个不同点：当 map 和 slice 作为函数参数时，在函数参数内部对 map 的操作会影响 map 自身；而对 slice 却不会（之前讲 slice 的文章里有讲过）。

主要原因：一个是指针（`*hmap`），一个是结构体（`slice`）。Go 语言中的函数传参都是值传递，在函数内部，参数会被 copy 到本地。`*hmap`指针 copy 完之后，仍然指向同一个 map，因此函数内部对 map 的操作会影响实参。而 slice 被 copy 后，会成为一个新的 slice，对它进行的操作不会影响到实参。