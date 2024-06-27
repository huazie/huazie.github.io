---
title: Go语言学习8-接口类型
date: 2016-07-06 21:15:19
updated: 2024-03-19 17:21:57
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 接口类型
---



![](/images/go-logo.png)

# 引言
上篇我们了解了Go语言的[《函数类型》](../../../../../../2016/07/05/go/go-learning/go-learning7/)，本篇主要了解接口类型。主要如下：

<!-- more -->

# 主要内容

一个Go语言的接口由一个方法的集合代表。只要一个数据类型（或与其对应的指针类型）附带的方法集合是某一个接口的方法集合的**超集**，那么就可以判定该类型实现了这个接口。

## 1. 类型表示法

接口类型的声明由若干个方法的声明组成。方法的声明由**方法名称**和**方法签名**构成。在一个接口类型的声明中不允许出现重复的方法名称。

**接口类型**是所有自定义的接口类型的统称。以标准库代码包 **sort** 中的接口类型 **Interface** 为例，声明如下：

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(I, j int)
}
```
在Go语言中可以将一个接口类型嵌入到另一个接口类型中。如下接口类型声明：

```go
type Sortable interface {
    sort.Interface
    Sort()
}
```

如上接口类型 **Sortable** 实际包含了4个方法声明，分别是**Len**、**Less**、**Swap** 和 **Sort**。

Go语言并不提供典型的类型驱动的子类化方法，但是却利用这种嵌入的方式实现了同样的效果。类型嵌入同样体现了非嵌入式的风格，同样适用于下面要讲的结构体类型。

>**注意**：一个接口类型只接受其他接口类型的嵌入。

对于接口的嵌入，一个约束就是**不能嵌入自身**，包括**直接嵌入**和**间接嵌入**。

**直接嵌入**如下：

```go
type Interface1 interface {
    Interface1
}
```

而**间接嵌入**如下：

```go
type Interface2 interface {
    Interface3
}

type Interface3 interface {
    Interface2
}
```


错误的接口嵌入会造成编译错误。另外，当前接口类型中声明的方法也不能与任何被嵌入其中的接口类型的方法重名，否则也会造成编译错误。

至于Go语言的自身定义的一个特殊的接口类型----空类型 **interface{}**，前面也提到过，就是不包含任何方法声明的接口。并且，Go语言中所有数据类型都是它的实现。

## 2. 值表示法

Go语言的接口类型没有相应的值表示法，因为接口是规范而不是实现。但一个接口类型的变量可以被赋予任何实现了这个接口类型的数据类型的值，因此接口类型的值可以由任何实现了这个接口类型的其他数据类型的值来表示。

## 3. 属性和基本操作

接口的最基本属性就是它们的方法集合。

实现一个接口类型的可以是任何自定义的数据类型，只要这个数据类型附带的方法集合是该接口类型的方法集合的超集。编写一个自定义的数据类型 **SortableStrings** ,如下：

```go
type SortableStrings [3]string
```
如上这个自定义的数据类型相当于 **[3]string** 类型的一个别名类型。现在想让这个自定义数据类型实现 **sort.Interface** 接口类型，就需要实现**sort.Interface** 中声明的全部方法，这些方法的实现都需要以类型 **SortableStrings** 为接收者的类型。这些方法的声明如下：

```go
func (self SortableStrings) Len() int {
    return len(self)
}

func (self SortableStrings) Less(i, j int) bool {
    return self[i] < self[j]
}

