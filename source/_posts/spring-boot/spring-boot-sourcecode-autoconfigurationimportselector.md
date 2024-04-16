---
title: 【Spring Boot 源码学习】走近 AutoConfigurationImportSelector
date: 2023-07-30 21:36:51
updated: 2023-09-12 10:30:16
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - AutoConfigurationImportSelector
  - ImportSelector
  - 自动装配逻辑
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)


# 引言

上篇博文我们了解了 [@EnableAutoConfiguration 注解](/2023/07/22/spring-boot/spring-boot-sourcecode-enableautoconfiguration/)，其中真正实现自动配置功能的核心实现者 `AutoConfigurationImportSelector` 还没有详细说明，本篇将从它的源码入手来重点介绍。

# 主要内容
在介绍  `AutoConfigurationImportSelector` 之前，有必要了解下它所实现的 `ImportSelector` 接口 ，如下所示：
## 1. ImportSelector 接口

在上篇博文中，我们介绍过 `@Import` 注解，它的许多功能其实是需要 `ImportSelector` 接口来实现，`ImportSelector` 接口决定可引入哪些 `@Configuration`。

下面我们来看一下 `ImportSelector` 接口的源码【**spring-context:5.3.25**】：

```java
/**
 * 实现了确定基于给定选择条件应该导入哪些 @Configuration 类的类型的接口，
 * 通常是一个或多个注解属性。
 * 
 * 一个 ImportSelector 可以实现以下任意一个 Aware 接口，并在调用 selectImports 
 * 方法之前调用其对应的方法：
 *    EnvironmentAware
 *    BeanFactoryAware
 *    BeanClassLoaderAware
 *    ResourceLoaderAware
 * 
 * 另外，该类也可以提供一个带有以下支持的参数类型的单个构造函数：
 *    Environment
 *    BeanFactory
 *    ClassLoader
 *    ResourceLoader
 * 
 * ImportSelector 实现通常与普通的 @Import 注解一样进行处理。
 * 然而，还可以推迟选择要导入的内容，直到所有 @Configuration 
 * 类都被处理完毕（详见 DeferredImportSelector）。
 *
 */
public interface ImportSelector {

    /**
     * 根据导入的 @Configuration 类的 AnnotationMetadata（注解元数据），
     * 选择并返回应该导入的类名称。
     * 
     * @return 返回类名的数组，如果没有则返回空数组。
     */
    String[] selectImports(AnnotationMetadata importingClassMetadata);

    /**
     * 返回一个用于从导入的候选类中排除类的断言函数，
     * 该函数会递归地应用于通过此选择器的导入项找到的所有类。
     * 
     * 如果对于给定的完全限定类名，该断言函数返回 true，
     * 则该类将不被视为被导入的配置类，从而跳过类文件加载和元数据检查。
     * 
     * @return 返回一个用于完全限定的候选类名的筛选断言函数，该函数适用于递归导入的配置类。
     * 如果没有筛选断言函数，则返回 null。
     * @since 5.2.4
     */
    @Nullable
    default Predicate<String> getExclusionFilter() {
        return null;
    }

}
```

通过阅读上述源码，我们可以看到 `ImportSelector` 接口提供了一个参数为 `AnnotationMetadata`【它里面包含了被 `@Import` 注解的类的注解信息，即注解元数据】 的方法 `selectImports` ，并返回了一个字符串数组【可以根据具体实现决定返回哪些配置类的全限定名】。

源码注释里也提到了，如果我们实现了 `ImportSelector` 接口的同时又实现了以下的 **4** 个 `Aware` 接口，那么 `Spring` 保证在调用 `ImportSelector` 之前会先调用 `Aware` 接口的方法。

这 4 个 Aware 接口分别是：

 -    `EnvironmentAware`
 -    `BeanFactoryAware`
 -    `BeanClassLoaderAware`
 -    `ResourceLoaderAware`

我们本篇要重点进行源码解析的 `AutoConfigurationImportSelector` 就实现了上述 **4** 个 `Aware` 接口，部分源码如下所示：

