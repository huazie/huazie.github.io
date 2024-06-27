---
title: Sublime Text 3配置C/C++开发环境
date: 2023-12-15 09:39:11
updated: 2024-01-29 15:53:57
categories:
  - [开发工具]
tags:
  - C
  - C++
  - Sublime Text 3
  - 开发环境搭建
---

![](/images/cplus-logo.png)


# 一、引言
**C/C++** 的开发工具有很多，相信大家多多少少都有所了解，比如 **Visual C++**、**Codeblocks**、**VSCode**、**Dev C++** 等等；本篇 **Huazie** 介绍一个比较轻量级的开发环境 **Sublime Text 3**，并用它来配置 **C/C++** 开发环境。

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
## 2.2 接入 mingw-w64 

**Mingw-w64** 是一个用于在 **Windows** 系统上支持 **GCC** 编译器的完整运行时环境。这个项目是原始 **Mingw.org** 项目的进一步发展，旨在支持 **64位 和 32位 Windows** 操作系统原生的二进制文件。在 **2007** 年，**Mingw-w64** 从原始 **Mingw.org**项目分叉出来，提供对 **64位**和 **新API** 的支持。从那时起，它逐渐得到了广泛的应用和传播。

**Mingw-w64** 提供了一百万行以上的头文件、库和运行时，用于在 **Windows** 上链接和运行代码。它包括 **Winpthreads** 库（用于 **C++11** 线程支持）和 **Winstorecompat**库（一个便利库，可简化与 **Windows** 应用商店的一致性）。与**VisualStudio** 相比，**Mingw-w64** 在数学支持方面更为完善，并且执行速度更快。
### 2.2.1 下载 mingw-w64 
参考这个帖子内容，下载一下 **mingw-w64**：

[https://tieba.baidu.com/p/5487544851](https://tieba.baidu.com/p/5487544851)


我目前用的 **MinGW-w64 9.0.0**，有需要的可以在评论回复 或者 私信我。
### 2.2.2 环境变量配置

不管是安装版的，还是免安装版的，都会有类似如下的目录：

![](mingw.png)

上图中的 **bin** 目录，我们点进去看下：

![](mingw-bin.png)

实际上这里的 **gcc.exe** 和 **g++.exe** 就是 **C/C++** 程序编译和运行的关键所在，当然这里还有其他的可执行程序，感兴趣的朋友可以自行搜索了解，这里不展开了。

现在，**Huazie** 以 **window 11** 系统为例，介绍下配置环境变量，如下：

右击 **Window** 图标，打开下图并选择 **系统**：

![](/images/dev-tool/windows-system.png)


点击 **高级系统设置**，打开系统属性页面，点击 **环境变量** ：

![](/images/dev-tool/windows-env-config.png)

找到 **Path** 环境变量，配置上面你的 **mingw64** 的 **bin** 目录进去：

![](mingw-bin-env-config.png)

## 2.3 配置 C/C++ 开发环境

打开 **Sublime Text 3**，我们可以看到如下的界面：

![](/images/dev-tool/sublime-text3-default-page.png)

菜单栏选择 **Tools** => **Build System** => **New Build System**

![](/images/dev-tool/sublime-text3-new-buildsystem.png)

### 2.3.1 C Build System 配置

```c
{
    "windows": {
        "cmd": ["gcc", "-std=c11", "${file}", "-o", "${file_base_name}.exe"],
    },
    "cmd": ["gcc", "-std=c11", "${file}", "-o", "${file_base_name}"],
    "file_regex": "^(.*)\\(([0-9]+),([0-9]+)\\) (Error|Fatal): (.*)$",
    "working_dir": "${file_path}",
    "selector": "source.c",
    "encoding": "cp936",
    "variants": [
    {
        "name": "Run",
        "shell": true,
        "windows": {
            "shell_cmd" : "start cmd /c \"\"${file_base_name}.exe\" & echo. & pause\""
        }
    },
    {
        "name": "Build and Run",
        "shell": true,
        "windows": {
            "shell_cmd": "gcc -std=c11 \"${file}\" -o \"${file_base_name}.exe\" && start cmd /c \"\"${file_base_name}.exe\" & echo. & pause\""
        },
    }
    ]
}
```

上述内容保存在，前面打开的 **New Build System** 中，并命名为 **C.sublime-build**：

![](c-sublime-build.png)
### 2.3.2 C++ Build System 配置

```c
{  
    "cmd": ["g++", "${file}", "-o", "${file_path}/${file_base_name}", "-Wall" ,"&&","start","cb_console_runner.exe","${file_path}/${file_base_name}"],  
    "file_regex": "^(..[^:]*):([0-9]+):?([0-9]+)?:? (.*)$",  
    "working_dir": "${file_path}",  
    "selector": "source.c, source.c++",  
    "shell": true,  
    "encoding": "cp936",  
    "variants":  
    [  
        {  
            "name": "Run",  
            "cmd": ["start","cb_console_runner.exe","${file_path}/${file_base_name}"]  
        }  
    ]  
}  
```

重新 **New Build System**，将上述内容保存其中，并命名为 **CPP.sublime-build**

![](cplus-sublime-build.png)
## 2.4 编写 C/C++ 代码

现在让我们开始编写吧！！！

### 2.4.1 第一个 C 代码【helloworld.c】

```c
#include<stdio.h>

int main() 
{
    printf("hello world!\n");
    printf("[C]Author: Huazie");
    return 0;
}

```
### 2.4.2 第一个 C++ 代码【helloworld.cpp】

```cpp
#include <iostream>
using namespace std;

int main()
{
    cout << "Hello World" << endl;
    cout << "[C++]Author: Huazie";
    return 0;
}
```

## 2.5 运行 C/C++ 代码

### 2.5.1 运行 C 代码

菜单栏 **Tools** => **Build System** ，然后 选择 **C**，就是前面的 **C.sublime-build**。

然后直接 **Ctrl + B**，编译运行当前的程序，运行截图如下所示：

![](c-result.png)
### 2.5.2 运行 C++ 代码

菜单栏 **Tools** => **Build System** ，然后 选择 **CPP**，就是前面的 **CPP.sublime-build**。

然后直接 **Ctrl + B**，编译运行当前的程序，运行截图如下所示：

![](cplus-result.png)

# 三、总结

本篇 **Huazie** 介绍了 **Sublime Text 3** 配置 **C/C++** 的相关内容，感兴趣的朋友赶紧配置起来，有任何问题可以随时评论区沟通。

# 四、更新

经过上面的配置，大家都能写自己的第一个 **C/C++** 代码了，但是有些小伙伴慢慢使用发现，如果输出的内容包含中文，打印出来的信息是乱码的。有关这个问题，请查看笔者的另一篇博文 [《Sublime Text 3 解决中文乱码问题》](../../../../../2023/12/15/dev-tool/sublime-text3-convert-to-utf8/)
