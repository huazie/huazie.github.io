---
title: Sublime Text 3中的代码格式化技巧大揭秘
date: 2024-11-22 23:40:21
updated: 2024-11-22 23:40:21
categories:
  - 开发工具
tags:
  - Sublime Text 3
  - 代码格式化
  - SublimeJsPrettier
  - Prettier
  - Beautify
  - Javascript Beautify
  - HTMLBeautiry
  - clang-format
  - black
---

![](/images/sublime-text3-logo.png)

# 前言

最近 **Huazie** 一直在用 **Sublime Text 3** 开发 [Hexo Diversity](https://github.com/huazie/hexo-theme-diversity) 主题开源项目，当找到一些解决方案的代码并拷贝过来时，总会遇到缩进和换行问题，此时复制的代码就显得杂乱无章的。开始笔者一般选择在线代码格式化工具处理过再复制过来，慢慢用的多了，就想要能够在 **Sublime Text 3** 中直接可以格式化，以此来提升效率。

讲到这，就引申出来本篇要介绍的内容了。

<!-- more -->

# 1. 使用内置的格式化功能

**Sublime Text 3** 自带了一些基本的格式化功能，适用于 `HTML`、`CSS` 和 `JavaScript` 等语言，我们可以通过如下操作来使用这些功能：

1.  打开你想要格式化的文件 或者 选中你要格式化的代码片段；
2.  选择菜单栏中的 `Edit` -> `Line` -> `Reindent`。

下面让我们来格式化一下 `JavaScript` 代码，如下动图：

![](reindent.gif)

# 2. 安装并使用第三方插件

**Sublime Text 3** 有丰富的插件生态，下面 **Huazie** 介绍一些常用的格式化插件供大家使用。

## 2.1 SublimeJsPrettier

**SublimeJsPrettier** 是一个为 **Sublime Text** 开发的插件，它集成了 **Prettier** 的代码格式化功能。

> **Prettier** 是一个流行的代码格式化工具，支持多种编程语言，包括但不限于`JavaScript、JSX、TypeScript、CSS、SCSS、HTML`等。

### 2.1.1 全局安装 Prettier

在使用 **SublimeJsPrettier** 之前需要全局安装 **Prettier**，可以通过 `npm` 进行安装，如下所示：

```bash
npm install -g prettier
```

> 当然这里执行 `npm` 命令之前，我们需要首先安装 `Node.js` ，并配置相关环境变量，这里可以参考笔者的[《【实操】基于 GitHub Pages + Hexo 搭建个人博客》](https://juejin.cn/post/7373226679731240970#heading-1) 中的第二章节。

![](prettier.png)

### 2.1.2 安装 Package Control

这个如果已经安装过的朋友可以直接跳过。关于 `Package Control` 的安装可以参考[《Sublime Text 3 中安装Package Control并配置》](https://zhuanlan.zhihu.com/p/349113898)

### 2.1.3 安装 SublimeJsPrettier

1.  首先打开我们的 **Sublime Text** 编辑器；

![](/images/dev-tool/sublime-text3-default-page.png)

2.  按下 `Ctrl+Shift+P`（Windows/Linux）或 `Cmd+Shift+P`（Mac）打开命令面板；

![](/images/dev-tool/sublime-text3-ctrl-shift-p.png)

3.  在输入框中，输入 **Package Control: Install Package** 并选择，弹出另一个搜索框；

![](/images/dev-tool/sublime-text3-install-package.png)

4.  继续在搜索框中，输入 **JsPrettier**，然后选择并开始安装；

![](sublime-text3-jsprettier.png)

5.  安装成功后，会显示如下页面。

![](jsprettier-success.png)

### 2.1.4 使用 SublimeJsPrettier

有如下三种使用 **SublimeJsPrettier** 格式化代码的方法：

*   **通过命令面板**，即按下 `Ctrl+Shift+P`（Windows/Linux）或 `Cmd+Shift+P`（Mac），然后输入 `JsPrettier Format Code`；

![](jsprettier-format-code.gif)

*   **通过上下文菜单**，在当前要格式化的文件视图中右键点击，然后选择 `JsPrettier Format Code`;

![](jsprettier-format-code-1.gif)

*   **通过键盘绑定**，由于该插件没有默认的键盘绑定，因此需要手动添加；打开 **Preferences（偏好设置**）-> **Key Bindings...（键盘绑定……）**，并为 `js_prettier` 添加一个条目。例如【这里可以将 `ctrl+alt+l` 替换为大家喜欢的键盘组合】:

```json
{ "keys": ["ctrl+alt+l"], "command": "js_prettier" }
```

![](jsprettier-key-bindings.png)

## 2.2 Beautify

**Beautify** 是 **Sublime Text** 中用于代码美化的工具，可以帮助开发者自动整理和美化代码。它能够根据代码的语法规则自动调整缩进、排序标签、整理注释，并去除不必要的空格和空白。支持众多的语言，如`HTML、CSS、JavaScript、PHP、SQL`等。

### 2.2.1 安装 Beautify

1.  打开 **Sublime Text**，按下 `Ctrl+Shift+P`（Windows/Linux）或`Cmd+Shift+P`（Mac）打开命令面板；

![](/images/dev-tool/sublime-text3-ctrl-shift-p.png)

2.  在命令面板中，输入`Package Control: Install Package`并回车，进入插件搜索和安装界面；

![](/images/dev-tool/sublime-text3-install-package.png)

3.  在搜索框中，输入 `Beautify`，选择需要的插件并安装，如`Javascript Beautify`、 `HTMLBeautiry`等。

![](sublime-text3-beautify.png)

### 2.2.2 使用 Beautify

**Huazie** 选择安装了 `Javascript Beautify`，那这里可以：

* 直接在当前要格式化的文件视图中右键点击，然后选择 `Javascript Beautify`；
  或者直接选中要格式化的代码，然后右键点击，然后选择 `Javascript Beautify` ；

![](javascript-beautify.gif)

* 使用快捷键【`Ctrl+Alt+F`】，这里如果与其他应用有冲突，需要先关掉对应的应用再试试。

# 3. 集成外部代码格式化工具

## 3.1 集成 clang-format

### 3.1.1 安装 clang-format

* 访问 **LLVM** 的官方网站（例如：[LLVM官方GitHub发布页面](https://github.com/llvm/llvm-project/releases) 或 [官方下载页面](https://releases.llvm.org/)）或  [LLVM Snapshot Builds](https://llvm.org/builds/) ，下载适合你操作系统的**LLVM**安装包；
* 安装 **LLVM**，这里面包含 **clang-format** 工具；
* 上述安装好之后，在命令行中输入 `clang-format --version`，以验证 **clang-format** 是否正确安装。

![](clang-format-version.png)

### 3.1.2 配置 clang-format

1.  打开 **Sublime Text**，选择菜单栏中的 `Tools` -> `Build System` -> `New Build System`，打开如下页面：

![](/images/dev-tool/sublime-text3-new-buildsystem.png)

2.  在上述打开文件中，添加如下内容：

```sublime-build
{
    "cmd": ["clang-format", "-i", "$file"],
    "selector": "source.c, source.cpp, source.objc, source.objcpp"
}
```

3.  保存文件，并命名为 `ClangFormat.sublime-build`

### 3.1.3 使用 clang-format

1.  选择菜单栏中的 `Tools` -> `Build System` -> `ClangFormat`;
2.  打开你需要格式化的 **C/C++** 源码文件，并按住 `Ctrl + B`，稍等一会，即可看到格式化后的效果。

![](clang-format.gif)

## 3.2 集成 black

### 3.2.1 安装 Black

1.  使用 `pip` 命令，安装 **black**，如下：

![](pip-install-black.png)

2.  查看 **black** 位置，并输出 **black** 版本验证已安装并配置；

```bash
# Windows
where black

black --version
```

![](where-black-version.png)

### 3.2.2 安装 python-black

1.  首先打开我们的 **Sublime Text** 编辑器，并按下 `Ctrl+Shift+P`（Windows/Linux）或 `Cmd+Shift+P`（Mac）打开命令面板；
2.  在输入框中，输入 **Package Control: Install Package** 并选择，弹出另一个搜索框；
3.  继续在搜索框中，输入 **Black**，然后选择如下并开始安装；

![](black-formatter.png)

### 3.2.3 使用 python-black

* **python** 代码修改保存后会自动格式化【个人使用还是觉得不是什么缩进都能格式化，有懂行的可以评论区讨论一下】
* 在打开的 **python** 代码处，右击鼠标，选择 `Black:Format`

# 总结

通过内置的格式化功能、第三方插件或集成外部代码格式化工具，**Sublime Text 3** 提供了多种方式来满足开发者的代码格式化需求。我们选择哪种方法取决于具体的编程语言、个人偏好以及项目的需求。

如果你还有更好的方法，不妨在评论区和大家分享一下吧！

