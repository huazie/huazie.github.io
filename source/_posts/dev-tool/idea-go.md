---
title: Intellij IDEA 配置Go语言开发环境
date: 2021-04-28 17:04:15
updated: 2024-01-14 21:00:22
categories:
  - 开发语言-Go
tags:
  - Go
  - Intellij IDEA
  - 开发环境搭建
---

[《开发工具系列》](/categories/开发工具/)  [《开发语言-Go》](/categories/开发语言-Go/)

![](/images/go-logo.png)

# 1. Go语言环境搭建
本篇博文是在读者 `Go` 自身环境已经搭好，`Intellij IDEA` 环境也已装好的基础上所总结而来。

`Go` 语言环境搭建可参考笔者的另一篇博文 [Go语言学习1-基础入门](/2016/06/27/go/go-learning1/)。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 2. IDEA配置Go语言开发环境

## 2.1 添加Go插件
（1）首先，启动 `Intellij IDEA`，点击 File --> Settings --> Plugins，如下：

![](go-plugins.png)

（2）然后，点击 `Browse repositories`，打开 `Browse Repositories`，并搜索 `go`，这时候出现很多相关的结果，选择如下即可。

![](go-plugins1.png)

（3）点击 Install，等安装好了，提示重启IDEA即可。
![](install.png)
## 2.2 配置GOROOT
点击 File --> Settings --> Languages & Frameworks --> Go --> GOROOT，如下所示：

![](GOROOT.png)

## 2.3 配置GOPATH
点击 File --> Settings --> Languages & Frameworks --> Go --> GOPATH，如果按1中 `Go` 语言环境搭建的步骤，相信到这边的 Global GOPATH 就有了如下截图所示的内容，在下面的 `Project GOPATH` 可以添加我们自己的工程路径。

![](GOPATH.png)

# 3. 新建Go项目
点击 File --> New --> Project，打开 New Project页面，如下截图：

![](new-go-project.png)

选择 **Go**，点击 **Next** 按钮，进入如下页面，填写 项目名称 和 项目路径

![](new-go-project1.png)

点击 `Finish`，选择 以 新窗口 打开新建工程，如下所示：

![](new-go-project2.png)

按照2中配置新建项目的 `GOPATH`，如下截图：

![](GOPATH1.png)

# 4. 编写Go代码
现在可以编写 `Go` 代码了，可以看到如下截图，拥有的代码提示功能，很大程度上方便了开发。

![](code-go.png)

# 5. 运行Go代码
简单编写打印输出代码，然后右键 运行，如下截图所示：

![](run-go.png)

运行结果如下：

![](run-go-result.png)

# 总结
`Intellij IDEA` 配置 `Go` 语言开发环境到此完成【适用于 `GOPATH`】，欢迎大家尝试 ！！！

# 拓展

 从 `Go 1.11` 及其更高版本，`Go` 语言支持 `go mod`，它是 `Go` 语言提供的一个官方包管理工具，用于管理 `Go` 项目中的依赖关系和版本号。通过  `go mod`，开发者可以很方便地管理自己的项目，并且不需要再向 `GOPATH` 中添加第三方的依赖包。在使用 `go mod` 时，开发者可以将项目代码保存在本地文件系统中，不再需要克隆到 `GOPATH` 的指定目录下。同时，`go mod` 还可以从网络上下载并管理所需的依赖包，非常方便快捷。

针对 包管理方式的，后续有机会将讲解，敬请期待！

