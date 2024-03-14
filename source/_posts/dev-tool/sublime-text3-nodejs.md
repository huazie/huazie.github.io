---
title: Sublime Text 3配置 Node.js 开发环境
date: 2024-02-07 17:59:34
updated: 2024-02-07 17:59:34
categories:
  - [开发工具]
tags:
  - Node.js
  - JavaScript
  - Sublime Text 3
  - 开发环境搭建
---


[《开发工具系列》](/categories/开发工具/) 

![](/images/nodejs-logo.png)

# 一、引言

**Node.js** 是基于 **Chrome JavaScript** 运行时建立的一个平台，它简单理解就是运行在服务端的 **JavaScript**。它的开发环境有很多，比如 **VS Code**、**Atom** 等等，相信大家多多少少都有接触过；而本篇 **Huazie** 将要介绍一个比较轻量级的开发工具 **Sublime Text 3**，并用它来配置 **Node.js** 的开发环境。

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

## 2.2 初识 Node.js

**Node.js** 是一个开源的、跨平台的 **JavaScript** 运行环境，它使得开发者可以使用 **JavaScript** 来编写服务器端的应用程序。

它的主要特点如下：

- **单线程** ：**Node.js** 是单线程的，这意味着它一次只能做一件事。但是，由于 **JavaScript** 是事件驱动的，**Node.js** 可以处理大量并发请求，而不会因为某个请求耗时过长而阻塞其他请求。
- **非阻塞 I/O**：**Node.js** 使用事件驱动和非阻塞 **I/O** 模型，使得它能够处理大量并发请求。这使得 **Node.js** 在处理大量并发请求时比传统的多线程服务器更加高效。
- **跨平台**：**Node.js** 可以运行在 **Windows**、**Mac OS X**、**Linux** 等操作系统上。
- **强大的社区支持**：**Node.js** 拥有庞大的社区和丰富的第三方库，使得开发者可以快速地构建各种应用程序。
- **与浏览器无缝集成**：由于 **Node.js** 是基于 **Chrome** 的 **V8 JavaScript** 引擎构建的，因此它与浏览器中的 **JavaScript** 有着相同的 **API** 和语法。这意味着 **Node.js** 可以直接使用许多浏览器中的 **JavaScript** 库和框架，如 **Express**、**Mongoose**、**Socket.IO** 等。
## 2.3 接入 Node.js 

### 2.3.1 下载并安装 Node.js 

[Node.js 官方下载地址](https://nodejs.org/en/download/)

笔者本地下载的是 **20.11.0 LTS**，这对大多数用户来说已经足够了

![](/images/dev-tool/nodejs-download.png)
 
笔者的 **Windows** 系统，下载完了是如下的 msi 安装包【其他系统自行去官网下载即可】：
![](/images/dev-tool/node-install-package.png)

这里直接双击安装即可，安装完了就可以去配置相关的环境变量了。

当然其他平台，如 Linux，MacOS 可以自行参考[《Node.js 安装配置》](https://www.runoob.com/nodejs/nodejs-install-setup.html)

### 2.3.2 环境变量配置

现在，**Huazie** 以 **windows 11** 系统为例，介绍下配置环境变量，如下：

右击 **Window** 图标，打开下图并选择 **系统**：

![](/images/dev-tool/windows-system.png)

点击 **高级系统设置**，打开系统属性页面，点击 **环境变量** ：

![](/images/dev-tool/windows-env-config.png)

找到 **Path** 系统环境变量，配置上面你的 **Node.js** 的安装目录进去：

![](/images/dev-tool/nodejs-env-config.png)

环境变量配置好之后，我们就可以通过 CMD 命令行，检查：

- `npm -v` ：查看当前安装的 **npm** 的版本号
![](/images/dev-tool/npm-v.png)

- `node -v` : 查看当前安装的 **Node.js** 的版本号
![](/images/dev-tool/node-v.png)

## 2.4 配置 Node.js 开发环境

初次打开 **Sublime Text 3**，我们可以看到如下的界面：

![](/images/dev-tool/sublime-text3-default-page.png)

在菜单栏选择 **Tools** => **Build System** => **New Build System**，打开如下页面

![](/images/dev-tool/sublime-text3-new-buildsystem.png)

```bash
{
	"cmd": ["node",  "$file_name"], 
    "file_regex": "^(..[^:]*):([0-9]+):?([0-9]+)?:? (.*)$", 
    "working_dir": "${file_path}",
    "selector": "source.js",
    "variants": [
	    {
	        "name": "Run",
	        "shell": true,
	        "windows": {
	            "shell_cmd" : "start cmd /c \"node $file_name & echo. & pause\""
	        }
	    }
    ]
}

```

将上述内容保存在，前面打开的 **New Build System** 中，并命名为 **Node.sublime-build**。

## 2.5 编写 Node.js 代码 

现在让我们开始编写第一个 **Node.js** 代码吧！

```javascript
var author = "huazie";
console.log("Hello, World!");
console.log("Author: " + author);
```

针对上述 **Node.js** 代码，我们会新建一个 **helloworld.js** 文件进行保存。

> **注意：** **Node.js** 源码文件也就是 **JavaScript** 源码文件，它的后缀为 `js`
## 2.6 运行 Node.js 代码


菜单栏 **Tools** => **Build System** ，然后 选择 **Node**，就是前面的 **Node.sublime-build**。

然后按住 **Ctrl + Shift + B**，选择 `Node`，直接运行当前的代码，并在下面输出结果，如下所示：

![](nodejs-result.png)

如果按住 **Ctrl + Shift + B**，选择 `Node Run`，则运行当前代码，并弹出命令窗口输出结果，如下所示：


![](nodejs-result-1.png)

通过上面操作之后，我们也可以直接使用 **Ctrl + B**【这里复用上一次 **Ctrl + Shift + B** 选择的 **Build System**】，来直接运行我们的 **JS** 源代码并输出结果。

# 三、总结

本篇 **Huazie** 介绍了 **Sublime Text 3** 配置 **Node.js** 开发环境的相关内容，感兴趣的朋友赶紧配置起来，有任何问题可以随时评论区沟通。




