---
title: 【Spring Boot 源码学习】ConditionEvaluationReport 日志记录上下文初始化器
date: 2024-03-31 23:22:19
updated: 2024-04-01 10:54:37
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - ConditionEvaluationReportLoggingListener
  - ConditionEvaluationReport
  - ConditionEvaluationReportListener
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
上篇博文[《共享 MetadataReaderFactory 上下文初始化器》](/2024/03/24/spring-boot/spring-boot-sourcecode-sharedmetadatareaderfactorycontextinitializer/)，**Huazie** 带大家详细分析了 
`SharedMetadataReaderFactoryContextInitializer` 。而在 **spring-boot-autoconfigure** 子模块中预置的上下文初始化器中，除了共享 `MetadataReaderFactory` 上下文初始化器，还有一个尚未分析。

那么本篇就来详细分析一下 `ConditionEvaluationReportLoggingListener` 【即 `ConditionEvaluationReport` 日志记录上下文初始化器】。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
    <tr>
        <td rowspan="29" align="left">
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
</table>

# 三、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

## 3.1 源码初识

我们先来看看 `ConditionEvaluationReportLoggingListener` 的部分源码，如下：

```java
public class ConditionEvaluationReportLoggingListener
        implements ApplicationContextInitializer<ConfigurableApplicationContext> {

    private final Log logger = LogFactory.getLog(getClass());

    private ConfigurableApplicationContext applicationContext;

    private ConditionEvaluationReport report;

    private final LogLevel logLevelForReport;

    public ConditionEvaluationReportLoggingListener() {
        this(LogLevel.DEBUG);
    }

    public ConditionEvaluationReportLoggingListener(LogLevel logLevelForReport) {
        Assert.isTrue(isInfoOrDebug(logLevelForReport), "LogLevel must be INFO or DEBUG");
        this.logLevelForReport = logLevelForReport;
    }
    
    // 省略。。。

    @Override
    public void initialize(ConfigurableApplicationContext applicationContext) {
        this.applicationContext = applicationContext;
        applicationContext.addApplicationListener(new ConditionEvaluationReportListener());
        if (applicationContext instanceof GenericApplicationContext) {
            // Get the report early in case the context fails to load
            this.report = ConditionEvaluationReport.get(this.applicationContext.getBeanFactory());
        }
    }
    
    // 省略。。。
}
```
从上述源码中，我们可以看出 `ConditionEvaluationReportLoggingListener` 实现了 `ApplicationContextInitializer<ConfigurableApplicationContext>` 【即应用上下文初始化器接口】，有关 `ApplicationContextInitializer` 的详细介绍，请查看[《ApplicationContextInitializer 详解》](/2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/)。

它有三个成员变量，分别是：

- `ConfigurableApplicationContext applicationContext` ： 应用上下文对象
- `ConditionEvaluationReport report` ：条件评估报告对象，用于报告和记录条件评估详细信息。
- `LogLevel logLevelForReport` ：条件评估报告的日志级别，包含 `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`, `OFF`

再来看看构造方法，它有两个：

- **无参构造方法**：初始化日志级别为 `DEBUG`【默认通过它实例化该上下文初始化器】
- **带 `LogLevel` 参数的构造方法**：`Assert.isTrue` 是用于验证一个条件是否为真。通过 `isInfoOrDebug` 来判断日志级别参数 `logLevelForReport`  是否是 **INFO** 或 **DEBUG** 级别，如果不是，则会抛出一个 `IllegalArgumentException` 异常并显示错误信息 **"LogLevel must be INFO or DEBUG"**。

我们继续查看 `initialize` 方法，可以看到 ：
- 首先，初始化成员变量应用上下文对象 `applicationContext`，便于后续使用。
- 然后，向应用上下文对象中添加一个应用监听器实现【即 `ConditionEvaluationReportListener`】，这里可查看 **3.2** 小节的内容。
- 最后，如果 `applicationContext` 是 `GenericApplicationContext` 的一个实例对象，则通过 `ConditionEvaluationReport` 的静态方法 `get` 来获取指定 **Bean** 工厂中的 **条件评估报告** 实例对象。
## 3.2 ConditionEvaluationReport 监听器
下面继续来查看 `ConditionEvaluationReportListener` 的源码：

![](ConditionEvaluationReportListener.png)

阅读上述源码，可以看到 `ConditionEvaluationReportListener` 实现了 `GenericApplicationListener` 接口，继续翻看 `GenericApplicationListener` 接口源码：

![](GenericApplicationListener.png)

继续翻看 `SmartApplicationListener` 接口源码：

![](SmartApplicationListener.png)

从上述源码中，我们发现 `GenericApplicationListener` 继承了 `SmartApplicationListener`，而 `SmartApplicationListener` 则继承了 `ApplicationListener<ApplicationEvent>`。

