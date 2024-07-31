---
title: Go语言学习15-基本流程控制
date: 2016-07-18 22:10:15
updated: 2024-03-24 21:21:11
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 基本流程控制
  - if
  - switch
  - for
  - goto
---



![](/images/go-logo.png)

# 引言
**Go** 语言在流程控制结构方面有些像C语言，但是在很多方面都与**C**不同。特点如下：

<!-- more -->

 - 在**Go**语言中没有 `do` 和 `while` 循环，只有一个更加广义的 `for` 语句。
 
 - **Go**语言中的 `switch` 语句更加灵活多变。**Go**语言的 `switch` 语句还可以被用于进行类型判断。
 
 - 与 `for` 语句类似，**Go**语言中的 `if` 语句和 `switch` 语句都可以接受一个可选的初始化子语句。
 
 - **Go**语言支持在 `break` 语句和 `continue` 语句之后跟一个可选的标记（Label）语句，以标识需要终止或继续的代码块。
 
 - **Go**语言中还有一个类似于多路转接器的 `select` 语句。
 
 - **Go**语言中的 `go` 语句可以被用于灵活地启用 `Goroutine`。
 
 - **Go**语言中的 `defer` 语句可以使我们更加方便地执行异常捕获和资源回收任务。

# 主要内容

## 1. 代码块和作用域

代码块就是一个由花括号 “{” 和 “}” 括起来的若干表达式和语句的序列。代码块中也可以不包含任何内容，即为空代码块。

在**Go**语言的源代码中，除了显式的代码块之外，还有一些隐式的代码块，如下：

 - 所有**Go**语言源代码形成了一个最大的代码块。这个最大的代码块也被称为全域代码块。
 
 - 每一个代码包都是一个代码块，即代码包代码块。它们分别包含了当前代码包内的所有**Go**语言源代码。
 
 - 每一个源码文件都是一个代码块，即源码文件代码块。它们分别包含了当前文件内的所有**Go**语言源码。
 
 - 每一个 `if` 语句、`for` 语句、`switch` 语句和 `select` 语句都是一个代码块。
 
 - 每一个在 `switch` 或 `select` 语句中的子句都是一个代码块。


在**Go**语言中，每一个标识符都有它的作用域。使用代码块表示词法上的作用域范围，规则如下：

 - 一个预定义标识符的作用域是全域代码块。
 
 - 代表了一个常量、类型、变量或函数（不包括方法）的，被声明在顶层的（即在任何函数声明之外被声明的）标识符的作用域是代码包代码块。
 
 - 一个被导入的代码包的名称的作用域是包含该代码包导入语句的源码文件代码块。
 
 - 一个代表了方法接收者、方法参数或方法结果的标识符的作用域是方法代码块。
 
 - 对于一个代表了常量或变量的标识符，如果它被声明在函数内部，那么它的作用域总是包含它的声明的那个最内层的代码块。
 
 - 对于一个代表了类型的标识符，如果它被声明在函数内部，那么它的作用域就是包含它的声明的那个最内层的代码块。
 
在**Go**语言中，可以在某个代码块中对一个已经在包含它的外层代码块中声明过的标识符进行重声明。这种情况下，在外层代码块中声明的那个同名标识符被屏蔽了。例如：

```go
package main

import (
    "fmt"
)

var v string = "1, 2, 3"

func main(){
    v := []int{1, 2, 3}
    if v != nil {  // 此时v代表的是一个切片类型值，因此可以与空值nil进行判等
        var v int = 123
        fmt.Printf("%v\n",v)
    }
}
```
运行结果截图如下：

 ![](result.png)

## 2. if 语句

**Go**语言的 `if` 语句总是以关键字 `if` 开始。之后，可以后跟一条简单语句（当然也可以没有），然后是一个作为条件判断的布尔类型的表达式以及一个用花括号 “{” 和 “}” 括起来的代码块。`if` 语句也可以由 `else` 分支，它是 `else` 关键字和一个用花括号 “{” 和 “}” 括起来的代码块。

