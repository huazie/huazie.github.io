---
title: 【Spring Boot 源码学习】OnClassCondition 详解
date: 2023-09-11 21:58:39
updated: 2024-02-04 21:32:13
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - OnClassCondition
  - ConditionalOnClass
  - ConditionalOnMissingClass
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
上篇博文带大家从源码深入了自动配置过滤匹配父类 **FilteringSpringBootCondition**，那么笔者接下来的博文将要介绍它的三个子类 `OnClassCondition`、`OnBeanCondition` 和 `OnWebApplicationCondition` 的实现。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="9" align="left" > 
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
</table>


# 主要内容
话不多说，我们开始本篇的内容，重点详解 `OnClassCondition` 的实现。

## 1. getOutcomes 方法

`OnClassCondition` 也是 `FilteringSpringBootCondition` 的子类，我们首先从 `getOutcomes` 方法源码来分析【**Spring Boot 2.7.9**】：

```java
// OnClassCondition 用于检查是否存在特定类
@Order(Ordered.HIGHEST_PRECEDENCE)
class OnClassCondition extends FilteringSpringBootCondition {

    @Override
    protected final ConditionOutcome[] getOutcomes(String[] autoConfigurationClasses,
        AutoConfigurationMetadata autoConfigurationMetadata) {
        // 如果有多个处理器可用，则拆分工作并在后台线程中执行一半。
        // 使用单个附加线程似乎可以提供最佳性能。
        // 线程越多，情况就越糟。
        if (autoConfigurationClasses.length > 1 && Runtime.getRuntime().availableProcessors() > 1) {
            return resolveOutcomesThreaded(autoConfigurationClasses, autoConfigurationMetadata);
        } else {
            OutcomesResolver outcomesResolver = new StandardOutcomesResolver(autoConfigurationClasses, 0,
                autoConfigurationClasses.length, autoConfigurationMetadata, getBeanClassLoader());
            return outcomesResolver.resolveOutcomes();
        }
    }
    // ...
}
```

上述 `getOutcomes` 方法中，如果有多个处理器可用，则拆分工作并在后台线程中执行一半，使用单个附加线程似乎可以提供最佳性能【不过线程越多，情况就越糟】；否则，直接新建 `StandardOutcomesResolver` 来处理。

## 2. 多处理器拆分处理
先来看看 `resolveOutcomesThreaded` 的源码【**Spring Boot 2.7.9**】：

```java
private ConditionOutcome[] resolveOutcomesThreaded(String[] autoConfigurationClasses,
      AutoConfigurationMetadata autoConfigurationMetadata) {
    int split = autoConfigurationClasses.length / 2;
    OutcomesResolver firstHalfResolver = createOutcomesResolver(autoConfigurationClasses, 0, split, autoConfigurationMetadata);
    OutcomesResolver secondHalfResolver = new StandardOutcomesResolver(autoConfigurationClasses, split,
        autoConfigurationClasses.length, autoConfigurationMetadata, getBeanClassLoader());
    ConditionOutcome[] secondHalf = secondHalfResolver.resolveOutcomes();
    ConditionOutcome[] firstHalf = firstHalfResolver.resolveOutcomes();
    ConditionOutcome[] outcomes = new ConditionOutcome[autoConfigurationClasses.length];
    System.arraycopy(firstHalf, 0, outcomes, 0, firstHalf.length);
    System.arraycopy(secondHalf, 0, outcomes, split, secondHalf.length);
    return outcomes;
}
```

进入 `resolveOutcomesThreaded` 方法，我们可以看到这里主要采用了分半处理的方法来提升处理效率【单个附加线程处理一半数据，主线程处理一半数据】。

我们来仔细分析一下：

- 首先，获取自动配置类数组的一半长度，用于后续分半处理。
- 然后，通过调用 `createOutcomesResolver` 方法【入参表示要处理自动配置类数组的前面一半的数据】创建了一个`OutcomesResolver` 对象 `firstHalfResolver`；进入 **createOutcomesResolver** 方法，我们可以看到这里是先新建了一个 `StandardOutcomesResolver`，然后将其作为构造函数入参，返回一个 `ThreadedOutcomesResolver` 对象，通过翻看代码，发现就是这里面会新启动一个线程来处理数据。
  ```java
  private OutcomesResolver createOutcomesResolver(String[] autoConfigurationClasses, int start, int end,
      AutoConfigurationMetadata autoConfigurationMetadata) {
      OutcomesResolver outcomesResolver = new StandardOutcomesResolver(autoConfigurationClasses, start, end,
          autoConfigurationMetadata, getBeanClassLoader());
      try {
          return new ThreadedOutcomesResolver(outcomesResolver);
      }
      catch (AccessControlException ex) {
          return outcomesResolver;
      }
  }
  ```
  ![](ThreadedOutcomesResolver.png)

