---
title: 【Spring Boot 源码学习】BootstrapRegistry 初始化器实现
date: 2024-03-02 12:35:28
updated: 2024-03-02 12:35:28
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - BootstrapRegistry
  - BootstrapContext
  - BootstrapRegistryInitializer 
---



![](/images/spring-boot-logo.png)

# 一、引言
前面的博文[《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)，Huazie 带大家一起详细分析了 **Spring Boot** 启动时加载并初始化 `BootstrapRegistryInitializer` 及其相关的类的逻辑。本篇就让我们自定义 `BootstrapRegistryInitializer` 接口实现，以此来执行自定义的初始化操作【如注册自定义的 **Bean**、添加 **BootstrapContext** 关闭监听器】。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="25" align="left" > 
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
</table>

# 三、主要内容
> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 3.1 BootstrapRegistry
在[《BootstrapRegistry 详解》](/2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/) 中，**Huazie** 详细介绍了 `BootstrapRegistry` 的源码，这有助于下面介绍的 `BootstrapRegistry` 初始化器的实现逻辑，有不知道的朋友们直接查看即可，这里不再赘述。
## 3.2 BootstrapRegistryInitializer
在[《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 中，**Huazie** 详解分析了加载并初始化 `BootstrapRegistryInitializer` 的逻辑，这同样有助于
理解下面将要讲解的内容，还不熟悉的朋友们赶紧花点时间了解下，这里不再赘述。
## 3.3 BootstrapRegistry 初始化器实现

### 3.3.1 定义 DemoBootstrapper

下面我们来定义一个类 `DemoBootstrapper`，该类实现 `BootstrapRegistryInitializer` 接口，如下：

```java
public class DemoBootstrapper implements BootstrapRegistryInitializer {
    @Override
    public void initialize(BootstrapRegistry registry) {
        // 注册一些自定义的对象 或者 加载自定义的一些配置
    }
}
```

在[《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 的 **3.2** 小节，**Huazie** 介绍了 `BootstrapRegistryInitializer` 的加载过程，上述我们自定义的 `DemoBootstrapper`  也会在 **Spring Boot** 启动引导阶段进行加载并初始化。

上述定义中，我们只是展示了一个空实现的类，其中的 `initialize` 方法还未做处理。

至于 `initialize` 方法中该添加哪些逻辑，这就要看它的参数 `BootstrapRegistry` 接口了。

下面代码，Huazie 演示了如何 **注册自定义的对象**，以及添加 **引导上下文关闭事件监听器**。

```java
// 注册 User对象，它就是一个简单的 POJO 类，含两个成员变量：名称 name 和年龄 age
registry.register(User.class, context -> new User("Huazie", 18));
// 添加 BootstrapContext关闭 监听器
registry.addCloseListener(new DemoBootstrapContextClosedListener());
```

`DemoBootstrapContextClosedListener` ，即引导上下文关闭事件监听器，相关演示代码如下：

```java
public class DemoBootstrapContextClosedListener implements ApplicationListener<BootstrapContextClosedEvent>, Ordered {
    @Override
    public void onApplicationEvent(BootstrapContextClosedEvent event) {
        BootstrapContext bootstrapContext = event.getBootstrapContext();
        if (bootstrapContext.isRegistered(User.class)) {
            System.out.println("BootstrapContext关闭时获取User：" + bootstrapContext.get(User.class));
        }
    }

    @Override
    public int getOrder() {
        return 1;
    }
}
```
虽然上面添加了引导上下文关闭事件监听器，但是我们还不知道什么时候它会被执行。

在 `DemoBootstrapContextClosedListener` 中，我们看到了 `BootstrapContext`  的使用，显然这里涉及到了引导上下文的实际使用场景，由于篇幅受限，将在下篇介绍，大家不妨期待一下。

另外，在 `DemoBootstrapContextClosedListener` 中，还看到它实现了 `Ordered` 接口【**spring-core** 包中的接口】。

**那么在事件监听器中，这个 `Ordered` 接口是用来做什么的呢？**

在回答这个问题之前，我们先来看看 `Ordered` 接口的源码：

```java
public interface Ordered {
    int HIGHEST_PRECEDENCE = Integer.MIN_VALUE;
    int LOWEST_PRECEDENCE = Integer.MAX_VALUE;
    int getOrder();
}
```
`Ordered` 接口定义了两个常量和一个方法：
- `HIGHEST_PRECEDENCE`：最高优先级值的有用常数【最小的 `Integer` 值】
- `LOWEST_PRECEDENCE` ：最低优先级值的有用常数【最大的 `Integer` 值】
- `int getOrder()` ：获取当前对象的优先级值【值越小，优先级越高】

源码中对于 `Ordered` 接口是这样说的：

它可以被需要排序的对象实现，例如在集合中。实际的排序可以被解读为优先级排序，其中第一个对象（即有着最低的排序值）具有最高的优先级。

当然，`Ordered` 接口还有一个扩展接口，即优先级标记接口 `PriorityOrdered`。`PriorityOrdered` 对象总是优先于普通 `Ordered` 对象，无论它们的排序值如何。当对一组 `Ordered` 对象进行排序时，`PriorityOrdered` 对象和普通 `Ordered` 对象实际上被视为两个独立的子集，`PriorityOrdered` 对象子集先于普通 `Ordered` 对象子集，并在这些子集内部应用相对排序。

上述排序逻辑请查看 **spring-core** 包中的 `AnnotationAwareOrderComparator` 类 和 `OrderComparator` 类，这里不再赘述了。

现在可以回答上面的问题了：在事件监听器中实现 `Ordered` 接口，可以用来确保 **多个监听同一事件的监听器** 可以按照我们 **预定的顺序执行**。

### 3.3.2 添加 DemoBootstrapper
不过，要想能够加载到自定义的 `DemoBootstrapper` ，我们还需要将它添加到 `bootstrapRegistryInitializers` 中才可以。

```java
// SpringApplication的私有变量
private List<BootstrapRegistryInitializer> bootstrapRegistryInitializers;
```

**那么，我们该如何添加呢？**

通过阅读 `SpringApplication` 的源码，可以总结如下的两种方式：

-  在 `META-INF/spring.factories` 中添加 `org.springframework.boot.BootstrapRegistryInitializer` 的配置。这种方式，我们从 [《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 的 **3.2** 小节可见一斑。
  ```java
  org.springframework.boot.BootstrapRegistryInitializer=com.example.demo.DemoBootstrapper
  ```

- 通过 `SpringApplication` 中的 `addBootstrapRegistryInitializer` 方法添加。其实这里在笔者的[《SpringApplication 的定制化介绍》](/2024/01/07/spring-boot/spring-boot-sourcecode-springapplication-customization/)中的 **1.5** 小节也提及过。
  ```java
  SpringApplication springApplication = new SpringApplication(DemoApplication.class);
  springApplication.addBootstrapRegistryInitializer(new DemoBootstrapper());
  // 其他省略。。。
  ```

# 四、总结
本篇 **Huazie** 介绍了如何自定义 `BootstrapRegistry` 初始化器实现，其中演示如何在**引导上下文**中注册了**自定义的对象**以及如何在**引导上下文**中添加**引导上下文关闭事件监听器**。