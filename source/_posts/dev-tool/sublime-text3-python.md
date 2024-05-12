---
title: Sublime Text 3配置 Python 开发环境
date: 2024-01-12 11:59:47
updated: 2024-01-29 15:48:08
categories:
  - [开发工具]
tags:
  - Python
  - Sublime Text 3
  - 开发环境搭建
---

[《开发工具系列》](/categories/开发工具/)

![](/images/python-logo.png)

# 一、引言
**Python** 是一种简洁但功能强大的面向对象编程语言。它的常用开发工具有很多，相信大家多多少少都有所了解，比如 **PyCharm**、**Visual Studio Code**、**IDLE** 等等；本篇 **Huazie** 介绍一个比较轻量级的开发环境 **Sublime Text 3**，并用它来配置 **Python** 开发环境。

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

## 2.2 初识 Python

**Python** 是一种高级编程语言，它具有如下的特点：

- 优雅的语法，编写和阅读都很简单。
- 简单好用，轻松写程序。这个特点让 **Python** 做以下工作很方便：开发一个原型或其他特定的小任务，而不用太费劲维护。
- 内置庞大的标准库，包含常见的编程任务，比如连接网页服务器，用正则表达式搜索文本，读写文件。
- **Python** 交互模式可以轻松测试代码片段，也可以使用 **IDLE** 的集成开发环境。
- 也可以轻松扩展用 **C或C++** 编译出来的新模块。
- 可以嵌入软件系统来提供编程接口。
- 跨平台运行，包括 **Mac OS X，Windows，Linux和Unix**，在 **安卓和iOS** 上也有非官方实现。

**Python** 是自由软件，包括两种含义：

 1. 下载或使用 **Python** 是免费的
 2. 虽然 **Python** 编程语言有版权，但可以自由修改和分发。

**Python** 语言的编程语言特性：

- 多样的数据类型：数（浮点数、复数和无限长整数），字符串（**ASCII** 和 **Unicode**）及字典
- **Python** 通过类和多重继承来支持面向对象编程
- 代码可以用模块和包来组织
- 支持抛出和捕捉异常，用于干净的错误处理
- 数据是强类型、动态类型，不兼容数据操作会抛出异常（比如试图把字符串和数字加起来），这样能马上发现错误。
- 包含生成器（**generators**）和列表操作（**list comprehesions**）这样的高级特性
- 自动管理内存，避免你在自己的代码里费心申请释放内存
## 2.3 接入 Python

### 2.3.1 下载 

[Python 官网下载](https://www.python.org/downloads/)

![](python-download.png)


> **Python 3.12** 支持 **Windows 8.1 及其后的版本**。 如果你需要 **Windows 7** 支持，请安装 **Python 3.8**。

### 2.3.2 安装和使用 python

[Using Python on Unix platforms](https://docs.python.org/3/using/unix.html)

[Using Python on Windows](https://docs.python.org/3/using/windows.html)  【**博主下面演示的环境**】

[Using Python on a Mac](https://docs.python.org/3/using/mac.html)


笔者这里下载的是 **Python 3.12.1**，安装界面如下图： 

![](python-install-page.png)

从上图，我们可以看到两个选项：

- **Install Now（立即安装）：** 
	- 你不需要成为管理员（除非需要更新 **C** 运行时库的系统，或者为所有用户安装 **Windows** 的 **Python** 启动器）

	- **Python** 将安装在你的用户目录下

	- **Windows** 的 **Python** 启动器将根据第一页底部的选项进行安装

	- 标准库、测试套件、启动器和 **pip** 将被安装

	- 如果选择，安装目录将被添加到你的 **PATH** 中

	- 快捷方式只对当前用户可见

- **Customize installation（自定义安装）：** 
	将允许大家选择：要安装的功能、安装位置、其他选项或安装后的操作。如果要安装调试符号或二进制文件，我们需要使用此选项。
	如要为全部用户安装，应选择 **“自定义安装”**：

	- 您可能需要提供管理员凭据或批准

	- **Python** 将安装到 **Program Files** 目录

	- **Windows** 的 **Python** 启动器将安装到 **Windows** 目录

	- 安装过程中可以选择可选功能

	- 标准库可以预先编译为字节码

	- 如果选择，安装目录将被添加到系统 **PATH** 中

	- 快捷方式对所有用户可用

### 2.3.3 环境变量配置


现在，**Huazie** 以 **windows 11** 系统为例，介绍下配置环境变量，如下：

右击 **Window** 图标，打开下图并选择 **系统**：

![](/images/dev-tool/windows-system.png)


点击 **高级系统设置**，打开系统属性页面，点击 **环境变量** ：

![](/images/dev-tool/windows-env-config.png)

找到 **Path** 环境变量，配置上面你的 **Python** 安装目录进去：

![](python-env-config.png)
> **注意：** **Python** 安装目录中的 **Scripts** 目录中包含了 **Python** 的包管理工具，该环境变量按需配置即可。


然后 **Win + R**，打开如下窗口，输入 **cmd**，点击确认打开命令行窗口。

![](/images/dev-tool/win-r-cmd.png)

在命令行窗口内，输入 `python -V` 查看，如下图所示即为安装成功：

![](python-v.png)


## 2.4 配置 Python 开发环境

初次打开 **Sublime Text 3**，我们可以看到如下的界面：

![](/images/dev-tool/sublime-text3-default-page.png)

菜单栏选择 **Tools** => **Build System**，可能会看到如下的 **Python**，我们直接选择这个来编译下面的代码。

![](python-build-system.png)

因为我们是在 **Windows** 下安装的 **Python3** ，所以也可以手动添加 **Python3** 相关的 **Build System** 配置。

在菜单栏选择 **Tools** => **Build System** => **New Build System**，打开如下页面

![](/images/dev-tool/sublime-text3-new-buildsystem.png)

```bash
{
	"cmd": ["python", "-u", "$file"], 
	"file_regex": "^(..[^:]*):([0-9]+):?([0-9]+)?:? (.*)$",
	"selector": "source.python"
}
```

将上述内容保存在，前面打开的 **New Build System** 中，并命名为 **Python3.sublime-build**。


当然，也可以按住 **Ctrl+Shift+p** ，弹出如下输入窗口，在其中输入 **install package**，并选中红框内的列表。

![](install-package.png)

稍等一会儿，在打开的窗口里，输入 **Python** ，点击如下红框内的 **Python 3**：

![](install-package-python3.png)

通过上述这种方式，我们也可以安装 **Python3** 对应的 **Build System**，也就不需要手动添加了【安装成功后，可在 **Tools** => **Build System** 查看到】


## 2.5 编写 Python 代码 

现在让我们开始编写第一个 **python** 代码吧！

```python
# Python 3
author = "huazie";
print("Hello, World!");
print("Author: " + author);
```

针对上述 **python** 代码，我们会新建一个 **helloworld.py** 文件进行保存。

> **注意：** **python** 源码文件名的后缀为 `py`
## 2.6 运行 Python 代码


菜单栏 **Tools** => **Build System** ，然后 选择 **Python3**，就是前面的 **Python3.sublime-build**。

然后直接 **Ctrl + B**，编译运行当前的程序，运行截图如下所示：

![](python-result.png)


# 三、总结

本篇 **Huazie** 介绍了 **Sublime Text 3** 配置 **Python** 开发环境的相关内容，感兴趣的朋友赶紧配置起来，有任何问题可以随时评论区沟通。



