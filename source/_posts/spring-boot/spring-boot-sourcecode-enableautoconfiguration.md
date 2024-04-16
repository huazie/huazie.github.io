---
title: 【Spring Boot 源码学习】@EnableAutoConfiguration 注解
date: 2023-07-22 21:10:13
updated: 2024-01-23 11:23:05
categories:
- 开发框架-Spring Boot
tags:
- Spring Boot
- EnableAutoConfiguration注解
- Import注解
- AutoConfigurationPackage注解
- 自动配置注解
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言

在 **Huazie** 的上篇博文中，我们详细了解了关于 [@SpringBootApplication 注解](/2023/07/16/spring-boot/spring-boot-sourcecode-springbootapplication/)的一些内容，文章最后提到了 `@EnableAutoConfiguration` 注解，用来开启 `Spring Boot` 的自动配置功能，这将是本篇将要重点讲解的内容。

# 主要内容
## 1. @EnableAutoConfiguration 功能解析

我们知道，在没有使用 `Spring Boot` 的情况下，`Bean` 的生命周期都是由 `Spring` 来管理的，并且 `Spring` 是无法自动配置 `@Configuration` 注解的类。`Spring Boot` 的核心功能之一就是根据约定自动管理 `@Configuration` 注解的类 ，其中 `@EnableAutoConfiguration` 注解就是实现该功能的组件之一。

`@EnableAutoConfiguration` 注解 位于 `spring-boot-autoconfigure` 包内，当使用 `@SpringBootApplication` 注解时，它也就会自动生效。
![](spring-boot-autoconfigure.png)

结合上面的内容，我们很容易猜到 `@EnableAutoConfiguration` 注解是用来启动 `Spring` 应用程序上下文时进行自动配置，它会尝试猜测和配置项目可能需要的 `Bean`。

自动配置通常是根据项目中引入的类和已定义的 `Bean` 来实现的。在自动配置过程中，会检查项目的`classpath`（类路径）中引入的类以及项目依赖的 `jar` 包中的组件。

### 1.1 常见的自动配置示例

下面我们来看看，常见的自动配置的示例，如下所示：

- **数据库连接池：** 假设项目中引入了 `Spring Boot` 的 `JDBC Starter` 依赖，它会根据类路径中的相关库（如 `HikariCP、Durid、Tomcat JDBC`等）自动配置数据库连接池。我们只需在配置文件中提供数据库连接的信息，`Spring Boot` 将会自动创建并配置连接池。

- **Web应用程序：** 当引入了 `Spring Boot` 的 `Web Starter` 依赖时，它会自动配置嵌入式的 `Web` 服务器（如 `Tomcat、Jetty、 Undertow`等），并为我们提供默认的 `Web` 应用程序上下文和基本的 `Web` 配置，例如 `Servlet、Filter、Listener` 等。

- **Spring MVC：** 如果在项目中引入了 `Spring MVC` 的相关依赖，`Spring Boot` 会自动配置 **基于注解的控制器**、**视图解析器**、**异常处理** 等，使得开发 `Web` 应用变得更加简单。

- **持久化框架集成：** 当引入了特定的持久化框架（如 `Hibernate、MyBatis` 等）的相关依赖时，`Spring Boot` 会自动配置相应的 `SessionFactory`、**事务管理器** 等组件，以帮助你进行数据库操作。

- **安全框架：** 当引入了 `Spring Security` 的相关依赖时，`Spring Boot` 会自动配置基本的 **安全过滤器链**、**用户认证和授权** 等，提供基本的应用程序安全性。

### 1.2 源码介绍

下面我们来看看 `@EnableAutoConfiguration` 注解的源码【版本：2.7.9】：

