---
title: Go语言学习16-特殊流程控制
date: 2016-07-20 20:14:16
updated: 2024-03-24 21:21:11
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 特殊流程控制
  - defer
  - error
  - panic
  - recover
---



![](/images/go-logo.png)

# 引言
上一篇博文介绍了 `Go` 语言的[《基本流程控制》](../../../../../../2016/07/18/go/go-learning/go-learning15/)，本篇我们介绍 `Go` 语言的特殊流程控制。

<!-- more -->

# 主要内容

## 1. defer语句

**defer** 语句被用于预定对一个函数的调用。被 **defer** 语句调用的函数称为延迟函数。**defer** 语句只能出现在函数或方法的内部。

一条 **defer** 语句总是以关键字 **defer** 开始。在 **defer** 的右边还会有一条表达式语句，且它们之间要以空格分隔。例如：

```go
defer fmt.Println("The finishing touches.")
```

如上的表达式语句必须代表一个函数或方法的调用。但是像针对各种内建函数的那些调用表达式，因为它们并不能称为表达式语句，所以不允许出现在这里。同时这个位置出现的表达式语句是不能被圆括号括起来的。

**defer** 语句的执行时机总是在直接包含它的那个函数（简称**外围函数**）把流程控制权交还给它的调用方的前一刻，无论 **defer** 语句出现在外围函数的函数体中的哪一个位置上。具体分为：

 - 当外围函数的函数体中的相应语句全部被正常执行完毕的时候，只有在该函数中的所有 **defer** 语句都被执行完毕之后该函数才会真正地结束执行。
 
 - 当外围函数的函数体中的 **return** 语句被执行的时候，只有在该函数中的所有 **defer** 语句都被执行完毕之后该函数才会真正地返回。
 
 - 当在外围函数中有运行时恐慌发生的时候，只有在该函数中的所有 **defer** 语句都被执行完毕之后该运行时恐慌才会真正地被扩散至该函数的调用方。

也就是说，外围函数的执行的结束会由于其中的 **defer** 语句的执行而被推迟。例如：

```go
func isPositiveEnenNumber(number int) (result bool){
    defer fmt.Println("done.");
    if number < 0 {
        panic(errors.New("The number is a negative number!"))
    }
    if number % 2 ==0 {
        return true
    }
    return
}
```
上述示例中，无论参数 **number** 是怎样的值，以及该函数的执行会以怎样的方式结束，在该函数的调用方重获流程控制权之前标准输出上都一定会打印 ```done.```

综上总结，使用 **defer** 语句的优势有两个：

 - 收尾任务总会被执行，这样就不会因粗心大意而造成资源的浪费。
 
 - 可以把它们放到外围函数的函数体中的任何地方（一般是**函数体开始处**或紧跟在申请资源的语句的后面），而不是只能放在函数体的最后。

在 **defer** 语句中，调用的函数不但可以是已声明的命名函数，还可以是临时编写的匿名函数。例如：

```go
defer func(){
    fmt.Println("The finishing touches.")
}()
```

>**注意：** 一个针对匿名函数的调用表达式是由一个**函数字面量**和一个代表了调用操作的 **一对圆括号** 组成的。

无论在 **defer** 关键字右边的是命名函数还是匿名函数，都可以称为 **延迟函数**。因为它总是会被延迟到外围函数执行结束前一刻才被真正地调用。每当 **defer** 语句被执行的时候，传递给延迟函数的参数都会以通常的方式被求值。如下：

```go
func begin(funcName string) string {
    fmt.Printf("Enter function %s.\n", funcName)
    return funcName
}

func end(funcName string) string {
    fmt.Printf("Exit function %s.\n", funcName)
    return funcName
}

func record(){
    defer end(begin("record"))
    fmt.Println("In function record")
}
```

对函数 **record** 进行调用之后，运行截图如下：
 
![](result.png)

出于同一条 **defer** 语句可能会被多次执行的考虑，如下：

