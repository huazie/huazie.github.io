---
title: 【Spring Boot 源码学习】SpringApplication 的定制化介绍
date: 2024-01-07 22:35:03
updated: 2024-01-17 08:32:03
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - SpringApplication
  - 基础配置
  - 数据源配置
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
前面的博文，**Huazie** 带大家从 **Spring Boot** 的启动类 `SpringApplication` 上入手，了解了 `SpringApplication` 的实例化过程。这实例化构造过程中包含了各种初始化的操作，都是 **Spring Boot** 默认配置的。如果我们需要定制化配置，`SpringApplication` 也提供了相关的入口，且看下面的介绍。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="22" align="left" > 
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
</table>

# 三、主要内容
针对 `SpringApplication` 的定制化配置，**Spring Boot** 中也提供了不同的方式，比如通过入口类、配置文件、环境变量、命令行参数等等。

> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 1. 基础配置

所谓基础配置，即是可以直接通过 `set` 或 `add` 方法来进行参数的配置，这些 `set` 或 `add` 方法修改的配置都是 **Spring Boot** 预定义的一些参数，有些也可以在 **application.properties** 文件中进行配置。

### 1.1   设置关闭 Banner
在笔者的[《Banner 信息打印流程》](/2023/11/19/spring-boot/spring-boot-sourcecode-banner-printer/) 中，第 2 小节就介绍了如何关闭 **Banner** 信息打印。

通过 `SpringApplication` 提供的 `setBannerMode` 方法，我们就可以在启动入口类中，这样来编写：


```java
@SpringBootApplication
public class DemoApplication {
    public static void main(String[] args) {
        SpringApplication springApplication = new SpringApplication(DemoApplication.class);
        springApplication.setBannerMode(Banner.Mode.OFF);
        springApplication.run(args);
    }
}
```

### 1.2 设置自定义 Banner 打印对象

在笔者的[《自定义 Banner 信息打印》](/2023/11/24/spring-boot/spring-boot-sourcecode-custom-banner-printer/) 中，第 4 小节就介绍了如何自定义 **Banner** 接口实现。

通过 `SpringApplication` 提供的 `setBanner` 方法，我们可以修改入口类，如下：

```java
@SpringBootApplication
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication springApplication = new SpringApplication(DemoApplication.class);
        springApplication.setBanner(new CustomBanner());
        springApplication.run(args);
    }
}
```

### 1.3 设置应用程序主入口类

在笔者的[《初识 SpringApplication》](/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/) 中，2.6 小节介绍了 `SpringApplication` 实例化时，会自动推断应用入口类，最终会被赋值给 `SpringApplication` 的成员变量 `mainApplicationClass`。

当然，通过 `SpringApplication` 提供的 `setMainApplicationClass` 方法，设置一个特定的主应用程序类，该类将用作日志源并获取版本信息。默认情况下，`SpringApplication` 实例化时，会自动推断主应用程序类。如果没有明确的应用程序类，我们可以设置为 **null**。

### 1.4 设置用于创建应用程序上下文的工厂

```java
public void setApplicationContextFactory(ApplicationContextFactory applicationContextFactory) {
    this.applicationContextFactory = (applicationContextFactory != null) ? applicationContextFactory
        : ApplicationContextFactory.DEFAULT;
}
```

通过 `SpringApplication` 提供的 `setApplicationContextFactory` 方法，我们可以用于创建应用程序上下文的工厂。如果没有设置，默认为一个工厂【即 `DefaultApplicationContextFactory`】，该工厂将为 **Servlet Web 应用程序** 创建 `AnnotationConfigServletWebServerApplicationContext`，为**响应式 Web 应用程序**  创建 `AnnotationConfigReactiveWebServerApplicationContext`，并为 **非 Web 应用程序** 创建 `AnnotationConfigApplicationContext`。

![](DefaultApplicationContextFactory.png)

### 1.5 添加 BootstrapRegistry 初始化器实现

