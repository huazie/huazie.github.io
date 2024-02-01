---
title: 【Spring Boot 源码学习】@SpringBootApplication 注解
date: 2023-07-16 22:56:45
updated: 2024-01-15 22:40:10 
categories:
- 开发框架-Spring Boot
tags:
- Spring Boot
- SpringBootApplication注解
- 自动配置
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言
在 Huazie 前面的博文 [《Spring Boot 核心运行原理介绍》](/2023/07/13/spring-boot/spring-boot-core-operating-principle/)中，我们初步了解了 `Spring Boot` 核心运行原理，知道了 `@EnableAutoConfiguration` 是用来开启自动配置的注解。但创建过 Spring Boot 项目的读者肯定会说，我们并没有直接看到这个注解，实际上前面我也提到，它是由组合注解 `@SpringBootApplication` 引入的。

至于  `@EnableAutoConfiguration` 的讲解，我将放到后面再深入源码了解，本篇先介绍 组合注解 `@SpringBootApplication`。

好了，废话不多说，看下面介绍：
# 主要内容
## 1. 创建 Spring Boot 项目

首先，我们打开 `Intellij IDEA` 开发工具，选择 `File` -> `New` -> `Project`

![](new-project.png)
然后打开 `New Project` ，选择 `Spring Initializr`【这是用来创建 `Spring Boot` 项目】：

![](new-project1.png)

选择 `Next`，打开 `Project Metadata`【这里可以配置项目的一些基础信息】：

![](new-project2.png)


继续 `Next`，打开 `Dependencies`【有需要其他依赖可以添加】：

![](new-project3.png)

继续 `Next`，选择你的 `Spring Boot` 项目位置：

![](new-project4.png)

最后点 **Finish** 完成。

> **注意：** 刚才新建的项目中，如果因为 `Spring Boot` 的版本问题导致项目报错，那就换个版本再试试。

## 2. Spring Boot 入口类

上述新建好的项目创建完成默认会生成一个 `XXXApplication` 的入口类。默认情况下，这个类的命名规则都是 **artifactId + Application**【**artifactId** 首字母大写】：

![](demo-artifactId.png)
![](demo-main-class.png)
如上图中的  `DemoApplication` 就是我们这里 `Spring Boot` 项目的入口类，代码如下所示：

```java
package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }

}
```

在上述的 `Spring Boot` 入口类 `DemoApplication` 中，唯一的注解就是 `@SpringBootApplication`。前面我们也提到了，它是通过其内部组合的 `@EnableAutoConfiguration` 来开启自动配置功能。

下面我们来详细介绍下 `@SpringBootApplication` 注解：
## 3. @SpringBootApplication 介绍


先来看 `@SpringBootApplication` 的源码【**版本：2.7.9**】：


