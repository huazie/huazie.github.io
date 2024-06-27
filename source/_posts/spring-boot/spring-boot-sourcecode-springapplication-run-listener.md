---
title: 【Spring Boot 源码学习】SpringApplication 的 run 方法监听器
date: 2024-04-28 18:01:52
updated: 2024-04-28 18:01:52
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - SpringApplication
  - SpringApplicationRunListeners
  - SpringApplicationRunListener
  - EventPublishingRunListener
  - AvailabilityChangeEvent
---



![](/images/spring-boot-logo.png)


# 一、引言

书接前文[《SpringApplication 的 run 方法核心流程介绍》](../../../../../2024/04/12/spring-boot/spring-boot-sourcecode-springapplication-run-listener/)，**Huazie** 围绕 `SpringApplication` 的  `run` 方法，带大家一起初步了解了 **Spring Boot** 的核心运行流程。其中有关运行流程监听器的内容出现最多，但还未细讲。那么本篇就深入了解下 `SpringApplication` 的 `run` 方法监听器。


<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)


# 二、主要内容
> 注意： 以下涉及 **Spring Boot** 源码 均来自版本 **2.7.9**，其他版本有所出入，可自行查看源码。

## 2.1 SpringApplicationRunListeners 

`SpringApplication` 的 `run(String... args)` 方法中，我们可以看到如下代码：

```c
SpringApplicationRunListeners listeners = getRunListeners(args);
```

这里是获取一个 `SpringApplicationRunListeners` 对象，它管理了一组 `SpringApplicationRunListener` 的监听器集合。

继续查看 `getRunListeners` 方法：

```java
private SpringApplicationRunListeners getRunListeners(String[] args) {
    Class<?>[] types = new Class<?>[] { SpringApplication.class, String[].class };
    return new SpringApplicationRunListeners(logger,
      getSpringFactoriesInstances(SpringApplicationRunListener.class, types, this, args),
      this.applicationStartup);
}
```

可以看到 `getRunListeners` 方法里，直接 `new` 了一个 `SpringApplicationRunListeners` 对象并返回，它的构造函数有三个参数：

![](three-paramters.png)
- `Log log` ：日志对象
- `List<SpringApplicationRunListener> listeners` ：监听器集合
- `ApplicationStartup applicationStartup` ：应用启动指标对象，通过步骤来记录应用程序启动阶段的情况。核心容器及其基础设施组件可以使用 `ApplicationStartup` 来标记应用程序启动期间的步骤，并收集有关执行上下文或它们处理时间的数据。

`getSpringFactoriesInstances(SpringApplicationRunListener.class, types, this, args)` 方法是获取 `SpringApplicationRunListener` 的监听器集合，如果看过笔者前面的系列文章的朋友，应该对该方法并不陌生。

我们进入 `getSpringFactoriesInstances` 方法，查看如下：

![](/images/springboot/getSpringFactoriesInstances.png)

可以看到了如下的代码 ：

```java
SpringFactoriesLoader.loadFactoryNames(type, classLoader);
```

这里是通过 `SpringFactoriesLoader` 类的 `loadFactoryNames` 方法来获取 `META-INF/spring.factories` 中配置 key 为 `org.springframework.boot.SpringApplicationRunListener` 的数据；

![](SpringApplicationRunListener-Config.png)

有关实现类 `EventPublishingRunListener` 请查看 **2.3** 小节。

继续查看 `SpringApplicationRunListeners` 方法，可以看到：

![](SpringApplicationRunListeners.png)

上述标红的方法对应了 **Spring Boot** 运行流程的不同阶段，这些在[《SpringApplication 的 run 方法核心流程介绍》](../../../../../2024/04/12/spring-boot/spring-boot-sourcecode-springapplication-run-listener/)都有介绍过。

以 `starting` 方法为例：

![](starting.png)
![](doWithListeners.png)

`(listener) -> listener.starting(bootstrapContext)` 是 **Java 8** 中引入的 **Lambda** 表达式写法【一种简洁的表示匿名函数（没有名称的函数）的方式】。它表示一个接受 `SpringApplicationRunListener` 类型参数 `listener` 并且不返回任何结果的函数，函数体内部调用了 `listener` 的 `starting` 方法。

`Consumer` 是 **Java 8** 中的一个函数式接口，它表示接受一个输入参数并且不返回结果的操作。也就是说，上述的 **Lambda** 表达式在这里被用来创建 `Consumer` 接口的一个实例。

![](Consumer.png)

