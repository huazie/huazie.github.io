---
title: 【Spring Boot 源码学习】初识 ConfigurableEnvironment
date: 2024-07-01 11:50:13
updated: 2024-07-01 14:10:16
categories:
- 开发框架-Spring Boot
tags:
- Spring Boot
- ConfigurableEnvironment
- Environment
- ConfigurablePropertyResolver
---

![](/images/spring-boot-logo.png)

# 一、引言

上篇博文，**Huazie** 带大家深入分析下 `ApplicationArguments` 接口及其默认实现。在初始化完 `ApplicationArguments` 之后，**Spring Boot** 就开始通过 `prepareEnvironment` 方法对 `ConfigurableEnvironment` 对象进行初始化操作。在介绍 `ConfigurableEnvironment` 的初始化之前，我们有必要先认识一下 `ConfigurableEnvironment` 接口。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)


# 二、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

下面贴出 `ConfigurableEnvironment` 的源码：

```java
public interface ConfigurableEnvironment extends Environment, ConfigurablePropertyResolver {
    void setActiveProfiles(String... profiles);
    void addActiveProfile(String profile);
    void setDefaultProfiles(String... profiles);
    MutablePropertySources getPropertySources();
    Map<String, Object> getSystemProperties();
    Map<String, Object> getSystemEnvironment();
    void merge(ConfigurableEnvironment parent);
}
```

从上述源码，可以看出 `ConfigurableEnvironment` 接口继承了 `Environment` 和 `ConfigurablePropertyResolver` 接口，它们其实最终都继承自 `PropertyResolver` 接口。

## 2.1 Environment
![](Environment.png)

`org.springframework.core.env.Environment` 表示当前应用程序在其中运行的环境。它模拟了应用程序环境的两个关键方面：

### 2.1.1 配置文件（profiles）
`Profile` 是一个命名的、逻辑上的 **bean** 定义组，这些定义只有在给定的配置文件处于活动状态时才会被注册到容器中。通过 `Environment` 可以确定哪些配置文件（如果有）当前是活动的，以及哪些配置文件（如果有）应该默认是活动的。可以使用 `@Profile` 注解来指定 **bean** 应该在哪个配置文件下被注册。
### 2.1.2 属性（properties）
属性有各种来源，如**属性文件**、**JVM 系统属性**、**系统环境变量**、**JNDI**、**servlet 上下文参数**、**临时的 Properties 对象**、**Map** 等。`Environment` 对象为用户提供了一个方便的服务接口，用于配置属性源并从这些源中解析属性。通过 `Environment`，可以方便地访问和解析这些属性，而无需直接操作这些源。

此外，`Environment` 接口还继承了 `PropertyResolver` 接口【用于根据任何底层源解析属性的接口】，这意味着它还提供了与属性访问相关的功能。

## 2.2 ConfigurablePropertyResolver
![](ConfigurablePropertyResolver.png)

`org.springframework.core.env.ConfigurablePropertyResolver` 也继承了 `PropertyResolver` 接口，并在其基础上增加了更多的配置选项：

### 2.2.1 属性类型转换配置
`ConfigurablePropertyResolver` 提供了基于 `org.springframework.core.convert.ConversionService` 的属性类型转换功能。`ConversionService` 是 **Spring** 中用于类型转换的接口，它允许将一种类型的对象转换为另一种类型的对象。

与之关联的方法如下：

- `getConversionService()` : 获取当前用于类型转换的 `ConfigurableConversionService` 实例。
- `setConversionService(ConfigurableConversionService conversionService)`: 设置用于类型转换的 `ConfigurableConversionService` 实例。这允许用户自定义类型转换的逻辑，以满足特定的应用程序需求。

### 2.2.2 占位符配置
`ConfigurablePropertyResolver` 允许开发者配置占位符的前缀和后缀。默认情况下，前缀是 `${`，后缀是 `}`。占位符的值本身也可以包含其他占位符，形成嵌套占位符，`ConfigurablePropertyResolver` 支持嵌套占位符的解析。

与之相关的方法如下：

- `setPlaceholderPrefix(String placeholderPrefix)` : 设置占位符的前缀。在解析属性时，这些前缀将被用来识别需要替换的占位符。
- `setPlaceholderSuffix(String placeholderSuffix)` : 设置占位符的后缀。与前缀一起，它们定义了占位符的完整格式。 
- `setIgnoreUnresolvableNestedPlaceholders(boolean ignoreUnresolvableNestedPlaceholders)` : 设置是否忽略无法解析的嵌套占位符。如果设置为 `true`，则当遇到无法解析的嵌套占位符时，解析器将不会抛出异常，而是继续执行。

