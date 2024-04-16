---
title: 【Spring Boot 源码学习】自动装配流程源码解析（下）
date: 2023-08-21 08:00:00
updated: 2024-02-01 17:27:01
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - 自动装配流程
  - 排除指定自动配置组件
  - 过滤自动配置组件
  - 触发自动配置事件
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

# 引言
上篇博文，笔者带大家了解了自动装配流程中有关自动配置加载的流程；

本篇将介绍自动装配流程剩余的内容，包含了自动配置组件的排除和过滤、触发自动配置事件。

# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
    <tr>
        <td rowspan="7" align="left"> 
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
</table>

# 主要内容
书接上篇，本篇继续从源码分析自动装配流程：
## 4. 排除指定自动配置组件

如果我们在实际使用时，并不需要其中的某些组件，那就可以通过 `@EnableAutoConfiguration` 注解的 `exclude` 或 `excludeName` 属性来进行有针对性的排除 或者  在Spring Boot 的配置文件进行排除。

下面我们来分析一下排除逻辑的源码：

```java
Set<String> exclusions = getExclusions(annotationMetadata, attributes);

protected Set<String> getExclusions(AnnotationMetadata metadata, AnnotationAttributes attributes) {
    Set<String> excluded = new LinkedHashSet<>();
    // 获取 exclude 属性 配置的 待排除的自动配置组件
    excluded.addAll(asList(attributes, "exclude"));
    // 获取 excludeName 属性 配置的 待排除的自动配置组件
    excluded.addAll(asList(attributes, "excludeName"));
    // 获取 Spring Boot 配置文件中 配置的 待排除的自动配置组件
    excluded.addAll(getExcludeAutoConfigurationsProperty());
    return excluded;
}

protected List<String> getExcludeAutoConfigurationsProperty() {
    Environment environment = getEnvironment();
    if (environment == null) {
        return Collections.emptyList();
    }
    if (environment instanceof ConfigurableEnvironment) {
        Binder binder = Binder.get(environment);
        return binder.bind(PROPERTY_NAME_AUTOCONFIGURE_EXCLUDE, String[].class)
            .map(Arrays::asList)
            .orElse(Collections.emptyList());
    }
    String[] excludes = environment.getProperty(PROPERTY_NAME_AUTOCONFIGURE_EXCLUDE, String[].class);
    return (excludes != null) ? Arrays.asList(excludes) : Collections.emptyList();
}
```

上面的代码也挺好理解，分别从注解属性 exclude 、 excludeName 以及配置文件中获取待排除的自动配置组件。

下面我们来演示一下该如何配置，从而排除我们不需要的自动配置组件：

- 添加注解属性 exclude 和 excludeName
![](addexcludeorexcludename.png)
- 添加配置文件属性
![](addexcludeorexcludename1.png)
- 我们启动先前建的 **Spring Boot** 项目的应用类，分别查看到如下的信息：
![](getExclusions.png)
![](getExclusions1.png)
![](getExclusions2.png)

当上面获取了被排除的自动配置组件之后，需要对待排除的类进行检查，如下所示：

```java
checkExcludedClasses(configurations, exclusions);

private void checkExcludedClasses(List<String> configurations, Set<String> exclusions) {
    List<String> invalidExcludes = new ArrayList<>(exclusions.size());
    for (String exclusion : exclusions) {
        // 如果待排除的自动配置类存在且可以加载
        // 并且已去重过的自动配置组件中不存在该待排除的自动配置类
        if (ClassUtils.isPresent(exclusion, getClass().getClassLoader()) && !configurations.contains(exclusion)) {
            // 添加到非法的排除列表中
            invalidExcludes.add(exclusion);
        }
    }
    // 如果存在非法的排除项，则抛出相应的异常信息
    if (!invalidExcludes.isEmpty()) {
        handleInvalidExcludes(invalidExcludes);
    }
}

protected void handleInvalidExcludes(List<String> invalidExcludes) {
    StringBuilder message = new StringBuilder();
    for (String exclude : invalidExcludes) {
        message.append("\t- ").append(exclude).append(String.format("%n"));
    }
    throw new IllegalStateException(String.format(
            "The following classes could not be excluded because they are not auto-configuration classes:%n%s",
            message));
}
```

