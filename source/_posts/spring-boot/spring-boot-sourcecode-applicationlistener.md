---
title: 【Spring Boot 源码学习】ApplicationListener 详解
date: 2023-12-10 22:10:55
updated: 2024-01-22 17:39:52
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - ApplicationListener
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言

书接前文[《初识 SpringApplication》](/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)，我们从 **Spring Boot** 的启动类 `SpringApplication` 上入手，了解了 `SpringApplication` 实例化过程。其中，[《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)  和 [《ApplicationContextInitializer 详解》](/2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/)博文中，Huazie 已经带大家详细分析了 `BootstrapRegistryInitializer` 和 `ApplicationContextInitializer`  的加载和初始化过程，如下还有 **2.5** 还未详细分析：
![](/images/springboot/loader.png)

那本篇博文就主要围绕 **2.5** 的内容展开，详细分析一下 `ApplicationListener ` 的加载和处理应用程序事件的逻辑。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="21" align="left" > 
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
</table>

# 主要内容
> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 1. 初识 ApplicationListener 
我们先来看看 `ApplicationListener ` 接口的源码【**spring-context-5.3.25.jar**】：

```java
@FunctionalInterface
public interface ApplicationListener<E extends ApplicationEvent> extends EventListener {

  void onApplicationEvent(E event);

  static <T> ApplicationListener<PayloadApplicationEvent<T>> forPayload(Consumer<T> consumer) {
    return event -> consumer.accept(event.getPayload());
  }
}
```

从上述代码，我们可以看到 `ApplicationListener` 接口被  `@FunctionalInterface` 注解修饰。

> **知识点：** `@FunctionalInterface` 是 **Java 8** 中引入的一个注解，用于标识一个函数式接口。函数式接口是只有一个抽象方法的接口，常用于实现 **Lambda** 表达式和方法引用。
> 使用 `@FunctionalInterface` 注解可以向编译器指示该接口是一个函数式接口，从而在编译时进行类型检查，确保该接口 **只包含一个抽象方法**。此外，该注解还可以为函数式接口生成特殊的方法，如默认方法（default method）和 静态方法（static method），这些方法可以在接口中提供更多的功能，这里就不赘述了，感兴趣的朋友可以自行查阅相关函数式接口的资料。

`ApplicationListener` 是 **Spring** 中应用程序事件监听器实现的接口。它基于观察者设计模式的`java.util.EventListener` 接口的标准。在注册到 **Spring ApplicationContext** 时，事件将进行相应的过滤，只有匹配的事件对象才会使该监听器被调用。

在 `ApplicationListener` 接口中，我们可以看到它定义了一个 `onApplicationEvent(E event)` 方法，当监听事件被触发时，`onApplicationEvent` 方法就会被调用执行。`onApplicationEvent` 方法一般用于处理应用程序事件，参数 `event` 为 `ApplicationEvent` 的子类，也就是具体要响应处理的各种类型的应用程序事件。例如，当某个特定事件发生时，你可能想要记录日志、更新数据库、发送电子邮件等等。

另外，`ApplicationListener` 接口还提供了一个静态方法 `forPayload(Consumer<T> consumer)`，用于创建一个新的 `ApplicationListener` 实例。这个方法接受一个 `Consumer<T>` 类型的参数，这个参数是一个函数接口，它接受一个泛型参数 `T`，并对其执行一些操作。通过这个方法，你可以将一个 `Consumer` 函数作为参数，然后返回一个对应的事件监听器。这个监听器会在事件发生时，调用 `Consumer` 函数处理事件的有效载荷【即事件中包含的有效信息或数据】。

## 2. 加载 ApplicationListener 

```java
setListeners((Collection) getSpringFactoriesInstances(ApplicationListener.class));
```

上述代码是 `SpringApplication` 的核心构造方法中的逻辑，它用于加载实现了 `ApplicationListener` 接口的监听器实例集合，并将该监听器实例集合设置到 `SpringApplication` 的 `listeners` 变量中。

