---
title: flea-cache使用之Redis哨兵模式接入
date: 2025-02-15 22:59:31
updated: 2025-02-15 22:59:31
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - flea-framework
  - flea-cache
  - Redis
  - Redis哨兵模式
---

![](/images/cache.png)

# 1. 参考
[flea-cache使用之Redis哨兵模式接入 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-cache)

<!-- more -->

![](flea-cache-redissentinel.png)

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 2. 依赖
[jedis-3.0.1.jar](https://mvnrepository.com/artifact/redis.clients/jedis/3.0.1)
```xml
<!-- Java redis -->
<dependency>
     <groupId>redis.clients</groupId>
     <artifactId>jedis</artifactId>
     <version>3.0.1</version>
</dependency>
```
[spring-context-4.3.18.RELEASE.jar](https://mvnrepository.com/artifact/org.springframework/spring-context/4.3.18.RELEASE)
```xml
<!-- Spring相关 -->
<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-context</artifactId>
    <version>4.3.18.RELEASE</version>
</dependency>
```
[spring-context-support-4.3.18.RELEASE.jar](https://mvnrepository.com/artifact/org.springframework/spring-context-support/4.3.18.RELEASE)
```xml
<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-context-support</artifactId>
    <version>4.3.18.RELEASE</version>
</dependency>
```

# 3. 基础接入
## 3.1 定义Flea缓存接口 
 [IFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/IFleaCache.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 3.2 定义抽象Flea缓存类
[AbstractFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCache.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 3.3 定义Redis客户端接口类 
[RedisClient](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClient.java) 可参考笔者的这篇博文 [Redis集群模式接入](../../../../../../2021/11/25/flea-framework/flea-cache/flea-cache-rediscluster/)，不再赘述。

## 3.4 定义Redis客户端命令行
[RedisClientCommand](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClientCommand.java) 可参考笔者的这篇博文 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，不再赘述。

## 3.5 定义哨兵模式Redis客户端实现类
[FleaRedisSentinelClient](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/FleaRedisSentinelClient.java) 即**Flea**哨兵模式**Redis**客户端实现，封装了**Flea框架**操作**Redis**缓存的基本操作。它内部具体操作**Redis**缓存的功能，由**Jedis**哨兵连接池获取**Jedis**实例对象完成，包含读、写、删除**Redis**缓存的基本操作方法。

哨兵模式下，单个缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = new FleaRedisSentinelClient.Builder().build();
   // 执行读，写，删除等基本操作
   redisClient.set("key", "value"); 
```
哨兵模式下，整合缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = new FleaRedisSentinelClient.Builder(poolName).build();
   // 执行读，写，删除等基本操作
   redisClient.set("key", "value");  
```

当然每次都新建**Redis**客户端显然不可取，我们可通过**Redis**客户端工厂获取**Redis**客户端。
哨兵模式下，单个缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = RedisClientFactory.getInstance(CacheModeEnum.SENTINEL);
```

哨兵模式下，整合缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = RedisClientFactory.getInstance(poolName, CacheModeEnum.SENTINEL); 
```

```java
public class FleaRedisSentinelClient extends FleaRedisClient {

    private JedisSentinelPool jedisSentinelPool;

    private int maxAttempts;

    /**
     * Redis哨兵客户端构造方法 (默认)
     *
     * @since 2.0.0
     */
    private FleaRedisSentinelClient() {
        this(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * Redis哨兵客户端构造方法（指定连接池名）
     *
     * @param poolName 连接池名
     * @since 2.0.0
     */
    private FleaRedisSentinelClient(String poolName) {
        super(poolName);
        init();
    }

    /**
     * 初始化Jedis哨兵实例
     *
     * @since 2.0.0
     */
    private void init() {
        if (CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME.equals(getPoolName())) {
            jedisSentinelPool = RedisSentinelPool.getInstance().getJedisSentinelPool();
            maxAttempts = RedisSentinelConfig.getConfig().getMaxAttempts();
        } else {
            jedisSentinelPool = RedisSentinelPool.getInstance(getPoolName()).getJedisSentinelPool();
            maxAttempts = CacheConfigUtils.getMaxAttempts();
        }
    }

    @Override
    public String set(final String key, final Object value) {
        return new RedisClientCommand<String, JedisSentinelPool, Jedis>(this.jedisSentinelPool, this.maxAttempts) {
            @Override
            public String execute(Jedis connection) {
                if (value instanceof String)
                    return connection.set(key, (String) value);
                else
                    return connection.set(SafeEncoder.encode(key), ObjectUtils.serialize(value));
            }
        }.run();
    }

    // 省略。。。。。。

    /**
     * 内部建造者类
     *
     * @author huazie
     * @version 2.0.0
     * @since 2.0.0
     */
    public static class Builder {

        private String poolName; // 连接池名

        /**
         * 默认构造器
         *
         * @since 2.0.0
         */
        public Builder() {
        }

        /**
         * 指定连接池的构造器
         *
         * @param poolName 连接池名
         * @since 2.0.0
         */
        public Builder(String poolName) {
            this.poolName = poolName;
        }

        /**
         * 构建Redis哨兵客户端对象
         *
         * @return Redis哨兵客户端
         * @since 2.0.0
         */
        public RedisClient build() {
            if (StringUtils.isBlank(poolName)) {
                return new FleaRedisSentinelClient();
            } else {
                return new FleaRedisSentinelClient(poolName);
            }
        }
    }
}
```
该类的构造函数初始化逻辑，可以看出我们使用了 `RedisSentinelPool`， 下面来介绍一下。

## 3.6 定义Redis哨兵连接池
我们使用 [RedisSentinelPool](https://github.com/Huazie/flea-framework/blob/main/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisSentinelPool.java) 来初始化**Jedis**哨兵连接池实例，其中重点是获取分布式**Jedis**连接池 `ShardedJedisPool` ，该类其中一个构造方法如下：

```java
public JedisSentinelPool(String masterName, Set<String> sentinels, final GenericObjectPoolConfig poolConfig, 
    final int connectionTimeout, final int soTimeout, final String password, final int database, final String clientName) {
       
}
```

**Redis**哨兵连接池，用于初始化**Jedis**哨兵连接池实例。

针对单独缓存接入场景，采用默认连接池初始化的方式； 可参考如下：
```java
    // 初始化默认连接池
    RedisSentinelPool.getInstance().initialize();
``` 
针对整合缓存接入场景，采用指定连接池初始化的方式； 可参考如下：
```java
    // 初始化指定连接池
    RedisSentinelPool.getInstance(group).initialize(cacheServerList);
```

```java
public class RedisSentinelPool {

    private static final ConcurrentMap<String, RedisSentinelPool> redisPools = new ConcurrentHashMap<>();

    private static final Object redisSentinelPoolLock = new Object();

    private String poolName; // 连接池名

    private JedisSentinelPool jedisSentinelPool; // Jedis哨兵连接池

    private RedisSentinelPool(String poolName) {
        this.poolName = poolName;
    }

    /**
     * 获取Redis哨兵连接池实例 (默认连接池)
     *
     * @return Redis哨兵连接池实例对象
     * @since 2.0.0
     */
    public static RedisSentinelPool getInstance() {
        return getInstance(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * 获取Redis哨兵连接池实例 (指定连接池名)
     *
     * @param poolName 连接池名
     * @return Redis哨兵连接池实例对象
     * @since 2.0.0
     */
    public static RedisSentinelPool getInstance(String poolName) {
        if (!redisPools.containsKey(poolName)) {
            synchronized (redisSentinelPoolLock) {
                if (!redisPools.containsKey(poolName)) {
                    RedisSentinelPool redisSentinelPool = new RedisSentinelPool(poolName);
                    redisPools.put(poolName, redisSentinelPool);
                }
            }
        }
        return redisPools.get(poolName);
    }

    /**
     * 默认初始化
     *
     * @since 2.0.0
     */
    public void initialize(int database) {
        // 省略。。。。。。
    }

    /**
     * 初始化 (非默认连接池)
     *
     * @param cacheServerList 缓存服务器集
     * @since 2.0.0
     */
    public void initialize(List<CacheServer> cacheServerList) {
        // 省略。。。。。。
    }

    /**
     * Jedis哨兵连接池
     *
     * @return Jedis哨兵连接池
     * @since 2.0.0
     */
    public JedisSentinelPool getJedisSentinelPool() {
        if (ObjectUtils.isEmpty(jedisSentinelPool)) {
            ExceptionUtils.throwFleaException(FleaCacheConfigException.class, "获取Jedis哨兵连接池失败：请先调用initialize初始化");
        }
        return jedisSentinelPool;
    }
}
```

## 3.7 定义Redis哨兵配置文件
**flea-cache** 读取 [redis.sentinel.properties](https://github.com/Huazie/flea-framework/blob/main/flea-config/src/main/resources/flea/cache/redis.sentinel.properties)（Redis哨兵配置文件），用作初始化 **RedisSentinelPool** 

```bash
# Redis哨兵配置
redis.sentinel.switch=1

redis.systemName=FleaFrame

redis.sentinel.masterName=mymaster

redis.sentinel.server=127.0.0.1:36379,127.0.0.1:36380,127.0.0.1:36381

#redis.sentinel.password=huazie123

redis.sentinel.connectionTimeout=2000

redis.sentinel.soTimeout=2000

# Redis哨兵客户端连接池配置
redis.pool.maxTotal=100

redis.pool.maxIdle=10

redis.pool.minIdle=0

redis.pool.maxWaitMillis=2000

redis.maxAttempts=5

redis.nullCacheExpiry=10
```

- `redis.sentinel.switch` : **Redis**哨兵配置开关（1：开启 0：关闭），如果不配置也默认开启
- `redis.systemName` : **Redis**缓存所属系统名
- `redis.sentinel.masterName` : **Redis**主服务器节点名称
- `redis.sentinel.server` : **Redis**哨兵节点的地址集合
- `redis.sentinel.password` : **Redis**主从服务器节点登录密码（各节点配置同一个）
- `redis.sentinel.connectionTimeout` : **Redis**哨兵客户端**socket**连接超时时间（单位：ms）
- `redis.sentinel.soTimeout` : **Redis**哨兵客户端**socket**读写超时时间（单位：ms）
- `redis.pool.maxTotal` : **Jedis**连接池最大连接数
- `redis.pool.maxIdle` : **Jedis**连接池最大空闲连接数
- `redis.pool.minIdle` : **Jedis**连接池最小空闲连接数
- `redis.pool.maxWaitMillis` : **Jedis**连接池获取连接时的最大等待时间（单位：ms）
- `redis.maxAttempts` : **Redis**客户端操作最大尝试次数【包含第一次操作】
- `redis.nullCacheExpiry` : 空缓存数据有效期（单位：s）

## 3.8 定义Redis Flea缓存类
[RedisFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/RedisFleaCache.java) 可参考笔者的这篇博文 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，不再赘述。

## 3.9 定义抽象Flea缓存管理类 
[AbstractFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCacheManager.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 3.10 定义Redis哨兵模式Flea缓存管理类
[RedisSentinelFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/manager/RedisSentinelFleaCacheManager.java)  继承抽象Flea缓存管理类 `AbstractFleaCacheManager`，用于接入**Flea**框架管理**Redis**缓存。

它的默认构造方法，用于初始化哨兵模式下默认连接池的**Redis**客户端【默认**Redis**数据库索引为0】, 这里需要先初始化**Redis**哨兵连接池，默认连接池名为【`default`】； 然后通过 `RedisClientFactory` 获取哨兵模式下默认连接池的Redis客户端 `RedisClient`，可在 **3.11** 查看。

它的带参数的构造方法，用于初始化哨兵模式下默认连接池的Redis客户端【指定Redis数据库索引】。

方法 `newCache` 用于创建一个 `RedisFleaCache` 的实例对象，它里面包含了 读、写、删除 和 清空 缓存的基本操作，每一类 **Redis** 缓存数据都对应了一个 `RedisFleaCache` 的实例对象。

```java
public class RedisSentinelFleaCacheManager extends AbstractFleaCacheManager {

    private RedisClient redisClient; // Redis客户端

    /**
     * 默认构造方法，初始化哨兵模式下默认连接池的Redis客户端
     *
     * @since 2.0.0
     */
    public RedisSentinelFleaCacheManager() {
        this(0);
    }

    /**
     * 初始化哨兵模式下默认连接池的Redis客户端，指定Redis数据库索引
     *
     * @since 2.0.0
     */
    public RedisSentinelFleaCacheManager(int database) {
        if (!RedisSentinelConfig.getConfig().isSwitchOpen()) return;
        // 初始化默认连接池
        RedisSentinelPool.getInstance().initialize(database);
        // 获取哨兵模式下默认连接池的Redis客户端
        redisClient = RedisClientFactory.getInstance(CacheModeEnum.SENTINEL);
    }

    @Override
    protected AbstractFleaCache newCache(String name, int expiry) {
        int nullCacheExpiry = RedisSentinelConfig.getConfig().getNullCacheExpiry();
        if (RedisSentinelConfig.getConfig().isSwitchOpen())
            return new RedisFleaCache(name, expiry, nullCacheExpiry, CacheModeEnum.SENTINEL, redisClient);
        else
            return new EmptyFleaCache(name, expiry, nullCacheExpiry);
    }
}
```

## 3.11 定义Redis客户端工厂类 
[RedisClientFactory](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClientFactory.java) ，有四种方式获取 **Redis** 客户端：
 - 一是获取分片模式下默认连接池的 **Redis** 客户端，应用在单个缓存接入场景；
 - 二是获取指定模式下默认连接池的 **Redis** 客户端，应用在单个缓存接入场景【**3.10** 采用】；
 - 三是获取分片模式下指定连接池的 **Redis** 客户端，应用在整合缓存接入场景；
 - 四是获取指定模式下指定连接池的 **Redis** 客户端，应用在整合缓存接入场景。

```java

/**
 * Redis客户端工厂，用于获取Redis客户端。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public class RedisClientFactory {

    private static final ConcurrentMap<String, RedisClient> redisClients = new ConcurrentHashMap<>();

    private static final Object redisClientLock = new Object();

    private RedisClientFactory() {
    }

    /**
     * 获取分片模式下默认连接池的Redis客户端
     *
     * @return 分片模式的Redis客户端
     * @since 1.0.0
     */
    public static RedisClient getInstance() {
        return getInstance(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * 获取指定模式下默认连接池的Redis客户端
     *
     * @param mode 缓存模式
     * @return 指定模式的Redis客户端
     * @since 1.1.0
     */
    public static RedisClient getInstance(CacheModeEnum mode) {
        return getInstance(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME, mode);
    }

    /**
     * 获取分片模式下指定连接池的Redis客户端
     *
     * @param poolName 连接池名
     * @return 分片模式的Redis客户端
     * @since 1.0.0
     */
    public static RedisClient getInstance(String poolName) {
        return getInstance(poolName, CacheModeEnum.SHARDED);
    }

    /**
     * 获取指定模式下指定连接池的Redis客户端
     *
     * @param poolName 连接池名
     * @param mode     缓存模式
     * @return 指定模式的Redis客户端
     * @since 1.1.0
     */
    public static RedisClient getInstance(String poolName, CacheModeEnum mode) {
        String key = StringUtils.strCat(poolName, CommonConstants.SymbolConstants.UNDERLINE, StringUtils.valueOf(mode.getMode()));
        if (!redisClients.containsKey(key)) {
            synchronized (redisClientLock) {
                if (!redisClients.containsKey(key)) {
                    RedisClientStrategyContext context = new RedisClientStrategyContext(poolName);
                    redisClients.putIfAbsent(key, FleaStrategyFacade.invoke(mode.name(), context));
                }
            }
        }
        return redisClients.get(key);
    }
}

```
在上面 的 `getInstance(String poolName, CacheModeEnum mode)` 方法中，使用了 **RedisClientStrategyContext** ，用于定义 **Redis** 客户端策略上下文。根据不同的缓存模式，就可以找到对应的 **Redis** 客户端策略。

## 3.12 定义Redis客户端策略上下文
[RedisClientStrategyContext](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/common/RedisClientStrategyContext.java)  可参考笔者的这篇博文 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，不再赘述。

## 3.13 定义哨兵模式Redis客户端策略
 [RedisSentinelClientStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/strategy/RedisSentinelClientStrategy.java) 用于新建一个 `Flea Redis` 哨兵客户端。

```java
/**
 * 哨兵模式Redis客户端 策略
 *
 * @author huazie
 * @version 2.0.0
 * @since 2.0.0
 */
public class RedisSentinelClientStrategy implements IFleaStrategy<RedisClient, String> {

    @Override
    public RedisClient execute(String poolName) throws FleaStrategyException {
        RedisClient originRedisClient;
        // 新建一个Flea Redis哨兵客户端类实例
        if (CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME.equals(poolName)) {
            originRedisClient = new FleaRedisSentinelClient.Builder().build();
        } else {
            originRedisClient = new FleaRedisSentinelClient.Builder(poolName).build();
        }
        return originRedisClient;
    }
}
```
好了，到这里我们可以来测试 Redis 哨兵模式。

## 3.14 Redis集群模式接入自测
单元测试类  [FleaCacheTest](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/test/java/com/huazie/fleaframework/cache/FleaCacheTest.java)

首先，这里需要按照 **Redis哨兵配置文件** 中的地址部署相应的 **Redis哨兵服务** 和 **Redis主从服务**，后续有机会我再出一篇简单的Redis主从 + 哨兵的搭建博文。

```java
    @Test
    public void testRedisSentinelFleaCache() {
        try {
            // 哨兵模式下Flea缓存管理类，复用原有获取方式
//            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(CacheEnum.RedisSentinel.getName());
            // 哨兵模式下Flea缓存管理类，指定数据库索引
            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(0);
            AbstractFleaCache cache = manager.getCache("fleajerseyresource");
            LOGGER.debug("Cache={}", cache);
            //#### 1.  简单字符串
            cache.put("author", "huazie");
            cache.put("other", null);
//            cache.get("author");
//            cache.get("other");
//            cache.delete("author");
//            cache.delete("other");
//            cache.clear();
            cache.getCacheKey();
            LOGGER.debug(cache.getCacheName() + ">>>" + cache.getCacheDesc());
        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
```

# 4. 进阶接入
## 4.1 定义抽象Spring缓存 
[AbstractSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCache.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 4.2 定义Redis Spring缓存类
 [RedisSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/RedisSpringCache.java) 可参考笔者的这篇博文 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，不再赘述。

## 4.3 定义抽象Spring缓存管理类
[AbstractSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCacheManager.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 4.4 定义Redis哨兵模式Spring缓存管理类
 [RedisSentinelSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/manager/RedisSentinelSpringCacheManager.java) 继承抽象 **Spring** 缓存管理类 `AbstractSpringCacheManager`，用于接入**Spring**框架管理**Redis**缓存; 基本实现同 **RedisSentinelFleaCacheManager**，唯一不同在于 **newCache** 的实现。

它的默认构造方法，用于初始化哨兵模式下默认连接池的**Redis**客户端【默认**Redis**数据库索引为0】, 这里需要先初始化**Redis**哨兵连接池，默认连接池名为【`default`】； 然后通过 `RedisClientFactory` 获取哨兵模式下默认连接池的Redis客户端 `RedisClient`，可在 **3.11** 查看。

它的带参数的构造方法，用于初始化哨兵模式下默认连接池的Redis客户端【指定Redis数据库索引】。

方法【`newCache`】用于创建一个**Redis Spring**缓存， 而它内部是由**Redis Flea**缓存实现具体的 读、写、删除 和 清空 缓存的基本操作。

```java
public class RedisSentinelSpringCacheManager extends AbstractSpringCacheManager {

    private RedisClient redisClient; // Redis客户端

    /**
     * 默认构造方法，初始化哨兵模式下默认连接池的Redis客户端
     *
     * @since 2.0.0
     */
    public RedisSentinelSpringCacheManager() {
        this(0);
    }

    /**
     * 初始化哨兵模式下默认连接池的Redis客户端，指定Redis数据库索引
     *
     * @since 2.0.0
     */
    public RedisSentinelSpringCacheManager(int database) {
        if (!RedisSentinelConfig.getConfig().isSwitchOpen()) return;
        // 初始化默认连接池
        RedisSentinelPool.getInstance().initialize(database);
        // 获取哨兵模式下默认连接池的Redis客户端
        redisClient = RedisClientFactory.getInstance(CacheModeEnum.SENTINEL);
    }

    @Override
    protected AbstractSpringCache newCache(String name, int expiry) {
        int nullCacheExpiry = RedisSentinelConfig.getConfig().getNullCacheExpiry();
        if (RedisSentinelConfig.getConfig().isSwitchOpen())
            return new RedisSpringCache(name, expiry, nullCacheExpiry, CacheModeEnum.SENTINEL, redisClient);
        else
            return new RedisSpringCache(name, new EmptyFleaCache(name, expiry, nullCacheExpiry));
    }

}
```

## 4.5 spring 配置

```xml
    <!--
        配置缓存管理 redisSentinelSpringCacheManager
        配置缓存时间 configMap (key缓存对象名称 value缓存过期时间)
    -->
    <bean id="redisSentinelSpringCacheManager" class="com.huazie.fleaframework.cache.redis.manager.RedisSentinelSpringCacheManager">
        <!-- 使用带参数的构造函数实例化，指定Redis数据库索引 -->
        <!--<constructor-arg index="0" value="0"/>-->
        <property name="configMap">
            <map>
                <entry key="fleajerseyi18nerrormapping" value="86400"/>
                <entry key="fleajerseyresservice" value="86400"/>
                <entry key="fleajerseyresclient" value="86400"/>
                <entry key="fleajerseyresource" value="86400"/>
            </map>
        </property>
    </bean>

    <!-- 开启缓存 -->
    <cache:annotation-driven cache-manager="redisSentinelSpringCacheManager" proxy-target-class="true"/>
```

## 4.6 缓存自测

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class SpringCacheTest {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(SpringCacheTest.class);

    @Autowired
    @Qualifier("redisSentinelSpringCacheManager")
    private AbstractSpringCacheManager redisSentinelSpringCacheManager;

    @Test
    public void testRedisSentinelSpringCache() {
        try {
            // 哨兵模式下Spring缓存管理类
            AbstractSpringCache cache = redisSentinelSpringCacheManager.getCache("fleajerseyresource");
            LOGGER.debug("Cache = {}", cache);

            //#### 1.  简单字符串
            cache.put("menu1", "huazie");
            cache.put("menu2", null);
//            cache.get("menu1");
//            cache.get("menu2");
//            cache.getCacheKey();
//            cache.delete("menu1");
//            cache.delete("menu2");
//            cache.clear();
            cache.getCacheKey();
            AbstractFleaCache fleaCache = (AbstractFleaCache) cache.getNativeCache();
            LOGGER.debug(fleaCache.getCacheName() + ">>>" + fleaCache.getCacheDesc());
        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
}
```

# 结语

到这一步，**Redis** 哨兵模式单独接入的内容终于搞定了，有关整合接入**Redis**哨兵模式的，请查看笔者的[《整合Memcached和Redis接入》](../../../../../../2019/08/23/flea-framework/flea-cache/flea-cache-corecache/)。