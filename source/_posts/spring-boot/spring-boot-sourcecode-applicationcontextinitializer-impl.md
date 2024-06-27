---
title: 【Spring Boot 源码学习】深入 ApplicationContext 初始化器实现
date: 2024-03-17 18:57:56
updated: 2024-03-17 18:57:56
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - ApplicationContextInitializer
  - ConfigurationWarningsApplicationContextInitializer
  - ContextIdApplicationContextInitializer
  - DelegatingApplicationContextInitializer
  - RSocketPortInfoApplicationContextInitializer
  - ServerPortInfoApplicationContextInitializer
  - SharedMetadataReaderFactoryContextInitializer
  - ConditionEvaluationReportLoggingListener
  - 应用上下文初始化器实现
---



![](/images/spring-boot-logo.png)

# 一、引言
前面的博文[《ApplicationContextInitializer 详解》](../../../../../2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/)，**Huazie** 带大家详细分析了 `ApplicationContextInitializer` 的加载和初始化的逻辑，不过有关 `ApplicationContextInitializer` 接口的实现尚未提及 。

那本篇 **Huazie** 就带大家一起分析 **Spring Boot** 中预置的应用上下文初始化器实现【即 `ApplicationContextInitializer` 接口实现类】的源码，了解在 **Spring** 容器刷新之前初始化应用程序上下文的一些具体操作。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="27" align="left" > 
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
    <td align="left" > 
      <a href="/2023/08/06/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-1/">【Spring Boot 源码学习】自动装配流程源码解析（上）</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/08/21/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-2/">【Spring Boot 源码学习】自动装配流程源码解析（下）</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/09/08/spring-boot/spring-boot-sourcecode-filteringspringbootcondition/">【Spring Boot 源码学习】深入 FilteringSpringBootCondition</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/09/11/spring-boot/spring-boot-sourcecode-onclasscondition/">【Spring Boot 源码学习】OnClassCondition 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/09/21/spring-boot/spring-boot-sourcecode-onbeancondition/">【Spring Boot 源码学习】OnBeanCondition 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/10/06/spring-boot/spring-boot-sourcecode-onwebapplicationcondition/">【Spring Boot 源码学习】OnWebApplicationCondition 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/10/15/spring-boot/spring-boot-sourcecode-conditional/">【Spring Boot 源码学习】@Conditional 条件注解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/10/22/spring-boot/spring-boot-sourcecode-httpencodingautoconfiguration/">【Spring Boot 源码学习】HttpEncodingAutoConfiguration 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/10/29/spring-boot/spring-boot-sourcecode-redisautoconfiguration/">【Spring Boot 源码学习】RedisAutoConfiguration 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/11/05/spring-boot/spring-boot-sourcecode-jedisconnectionconfiguration/">【Spring Boot 源码学习】JedisConnectionConfiguration 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/">【Spring Boot 源码学习】初识 SpringApplication</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/11/19/spring-boot/spring-boot-sourcecode-banner-printer/">【Spring Boot 源码学习】Banner 信息打印流程</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/11/24/spring-boot/spring-boot-sourcecode-custom-banner-printer/">【Spring Boot 源码学习】自定义 Banner 信息打印</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/">【Spring Boot 源码学习】BootstrapRegistryInitializer 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/">【Spring Boot 源码学习】ApplicationContextInitializer 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/12/10/spring-boot/spring-boot-sourcecode-applicationlistener/">【Spring Boot 源码学习】ApplicationListener 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2024/01/07/spring-boot/spring-boot-sourcecode-springapplication-customization/">【Spring Boot 源码学习】SpringApplication 的定制化介绍</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/">【Spring Boot 源码学习】BootstrapRegistry 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2024/02/25/spring-boot/spring-boot-sourcecode-bootstrapcontext/">【Spring Boot 源码学习】深入 BootstrapContext 及其默认实现</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2024/03/02/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer-impl/">【Spring Boot 源码学习】BootstrapRegistry 初始化器实现</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2024/03/09/spring-boot/spring-boot-sourcecode-bootstrapcontext-actual-usage-scenario/">【Spring Boot 源码学习】BootstrapContext的实际使用场景</a> 
    </td>
  </tr>
</table>

# 三、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 3.1 spring-boot 子模块中内置的实现类
我们先来看一张截图：

![](applicationcontextinitializer.png)

从上图中可以看出，**spring-boot** 子模块中配置的 `ApplicationContextInitializer` 实现一共有 **5** 个，下面我们一一来介绍下：