在 **Huazie** 的[《BootstrapRegistryInitializer 详解》](/2023/11/30/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer/)中，介绍了 加载和初始化 `BootstrapRegistryInitializer` 的逻辑，有需要的小伙伴可以去瞅一眼。

**那除了默认的加载过程，还有啥办法手动添加 `BootstrapRegistryInitializer` 呢？**

通过 `SpringApplication` 的 `addBootstrapRegistryInitializer` 方法，我们可以在 `bootstrapRegistryInitializers` 中添加额外的 `BootstrapRegistry` 初始化器实现。

```java
public void addBootstrapRegistryInitializer(BootstrapRegistryInitializer bootstrapRegistryInitializer) {
    Assert.notNull(bootstrapRegistryInitializer, "BootstrapRegistryInitializer must not be null");
    this.bootstrapRegistryInitializers.addAll(Arrays.asList(bootstrapRegistryInitializer));
}
```


### 1.6 设置或添加 ApplicationContext 初始化器实现

在 **Huazie** 的[《ApplicationContextInitializer 详解》](/2023/12/03/spring-boot/spring-boot-sourcecode-applicationcontextinitializer/)中，介绍了加载和初始化 `ApplicationContextInitializer ` 的逻辑，大家可以自行去回顾下。

**除了默认的加载过程，我们还可以通过 `SpringApplication` 自身进行设置或添加**

通过 `SpringApplication` 的 `setInitializers` 方法，我们可以重新设置 initializers【**注意：** 调用 `setInitializers` 方法后，该变量之前的赋值都将丢失】。

```java
private List<ApplicationContextInitializer<?>> initializers;

public void setInitializers(Collection<? extends ApplicationContextInitializer<?>> initializers) {
    this.initializers = new ArrayList<>(initializers);
}
```

通过 `SpringApplication` 的 `addInitializers` 方法，我们可以在 `initializers` 中添加额外的 `ApplicationContextInitializer` 数组。 

```java
public void addInitializers(ApplicationContextInitializer<?>... initializers) {
    this.initializers.addAll(Arrays.asList(initializers));
}
```

### 1.7 设置 ApplicationListener 实现

在 **Huazie** 的[《ApplicationListener 详解》](/2023/12/10/spring-boot/spring-boot-sourcecode-applicationlistener/)中，我们详细分析了 ApplicationListener 的加载和处理应用程序事件的逻辑，有需要可以去回顾下。

那除了默认的加载流程，我们还可以通过 `SpringApplication` 的 `setListeners` 方法，重新设置 `listeners`【**注意：** 调用 `setListeners` 方法后，`listeners` 之前的赋值都将丢失】。

```java
private List<ApplicationListener<?>> listeners;

public void setListeners(Collection<? extends ApplicationListener<?>> listeners) {
    this.listeners = new ArrayList<>(listeners);
}
```

当然也可以通过 `SpringApplication` 的 `addListeners` 方法，在 `listeners ` 中添加额外的 `ApplicationListener` 数组。 

```java
public void addListeners(ApplicationListener<?>... listeners) {
    this.listeners.addAll(Arrays.asList(listeners));
}
```

### 1.8 设置要运行的Web应用程序的类型

在 **Huazie** 的[《初识 SpringApplication》](/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)中，2.2 小节分析了 `SpringApplication` 构造函数中默认的 **Web** 应用类型推断的逻辑。

当然，我们也可以通过 `SpringApplication` 的 `setWebApplicationType` 方法，手动设置要运行的 **Web** 应用程序的类型。

```java
public void setWebApplicationType(WebApplicationType webApplicationType) {
    Assert.notNull(webApplicationType, "WebApplicationType must not be null");
    this.webApplicationType = webApplicationType;
}
```

### 1.9 设置 SpringApplication 中 各种 boolean 变量

#### 1.9.1 setAllowBeanDefinitionOverriding 方法

```java
public void setAllowBeanDefinitionOverriding(boolean allowBeanDefinitionOverriding) {
    this.allowBeanDefinitionOverriding = allowBeanDefinitionOverriding;
}
```

