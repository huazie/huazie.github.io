---
title: 【Spring Boot 源码学习】深入 FilteringSpringBootCondition
date: 2023-09-08 08:00:00
updated: 2024-02-02 17:18:53
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - FilteringSpringBootCondition
  - ClassNameFilter
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
前两篇博文笔者带大家从源码深入了解了 **Spring Boot** 的自动装配流程，其中自动配置过滤的实现由于篇幅限制，还未深入分析。

那么从本篇开始，Huazie 就带大家走近 `AutoConfigurationImportFilter`，一起从源码解析 `FilteringSpringBootCondition`、`OnBeanCondition`、`OnClassCondition`、`OnWebApplicationCondition` 的实现。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
    <tr>
        <td rowspan="8" align="left" > 
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
</table>

# 主要内容
在开始本篇内容之前，我们再次来回顾一下上篇博文介绍的 `AutoConfigurationImportFilter` 的源码和相关的类图：

```java
@FunctionalInterface
public interface AutoConfigurationImportFilter {
    // 自动配置组件的过滤匹配
    boolean[] match(String[] autoConfigurationClasses, AutoConfigurationMetadata autoConfigurationMetadata);
}
```

![](autoconfigurationimportfilter.png)

通过上面的关联类图，我们可以看到 `AutoConfigurationImportFilter` 接口实际上是由抽象类 `FilteringSpringBootCondition` 来实现的，另外翻看它的源码，该抽象类还定义了一个抽象方法 `getOutcomes` ，然后 `OnBeanCondition`、`OnClassCondition`、`OnWebApplicationCondition` 继承该抽象类，实现 `getOutcomes` 方法，完成实际的过滤匹配操作。

本篇，我们就从源码入手重点介绍 `FilteringSpringBootCondition` ：
## 1.  match 方法

上一篇博文我们已经从 `FilteringSpringBootCondition` 的部分源码进行了分析，它的 `match` 方法主要是调用 `getOutcomes` 方法，并将其返回的结果转换成布尔数组。而这个 `getOutcomes` 方法是过滤匹配的核心功能，由抽象类 FilteringSpringBootCondition 的子类来实现它。

这里再简单回顾一下 `match` 方法的处理逻辑：

```java
@Override
public boolean[] match(String[] autoConfigurationClasses, AutoConfigurationMetadata autoConfigurationMetadata) {
    ConditionEvaluationReport report = ConditionEvaluationReport.find(this.beanFactory);
    // 调用 由子类实现的 getOutcomes 方法，完成实际的过滤匹配操作
    ConditionOutcome[] outcomes = getOutcomes(autoConfigurationClasses, autoConfigurationMetadata);
    boolean[] match = new boolean[outcomes.length];
    // 将 getOutcomes 方法返回结果转换成布尔数组
    for (int i = 0; i < outcomes.length; i++) {
        match[i] = (outcomes[i] == null || outcomes[i].isMatch());
        if (!match[i] && outcomes[i] != null) {
            logOutcome(autoConfigurationClasses[i], outcomes[i]);
            if (report != null) {
                report.recordConditionEvaluation(autoConfigurationClasses[i], this, outcomes[i]);
            }
        }
    }
    return match;
}
```

上述代码中，我们可以看到，将 `getOutcomes` 方法返回结果转换成布尔数组的循环逻辑中有一段代码如下：

```java
match[i] = (outcomes[i] == null || outcomes[i].isMatch());
```
这里是将返回结果转换成布尔值，分别是：
- 如果匹配结果为 `null` ，认为符合匹配要求， 设置 `match[i] = true`；
- 如果匹配结果不为 `null`，并且 匹配对象的 `isMatch == true`，也认为符合匹配要求， 设置 `match[i] = true`；

这个时候，我们就能理解 上篇博文讲到的 **不符合过滤匹配要求，则清空当前的自动配置组件** 的逻辑：

![](filtermethod.png)

