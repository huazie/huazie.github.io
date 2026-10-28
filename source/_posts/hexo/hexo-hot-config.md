---
title: 改配置不想重启？Hexo Hot Config 插件让你的主题配置热更新
date: 2026-10-27 10:00:00
updated: 2026-10-27 10:00:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-plugins
  - hexo-hot-config
  - 热更新
---

![](/images/diversity-logo.png)

# 📖 前言

用 **Hexo** 搭建过博客的朋友，对下面这个场景一定不陌生：本地 `hexo server` 预览时，想微调一下主题配置，比如改个标题、调个文案，保存之后发现页面纹丝不动，只能 `Ctrl + C` 停掉服务器再重新启动，等待那熟悉的编译过程走完……

改一行，重启一次；再改一行，再重启一次。一来二去，调配置的兴致先被磨掉了一半。

为了解决这个问题，**Huazie** 在 [diversity-plugins](https://github.com/huazie/diversity-plugins) 项目中抽离出了 **hexo-hot-config** 插件，让主题配置支持**热更新**。接下来就跟着 **Huazie** 一起来看看吧！

<!-- more -->

# 🔄 插件简介

[hexo-hot-config](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-hot-config) 是一款 **Hexo 热配置插件**，支持主题配置热更新，核心能力就两条：

- 监听 `_config.{theme}.yml` 文件变化，自动重新生成页面
- 无需重启 Hexo 服务器即可实时预览主题配置更改

> 其中 `{theme}` 是你当前使用的主题名称，比如 **Diversity** 主题对应的配置文件就是 `_config.diversity.yml`。

# ✨ 功能特性

| 特性 | 描述 | 优势 |
|------|------|------|
| **配置热更新** | 监听主题配置文件变化，自动重新生成页面 | 保存即生效，所见即所得 |
| **零配置启用** | 安装后无需任何额外配置 | 装上就能用，无学习成本 |
| **实时预览** | 无需重启 Hexo 服务器 | 调配置不再来回折腾 |

# 🚀 快速上手

## 1. 安装插件

在你的 **Hexo** 项目根目录执行：

```bash
npm install hexo-hot-config --save
```

## 2. 启动预览

安装完成后**无需任何额外配置**，直接运行 Hexo 服务器命令即可自动启用：

```bash
hexo server
```

## 3. 体验热更新

启动后，随手修改项目根目录的 `_config.{theme}.yml` 文件并保存，插件会自动检测变更并重新生成页面，最新的配置效果立刻就能预览到，全程**无需手动重启服务器**。

# ⚙️ 工作原理

插件的处理流程并不复杂，主要分为四步：

1. 插件会在 **Hexo 服务器启动时**初始化配置监听器
2. 监听独立的主题配置文件 `_config.{theme}.yml` 的变更
3. 当检测到文件变更时，**重新加载配置**并合并到 Hexo 主题配置中
4. 触发页面**重新生成**，实时更新预览效果

# ⚠️ 注意事项

- 仅在运行 `hexo server` 命令时生效（它本来就是为本地预览服务的）
- 如果 `_config.{theme}.yml` 不存在于 Hexo 项目根目录，则该插件将不会执行监控
- 配置文件必须是有效的 **YAML** 格式，否则可能导致加载失败

# 📝 小结

本文介绍了 **diversity-plugins** 中的 [hexo-hot-config](https://www.npmjs.com/package/hexo-hot-config) 插件，一句话总结它的价值：**让"改配置 → 看效果"从分钟级变成秒级**。

如果你的博客也在使用独立的主题配置文件（`_config.{theme}.yml`），强烈建议装上它，本地调配置的体验会顺滑很多。

更多详细信息，请查阅插件的 [README 文档](https://github.com/huazie/diversity-plugins/tree/main/packages/hexo-hot-config/README.md)。
