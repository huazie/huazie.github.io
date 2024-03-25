---
title: Go语言实战2-自定义OrderedMap
date: 2016-08-17 19:35:44
updated: 2024-03-04 11:11:40
categories:
  - [开发语言-Go,Go语言实战]
tags:
  - Go
  - OrderedMap
---

[《开发语言-Go》](/categories/开发语言-Go/) [《Go语言实战》](/categories/开发语言-Go/Go语言实战/) 

![](/images/go-logo.png)

# 引言

在Go语言中，字典类型的元素值的迭代顺序是不确定的。想要实现有固定顺序的Map就需要让自定义的 **OrderedMap** 实现 **sort.Interface** 接口类型。该接口类型中的方法 **Len** 、**Less** 和 **Swap** 的含义分别是获取元素的数量、比较相邻元素的大小以及交换它们的位置。


# 主要内容
## 1. 定义 OrderedMap

想要实现自定义一个有序字典类型，仅仅基于字典类型是不够的，需要使用一个元素有序的数据类型值作为辅助。如下声明了一个名为 **OrderedMap** 的结构体类型：

```go
type OrderedMap struct {
	keys 	[]interface{}
	m 		map[interface{}]interface{}
}
```


如上在 **OrderedMap** 类型中，除了一个字典类型的字段 **m**，还有一个切片类型的字段。

## 2. 实现 sort.Interface

这里需要添加为 **OrderedMap** 类型如下方法：

```go
// 获取键值对的数量
func (omap *OrderedMap) Len() int {
	return len(omap.keys)
}

func (omap *OrderedMap) Less(i, j, int) bool {
	// 省略若干条语句
}

func (omap *OrderedMap) Swap(i, j int) {
	omap.keys[i], omap.keys[j] = omap.keys[j], omap.keys[i];
}

```

如上 ***OrderMap** 类型（不是 **OrderedMap** 类型）就是一个 **sort.Interface** 接口类型的实现类型了。可以看出，**Len** 方法中是以 **keys** 字段的值的长度作为结果值。而在 **Swap** 方法中，使用平行赋值语句交换的两个元素值也是操作的 **keys** 字段的值。

现在考虑 **Less** 方法的实现，方法 **Less** 的功能是比较相邻的两个元素值得大小病返回判断结果。读者如果了解Go语言值的可比性与有序性（详细的各种数据类型的比较方法可以参考笔者的 [Go语言学习12-数据的使用](/2016/07/13/go/go-learning/go-learning12/) 中所述），就知道在Go语言中，只有当值的类型具备有序性的时候，它才可能与其他的同类型值比较大小。另外Go语言规定，字典类型的键类型的值必须是可比的（即可以判定两个该类型的值是否相等）。但是，在Go语言中具有可比性的数据类型中只有一部分同时具备有序性，因此只是依靠Go语言本身对字典类型的键类型的约束是不够的。在Go语言中，具备有序性的预定义数据类型只有 **整数类型**、**浮点数类型** 和 **字符串类型**。

类型 **OrderedMap** 中的字段 **keys** 是 **[]interface{}** 类型的，因此想要有序性，总是需要比较两个 **[]interface{}** 类型值的大小。因为接口类型的值只具有可比性而不具备有序性，所以上面声明的**OrderedMap** 类型在这里显得不适用了。如果将 **keys** 字段的元素类型改为某一个具体的数据类型（整数类型、浮点数类型和字符串类型），又会使得 **OrderedMap** 类型的使用很局限。

因此这里需要对 **keys** 字段进行详细的分析，如下需求列表：

 - 字段keys的值中的元素值应该都是有序的，应该可以方便地比较它们之间的大小。
 - 字段keys的值的元素类型不应该是一个具体的类型，应该可以在运行时再确定它的元素类型。
 - 字段keys的值应该可以方便地进行添加元素值、删除元素值以及获取元素值等操作。
 - 字段keys的值中的元素值应该可以被依照固定的顺序获取。
 - 字段keys的值中的元素值应该能够被自动地排序。
 - 字段keys的值总是已排序的，应该能够确定某一个元素值的具体位置。
 - 字段keys既然可以在运行时决定它的值的元素类型，那么就可以在运行时获知这个元素类型。
 - 字段keys的值中的不同元素值比较大小的具体方法，应该可以在运行时获知到。