当然 `FilteringSpringBootCondition` 内还有其他的内容，这些内容在它的子类中也将使用到，我们先提前了解下，以便后续能更好地理解子类的功能实现。
## 2. ClassNameFilter 枚举类


首先查看 `ClassNameFilter` 枚举类的源码【**Spring Boot 2.7.9**】：

```java
protected enum ClassNameFilter {

    PRESENT {

        @Override
        public boolean matches(String className, ClassLoader classLoader) {
            return isPresent(className, classLoader);
        }

    },

    MISSING {

        @Override
        public boolean matches(String className, ClassLoader classLoader) {
            return !isPresent(className, classLoader);
        }

    };

    abstract boolean matches(String className, ClassLoader classLoader);

    // ....

}
```

`ClassNameFilter` 枚举类包含两个枚举常量，分别是 `PRESENT` 和 `MISSING`；这两个枚举常量都实现了 `ClassNameFilter` 枚举类定义的 `matches` 的抽象方法，其中
-  `PRESENT` 中的 `matches` 返回 `isPresent(className, classLoader);` 
-  `MISSING`  中的 `matches` 返回 `!isPresent(className, classLoader);`

我们继续看 isPresent 方法，分析一下它的功能：

```java
static boolean isPresent(String className, ClassLoader classLoader) {
    if (classLoader == null) {
        classLoader = ClassUtils.getDefaultClassLoader();
    }
    try {
        resolve(className, classLoader);
        return true;
    }
    catch (Throwable ex) {
        return false;
    }
}

protected static Class<?> resolve(String className, ClassLoader classLoader) throws ClassNotFoundException {
    if (classLoader != null) {
        return Class.forName(className, false, classLoader);
    }
    return Class.forName(className);
}
```

上述 `isPresent` 方法的逻辑其实也并不复杂，就是通过类加载器去加载指定的类【即 **className** 字符串对应的类】：
- 如果指定的类加载成功，则直接返回 `true`；
- 如果指定的类加载失败，则要抛出异常，捕获异常后，返回 `false`。

那显然 `ClassNameFilter.PRESENT.matches(className, classLoader)`  **用于校验指定的类是否加载成功**：
- 如果指定的类加载成功，则返回 `true`；
- 如果指定的类加载失败，则返回 `false`。

而 `ClassNameFilter.MISSING.matches(className, classLoader)` **用于校验指定的类是否加载失败**：
- 如果指定的类加载失败，则返回 `true`；
- 如果指定的类加载成功，则返回 `false`。

## 3. filter 方法

继续翻看 `FilteringSpringBootCondition` 源码，还有一个 `filter` 方法需要重点介绍下：
```java
protected final List<String> filter(Collection<String> classNames, ClassNameFilter classNameFilter, ClassLoader classLoader) {
    if (CollectionUtils.isEmpty(classNames)) {
        return Collections.emptyList();
    }
    List<String> matches = new ArrayList<>(classNames.size());
    for (String candidate : classNames) {
        if (classNameFilter.matches(candidate, classLoader)) {
            matches.add(candidate);-
        }
    }
    return matches;
}
```

结合上面的 `ClassNameFilter` 枚举类，我们可以很容易理解上面的代码逻辑。

- 如果 `classNameFilter` 是 `ClassNameFilter.PRESENT`，则 `filter` 方法获取指定的类集合中加载成功的类集合【即匹配成功的类集合】；
- 如果 `classNameFilter` 是 `ClassNameFilter.MISSING`，则 `filter` 方法获取指定的类集合中加载失败的类集合【即匹配失败的类集合】。

# 总结

本篇 **Huazie** 带大家介绍了自动配置过滤匹配的核心父类 `FilteringSpringBootCondition`，这对于笔者后续博文详解它的三个子类【`OnBeanCondition`、`OnClassCondition`、`OnWebApplicationCondition`】非常重要，敬请期待！！！。


