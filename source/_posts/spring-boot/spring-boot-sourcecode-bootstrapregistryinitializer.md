---
title: 【Spring Boot 源码学习】 BootstrapRegistryInitializer 详解
date: 2023-11-30 07:00:00
updated: 2024-01-30 00:07:58
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - BootstrapRegistryInitializer
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言

书接前文[《初识 SpringApplication》](/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)，我们从 **Spring Boot** 的启动类 `SpringApplication` 上入手，了解了 `SpringApplication` 实例化过程。其中，有如下三块内容还未详细分析：
![](/images/springboot/loader.png)

本篇博文就主要围绕 **2.3** 的内容展开，详细分析一下加载并初始化 `BootstrapRegistryInitializer` 及其相关的类的逻辑。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="19" align="left" > 
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
</table>

# 三、主要内容
> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 3.1 初识 BootstrapRegistryInitializer
废话不多说，我们直接来看 `BootstrapRegistryInitializer` 接口的源码：

```java
@FunctionalInterface
public interface BootstrapRegistryInitializer {
  void initialize(BootstrapRegistry registry);
}
```

阅读上述代码，可以看到 `BootstrapRegistryInitializer` 接口被  `@FunctionalInterface` 注解修饰。

> `@FunctionalInterface` 是 **Java 8** 中引入的一个注解，用于标识一个函数式接口。函数式接口是只有一个抽象方法的接口，常用于实现 **Lambda** 表达式和方法引用。
> 使用 `@FunctionalInterface` 注解可以向编译器指示该接口是一个函数式接口，从而在编译时进行类型检查，确保该接口 **只包含一个抽象方法**。此外，该注解还可以为函数式接口生成特殊的方法，如默认方法（default method）和 静态方法（static method），这些方法可以在接口中提供更多的功能，这里就不赘述了，感兴趣的朋友可以自行查阅相关函数式接口的资料。

 `BootstrapRegistryInitializer` 接口只定义了一个 `initialize` 方法，该方法只有一个参数是 `BootstrapRegistry`；

`BootstrapRegistry` 是一个简单的对象注册表，它在启动和环境后处理期间都可用，直到`ApplicationContext` 准备好为止。它可用于注册可能创建成本较高或在 `ApplicationContext` 可用之前需要共享的实例。它的一个默认实现是 `DefaultBootstrapContext` ，后面我们会看到。

注册表使用 `Class` 作为键，这意味着每个给定类型只能存储一个实例。

`addCloseListener(ApplicationListener)` 方法可用于添加监听器，当 `BootstrapContext` 已关闭并且 `ApplicationContext` 完全准备好时，该监听器可以执行操作。例如，一个实例可能选择将自己注册为常规的 **Spring bean**，以便应用程序可以使用它。

> 简而言之，`BootstrapRegistry` 是一个用于存储和共享对象的注册表，这些对象在`ApplicationContext` 准备好之前就可能已经被创建并需要被共享。

在 **Spring Cloud Config** 中，客户端通过向配置中心（**Config Server**）发送请求来获取应用程序的配置信息。而 `BootstrapRegistryInitializer` 就是负责将配置中心的相关信息注册到 **Spring** 容器中的。
## 3.2 加载 BootstrapRegistryInitializer
```java
this.bootstrapRegistryInitializers = new ArrayList<>(
        getSpringFactoriesInstances(BootstrapRegistryInitializer.class));
```

上述代码是 `SpringApplication` 的核心构造方法中的逻辑，它用于加载实现了 `BootstrapRegistryInitializer` 接口的类的实例集合。

我们进入 `getSpringFactoriesInstances` 方法，查看如下：

![](/images/springboot/getSpringFactoriesInstances.png)

我们看到了如下的代码 ：

```java
SpringFactoriesLoader.loadFactoryNames(type, classLoader);
```

这里是通过 `SpringFactoriesLoader` 类的 `loadFactoryNames` 方法来获取 `META-INF/spring.factories` 中配置 key 为 `org.springframework.boot.BootstrapRegistryInitializer` 的数据；

当然这些配置不在 **Spring Boot** 的 `META-INF/spring.factories` 中，我们上面提到  **Spring Cloud Config** 就是用 `BootstrapRegistryInitializer` 将配置中心的相关信息注册到 **Spring** 容器中，那我们就来看看 **Spring Cloud Config** 相关的配置。

如下是 **Spring Cloud Config** 的 `Starter` 依赖：

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-config</artifactId>
    <version>4.0.4</version>
</dependency>
```

导入上述依赖之后，相关的 `META-INF/spring.factories` 配置，我们发现是在如下的 jar 包里面【**spring-cloud-config-client-4.0.4.jar**】：

![](spring-cloud-config-client.png)

查看 `META-INF/spring.factories` 配置文件，我们可以看到如下：

![](spring-cloud-config-client-spring-factories.png)
```bash
# Spring Boot BootstrapRegistryInitializers
org.springframework.boot.BootstrapRegistryInitializer=\
org.springframework.cloud.config.client.ConfigClientRetryBootstrapper
```

有关 **Spring Cloud Config** 的内容，这里就不展开介绍了，感兴趣的小伙伴自行查阅 [**Spring Cloud Config** 的官方文档](https://docs.spring.io/spring-cloud-config/docs/current/reference/html/)。

## 3.3 BootstrapRegistryInitializer 的初始化

这里我们需要查看 `SpringApplication` 的 `run(String... args)` 方法，如下所示：

![](createBootstrapContext.png)

在上述的 `createBootstrapContext` 方法中，就对 `BootstrapRegistryInitializer` 进行初始化，我们继续往下看：

![](createBootstrapContext1.png)

从上图中，我们可以看到这样一行代码：

```java
this.bootstrapRegistryInitializers.forEach((initializer) -> initializer.initialize(bootstrapContext));
```

这里涉及如下的知识点：
- `this.bootstrapRegistryInitializers.forEach()` : **Java 8** 的 **Stream API**，它用于遍历列表中的每个元素，并执行给定的操作【即 `initializer.initialize(bootstrapContext)` 】。
- `(initializer) -> initializer.initialize(bootstrapContext)`  ： **Lambda** 表达式，这是 **Java 8** 引入的一个新特性，允许以更简洁的方式表示匿名方法。它表示一个接受`BootstrapRegistryInitializer` 类型参数 `initializer`，并调用其 `initialize(bootstrapContext)` 方法的功能。

简而言之，对于 `this.bootstrapRegistryInitializers` 列表中的每个 `BootstrapRegistryInitializer`，使用当前的 `bootstrapContext` 初始化它。这里的 `bootstrapContext` 其实就是 `BootstrapRegistry` 注册表的一个默认实现 `DefaultBootstrapContext` 。

从上述的  `SpringApplication` 的 `run(String... args)` 方法源码中，我们也可以看出 `BootstrapRegistryInitializer` 的初始化是在 **Spring Boot** 应用启动一开始进行的。

我们通过实现 `BootstrapRegistryInitializer` 接口并定义 `initialize` 方法，可以将自定义的 **Bean** 初始化器注册到 `ApplicationContext` 中。这样，在 **Spring Boot** 应用启动时，这些初始化器会被自动加载并执行，从而完成一些必要的初始化配置。
# 四、总结
本篇 **Huazie** 带大家详细分析了加载并初始化 `BootstrapRegistryInitializer` 的逻辑，这对于后续的 `SpringApplication` 运行流程的理解至关重要。

