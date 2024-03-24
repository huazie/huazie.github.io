---
title: 【Spring Boot 源码学习】共享 MetadataReaderFactory 上下文初始化器
date: 2024-03-24 20:26:00
updated: 2024-03-24 20:26:00
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - SharedMetadataReaderFactoryContextInitializer
  - CachingMetadataReaderFactoryPostProcessor
  - BeanDefinitionRegistryPostProcessor
  - ConfigurationClassPostProcessor
  - ConfigurationClassPostProcessorCustomizingSupplier 
  - SharedMetadataReaderFactoryBean
  - FactoryBean
  - BeanClassLoaderAware
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
上篇博文[《深入应用上下文初始化器实现》](/2024/03/17/spring-boot/spring-boot-sourcecode-applicationcontextinitializer-impl/)，**Huazie** 带大家详细分析了 分析 **Spring Boot** 中预置的应用上下文初始化器实现【即 `ApplicationContextInitializer` 接口实现类】的源码，了解了在 **Spring** 容器刷新之前初始化应用程序上下文的一些具体操作。

当然其中有些实现源码比较复杂，还没有深入分析。那本篇就来对其中的 
`SharedMetadataReaderFactoryContextInitializer` 【即 **共享 `MetadataReaderFactory` 上下文初始化器**】详细分析下。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="28" align="left" > 
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
  <tr>
    <td align="left" > 
      <a href="/2024/03/17/spring-boot/spring-boot-sourcecode-applicationcontextinitializer-impl/">【Spring Boot 源码学习】深入应用上下文初始化器实现</a> 
    </td>
  </tr>
</table>

# 三、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

## 3.1 源码初识

我们先来看看 `SharedMetadataReaderFactoryContextInitializer` 的部分源码，如下：

```java
class SharedMetadataReaderFactoryContextInitializer
    implements ApplicationContextInitializer<ConfigurableApplicationContext>, Ordered {
    
  // 其他省略。。。
  
  @Override
  public void initialize(ConfigurableApplicationContext applicationContext) {
    BeanFactoryPostProcessor postProcessor = new CachingMetadataReaderFactoryPostProcessor(applicationContext);
    applicationContext.addBeanFactoryPostProcessor(postProcessor);
  }
  
  @Override
  public int getOrder() {
    return 0;
  }
  // 其他省略。。。
}
```
从上述源码中，我们可以看出 `SharedMetadataReaderFactoryContextInitializer` 实现了 `ApplicationContextInitializer<ConfigurableApplicationContext>` 和 `Ordered` 接口：
- `ApplicationContextInitializer<ConfigurableApplicationContext>` ：应用上下文初始化器接口类，有关该类的详细介绍，请查看[《ApplicationContextInitializer 详解》](/2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/)。
- `Ordered` ：实现该接口可以控制应用上下文初始化器实现类的执行顺序，有关这点我们可以查看 `SpringApplication` 的 `getInitializers` 方法。
 ![](getInitializers.png)
 ![](asUnmodifiableOrderedSet.png)

这里排序的关键就是 **spring-core** 包中提供的 `org.springframework.core.annotation.AnnotationAwareOrderComparator`，它能够对实现了 `PriorityOrdered` 接口、`Ordered` 接口或被 `@Order` 注解修饰的类进行统一的排序。


我们继续查看上述 `initialize` 方法，可以看到这里向应用上下文中添加了一个 `BeanFactoryPostProcessor`【 即 `CachingMetadataReaderFactoryPostProcessor`】。

## 3.2 CachingMetadataReaderFactoryPostProcessor
我们继续查看 `CachingMetadataReaderFactoryPostProcessor` 的源码，如下：

![](CachingMetadataReaderFactoryPostProcessor.png)

从上述截图中，我们可以看出 `CachingMetadataReaderFactoryPostProcessor` 是一个静态内部类，它同时实现了 `PriorityOrdered` 和 `BeanDefinitionRegistryPostProcessor` 接口，有关这两个接口的作用，可以查看 [《深入应用上下文初始化器实现》](/2024/03/17/spring-boot/spring-boot-sourcecode-applicationcontextinitializer-impl/)中的 **3.1.1** 小节 ，这里不再赘述。

我们继续查看 `postProcessBeanDefinitionRegistry` 方法，发现这里调用了 `register` 方法 和 `configureConfigurationClassPostProcessor` 方法，下面一一介绍：

> `postProcessBeanDefinitionRegistry` 方法是 `BeanDefinitionRegistryPostProcessor` 接口中定义的方法，它用于在标准初始化之后修改应用上下文的内部 **bean** 定义注册表。所有的常规 **bean** 定义都将已经被加载，但还没有实例化任何 **bean**。这允许在下一个后处理阶段开始之前添加更多的 **bean** 定义。
### 3.2.1  register 方法
首先，进入 `register` 方法，如下所示：

