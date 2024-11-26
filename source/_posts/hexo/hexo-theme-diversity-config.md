---
title: 来花个几分钟，轻松掌握 Hexo Diversity 主题配置内容
date: 2024-11-15 20:03:15
updated: 2024-11-15 20:03:15
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - hexo-theme-diversity
  - 主题配置
---

![](/images/diversity-logo.png)

# 一、引言

看到这的读者们，如果还没有接入 **Hexo Diversity** 主题，请查看笔者的[《一篇搞定 **Hexo Diversity** 主题接入！支持多主题自由切换！》](../../../../../2024/11/07/hexo/hexo-theme-diversity-integration/)；如果已经接入了，那么为了更好地应用 **Hexo Diversity** 主题，有必要深入了解下与它相关的配置内容。

<!-- more -->

# 二、配置

本篇的配置内容主要依据 V2 版本而来，相比于 V1 版本，多了一些配置。

## 2.1 基础配置

`themes/diversity` 目录下的 `_config.yml` 包含如下配置：

```yml
title: Diversity

description: 博客多样性，一款多主题自由切换的Hexo主题

image: /images/diversity.png

favicon: /images/diversity.ico

back_image: /images/back.jpg

path:
  landscape: /images/landscape.jpg
  phase: /images/phase.png
  light: /images/light.jpg

source:
  landscape: https://github.com/hexojs/hexo-theme-landscape
  phase: https://github.com/hexojs/hexo-theme-phase
  light: https://github.com/hexojs/hexo-theme-light

page:
  blog_scroll_height: 200

back2top:
  enable: true
  enable_scroll_percent: false
  scroll_percent: 5
  position: right
  color: "#fc6423"
  exclude: [next]
```

*   **title** - Diversity主题默认标题
*   **description** - Diversity主题默认描述
*   **image** - 当网页链接被分享到社交平台时显示的图片URL
*   **favicon** - Favicon路径【一个小型图标，用于在浏览器的标签页、地址栏或书签栏中标识和区分不同的网站】
*   **back\_image** - 主题图片翻转后的背景图片
*   **path** - 多主题图片路径【主题名 + 图片路径】。 以 `landscape` 主题举例：
    *   如果该图片路径未配置，默认取 `/images/default.png`
*   **source** - 主题项目来源【用于主题来源按钮点击跳转】
*   **page** - 页面配置
    *   **blog\_scroll\_height** - 博客页滚动高度【单位：`px`】
        *   滚动页面高度大于等于配置高度，隐藏菜单导航栏
        *   滚动页面高度小于配置高度，显示菜单导航栏
*   **back2top** - 返回顶部按钮配置
    *   **enable** - 是否启用，可选值： `true` | `false`
    *   **enable\_scroll\_percent** - 返回顶部按钮中是否启用展示滚动百分比，可选值： `true` | `false`
    *   **scroll\_percent** - 展示返回顶部按钮的最少滚动百分比，建议值： `2 | 3 | 4 | 5`
    *   **position** - 返回顶部按钮展示位置，可选值： `left` | `right`
    *   **color** - 鼠标悬浮或用户触摸时，返回顶部按钮的内容所展示的颜色
    *   **exclude** - 被排除主题，配置中的主题不展示返回顶部按钮

将 `themes/diversity` 目录下的 `_config.diversity.yml`，添加到你的 **Hexo** 项目根目录。

它相较于 `_config.yml`，多了如下配置：

```yml
themes: [landscape,light,phase]

#ports: [5000,5001,5002]
```

*   **themes** - 多主题列表【这里配置主题页面展示的可用于切换的主题】
*   **ports** - 多主题服务器端口列表（不配置，默认从**4001**开始），用于本地 `hexo server` 启动各主题对应的HTTP服务

## 2.2 国际化配置

`themes/diversity` 目录下的 `languages` 目录中的 `zh-CN.yml` 包含如下配置：

```yml
menu:
  blog: 博客
  theme: 主题

button:
  theme-default: 设为默认
  cancel-defalut: 取消默认
  theme-redirect: 主题直达
  theme-source: 主题来源
  back-to-top: 返回顶部

gritter:
  title-theme: 主题【{0}】
  text-configured: 已设置
  text-canceled: 已取消
  text-click-to-jump: 点击跳转

no-theme:
  tip-text: 您还没有设置默认主题！点击下方按钮前往设置
  btn-text: 主题选择

introduction:
  landscape: Hexo 中的一个全新的默认主题，需要 Hexo 2.4 或者 更高的版本。
  phase: 通过 Phase，感受时间流逝，它是 Hexo 最美丽的主题。
  light: Hexo 中的一个简约主题。
```

*   **menu** - 导航栏菜单展示名称
*   **button** - 按钮文本
*   **gritter** - 主题选择页的提示文本
*   **no-theme** - 无主题页的文本
*   **introduction** - 主题介绍【如果没有配置，则不展示介绍】。新增一个主题接入，这里的主题介绍需要对应新增。如果新增主题没有配置介绍，则主题选择页的卡片不展示主题介绍。

## 2.3 多主题相关配置

在我们的 **Hexo** 项目根目录中，添加 **config** 目录，并为上述多主题列表中的每个主题添加一个对应主题名的配置目录，同时在该主题名的配置目录下添加对应的 `_config.yml` 【它可以直接从你原来项目根目录下的 `_config.yml` 复制过来即可】，形如：

```pre
├─config
│  ├─landscape
│  │  ├─_config.yml
│  ├─light
│  │  ├─_config.yml
│  ├─phase
│  │  ├─_config.yml
```

> 注意：这里除了将各主题的配置独立开来，同时也为了将自动生成的 `db.json` 独立开来，保证各**Hexo实例**运行时互不干扰。

修改上述各主题配置目录下的 `_config.yml`，以 **landscape** 举例：

```diff
_config.yml
- url: http://example.com
+ url: http://example.com/landscape

- public_dir: public
+ public_dir: public/landscape

- theme: other-theme
+ theme: landscape
```

在你的 **Hexo** 项目根目录下，我们依旧可以添加不同主题独立的 `_config.[theme].yml` 文件，更多了解请查看官方[《配置》](https://hexo.io/zh-cn/docs/configuration)

针对不同主题，可在各自配置中启用分类和标签生成配置

```yml
category_generator:
  enable_index_page: true
  layout: category-index
  per_page: 10
  order_by: -date
```

*   **category\_generator** - 分类生成配置
    *   **enable\_index\_page** - `true` 【启用分类首页生成, 通常是 `/categories/index.html`]
    *   **layout** - 分类首页布局。 如果不配置，则默认为 `category-index`
    *   **per\_page** - 每页展示条数
    *   **order\_by** - 默认按日期降序排列（新到旧）

```yml
tag_generator:
  enable_index_page: true
  layout: tag-index
  per_page: 100
  order_by: -date
```

*   **tag\_generator** - 标签生成配置
    *   **enable\_index\_page** - `true` 【启用标签首页生成, 通常是 `/tags/index.html`]
    *   **layout** - 标签首页布局。 如果不配置，则默认为 `tag-index`
    *   **per\_page** - 每页展示条数
    *   **order\_by** - 默认按日期降序排列（新到旧）

# 三、结语

目前 **Diversity** 主题涉及的配置不多，花个几分钟，大家基本都能轻松掌握。

如果您在阅读配置过程中，有啥疑问，欢迎来评论区和我讨论！！！

