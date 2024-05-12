---
title: 【Spring Boot 源码学习】自定义 Banner 信息打印
date: 2023-11-24 20:14:37
updated: 2024-01-25 09:59:09
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - ImageBanner
  - ResourceBanner
  - 自定义Banner接口实现
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
上篇博文，**Huazie** 带大家了解了完整的 **Banner** 信息打印流程。相信大家都跃跃一试了，那么本篇就以这些基础的知识，来自定义 **Banner** 信息打印。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="18" align="left" > 
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
</table>

# 主要内容

> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

## 1. ResourceBanner 打印
通过 `SpringApplicationBannerPrinter##getTextBanner` 方法的源码了解，我们现在可以进行如下的操作：
### 1.1 添加默认的 banner.txt 资源文件
当没有配置 `spring.banner.location` 属性，`Spring Boot` 默认就会加载资源根目录的 **banner.txt** 文件，如果存在该资源文件，则会使用 `ResourceBanner` 打印 **Banner** 信息。

下面我们在新建的 **demo** 项目的资源根目录添加名为 **banner.txt** 的资源文件，如下图所示：

![](banner-txt.png)

然后，直接运行我们的 `DemoApplication` 入口类，可见如下运行截图：

![](show-banner-txt.png)
### 1.2 指定任意路径的资源文件
现在我们把上面的 **banner.txt** 移到资源根目录新建的 **banner** 目录里，并更名为 **mybanner.txt**，如下图所示：

![](mybanner-txt.png)

接着，在 **application.properties** 中配置如下：

```bash
# Banner 资源文件路径
spring.banner.location=classpath:banner/mybanner.txt
```

然后，还是直接运行我们的 `DemoApplication` 入口类，可见如下运行截图：

![](show-mybanner-txt.png)
### 1.3 添加自定义的信息

查看 `ResourceBanner` 的源码，我们可以看到如下的代码：

![](ResourceBanner.png)
![](ResourceBanner1.png)
![](ResourceBanner2.png)

这里就不得不提 `PropertyResolver`，它是 **Spring** 框架中的一个组件，主要用于解析各种属性源的属性值。它能够处理多种类型的底层源，包括**properties** 文件、**yaml** 文件，甚至是一些 **nosql** 数据库【因为这些数据源同样采用 **key-value** 形式存储数据】。

查看 `PropertyResolver` 的 **API** 中，我可以看到它定义了一系列读取、解析和判断是否包含指定属性的方法。此外，它还支持以 `${propertyName:defaultValue}` 格式的属性占位符，替换为实际的值的功能，这在动态配置中非常有用。

接下来，我们在 `application.properties` 中配置如下：

![](defaults-property.png)

然后，我们在 **banner.txt** 中可以添加如下属性占位符：

![](banner-txt-new.png)

最后，运行我们的 `DemoApplication` 入口类，可见如下运行截图：

![](show-banner-txt-new.png)

## 2. ImageBanner  打印
通过 `SpringApplicationBannerPrinter##getImageBanner` 方法的源码了解，我们现在可以进行如下的操作：
### 2.1 添加默认的图像资源文件
当没有配置 `spring.banner.image.location` 属性，`Spring Boot` 默认就会加载资源根目录的 **banner.gif** 或 **banner.jpg** 或 **banner.png** 等文件，如果存在其中某个资源文件，则会使用 `ImageBanner` 打印 **Banner** 信息。

下面我们在新建的 **demo** 项目的资源根目录添加名为 **banner.gif** 的资源文件，如下图所示：

![](banner-gif.png)

然后，同样运行我们的 `DemoApplication` 入口类，可见如下运行截图：

![](show-banner-gif.png)

换成另外两个 **banner.jpg** 或 **banner.png** 也是能够加载的，如下：
![](show-banner-jpg.png)
![](show-banner-png.png)

默认 **Banner** 图像资源的加载逻辑：
- 存在 **banner.gif**，则只加载 **banner.gif**；
- 不存在 **banner.gif**，存在 **banner.jpg**，则只加载 **banner.jpg**
- 不存在 **banner.gif**，也不存在 **banner.jpg**，则加载 **banner.png**

### 2.2 指定任意路径的图像资源文件

现在我们把上面的 **banner.png** 移到资源根目录新建的 **banner** 目录里，并更名为 **mybanner.png**，如下图所示：

![](mybanner-png.png)

接着，在 **application.properties** 中配置如下：

```bash
# Banner 图像资源文件路径
spring.banner.image.location=classpath:banner/mybanner.png
```

然后，我们运行 `DemoApplication` 入口类，可见如下运行截图：

![](show-mybanner-png.png)
### 2.3 添加自定义的图像显示信息

查看 `ImageBanner` 的源码，我们可以看到如下的代码：

![](ImageBanner.png)
![](ImageBanner1.png)
![](ImageBanner2.png)

从上述源码，我们看到 `ImageBanner` 里面可以自定义一些图像的显示属性，比如：

- `spring.banner.image.width` ：设置 **banner** 图像的宽度，默认为 76 像素
- `spring.banner.image.height` ：设置 **banner** 图像的高度，默认按照宽度计算缩放比例，重新计算新图像的高度。
  ![](ImageBanner3.png)

