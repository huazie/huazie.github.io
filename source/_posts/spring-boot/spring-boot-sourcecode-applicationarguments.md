---
title: 【Spring Boot 源码学习】深入 ApplicationArguments 接口及其默认实现
date: 2024-05-10 13:29:42
updated: 2024-05-10 13:29:42
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言

在 [《SpringApplication 的 run 方法核心流程介绍》](/2024/04/28/spring-boot/spring-boot-sourcecode-springapplication-run-listener/) 博文中，我们知道了 `ApplicationArguments` 是 **Spring Boot** 中用于获取 **应用程序启动参数** 的接口，其默认实现是 `DefaultApplicationArguments`。

不过有关内容尚未详细介绍，本篇就带大家深入分析下 `ApplicationArguments` 接口及其默认实现。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

## 2.1 ApplicationArguments

首先来看 **应用程序启动参数接口类** `ApplicationArguments` 的源码：

```java
public interface ApplicationArguments {
  String[] getSourceArgs();
  Set<String> getOptionNames();
  boolean containsOption(String name);
  List<String> getOptionValues(String name);
  List<String> getNonOptionArgs();
}
```

`ApplicationArguments` 接口共包含 `5` 个方法，均用于运行 `SpringApplication` 的参数的访问：

- `getSourceArgs`：该方法返回传递给应用程序的原始未处理参数。
- `getOptionNames`：该方法返回所有选项参数的名称，如果没有则返回一个空集。例如，如果参数是 `"--foo=bar --debug"`，则返回值应为 `["foo", "debug"]`。
- `containsOption`：该方法返回从参数中解析出的选项参数集合中是否包含给定名称的选项。如果参数中包含给定名称的选项，则返回 `true`。
- `getOptionValues`：该方法返回与给定名称的参数选项关联的值集合。
  - 如果选项存在且没有参数（例如：`"--foo"`），则返回空集合（`[]`）。
  - 如果选项存在且有一个值（例如：`"--foo=bar"`），则返回包含一个元素的集合（`["bar"]`）。
  - 如果选项存在且有多个值（例如：`"--foo=bar --foo=baz"`），则返回一个包含每个值的元素的集合（`["bar", "baz"]`）。
  - 如果选项不存在，则返回 `null`。
- `getNonOptionArgs`：该方法返回解析出的非选项参数的集合，如果没有则返回一个空列表。

## 2.2 DefaultApplicationArguments
在 `SpringApplication` 的 `run` 方法中，我们可以看到如下标红的内容：

![](run.png)

`DefaultApplicationArguments` 就是 `ApplicationArguments` 接口的一个默认实现，还是来看看相关源码：
![](DefaultApplicationArguments.png)
### 2.2.1 成员变量
`DefaultApplicationArguments` 的成员变量有两个：

- `Source source` ：私有的静态内部类，继承自 `org.springframework.core.env.SimpleCommandLinePropertySource`【**spring-core** 包中的一个类，旨在提供解析命令行参数的最简单方法】
- `String[] args`：原始的命令行参数数组

![](Source.png)

这个 `Source` 类是对 `SimpleCommandLinePropertySource` 的一个简单封装，可以看到它这里只是简单地调用了父类的实现，没有添加任何新的功能或逻辑。当然就目前而言，这个类似乎是多余的，**Huazie** 猜测这也许是为了将来的扩展吧。

### 2.2.2 构造方法

```java
public DefaultApplicationArguments(String... args) {
  Assert.notNull(args, "Args must not be null");
  this.source = new Source(args);
  this.args = args;
}
```

构造方法主要用来初始化上述两个成员变量，它接受一个可变长度的字符串数组作为参数，该参数就是运行 `SpringApplication` 的命令行参数。
### 2.2.3 成员方法
![](method.png)

阅读上述源码，可以看到实现的 **5** 个方法中都跟成员变量 `source` 有关系。