## 3. 定义Keys接口类型

为了满足上面的需求列表，可以定义如下的名为 **Keys** 的接口类型：

```go
type Keys interface {
	sort.Interface
	Add(k interface{}) bool
	Remove(k interface{}) bool
	Clear()
	Get(index int) interface{}
	GetAll() []interface{}
	Search(k interface{})(index int, contains bool)
	ElemType() reflect.Type
	CompareFunc() func(interface{}, interface{}) int8
}
```

在 **Keys** 接口类型中嵌入了 **sort.Interface** 接口类型，也就是说 **Keys** 类型的值一定是可排序的。**Add**，**Remove**，**Clear** 和 **Get** 这4个方法使得可以对 **Keys** 的值进行添加，删除，清除和获取元素值的操作。**GetAll** 方法可以获取一个与 **Keys** 类型值有着相同元素值集合和元素迭代顺序的切片值。**Search** 方法确定某一个元素值的具体位置，**CompareFunc** 返回一个比较大小的具体方法，**ElemType** 方法返回一个 **reflect.Type** 类型的结果值。

> **注意：** 实际上，**reflect** 包中的程序实体提供了Go语言运行时的反射机制。通过这些程序实体，可以编写出一些代码来动态的操纵任意类型的对象（如 **TypeOf** 函数用于获取一个 **interface{类型的值的动态类型信息}**）。

## 4. Keys接口类型的实现类型

**Keys** 接口类型的定义并没有体现需求列表中的第**1 , 2 , 5**项所描述的功能。既然 **Keys** 接口类型的值必须是 **sort.Interface** 接口的一个实体，通过 **sort** 代码包中的程序实体实现元素自动排序的功能应该不难。

为了能够动态地决定元素类型，需要在这个 **Keys** 的实现类型中声明一个 **[]interface{}** 类型的字段，以作为存储被添加到 **Keys** 类型值中的元素值得底层数据结构：

```go
Container []interface{}
```

另外由于Go语言本身并没有对自定义泛型提供支持，因此需要这个字段的值存储某一个数据类型的元素值。但是接口类型的值不具备有序性（即不能比较大小）。尽管这样，我们可以让具体使用者去实现一个比较大小的方法，还需添加如下字段：

```go
compareFunc func(interface{}, interface{}) int8
```

这是一个函数类型的字段，这个函数返回一个 **int8** 类型的结果值，对结果值做出如下规定：
- 当第一个参数值小于第二个参数值时，结果值应该小于0
- 当第一个参数值大于第二个参数值时，结果值应该大于0
- 当第一个参数值等于第二个参数值时，结果值应该等于0

现在，通过将比较两个元素值大小的问题抛给使用者，既解决了需要动态确定元素类型的问题，又明确了比较两个元素值大小的解决方法。不过，由于 **container** 字段是 **[]interface{}** 类型的，常常不能够方便地在运行时获取到它的实际元素类型（比如在它的值中还没有任何元素值的时候）。这里需要一个明确 **container** 字段的实际元素类型的字段，这个字段的值所代表的类型也应该是当前的Keys类型值的实际元素类型。如下 **Keys** 接口类型的实现类型的声明如下：

```go
type myKeys struct {
	container 		[]interface{}
	compareFunc 	func(interface{}, interface{}) int8
	elemType 		reflect.Type
}
```


现在使用一个 ***myKeys** 类型的值来存储 **int64** 类型的元素值，应该如下来初始化它：

```go
int64Keys := &myKeys{
	container : make([]interface{}, 0),
	compareFunc : func(e1 interface{}, e2 interface{}) int8 {
		k1 := e1.(int64)
		k2 := e2.(int64)
		if k1 < k2 {
			return -1
		} else if k1 > k2 {
			return 1
		} else {
			return 0
		}
	},
	elemType : reflect.Typeof(int64(1))
}
```
> **注意：** **compareFunc** 字段的值中的那两个类型断言表达式的目标类型一定要与 **elemType** 字段的值所代表的类型保持一致。**elemType** 字段的值所代表的类型其实就是调用 **reflect.TypeOf** 函数时传入的那个参数值的类型，即 **int64**。