```java
/**
 * 启用Spring应用程序上下文的自动配置，尝试猜测和配置可能需要的Bean。
 * 自动配置类通常基于你的类路径和你已定义的Bean来应用。
 * 例如，如果你在类路径中引入了tomcat-embedded.jar，那么很可能希望有一个
 * TomcatServletWebServerFactory（除非你已经定义了自己的ServletWebServerFactory Bean）。
 * 
 * 当使用@SpringBootApplication注解时，上下文的自动配置会自动启用，因此添加此注解没有额外的效果。
 * 
 * 自动配置试图尽可能智能，并且随着你定义更多自己的配置而退避。
 * 你可以手动使用exclude()方法排除任何你不想应用的配置（如果无法访问它们，
 * 则可以使用excludeName()方法）。你还可以通过spring.autoconfigure.exclude属性来排除它们。
 * 自动配置总是在用户自定义的Bean注册之后应用。
 * 
 * 使用@EnableAutoConfiguration注解标注的类所在的包通常具有特殊意义，并且经常被用作"默认"。
 * 例如，在扫描@Entity类时将使用该包。通常建议将@EnableAutoConfiguration（如果你没有使用
 * @SpringBootApplication）放在根包中，以便可以搜索所有子包和类。
 * 
 * 自动配置类是常规的Spring @Configuration Bean。它们是通过ImportCandidates 和
 * SpringFactoriesLoader 机制（针对这个类进行索引）来定位的。通常，自动配置Bean是
 * @Conditional Bean（通常使用@ConditionalOnClass和@ConditionalOnMissingBean注解）。
 *
 * @author Phillip Webb
 * @author Stephane Nicoll
 * @since 1.0.0
 * @see ConditionalOnBean
 * @see ConditionalOnMissingBean
 * @see ConditionalOnClass
 * @see AutoConfigureAfter
 * @see SpringBootApplication
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Inherited
@AutoConfigurationPackage
@Import(AutoConfigurationImportSelector.class)
public @interface EnableAutoConfiguration {

    /**
     * 可以用于覆盖自动配置是否启用的环境属性
     */
    String ENABLED_OVERRIDE_PROPERTY = "spring.boot.enableautoconfiguration";

    /**
     * 排除特定的自动配置类，以使它们永远不会应用
     * @return 要排除的类
     */
    Class<?>[] exclude() default {};

    /**
     * 排除特定的自动配置类名，以使它们永远不会应用
     * @return 要排除的类名
     * @since 1.3.0
     */
    String[] excludeName() default {};

}
```

通过查看源码，我们可以看到 `@EnableAutoConfiguration` 注解提供了一个常量 和 两个成员变量：

- `ENABLED_OVERRIDE_PROPERTY` :  用于覆盖自动配置是否启用的环境属性
- `exclude` ：排除特定的自动配置类
- `excludeName` ：排除特定的自动配置类名

正如前面所说， `@EnableAutoConfiguration` 会尝试猜测并配置你可能需要的 `Bean`，但实际情况如果是我们不需要这些预配置的 `Bean`，那么也可以通过它的两个成员变量 `exclude` 和 `excludeName` 来排除指定的自动配置。

```java
// 通过 @SpringBootApplication 排除 DataSourceAutoConfiguration
@SpringBootApplication(exclude = DataSourceAutoConfiguration.class)
public class DemoApplication {

}
```
或者：

```java
// 通过 @EnableAutoConfiguration 排除 DataSourceAutoConfiguration
@Configuration
@EnableAutoConfiguration(exclude = DataSourceAutoConfiguration.class)
public class DemoConfiguration {
}
```

> **注意：**
> `Spring Boot` 在进行实体类扫描时，会从 `@EnableAutoConfiguration` 注解标注的类所在的包开始扫描。这也是在使用 `@SpringBootApplication` 注解时需要将被注解的类放在顶级 `package` 下的原因，如果放在较低层级，它所在 `package` 的同级或上级中的类就无法被扫描到，从而无法正常使用相关注解（如 `@Entity`）。

从我们上篇博文中新建的 Spring Boot 项目可知，`@SpringBootApplication` 注解通常用于标记 `Spring Boot` 应用程序的入口类。它会自动启用 `Spring Boot` 的自动配置和组件扫描等功能。但是，如果你希望将自动配置应用于其他类，而不是入口类本身，那么你可以将 `@SpringBootApplication` 注解添加到这些类上。同样地，`@EnableAutoConfiguration` 注解也可以用于其他类，而不仅限于入口类。这个注解用于启用`Spring` 的自动配置功能，并根据类路径和已定义的Bean来自动配置应用程序上下文。

因此，在 `Spring Boot` 应用程序中，入口类只是一个用来引导应用程序的类，而真正的自动配置和功能开启是通过 `@SpringBootApplication` 和 `@EnableAutoConfiguration` 注解所用的其他类完成的。

## 2. @Import 注解介绍
从上面 `@EnableAutoConfiguration` 注解的源码可知，`@Import(AutoConfigurationImportSelector.class)` 也是`@EnableAutoConfiguration` 注解的组成部分，这也是自动配置功能的核心实现者。

下面我们重点讲解一下 `@Import` 注解，至于它对应的 `ImportSelector` ，我们将在后续的博文中详细介绍。

