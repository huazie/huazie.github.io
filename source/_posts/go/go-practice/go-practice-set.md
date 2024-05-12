---
title: Go语言实战1-自定义集合Set
date: 2016-07-31 07:43:58
updated: 2024-03-04 11:11:40
categories:
  - [开发语言-Go,Go语言实战]
tags:
  - Go
  - Set
  - HashSet
---

[《开发语言-Go》](/categories/开发语言-Go/) [《Go语言实战》](/categories/开发语言-Go/Go语言实战/) 

![](/images/go-logo.png)


# 引言

在Go语言中有作为 **Hash Table** 实现的字典（**Map**）类型，但标准数据类型中并没有集合（Set）这种数据类型。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

比较 **Set** 和 **Map** 的主要特性，有类似特性如下：

 - 它们中的元素都是不可重复的。
 
 - 它们都只能用迭代的方式取出其中的所有元素。
 
 - 对它们中的元素进行迭代的顺序都是与元素插入顺序无关的，同时也不保证任何有序性。

但是，它们之间也有一些区别，如下：

 - **Set** 的元素是一个单一的值，而 **Map** 的元素则是一个键值对。
 
 - **Set** 的元素不可重复指的是不能存在任意两个单一值相等的情况。Map的元素不可重复指的是任意两个键值对中的键的值不能相等。

从上面的特性可知，可以把集合类型（Set）作为字典类型（Map）的一个简化版本。也就是说，可以用 **Map** 来编写一个 **Set** 类型的实现。

实际上，在Java语言中，**java.util.HashSet** 类就是用 **java.util.HashMap** 类作为底层支持的。所以这里就从HashSet出发，逐步抽象出集合Set。
# 主要内容
## 1. 定义 HashSet

首先，在工作区的 **src** 目录的代码包 **basic/set**（可以自行定义，但后面要保持一致）中，创建一个名为 **hash_set.go** 的源码文件。

根据代码包 **basic/set** 可知，源码文件 **hash_set.go** 的包声明语句（这里可以查看[《Go语言学习1-基础入门》中的 1.5.1 包声明](/2016/06/27/go/go-learning/go-learning1/)）如下：

```go
package set
```

上面提到可以将集合类型作为字典类型的一个简化版本。现在我们的 **HashSet** 就以字典类型作为其底层的实现。**HashSet** 声明如下：

```go
type HashSet struct {
	m map[interface{}]bool
}
```

如上声明 **HashSet** 类型中的唯一的字段的类型是 **map[interface{}]bool**。选择这样一个字典类型是因为通过将字典 **m** 的键类型设置为 **interface{}**，让 **HashSet** 的元素可以是任何类型的，因为这里需要使用 **m** 的值中的键来存储 **HashSet** 类型的元素值。那使用 **bool** 类型作为 **m** 的值的元素类型的好处如下：

 - 从值的存储形式的角度看，**bool** 类型值只占用一个字节。
 
 - 从值的表示形式的角度看，**bool** 类型的值只有两个 --- **true** 和 **false**。并且，这两个值度都是预定义常量。

把 **bool** 类型作为值类型更有利于判断字典类型值中是否存在某个键。例如：如果在向 **m** 的值添加键值对的时候总是以 **true** 作为其中的元素的值，则索引表达式 ```m["a"]```的结果值总能体现出在m的值中是否包含键为 ```"a"```的键值对。对于 **map[interface{}]bool** 类型的值来说，如下：

```go
if m["a"] { // 判断是否m中包含键为“a”的键值对
	// 省略其他语句
}
```

如上 **HashSet** 类型的基本结构已确定了，现在考虑如何初始化 **HashSet** 类型值。由于字典类型的零值为 **nil** ，而用 **new** 函数来创建一个 **HashSet** 类型值，也就是 **new(HashSet).m** 的求值结果将会是一个 **nil** (关于 **new** 函数可以查阅本人另一篇博文 [Go语言学习14-内建函数](/2016/07/15/go/go-learning/go-learning14/))。因此，这里需要编写一个专门用于创建和初始化 **HashSet** 类型值的函数，该函数声明如下：

```go
func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}
```

如上可以看到，使用 **make** 函数对字段 **m** 进行了初始化(关于 **make** 函数也可查阅本人另一篇博文 [Go语言学习14-内建函数](/2016/07/15/go/go-learning/go-learning14/))。同时注意观察函数 **NewHashSet** 的结果声明的类型是 ***HashSet** 而不是 **HashSet**，目的是让这个结果值的方法集合中包含调用接收者类型为 **HashSet** 或 ***HashSet** 的所有方法。这样做的好处将在后面编写 **Set** 接口类型的时候再予以说明。

