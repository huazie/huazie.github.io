---
title: 从单一到多元：揭秘 Hexo Diversity 主题的运行原理
date: 2024-11-01 21:03:01
updated: 2024-11-01 21:03:01
categories:
  - 博客框架-Hexo
tags:
  - Hexo
  - Diversity
  - hexo-theme-diversity
comments: true
---

![](/images/hexo-logo.png)

# 一、 引言

众所周知，**Hexo** 是一个快速、简洁且高效的博客框架。 它使用 **Markdown** 或 **其他标记语言** 解析文章，在几秒内，即可利用靓丽的主题生成静态网页。

目前 **Hexo** 拥有 **400+** 的主题，使用者可以在众多的主题中进行选择；当然我们也可以使用任何兼容的模板引擎创建自己的主题，从而应用到个人博客中。

<!-- more -->

不过比较遗憾的是，这么多的美观、强大、可定制的主题之中，我们只能选择其中一个主题来生成自己的静态博客。

笔者在接触 **Hexo** 之后，也被它众多的主题所吸引，同时也在考虑如何让个人博客可以在众多的主题间切换呢？【小孩子才会选择，大人全都要】。

# 二、 Diversity 主题

随着接触越多的 **Hexo** 的核心 **API** 和源码，一个在多主题中自由切换的 **Hexo 主题** 应运而生，取名 [Diversity](https://github.com/huazie/hexo-theme-diversity)，即多样性。

## 2.1 Hexo 控制台命令

在介绍 **Diversity** 主题的运行原理之前，我们先来看看 **Hexo** 的一些关键控制台命令，比如：

- `hexo generate` ： 生成静态文件（一般在`public`目录下）。
- `hexo server` ：启动服务器。 默认情况下，访问网址为： `http://localhost:4000/`。
- `hexo clean` ：清除缓存文件(`db.json`) 和已生成的静态文件 (一般在`public`目录下)。


## 2.2 Hexo 核心 API

翻看官方网站的 [核心API概述](https://hexo.io/zh-cn/api/) ，可以看到一个完整的 Hexo 实例运行流程如下：

- **初始化** ：首先，需要新建一个 `Hexo` 实例。 一个新的实例需要两个参数：网站根目录 `base_dir`，以及包含初始化选项的对象。 接着执行 `init` 方法后，**Hexo** 会加载插件及配置文件。
```js
var Hexo = require("hexo");
var hexo = new Hexo(process.cwd(), {});

hexo.init().then(function () {
    // ...
});
```

- **加载文件** ：**Hexo** 提供了两种方法来加载文件：`load` 和 `watch`。 `load` 用于加载 `source` 文件夹中的所有文件以及主题数据。 `watch` 执行与 `load` 相同的操作，但还会开始连续监视文件更改。

- **执行指令** ：任何控制台命令都可以通过在 **Hexo** 实例上使用 `call` 方法来显式调用。这样的调用需要两个参数：**控制台命令的名称**和 **一个选项参数**。不同的控制台命令有不同的选项可用。
```javascript
hexo.call("generate", {}).then(function () {
    // ...
});
```

- **退出** ：无论控制台命令完成与否，都应调用 `exit` 方法。 这样 **Hexo** 就能优雅地退出，并完成保存数据库等重要工作。
```js
hexo
  .call("generate", args)
  .then(() => hexo1.exit())
  .catch(err => hexo1.exit(err));
```

## 2.3 运行原理

通过上述 **Hexo** 核心 **API** 的了解，相信大家已经比较清楚 **Hexo** 实例的运行流程了。

**那 Diversity 主题是如何实现多个主题自由切换的呢？**

笔者的想法就是在我们运行**Hexo**的控制台命令时，可以针对**Diversity**主题中配置的**多主题目录列表**的每个主题，都创建一个 **Hexo** 实例并执行当前指令。

如下是 `_config.diversity.yml` 中的配置：

```yml
# 多主题目录列表
themes: [phase,landscape,light]
# 多主题服务器端口列表
# 不配置，默认从4001开始
#ports: [5000,5001,5002]
```

**这个时候读者可能就要问了，那不同的主题是如何区分处理的呢？**

### 2.3.1 多主题配置相关

不同主题的根配置 `_config.xml` 需要单独配置，在你的 **Hexo** 项目根目录，添加 **config** 目录，为上述多主题列表中的每个主题添加一个对应主题名的配置目录，并在该配置目录下添加对应的 `_config.yml` 【直接从你原来项目根目录下的 `_config.yml` 复制一份即可】，形如：

```txt
├─config
│  ├─landscape
│  │  ├─_config.yml
│  ├─light
│  │  ├─_config.yml
│  ├─phase
│  │  ├─_config.yml
```

修改上述各主题配置目录下的 `_config.yml`，以 **landscape** 举例：

```diff
- url: http://example.com
+ url: http://example.com/landscape

- public_dir: public
+ public_dir: public/landscape

- theme: other-theme
+ theme: landscape
```

只要按上述添加好配置，

当运行 `generate` 指令时，**Diversity** 主题可以针对不同的主题生成不同的静态页面【例如：`landscape` 的静态页面就将生成在 `public/landscape`】。

当运行 `server` 指令时，**Diversity** 主题可以针对不同的主题启动不同的Http服务【例如：`landscape` 的本地地址 `http://localhost:4002/landscape`】。

另外在你的 **Hexo** 项目根目录下，我们依旧可以添加不同主题独立的 `_config.[theme].yml` 文件，更多了解请查看官方[《配置》](https://hexo.io/zh-cn/docs/configuration)

![](config-yml.png)


### 2.3.2 多主题执行指令

在 **Diversity** 主题的 `script` 目录新增 `index.js`，用于处理**多主题执行指令**。

> 主题目录中的 `script` 目录为脚本文件夹。 在启动时，**Hexo** 会加载此文件夹内的 **JavaScript** 文件。 请参见 [plugins](https://hexo.io/zh-cn/docs/plugins). 以获得更多信息。

我们来看看核心的处理逻辑：

- 首先，从**Hexo**的主题配置中获取多主题目录列表；
```js
const themeConfig = hexo.config.theme_config;
const themes = themeConfig.themes;
```

- 然后，获取hexo执行的控制台命令；
```js
// 获取控制台命令的别名
const { alias } = hexo.extend.console;
// 获取hexo执行命令
const cmd = alias[hexo.env.cmd];
// 当前项目根目录
const cwd = process.cwd();
hexo.log.info('Cmd =', cmd);

themeConfig.cmd = cmd;
```

- 最后，遍历多主题目录列表 `themes`，针对每个主题创建一个 **Hexo** 实例并执行当前的控制台命令。每个主题都有对应的 `_config.yml` ，可从上面的介绍了解到，位于项目根目录的 `config` 目录中的主题名目录下。
```js
// 多主题目录配置的数组索引
let index = 0;
// 循环处理配置的多主题列表
themes.forEach(function(theme) {
    themeConfig.index = index;
    if (Util.isMatchCmd(cmd)) {
        hexo.log.info('Theme', (index + 1), '=', theme);
        const {args} = hexo.env;
        const fileName = '_config.yml';
        args.output = path_1.join(cwd, 'config', theme);
        if (!Util.isExist(args.output, fileName)) {
            hexo.log.error('Please add the [' + fileName + '] file in [' + args.output + '].');
            return;
        }
        args.config = path_1.join(args.output, fileName);
        const hexo1 = new Hexo(cwd, args);
        require('./config')(hexo1, themeConfig);
        hexo1.init()
            .then(() => require('./helper')(hexo1, themeName))
            .then(() => require('./generator')(hexo1, themeName))
            .then(() => hexo1.call(cmd, args))
            .then(() => hexo1.exit())
            .catch(err => hexo1.exit(err));
    } 
    // 下一个主题
    index++;
});
```

## 2.4 版本演进

### 2.4.1 V1版本

主题选择页：

![](v1.png)


鼠标左右拖拽滚动展示：


![](v1-1.gif)


鼠标滚轮前后滚动展示：


![](v1.gif)


已配置的主题，鼠标悬停自动翻转展示详情：

![](v1-2.gif)


**主题直达**按钮，用于跳转对应主题博客页面：


![](v1-3.gif)

**主题来源**按钮，用于跳转对应的主题开源项目：


![](v1-4.gif)


### 2.4.2 V2版本

#### 2.4.2.1 PC 端

个人博客页【尚未设置过默认主题】：

![](v2-no-theme.png)


主题选择页：


![](v2-theme.png)


设置默认主题【当然`V1` 版本中的拖动和滚动展示同样支持，这里就不展示了】：

![](v2-1.gif)


取消默认主题，并重新设置新的默认主题：


![](v2-2.gif)


已有默认主题，直接重新设置默认主题：


![](v2-3.gif)


#### 2.4.2.2 Phone 端

个人博客页【尚未设置过默认主题】：

![](v2-phone-no-theme.png)


主题选择页：


![](v2-phone-theme.png)


**左右箭头切换**滚动展示，**左滑右滑切换**滚动展示：

![](v2-phone-1.gif)


设置默认主题：

![](v2-phone-2.gif)

## 2.5 后续展望

本篇博文截止，[Diversity](https://github.com/huazie/hexo-theme-diversity/releases/tag/v2.2.6) 主题的版本是 **2.2.6**。

后续该项目将会持续更新中，包括但不限于首页，友链，留言，关于等等。

# 三、 总结

**Diversity 主题** 解决了 **Hexo** 站点只能在线接入一个主题运行的痛点，使用者可以在配置的主题中自由切换展示。

如果你也有兴趣，不妨来试试接入 [Diversity](https://github.com/huazie/hexo-theme-diversity) 吧！！！

有任何的问题，欢迎来评论区和我讨论！





