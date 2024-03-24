---
title: Go语言学习9-结构体类型
date: 2016-07-07 23:20:45
updated: 2024-03-19 17:21:57
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 函数类型
---

[《开发语言-Go》](/categories/开发语言-Go/) [《Go语言学习》](/categories/开发语言-Go/Go语言学习/) 

![](/images/go-logo.png)

# 引言

书接上篇，我们了解了Go语言的[《接口类型》](/2016/07/06/go/go-learning/go-learning8/)，现在介绍Go语言的结构体类型。主要如下：

# 主要内容

结构体类型既可以包含若干个命名元素（又称字段），又可以与若干个方法相关联。

## 1. 类型表示法

结构体类型的声明可以包含若干个字段的声明。字段声明左边的标识符表示了该字段的名称，右边的标识符代表了该字段的类型，这两个标识符之间用空格分隔。

结构体类型声明中的每个字段声明都独占一行。同一个结构体类型声明中的字段不能出现重名的情况。

结构体类型也分为**命名结构体类型**和**匿名结构体类型**。

**命名结构体类型**

命名结构体类型以关键字**type**开始，依次包含结构体类型的**名称**、关键字**struct**和由**花括号**括起来的**字段声明列表**。如下:

```go
type Sequence struct {
    len int
    cap int
    Sortable
    sortableArray  sort.Interface
}
```
结构体类型的字段的类型可以是任何数据类型。当字段名称的首字母是大写字母时，我们就可以在任何位置（包括其他代码包）上通过其所属的结构体类型的值（以下简称结构体值）和选择表达式访问到它们。否则当字段名称的首字母是小写，这些字段就是包级私有的（只有在该结构体声明所属的代码包中才能对它们进行访问或者给它们赋值）。

如果一个字段声明中只有类型而没有指定名称，这个字段就叫做**匿名字段**。如上结构体 **Sequence** 中的 **Sortable** 就是一个匿名字段。**匿名字段**有时也被称为嵌入式的字段或结构体类型的嵌入类型。

匿名字段的类型必须由一个数据类型的名称或者一个与非接口类型对应的指针类型的名称代表。代表匿名字段类型的非限定名称将被隐含地作为该字段的名称。如果匿名字段是一个指针类型的话，那么这个指针类型所指的数据类型的**非限定名称**（由非限定标识符代表的名称）就会被作为该字段的名称。**非限定标识符**就是不包含代码包名称和点的标识符。

匿名类型的**隐含名称**的实例，如下：

```go
type Anonymities struct {
    T1
    *T2
    P.T3
    *P.T4
}
```

这个名为 **Anonymities** 的结构体类型包含了4个匿名字段。其中，**T1** 和 **P.T3** 为非指针的数据类型，它们隐含的名称分别为 **T1** 和 **T3**；***T2** 和 ***P.T4** 为指针类型，它们隐含的名称分别为 **T2** 和 **T4**。

>**注意**：匿名字段的隐含名称也不能与它所属的结构体类型中的其他字段名称重复。

结构体类型中的嵌入字段的类型所附带的方法都会成为该结构体类型的方法，结构体类型自动实现了它包含的所有嵌入类型所实现的接口类型。但是嵌入类型的方法的接收者类型仍然是该嵌入类型，而不是被嵌入的结构体类型。当在结构体类型中调用实际上属于嵌入类型的方法的时候，这一调用会被自动转发到这个嵌入类型的值上。

现在对 **Sequence** 的声明进行改动，如下：

```go
type Sequence struct {
    Sortable
    sorted bool
}
```

上面的 **Sequence** 中的匿名字段 **Sortable** 用来存储和操作可排序序列，布尔类型的字段 **sorted** 用来表示类型值是否已经被排序。

假设有一个 **Sequence** 类型的值 **seq**，调用 **Sortable** 接口类型中的方法 **Sort**，如下：

```go
seq.Sort()
```