设置是否允许通过注册一个与现有定义具有相同名称的定义来覆盖 **bean** 定义。默认为 `false`。

具体可见 `DefaultListableBeanFactory#setAllowBeanDefinitionOverriding(boolean)`


#### 1.9.2 setAllowCircularReferences 方法

```java
public void setAllowCircularReferences(boolean allowCircularReferences) {
    this.allowCircularReferences = allowCircularReferences;
}
```

设置是否允许 **bean** 之间的循环引用，并自动尝试解析它们。默认为 `false`。

具体可见 `AbstractAutowireCapableBeanFactory#setAllowCircularReferences(boolean)`

#### 1.9.3 setLazyInitialization 方法

```java
public void setLazyInitialization(boolean lazyInitialization) {
    this.lazyInitialization = lazyInitialization;
}
```

设置是否应延迟初始化 **beans**。默认为 `false`。

具体可见 `BeanDefinition#setLazyInit(boolean)`

#### 1.9.4 setHeadless 方法

```java
public void setHeadless(boolean headless) {
    this.headless = headless;
}
```

设置应用程序是否为无头模式，即是否不应实例化 **AWT**。默认为 `true`，以防止出现 **Java** 图标。


#### 1.9.5 setRegisterShutdownHook 方法

```java
public void setRegisterShutdownHook(boolean registerShutdownHook) {
    this.registerShutdownHook = registerShutdownHook;
}
```

设置是否应注册一个关闭钩子（shutdown hook）到创建的 `ApplicationContext`。默认为 `true`，以确保 **JVM** 关闭时能够优雅地处理。

#### 1.9.6 setLogStartupInfo 方法

```java
public void setLogStartupInfo(boolean logStartupInfo) {
    this.logStartupInfo = logStartupInfo;
}
```

设置在应用程序启动时是否应记录应用程序信息。默认为 `true`。


#### 1.9.7 setAddCommandLineProperties 方法

```java
public void setAddCommandLineProperties(boolean addCommandLineProperties) {
    this.addCommandLineProperties = addCommandLineProperties;
}
```

设置是否应将 `CommandLinePropertySource` 添加到应用程序上下文中，以便暴露参数。默认为 `true`。


#### 1.9.8 setAddConversionService 方法

```java
public void setAddConversionService(boolean addConversionService) {
    this.addConversionService = addConversionService;
}
```
设置是否应将 `ApplicationConversionService` 添加到应用程序上下文的环境中。

在 **Spring Boot** 中，`ApplicationConversionService` 是一个重要的类型转换服务，用于实现应用程序中的数据转换。它提供了一种将一种类型的数据转换为另一种类型数据的方法，使得在不同组件或服务之间能够进行数据交互和集成。

###  1.10 设置默认的环境属性

```java
public void setDefaultProperties(Map<String, Object> defaultProperties) {
    this.defaultProperties = defaultProperties;
}

public void setDefaultProperties(Properties defaultProperties) {
    this.defaultProperties = new HashMap<>();
    for (Object key : Collections.list(defaultProperties.propertyNames())) {
      this.defaultProperties.put((String) key, defaultProperties.get(key));
    }
}
```

上述方法用于设置默认的环境属性，这些属性将在现有环境属性的基础上进行添加。


### 1.11 设置要使用的额外的配置文件值

```java
public void setAdditionalProfiles(String... profiles) {
    this.additionalProfiles = Collections.unmodifiableSet(new LinkedHashSet<>(Arrays.asList(profiles)));
}
```

该方法用于设置要使用的额外的配置文件值，这些值将在系统或命令行属性的基础上进行补充。


### 1.12 设置  bean 名称生成器

```java
public void setBeanNameGenerator(BeanNameGenerator beanNameGenerator) {
    this.beanNameGenerator = beanNameGenerator;
}
```

设置在生成 **bean** 名称时应该使用的 **bean** 名称生成器。


### 1.13 设置底层环境

```java
public void setEnvironment(ConfigurableEnvironment environment) {
    this.isCustomEnvironment = true;
    this.environment = environment;
}
```