```go
func printNumbers(){
    for i := 0; i < 5; i++ {
        defer fmt.Printf("%d ", i)
    }
}
```

对函数 **printNumbers** 进行调用之后，运行截图如下：
 
![](result-1.png)

如上的函数 **printNumbers** 有两点需要关注：

 - 在for语句的每次迭代的过程中都会执行一次其中的defer语句。Go语言会把代入参数值之后的调用表达式另行存储，以此类推，后面几次迭代所产生的延迟函数调用表达式依次为：
 
    ```go
    fmt.Printf("%d ", 0)
    fmt.Printf("%d ", 1)
    fmt.Printf("%d ", 2)
    fmt.Printf("%d ", 3)
    fmt.Printf("%d ", 4)
    ```
 - 对延迟函数调用表达式的求值顺序是与它们所在的defer语句被执行的顺序完全相反的。每当Go语言把已带入参数值的延迟函数调用表达式另行存储之后，还会把它追加到一个专门为当前外围函数存储延迟函数调用表达式的列表（也就是栈）当中，而该列表总是先进后出。因此这些延迟函数调用表达式的求值顺序会是：
 
    ```go
    fmt.Printf("%d ", 4)
    fmt.Printf("%d ", 3)
    fmt.Printf("%d ", 2)
    fmt.Printf("%d ", 1)
    fmt.Printf("%d ", 0)
    ```

我们再看看下面的例子，如下：

```go
func appendNumber(ints []int) (result []int) {
    result = append(ints, 1)
    defer func(){
        result = append(result, 2)
    }()
    result = append(result, 3)
    defer func(){
        result = append(result, 4)
    }()
    result = append(result, 5)
    defer func(){
        result = append(result, 6)
    }()
    return result
}
func main(){
    res := appendNumber([]int{0})
    fmt.Printf("result: %v\n", res)
}
```
运行结果截图如下【大家可以试着按上面的两点去分析下】：

![](result-2.png)

现在考虑一个问题，把 **printNumbers** 函数的声明修改为如下：

```go
func printNumbers(){
    for i := 0; i < 5; i++ {
        defer func() {
            fmt.Printf("%d ", i)
        }()
    }
}
```

运行结果截图如下：
 
![](result-3.png)

现在我们对运行结果进行分析可知：

在 **for** 语句被执行完毕的时候，共有 **5** 个相同的延迟函数调用表达式被存储在专属列表(**栈**)中，例如：

```go
func() {
    fmt.Printf("%d ", i)
}()
```
这时的变量 **i** 已经被修改为了 **5**，对 **5** 个相同的调用表达式的求值都会使标准输出打印出 **5** 。

针对上面的情况，可以修改如下：

```go
defer func(i int) {
    fmt.Printf("%d ", i)
}(i) // 在defer语句被执行的时候，传递给延迟函数的这个参数i就会被求值。
```

运行结果截图如下（这个和第一个版本的 **printNumbers** 函数执行效果是相同的）：
 
![](result-4.png)

如果 **延迟函数** 是一个匿名函数，并且在 **外围函数** 的声明中存在命名的结果声明，那么在延迟函数中的代码使可以对命名结果的值进行访问和修改的。例如：

```go
func modify(n int) (number int) {
    defer func(){
        number += n
    }()
    number++
    return
}
```

在 **延迟函数** 的声明中可以包含结果声明，但是其返回的结果值会在它被执行完毕时被丢弃。因此在编写延迟函数的声明的时候不会为其添加结果声明。另外，**推荐以传参的方式提供延迟函数所需的外部值**。例如：

```go
// 传入参数为1时，modify函数的结果值是5
func modify(n int) (number int) {
    defer func(plus int) (result int){
        result = n + plus
        number += result
        return // 此处虽然返回了结果，但是却并不会产生任何效果。
    }(3) // 延迟函数调用时直接传外部参数
    number++
    return
}
```

## 2. 异常处理

