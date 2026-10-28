---
title: 用 GitHub Discussions 搭建现代化评论区！Hexo Comments Giscus 插件实战
date: 2026-10-27 10:40:00
updated: 2026-10-27 10:40:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-comments-giscus
  - Giscus
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 前言

前面两篇介绍了基于 GitHub Issues 的 Utterances 与 Gitalk。本篇登场的 **Giscus** 则更进一步 —— 它基于 **GitHub Discussions（讨论区）** 实现，是这几款 GitHub 系评论系统中最"现代化"的一员。

跟着 **Huazie** 一起上手 [hexo-comments-giscus](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-giscus) 吧！

<!-- more -->

# 🔍 认识 Giscus

[`Giscus`](https://giscus.app) 是一个利用 **GitHub Discussions** 实现的评论系统。它所有的评论数据存储于 **GitHub** 公开仓库的 **Discussions**（讨论区）中，每个网页对应一个 **Discussion**，**评论即 Discussion 下的回复**。

与 Issues 相比，Discussions 更适合承载"讨论"：支持分类管理、楼中楼回复、表情反应（Reactions），社区氛围更浓。

# ✨ 插件功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **GitHub Discussions** | 基于 GitHub Discussions，无需数据库 | 零维护成本，高可用性 |
| **安全可靠** | 完全开源，无广告无追踪 | 保护用户隐私，透明可信 |
| **主题切换** | 支持亮色/暗色主题自动切换 | 完美适配各种主题风格 |
| **响应式设计** | 适配各种设备屏幕 | 移动端友好的用户体验 |
| **即时加载** | 支持懒加载和加载动画 | 优化页面性能 |
| **多语言支持** | 支持多种界面语言 | 国际化友好 |
| **易于配置** | 简单的 YAML 配置 | 快速上手，灵活定制 |

# 🚀 快速上手

## 1. 安装插件

在 Hexo 项目根目录执行：

```bash
# 1. 安装多评论系统核心插件（必需）
npm install hexo-generator-comments --save

# 2. 安装 Giscus 评论插件
npm install hexo-comments-giscus --save
```

> **提示**：`hexo-comments-giscus` 需要与 `hexo-generator-comments` 搭配使用。

## 2. 添加配置

在 Hexo 站点配置 `_config.yml` 或主题配置 `_config.yml`、`_config.[theme].yml` 中添加以下内容：

```yaml
giscus:
  # 是否启用 Giscus 评论系统
  enable: false
  # 是否启用加载提示，可选值：`true` | `false`
  loading: true
  # GitHub 仓库，格式为 username/repo-name
  repo: your-username/your-repo-name
  # GitHub 仓库的唯一 ID
  repo_id: your-repo-id
  # GitHub Discussions 分类名称
  category: General
  # Discussions 分类的唯一 ID
  category_id: your-category-id
  # 指定 GitHub discussion 的匹配规则
  mapping: pathname
  # 当 mapping 为 specific 时的特定术语
  term:
  # 是否启用严格的标题匹配
  strict: 0
  # 是否启用主帖子上的反应
  reactions_enabled: 1
  # 是否输出 discussion 的元数据
  emit_metadata: 0
  # 亮色主题
  theme: light
  # 暗色主题
  dark: dark
  # 评论区的语言
  lang:
  # 评论输入框的位置
  input_position: bottom
  # 懒加载评论
  data_loading: lazy
```

> **重要**：请将 `repo`、`repo_id`、`category_id` 等占位符替换为你的实际仓库信息。

## 3. 获取 repo_id 与 category_id

这是 Giscus 接入中最容易卡住的一步。推荐直接打开 [`giscus.app`](https://giscus.app/zh-CN)，按页面向导依次填入仓库名、选择 Discussion 分类，网站会**自动生成** `repo_id` 与 `category_id`，复制过来即可，无需手工从 GitHub API 挖取。

# ⚠️ 注意事项

* **仓库需公开**并安装 [giscus App](https://github.com/apps/giscus)，评论功能才能正常读写 Discussions。
* **mapping 推荐 `pathname`**：与 Utterances 的 `issue_term` 类似，一旦选定尽量别改，避免历史 Discussion 与页面对不上。
* **strict 严格匹配**：设为 `1` 后将启用更严格的标题匹配，可降低映射冲突的概率。
* **reacts_enabled**：控制主帖子上的表情反应开关，`1` 开启、`0` 关闭。
* **懒加载**：`data_loading: lazy` 可让评论区在滚动到可视区域时才加载，进一步提升页面性能。

# 📝 小结

本文带大家上手了 **diversity-plugins** 中的 [hexo-comments-giscus](https://www.npmjs.com/package/hexo-comments-giscus) 插件。如果你的博客读者以开发者为主，又希望评论区具备**分类、楼中楼、表情反应**等更现代的互动能力，Giscus 值得一试。

下一篇我们换个口味 —— 不依赖 GitHub、数据完全自主可控的 **Twikoo**，敬请期待！