### 3.1.1 ConfigurationWarningsApplicationContextInitializer
该类用于报告常见配置错误的警告，我们来看看相关源码：

```java
public class ConfigurationWarningsApplicationContextInitializer
    implements ApplicationContextInitializer<ConfigurableApplicationContext> {

    @Override
    public void initialize(ConfigurableApplicationContext context) {
        context.addBeanFactoryPostProcessor(new ConfigurationWarningsPostProcessor(getChecks()));
    }

    // 省略其他。。。
}
```

阅读上述源码，我们可以看到 `initialize` 方法里，通过 `ConfigurableApplicationContext` 的 `addBeanFactoryPostProcessor` 方法，在 **应用程序上下文** 中添加了一个 `BeanFactoryPostProcessor` 实现【 该实现类为 `ConfigurationWarningsPostProcessor`】。

`BeanFactoryPostProcessor` 是 **Spring** 框架中的一个接口，它的作用是在 **Spring** 容器初始化时对 **Bean** 的定义进行修改或增强【添加属性、设置依赖关系等】。

在介绍 `ConfigurationWarningsPostProcessor` 之前，先来看看 `getChecks` 方法：

```java
protected Check[] getChecks() {
    return new Check[] { new ComponentScanPackageCheck() };
}
```

我们继续查看 `ComponentScanPackageCheck` ，由于篇幅受限，**Huazie** 贴下截图：

![](componentscanpackagecheck.png)

`ComponentScanPackageCheck` 是 `ConfigurationWarningsApplicationContextInitializer` 中的一个静态内部类，它的目的是在**Spring Boot** 应用启动时，检查 `@ComponentScan` 的使用情况，确保没有错误或不推荐的配置方式。通过 `ComponentScanPackageCheck` 的 `getWarning` 方法的检查，如果发现了不恰当的 `@ComponentScan` 使用，它会生成相应的警告信息，帮助开发者及时发现并修正潜在的配置问题。


下面我们可以来分析下 `ConfigurationWarningsPostProcessor`，如下截图：

![](configurationwarningspostprocessor.png)

该类也是一个静态内部类，它同时实现了 `PriorityOrdered` 和 `BeanDefinitionRegistryPostProcessor` 接口：

- `PriorityOrdered` ：实现该接口是用于提高其在多个 `BeanFactoryPostProcessor` 处理中的执行优先级。
- `BeanDefinitionRegistryPostProcessor`：它是对 `BeanFactoryPostProcessor` 的扩展，允许在常规的 `BeanFactoryPostProcessor` 检测启动之前注册更多的 **bean** 定义，这些定义反过来定义了`BeanFactoryPostProcessor` 实例【可查看 `PostProcessorRegistrationDelegate` 了解】。

从上述截图中，我们可以看到 `postProcessBeanFactory` 方法【`BeanFactoryPostProcessor` 接口定义的方法】是空实现，而`postProcessBeanDefinitionRegistry` 方法【`BeanDefinitionRegistryPostProcessor` 接口定义的方法】里，遍历了 `checks` 数组中的每个检查项，并调用 `check.getWarning(registry)` 方法获取警告信息。如果警告信息不为空，则调用私有方法 `warn(message)` 打印警告信息。


### 3.1.2 ContextIdApplicationContextInitializer

先来看看 `ContextIdApplicationContextInitializer` 的源码，如下：

![](contextidapplicationcontextinitializer.png)
![](contextid.png)

通过阅读上述源码，可以看出 `ContextIdApplicationContextInitializer` 是一个用于设置 **Spring ApplicationContext ID** 的应用上下文初始化器。其中，`spring.application.name` 属性用于创建 **ID**。如果该属性未设置，则使用 **application**。

我们在 `initialize` 方法中，还看到了如下的代码：

```java
applicationContext.getBeanFactory().registerSingleton(ContextId.class.getName(), contextId);
```

这里就是将一个名为 `ContextId` 的类注册为单例对象，并将其存储在 **Spring** 的 `ApplicationContext` 中。然后我们就可以在应用程序的不同部分共享和重用同一个 `ContextId` 实例，而无需每次都创建新的实例。

### 3.1.3 DelegatingApplicationContextInitializer
同样先来看看 `DelegatingApplicationContextInitializer` 的源码，如下截图：

![](delegatingapplicationcontextinitializer.png)

通过阅读该类的 `initialize` 方法，我们可以看出 `DelegatingApplicationContextInitializer` 初始化工作是委托给其他在 `context.initializer.classes` 环境属性下指定的应用上下文初始化器进行的。