```java
private void register(BeanDefinitionRegistry registry) {
  if (!registry.containsBeanDefinition(BEAN_NAME)) {
    BeanDefinition definition = BeanDefinitionBuilder
      .rootBeanDefinition(SharedMetadataReaderFactoryBean.class, SharedMetadataReaderFactoryBean::new)
      .getBeanDefinition();
    registry.registerBeanDefinition(BEAN_NAME, definition);
  }
}
```
其中 `BEAN_NAME` 如下截图所示：

![](BEAN_NAME.png)

`register` 方法逻辑简单，它的功能是检查 `BeanDefinitionRegistry` 中是否已存在名为 `BEAN_NAME` 的 `BeanDefinition`，如果不存在，则创建一个 `SharedMetadataReaderFactoryBean` 的 `BeanDefinition` 并将其注册到 `registry` 中。【有关 `SharedMetadataReaderFactoryBean` ，可以查看 **3.4** 小节】

> **知识点：** `BeanDefinitionRegistry` 是 **Spring** 中一个接口，它可以被看作是一个用来管理 `BeanDefinition` 的注册表。`BeanDefinition` 可以被理解为 **Spring** 中 **Bean** 的配置描述，它包含了 **Bean** 的元数据，如**类名**、**是否为抽象类**、**构造函数**、**属性值** 等相关信息。这些信息将会告诉 **Spring** 如何创建和初始化相应的 **Bean**。

### 3.2.1 configureConfigurationClassPostProcessor 方法
接着，进入 `configureConfigurationClassPostProcessor` 方法，可见如下截图：

![](configureConfigurationClassPostProcessor.png)
![](CONFIGURATION_ANNOTATION_PROCESSOR_BEAN_NAME.png)
在 `configureConfigurationClassPostProcessor(BeanDefinitionRegistry)` 方法中，先从 `BeanDefinitionRegistry` 中获取名为 `org.springframework.context.annotation.internalConfigurationAnnotationProcessor` 【即内部管理的 `Configuration` 注解处理器的 bean 名称】的 `BeanDefinition`；如果找不到对应的 `BeanDefinition`，则捕获 `NoSuchBeanDefinitionException` 异常后不做任何处理。

![](BeanDefinition.png)

> **知识点：** `ConfigurationClassPostProcessor` 是 **Spring** 框架中的一个核心类，它实现了 `BeanPostProcessor` 的子接口 `BeanDefinitionRegistryPostProcessor`。它的主要作用是解析被 `@Configuration` 注解的类，并将解析到的 **Bean** 封装为 `BeanDefinition` 注册到 **Spring** 容器中，以供后续步骤进行统一的实例化。此外，`ConfigurationClassPostProcessor` 还会处理其他与配置相关的注解，如 `@Component、@PropertySources、@ComponentScan、@Import`等。


继续进入 `configureConfigurationClassPostProcessor(BeanDefinition)` 方法中：

![](configureConfigurationClassPostProcessor.png)

从上述截图中，可以看到：

- 如果 `definition` 是 `AbstractBeanDefinition 的实例`，则调用 `configureConfigurationClassPostProcessor(AbstractBeanDefinition)` 方法：
  - 首先，从 `AbstractBeanDefinition` 中获取 **bean** 的实例提供者 `instanceSupplier`。
  - 接着，判断该实例提供者是否为空 ？
    - 如果不为空，则
      - 重新设置 `definition` 的实例提供者为 **3.3** 中的自定义供应者 `ConfigurationClassPostProcessorCustomizingSupplier`。
      - 直接返回即可。
    - 如果为空，则调用  `configureConfigurationClassPostProcessor(MutablePropertyValues)` 方法。

- 如果 `definition` 不是 `AbstractBeanDefinition 的实例`，则直接调用 `configureConfigurationClassPostProcessor(MutablePropertyValues)` 方法。

我们继续查看 `configureConfigurationClassPostProcessor(MutablePropertyValues)` 方法，里面只有一行代码，功能是向 `propertyValues` 中添加一个新的属性 `"metadataReaderFactory"`，其值为一个指向当前上下文中名为 `org.springframework.boot.autoconfigure.internalCachingMetadataReaderFactory` 的  **Bean** 的引用【`RuntimeBeanReference`】。

**知识点：** 
- `MutablePropertyValues` 是 **Spring** 框架中的一个类，主要用于封装类属性的集合。它是一个 `List` 容器，包装了多个 `PropertyValue` 对象，每个 `PropertyValue` 对象封装了一个属性及其对应的值。当需要在 `BeanDefinition` 中修改某个类里面的属性时，就可以使用`MutablePropertyValues` 类。

- `RuntimeBeanReference` 是 **Spring** 框架中用于表示运行时 **Bean** 引用的一个对象。在 **Spring** 的 **Bean** 解析阶段，当解析器遇到需要依赖其他 **Bean** 的情况时，它会依据依赖 **Bean** 的名称创建一个`RuntimeBeanReference` 对象，并将这个对象放入 `BeanDefinition` 的`MutablePropertyValues` 中。这个 `RuntimeBeanReference` 对象是对实际 **Bean** 的引用，它会在运行时被解析成实际的 **Bean** 对象。

