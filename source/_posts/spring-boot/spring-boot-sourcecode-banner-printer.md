---
title: 【Spring Boot 源码学习】Banner 信息打印流程
date: 2023-11-19 12:11:40
updated: 2024-01-26 09:39:09
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - SpringApplicationBannerPrinter
  - Banners
  - ImageBanner
  - ResourceBanner
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
上篇博文，**Huazie** 带大家初步了解了 `SpringApplication` 的实例化过程。在介绍 `SpringApplication` 的核心构造函数的第一个参数 `ResourceLoader` 时，简单提及了它用于 **Spring Boot** 在启动时打印对应的 **Banner**  信息。这里就引申出了本篇将要介绍的 **Banner** 信息打印流程。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="17" align="left" > 
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
</table>

# 三、主要内容

## 1. printBanner 方法

话不多说，我们先来看 `SpringApplication` 中，有关打印 **banner** 信息的源码：

![](printBanner.png)


上述代码其实是 `SpringApplication` 实例化之后，在其 `public ConfigurableApplicationContext run(String... args)` 方法中被调用的，以下为源码截图，暂不展开介绍，下面重点介绍 `printBanner` 方法的执行流程。

![](printBanner1.png)

## 2. 关闭 Banner 信息打印

我们进入 `printBanner` 方法，先看到了如下判断：

```java
if (this.bannerMode == Banner.Mode.OFF) {
    return null;
}
```

如果 `bannerMode` 是关闭模式，则直接返回 `null`，即不打印 **banner** 信息。

**那我们如何设置 bannerMode 为关闭模式呢？**

在 `SpringApplication` 中，提供了如下的 `setXXX` 方法进行设置：

```java
public void setBannerMode(Banner.Mode bannerMode) {
    this.bannerMode = bannerMode;
}
```

那么我们就可以在启动入口类中，这样来编写：

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
将上述 `setBannerMode` 调用注释掉，运行 `DemoApplication` 类，可见如下截图：

![](setBannerMode.png)
把 `setBannerMode` 调用注释放开，继续运行 `DemoApplication` 类，可见如下截图：

![](setBannerMode1.png)
## 3. SpringApplicationBannerPrinter 类

我们继续往下看源码：

```java
ResourceLoader resourceLoader = (this.resourceLoader != null) ? this.resourceLoader : new DefaultResourceLoader(null);
SpringApplicationBannerPrinter bannerPrinter = new SpringApplicationBannerPrinter(resourceLoader, this.banner);
```
第一行代码是获取资源加载类 `ResourceLoader` ，如果 `SpringApplication` 中的 `resourceLoader` 变量不为 `null`，则使用该变量对应的类作为资源加载类，否则新建一个 `DefaultResourceLoader` 作为默认的资源加载类；

第二行代码是实例化 **Banner** 打印类 `SpringApplicationBannerPrinter` ，它的构造函数分别是：
- `ResourceLoader resourceLoader` ： 资源加载类
- `Banner fallbackBanner` ：**banner** 信息打印接口类![](Banner.png)


最后根据 `bannerMode` 的值不同，有下面两种 **banner** 信息打印模式：
### 3.1 LOG 模式打印

**LOG** 模式打印，可见如下：

```java
bannerPrinter.print(environment, this.mainApplicationClass, logger);
```

继续查看 `SpringApplicationBannerPrinter` 类的 `print` 方法，如下图所示：

![](print.png)
#### 3.1.1 getBanner 方法

首先，这里先调用了 `getBanner` 方法，它用于获取一个 Banner 接口，该接口对应实际要打印的 `Banner` 信息的实现类。

我们查看其相关源码，如下图所示：

![](getBanner.png)

##### 3.1.1.1 新建 Banners 

`Banners` 是 `SpringApplicationBannerPrinter` 的私有静态内部类，它实现了 **Banner** 接口，可以添加多个不同的 `Banner` 实现，也就是它组合了多个不同的 `Banner` 实现，其 `printBanner` 方法就是将内部的不同的 `Banner` 实现按添加顺序依次调用它们自己的 `printBanner`  方法进行打印。

`Banners` 的相关源码，如下图所示：

![](Banners.png)

新建了 `Banners` 的对象 banners 之后， 我们继续往下看：

##### 3.1.1.2 添加 ImageBanner

```java
banners.addIfNotNull(getImageBanner(environment));
```

这里往 `banners` 中添加了一个 `ImageBanner` ，该类也是 `Banner` 接口的一个实现，用于打印从图像资源生成的 **ASCII艺术**【它是一种使用标准 **ASCII** 字符集创建的视觉艺术形式】。当然如果 `getImageBanner` 方法返回 `null`，那么 `banners` 的 `addIfNotNull` 也不会处理。

我们继续看 `getImageBanner` 方法，如下图所示：

![](getImageBanner.png)

下面我们来仔细分析一下上述逻辑：

