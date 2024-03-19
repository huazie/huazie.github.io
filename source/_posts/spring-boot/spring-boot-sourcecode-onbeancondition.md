---
title: 【Spring Boot 源码学习】OnBeanCondition 详解
date: 2023-09-21 20:00:00
updated: 2024-02-04 21:50:31
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - OnBeanCondition
  - ConditionalOnBean
  - ConditionalOnSingleCandidate
  - ConditionalOnMissingBean
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
上篇博文带大家从 `Spring Boot` 源码深入详解了 **OnClassCondition**，那本篇也同样从源码入手，带大家深入了解 **OnBeanCondition** 的过滤匹配实现。
# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="10" align="left" > 
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
</table>


# 主要内容

话不多说，马上进入正题，我们开始本篇的内容，重点详解 `OnBeanCondition` 的实现。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
## 1. getOutcomes 方法

`OnBeanCondition` 同样也是 `FilteringSpringBootCondition` 的子类，我们依旧是从 `getOutcomes` 方法源码来分析【**Spring Boot 2.7.9**】：

```java
@Order(Ordered.LOWEST_PRECEDENCE)
class OnBeanCondition extends FilteringSpringBootCondition implements ConfigurationCondition {

  // ...

  @Override
  protected final ConditionOutcome[] getOutcomes(String[] autoConfigurationClasses,
      AutoConfigurationMetadata autoConfigurationMetadata) {
    ConditionOutcome[] outcomes = new ConditionOutcome[autoConfigurationClasses.length];
    for (int i = 0; i < outcomes.length; i++) {
      String autoConfigurationClass = autoConfigurationClasses[i];
      if (autoConfigurationClass != null) {
        Set<String> onBeanTypes = autoConfigurationMetadata.getSet(autoConfigurationClass, "ConditionalOnBean");
        outcomes[i] = getOutcome(onBeanTypes, ConditionalOnBean.class);
        if (outcomes[i] == null) {
          Set<String> onSingleCandidateTypes = autoConfigurationMetadata.getSet(autoConfigurationClass,
              "ConditionalOnSingleCandidate");
          outcomes[i] = getOutcome(onSingleCandidateTypes, ConditionalOnSingleCandidate.class);
        }
      }
    }
    return outcomes;
  }
  // ...
}
```

上述 `getOutcomes` 方法中针对 自动配置数据的循环处理逻辑，大致可总结为如下两种：

- 通过调用 `AutoConfigurationMetadata` 接口的 `getSet(String className, String key)` 方法来获取与`autoConfigurationClass` 关联的名为 **`"ConditionalOnBean"`** 的条件属性值，可能含多个，存入 `Set` 集合 `onBeanTypes` 变量中；接着调用 `getOutcome(Set<String> requiredBeanTypes, Class<? extends Annotation> annotation)` 方法来获取过滤匹配结果，并赋值给 `outcomes[i]`。

  我们以 **RedisCacheConfiguration** 为例，可以看到如下配置：
  ![](RedisCacheConfiguration.png)

- 如果上述过滤匹配结果 `outcomes[i]` 为 `null`，则通过调用 `AutoConfigurationMetadata` 接口的 `getSet(String className, String key)` 方法来获取与`autoConfigurationClass` 关联的名为 **`"ConditionalOnSingleCandidate"`** 的条件属性值，可能含多个，存入 `Set` 集合 `onSingleCandidateTypes` 变量中；接着调用 `getOutcome(Set<String> requiredBeanTypes, Class<? extends Annotation> annotation)` 方法来获取过滤匹配结果，并赋值给 `outcomes[i]`。

  我们以 **MongoDatabaseFactoryConfiguration** 为例，可以看到如下配置：
  ![](MongoDatabaseFactoryConfiguration.png)

> 有关 `AutoConfigurationMetadata` 接口的 `get(String className, String key)` 方法的逻辑，请查看 **Huazie** 的 上一篇博文[【Spring Boot 源码学习】OnClassCondition 详解](/2023/09/11/spring-boot/spring-boot-sourcecode-onclasscondition/)，这里不再赘述。

