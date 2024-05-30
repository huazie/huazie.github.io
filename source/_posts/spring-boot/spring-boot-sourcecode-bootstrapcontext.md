---
title: 【Spring Boot 源码学习】深入 BootstrapContext 及其默认实现
date: 2024-02-25 16:41:42
updated: 2024-02-25 16:41:42
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - BootstrapContext
  - ConfigurableBootstrapContext
  - DefaultBootstrapContext
  - 
---



![](/images/spring-boot-logo.png)

# 一、引言
书接前文[《BootstrapRegistry 详解》](/2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/)，在介绍 `BootstrapRegistry`  的内部类 `InstanceSupplier` 的 `get` 方法时，看到了它的唯一参数 `BootstrapContext` 接口【即引导上下文】。而这个接口及其默认实现就是本篇要重点介绍的对象，且听我娓娓道来。


<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="24" align="left" > 
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
</table>

# 三、主要内容

> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 3.1 BootstrapContext

### 3.1.1 源码初识

```java
public interface BootstrapContext {
    <T> T get(Class<T> type) throws IllegalStateException;
    <T> T getOrElse(Class<T> type, T other);
    <T> T getOrElseSupply(Class<T> type, Supplier<T> other);
    <T, X extends Throwable> T getOrElseThrow(Class<T> type, Supplier<? extends X> exceptionSupplier) throws X;
    <T> boolean isRegistered(Class<T> type);
}
```

`BootstrapContext` 是一个简单的引导上下文，它在启动期间以及环境后处理过程中可用，直到应用上下文 `ApplicationContext` 准备就绪。

它提供了对可能创建成本高昂的单例的延迟访问，或者在 `ApplicationContext` 可用之前需要共享的单例。

它一共包含 **5** 个方法，下面分别来介绍下：

### 3.1.2 get 方法

`get` 方法，只有一个参数：

- `Class<T> type` ：实例类型

该方法用于返回一个指定类型的实例对象。如果类型已在上下文中注册，则从上下文中返回一个实例。如果之前未访问过该实例，则会创建它。

有关其具体实现，可查看 **3.3** 小节【`BootstrapContext` 的默认实现】

### 3.1.3 getOrElse 方法
`getOrElse` 方法，包含两个参数：

- `Class<T> type` ：实例类型
- `T other` ：如果上述类型还未注册，则使用该实例进行返回

该方法用于返回一个指定类型的实例对象。如果类型已在上下文中注册，则从上下文中返回一个实例。如果之前未注册过该实例，则直接用第二个参数 other 进行返回【**这里跟 `get` 方法有所区别**】。

有关其具体实现，可查看 **3.3** 小节【`BootstrapContext` 的默认实现】
### 3.1.4 getOrElseSupply 方法
`getOrElseSupply` 方法，也包含两个参数：

- `Class<T> type` ：实例类型
- `Supplier<T> other` ：如果上述类型还未注册，则使用该提供者返回指定实例对象

该方法用于返回一个指定类型的实例对象。如果类型已在上下文中注册，则从上下文中返回一个实例。如果之前未注册过该实例，则用 `other.get()` 进行返回【**这里类似 `getOrElse` 方法，其实默认实现中 `getOrElse` 就是调用 `getOrElseSupply` 进行返回的**】。

有关其具体实现，可查看 **3.3** 小节【`BootstrapContext` 的默认实现】

### 3.1.5 getOrElseThrow 方法
`getOrElseThrow` 方法，同样也包含两个参数：

- `Class<T> type` ：实例类型
- `Supplier<? extends X> exceptionSupplier` ：如果上述类型还未注册，则使用该提供者抛出指定的异常

`X` 是 `Throwable` 的子类，如果上述类型还未注册过，则将抛出 `X` 或者 `X` 的子类。

该方法用于返回一个指定类型的实例对象。如果类型已在上下文中注册，则从上下文中返回一个实例。如果之前未注册过该实例，则通过 `throw exceptionSupplier.get()` 将指定异常抛出【这个在 默认实现 `DefaultBootstrapContext` 中即可看到】。
### 3.1.6 isRegistered 方法
`isRegistered` 方法，只有一个参数：

- `Class<T> type` ：实例类型

该方法用于判断指定的类型是否已经被注册过。如果已经在上下文中注册过了，则返回 true；否则，返回false。

有关其具体实现，可查看 **3.3** 小节【`BootstrapContext` 的默认实现】

## 3.2 ConfigurableBootstrapContext

```java
public interface ConfigurableBootstrapContext extends BootstrapRegistry, BootstrapContext {

}
```
通过阅读 `ConfigurableBootstrapContext` 源码，我们可以看到它继承了 `BootstrapRegistry` 和 `BootstrapContext` 接口。这也就意味着 `ConfigurableBootstrapContext` 接口同时拥有了这两者的所有功能，即它是一个可配置的引导上下文。

