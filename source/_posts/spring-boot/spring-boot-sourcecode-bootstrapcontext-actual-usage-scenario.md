---
title: 【Spring Boot 源码学习】BootstrapContext的实际使用场景
date: 2024-03-09 10:00:00
updated: 2024-03-09 10:00:00
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - BootstrapContext
  - 

---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
上一篇博文[《BootstrapRegistry 初始化器实现》](/2024/03/02/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer-impl/)，**Huazie** 向大家介绍了如何自定义 `BootstrapRegistryInitializer` 接口实现，并以此来执行自定义的初始化操作【如注册自定义的 **Bean**、添加 **BootstrapContext** 关闭监听器】。其中涉及到了 `BootstrapContext` 的部分使用场景，那本篇就向大家演示下 **Spring Boot** 启动过程中如何使用引用上下文 `BootstrapContext` 及其默认实现 。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="26" align="left" > 
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
</table>

# 三、主要内容
> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。
## 3.1 BootstrapContext
在 [《深入 BootstrapContext 及其默认实现》](/2024/02/25/spring-boot/spring-boot-sourcecode-bootstrapcontext/) 中，**Huazie** 详细介绍了引导上下文 `BootstrapContext` 及其默认实现 `DefaultBootstrapContext`，在继续下面的内容之前，有不知道的朋友们可以去回顾一下，这里不再赘述。
## 3.2 BootstrapRegistry 初始化器实现
在开始讲解 `BootstrapContext` 的实际使用场景之前，我们需要首先通过 **`BootstrapRegistry` 初始化器实现类** 注册自定义的对象，以便后续在实际的场景中通过 `BootstrapContext` 来获取。

这块内容，有需要了解的朋友，请翻看 **Huazie** 的上一篇博文[《BootstrapRegistry 初始化器实现》](/2024/03/02/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer-impl/)，这里不再赘述。
## 3.3 BootstrapContext 的实际使用场景

首先我们需要通过源码来明确下需要 **添加哪些内容，哪些场景和引导上下文有关？**

先来看如下的截图【`SpringApplication##run`】：

![](run-sourcecode.png)

从上面可以看到 `BootstrapContext` 实际上有 **3** 处使用场景，分别是：

- 早期启动时
- 环境配置准备完成时
- 应用上下文准备完成后关闭 `BootstrapContext` 
### 3.3.1 早期启动时
首先我们来看看 **早期启动时** 的源码截图：

![](starting.png)
![](starting1.png)
![](starting2.png)

从上述截图可知，这里多播了 `ApplicationStartingEvent` 事件，我们如果想要监听这个事件，只需要实现对应的事件监听器，并添加到监听器列表 `listeners` 中即可。

下面我们来自定义有关 `ApplicationStartingEvent` 事件的监听器类：

```java
public class DemoStartingListener implements ApplicationListener<ApplicationStartingEvent>, Ordered {

    @Override
    public void onApplicationEvent(ApplicationStartingEvent event) {
        ConfigurableBootstrapContext bootstrapContext = event.getBootstrapContext();
        User user = bootstrapContext.get(User.class);
        System.out.println("启动时获取User：" + user);
        user.setName("Huazie_1");
        user.setAge(19);
    }

    @Override
    public int getOrder() {
        return 1;
    }
}
```
上述代码中的 `User` 类只是一个简单的 **POJO** 对象，这里源码就不列出来了，可以自行定义即可。

**那么上述自定义的监听器该如何添加到监听器列表 `listeners` 中呢？**

通过阅读相关的源码，可总结如下的两种方式：

-  在 `META-INF/spring.factories` 中添加 `org.springframework.context.ApplicationListener` 的配置。这种方式，我们从 [《ApplicationListener 详解》](/2023/12/10/spring-boot/spring-boot-sourcecode-applicationlistener/) 的 **3.2** 小节可见一斑。
- 通过 `SpringApplication` 中的 `addListeners` 方法添加。这里其实在笔者的[《SpringApplication 的定制化介绍》](/2024/01/07/spring-boot/spring-boot-sourcecode-springapplication-customization/)中的 **1.7** 小节也提及过。