## 2. 实现HashSet的基本功能

依据其他编程语言中的 **HashSet** 类型可知，它们大部分应该提供的基本功能如下：

 - 添加元素值。
 
 - 删除元素值。
 
 - 清除所有元素值。
 
 - 判断是否包含某个元素值。
 
 - 获取元素值的数量。
 
 - 判断与其他HashSet类型值是否相同。
 
 - 获取所有元素值，即生成可迭代的快照。
 
 - 获取自身的字符串表示形式。

现在对这些功能一一实现，读者可自行实现，以下仅供参考。

### 2.1 添加元素值

```go
// 方法Add会返回一个bool类型的结果值，以表示添加元素值的操作是否成功。
// 方法Add的声明中的接收者类型是*HashSet。
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] { // 当前的m的值中还未包含以e的值为键的键值对
		set.m[e] = true // 将键为e(代表的值)、元素为true的键值对添加到m的值当中
		return true // 添加成功
	}
	return false // 添加失败
}
```

这里使用 ***HashSet** 而不是 **HashSet**，主要是从节约内存空间的角度出发，分析如下：

 - 当 **Add** 方法的接收者类型为 **HashSet** 的时候，对它的每一次调用都需要对当前 **HashSet** 类型值进行一次复制。虽然在 **HashSet** 类型中只有一个引用类型的字段，但是这也是一种开销。而且这里还没有考虑 **HashSet** 类型中的字段可能会变得更多的情况。
 - 当 **Add** 方法的接收者类型为 ***HashSet** 的时候，对它进行调用时复制的当前 ***HashSet** 的类型值只是一个指针值。在大多数情况下，一个指针值占用的内存空间总会被它指向的那个其他类型的值所占用的内存空间小。无论一个指针值指向的那个其他类型值所需的内存空间有多么大，它所占用的内存空间总是不变的。


### 2.2 删除元素值

```go
// 调用delete内建函数删除HashSet内部支持的字典值
func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e) // 第一个参数为目标字典类型，第二个参数为要删除的那个键值对的键
}
```

### 2.3 清除所有元素

```go
// 为HashSet中的字段m重新赋值
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}
```

如果接收者类型是 **HashSet**，该方法中的赋值语句的作用只是为当前值的某个复制品中的字段m赋值而已，而当前值中的字段 **m** 则不会被重新赋值。方法 **Clear** 中的这条赋值语句被执行之后，当前的 **HashSet** 类型值中的元素就相当于被清空了。已经与字段 **m** 解除绑定的那个旧的字典值由于不再与任何程序实体存在绑定关系而成为了无用的数据。它会在之后的某一时刻被Go语言的垃圾回收器发现并回收。

### 2.4 判断是否包含某个元素值

```go
// 方法Contains用于判断其值是否包含某个元素值。
// 这里判断结果得益于元素类型为bool的字段m
func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}
```

当把一个 **interface{}** 类型值作为键添加到一个字典值的时候，Go语言会先获取这个 **interface{}** 类型值的实际类型（即动态类型），然后再使用与之对应的 **hash** 函数对该值进行 **hash** 运算，也就是说，**interface{}** 类型值总是能够被正确地计算出 **hash** 值。但是字典类型的键**不能**是函数类型、字典类型或切片类型，否则会引发一个运行时恐慌，并提示如下：
**panic: runtime error: hash of unhashable type <某个函数类型、字典类型或切片类型的名称>**


### 2.5 获取元素值的数量

```go
// 方法Len用于获取HashSet元素值数量
func (set *HashSet) Len() int {
	return len(set.m)
}
```

### 2.6 判断与其他HashSet类型值是否相同

```go
// 方法Same用来判断两个HashSet类型值是否相同
func (set *HashSet) Same(other *HashSet) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}
```

两个 **HashSet** 类型值相同的必要条件是，它们包含的元素应该是完全相同的。由于 **HashSet** 类型值中的元素的迭代顺序总是不确定的，所以也就不用在意两个值在这方面是否一致。如果要判断两个 **HashSet** 类型值是否是同一个值，就需要利用指针运算进行内存地址的比较。


### 2.7 获取所有元素值，即生成可迭代的快照