下面我们继续查看 `getOutcome(Set<String> requiredBeanTypes, Class<? extends Annotation> annotation)` 方法的逻辑：

```java
private ConditionOutcome getOutcome(Set<String> requiredBeanTypes, Class<? extends Annotation> annotation) {
  List<String> missing = filter(requiredBeanTypes, ClassNameFilter.MISSING, getBeanClassLoader());
  if (!missing.isEmpty()) {
    ConditionMessage message = ConditionMessage.forCondition(annotation)
      .didNotFind("required type", "required types")
      .items(Style.QUOTE, missing);
    return ConditionOutcome.noMatch(message);
  }
  return null;
}
```

进入 `getOutcome` 方法，可以看到：
- 首先调用父类 `FilteringSpringBootCondition` 中的 `filter` 方法，来获取给定的类集合 `requiredBeanTypes` 中加载失败的类集合 `missing`【即当前类加载器中不存在的类集合】；
- 如果 `missing` 不为空，说明存在加载失败的类，则返回 不满足过滤匹配的结果【即 ConditionOutcome.noMatch，其中没有找到 `missing` 中需要的类型】；
- 如果  `missing` 为空，直接返回 null 即可。

## 2. getMatchOutcome 方法

同 `OnClassCondition` 一样，`OnBeanCondition` 同样实现了 `FilteringSpringBootCondition` 的父类 `SpringBootCondition` 中的抽象方法 `getMatchOutcome` 方法。

> 有关 `SpringBootCondition` 的介绍，这里不赘述了，请查看笔者的 [【Spring Boot 源码学习】OnClassCondition 详解](/2023/09/11/spring-boot/spring-boot-sourcecode-onclasscondition/)。

通过查看 `getMatchOutcome` 方法源码，可以看到针对 `ConditionalOnBean` 注解、`ConditionalOnSingleCandidate` 注解 和 `ConditionalOnMissingBean` 注解的三块处理逻辑，下面来一一讲解：

```java
@Override
public ConditionOutcome getMatchOutcome(ConditionContext context, AnnotatedTypeMetadata metadata) {
  ConditionMessage matchMessage = ConditionMessage.empty();
  MergedAnnotations annotations = metadata.getAnnotations();
  // ConditionalOnBean 注解处理
  // ConditionalOnSingleCandidate 注解处理
  // ConditionalOnMissingBean 注解处理
  return ConditionOutcome.match(matchMessage);
}
```

### 2.1 ConditionalOnBean 注解处理
我们来看看 `ConditionalOnBean` 注解处理逻辑的源码：

```java
  if (annotations.isPresent(ConditionalOnBean.class)) {
    Spec<ConditionalOnBean> spec = new Spec<>(context, metadata, annotations, ConditionalOnBean.class);
    MatchResult matchResult = getMatchingBeans(context, spec);
    if (!matchResult.isAllMatched()) {
      String reason = createOnBeanNoMatchReason(matchResult);
      return ConditionOutcome.noMatch(spec.message().because(reason));
    }
    matchMessage = spec.message(matchMessage)
      .found("bean", "beans")
      .items(Style.QUOTE, matchResult.getNamesOfAllMatches());
  }
```

针对上述代码，且听分析如下：

- 首先调用 `MergedAnnotations` 接口的  `isPresent(Class<A> annotationType)` 方法判断指定的注解类型是直接存在或者元存在【这里相当于调用 `get(annotationType).isPresent()`】，如果返回 true，表示存在指定的注解类型。
- 如果存在 `@ConditionalOnBean`，则
  - 创建一个条件规范 `Spec` 对象，该类是从底层的注解中提取的搜索规范；
  - 接着，调用 `getMatchingBeans` 方法，并从上下文【`context`】中获取与条件规范【`spec`】匹配的 **Spring Beans** 的结果【`MatchResult`】；
  - 然后，检查匹配结果，如果不是所有的条件都匹配，则继续如下：
    - 调用 `createOnBeanNoMatchReason` 方法，创建一个描述条件不匹配原因的字符串并返回；
    - 返回一个表示未匹配条件的 `ConditionOutcome` 对象【其中包含了条件规范的消息以及不匹配的原因】；
  - 否则，更新匹配消息，并记录 找到了所有匹配的 **Spring Beans**。


