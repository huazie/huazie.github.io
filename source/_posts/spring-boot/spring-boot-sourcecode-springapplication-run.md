---
title: 【Spring Boot 源码学习】SpringApplication 的 run 方法核心流程介绍
date: 2024-04-12 11:27:53
updated: 2024-04-12 11:27:53
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - SpringApplication
  - 运行核心流程
---



![](/images/spring-boot-logo.png)

# 一、引言

在前面的博文[《初识 SpringApplication》](../../../../../2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)中，**Huazie** 带大家一起分析了 `SpringApplication` 类实例化的逻辑。当 `SpringApplication` 对象被创建之后，我们就可以调用它的 `run` 方法来启动和运行 **Spring Boot** 项目。

本篇博文将围绕 `SpringApplication` 的  `run` 方法展开，带大家一起从源码分析 **Spring Boot** 的运行流程。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
    <tr>
        <td rowspan="30" align="left" > 
            <a href="/categories/开发框架-Spring-Boot/">Spring Boot 源码学习</a> 
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/02/19/spring-boot/spring-boot-project-introduction/">Spring Boot 项目介绍</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/07/13/spring-boot/spring-boot-core-operating-principle/">Spring Boot 核心运行原理介绍</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/07/16/spring-boot/spring-boot-sourcecode-springbootapplication/">【Spring Boot 源码学习】@SpringBootApplication 注解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/07/22/spring-boot/spring-boot-sourcecode-enableautoconfiguration/">【Spring Boot 源码学习】@EnableAutoConfiguration 注解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/07/30/spring-boot/spring-boot-sourcecode-autoconfigurationimportselector/">【Spring Boot 源码学习】走近 AutoConfigurationImportSelector</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/08/06/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-1/">【Spring Boot 源码学习】自动装配流程源码解析（上）</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/08/21/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-2/">【Spring Boot 源码学习】自动装配流程源码解析（下）</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/09/08/spring-boot/spring-boot-sourcecode-filteringspringbootcondition/">【Spring Boot 源码学习】深入 FilteringSpringBootCondition</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/09/11/spring-boot/spring-boot-sourcecode-onclasscondition/">【Spring Boot 源码学习】OnClassCondition 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/09/21/spring-boot/spring-boot-sourcecode-onbeancondition/">【Spring Boot 源码学习】OnBeanCondition 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/10/06/spring-boot/spring-boot-sourcecode-onwebapplicationcondition/">【Spring Boot 源码学习】OnWebApplicationCondition 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/10/15/spring-boot/spring-boot-sourcecode-conditional/">【Spring Boot 源码学习】@Conditional 条件注解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/10/22/spring-boot/spring-boot-sourcecode-httpencodingautoconfiguration/">【Spring Boot 源码学习】HttpEncodingAutoConfiguration 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/10/29/spring-boot/spring-boot-sourcecode-redisautoconfiguration/">【Spring Boot 源码学习】RedisAutoConfiguration 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/11/05/spring-boot/spring-boot-sourcecode-jedisconnectionconfiguration/">【Spring Boot 源码学习】JedisConnectionConfiguration 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/">【Spring Boot 源码学习】初识 SpringApplication</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/11/19/spring-boot/spring-boot-sourcecode-banner-printer/">【Spring Boot 源码学习】Banner 信息打印流程</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/11/24/spring-boot/spring-boot-sourcecode-custom-banner-printer/">【Spring Boot 源码学习】自定义 Banner 信息打印</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/">【Spring Boot 源码学习】BootstrapRegistryInitializer 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/">【Spring Boot 源码学习】ApplicationContextInitializer 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2023/12/10/spring-boot/spring-boot-sourcecode-applicationlistener/">【Spring Boot 源码学习】ApplicationListener 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/01/07/spring-boot/spring-boot-sourcecode-springapplication-customization/">【Spring Boot 源码学习】SpringApplication 的定制化介绍</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/">【Spring Boot 源码学习】BootstrapRegistry 详解</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/02/25/spring-boot/spring-boot-sourcecode-bootstrapcontext/">【Spring Boot 源码学习】深入 BootstrapContext 及其默认实现</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/03/02/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer-impl/">【Spring Boot 源码学习】BootstrapRegistry 初始化器实现</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/03/09/spring-boot/spring-boot-sourcecode-bootstrapcontext-actual-usage-scenario/">【Spring Boot 源码学习】BootstrapContext的实际使用场景</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/03/17/spring-boot/spring-boot-sourcecode-applicationcontextinitializer-impl/">【Spring Boot 源码学习】深入 ApplicationContext 初始化器实现</a>
        </td>
    </tr>
    <tr>
        <td align="left">
            <a href="/2024/03/24/spring-boot/spring-boot-sourcecode-sharedmetadatareaderfactorycontextinitializer/">【Spring Boot 源码学习】共享 MetadataReaderFactory 上下文初始化器</a>
        </td>
    </tr>
    <tr>
        <td align="left" > 
            <a href="2024/03/31/spring-boot/spring-boot-sourcecode-conditionevaluationreportlogginglistener/">【Spring Boot 源码学习】ConditionEvaluationReport 日志记录上下文初始化器</a> 
        </td>
    </tr>