如果 **Sequence** 类型中也包含了一个与 **Sortable** 接口类型中的方法 **Sort** 的名称和签名相同的方法，那么上面的调用一定是对 **Sequence** 类型值自身附带的 **Sort** 方法的调用，而嵌入类型 **Sortable** 的方法 **Sort** 被隐藏了。

如果需要在原有的排序操作上添加一些额外功能，可以这样声明一个同名的方法：

```go
func (self *Sequence) Sort() {
    self.Sortable.Sort()
    self.sorted = true
}
```

这样声明的方法实现了对于匿名字段 **Sortable** 的 **Sort** 方法的功能进行无缝扩展的目的。

如果两个 **Sort** 方法的名称相同但签名不同，那么嵌入类型 **Sortable** 的方法 **Sort** 也同样会被隐藏。这时，在 **Sequence** 的类型值上调用 **Sort** 方法的时候，必须依据该 **Sequence** 结构体类型的 **Sort** 方法的签名来编写调用表达式。如下声明 **Sequence** 类型附带的名为 **Sort** 的方法：

```go
func (self *Sequence) Sort(quicksort bool) {
    //省略若干语句
}
```

但是调用表达式 **seq.Sort()** 就会造成一个编译错误，因为 **Sortable** 的无参数的 **Sort** 方法已经被隐藏了，只能通过 **seq.Sort(true)** 或 **seq.Sort(false)** 来对 **Sequence** 的 **Sort** 方法进行调用。

> **注意**：无论被嵌入类型是否包含了同名的方法，调用表达式 **seq.Sortable.Sort()** 总是可以来调用嵌入类 **Sortable** 的 **Sort** 方法。

现在，区别一下**嵌入类型**是一个**非指针的数据类型**还是一个**指针类型**，假设有结构体类型 **S** 和非指针类型的数据类型 **T**，那么 ***S** 表示指向 **S** 的指针类型，***T** 表示指向 **T** 的指针类型，则：

1.    如果在 **S** 中包含了一个嵌入类型 **T**，那么 **S** 和 ***S** 的方法集合中都会包含接收者类型为 **T** 的方法。除此之外，***S** 的方法集合中还会包含接收者类型为 ***T** 的方法。

2.    如果在 **S** 中包含了一个嵌入类型 ***T**，那么 **S** 和 ***S** 的方法集合中都会包含接收者类型为 **T** 和 ***T** 的所有方法。

现在再讨论另一个问题。假设，我们有一个名为 **List** 的结构体类型，并且在它的声明中嵌入了类型 **Sequence**，如下：

```go
type List struct {
    Sequence
}

```

假设有一个 **List** 类型的值 **list**，调用嵌入的 **Sequence** 类型值的字段 **sorted**，如下：

```go
list.sorted
```

如果 **List** 类型也有一个名称为 **sorted** 的字段的话，那么其中的 **Sequence** 类型值的字段 **sorted** 就会被隐藏。

注意: 选择表达式 **list.sorted** 只代表了对 **List** 类型的 **sorted** 字段的访问，不论这两个名称为 **sorted** 的字段的类型是否相同。和上面的类似，这里选择表达式 **list.Sequence.sorted** 总是可以访问到嵌入类型 **Sequence** 的值的 **sorted** 字段。

对于结构体类型的多层嵌入的规则，有两点需要说明：

1.    可以在被嵌入的结构体类型的值上像调用它自己的字段或方法那样调用任意深度的嵌入类型值的字段或方法。唯一的前提条件就是这些嵌入类型的字段或方法没有被隐藏。如果它们被隐藏，也可以通过类似 **list. Sequence.sorted** 这样的表达式进行访问或调用它们。

2.    被嵌入的结构体类型的字段或方法可以隐藏任意深度的嵌入类型的同名字段或方法。任何较浅层次的嵌入类型的字段或方法都会隐藏较深层次的嵌入类型包含的同名的字段或方法。注意，这种隐藏是可以交叉进行的，即字段可以隐藏方法，方法也可以隐藏字段，只要它们的名称相同即可。

