---
title: Go语言学习14-内建函数
date: 2016-07-15 20:18:49
updated: 2024-03-24 21:21:11
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 内建函数
  - close函数
  - len函数
  - cap函数
  - new函数
  - make函数
  - append函数
  - copy函数
  - delete函数
  - complex函数
  - real函数
  - imag函数
  - panic函数
  - recover函数
  - print函数
  - println函数
---

[《开发语言-Go》](/categories/开发语言-Go/) [《Go语言学习》](/categories/开发语言-Go/Go语言学习/) 

![](/images/go-logo.png)

# 引言
上一篇博文，我们介绍了 `Go` 语言的 [《类型转换》](/2016/07/14/go/go-learning/go-learning13/)；本篇博文我们重点介绍 `Go` 语言的 内建函数
# 主要内容

所谓内建函数，就是Go语言内部预定义的函数。调用它们的方式与调用普通函数并无差异，并且在使用它们之前也不需要导入任何代码包。这里并不能把内建函数当做值来使用。因为它们并不像普通函数那样有隶属的Go语言数据类型。

## 1. close函数

内建函数 **close** 只接受通道类型的值（简称通道）作为参数。例如：

```go
ch := make(chan int, 1)
close(ch)
```

调用这个 **close** 函数之后，会使作为参数的通道无法再接受任何元素值。若试图关闭一个仅能接受元素值的通道，则会造成一个编译错误。在通道被关闭之后，再向它发送元素值或者试图再次关闭它的时候，都会引发一个运行时恐慌。试图关闭一个为 **nil** 的通道值也会引发一个运行时恐慌。

试图调用 **close** 函数关闭一个通道，并不会影响到在此调用之前已经发送的那些元素值，它们会被正常接收（如果存在接收操作的话）。但是，在此调用之后，所有的接收操作都会立即返回一个该通道的元素类型的零值。

## 2. len函数与cap函数

### 2.1 len函数的使用方法

|参数类型     |    结果          |        备注         |
|:---------- |:-----------------|:-------------------|
|string    |string类型值的字节长度  |  无        |
|[n]T或*[n]T  |数组类型值长度，它等于n  | n代表了数组类型的长度，T代表了数组类型的元素类型        |
|[]T      |切片类型值的长度   |  T代表了切片类型的元素类型       |
|map[K]T  |字典类型值的长度，其中已包含的键的数量   | K代表了字典类型的键类型，T代表了字典类型的元素类型     |
|chan T    |通道类型值当前包含的元素的数量   |  T代表了通道类型的元素类型        |              


 
### 2.2 cap函数的使用方法

|参数类型     |    结果          |        备注         |
|:---------- |:----------------|:-------------------|
|[n]T或*[n]T |数组类型值长度，它等于n  | n代表了数组类型的长度，T代表了数组类型的元素类型        |
|[]T        |切片类型值的容量    | T代表了切片类型的元素类型  |  
|chan T     |通道类型值的容量    | T代表了通道类型的元素类型  |
         

对于一个切片类型值来说，它的长度和容量的关系：

**0 <= len(s) <= cap(s)**

一个切片值的容量就是它拥有的那个底层数组的长度。这个底层数组的长度必定不会小于该切片值的长度。

值为 **nil** 的 **切片类型值**、**字典类型** 和 **通道类型值** 的长度都是 **0** 。值为 **nil** 的 **切片类型值** 和 **通道类型值** 的容量也都是 **0** 。

如果 **s** 是一个 **string** 类型的常量，那么表达式 **len(s)** 和 **cap(s)** 也都等同于常量。**len(s)** 所代表的值在编译期间就会被计算出来。如果 **s** 是一个表达式，且其类型是数组类型或指向数组类型的指针类型，那么只要该表达式中不包含通道接收操作和函数调用操作，它就不会被求值。因为 **s** 的类型中已经包含了它的长度信息。在对表达式 **len(s)** 和 **cap(s)** 进行求值的时候并不需要求得 **s** 的结果值而只需要从 **s** 的类型中取得其长度即可。在这种情况下，这两个表达式也会等同于常量。

## 3. new函数和make函数
这两个函数是专门用于数据初始化。这里的数据初始化是指对某个数据类型的值或变量的初始化。在Go语言中，几乎所有的数据类型的值都可以使用字面量来进行表示和初始化。在大多数情况下，使用字面量就可以满足初始化值或变量的要求。

这部分内容，可参考这篇博文 [**《Go语言学习11-数据初始化》**](/2016/07/09/go/go-learning/go-learning11/)

## 4. append函数和copy函数

