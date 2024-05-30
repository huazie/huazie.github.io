---
title: Sublime Text 3配置 Java 开发环境
date: 2024-01-26 17:33:20
updated: 2024-01-29 15:46:24
categories:
  - [开发工具]
  - [开发语言-Java]
tags:
  - Java
  - Sublime Text 3
  - 开发环境搭建
---


![](/images/java-logo.png)

# 一、引言
**Java** 是一种跨平台、面向对象、功能强大且安全可靠的编程语言。它有很多常用的开发工具，比如 **Eclipse**、**IDEA** 等等，相信大家多多少少都有所涉猎；而本篇 **Huazie** 将要介绍一个比较轻量级的开发工具 **Sublime Text 3**，并用它来配置 **Java** 开发环境。

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

## 2.2 初识 Java

**Java** 是一种面向对象的编程语言，它诞生于 **1995** 年，由**Sun Microsystems** 公司（现已被甲骨文公司收购）开发，最初是用于智能家电平台上运行的 **OAK** 语言，后来发展成为一种功能强大的编程语言。

**Java** 语言的设计初衷是 **“一次编写，到处运行”**，即源代码只需编写一次，就可以在任何支持 **Java** 的平台上运行。

**Java** 语言拥有丰富的API库和工具，广泛应用于 **Web** 开发、移动应用开发、游戏开发、大数据和云计算等领域。


## 2.3 接入 Java

在开始接入 **Java** 之前，还有个概念需要明确一下，那就是 **JDK**（全称是 **Java Development Kit**），即 **Java 开发工具包**。

它是用于开发 **Java** 程序的一套工具和库，其中包含了如下的内容：
- **Java编译器（javac）**：通过 **javac** 命令，将 **Java** 源代码（**.java** 文件）编译成字节码文件（**.class** 文件）。
- **Java虚拟机（JVM）**：通过 **java** 命令，启动 **JVM**，并加载执行 **Java** 字节码文件。从 **JDK1.8** 开始，可以直接使用 `java 主类名.java` 运行 **Java** 文件【其中包含了编译源代码、执行字节码等步骤】。
- **Java基础类库（Java API）**：包括 `java.util、java.io、java.net` 等常用类库，用于支持各种常见的编程任务。
- **其他工具和实用程序**：如 **Java调试器（jdb）**、**Java反编译器（javap）**、**Java文档生成器（javadoc）** 等。

### 2.3.1 JDK 下载