如果在同一嵌入层次中的两个嵌入类型拥有同名的字段或方法，那么涉及它们的选择表达式或调用表达式会因为编译器不能确定被选择或调用的目标而造成一个编译错误。

**匿名结构体类型**

匿名结构体类型比命名结构体类型少了关键字**type**和**类型名称**，声明如下：

```go
struct {
    Sortable
    sorted bool
}
```

可以在数组类型、切片类型或字典类型的声明中，将一个匿名的结构体类型作为他们的元素的类型。还可以将匿名结构体类型作为一个变量的类型，例如：

```go
var anonym struct {
    a int
    b string
}
```

不过对于上面，更常用的做法就是在声明以匿名结构体类型为类型的变量的同时对其初始化，例如：

```go
anonym := struct {
    a int
    b string
}{0, "string"}
```

与命名结构体类型相比，匿名结构体类型更像是“一次性”的类型，它不具有通用性，常常被用在临时数据存储和传递的场景中。

在Go语言中，可以在结构体类型声明中的字段声明的后面添加一个字符串字面量标签，以作为对应字段的附加属性。例如：

```go
type Person struct {
    Name    string `json:"name"`
    Age     uint8 `json:"age"`
    Address string `json:"addr"`
}
```

如上的字段的字符串字面量标签一般有两个**反引号**包裹的任意字符串组成。并且，它应该被添加但在与其对应的字段的同一行的最右侧。

这种标签对于使用该结构体类型及其值的代码来说是不可见的。但是，可以用标准库代码包 **reflect** 中提供的函数查看到结构体类型中字段的标签。这种标签常常会在一些特殊应用场景下使用，比如，标准库代码包 **encoding/json** 中的函数会根据这种标签的内容确定与该结构体类型中的字段对应的 **JSON** 节点的名称。

## 2. 值表示法

结构体值一般由复合字面量（类型字面量和花括号构成）来表达。在Go语言中，常常将用于表示结构体值的复合字面量简称为结构体字面量。在同一个结构体字面量中，一个字段名称只能出现一次。例如：

```go
Sequence{Sortable: SortableStrings{"3", "2", "1"}, sorted: false}
```

类型 **SortableStrings** 实现了接口类型 **Sortable**，这个可以在Go语言学习笔记4中了解到。这里就可以把一个 **SortableStrings** 类型的值赋给 **Sortable** 字段。

编写结构体字面量，还可以忽略字段的名称，但有如下的两个限制：

1.    如果想要省略其中某个或某些键值对的键，那么其他的键值对的键也必须省略。

    ```go
    Sequence{ SortableStrings{"3", "2", "1"}, sorted: false} // 这是不合法的
    ```
    
2.    多个字段值之间的顺序应该与结构体类型声明中的字段声明的顺序一致，并且不能够省略掉任何一字段的赋值。但是不省略字段名称的字面量却没有此限制。例如：
    ```go
    Sequence{ sorted: false , Sortable: SortableStrings{"3", "2", "1"}} // 合法
    Sequence{SortableStrings{"3", "2", "1"}, false} // 合法 
    Sequence{ Sortable: SortableStrings{"3", "2", "1"}} // 合法，未被明确赋值的字段的值会被其类型的零值填充。
    Sequence{ false , SortableStrings{"3", "2", "1"}} // 不合法，顺序不一致，会编译错误
    Sequence{ SortableStrings{"3", "2", "1"}} // 不合法，顺序不一致，会编译错误
    ```

在Go语言中，可以在结构体字面量中不指定任何字段的值。例如：

```go
Sequence{} // 这种情况下，两个字段都被赋予它们所属类型的零值。
```

与数组类型相同，结构体类型属于**值类型**。结构体类型的零值就是如上的不为任何字段赋值的结构体字面量。

## 3. 属性和基本操作

一个结构体类型的属性就是它所包含的字段和与它关联的方法。在访问权限允许的情况下，我们可以使用选择表达式访问结构体值中的字段，也可以使用调用表达式调用结构体值关联的方法。