### 2.2 ConditionalOnSingleCandidate 注解处理
我们继续查看 `ConditionalOnSingleCandidate` 注解处理逻辑的源码：

```java
  if (metadata.isAnnotated(ConditionalOnSingleCandidate.class.getName())) {
    Spec<ConditionalOnSingleCandidate> spec = new SingleCandidateSpec(context, metadata, annotations);
    MatchResult matchResult = getMatchingBeans(context, spec);
    if (!matchResult.isAllMatched()) {
      return ConditionOutcome.noMatch(spec.message().didNotFind("any beans").atAll());
    }
    Set<String> allBeans = matchResult.getNamesOfAllMatches();
    if (allBeans.size() == 1) {
      matchMessage = spec.message(matchMessage).found("a single bean").items(Style.QUOTE, allBeans);
    }
    else {
      List<String> primaryBeans = getPrimaryBeans(context.getBeanFactory(), allBeans,
          spec.getStrategy() == SearchStrategy.ALL);
      if (primaryBeans.isEmpty()) {
        return ConditionOutcome
          .noMatch(spec.message().didNotFind("a primary bean from beans").items(Style.QUOTE, allBeans));
      }
      if (primaryBeans.size() > 1) {
        return ConditionOutcome
          .noMatch(spec.message().found("multiple primary beans").items(Style.QUOTE, primaryBeans));
      }
      matchMessage = spec.message(matchMessage)
        .found("a single primary bean '" + primaryBeans.get(0) + "' from beans")
        .items(Style.QUOTE, allBeans);
    }
  }
```

同样针对上述代码，跟着 **Huazie** 来一步步分析下：
- 首先调用 `AnnotatedTypeMetadata` 接口的 `isAnnotated(String annotationName)` 方法判断元数据中是否存在指定注解。如果返回 `true`，表示元数据中存在指定注解。
- 如果元数据中存在 `@ConditionalOnSingleCandidate` 注解，则
  - 创建了一个 `SingleCandidateSpec` 的对象 `spec` ，并传入上下文 【`context`】、元数据 【`metadata`】 和注解信息 【`annotations`】 ，该类是专门针对 `@ConditionalOnSingleCandidate` 注解的条件规范。
  - 接着调用 `getMatchingBeans` 方法对 `context` 中的所有 `bean` 进行匹配，并将与条件规范【`spec`】匹配的 **Spring Beans** 的结果存储在 `matchResult` 变量中；
  - 如果没有匹配的 `bean`，则返回表示未匹配条件的 `ConditionOutcome` 对象【其中记录了 **没有找到任何 `bean`** 的信息】；
  - 否则，获取匹配的所有 `bean` 名称并存储在 `allBeans` 变量中。
    - 如果仅有一个匹配的 `bean`，则更新匹配消息，并记录找到了 单个 `bean` 的信息；
    - 否则，获取首选 `bean` 名称列表，并检查列表是否为空；
      - 如果列表为空，则返回表示未匹配条件的 `ConditionOutcome` 对象【其中记录了 **一个首选 `bean` 也没有找到** 的信息】；
      - 如果首选 `bean` 名称列表包含多个 `bean`，则返回表示未匹配条件的 `ConditionOutcome` 对象【其中记录了 **找到了多个首选 `bean`** 的信息】；
      - 否则，更新匹配消息，并记录 **找到了首选 `bean`** 的信息。

### 2.3 ConditionalOnMissingBean 注解处理
我们继续查看 `ConditionalOnMissingBean` 注解处理逻辑的源码：

