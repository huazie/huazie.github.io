---
title: 【Spring Boot 源码学习】@Conditional 条件注解
date: 2023-10-15 23:19:30
updated: 2024-01-20 17:45:14
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - Conditional
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
前面的博文，Huazie 带大家从 `Spring Boot` 源码深入了解了自动配置类的读取和筛选的过程，然后又详解了**OnClassCondition**、 **OnBeanCondition**、**OnWebApplicationCondition** 这三个自动配置过滤匹配子类实现。

在上述的博文中，我们其实已经初步涉及到了像 `@ConditionalOnClass`、`@ConditionalOnBean`、`@ConditionalOnWebApplication` 这样的条件注解，并且这些条件注解里面，我们都能看到 `@Conditional` 注解。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="12" align="left" > 
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
</table>



# 主要内容
本篇我们重点介绍 `@Conditional` 条件注解，参见如下：

## 1. 初识 @Conditional
我们先来看看 `@Conditional` 注解的源码【**Spring Context 5.3.25**】：

```java
/**
 * 表示组件仅在所有指定条件匹配时才有资格注册。
 * 
 * 条件是在bean定义即将注册之前可以通过编程确定的任何状态（有关详细信息，请参阅Condition）。
 * 
 * @Conditional注解可以以以下任意方式使用：
 *  作为类型级别的注释直接或间接地应用于带有@Component的任何类，包括@Configuration类
 *  作为元注释，用于组合自定义注释标签
 *  作为@Bean方法上的注释级别注解
 * 
 * 如果一个@Configuration类被标记为@Conditional，则该类的所有@Bean方法、@Import注解和@ComponentScan注解都将受到条件约束。
 *
 * @author Phillip Webb
 * @author Sam Brannen
 * @since 4.0
 * @see Condition
 */
@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface Conditional {

  /**
   * 必须匹配才能注册组件的所有条件类
   */
  Class<? extends Condition>[] value();

}
```

翻看上述源码，可以看到 `@Conditional` 条件注解是从 **Spring 4.0** 开始引入的，它表示组件仅在所有指定条件匹配时才有资格注册。比如，当类加载器下存在某个指定的类的时候才会对注解的类进行实例化操作。

它唯一的元素属性是接口 `Condition` 的数组，只有数组中指定的所有 `Condition` 的 `matches` 方法都返回 `true` 的情况下，被注解的类才会被加载。我们前面讲到的 `OnClassCondition` 等类就是 `Condition` 的子类之一。

```java
/**
 * 一个必须匹配才能注册的单个 Condition。
 *
 * <p> 在 bean 定义即将被注册之前立即进行检查，并可以根据在该点可以确定的任何标准自由否决注册。
 *
 * <p> 条件必须遵循与 BeanFactoryPostProcessor 相同的限制，并确保不要与 bean 实例进行交互。
 * 对于与 @Configuration beans交互的更细粒度的控制，请考虑实现 ConfigurationCondition 接口。
 *
 * @author Phillip Webb
 * @since 4.0
 * @see ConfigurationCondition
 * @see Conditional
 * @see ConditionContext
 */
@FunctionalInterface
public interface Condition {

  /**
   * 确定条件是否匹配。
   * @param context 条件上下文
   * @param metadata 正在检查的 AnnotationMetadata 或 MethodMetadata 的元数据
   * @return 如果条件匹配并且可以注册组件，则返回 true；否则返回 false，否决带有注解的组件的注册。
   */
  boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata);

}
```

上述就是 `Condition` 接口的源码，它的 `matches` 方法用来确定条件是否匹配，其中两个参数分别如下：

- `ConditionContext` ：条件上下文，可通过该接口提供的方法来获得 **Spring** 应用的上下文信息，接口定义如下：
  
  ```java
  public interface ConditionContext {

    /**
     * 返回一个 BeanDefinitionRegistry 对象，该对象将包含如果条件匹配时应该持有的bean定义。
     * 如果没有可用的注册表（这种情况很少见：只有当使用 ClassPathScanningCandidateComponentProvider 时才会出现），
     * 则会抛出IllegalStateException异常。
     */
    BeanDefinitionRegistry getRegistry();
  
    /**
     * 返回一个 ConfigurableListableBeanFactory 对象，该对象将包含如果条件匹配时应该持有的bean定义，
     * 或者 如果bean工厂不可用（或者无法向下转型为 ConfigurableListableBeanFactory），则返回null。
     */
    @Nullable
    ConfigurableListableBeanFactory getBeanFactory();
  
    /**
     * 返回当前应用程序正在运行的环境。
     */
    Environment getEnvironment();
  
    /**
     * 返回当前正在使用的资源加载器。
     */
    ResourceLoader getResourceLoader();
  
    /**
     * 返回应该用来加载额外类的 ClassLoader。如果系统类加载器不可访问，则返回null。
     */
    @Nullable
    ClassLoader getClassLoader();

  }
  ```