**append** 函数和 **copy** 函数都被用于辅助在切片类型值之上的操作。

这部分内容，可参考这篇博文 [**《Go语言学习5-切片类型》**](/2016/07/03/go/go-learning/go-learning5/)。

## 5. delete函数

内建函数 **delete** 专用于删除一个字典类型值中的某个键值对。它接受两个参数，第一个参数是作为目标的字典类型值，而第二个参数则是能够代表要删除的那个键值对的键。例如：

```go
delete(m, k)
```

这里有两点需要注意：

 - 第二个参数 **k** 与 **m** 的键的类型之间必须满足赋值规则。
 - 当 **m** 的值是 **nil** 或者 **k** 所代表的键值对并不存在于 **m** 中的时候，**delete(m, k)** 不会做任何操作。在没有可删除的目标的时候，删除操作将被忽略。这种删除失败不会被反馈。

## 6. complex 函数、real 函数和 imag 函数

这3个内建函数都是专用于操作复数类型值的。

**complex** 函数被用于根据浮点数类型的实部和虚部来构造复数类型。例如：

```go
var cplx128 complex128 = complex(2, -2)
```

内建函数 **real** 和 **imag** 则分别被用于从一个复数类型值中抽取浮点数类型的**实部**部分和浮点数类型的**虚部**部分。例如：

```go
var im64 = imag(cplx128)
var r64 = real(cplx128)
```

对于 **complex** 函数来说，两个参数的类型必须是同一种浮点数类型，并且其结果类型与参数类型对应。在Go语言中，对于复数有一个恒等式：

```go
z == complex(real(z), imag(z)) // z是一个复数类型的变量
```

>**注意：** 如果 **complex** 函数的两个参数都没有显示的类型，那么该函数的结果的类型将会是 **complex128** 类型的。如果 **complex** 函数的参数都是常量，那么它的结果值也必是常量。

## 7. panic函数和recover函数

内建函数 **panic** 函数和 **recover** 函数分别被用于报告和处理运行时恐慌。

函数 **panic** 只接受一个参数。这个参数可以是任意类型的值。按照惯例，**panic** 函数的实际参数的类型常常是接口类型 **error** 的某个实现类型。**panic** 函数的参数都应该足以表示恐慌发生时的异常情况。

函数 **recover** 不接受任何参数，但是返回一个 **interface{}** 类型的结果值。**interface{}** 就是空接口。所有的数据类型都是它的实现类型。因此，**recover** 函数的结果值可能是任何类型的。这是与 **panic** 函数的那个唯一参数相对应的，它们都是 **interface{}** 类型的。如果运行时恐慌的报告是通过调用 **panic** 函数来进行的话，那么之后调用 **recover** 函数所得的结果值就应该是先前 **panic** 函数在被调用时接受的那个参数值。**recover** 函数的结果值也有可能 **nil** 。如果是 **nil**，属于以下的情况：

 - 传递给panic函数的参数值就是nil。
 - 运行时恐慌根本就没有发生。狭义的讲，panic函数没有被调用。
 - 函数recover并没有在defer语句中被调用。

在任何情况下任何位置上调用 **recover** 函数都不会产生任何副作用。如果不用它来处理运行时恐慌，那么对它的调用也就没有任何意义了。**panic** 函数和 **recover** 函数之间肯定是存在着某种联系，在后面的博文中将对异常报告和处理的更多细节进行讲解。

## 8. print函数和println函数

函数 **print** 的作用是依次（即从左到右）打印出传递给它的参数值，每个参数值对应的打印内容都由它们的具体实现决定。而函数 **println** 函数则会在 **print** 函数打印的内容的基础上再在每个参数之间加入空格，并在最后加入换行符。例如：

```go
print("A", 12.4, 'R', "C")
println("A", 12.4, 'R', "C")
```

调用表达式被求值之后，出现内容：

**A+1.240000e+00182CA +1.240000e+001 82 C**

对于上面的这两个函数，有以下需要注意：

 - 它们接受的参数只能是有限的数据类型的值。并且，在这些受支持的数据类型当中。大部分都是Go语言的基础数据类型。
 
 - 这两个函数针对于每种受支持的数据类型的打印格式都是固定的，无法自定义。

Go语言并不保证会在以后的版本中一直保留这两个函数。因此，尽量不要在程序中使用这两个函数，尤其是用于生产环境的程序。应该使用标准库代码包 **fmt** 中的函数 **Print** 和 **Println** 来替代它们。

# 结语
至此，**Go语言数据的使用** 就讲完了，下篇博文将要介绍 **Go语言流程控制方法**，敬请期待！！！