```java
  if (metadata.isAnnotated(ConditionalOnMissingBean.class.getName())) {
    Spec<ConditionalOnMissingBean> spec = new Spec<>(context, metadata, annotations,
        ConditionalOnMissingBean.class);
    MatchResult matchResult = getMatchingBeans(context, spec);
    if (matchResult.isAnyMatched()) {
      String reason = createOnMissingBeanNoMatchReason(matchResult);
      return ConditionOutcome.noMatch(spec.message().because(reason));
    }
    matchMessage = spec.message(matchMessage).didNotFind("any beans").atAll();
  }
```

经过上述两种处理逻辑的分析，相信大家应该可以看懂第三种处理逻辑的分析：

- 首先调用 `AnnotatedTypeMetadata` 接口的 `isAnnotated(String annotationName)` 方法判断元数据中是否存在指定注解。如果返回 `true`，表示元数据中存在指定注解。
- 如果存在 `@ConditionalOnMissingBean` 注解，则
  - 创建一个条件规范 `Spec` 对象，该类是从底层的注解中提取的搜索规范；
  - 接着，调用 `getMatchingBeans` 方法，并从上下文【`context`】中获取与条件规范【`spec`】匹配的 **Spring Beans** 的结果【`MatchResult`】；
  - 如果存在任何一个匹配的 `bean`，则
    - 调用 `createOnMissingBeanNoMatchReason` 方法，创建一个描述条件不匹配原因的字符串并返回；
    - 返回一个表示未匹配条件的 `ConditionOutcome` 对象【其中包含了条件规范的消息以及不匹配的原因】；
  -  否则，更新匹配消息，并记录 **找不到指定类型的 `bean`** 的信息。



## 3. getMatchingBeans 方法
上述三种注解处理逻辑中，我们都看到了调用 `getMatchingBeans` 方法，下面重点来讲解一下:

```java
protected final MatchResult getMatchingBeans(ConditionContext context, Spec<?> spec) {
  // ...
}
```
我们可以看到 `getMatchingBeans` 方法，有两个参数，它们分别是 上下文 【`context`】和 条件规范【`spec`】；

继续看 `getMatchingBeans` 方法内部逻辑：

```java
  ClassLoader classLoader = context.getClassLoader();
  ConfigurableListableBeanFactory beanFactory = context.getBeanFactory();
```
这里从上下文【`context`】中获取 `ClassLoader` 和 `ConfigurableListableBeanFactory` ；

> **知识拓展：**
> - `ClassLoader` 是 **Java** 中的一个接口，用于加载类。它是 **Java** 类加载机制的核心部分，负责将 **.class** 文件转换为 **Java** 类实例。`ClassLoader` 可以从不同的来源（如文件系统、网络、数据库等）加载类，也可以实现自定义的类加载逻辑。
> -  `ConfigurableListableBeanFactory` 是 **Spring** 框架中的一个核心接口，它扩展了`ListableBeanFactory` 接口，提供了更多的配置和扩展功能。它是一个 `bean` 工厂的抽象概念，用于管理 **Spring** 容器中的 `bean` 对象。`ConfigurableListableBeanFactory` 提供了添加、移除、注册和查找 `bean` 的方法，以及设置和获取 `bean` 属性值的功能。它还支持`bean` 的后处理和事件传播。


```java
  boolean considerHierarchy = spec.getStrategy() != SearchStrategy.CURRENT;
```
这里根据 `Spec` 对象的 `SearchStrategy` 属性来确定是否考虑 `bean` 的层次结构。如果 `SearchStrategy` 是 `CURRENT`，则不考虑层次结构【即 `considerHierarchy 为 false`】；否则，考虑层次结构【即 `considerHierarchy 为 true`】。


```java
  Set<Class<?>> parameterizedContainers = spec.getParameterizedContainers();
```
这里获取 `Spec` 对象的 `parameterizedContainers` 属性，这是一个包含参数化容器类型的集合