被用于实现 **sort.Interface** 接口类型的方法的声明如下：

```go
func (keys *myKeys) Len() int {
	return len(keys.container)
}

//该方法中，比较两个元素值的操作全权交给了compareFunc字段所代表的那个函数
func (keys *myKeys) Less(i, j int)bool {
	return keys.compareFunc(keys.container[i], keys.container[j]) == -1
}

func (keys *myKeys) Swap(i, j int) {
	keys.container[i], keys.container[j] = keys.container[j], keys.container[i];
}
```

如上这3个方法的接收者类型都是 ***myKeys**，所以事先 **sort.Interface** 接口类型的类型是 ***myKeys** 而不是 **myKeys**。

## 5. 实现Add方法

现在考虑实现 **Add** 方法，但在真正向字段 **container** 的值添加元素值之前，需要先判断这个元素值的类型是否符合要求。当然，这需要使用字段 **elemType** 的值，它代表了可接受的元素值的类型。现在使用一个独立的方法来实现这个判断，如下：

```go
func (keys *myKeys) isAcceptableElem(k interface{}) bool {
	if k == nil {
		return false
	}
	if reflect.TypeOf(k) != keys.elemType {
		return false
	}
	return true
}
```

在 **Add** 方法中，使用 **isAcceptableElem** 方法来判定元素值的类型是否可被接收。如果结果是否定的，直接返回 **false** ;如果结果是肯定的，就向 **container** 字段的值添加这个元素值。在添加之后，应该对 **container** 的值中的元素值进行一次排序。这需要用到 **sort** 代码包中的排序函数 **sort.Sort**，它的声明如下：

```go
func Sort(data Interface) {
	// 省略若干语句
}
```

函数 **sort.Sort** 的签名中的参数类型 **Interface** 其实就是接口类型 **sort.Interface**，并且这两个程序实体处在同一个代码包中。
>**知识点**：**sort.Sort** 函数使用的排序算法是一种由三向切分的快速排序算法，堆排序算法和插入排序算法组成的混合算法。虽然快速排序是最快的通用排序算法，但在元素值很少的情况下它比插入顺序要慢一些。而堆排序的空间复杂度是常数级别的，且它的时间复杂度在大多数情况下只略逊于其他两种排序算法，所以在快速排序中的递归达到一定深度的时候，切换至堆排序来节约空间是值得的。这样的算法组合使得**sort.Sort** 函数的时间复杂度在最坏的情况下是 **O(N*logN)** 的，并且能够有效地控制对空间的使用，但是不提供稳定性的保证（即在排序过程中不保留数组或切片值中重复元素的相对位置）。

现在实现 **Add** 方法，声明如下：

```go
func (keys *myKeys) Add(k interface{}) bool {
	ok := keys.isAcceptableElem(k)
	if !ok {
		return false
	}
	keys.container = append(keys.container, k)
	// sort.Sort函数会通过对keys的值的Len、Less和Swap方法的调用来完成排序。
	// 而在Less方法中，通过compareFunc函数对相邻的元素值进行比较的。
	sort.Sort(keys)
	return ture
}
```

## 6. 实现Search方法

现在考虑实现 **Remove** 方法，但在实现之前需要使用 **Search** 方法找到指定删除的元素所处的位置，因此先实现 **Search** 方法。在 **Search** 方法中，要搜索参数 **k** 代表的值在 **container** 中对应的索引值。由于 **k** 的类型是 **interface{}** 的，所以需要先使用 **isAcceptableElem** 方法对它进行判定，然后可以通过调用 **sort.Search** 函数来实现搜索元素值的核心逻辑。**sort.Search** 函数的声明如下：

```go
func Search(n int, f func(int) bool) int {
	// 省略若干条语句
}
```

