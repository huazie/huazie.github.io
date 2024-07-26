---
title: 【Spring Boot 源码学习】OnWebApplicationCondition 详解
date: 2023-10-06 23:10:28
updated: 2024-02-04 22:11:33
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - OnWebApplicationCondition
  - ConditionalOnWebApplication
---



![](/images/spring-boot-logo.png)

# 引言
上篇博文带大家从 `Spring Boot` 源码深入详解了 **OnBeanCondition**，那本篇也同样从源码入手，带大家深入了解 **OnWebApplicationCondition** 的过滤匹配实现。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 主要内容
本篇我们重点详解 `OnWebApplicationCondition` 的实现，参见如下：

## 1. getOutcomes 方法

鉴于前面博文的了解，我们知道 `OnWebApplicationCondition` 也是 `FilteringSpringBootCondition` 的子类，所以这里同样也是从 `getOutcomes` 方法源码来分析【**Spring Boot 2.7.9**】：

```java
@Order(Ordered.HIGHEST_PRECEDENCE + 20)
class OnWebApplicationCondition extends FilteringSpringBootCondition {

    // ...

    @Override
    protected ConditionOutcome[] getOutcomes(String[] autoConfigurationClasses,
        AutoConfigurationMetadata autoConfigurationMetadata) {
        ConditionOutcome[] outcomes = new ConditionOutcome[autoConfigurationClasses.length];
        for (int i = 0; i < outcomes.length; i++) {
            String autoConfigurationClass = autoConfigurationClasses[i];
            if (autoConfigurationClass != null) {
                outcomes[i] = getOutcome(
                    autoConfigurationMetadata.get(autoConfigurationClass, "ConditionalOnWebApplication"));
            }
        }
        return outcomes;
    }
    // ...
}
```

上述逻辑很容易理解，遍历自动配置数组 `autoConfigurationClasses` ，循环如下：
- 首先，从 `autoConfigurationClasses` 中获取自动配置数据 `autoConfigurationClass`；
- 然后，调用 `AutoConfigurationMetadata` 接口的 `get(String className, String key)` 方法来获取与 `autoConfigurationClass` 关联的名为 `"ConditionalOnWebApplication"` 的条件属性值【即应用类型枚举值】；

  应用类型枚举可以查看 `@ConditionalOnWebApplication` 注解获取，如下所示：
  ![](ConditionalOnWebApplication.png)
- 最后，调用 `getOutcome` 方法，并传入上述获取的应用类型枚举值 `type`：
  ![](ConditionalOnWebApplication-getOutcome.png)      
  - 如果 `type` 是 `SERVLET`， 则判断 `org.springframework.web.context.support.GenericWebApplicationContext` 是否存在；
  如果不存在，则返回一个未满足过滤匹配条件的 `ConditionOutcome` 对象【其中包含 **did not find servlet web application classes** 的信息 】。
  - 如果 `type` 是 `REACTIVE`，则判断 `org.springframework.web.reactive.HandlerResult` 是否存在；
  如果不存在，则返回一个未满足过滤匹配条件的 `ConditionOutcome` 对象【其中包含 **did not find reactive web application classes** 的信息 】。
  - 如果 `org.springframework.web.context.support.GenericWebApplicationContext` 不存在且 `org.springframework.web.reactive.HandlerResult` 也不存在，则返回一个未满足过滤匹配条件的 `ConditionOutcome` 对象【其中包含 **did not find reactive or servlet web application classes** 的信息 】。
  - 如果都存在，则直接返回 `null`。

## 2. getMatchOutcome 方法

同 `OnClassCondition` 一样，`OnWebApplicationCondition` 同样实现了 `FilteringSpringBootCondition` 的父类 `SpringBootCondition` 中的抽象方法 `getMatchOutcome` 方法。

> 有关 `SpringBootCondition` 的介绍，这里不赘述了，请查看笔者的 [【Spring Boot 源码学习】OnClassCondition 详解](../../../../../2023/09/11/spring-boot/spring-boot-sourcecode-onclasscondition/)。

那么，我们进入 `getMatchOutcome` 方法中查看如下源码【**Spring Boot 2.7.9**】：

```java
@Override
public ConditionOutcome getMatchOutcome(ConditionContext context, AnnotatedTypeMetadata metadata) {
    boolean required = metadata.isAnnotated(ConditionalOnWebApplication.class.getName());
    ConditionOutcome outcome = isWebApplication(context, metadata, required);
    if (required && !outcome.isMatch()) {
        return ConditionOutcome.noMatch(outcome.getConditionMessage());
    }
    if (!required && outcome.isMatch()) {
        return ConditionOutcome.noMatch(outcome.getConditionMessage());
    }
    return ConditionOutcome.match(outcome.getConditionMessage());
}
```