常用的简单语句包括**短变量声明**、**赋值语句**和**表达式语句**。除了特殊的内建函数和代码包 `unsafe` 中的函数，针对其他函数和方法的调用表达式和针对通道类型的接收表达式都可以出现在语句上下文中。在必要时，还可以使用圆括号将它们括起来。其他的简单语句还包括**发送语句**、**自增/自减语句**和**空语句**。

**Go**语言 **if** 语句的举例：

```go
if 100 < number {
    number++
} else {
    number--
}
```

在上面的 `if` 语句的条件表达式 `100 < number` 并没有被圆括号括起来，这是**Go**语言的流程控制语句的特点之一。同时，强调一点是跟在**条件表达式**和 `else` 关键字之后的两个代码块必须由花括号 “{” 和 “}” 括起来，不论代码块中包含几条语句以及是否包含语句。

`if` 语句可以接受一条**初始化子语句**，常常用它来初始化一个变量如下：

```go
if diff := 100 – number; 100 < diff { // 初始化子句和条件表达式之间是需要用分号“;”分隔的
    number++
} else if 200 < diff {
    number--
} else {
    number -= 2
}

```

由于在**Go**语言中一个函数可以返回多个结果，因此常常会把在函数执行期间出现的常规错误也作为结果之一。例如，标准库代码包 `os` 中的函数 `Open`，它的声明如下：

```go
func Open(name string) (file *File, err error)
```

函数 `os.Open` 返回的第一个结果是与已经被“打开”的文件相对应的 `*File` 类型的值，而第二个结果是代表了常规错误的 `error` 类型的值。`error` 是一个预定义的接口类型，所有实现它的类型都应该被用于描述一个常规错误。
在导入代码包 `os` 之后，调用 `Open` 函数：

```go
f, err := os.Open(name)
if err != nil {
    return err
}
```

如上调用后，先检查 `err` 的值是否为 `nil`。如果变量 `err` 的值不为 `nil`，那么说明 `os.Open` 函数在被执行过程中发生了错误，这时变量 `f` 的的值肯定是不可用的。


在**Go**语言中，`if` 语句常被作为**卫述语句**。**卫述语句**是指被用来检查关键的先决条件的合法性并在检查未通过的情况下立即终止当前代码块的执行的语句。上面调用 `Open` 函数之后检查的 `if` 语句就是**卫述语句**的一种。

```go
func update (id int, department string) bool {
    if id <= 0 {
        return false
    }
    // 省略若干语句
    return true
} // 在函数update开始处的那条if语句就属于卫述语句。
```

对函数稍加改造如下：

```go
func update (id int, department string) error{ // 需要事先导入标准库的代码包errors
    if id <= 0 {
        return errors.New("The id is INVALID!")
    }
    // 省略若干语句
    return nil
} // update函数返回的结果不但可以表示在函数执行期间是否发生了错误，而且还可以体现出错误的具体描述。
```

## 3. switch 语句

语句 `switch` 可以使用表达式或者类型说明符作为 `case` 判定方法。`switch` 语句也就可以被分为两类：**表达式switch语句** 和 **类型switch语句**。

### 3.1 表达式switch语句

在表达式 `switch` 语句中，`switch` 表达式和 `case` 携带的表达式（也称为 `case` 表达式）都会被求值。对这些表达式的求值是自左向右、自上而下进行的。如果在 `switch` 语句中没有显示的 `switch` 表达式，那么 `true` 将会被作为 `switch` 表达式。例如：

```go
switch content {
default:
    fmt.Println("Unknown language")
case "Python":
    fmt.Println("A interpreted language")
case "**Go**":
    fmt.Println("A compiled language")
}
```

`switch` 关键字之后会紧跟一个 `switch` 表达式。`switch` 表达式中涉及的标识符都必须是已经被声明过的。同时还可以在 `switch` 和 `switch` 表达式之间插入一条简单语句，如下：