### 2.2.3 值分隔符配置

**值分隔符**是指在解析属性值时，用于分隔占位符与其关联默认值的字符设置。

比如，在配置文件中有这样的属性值：`${propertyName:defaultValue}`。

在这里 `propertyName` 是占位符，而 `defaultValue` 是在 `propertyName` 无法解析时使用的默认值。那显然在上述示例中，`:` 就是 **值分隔符**。

与之相关的方法如下：
- `setValueSeparator(@Nullable String valueSeparator)` : 设置值分隔符。在某些情况下，属性值可能包含多个值，这些值由分隔符分隔。此方法允许用户指定分隔符。

### 2.2.4 必需属性验证配置
必需属性验证配置是 **Spring** 框架中用于确保应用程序配置中包含某些关键属性的一种机制。

与之相关的方法如下：

- `setRequiredProperties(String... requiredProperties)`: 设置必需的属性。这些属性必须在解析过程中存在，否则验证将失败。
- `validateRequiredProperties() throws MissingRequiredPropertiesException` : 验证是否所有必需的属性都已设置。如果任何必需属性缺失，此方法将抛出 `MissingRequiredPropertiesException` 异常。

## 2.3 ConfigurableEnvironment
了解了 `Environment` 和 `ConfigurablePropertyResolver`，我们再来看看 `ConfigurableEnvironment` 。

### 2.3.1 接口方法
`ConfigurableEnvironment` 代表了一个可配置的环境，其定义了如下的方法：

- `setActiveProfiles(String... profiles)` ：设置当前激活的 `Profile` 组集合。在**Spring** 中，`Profile` 允许用户根据特定的环境（如开发、测试、生产）加载不同的配置。通过传递一个或多个 `Profile` 名称作为参数，你可以激活这些 `Profile`。
- `addActiveProfile(String profile)` ：向当前激活的 `Profile` 组集合中添加一个 `Profile` 组。
- `setDefaultProfiles(String... profiles)` ：设置默认激活的 `Profile` 组集合。激活的 `Profile` 组集合为空时，会默认实用默认的 `Profile` 组集合。
- `getPropertySources()` ：返回当前环境的 `MutablePropertySources` 对象。`PropertySources` 是一个包含多个 `PropertySource` 的列表，每个`PropertySource` 都可以提供属性。`MutablePropertySources` 允许你添加、替换或删除 `PropertySource`。
- `getSystemProperties()` ：返回 **Java** 系统属性的映射。这些属性是 **JVM** 启动时通过 `-D` 参数或在代码中使用 `System.setProperty(key, value)` 进行设置。
- `getSystemEnvironment()` ：返回操作系统环境变量的映射。这些变量通常包含关于系统配置和运行时的信息。
- `merge(ConfigurableEnvironment parent)` ：将父 `ConfigurableEnvironment` 的属性源合并到当前环境中。合并时，父环境的属性源将添加到当前环境的属性源列表的开头，从而允许它们覆盖当前环境的任何同名属性。

### 2.3.2 具体实现
`org.springframework.core.env.AbstractEnvironment` 是一个抽象类，实现了 `ConfigurableEnvironment` 接口，为环境配置（如属性源和 Profile 文件管理）提供了基本的支持。

`org.springframework.core.env.StandardEnvironment` 继承自 `AbstractEnvironment`，应用于**非 Web 环境**。它是 **Spring** 中默认的环境配置类，负责读取系统属性、环境变量以及配置文件中的配置信息，并将其封装在一个 `PropertySources` 对象中供 **Spring** 应用程序使用。

`org.springframework.web.context.support.StandardServletEnvironment` 继承自 `StandardEnvironment`，它是基于 **Servlet** 的 **Web** 应用程序要使用的 **Environment** 实现。所有基于 **Servlet** 的 **Web** 相关的 `ApplicationContext` 类都会默认初始化一个实例。提供 `ServletConfig`、`ServletContext` 和基于 **JNDI** 的 `PropertySource` 实例。在初始化过程中，会根据 `ServletContext` 和 `ServletConfig` 的可用性来初始化和配置属性源。通过 `customizePropertySources()` 方法，可以自定义属性源的添加顺序和配置方式。

`org.springframework.mock.env.MockEnvironment` 继承自 `AbstractEnvironment`，它用于测试目的，可以模拟环境变量和系统属性的值。
# 三、总结

本篇博文 **Huazie** 同大家一起了解了 `ConfigurableEnvironment` 接口和其父接口，这些对于后续理解 `ConfigurableEnvironment`  的初始化操作至关重要。接下来的博文将会继续聚焦 **Spring Boot** 启动运行阶段，敬请期待！！！