下面的 **3.3** 小节，我们在自定义 **`ApplicationContext` 初始化器实现** 时就会用到。
### 3.1.4 RSocketPortInfoApplicationContextInitializer

无需多言，直接查看 `RSocketPortInfoApplicationContextInitializer` 的源码，如下：

```java
public class RSocketPortInfoApplicationContextInitializer
    implements ApplicationContextInitializer<ConfigurableApplicationContext> {
    @Override
    public void initialize(ConfigurableApplicationContext applicationContext) {
        applicationContext.addApplicationListener(new Listener(applicationContext));
    }
    
    // 省略。。。
}
```

阅读上述的 initialize 方法，可以看到这里向应用上下文中添加了一个 `ApplicationListener`，而这个 `Listener` 是 `RSocketPortInfoApplicationContextInitializer` 中的一个静态内部类。

继续阅读 `Listener` 的源码：

![](listener.png)

`Listener` 是用来监听 `RSocketServerInitializedEvent` 事件，该事件是在应用程序上下文刷新且 `RSocketServer` 准备就绪后发布的。

继续查看 `onApplicationEvent` 方法，我们可以看出该监听器是用来设置 `RSocketServer` 服务器实际监听的端口的环境属性。属性 `local.rsocket.server.port`  可以直接使用 `@Value` 注入到测试中，也可以通过 `Environment` 获取。另外该属性会自动向上传播到任何父上下文。
### 3.1.5 ServerPortInfoApplicationContextInitializer
同样还是从 `ServerPortInfoApplicationContextInitializer` 源码入手，如下所示：

![](serverportinfoapplicationcontextinitializer.png)

通过阅读上面的 `initialize` 方法，可以看到这里也是比较简单，直接向应用上下文中添加了一个 `ApplicationListener` ，当然这个应用事件监听器比较特殊，就是其本身，因为 `ServerPortInfoApplicationContextInitializer` 实现了 `ApplicationListener` 接口。

该 `ApplicationListener` 监听的事件是 `WebServerInitializedEvent`，它是一个在 `WebServer` 准备就绪时发布的事件。

我们继续阅读 `onApplicationEvent` 方法的源码：

![](onapplicationevent.png)

我们来简单总结如下：

该应用事件监听器用于设置 **WebServer** 服务器实际监听的端口的环境属性。属性 `local.server.port` 【如果 `WebServerInitializedEvent` 有一个服务器命名空间，它将被用来构造属性名称。例如，**“management” actuator** 上下文将具有属性名称 `local.management.port`】可以直接使用 `@Value` 注入到测试中，也可以通过 `Environment` 获取。该属性同样会自动向上传播到任何父上下文。

> **Actuator** 是 **Spring Boot** 提供的一个开发库，它允许开发人员在运行时监控和管理应用程序。通过 **Actuator**，你可以查看应用程序的运行状况、性能指标、日志信息等。同时，它也提供了一些内置的管理端点，如健康检查、环境信息、应用信息等，方便开发人员进行调试和监控。**Actuator** 还提供了扩展机制，允许你自定义管理端点，以满足特定的需求。

## 3.2 spring-boot-autoconfigure 子模块中内置的实现类
同样我们先看截图：

![](applicationcontextinitializer-1.png)

从上图中可以看出，**spring-boot-autoconfigure** 子模块中配置的 `ApplicationContextInitializer` 实现有 **2** 个，下面来简单介绍下：

### 3.2.1 SharedMetadataReaderFactoryContextInitializer

`SharedMetadataReaderFactoryContextInitializer` 是一个应用上下文初始化器，主要作用是在 **Spring** 应用程序上下文创建之初，初始化一个共享的 `MetadataReaderFactory` 实例到在 **Spring** 应用上下文中。这样，在整个应用程序生命周期内，不同的组件在需要读取类的元数据时，都可以使用一个共享的 `MetadataReaderFactory` 实例，而无需每次都创建新的实例。

> 在 **Spring** 中，元数据（**metadata**）是用来描述 **Bean** 信息的数据，例如类名、方法名、属性名等。在应用程序运行时，**Spring** 会读取这些元数据来创建和管理 **Bean**。而 `MetadataReaderFactory` 就是负责读取和解析类的元数据，比如注解、类属性等。


这块的逻辑比较复杂，**Huazie** 后续将再出一篇博文详细分析，敬请期待！
### 3.2.2 ConditionEvaluationReportLoggingListener
`ConditionEvaluationReportLoggingListener` 是一个用于将 `ConditionEvaluationReport` 写入日志的应用上下文初始化器，该应用上下文初始化器并不打算在多个应用程序上下文实例之间共享。

