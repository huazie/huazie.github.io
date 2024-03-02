---
title: Sublime Text 3 解决中文乱码问题
date: 2023-12-15 09:42:44
updated: 2024-01-29 10:58:09
categories:
  - 开发工具
tags:
  - Sublime Text 3
  - 中文乱码
---

[《开发工具系列》](/categories/开发工具/)

![](/images/sublime-text3-logo.png)

# 一、引言

在 [《Sublime Text 3配置C/C++开发环境》](/2023/12/15/dev-tool/sublime-text3-c-cplus/) 博文中，**Huazie** 带大家利用 **Sublime Text 3** 配置了 **C/C++** 开发环境，相信大家都已经开始使用 **Sublime Text 3** 运行 **C/C++** 了，但是慢慢地使用过程中，大家可能发现，如果输出内容包含中文，打印出来的信息确是乱码的，如下图所示：

![](chinese-garbled-code.png)

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、主要内容

## 1. 初识 ConvertToUTF8 插件

**那么上述中文乱码的问题，我们该如何解决呢？**

这里就不得不提到一个 **Sublime Text** 下的插件 --- **ConvertToUTF8**，它是一个用于将文件编码转换为 **UTF-8** 格式的插件。

使用此插件，可以编辑和保存 **Sublime Text** 当前不支持的编码文件，如 **GB2312、GBK、BIG5、EUC-KR、EUC-JP** 等。


## 2. 安装 ConvertToUTF8 插件

菜单栏选择 **Preferences** => **Package Control** 或者 **按住 `Ctrl+Shift+p`**，弹出如下输入窗口，在其中输入 `install package`，并选中红框内的列表。

> 如果安装过 **Package Control** 可以忽略，没有安装过的朋友，请参考[《Sublime Text 3 中安装Package Control并配置》](https://zhuanlan.zhihu.com/p/349113898)

![](install-package.png)

点击之后，正在加载插件库【加载过程缓慢，耐心等待一会】：

![](install-package-1.png)


然后在弹出的输入窗口中，输入 `ConvertToUTF8` 并回车，点击红框处，即可开始安装；

![](convert-to-utf8.png)

在左下角会显示正在安装  `ConvertToUTF8` 中【耐心等候一会儿】：

![](convert-to-utf8-1.png)

弹出如下页面，即表示安装成功，接着我们直接重启 **Sublime Text 3**即可。

![](convert-to-utf8-2.png)

我们在菜单栏 **Preferences** => **Package Setting** => **CovertToUTF8** 下可以添加或修改 **CovertToUTF8** 插件相关的设置：

- **Default** ： 默认设置
- **User** : 个人自定义设置

![](settings.png)

## 3. 中文乱码问题解决

我们打开之前新建的第一个 **C** 文件，修改如下：

```c
#include<stdio.h>

int main() 
{
	printf("hello world!\n");
	printf("[C]作者: Huazie");
	return 0;
}
```

这时如果选择直接编译运行的话，那结果还是乱码的，如下：

![](c-result.png)

这个时候，我们需要在菜单栏点击 **File**，选择 **Reload with Encoding** ，再选择 **GBK**：

![](convert-to-gbk.png)

这时，会出现如下的弹出框，点击 **OK** 即可

![](convert-to-gbk-1.png)

点完之后，我们的文件编码已经改变，包括中文都已经显示为乱码，如下：

![](convert-to-gbk-2.png)

> **注意：** 在使用 **ConvertToUTF8**  插件之前，建议先备份原始文件，以防转换过程中出现问题。

由于我们这里比较简单，只需要重新将上面的代码复制过来，保存之后，选择 **C Build System**，按住 **Ctrl + B**，直接运行当前的 **C** 代码，运行结果如下图所示：

![](c-result-1.png)

从上图可以看出，这里的中文已经能够正常输出了，到这一步，中文乱码问题算是彻底解决了。
# 三、总结

上述中文乱码问题的解决，不仅仅适用于 **C/C++** 代码，也适用于其他任何 **Sublime Text 3** 集成的开发环境。如果你也有同样的问题，不妨装上 **ConvertToUTF8** 插件试试吧！