在前面的博文中已经涉及了Go语言的异常处理的知识，比如 **接口类型error**、**内建函数panic** 和 **标准库代码包errors**。本小节将对Go语言的各种异常处理方法进行系统的讲解。

### 2.1 error

在Go语言标准库代码包中的很多函数和方法会返回 **error** 类型值来表明错误状态及其详细信息。**error** 是一个预定义标识符，它代表了一个Go语言内建的接口类型。该接口类型声明如下：

```go
type error interface {
    Error() string
}
```

其中，**Error** 方法声明的意义就在于为方法调用方提供当前错误状态的详细信息。任何数据类型只要实现了这个可以返回 **string** 类型值的 **Error** 方法就可以成为一个 **error** 接口类型的实现。但在通常情况下，不需要自己编写一个 **error** 的实现类型，Go语言的标准库代码包 **errors** 为我们提供了一个用于创建 **error** 类型值的函数 **New**，声明如下：

```go
func New(text string) error {
    return &errorString(text) // 返回一个error类型值，它的动态类型就是errors.errorString类型
}
```
从 **errors.errorString**的名称上可知，**errorString** 的首字母小写，该类型是一个包级私有的类型。它只是**errors** 包的内部实现的一部分，而非公开的 **API** 。**errors.errorString** 类型及其方法的声明如下：

```go
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

打印**error**类型值所代表的错误的详细信息。

```go
var err error = errors.New("A normal error example")
fmt.Println(err)
fmt.Printf("%s\n", err)
```

另一个可以生成 **error** 类型值的方法是调用 **fmt** 包中的 **Errorf** 函数，调用类似如下代码：

```go
// 初始化一个error类型值并作为该函数的结果值返回给调用方。
err2 := fmt.Errorf("%s\n","A normal error")
```


在 **fmt.Errorf** 函数的内部，创建和初始化 **error** 类型值的操作正是通过调用 **errors.New** 函数来完成的。

结构体类型 **os.PathError** 是一个 **error** 接口类型的实现类型。声明及其方法如下：

```go
// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op  string  // “open” , ”unlink”, etc
    Path string // The associated file
    Err  error  // Returned by the system call
}

func (e *PathError) Error() string { 
    return e.Op + " " + e.Path + ": " + e.Err.Error() 
}
```

先判定获取到的 **error** 类型值的动态类型，再依次来进行必要的类型转换和后续操作。例如：

```go
file , err := os.Open("E:\\Software\\lgh.txt")
if err != nil {
    if pe, ok := err.(*os.PathError); ok {// 判断err是否为*os.PathError类型
        fmt.Printf("Path Error: %s \n(op=%s,path=%s)\n", pe.Err, pe.Op, pe.Path)    
    } else {
        fmt.Printf("Unknown Error: %s\n",err)
    }
}
```

如上**Open** 的参数的文件路径不存在，运行截图如下：
 
![](result-5.png)

在上面示例中的 **os.Open** 函数在执行过程中没有发生任何错误，就可以对变量 **file** 的内容进行读取了。例如：

```go
reader := bufio.NewReader(file) // 创建一个可以读取文件内容的读取器
var buf bytes.Buffer // 缓存从文件读取出来的内容
for {
    // reader读取器，返回3个结果值。reader类型所属的方法如下：
    // func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
    // 当读取器读到file所代表的文件的末尾时，ReadLine方法会直接将变量io.EOF的值作为它的第三个结果值err返回。
    byteArray, _, err1 := reader.ReadLine()
    if err1 != nil {
        // io.EOF是error类型的变量，在标准库代码包io中，它的声明如下：
        // var EOF = errors.New("EOF") ,EOF是文件结束符(End Of File)的缩写。
        // 严格来说，EOF并不应该算作一个真正的错误，而仅仅属于一种"错误信号"
        if err1 == io.EOF { // 判断读取器是否已经读到了文件的末尾
            break
        } else {
            fmt.Printf("Read Error: %s\n", err1)
            break
        }
    } else {
        buf.Write(byteArray)
    }
    fmt.Printf("%s\n", byteArray)
}
```

实现 **error** 接口类型的另一个技巧是，可以通过把 **error** 接口类型嵌入到新的接口类型中来对它进行扩展。标准库代码包 **net** 中的 **Error** 接口类型，声明如下：

```go
//An Error represents a network error
type Error interface {
    error
    Timeout() bool // Is the error a timeout?
    Temporary() bool // Is the error temporary?
}
```

一些在 **net** 包中声明的函数会返回动态类型为 **net.Error** 的 **error** 类型值。如果变量 **err** 的动态类型是 **net.Error**，可以根据它的 **Temporary** 方法的结果值来判断当前的错误状态是否临时的：

```go
if netErr, ok := err.(net.Error); ok && netErr.Temporary(){
    // 省略若干语句
}
```

如果是临时的，那么就可以间隔一段时间之后再对之前的操作进行重试，否则就记录错误状态的信息并退出。

### 2.2 panic

Go语言内建的一个专用函数，目的使编程人员能够在自己的程序中报告运行期间的，不可恢复的错误状态。**panic** 函数被用于停止当前的控制流程的执行并报告一个运行时恐慌。它可以接受一个任意类型的参数值，这个参数常常是一个 **string** 类型值或者 **error** 类型值。例如：

```go
package main

