---
title: GitHub 账号直接评论！Hexo Comments Gitalk 插件上手指南
date: 2026-10-27 10:30:00
updated: 2026-10-27 10:30:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-comments-gitalk
  - Gitalk
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 一、前言

上一篇我们认识了基于 GitHub Issues 的轻量评论插件 [Utterances](../../../../../2026/10/27/hexo/hexo-comments-utterances/)。本篇登场的 **Gitalk** 同样基于 GitHub Issues，但功能更丰富 —— 它支持 **GitHub OAuth 安全登录**，访客体验更完整。

跟着 **Huazie** 一起来看看 [hexo-comments-gitalk](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-gitalk) 怎么用吧！

<!-- more -->

# 🔍 二、认识 Gitalk

[`Gitalk`](https://gitalk.github.io) 是一个基于 **GitHub Issue** 和 **Preact** 开发的现代评论插件。它允许网站访客使用 **GitHub** 账号登录并发表评论，所有评论数据都储存在相应的 **GitHub** 仓库中。

与 Utterances 一样，它无需自建数据库；不同的是，Gitalk 走的是完整的 **OAuth 认证** 流程，登录、评论、查看都在评论区内完成。

# ✨ 三、插件功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **GitHub 集成** | 基于 GitHub Issues，无需数据库 | 零维护成本，高可用性 |
| **OAuth 认证** | 支持 GitHub OAuth 安全登录 | 保护用户隐私，安全可靠 |
| **主题切换** | 支持亮色/暗色主题自动切换 | 完美适配各种主题风格 |
| **响应式设计** | 适配各种设备屏幕 | 移动端友好的用户体验 |
| **多语言支持** | 支持多种界面语言 | 国际化友好 |
| **易于配置** | 简单的 YAML 配置 | 快速上手，灵活定制 |

# 🚀 四、快速上手

## 1. 申请 GitHub OAuth Application

Gitalk 需要 GitHub OAuth 应用授权。前往 **GitHub → Settings → Developer settings → OAuth Apps** 新建一个 OAuth App，`Homepage URL` 与 `Authorization callback URL` 填写你的博客地址，创建后即可获得 `Client ID` 和 `Client Secret`。

## 2. 安装插件

在 Hexo 项目根目录执行：

```bash
# 1. 安装多评论系统核心插件（必需）
npm install hexo-generator-comments --save

# 2. 安装 Gitalk 评论插件
npm install hexo-comments-gitalk --save
```

> **提示**：`hexo-comments-gitalk` 需要与 `hexo-generator-comments` 搭配使用。

## 3. 添加配置

在 Hexo 站点配置 `_config.yml` 或主题配置 `_config.yml`、`_config.[theme].yml` 中添加以下内容：

```yaml
gitalk:
  # 是否启用 Gitalk 评论系统
  enable: false
  # 是否启用加载提示（评论加载时显示加载动画）
  loading: true
  # GitHub Application Client ID
  client_id: your-client-id
  # GitHub Application Client Secret
  client_secret: your-client-secret
  # GitHub 仓库名称
  repo: your-repo-name
  # GitHub 仓库所有者
  github_id: your-github-id
  # GitHub 仓库管理员
  admin_user: your-github-id
  # 页面唯一标识
  issue_term: pathname
  # 语言设置
  language:
  # 是否启用无干扰模式
  distraction_free_mode: true
  # 代理 URL（可选）
  proxy: https://cors-anywhere.azm.workers.dev/https://github.com/login/oauth/access_token
```

> **重要**：请将配置中的占位符替换为你的实际 GitHub 应用信息。

## 4. 关键配置项说明

| 选项 | 描述 |
|------|------|
| `client_id` / `client_secret` | 第 1 步申请的 GitHub Application 凭证 |
| `repo` | 用于存储评论的 GitHub 仓库名称 |
| `github_id` | GitHub 仓库所有者 |
| `admin_user` | GitHub 仓库管理员（一般与仓库所有者一致） |
| `issue_term` | 页面唯一标识，决定文章与 Issue 的对应关系 |
| `language` | 界面语言，如 `zh-CN`、`en` |
| `distraction_free_mode` | 无干扰模式开关 |
| `proxy` | OAuth 访问令牌的代理 URL，可缓解跨域问题 |

完整的配置选项详解，请查阅插件的 [README 文档](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-gitalk/README.md)。

# ⚠️ 五、注意事项

* **凭证保密**：`client_secret` 属于敏感信息，注意不要把含真实凭证的配置文件提交到公开仓库。
* **评论初始化**：每个页面的评论对应仓库中的一个 Issue，首次访问时会自动创建。
* **管理员配置**：`admin_user` 需填写仓库管理员的 GitHub ID，便于对评论区 Issue 进行管理。
* **跨域代理**：部分网络环境下获取访问令牌可能受限，可通过 `proxy` 配置自定义代理地址。

# 📝 六、小结

本文带大家上手了 **diversity-plugins** 中的 [hexo-comments-gitalk](https://www.npmjs.com/package/hexo-comments-gitalk) 插件。相比 Utterances，Gitalk 提供了更完整的 OAuth 登录与评论管理体验，适合希望**评论区功能更丰富**的博客。

如果你觉得维护 OAuth 应用有点重，下一篇的 **Giscus**（基于 GitHub Discussions）或许会让你眼前一亮，敬请期待！