由于 **sort.Search** 函数使用二分查找算法在切片值中搜索指定的元素值。该搜索算法有着稳定的 **O(logN)** 的时间复杂度，但它要求被搜索的数组或切片值必须是有序的，而这里在添加元素的时候已经保证了**container** 字段的值中的元素值是已被排序过的。

从上面的声明中可知，**sort.Search** 函数有两个参数。第一个参数接受的是欲排序的切片值得长度，而第二个参数接受的是一个函数值。这个函数值的含义是：对于一个给定的索引值，判定与之对应的元素值是否等于欲查找的元素值或者应该排在欲查找的元素值的右边。对于参数f的值如下：

```go
func(i int) bool {
	return keys.compareFunc(keys.container[i], k) >= 0
}
```

如上这个参数 **f** 应该怎样理解呢？这里先假设有这样一个切片值：

```go
[]int{1, 3, 5, 7, 9, 11, 13, 15}
```

现在要查找的元素值是 **7**，依据二分查找算法，**sort.Search** 函数内部会在第三次折半的时候使用 **7** 的索引值 **3** 作为函数f的参数值，函数f的结果值应该是 **true**，**sort.Search** 函数的执行会结束并返回 **7** 的索引值 **3** 作为它的结果值。但是，还有一种情况就是，要查找的元素根本就不在这个切片值里，比如 **6** 或者 **8** 等等，**sort.Search** 函数的执行也会在 **f(3)** 被求值之后结束，且它的结果值会是 **4** 或者 **3**，对应的元素都不是要查找的。

**sort.Search** 函数的结果值总会在 **[0, n]** 的范围内，但结果值并不一定就是欲查找的元素值所对应的索引值。因此，需要在调用 **sort.Search** 函数的结果值之后再进行一次判断,如下：

```go
if index < keys.Len() && keys.container[index] == k {
	contains = true
}
```

其中 **index** 代表了 **sort.Search** 函数的结果值，这里需要先检查结果值是否在有效的索引范围之内，然后还需要判断它所对应的元素值是否就是要查找的。经过这些分析之后，相信不难实现 ***myKeys** 类型的 **Search** 函数。

## 7. 实现Remove方法

实现了 **Search** 函数就来看看 **Remove** 函数，通过调用 ***myKeys** 类型的 **Search** 函数，就可以获取欲删除的元素值对应的索引值和它是否被包含在 **container** 中的判断结果。如果第二个结果值是 **false**，就可以直接返回 **false**，否则就从 **container** 中删除掉这个元素值。从切片值中删除一个元素值有很多种方法，比如使用 **for** 语句、**copy** 语句或 **append** 函数等等。这里选择用 **append** 函数来实现，因为它可以在不增加时间复杂度和空间复杂度的情况下使用更少的代码来完成功能，且不降低可读性。如下实现了删除一个元素值的功能：

```go
keys.container = append(keys.container[0: index],keys.container[index+1: ]…)
```

如上代码充分地使用了切片表达式和 **append** 函数，使用如下切片表达式：

```go
keys.container[0: index] // 取出container字段的值中的在欲删除元素值之前的子元素序列
keys.container[index+1: ] // 取出container字段的值中的在欲删除元素值之后的子元素序列
```

接着通过 **append** 函数将两个元素子序列拼接起来，可以在第二个参数值之后添加“…”以表示把第二个参数值中的每个元素值都作为传给 **append** 函数的独立参数，这样就把第二个子序列中的所有元素值逐个追加到了第一个子序列的尾部，最后把拼接后的元素序列赋值给了 **container** 字段。

## 8. 实现Clear方法

```go
func (keys *myKeys) clear() {
	keys.container = make([]interface{}, 0)
}
```

## 9. 实现Get方法

```go
func (keys *myKeys) Get(index int) interface{} {
	if index >= keys.Len() {
		return nil
	}
	return keys.container[index]
}

```

## 10. 实现GetAll方法

```go
func (keys *myKeys) GetAll() []interface{} {
	initialLen := len(keys.container)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for _, key := range keys.container {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}
```

