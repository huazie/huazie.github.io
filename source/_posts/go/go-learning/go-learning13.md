---
title: Go语言学习13-类型转换
date: 2016-07-14 21:34:20
updated: 2024-03-24 21:21:11
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 类型转换
---



![](/images/go-logo.png)

# 引言
在上一篇博文中，我们介绍了 `Go` 语言的 [《数据的使用》](../../../../../../2016/07/13/go/go-learning/go-learning12/)；本篇博文，我们将介绍 `Go` 语言的类型转换。

<!-- more -->

# 主要内容
## 1. 概念
类型转换是把一个类型的值转换成为另一个类型的值。把这个值原来的类型称为**源类型**，而这个值被转换后的类型称为**目标类型**。

如果 **T** 是值 **x** 的目标类型，那么相应的类型转换表达式如下：

```go
T(x) // x可以是一个表达式，不过这个表达式的结果值只能有一个
```

如果代表目标类型的字面量始于操作符 * 或 **<-** ，后者它是没有结果声明列表的函数类型，那么往往需要用圆括号括起来，以避免歧义的产生。例如：

```go
*string(v)   // 等同于 *(string(v))，先将变量v代表的值转换为string类型的值，然后再获取指向它的指针类型值。
(*string)(v) // 把变量v的值转换为指针类型*string的值。

<-chan int(v)   // 等同于 <-(chan int(v)),先将变量v代表的值转换为chan int类型的值，然后再从此通道类型值中接收一个int类型的值。
(<-chan int)(v) // 把变量v的值转换为通道类型<-chan int的值。

func()(v)     // Go语言理解为任何无参数声明但有一个结果声明的匿名函数。
(func())(v)   // 把变量v的值转换为函数类型func()的值。
func() int(v) // 等同于(func() int)(v),把v的值转换成一个有结果声明的函数的类型。
```

对于常量 **x**，如果它能够被转换为类型 **T** 的值，那么它们符合如下情况：

 - **x** 可以被类型T的值代表。例如，**iota** 可以表示一个大于等于零的整数常量。他可以把 **uint** 类型的值代表。类型表达式 **uint(iota)** 是合法的，它的结果值会是一个 **uint** 类型的常量。
 
 - **x** 是一个浮点数常量，**T** 是一个浮点数类型，并且 **x** 在（根据IEEE-754标准中描述的向偶数舍入规则）被舍入之后可以被类型 **T** 的值代表。例如：
 
	```go
	float32(0.49999998) // 求值结果是一个float32类型的常量0.5
	```
 - 如果 **x** 是一个整数常量，并且 **T** 是一个 **string** 类型，那么将会遵循一套规则来决定类型转换的结果。它同样适合于非常量的值。这将在后面的 **与 string类型相关的转换** 处讲解。
 
 - 对于非常量 **x**，它能够被转换为类型 **T** 的值，那么它们符合如下情况：
值 **x** 可以被赋值给类型 **T** 的变量。例如：

	```go
	type Computer interface {
		CpuType() string
	}

	type Laptop struct {
		cpuType string
	}

	func (self Laptop) CpuType() string {
		return self.cpuType
	}
	```
   **类型转换表达式:**

	```go
	// 合法，求值结果会是一个Computer类型的值。因为类型Laptop是接口类型Computer的一个实现类型
	Computer(Laptop{cpuType: "Intel Core i5"}) 
	```

 - 值 **x** 的类型和类型 **T** 的潜在类型是相等的。例如：
 
	```go
	type MyString string
	// 类型转换表达式
	MyString("Huazie") // 合法，类型MyString的潜在类型就是string类型。
	```

 - 值 **x** 的类型和类型 **T** 都是未命名的指针类型，并且它们的基本类型（指向那个值的类型）的潜在类型是相等的。例如:

	```go
	var str1 string
	// 类型转换表达式
	(*MyString)(&str1) // 合法，求值结果是一个*MyString类型的值。
	```

 - 值 **x** 的类型和类型 **T** 都是整数类型或都是浮点数类型。例如:
 
	```go
	var i32 uint32
	var f32 folat32
	// 类型转换表达式
	int64(i32) // 合法。
	float64(f32) // 合法。
	```

 - 值 **x** 的类型和类型 **T** 都是复数类型。例如:
 
	```go
	var comp64 complex64
	// 类型转换表达式
	complex128(comp64) // 合法。
	```

 - 值 **x** 是一个整数类型值或是一个元素类型为 **byte** 或 **rune** 的切片类型值，且 **T** 是一个 **string** 类型。例如：
 
	```go
	string([]byte{'a'}) // 合法，求值结果是string类型值"a"。
	```

 - 值 **x** 是一个 **string** 类型值，且T是一个元素类型为 **byte** 或 **rune** 的切片类型。
 
	```go
	[]rune("Huazie") // 合法
	```

## 2. 数值类型之间的转换

可以通过常量声明或者数据类型转换把一个无类型的常量类型化，例如：