在 `getOptionNames` 方法中，可以看到返回的 `Set` 集合是通过 `Collections.unmodifiableSet` 包装过的 `Set` 集合的不可修改视图。`Collections.unmodifiableSet` 允许模块向用户提供对内部集合的“只读”访问权限。对返回集合的查询操作将“穿透”到指定的集合，而尝试修改返回的集合（无论是直接修改还是通过其迭代器）都将导致`UnsupportedOperationException` 异常。如果指定的集合是可序列化的，那么返回的集合也将是可序列化的。

同理，在 `getOptionValues` 方法中的 `Collections.unmodifiableList` 方法返回的是 `List` 集合的不可修改视图。
## 2.3 SimpleCommandLinePropertySource
上述 2.2.3 中，我们可以看到最终成员方法的处理都是来自 `SimpleCommandLinePropertySource` 类中的实现方法。

我们来看看相关源码：

![](SimpleCommandLinePropertySource.png)

`SimpleCommandLinePropertySource` 继承自 `CommandLinePropertySource<CommandLineArgs>`；

`CommandLinePropertySource<T>` 是一个抽象基类，用于实现由命令行参数支持的`PropertySource`。泛型类型 `T` 代表命令行选项的底层来源。

在 `SimpleCommandLinePropertySource` 中，`T` 是 `CommandLineArgs`，它是命令行参数的简单表示，分为 **"选项参数"** 和 **"非选项参数"**。

![](CommandLineArgs.png)

继续看 `SimpleCommandLinePropertySource` 的构造方法，可以看到 命令行参数 `CommandLineArgs` 是 通过 `SimpleCommandLineArgsParser` 的 `parse(String... args)` 来解析的。

```java
new SimpleCommandLineArgsParser().parse(args);
```
![](SimpleCommandLineArgsParser.png)
## 2.4 应用场景

有关 `ApplicationArguments` 的应用场景，我们一步步跟着源码来看：
### 2.4.1 准备和配置应用环境
首先是 `run` 方法中，先创建一个 `DefaultApplicationArguments` 对象，并赋值给 `applicationArguments` 变量；

接着调用 `prepareEnvironment` 方法来准备和配置应用环境，并传入 `applicationArguments` 参数；

![](run-1.png)

进入 `prepareEnvironment` 方法，可以看到如下：

![](prepareEnvironment.png)

这里是获取了原始的参数数组，并作为入参，传入 `configureEnvironment` 方法中；

![](configureEnvironment.png)

进入 `configurePropertySources` 方法，它是用来配置应用程序的环境属性源，可以看到最终 args 参数都会被构建成 `SimpleCommandLinePropertySource` 的属性源。


![](configurePropertySources.png)

至于 `configureProfiles` 方法，目前是空实现，留着未来扩展。

### 2.4.2 准备和配置应用上下文
调用 `prepareContext` 方法准备和配置应用上下文，并传入 `applicationArguments` 参数；

![](run-2.png)

进入 `prepareContext` 方法，可以看到：

![](prepareContext.png)

如上标红处，将 `applicationArguments` 注册为单例的 **Bean** 实例，其名称为 `springApplicationArguments`，后续其他地方需要时，就可以通过该名称从 **Spring** 容器中获取 `applicationArguments`。

### 2.4.3 刷新应用上下文之后【afterRefresh方法】

`afterRefresh` 方法的实现默认为空，可由开发人员自行扩展。

### 2.4.4 调用运行器【callRunners方法】
`callRunners` 方法里，主要是调用 `ApplicationRunner` 和 `CommandLineRunner` 的运行方法。

![](callRunner.png)
![](ApplicationRunner.png)
![](CommandLineRunner.png)



# 三、总结

本篇博文 **Huazie** 同大家一起深入分析了 `ApplicationArguments` 接口及其默认实现，相信这些可以进一步加深大家对于 **Spring Boot** 启动运行阶段中命令参数获取和使用的理解。接下来的博文将会继续聚焦 **Spring Boot** 启动运行阶段，敬请期待！！！