```go
switch content := getContent(); content { // content := getContent()会在switch表达式content被求值之前被执行
default:
    fmt.Println("Unknown language")
case "Python":
    fmt.Println("A interpreted language")
case "Go":
    fmt.Println("A compiled language")
}
```

一条 `case` 语句由一个 `case` 表达式和一个语句列表组成，并且这两者之间需要用冒号 `:` 分隔。一个 `case` 表达式由一个 `case` 关键字和一个表达式列表（可以包含多个表达式）组成。例如：

```go
switch content := getContent(); content {
default:
    fmt.Println("Unknown language")
case "Python", "Ruby":
    fmt.Println("A interpreted language")
case "Go", "C", "Java":
    fmt.Println("A compiled language")
}
```

在一条 `case` 语句中的语句列表的最后一条语句可以是 `fallthrough` 语句。一条 `fallthrough` 语句会将流程控制权转移下一条 `case` 语句上。例如：

```go
switch content := getContent(); content {
default:
    fmt.Println("Unknown language")
case "Ruby":
    fallthrough
case "Python":
    fmt.Println("A interpreted language")
case "Go", "C", "Java":
    fmt.Println("A compiled language")
}
```

如上当变量 `content` 的值与 `"Ruby"` 相等的时候，在标准输出上打印出的内容会是 `A interpreted language`。`fallthrough` 语句只能够作为 `case` 语句中的语句列表的最后一条语句, `fallthrough` 语句不能出现在最后一条 `case` 语句的语句列表中。

`break` 语句也可以出现在 `case` 语句列表中。一条 `break` 语句由一个 `break` 关键字和一个**可选的标记**组成，这两者之间用空格分隔。例如：

```go
switch content := getContent(); content {
default:
    fmt.Println("Unknown language")
case "Ruby":
    break
case "Python":
    fmt.Println("A interpreted language")
case "Go", "C", "Java":
    fmt.Println("A compiled language")
}
```

### 3.2 类型switch语句

类型 `switch` 语句将对类型进行判定，而不是值。它的 `switch` 表达式的表现形式与类型断言表达式相似，但与类型断言表达式不同的是，它使用关键字 `type` 来充当欲判定的类型，而不是使用一个具体的类型字面量。例如：

```go
switch v.(type){
    case string:
        fmt.Printf("The string is '%s'.\n", v.(string)) // v.(string)把v的值转换成了string类型的值
    case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64: // v是byte类型或rune类型，也会执行下列分支
        fmt.Printf("The string is %d.\n", v)
    default:
        fmt.Printf("Unsupported value. (type=%T)\n", v)
    }
}
```

在类型 `switch` 语句中，`case` 表达式中的类型字面量可以是 `nil`，如果 `v` 的值是 `nil`，那么表达式 `v.(type)` 的结果值也会是 `nil`。与表达式 `switch` 语句不同的是，`fallthrough` 语句不允许出现在类型 `switch` 语句中。

对**类型 switch 语句**的 `switch` 表达式进行变形：

```go
switch i:= v.(type){
    case string:
        fmt.Printf("The string is '%s'.\n", i)
    case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
        fmt.Printf("The integer is %d.\n", i)
    default:
        fmt.Printf("Unsupported value. (type=%T)\n", i)
    }
}
```

第一个`case`语句相当于：

```go
case string:
    i := v.(string)
    fmt.Printf("The string is '%s'.\n", i)
```

对于包含多个类型字面量的 `case` 表达式，比如第二个 `case` 语句。例如，如果上面v的动态类型是 `uint16` 类型，那么第二个 `case` 语句相当于：

```go
case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
    i := v.(uint16)
    fmt.Printf("The integer is %d.\n", i)
```

如上通过这种方式后，不需要在每个 `case` 语句中分别对那个欲判定类型的值进行显示地类型转换了。


在 `switch` 表达式缺失的情况下，该 `switch` 语句的判定目标会被视为布尔类型 `true`，也就是所有 `case` 表达式的结果值都应该是布尔类型。例如：

```go
switch {
    case number < 100:
        number++
    case number < 200:
        number--
    default:
        number -= 2
    }
}
```

