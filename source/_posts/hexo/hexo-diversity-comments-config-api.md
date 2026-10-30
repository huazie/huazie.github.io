---
title: 从默认值到页面标识解析！Diversity Comments SDK 配置与 API 详解
date: 2026-10-28 10:20:00
updated: 2026-10-28 10:20:00
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - diversity-comments-site
  - DiversityComments
  - 留言评论
  - 评论系统
---

![](/images/diversity-logo.png)

# 📖 一、前言

上一篇 [《一套 API 集成六大评论系统！Diversity Comments SDK 快速上手》](../../../../../2026/10/28/hexo/hexo-diversity-comments-quickstart/) 讲了三步接入。本篇深入两个日常开发离不开的主题：**全部配置项**（含 SDK 内部的页面标识解析规则）与 **API 方法**。

文中所有默认值均提取自 SDK 源码的 `DEFAULT_CONFIGS`（父页面与 iframe 共享同一套默认值，单一来源），无臆造参数。

<!-- more -->

# 🗂 二、配置总结构

```javascript
{
    server: '',                    // 评论服务地址（不传则取脚本加载 origin）
    container: '#diversity-comments',
    comments: { ... },             // 通用配置（对所有系统生效）
    utterances: { ... },           // 六大系统各自独立配置
    gitalk:     { ... },
    giscus:     { ... },
    twikoo:     { ... },
    gitment:    { ... },
    waline:     { ... },
    onReady: null,                 // 回调函数
    onError: null,
    onActiveChange: null
}
```

`init(options)` 内部使用 **deepMerge** 与默认配置合并：对象递归合并、数组与原始值直接覆盖——因此只需传入要改的字段。

# ⚙️ 三、comments 通用配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `pageId` | string | 不自动设置 | 页面唯一标识。传入则所有系统统一使用；未传时各插件按自身模式从父页面取值（见下文解析规则） |
| `style` | string | `'tabs'` | 导航模式：`tabs`（页签）/ `dropdown`（下拉）/ `both`（两者共存，右上角可切换） |
| `active` | string | `'utterances'` | 默认激活的评论系统 |
| `lazyload` | boolean | `true` | 进入视口才加载评论 |
| `storage` | boolean | `true` | 记住用户选择的评论系统（localStorage key: `selected_comment`） |
| `darkMode` | string | `'auto'` | `auto`（跟随系统 prefers-color-scheme）/ `light` / `dark` |
| `lang` | string | `''` | 统一覆盖各系统语言；为空则各插件使用自身默认。**Utterances 无语言配置，不适用** |
| `nav` | object | 见下 | 各系统在导航栏中的文案与排序 |

## nav 导航配置

每个系统可配置 `text`（显示文案）与 `order`（排序，默认 1~6 按 Utterances/Gitalk/Giscus/Twikoo/Gitment/Waline）：

```javascript
comments: {
    nav: {
        utterances: { text: 'Utterances', order: 1 },
        gitalk:    { text: 'Gitalk',    order: 2 },
        giscus:    { text: 'Giscus',    order: 3 },
        twikoo:    { text: 'Twikoo',    order: 4 },
        gitment:   { text: 'Gitment',   order: 5 },
        waline:    { text: 'Waline',    order: 6 }
    }
}
```

# 🔍 四、pageId 解析规则

SDK 收到父页面配置后，会按下面规则把「关键词模式」翻译成真实的页面标识。优先级：**pageId 显式传入 > 按模式从父页面取值**。

| 系统 | 处理字段 | 关键词 → 取值 |
|------|----------|---------------|
| Utterances | `issue_term` | `pathname` → 父页面 pathname；`url` → 完整 URL；`title` → 页面标题；`og:title` → og:title meta |
| Gitalk | `issue_term` | 同上 |
| Gitment | `issue_term` | 同上 |
| Waline | `path` | 同上 |
| Giscus | `mapping` | 关键词模式时改写为 `mapping: 'specific'` + `term: <取值>` |
| Twikoo | `path` | 为空时取 `pageId \|\| 父页面 pathname` |

两个细节：

- `pageId` 传入时，上述所有关键词模式统一替换为 `pageId`
- 用户已写死的自定义值（数字、具体字符串）**不会被覆盖**

## lang 统一应用规则

`comments.lang` 非空时，统一写入各启用系统的语言字段：

| 系统 | 覆盖字段 | 说明 |
|------|----------|------|
| Gitalk | `language` | 字段名不同 |
| Giscus / Twikoo / Waline / Gitment | `lang` | |
| Utterances | — | **无语言配置，不适用** |

# 🧱 五、六大系统配置默认值

以下为各系统的默认配置（各系统默认 `enable: true`，未显式配置的项走默认值）：

## 1. Utterances（基于 GitHub Issues，免费）

```javascript
utterances: {
    enable: true, loading: true,
    repo: 'user-name/repo-name',
    issue_term: 'pathname',
    theme: 'github-light',      // 亮色主题
    dark: 'github-dark'         // 暗色主题
}
```

## 2. Gitalk（基于 GitHub Issues，支持 OAuth）

```javascript
gitalk: {
    enable: true, loading: true,
    github_id: 'user-name', repo: 'repo-name',
    client_id: 'your-client-id', client_secret: 'your-client-secret',
    admin_user: '', distraction_free_mode: true,
    issue_term: 'pathname', language: 'zh-CN'
}
```

