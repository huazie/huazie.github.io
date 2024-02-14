---
title: 【Spring Boot 源码学习】 BootstrapRegistry 详解
date: 2024-01-31 23:32:34
updated: 2024-02-01 10:11:18
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - BootstrapRegistry
  - BootstrapContext
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
前面的博文[《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)，Huazie 带大家一起详细分析了 **Spring Boot** 启动时加载并初始化 `BootstrapRegistryInitializer` 及其相关的类的逻辑。其中有个 `BootstrapRegistry` 接口只是简单提及，本篇就详细分析一下 `BootstrapRegistry` 接口，这对于我们后续理解 《`BootstrapRegistry` 初始化器实现》的内容至关重要。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="23" align="left" > 
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
</table>

# 三、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

在 [《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 的 3.1 小节，我们对 `BootstrapRegistry` 进行了初步的介绍：它是一个用于存储和共享对象的注册表，这些对象在 `ApplicationContext` 准备好之前就可能已经被创建并需要被共享。
## 3.1 源码初识
首先让我们来看看 `BootstrapRegistry` 的源码：

```java
public interface BootstrapRegistry {

  <T> void register(Class<T> type, InstanceSupplier<T> instanceSupplier);

  <T> void registerIfAbsent(Class<T> type, InstanceSupplier<T> instanceSupplier);

  <T> boolean isRegistered(Class<T> type);

  <T> InstanceSupplier<T> getRegisteredInstanceSupplier(Class<T> type);

  void addCloseListener(ApplicationListener<BootstrapContextClosedEvent> listener);

  @FunctionalInterface
  interface InstanceSupplier<T> {

    T get(BootstrapContext context);

    default Scope getScope() {
      return Scope.SINGLETON;
    }

    default InstanceSupplier<T> withScope(Scope scope) {
      // 。。。
    }

    static <T> InstanceSupplier<T> of(T instance) {
      return (registry) -> instance;
    }

    static <T> InstanceSupplier<T> from(Supplier<T> supplier) {
      return (registry) -> (supplier != null) ? supplier.get() : null;
    }

  }

  enum Scope {
    SINGLETON,
    PROTOTYPE
  }

}

```

它包含了 5 个方法，1个内部接口类，1个内部枚举类，下面我们一一来介绍下：
##  3.2 register 方法
`register` 方法，包含两个参数：
- `Class<T> type` ：实例类型
- `InstanceSupplier<T> instanceSupplier` ：实例供应者

该方法用于将特定类型注册到注册表中。如果指定的类型已经被注册，并且尚未作为单例获取，那么它将被替换。
## 3.3 registerIfAbsent 方法
`registerIfAbsent` 方法，包含两个参数：
- `Class<T> type` ：实例类型
- `InstanceSupplier<T> instanceSupplier` ：实例供应者

如果尚未存在特定类型的注册，则向注册表中注册该类型。
## 3.4 isRegistered 方法
`isRegistered` 方法，只有一个参数：
- `Class<T> type` ：实例类型

该方法用于返回给定类型是否已注册。如果给定类型已经注册，则返回 `true`，否则，返回 `false`。

## 3.5 getRegisteredInstanceSupplier 方法
`getRegisteredInstanceSupplier` 方法，也只有一个参数：

- `Class<T> type` ：实例类型

该方法返回给定类型的任何现有的 `BootstrapRegistry.InstanceSupplier`。
## 3.6 addCloseListener 方法
`addCloseListener` 方法，只有一个参数：

- `ApplicationListener<BootstrapContextClosedEvent> listener` ：待添加的监听器

该方法用于添加一个 `ApplicationListener`，当 `BootstrapContext` 关闭并且 `ApplicationContext` 已经准备就绪时，该监听器将与 `BootstrapContextClosedEvent` 一起被调用。
## 3.7 InstanceSupplier 内部接口类
`InstanceSupplier` 内部接口类是用于提供实际实例的供应者。

它定义了一个 1 个普通方法，2 个默认方法，2 个静态方法。

**知识拓展：**

从 **Java 8** 开始，支持在接口中定义默认方法和静态方法。
- 默认方法（**Default Method**）允许你在接口中添加一个有默认实现的非抽象方法。这使得接口可以更加灵活地扩展，而不需要破坏与现有代码的兼容性。默认方法使用关键字 default 进行声明，并提供了具体的实现。
- 静态方法（**Static Method**）允许你在接口中定义一个静态方法，该方法可以在不创建接口实例的情况下调用。静态方法使用关键字 static 进行声明，并可以直接通过接口名来调用。

### 3.7.1 get 方法
`get` 方法，只包含一个参数：

- `BootstrapContext context` ：BootstrapContext 是一个用于获取其他引导实例的上下文

该方法是工厂方法，用于在需要时创建实例，后续我们在讲解 `DefaultBootstrapContext` 时也会涉及。

### 3.7.2 getScope 默认方法
`getScope` 默认方法，用于返回提供的实例的作用域；如果该方法没有被重写，则默认返回 `Scope.SINGLETON` 。

### 3.7.3 withScope 默认方法

```java
default InstanceSupplier<T> withScope(Scope scope) {
  Assert.notNull(scope, "Scope must not be null");
  InstanceSupplier<T> parent = this;
  return new InstanceSupplier<T>() {
    @Override
    public T get(BootstrapContext context) {
      return parent.get(context);
    }

    @Override
    public Scope getScope() {
      return scope;
    }
  };
}
```
通过阅读上述代码可知，该方法根据其参数 `scope` ，返回一个指定作用域的新的 `BootstrapRegistry.InstanceSupplier` 的匿名对象，该匿名对象重写了 `get` 和 `getScope` 方法。这里使用匿名对象的好处就是可以在不定义新类的情况下快速地创建一个具有特定行为的对象。

细心的读者们，可能发现了匿名对象的 `get` 方法中，使用了 `withScope` 方法中定义的变量 `parent`，它被用来存储当前对象的引用 `this`。

**那么这里的 `parent` ，能不能直接替换成 `this` 呢？**

显然是不可以的，`this` 关键字用在匿名内部类中，指代的是该匿名内部类本身，而不是外部对象。而匿名对象这里的重写的 `get` 方法，实际上需要调用 `withScope` 方法所在的对象的 `get` 方法来实现功能。如果这里用 `this`，实际上就是自己调自己，一直无限递归调用，最终导致栈溢出错误。

### 3.7.4 of 静态方法

该静态方法是一个工厂方法，用于为给定实例创建一个`BootstrapRegistry.InstanceSupplier`。

```java
return (registry) -> instance;
```

这里采用了 **Java 8** 的 `Lambda` 表达式，也相当于如下的写法：

```java
return new InstanceSupplier<T>() {
    @Override
    public T get(BootstrapContext registry) {
        return instance;
    }
};
```
### 3.7.5 from 静态方法
该静态方法也是一个工厂方法，用于通过一个 `Supplier` 创建`BootstrapRegistry.InstanceSupplier`。

```java
return (registry) -> (supplier != null) ? supplier.get() : null;
```
这里也是用了 **Java 8** 的 `Lambda` 表达式，相当于如下的写法：

```java
return new InstanceSupplier<T>() {
    @Override
    public T get(BootstrapContext registry) {
        return (supplier != null) ? supplier.get() : null;
    }
};
```

> **知识点：**`Supplier` 是 **Java 8** 开始引入，作为 `java.util.function` 包的一部分，它与 `Lambda` 表达式一起被引入，以支持函数式编程范式。该接口是为了简化无参数方法的表示，特别是在需要延迟执行或创建对象时。`Supplier` 接口只有一个抽象方法 `get`()，它不接受任何参数，但返回一个通用类型的值。 

## 3.8 Scope 内部枚举类

Scope 表示一个实例的作用域，它只包含两个枚举变量，分别是：

- `SINGLETON` ：单例实例。`InstanceSupplier` 只会被调用一次，并且每次调用都会返回相同的实例。
- `PROTOTYPE` ：原型实例。`InstanceSupplier` 将在需要实例时被调用。


# 四、总结
本篇 **Huazie** 带大家通读了 `BootstrapRegistry` 的相关源码，这些内容对于后面的源码学习至关重要。

