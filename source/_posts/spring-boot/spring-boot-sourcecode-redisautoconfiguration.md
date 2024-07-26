---
title: 【Spring Boot 源码学习】RedisAutoConfiguration 详解
date: 2023-10-29 23:11:28
updated: 2024-01-29 00:18:33
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - RedisAutoConfiguration
  - RedisTemplate
---



![](/images/spring-boot-logo.png)

# 引言
上篇博文，**Huazie** 带大家从源码角度分析了 **Spring Boot** 内置的 `http` 编码功能，进一步熟悉了自动配置的装配流程。本篇趁热打铁，继续带大家分析 **Spring Boot** 内置的有关 **Redis** 的自动配置类【`RedisAutoConfiguration`】。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 主要内容

## 1. Spring Data Redis

**Spring Data Redis** 是 **Spring Data** 家族的一部分，它提供了从 **Spring** 应用程序中轻松配置和访问 **Redis** 的功能。

我们来看看官方介绍的特性：

- 连接包作为多个 **Redis** 驱动程序（ **Lettuce** 和 **Jedis** ）的低级别抽象。
- 将 **Redis** 驱动程序异常转换为 **Spring** 的可移植数据访问异常层次结构。
- 提供各种 **Redis** 操作、异常转换和序列化支持的 **RedisTemplate**。
- 支持发布订阅（例如用于消息驱动 **POJO** 的消息监听器容器）。
- 支持 **Redis Sentinel** 和 **Redis Cluster**。
- 使用 **Lettuce** 驱动程序的响应式 **API**。
- 支持 **JDK**、**String**、**JSON**和 **Spring** 对象 / **XML** 映射序列化器。
- 在 **Redis** 上实现 **JDK** 集合。
- 支持原子计数器类。
- 支持排序和管道功能。
- 专用于 **SORT**、**SORT**/**GET**模式和支持返回批量值的功能。
- 为 **Spring** 缓存抽象提供 **Redis** 实现。
- 自动实现 **Repository** 接口，包括使用 `@EnableRedisRepositories` 支持自定义查询方法。
- 对存储库提供 **CDI** 支持。

在 **Spring Data Redis** 中，我们可以直接使用 `RedisTemplate` 及其相关的类来操作 **Redis**。虽然 `RedisConnection` 提供了接受和返回二进制值（字节数组）的低级方法，但 `RedisTemplate` 负责序列化和连接管理，使用户可以无需处理这些细节。

`RedisTemplate` 还提供了操作视图（按照 **Redis** 命令参考进行分组），这些视图提供了丰富、通用的接口，用于针对特定类型或特定键进行操作（通过 `KeyBound` 接口实现），如下表所示：

|接口| 描述  |
|:--|:--|
| GeoOperations | Redis地理空间操作，例如GEOADD、GEORADIUS等。 |
| HashOperations | Redis哈希操作 |
| HyperLogLogOperations | Redis键绑定哈希操作 |
| ListOperations | Redis列表操作 |
| SetOperations | Redis集合操作 |
| ValueOperations |Redis字符串（或值）操作  |
|ZSetOperations  |Redis有序集合操作  |
| BoundGeoOperations | Redis键绑定地理空间操作 |
| BoundHashOperations | Redis键绑定哈希操作 |
| BoundKeyOperations | Redis键绑定操作 |
| BoundListOperations | Redis键绑定列表操作 |
|BoundSetOperations |Redis键绑定集合操作 |
|BoundValueOperations |Redis键绑定字符串（或值）操作 |
|BoundZSetOperations |Redis键绑定有序集合操作 |


下面我们来看看相关的 **Spring** 配置：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:p="http://www.springframework.org/schema/p"
    xsi:schemaLocation="http://www.springframework.org/schema/beans https://www.springframework.org/schema/beans/spring-beans.xsd">

    <bean id="jedisConnFactory" class="org.springframework.data.redis.connection.jedis.JedisConnectionFactory" p:use-pool="true"/>
    <!-- redis 模板定义 -->
    <bean id="redisTemplate" class="org.springframework.data.redis.core.RedisTemplate" p:connection-factory-ref="jedisConnFactory"/>

</beans>
```

一旦配置完成，Redis 模板就是线程安全的，并且可以在多个实例之间重用。

`RedisTemplate` 使用基于 **Java** 的序列化器进行大部分操作。也就意味着通过模板写入或读取的任何对象都是通过 **Java** 进行序列化和反序列化的。

我们也可以更改模板上的序列化机制，可以添加如下配置：

```xml
<bean id="redisTemplate" class="org.springframework.data.redis.core.RedisTemplate">
    <property name="connectionFactory" ref="redisConnectionFactory"/>
    <property name="keySerializer">
        <bean class="org.springframework.data.redis.serializer.StringRedisSerializer"/>
    </property>
    <property name="valueSerializer">
        <bean class="org.springframework.data.redis.serializer.JdkSerializationRedisSerializer"/>
    </property>
    <property name="hashKeySerializer">
        <bean class="org.springframework.data.redis.serializer.StringRedisSerializer"/>
    </property>
    <property name="hashValueSerializer">
        <bean class="org.springframework.data.redis.serializer.JdkSerializationRedisSerializer"/>
    </property>
</bean>
```

而 Redis 模块提供了几个序列化器的实现，有关这些实现大家可以查看 `org.springframework.data.redis.serializer` 包。

![](redis-serializer.png)

还可以将任何序列化程序设置为 `null`，并通过设置 `enableDefaultSerializer` 属性为 `false` 来使用`RedisTemplate` 与原始字节数组一起使用。

> **注意**： 模板要求所有键都不为空。但是，只要底层序列化程序接受值，值就可以为空。



下面我们可以注入 `RedisTemplate`，并调用 `RedisTemplate` 的方法进行存储、查询、删除等操作。

```java
@Autowired
private RedisTemplate<String, Object> redisTemplate;


// 存储数据
redisTemplate.opsForValue().set("key", "value");
// 查询数据
Object value = redisTemplate.opsForValue().get("key");
// 删除数据
redisTemplate.delete("key");
```


对于需要特定模板视图的情况，声明视图作为依赖项并注入模板。容器会自动执行转换，消除`opsFor[X]` 调用，如下所示的示例：

```java
public class Example {  
    // inject the template as ListOperations
    @Resource(name="redisTemplate")
    private ListOperations<String, String> listOps;
    
    public void addLink(String userId, URL url) {
        listOps.leftPush(userId, url.toExternalForm());
    }
}
```

当然 **Spring Data Redis** 肯定不止上述这些，有需要深入了解的读者们，请看如下：

> **参考：** [Spring Data Redis 官方文档](https://spring.io/projects/spring-data-redis)

## 2. RedisAutoConfiguration


**那么 Spring Data Redis 的 `RedisTemplate` 的自动配置在 **Spring Boot** 是如何实现的呢？**

**Spring Boot**  是通过内置的 `RedisAutoConfiguration` 配置类来完成这一功能。下面我们具体分析一下：

> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 `2.7.9`，其他版本有所出入，可自行查看源码。

### 2.1 加载自动配置组件
从之前的[《【Spring Boot 源码学习】自动装配流程源码解析（上）》](../../../../../2023/08/06/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-1/)中，我们知道 **Spring Boot** 内部针对自动配置类，会读取如下两个配置文件：

- `META-INF/spring.factories`
- `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports`

![](autoconfiguration.png)

实际上 在 `Spring Boot 2.7.9` 版本中， **Spring Boot** 自己内部的 `META-INF/spring.factories` 中有关自动配置的注册类的配置信息已经被去除掉了，不过其他外围的 **jar** 中可能有自己的 `META-INF/spring.factories` 文件，它里面也有关于自动配置注册类的配置信息；

而 Spring Boot 内置的 `RedisAutoConfiguration` 配置类，则是配置在上述的第二个配置文件 `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports` 中。

![](redisautoconfiguration.png)


### 2.2 过滤自动配置组件

上述自动配置加载完之后，就来到了 [《【Spring Boot 源码学习】自动装配流程源码解析（下）》](../../../../../2023/08/21/spring-boot/spring-boot-sourcecode-autoconfigurationdetail-2/) 介绍的 **过滤自动配置组件** 逻辑。

这部分数据对应的配置内容在 `META-INF/spring-autoconfigure-metadata.properties` 文件中：

```java
org.springframework.boot.autoconfigure.data.redis.RedisAutoConfiguration=
org.springframework.boot.autoconfigure.data.redis.RedisAutoConfiguration.ConditionalOnClass=org.springframework.data.redis.core.RedisOperations
```

显然这里涉及到了 `ConditionalOnClass` 注解，我们翻看 `RedisAutoConfiguration` 配置类的源码，如下：

```java
@AutoConfiguration
@ConditionalOnClass(RedisOperations.class)
@EnableConfigurationProperties(RedisProperties.class)
@Import({ LettuceConnectionConfiguration.class, JedisConnectionConfiguration.class })
public class RedisAutoConfiguration {

    @Bean
    @ConditionalOnMissingBean(name = "redisTemplate")
    @ConditionalOnSingleCandidate(RedisConnectionFactory.class)
    public RedisTemplate<Object, Object> redisTemplate(RedisConnectionFactory redisConnectionFactory) {
        // 。。。
    }

    @Bean
    @ConditionalOnMissingBean
    @ConditionalOnSingleCandidate(RedisConnectionFactory.class)
    public StringRedisTemplate stringRedisTemplate(RedisConnectionFactory redisConnectionFactory) {
        // 。。。
    }

}
```

#### 2.2.1 涉及注解 

我们先来看看上述 `RedisAutoConfiguration` 配置类涉及到的注解，如下：
- `@AutoConfiguration` : 该类是一个自动配置类，**Spring Boot** 会根据项目中的依赖自动配置这个类的实例。
- `@ConditionalOnClass(RedisOperations.class)` ：只有在项目中引入了 `RedisOperations` 类（通常由 **spring-data-redis** 库提供）的情况下，才会加载这个配置类。
- `@EnableConfigurationProperties(RedisProperties.class)` ：启用`RedisProperties` 类作为配置属性。这样，我们就可以在 **application.properties** 或**application.yml** 文件中定义 **Redis** 的相关配置。
- `@Import({ LettuceConnectionConfiguration.class, JedisConnectionConfiguration.class })` ：导入注解，表示导入 `LettuceConnectionConfiguration` 和 `JedisConnectionConfiguration` 这两个类。这两个类通常用于配置 **Redis** 连接的具体实现，例如使用 **Lettuce** 还是 **Jedis** 等。
- `@Bean` ：用于声明一个方法创建的对象是一个 **Spring** 管理的 **Bean**。**Spring** 容器会自动管理这个 **Bean** 的生命周期，包括依赖注入、初始化和销毁等。
- `@ConditionalOnMissingBean` ：只有在当前 **Spring** 容器中不存在指定类型的 **Bean** 时，才会执行被注解的方法。这样可以用于确保在需要的时候才创建某个 **Bean**，避免重复创建。
- `@ConditionalOnSingleCandidate`：只有在当前上下文中存在且只有一个指定类型的 **bean** 候选者时，才会创建这个 **bean**。

#### 2.2.2 RedisProperties

其中 `RedisProperties` 类的属性值对应着 `application.yml` 或 `application.properties` 中的配置，通过注解`@ConfigurationProperties(prefix = "spring.redis")` 实现的属性注入。

有关属性注入的内容后续笔者会另外介绍，我们先来看看`RedisProperties` 类相关的部分源码 和 对应的配置参数：

```java
@ConfigurationProperties(prefix = "spring.redis")
public class RedisProperties {

    // 。。。

    // Redis 服务器主机地址.
    private String host = "localhost";
    
    // 。。。

    // Redis 服务器的端口
    private int port = 6379;

    private Sentinel sentinel;

    private Cluster cluster;

    private final Jedis jedis = new Jedis();

    private final Lettuce lettuce = new Lettuce();

    // Redis 连接池配置
    public static class Pool {
        // 。。。
    }
    // Redis 集群配置
    public static class Cluster {
        // 。。。
    }
    // Redis 哨兵配置
    public static class Sentinel {
        // 。。。
    }
    // Jedis 客户端配置
    public static class Jedis {

        // Jedis 连接池配置
        private final Pool pool = new Pool();
    }
    // Lettuce 客户端配置
    public static class Lettuce {
        // Lettuce 连接池配置
        private final Pool pool = new Pool();

        private final Cluster cluster = new Cluster();
    }
}
```

然后在 `application.properties` 中，我们就可以添加类似如下的配置：

```bash
# Redis 单机配置
spring.redis.host=127.0.0.1
spring.redis.port=31113

# Redis 集群配置
# nodes属性是Redis集群节点的地址和端口，用逗号分隔。
spring.redis.cluster.nodes=192.168.1.1:7000,192.168.1.2:7001,192.168.1.3:7002
# max-redirects属性是最大重定向次数，用于处理节点故障的情况。
spring.redis.cluster.max-redirects=3

# mymaster是哨兵模式下的主节点名称。
spring.redis.sentinel.master=mymaster
# nodes是哨兵模式下的从节点地址和端口。
spring.redis.sentinel.nodes=192.168.1.1:26379,192.168.1.2:26379,192.168.1.3:26379

# ...其他配置省略
```

### 2.3 redisTemplate 方法
先来看看 `redisTemplate` 方法的源码【Spring Boot 2.7.9】：

```java
@Bean
@ConditionalOnMissingBean(name = "redisTemplate")
@ConditionalOnSingleCandidate(RedisConnectionFactory.class)
public RedisTemplate<Object, Object> redisTemplate(RedisConnectionFactory redisConnectionFactory) {
    RedisTemplate<Object, Object> template = new RedisTemplate<>();
    template.setConnectionFactory(redisConnectionFactory);
    return template;
}
```

上述逻辑表示只有在当前上下文中不存在名为 `"redisTemplate"` 的 **Bean** 时，才会创建一个名为 `redisTemplate` 的 **RedisTemplate Bean**，并将其与一个可用的 **Redis** 连接工厂关联起来。


### 2.4 stringRedisTemplate 方法

我们再来看看 `stringRedisTemplate` 方法的源码【Spring Boot 2.7.9】：

```java
@Bean
@ConditionalOnMissingBean
@ConditionalOnSingleCandidate(RedisConnectionFactory.class)
public StringRedisTemplate stringRedisTemplate(RedisConnectionFactory redisConnectionFactory) {
    return new StringRedisTemplate(redisConnectionFactory);
}
```

上述逻辑也好理解，它表示只有在当前上下文中不存在名为 `"stringRedisTemplate"` 的 **Bean** 时，才会创建一个名为`stringRedisTemplate`的 **StringRedisTemplate Bean**，并将其与一个可用的 **Redis** 连接工厂关联起来。

`StringRedisTemplate` 是 `RedisTemplate` 的子类，专门用于处理字符串类型的数据。

`StringRedisTemplate` 使用的是 `StringRedisSerializer`，它在存入数据时会将数据先序列化成字节数组。

默认情况下，`StringRedisTemplate` 采用的序列化策略有两种：
- `String` 的序列化策略，
- `JDK` 的序列化策略。


# 总结
本篇我们深入分析了 `RedisAutoConfiguration ` 配置类的相关内容，进一步加深了对自动配置装配流程的了解。其中有关 `LettuceConnectionConfiguration` 和 `JedisConnectionConfiguration` 这两个用于配置 **Redis** 连接的具体实现，笔者后面有时间再带大家详细分析一下。