```java
  if (spec.getStrategy() == SearchStrategy.ANCESTORS) {
    BeanFactory parent = beanFactory.getParentBeanFactory();
    Assert.isInstanceOf(ConfigurableListableBeanFactory.class, parent,
        "Unable to use SearchStrategy.ANCESTORS");
    beanFactory = (ConfigurableListableBeanFactory) parent;
  }
```
如果 `Spec` 对象的 `SearchStrategy` 属性是 `SearchStrategy.ANCESTORS`，则调用 `getParentBeanFactory` 方法获取其父工厂，并将其转换为 `ConfigurableListableBeanFactory` 类型。

```java
  MatchResult result = new MatchResult();
```
新建一个 `MatchResult` 对象，用于存储匹配结果；

```java
  Set<String> beansIgnoredByType = getNamesOfBeansIgnoredByType(classLoader, beanFactory, considerHierarchy,
        spec.getIgnoredTypes(), parameterizedContainers);
```
调用 `getNamesOfBeansIgnoredByType` 方法，获取被忽略类型的 `bean` 名称集合 `beansIgnoredByType` ；

```java
  for (String type : spec.getTypes()) {
    Collection<String> typeMatches = getBeanNamesForType(classLoader, considerHierarchy, beanFactory, type,
        parameterizedContainers);
    Iterator<String> iterator = typeMatches.iterator();
    while (iterator.hasNext()) {
      String match = iterator.next();
      if (beansIgnoredByType.contains(match) || ScopedProxyUtils.isScopedTarget(match)) {
        iterator.remove();
      }
    }
    if (typeMatches.isEmpty()) {
      result.recordUnmatchedType(type);
    }
    else {
      result.recordMatchedType(type, typeMatches);
    }
  }
```
遍历 `Spec` 对象的 `types` 属性，它是一个 `Set<String>` 集合
- 首先，针对每个类型 `type`，调用 `getBeanNamesForType` 方法获取匹配的 `bean` 名称集合 `typeMatches` 。
- 然后，使用迭代器遍历这个集合，如果集合中的某个元素在被忽略类型的集合中，就将其从迭代器中移除。
- 最后，如果 `typeMatches` 集合为空，则记录未匹配的类型；否则，记录匹配的类型。

```java
  for (String annotation : spec.getAnnotations()) {
    Set<String> annotationMatches = getBeanNamesForAnnotation(classLoader, beanFactory, annotation,
        considerHierarchy);
    annotationMatches.removeAll(beansIgnoredByType);
    if (annotationMatches.isEmpty()) {
      result.recordUnmatchedAnnotation(annotation);
    }
    else {
      result.recordMatchedAnnotation(annotation, annotationMatches);
    }
  }
```
遍历 `Spec` 对象的 `annotations` 属性：
- 首先，针对每个注解 `annotation`，调用 `getBeanNamesForAnnotation` 方法获取匹配的 `bean` 名称集合 `annotationMatches` 。
- 然后，从 `annotationMatches`  集合中移除被忽略类型的集合。
- 最后，如果 `annotationMatches` 集合为空，则记录未匹配的注解；否则，记录匹配的注解。


```java
  for (String beanName : spec.getNames()) {
    if (!beansIgnoredByType.contains(beanName) && containsBean(beanFactory, beanName, considerHierarchy)) {
      result.recordMatchedName(beanName);
    }
    else {
      result.recordUnmatchedName(beanName);
    }
  }
```

遍历 `Spec` 对象的 `names` 属性，对于每个 `bean` 名称，如果它不在被忽略类型的集合中，并且它在 `bean` 工厂中存在，就记录匹配的名称；否则，记录未匹配的名称。

# 总结

本篇 **Huazie** 带大家介绍了自动配置过滤匹配子类 `OnBeanCondition` ，内容较多，感谢大家的支持；笔者接下来的博文还将详解 `OnWebApplicationCondition` 的实现，敬请期待！！！