设置与创建的应用程序上下文一起使用的底层环境。

### 1.14 设置资源加载器

```java
public void setResourceLoader(ResourceLoader resourceLoader) {
    Assert.notNull(resourceLoader, "ResourceLoader must not be null");
    this.resourceLoader = resourceLoader;
}
```

设置在加载资源时应使用的 `ResourceLoader`。

### 1.15 设置环境配置属性前缀

```java
public void setEnvironmentPrefix(String environmentPrefix) {
    this.environmentPrefix = environmentPrefix;
}
```

设置从系统环境中获取配置属性时应使用的前缀。


### 1.16 设置应用启动指标对象

```java
public void setApplicationStartup(ApplicationStartup applicationStartup) {
    this.applicationStartup = (applicationStartup != null) ? applicationStartup : ApplicationStartup.DEFAULT;
}
```

设置用于收集启动指标的 `ApplicationStartup`。如果没有指定，则默认使用 `DefaultApplicationStartup` 。

## 2. 数据源配置

除了上述直接通过 `set` 或 `add` 方法来进行参数的配置，`SpringApplication` 中还提供了可以通过设置配置源参数对整个配置文件或配置类进行配置。

### 2.1 通过 SpringApplication 构造方法参数

在 Huazie 的[《初识 SpringApplication》](/2023/11/12/spring-boot/spring-boot-sourcecode-springapplication/)中的 2.1 小节就介绍了可以通过其构造参数 `primarySources` 来配置普通类或指定某个配置类，但这种方式有其局限性，它无法指定 **XML** 配置和基于 **package** 的配置。

### 2.2 通过 setSources 方法

先来看看相关的源码：

```java
private Set<String> sources = new LinkedHashSet<>();

public void setSources(Set<String> sources) {
    Assert.notNull(sources, "Sources must not be null");
    this.sources = new LinkedHashSet<>(sources);
}
```

该方法的参数为 `String` 类型的 `Set` 集合，可以传类名、package 名 和 XML 配置文件资源。

下面我们来演示一下：

- 首先，我们在 `application.properties` 中添加如下配置：

  ```java
  author=huazie
  ```
- 然后，新增一个普通类 `CustomConfiguration`，如下：

  ```java
  public class CustomConfiguration {

      @Value("${author}")
      private String author;
  
      public CustomConfiguration() {
          System.out.println("CustomConfiguration已被创建");
      }
  
      public String getAuthor() {
          return author;
      }
  
      public void setAuthor(String author) {
          this.author = author;
      }
  }
  ```
- 接着，我们重新编写 `DemoApplication` 类，如下所示：

  ```java
  @SpringBootApplication
  public class DemoApplication {
  
      public static void main(String[] args) {
          SpringApplication springApplication = new SpringApplication(DemoApplication.class);
          Set<String> sources = new HashSet<>();
          sources.add(CustomConfiguration.class.getName());
          springApplication.setSources(sources);
          ConfigurableApplicationContext context = springApplication.run(args);
          CustomConfiguration customConfiguration = context.getBean(CustomConfiguration.class);
          System.out.println(customConfiguration.getAuthor());
      }
  }
  ```

- 最后，我们运行 `DemoApplication` 中的 `main` 方法，从如下截图中可以看出这里已经打印了自定义类的属性值：
  ![](showCustomConfiguration.png)


### 2.3 合并配置源信息

无论是通过构造参数，还是通过 `setSources` 方法，对配置源信息进行指定，在 **Spring Boot** 中都会将其合并。

这里我们就不得不提，`SpringApplication` 提供的 `getAllSources` 方法，

![](getAllSources.png)

该方法将构造函数中指定的任何主要源 `primarySources` 与已显式设置的任何其他源 `sources` 组合在一起。


# 四、总结
**23 年 7 月** ，**Huazie** 正式开启了【**Spring 源码学习**】系列，一路下来也有 **22** 篇文章了。**24** 年该系列将继续更新下去，希望阅读和订阅多多益善！！！
