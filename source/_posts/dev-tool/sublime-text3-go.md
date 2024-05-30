---
title: Sublime Text 3配置Go语言开发环境
date: 2017-05-23 16:45:04
updated: 2024-01-14 20:30:04 
categories:
  - [开发工具]
  - [开发语言-Go]
tags:
  - Go
  - Sublime Text 3
  - GoSublime
  - 开发环境搭建
---

![](/images/go-logo.png)

# 1. 初识 Sublime Text 3 

**Sublime Text 3** 是一款流行的文本编辑器，它的特点是体积小巧、启动速度快、界面简洁美观。它具有强大的代码编辑功能，支持多种编程语言。此外，**Sublime Text 3** 还具有丰富的插件生态系统，用户可以根据自己的需求安装各种插件来扩展其功能。

<!-- more -->

**Sublime Text 3** 的一些主要特点，如下所示：

- **强大的代码编辑功能**：**Sublime Text 3** 提供了许多实用的代码编辑功能，如自动完成、代码高亮、代码片段等，大大提高了编程效率。

- **支持多种编程语言**：**Sublime Text 3** 支持多种编程语言，包括 **HTML、CSS、JavaScript、Python、Ruby、PHP** 等，用户可以根据需要选择不同的语言模式。

- **插件生态系统**：**Sublime Text 3** 拥有丰富的插件生态系统，用户可以通过安装插件来扩展其功能，如 **Emmet**（用于编写 **HTML** 和 **CSS**）、**Package Control**（用于安装和管理插件）等。

- **自定义快捷键**：**Sublime Text 3** 允许用户自定义快捷键，以便更快速地执行常用操作。

- **多窗口编辑**：**Sublime Text 3** 支持多窗口编辑，用户可以同时打开多个文件进行编辑，方便进行代码对比和复制粘贴操作。

- **跨平台支持**：**Sublime Text 3** 支持 **Windows**、**Mac** 和 **Linux** 操作系统，用户可以在不同的平台上使用相同的设置和插件。

- **版本控制集成**：**Sublime Text 3** 可以与版本控制系统（如 **Git**）集成，方便用户进行代码版本管理。

# 2. Go语言环境搭建
本篇博文是在读者Go自身环境已经搭好，Sublime Text 3也已装好的基础上所总结而来。

Go语言环境搭建可参考笔者的另一篇博文 [Go语言学习1-基础入门](/2016/06/27/go-learning1)。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 3. GoSublime安装和配置
首先我们需要下载Sublime下的一个重要的 `Go` 语言插件 **GoSublime**
 
 ## 3.1 安装步骤
 具体安装步骤如下：
  
（1）安装 `Package Control`，这个如果已经安装过的朋友可以直接跳过。关于 `Package Control` 的安装可以参考[《Sublime Text 3 中安装Package Control并配置》](https://zhuanlan.zhihu.com/p/349113898)

（2）打开 `Sublime Text 3` ，按住 `Ctrl+Shift+p` ，弹出如下输入窗口，在其中输入 `install package`，并选中红框内的列表。

![](install-package.png)

（3）然后在弹出的输入窗口中，输入 `GoSublime`，回车即可，此时 `GoSublime` 安装成功后，重启 `Sublime Text 3` 即可。

![](install-gosublime.png)
这一步很多朋友如果搜索不到，可以手动安装 `GoSublime` 插件，具体可参考 [Sublime Text 3 安装Go语言相关插件gosublime时 搜不到gosublime](https://www.cnblogs.com/chengxuyuan326260/p/10095914.html) ；建议`GoSublime` 采用 `development` 分支，看 `master` 分支好像已不维护了
![](gosublime-master-not-supported.png)

## 3.2 代码开发
接下来就可以开发代码了，终于可以编写 `Go` 语言程序了。在编写代码过程中可以体会到 `GoSublime` 拥有的代码提示功能，很大地方便了开发：

![](go-code-writing.png)


## 3.3 编译运行
按住 `Ctrl+B` 就可以编译你的命令源码文件，运行结果将会展示在下面：
![](go-code-running.png)

如果 `Ctrl+B` 没有效果，就需要到工具栏 `Tools->Build System->New Build System`，在新打开的文本中输入如下文本 ：
```
{ 
	"cmd": ["go", "run", "$file_name"], 
	"file_regex": "^[ ]*File \"(…*?)\", line ([0-9]*)", 
	"working_dir": "$file_path", 
	"selector": "source.go" 
}
```
保存，命名为 **go.sublime-build** 就可。然后在 `Tools->Build System` 中选中 **go** 即可，这个时候在进行（4）中的操作就能够得到结果了。
> **注意**：有一点要关注 我这边选择 **go**，是因为我上面 **sublime-build**文件的命名为 **go** ，如果是其他命名可以自行选择。



