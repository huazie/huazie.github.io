---
title: 【Spring Boot 源码学习】自动装配流程源码解析（上）
date: 2023-08-06 22:24:55 
updated: 2024-02-01 17:26:48
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - 自动装配流程
  - 加载自动配置组件
---



![](/images/spring-boot-logo.png)

# 引言
上篇博文，笔者带大家从整体上了解了[AutoConfigurationImportSelector](/2023/07/30/spring-boot/spring-boot-sourcecode-autoconfigurationimportselector/) 自动装配逻辑的核心功能及流程，由于篇幅有限，更加细化的功能及流程详解还没有介绍。本篇开始将从其源码入手，重点解析细化后的自动装配流程源码。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
    <tr>
        <td rowspan="6" align="left"> 
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
</table>

# 主要内容
下面就让我们从 `AutoConfigurationImportSelector` 的 `selectImports` 方法源码入手，开始了解自动装配流程。

下面来看一下 selectImports 的相关源码【Spring Boot 2.7.9】：

```java
@Override
public String[] selectImports(AnnotationMetadata annotationMetadata) {
    if (!isEnabled(annotationMetadata)) {
        return NO_IMPORTS;
    }
    AutoConfigurationEntry autoConfigurationEntry = getAutoConfigurationEntry(annotationMetadata);
    return StringUtils.toStringArray(autoConfigurationEntry.getConfigurations());
}

protected AutoConfigurationEntry getAutoConfigurationEntry(AnnotationMetadata annotationMetadata) {
    if (!isEnabled(annotationMetadata)) {
        return EMPTY_ENTRY;
    }
    AnnotationAttributes attributes = getAttributes(annotationMetadata);
    List<String> configurations = getCandidateConfigurations(annotationMetadata, attributes);
    configurations = removeDuplicates(configurations);
    Set<String> exclusions = getExclusions(annotationMetadata, attributes);
    checkExcludedClasses(configurations, exclusions);
    configurations.removeAll(exclusions);
    configurations = getConfigurationClassFilter().filter(configurations);
    fireAutoConfigurationImportEvents(configurations, exclusions);
    return new AutoConfigurationEntry(configurations, exclusions);
}
```

## 1. 自动配置开关
检查自动配置是否开启的逻辑位于 `AutoConfigurationImportSelector` 的 `selectImports` 方法中的第一段代码。如果开启自动配置，则继续执行后续操作；否则就返回空数组。

检查自动配置是否开启的源码，如下所示：
```java
@Override
public String[] selectImports(AnnotationMetadata annotationMetadata) {
    if (!isEnabled(annotationMetadata)) {
          return NO_IMPORTS;
       }
    // ...省略
}
```

从上面的源码可以看出，这里主要使用了 `isEnabled` 方法来判断自动配置是否开启；其中该方法返回 true，表示开启自动配置，返回 false，表示关闭自动配置。

我们来看一下它的源码，如下所示：

```java
protected boolean isEnabled(AnnotationMetadata metadata) {
    if (getClass() == AutoConfigurationImportSelector.class) {
        return getEnvironment().getProperty(EnableAutoConfiguration.ENABLED_OVERRIDE_PROPERTY, Boolean.class, true);
    }
    return true;
}
```

通过阅读上述 `isEnabled` 方法源码，我们可以看出，如果当前类为 `AutoConfigurationImportSelector`，会从环境中获取 `key` 为`EnableAutoConfiguration.ENABLED_OVERRIDE_PROPERTY` 的配置属性，而笔者前面的系列博文如果大家看过的话，介绍 `EnableAutoConfiguration` 注解时，就说了这个常量的值为 `spring.boot.enableautoconfiguration` 。

![](enabled_override_property.png)


如果该属性的配置获取不到，则默认为 `true`，也就是默认会开启自动配置。如果当前类为其他类，也则默认直接返回 `true`。

如果想覆盖或重置该配置，我们可以在 **application.properties** 或 **application.yml** 中针对 `spring.boot.enableautoconfiguration` 参数进行配置。

以 `application.properties` 配置关闭自动配置为例，如下所示 ：

```
spring.boot.enableautoconfiguration=false
```

## 2. 加载自动配置组件

接下来，我们看到调用 `getAutoConfigurationEntry` 的代码，它是用来封装将被引入的自动配置信息：

```java
AutoConfigurationEntry autoConfigurationEntry = getAutoConfigurationEntry(annotationMetadata);
```

然后我们进入 `getAutoConfigurationEntry` 方法，看到了获取注解属性的逻辑，如下所示：

```java
AnnotationAttributes attributes = getAttributes(annotationMetadata);

// 从 AnnotationMetadata 返回适当的 AnnotationAttributes。默认情况下，此方法将返回 getAnnotationClass() 的属性。
protected AnnotationAttributes getAttributes(AnnotationMetadata metadata) {
    String name = getAnnotationClass().getName();
    AnnotationAttributes attributes = AnnotationAttributes.fromMap(metadata.getAnnotationAttributes(name, true));
    Assert.notNull(attributes, () -> "No auto-configuration attributes found. Is " + metadata.getClassName()
            + " annotated with " + ClassUtils.getShortName(name) + "?");
    return attributes;
}
```