`@Import` 注解位于 `spring-context` 项目内，主要提供导入配置类的功能。在后续的学习源码的过程中，我们会发现有大量的 `EnableXXX` 类使用了`@Import` 注解。

下面我们来看一下 `@Import` 的源码【版本 **spring-context-5.3.25**】：

```java
/**
 * 指示导入一个或多个组件类，通常是@Configuration类。
 * 
 * 提供与Spring XML中的<import/>元素相当的功能。允许导入@Configuration类、
 * ImportSelector 和 ImportBeanDefinitionRegistrar 实现，以及普通的组件类
 * （从4.2开始；类似于AnnotationConfigApplicationContext.register）。
 * 
 * 在导入的 @Configuration 类中声明的@Bean定义应该通过@Autowired注入来访问。
 * 可以将bean本身进行自动装配，也可以将声明bean的配置类实例进行自动装配。
 * 后一种方法允许在@Configuration类方法之间进行显式且友好的导航（适用于IDE）。
 * 
 * 可以在类级别或作为元注解进行声明。
 * 
 * 如果需要导入XML或其他非 @Configuration 的bean定义资源，请使用 @ImportResource 注解来实现。
 *
 * @author Chris Beams
 * @author Juergen Hoeller
 * @since 3.0
 * @see Configuration
 * @see ImportSelector
 * @see ImportBeanDefinitionRegistrar
 * @see ImportResource
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface Import {

    /**
     * @Configuration、ImportSelector、ImportBeanDefinitionRegistrar 或常规组件类可以用来进行导入。
     */
    Class<?>[] value();

}
```

上面的源码注释中，已经将 `@Import` 注解上的英文注释翻译成了中文注释，大家可以阅读了解下，这里就不再展开介绍了。

## 3. @AutoConfigurationPackage 注解介绍

细心的朋友可能发现了，在 `@EnableAutoConfiguration` 注解的源码中，还有一个 `@AutoConfigurationPackage` 注解。

那么 `@AutoConfigurationPackage` 注解有啥作用呢？

在解答之前，我们先来看看 @AutoConfigurationPackage 注解的源码【版本 2.7.9 】：

```java
/**
 * 将包注册到 AutoConfigurationPackages 中。
 * 当没有指定基础包或基础包类时，将会注册带有注解类的包。
 *
 * @author Phillip Webb
 * @since 1.3.0
 * @see AutoConfigurationPackages
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Inherited
@Import(AutoConfigurationPackages.Registrar.class)
public @interface AutoConfigurationPackage {

    /**
     * 应该注册到 AutoConfigurationPackages 的基础包。
     * 使用 basePackageClasses 作为基于类型安全的替代方法，而不是基于字符串的包名。
     * 
     * @since 2.3.0
     */
    String[] basePackages() default {};

    /**
     * @AutoConfigurationPackage 提供了一种类型安全的替代方案，用于指定要注册到 AutoConfigurationPackages 的包。
     * 考虑在每个包中创建一个特殊的无操作标记类或接口，除了被此属性引用外，不具备任何其他功能。
     * 
     * @since 2.3.0
     */
    Class<?>[] basePackageClasses() default {};

}
```

通过上述阅读源码，我们可以看到 `@AutoConfigurationPackage` 上有如下这段代码：

```java
@Import(AutoConfigurationPackages.Registrar.class)
```


通过上面的  `@Import` 注解介绍，我们可以知道，这段代码的作用其实就是通过导入`AutoConfigurationPackages.Registrar` 类【其中 `ImportBeanDefinitionRegistrar` 用于存储导入配置的基础包信息】，将基础包及其子包注册到 `AutoConfigurationPackages` 中，以便实现自动配置的功能。

通常情况下，`Spring Boot` 应用程序会将主配置类（例如使用 `@SpringBootApplication` 注解的类）置于根包中。这样做的话，根包会作为默认的扫描路径，用于自动发现和注册 `Spring` 组件（如`@Controller、@Service、@Repository` 等）。

当使用 `@AutoConfigurationPackage` 注解时，它会将指定类所在的包及其子包中的组件自动注册到Spring应用程序上下文中，即自动装配这些组件，从而简化了组件的配置和使用。

# 总结

本篇笔者介绍了 `@EnableAutoConfiguration` 注解的相关功能，当然其中真正实现自动配置功能的核心实现者 `AutoConfigurationImportSelector` 还没有详细说明，那下一篇博文我们将重点对 `AutoConfigurationImportSelector` 的源码进行学习，敬请期待！！！




