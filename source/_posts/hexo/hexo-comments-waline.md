---
title: 50kB 搞定全功能评论系统！Hexo Comments Waline 插件实战
date: 2026-10-27 11:10:00
updated: 2026-10-27 11:10:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-comments-waline
  - Waline
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 一、前言

系列最后一篇，登场的评论插件是 **Waline** —— 六款插件中功能最丰富的一员：评论、阅读量、点赞、表情反应一站配齐，而前端只需 **50kB**。

跟着 **Huazie** 一起上手 [hexo-comments-waline](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-waline)，给这个系列收个尾！

<!-- more -->

# 🔍 二、认识 Waline

[`Waline`](https://waline.js.org) 是一个简单但功能强大的评论解决方案，定位是"轻量后端起评"：

* **轻量简洁**：仅需 50kB 大小，快速加载，优化性能
* **多部署方式**：支持 Vercel、Netlify、Docker、自建服务器，灵活选择
* **功能丰富**：支持评论、阅读量、点赞、表情反应，一站式互动
* **多语言支持**：内置 30+ 语言，国际化友好

与 Twikoo 类似，Waline 也需要一个**服务端**来存取评论数据，好消息是 Vercel 免费额度即可跑起来。

# ✨ 三、插件功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **轻量简洁** | 仅需 50kB 大小 | 快速加载，优化性能 |
| **多部署方式** | 支持 Vercel、Netlify、Docker、自建服务器 | 灵活选择，易于部署 |
| **安全可靠** | 完全开源，支持多种登录方式 | 保护用户隐私，透明可信 |
| **功能丰富** | 支持评论、阅读量、点赞、表情反应 | 一站式互动解决方案 |
| **主题切换** | 支持亮色/暗色主题自动切换 | 完美适配各种主题风格 |
| **多语言支持** | 内置 30+ 语言支持 | 国际化友好 |
| **易于配置** | 简单的 YAML 配置 | 快速上手，灵活定制 |

# 🚀 四、快速上手

## 1. 部署 Waline 服务端

按 [Waline 官方文档](https://waline.js.org) 选择一种方式部署服务端（Vercel / Netlify / Docker / 自建服务器均可），部署完成后即可获得**服务端地址（serverURL）**。

## 2. 安装插件

在 Hexo 项目根目录执行：

```bash
# 1. 安装多评论系统核心插件（必需）
npm install hexo-generator-comments --save

# 2. 安装 Waline 评论插件
npm install hexo-comments-waline --save
```

> **提示**：`hexo-comments-waline` 需要与 `hexo-generator-comments` 搭配使用。

## 3. 添加配置

在 Hexo 站点配置 `_config.yml` 或主题配置 `_config.yml`、`_config.[theme].yml` 中添加以下内容：

```yaml
waline:
  # 是否启用 Waline 评论系统
  enable: false
  # 是否启用加载提示
  loading: true
  # Waline 服务端地址，必须配置
  server_url: https://your-waline-server.example.com
  # Waline JS 文件 URL
  js_url: https://unpkg.com/@waline/client@v3/dist/waline.js
  # Waline CSS 文件 URL
  css_url: https://unpkg.com/@waline/client@v3/dist/waline.css
  # 当前文章的唯一标识，用于区分不同页面的评论
  path: pathname
  # 评论区语言
  lang: zh-CN
  # 表情包设置，支持数组或 false 禁用
  emoji:
  # 暗黑模式适配
  dark: auto
  # 评论时需要填写的信息字段
  meta: ['nick', 'mail', 'link']
  # 必填的字段数组
  required_meta: []
  # 评论排序方式
  comment_sorting: latest
  # 登录配置
  login: enable
  # 评论字数限制
  word_limit: false
  # 每页评论条数
  page_size: 10
  # 表情包搜索功能
  search: false
  # 是否隐藏页脚版权信息
  no_copyright: false
  # 是否隐藏 RSS 订阅链接
  no_rss: false
  # 表情反应功能
  reaction: true
```

> **重要**：请将 `server_url` 替换为你的实际 Waline 服务端地址。

## 4. 关键配置项说明

| 选项 | 类型 | 默认值 | 必填 | 描述 |
|------|------|--------|------|------|
| `enable` | Boolean | `false` | 是 | 是否启用 Waline 评论系统 |
| `loading` | Boolean | `true` | 否 | 是否启用加载提示（评论加载时显示加载动画） |
| `server_url` | String | - | 是 | Waline 服务端地址（Vercel/Docker/自建） |
| `js_url` / `css_url` | String | unpkg CDN | 否 | Waline 客户端 JS / CSS 资源地址 |
| `path` | String | `pathname` | 否 | 评论区路径，用于区分不同的页面 |
| `lang` | String | `zh-CN` | 否 | 评论区的语言（界面文本） |
| `dark` | String | `auto` | 否 | 暗黑模式适配，`auto` 为自动适配 |
| `page_size` | Number | `10` | 否 | 每页评论条数 |
| `reaction` | Boolean | `true` | 否 | 表情反应功能开关 |

完整配置选项请查阅插件的 [README 文档](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-waline/README.md)。

# ⚠️ 五、注意事项

* **先部署后接入**：`server_url` 来自已部署的 Waline 服务端，请确保服务可用后再启用评论。
* **meta 与 required_meta 搭配**：`meta` 控制评论时展示哪些信息字段（昵称/邮箱/链接），`required_meta` 从中挑选必填项。
* **表情与搜索**：`emoji` 支持传入表情包数组，设为 `false` 可禁用；`search` 可开启表情包搜索功能。
* **客户端资源可自持**：`js_url` / `css_url` 默认走 unpkg CDN，必要时可替换为自部署地址，避免 CDN 波动。

# 📝 六、小结

至此，**diversity-plugins** 的六款评论插件就全部介绍完了。作为收官的 [hexo-comments-waline](https://www.npmjs.com/package/hexo-comments-waline)，它以 50kB 的身躯提供了**评论、阅读量、点赞、表情反应**的一站式互动能力，是功能党的不二之选。

最后帮你快速决策：读者都是开发者 → **Utterances / Giscus**；想要完整评论管理 → **Gitalk**；数据自主可控 → **Twikoo**；功能最全 → **Waline**。配上核心插件 `hexo-generator-comments`，还可以让访客在多系统间自由切换 —— 这正是这个插件家族的设计初衷。快去给你的博客挑一款吧！