有关监听器实现 `Ordered` 接口，这里再次提及下：**它可以用来确保多个监听同一事件的监听器可以按照我们预定的顺序执行。**
### 3.3.2 环境配置准备完成时
接着，我们来看看 **准备环境配置** 的源码截图：

![](environmentPrepared.png)
![](environmentPrepared1.png)
![](environmentPrepared2.png)
![](environmentPrepared3.png)

从上述截图可知，这里显然在环境配置准备完成之后，多播了 `ApplicationEnvironmentPreparedEvent` 事件，我们如果想要监听这个事件，只需要实现对应的事件监听器，并添加到监听器列表 `listeners` 中即可。

下面我们来自定义有关 `ApplicationEnvironmentPreparedEvent` 事件的监听器类：

```java
public class DemoEnvironmentPreparedListener implements ApplicationListener<ApplicationEnvironmentPreparedEvent>, Ordered {
    @Override
    public void onApplicationEvent(ApplicationEnvironmentPreparedEvent event) {
        ConfigurableBootstrapContext bootstrapContext = event.getBootstrapContext();
        if (bootstrapContext.isRegistered(User.class)) {
            User user = bootstrapContext.get(User.class);
            System.out.println("环境准备时获取User：" + user);
            user.setName("Huazie_2");
            user.setAge(20);
        }
    }

    @Override
    public int getOrder() {
        return 1;
    }
}
```

至于该监听器如何添加到监听器列表 `listeners` 中，显然跟 3.4.1 中的 `DemoStartingListener` 一样，等下会通过 `SpringApplication` 进行添加演示。

### 3.3.3 应用上下文准备完成后关闭 `BootstrapContext`
最后，我们看看准备应用上下文的源码截图：

![](close.png)
![](close1.png)

从上述截图中，我们可以看出的确是在应用上下文准备完成后，调用了 `DefaultBootstrapContext` 的 `close` 方法，多播了 `BootstrapContextClosedEvent` 事件。我们如果想要监听这个事件，只需要实现对应的事件监听器，不过添加该监听器就不像 **3.4.1** 和 **3.4.2** 那样了。其实在 **3.2** 小节介绍的[《BootstrapRegistry 初始化器实现》](/2024/03/02/spring-boot/spring-boot-sourcecode-bootstrapregistryinitializer-impl/)中，我们已经介绍了如何添加 **`BootstrapContext` 关闭事件 监听器**，这里就不再赘述。

## 3.4 实际使用演示
`BootstrapContext` 的实际使用场景已经在 **3.3** 中介绍，下面 Huazie 就带大家实操下。

首先，**Spring Boot** 启动类中需要修改如下：

```java
@SpringBootApplication
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication springApplication = new SpringApplication(DemoApplication.class);
        // 关闭 Banner打印
        springApplication.setBannerMode(Banner.Mode.OFF);
        // 添加 BootstrapRegistry初始化器实现
        springApplication.addBootstrapRegistryInitializer(new DemoBootstrapper());
        // 添加 监听器实现
        springApplication.addListeners(new DemoStartingListener(), new DemoEnvironmentPreparedListener());
        springApplication.run(args);
    }
}
```

修改好启动类的代码，我们就可以来运行了，如下截图：

![](result.png)
从上述截图中，我们可以看到 **3.3** 中介绍的 **3** 个实际使用场景，已经全部打印日志信息了，说明定义的监听器已经执行了。
# 四、总结
本篇 **Huazie** 通过介绍 `BootstrapContext` 的实际使用场景，并演示了引导上下文在这些场景的实际使用，加深了大家对于 **Spring Boot** 的启动引导过程的了解，为后续的源码分析打下基础。

后续的博文，**Huazie** 就将从 `SpringApplication` 的 `run` 方法入手，开始介绍 **Spring Boot** 的运行流程，敬请期待！
