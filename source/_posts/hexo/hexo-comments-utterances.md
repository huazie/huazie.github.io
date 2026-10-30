---
title: 把 GitHub Issues 变成评论区！Hexo Comments Utterances 插件上手指南
date: 2026-10-27 10:20:00
updated: 2026-10-27 10:20:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-comments-utterances
  - Utterances
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 一、前言

上一篇 [《一套配置通吃六大评论系统！深入解读 Hexo Generator Comments 核心插件》](../../../../../2026/10/27/hexo/hexo-generator-comments/) 介绍了负责统一评论界面的核心插件。而真正让某一款评论系统"跑起来"的，是与其搭配的独立评论插件。

从本篇开始，**Huazie** 会逐一带大家上手这些评论插件。首先登场的，是其中最轻量的一员 —— [hexo-comments-utterances](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-utterances)。

<!-- more -->

# 🔍 二、认识 Utterances

[`Utterances`](https://utteranc.es) 是一个由开源社区构建的评论插件，它提供了一种在博客、文章或任何静态网站上添加互动式评论功能的高效解决方案。

它的核心思路非常巧妙：**所有的评论数据都存储于你指定的 GitHub 公开仓库的 Issues 中**，每个网页对应一个 Issue，**评论即 Issue 下的留言**。

这也让它天然具备了几大优势：

* **无需独立服务器支持**：数据托管在 GitHub，零部署成本
* **零维护成本**：不怕丢数据，GitHub 替你保管
* **界面纯净无广告**：开源透明，无追踪
* **访客使用 GitHub 账号即可留言**

# ✨ 三、插件功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **轻量级** | 基于 GitHub Issues，无需数据库 | 零维护成本，高可用性 |
| **安全可靠** | 完全开源，无广告无追踪 | 保护用户隐私，透明可信 |
| **主题切换** | 支持亮色/暗色主题自动切换 | 完美适配各种主题风格 |
| **响应式设计** | 适配各种设备屏幕 | 移动端友好的用户体验 |
| **即时加载** | 支持懒加载和加载动画 | 优化页面性能 |
| **易于配置** | 简单的 YAML 配置 | 快速上手，灵活定制 |

# 🚀 四、快速上手

## 1. 安装插件

在 Hexo 项目根目录执行（核心插件必须一起安装）：

```bash
# 1. 安装多评论系统核心插件（必需）
npm install hexo-generator-comments --save

# 2. 安装 Utterances 评论插件
npm install hexo-comments-utterances --save
```

> **提示**：`hexo-comments-utterances` 需要与 `hexo-generator-comments` 搭配使用。

## 2. 添加配置

在 Hexo 站点配置 `_config.yml` 或主题配置 `_config.yml`、`_config.[theme].yml` 中添加以下内容：

```yaml
utterances:
  # 是否启用 Utterances 评论系统
  enable: false
  # 是否启用加载提示，可选值：`true` | `false`
  loading: true
  # GitHub 仓库所有者和名称，格式为 username/repo-name
  repo: your-username/your-repo-name
  # 指定 issue 的匹配规则
  issue_term: pathname
  # 亮色主题
  theme: github-light
  # 暗色主题
  dark: github-dark
```

> **重要**：请将 `your-username/your-repo-name` 替换为你的实际 GitHub 仓库地址。

## 3. 配置项说明

| 选项 | 类型 | 默认值 | 必填 | 描述 |
|------|------|--------|------|------|
| `enable` | Boolean | `false` | 是 | 是否启用 Utterances 评论系统 |
| `loading` | Boolean | `true` | 否 | 是否启用加载提示（评论加载时显示加载动画） |
| `repo` | String | - | 是 | GitHub 仓库，格式：`username/repo-name` |
| `issue_term` | String | `pathname` | 是 | 映射评论到 issue 的方式 |
| `theme` | String | `github-light` | 是 | 亮色模式主题 |
| `dark` | String | `github-dark` | 是 | 暗色模式主题 |

其中 `issue_term` 决定了"网页 → Issue"的映射规则，可选值包括 `pathname` | `url` | `title` | `og:title` | `issue number` | `specific term`。一般**推荐使用默认的 `pathname`**，即用页面路径作为 issue 标题，适合绝大多数场景。完整映射规则说明请查阅插件的 [README 文档](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-utterances/README.md)。

# ⚠️ 五、注意事项

* **仓库需公开**：Utterances 的评论数据存放在 GitHub 仓库的 Issues 中，需绑定一个公开的 GitHub 仓库。
* **issue_term 一旦选定尽量别改**：修改映射方式会导致已生成的 Issue 与页面对不上，历史评论可能"找不到家"。
* **主题适配**：`theme` 与 `dark` 分别控制亮色、暗色模式下的评论样式，可选值如 `github-light`、`github-dark`、`preferred-color-scheme` 等。

# 📝 六、小结

本文带大家上手了 **diversity-plugins** 中的 [hexo-comments-utterances](https://www.npmjs.com/package/hexo-comments-utterances) 插件：**一个公开仓库 + 几行 YAML 配置，就能让你的博客拥有基于 GitHub Issues 的评论区**。

如果你的读者群体以开发者为主（人人都有 GitHub 账号），Utterances 是非常划算的选择。下一篇我们继续看同为 GitHub Issues 系、但功能更丰富的 Gitalk，敬请期待！