同样可以在 `switch` 关键字和 `switch` 表达式中添加一条简单语句，例如：

```go
switch number := 123; { // 这里switch表达式缺失，默认switch的判定目标为布尔类型
    case number < 100:
        number++
    case number < 200:
        number--
    default:
        number -= 2
    }
}
```

## 4. for 语句

### 4.1 for 子句

`for` 子句的3个部分是由固定顺序组成，即**初始化子句在左**，**条件在中**，**后置子句在右**，并且在它们之间需要用分号“;”来分隔。可以在编写 `for` 子句的时候省略掉其中的任何部分，为了避免歧义，与省略部分相邻的分隔符“;”也必须被保留。初始化子句总会在充当条件的表达式被第一次求值之前执行，且只会执行一次，而后置子句的执行总会在每次代码块执行完之后紧接着进行。后置子句一定不能是短变量声明。例如：

```go
for i := 0; i < 100; i++ {
    number++
}

var j uint = 1
for ; j%5 != 0; j *= 3 { // 省略初始化子句
    number++
}

for k != 1; k%5 != 0; { // 省略后置子句
    k *= 3
    number++
}
```

在 `for` 子句的初始化子句和后置子句同时被省略或者其中的所有部分都被省略的情况下，分隔符 `;` 可以被省略。例如：

```go
// number是一个int类型的变量
for number < 200 {
    number += 2
}
```

当 `for` 子句的3个部分都省略，就陷入了死循环。例如：

```go
for {
    number++
}
```

### 4.2 range 子句

`for` 语句使用 `range` 子句可以迭代出一个**数组**或**切片**值中的每个元素，一个**字符串**值中的每个字符或者一个**字典**值中的每个键值对，甚至可以被用于持续接收一个**通道类型**值中的元素。例如：

```go
ints := []int{1, 2, 3, 4, 5}
for i, d := range ints {
    fmt.Printf("%d: %d\n", i, d)
}
```

事先声明了标识符，例如：

```go
var i, d int
ints := []int{1, 2, 3, 4, 5}
for i, d = range ints {
    fmt.Printf("%d: %d\n", i, d)
}
```
运行截图如下：
 
![](result-1.png)

### 4.3 range 子句的迭代产出

|`range` 表达式的类型  |  第一个产出值| 第二个产出值（若显示获取）| 备注|
|:------------------------------|:----------------------|:--------------------------------------------|:--------|
|a：[n]T、*[n]T或[]T |i：int类型的元素索引值|与索引对应的元素的值a[i],类型为T|a为range表达式的结果值,N为数组类型的长度,T为数组类型或切片类型的元素类型|
|s：string类型 |i：int类型的元素索引值|与索引对应的元素的值s[i],类型为rune|s为range表达式的结果值|
|m：map[K]T |k：键值对中的键的值，类型为K|与键对应的元素值m[k],类型为T|m为range表达式的结果值,K为字典类型的键的类型,V为字典类型的元素类型|
|c：chan T|e： 元素的值，类型为T|  无   |c为range表达式的结果值，T为通道类型的元素类型|



对于所有可迭代的数据类型的值来说，可以要求每次迭代只产出第一个迭代值。例如：

```go
m := map[uint]string{1: "A", 6: "C", 7: "B"}
var maxKey uint
for k := range m{
    fmt.Printf("k: %d\n", k)
    if k > maxKey {
        maxKey = k
    }
}
fmt.Printf("maxKey: %d\n", maxKey)
```
运行截图如下：
 
![](result-2.png)


忽略第一个迭代值而只使用第二个迭代值的方法，如下：

```go
m := map[uint]string{1: "A", 6: "C", 7: "B"}
var values []string
for _, v := range m {
    values = append(values, v)
}
fmt.Printf("values: %v\n", values)
```
运行截图如下：
 
![](result-3.png)

在 `for` 语句中，可以使用 `break` 语句来终止 `for` 语句的执行。例如：

