---
title: 【Spring Boot 源码学习】 JedisConnectionConfiguration 详解
date: 2023-11-05 23:35:22
updated: 2024-01-27 23:27:14
categories:
  - 开发框架-Spring Boot
tags:
  - Spring Boot
  - JedisConnectionConfiguration
  - RedisConnectionConfiguration
  - 单机模式
  - 哨兵模式
  - 集群模式
---

[《Spring Boot 源码学习系列》](/categories/开发框架-Spring-Boot/)

![](/images/spring-boot-logo.png)

# 一、引言
上篇博文，**Huazie** 带大家从源码角度分析了 **Spring Boot** 内置的有关 **Redis** 的自动配置类【`RedisAutoConfiguration`】，其中有关 `LettuceConnectionConfiguration` 和 `JedisConnectionConfiguration` 这两个用于配置 **Redis** 连接的具体实现还未介绍。本篇就以我们常用的 **Jedis** 实现 为例，带大家详细分析一下 `JedisConnectionConfiguration` 配置类。

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)
# 二、往期内容
在开始本篇的内容介绍之前，我们先来看看往期的系列文章【有需要的朋友，欢迎关注系列专栏】：

<table>
  <tr>
    <td rowspan="15" align="left" > 
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
  <tr>
    <td align="left" > 
      <a href="/2023/10/15/spring-boot/spring-boot-sourcecode-conditional/">【Spring Boot 源码学习】@Conditional 条件注解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/10/22/spring-boot/spring-boot-sourcecode-httpencodingautoconfiguration/">【Spring Boot 源码学习】HttpEncodingAutoConfiguration 详解</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2023/10/29/spring-boot/spring-boot-sourcecode-redisautoconfiguration/">【Spring Boot 源码学习】RedisAutoConfiguration 详解</a> 
    </td>
  </tr>
</table>

# 三、主要内容

## 1. RedisConnectionFactory

`RedisConnectionFactory` 是 **Spring Data Redis** 中的一个接口，它提供了创建和管理 **Redis** 连接的方法。使用 `RedisConnectionFactory` 可以获取到 **Redis** 连接对象，然后通过该对象对 **Redis** 进行存储、查询、删除等操作。

我们来看看 `RedisConnectionFactory` 的相关的源码：

```java
// 线程安全的 Redis 连接工厂
public interface RedisConnectionFactory extends PersistenceExceptionTranslator {
  RedisConnection getConnection();

  RedisClusterConnection getClusterConnection();

  boolean getConvertPipelineAndTxResults();

  RedisSentinelConnection getSentinelConnection();
}
```

我们简单分析一下 **Redis** 连接工厂中的方法：

- `RedisConnection getConnection()` ：提供与 **Redis** 交互的合适连接。如果连接工厂需要初始化，但工厂尚未初始化，则抛出 `IllegalStateException`。
- `RedisClusterConnection getClusterConnection()` ：提供与 **Redis Cluster** 交互的合适连接。如果连接工厂需要初始化，但工厂尚未初始化，则抛出 `IllegalStateException`。
- `boolean getConvertPipelineAndTxResults()` ：指定是否应将管道结果转换为预期的数据类型。如果为 `false`，`RedisConnection.closePipeline()` 和`RedisConnection#exec()` 的结果将是底层驱动程序返回的类型。
- `RedisSentinelConnection getSentinelConnection()` ：提供与 **Redis Sentinel** 交互的合适连接。如果连接工厂需要初始化，但工厂尚未初始化，则抛出 `IllegalStateException`。

以常用的 **Jedis** 实现为例，我们介绍一下 **Redis** 连接工厂的 **Jedis** 实现，即 `JedisConnectionFactory`。由于该类这是 **Spring Data Redis** 中的代码，本篇不详细展开了，感兴趣的朋友可以自行翻阅 **Spring** 源码进行查看。

**那 `JedisConnectionFactory` 主要有哪些内容呢 ？**

- **创建 Jedis 连接** ：通过调用 `createXXX()` 方法，可以创建一个 **Jedis** 连接对象，用于与 **Redis** 服务器进行通信。当然，在获取连接之前，我们必须先初始化该连接工厂。