</table>

# 三、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 3.1 run 方法源码初识
![](run.png)

上述截图就是 `SpringApplication` 的 `run` 方法核心代码。

下面 **Huazie** 将带着大家一起通读这块源码，从整体上了解下 `run` 方法核心流程。
## 3.2 引导上下文 BootstrapContext

```java
DefaultBootstrapContext bootstrapContext = createBootstrapContext();
```

翻看 `DefaultBootstrapContext` 的源码可知，从 **Spring Boot 2.4.0** 版本开始支持引导上下文。

![](DefaultBootstrapContext.png)

在[《BootstrapRegistryInitializer 详解》](../../../../../2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)中，Huazie 带大家详细分析了加载并初始化 `BootstrapRegistryInitializer` 的逻辑。而这里的 `createBootstrapContext` 方法就是用于创建默认的引导上下文对象 `DefaultBootstrapContext`，并利用 `BootstrapRegistry` 初始化器初始化该引导上下文对象。

想深入了解的朋友们，可查看 **Huazie** 下面列出的博文：

- [《BootstrapRegistry 详解》](../../../../../2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/)
- [《深入 BootstrapContext 及其默认实现》](../../../../../2024/02/25/spring-boot/spring-boot-sourcecode-bootstrapcontext/)
- [《BootstrapRegistry 初始化器实现》](../../../../../2024/03/02/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer-impl/)
- [《BootstrapContext的实际使用场景》](../../../../../2024/03/09/spring-boot/spring-boot-sourcecode-bootstrapcontext-actual-usage-scenario/)
## 3.3  系统属性【java.awt.headless】

```java
private static final String SYSTEM_PROPERTY_JAVA_AWT_HEADLESS = "java.awt.headless";

// 配置 java.awt.headless 系统属性
private void configureHeadlessProperty() {
    System.setProperty(SYSTEM_PROPERTY_JAVA_AWT_HEADLESS,
    System.getProperty(SYSTEM_PROPERTY_JAVA_AWT_HEADLESS, Boolean.toString(this.headless)));
}
```

`java.awt.headless` 是 **Java** 中的一个系统属性，用于指示 **Java** 应用程序是否运行在 **Headless** 模式下。**Headless** 模式是指系统缺少显示设备、键盘或鼠标的状态，通常应用于服务器环境，如应用集群、数据库集群等，这些环境通常通过网络远程操作，没有实际的显示设备。

在 **Java** 中，**AWT（Abstract Window Toolkit）** 是用于构建图形用户界面（**GUI**）应用的标准 **API** 接口。

**Java** 为 **AWT** 提供了两种模式实现以适应不同的运行环境：
- **标准模式**，适用于具有可用显示设备、驱动和图形用户界面的环境。
- **Headless 模式** ，适用于没有显示设备、驱动或图形用户界面的环境，例如服务器。

> **注意：** 设置 `java.awt.headless` 属性为 `true` 会使 **Java AWT** 工具包在 **headless** 模式下运行，这意味着它将不会尝试加载或访问与图形用户界面相关的资源或功能。
## 3.4 早期启动阶段

```java
SpringApplicationRunListeners listeners = getRunListeners(args);
listeners.starting(bootstrapContext, this.mainApplicationClass);
```

