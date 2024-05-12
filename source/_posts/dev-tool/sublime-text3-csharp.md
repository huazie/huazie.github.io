---
title: Sublime Text 3配置 C# 开发环境
date: 2023-12-22 11:50:03
updated: 2024-01-29 15:53:57
categories:
  - [开发工具]
tags:
  - C#
  - Sublime Text 3
  - 开发环境搭建
---


[《开发工具系列》](/categories/开发工具/)

![](/images/csharp-logo.png)

# 一、引言
**C#** 是一种面向对象的编程语言，由微软公司开发。它的常用开发工具，相信大家多多少少都有所了解，比如 **Visual Studio**、**Visual Studio Code**；本篇 **Huazie** 介绍一个比较轻量级的开发环境 **Sublime Text 3**，并用它来配置 **C#** 开发环境。

<!-- more -->

# 二、主要内容
## 2.1 初识 Sublime Text 3 

**Sublime Text 3** 是一款流行的文本编辑器，它的特点是体积小巧、启动速度快、界面简洁美观。它具有强大的代码编辑功能，支持多种编程语言。此外，**Sublime Text 3** 还具有丰富的插件生态系统，用户可以根据自己的需求安装各种插件来扩展其功能。

**Sublime Text 3** 的一些主要特点，如下所示：

- **强大的代码编辑功能**：**Sublime Text 3** 提供了许多实用的代码编辑功能，如自动完成、代码高亮、代码片段等，大大提高了编程效率。

- **支持多种编程语言**：**Sublime Text 3** 支持多种编程语言，包括 **HTML、CSS、JavaScript、Python、Ruby、PHP** 等，用户可以根据需要选择不同的语言模式。

- **插件生态系统**：**Sublime Text 3** 拥有丰富的插件生态系统，用户可以通过安装插件来扩展其功能，如 **Emmet**（用于编写 **HTML** 和 **CSS**）、**Package Control**（用于安装和管理插件）等。

- **自定义快捷键**：**Sublime Text 3** 允许用户自定义快捷键，以便更快速地执行常用操作。

- **多窗口编辑**：**Sublime Text 3** 支持多窗口编辑，用户可以同时打开多个文件进行编辑，方便进行代码对比和复制粘贴操作。

- **跨平台支持**：**Sublime Text 3** 支持 **Windows**、**Mac** 和 **Linux** 操作系统，用户可以在不同的平台上使用相同的设置和插件。

- **版本控制集成**：**Sublime Text 3** 可以与版本控制系统（如 **Git**）集成，方便用户进行代码版本管理。

## 2.2 初识 C#

**C#** （发音为 **“C-Sharp”** ）是微软开发的一种面向对象的编程语言，它是 **.NET** 框架的重要组成部分。

**C#** 的主要特点包括：

- **类型安全：** **C#** 是一种强类型的语言，这意味着它会在编译时检查类型错误，而不是在运行时。这有助于提高代码的稳定性和可维护性。
- **面向对象：** **C#** 支持面向对象编程，包括类、接口、继承、多态等概念。这使得代码更加模块化、可重用和易于维护。
- **简洁的语法：** **C#** 的语法相对简洁，易于学习和使用。它支持许多现代编程语言的特性，如 **LINQ（Language Integrated Query）**、异步编程、**Lambda** 表达式等。
- **强大的库支持：** **C#** 有强大的标准库和第三方库支持，可以方便地访问数据库、文件系统、网络等资源。
- **与.NET框架集成：** **C#** 是 **.NET** 框架的一部分，可以方便地使用 **.NET** 框架提供的类库和功能。
## 2.3 接入 .NET Framework

**.NET Framework** 是微软推出的一种开发框架，用于构建多种类型的应用程序，包括传统的**Windows** 应用程序、基于 **Web** 的应用程序、移动应用程序和云服务。它提供了一个公共的面向对象的编程环境，支持多种编程语言，如 **C#、VB.NET、F#** 等。

**.NET Framework** 具有两个主要组件：**公共语言运行库** 和 **.NET Framework类库**。
- 公共语言运行库是 **.NET Framework** 的基础，类似于 **Java** 的虚拟机，它负责代码的编译、执行和内存管理。
- **.NET Framework** 类库是一个综合性的面向对象的可重用类型集合，提供了丰富的类和方法，用于处理各种任务，如数学计算、字符操作、数据库操作等。

**.NET Framework** 的目标是实现代码的可移植性、安全性和可执行性。它提供了一个一致的面向对象的编程环境，无论对象代码是在本地存储和执行，还是在本地执行但在 **Internet** 上分布，或者是在远程执行的。此外，它还提供了一个将软件部署和版本控制冲突最小化的代码执行环境，以及一个可提高代码（包括由未知的或不完全受信任的第三方创建的代码）执行安全性的代码执行环境。


### 2.3.1 下载  .NET Framework