## 3.3 ConfigurationClassPostProcessor 的自定义供应者
还是一样，先来看看源码：

![](ConfigurationClassPostProcessorCustomizingSupplier.png)

`ConfigurationClassPostProcessorCustomizingSupplier` 是一个实现了 `Supplier<Object>` 接口的自定义供应者类，其包含两个成员变量：
- `ConfigurableApplicationContext context`：应用上下文对象
- `Supplier<?> instanceSupplier`：原始的 `Supplier`，用于获取`ConfigurationClassPostProcessor` 的实例。

通过阅读上述 `get` 方法，我们可以看到该类在不改变原始 `Supplier` 逻辑的情况下，对提供的 `ConfigurationClassPostProcessor` 实例重新设置了 `metadataReaderFactory` 属性值，而该值是通过调用 `context.getBean(BEAN_NAME, MetadataReaderFactory.class)` 从 **Spring** 上下文中获取的一个 `MetadataReaderFactory` 的 `Bean` 对象。

## 3.4 共享 MetadataReaderFactory 的 FactoryBean
话不多说，先来看看相关源码截图：

![](SharedMetadataReaderFactoryBean.png)

`SharedMetadataReaderFactoryBean` 也是一个静态内部类，它实现了 `FactoryBean<ConcurrentReferenceCachingMetadataReaderFactory>`、`BeanClassLoaderAware` 和`ApplicationListener<ContextRefreshedEvent>` 这三个接口，下面来详细分析下：
### 3.4.1 FactoryBean 接口

`FactoryBean` 是 **Spring** 框架中用于创建复杂 **Bean** 的接口。它包含如下三个方法：

```java
T getObject() throws Exception;

Class<?> getObjectType();

default boolean isSingleton() {
  return true;
}
```

- `getObject()`：该方法用于返回由 **FactoryBean** 创建的对象实例。这是最核心的方法，通过实现这个方法，可以定义如何创建和返回所需的对象。在 `SharedMetadataReaderFactoryBean` 中，该方法返回 `ConcurrentReferenceCachingMetadataReaderFactory` 的实例 `metadataReaderFactory`
-  `getObjectType()`：该方法用于返回由 **FactoryBean** 创建的对象的类型，如果事先不知道，则返回 `null`。在 `SharedMetadataReaderFactoryBean` 中，该方法返回`CachingMetadataReaderFactory.class`，虽然实际的类型是`ConcurrentReferenceCachingMetadataReaderFactory`【**这里暂且打个问号，有清楚的朋友可以评论区讨论下**】
- `isSingleton()`：该方法用于判断由 **FactoryBean** 创建的对象是否为单例。如果返回 `true`，则表示创建的对象在 **Spring IoC** 容器中是单例的，即整个应用程序中只有一个实例；如果返回 `false`，则表示每次请求都会创建一个新的实例。在 `SharedMetadataReaderFactoryBean` 中，该方法返回 `true`。

### 3.4.2 BeanClassLoaderAware 接口

`BeanClassLoaderAware` 是 **Spring** 框架中的一个 `Aware` 接口，它的主要作用是允许 `Bean` 在初始化时获取关于自身类加载器的信息，以便执行一些特定的操作，比如动态加载其他类、访问资源等。

该接口只有一个方法：

```java
void setBeanClassLoader(ClassLoader classLoader);
```

`SharedMetadataReaderFactoryBean` 实现了该接口，并重写了 `setBeanClassLoader` 方法，并在该方法中，使用传入的 `classLoader` 来创建一个新的 `ConcurrentReferenceCachingMetadataReaderFactory` 实例，然后将其赋值给成员变量 `metadataReaderFactory`【如上 **3.4** 中源码截图中可见 `SharedMetadataReaderFactoryBean##getObject()` 返回的就是该变量】。

### 3.4.3 ApplicationListener 接口

 监听器接口 `ApplicationListener`，在之前的博文已经介绍过。`SharedMetadataReaderFactoryBean` 实现了该接口，实现 `onApplicationEvent` 方法，并监听 `ContextRefreshedEvent` 事件。当接收到 `ContextRefreshedEvent` 事件时，就会回调 `onApplicationEvent` 方法，然后 `onApplicationEvent` 方法里调用 `metadataReaderFactory` 的 `clearCache` 方法来清除缓存。这是为了在应用上下文刷新后确保 `MetadataReader` 缓存是最新的。

![](ConcurrentReferenceCachingMetadataReaderFactory.png)
# 四、总结
本篇 **Huazie** 带大家一起分析了 **spring-boot-autoconfigure** 子模块中预置的 应用上下文初始化器实现 `SharedMetadataReaderFactoryContextInitializer` 。其中涉及了很多 **Spring** 的知识，由于篇幅受限没有细说，大家可以查看相关 **Spring** 文档，并运行 **Spring Boot** 项目进一步加深理解。

