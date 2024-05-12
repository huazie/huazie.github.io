---
title: Go语言学习20-测试运行记录和覆盖率
date: 2016-07-27 21:32:27
updated: 2024-04-13 21:25:25
categories:
  - [开发语言-Go,Go语言学习]
tags:
  - Go
  - 测试运行记录
  - 测试覆盖率
---

[《开发语言-Go》](/categories/开发语言-Go/) [《Go语言学习》](/categories/开发语言-Go/Go语言学习/) 

![](/images/go-logo.png)

# 1. 测试运行记录

在硬件环境方面，主要考察计算机的 **负载状况**，比如 **CPU使用率**、**内存使用率**、**磁盘使用情况** 等。
在软件系统方面，主要包括 **内存分配**、**并发处理数量** 及 **死锁** 等情况。

<!-- more -->

## 1.1 收集资源使用情况

*与收集资源使用情况有关的标记*
|标记名称       |     标记描述           |
|:-------------|:---------------------|
|-cpuprofile cpu.out|记录CPU使用情况，并写到指定的文件中，直到测试退出。<br />cpu.out作为指定文件的文件名可以被其他任何名称代替|
|-memprofile mem.out|记录内存使用情况，并在所有测试通过后将内存使用概要写到指定的文件中。<br />mem.out作为指定文件的文件名可以被其他任何名称代替|
|-memprofilerate n|此标记控制着记录内存分配操作的行为，这些记录将会被写到内存使用概要文件中。<br />N代表着分析器的取样间隔，单位为字节。每当有n个字节的内存被分配时，分析器就会取样一次 |



现在执行 **-cpuprofile** 标记的 **go test** 命令来运行标准库的 **net** 代码包中的测试。例如：

```bash
E:\Software\Go\goc2p\pprof>go test -cpuprofile cpu.out net
```

**go test** 命令程序不允许在多个代码包测试的时候使用 **-cpuprofile** 标记，如下：

![](result.png)

在上面的测试运行完成后，在当前目录中生成了一个 **net.test.exe** 的可执行文件和 **cpu.out** 的文件，可以使用 **go tool pprof** 命令来交互式的对这个概要文件进行查阅，如下：

![](result-1.png)
 

当 **-memprofile** 标记有效时，测试运行程序会在运行测试的同时记录它们对内存的使用情况（程序运行期间的堆内存的分配情况）。

使用 **-memprofilerate** 标记来设置分析器的取样间隔，单位是字节。它的值越小就意味着取样效果越好，因为取样间隔会更短。在 **testing** 包内部，**-memprofilerate** 标记的值会赋给 **runtime** 包中的 **int** 类型的变量 **MemProfileRate**，该变量默认值为 **512 * 1024**，即 **512K** 字节。如果设置为 **0**，则代表 **停止取样**。

**-memprofile** 标记和 **-memprofilerate** 标记的测试命令语句，如下：

```bash
E:\Software\Go\goc2p\pprof>go test -memprofile mem.out -memprofilerate 10 net
```


执行上面命令之后，在当前目录中生成了一个 **net.test.exe** 的可执行文件和 **cpu.out** 的文件,可以使用 **go tool pprof** 命令来交互式的对这个概要文件进行查阅，如下：
 
![](result-2.png)

## 1.2 记录程序阻塞事件

*与记录程序阻塞事件有关的标记*
|标记名称           |     标记描述         |
|:----------------|:---------------------|
|-blockprofile block.out |记录Goroutine阻塞事件，并在所有测试通过后将概要信息写到指定的文件中。<br />block.out作为指定文件的文件名可以被其他任何名称代替|
|-blockprofilerate n   |该标记用于控制记录Goroutine阻塞事件的时间间隔，n的单位为次，默认值为1|


**-blockprofile** 和 **-blockprofilerate** 标记的测试命令语句，如下：

```bash
E:\Software\Go\goc2p\pprof>go test -blockprofile block.out -blockprofilerate 100 net
```

**-blockprofilerate** 标记的值也是通过标准库代码包 **runtime** 中的 **API**----函数 **SetBlockProfileRate** ----传递给Go运行时系统的。当传入的参数为 **0** 时，就相当于彻底取消记录操作。当传入的参数为 **1** 时，每个阻塞事件都将被记录。**-blockprofilerate** 标记的默认值为 **1**。


# 2. 测试覆盖率

测试覆盖率是指，作为被测试对象的代码包中的代码有多少在刚刚执行的测试中被使用到。如果执行的该测试致使当前代码包中的90%的语句都被执行了，那么该测试的测试覆盖率就是90%。