- 首先，这里从环境配置中获取`BANNER_IMAGE_LOCATION_PROPERTY`  属性值【即 `spring.banner.image.location` 属性对应的值】；
- 接着，**判断上述 Banner 图像位置属性是否存在？？？**
- 如果 **Banner** 图像位置属性存在，则通过资源加载类 `resourceLoader` 获取对应路径的资源对象 `resource`。
  - 如果资源存在，则直接返回一个 `ImageBanner` 的实例对象，构造参数传入上述获取的资源对象 `resource`；
  - 如果资源不存在，则直接返回 `null`。
- 如果 **Banner** 图像位置属性不存在，则依次通过资源加载类 `resourceLoader` 获取 **banner.gif**、**banner.jpg**、**banner.png** 的图像资源 `resource`，如果先匹配到其中一个，则直接返回一个 `ImageBanner` 的实例对象，构造参数传入上述获取的图像资源 `resource`；
- 最后，上述都没有资源存在，则返回 `null`。


##### 3.1.1.3 添加 ResourceBanner

```java
banners.addIfNotNull(getTextBanner(environment));
```
这里往 `banners` 中添加了一个 `ResourceBanner`，该类同样也是 `Banner` 接口的一个实现，用于从源文本资源中打印 `Banner` 信息。当然如果 `getTextBanner` 方法返回 `null`，那么 `banners` 的 `addIfNotNull` 也不会处理。

我们继续看 `getTextBanner` 方法，如下图所示：

![](getTextBanner.png)

相比 `getImageBanner` 方法，`getTextBanner` 方法的逻辑就比较简单：

- 首先，这里依旧是从环境配置中获取 `BANNER_LOCATION_PROPERTY`  属性对应的资源位置值【即 `spring.banner.location` 属性对应的值】；如果未获取到，则默认的资源位置就是 `DEFAULT_BANNER_LOCATION`【即 **banner.txt**】；
- 接着，通过资源加载类 `resourceLoader` 获取指定的位置资源对象 `resource`；
- 然后，检查资源 `resource` 是否存在且其URL不包含`"liquibase-core"`字符串？
  - 如果满足条件，则创建一个 `ResourceBanner` 对象并返回。
  - 如果在尝试访问资源时发生 `IOException` 异常，将在 `catch` 语句块中忽略该异常。
- 最后，如果没有找到符合条件的资源或发生异常，最终将返回 `null`。

##### 3.1.1.4 确认并返回 Banner 实现

```java
if (banners.hasAtLeastOneBanner()) {
    return banners;
}
if (this.fallbackBanner != null) {
    return this.fallbackBanner;
}
return DEFAULT_BANNER;
```

如果上述的 `banners` 中至少存在一个 Banner 实现，则直接返回 `banners` 对象。

反之，如果实例化 `SpringApplicationBannerPrinter` 时，构造函数传入的 `fallbackBanner ` 不为空，则直接返回 `fallbackBanner` 作为 最终的 `Banner` 实现。

如果上述都不符合要求，则返回默认的 `Banner` 实现 `DEFAULT_BANNER` 【即 `SpringBootBanner`】。

```java
private static final Banner DEFAULT_BANNER = new SpringBootBanner();
```

`SpringBootBanner` 其实就是我们启动 **Spring Boot** 打印出来的信息，如下所示：

![](showBannerInfo.png)
#### 3.1.2 以日志模式打印

```java
try {
    logger.info(createStringFromBanner(banner, environment, sourceClass));
}
catch (UnsupportedEncodingException ex) {
    logger.warn("Failed to create String for banner", ex);
}
```

上述 3.1.1 中的截图就是利用日志对象，打印 **INFO** 级别的 **Banner** 信息，最终会被输出到日志文件中。

**Banner** 的具体信息，可见 `createStringFromBanner` 方法，我们继续进入查看：

![](createStringFromBanner.png)

这里逻辑并不复杂，总结如下：
- 首先，创建一个字节数组输出流 `baos`，用于接收要打印的 Banner 信息；
- 接着，调用 **3.1.1.4** 中获取的 `Banner` 实现的 `printBanner` 方法，将要打印的 **Banner** 信息输出到 `baos` 中【这里具体看不同的 `Banner` 实现】；
- 然后，从环境配置中获取 **Banner** 字符集属性值【即 `spring.banner.charset` 属性对应的值】；如果无法获取，则默认是 **UTF-8**；
- 最后，将字节数组输出流转换为指定字符集的字符串，并返回

### 3.2 CONSOLE 模式打印

默认情况下就是 **CONSOLE** 模式打印，可见如下：

```java
bannerPrinter.print(environment, this.mainApplicationClass, System.out)
```

继续查看  `SpringApplicationBannerPrinter` 类的 另一个 `print` 方法，如下图所示：

![](print1.png)

#### 3.2.1 getBanner 方法

详见 **3.1.1** ，这里不再赘述。

#### 3.2.2 以控制台模式打印

由于上面是将 `System.out` 传入到 `PrintStream` 中，所以最终是将 **Banner** 信息直接输出到控制台，可见如下截图：

![](showConsoleBannerInfo.png)
# 四、总结
本篇 **Huazie** 带大家通读了 **Banner** 信息打印的源码，相信如果上面的内容都看下来的话，完全熟悉  **Banner** 信息打印流程不再是个问题。有了这些基础的知识，我们就可以来自定义  **Banner** 信息打印，敬请期待下篇博文！！！