`SpringApplicationRunListeners` 中包含了一组 `SpringApplicationRunListener` 的集合。`SpringApplicationRunListener` 是 `SpringApplication` 的 `run` 方法的监听器，它用来监听 **Spring Boot** 应用的不同启动阶段，这些阶段都会发布对应的事件。

这里 `starting` 方法，就对应了最早期的启动阶段，它在 `run` 方法刚开始执行时就被立即调用。`starting` 方法里会发布 `ApplicationStartingEvent` 事件，通过监听该事件，应用可以执行一些非常早期的初始化工作，比如配置系统属性、初始化基础组件等等。
## 3.5 准备和配置应用环境

```java
ApplicationArguments applicationArguments = new DefaultApplicationArguments(args);
ConfigurableEnvironment environment = prepareEnvironment(listeners, bootstrapContext, applicationArguments);
configureIgnoreBeanInfo(environment);
```

`ApplicationArguments` 是 **Spring Boot** 中用于获取命令行参数的接口，其默认实现是 `DefaultApplicationArguments`。

`prepareEnvironment` 方法用于准备和配置应用程序的运行时环境，它会发布 `ApplicationEnvironmentPreparedEvent` 事件，通过监听该事件，应用程序可以执行一系列操作来准备和配置其运行环境。其返回的 `ConfigurableEnvironment` 对象，包含了应用程序的所有配置信息。

通过 `ConfigurableEnvironment` 对象，我们可以获取特定配置属性的值，也可以在运行时动态修改配置属性。

我们来看看 `configureIgnoreBeanInfo` 方法：

![](configureIgnoreBeanInfo.png)

![](spring-beaninfo-ignore.png)

在 `configureIgnoreBeanInfo` 方法中，可以看到如下代码：

```java
Boolean ignore = environment.getProperty(CachedIntrospectionResults.IGNORE_BEANINFO_PROPERTY_NAME, Boolean.class, Boolean.TRUE);
```

从上述代码中，可以看到通过 `environment` 变量获取属性名为 `spring.beaninfo.ignore` 的属性值，其 `getProperty` 方法有三个参数：
- 第一个参数是属性名。
- 第二个参数是期望返回的属性值的类型，这里是 `Boolean.class`。
- 第三个参数是默认值，如果找不到属性或者属性不能被转换为 `Boolean` 类型，则使用 `Boolean.TRUE` 作为默认值。

系统属性 `spring.beaninfo.ignore`  用于指示 **Spring** 在调用 **JavaBeans Introspector** 时使用`Introspector.IGNORE_ALL_BEANINFO` 模式。如果此属性的值为 `true`，则 **Spring** 会跳过搜索 `BeanInfo` 类（通常适用于以下情况：应用程序中的 **beans** 从一开始就没有定义这样的类）。

默认值是 `false`，表示 **Spring** 会考虑所有的 `BeanInfo` 元数据类，就像标准 `Introspector.getBeanInfo(Class)` 调用那样。如果在启动时或延迟加载时，反复访问不存在的 `BeanInfo` 类开销很大，可以考虑将此标志切换为 `true`。

**请注意**：如果存在反复访问不存在的 `BeanInfo` 类，可能也表明缓存未奏效。最好将 **Spring** 的 **jar** 包与应用类放在同一个 `ClassLoader` 中，这样可以在任何情况下与应用程序的生命周期一起进行干净的缓存。对于 **Web** 应用程序，如果采用多 `ClassLoader` 布局，可以考虑在 **web.xml** 中声明一个本地的 `org.springframework.web.util.IntrospectorCleanupListener`，这也可以实现有效的缓存。

## 3.6 打印 Banner 信息

```java
Banner printedBanner = printBanner(environment);
```

`printBanner` 方法用于 **Spring Boot** 启动时的 **Banner** 信息打印。

想要深入了解 **Banner** 打印的读者们，请查看如下博文：

- [《Banner 信息打印流程》](../../../../../2023/11/19/spring-boot/spring-boot-sourcecode-banner-printer/)
- [《自定义 Banner 信息打印》](../../../../../2023/11/24/spring-boot/spring-boot-sourcecode-custom-banner-printer/)

## 3.7 新建应用上下文

```java
ConfigurableApplicationContext context = createApplicationContext();

protected ConfigurableApplicationContext createApplicationContext() {
    return this.applicationContextFactory.create(this.webApplicationType);
}
```