上述代码中对于待排除类的检查逻辑也好理解，如果待排除的自动配置类存在且可以加载【即存在于当前的ClassLoader中】，并且已去重过的自动配置组件中不存在该待排除的自动配置类，则认为待排除的自动配置类是非法的，抛出相关异常。

我们下面通过示例来验证一下：

- 在我们的示例项目中添加一个自动配置类【注意这里只做演示，无其他意义】
![](addautoconfiguration.png)


- 配置文件添加项目中的一个自动配置类
![](addautoconfiguration1.png)
- 我们启动先前建的 **Spring Boot** 项目的应用类，可以看到如下的启动异常报错：
![](addautoconfiguration2.png)

如果上述检查通过，则说明待排除的自动配置类都符合要求，则调用如下代码从自动配置集合中移除上面获取的待排除的自动配置类信息。

```java
configurations.removeAll(exclusions);
```

## 5. 过滤自动配置组件

经过上面的自动配置组件排除逻辑之后，接下来就要过滤自动配置组件了，而过滤逻辑主要是通过检查配置类的注解是否符合 `spring.factories` 文件中 `AutoConfigurationImportFilter` 指定的注解检查条件，来决定该过滤哪些自动配置组件。

下面开始分析相关代码，如下所示【**Spring Boot 2.7.9**】：

```java
configurations = getConfigurationClassFilter().filter(configurations);
```
进入 `getConfigurationClassFilter` 方法，如下所示：

```java
private ConfigurationClassFilter getConfigurationClassFilter() {
    if (this.configurationClassFilter == null) {
        List<AutoConfigurationImportFilter> filters = getAutoConfigurationImportFilters();
        for (AutoConfigurationImportFilter filter : filters) {
            invokeAwareMethods(filter);
        }
        this.configurationClassFilter = new ConfigurationClassFilter(this.beanClassLoader, filters);
    }
    return this.configurationClassFilter;
}
```
`getConfigurationClassFilter` 方法返回一个 `ConfigurationClassFilter` 实例，用来过滤掉不必要的配置类。

继续看 `getAutoConfigurationImportFilters` 方法，如下所示：

```java 
protected List<AutoConfigurationImportFilter> getAutoConfigurationImportFilters() {
    return SpringFactoriesLoader.loadFactories(AutoConfigurationImportFilter.class, this.beanClassLoader);
}
```
它通过 `SpringFactoriesLoader` 类的 `loadFactories` 方法来获取 `META-INF/spring.factories` 中配置 `key` 为 `AutoConfigurationImportFilter` 的 `Filters` 列表；

我们可以查看相关配置了解一下，如下所示：

```java
# Auto Configuration Import Filters
org.springframework.boot.autoconfigure.AutoConfigurationImportFilter=\
org.springframework.boot.autoconfigure.condition.OnBeanCondition,\
org.springframework.boot.autoconfigure.condition.OnClassCondition,\
org.springframework.boot.autoconfigure.condition.OnWebApplicationCondition
```

如上所示，在 `spring-boot-autoconfigure` 中默认配置了三个筛选条件：`OnBeanCondition`、`OnClassCondition`、`OnWebApplicationCondition`，它们均实现了 `AutoConfigurationImportFilter` 接口。

相关类图如下所示：

![](autoconfigure-condition.png)

我们继续往下看 **invokeAwareMethods**，如下所示：

```java
private void invokeAwareMethods(Object instance) {
    if (instance instanceof Aware) {
        if (instance instanceof BeanClassLoaderAware) {
            ((BeanClassLoaderAware) instance).setBeanClassLoader(this.beanClassLoader);
        }
        if (instance instanceof BeanFactoryAware) {
            ((BeanFactoryAware) instance).setBeanFactory(this.beanFactory);
        }
        if (instance instanceof EnvironmentAware) {
            ((EnvironmentAware) instance).setEnvironment(this.environment);
        }
        if (instance instanceof ResourceLoaderAware) {
            ((ResourceLoaderAware) instance).setResourceLoader(this.resourceLoader);
        }
    }
}
```
这里先判断传入的 `instance` 对象是否是 `Aware` 接口？

如果是 `Aware` 接口，则判断是否是它的 `BeanClassLoaderAware`、 `BeanFactoryAware`、`EnvironmentAware` 和 `ResourceLoaderAware` 这 4 个子接口实现？ 

如果是，则调用对应的回调方法设置相应参数。