所谓 **快照**，就是目标值在某一个时刻的映像。对于一个 **HashSet** 类型值来说，它的快照中的元素迭代顺序总是可以确定的，快照只反映了该 **HashSet** 类型值在某一个时刻的状态。另外，还需要从元素可迭代且顺序可确定的数据类型中选取一个作为快照的类型。这个类型必须是以单值作为元素的，所以字典类型最先别排除。又由于 **HashSet** 类型值中的元素数量总是不固定的，所以无法用一个数组类型的值来表示它的快照。如上分析可知，Go语言中可以使用的快照的类型应该是一个切片类型或者通道类型。

```go
// 方法Elements用于生成快照
func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m) // 获取HashSet中字段m的长度，即m中包含元素的数量
	// 初始化一个[]interface{}类型的变量snapshot来存储m的值中的元素值
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	// 按照既定顺序将迭代值设置到快照值(变量snapshot的值)的指定元素位置上,这一过程并不会创建任何新值。
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else { // m的值中的元素数量有所增加，使得实际迭代的次数大于先前初始化的快照值的长度
			snapshot = append(snapshot, key) // 使用append函数向快照值追加元素值。
		}
		actualLen++ // 实际迭代的次数
	}
	// 对于已被初始化的[]interface{}类型的切片值来说，未被显示初始化的元素位置上的值均为nil。
	// m的值中的元素数量有所减少，使得实际迭代的次数小于先前初始化的快照值的长度。
	// 这样快照值的尾部存在若干个没有任何意义的值为nil的元素，
	// 可以通过snapshot = snapshot[:actualLen]将无用的元素值从快照值中去掉。
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}
```


>**注意：** 在 **Elements** 方法中针对并发访问和修改 **m** 的值的情况采取了一些措施。但是由于m的值本身并不是并发安全的，所以并不能保证 **Elements** 方法的执行总会准确无误。要做到真正的并发安全，还需要一些辅助的手段，比如读写互斥量。


### 2.8 获取自身的字符串表示形式

```go
// 这个String方法的签名算是一个惯用法。
// 代码包fmt中的打印函数总会使用参数值附带的具有如此签名的String方法的结果值作为该参数值的字符串表示形式。
func (set *HashSet) String() string {
	var buf bytes.Buffer // 作为结果值的缓冲区
	buf.WriteString("HashSet{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	//n := 1
	//for key := range set.m {
	//	buf.WriteString(fmt.Sprintf("%v", key))
	//	if n == len(set.m) {//最后一个元素的后面不添加逗号
	//		break;
	//	} else {
	//		buf.WriteString(",")
	//	}
	//	n++;
	//}
	buf.WriteString("}")
	return buf.String()  
}
```

如上已经完整地编写了一个具备常用功能的Set的实现类型，后面将讲解更多的高级功能来完善它。

## 3. 高级功能

集合 **Set** 的真包含的判断功能。根据集合代数中的描述，如果集合 **A** 真包含了集合 **B**，那么就可以说集合 **A** 是集合 **B** 的一个超集。
```go
// 判断集合 set 是否是集合 other 的超集  
func (set *HashSet) IsSuperset(other *HashSet) bool {
	if other == nil {//如果other为nil，则other不是set的子集
		return false
	}
	setLen := set.Len()//获取set的元素值数量
	otherLen := other.Len()//获取other的元素值数量
	if setLen == 0 || setLen == otherLen {//set的元素值数量等于0或者等于other的元素数量
		return false
	}
	if setLen > 0 && otherLen == 0 {//other为元素数量为0，set元素数量大于0，则set也是other的超集
		return true
	}
	for _, v := range other.Elements() {
		if !set.Contains(v) {//只要set中有一个包含other中的数据，就返回false
			return false
		}
	}
	return true
}

```

集合的运算包括**并集**、**交集**、**差集** 和 **对称差集**。
**并集运算** 是指把两个集合中的所有元素都合并起来并组合成一个集合。
**交集运算** 是指找到两个集合中共有的元素并把它们组成一个集合。
集合 **A** 对集合 **B** 进行**差集运算**的含义是找到只存在于集合 **A** 中但不存在于集合 **B** 中的元素并把它们组成一个集合。
对称差集运算与差集运算类似但有所区别。对称差集运算是指找到只存在于集合 **A** 中但不存在于集合 **B** 中的元素，再找到只存在于集合 **B** 中但不存在于集合 **A** 中的元素，最后把它们合并起来并组成一个集合。