可以直达 [官网下载](https://dotnet.microsoft.com/zh-cn/download/dotnet-framework) .NET Framework

![](dot-net-download.png)


实际上 **Windows** 系统基本上都集成了 **.NET Framework** ，我们可以从[《安装面向开发人员的 .NET Framework》](https://learn.microsoft.com/zh-cn/dotnet/framework/install/guide-for-developers)查看不同 Windows 系统适配的版本。


![](dot-net-rel-system.png)

有关 **.NET Framework** 的更多内容，请查看 [官方文档](https://learn.microsoft.com/zh-cn/dotnet/framework/)
### 2.3.2 环境变量配置

我们可以到 `C:\Windows\Microsoft.NET\Framework64` 查看 .NET 的不同版本：

![](dot-net-version.png)
![](csc-exe.png)

上图中的 **csc.exe** 其实就是本次配置的关键。

> **知识点：** `csc.exe` 是 **C#** 的命令行编译器，全称为 **CSharpCompiler**。它是微软 **.NET Framework** 中的一个重要组件，用于将 **C#** 源代码【后缀为 **cs** 的文件】编译成可执行程序或库文件。


现在，**Huazie** 以 **window 11** 系统为例，介绍下配置环境变量，如下：

右击 **Window** 图标，打开下图并选择 **系统**：

![](/images/dev-tool/windows-system.png)


点击 **高级系统设置**，打开系统属性页面，点击 **环境变量** ：

![](/images/dev-tool/windows-env-config.png)

找到 **Path** 环境变量，配置上面你的 .NET 框架指定目录配置进去：

![](dot-net-env-config.png)



## 2.4 配置 C# 开发环境

初次打开 **Sublime Text 3**，我们可以看到如下的界面：

![](/images/dev-tool/sublime-text3-default-page.png)

菜单栏选择 **Tools** => **Build System** => **New Build System**

![](/images/dev-tool/sublime-text3-new-buildsystem.png)

**C# Build System 配置**

注意看下面的 **shell_cmd** 是关键。

```bash
{
	"shell_cmd": "csc /out:\"${file_path}/${file_base_name}.exe\" \"${file}\"",
	"file_regex": "^(..[^:]*):([0-9]+):?([0-9]+)?:? (.*)$",
	"working_dir": "${file_path}",
	"selector": "source.cs",
	"variants":
		[
			{
				"name": "Build & Run",
				"shell_cmd": "csc /out:\"${file_path}/${file_base_name}.exe\" \"${file}\" && start \"${file_base_name}.exe\" /d \"${file_path}\" \"${file_base_name}.exe\"",
				"working_dir": "${file_path}"
			},
			{
				"name": "Run",
				"shell_cmd": "start \"${file_base_name}.exe\" /d \"${file_path}\" \"${file_base_name}.exe\"",
				"working_dir": "${file_path}"
			},
			{
				"name": "Build (Form)",
				"shell_cmd": "csc /t:winexe /r:System.Windows.Forms.dll;System.Drawing.dll /out:\"${file_path}/${file_base_name}.exe\" \"${file}\"",
				"working_dir": "${file_path}"
			},
			{
				"name": "Build & Run (Form)",
				"shell_cmd": "csc /t:winexe /r:System.Windows.Forms.dll;System.Drawing.dll /out:\"${file_path}/${file_base_name}.exe\" \"${file}\" && start \"${file_base_name}.exe\" /d \"${file_path}\" \"${file_base_name}.exe\"",
				"working_dir": "${file_path}"
			},
			{
				"name": "Run (Form)",
				"shell_cmd": "start \"${file_base_name}.exe\" /d \"${file_path}\" \"${file_base_name}.exe\"",
				"working_dir": "${file_path}"
			},
		]
}
```

上述内容保存在，前面打开的 **New Build System** 中，并命名为 **C#.sublime-build**。



## 2.5 编写 C# 代码

现在让我们开始编写第一个 **C#** 代码吧！

```csharp
using System;
namespace HelloWorldApp
{
    class HelloWorld
    {

        static void Main(string[] args)
        {
            Console.WriteLine("hello world!");
            Console.WriteLine("[C#]Author:{0}", "Huazie");
            Console.ReadKey();
        }
    }
}
```
上述 **C#** 代码，我们会新建一个 **helloworld.cs** 文件进行保存。

> **注意：** **C#** 源码文件名的后缀为 `cs`
## 2.6 运行 C# 代码


菜单栏 **Tools** => **Build System** ，然后 选择 **C#** ，就是前面的 **C#.sublime-build**。

然后直接 **Ctrl + B**，编译运行当前的程序，运行截图如下所示：

![](csharp-result.png)


# 三、总结

本篇 **Huazie** 介绍了 **Sublime Text 3** 配置 **C#** 的相关内容，感兴趣的朋友赶紧配置起来，有任何问题可以随时评论区沟通。



