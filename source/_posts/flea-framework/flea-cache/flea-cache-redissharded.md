---
title: flea-cache使用之Redis分片模式接入
date: 2021-11-18 23:12:25 
updated: 2023-12-20 16:58:26
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - flea-framework
  - flea-cache
  - Redis
  - Redis分片模式
---

![](/images/cache.png)

# 1. 参考
[flea-cache使用之Redis分片模式接入 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-cache)
<!-- more -->

![](flea-cache-redissharded.png)

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
[AbstractFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCache.java)  可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 3.3 定义Redis客户端接口类 
[RedisClient](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClient.java) , 注意该版，相比《[flea-frame-cache使用之Redis接入【旧】](../../../../../../2019/08/19/flea-framework/flea-cache/flea-frame-cache-redis/)》博文中，废弃如下与 `ShardedJedis` 有关的方法：
```java
    ShardedJedisPool getJedisPool();
    
    void setShardedJedis(ShardedJedis shardedJedis);
    
    ShardedJedis getShardedJedis();
```
《[flea-frame-cache使用之Redis接入【旧】](../../../../../../2019/08/19/flea-framework/flea-cache/flea-frame-cache-redis/)》博文中 提到了使用 **Redis**客户端代理方式 访问 `RedisClient`， 在这版为了实现 **Redis** 访问异常后的重试机制，废弃了代理模式，采用了命令行模式，可参考下面的 `RedisClientCommand` 。

## 3.4 定义Redis客户端命令行
[RedisClientCommand](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClientCommand.java) 封装了使用`ShardedJedis`操作**Redis**缓存的公共逻辑，如果出现异常可以重试 `maxAttempts` 次。

抽象方法 `execute` ，由子类或匿名类实现。在实际调用前，需要从分布式**Jedis**连接池中获取分布式**Jedis**对象；调用结束后， 关闭分布式**Jedis**对象，归还给分布式**Jedis**连接池。

```java
public abstract class RedisClientCommand<T> {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(RedisClientCommand.class);

    private final ShardedJedisPool shardedJedisPool; // 分布式Jedis连接池

    private final int maxAttempts; // Redis客户端操作最大尝试次数【包含第一次操作】

    public RedisClientCommand(ShardedJedisPool shardedJedisPool, int maxAttempts) {
        this.shardedJedisPool = shardedJedisPool;
        this.maxAttempts = maxAttempts;
    }

    public abstract T execute(ShardedJedis connection);

    /**
     * 执行分布式Jedis操作
     *
     * @return 分布式Jedis对象操作的结果
     * @since 1.0.0
     */
    public T run() {
        return runWithRetries(this.maxAttempts);
    }

    /**
     * 执行分布式Jedis操作，如果出现异常，包含第一次操作，可最多尝试maxAttempts次。
     *
     * @param attempts 重试次数
     * @return 分布式Jedis对象操作的结果
     * @since 1.0.0
     */
    private T runWithRetries(int attempts) {
        if (attempts <= 0) {
            throw new FleaCacheMaxAttemptsException("No more attempts left.");
        }
        ShardedJedis connection = null;
        try {
            connection = shardedJedisPool.getResource();
            Object obj = null;
            if (LOGGER.isDebugEnabled()) {
                obj = new Object() {};
                LOGGER.debug1(obj, "Get ShardedJedis = {}", connection);
            }
            T result = execute(connection);
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug1(obj, "Result = {}", result);
            }
            return result;
        } catch (JedisConnectionException e) {
            // 在开始下一次尝试前，释放当前分布式Jedis的连接，将分布式Jedis对象归还给分布式Jedis连接池
            releaseConnection(connection);
            connection = null; // 这里置空是为了最后finally不重复操作
            if (LOGGER.isErrorEnabled()) {
                Object obj = new Object() {};
                LOGGER.error1(obj, "Redis连接异常：", e);
                int currAttempts = this.maxAttempts - attempts + 1;
                LOGGER.error1(obj, "第 {} 次尝试失败，开始第 {} 次尝试...", currAttempts, currAttempts + 1);
            }
            return runWithRetries(attempts - 1);
        } finally {
            releaseConnection(connection);
        }
    }

    /**
     * 释放指定分布式Jedis的连接，将分布式Jedis对象归还给分布式Jedis连接池
     *
     * @param connection 分布式Jedis实例
     * @since 1.0.0
     */
    private void releaseConnection(ShardedJedis connection) {
        if (ObjectUtils.isNotEmpty(connection)) {
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug1(new Object() {}, "Close ShardedJedis");
            }
            connection.close();
        }
    }
}

```
## 3.5 定义分片模式Redis客户端实现类
[FleaRedisShardedClient](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/FleaRedisShardedClient.java) 主要使用 `ShardedJedis` 来操作 **Redis** 数据，它封装了**Flea框架**操作**Redis**缓存的基本操作。