[JDK 官网下载](https://www.oracle.com/java/technologies/downloads/)，目前最新版为 **JDK21**。

大家可以按照自己系统，选择相应的版本进行下载：

![](jdk-download.png)

以 **Windows** 为例：

- `x64 Compressed Archive` ： JDK的免安装版本
- `x64 Installer` ：JDK的离线安装版本
- `x64 MSI Installer` ：JDK的离线安装版本
### 2.3.2 安装和使用 java

以 **Windows** 为例：

- 如果是免安装版本，只需要解压之后，将对应的文件复制到指定的文件夹，比如 `C:\java`

- 如果是安装版本，那就下载之后，双击安装，同样选择一个指定的文件夹进行安装即可。


**Huazie** 的 **Windows** 系统上就安装了多个 **JDK** 版本，大家按照自身需要，自行选择安装和使用。

![](jdk-directory.png)

### 2.3.3 环境变量配置


现在，**Huazie** 以 **windows 11** 系统为例，介绍下配置环境变量，如下：

右击 **Window** 图标，打开下图并选择 **系统**：

![](/images/dev-tool/windows-system.png)

点击 **高级系统设置**，打开系统属性页面，点击 **环境变量** ：

![](/images/dev-tool/windows-env-config.png)

新增 **JAVA_HOME** 的环境变量【如果要更换当前的 **JDK** 版本，修改这里即可】：

![](java-home.png)

新增 **CLASSPATH** 的环境变量【用于告诉 **JVM** 在哪些目录下查找类文件】：

![](classpath.png)
- **当前目录（.）**：表示在当前目录下查找类文件。
- **Java类库路径**：包括Java运行环境提供的类库（如rt.jar、tools.jar等）和第三方类库。这些类库文件通常以.jar或.zip格式存在，并需要指定它们所在的目录路径。
- **自定义类文件路径**：如果开发者编写了自己的 **Java** 类文件，也可以将这些类文件所在的目录路径也加入到 **CLASSPATH** 中。



找到 **Path** 环境变量，配置上面你的 **Java** 安装目录的 **lib** 目录进去：

![](java-env-config.png)



如果上面是 **JDK** 离线安装版，**Path** 里面可能存在如下环境变量，则需要手动删除，以免影响上面的环境变量配置。

![](javaoraclepath.png)



然后 **Win + R**，打开如下窗口，输入 **cmd**，点击确认打开命令行窗口。

![](/images/dev-tool/win-r-cmd.png)

在命令行窗口内，输入 `java -version` 查看，如下图所示即为安装成功：

![](java-version.png)



## 2.4 配置 Java 开发环境

初次打开 **Sublime Text 3**，我们可以看到如下的界面：

![](/images/dev-tool/sublime-text3-default-page.png)

在菜单栏选择 **Tools** => **Build System** => **New Build System**，打开如下页面

![](/images/dev-tool/sublime-text3-new-buildsystem.png)

```bash
{
    "cmd": ["javac",  "$file_name"], 
    "file_regex": "^(..[^:]*):([0-9]+):?([0-9]+)?:? (.*)$", 
    "working_dir": "${file_path}",
    "selector": "source.java",
    "variants": [
        {
            "name": "Run",
            "shell": true,
            "windows": {
                "shell_cmd" : "start cmd /c \"java ${file_base_name} & echo. & pause\""
            }
        },
        {
            "name": "Build and Run",
            "shell": true,
            "windows": {
                "shell_cmd": "javac \"${file_name}\" && start cmd /c \"java ${file_base_name} & echo. & pause\""
            }
        }
    ]
}
```

将上述内容保存在，前面打开的 **New Build System** 中，并命名为 **Java8.sublime-build**【这里命名可以按自己的 **JDK** 版本来】。


## 2.5 编写 Java 代码 

现在让我们开始编写第一个 **Java** 代码吧！

```java
public class HelloWorld {
    public static void main(String[] args) {
        String name = "Huazie";
        System.out.println("Hello World!");
        System.out.println("Author:" + name);
    }
}
```

针对上述 **java** 代码，我们会新建一个 **HelloWorld.java** 文件进行保存。

> **注意：** **java** 源码文件名的后缀为 `java`，且文件名必须和类名保持一致。
## 2.6 编译和运行 Java 代码

上述 **HelloWorld.java** 我们也可以打开命令行窗口进行编译运行，如下图：

![](java-cmd-result.png)

当然，对于初学者，使用这种方式理解 **Java** 的编译和运行过程，还是可以的。

不过既然我们已经使用了 **Sublime Text 3** 的开发环境了，那就在菜单栏 **Tools** => **Build System** ，然后 选择 **Java8**，就是前面的 **Java8.sublime-build**。

然后直接按住 **Ctrl + Shift + B**，会弹出如下界面：

![](java-build-system.png)
选择 `Java8`，我们可以直接编译当前 **HelloWorld.java** 源码文件，并在当前目录生成对应的字节码文件，如下图：

![](java-build-system-1.png)

接着，还是按住 **Ctrl + Shift + B**，选择 `Java8 - Run`，就可以在 **CMD** 窗口中运行我们上面编译好的 **Java 字节码文件**，并输出相关内容。

如下图所示：

![](java-build-system-result.png)

当前上面还是需要两次操作，我们把两次整合一下，还是按住 **Ctrl + Shift + B**，选择 `Java8 - Build And Run`，这一次就可以直接编译和运行一起【这一步为了看到效果，先将之前生成的字节码文件删掉，然后再操作即可】。

通过上面操作之后，我们就可以直接使用 **Ctrl + B**【这里复用上一次 **Ctrl + Shift + B** 选择的 `Java8 - Build And Run`】，这样就可以直接编译和运行我们的 **Java** 代码。

> **注意：**  这里的编译执行不适合有包的情况，涉及到比较复杂的逻辑，还是采用 **Eclipse** 和 **IDEA** 这些专门的开发工具较为适合。

## 2.7 乱码问题
经过上面的配置，相信大家都能编译和运行第一个 **Java** 代码了，但是有些小伙伴慢慢使用发现，如果输出的内容包含中文，打印出来的信息是乱码的。有关这个问题，请查看笔者的另一篇博文 [《Sublime Text 3 解决中文乱码问题》](/2023/12/15/dev-tool/sublime-text3-convert-to-utf8/)

# 三、总结

本篇 **Huazie** 介绍了 **Sublime Text 3** 配置 **Java** 开发环境的相关内容，感兴趣的朋友赶紧配置起来，有任何问题可以随时评论区沟通。