数值常量**1024**是**无类型**的，可以把它赋给一个 `int` 类型的变量：

```go
var number int = 1024
```
或把它向 `int` 类型转换：

```go
int(1024)
```

对于**非常量的数值类型值**，规则如下：

 - 当把一个整数类型值从需要较少二进制位表示的整数类型转换到需要较多二进制位表示的整数类型（比如从 **int8** 类型转换到 **int16** 类型）的时候 : 如果这个整数类型值是有符号的，那么该符号位上的（最左边的）那个二进制值将作为扩展项填充在转换过程中新增的那些二进制位上，否则将会把 **0** 作为扩展项进行填充。这种扩展方式是针对整数类型值的补码而言的。例如：**int16** 类型值 **-32767** 的十六进制表示是 **0xffff** 。它的补码是 **0x8001** 。此补码最左边的二级制位上的二级制值是 **1**。如果要把这个 **int16** 类型值转换为 **int32** 类型值，就需要用最左边的这个值 **1** 填充在高位一侧新增的那**16**个二进制位上。类型转换之后的补码是 **0xffff8001** 。在这个补码之上再求其补码以得其原码，即 **0x80007fff** 。此原码表示的就是十进制数 **-32767** ，与类型转换前的那个数值相等。
 
 - 当把一个整数类型值从需要较多二级制位表示的整数类型转换到需要较少二级制位表示的整数类型的时候，需要把多余的若干个较高位置的二进制值裁掉，而只保留与目标类型所需二进制位数相当的若干个较低位置的二进制值。例如，**int16** 类型值 **-32767**，如果要把它转换为一个 **int8** 类型值，就需要对其补码 **0x8001** 截取较低 **8** 为的二进制值，得到 **0x01**。由于此值的最左边的二进制位上是 **0**，所以它本身就是类型转换总会得到一个有效的数值。但对于整数常量来说，这样的类型转换就会造成一个编译错误。例如，类型转换表达式 **int8(-32767)** 会使编译器报错，因为整数常量 **-32767** 超出了 **int8** 类型所能表示的数值范围。
 
  - 当把一个浮点数类型值向整数类型值进行转换的时候，该浮点数类型值的小数部分将被抹去。例如，如果有一个 **float32** 类型的变量 **f32** 且其值为 **-32767.345** ,那么类型表达式 `int32(f32)` 的求值结果为 **-32767** 。如果浮点数类型值在被抹去小数位之后超出了目标整数类型的表示范围，那么该值还会被截短。例如，在类型表达式 `int8(f32)` 被求值的过程中会首先 **float32** 类型值 **-32767.345** 的小数部分却去掉，然后再将其中较高的 **24** 位的二进制值截掉，最终得到结果 **1** 。
 
 - 当把一个整数或浮点数转换为一个浮点数类型的值或者把一个复数转换为一个复数类型的值的时候，该值将会被依据目标类型的精度进行舍入操作。例如，在 **float32** 类型的变量 **x** 中存储的值可能会超出 **IEEE-754** 标准中规定的 **32** 位（二进制值代表的）浮点数的精度。但是，类型表达式 `float32(x)` 的求值结果一定会是 **x** 的值向**32**位浮点数的精度转化之后的值。算术表达式 `x + 0.1` 的结果值可能会超出**32**位浮点数的精度，但是类型转换表达式 `float32(x + 0.1)` 的求值结果却不会这样。

在非常量的**浮点数类型**值或**复数类型值**的类型转换中，当目标类型的精度不能够满足被转换的值的需要的时候，虽然转换会成功，但其结果将是不确定的，这依赖于不同平台的Go语言的具体实现。