import (
    "errors"
)

func main(){
    outerFunc()
}

func outerFunc(){
    innerFunc()
}

func innerFunc(){
    panic(errors.New("A intended fatal error!"))
}
```

当在函数 **innerFunc** 中调用 **panic** 函数之后，函数 **innerFunc** 的执行会被停止。然后，流程控制权会被交回给函数 **innerFunc** 的调用方 **outerFunc** 函数，此时，**outerFunc** 函数的执行也将被停止。运行时恐慌就这样沿着调用栈反方向进行传达，直至到达当前 **Goroutine**（也被称为Go程，可以看作是一个能够独占一个系统线程并在其中运行程序的独立环境）调用栈的最顶层。这时，当前 **Goroutine** 的调用栈中的所有函数的执行都已经被停止了，意味着程序已经崩溃了。

运行时恐慌并不都是通过调用 **panic** 函数的方式引发的。它也可以由Go语言的运行时系统来引发。例如：

```go
myIndex := 4
ia := [3]int{1, 2, 3}
_ = ia[myIndex] // 产生了一个数组访问越界的运行时错误，会引发一个运行时恐慌。
```

如上这个运行时恐慌由运行时系统报告的，它相当于显示地调用 **panic** 函数并传入一个 **runtime.Error** 类型的参数值，该类型的声明如下：

```go
type Error interface {
    error 
    // RuntimeError is a no-op function but serves to distinguish types that are runtime errors
    // from ordinary errors: a type is a runtime error if it has a RuntimeError method.
    RuntimeError()
}
```

### 2.3 recover

运行时恐慌一旦被引发就会向调用方传递直至程序崩溃。Go语言提供了专用于“拦截”运行时恐慌的内建函数--- **recover**。它可以使当前的程序从运行时恐慌的状态中恢复并重新获得流程控制权。**recover** 函数有一个 **interface{}** 类型的结果值，如果当前的程序正处于运行时恐慌的状态下，那么调用 **recover** 函数将会得到一个 **非nil** 的 **interface{}** 类型值。如果当时的运行时恐慌是由Go语言的运行时程序引发的，就会获得一个 **runtime.Error** 类型的值。

只有在 **defer** 语句的延迟函数中调用 **recover** 函数才能够起到“拦截”运行时恐慌的作用。例如:

```go
defer func(){
    if r := recover(); r != nil {
        fmt.Printf("Recovered panic: %s\n", r)
    }
}()
```

再看如下一个示例，有助于理解 **panic** 函数、**recover** 函数和 **defer** 语句有关的运行机制。例如：

```go
package main

import (
    "fmt"
)