```java

@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Inherited
@SpringBootConfiguration
@EnableAutoConfiguration
@ComponentScan(excludeFilters = { 
		@Filter(type = FilterType.CUSTOM, classes = TypeExcludeFilter.class),
		@Filter(type = FilterType.CUSTOM, classes = AutoConfigurationExcludeFilter.class) })
public @interface SpringBootApplication {

	/**
	 * 排除特定的自动配置类，使它们永远不会被应用
	 * @return 要排除（禁用）的类
	 */
	@AliasFor(annotation = EnableAutoConfiguration.class)
	Class<?>[] exclude() default {};

	/**
	 * 排除特定的自动配置类名称，以确保它们永远不会被应用
	 * @return 要排除的自动配置类名称
	 * @since 1.3.0
	 */
	@AliasFor(annotation = EnableAutoConfiguration.class)
	String[] excludeName() default {};

	/**
	 * 用于扫描带注解组件的基础包。使用 {@link #scanBasePackageClasses} 可以使用类型安全的方式替代基于字符串的包名。
	 * <p>
	 * 注意：该设置仅对 {@code @ComponentScan} 注解有效，不影响 {@code @Entity} 扫描或 Spring Data 的 {@link Repository} 扫描。
	 * 对于这些情况，你应该添加 {@link org.springframework.boot.autoconfigure.domain.EntityScan @EntityScan} 和
	 * {@code @Enable...Repositories} 注解。
	 * @return 要进行扫描的基础包
	 * @since 1.3.0
	 */
	@AliasFor(annotation = ComponentScan.class, attribute = "basePackages")
	String[] scanBasePackages() default {};

	/**
	 * 用于指定要扫描带注解组件的包的类型安全方式。将扫描指定类所在的包。
	 * <p>
	 * 考虑在每个包中创建一个特殊的空类或接口，只用于作为此属性引用的标记类。
	 * <p>
	 * 注意：该设置仅对 {@code @ComponentScan} 注解有效，不影响 {@code @Entity} 扫描或 Spring Data 的 {@link Repository} 扫描。
	 * 对于这些情况，你应该添加 {@link org.springframework.boot.autoconfigure.domain.EntityScan @EntityScan} 和
	 * {@code @Enable...Repositories} 注解。
	 * @return 要进行扫描的基础包
	 * @since 1.3.0
	 */
	@AliasFor(annotation = ComponentScan.class, attribute = "basePackageClasses")
	Class<?>[] scanBasePackageClasses() default {};

	/**
	 * 用于在 Spring 容器中为检测到的组件命名的 {@link BeanNameGenerator} 类。
	 * <p>
	 * {@link BeanNameGenerator} 接口本身的默认值表示应该使用处理此 {@code @SpringBootApplication} 注解的扫描器的继承的 bean 名称生成器，
	 * 例如默认的 {@link AnnotationBeanNameGenerator} 或在引导时提供给应用程序上下文的任何自定义实例。
	 * @return 要使用的 {@link BeanNameGenerator}
	 * @see SpringApplication#setBeanNameGenerator(BeanNameGenerator)
	 * @since 2.3.0
	 */
	@AliasFor(annotation = ComponentScan.class, attribute = "nameGenerator")
	Class<? extends BeanNameGenerator> nameGenerator() default BeanNameGenerator.class;

	/**
	 * 指定是否应代理 {@link Bean @Bean} 方法以强制执行 bean 的生命周期行为，例如，即使在用户代码中直接调用 {@code @Bean} 方法，
	 * 也能返回共享的单例 bean 实例。此功能需要方法拦截，通过运行时生成的 CGLIB 子类来实现，其中包括一些限制，例如配置类及其方法不允许声明为 {@code final}。
	 * <p>
	 * 默认值为 {@code true}，允许在配置类内部进行 'inter-bean references'，同时允许从另一个配置类中调用此配置的 {@code @Bean} 方法。
	 * 如果每个特定配置的 {@code @Bean} 方法都是自包含的并且设计为容器使用的普通工厂方法，则可以将此标志切换为 {@code false}，
	 * 以避免 CGLIB 子类处理。
	 * <p>
	 * 关闭 bean 方法拦截将实际上像在非 {@code @Configuration} 类上声明的那样单独处理 {@code @Bean} 方法，也就是 "@Bean Lite Mode"
	 * （参见 {@link Bean @Bean 的文档}）。因此，在行为上等同于删除 {@code @Configuration} 注解。
	 * @since 2.2
	 * @return 是否代理 {@code @Bean} 方法
	 */
	@AliasFor(annotation = Configuration.class)
	boolean proxyBeanMethods() default true;

}
```

通过查看上述源码，我们可以看到 `@SpringBootApplication` 注解提供如下的成员属性【这里默认大家都是知道 **注解中的成员变量是以方法的形式存在的**】：