```go
// 该变量的值包含了某个网络的所有用户昵称及其重复次数
// 这个字典的键表示用户昵称，而值则代表了使用该昵称的用户数量。
var namesCount map[string]int = map[string]int{"霓虹": 3,"Huazie": 1, "Tom": 2, "诗": 4}
// 存储查找到所有的只包含中文的用户昵称的计数信息。
targetsCount := make(map[string]int)
for k, v := range namesCount {
    matched := true
    for _, r := range k {
        if r < '\u4e00' || r > '\u9fbf' { // 用户昵称中包含了非中文字符
            matched = false
            break // 只会终止直接包含它的那条for语句的执行
        }
    }
    if matched {
        targetsCount[k] = v
    }
}
fmt.Printf("targetsCount: %v\n", targetsCount)
```
运行截图如下：

![](result-4.png)

### 4.4 标记语句

一条标记语句可以成为 `goto` 语句、`break` 语句或 `continue` 语句的目标。标记语句中的标记只是一个标识符，它可以被放置在任何语句的左边以作为这个语句的标签。标记和被标记的语句之间需要用冒号来分隔。例如：

**（1）break和标记语句的使用**

```go
    var namesCount map[string]int = map[string]int{"霓虹": 3,"Huazie": 1, "Tom": 2, "诗": 4}
    targetsCount := make(map[string]int)
L:
    for k, v := range namesCount {
        for _, r := range k {
            if r < '\u4e00' || r > '\u9fbf' {
                break L // 发现第一个非全中文的用户昵称的时候停止查找
            }
        }
        targetsCount[k] = v
    }
    fmt.Printf("targetsCount: %v\n", targetsCount)
```

运行截图如下：
  
![](result-5.png)

在**Go**语言中 `continue` 语句只能在 `for` 语句中被使用。`continue` 语句会使直接包含它的那个 `for` 循环直接进入下一次迭代。

**（2）continue与标记语句的使用**

```go
    var namesCount map[string]int = map[string]int{"霓虹": 3,"Huazie": 1, "Tom": 2, "诗": 4}
    targetsCount := make(map[string]int)
L:
    for k, v := range namesCount {
        for _, r := range k {
            if r < '\u4e00' || r > '\u9fbf' {
                continue L // L标记代表的那个for循环直接进入下一次迭代
            }
        }
        targetsCount[k] = v
    }
    fmt.Printf("targetsCount: %v\n", targetsCount)
```

运行截图如下：
 
![](result-6.png)

使用**Go**语言的 `for` 语句写出反转一个切片类型值中的所有元素值的代码。（**不使用** 在 `for` 语句之外声明的任何变量作为辅助）：

```go
var numbers []int = []int{1,2,3,4,5}
fmt.Printf("before numbers: %v\n", numbers)
for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
    numbers[i], numbers[j] = numbers[j], numbers[i]
}
fmt.Printf("after numbers: %v\n", numbers)
```

运行截图如下：
 
![](result-7.png)

**初始化子句** 和 **后置子句** 的只能是单一语句而不能是多个语句，但是可以使用**平行语句**来丰富两个子句的语义。

## 5. goto 语句

一条 `goto` 语句会把流程控制权无条件地转移到它右边的标记所代表的语句上。`goto` 语句只能与标记语句连用,并且在它的右边必须要出现一个标记。

`goto` 语句使用的注意点：

(1) 不允许因使用 `goto` 语句而使任何本不在当前作用域中的变量进入该作用域。例如：

```go
    goto L
    v := "B"
L:
    fmt.Printf("V: %v\n", v)
```

这段代码会造成一个编译错误，主要原因是语句 `goto L` 恰恰使变量 `v` 的声明语句被跳过了。

修改上面的代码，保证顺利通过编译，如下：

```go
    v := "B"
    goto L
L:
    fmt.Printf("V: %v\n", v)
```

把某条 `goto` 语句的直属代码块叫作代码块 `A`，而把该条 `goto` 语句右边的标记所指代的那条标记语句的直属代码块叫作代码块 `B`，那么只要代码块 `B` 不是代码块 `A` 的外层代码块，这条 `goto` 语句就是不合法的。例如：