- **管理连接池** ：它内部维护了一个连接池，用于管理和复用 **Jedis** 连接。当需要创建一个新的 **Jedis** 连接时，首先会检查连接池中是否有可用的连接，如果有则直接使用，否则创建一个新的连接。这样可以提高性能，减少频繁创建和关闭连接带来的开销。

  ```java
  protected Jedis fetchJedisConnector() {
    try {

      if (getUsePool() && pool != null) {
        return pool.getResource();
      }

      Jedis jedis = createJedis();
      // force initialization (see Jedis issue #82)
      jedis.connect();

      return jedis;
    } catch (Exception ex) {
      throw new RedisConnectionFailureException("Cannot get Jedis connection", ex);
    }
  }
  ```
- **配置连接参数** ：允许用户自定义连接参数，例如 **超时时间**、**最大连接数**等。这些参数可以在创建连接时通过构造函数传入，也可以在创建连接后，通过 `JedisPoolConfig` 或者下面的三种连接类型的配置类进行修改。

- **支持多种连接类型** ：包括 **单机连接**、**哨兵连接** 和 **集群连接**。这些连接类型的配置如下：
  - `RedisStandaloneConfiguration`（单机配置）
  - `RedisSentinelConfiguration`（哨兵配置）
  - `RedisClusterConfiguration`（集群配置）
  
### 1.1 单机连接
单机连接，我们需要使用到 `RedisStandaloneConfiguration` ，可见如下示例：

```java
@Configuration
public class RedisConfig {

    @Bean
    public JedisConnectionFactory jedisConnectionFactory() {
        RedisStandaloneConfiguration configuration = new RedisStandaloneConfiguration("localhost", 6379);
        JedisConnectionFactory jedisConnectionFactory = new JedisConnectionFactory(configuration);
        return jedisConnectionFactory;
    }
}
```

### 1.2 集群连接
集群连接，我们需要使用到 `RedisClusterConfiguration` ，示例如下：

```java
@Configuration
public class RedisConfig {

    @Bean
    public JedisConnectionFactory jedisConnectionFactory() {
        Set<Host> nodes = new HashSet<>();
        nodes.add(new Host("127.0.0.1", 20011));
        nodes.add(new Host("127.0.0.1", 20012));
        nodes.add(new Host("127.0.0.1", 20013));
        RedisClusterConfiguration clusterConfiguration = new RedisClusterConfiguration();
        clusterConfiguration.setClusterNodes(nodes);
        JedisConnectionFactory jedisConnectionFactory = new JedisConnectionFactory(clusterConfiguration);
        return jedisConnectionFactory;
    }
}

```

### 1.3 哨兵连接
哨兵连接，我们需要使用到 `RedisSentinelConfiguration` ，参考如下：

```java
@Configuration
public class RedisConfig {

    @Bean
    public JedisConnectionFactory jedisConnectionFactory() {
        Set<Host> sentinels = new HashSet<>();
        sentinels.add(new Host("127.0.0.1", 30001));
        sentinels.add(new Host("127.0.0.1", 30002));
        sentinels.add(new Host("127.0.0.1", 30003));
        RedisSentinelConfiguration sentinelConfiguration = new RedisSentinelConfiguration();
        sentinelConfiguration.setMasterName("mymaster");
        sentinelConfiguration.setSentinels(sentinels);
        JedisConnectionFactory jedisConnectionFactory = new JedisConnectionFactory(sentinelConfiguration);
        return jedisConnectionFactory;
    }
}
```
## 2. JedisConnectionConfiguration


**那么 Spring Data Redis 的 `JedisConnectionFactory` 的自动配置在 **Spring Boot** 是如何实现的呢？**

**Spring Boot**  是通过内置的 `JedisConnectionConfiguration` 配置类来完成这一功能。下面我们具体分析一下：

> **注意：** 以下涉及 **Spring Boot** 源码 均来自版本 `2.7.9`，其他版本有所出入，可自行查看源码。

### 2.1 RedisConnectionConfiguration
翻看 `JedisConnectionConfiguration` 的源码，我们发现它继承了 `RedisConnectionConfiguration` 类，该类的部分源码如下：

```java
abstract class RedisConnectionConfiguration {
  private static final boolean COMMONS_POOL2_AVAILABLE = ClassUtils.isPresent("org.apache.commons.pool2.ObjectPool",
      RedisConnectionConfiguration.class.getClassLoader());
  // 。。。

  protected final RedisStandaloneConfiguration getStandaloneConfig() {
    // 。。。
  }

  protected final RedisSentinelConfiguration getSentinelConfig() {
    // 。。。
  }

  protected final RedisClusterConfiguration getClusterConfiguration() {
    // 。。。
  }

  protected final RedisProperties getProperties() {
    // 。。。
  }

  protected boolean isPoolEnabled(Pool pool) {
    Boolean enabled = pool.getEnabled();
    return (enabled != null) ? enabled : COMMONS_POOL2_AVAILABLE;
  }

  private List<RedisNode> createSentinels(RedisProperties.Sentinel sentinel) {
    // 。。。
  }

  protected ConnectionInfo parseUrl(String url) {
    // 。。。
  }

  static class ConnectionInfo {

    private final URI uri;

    private final boolean useSsl;

    private final String username;

    private final String password;

    // 。。。
  }
}
```