上述 `createApplicationContext` 方法的功能是：根据给定的 **Web** 应用程序类型 `webApplicationType` 创建一个可配置的应用上下文对象 `ConfigurableApplicationContext` 。

在[《初识 SpringApplication》](../../../../../2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)这篇博文的 **2.2 小节**【**Web 应用类型推断**】中，大家可以看到 **Web** 应用程序类型 `webApplicationType` 是如何获取的，这里不赘述了，感兴趣的可以自行查看。
## 3.8 准备和配置应用上下文
```java
context.setApplicationStartup(this.applicationStartup);
prepareContext(bootstrapContext, context, environment, listeners, applicationArguments, printedBanner);
```

`setApplicationStartup` 方法用于设置当前应用上下文的 `ApplicationStartup` ，这允许应用上下文在启动期间记录指标。

`prepareContext` 方法用于准备和配置应用程序上下文，这里会依次发布如下事件：
1.  `ApplicationContextInitializedEvent`：当 `SpringApplication` 启动并且 `ApplicationContext` 已准备好，且 `ApplicationContextInitializer` 集合已被调用，但在加载任何 **bean** 定义之前，将发布该事件。
2.  `ApplicationPreparedEvent` ：当 `SpringApplication` 启动并且 `ApplicationContext` 已经完全准备好但尚未刷新时，将发布事件。在此阶段，**bean** 定义将被加载，环境已经准备好可以使用。

## 3.9 刷新应用上下文

```java
static final SpringApplicationShutdownHook shutdownHook = new SpringApplicationShutdownHook();

refreshContext(context);

private void refreshContext(ConfigurableApplicationContext context) {
    if (this.registerShutdownHook) {
        shutdownHook.registerApplicationContext(context);
    }
    refresh(context);
}

protected void refresh(ConfigurableApplicationContext applicationContext) {
    applicationContext.refresh();
}
```

`registerShutdownHook` 变量表示是否应注册一个关闭钩子，默认为 `true`。

`SpringApplicationShutdownHook` 是一个用于执行 **Spring Boot** 应用程序优雅关闭的 `Runnable` 关机钩子。这个钩子跟踪已注册的应用程序上下文以及通过 `SpringApplication.getShutdownHandlers()` 注册的任何操作。

`refreshContext` 方法里面可以看到调用 `refresh` 方法，`refresh` 方法里面则是调用 `ConfigurableApplicationContext`【实现类是 `AbstractApplicationContext` ，该类属于 **spring-context** 包】的 `refresh` 方法，该方法是用来刷新底层的应用上下文。

它会加载或刷新配置的持久化表示，这可能来自基于 **Java** 的配置、**XML** 文件、属性文件、关系数据库模式或其他某种格式。调用此方法后，要么实例化所有单例对象，要么不实例化任何单例对象。

它最后会发布 `ContextRefreshedEvent` 事件，通过监听该事件，可以执行一些应用上下文初始化或刷新后需要进行的操作。

## 3.10 afterRefresh 方法
刷新应用上下文之后，调用 `afterRefresh` 方法。该方法的实现默认为空，可由开发人员自行扩展。

![](afterRefresh.png)

## 3.11 打印启动日志

```java
long startTime = System.nanoTime();

Duration timeTakenToStartup = Duration.ofNanos(System.nanoTime() - startTime);
if (this.logStartupInfo) {
    new StartupInfoLogger(this.mainApplicationClass).logStarted(getApplicationLog(), timeTakenToStartup);
}
```

`System.nanoTime()` 用于获取当前时间的纳秒值。

`Duration.ofNanos()` 用于将纳秒数转换为 `Duration` 对象，`timeTakenToStartup` 表示 **Spring Boot** 应用启动所需的时间。

`logStartupInfo` 表示是否需要记录启动信息，如果为 `true`，则需要记录启动信息。

`StartupInfoLogger` 类用于在应用程序启动时记录应用信息，其中 `logStarted` 方法用于以 `INFO` 日志级别打印应用启动时间。

![](logStarted.png)
![](getStartMessage.png)

实际运行日志信息类似如下：

![](result.png)
## 3.12 Spring 容器启动完成

```java
listeners.started(context, timeTakenToStartup);
```