```java
private List<ApplicationContextInitializer<?>> initializers;
```

我们进入 `getSpringFactoriesInstances` 方法，查看如下：

![](/images/springboot/getSpringFactoriesInstances.png)

我们看到了如下的代码 ：

```java
SpringFactoriesLoader.loadFactoryNames(type, classLoader);
```

这里是通过 `SpringFactoriesLoader` 类的 `loadFactoryNames` 方法来获取 `META-INF/spring.factories` 中配置 key 为 `org.springframework.context.ApplicationListener` 的数据；

我们以 **spring-boot-autoconfigure-2.7.9.jar** 为例：

![](ApplicationListener.png)

```bash
# Application Listeners
org.springframework.context.ApplicationListener=\
org.springframework.boot.autoconfigure.BackgroundPreinitializer
```
## 3. 响应应用程序事件

这里我们需要查看 `SpringApplication` 的 `run(String... args)` 方法，如下所示：

![](listeners.png)

我们看上面的 `SpringApplicationRunListeners` ，其内的 `listeners` 变量是 `SpringApplicationRunListener` 接口的集合，如下所示：

![](listeners1.png)

而 `SpringApplicationRunListener` 接口的一个实现就是 `EventPublishingRunListener` 类，该类的作用就是根据 **Spring Boot** 程序启动过程的 **不同阶段** 发布对应的事件，然后由不同的实现 `ApplicationListener` 接口的应用程序监听器，来处理对应的事件【有关 `SpringApplicationRunListener` 监听器的内容，我们后续博文中会详细介绍，这里不展开了】。


如下图是 `SpringApplicationRunListeners` 类中的方法，它们分别对应了 **Spring Boot** 程序启动过程中要发布的**不同阶段**的事件的逻辑。

![](SpringApplicationRunListeners.png)
- `starting` ：当 `run` 方法第一次被执行时，该方法会立即被调用，可用于非常早期的初始化工作
- `environmentPrepared` ：当 `environment` 准备完成，在 `ApplicationContext` 创建之前，该方法被调用
- `contextPrepared` ：当 `ApplicationContext` 构建完成，资源还未被加载时，该方法被调用
- `contextLoaded` ：当 `ApplicationContext` 加载完成，未被刷新之前，该方法被调用
- `started` ：当 `ApplicationContext` 刷新并启动之后，`CommandLineRunner` 和 `ApplicationRunner` 未被调用之前，该方法被调用
- `ready` ：当所有准备工作就绪，`run` 方法执行完成之前，该方法被调用
- `failed` ：当应用程序出现错误时，该方法被调用


我们以 `starting` 方法的逻辑为例，看一下 `ApplicationStartingEvent` 事件发布并被处理的过程。

```java
void starting(ConfigurableBootstrapContext bootstrapContext, Class<?> mainApplicationClass) {
  doWithListeners("spring.boot.application.starting", (listener) -> listener.starting(bootstrapContext),
      (step) -> {
        if (mainApplicationClass != null) {
          step.tag("mainApplicationClass", mainApplicationClass.getName());
        }
      });
}
```

我们继续看 `doWithListeners` 方法：

![](doWithListeners.png)

结合上面的截图，我们重点看下这行：

```java
(listener) -> listener.starting(bootstrapContext)
```

这里时调用了 `SpringApplicationRunListener` 接口的 `starting` 方法：

![](starting.png)

这里的 `multicastEvent` 方法就是用来发布一个指定的应用程序事件，比如这里发布的就是 `ApplicationStartingEvent` 事件。

![](multicastEvent.png)
![](invokeListener.png)



# 总结
本篇 **Huazie** 带大家详细分析了 `ApplicationListener ` 的加载和处理应用程序事件，这对于后续的 `SpringApplication` 运行流程的理解至关重要。