## 2.1 与测试覆盖率相关标记
**go test** 命令可接受的与测试覆盖率有关的标记
|   标记名称   | 使用示例    |      说明     |
|:------------|:----------|:--------------|
|-cover        |-cover         |启用测试覆盖率分析|
|-covermode     |-covermode=set  |自动添加 **-cover** 标记并设置不同的测试覆盖率统计模式，支持的模式共有以下3个。<br />**set**：只记录语句是否被执行过 <br /> **count**: 记录语句被执行的次数 <br /> **atomic**: 记录语句被执行的次数，并保证在并发执行时也能正确计数，但性能会受到一定的影响 <br /> 这几个模式不可以被同时使用，在默认情况下，测试覆盖率的统计模式是**set**|
|-coverpkg|-coverpkg bufio,net|自动添加 **-cover** 标记并对该标记后罗列的代码包中的程序进行测试覆盖率统计。<br />在默认情况下，测试运行程序会只对被直接测试的代码包中的程序进行统计。<br />该标记意味着在测试中被间接使用到的其他代码包中的程序也可以被统计。<br />另外，代码包需要由它的导入路径指定，且多个导入路径之间以逗号“，”分隔。|
|-coverprofile|-coverprofile cover.out|自动添加 **-cover** 标记并把所有通过的测试的覆盖率的概要写入指定的文件中|


在 **go test** 命令后加入 **-cover** 标记以开启测试覆盖率统计。如下：

 ![](result-3.png)

标记 **-coverpkg** 可以获得间接被使用的代码包中的程序在测试期间的执行率。如下：

![](result-4.png) 

在使用 **-coverpkg** 标记来指定要被统计的代码包之后，未被指定的代码则肯定不会被统计，即使是被直接测试的那个代码包。

标记 **-coverprofile** 会使测试运行程序把测试覆盖率的统计信息写入到指定的文件中。该文件会被放在执行 **go test** 命令时的那个目录下。如下：

![](result-5.png) 

对于 **cover.out** 的内容可以使用Go语言提供的 **cover** 工具查看，通过 **go tool cover** 命令来运行它。**cover** 工具的两个功能：

 1. 根据指定的规则重写某一个源码文件中的代码，并输出指定的目标上。
 
 2. 读取测试覆盖率的统计信息文件，并以指定的方式呈现。

对于第一个功能，运行附带了测试覆盖率相关标记的 **go test** 命令之后，测试运行程序会使用 **cover** 工具在被测试的代码包中的非测试源码文件被编译之前重写它们。重写的方式会由运行 **go test** 命令时的 **-covermode** 标记所指定的测试覆盖率统计模式决定。在默认情况下，这个统计模式为 **set** 。

## 2.2 set 统计模式

在 **testing/ct** 下，有如下的源码文件 **ct_demo.go**，如下：

```go
package ct

func TypeCategoryOf(v interface{}) string {
    switch v.(type) {
    case bool:
        return "boolean"
    case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
        return "integer"
    case float32, float64:
        return "float"
    case complex64, complex128:
        return "complex"
    case string:
        return "string"
    }
    return "unknown"//如果只是非基本类型，统一返回“unknown”
}
```

在运行 **go test –cover** 命令之后,测试运行程序 **cover** 工具会把 **TypeCategoryOf** 函数重写成如下：

```go
func TypeCategoryOf(v interface{}) string {
    GoCover.Count[0] = 1
    switch v.(type) {
    case bool:
        GoCover.Count[2] = 1
        return "boolean"
    case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
        GoCover.Count[3] = 1
        return "integer"
    case float32, float64:
        GoCover.Count[4] = 1
        return "float"
    case complex64, complex128:
        GoCover.Count[5] = 1
        return "complex"
    case string:
        GoCover.Count[6] = 1
        return "string"
    }
    GoCover.Count[1] = 1
    return "unknown"
}
```

原先的源代码中的每一个流程分支上都会被安插一个计数器（由 **GoCover.Count** 代表）。每当某个流程分支被执行的时候，相应的计数器就会被赋值为 **1**，这样的方式只能表示该分支是否被执行过。

变量 **GoCover** 代表一个匿名的结构体类型的值，它的值用于存储每个分支的计数值、每个分支的起始和结束位置在当前源码文件中的行数以及每个分支中的语句的条数等信息。变量 **GoCover** 暂且可以称为执行计数变量，它的名称是可以通过 **-var** 标记告知给 **cover** 工具。

**go test** 命令的测试目标可以是一个或多个代码包，而 **cover** 工具每次只能重写一个源码文件。如果被测试的源码文件有多个，那么 **go test** 命令会对每一个文件都运行一次 **cover** 工具。为了使被重写的每一个源码文件中的执行计数变量的名称不重复，**go test** 命令针对不同的源码文件传递给 **cover** 工具的 **-var** 标记的值也是不同的。