```java
// 其他导入语句省略
import org.springframework.beans.factory.BeanClassLoaderAware;
import org.springframework.beans.factory.BeanFactoryAware;
import org.springframework.context.EnvironmentAware;
import org.springframework.context.ResourceLoaderAware;
import org.springframework.context.annotation.DeferredImportSelector;
import org.springframework.core.Ordered;

/**
 * DeferredImportSelector 用于处理自动配置的延迟导入选择器。
 * 如果需要自定义的 @EnableAutoConfiguration 变体，也可以通过继承这个类来实现。
 */
public class AutoConfigurationImportSelector implements DeferredImportSelector, BeanClassLoaderAware,
        ResourceLoaderAware, BeanFactoryAware, EnvironmentAware, Ordered {
    // 其他省略
}
```

从上面的 类定义中，我们可以看到 `AutoConfigurationImportSelector` 并没有直接实现 `ImportSelector` 接口，而是实现了 `DeferredImportSelector` 接口【它是 `ImportSelector` 的子接口 】。

## 2. DeferredImportSelector 接口

**那 `AutoConfigurationImportSelector` 为啥不直接实现 `ImportSelector` 接口，而是实现了 `DeferredImportSelector` 接口呢？它们俩有什么区别呢？**

在讲解清楚之前，我们先来看看  `DeferredImportSelector` 接口的源码【**spring-context:5.3.25**】：

```java
/**
 * 一种在所有 @Configuration bean 处理完毕后运行的 ImportSelector 变体。
 * 这种类型的选择器在所选的导入项带有条件时特别有用。
 * 
 * 实现类可以扩展 org.springframework.core.Ordered 接口或使用 
 * org.springframework.core.annotation.Order 注解来指定与
 * 其他 DeferredImportSelectors 的优先级。
 * 
 * 实现类还可以提供一个导入组（import group），
 * 它可以在不同的选择器之间提供额外的排序和过滤逻辑。
 */
public interface DeferredImportSelector extends ImportSelector {

    /**
     * 返回一个特定的导入组。
     * 默认实现会在不需要分组的情况下返回 null。
     * 
     * @return 导入组的类，如果没有则返回 null。
     * @since 5.0
     */
    @Nullable
    default Class<? extends Group> getImportGroup() {
        return null;
    }


    /**
     * 用于将来自不同导入选择器的结果进行分组的接口。
     * 
     * @since 5.0
     */
    interface Group {

        /**
         * 使用指定的 DeferredImportSelector 处理导入的 @Configuration 类的 AnnotationMetadata。
         */
        void process(AnnotationMetadata metadata, DeferredImportSelector selector);

        /**
         * 返回此组应该导入的类的条目
         */
        Iterable<Entry> selectImports();


        /**
         * 一个条目，包含导入的配置类的 AnnotationMetadata 和要导入的类名。
         */
        class Entry {

            private final AnnotationMetadata metadata;

            private final String importClassName;

            public Entry(AnnotationMetadata metadata, String importClassName) {
                this.metadata = metadata;
                this.importClassName = importClassName;
            }

            /**
             * 返回导入的配置类的 AnnotationMetadata【注解元数据】
             */
            public AnnotationMetadata getMetadata() {
                return this.metadata;
            }

            /**
             * 返回要导入的类的完全限定名称。
             */
            public String getImportClassName() {
                return this.importClassName;
            }

            // 省略。。。
        }
    }

}

```

通过阅读上述源码，可以了解到之所以  `AutoConfigurationImportSelector`  没有直接实现 `ImportSelector` 接口，而是实现了 `DeferredImportSelector` 接口，是因为通过`DeferredImportSelector` 接口能够在处理自动配置时，拥有更高的灵活性和可定制性。

**总结来讲，它们的区别主要是如下几个方面：**

- **延迟导入**：`DeferredImportSelector` 具有延迟导入的能力，可以在所有的 `@Configuration` 类都被处理完毕之后再进行选择和导入。这样可以在整个配置加载过程完成后再根据某些条件或规则来决定要导入哪些类，从而实现更加动态和灵活的自动配置机制。

- **筛选导入**：`DeferredImportSelector` 提供了一个用于筛选候选类名的断言函数，可以根据一定的条件来排除某些类的导入。这样可以对自动配置的候选类进行进一步的过滤和控制，使得只有符合特定条件的类才会被真正导入。