在Go语言中，只存在嵌入而不存在继承的概念。不能把前面声明的 **List** 类型的值赋给一个 **Sequence** 类型的变量，这样的赋值语句会造成一个编译错误。在一个结构体类型的别名类型的值上，既不能调用那个结构体类型的方法，也不能调用与那个结构体类型对应的指针类型的方法。别名类型不是它源类型的子类型，但别名类型内部的结构会与它的源类型一致。

对于一个结构体类型的别名类型来说，它拥有源类型的全部字段，但这个别名类型并没有继承与它的源类型关联的任何方法。

如果只是将 **List** 类型作为 **Sequence** 类型的一个别名类型，那么声明如下：

```go
type List Sequence
```

此时，**List** 类型的值的表示方法与 **Sequence** 类型的值的表示方法一样，如下：

```go
List{ SortableStrings{"4", "5", "6"}, false}
```

如果有一个 **List** 类型的值 **List**，那么选择表达式 **list.sorted** 访问的就是这个 **List** 类型的值的 **sorted** 字段，同样，我们也可以通过选择表达式 **list.Sortable** 访问这个值的嵌入字段 **Sortable**。但是这个 **List** 类型目前却不包含与它的源类型 **Sequence** 关联的方法。

在Go语言中，虽然很多预定义类型都属于泛型类型（比如数组类型、切片类型、字典类型和通道类型），但却不支持自定义的泛型类型。为了使 **Sequence** 类型能够部分模拟泛型类型的行为特征，只是向它嵌入 **Sortable** 接口类型是不够的，需要对 **Sortable** 接口类型进行拓展。如下：

```go
type GenericSeq interface {
    Sortable
    Append(e interface{}) bool
    Set(index int, e interface{}) bool
    Delete(index int) (interface{}, bool)
    ElemValue(index int) interface{}
    ElemType() reflect.Type
    value() interface{}
}
```

如上的接口类型 **GenericSeq** 中声明了用于添加、修改、删除、查询元素，以及获取元素类型的方法。一个数据类型要实现 **GenericSeq** 接口类型，也必须实现 **Sortable** 接口类型。

现在，将嵌入到 **Sequence** 类型的 **Sortable** 接口类型改为 **GenericSeq** 接口类型，声明如下：

```go
type Sequence struct {
    GenericSeq
    sorted bool
    elemType reflect.Type
}
```

在如上的类型声明中，添加了一个 **reflect.Type** 类型（即标准库代码包 **reflect** 中的 **Type** 类型）的字段 **elemType**，目的用它来缓存 **GenericSeq** 字段中存储的值的元素类型。

为了能够在改变 **GenericSeq** 字段存储的值的过程中及时对字段 **sorted** 和 **elemType** 的值进行修改，如下还创建了几个与 **Sequence** 类型关联的方法。声明如下：

```go
func (self *Sequence) Sort() {
    self.GenericSeq.Sort()
    self.sorted = true
}

func (self *Sequence) Append(e interface{}) bool{
    result := self. GenericSeq.Append(e)
    //省略部分代码
    self.sorted = true
    //省略部分代码
    return result
}

func (self *Sequence) Set(index int, e interface{}) bool {
    result := self. GenericSeq.Set(index, e)
    //省略部分代码
    self.sorted = true
    //省略部分代码
    return result
}

func (self *Sequence) ElemType() reflect.Type {
    //省略部分代码
    self.elemType = self.GenericSeq.ElemType()
    //省略部分代码
    return self.elemType
}
```

如上的这些方法分别与接口类型 **GenericSeq** 或 **Sortable** 中声明的某个方法有着相同的方法名称和方法签名。通过这种方式隐藏了 **GenericSeq** 字段中存储的值的这些同名方法，并对它们进行了无缝扩展。

# 附录

**GenericSeq** 接口类型的实现类型以及 **Sequence** 类型的完整实现代码 [点击这里](https://github.com/hyper0x/goc2p/blob/master/src/basic/seq.go)