简单阅读上述的源码，我们可以很快总结一下：

- `getStandaloneConfig()` ：返回一个 `RedisStandaloneConfiguration` 对象，用于配置单机模式的 **Redis** 连接。
- `getSentinelConfig()` ：返回一个 `RedisSentinelConfiguration` 对象，用于配置哨兵模式的 **Redis** 连接。
- `getClusterConfiguration()` ：返回一个 `RedisClusterConfiguration` 对象，用于配置集群模式的 **Redis** 连接。
- `getProperties()` ：返回一个 `RedisProperties` 对象，用于获取 **Redis** 连接的相关配置信息。
- `isPoolEnabled(Pool pool)` ：判断给定的连接池是否启用。如果连接池的`enabled` 属性不为 `null`，则返回该属性值；否则返回`COMMONS_POOL2_AVAILABLE` 常量【如果`org.apache.commons.pool2.ObjectPool` 类存在，那么 `COMMONS_POOL2_AVAILABLE` 将被设置为 `true`，否则将被设置为 `false`】。
- `createSentinels(RedisProperties.Sentinel sentinel)` ：根据给定的哨兵配置创建一个 `RedisNode` 列表，用于配置哨兵模式的 `Redis` 连接。
- `parseUrl(String url)`：解析给定的 **URL** 字符串，并返回一个包含连接信息的 `ConnectionInfo` 对象。
其中内部静态类 `ConnectionInfo`，用于存储解析后的连接信息，包括：
  - `uri`：连接的 **URI**。
  - `useSsl`：是否使用 **SSL** 加密。
  - `username`：连接的用户名。
  - `password`：连接的密码。

### 2.2 导入自动配置
上篇博文中，我们已经知道了 `JedisConnectionConfiguration` 是在 `RedisAutoConfiguration` 中通过 `@Import({ LettuceConnectionConfiguration.class, JedisConnectionConfiguration.class })` 导入的。
### 2.3 相关注解介绍

我们在 `META-INF/spring-autoconfigure-metadata.properties` 文件中，发现了有关 `JedisConnectionConfiguration` 的相关配置：

```java
org.springframework.boot.autoconfigure.data.redis.JedisConnectionConfiguration=
org.springframework.boot.autoconfigure.data.redis.JedisConnectionConfiguration.ConditionalOnClass=org.apache.commons.pool2.impl.GenericObjectPool,redis.clients.jedis.Jedis,org.springframework.data.redis.connection.jedis.JedisConnection
```

显然这里涉及到了 `ConditionalOnClass` 注解，我们翻看 `JedisConnectionConfiguration` 配置类的源码，如下：

```java
@Configuration(proxyBeanMethods = false)
@ConditionalOnClass({ GenericObjectPool.class, JedisConnection.class, Jedis.class })
@ConditionalOnMissingBean(RedisConnectionFactory.class)
@ConditionalOnProperty(name = "spring.redis.client-type", havingValue = "jedis", matchIfMissing = true)
class JedisConnectionConfiguration extends RedisConnectionConfiguration {

  // 。。。

  @Bean
  JedisConnectionFactory redisConnectionFactory(
      ObjectProvider<JedisClientConfigurationBuilderCustomizer> builderCustomizers) {
    return createJedisConnectionFactory(builderCustomizers);
  }

  // 。。。
}
```