它内部具体操作**Redis**缓存的功能，由分布式**Jedis**对象完成， 包含读、写、删除**Redis**缓存的基本操作方法。

分片模式下，单个缓存接入场景，可通过如下方式使用：
```java
    RedisClient redisClient = new FleaRedisShardedClient.Builder().build();
    // 执行读，写，删除等基本操作
    redisClient.set("key", "value"); 
```

分片模式下，整合缓存接入场景，可通过如下方式使用：
```java
    RedisClient redisClient = new FleaRedisShardedClient.Builder(poolName).build();
    // 执行读，写，删除等基本操作
    redisClient.set("key", "value"); 
```

当然每次都新建**Redis**客户端显然不可取，我们可通过**Redis**客户端工厂获取**Redis**客户端。
分片模式下，单个缓存接入场景，可通过如下方式使用：
```java
    RedisClient redisClient = RedisClientFactory.getInstance();
    // 或者
    RedisClient redisClient = RedisClientFactory.getInstance(CacheModeEnum.SHARDED);
```

分片模式下，整合缓存接入场景，可通过如下方式使用：
```java
    RedisClient redisClient = RedisClientFactory.getInstance(poolName);
    // 或者
    RedisClient redisClient = RedisClientFactory.getInstance(poolName, CacheModeEnum.SHARDED);
```
    

```java
public class FleaRedisShardedClient extends FleaRedisClient {

    private ShardedJedisPool shardedJedisPool; // 分布式Jedis连接池

    private int maxAttempts; // Redis客户端操作最大尝试次数【包含第一次操作】

    /**
     * <p> Redis客户端构造方法 (默认连接池名) </p>
     *
     * @since 1.0.0
     */
    private FleaRedisShardedClient() {
        this(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * <p> Redis客户端构造方法（指定连接池名） </p>
     *
     * @param poolName 连接池名
     * @since 1.0.0
     */
    private FleaRedisShardedClient(String poolName) {
        super(poolName);
        init();
    }

    /**
     * <p> 初始化分布式Jedis连接池 </p>
     *
     * @since 1.0.0
     */
    private void init() {
        if (CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME.equals(getPoolName())) {
            shardedJedisPool = RedisShardedPool.getInstance().getJedisPool();
            maxAttempts = RedisShardedConfig.getConfig().getMaxAttempts();
        } else {
            shardedJedisPool = RedisShardedPool.getInstance(getPoolName()).getJedisPool();
            maxAttempts = CacheConfigUtils.getMaxAttempts();
        }
    }

    @Override
    public String set(final String key, final Object value) {
        return new RedisClientCommand<String>(this.shardedJedisPool, this.maxAttempts) {
            @Override
            public String execute(ShardedJedis connection) {
                if (value instanceof String)
                    return connection.set(key, (String) value);
                else
                    return connection.set(SafeEncoder.encode(key), ObjectUtils.serialize(value));
            }
        }.run();
    }

    @Override
    public String set(final byte[] key, final byte[] value) {
        // 省略。。。。。。
    }

    @Override
    public String set(final String key, final Object value, final int expiry) {
        // 省略。。。。。。
    }

    @Override
    public String set(final byte[] key, final byte[] value, final int expiry) {
        // 省略。。。。。。
    }

    @Override
    public String set(final String key, final Object value, final long expiry) {
        // 省略。。。。。。
    }

    @Override
    public String set(final byte[] key, final byte[] value, final long expiry) {
        // 省略。。。。。。
    }

    @Override
    public String set(final String key, final Object value, final SetParams params) {
        // 省略。。。。。。
    }

    @Override
    public String set(final byte[] key, final byte[] value, final SetParams params) {
        // 省略。。。。。。
    }

    @Override
    public byte[] get(final byte[] key) {
        // 省略。。。。。。
    }

    @Override
    public Long del(final String key) {
        // 省略。。。。。。
    }

    /**
     * <p> 获取客户端类 </p>
     *
     * @param key 数据键
     * @return 客户端类
     * @since 1.0.0
     */
    @Override
    protected Client getClientByKey(final Object key) {
        // 省略。。。。。。
    }

    /**
     * <p> 内部建造者类 </p>
     */
    public static class Builder {
        // 省略。。。。。。
    }
}

```
该类的构造函数初始化逻辑，可以看出我们使用了 `RedisShardedPool` ， 下面来介绍一下。