我们来分析一下相关逻辑：

- 首先，通过调用 `AnnotatedTypeMetadata` 接口的 `isAnnotated` 方法，判断元数据中是否存在 `@ConditionalOnWebApplication` 注解【当应用程序为 Web 应用程序时，该条件注解用来匹配】。如果返回 `true`，表示元数据中存在指定注解；否则，返回  `false`。

- 然后，调用 `isWebApplication` 方法来获取条件匹配结果 `outcome`【有关内容查看 第 3 小节】；
- 如果 `required` 为 `true`【即存在 `@ConditionalOnWebApplication` 注解】，并且 条件结果不匹配，则返回一个新的 `ConditionOutcome` 对象，标记为不匹配，并带有原始的消息。
- 如果 `required` 为 `false`【即不存在 `@ConditionalOnWebApplication` 注解】，并且 条件结果匹配，则同样返回一个新的 `ConditionOutcome` 对象，标记为不匹配，并带有原始的消息。
- 最后，上述两个条件判断都不满足，则将返回一个匹配的 `ConditionOutcome` 对象，并带有原始的消息。
## 3. isWebApplication 方法

下面，我们进入 `isWebApplication` 方法中：

```java
private ConditionOutcome isWebApplication(ConditionContext context, AnnotatedTypeMetadata metadata,
      boolean required) {
    switch (deduceType(metadata)) {
        case SERVLET:
            return isServletWebApplication(context);
        case REACTIVE:
            return isReactiveWebApplication(context);
        default:
            return isAnyWebApplication(context, required);
    }
}
```

上述的逻辑也很简单：
- 首先，通过 `deduceType` 方法获取可获取的应用类型；查看其源码可知，如果存在 `@ConditionalOnWebApplication` 注解，则获取其对应的 `type` 属性；否则默认返回 `Type.ANY`【即任何 `Web` 应用程序都将匹配】。

  ```java
  private Type deduceType(AnnotatedTypeMetadata metadata) {
      Map<String, Object> attributes = metadata.getAnnotationAttributes(ConditionalOnWebApplication.class.getName());
      if (attributes != null) {
          return (Type) attributes.get("type");
      }
      return Type.ANY;
  }
  ```
- 如果是 `Type.SERVLET`，则调用 `isServletWebApplication` 方法返回条件匹配结果。

- 如果是 `Type.REACTIVE`，则调用 `isReactiveWebApplication` 方法返回条件匹配结果。

- 如果不是上述两个应用类型，则默认调用 `isAnyWebApplication` 方法返回条件匹配结果。

### 3.1 isServletWebApplication 方法

我们直接翻看 `isServletWebApplication` 方法的源码，如下：

```java
private ConditionOutcome isServletWebApplication(ConditionContext context) {
    ConditionMessage.Builder message = ConditionMessage.forCondition("");
    if (!ClassNameFilter.isPresent(SERVLET_WEB_APPLICATION_CLASS, context.getClassLoader())) {
        return ConditionOutcome.noMatch(message.didNotFind("servlet web application classes").atAll());
    }
    if (context.getBeanFactory() != null) {
        String[] scopes = context.getBeanFactory().getRegisteredScopeNames();
        if (ObjectUtils.containsElement(scopes, "session")) {
          return ConditionOutcome.match(message.foundExactly("'session' scope"));
        }
    }
    if (context.getEnvironment() instanceof ConfigurableWebEnvironment) {
        return ConditionOutcome.match(message.foundExactly("ConfigurableWebEnvironment"));
    }
    if (context.getResourceLoader() instanceof WebApplicationContext) {
        return ConditionOutcome.match(message.foundExactly("WebApplicationContext"));
    }
    return ConditionOutcome.noMatch(message.because("not a servlet web application"));
}
```

我们来详细分析一下：

- 首先，检查类加载器中是否存在 `org.springframework.web.context.support.GenericWebApplicationContext` ？
  - 如果没有，那么将返回不匹配的结果，并附带消息 `"did not find servlet web application classes"`。
- 如果条件上下文 `context` 中 `BeanFactory` 不为空，则获取所有注册的 `scope` 名称，并检查其中是否包含 `"session"`。如果包含，则返回匹配的结果，并附带消息 `"found session scope"`。
- 如果条件上下文 `context` 中 `Environment` 是 `ConfigurableWebEnvironment` 的实例，则将返回匹配的结果，并附带消息 `"found ConfigurableWebEnvironment"`。
- 如果条件上下文 `context` 中 `ResourceLoader` 是 `WebApplicationContext` 的实例，那么将返回匹配的结果，并附带消息 `"found WebApplicationContext"`。
- 如果上述的条件都不满足，则最后将返回不匹配的结果，并附带消息 `"not a servlet web application"`。