> `Aware` 接口是一个一个标记超接口，它表示一个 `bean` 有资格通过回调方式从 `Spring` 容器中接收特定框架对象的通知。具体的方法签名由各个子接口确定，但通常应该只包括一个接受单个参数并返回 `void` 的方法。


继续往下翻看源码，在 `getConfigurationClassFilter` 方法最后，我们可以看到它返回了一个内部类 `ConfigurationClassFilter` 的实例对象。

有了内部类 `ConfigurationClassFilter` ，接下来就可以开始自动配置组件的过滤操作，主要是通过内部类 `ConfigurationClassFilter` 的 `filter` 方法来实现过滤自动配置组件的功能。


不过在分析 `filter` 方法之前，我们先了解下内部类 `ConfigurationClassFilter` 中两个成员变量 ：
- `List<AutoConfigurationImportFilter> filters`  ： 上面已介绍，它是 `META-INF/spring.factories` 中配置的 **key** 为 `AutoConfigurationImportFilter` 的 `Filters` 列表
- `AutoConfigurationMetadata autoConfigurationMetadata` ：元数据文件 `META-INF/ spring-autoconfigure-metadata.properties` 中配置对应实体类，详细分析请看下面。

`AutoConfigurationMetadata` 自动配置元数据，这个前面没有涉及到，从内部类 `ConfigurationClassFilter` 的构造函数中，我们可以看到如下：

```java
this.autoConfigurationMetadata = AutoConfigurationMetadataLoader.loadMetadata(classLoader);
```

详细代码，由于篇幅受限，这里就不贴了，大家可以自行查看相关源码，从如下的截图中，我们也可以直观了解下。

![](autoconfigurationmetadataloader.png)
![](spring-autoconfigure-metadata.png)

好了，现在我们进入 `filter` 方法中，最关键的就是下面 的双层 for 循环处理：

```java
List<String> filter(List<String> configurations) {
    long startTime = System.nanoTime();
    String[] candidates = StringUtils.toStringArray(configurations);
    boolean skipped = false;
    // 具体的过滤匹配操作
    for (AutoConfigurationImportFilter filter : this.filters) {
        boolean[] match = filter.match(candidates, this.autoConfigurationMetadata);
        for (int i = 0; i < match.length; i++) {
            if (!match[i]) {
                // 不符合过滤匹配要求，则清空当前的自动配置组件
                candidates[i] = null;
                skipped = true;
            }
        }
    }
    // 如果匹配完了，都无需跳过，直接返回当前配置即可
    if (!skipped) {
        return configurations;
    }
    // 有一个不满足过滤匹配要求，都重新处理并返回符合要求的自动配置组件
    List<String> result = new ArrayList<>(candidates.length);
    for (String candidate : candidates) {
        // 如果当前自动配置组件不满足过滤匹配要求，则上面会被清空
        // 因此这里只需判断即可获取符合要求的自动配置组件
        if (candidate != null) {
            result.add(candidate);
        }
    }
    if (logger.isTraceEnabled()) {
        int numberFiltered = configurations.size() - result.size();
        logger.trace("Filtered " + numberFiltered + " auto configuration class in "
                + TimeUnit.NANOSECONDS.toMillis(System.nanoTime() - startTime) + " ms");
    }
    return result;
}
```

翻看上面的 `filter` 方法源码，我们可以很明显地看到，Spring Boot 就是通过如下的代码来实现具体的过滤匹配操作。

```java
boolean[] match = filter.match(candidates, this.autoConfigurationMetadata);
```

在介绍如何实现具体的过滤匹配操作之前，先来看一下 `AutoConfigurationImportFilter` 接口的源码：

```java
@FunctionalInterface
public interface AutoConfigurationImportFilter {
    boolean[] match(String[] autoConfigurationClasses, AutoConfigurationMetadata autoConfigurationMetadata);
}
```

上面的 `match` 方法就是实现具体的过滤匹配操作；

**参数：**
- `String[] autoConfigurationClasses` ：待过滤的自动配置类数组
- `AutoConfigurationMetadata autoConfigurationMetadata` ：自动配置的元数据信息

**返回值：**

过滤匹配后的结果布尔数组，数组的大小与 `autoConfigurationClasses` 一致，如果自动配置组件需过滤掉，则设置布尔数组对应值为 `false`。

