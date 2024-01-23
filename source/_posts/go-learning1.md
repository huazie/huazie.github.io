---
title: Go语言学习1-基础入门
date: 2016-06-27 09:28:36
updated: 2024-01-16 22:11:19
categories:
- [开发语言-Go]
tags:
- Go
- Go语言安装和配置
- Go工作区
- Go源码文件
- Go代码包
---

[《开发语言-Go》](/categories/开发语言-Go/)

![](/images/go-logo.png)

# 引言
本篇介绍 `Go` 语言的基础入门知识，如下贴出了 `Go` 语言相关的网站，方便大家学习

Go语言官方网站（[http://golang.org](http://golang.org)）
代码包文档网站（[http://godoc.org](http://godoc.org)）
Go语言中文网（[http://studygolang.com](http://studygolang.com)）

Go语言开发包下载路径：
[https://golang.google.cn/dl/](https://golang.google.cn/dl/)

《Go并发编程实战》所用到的源码实例下载路径：
[https://github.com/hyper-carrot/goc2p](https://github.com/hyper-carrot/goc2p)

# 1. Go语言配置环境变量

**windows** 下：

```
GOROOT={你的Go语言的根目录}
# 在环境变量PATH后追加
;%GOROOT%\bin
```

**linux** 下：
**Go** 语言官方建议把 **go** 文件夹复制到 `/usr/local`目录中，但也可以复制到其他目录；编辑 `/etc/profile` 文件，如下：

```bash
export GOROOT=/usr/local/go
export PATH=\$PATH:\$GOROOT/bin
```

保存 `/etc/profile` 文件，使用 **source** 命令使配置生效。
```bash
source /etc/profile
```

> 注意:  路径连接符 **windows** 下是 `\`，linux下是 `/`

**Go** 语言还有两个隐含的环境变量---- **GOOS** 和 **GORACH**
- **GOOS** 代表程序构建环境的目标操作系统，其值可以是 **liunx**，**windows**，**freebsd**，**darwin**；
- **GORACH** 代表程序构建环境的目标计算架构，其值可以是**386**，**amd64** 或 **arm**；

之后提到的 **平台相关目录** 是通过 **\${GOOS}_\${GORACH}** 的方式来命名的。（如 **Go** 归档文件的存放路径就是根据 **“平台相关目录”** 来指定的）

设置好环境变量后，在命令行中输入 **go** 出现如下信息，表示成功。

```cmd
Go is a tool for managing Go source code.

Usage:

        go command [arguments]

The commands are:

        build       compile packages and dependencies
        clean       remove object files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         run go tool fix on packages
        fmt         run gofmt on package sources
        generate    generate Go files by processing source
        get         download and install packages and dependencies
        install     compile and install packages and dependencies
        list        list packages
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         run go tool vet on packages

Use "go help [command]" for more information about a command.

Additional help topics:

        c           calling between Go and C
        buildmode   description of build modes
        filetype    file types
        gopath      GOPATH environment variable
        environment environment variables
        importpath  import path syntax
        packages    description of package lists
        testflag    description of testing flags
        testfunc    description of testing functions

Use "go help [topic]" for more information about that topic.

```

# 2. Go语言安装目录介绍
不管是安装版还是免安装版，**Go** 语言安装目录的文件夹名都是相同，下面我们一一来介绍下【如下图是 **go1.6.3.windows-amd64** 的】：

![](go-directory.png)


| 文件夹名 |  描述|
|--|:--|
| **api**  |  存放 **Go API** 检查器的辅助文件。<br/>其中，**go1.1.txt**、**go1.2.txt**、**go1.3.txt**和 **go1.txt** 等文件分别罗列了不同版本的 **Go** 语言的全部 **API** 特征; <br/>**except.txt** 文件中罗列了一些(在不破坏兼容性的前提下）可能会消失的 **API** 特性; <br/>**next.txt** 文件则列出了可能在下一个版本中添加的新 **API** 特性。<br/>![](go-version.png)|
| **bin** | 存放所有由官方提供的 **Go** 语言相关工具的可执行文件。<br/>默认情况下，该目录会包含 **go**、**godoc** 和 **gofmt** 这3个工具。<br/>![](go-tools.png)|
|**blog**|  用于存放官方博客中的所有文章，这些文章都是 **Markdown** 格式的。|
|**doc**  |  存放 **Go** 语言几乎全部的 **HTML** 格式的官方文档和说明，方便开发者在离线时查看。|
|**lib**| 用于存放一些特殊的库文件【如时区相关】。|
| **misc** | 存放各类编辑器或 **IDE**(集成开发环境）软件的插件，辅助它们查看和编写Go代码。<br/>有经验的软件开发者定会在该文件夹中看到很多熟悉的工具。 |
| **pkg** | 用于在构建安装后，保存 **Go** 语言标准库的所有归档文件。<br/>**pkg** 文件夹包含一个与 **Go** 安装平台相关的子目录，我们称之为“平台相关目录”。<br/> 例如，在针对 **Linux 32bit** 操作系统的二进制安装包中，平台相关目录的名字就是 **linux_386**；<br/>而在针对**Windows 64bit** 操作系统的安装包中，平台相关目录的名字则为 **windows amd64** <br/> **Go** 源码文件对应于以 **“.a”** 为结尾的归档文件，它们就存储在 **pkg** 文件夹下的平台相关目录中。<br/> **pkg** 文件夹下有一个名叫 **tool** 的子文件夹，该子文件夹下也有一个平台相关目录，<br/>其中存放了很多可执行文件【可参见 **1.6 标准命令概述**】。 |
|**src**  | 存放所有标准库、**Go** 语言工具，以及相关底层库( **C** 语言实现）的源码。<br/>通过查看这个文件夹，可以了解到 **Go** 语言的方方面面。本书的后续章节会适时地对其中的部分文件进行说明。 |
| **test** | 存放测试 **Go** 语言自身代码的文件。通过阅读这些测试文件，可大致了解 **Go** 语言的一些特性和使用方法。 |



# 3. 工作区
**Go** 代码必须放在工作区中，工作区其实就是一个对应于特定工程的目录，它包含 **3** 个子目录 **src** 目录，**pkg** 目录和 **bin** 目录。

- **src目录**
用于以代码包的形式组织并保存 **Go** 源码文件。这里的代码包，与 **src** 下的子目录一一对应。例如，若一个源码文件被声明为属于代码包logging，那么它就应当被保存在 **src**目录下名为 **logging** 的子目录中。当然，我们也可以把 **Go** 源码文件直接放于 **src** 目录下，但这样的 **Go** 源码文件就只能被声明为属于 **main** 代码包了。除非用于临时测试或演示，一般还是建议把**Go** 源码文件放入特定的代码包中。

- **pkg目录**
用于存放经由 `go install` 命令构建安装后的代码包（包含 **Go** 库源码文件）的 ***.a** 归档文件。该目录与 **GOROOT** 目录下的 **pkg** 功能类似，区别在于，工作区中的 **pkg** 目录专门用来存放用户代码的归档文件。构建和安装用户源码的过程一般会以代码包为单位进行，比如 **logging** 包被编译安装后，将生成一个名为 **logging.a** 的归档文件，并存放在当前工作区的 **pkg** 目录下的平台相关目录中。

- **bin目录**
与 **pkg** 目录类似，在通过 `go install` 命令完成安装后，保存由 **Go** 命令源码文件生成的可执行文件。在 **Linux** 操作系统下，这个可执行文件一般是一个与源码文件同名的文件。在 **Windows** 操作系统下，这个可执行文件的名称是源码文件名称加 **.exe** 后缀。

# 4. GOPATH
工作区的目录路径需要添加到环境变量 **GOPATH** 中。否则，即使处于同一个工作区（事实上，未被加入到环境变量 **GOPATH** 中的目录不应该称作工作区），代码之间也无法通过绝对代码包路径完成调用。在实际开发中，工作区往往有多个，这些工作区的目录路径都需要添加至 **GOPATH** 中。
如 **Linux** 下有两个工作区：

```bash
~/Go/lib
~/Go/goc2p
```

修改/etc/profile文件，添加环境变量GOPATH的内容：

```bash
export GOPATH=\$HOME/Go/lib:\$HOME/Go/goc2p
```

保存 `/etc/profile` 文件，并用 **source** 命令使配置生效。

**注意：**
- **GOPATH** 中不要包含环境变量 **GOROOT** 的值（即 **Go** 的安装目录路径），将 **Go** 语言本身的工作区和用户的工作区严格地分开；
- 通过 **Go** 工具中的代码获取命令 **go get**，可以将指定项目的源码下载到我们在环境变量 **GOPATH** 中设定的第一个工作区中，并在其中完成构建和安装的过程。

**Windows** 下直接在系统变量中添加 **GOPATH** 环境变量即可，其中值为你的工作区的根目录。

# 5. 源码文件
**Go** 语言的源码文件分为 **3** 类：
- **Go** 库源码文件
- **Go** 命令源码文件
- **Go** 测试源码文件

## 5.1 命令源码文件

声明为属于 **main** 代码包，并且包含 **无参数声明** 和 **结果声明** 的 **main** 函数的源码文件。这类文件可以独立运行（使用 `go run` 命令），也可以被 `go build` 或 `go install` 命令转换为可执行文件。

同一个代码包中的所有源码文件，其所属代码包的名称必须一致。如果命令源码文件和库源码文件处于同一代码包中，该包就无法正确执行 `go build` 和 `go install` 命令。换句话说，这些源码文件也就无法被编译和安装。因此，命令源码文件通常会单独放在一个代码包中。一般情况下，一个程序模块或软件的启动入口只有一个。

同一个代码包中可以有多个命令源码文件，可通过 `go run` 命令分别运行它们。但通过 `go build` 和 `go install` 命令无法编译和安装该代码包。所以一般情况下，不建议把多个命令源码文件放在同一个代码包中。

当代码包中有且仅有一个命令源码文件时，在文件所在目录中执行 `go build` 命令，即可在该目录下生成一个与目录同名的可执行文件；若使用 `go install` 命令，则可在当前工作区的 **bin** 目录下生成相应的可执行文件。

## 5.2 库源码文件

存在于某个代码包中的普通源码文件。通常，库源码文件声明的包名会与它实际所属的代码包（目录）名一致，且库源码文件中不包含 **无参数声明** 和 **无结果声明** 的 **main** 函数。如在 `basic/set` 目录下执行 `go install` 命令，成功地安装了 `basic/set` 包，并生成一个名为 **set.a** 的归档文件。归档文件的存放目录由以下规则产生：

1. 安装库源码文件时所生成的归档文件会被存放到当前工作区的 **pkg** 目录中。
2. 根据被编译的目标计算机架构，归档文件会被放置在 **pkg** 目录下的平台相关目录中。如上的 **set.a** 在我的 **64** 位 **window** 系统上就是**pkg\windows_amd64** 目录中。
3. 存放归档文件的目录的相对路径与被安装的代码包的上一级代码包的相对路径是一致的。第一个相对路径就是相对于工作区的 **pkg** 目录下的平台相关目录而言的，而第二个相对路径是相对于工作区的 **src** 目录而言的。如果被安装的代码包没有上一级代码包（也就是说它的父目录就是工作的 **src** 目录），那么它的归档文件就会被直接存放到当前工作区的 **pkg** 目录的平台相关目录下。如 **basic** 包的归档文件 **basic.a** 总会被直接存放到 **pkg\windows_amd64** 目录下，而 `basic/set` 包的归档文件 **set.a** 则会被存放到 **pkg\ windows_amd64\basic** 目录下。

## 5.3 测试源码文件

 这是一种特殊的库文件，可以通过执行 `go test` 命令运行当前代码包下的所有测试源码文件。成为测试源码文件的充分条件有两个：
1. 文件名需要以 **_test.go** 结尾
 2. 文件中需要至少包含该一个名称为 **Test** 开头或 **Benchmark** 开头，拥有一个类型为 **testing.T** 或 **testing.B** 的参数的函数。类型 **testing.T** 或 **testing.B** 分别对应功能测试和基础测试所需的结构体。

当在某个代码包中执行 `go test` 命令，该代码包中的所有测试源码文件就会被找到并运行。

> **注意**：存储 **Go** 代码的文本文件需要以 **UTF-8** 编码存储。如果源码文件中出现了非 **UTF-8** 编码的字符，则在运行、编译或安装时，**Go** 会抛出 **illegal UTF-8 sequence** 的错误。

# 6. 代码包
**Go** 语言中的代码包是对代码进行构建和打包的基本单元。
	
## 6.1 包声明

在 **Go** 语言中，代码包中的源码文件名可以是任意的；这些任意名称的源码文件都必须以包声明语句作为文件中代码的第一行。比如 **src** 目录下的代码包 `basic/set` 包中的所有源码文件都要先声明自己属于 `basic/set` 包：
```go
	package set
```
**package** 是 **Go** 语言中用于包声明语句的关键字。**Go** 语言规定包声明中的包名为代码包路径的最后一个元素。如上，`basic/set` 包的包路径为`basic/set`，而包声明中的包名则为 **set**。除了命令源码文件不论存放在哪个包中，都必须声明为属于 **main** 包。

## 6.2 包导入

代码包的导入使用代码包导入路径。代码包导入路径就是代码包在工作区的 **src** 目录下的相对路径，比如 **basic** 的绝对路径为 `E:\Go\goc2p\src\basic\set`，而 `E:\Go\goc2p` 是被包含在环境变量 **GOPATH** 中的工作区目录路径，则其代码包导入路径就是 `basic/set`。
```go
	import basic/set
```
	
当导入多个代码包时，需要用圆括号括起它们，且每个代码包名独占一行。在调用被导入代码包中的函数或使用其中的结构体、变量或常量时，需要使用包路径的最后一个元素加 **.** 的方式指定代码所在的包。
	
如果我们有两个包 `logging` 和 `go_lib/logging`，并且有相同的方法`NewSimpleLogger()`，且有一个源码文件需要导入这两个包：
```go
	import (
		"logging"
		"go_lib/logging"
	)
```
则这句代码 `logging.NewSimpleLogger()` 就会引起冲突，**Go** 语言无法知道`logging.` 代表的是哪一个包。所以，在 **Go** 语言中，如果在同一个源码文件中导入多个代码包，那么代码包路径的最后一个元素不可以重复。

如果用这段代码包导入代码，在编译代码时，**Go** 语言会抛出 **”logging redeclared as imported package name”** 的错误。如果确实需要导入，当有这类重复时，我们可以给它们起个别名来区别：
```go
	import (
		la "logging"
		lb "go_lib/logging"
	)
```

调用包中的代码：
```go
	var logger la.Logger = la.NewSimpleLogger()
```

这里不必给每个引起冲突的代码包都起一个别名，只要能够区分它们就可以了。

如果我们想直接调用某个依赖包的程序，就可以用 **.** 来代替别名。
```go
	import (
		. "logging"
		lb "go_lib/logging"
	)
```
	
在当前源码文件中，可以直接进行代码调用了：
```go
	var logger Logger = NewSimpleLogger()
```
	
**Go** 语言把变量、常量、函数、结构体和接口统称为程序实体，而把它们的名字统称为标识符。标识符可以是任何 **Unicode** 编码可以表示的字母字符、数字以及下划线 **_**，并且，首字母不能是数字或下划线。

标识符的首字母的大小写控制着对应程序实体的访问权限。如果标识符的首字母是大写的，那么它对应的程序实体就可以被本代码包之外的代码访问到，也可以称其为可导出的。否则对应的程序实体就只能被本包内的代码访问。当然，还需要有以下两个额外条件：
1. 程序实体必须是非局部的。局部程序实体是被定义在函数或结构体的内部。
2. 代码包所在的目录必须被包含在环境变量 **GOPATH** 中的工作区目录中。

如果代码包 **logging** 中有一个叫做 **getSimpleLogger** 的函数，那么光从这个函数的名字上我们就可以看出，这个函数是不能被包外代码调用的。

如果我们只想初始化某个代码包而不需要在当前源码文件中使用那个代码包中的任何代码，既可以用 **_** 来代替别名
```go
	import (
		_ "logging"
	)
```
	
	
## 6.3 包初始化

在 **Go** 语言中，可以有专门的函数负责代码包初始化。这个函数需要无参数声明和结果声明，且名称必须为 **init**，如下：
```go
	func init() {
		println("Initialize")
	}
```
	
Go语言会在程序真正执行前对整个程序的依赖进行分析，并初始化相关的代码包。也就是说，所有的代码包初始化函数都会在 **main** 函数（命令源码文件中的入口函数）之前执行完成，而且只会执行一次。并且，当前代码包中的所有全局变量的初始化都会在代码包初始化函数执行前完成。这就避免了在代码包初始化函数对某个变量进行赋值之后又被该变量声明中赋予的值覆盖掉的问题。

这里举出 **《Go并发编程实战》** 中的例子，帮助理解上面的包初始化，如下：
```go
	package main // 命令源码文件必须在这里声明自己属于main包

	import ( // 引入了代码包fmt和runtime
		"fmt"
		"runtime"
	)
	
	func init() { // 包初始化函数
		fmt.Printf("Map: %v\n", m) // 先格式化再打印
		// 通过调用runtime包的代码获取当前机器所运行的操作系统以及计算架构
		// 而后通过fmt包的Sprintf方法进行字符串格式化并赋值给变量info
		info = fmt.Sprintf("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
	}

	// 非局部变量，map类型，且已初始化
	var m map[int]string = map[int]string{1: "A", 2: "B", 3: "C"}
	var info string // 非局部变量，string类型，未被初始化

	func main() { // 命令源码文件必须有的入口函数
		fmt.Println(info) // 打印变量info
	}
```
	
命名源码文件名为 **initpkg_demo.go**，并保存到工作区的 `basic/pkginit` 包中。如下图为我本机运行的结果：

 ![](go-result.png)
	
在同一个代码包中，可以存在多个代码包初始化函数，甚至代码包内的每一个源码文件都可以定义多个代码包初始化函数。**Go** 语言编译器不能保证同一个代码包中的多个代码包初始化函数的执行顺序。如果要求按特定顺序执行的话，可以考虑使用 **Channel**（**Go**语言并发编程模型的一员）进行控制。

**Go** 语言认可两个特殊的代码包名称----**all** 和 **std**。**all** 代表了环境变量 **GOPATH** 中包含的所有工作区中的所有代码包，而 **std** 则代表了 **Go** 语言标准库中的所有代码包。

# 7. 标准命令概述 
- **build**    	编译给定的代码包或Go语言源码文件及其依赖包
- **clean**    	清除执行其他go命令后遗留的目录和文件
- **doc** 		执行godoc命令以打印指定代码包。
- **env** 		打印Go语言环境信息
- **fix**	    执行go tool fix命令以修正给定代码包的源码文件中包含的过时语法和代码调用
- **fmt**       执行gofmt命令以格式化戈丁代码包中的源码文件。
- **generate**  generate Go files by processing source
- **get**       下载和安装给定的代码包及其依赖包
- **install**   编译和安装给定的代码包及其依赖包
- **list**      显示给定代码包的信息
- **run**       编译并运行给定的命令源码文件
- **test**      测试给定的代码包
- **tool**      运行Go语言的特殊工具
- **version**   显示当前安装的Go语言的版本信息
- **vet**       run go tool vet on packages
	
在执行上述命令的时候可以通过附加一些额外的标记来定制命令的执行过程，这些标记可以看做是命令的特殊参数，这些特殊参数可以添加到命令名称和命令的真正参数中间，如下： 
- **-a** 强行重新构建所有涉及的Go语言代码包（包括Go语言标准库中的代码包），即使它们已经是最新的了。
- **-n** 使命令仅打印在执行期间使用到的所有命令，而不真正执行它们。
- **-v** 打印出命令执行过程中涉及的Go语言代码包的名字。这些代码包一般包括我们自己给定的目标代码包，有时候还会包括该代码包直接或间接依赖的代码包。
- **-work** 打印出命令执行时生成和使用的临时工作目录的名字，且命令执行完成后不对它进行删除。
- **-x** 打印出命令执行期间使用到的所有命令。	

常用的Go语言的特殊工具，如下：
- **fix**  可以把给定代码包的所有Go语言源码文件中的旧版本代码修改为新版本。它是我们升级Go语言版本后会使用到的工具。
- **vet**  用于检查Go语言源码中静态错误的简单工具。可以使用它检测一些常见的Go语言代码编写错误。
- **pprof**  用于以交互的方式访问一些性能概要文件。命令将会分析给定的概要文件，并根据要求提供高可读性的输出信息。这个工具可以分析的概要文件包括CPU概要文件、内存概要文件和程序阻塞概要文件。这些内含Go语言运行时信息的概要文件可以通过标准库代码包runtime和runtime/pprof中的程序来生成。
- **cgo**  用于帮助Go语言代码使用C语言代码库，以及使用Go语言代码可以被C语言代码引用。

# 结语
最后附上 **《Go并发编程实战》** 作者郝林托管到 **GitHub** 的 **Go** 命令教程，里面涉及了 **Go** 命令和工具的详细用法：[https://github.com/hyper-carrot/go_command_tutorial
](https://github.com/hyper-carrot/go_command_tutorial)
Go语言学习的第一天，以后持续更新。。。