- `exclude` ：根据类（**Class**）排除指定的自动配置，该成员属性覆盖了 `@SpringBootApplication` 中组合的 `@EnableAutoConfiguration` 中定义的 **exclude** 成员属性。
- `excludeName` ：根据类名排除指定的自动配置，覆盖了 `@EnableAutoConfiguration` 中定义的 **excludeName** 成员属性。
- `scanBasePackages` ：指定扫描的基础 **package**，用于扫描带注解组件的基础包，例如包含 `@Component` 等注解的组件。
- `scanBasePackageClasses` ：指定扫描的类，用于相关组件的初始化。
- `nameGenerator` ：用于在 `Spring` 容器中为检测到的组件命名的 `BeanNameGenerator` 类。
- `proxyBeanMethods` ：指定是否代码 `@Bean` 方法以强制执行 **bean** 的生命周期行为。该功能需要通过运行时生成 `CGLIB` 子类来实现方法拦截。不过它包括一定的限制，例如配置类及其方法不允许声明为 `final` 等。`proxyBeanMethods` 的默认值为 `true`，允许配置类中进行 **inter-bean references**（bean 之间的引用）以及对该配置的 `@Bean` 方法的外部调用。如果 `@Bean` 方法都是自包含的，并且仅提供了容器使用的普通工程方法的功能，则可设置为 `false`，避免处理 `CGLIB` 子类。另外我们从源码中 `@since 2.2` 处也可以看出来，该属性是在 Spring Boot 2.2 版本新增的。


细心的读者，可能看过上面的源码会发现，`@SpringBootApplication` 注解的成员属性上大量使用了 `@AliasFor` 注解，**那该注解有什么作用呢？**

`@AliasFor` 注解用于桥接其他注解，它的 `annotation` 属性中指定了所桥接的注解类。如果我们点到 `annotation` 属性配置的注解中，可以看出 `@SpringBootApplication` 注解的成员属性其实已经在其他注解中定义过了。那之所以使用 `@AliasFor` 注解并重新在 `@SpringBootApplication` 中定义，主要就是为了减少用户在使用多注解上的麻烦。

> **知识拓展：**
> 简单总结一下 `@AliasFor` 的作用：
> - **定义别名关系**：通过在注解属性上使用 `@AliasFor` 注解，可以将一个属性与另一个属性建立别名关系。这意味着当使用注解时，你可以使用别名属性来设置目标属性的值。
> - **属性互通**：通过在两个属性上使用 `@AliasFor` 注解，并且将它们的 **attribute** 属性分别设置为对方，可以实现属性之间的双向关联。这意味着当设置其中一个属性的值时，另一个属性也会自动被赋予相同的值。
> - **注解继承**：当一个注解 **A** 使用 `@AliasFor` 注解指定了另一个注解 **B** 的属性为自己的别名属性时，如果类使用了注解 **A**，那么注解 **B** 的相关属性也会得到相应的设置。

在 `Spring Boot` 早期的版本中并没有 `@SpringBootConfiguration` 注解，它是后面的版本新加的，其内组合了 `@Configuration` 注解，如下图所示：

![](SpringBootConfiguration.png)

`@EnableAutoConfiguration` 注解组合了 `@AutoConfigurationPackage` 注解，如下图所示：

![](EnableAutoConfiguration.png)

除了一些元注解和基础注解，我们用一张类图来描述下 `@SpringBootApplication` 注解的组合结构：

![](SpringBootApplication.png)

从上图中，我们可以总结一下 `@SpringBootApplication` 注解的核心作用，如下：

- 组合 `@EnableAutoConfiguration`，用于开启 `Spring Boot` 的自动配置功能；
- 组合 `@ComponentScan`，用于激活 `@Component` 等注解类的初始化；
- 组合 `@SpringBootConfiguration`，用于标识一个类为配置类，以便在 `Spring` 应用程序上下文中进行配置。



# 总结

本篇通过查看 `@SpringBootApplication` 注解的源码，介绍其成员属性和组合注解的相关内容，这些内容将为我们后续的源码学习打下基础。下一篇博文，我们将开始介绍 `@EnableAutoConfiguration` 注解，开启 `Spring Boot` 的自动配置功能，敬请期待！！！