- 接着，先新建了一个 `StandardOutcomesResolver`【其构造方法入参表示要处理自动配置类数组的后面一半的数据】，并赋值给 一个 `OutcomesResolver` 对象 `secondHalfResolver`；
- 最后，调用 `firstHalfResolver` 和 `secondHalfResolver` 的 `resolveOutcomes` 方法来处理自动配置类数据，并将处理结果合并到 `outcomes` 中返回。

 通过上面分析，我们发现不论是 单个附加线程处理一半数据，还是 主线程处理一半数据，其核心还是 `StandardOutcomesResolver` 这个类。

## 3. StandardOutcomesResolver 内部类

下面我们来看看内部类 `StandardOutcomesResolver` 中的 resolveOutcomes 方法的实现代码【**Spring Boot 2.7.9**】：

```java
private static final class StandardOutcomesResolver implements OutcomesResolver {
    // ...省略

    @Override
    public ConditionOutcome[] resolveOutcomes() {
        return getOutcomes(this.autoConfigurationClasses, this.start, this.end, this.autoConfigurationMetadata);
    }
}
```

进入 `resolveOutcomes` 方法，我们可以看到这里直接调用了 `getOutcomes` 方法并返回处理结果，如下所示：

```java
    private ConditionOutcome[] getOutcomes(String[] autoConfigurationClasses, int start, int end,
        AutoConfigurationMetadata autoConfigurationMetadata) {
        ConditionOutcome[] outcomes = new ConditionOutcome[end - start];
        for (int i = start; i < end; i++) {
            String autoConfigurationClass = autoConfigurationClasses[i];
            if (autoConfigurationClass != null) {
                String candidates = autoConfigurationMetadata.get(autoConfigurationClass, "ConditionalOnClass");
                if (candidates != null) {
                    outcomes[i - start] = getOutcome(candidates);
                }
            }
        }
        return outcomes;
    }
```