### 3.2 isReactiveWebApplication 方法

同样，我们也先来查看下 `isReactiveWebApplication` 方法的源码，如下：

```java
private ConditionOutcome isReactiveWebApplication(ConditionContext context) {
    ConditionMessage.Builder message = ConditionMessage.forCondition("");
    if (!ClassNameFilter.isPresent(REACTIVE_WEB_APPLICATION_CLASS, context.getClassLoader())) {
        return ConditionOutcome.noMatch(message.didNotFind("reactive web application classes").atAll());
    }
    if (context.getEnvironment() instanceof ConfigurableReactiveWebEnvironment) {
        return ConditionOutcome.match(message.foundExactly("ConfigurableReactiveWebEnvironment"));
    }
    if (context.getResourceLoader() instanceof ReactiveWebApplicationContext) {
        return ConditionOutcome.match(message.foundExactly("ReactiveWebApplicationContext"));
    }
    return ConditionOutcome.noMatch(message.because("not a reactive web application"));
}
```
通过上述 `isServletWebApplication` 方法中的分析，我们可以很快总结下： 
- 首先，检查类加载器中是否存在 `org.springframework.web.reactive.HandlerResult` ？
  - 如果没有，那么将返回不匹配的结果，并附带消息 `"did not find reactive web application classes"`。
- 如果条件上下文 `context` 中 `Environment` 是 `ConfigurableReactiveWebEnvironment` 的实例，则将返回匹配的结果，并附带消息 `"found ConfigurableReactiveWebEnvironment"`。
- 如果条件上下文 `context` 中 `ResourceLoader` 是 `ReactiveWebApplicationContext` 的实例，那么将返回匹配的结果，并附带消息 `"found ReactiveWebApplicationContext"`。
- 如果上述的条件都不满足，则最后将返回不匹配的结果，并附带消息 `"not a reactive web application"`。

### 3.3 isAnyWebApplication 方法
还是一样，我们先来看看 `isAnyWebApplication` 方法的源码，如下：

```java
private ConditionOutcome isAnyWebApplication(ConditionContext context, boolean required) {
    ConditionMessage.Builder message = ConditionMessage.forCondition(ConditionalOnWebApplication.class,
        required ? "(required)" : "");
    ConditionOutcome servletOutcome = isServletWebApplication(context);
    if (servletOutcome.isMatch() && required) {
        return new ConditionOutcome(servletOutcome.isMatch(), message.because(servletOutcome.getMessage()));
    }
    ConditionOutcome reactiveOutcome = isReactiveWebApplication(context);
    if (reactiveOutcome.isMatch() && required) {
        return new ConditionOutcome(reactiveOutcome.isMatch(), message.because(reactiveOutcome.getMessage()));
    }
    return new ConditionOutcome(servletOutcome.isMatch() || reactiveOutcome.isMatch(),
        message.because(servletOutcome.getMessage()).append("and").append(reactiveOutcome.getMessage()));
}
```
这里就更简单了，总结如下：

- 首先，通过调用 isServletWebApplication 方法获取条件匹配结果；
  - 如果 `Servlet Web` 应用程序的条件结果匹配并且 `required` 为 `true`，则返回一个包含匹配状态和相关消息的 `ConditionOutcome` 对象。
-  接着，通过调用 isReactiveWebApplication 方法获取条件匹配结果；
  - 如果 `Reactive Web` 应用程序的条件结果匹配并且 `required` 为 `true`，则同样返回一个包含匹配状态和相关消息的 `ConditionOutcome` 对象。
- 最后，如果上述两种情况都不满足 或者  `required` 为 `false` ，则返回一个新的 `ConditionOutcome` 对象，它包含 `servletOutcome.isMatch() || reactiveOutcome.isMatch()`  的匹配状态 和 `servletOutcome` 与 `reactiveOutcome` 两者拼接的消息。

# 总结

本篇 **Huazie** 带大家从源码角度深入了解了自动配置过滤匹配子类  `OnWebApplicationCondition` ，至此 **Spring Boot** 中有关自动配置过滤匹配的三个实现已经介绍完毕，当然有关过滤匹配条件的内容还没结束，下一篇笔者将介绍 `@Conditional` 条件注解。