`GenericApplicationListener` 是 **Spring** 框架中的一个接口，它扩展了 `ApplicationListener` 接口，暴露了更多的元数据，如支持的事件和源类型。在 **Spring Framework 4.2** 及更高版本中，`GenericApplicationListener` 替代了基于类的 `SmartApplicationListener`，允许你使用 `ResolvableType` 来指定支持的事件类型，而不仅仅是 `Class` 类型，这样就可以在运行时更准确地解析和匹配事件类型。

> **知识点：**  `ResolvableType` 是 **Spring** 框架中提供的一个工具类，它用于在运行时解析和处理 **Java** 泛型信息。在 **Java 5** 引入泛型之后，为了支持泛型，新增了 `Type` 类来表示 **Java** 中的某种类型。然而，反射包中提供的方法在获取泛型类型时，通常返回的是 `Type` 或其子类的实例，使用时可能需要进行繁琐的强制类型转换。`ResolvableType` 的出现就是为了简化对泛型信息的获取和处理。它能够将 `Class`、`Field`、`Method` 等描述为 `ResolvableType`（即转换为 `Type`），从而方便地进行泛型的解析和操作。通过使用 `ResolvableType`，你可以轻松地获取 **类、接口、属性、方法** 等的泛型信息，而无需进行复杂的类型转换或编写繁琐的代码。

现在我们再来看看 `ConditionEvaluationReportListener` 中重写的 `supportsEventType(ResolvableType)`  方法：

![](supportsEventType.png)

也就是说，该监听器实际上监听是如下两个事件：
-  `ContextRefreshedEvent` ：上下文刷新事件。该事件会在 `ApplicationContext` 完成初始化或刷新时发布。
- `ApplicationFailedEvent` ：应用启动失败事件。该事件是在 **Spring Boot** 应用启动失败时触发，一般发生在 `ApplicationStartedEvent` 事件之后。

我们继续查看 `ConditionEvaluationReportListener` 的核心方法 `onApplicationEvent` ，发现它直接调用了 `ConditionEvaluationReportLoggingListener` 中的 `onApplicationEvent` 方法，来实现条件评估报告的日志打印功能。
## 3.3  onApplicationEvent 方法
我们继续查看 `ConditionEvaluationReportLoggingListener` 中的 `onApplicationEvent` 方法：

![](onApplicationEvent.png)


从上图中，可以看到这里针对 **3.2** 中监听器监听的两个事件分别进行了处理，而这里的核心方法就是 `logAutoConfigurationReport(boolean)` 方法。

继续查看 `logAutoConfigurationReport(boolean)` 方法：

![](logAutoConfigurationReport.png)

从上图中，我们可以简单总结一下：

- 首先，如果条件评估报告 `report` 为空，则通过 `ConditionEvaluationReport` 的静态方法 `get` 来获取当前应用上下文指定的 **Bean** 工厂中的 **条件评估报告** 实例对象。
- 判断 `report` 中的条件评估结果是否为空？
    - 如果不为空，判断条件评估报告的日志级别 
        - 如果是 `INFO` 级别 ，则继续
            - 如果当前允许记录 `INFO` 级别日志，则按 `INFO` 级别输出相关的条件评估结果的日志信息。
        - 如果是 `DEBUG` 级别，则继续
            - 如果当前允许记录 `DEBUG` 级别日志，则按 `DEBUG` 级别输出相关的条件评估结果的日志信息。

## 3.4 条件评估报告的打印展示

首先，我们在当前 **Spring Boot** 项目中设置当前的日志级别为 `DEBUG`【当然还可以指定其他日志配置文件，这里不展开讲了】：

![](debug.png)

运行我们的自测类或者应用主类，可以看到如下的运行结果：

![](result-1.png)
![](result-2.png)
![](result-3.png)
![](result-4.png)

从上述运行结果中，可以看出条件评估报告中包含如下的内容：
- `Positive matches`：**正匹配**，即 `@Conditional` 条件为真时，相关的配置类被`Spring` 容器加载，配置类中定义的 **bean** 和其他组件将被创建并添加到 `Spring` 的应用上下文中。
- `Negative matches`：**负匹配**，即 `@Conditional` 条件为假时，相关的配置类未被 **Spring** 容器加载。尽管相关的配置类存在于项目中，但由于某些条件不满足（如缺少必要的依赖或配置），**Spring** 容器不会创建该配置类中定义的 **bean**。
- `Exclusions`：**排除**，即明确要排除的配置类，这些被排除的自动配置类中的组件将不会被创建。
- `Unconditional classes`：**无条件类**，即自动配置类不包含任何类级别的条件。与 **Positive matches** 和 **Negative matches** 不同，这些类不依赖于任何特定的条件来决定是否加载。它们总是会被 **Spring** 容器处理，无论其他条件如何。

# 四、总结
本篇 **Huazie** 带大家一起分析了 **spring-boot-autoconfigure** 子模块中预置的另一个应用上下文初始化器实现 `ConditionEvaluationReportLoggingListener` ，它实现了条件评估报告的打印记录功能，极大地方便了开发者定位配置类加载问题。