func (self SortableStrings) Swap(i, j int) {
    self[i], self[j] = self[j], self[i]
}
```

有了上面三个方法的声明，**SortableStrings**类型就已经是一个**sort.Interface**接口类型的实现了。使用 [**Go语言学习2-基本词法** ](../../../../../../2016/06/27/go/go-learning/go-learning2/)中讲的**类型断言表达式验证**，编写代码如下：

```go
_, ok := interface{}(SortableStrings{}).(sort.Interface)
```

>**注意**: 想要让这条语句编译通过，首先需要导入标准代码包sort。

在如上赋值语句的右边是一个类型断言表达式，左边的两个标识符代表了这个表达式的求值结果。这里不关心转换后的结果，只关注类型转换是否成功，因此第一个标识符为**空标识符“_”**；第二个标识符 **ok** 代表了一个布尔类型的变量，**true** 表示转换成功，**false** 表示转换失败。如下图，显示 **ok** 的结果为 **true**，因为 **SortableStrings** 类型确实实现了接口类型 **sort.Interface** 中声明的所有方法。

![](result.png)


一个接口类型可以被任意数量的数据类型实现。一个数据类型也可以同时实现多个接口类型。

如上的自定义数据类型 **SortableStrings** 也可以实现接口类型 **Sortable**，如下再编写一个方法声明：

```go
func (self SortableStrings) Sort() {
    sort.Sort(self)
}
```

现在，**SortableStrings** 类型在实现了接口类型 **sort.Interface** 的同时也实现了接口类型 **Sortable**。类型断言表达式验证如下：

```go
_, ok2 := interface{}(SortableStrings{}).(Sortable)
```

ok2的结果为true，如下图：

![](result-1.png)

现在，把 **SortableStrings** 类型包含的 **Sort** 方法中的接收者类型由 **SortableStrings** 改为 ***SortableStrings**，如下：

```go
func (self *SortableStrings) Sort() {
    sort.Sort(self)
}
```

这个函数的接收者类型改为了与 **SortableStrings** 类型对应的指针类型。方法**Sort**不再是一个值方法了，已经变成了一个指针方法。只有与 **SortableStrings** 类型的值对应的指针值才能够通过上面的类型断言，如下：

```go
_, ok3 := interface{}(&SortableStrings{}).(Sortable)
```

这时 **ok2** 的值为 **false**，**ok3** 的值为 **true**，如下图：

![](result-2.png)

再添加如下测试代码:

```go
ss := SortableStrings("2", "3", "1")
ss.Sort()
fmt.Printf("Sortable strings: %v\n", ss)
```
以上出现的关于标准库代码包 **fmt** 的用法，大家可以参考：    http://docscn.studygolang.com/pkg/fmt

测试结果如下图：

![](result-3.png)

上面打印的信息中的 [2, 3, 1] 是 **SortableStrings** 类型值的字符串表示，从上面的结果可以看见，变量 **ss** 的值并没有排序，但在打印前已经调用了 **Sort** 方法。

下面且听解释：
: 上面讲到，在值方法中，对接收者的值的改变在该方法之外是不可见的。**SortableStrings** 类型的 **Sort** 方法实际上是通过函数 **sort.Sort** 来对接收者的值进行排序的。**sort.Sort** 函数接受一个类型为 **sort.Interface** 的参数值，并利用这个值的方法**Len**、**Less** 和 **Swap**来修改其参数中的各个元素的位置以完成排序工作。对于 **SortableStrings** 类型，虽然它实现了接口类型 **sort.Interface** 中声明的全部方法，但是这些方法都是值方法，从而这些方法中对接收者值的改变并不会影响到它的源值，只是改变了源值的复制品而已。

对于上面的问题，目前的解决方案是将 **SortableStrings** 类型的方法**Len**、**Less** 和 **Swap**的接收者类型都改为 ***SortableStrings**，如下图展示的运行结果：

![](result-4.png)

但这时的 **SortableStrings** 类型就不再是接口类型 **sort.Interface** 的实现，***SortableStrings** 才是接口类型 **sort.Interface** 的实现，如上图中 **ok** 的值为 **false**。

现在我们再考虑一种方案，对 **SortableStrings** 类型的声明稍作改动：

```go
type SortableStrings []string //去掉了方括号中的3
```

这个时候实际上是将 **SortableStrings** 有数组类型的别名类型改为了**切片类型**的别名类型，但是又使得现在与之相关的方法无法通过编译。主要的错误如下图：

![](result-5.png)

上面显示的主要错误有两个，一是内建函数 **len** 的参数不能是指向切片值的指针类型值；二是**索引表达式**不能被应用在指向切片值的指针类型值上。

下面对于此的解决方法就是将方法**Len**、**Less**、**Swap** 和 **Sort** 的接收者类型都由***SortableStrings**改回**SortableStrings**。这里是因为改动后的**SortableStrings**是切片类型，而**切片类型**是引用类型；对于**引用类型**来说，**值方法**对接收者值的改变也会反映在其源值上。如下图为修改过的结果：

![](result-6.png)

# 结语

对于上面的接口出现的代码，可以点击下载 [Go源码文件](http://download.csdn.net/download/u012855229/9565795)，自己修改修改，好好体会体会接口的用法。只需要在自己的工作区的src目录中的任意包中（这些包有意义即可）放入以下源码文件，进入命令行该文件目录输入上面的命令即可，当然首先你的Go语言环境变量要配好。

本篇就聊到这里，下篇继续未完的Go语言数据类型…