---
title: 【Spring Boot 源码学习】ApplicationContextInitializer 详解
date: 2023-12-03 13:37:52
updated: 2024-01-23 09:41:12
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - ApplicationContextInitializer
---



![](/images/spring-boot-logo.png)

# 一、引言

书接前文[《初识 SpringApplication》](../../../../../2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)，我们从 **Spring Boot** 的启动类 `SpringApplication` 上入手，了解了 `SpringApplication` 实例化过程。其中，[《BootstrapRegistryInitializer 详解》](../../../../../2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)  博文中，Huazie 已经带大家详细分析了 `BootstrapRegistryInitializer` 的加载和初始化过程，如下还有 **2.4** 和 **2.5** 这两处还未详细分析：

<!-- more -->

![](/images/springboot/loader.png)

那本篇博文就主要围绕 **2.4** 的内容展开，详细分析一下`ApplicationContextInitializer` 加载和初始化逻辑。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、主要内容
> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 2.1 初识 ApplicationContextInitializer
我们先来看看 `ApplicationContextInitializer` 接口的源码：

```java
@FunctionalInterface
public interface ApplicationContextInitializer<C extends ConfigurableApplicationContext> {
    void initialize(C applicationContext);
}
```

从上述代码，我们可以看到 `ApplicationContextInitializer` 接口被  `@FunctionalInterface` 注解修饰。

> **知识点：** `@FunctionalInterface` 是 **Java 8** 中引入的一个注解，用于标识一个函数式接口。函数式接口是只有一个抽象方法的接口，常用于实现 **Lambda** 表达式和方法引用。
> 使用 `@FunctionalInterface` 注解可以向编译器指示该接口是一个函数式接口，从而在编译时进行类型检查，确保该接口 **只包含一个抽象方法**。此外，该注解还可以为函数式接口生成特殊的方法，如默认方法（default method）和 静态方法（static method），这些方法可以在接口中提供更多的功能，这里就不赘述了，感兴趣的朋友可以自行查阅相关函数式接口的资料。

`ApplicationContextInitializer` 是个回调接口，它只包含一个 `initialize` 方法，该方法用来初始化给定的应用程序上下文，即它的唯一参数 `applicationContext`。

`ApplicationContextInitializer` 的主要用途是在 `ConfigurableApplicationContext` 类型（或其子类型）的应用程序上下文刷新之前，允许用户初始化 **Spring** `ConfigurableApplicationContext` 对象实例。通常用于需要在应用程序上下文中进行一些程序化初始化的 **Web** 应用程序。例如，注册属性源或激活与上下文环境相关的配置文件。请参阅 `ContextLoader` 和`FrameworkServlet` 支持，它们分别支持声明 `"contextInitializerClasses"` 上下文参数和初始化参数。建议使用 `ApplicationContextInitializer` 处理器检测是否实现了 **Spring** 的 `Ordered` 接口或者是否存在`@Order` 注解，并在调用之前根据这些信息对实例进行排序。

## 2.2 加载 ApplicationContextInitializer

```java
setInitializers((Collection) getSpringFactoriesInstances(ApplicationContextInitializer.class));
```

上述代码是 `SpringApplication` 的核心构造方法中的逻辑，它用于加载实现了 `ApplicationContextInitializer` 接口的类的实例集合，并将该实例集合设置到 `SpringApplication` 的 `initializers` 变量中。

```java
private List<ApplicationContextInitializer<?>> initializers;
```

我们进入 `getSpringFactoriesInstances` 方法，查看如下：

![](/images/springboot/getSpringFactoriesInstances.png)

我们看到了如下的代码 ：

```java
SpringFactoriesLoader.loadFactoryNames(type, classLoader);
```

这里是通过 `SpringFactoriesLoader` 类的 `loadFactoryNames` 方法来获取 `META-INF/spring.factories` 中配置 key 为 `org.springframework.context.ApplicationContextInitializer` 的数据；

我们以 **spring-boot-autoconfigure-2.7.9.jar** 为例：

![](ApplicationContextInitializer.png)
```bash
# Initializers
org.springframework.context.ApplicationContextInitializer=\
org.springframework.boot.autoconfigure.SharedMetadataReaderFactoryContextInitializer,\
org.springframework.boot.autoconfigure.logging.ConditionEvaluationReportLoggingListener
```
## 2.3 ApplicationContextInitializer 的初始化

这里我们需要查看 `SpringApplication` 的 `run(String... args)` 方法，如下所示：

![](prepareContext.png)

在上述的 `prepareContext` 方法中，就能看到 `ApplicationContextInitializer` 的初始化；而在 `prepareContext` 方法的下面，就是 `refreshContext` 方法，正好解释了上面说的 `ApplicationContextInitializer` 的主要用途。

我们继续往下看：

![](prepareContext1.png)

上述截图中，我们继续看 `applyInitializers` 方法：

![](applyInitializers.png)

到这步，已经很清楚了，我们来总结一下上述 `applyInitializers` 方法的逻辑：

*   通过 `getInitializers` 方法，获取了 `SpringApplication` 的 `initializers` 变量，即实现了 `ApplicationContextInitializer` 接口的集合。
    * 继续查看 `getInitializers` 方法，我们可以看到如下：
    ```
    public Set<ApplicationContextInitializer<?>> getInitializers() {
        return asUnmodifiableOrderedSet(this.initializers);
    }
    
    private static <E> Set<E> asUnmodifiableOrderedSet(Collection<E> elements) {
        List<E> list = new ArrayList<>(elements);
        list.sort(AnnotationAwareOrderComparator.INSTANCE);
        return new LinkedHashSet<>(list);
    }
    ```
    注意看 `asUnmodifiableOrderedSet` 方法中的 `AnnotationAwareOrderComparator` 类，它是 **Spring** 框架中的一个核心排序组件，其作用是对实现了 `Ordered` 接口、`PriorityOrdered`接口以及被 `@Order` 或 `@Priority` 注解修饰的类进行统一的排序。

*   遍历 `ApplicationContextInitializer` 接口的集合，循环操作 `initializer` 的初始化；
    *   通过 `GenericTypeResolver##resolveTypeArgument` 方法，来解析 `initializer` 对象中的泛型类型参数，并赋值给 `requiredType` 变量；
    
    *   通过 `Assert##isInstanceOf` 方法，来检查 `context` 对象是否是`requiredType` 类型的实例。如果不是，那么会抛出一个异常，异常信息为 **"Unable to call initializer."**；
    *   调用 `ApplicationContextInitializer` 接口的 `initialize` 方法，初始化给定的应用上下文对象 `context`。

# 三、总结

本篇 **Huazie** 带大家详细分析了 **Spring Boot** 中 `ApplicationContextInitializer`  的加载和初始化的逻辑，这对于后续的 `SpringApplication` 运行流程的理解至关重要。