## 3.6 定义Redis分片连接池
[RedisShardedPool](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisShardedPool.java) ，上个版本我们使用 `RedisPool` 初始化**Redis**相关配置信息，为了体现**Redis**分片模式，这个版本里面，我们使用 `RedisShardedPool` 用于**Redis**相关配置信息的初始化，其中重点是获取分布式**Jedis**连接池 `ShardedJedisPool` ，该类其中一个构造方法如下：
```java
/**
 * @param poolConfig 连接池配置信息
 * @param shards Jedis分布式服务器列表
 * @param algo 分布式算法
 */
public ShardedJedisPool(final GenericObjectPoolConfig poolConfig, List<JedisShardInfo> shards,
      Hashing algo) 
```

**Redis**分片连接池，用于初始化分布式 **Jedis** 连接池。

针对单独缓存接入场景，采用默认连接池初始化的方式； 可参考如下：
```java
    // 初始化默认连接池
    RedisShardedPool.getInstance().initialize();
``` 
针对整合缓存接入场景，采用指定连接池初始化的方式； 可参考如下：
```java
    // 初始化指定连接池
    RedisShardedPool.getInstance(group).initialize(cacheServerList); 
```

```java
public class RedisShardedPool {

    private static final ConcurrentMap<String, RedisShardedPool> redisPools = new ConcurrentHashMap<>();

    private String poolName; // 连接池名

    private ShardedJedisPool shardedJedisPool; // 分布式Jedis连接池

    private RedisShardedPool(String poolName) {
        this.poolName = poolName;
    }

    /**
     * <p> 获取Redis连接池实例 (默认连接池) </p>
     *
     * @return Redis连接池实例对象
     * @since 1.0.0
     */
    public static RedisShardedPool getInstance() {
        return getInstance(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * <p> 获取Redis连接池实例 (指定连接池名) </p>
     *
     * @param poolName 连接池名
     * @return Redis连接池实例对象
     * @since 1.0.0
     */
    public static RedisShardedPool getInstance(String poolName) {
        if (!redisPools.containsKey(poolName)) {
            synchronized (redisPools) {
                if (!redisPools.containsKey(poolName)) {
                    RedisShardedPool redisShardedPool = new RedisShardedPool(poolName);
                    redisPools.putIfAbsent(poolName, redisShardedPool);
                }
            }
        }
        return redisPools.get(poolName);
    }

    /**
     * <p> 默认初始化 </p>
     *
     * @since 1.0.0
     */
    public void initialize() {
        // 省略。。。。。。
    }

    /**
     * <p> 初始化 (非默认连接池) </p>
     *
     * @param cacheServerList 缓存服务器集
     * @since 1.0.0
     */
    public void initialize(List<CacheServer> cacheServerList) {
        // 省略。。。。。。
    }

    /**
     * <p> 获取当前连接池名 </p>
     *
     * @return 连接池名
     * @since 1.0.0
     */
    public String getPoolName() {
        return poolName;
    }

    /**
     * <p> 分布式Jedis连接池 </p>
     *
     * @return 分布式Jedis连接池
     * @since 1.0.0
     */
    public ShardedJedisPool getJedisPool() {
        if (ObjectUtils.isEmpty(shardedJedisPool)) {
            throw new FleaCacheConfigException("获取分布式Jedis连接池失败：请先调用initialize初始化");
        }
        return shardedJedisPool;
    }
}
```
## 3.7 Redis配置文件
**flea-cache**读取 [redis.properties](https://github.com/Huazie/flea-framework/blob/dev/flea-config/src/main/resources/flea/cache/redis.properties)（**Redis**配置文件），用作初始化 `RedisShardedPool`。

```bash
# Redis配置
redis.switch=0

redis.systemName=FleaFrame

redis.server=127.0.0.1:10001,127.0.0.1:10002,127.0.0.1:10003

redis.password=huazie123,huazie123,huazie123

redis.weight=1,1,1

redis.connectionTimeout=2000

redis.soTimeout=2000

redis.hashingAlg=1

# Redis客户端连接池配置
redis.pool.maxTotal=100

redis.pool.maxIdle=10

redis.pool.minIdle=0

redis.pool.maxWaitMillis=2000

redis.maxAttempts=5

redis.nullCacheExpiry=10
```

- `redis.switch` : **Redis**分片配置开关（1：开启 0：关闭），如果不配置也默认开启
- `redis.systemName` : **Redis**缓存所属系统名
- `redis.server` : **Redis**服务器地址
- `redis.password` : **Redis**服务登录密码
- `redis.weight` : **Redis**服务器权重分配
- `redis.connectionTimeout` : **Redis**客户端**socket**连接超时时间（单位：ms）
- `redis.soTimeout` : **Redis**客户端**socket**读写超时时间（单位：ms）
- `redis.hashingAlg` : **Redis**分布式hash算法【1 : MURMUR_HASH 2 : MD5】
- `redis.pool.maxTotal` : **Jedis**连接池最大连接数
- `redis.pool.maxIdle` : **Jedis**连接池最大空闲连接数
- `redis.pool.minIdle` : **Jedis**连接池最小空闲连接数
- `redis.pool.maxWaitMillis` : **Jedis**连接池获取连接时的最大等待时间（单位：ms）
- `redis.maxAttempts` : **Redis**客户端操作最大尝试次数【包含第一次操作】
- `redis.nullCacheExpiry` : 空缓存数据有效期（单位：s）

## 3.8 定义Redis Flea缓存类
[RedisFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/RedisFleaCache.java) 继承抽象**Flea**缓存类 `AbstractFleaCache` ，实现了以**Flea**框架操作**Redis**缓存的基本操作方法，其构造方法可见如需要传入**Redis**客户端 `RedisClient` ，相关使用下面介绍。

在上述基本操作方法中，实际使用**Redis**客户端【`redisClient`】 读、写和删除**Redis**缓存。其中写缓存方法【`putNativeValue`】在 添加的数据值为【`null`】时，默认添加空缓存数据【`NullCache`】 到Redis中，有效期取初始化参数【`nullCacheExpiry`】。

- 单个缓存接入场景，有效期配置可查看【`redis.properties`】中的配置参数 【`redis.nullCacheExpiry`】

- 整合缓存接入场景，有效期配置可查看【`flea-cache-config.xml`】中的缓存参数 【`<cache-param key="fleacore.nullCacheExpiry" desc="空缓存数据有效期（单位：s）">300</cache-param>`】

```java
public class RedisFleaCache extends AbstractFleaCache {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(RedisFleaCache.class);

    private RedisClient redisClient; // Redis客户端

    private CacheModeEnum cacheMode; // 缓存模式【分片模式和集群模式】

    /**
     * <p> 带参数的构造方法，初始化Redis Flea缓存类 </p>
     *
     * @param name            缓存数据主关键字
     * @param expiry          缓存数据有效期（单位：s）
     * @param nullCacheExpiry 空缓存数据有效期（单位：s）
     * @param cacheMode       缓存模式【分分片模式和集群模式】
     * @param redisClient     Redis客户端
     * @since 1.0.0
     */
    public RedisFleaCache(String name, int expiry, int nullCacheExpiry, CacheModeEnum cacheMode, RedisClient redisClient) {
        super(name, expiry, nullCacheExpiry);
        this.cacheMode = cacheMode;
        this.redisClient = redisClient;
        if (CacheUtils.isClusterMode(cacheMode))
            cache = CacheEnum.RedisCluster; // 缓存实现之Redis集群模式
        else
            cache = CacheEnum.RedisSharded; // 缓存实现之Redis分片模式
    }

    @Override
    public Object getNativeValue(String key) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(new Object() {}, "KEY = {}", key);
        }
        return redisClient.get(key);
    }

    @Override
    public Object putNativeValue(String key, Object value, int expiry) {
        if (LOGGER.isDebugEnabled()) {
            Object obj = new Object() {};
            LOGGER.debug1(obj, "REDIS FLEA CACHE, KEY = {}", key);
            LOGGER.debug1(obj, "REDIS FLEA CACHE, VALUE = {}", value);
            LOGGER.debug1(obj, "REDIS FLEA CACHE, EXPIRY = {}s", expiry);
            LOGGER.debug1(obj, "REDIS FLEA CACHE, NULL CACHE EXPIRY = {}s", getNullCacheExpiry());
        }
        if (ObjectUtils.isEmpty(value)) {
            return redisClient.set(key, new NullCache(key), getNullCacheExpiry());
        } else {
            if (expiry == CommonConstants.NumeralConstants.INT_ZERO) {
                return redisClient.set(key, value);
            } else {
                return redisClient.set(key, value, expiry);
            }
        }
    }

    @Override
    public Object deleteNativeValue(String key) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(new Object() {}, "KEY = {}", key);
        }
        return redisClient.del(key);
    }

    @Override
    public String getSystemName() {
        if (CacheUtils.isClusterMode(cacheMode))
            // 集群模式下获取缓存归属系统名
            return RedisClusterConfig.getConfig().getSystemName();
        else
            // 分片模式下获取缓存归属系统名
            return RedisShardedConfig.getConfig().getSystemName();
    }
}
```

## 3.9 定义抽象Flea缓存管理类
[AbstractFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCacheManager.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 3.10 定义Redis分片模式Flea缓存管理类
[RedisShardedFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/manager/RedisShardedFleaCacheManager.java)  继承抽象**Flea**缓存管理类 `AbstractFleaCacheManager`，用于接入**Flea**框架管理**Redis**缓存。

它的默认构造方法，用于初始化分片模式下默认连接池的**Redis**客户端, 这里需要先初始化**Redis**连接池，默认连接池名为【`default`】； 然后通过 `RedisClientFactory` 获取分片模式下默认连接池的Redis客户端 `RedisClient`，可在 **3.11** 查看。

`newCache` 用于创建一个**Redis Flea**缓存， 它里面包含了 读、写、删除 和 清空 缓存的基本操作。 该方法返回的是 `RedisFleaCache` 的实例对象，每一类 **Redis** 缓存数据都对应了一个  `RedisFleaCache` 的实例对象。

```java
public class RedisShardedFleaCacheManager extends AbstractFleaCacheManager {

    private RedisClient redisClient; // Redis客户端

    /**
     * <p> 默认构造方法，初始化分片模式下默认连接池的Redis客户端 </p>
     *
     * @since 1.0.0
     */
    public RedisShardedFleaCacheManager() {
        // 初始化默认连接池
        RedisShardedPool.getInstance().initialize();
        // 获取分片模式下默认连接池的Redis客户端
        redisClient = RedisClientFactory.getInstance();
    }

    @Override
    protected AbstractFleaCache newCache(String name, int expiry) {
        int nullCacheExpiry = RedisShardedConfig.getConfig().getNullCacheExpiry();
        return new RedisFleaCache(name, expiry, nullCacheExpiry, CacheModeEnum.SHARDED, redisClient);
    }
}
```

## 3.11 定义Redis客户端工厂类
[RedisClientFactory](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClientFactory.java)  ，有四种方式获取 **Redis** 客户端：
 - 一是获取分片模式下默认连接池的 **Redis** 客户端，应用在单个缓存接入场景【**3.10** 采用】；
 - 二是获取指定模式下默认连接池的 **Redis** 客户端，应用在单个缓存接入场景；
 - 三是获取分片模式下指定连接池的 **Redis** 客户端，应用在整合缓存接入场景；
 - 四是获取指定模式下指定连接池的 **Redis** 客户端，应用在整合缓存接入场景。

```java
public class RedisClientFactory {

    private static final ConcurrentMap<String, RedisClient> redisClients = new ConcurrentHashMap<>();

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
            synchronized (redisClients) {
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
在上面 的 `getInstance(String poolName, CacheModeEnum mode)` 方法中，使用了 `RedisClientStrategyContext` ，用于定义 **Redis** 客户端策略上下文。根据不同的缓存模式，就可以找到对应的 **Redis** 客户端策略。

## 3.12 定义 Redis 客户端策略上下文
[RedisClientStrategyContext](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/common/RedisClientStrategyContext.java) 包含了 **Redis** 分片 和 **Redis** 集群 相关的客户端策略。
```java
/**
 * Redis客户端策略上下文
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class RedisClientStrategyContext extends FleaStrategyContext<RedisClient, String> {

    public RedisClientStrategyContext() {
        super();
    }

    public RedisClientStrategyContext(String contextParam) {
        super(contextParam);
    }

    @Override
    protected Map<String, IFleaStrategy<RedisClient, String>> init() {
        Map<String, IFleaStrategy<RedisClient, String>> fleaStrategyMap = new HashMap<>();
        fleaStrategyMap.put(CacheModeEnum.SHARDED.name(), new RedisShardedClientStrategy());
        fleaStrategyMap.put(CacheModeEnum.CLUSTER.name(), new RedisClusterClientStrategy());
        return Collections.unmodifiableMap(fleaStrategyMap);
    }
}
```
## 3.13 定义分片模式 Redis 客户端策略
[RedisShardedClientStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/strategy/RedisShardedClientStrategy.java) 用于新建分片模式**Redis**客户端

```java
public class RedisShardedClientStrategy implements IFleaStrategy<RedisClient, String> {

    @Override
    public RedisClient execute(String poolName) throws FleaStrategyException {
        RedisClient originRedisClient;
        // 新建一个Flea Redis客户端类实例
        if (CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME.equals(poolName)) {
            originRedisClient = new FleaRedisShardedClient.Builder().build();
        } else {
            originRedisClient = new FleaRedisShardedClient.Builder(poolName).build();
        }
        return originRedisClient;
    }
}
```
好了，到这里我们可以来测试 **Redis** 分片模式。

## 3.14 Redis接入自测
单元测试类详见 [FleaCacheTest](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/test/java/com/huazie/fleaframework/cache/FleaCacheTest.java)

首先，这里需要按照 **Redis** 配置文件中的地址部署相应的 **Redis** 服务，可参考笔者的 [这篇博文](../../../../../../2019/08/30/flea-framework/flea-cache/flea-cache-windows-more-services/)。
```java
    @Test
    public void testRedisShardedFleaCache() {
        try {
            // 分片模式下Flea缓存管理类
            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(CacheEnum.RedisSharded.getName());
            AbstractFleaCache cache = manager.getCache("fleaparadetail");
            LOGGER.debug("Cache={}", cache);
            //## 1.  简单字符串
            cache.put("menu1", "huazie");
            cache.put("menu2", null);
//            cache.get("menu1");
//            cache.get("menu2");
//            cache.delete("menu1");
//            cache.delete("menu2");
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
[RedisSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/RedisSpringCache.java) 继承抽象 **Spring** 缓存 `AbstractSpringCache` 的读、写、删除 和 清空 缓存的基本操作方法，由**Redis Spring**缓存管理类初始化，用于对接 **Spring**； 

它的构造方法中，必须传入一个具体**Flea**缓存实现类，这里我们使用 **Redis Flea**缓存【`RedisFleaCache`】。

```java
public class RedisSpringCache extends AbstractSpringCache {

    /**
     * <p> 带参数的构造方法，初始化Redis Spring缓存类 </p>
     *
     * @param name      缓存数据主关键字
     * @param fleaCache 具体Flea缓存实现
     * @since 1.0.0
     */
    public RedisSpringCache(String name, IFleaCache fleaCache) {
        super(name, fleaCache);
    }

    /**
     * <p> 带参数的构造方法，初始化Redis Spring缓存类 </p>
     *
     * @param name            缓存数据主关键字
     * @param expiry          缓存数据有效期（单位：s）
     * @param nullCacheExpiry 空缓存数据有效期（单位：s）
     * @param cacheMode       缓存模式【分片模式和集群模式】
     * @param redisClient     Redis客户端
     * @since 1.0.0
     */
    public RedisSpringCache(String name, int expiry, int nullCacheExpiry, CacheModeEnum cacheMode, RedisClient redisClient) {
        this(name, new RedisFleaCache(name, expiry, nullCacheExpiry, cacheMode, redisClient));
    }

}
```
## 4.3 定义抽象Spring缓存管理类 

[AbstractSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCacheManager.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 4.4 定义Redis分片模式Spring缓存管理类
[RedisShardedSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/manager/RedisShardedSpringCacheManager.java) 继承抽象 **Spring** 缓存管理类 `AbstractSpringCacheManager`，用于接入**Spring**框架管理**Redis**缓存; 

它的默认构造方法，用于初始化分片模式下默认连接池的**Redis**客户端, 这里需要先初始化**Redis**连接池，默认连接池名为【`default`】； 然后通过**Redis**客户端工厂类来获取**Redis**客户端。

它的基本实现同 `RedisShardedFleaCacheManager`，唯一不同在于 `newCache` 的实现。方法【`newCache`】用于创建一个**Redis Spring**缓存， 而它内部是由**Redis Flea**缓存实现具体的 读、写、删除 和 清空 缓存的基本操作。

```java
public class RedisShardedSpringCacheManager extends AbstractSpringCacheManager {

    private RedisClient redisClient; // Redis客户端

    /**
     * <p> 默认构造方法，初始化分片模式下默认连接池的Redis客户端 </p>
     *
     * @since 1.0.0
     */
    public RedisShardedSpringCacheManager() {
        // 初始化默认连接池
        RedisShardedPool.getInstance().initialize();
        // 获取分片模式下默认连接池的Redis客户端
        redisClient = RedisClientFactory.getInstance();
    }

    @Override
    protected AbstractSpringCache newCache(String name, int expiry) {
        int nullCacheExpiry = RedisShardedConfig.getConfig().getNullCacheExpiry();
        return new RedisSpringCache(name, expiry, nullCacheExpiry, CacheModeEnum.SHARDED, redisClient);
    }
}
```
## 4.5 spring 配置

如下用于配置缓存管理 `redisShardedSpringCacheManager`，其中 `configMap` 为缓存时间(`key`缓存对象名称 `value`缓存过期时间)

```xml
    <bean id="redisShardedSpringCacheManager" class="com.huazie.frame.cache.redis.manager.RedisShardedSpringCacheManager">
        <property name="configMap">
            <map>
                <entry key="fleaconfigdata" value="86400"/>
            </map>
        </property>
    </bean>

    <!-- 开启缓存 -->
    <cache:annotation-driven cache-manager="redisShardedSpringCacheManager" proxy-target-class="true"/>
```
## 4.6 缓存自测

```java
    private ApplicationContext applicationContext;

    @Before
    public void init() {
        applicationContext = new ClassPathXmlApplicationContext("applicationContext.xml");
        LOGGER.debug("ApplicationContext={}", applicationContext);
    }

    @Test
    public void testRedisShardedSpringCache() {
        try {
            // 分片模式下Spring缓存管理类
            AbstractSpringCacheManager manager = (RedisShardedSpringCacheManager) applicationContext.getBean("redisShardedSpringCacheManager");
            LOGGER.debug("RedisCacheManager={}", manager);

            AbstractSpringCache cache = manager.getCache("fleaconfigdata");
            LOGGER.debug("Cache={}", cache);

            Set<String> cacheKey = cache.getCacheKey();
            LOGGER.debug("CacheKey = {}", cacheKey);

            //## 1.  简单字符串
//			cache.put("menu1", "huazie");
//            cache.get("menu1");
//            cache.get("menu1", String.class);

            //## 2.  简单对象(要是可以序列化的对象)
//			String user = new String("huazie");
//			cache.put("user", user);
//			LOGGER.debug(cache.get("user", String.class));
//            cache.get("FLEA_RES_STATE");
//            cache.clear();

            //## 3.  List塞对象
//			List<String> userList = new ArrayList<>();
//			userList.add("huazie");
//			userList.add("lgh");
//			cache.put("user_list", userList);

//			LOGGER.debug(cache.get("user_list",userList.getClass()).toString());

        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
```

# 结语
**Redis** 接入重构工作已经全部结束，当前版本为 **Redis** 分片模式。下一篇博文，我将要介绍 [Redis 集群模式的接入](../../../../../../2021/11/25/flea-framework/flea-cache/flea-cache-rediscluster/)工作，敬请期待！！！