上述逻辑也好理解，那就是遍历并处理自动配置类数组 `autoConfigurationClasses` 在 索引 `start` 到 `end - 1` 之间的数据。其中循环里面：
- 首先，获取要处理的自动配置类 `autoConfigurationClass` ；
- 然后，通过调用 `AutoConfigurationMetadata` 接口的 `get(String className, String key)` 方法来获取与`autoConfigurationClass` 关联的名为 `"ConditionalOnClass"` 的条件属性值。而该 `get` 方法的具体实现可见 `AutoConfigurationMetadataLoader` 类，这个我们在上一篇博文中也提及到，它会加载 `META-INF/spring-autoconfigure-metadata.properties` 中的配置。
![](AutoConfigurationMetadataLoader.png)
  ```java
  final class AutoConfigurationMetadataLoader {
      // ... 省略

      private static class PropertiesAutoConfigurationMetadata implements AutoConfigurationMetadata {
          // ... 省略
          @Override
          public String get(String className, String key) {
              return get(className, key, null);
          }

          @Override
          public String get(String className, String key, String defaultValue) {
              String value = this.properties.getProperty(className + "." + key);
              return (value != null) ? value : defaultValue;
          }
      }
  }
  ```
  通过上述截图和代码，我们可以看到 `AutoConfigurationMetadataLoader` 的内部类`PropertiesAutoConfigurationMetadata` 实现了 `AutoConfigurationMetadata` 接口的具体方法，其中就包含上述用到的 `get(String className, String key)` 方法。
  
  仔细查看 `get` 方法的实现，我们不难发现上述 `getOutcomes` 方法中获取的 `candidates`，其实就是 `META-INF/spring-autoconfigure-metadata.properties` 文件中配置的 `key` 为 **自动配置类名.ConditionalOnClass** 的字符串，而 **value** 为其获得的值。
 
  我们以 `RedisCacheConfiguration` 为例，可以看到如下配置：
 
  ![](RedisCacheConfiguration.png)
 
 - 最后，调用 `getOutcome(String candidates)` 方法来完成最后的过滤匹配工作。
   
   下面来看看相关的源码实现：
  ```java
  private ConditionOutcome getOutcome(String candidates) {
      try {
          if (!candidates.contains(",")) {
              return getOutcome(candidates, this.beanClassLoader);
          }
          for (String candidate : StringUtils.commaDelimitedListToStringArray(candidates)) {
              ConditionOutcome outcome = getOutcome(candidate, this.beanClassLoader);
              if (outcome != null) {
                  return outcome;
              }
          }
      }
      catch (Exception ex) {
          // We'll get another chance later
      }
      return null;
  }
  ```
  如果 `candidates` 不包含逗号，说明只有一个，直接调用 `getOutcome(String className, ClassLoader classLoader)` 返回过滤匹配结果；否则就是包含多个，调用 `StringUtils.commaDelimitedListToStringArray(candidates)` 将逗号分隔的字符串（如candidates）转换为一个字符串数组，然后遍历处理，还是调用 `getOutcome(String className, ClassLoader classLoader)` 过滤匹配结果，如果 `outcome` 不为空，则直接返回 `outcome` 。
  > `StringUtils.commaDelimitedListToStringArray(candidates)`  它会根据逗号来分割输入的字符串，并移除每个元素中的空格。返回的字符串数组包含了被分割后的各个元素
    
  下面我们直接进入 `getOutcome(String className, ClassLoader classLoader)` 方法查看其源码：
  ```java
  private ConditionOutcome getOutcome(String className, ClassLoader classLoader) {
      if (ClassNameFilter.MISSING.matches(className, classLoader)) {
          return ConditionOutcome.noMatch(ConditionMessage.forCondition(ConditionalOnClass.class)
              .didNotFind("required class")
              .items(Style.QUOTE, className));
      }
      return null;
  }
  ```
  我们这里可以看到上面介绍过的 `ClassNameFilter.MISSING` ，它是用于校验指定的类是否加载失败。而这里意思就是如果 `className` 对应的类不存在，则返回没有满足过滤匹配的结果【即 `ConditionOutcome.noMatch.didNotFind` ，其中不存在需要的类】；否则返回 `null`。
  
结合 `FilteringSpringBootCondition` 的介绍，我们知道了 `OnClassCondition` 类 `getOutComes` 方法判断的是 自动配置类关联的 **OnClassCondition** 配置属性对应的类，如果它存在，则后面处理时保留自动配置类；否则，后面会清空自动配置类；

## 4. getMatchOutcome 方法

通过翻看源码，我们其实也可以发现，`OnClassCondition` 类还实现了 `FilteringSpringBootCondition` 的父类 `SpringBootCondition` 中的抽象方法。

如下是 `SpringBootCondition` 类的部分源码【**Spring Boot 2.7.9**】：

```java
public abstract class SpringBootCondition implements Condition {

    // ...

    @Override
    public final boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        String classOrMethodName = getClassOrMethodName(metadata);
        try {
            ConditionOutcome outcome = getMatchOutcome(context, metadata);
            // ...
            return outcome.isMatch();
        }
        catch (NoClassDefFoundError ex) {
            // ...
        }
        catch (RuntimeException ex) {
            // ...
        }
    }

    public abstract ConditionOutcome getMatchOutcome(ConditionContext context, AnnotatedTypeMetadata metadata);

}
```

在 `SpringBootCondition` 中有个最终方法 `matches`，该方法逻辑很简单，就是调用  `getMatchOutcome` 方法获取过滤匹配结果，然后通过 `outcome.isMatch()` 返回过滤匹配结果值【**true：满足过滤匹配    false：不满足过滤匹配**】 


简单了解上述内容之后，我们继续看 `OnClassCondition` 中 `getMatchOutcome` 的完整实现：

