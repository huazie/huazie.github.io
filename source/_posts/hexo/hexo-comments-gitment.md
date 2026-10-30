---
title: GitHub Issues 轻量评论方案！Hexo Comments Gitment 插件上手指南
date: 2026-10-27 11:00:00
updated: 2026-10-27 11:00:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-comments-gitment
  - Gitment
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 一、前言

GitHub Issues 系评论插件的第三位成员登场 —— **Gitment**。它比 Gitalk 更轻量，同时保留 OAuth 登录与 Markdown 支持。不过在接入前有一个**必须知道**的坑：它的默认 OAuth 代理已经失效，需要自行部署。本篇 **Huazie** 会把这个坑讲清楚。

一起来上手 [hexo-comments-gitment](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-gitment) 吧！

<!-- more -->

# 🔍 二、认识 Gitment

[`Gitment`](https://github.com/imsun/gitment) 是一个基于 **GitHub Issues** 的轻量级评论系统，访客通过 GitHub OAuth 登录后即可评论，评论数据存储在你指定仓库的 Issues 中。

它支持 **Markdown/GFM 语法**，代码高亮、格式丰富的评论都能完美呈现，非常契合技术博客的场景。

# ✨ 三、插件功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **GitHub 集成** | 基于 GitHub Issues，无需数据库 | 零维护成本，高可用性 |
| **OAuth 认证** | 支持 GitHub OAuth 安全登录 | 保护用户隐私，安全可靠 |
| **主题切换** | 支持亮色/暗色主题自动切换 | 完美适配各种主题风格 |
| **响应式设计** | 适配各种设备屏幕 | 移动端友好的用户体验 |
| **Markdown 支持** | 支持 Markdown/GFM 语法 | 代码高亮，格式丰富 |
| **易于配置** | 简单的 YAML 配置 | 快速上手，灵活定制 |

# 🚀 四、快速上手

## 1. 安装插件

在 Hexo 项目根目录执行：

```bash
# 1. 安装多评论系统核心插件（必需）
npm install hexo-generator-comments --save

# 2. 安装 Gitment 评论插件
npm install hexo-comments-gitment --save
```

> **提示**：`hexo-comments-gitment` 需要与 `hexo-generator-comments` 搭配使用。

## 2. 添加配置

在 Hexo 站点配置 `_config.yml` 或主题配置 `_config.yml`、`_config.[theme].yml` 中添加以下内容：

```yaml
gitment:
  # 是否启用 Gitment 评论系统
  enable: false
  # 是否启用加载提示（评论加载时显示加载动画）
  loading: true
  # GitHub 仓库所有者
  owner: your-github-id
  # 用于存储评论的 GitHub 仓库名
  repo: your-repo-name
  # GitHub Application Client ID
  client_id: your-client-id
  # GitHub Application Client Secret
  client_secret: your-client-secret
  # 指定 GitHub Issue 的匹配规则
  # 其中 pathname | url | title 用来匹配 issue 的标签，`自定义字符串` 是自定义其他页面唯一标识
  issue_term: pathname
  # 页面唯一标识（可选，默认使用 window.location.href）
  id:
  # 页面标题（可选，默认使用 document.title）
  title:
  # 页面链接（可选，默认使用 window.location.href）
  link:
  # 页面描述（可选）
  desc:
  # Issue 标签列表（可选）
  labels: []
  # 评论每页显示数量（可选，默认 20）
  per_page: 20
  # 评论最大高度限制（px），超过则折叠（可选，默认 250）
  max_comment_height: 250
  # OAuth 代理地址（可选）
  proxy:
  # 是否启用 Gitmint 增强（可选）
  gitmint: false
```

> **重要**：`client_id` 与 `client_secret` 来自你创建的 GitHub OAuth Application（申请方式同 Gitalk，见 [上一篇](../../../../../2026/10/27/hexo/hexo-comments-gitalk/)）。

## 3. 关键配置项说明

| 选项 | 描述 |
|------|------|
| `owner` / `repo` | GitHub 仓库所有者与仓库名，评论数据就存在这里 |
| `client_id` / `client_secret` | GitHub OAuth Application 凭证 |
| `issue_term` | Issue 匹配规则，支持 `pathname` / `url` / `title` / 自定义字符串 |
| `per_page` | 评论每页显示数量，默认 20 |
| `max_comment_height` | 评论最大高度（px），超过则折叠，默认 250 |
| `proxy` | OAuth 代理地址，**重点见下方注意事项** |
| `gitmint` | 是否启用 Gitmint 增强 |

# ⚠️ 五、注意事项

* **默认 OAuth 代理已失效**：Gitment 项目默认代理 `gh-oauth.imsun.net` 已经失效，如遇登录异常，需通过 `proxy` 配置**自行部署**的 OAuth 代理服务。这是选择 Gitment 前必须评估的成本。
* **Gitmint 增强选项**：`gitmint: true` 会启用 [Gitmint](https://github.com/imsun/gitment) 的兼容替代品 —— 它修复了 Gitment 的多个问题；启用后 CSS/JS/CommentClass 将完全替换为 Gitmint，**与 Gitment 互斥**。
* **凭证保密**：`client_secret` 属于敏感信息，注意不要提交到公开仓库。
* **高度折叠**：`max_comment_height` 默认 250px，长评论会折叠展示，可按需调整。

# 📝 六、小结

本文带大家上手了 **diversity-plugins** 中的 [hexo-comments-gitment](https://www.npmjs.com/package/hexo-comments-gitment) 插件。Gitment 轻量、支持 Markdown，但**默认 OAuth 代理失效**这一点决定了它更适合愿意自行部署代理的用户。

如果不想折腾，同系的 Gitalk / Giscus 或下一篇的 **Waline** 都是更省心的选择，敬请期待！
