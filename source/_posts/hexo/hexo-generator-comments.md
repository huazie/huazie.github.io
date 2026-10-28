---
title: 一套配置通吃六大评论系统！深入解读 Hexo Generator Comments 核心插件
date: 2026-10-27 10:10:00
updated: 2026-10-27 10:10:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-generator-comments
  - 留言评论
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 前言

之前 **Huazie** 在 [《个人博客支持留言评论，来看看 Hexo Diversity 主题如何实现的！》](../../../../../2025/05/08/hexo/hexo-theme-diversity-comments/) 一文中，介绍过 **Diversity** 主题是如何集成多种评论系统的。

当时这套评论能力还"长"在主题里，只有使用 **Diversity** 主题的朋友才能享受。现在，它已经被抽离成了独立插件 —— [hexo-generator-comments](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-generator-comments)，**任何 Hexo 主题都能用了**！

接下来就跟着 **Huazie** 一起深入了解这个多评论系统的核心插件。

<!-- more -->

# 🧩 插件简介

[hexo-generator-comments](https://www.npmjs.com/package/hexo-generator-comments) 是一款 **Hexo 多评论系统生成插件**，支持多种评论系统的集成与切换，提供统一的评论界面。

在整个 **diversity-plugins** 插件家族中，它是名副其实的"核心担当"：

- 六款评论插件（Utterances / Gitalk / Giscus / Twikoo / Gitment / Waline）都**必须与它搭配使用**
- 它负责生成统一的评论页面（默认路径 `/comments/`），处理多系统的展示、切换与记忆

# ✨ 功能特性

| 特性 | 描述 |
|------|------|
| **多评论系统支持** | 同时集成多种评论系统（Utterances、Gitalk、Giscus、Twikoo、Gitment、Waline 等） |
| **多展示模式** | 支持选项卡（tabs）、下拉菜单（dropdown）及双模式切换（both），灵活适配不同场景 |
| **用户偏好记忆** | 智能记住访客选择的评论系统，提升用户体验 |
| **懒加载支持** | 可选懒加载机制，显著提高页面加载速度 |
| **主题无关性** | 完美兼容任何 Hexo 主题，无缝集成 |
| **自定义布局** | 灵活的布局和样式自定义选项 |
| **深色模式支持** | 内置深色模式样式，自动适应系统主题 |

# ⚙️ 基本配置

安装插件后，在 Hexo 站点根目录的 `_config.yml` 中添加以下配置：

```yaml
comments:
  title:
  layout:
  path:
  darkclass:
  style: tabs
  active:
  storage: true
  lazyload: false
  nav:
    #utterances:
    #  text: Utterances
    #  order: 0
    #gitalk:
    #  order: 1
```

**配置项说明**

* `title`：自定义评论标题【选填】，默认为 "Comments"
* `layout`：自定义评论布局文件名【选填】，不包含扩展名
* `path`：自定义评论页面路径【选填】，默认为 `comments`
* `darkclass`：深色主题类名【选填】。当主题通过在 `<html>`（documentElement）上添加该类名来标识深色模式时，插件会自动识别并同步评论区域的明暗状态
* `style`：多个评论系统启用时的默认展示风格，可选值：`tabs`【选项卡】 | `dropdown`【下拉菜单】 | `both`【双模式，支持用户一键切换】
* `active`：默认显示的评论系统，需与对应评论插件配置的系统名一致，如 `utterances` | `gitalk` | `giscus` | `twikoo` | `gitment` | `waline` 等
* `storage`：是否记住访客选择的评论系统，可选值：`true` | `false`
* `lazyload`：是否懒加载评论系统，可选值：`true` | `false`
* `nav`：调整导航元素的展示文本或顺序
  * `text`：导航元素的展示文本【选填】，默认为评论系统名
  * `order`：导航元素的展示顺序【数字越大，展示越靠后】

# 🗨️ 已支持的评论系统

以下是目前已支持的评论系统一览：

| 评论系统 | 特点 | 适用场景 |
|----------|------|----------|
| **Utterances** | 基于 GitHub Issues，轻量级 | 技术博客、开源项目 |
| **Gitalk** | 基于 GitHub Issues，功能丰富 | 个人博客、技术分享 |
| **Giscus** | 基于 GitHub Discussions，现代化 | 社区讨论、互动博客 |
| **Twikoo** | 免费私有部署，简洁安全 | 数据自主可控、个人博客 |
| **Gitment** | 基于 GitHub Issues，轻量级 | 个人博客、技术分享 |
| **Waline** | 轻量后端起评，功能丰富 | 个人博客、互动站点 |

每款评论系统都对应一个独立的插件包，**Huazie** 会在本系列后续文章中逐一带大家上手。

# 🔗 与评论插件的搭配

核心插件与评论插件的关系，可以概括为一个"安装公式"：

```bash
# 核心插件（必装）
npm install hexo-generator-comments --save

# 按需挑选一款或多款评论插件
npm install hexo-comments-utterances --save
npm install hexo-comments-giscus --save
```

装好后在 `_config.yml` 中将对应评论系统的 `enable` 设为 `true`，访问 `/comments/` 路径即可查看效果。同时启用多款评论系统时，插件会按 `style` 配置以选项卡或下拉菜单的形式呈现，访客可自由切换。

> 另外提一句：**Diversity** 主题已内置对该插件体系的集成，使用 Diversity 的朋友只需按主题配置说明开启即可。

# 📝 小结

本文介绍了 **diversity-plugins** 的核心插件 [hexo-generator-comments](https://www.npmjs.com/package/hexo-generator-comments)：

- **多系统**：一套配置同时集成六款评论系统
- **多模式**：选项卡 / 下拉菜单 / 双模式随心切换
- **零绑定**：不挑主题，任何 Hexo 站点都能用

如果你正打算给博客添加评论功能，或者想从单一评论系统升级为多系统切换，不妨从它开始。更多详细信息，请查阅插件的 [README 文档](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-generator-comments/README.md)。