## 2.3 count 统计模式

在**set**统计模式下的如下语句：

```go
GoCover.Count[0] = 1
```

会被写为

```go
GoCover.Count[0]++
```

这就相当于对分支的每一次执行都进行了记录。

## 2.4 atomic统计模式

**set**统计模式的语句会被写成如下：

```go
_cover_atomic_.AddUint32(&CoverVar.Count[1], 1)
```

其中 **_cover_atomic_** 是代码包 **sync/atomic** 的别名，需要在源代码中导入如下语句：

```go
import _cover_atomic_ "sync/atomic"
```

**atomic** 统计模式一般只在存在并发执行的应用场景下才被使用，因为原子操作的执行会带来一定的性能损耗。而 **set** 和 **count** 统计模式则比它通用的多，它们对于原先的程序的执行成本的影响也小很多。

通过 **-mode** 标记把需要使用的测试覆盖率统计模式直接传递给 **cover** 工具。它的用法和含义与 **go test** 命令的 **-covermode** 标记一致。实际上，**go test** 命令会把 **-covermode** 标记的值原封不动地作为运行 **cover** 工具是传递给它的 **-mode** 标记的值。但是，**-mode** 标记并没有默认值，因此使用 **cover** 工具对某个源代码文件进行重写的时候必须添加 **-mode** 标记。

在默认情况下，被重写的源代码会输出到标准输出上。但是也可以使用 **-o** 标记把这些代码存放到指定的文件中。

如下使用 **cover** 工具对工作区 **src** 目录下的 **testing\ct** 代码包下的 **ct_demo.go** 进行重写：

```bash
E:\Software\Go\goc2p\src\testing\ct>go tool cover -mode=set -var="GoCover" -o ct _out.go ct_demo.go
```

这时 **ct_out.go** 就是 **ct_demo.go** 被重写后的源码文件，各位可以自行运行查看。

对于 **cover** 工具的第二个功能：
查看之前 **go test cnet/ctcp -coverprofile=cover.out** 命令为 **cnet/ctcp** 代码包生成了一个测试覆盖率的概要文件 **cover.out**，如下运行截图：

![](result-6.png)

标记 **-func** 可以让 **cover** 工具把概要文件中包含的每个函数的测试覆盖率概要信息打印到标准输出上。

在上面的这段输出内容中，除了最后一行，每一行的内容都包括了3项信息，分别是：函数所在的源码文件的相对路径、函数名称和函数中被测试到的语句的数量占该函数的语句总数的百分比。而最后一行内容中的百分比则是被测试代码包中的所有测试到的语句的数量占该代码包中的语句总数的百分比。

现在使用 **-html** 标记，可以看到更加图形化的信息来直观的反应统计情况。如下：

```bash
go tool cover -html=cover.out
```

该命令会立即返回并且在标准输出上也并不会出现任何内容。取而代之的是当前操作系统的默认浏览器会被启动并显示 **cover** 工具刚刚根据概要文件生成的 **HTML** 格式的页面文件。如下：

![](result-7.png)

在这个 **HTML** 页面中，被测试到的语句以绿色显示，未被测试到的语句以红色显示，而未参加测试覆盖率计算的语句则用灰色表示。左上角还可以通过下拉框选择被测试代码包中的不同源码文件以查看它们的测试覆盖率情况。把鼠标的光标停留在绿色的语句上面的时候，在光标附近还会出现该语句执行次数的数字，各位可以自己试试。

上面的 **HTML** 页面文件展示的是在 **set** 统计模式下生成的概要文件。在 **count** 和 **atomic** 统计模式下，各位可以自己试试，参照如下：
 
![](result-8.png)


**cover工具可接受的标记**

|   标记名称   |     使用示例    |      说明     |
|:------------|:--------------|:-------------|
|-func         |-func=cover.out  |根据概要文件（即cover.out）中的内容，输出每一个被测试函数的测试覆盖率概要信息|
|-html         |-html=cover.out  |把概要文件中的内容转换成HTML格式的文件，并使用当前操作系统中的默认网络浏览器查看它|
|-mode        |-mode=count    |被用于设置测试概要文件的统计模式，详见go test命令的-covermode标记|
|-o           |-o=cover.out     |把重写后的源代码输出到指定文件中，如果不添加此标记，那么重写后的源代码会输出到标准输出上|
|-var          |-var=GoCover    |设置被添加到原先的源代码中的额外变量的名称|


# 总结
本篇介绍了 `Go` 语言的测试运行记录和覆盖率，下一篇介绍 `Go` 语言的程序文档，敬请期待！