由于 `Consumer` 接口只有一个抽象方法 `accept`，上述的 **Lambda** 表达式将自动实现了这个方法。在 **Lambda** 表达式中，`(listener)` 对应 `accept` 方法的参数，而 `-> listener.starting(bootstrapContext)` 则定义了当 `accept` 方法被调用时应该执行的操作。

```java
this.listeners.forEach(listenerAction);
```

这里遍历了监听器集合中的每个监听器，并执行上述 **Lambda** 表达式定义的函数。

![](forEach.png)
## 2.2 SpringApplicationRunListener
`SpringApplicationRunListener` 提供了一系列运行流程中回调的方法，如下图所示：

![](SpringApplicationRunListener.png)

下面来逐一介绍下【其中一些标注了 `@Deprecated` 的方法，即表示当前版本废弃的方法】：

- **`starting`**：当 `run` 方法第一次被执行时，会被立即调用，可用于非常早期的初始化工作。
- **`environmentPrepared`**：当 `environment` 准备完成，在 `ApplicationContext` 创建之前，该方法被调用。
- **`contextPrepared`**：当 `ApplicationContext` 构建完成，资源还未被加载时，该方法被调用。
- **`contextLoaded`**：当 `ApplicationContext` 加载完成，未被刷新之前，该方法被调用。
- **`started`**：当 `ApplicationContext` 刷新并启动之后， `CommandLineRunner` 和 `ApplicationRunner` 未被调用之前，该方法被调用。
- **`ready`**：当应用程序上下文已刷新，并且所有的 `CommandLineRunner` 和 `ApplicationRunner` 都已被调用后，`run` 方法完成之前，该方法被立即调用。
- **`running`**：同 `ready` 方法，在当前版本已被废弃。
- **`failed`**：当应用程序出现错误时，该方法被调用。

## 2.3 实现类 EventPublishingRunListener
`EventPublishingRunListener` 是 **Spring Boot** 中 `SpringApplicationRunListener` 接口的唯一内部实现。
### 2.3.1 成员变量和构造方法



先来看看它的成员变量和构造方法：

![](three-paramters-1.png)

`EventPublishingRunListener` 有三个成员变量：
 
- `SpringApplication application`：当前运行的 `SpringApplication` 实例
- `String[] args`：启动应用程序时的命令参数
- `SimpleApplicationEventMulticaster initialMulticaster`：事件广播器的简单实现，它会将所有事件多播给所有已注册的监听器，由监听器自身决定是否忽略它们不感兴趣的事件。

`EventPublishingRunListener` 的构造方法中初始化上述三个变量之后，就会遍历 `SpringApplication` 中的所有 `ApplicationListener` 实例，并将它们和 `SimpleApplicationEventMulticaster` 进行关联，以便后续将事件多播给所有已注册的监听器。

### 2.3.2 成员方法
![](method.png)
![](method-1.png)
#### 2.3.2.1 不同阶段的事件处理

通过阅读上述源码，可以大致总结一下 **Spring Boot** 启动运行的不同阶段的事件处理流程：
- 首先，**Spring Boot** 应用程序启动的某个阶段，调用 `EventPublishingRunListener` 的某个方法；
- 然后，在这些方法中，将 `application` 参数和 `args` 参数【某些事件还有其他参数】封装到对应的事件中，这些事件都是 `SpringApplicationEvent` 的实现类；
- 接着，通过成员变量 `initialMulticaster` 的 `multicastEvent` 方法对指定事件进行广播 或者 通过当前方法的应用上下文参数 `context` 的 `publishEvent` 方法来对事件进行发布；
- 最后，对应的事件监听器 `ApplicationListener` 被触发，执行相应的业务逻辑。

细心的朋友，可能发现了，上述 `EventPublishingRunListener` 的某些方法是通过成员变量 `initialMulticaster` 的 `multicastEvent` 方法对指定事件进行广播，而某些方法是通过当前方法的应用上下文参数 `context` 的 `publishEvent` 方法来对事件进行发布。

**那这是为啥呢？**

在解答这个问题之前，我们先来看下 `EventPublishingRunListener` 的 `contextLoaded` 方法，如下所示：

![](contextLoaded.png)

大致总结一下 `contextLoaded` 方法的处理逻辑：

- 遍历 `application` 的所有监听器实现类，如果该实现类还实现了 `ApplicationContextAware` 接口，则将上下文信息设置到该监听器内；
- 将 `application` 中的监听器实现类都添加到应用上下文中；
- 通过成员变量 `initialMulticaster` 的 `multicastEvent` 方法对 ApplicationPreparedEvent 事件进行广播。

