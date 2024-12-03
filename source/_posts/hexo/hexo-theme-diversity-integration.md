---
title: 一篇搞定 Hexo Diversity 主题接入！支持多主题自由切换！
date: 2024-11-07 21:08:45
updated: 2024-11-07 21:08:45
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - hexo-theme-diversity
---

![](/images/diversity-logo.png)

# 一、引言

提到博客框架，就不得不说 **Hexo**，它是一个快速、简洁且高效的博客框架。目前 **Hexo** 拥有 **400+** 的主题，使用者可以在众多的主题中选择一个应用到自己的框架中。

因为框架限制，我们在线运行时只能选择其中一个主题来展示自己的博客。当然也可以部署多套不同主题的环境，但这不是最好的方案。

下面我们要接入的 **Hexo Diversity** 主题，就是一个可以实现多主题自由切换的解决方案。

<!-- more -->

# 二、主要内容

## 2.1 Hexo 项目搭建

由于本篇主要介绍 **Hexo Diversity** 主题的接入工作，因此需要各位将自己的 **Hexo** 项目先搭建好。

有关 **Hexo** 项目的搭建：

*   大家可以参考[官方文档](https://hexo.io/zh-cn/docs/)介绍
*   或者也可以查看 **Huazie** 的[【实操】基于 GitHub Pages + Hexo 搭建个人博客](../../../../../2024/01/17/hexo/hexo-github-pages)

## 2.2 Hexo Diversity 主题接入

本篇博文截止，[Diversity](https://github.com/huazie/hexo-theme-diversity/releases/tag/v2.2.7) 主题的版本是 **2.2.7**。

从整体上看，目前 **Diversity** 主题一共两大版本，下面分别介绍下：

### 2.2.1 V1版本接入

**V1**版本只提供一个主题选择页面，它实现了多个 **Hexo** 主题自由切换的核心功能。

主题选择页包含如下：

*   支持鼠标左右拖拽滚动展示【**phone**端支持左右滑动展示】
*   支持鼠标滚轮前后滚动展示
*   已配置的主题图片，支持鼠标悬停展示详情页【**phone**端触摸展示详情页】，可用于跳转对应主题博客页面，主题来源项目

**下面跟着 Huazie 来看看接入 V1 版本都需要做哪些操作？**

*   首先，在你的**Hexo**项目根目录下，执行以下命令，并将 `_config.yml`中的`theme`修改为`diversity`;

```shell
git clone --depth 1 https://github.com/huazie/hexo-theme-diversity/tree/v1 themes/diversity
```

```diff
_config.yml
- theme: other-theme
+ theme: diversity
```

*   然后，将`themes/diversity`目录下的`_config.diversity.yml`，移动到你的**Hexo**项目根目录。该文件主要配置 **多主题列表** 和 **多主题服务器端口列表**;

```yml
themes: [landscape,light,phase]

#ports: [5000,5001,5002]
```

*   接着，在你的**Hexo**项目根目录，添加**config**目录，为上述多主题列表中的每个主题添加一个对应主题名的配置目录，并在该配置目录下添加对应的`_config.yml`【直接从你原来项目根目录下的`_config.yml`复制一份即可】，形如：

```txt
├─config
│  ├─landscape
│  │  ├─_config.yml
│  ├─light
│  │  ├─_config.yml
│  ├─phase
│  │  ├─_config.yml
```

修改上述各主题目录下的`_config.yml`，以**landscape**举例：

```diff
_config.yml
- url: http://example.com
+ url: http://example.com/landscape

- public_dir: public
+ public_dir: public/landscape

- theme: other-theme
+ theme: landscape
```

将上述各主题相关的源码下载到项目根目录的 `themes` 下，并各自命名，如下所示【其中`diversity`是第一步所下载的】：

```txt
├─themes
│  ├─diversity
│  ├─light
│  ├─phase
```

其中 `landscape` 主题通过项目依赖导入，如下：

```json
"dependencies": {
    "hexo-theme-landscape": "^1.0.0"
 }
```

*   最后，本地 **Hexo** 站点，执行如下命令：命令执行完后，通过浏览器访问 <http://localhost:4000> 即可展示你的站点【更多配置了解查看 [V1版本](https://github.com/huazie/hexo-theme-diversity/tree/v1)】。

![](v1-hexo-clean.png)

![](v1-hexo-server.png)

### 2.2.2 V2 版本

**V2**版本提供导航栏菜单，目前包含 **博客** 和 **主题** 两大导航栏菜单。

*   **博客** 导航菜单用于展示设置的默认主题博客页面，如果没有设置默认主题，则展示无主题页【用于跳转主题选择页】
*   **主题** 导航菜单除了V1版本中的功能，还可以设置默认主题【弹出的提示信息中可点击跳转个人博客页】、取消默认主题。

在V1版本的基础上，针对主题选择页做了如下优化：

*   **电脑端**

    *   支持鼠标左右拖拽滚动展示
    *   支持鼠标滚轮前后滚动展示

*   **手机端**

    *   支持左右箭头切换滚动展示
    *   支持左滑右滑切换滚动展示

**那么现在，让我们看看 V2 版本该如何接入你的 Hexo 项目中 ？**

*   首先，同样是在你的**Hexo**项目根目录下，执行以下命令，并将 `_config.yml`中的`theme`修改为`diversity`;

```shell
git clone --depth 1 https://github.com/huazie/hexo-theme-diversity/tree/v2 themes/diversity
```

*   然后，将`themes/diversity`目录下的`_config.diversity.yml`，移动到你的**Hexo**项目根目录。与 **V1** 版本不同的是，**V2** 版本中，我们还需要将`themes/diversity`目录下的`other`目录中的目录和文件复制或移动到你的**Hexo**项目根目录。

*   接着，可以直接复用 **V1** 版本中新建的 `config` 目录和其他主题的源码; 针对不同主题，可在各自主题配置中启用分类和标签生成配置【更多配置了解查看 [V2版本](https://github.com/huazie/hexo-theme-diversity/tree/v2)】。

```txt
├─themes
│  ├─diversity
│  ├─light
│  ├─next
│  ├─phase
```

*   最后，本地 **Hexo** 站点，执行如下命令：命令执行完后，通过浏览器访问 <http://localhost:4000> 即可展示你的站点。

![](v2-hexo-clean.png)

![](v2-hexo-server.png)

# 三、总结

**Hexo Diversity** 主题只需要简单几步操作就可以接入，快来试试吧！！！

在接入过程中，有任何的问题，欢迎来评论区和我讨论！！！