对于开发人员来讲，只需要实现这个接口，并编写相应实现代码，就可以来配置和管理应用程序的引导过程。当然 **Spring Boot** 显然已经帮我们考虑了，这也就是下面 **Huazie** 将要介绍的引导上下文的默认实现 `DefaultBootstrapContext`。

## 3.3 DefaultBootstrapContext

在 [《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/) 的 3.1 小节，我们提到了 `BootstrapRegistry` 的一个默认实现 `DefaultBootstrapContext` ，下面我们就来深入分析一下。

### 3.3.1 源码初识
话不多说，直接翻看对应的源码：

```java
public class DefaultBootstrapContext implements ConfigurableBootstrapContext {

    private final Map<Class<?>, InstanceSupplier<?>> instanceSuppliers = new HashMap<>();

    private final Map<Class<?>, Object> instances = new HashMap<>();

    private final ApplicationEventMulticaster events = new SimpleApplicationEventMulticaster();

    // 实现 BootstrapRegistry 接口中的方法

    // 实现 BootstrapContext 接口中的方法

    public void close(ConfigurableApplicationContext applicationContext) {
        this.events.multicastEvent(new BootstrapContextClosedEvent(this, applicationContext));
    }
}
```

上述源码中，**Huazie** 省略了有关实现 `BootstrapRegistry` 或 `BootstrapContext` 接口中的方法，这些内容将会在下面的小节详细深入分析。

我们从类开头，一下子就能看到三个私有的，不可变的成员变量：

- `Map<Class<?>, InstanceSupplier<?>> instanceSuppliers` ：这是个 `HashMap`，它的键是 `Class` 对象，值是 `InstanceSupplier` 对象【它是 `BootstrapRegistry` 中一个内部接口类，用于提供实际的实例，具体内容可以查看 **Huazie** 的上一篇博文 】
- `Map<Class<?>, Object> instances` ：这同样也是个 `HashMap`，它的键是 `Class` 对象，值是 `Object` 对象【即对应的实际实例对象】
- `ApplicationEventMulticaster events` ：它是 `Spring` 框架中的一个组件，用于管理和广播应用程序事件。`SimpleApplicationEventMulticaster` 是其一个简单的实现。

> **注意：** `SimpleApplicationEventMulticaster` 会将所有的事件广播给所有已注册的监听器，而由监听器自行决定忽略它们不感兴趣的事件。监听器通常会在传入的事件对象上进行相应的 `instanceof` 检查。
默认情况下，所有的监听器都在调用线程中被调用。这允许存在一个恶意的监听器阻塞整个应用程序的风险，但增加了最小的开销。如果指定了替代的任务执行器，可以让监听器在不同的线程中执行，例如来自一个线程池。
### 3.3.2 实现 BootstrapRegistry 接口中的方法
在 [《BootstrapRegistry 详解》](/2024/01/31/spring-boot/spring-boot-sourcecode-bootstrapregistry/)中，我们已经了解相关的 5 个方法，下面直接看 `DefaultBootstrapContext` 中的实现： 

```java
    @Override
    public <T> void register(Class<T> type, InstanceSupplier<T> instanceSupplier) {
        register(type, instanceSupplier, true);
    }

    @Override
    public <T> void registerIfAbsent(Class<T> type, InstanceSupplier<T> instanceSupplier) {
        register(type, instanceSupplier, false);
    }

    private <T> void register(Class<T> type, InstanceSupplier<T> instanceSupplier, boolean replaceExisting) {
        Assert.notNull(type, "Type must not be null");
        Assert.notNull(instanceSupplier, "InstanceSupplier must not be null");
        synchronized (this.instanceSuppliers) {
            boolean alreadyRegistered = this.instanceSuppliers.containsKey(type);
            if (replaceExisting || !alreadyRegistered) {
                Assert.state(!this.instances.containsKey(type), () -> type.getName() + " has already been created");
                this.instanceSuppliers.put(type, instanceSupplier);
            }
        }
    }

    @Override
    public <T> boolean isRegistered(Class<T> type) {
        synchronized (this.instanceSuppliers) {
            return this.instanceSuppliers.containsKey(type);
        }
    }

    @Override
    @SuppressWarnings("unchecked")
    public <T> InstanceSupplier<T> getRegisteredInstanceSupplier(Class<T> type) {
        synchronized (this.instanceSuppliers) {
            return (InstanceSupplier<T>) this.instanceSuppliers.get(type);
        }
    }

    @Override
    public void addCloseListener(ApplicationListener<BootstrapContextClosedEvent> listener) {
        this.events.addApplicationListener(listener);
    }
```

翻看上述源码，我们可以看到除了 `addCloseListener` 方法，其他方法中都使用 `synchronized` 关键字了，而这里同步的对象就是上面提到的 `instanceSuppliers`。因为 `instanceSuppliers` 是 `HashMap`，它并不是线程安全的，为了防止多个线程同时修改 `instanceSuppliers` 对象，导致数据不一致的问题，这里就需要对该对象进行同步，保证在同一时刻只有一个线程可以访问该代码块。