仔细查看上述源码，我们发现在 `contextLoaded` 方法之前，都是通过 `initialMulticaster` 的 `multicastEvent` 方法进行事件广播的，而在 `contextLoaded` 方法之后均采用当前方法的应用上下文参数 `context` 的 `publishEvent` 方法来对事件进行发布的。

现在可以回答上面的疑问了。因为只有调用 `contextLoaded` 方法之后，应用上下文才算初始化完成，这时才可以通过它的 `publishEvent` 方法来进行事件的发布。


#### 2.3.2.2 可用性状态变化事件
在 `started` 和 `ready` 方法的实现中，我们还看到 `AvailabilityChangeEvent` 的调用。

`AvailabilityChangeEvent` ，即应用程序可用性状态变化事件。任何应用程序组件都可以发送这些事件以更新应用程序的状态。

```java
// started
AvailabilityChangeEvent.publish(context, LivenessState.CORRECT);
```
`started` 方法发布了一个 `LivenessState.CORRECT` 类型的可用性状态变化事件。`LivenessState` 是一个表示可用性状态的枚举类型，其中 `CORRECT` 表示应用程序正在运行，其内部状态是正确的。它还有一个 `BROKEN` 的枚举类型，表示应用程序正在运行，但其内部状态已损坏。

```java
// ready
AvailabilityChangeEvent.publish(context, ReadinessState.ACCEPTING_TRAFFIC);
```
`ready` 方法发布了一个 `ReadinessState.ACCEPTING_TRAFFIC` 类型的可用性状态变化事件。`ReadinessState` 也是一个可用性状态的枚举类型，其中 `ACCEPTING_TRAFFIC` 表示应用程序已准备好接收流量。它还有一个 `REFUSING_TRAFFIC` 的枚举类型，表示应用程序不愿意接收流量。


## 2.4 自定义 SpringApplicationRunListener 

了解了这么多关于 `SpringApplication` 的 `run` 方法监听器的内容，现在让我们来自定义 `SpringApplicationRunListener` 的接口实现看看，如下所示：

```java
import org.springframework.boot.ConfigurableBootstrapContext;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.SpringApplicationRunListener;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.core.env.ConfigurableEnvironment;

import java.time.Duration;

public class DemoRunListener implements SpringApplicationRunListener {

    public DemoRunListener(SpringApplication application, String[] args) {
        System.out.println("DemoRunListener的构造方法被调用");
    }

    @Override
    public void starting(ConfigurableBootstrapContext bootstrapContext) {
        System.out.println("DemoRunListener的 starting 方法被调用");
    }

    @Override
    public void environmentPrepared(ConfigurableBootstrapContext bootstrapContext, ConfigurableEnvironment environment) {
        System.out.println("DemoRunListener的 environmentPrepared 方法被调用");
    }

    @Override
    public void contextPrepared(ConfigurableApplicationContext context) {
        System.out.println("DemoRunListener的 contextPrepared 方法被调用");
    }

    @Override
    public void contextLoaded(ConfigurableApplicationContext context) {
        System.out.println("DemoRunListener的 contextLoaded 方法被调用");
    }

    @Override
    public void started(ConfigurableApplicationContext context, Duration timeTaken) {
        System.out.println("DemoRunListener的 started 方法被调用");
    }

    @Override
    public void ready(ConfigurableApplicationContext context, Duration timeTaken) {
        System.out.println("DemoRunListener的 ready 方法被调用");
    }

    @Override
    public void failed(ConfigurableApplicationContext context, Throwable exception) {
        System.out.println("DemoRunListener的 failed 方法被调用");
    }
}
```

上述 `DemoRunListener` 类定义好了之后，我们就可以在 **resources** 目录下的 **META-INF** 目录下的 **spring.factories** 文件中添加如下配置【**如果对应的目录和文件没有，则需要自行创建**】：

![](DemoRunListener-Config.png)

到这一步，我们就可以启动自己的 **Spring Boot** 项目，运行结果如下截图：

![](result.png)

从上图中，可以看到不同的启动运行阶段，分别打印了不同的日志出来，说明我们自定义的实现类生效了。

# 三、总结

本篇博文 **Huazie** 同大家一起深入分析了 `SpringApplication` 的 `run` 方法监听器，从配置的加载，接口定义，实现类等方面作了详细了解，最后通过自定义 `SpringApplicationRunListener` 接口实现并运行查看，进一步加深了理解。