我们先来看看上述 `JedisConnectionConfiguration ` 配置类涉及到的注解，如下：
- `@Configuration(proxyBeanMethods = false)` : 该类是一个配置类，用于定义和配置 **Spring** 容器中的 **bean**。`proxyBeanMethods = false`表示不使用 **CGLIB** 代理来创建 **bean** 的子类实例。
- `@ConditionalOnClass({ GenericObjectPool.class, JedisConnection.class, Jedis.class })` ：只有在项目中存在 `GenericObjectPool`、`JedisConnection` 和 `Jedis` 这三个类时，才会加载这个配置类。这可以确保项目依赖中包含了这些类，避免因为缺少依赖而导致的配置错误。
- `@ConditionalOnMissingBean(RedisConnectionFactory.class)` ：表示只有在项目中不存在 `RedisConnectionFactory` 这个 **bean** 时，才会加载这个配置类。这可以确保项目没有重复定义相同的 **bean**，避免冲突。
- `@ConditionalOnProperty(name = "spring.redis.client-type", havingValue = "jedis", matchIfMissing = true)` ：只有在项目的配置文件中指定了 `spring.redis.client-type` 属性值为 `"jedis"` 时，才会加载这个配置类。如果配置文件中指定该属性值不是 `"jedis"`，则不会加载这个配置类。`matchIfMissing = true` 表示如果没有找到匹配的属性值，也会加载这个配置类。
- `@Bean` ：用于声明一个方法创建的对象是一个 **Spring** 管理的 **Bean**。**Spring** 容器会自动管理这个 **Bean** 的生命周期，包括依赖注入、初始化和销毁等。


### 2.4 redisConnectionFactory 方法

通过翻看 `JedisConnectionConfiguration` 的源码，我们可以看到 `redisConnectionFactory` 方法是被 `@Bean` 注解标注的，意味着该方法创建的 `Jedis` 连接工厂将成为 `Spring` 管理的 `Bean` 对象。

该方法接受一个入参 `ObjectProvider<JedisClientConfigurationBuilderCustomizer>`，它是一个提供者（**Provider**），它可以提供一个或多个`JedisClientConfigurationBuilderCustomizer` 对象。该对象是一个用于定制 **Jedis** 客户端配置的接口。通过实现这个接口，可以自定义 **Jedis** 客户端的配置，例如设置连接池大小、超时时间、**SSL** 配置等。这样我们就可以根据实际的需求灵活地调整 **Jedis** 客户端的行为。

进入 `redisConnectionFactory` 方法，我们看到它直接调用了 `createJedisConnectionFactory` 方法并返回一个 `JedisConnectionFactory` 对象。  

我们继续查看 `createJedisConnectionFactory` 方法：

```java
private JedisConnectionFactory createJedisConnectionFactory(
      ObjectProvider<JedisClientConfigurationBuilderCustomizer> builderCustomizers) {
  JedisClientConfiguration clientConfiguration = getJedisClientConfiguration(builderCustomizers);
  if (getSentinelConfig() != null) {
    return new JedisConnectionFactory(getSentinelConfig(), clientConfiguration);
  }
  if (getClusterConfiguration() != null) {
    return new JedisConnectionFactory(getClusterConfiguration(), clientConfiguration);
  }
  return new JedisConnectionFactory(getStandaloneConfig(), clientConfiguration);
}
```
我们详细来分析一下上述代码：

