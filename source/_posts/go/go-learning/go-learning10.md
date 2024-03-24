---
title: Go语言学习10-指针类型
date: 2016-07-08 21:19:23
updated: 2024-03-24 21:21:11
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 指针类型
---

[《开发语言-Go》](/categories/开发语言-Go/) [《Go语言学习》](/categories/开发语言-Go/Go语言学习/) 

![](/images/go-logo.png)

# 引言
书接上篇，我们了解了Go语言的[《结构体类型》](/2016/07/07/go/go-learning/go-learning9/)，本篇介绍Go语言的指针类型。主要如下：

# 主要内容

指针是一个代表着某个内存地址的值。这个内存地址往往是在内存中存储的另一个变量的值的起始位置。Go语言既没有像Java语言那样取消了代码对指针的直接操作的能力，也避免了C/C++语言中由于对指针的滥用而造成的安全和可靠性问题。

Go语言的指针类型指代了指向一个给定类型的变量的指针。它常常被称为指针的基本类型。指针类型是Go语言的复合类型之一。

## 1. 类型表示法

可以通过在任何一个有效的数据类型的左边加入 **\*** 来得到与之对应的指针类型。例如，一个元素类型为 **int** 的切片类型所对应的指针类型是 **\*[]int** ,前面的结构体类型 **Sequence** 所对应的指针类型是 **\*Sequence**。

注意：如果代表类型的是一个限定标识符（如 **sort.StringSlice**），那么表示与其对应的指针类型的字面量应该是 **\*sort.StringSlice** ,而不是 **sort.\*StringSlice**。

在Go语言中，还有一个专门用于存储内存地址的类型 **uintptr**。而 **uintptr** 类型与 **int** 类型和 **uint** 类型一样，也属于整数类型。它的值是一个能够保存一个指针类型值（简称指针值）的位模式形式。

## 2. 值表示法

如果一个变量 **v** 的值是可寻址的，表达式 **&v** 就代表了指向变量 **v** 的值的指针值。

>**知识点：** 如果某个值确实被存储在了计算机中，并且有一个内存地址可以代表这个值在内存中存储的起始位置，那么就可以说这个值以及代表它的变量是**可寻址的**。

## 3. 属性和基本操作

指针类型属于**引用类型**，它的零值是 **nil**。

对指针的操作，从标准代码包 **unsafe** 讲起，如下为省略文档的 **unsafe** 包下面的 **unsafe.go** 的源码（可自行到Go安装包 **src** 目录查看详细内容）：

```go
package unsafe

type ArbitraryType int
type Pointer *ArbitraryType

func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
```

在代码包 **unsafe** 中，有一个名为 **ArbitraryType** 的类型。从类型声明上看，它是 **int** 类型的一个别名类型。但是，它实际上可以代表任意的Go语言表达式的结果类型。事实上，它也并不算是 **unsafe** 包的一部分，在这里声明它仅处于代码文档化的目的。另外 **unsafe** 还声明了一个名为 **Pointer** 的类型，它代表了**ArbitraryType** 类型的指针类型。

如下有4个与 **unsafe.Pointer** 类型相关的特殊转换操作：

1.  一个指向其他类型的指针值都可以被转换为一个unsafe.Pointer类型值。例如，如果有一个float32类型的变量f32，那么可以将与它的对应的指针值转换为一个unsafe.Pointer类型的值：

  ```go
  pointer := unsafe.Pointer(&f32)
  ```
  其中，在特殊标记 **:=** 右边就是用于进行转换操作的调用表达式。取值表达式 **&f32** 的求值结果是一个 ***float32** 类型的值。

2. 一个 **unsafe.Pointer** 类型值可以被转换为一个与任何类型对应的指针类型的值。例如：

  ```go
  vptr := (*int)(pointer)
  ```
  上面的代码用于将 **pointer** 的值转换为与指向int类型值的指针值，并赋值给变量 **vptr**。***int** 类型和 ***float32** 类型在内存中的布局是不同的，如果我们在它们之上直接进行类型转换（对应表达式** (*int)(&f32))** 是不行，这会产生一个编译错误。有了上面的 **unsafe.Pointer** 作为中转类型的时候，看起来操作没有问题，但在使用取值表达式 ***vptr** 的时候会出现问题，**int** 类型的值和 **float32** 类型的值解析得到的结果是完全不同的，这样会产生一个不正确的结果。比如，如果这里对变量 **vptr** 的赋值语句改为：

  ```go
  vptr := (*string)(pointer)
  ```
  取值表达式 ***vptr**的求值就会引发一个运行时恐慌。

3.  一个 **uintptr** 类型的值也可以被转换为一个 **unsafe.Pointer** 类型的值。例如：
  ```go
  pointer2 := unsafe.Pointer(uptr)
  ```


4.  一个 **unsafe.Pointer** 类型值可以被转换为一个 **uintptr** 类型的值。例如：

  ```go
  uptr := uintptr(pointer)
  ```
  
>**注意**：正是因为这些特殊的转换操作，**unsafe.Pointer** 类型可以使程序绕过Go语言的类型系统并在任意的内存地址上进行读写操作成为可能。**但这些操作非常危险，小心使用**。

现在用之前的结构体类型 **Person** 举例，如下：

```go
type Person struct {
    Name  string `json:"name"`
    Age   uint8 `json:"age"`
    Address string `json:"addr"`
}
```

初始化 **Person** 的值，并把它的指针值赋给变量 **p** :

```go
p := &Person(“Huazie”, 23, “Nanjing”)
```

下面利用上述特殊转换操作中的第一条和第三条获取这个结构体值在内存中的存储地址：

```go
var puptr = uintptr(unsafe.Pointer(p))
```

变量 **puptr** 的值就是存储上面那个 **Person** 类型值的内存地址。由于类型 **uintptr** 的值实际上是一个无符号整数，所以我们可以在该类型的值上进行任何算术运算。例如：

```go
// 变量np表示结构体中的Name字段值的内存地址。
var np uintptr = puptr + unsafe.Offsetof(p.Name) 
```

如上 **unsafe.Offsetof** 函数会返回作为参数的某字段（由相应的选择表达式表示）在其所属的结构体类型之中的存储偏移量。也就是，在内存中从存储这个结构体值的起始位置到存储其中某字段的值的起始位置之间的距离。这个存储偏移量（或者说距离）的单位是字节，它的值的类型是 **uintptr**。对于同一个结构体类型和它的同一个字段来说，这个存储偏移量总是相同的。

在获得存储 **Name** 字段值的内存地址之后，将它还原成指向这个 **Name** 字段值的指针类型值，如下：

```go
var name *string = (*string)(unsafe.Pointer(np))
```

获取这个 **Name** 字段的值：

```go
*name
```

只要获得了存储某个值的内存地址，就可以通过一定的算术运算得到存储在其他内存地址上的值甚至程序。如下一个恒等式显示上面的一些操作：

```go
uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f) == uintptr(unsafe.Pointer(&s.f))
```
# 结语
Go数据类型的知识就记到这，下一篇介绍 **Go语言数据的使用**。其中 **通道类型**，比较特殊，将会在后续的博文仔细讲解，敬请期待！！！

最后附上知名的Go语言开源框架： 
>**Skynet:** 一个分布式服务框架。它可以帮助我们构建起大规模的分布式应用系统。它的源码放置在https://github.com/skynetservices/skynet上。 
