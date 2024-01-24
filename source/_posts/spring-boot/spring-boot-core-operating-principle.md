---
title: Spring Boot 核心运行原理介绍
date: 2023-07-13 08:00:00
updated: 2024-01-15 22:30:50 
categories:
- 开发框架-Spring Boot
tags:
- Spring Boot
- 核心运行原理
- 自动配置
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 引言

还记得，笔者在前面的博文[《Spring Boot 项目介绍》](/2023/02/19/spring-boot-project-introduction/)中提到了，Spring Boot 最核心的功能就是自动配置，该功能的实现是基于 “约定由于配置” 的原则。

那很多读者就要问了，`Spring Boot` 它是如何来约定的呢？又是如何实现自动配置功能的呢？

从本篇开始，笔者将带领大家通过学习 `Spring Boot` 源码，来了解它核心的运行原理。后续的内容将会涉及自动配置的运作原理、核心功能模块、核心注解等等。

# 整体介绍

## 1. 核心运行原理图

在我们的项目中，接入 `Spring Boot` 其实是最简单的。我们只需要引入 `Spring Boot` 对应的 `Starters`，然后它启动时就会自动加载相关依赖，并配置相应的初始化参数，从而可以很方便地对第三方软件进行集成。

我们先从整体上来看一下 `Spring Boot` 实现上述自动配置机制的核心运行原理图：

![](core-operating-principle.png)


上图描述了 `Spring Boot` 自动配置功能运作过程中涉及的几个核心功能及其相互之间的关系，其中的内容将在第3小节介绍。


## 2. 自动配置的整体流程

从上面的 `Spring Boot` 自动配置功能核心运行原理图，我们可以了解它自动配置的整体流程，如下：


- 首先 `Spring Boot` 通过 `@EnableAutoConfiguration` 注解开启自动配置；
- 然后，加载 `spring.factories` 中注册的各种 `AutoConfiguration` 类；
- 接着，当某个 `AutoConfiguration` 类满足其注解 `@Conditional` 指定的生效条件（`Starters` 提供的依赖、配置或 `Spring` 容器中存在某个 `Bean` 等）时，实例化该 `AutoConfiguration` 类中定义的 `Bean`（组件等）；
- 最后，将这些组件都注入 `Spring` 容器中，完成依赖框架的自动配置功能。

## 3. 各核心功能和组件初步介绍

下面我们对上述涉及到的组件进行简单介绍：

- `@EnableAutoConfiguration` : 了解过的朋友，肯定疑惑，Spring Boot 的启动类明明没有引入该注解呀，其实该注解最终是由组合注解 @SpringBootApplication 引入，用来完成自动配置开启，扫描各个 jar 包下的 spring.factories 文件，并加载文件中注册的 AutoConfiguration 类等。

- `spring.factories` : 配置文件，位于 jar 包的 META-INF 目录下，按照指定格式注册了自动配置的 AutoConfiguration 类。spring.factories 也可以包含其他类型待注册的类。该配置文件不仅存在于 Spring Boot 项目中，也可以存在于自定义的自动配置（或 Starter）项目中。

- `AutoConfiguration 类` ：自动配置类，代表了 Spring Boot 中一类以 XXXAutoCofiguration 命名的自动配置类。其中定义了三方组件集成 Spring 所需初始化的 Bean 和条件。

- `@Conditional` ：条件注解及其衍生注解，在 AutoConfiguration 类上使用，当满足该条件注解时才会实例化 AutoConfiguration 类。

- `Starters` ：三方组件的依赖以及配置，Spring Boot 已经预置的组件。Spring Boot 默认的 Starters 项目往往只包含了一个 pom 依赖的项目。如果是自定义的 starter，该项目还需要包含 spring.factories 文件、AutoConfiguration 类和其他配置类。


# 总结

本篇我们从概念层面介绍了 Spring Boot 自动配置的核心运行原理和整理流程，后续的博文将围绕这些核心部分，从源码层面进行详细介绍，敬请期待！！！

# 参考
《**Spring Boot** 技术内幕-架构设计与实现原理》朱智胜 