- `spring.banner.image.margin` ：设置 **banner** 图像的外边距，默认为 2 像素
- `spring.banner.image.invert` ：设置是否反转图片的颜色。如果设置为 true，则颜色会被反转
- `spring.banner.image.bitdepth` ：设置图片的位深度，默认 4 位深度，还支持 8 位深度。位深度决定了图片的颜色精度，例如8位深度表示每个像素有256种颜色，不过大多数情况下，对于 **Banner** 图像输出到控制台，看起来基本没啥区别。
- `spring.banner.image.pixelmode` ：设置图片的的像素模式，有如下两个枚举值：
  - `TEXT` ：文本模式，适用于需要清晰、简洁的图像效果的情况。
  ![](pixelmode-text.png)
  - `BLOCK` ：块模式，适用于需要强调图像的某些部分或突出显示特定区域的情况。
  ![](pixelmode-block.png)

下面我们就来添加这些属性，来看看效果：

#### 2.3.1 添加 Banner 图像宽度属性

```bash
spring.banner.image.width=50
```

运行 `DemoApplication` 入口类，可见如下运行截图： 

![](banner-image-width.png)
#### 2.3.2 添加 Banner 图像高度属性

```bash
spring.banner.image.height=20
```

依旧运行 `DemoApplication` 入口类，可见如下运行截图：

![](banner-image-height.png)

#### 2.3.3 添加 Banner 图像外边距属性

```bash
spring.banner.image.margin=5
```

同样运行 `DemoApplication` 入口类，可见如下运行截图：

![](banner-image-margin.png)
#### 2.3.4 添加 Banner 图像是否反转图片颜色的属性

```bash
spring.banner.image.invert=true
```

继续运行 `DemoApplication` 入口类，可见如下运行截图：

![](banner-image-invert.png)
#### 2.3.5 添加 Banner 图像位深度的属性

```bash
spring.banner.image.bitdepth=8
```

然后运行 `DemoApplication` 入口类，可见如下运行截图：

![](banner-image-bitdepth.png)
我们发现上述好像设置了该属性，展示出来的图像并没有啥差异，事实上也的确如此，可能我们的图像比较简单。

#### 2.3.6 添加 Banner 图像像素模式的属性

```bash
spring.banner.image.pixelmode=block
```

运行 `DemoApplication` 入口类，可见如下运行截图：
![](banner-image-pixelmode-block.png)
![](banner-image-pixelmode-block1.png)

## 3. Banners 打印

`Banners` 是 `SpringApplicationBannerPrinter` 的私有静态内部类，它也实现了 `Banner` 接口，添加多个不同的 `Banner` 实现。在 `SpringApplicationBannerPrinter##getBanner` 方法中就能看到，新建 `Banners` 实例，并往其中添加了 `ImageBanner` 和 `ResourceBanner` 。

按照 Banners 的打印顺序，先添加进去的，先打印。

![](banners-getBanner.png)
![](banners-addIfNotNull.png)

我们看看 `ImageBanner` 和 `ResourceBanner` 同时生效的场景：

![](all-banner.png)

运行 `DemoApplication` 入口类，可见如下运行截图：

![](show-all-banner.png)

## 4. 自定义 Banner 接口实现

通过阅读 `SpringApplicationBannerPrinter` 的源码，我们知道如果 `Banners` 中没有 `ResourceBanner` 或者 `ImageBanner` 中的任何一个，就会判断自身的 `fallbackBanner` 变量是否存在，存在则直接返回。而该 `fallbackBanner` 变量实际上是 `SpringApplication` 中的 `banner` 变量。

![](SpringApplicationBannerPrinter.png)
![](SpringApplicationBannerPrinter1.png)
![](SpringApplicationBannerPrinter2.png)

而我们查看 SpringApplication 的源码，可以看到如下方法：

![](setBanner.png)

下面就需要我们来自定义 `Banner` 接口的实现：

```java
/**
 * 自定义 Banner 接口实现
 *
 * @author huazie
 * @version 2.0.0
 * @since 2.0.0
 */
public class CustomBanner implements Banner {

    @Override
    public void printBanner(Environment environment, Class<?> sourceClass, PrintStream out) {
        String author = environment.getProperty("author");

        out.println(" _   _                 _      ");
        out.println("| | | |_   _  __ _ ___(_) ___ ");
        out.println("| |_| | | | |/ _` |_  / |/ _ \\");
        out.println("|  _  | |_| | (_| |/ /| |  __/");
        out.println("|_| |_|\\__,_|\\__,_/___|_|\\___|");
        out.println("                              ");
        out.println(" 作者： " + author);
        out.println();
    }
}
```

接下来，修改入口类 DemoApplication，如下：

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

最后运行 `DemoApplication` 入口类，可见如下运行截图：

![](show-custom-banner.png)
# 总结
本篇 **Huazie** 带大家自定义  **Banner** 信息打印，再次加深了对 **Banner** 信息打印流程的理解。当然，这只是 **Spring Boot** 启动过程中的一个小插曲，后续的博文我们将继续深入讲解 `SpringApplication` 的其他内容，敬请期待！！！ 