- 首先，调用 `getJedisClientConfiguration` 方法返回一个 `JedisClientConfiguration` 配置类对象。
  - 继续进入 `getJedisClientConfiguration` 方法：  

    ```java
    private JedisClientConfiguration getJedisClientConfiguration(
      ObjectProvider<JedisClientConfigurationBuilderCustomizer> builderCustomizers) {
      JedisClientConfigurationBuilder builder = applyProperties(JedisClientConfiguration.builder());
      RedisProperties.Pool pool = getProperties().getJedis().getPool();
      if (isPoolEnabled(pool)) {
        applyPooling(pool, builder);
      }
      if (StringUtils.hasText(getProperties().getUrl())) {
        customizeConfigurationFromUrl(builder);
      }
      builderCustomizers.orderedStream().forEach((customizer) -> customizer.customize(builder));
      return builder.build();
    }
    ```
    - 首先，调用 `applyProperties` 方法，获取一个 `JedisClientConfigurationBuilder`  对象，用于构建 `JedisClientConfiguration` 对象。
    

      ```java
      private JedisClientConfigurationBuilder applyProperties(JedisClientConfigurationBuilder builder) {
        PropertyMapper map = PropertyMapper.get().alwaysApplyingWhenNonNull();
        map.from(getProperties().isSsl()).whenTrue().toCall(builder::useSsl);
        map.from(getProperties().getTimeout()).to(builder::readTimeout);
        map.from(getProperties().getConnectTimeout()).to(builder::connectTimeout);
        map.from(getProperties().getClientName()).whenHasText().to(builder::clientName);
        return builder;
      }
      ```
      该方法的主要目的是根据属性配置来定制 `builder` 对象。
      - 首先，创建一个 `PropertyMapper` 对象 `map`，并调用其 `alwaysApplyingWhenNonNull()` 方法，以便在非空情况下始终应用映射规则。
      - 接下来，使用 `map.from()` 方法设置映射规则。这里分别设置了以下映射规则：

        - 如果属性中的 `isSsl` 为 `true`，则调用 `builder::useSsl` 方法，将 `builder` 对象的 `useSsl` 属性设置为 `true`。
        - 将属性中的 `timeout` 值设置为 `builder` 对象的 `readTimeout` 属性。
        - 将属性中的 `connectTimeout` 值设置为 `builder` 对象的 `connectTimeout` 属性。
        - 如果属性中的 `clientName` 有文本内容，则调用 `builder::clientName` 方法，将 `builder` 对象的 `clientName` 属性设置为该文本内容。
      - 最后，返回经过配置的 `builder` 对象。
    
    - 接着，从 `RedisProperties` 中获取 `Jedis` 连接池的配置信息。
      - `enabled` : 是否启用连接池。如果可用，则自动启用。在 **Jedis** 中，哨兵模式下的连接池是隐式启用的，此设置仅适用于单节点设置。

      - `maxIdle` : 池中空闲连接的最大数量。使用负值表示无限数量的空闲连接。

      - `minIdle` : 池中保持最小空闲连接的目标数量。此设置仅在空闲连接和驱逐运行之间的时间都为正时才有效。

      - `maxActive` : 给定时间内，连接池可以分配的最大连接数。使用负值表示无限制。

      - `maxWait` : 当连接池耗尽时，连接分配应阻塞的最长时间。使用负值表示无限期阻塞。

      - `timeBetweenEvictionRuns` : 空闲对象驱逐线程的运行时间间隔。当值为正时，空闲对象驱逐线程开始运行，否则不执行空闲对象驱逐。
    - 然后，判断连接池是否启用 ？
      - 如果启用，则调用 `applyPooling` 方法，将连接池配置应用到 `builder` 对象上。

        ```java
        private void applyPooling(RedisProperties.Pool pool,
            JedisClientConfiguration.JedisClientConfigurationBuilder builder) {
          builder.usePooling().poolConfig(jedisPoolConfig(pool));
        }
      
        private JedisPoolConfig jedisPoolConfig(RedisProperties.Pool pool) {
          JedisPoolConfig config = new JedisPoolConfig();
          config.setMaxTotal(pool.getMaxActive());
          config.setMaxIdle(pool.getMaxIdle());
          config.setMinIdle(pool.getMinIdle());
          if (pool.getTimeBetweenEvictionRuns() != null) {
            config.setTimeBetweenEvictionRuns(pool.getTimeBetweenEvictionRuns());
          }
          if (pool.getMaxWait() != null) {
            config.setMaxWait(pool.getMaxWait());
          }
          return config;
        }
        ```
        - `usePooling()`  : 启用连接池功能
        - `poolConfig(jedisPoolConfig(pool))` ：将连接池的配置信息传递给 `builder` 对象

    - 判断属性中的 `spring.redis.url` 是否包含非空的文本内容？
      - 如果包含非空的文本内容 ，则调用 customizeConfigurationFromUrl 方法：
        ```java
        private void customizeConfigurationFromUrl(JedisClientConfiguration.JedisClientConfigurationBuilder builder) {
          ConnectionInfo connectionInfo = parseUrl(getProperties().getUrl());
          if (connectionInfo.isUseSsl()) {
            builder.useSsl();
          }
        }
        ```
        - 首先，调用 `parseUrl` 方法来解析 `URL`，并将结果存储在 `connectionInfo` 变量中。
        - 然后，调用 `connectionInfo#isUseSsl` 方法，判断是否需要使用 **SSL** 连接。如果需要使用 **SSL** 连接，则调用 `builder` 对象的 `useSsl()` 方法来启用 **SSL** 功能。
    - 遍历 `builderCustomizers` 中的所有自定义器，并对每个自定义器调用其 `customize` 方法，传入 `builder` 作为参数，用于进一步定制 **Jedis** 客户端的配置。

- 接着，获取哨兵模式配置，并判断是否为空，如果不为空，则直接根据哨兵模式的配置创建并返回一个连接工厂实例。
- 然后，获取集群模式配置，并判断是否为空，如果不为空，则直接根据集群模式的配置创建并返回一个连接工厂实例。
- 最后，获取单机模式配置，根据单机模式的配置创建并返回一个连接工厂实例。

# 四、总结
本篇我们深入分析了 `JedisConnectionConfiguration ` 配置类的相关内容，该类用于配置 **Redis** 连接的 **Jedis** 实现。