func main(){
    fetchDemo()
    // 由于运行时恐慌在将要被继续传递给fetchDemo函数的调用方的时候被“拦截”。
    // 因此fetchDemo函数的调用方（也就是main函数）得以重获流程控制权，下一条语句可以打印
    fmt.Println("The main function is executed.") 
}

func fetchDemo() {
    defer func() {
        if v := recover(); v != nil {
            fmt.Printf("Recovered a panic. [index=%d]\n", v)
        }
    }()
    ss := []string{"A", "B", "C"}
    fmt.Printf("Fetch the elements in %v one by one...\n", ss)
    fetchElement(ss, 0)
    fmt.Println("The elements fetching is done.")//上面的语句出现了运行时恐慌，因此不会执行
}

func fetchElement(ss []string, index int) (element string) {
    if index >= len(ss) {
        fmt.Printf("Occur a panic![index=%d]\n", index)
        panic(index)
    }
    fmt.Printf("Fetching the element... [index=%d]\n", index)
    element = ss[index]
    defer fmt.Printf("The elements is \"%s\". [index=%d]\n", element, index)
    fetchElement(ss, index + 1)
    return
}
```
如上命令源码文件运行结果截图：

![](result-6.png)

在Go语言标准库中，即使使用的某个程序实体的内部发生了运行时恐慌，这个运行时恐慌也会在被传递给我们编写的程序使用方之前被“平息”并以 **error** 类型值的形式返回给使用方。在这些标准库代码包中，往往都会有自己的 **error** 接口类型的实现。只有当调用 **recover** 函数得到的结果值的类型是它们自定义的 **error** 类型的实现类型的时候，才会去处理这个运行时恐慌，否则就会重新引发一个运行时恐慌（**re-panic**）并携带相同的值。

在标准库代码包 **fmt** 中 **scan.go** 的 **Token** 函数就是如下的这样处理运行时恐慌的。声明如下：

```go
func (s *ss) Token(skipSpace bool, f func(rune) bool) (tok []byte, err error) {
    defer func() {
        if e := recover(); e != nil {
            if se, ok := e.(scanError); ok {
                err = se.err
            } else {
                panic(e)
            }
        }
    }()
    // 省略若干条语句
}
```

在 **Token** 函数包含的延迟函数中，当运行时恐慌携带的值的类型是 **fmt.scanError** 类型的时候，这个值就会被赋值给代表结果值的变量 **err**，否则运行时恐慌就会被重新引发。

一个运行时恐慌无论重新引发几次，它所有的引发信息都依然会被提供在最终的程序崩溃报告中。重新引发一个运行时恐慌的时候使用如下：

```go
panic(e)
```

在使用Go语言编写程序时，在使用上面类似 **Token** 函数的惯用法之前应该明确和统一可以被立即处理和需要被重新引发的运行时恐慌的种类。一般情况下，如果携带的值是动态类型为 **runtime.Error** 的 **error** 类型值的话，这个运行时恐慌就应该被重新引发。从运行时恐慌的分类和处理决策角度看，在必要时自行定义一些 **error** 类型的实现类型是有好处的。

> **建议：** 对于运行时恐慌的引发，应该在遇到致命的、不可恢复的错误状态时才去引发一个运行时恐慌，否则，可以完全利用函数或方法的结果值来向程序使用方传达错误状态。另外，应该仅在程序模块的边界位置上的函数或方法中对运行时恐慌进行“拦截”和“平息”。

# 结语
本篇讲述了 **Go** 语言特殊流程控制方法 **defer** 、**error**、**panic**、**recover**，下篇开始我们了解 **Go** 语言程序测试的相关内容，敬请期待！！！！


最后附上知名的 **Go** 语言开源框架： 
**Gobot:** 一个非常有意思的开源项目。它旨在成为下一代自动机工程学框架。换句话说，我们可以用它来控制机器人！它已经支持了10个(或者更多)不同的硬件平台。这其中包括已经被国内的计算机硬件发烧友所熟知的 **Arduino** 。该开源项目的官方网址是 [http://gobot.io](http://gobot.io)
