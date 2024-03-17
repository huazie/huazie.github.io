---
title: Spring Boot 项目介绍
date: 2023-02-19 22:55:05
updated: 2024-01-15 22:32:10 
categories:
- 开发框架-Spring Boot
tags:
- Spring Boot
- 约定由于配置
---


[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)


# 引言

作为学习过 **Java** 的软件开发者，相信都知道 **Spring** 这一伟大的框架，它所拥有的强大功能之一就是可以集成各种开源软件。但随着互联网的高速发展，各种框架层出不穷，这就对系统架构的灵活性、扩展性、可伸缩性 和 高可用性都提出了新的要求。

随着项目的发展，**Spring** 慢慢地集成了更多的开源软件，引入大量配置文件，这会导致程序出错率高、运行效率低下的问题。为了解决这些状况，**Spring Boot** 应运而生。

**Spring Boot** 本身并不提供 **Spring** 的核心功能，而是作为 **Spring** 的脚手架框架，以达到快速构建项目、预置三方配置、开箱即用的目的。

# 项目介绍
## 1. 设计理念

**约定优于配置（Convention Over Configuration）**，又称**按约定编程**，是一种软件设计范式，旨在减少软件开发人员需要做决定的数量，执行起来简单而又不失灵活。**Spring Boot** 的核心设计就完美遵从了此范式。

**Spring Boot** 的功能从细节到整体都是基于 ”约定优于配置“ 来开发的，从基础框架的搭建、配置文件、中间件的集成、内置容器以及其生态中的各种 **Starters**，无不遵从此设计范式。

**Starter** 作为 **Spring Boot** 的核心功能之一，基于自动配置代码提供了自动配置模块及依赖，让软件集成变得简单、易用。我们也可以创建自己的 **Starter**，来使我们的应用接入 Spring Boot。

## 2. 设计目标

说到 **Spring Boot** 的设计目标，就不得不提到它的研发团队--**Pivotal** 公司。该公司的企业目标是 **"致力于改变世界构造软件的方式（We are transforming how the world builds software）"** 。

**Spring Boot** 框架的设计理念完美遵从了它所属企业的目标。**Spring Boot** 不是为已解决的问题提供新的解决方案，而是为平台和开发者带来一种全新的体验：

**整合成熟技术框架、屏蔽系统复杂性、简化已有技术的使用，从而降低软件的使用门槛，提升软件开发和运维的效率。**

## 3. 源代码的目录结构

[Spring Boot 源代码地址](https://github.com/spring-projects/spring-boot)

> 知识点： **Spring Boot 2.3.x** 系列版本开始用 Gradle 构建 ，**2.2.x 及之前** 的系列版本都用 Maven 构建。

**Spring Boot** 项目的目录结构分为两部分，一部分是整个开源项目的目录结构，另一部分是细化到 **jar** 包级别的目录结构。

下面我们从整体到局部一起了解下 **Spring Boot** 项目的目录结构。

**首先从整体出发**，如下图所示是 **Spring Boot** 在 **GitHub** 上 **3.0.2** 版本源代码顶层目录结构：

![](source-code-directory.png)


不同版本之间的 **Spring Boot** 源代码的顶层目录结构会有所变化，但并不影响其核心功能。

| 目录名 | 描述  |
|:--|:--|
| **spring-boot-project** |  **Spring Boot** 核心项目代码，包含核心、工具、安全、文档、starters等项目。|
| **spring-boot-system-tests** | **Spring Boot** 部署和镜像测试。 |
|**spring-boot-tests**| **Spring Boot** 集成和冒烟的测试。|


关于顶层目录结构，简单了解即可，从 **1.5.x** 到 **3.0.x** 版本，该层级的目录结构在不停地发生变化。

**然后从局部出发**，如下图所示是 **spring-boot-project** 在 **GitHub** 上 **3.0.2** 版本源代码的目录结构：

![](source-code-directory1.png)

**spring-boot-project** 目录是在 **Spring Boot 2.0** 版本发布后新增的目录层级，并将原来在 **Spring Boot 1.5.x** 版本中的一级模块作为 **spring-boot-project** 的子模块。该模块包含了 **Spring Boot** 所有的核心功能。

| 目录名 | 描述  |
|:--|:--|
| **spring-boot** |  **Spring Boot** 核心代码，也是入口类 **SpringApplication** 类所在项目。|
| **spring-boot-actuator** |  提供应用程序的监控、统计、管理及自定义等相关功能。|
| **spring-boot-actuator-autoconfigure** |  针对 **actuator** 提供的自动配置功能。|
| **spring-boot-autoconfigure** |  **Spring Boot** 自动配置核心功能，默认集成了多种常用框架的自动配置类等。|
| **spring-boot-cli** |  命令工具，提供快速搭建项目原型、启动服务、执行 **Groovy** 脚本等功能。**3.0.x** 版本开始被移除。|
| **spring-boot-dependencies** |  依赖和插件的版本信息。|
| **spring-boot-devtools** |  开发者工具，提供热部署、实时加载、禁用缓存等提升开发效率的功能。|
| **spring-boot-docs** |  参考文档相关内容。|
| **spring-boot-parent** |  **spring-boot-dependencies** 的子模块，是其他项目的父模块。|
| **spring-boot-properties-migrator** |  **Spring Boot 2.0** 版本新增的模块，支持升级版本配置属性的迁移。**3.0.x** 版本开始被移除。|
| **spring-boot-starters** |  **Spring Boot** 以预定义的方式集成了其他应用的 **starter** 集合。|
| **spring-boot-test** |  测试功能相关代码。|
| **spring-boot-test-autoconfigure** |  测试功能自动配置相关代码。|
| **spring-boot-tools** |  **Spring Boot** 工具支持模块，包含 **Ant**、**Maven**、**Gradle** 等构建工具。|

## 4. 整体架构
上面给大家介绍了 **Spring Boot** 的核心项目结构及功能，下面我们从架构层面了解一下 **Spring Boot**  的不同模块之间的依赖关系。

为了更清晰地表达 **Spring Boot** 各项目之间的关系，在下图中我们基于依赖的传递性，省略了部分依赖关系。
> 比如，**Spring Boot Starters** 不仅依赖了 **Spring Boot Autoconfigure** 项目，还依赖了 **Spring Boot** 和 **Spring**，而 **Spring Boot Autoconfigure** 项目又依赖了 **Spring Boot** ，**Spring Boot** 又依赖了 **Spring** 相关项目。因而在下图中就省略了 **Spring Boot Starters** 和 底层依赖的关联。<br/>
> 同样，**Spring Boot Parent** 是 **Spring Boot** 及图中依赖 **Spring Boot** 项目的 **Parent** 项目，为了结构清晰，图中不显示相关关联。
> 
![](spring-boot-dependency.png)

从上图中我们可以看出 **Spring Boot** 几乎完全基于 **Spring**，同时提供了 **Spring Boot** 和 **Spring Boot Autoconfigure** 两个核心的模块，而其他相关功能又都是基于这两个核心模块展开的。

# 总结

到目前为止，**Spring Boot** 的最新发布版为 **3.0.2**，该项目也在不断的更新和迭代中，上述介绍可能会有遗漏或差异，请至 官方 **GitHub** 查阅详情。

# 参考
《**Spring Boot** 技术内幕-架构设计与实现原理》朱智胜 