结合上面的关联类图，我们可以看到 `AutoConfigurationImportFilter` 接口实际上是由抽象类 `FilteringSpringBootCondition` 来实现的，另外该抽象类还定义了一个抽象方法 `getOutcomes` ，然后 `OnBeanCondition`、`OnClassCondition`、`OnWebApplicationCondition` 继承该抽象类，实现 getOutcomes 方法，完成实际的过滤匹配操作。

抽象类 `FilteringSpringBootCondition` 的相关源码如下【**Spring Boot 2.7.9**】：

```java
abstract class FilteringSpringBootCondition extends SpringBootCondition
        implements AutoConfigurationImportFilter, BeanFactoryAware, BeanClassLoaderAware {

    // 其他代码省略

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

    protected abstract ConditionOutcome[] getOutcomes(String[] autoConfigurationClasses,
            AutoConfigurationMetadata autoConfigurationMetadata);

    // 其他代码省略
}
```

通过上面源码可以看出，抽象类 `FilteringSpringBootCondition` 的 `match` 方法主要是调用 `getOutcomes` 方法，并将其返回的结果转换成布尔数组。而这个  `getOutcomes` 方法是过滤匹配的核心功能，由抽象类 `FilteringSpringBootCondition` 的子类来实现它。

有关 `OnBeanCondition`、`OnClassCondition` 和 `OnWebApplicationCondition` 的内容由于篇幅受限，后续 Huazie 会再通过一篇博文详细讲解。

## 6. 触发自动配置事件

经过上面的排除和过滤之后，我们需要的自动配置类集合已经可以返回了。不过在返回之前，还需要再进行最后一步，触发自动配置导入事件，用来通知所有注册的自动配置监听器进行相关处理。

```java
fireAutoConfigurationImportEvents(configurations, exclusions);
```
进入 `fireAutoConfigurationImportEvents` 方法，可以看到如下源码：

```java
private void fireAutoConfigurationImportEvents(List<String> configurations, Set<String> exclusions) {
    List<AutoConfigurationImportListener> listeners = getAutoConfigurationImportListeners();
    if (!listeners.isEmpty()) {
        AutoConfigurationImportEvent event = new AutoConfigurationImportEvent(this, configurations, exclusions);
        for (AutoConfigurationImportListener listener : listeners) {
            invokeAwareMethods(listener);
            listener.onAutoConfigurationImportEvent(event);
        }
    }
}
```

接着，我们进入 `getAutoConfigurationImportListeners` 方法里，它是通过`SpringFactoriesLoader` 类提供的 `loadFactories` 方法将 `spring.factories` 中配置的接口 `AutoConfigurationImportListener` 的实现类加载出来。

```java
protected List<AutoConfigurationImportListener> getAutoConfigurationImportListeners() {
    return SpringFactoriesLoader.loadFactories(AutoConfigurationImportListener.class, this.beanClassLoader);
}
```

`spring.factories` 中配置的自动配置监听器，如下所示：

```
# Auto Configuration Import Listeners
org.springframework.boot.autoconfigure.AutoConfigurationImportListener=\
org.springframework.boot.autoconfigure.condition.ConditionEvaluationReportAutoConfigurationImportListener
```

然后，将过滤出的自动配置类集合和被排除的自动配置类集合作为入参创建一个 `AutoConfigurationImportEvent` 事件对象；

> 其中 `invokeAwareMethods(listener);`  类似上面的 `invokeAwareMethods(filter);`  这里不再赘述了。

最后，调用上述自动配置监听器的 `onAutoConfigurationImportEvent` 方法，并传入上述获取的 `AutoConfigurationImportEvent` 事件对象，来通知所有注册的监听器进行相应的处理。

**那这样做有什么好处呢？**

通过触发 `AutoConfigurationImportEvent` 事件，来通知所有注册的监听器进行相应的处理，我们就可以在导入自动配置类之后，执行一些附加的自定义逻辑或修改自动配置行为。

# 总结

本篇 **Huazie** 带大家通读了 **Spring Boot** 自动装配逻辑的源码，详细分析了自动装配的后续流程，主要包含 自动配置的排除 和 过滤。超过万字，能够看到这的小伙伴，**Huazie** 在这感谢各位的支持。后续我将持续输出有关 **Spring Boot** 源码学习系列的博文，想要及时了解更新的朋友，[关注这里即可](/categories/开发框架-Spring-Boot/)。