从源码中，我们可以看出 `register` 和 `registerIfAbsent` 方法的区别：

`registerIfAbsent` 只会在该类型尚未注册过时，才注册该类型。而 `register` 即使该类型已经注册过了，也会重新注册该类型。

`isRegistered` 方法比较特殊，它在 `BootstrapRegistry` 和 `BootstrapContext` 接口中均可以看到。其实现也不难理解，通过 `Map` 的 `containsKey` 方法，判断给定类型是否已注册。如果给定类型已经注册，则返回 `true`，否则，返回 `false`。

`getRegisteredInstanceSupplier` 方法也比较简单，直接从 `instanceSuppliers` 中获取指定类型的供应者。

`addCloseListener` 方法，用于添加一个监听器，该监听器用于监听 `BootstrapContextClosedEvent`，这块后续 **Huazie** 会带大家实践一下。


### 3.3.3 实现 BootstrapContext 接口中的方法

```java
    @Override
    public <T> T get(Class<T> type) throws IllegalStateException {
        return getOrElseThrow(type, () -> new IllegalStateException(type.getName() + " has not been registered"));
    }

    @Override
    public <T> T getOrElse(Class<T> type, T other) {
        return getOrElseSupply(type, () -> other);
    }

    @Override
    public <T> T getOrElseSupply(Class<T> type, Supplier<T> other) {
        synchronized (this.instanceSuppliers) {
            InstanceSupplier<?> instanceSupplier = this.instanceSuppliers.get(type);
            return (instanceSupplier != null) ? getInstance(type, instanceSupplier) : other.get();
        }
    }

    @Override
    public <T, X extends Throwable> T getOrElseThrow(Class<T> type, Supplier<? extends X> exceptionSupplier) throws X {
        synchronized (this.instanceSuppliers) {
            InstanceSupplier<?> instanceSupplier = this.instanceSuppliers.get(type);
            if (instanceSupplier == null) {
              throw exceptionSupplier.get();
            }
            return getInstance(type, instanceSupplier);
        }
    }

    @SuppressWarnings("unchecked")
    private <T> T getInstance(Class<T> type, InstanceSupplier<?> instanceSupplier) {
        // 省略 。。。
    }
  
```

从 **3.1** 小节，我们了解了 `BootstrapContext` 的 **4** 个获取方法。通过查看上述的源码，我们看到这里只需要分析 `getOrElseSupply` 和 `getOrElseThrow` 的实现即可。

同样在  `getOrElseSupply` 和 `getOrElseThrow` 方法中，我们看到了 `synchronized (this.instanceSuppliers)`，这里同 **3.3.2** 中讲解的一样，都是为了防止多个线程同时修改 `instanceSuppliers` 对象，导致数据不一致的问题。

`getOrElseSupply` 方法的实现也比较简单。如果指定类型的供应者存在，则通过 `getInstance` 方法从这个供应者中获取对应类型的实例对象；否则，直接从提供者 `other` 参数中获取。

`getOrElseThrow` 方法的实现也好理解。如果指定类型的供应者不存在，则直接从异常供应者中获取一个异常类，并将该异常抛出去即可；否则，通过 `getInstance` 方法从这个供应者中获取对应类型的实例对象。

很显然，上述方法最终都需要使用 `getInstance` 方法，从供应者中获取对应类型的实例对象。我们来看看相关的源码：

```java
    T instance = (T) this.instances.get(type);
    if (instance == null) {
        instance = (T) instanceSupplier.get(this);
        if (instanceSupplier.getScope() == Scope.SINGLETON) {
            this.instances.put(type, instance);
        }
    }
    return instance;
```

简单总结如下：

- 首先，从 `instances` 中获取指定类型的实例对象 `instance` 。
- 接着，如果 `instance` 不存在，则继续。
  - 从实例供应者 `instanceSupplier` 中获取一个实例对象，并赋值给 `instance`；
  - 如果实例供应者 `instanceSupplier` 指定的作用域是单例，则将获取的实例对象添加到 `instances` 中，方便后续直接获取。
- 最后，直接返回指定类型的实例对象 `instance`。

### 3.3.4 close 方法

当 `BootstrapContext` 被关闭且 `ApplicationContext` 已准备好时，该方法将被调用【后续笔者讲解 **Spring Boot** 的启动引导过程会涉及到】。

通过阅读源码，我们可以看到这里触发了一个名为 `BootstrapContextClosedEvent` 的事件，该事件会多播给所有注册了该事件的监听器，而这些监听器就是通过 **3.3.2** 小节中提到的 `addCloseListener` 方法添加的【后续 **Huazie** 会带大家实操一下】。

# 四、总结
本篇 **Huazie** 带大家深入了解了 `BootstrapContext` 及其默认实现，这些内容对我们理解 **Spring Boot** 的启动引导过程至关重要。下篇 **Huazie** 将通过自定义 `BootstrapRegistry` 初始化器实现，来看看引导上下文在 **Spring Boot** 的启动引导过程中的作用。