- `AnnotatedTypeMetadata` ：该接口提供了访问特定类或方法的注解功能，并且不需要加载类，可以用来检查带有 `@Bean` 注解的方法上是否还有其他注解。

  下面我们来查看下它的源码【**spring-core 5.3.25**】:

  ```java
  public interface AnnotatedTypeMetadata {
    // 返回一个MergedAnnotations对象，表示该类型的注解集合。
      MergedAnnotations getAnnotations();
    // 检查是否存在指定名称的注解，如果存在则返回true，否则返回false。
      default boolean isAnnotated(String annotationName) {
          return this.getAnnotations().isPresent(annotationName);
      }
  
    // 下面的方法，都是用来获取指定名称注解的属性值
      @Nullable
      default Map<String, Object> getAnnotationAttributes(String annotationName) {
          return this.getAnnotationAttributes(annotationName, false);
      }
  
      @Nullable
      default Map<String, Object> getAnnotationAttributes(String annotationName, boolean classValuesAsString) {
          MergedAnnotation<Annotation> annotation = this.getAnnotations().get(annotationName, (Predicate)null, MergedAnnotationSelectors.firstDirectlyDeclared());
          return !annotation.isPresent() ? null : annotation.asAnnotationAttributes(Adapt.values(classValuesAsString, true));
      }
  
      @Nullable
      default MultiValueMap<String, Object> getAllAnnotationAttributes(String annotationName) {
          return this.getAllAnnotationAttributes(annotationName, false);
      }
  
      @Nullable
      default MultiValueMap<String, Object> getAllAnnotationAttributes(String annotationName, boolean classValuesAsString) {
          Adapt[] adaptations = Adapt.values(classValuesAsString, true);
          return (MultiValueMap)this.getAnnotations().stream(annotationName).filter(MergedAnnotationPredicates.unique(MergedAnnotation::getMetaTypes)).map(MergedAnnotation::withNonMergedAttributes).collect(MergedAnnotationCollectors.toMultiValueMap((map) -> {
              return map.isEmpty() ? null : map;
          }, adaptations));
      }
  }
  ```

## 2. @Conditional 的衍生注解
在 **Spring Boot** 的 **autoconfigure** 项目中提供了各类基于`@Conditional` 注解的衍生注解，它们均位于 **spring-boot-autoconfigure** 项目的 `org.springframework.boot.autoconfigure.condition` 包下，如下图所示：

![](autoconfigure-condition.png)

上述有好几个条件注解，我们已经接触过了，下面我们再仔细介绍一下：

- `@ConditionalOnBean`：当容器中有指定 **Bean** 的条件下。

- `@ConditionalOnClass`：当 **classpath** 类路径下有指定类的条件下。

- `@ConditionalOnCloudPlatform`：当指定的云平台处于 **active** 状态时。

- `@ConditionalOnExpression`：基于 **SpEL** 表达式的条件判断。

- `@ConditionalOnJava`：基于 **JVM** 版本作为判断条件。

- `@ConditionalOnJndi`：在 **JNDI** 存在的条件下查找指定的位置。

- `@ConditionalOnMissingBean`：当容器里没有指定 **Bean** 的条件。

- `@ConditionalOnMissingClass`：当类路径下没有指定类的条件下。

- `@ConditionalOnNotWebApplication`：当项目不是一个 **Web** 项目的条件下。

- `@ConditionalOnProperty`：当指定的属性有指定的值的条件下。

- `@ConditionalOnResource`：类路径是否有指定的值。

- `@ConditionalOnSingleCandidate`：当指定的 **Bean** 在容器中只有一个，或者有多个但是指定了首选的 **Bean**。
- `@ConditionalOnWarDeployment` ：当应用以 War 包形式部署时（例如在 Tomcat、Jetty 等 Web 服务器中）

- `@ConditionalOnWebApplication`：当项目是一个 **Web** 项目的条件下


如果我们仔细观察这些注解的源码，很快会发现它们其实都组合了`@Conditional` 注解，不同的是它们在注解中指定的条件（`Condition`）不同。

下面我们以前面博文中了解过的 `@ConditionalOnWebApplication` 为例，来对衍生条件注解进行一个简单的分析：

```java
/**
 * 用于条件性地匹配应用程序是否为Web应用程序。默认情况下，任何Web应用程序都会匹配，但可以通过type()属性进行缩小范围。
 *
 * @author Dave Syer
 * @author Stephane Nicoll
 * @since 1.0.0
 */
@Target({ ElementType.TYPE, ElementType.METHOD })
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Conditional(OnWebApplicationCondition.class)
public @interface ConditionalOnWebApplication {
  // 所需的web应用类型
  Type type() default Type.ANY;
  // 可选应用类型枚举
  enum Type { 
    // 任何类型
        ANY,
        // 基于servlet的web应用
        SERVLET,
        // 基于reactive的web应用
        REACTIVE
  }
}
```

通过查看 `@ConditionalOnWebApplication` 注解的源码，我们发现它的确组合了 `@Conditional` 注解，并且指定了对应的 **Condition** 为`OnWebApplicationCondition`。该类继承自 `SpringBootCondition` 并实现 `AutoConfigurationImportFilter` 接口。

有关 `OnWebApplicationCondition` 类的详细介绍，请查看笔者的[《【Spring Boot 源码学习】OnWebApplicationCondition 详解》](/2023/10/06/spring-boot/spring-boot-sourcecode-onwebapplicationcondition/)，

了解了条件类的相关内容后，我们可以用如下图来表示 `Condition` 接口相关功能及实现类：

![](Condition.png)


# 总结
本篇我们介绍 `@Conditional` 条件注解及其衍生注解，至此有关自动配置装配的流程已经基本介绍完毕。

虽然我们从源码角度对自动装配流程有了清晰的认识，但还是不能熟练地运用。那么下篇博文，我们将以 **Spring Boot** 内置的 `http` 编码功能为例来分析一下整个自动配置的过程。