## 11. 实现ElemType和CompareFunc方法

```
func (keys *myKeys) ElemType() reflect.Type {
	return keys.elemType
}

func (keys *myKeys) CompareFunc() CompareFunction {
	return keys.compareFunc
}
```

## 12. 实现String方法

```go
// String方法被用于生成可读性更好的接收者值的字符串表示形式。
func (keys *myKeys) String() string {
	var buf bytes.Buffer
	buf.WriteString("Keys<")
	buf.WriteString(keys.elemType.Kind().String())
	buf.WriteString(">{")
	first := true
	buf.WriteString("[")
	for _, key := range keys.container {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("]")
	buf.WriteString("}")
	return buf.String()
}
```

## 13.定义NewKeys函数

目前已经完成了 **myKeys** 类型以及相关方法的编写，还应该编写一个用于初始化 ***myKeys** 类型值的函数。这个函数名为 **NewKeys**，结果类型为 **Keys**，声明如下：

```go
func NewKeys ( compareFunc func(interface{}, interface{}) int8, elemType reflect.Type) Keys {
	// 因为只有*myKeys类型的方法集合中才包含了Keys接口类型中声明的所有方法，
	// 所以下面返回的是一个*myKeys类型值，而不是一个myKeys类型值。
	return &myKeys {
		container : 		make([]interface{}, 0),
		compareFunc : 	compareFunc,
		elemType : 		elemType,
	}
}
```

从上可以看出，在 **NewKeys** 函数的参数声明列表中没有与 **container** 字段相对应的参数声明，原因是 **container** 字段的值总应该是一个长度为 **0** 的 **[]interface{}** 类型值，它不必由 **NewKeys** 函数的调用方提供。另外，**NewKeys** 函数的**compareFunc** 参数和 **elemType** 参数之间的关系，也要满足之前上文提到的约束条件。

## 14. 重新理解OrderedMap

至此我们已经编写完成了 **OrderedMap** 类型所需要用到的最核心的数据类型 **Keys** 和 **myKeys**。
由于有了 **Keys** 接口类型，**OrderedMap** 类型的声明被修改为：

```go
type myOrderedMap struct {
	keys 		Keys
	elemType 	reflect.Type
	m			map[interface{}]interface{}
}
```

如上更改了该类型的名称，这里要声明一个接口类型来描述有序字典类型所提供的功能，而 **OrderedMap** 更适合作为这个接口类型的名称。声明如下：

```go
// 泛化的Map的接口类型
type OrderedMap interface {
	// 获取给定键值对应的元素值。若没有对应元素值则返回nil。
	Get(key interface{}) interface{}
	// 添加键值对，并返回与给定键值对应的旧的元素值。若没有旧元素值则返回(nil, true)。
	Put(key interface{}, elem interface{}) (interface{}, bool)
	// 删除与给定键值对应的键值对，并返回旧的元素值。若没有旧元素值则返回nil。
	Remove(key interface{}) interface{}
	// 清除所有的键值对。
	Clear()
	// 获取键值对的数量。
	Len() int
	// 判断是否包含给定的键值。
	Contains(key interface{}) bool
	// 获取已排序的键值所组成的切片值。
	Keys() []interface{}
	// 获取已排序的元素值所组成的切片值。
	Elems() []interface{}
	// 获取已包含的键值对所组成的字典值。
	ToMap() map[interface{}]interface{}
	// 获取键的类型。
	KeyType() reflect.Type
	// 获取元素的类型。
	ElemType() reflect.Type
}
```

这里要使 ***myOrderedMap** 类型成为 **OrderedMap** 接口类型的实现类型。虽然方法不多，但实现起来并不难，完成后还可以编写一个 **NewOrderedMap** 函数，用于将初始化好的 ***myOrderedMap** 类型值作为结果值返回。

# 参考

完整的 **myOrderedMap** 类型以及相关方法的实现，如下链接：

[https://github.com/hyper-carrot/goc2p/tree/master/src/basic/map1](https://github.com/hyper-carrot/goc2p/tree/master/src/basic/map1)