> OAuth 还支持 `proxy` / `redirect_uri`（不在默认配置中），详见 [运行原理与 OAuth 篇](../../../../../2026/10/28/hexo/hexo-diversity-comments-internals/)。

## 3. Giscus（基于 GitHub Discussions）

```javascript
giscus: {
    enable: true, loading: true,
    repo: 'your-username/your-repo-name',
    repo_id: 'R_xxx',
    category: 'Announcements', category_id: 'DIC_xxx',
    mapping: 'pathname', term: '', strict: 0,
    reactions_enabled: 1, emit_metadata: 0,
    theme: 'light', dark: 'dark',
    lang: 'zh-CN', input_position: 'bottom', data_loading: 'lazy'
}
```

## 4. Twikoo（无需申请，自建后端）

```javascript
twikoo: {
    enable: true, loading: true,
    env_id: 'xxx', path: '', lang: 'zh-CN',
    js: 'https://cdn.jsdelivr.net/npm/twikoo@1.7.9/dist/twikoo.min.js'
}
```

## 5. Gitment（基于 GitHub Issues，支持 OAuth）

```javascript
gitment: {
    enable: true, loading: true,
    owner: 'user-name', repo: 'repo-name',
    client_id: 'your-client-id', client_secret: 'your-client-secret',
    issue_term: 'pathname', id: '', title: '', link: '', desc: '',
    labels: [], per_page: 20, max_comment_height: 250,
    gitmint: true, lang: 'zh-CN'
}
```

## 6. Waline（功能丰富，自建后端）

```javascript
waline: {
    enable: true, loading: true,
    server_url: 'https://your-waline-server.netlify.app/.netlify/functions/comment',
    js_url: 'https://unpkg.com/@waline/client@v3/dist/waline.js',
    css_url: 'https://unpkg.com/@waline/client@v3/dist/waline.css',
    path: 'pathname', lang: 'zh-CN', emoji: null, dark: '.dark-theme',
    comment_sorting: 'latest', meta: ['nick', 'mail', 'link'], required_meta: [],
    login: 'enable', word_limit: false, page_size: 10,
    search: false, no_copyright: false, no_rss: false, reaction: false
}
```

> 各系统的 `enable` 全部为 `false` 时，评论区会显示「未启用评论」占位提示并隐藏整个导航。

# 🔌 六、API 参考

`DiversityComments` 全部公开成员：

| 方法 | 签名 | 说明 |
|------|------|------|
| `init` | `init(options)` | 初始化评论系统 |
| `switchTo` | `switchTo(name)` | 切换到指定评论系统 |
| `switchDisplayMode` | `switchDisplayMode(style)` | 切换导航模式 |
| `setColorScheme` | `setColorScheme(scheme)` | 切换明暗模式 |
| `refresh` | `refresh()` | 刷新当前评论 |
| `destroy` | `destroy()` | 销毁实例，移除 iframe 与事件监听 |
| `getConfig` | `getConfig() → Object` | 获取当前配置（深拷贝） |
| `getDefaultConfig` | `getDefaultConfig() → Object` | 获取 SDK 默认配置（深拷贝） |
| `version` | 属性 | 版本号，当前 `'1.0.0'` |

> **操作类方法**（`switchTo` / `switchDisplayMode` / `setColorScheme` / `refresh`）在初始化完成前调用只打警告不生效；`getConfig` / `getDefaultConfig` / `destroy` 不受此限制。

## 回调函数

| 回调 | 参数 | 触发时机 |
|------|------|----------|
| `onReady(iframe, active)` | iframe 元素、当前激活的系统名 | 配置注入完成、评论系统激活 |
| `onActiveChange(active)` | 新激活的系统名 | 任何来源的切换（API 或 iframe 内手动点击） |
| `onError(msg)` | 错误信息 | 容器缺失等；未配置时兜底 `console.error` |

## 方法要点

```javascript
// 切换评论系统（storage: true 时同步记忆）
DiversityComments.switchTo('gitalk');

// 切换导航模式：'tabs' | 'dropdown' | 'both'
DiversityComments.switchDisplayMode('both');

// 切换明暗模式：'light' | 'dark' | 'auto'
// 典型用法：跟随宿主站点主题联动
DiversityComments.setColorScheme(siteTheme === 'dark' ? 'dark' : 'light');

// 刷新当前评论（重新加载当前评论系统）
DiversityComments.refresh();

// 销毁实例（SPA 路由切换时销毁后对新容器重新 init）
DiversityComments.destroy();
```

# 🧪 七、典型场景

## 1. 页面加载后延迟初始化（等容器就绪）

```javascript
var timer = setInterval(function() {
    if (document.querySelector('#diversity-comments')) {
        clearInterval(timer);
        DiversityComments.init(config);
    }
}, 100);
```

## 2. SPA 路由变化：销毁重建

```javascript
DiversityComments.destroy();
config.comments.pageId = location.pathname;
DiversityComments.init(config);
```

## 3. 只读展示：禁用记忆、固定初始系统

```javascript
config.comments.storage = false;      // 不记忆用户选择
config.comments.active = 'giscus';    // 每次都从 giscus 开始
```

# 📚 八、小结

本篇是配置与 API 的完整速查。还想了解 SDK 内部如何运转（iframe 沙箱、消息协议、OAuth 登录），请移步 [运行原理与 OAuth 实战](../../../../../2026/10/28/hexo/hexo-diversity-comments-internals/)。
