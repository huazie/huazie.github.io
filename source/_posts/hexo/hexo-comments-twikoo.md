---
title: 数据自主可控的免费评论方案！Hexo Comments Twikoo 插件上手指南
date: 2026-10-27 10:50:00
updated: 2026-10-27 10:50:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-comments-twikoo
  - Twikoo
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 一、前言

前面几篇介绍的 Utterances、Gitalk、Giscus 都深度绑定 GitHub 生态，访客评论需要 GitHub 账号。如果你的博客受众更"大众"一些，或者希望**评论数据完全自主可控**，那么本篇的 **Twikoo** 会是一个理想选择。

跟着 **Huazie** 一起上手 [hexo-comments-twikoo](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-comments-twikoo) 吧！

<!-- more -->

# 🔍 二、认识 Twikoo

[`Twikoo`](https://twikoo.js.org) 是一个**简洁、安全、免费**的静态网站评论系统。

与 GitHub 系评论不同，Twikoo 走的是"**免费私有部署**"路线：你需要部署一个轻量的后端环境，评论数据就存在这个环境里 —— 数据完全归你掌控，不依赖任何第三方账号体系。

它支持 **腾讯云 CloudBase、Vercel、Netlify** 等多种部署方式，其中 Vercel / Netlify 均有免费额度，个人博客零成本即可运行。

# ✨ 三、插件功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **免费私有部署** | 支持腾讯云 CloudBase、Vercel、Netlify 等多种部署方式 | 数据自主可控，零成本运行 |
| **安全可靠** | 完全开源，无广告无追踪 | 保护用户隐私，透明可信 |
| **暗黑模式** | 自动适配系统颜色方案 | 完美融入各种主题风格 |
| **响应式设计** | 适配各种设备屏幕 | 移动端友好的用户体验 |
| **即时加载** | 支持懒加载和加载动画 | 优化页面性能 |
| **易于配置** | 简单的 YAML 配置 | 快速上手，灵活定制 |
| **多语言支持** | 内置多语言界面 | 国际化用户体验 |

# 🚀 四、快速上手

## 1. 部署 Twikoo 后端环境

先按 [Twikoo 官方文档](https://twikoo.js.org) 部署一个后端环境（腾讯云 CloudBase / Vercel / Netlify 任选其一），部署完成后会得到一个**环境 ID（envId）**，这是接入的关键凭证。

## 2. 安装插件

在 Hexo 项目根目录执行：

```bash
# 1. 安装多评论系统核心插件（必需）
npm install hexo-generator-comments --save

# 2. 安装 Twikoo 评论插件
npm install hexo-comments-twikoo --save
```

> **提示**：`hexo-comments-twikoo` 需要与 `hexo-generator-comments` 搭配使用。

## 3. 添加配置

在 Hexo 站点配置 `_config.yml` 或主题配置 `_config.yml`、`_config.[theme].yml` 中添加以下内容：

```yaml
twikoo:
  # 是否启用 Twikoo 评论系统，可选值：true 【启用】 | false 【禁用】
  enable: false
  # 是否启用加载提示，可选值：true | false
  loading: true
  # Twikoo 环境 ID (envId)，必填项
  env_id: your-env-id
  # 环境地域，默认为 ap-shanghai，腾讯云环境填 ap-shanghai 或 ap-guangzhou；Vercel 环境不填
  region:
  # 页面路径，用于区分不同页面的评论，默认使用 window.location.pathname
  path:
  # 评论区语言
  lang: zh-CN
  # Twikoo JS SDK CDN 地址
  js: https://cdn.jsdelivr.net/npm/twikoo@1.7.9/dist/twikoo.min.js
```

> **重要**：请将 `your-env-id` 替换为你实际的 Twikoo 环境 ID。

## 4. 配置项说明

| 选项 | 类型 | 默认值 | 必填 | 描述 |
|------|------|--------|------|------|
| `enable` | Boolean | `false` | 是 | 是否启用 Twikoo 评论系统 |
| `loading` | Boolean | `true` | 否 | 是否启用加载提示（评论加载时显示加载动画） |
| `env_id` | String | - | 是 | Twikoo 环境 ID，在部署后台系统中获取 |
| `region` | String | `ap-shanghai` | 否 | 环境地域，腾讯云环境填 `ap-shanghai` 或 `ap-guangzhou`；Vercel 环境不填 |
| `path` | String | - | 否 | 页面路径，用于区分不同页面的评论，默认使用 `window.location.pathname` |
| `lang` | String | `zh-CN` | 否 | 评论区语言 |
| `js` | String | CDN 地址 | 否 | Twikoo JS SDK 的 CDN 地址，可指定特定版本或自部署地址 |

# ⚠️ 五、注意事项

* **先部署后接入**：`env_id` 来自已部署的 Twikoo 后端环境，请确保环境可用后再启用评论。
* **region 别乱填**：仅腾讯云环境需要填写 `ap-shanghai` 或 `ap-guangzhou`；使用 Vercel 部署时**不要填**。
* **SDK 版本可控**：`js` 配置可固定 Twikoo SDK 版本，也可指向自部署的 SDK 地址，避免 CDN 波动影响评论区。

# 📝 六、小结

本文带大家上手了 **diversity-plugins** 中的 [hexo-comments-twikoo](https://www.npmjs.com/package/hexo-comments-twikoo) 插件。它的核心卖点是：**免费私有部署 + 数据自主可控**，访客无需 GitHub 账号也能直接留言，适合面向更广泛读者的博客。

下一篇我们再回到 GitHub Issues 系，聊聊其中"资历最老"的 **Gitment**，敬请期待！