```java
@Override
public ConditionOutcome getMatchOutcome(ConditionContext context, AnnotatedTypeMetadata metadata) {
    ClassLoader classLoader = context.getClassLoader();
    ConditionMessage matchMessage = ConditionMessage.empty();
    List<String> onClasses = getCandidates(metadata, ConditionalOnClass.class);
    if (onClasses != null) {
        List<String> missing = filter(onClasses, ClassNameFilter.MISSING, classLoader);
        if (!missing.isEmpty()) {
            return ConditionOutcome.noMatch(ConditionMessage.forCondition(ConditionalOnClass.class)
                .didNotFind("required class", "required classes")
                .items(Style.QUOTE, missing));
        }
        matchMessage = matchMessage.andCondition(ConditionalOnClass.class)
            .found("required class", "required classes")
            .items(Style.QUOTE, filter(onClasses, ClassNameFilter.PRESENT, classLoader));
    }
    List<String> onMissingClasses = getCandidates(metadata, ConditionalOnMissingClass.class);
    if (onMissingClasses != null) {
        List<String> present = filter(onMissingClasses, ClassNameFilter.PRESENT, classLoader);
        if (!present.isEmpty()) {
            return ConditionOutcome.noMatch(ConditionMessage.forCondition(ConditionalOnMissingClass.class)
              .found("unwanted class", "unwanted classes")
              .items(Style.QUOTE, present));
        }
        matchMessage = matchMessage.andCondition(ConditionalOnMissingClass.class)
            .didNotFind("unwanted class", "unwanted classes")
            .items(Style.QUOTE, filter(onMissingClasses, ClassNameFilter.MISSING, classLoader));
    }
    return ConditionOutcome.match(matchMessage);
}
```
上面的逻辑大致可以总结为如下两处：

- 获取自动配置类上的 `ConditionalOnClass` 注解配置的类，然后调用父类 `FilteringSpringBootCondition` 中的 `filter` 方法，获取匹配失败的类集合。
如果匹配失败的类集合不为空，则返回不满足过滤匹配的结果【即 `ConditionOutcome.noMatch.didNotFind`，其中不存在需要的类】

  ```java
  List<String> missing = filter(onClasses, ClassNameFilter.MISSING, classLoader);
  if (!missing.isEmpty()) {
      return ConditionOutcome.noMatch(ConditionMessage.forCondition(ConditionalOnClass.class)
          .didNotFind("required class", "required classes")
          .items(Style.QUOTE, missing));
  }
  ```
  如果匹配失败的集合为空，则添加满足过滤匹配的结果，并返回【即 `ConditionMessage.empty.andCondition.found`，其中找到了需要的类】。

  ```java
  matchMessage = matchMessage.andCondition(ConditionalOnClass.class)
      .found("required class", "required classes")
      .items(Style.QUOTE, filter(onClasses, ClassNameFilter.PRESENT, classLoader));

  return ConditionOutcome.match(matchMessage);
  ```

- 获取自动配置类上的 `ConditionalOnMissingClass` 注解配置的类，然后调用父类 `FilteringSpringBootCondition` 中的 `filter` 方法，获取匹配成功的类集合。
如果匹配成功的类集合不为空，则返回不满足过滤匹配的结果【即 `ConditionOutcome.noMatch.found`，其中存在不想要的类】

  ```java
  List<String> present = filter(onMissingClasses, ClassNameFilter.PRESENT, classLoader);
  if (!present.isEmpty()) {
      // 找到了不想要的类
      return ConditionOutcome.noMatch(ConditionMessage.forCondition(ConditionalOnMissingClass.class)
          .found("unwanted class", "unwanted classes")
          .items(Style.QUOTE, present));
  }
  ```
  如果匹配成功的类集合为空，则添加满足过滤匹配的结果【即 `ConditionMessage.empty.andCondition.didNotFind`，其中没有找到不想要的类】。

  ```java
  matchMessage = matchMessage.andCondition(ConditionalOnMissingClass.class)
      .didNotFind("unwanted class", "unwanted classes")
      .items(Style.QUOTE, filter(onMissingClasses, ClassNameFilter.MISSING, classLoader));

  return ConditionOutcome.match(matchMessage);
  ```


# 总结

本篇 **Huazie** 带大家介绍了自动配置过滤匹配子类 `OnClassCondition`，内容较多，感谢大家的支持；

笔者接下来的博文还将详解 `OnBeanCondition` 和 `OnWebApplicationCondition` 的实现，敬请期待！！！