这里表示上下文已经刷新，应用程序已经启动，但是 `CommandLineRunners` 和 `ApplicationRunners` 尚未被调用。

`SpringApplicationRunListeners` 的 `started` 方法里会发布 `ApplicationStartedEvent` 事件，通知监听器 **Spring** 容器启动完成。

## 3.13 callRunners 方法

```java
callRunners(context, applicationArguments);
```

`callRunners` 方法里面会调用 `ApplicationRunner` 和 `CommandLineRunner` 的运行方法。

![](callRunners.png)

通过阅读上述代码，可以总结如下：
- 首先，从 `context` 中获取类型为 `ApplicationRunner` 和 `CommandLineRunner` 的 **Bean**；
- 接着，将它们放入 `List` 列表中，并进行排序。
- 最后，再遍历排序后的 `ApplicationRunner` 和 `CommandLineRunner` 的 **Bean**，并调用它们的 `run` 方法。

**Spring Boot** 提供 `ApplicationRunner` 和 `CommandLineRunner` 这两种接口，是为了通过它们来实现在容器启动时执行一些操作。在同一个应用上下文中可以定义多个 `ApplicationRunner` 或 `CommandLineRunner` 的**bean**，并可以使用 `Ordered` 接口或 `@Order` 注解进行排序。

`ApplicationRunner` 和 `CommandLineRunner` 这两个接口都有一个 run 方法，但不同之处是：
- `ApplicationRunner` 中 `run` 方法的参数为 `ApplicationArguments`
- `CommandLineRunner` 中 `run` 方法的参数为 **字符串数组**

如果需要访问 `ApplicationArguments` 而不是原始的字符串数组，大家可以考虑使用 `ApplicationRunner`。

## 3.14 Spring 容器正在运行中

```java
Duration timeTakenToReady = Duration.ofNanos(System.nanoTime() - startTime);
listeners.ready(context, timeTakenToReady);
```

这里表示应用上下文已经刷新，所有的 `CommandLineRunners` 和 `ApplicationRunners` 都已被调用，应用程序已准备好处理请求。

`SpringApplicationRunListeners` 的 `ready` 方法里会发布 `ApplicationReadyEvent` 事件，通知监听器 **Spring** 容器正在运行中。

在 **Spring Boot 2.6.0** 版本之前，大家看到调用的是 `SpringApplicationRunListener` 的 `running` 方法。从 **Spring Boot 2.6.0** 版本开始，新增了 `ready` 方法替代 `running` 方法。在 **Spring Boot 3.0.0** 版本中正式去除 `running` 方法。

![](ready.png)

## 3.15 异常处理

```java
handleRunFailure(context, ex, listeners);
```

从 **3.5** 到 **3.13** 小节 ，如果出现异常，则会捕获后调用 `handleRunFailure` 进行异常处理。

**3.14** 小节，同样它如果出现异常，也会捕获后调用 `handleRunFailure` 进行异常处理。

`handleRunFailure` 方法里会发布 `ApplicationFailedEvent` 事件，通过监听该事件，开发人员可以实现如下的一些操作：

- **错误日志记录**：当应用启动失败时，可以记录详细的错误信息到日志文件中，便于后续的问题排查和分析。
- **通知发送**：在应用启动失败时，可以发送通知给相关的开发或运维人员，以便他们能够及时响应并处理问题。
- **数据备份**：如果应用在启动过程中出现异常，可能需要对某些关键数据进行备份，以防止数据丢失。
- **资源清理**：在应用启动失败的情况下，可能需要释放或清理已经分配的资源，如数据库连接、文件句柄等。
- **尝试自动恢复**：在某些情况下，可以尝试自动重启应用或者执行其他恢复操作，以减少人工干预的需求。
- **自定义处理逻辑**：根据具体的业务需求，实现自定义的错误处理逻辑，比如回滚事务、关闭网络连接等。

有关这块更详细的内容，后续 **Huazie** 将专门出一篇讲解，敬请期待！！！

# 四、总结

本篇 **Huazie** 向大家初步介绍了 `SpringApplication` 的 `run` 方法核心流程。由于篇幅受限，其中很多环节并未深入讲解，后续 **Huazie** 将会针对这些内容深入分析，和大家一起从源码详细了解 **Spring Boot** 的运行流程。