## 3. 与string类型相关的转换

 - 当把一个有符号整数值或无符号整数值向字符串类型转换的时候，将会产生出一个字符串类型值。被转换的整数值应该是一个有效的 **Unicode** 代码点的代表。在作为结果的字符串类型值中的就是那个 **Unicode** 代码点对应的字符。在底层，这个字符串类型值是由该 **Unicode** 代码点的 **UTF-8** 编码值表示的。如果被转换的整数值不能代表一个有效的 **Unicode** 代码点，那么转换结果将会是 `"\ufffd"`，即 **Unicode** 字符“�”。例如:

	```go
	string(0x4e2d) // 求值结果为“中”,其UTF-8编码为\xe4\xb8\xad
	string('国')   // 求值结果值为“国”，其UTF-8编码为\xe5\x9b\xbd
	string(-1)     // 求值结果为“�”，整数值-1不能代表一个有效的Unicode代码点。
	```
	如果有一个目标类型是 **string** 类型的别名类型是 **MyString**，那么可以将它视同为 **string** 类型。例如：

	```go
	MyString(0x4e2d) // 等同于string(0x4e2d)
	```

 - 当把一个元素类型为 **byte** 的切片类型值向字符串类型转换时，将会产生出一个字符串类型值。这个字符串类型值实际上就是由被转换的切片类型值中的每个字节类型值依次组合而成的。如果切片类型值为 **nil**，那么类型转换的结果将会是“”。例如：
 
	```go
	string([]byte{'g', '\x6f', '\x6c', '\x61', 'n', 'g'})//求值结果是"golang"
	```
	由于使用 **"\x"** 为前导并后跟两位十六进制数可以表示宽度为一个字节的值，因此一个字节类型的值也就可以由这种方法表示。如果源类型是一个 **[ ]byte** 类型的别名类型，那么可以将它视同为 **[ ]byte** 类型。

 - 当把一个元素类型为 **rune** 的切片类型值向字符串类型转换时，将会产生出一个字符串类型值。这个字符串类型值实际上就是依次串联每个 **rune** 类型值后的结果。如果切片类型值为 **nil**，那么类型转换的结果将会是“”。例如：
 
	```go
	string([]rune{ 0x4e2D, 0x56fd })//求值结果是"中国"
	```
	如果源类型是一个 **[ ]rune** 类型的别名类型，那么我们可以将它视同为 **[ ]rune** 类型。

 - 当把一个字符串类型值向 **[ ]byte** 类型转换时，其结果将会是把该字符串类型值按字节拆分后的结果。对于“”来说，转换后的结果一定是 **[ ]byte** 类型的空值 **nil** 。例如：
 
	```go
	[]byte("hello")//结果是[]byte{104, 101, 108, 108, 111}
	```
	在这个 **[ ]byte** 类型值中的每个元素都是对应字符的ASCII编码值的十进制表示形式。如果目标类型是一个 **[ ]byte** 类型的别名类型，那么可以将它视同为 **[ ]byte** 类型。

 - 当把一个字符串类型值向 **[ ]rune** 类型转换时，其结果将会是把该字符串类型值按字符拆分后的结果。对于 "" 来说，转换后的结果一定是 **[ ]rune** 类型的空值 **nil** 。
 
	```go
	[]rune("中国") // 结果是[]byte{20013, 22269}
	```
在这个 **[ ]rune** 类型值中的每个元素都是对应字符的 **Unicode** 代码点的十进制表示形式。如果目标类型是一个 **[ ]rune** 类型的别名类型，那么可以将它视同为 **[ ]rune** 类型。

**UTF-8** 这种编码方式会把一个字符编码为一个或多个字节。对于同一个字符串类型值来说，与它对应的字节序列和字符序列中的元素并不一定是一 一对应的。字节序列中的单个字节并不一定能代表一个完整的字符。例如，以字符串类型值“中国”为例：

```go
// 字节序列的前三个元素代表了字符'中'的UTF-8编码值，而后三个元素则代表了字符'国'的UTF-8编码值。
[]byte{228, 184, 173, 229, 155, 189}
// 这个字符序列中的第一个元素代表了字符'中'的Unicode代码点，而第二个元素则代表了字符'国'的Unicode代码点。
[]byte{20013, 22269}
```

对于每一个 **ASCII** 编码可表示的字符来说，它的 **Unicode** 代码点和 **UTF-8** 编码值与其 **ASCII** 编码值都分别是一致的，且它们都可以由一个字节类型值代表。对于一个包含了 **ASCII** 编码可表示的字符的字符串类型值来说，与它对应的字节序列和字符序列中的元素值必定也是一一对应的。

**byte** 类型值和 **rune** 类型值都属于整数值的一种。所有整数值都可以由十进制字面量、八进制字面量和十六进制字面量来代表。可以把任意一种方式表示的 **rune** 字面量赋给任何整数类型的变量，只要该 **rune** 字面量对应的 **Unicode** 代码点不超出那个整数类型的表示范围。例如：

```go
var nation int16 = '国' // '国' == 0x56fd == 22269
[]byte{ 'g', '\x6f', '0x6c', '\u0061', '\156', '\U00000067' } // 求值结果是"golang"
```

## 4. 别名类型值之间的转换

类型是 **MyString** 是 **string** 类型的别名类型。如果一个整数值分别转换为这两个类型的值，将会得到相同的结果。把一个字符串字面量赋给 **MyString** 类型的变量：

```go
var ms MyString = "中国"
```

在 **MyString** 类型的值之上应用切片操作：

```go
ms[1]
```
在某个数据类型和它的别名类型之间以及同一个数据类型的多个别名类型之间的类型转换是合法的。并且，在这种类型转换的过程中并不会创造出新的值，而仅仅是变换了一下那个已存在的值的所属类型。

# 结语
本篇主要介绍了 **Go** 语言数据使用中类型转换相关的内容，下一篇我们将会介绍 **Go** 语言的一些内建函数的使用，敬请期待！！！

最后附上知名的Go语言开源框架： 

**etcd:**  一个高可用的键值存储系统。它可被用于建立共享配置系统和服务发现系统。它的灵感来自于 **Apache ZooKeeper** 。我们可以在 https://github.com/coreos/etcd 上找到它的源码。



