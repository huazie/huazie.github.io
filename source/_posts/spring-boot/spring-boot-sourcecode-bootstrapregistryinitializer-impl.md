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
前面的博文[《BootstrapRegistryInitializer 详解》](../../../../../2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)，Huazie 带大家一起详细分析了 **Spring Boot** 启动时加载并初始化 `BootstrapRegistryInitializer` 及其相关的类的逻辑。本篇就让我们自定义 `BootstrapRegistryInitializer` 接口实现，以此来执行自定义的初始化操作【如注册自定义的 **Bean**、添加 **BootstrapContext** 关闭监听器】。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、主要内容
> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 2.1 BootstrapRegistry
在[《BootstrapRegistry 详解》](../../../../../2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/) 中，**Huazie** 详细介绍了 `BootstrapRegistry` 的源码，这有助于下面介绍的 `BootstrapRegistry` 初始化器的实现逻辑，有不知道的朋友们直接查看即可，这里不再赘述。
## 2.2 BootstrapRegistryInitializer
在[《BootstrapRegistryInitializer 详解》](../../../../../2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 中，**Huazie** 详解分析了加载并初始化 `BootstrapRegistryInitializer` 的逻辑，这同样有助于
理解下面将要讲解的内容，还不熟悉的朋友们赶紧花点时间了解下，这里不再赘述。
## 2.3 BootstrapRegistry 初始化器实现

### 2.3.1 定义 DemoBootstrapper

下面我们来定义一个类 `DemoBootstrapper`，该类实现 `BootstrapRegistryInitializer` 接口，如下：

```java
public class DemoBootstrapper implements BootstrapRegistryInitializer {
    @Override
    public void initialize(BootstrapRegistry registry) {
        // 注册一些自定义的对象 或者 加载自定义的一些配置
    }
}
```

在[《BootstrapRegistryInitializer 详解》](../../../../../2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 的 **3.2** 小节，**Huazie** 介绍了 `BootstrapRegistryInitializer` 的加载过程，上述我们自定义的 `DemoBootstrapper`  也会在 **Spring Boot** 启动引导阶段进行加载并初始化。

上述定义中，我们只是展示了一个空实现的类，其中的 `initialize` 方法还未做处理。

至于 `initialize` 方法中该添加哪些逻辑，这就要看它的参数 `BootstrapRegistry` 接口了。

下面代码，**Huazie** 演示了如何 **注册自定义的对象**，以及添加 **引导上下文关闭事件监听器**。

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

### 2.3.2 添加 DemoBootstrapper
不过，要想能够加载到自定义的 `DemoBootstrapper` ，我们还需要将它添加到 `bootstrapRegistryInitializers` 中才可以。

```java
// SpringApplication的私有变量
private List<BootstrapRegistryInitializer> bootstrapRegistryInitializers;
```

**那么，我们该如何添加呢？**

通过阅读 `SpringApplication` 的源码，可以总结如下的两种方式：

-  在 `META-INF/spring.factories` 中添加 `org.springframework.boot.BootstrapRegistryInitializer` 的配置。这种方式，我们从 [《BootstrapRegistryInitializer 详解》](../../../../../2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 的 **3.2** 小节可见一斑。
  ```java
  org.springframework.boot.BootstrapRegistryInitializer=com.example.demo.DemoBootstrapper
  ```

- 通过 `SpringApplication` 中的 `addBootstrapRegistryInitializer` 方法添加。其实这里在笔者的[《SpringApplication 的定制化介绍》](../../../../../2024/01/07/spring-boot/spring-boot-sourcecode-springapplication-customization/)中的 **1.5** 小节也提及过。
  ```java
  SpringApplication springApplication = new SpringApplication(DemoApplication.class);
  springApplication.addBootstrapRegistryInitializer(new DemoBootstrapper());
  // 其他省略。。。
  ```

# 三、总结
本篇 **Huazie** 介绍了如何自定义 `BootstrapRegistry` 初始化器实现，其中演示如何在**引导上下文**中注册了**自定义的对象**以及如何在**引导上下文**中添加**引导上下文关闭事件监听器**。