- **自定义扩展**：通过实现 `DeferredImportSelector` 接口，开发人员可以更方便地扩展和定制自动配置逻辑。可以根据实际需求重写相应方法，实现自定义的自动配置规则和行为。

上述源码注释中，也说明了 `DeferredImportSelector` 的加载顺序可以通过 `@Order` 注解 或 实现 `Ordered` 接口来指定。它还可以提供一个导入组，实现在不同的选择器之间提供额外的排序和过滤逻辑，从而实现自定义 `Configuration` 的加载顺序。

## 3. AutoConfigurationImportSelector 功能概述

好了到这里，我们终于可以开始正式介绍 `AutoConfigurationImportSelector` 了。

下面我们通过如下的流程图，从整体上来了解 `AutoConfigurationImportSelector` 的核心功能及流程【其中省略了外部通过 `@Import` 注解调用该类的部分】：

 ![](autoconfigurationimportselector.png)

当 `AutoConfigurationImportSelector` 被 @Import 注解引入之后，它的 `selectImports` 方法会被调用并执行其实现的自动装配逻辑。

下面我们来看看 `selectImports` 方法的源码，如下所示：

```java
    @Override
    public String[] selectImports(AnnotationMetadata annotationMetadata) {
        // 检查自动配置功能是否开启，默认为开启
        if (!isEnabled(annotationMetadata)) {
            return NO_IMPORTS;
        }
        // 封装将被引入的自动配置信息
        AutoConfigurationEntry autoConfigurationEntry = getAutoConfigurationEntry(annotationMetadata);
        // 返回符合条件的配置类的全限定名数组
        return StringUtils.toStringArray(autoConfigurationEntry.getConfigurations());
    }
    
    /**
     * 根据导入@Configuration类的AnnotationMetadata返回AutoConfigurationImportSelector.AutoConfigurationEntry。
     * @param 配置类的注解元数据。
     * @return 应该导入的自动配置。
     */
    protected AutoConfigurationEntry getAutoConfigurationEntry(AnnotationMetadata annotationMetadata) {
        if (!isEnabled(annotationMetadata)) {
            return EMPTY_ENTRY;
        }
        // 从AnnotationMetadata返回适当的AnnotationAttributes。默认情况下，此方法将返回getAnnotationClass()的属性。
        AnnotationAttributes attributes = getAttributes(annotationMetadata);
        // 通过 SpringFactoriesLoader 类提供的方法加载类路径中META-INF目录下的
        // spring.factories文件中针对 EnableAutoConfiguration 的注解配置类
        List<String> configurations = getCandidateConfigurations(annotationMetadata, attributes);
        // 对获得的注解配置类集合进行去重处理，防止多个项目引入同样的配置类
        configurations = removeDuplicates(configurations);
        // 获得注解中被 exclude 或 excludeName 所排除的类的集合
        Set<String> exclusions = getExclusions(annotationMetadata, attributes);
        // 检查被排除类是否可实例化，是否被自动注册配置所使用，不符合条件则抛出异常
        checkExcludedClasses(configurations, exclusions);
        // 从自动配置类集合中去除被排除的类
        configurations.removeAll(exclusions);
        // 检查配置类的注解是否符合 spring.factories 文件中 AutoConfigurationImportFilter 指定的注解检查条件
        configurations = getConfigurationClassFilter().filter(configurations);
        // 将筛选完成的配置类和排除的配置类构建为事件类，并传入监听器。监听器的配置在于 spring.factories 文件中，通过 AutoConfigurationImportListener 指定
        fireAutoConfigurationImportEvents(configurations, exclusions);
        // 创建并返回一个条目，其中包含了筛选完成的配置类和排除的配置
        return new AutoConfigurationEntry(configurations, exclusions);
    }
```



# 总结

通过阅读上述源码，对照相关的流程图，我们从整体上了解了 `AutoConfigurationImportSelector` 自动装配逻辑的核心功能及流程，由于篇幅有限，更加细化的功能及流程解析，笔者将在后续的博文中，带大家一起通过源码来一步步完成，敬请期待！！！