我们启动先前建的 **Spring Boot** 项目的应用类，在 `getAttributes` 方法最后 `return` 处打个断点，我们可以看到如下的逻辑：

![](getAttributes.png)

了解到这，可以开始加载自动配置的组件了，也就是下面的代码：

```java
// 通过 SpringFactoriesLoader 类提供的方法加载类路径中META-INF目录下的
// spring.factories文件中针对 EnableAutoConfiguration 的注解配置类
List<String> configurations = getCandidateConfigurations(annotationMetadata, attributes);
```

我们进入 `getCandidateConfigurations` 方法中， 相关源码如下所示：

```java
protected List<String> getCandidateConfigurations(AnnotationMetadata metadata, AnnotationAttributes attributes) {
    List<String> configurations = new ArrayList<>(
            SpringFactoriesLoader.loadFactoryNames(getSpringFactoriesLoaderFactoryClass(), getBeanClassLoader()));
    ImportCandidates.load(AutoConfiguration.class, getBeanClassLoader()).forEach(configurations::add);
    Assert.notEmpty(configurations,
            "No auto configuration classes found in META-INF/spring.factories nor in META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports. If you "
                    + "are using a custom packaging, make sure that file is correct.");
    return configurations;
}
```

通过阅读上述源码，我们可以看出这里使用了 `SpringFactoriesLoader` 类提供的`loadFactoryNames` 方法来加载的。其中 `loadFactoryNames` 方法的第一个参数为 `getSpringFactoriesLoaderFactoryClass` 方法返回的 `EnableAutoConfiguration.class`。

```java
protected Class<?> getSpringFactoriesLoaderFactoryClass() {
    return EnableAutoConfiguration.class;
}
```

通过翻看 `loadFactoryNames` 方法对应的源码，我们可以知道它是读取的 `META-INF/spring.factories` 中的配置，并且只会读取配置文件中针对自动配置的注册类【即 EnableAutoConfiguration 相关的配置信息】。

`SpringFactoriesLoader` 类的 `loadFactoryNames` 方法相关源码，由于篇幅有限，这里就不贴出来了，感兴趣的朋友可以自行查阅，**Spring Boot** 中后续会有很多地方用到它，比如用于监听的 `Listeners` 和用于过滤的 `Filters`。

实际上 在 **Spring Boot 2.7.9** 版本中， 它自己内部的 `META-INF/spring.factories` 中有关自动配置的注册类的配置信息已经被去除掉了，不过其他外围的 **jar** 中可能有自己的  `META-INF/spring.factories` 文件，它里面也有关于自动配置注册类的配置信息；

另外我们在  `getCandidateConfigurations` 方法中，也看到了另一行代码获取自动配置注册类的信息，如下所示：

```java
ImportCandidates.load(AutoConfiguration.class, getBeanClassLoader()).forEach(configurations::add);
```

这里的代码其实就是读取的如下截图的配置信息【同 `META-INF/spring.factories` 一样，下面的配置也可能存在于不同的 jar 包中 】：

![](autoconfiguration-imports.png)


我们启动先前建的 **Spring Boot** 项目的应用类，在 `getCandidateConfigurations` 方法 `ImportCandidates` 类调用处打个断点，我们可以看到如下的截图【这里 `configurations` 目前还是空数据，说明从 `META-INF/spring.factories` 没有获取到自动配置注册类的相关信息，因为我们只引入了 **Spring Boot** 项目，并且它内部的  `META-INF/spring.factories` 中的确删除了自动配置注册类的相关信息】：

![](getCandidateConfigurations.png)

在 `getCandidateConfigurations` 方法 最后 `return` 处打个断点，我们可以看到如下的截图【这里 `configurations` 中加载的都是自动配置的注册类，也就是 上述 `ImportCandidates##load` 加载的信息，这里读取的是 `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports` 中的配置信息】：

![](getCandidateConfigurations-return.png)

## 3. 自动配置组件去重

因为上述加载的自动配置注册类，默认加载的是 ClassLoader 下面的所有 `META-INF/spring.factories`  或 `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports` 文件中的配置，所以就有可能出现在不同的jar包中有重复配置的情况。

**Spring Boot** 中则使用了 `Java` 中 Set 集合数据不可重复的特点，来实现去重处理，如下所示：

```java
// 对获得的注解配置类集合进行去重处理，防止多个项目引入同样的配置类
configurations = removeDuplicates(configurations);

// 利用 Set 集合数据不可重复的特点，来实现去重处理
protected final <T> List<T> removeDuplicates(List<T> list) {
    return new ArrayList<>(new LinkedHashSet<>(list));
}    
```

# 总结

本篇 **Huazie** 带大家通读了 **Spring Boot** 自动装配逻辑的部分源码，详细分析了加载自动装配的流程，剩下排除和过滤自动配置的流程将在下一篇继续讲解。内容较多，能够看到这的小伙伴，**Huazie** 在这感谢各位的支持。后续我将持续输出有关 **Spring Boot** 源码学习系列的博文，想要及时了解更新的朋友，[关注这里即可](/categories/开发框架-Spring-Boot/)。