```go
var n int = 10
if n % 3 != 0 {
    goto L1
}
switch {
case n % 7 == 0:
    fmt.Printf("%v is a common multiple of 7 and 3.\n", n)
default:
L1:
    fmt.Printf("%v isn't multiple of 3.\n", n)
}
```

如上，标记 `L1` 所指代的标记语句的直属代码块是由 `switch` 语句代表的，而 `goto L1` 语句的直属代码块是由 `if` 语句代表的，并且前者并不是后者的直属代码块。因此，`goto L1` 是非法的。
上面的代码会出现编译错误。修正上面的编译错误，代码如下：

```go
    var n int = 10
    if n % 3 != 0 {
        goto L1
    }
    switch {
    case n % 7 == 0:
        fmt.Printf("%v is a common multiple of 7 and 3.\n", n)
    default:
    }
L1:
    fmt.Printf("%v isn't multiple of 3.\n", n)
```

利用 `goto` 语句跳出嵌套的流程控制语句的执行。例如：

```go
// 查找name中的第一个非法字符并返回
// 如果返回的是空字符串说明name中不包含任何非法字符
func findEvildoer(name string) string{
    var evildoer string    
    for _, r := range name{
        switch {
        case r >= '\u0041' && r <= '\u005a': // a-z
        case r >= '\u0061' && r <= '\u007a': // A-Z
        case r >= '\u4e00' && r <= '\u9fbf': // 中文字符
        default:
            evildoer = string(r)
            goto L2
        }
    }
    goto L3
L2:
    fmt.Printf("The first evildoer of name '%s' is '%s'!\n", name, evildoer)
L3:
    return evildoer
}
```

如下使用 `break` 和 `if` 语句替换上面的两条 `goto` 语句:

```go
func findEvildoer1(name string) string{
    var evildoer string    
L2:
    for _, r := range name{
        switch {
        case r >= '\u0041' && r <= '\u005a': // a-z
        case r >= '\u0061' && r <= '\u007a': // A-Z
        case r >= '\u4e00' && r <= '\u9fbf': // 中文字符
        default:
            evildoer = string(r)
            break L2
        }
    }
    if evildoer != "" {
        fmt.Printf("The first evildoer of name '%s' is '%s'!\n", name, evildoer)
    }
    return evildoer    
}
```

(2) 另一个适合使用 `goto` 语句的场景是集中式的错误处理。例如：

```go
func checkValidity(name string) error{
    var errDetail string    
    for i, r := range name{
        switch {
        case r >= '\u0041' && r <= '\u005a': // a-z
        case r >= '\u0061' && r <= '\u007a': // A-Z
        case r >= '\u4e00' && r <= '\u9fbf': // 中文字符
        case r == '_' || r == '-' || r == '.': // 其他允许的符号
        default:
            errDetail = "The name contains some illegal characters."
            goto L3
        }
        if i == 0 {
            switch r {
            case '_':
                errDetail = "The name can not begin with a '_'."
                goto L3
            case '-':
                errDetail = "The name can not begin with a '-'."
                goto L3
            case '.':
                errDetail = "The name can not begin with a '.'."
                goto L3
            }
        }
    }
    return nil
L3:
    return errors.New("Validity check failure: "+errDetail)    
}
```


**Go** 语言可以方便地从错综复杂的流程控制中跳出，但是 `goto` 语句的代码块的可读性大大下降。因此，要节制地使用 `goto` 语句。

# 结语
本篇讲述了 **Go** 语言的基本流程控制，下篇继续讲解 **Go** 语言流程控制方法中的一些特殊流程控制语句。

最后附上知名的 **Go** 语言开源框架： 

**Groupcache:** 著名的内存缓存系统 **Mencached** 的作者用**Go**语言编写的一个与前者功能类似的函数库。作者想用它作为 **Memcached** 的 **替代者**。其源码：https://github.com/golang/groupcache