### 3.1 实现并集运算
```go
// 生成集合 set 和集合 other 的并集
func (set *HashSet) Union(other *HashSet) *HashSet {
	if set == nil || other == nil {// set和other都为nil，则它们的并集为nil
		return nil
	}
	unionedSet := NewHashSet()//新创建一个HashSet类型值，它的长度为0，即元素数量为0
	for _, v := range set.Elements() {//将set中的元素添加到unionedSet中
		unionedSet.Add(v)
	}
	if other.Len() == 0 {
		return unionedSet
	}
	for _, v := range other.Elements() {//将other中的元素添加到unionedSet中，如果遇到相同，则不添加（在Add方法逻辑中体现）
		unionedSet.Add(v)
	}
	return unionedSet
}
```

### 3.2 实现交集运算

```go
// 生成集合 set 和集合 other 的交集
func (set *HashSet) Intersect(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的交集为nil
		return nil
	}
	intersectedSet := NewHashSet() // 新创建一个HashSet类型值，它的长度为0，即元素数量为0
	if other.Len() == 0 { // other的元素数量为0，直接返回intersectedSet
		return intersectedSet
	}
	if set.Len() < other.Len() { // set的元素数量少于other的元素数量
		for _, v := range set.Elements() { // 遍历set
			if other.Contains(v) { // 只要将set和other共有的添加到intersectedSet
				intersectedSet.Add(v)
			}
		}
	} else { // set的元素数量多于other的元素数量
		for _, v := range other.Elements() { // 遍历other
			if set.Contains(v) { // 只要将set和other共有的添加到intersectedSet
				intersectedSet.Add(v)
			}
		}
	}
	return intersectedSet
}
```

### 3.3 差集

```go
// 生成集合 set 对集合 other 的差集
func (set *HashSet) Difference(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的差集为nil
		return nil
	}
	differencedSet := NewHashSet() // 新创建一个HashSet类型值，它的长度为0，即元素数量为0
	if other.Len() == 0 { // 如果other的元素数量为0
		for _, v := range set.Elements() { // 遍历set，并将set中的元素v添加到differencedSet
			differencedSet.Add(v)
		}
		return differencedSet // 直接返回differencedSet
	}
	for _, v := range set.Elements() { // other的元素数量不为0，遍历set
		if !other.Contains(v) { // 如果other中不包含v，就将v添加到differencedSet中
			differencedSet.Add(v)
		}
	}
	return differencedSet
}

```


### 3.4 对称差集

```go
// 生成集合 one 和集合 other 的对称差集
func (set *HashSet) SymmetricDifference(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的对称差集为nil
		return nil
	}
	diffA := set.Difference(other) // 生成集合 set 对集合 other 的差集
	if other.Len() == 0 { // 如果other的元素数量等于0，那么other对集合set的差集为空，则直接返回diffA
		return diffA
	}
	diffB := other.Difference(set) // 生成集合 other 对集合 set 的差集
	return diffA.Union(diffB) // 返回集合 diffA 和集合 diffB 的并集
}

```



## 4. 进一步重构

目前所实现的 **HashSet** 类型提供了一些必要的集合操作功能，但是不同应用场景下可能会需要使用功能更加丰富的集合类型。当有多个集合类型的时候，应该在它们之上抽取出一个接口类型以标识它们共有的行为方式。依据 **HashSet** 类型的声明，可以如下声明 **Set** 接口类型：

```go
type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
}
```

**注意：** **Set** 中的 **Same** 方法的签名与附属于 **HashSet**类型的 **Same** 方法有所不同。这里不能再接口类型的方法的签名中包含它的实现类型。因此这里的改动如下：

```go
func (set *HashSet) Same(other Set) bool {
	//省略若干语句
}
```

修改了 **Same** 方法的签名，目的是让 ***HashSet** 类型成为 **Set** 接口类型的一个实现类型。


高级功能的方法应该适用于所有的实现类型，完全可以抽离出成为独立的函数。并且，也不应该在每个实现类型中重复地实现这些高级方法。如下为改造后的 **IsSuperset** 方法的声明：

```go
// 判断集合 one 是否是集合 other 的超集
// 读者应重点关注IsSuperset与附属于HashSet类型的IsSuperset方法的区别
func IsSuperset(one Set, other Set) bool {
	if one == nil || other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}
```

# 总结

自定义集合 **Set** 的内容，基本已介绍完毕。大家可以试着对上面附属于 **HashSet** 类型的高级方法进行修改，以实现更完善的集合 **Set**。