当 **Spring** 应用程序上下文初始化时，它会评估所有使用条件注解的 **bean** 定义和配置。这些条件可能基于类是否存在、特定的属性设置、其他 **bean** 是否存在等。`ConditionEvaluationReport` 记录了每个条件注解的评估结果，包括哪些条件通过了（即 **bean** 或配置被创建或执行了），哪些条件没有通过（即 **bean** 或配置被跳过了）。

`ConditionEvaluationReport` 的评估结果报告默认将以 `DEBUG` 级别进行记录。崩溃报告会触发 `info` 级别的输出，建议再次运行并启用 `debug` 级别以显示报告。

这块的逻辑也比较复杂，**Huazie** 后续也会出一篇博文详细介绍下，大家可以期待一下。

## 3.3 自定义应用上下文初始化器实现

上面 Huazie 同大家一起分析了 **Spring Boot** 中一些内置的应用上下文初始化器实现，相信对于如何实现 `ApplicationContextInitializer` 接口，已经有了较为深入的了解。

### 3.3.1 定义 DemoApplicationContextInitializer 

那下面就让我们自定义 `ApplicationContext` 初始化器实现，如下所示：

```java
public class DemoApplicationContextInitializer implements ApplicationContextInitializer<ConfigurableApplicationContext>, Ordered {

    private int order = 0;

    @Override
    public void initialize(ConfigurableApplicationContext applicationContext) {
        User user = new User("Huazie", 18);
        applicationContext.getBeanFactory().registerSingleton(User.class.getName(), user);
    }

    public void setOrder(int order) {
        this.order = order;
    }

    @Override
    public int getOrder() {
        return this.order;
    }

}
```

上述 `DemoApplicationContextInitializer` 的 `initialize` 方法中，我们注册了一个 `User` 类的单例 **Bean**。

### 3.3.2 添加 DemoApplicationContextInitializer

**现在自定义的应用上下文初始化器有了，我们该如何添加它呢？**

通过阅读 `SpringApplication` 的源码 和 本篇 **3.1.3** 小节的介绍，我们可以总结如下的三种方式：

-  在 `META-INF/spring.factories` 中添加 `org.springframework.context.ApplicationContextInitializer` 的配置。这种方式，我们从 [《ApplicationContextInitializer 详解》](../../../../../2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/) 的 **3.2** 小节可见一斑。
  ```bash
  org.springframework.context.ApplicationContextInitializer=com.example.demo.DemoApplicationContextInitializer
  ```

- 通过 `SpringApplication` 中的 `addInitializers` 方法添加。其实这里在笔者的[《SpringApplication 的定制化介绍》](../../../../../2024/01/07/spring-boot/spring-boot-sourcecode-springapplication-customization/)中的 **1.6** 小节也提及过。
  ```java
  SpringApplication springApplication = new SpringApplication(DemoApplication.class);
  springApplication.addInitializers(new DemoApplicationContextInitializer());
  // 其他省略。。。
  ```
- 在 **application.properties** 中添加 `context.initializer.classes` 的属性配置。这里实际上来源于 **3.1.3** 小节的 **DelegatingApplicationContextInitializer**。
  ```bash            
  # 逗号分隔的类名列表
  context.initializer.classes=com.example.demo.DemoApplicationContextInitializer
  ```
  在 **application.yml** 中添加 `context.initializer.classes` 的属性配置
  ```yml
  # 在 YAML 中，数组或列表元素使用 - 符号来定义
  context:  
      initializer:  
        classes:  
            - com.example.demo.DemoApplicationContextInitializer  
  ```

### 3.3.3 实际演示

我们采用第三种添加方式，配置截图如下：

![](config.png)

添加如下自测类【用来演示获取在 `DemoApplicationContextInitializer` 中注册的 `User` 类的单例 **Bean** 对象】：

```java
import com.example.demo.pojo.User;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
public class DemoApplicationTests {

    @Autowired
    private User user;

    @Test
    public void test() {
        System.out.println("User = " + user);
    }
}
```

我们来看看运行结果，如下所示：

![](result.png)

从上图可以看出，我们自定义的应用上下文初始化器实现显然已经执行了，并且成功注册了 `User` 类的单例 **Bean** 对象。

# 四、总结
本篇 **Huazie** 带大家一起分析了 **Spring Boot** 中预置的 `ApplicationContext` 初始化器实现，然后自定义了一个应用上下文初始化器实现类，进一步加深了对 **Spring Boot** 初始化应用上下文过程的了解，为后续的启动运行过程的理解打下了坚实的基础。

