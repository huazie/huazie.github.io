---
title: flea-cache使用之Redis集群模式接入
date: 2021-11-25 09:32:00
updated: 2023-12-27 09:33:41
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - flea-framework
  - flea-cache
  - Redis
  - Redis集群模式
---

![](/images/cache.png)

# 1. 参考
[flea-cache使用之Redis集群模式接入 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-cache)

<!-- more -->

![](flea-cache-rediscluster.png)

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
[RedisClient](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClient.java) 定义了 读、写、删除 Redis缓存的基本操作方法

```java
/**
 * Redis客户端接口，定义了 读、写、删除 Redis缓存的基本操作方法。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public interface RedisClient {

    String set(final String key, final Object value);

    String set(final byte[] key, final byte[] value);

    String set(final String key, final Object value, final int expiry);

    String set(final byte[] key, final byte[] value, final int expiry);

    String set(final String key, final Object value, final long expiry);

    String set(final byte[] key, final byte[] value, final long expiry);

    String set(final String key, final Object value, final SetParams params);

    String set(final byte[] key, final byte[] value, final SetParams params);

    Object get(final String key);

    byte[] get(final byte[] key);

    Long del(final String key);

    String getLocation(final String key);
    
    String getLocation(final byte[] key);

    String getHost(final String key);

    String getHost(final byte[] key);

    Integer getPort(final String key);

    Integer getPort(final byte[] key);

    Client getClient(final String key);

    Client getClient(final byte[] key);

    String getPoolName();
}
```

## 3.4 定义集群模式Redis客户端实现类
[FleaRedisClusterClient](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/FleaRedisClusterClient.java) 即**Flea**集群模式**Redis**客户端实现，封装了**Flea框架**操作**Redis**缓存的基本操作。它内部具体操作**Redis**集群缓存的功能，由**Jedis**集群实例对象 `JedisCluster` 完成， 包含读、写、删除**Redis**缓存的基本操作方法。

集群模式下，单个缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = new FleaRedisClusterClient.Builder().build();
   // 执行读，写，删除等基本操作
   redisClient.set("key", "value"); 
```
集群模式下，整合缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = new FleaRedisClusterClient.Builder(poolName).build();
   // 执行读，写，删除等基本操作
   redisClient.set("key", "value");  
```

当然每次都新建**Redis**客户端显然不可取，我们可通过**Redis**客户端工厂获取**Redis**客户端。
集群模式下，单个缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = RedisClientFactory.getInstance(CacheModeEnum.CLUSTER);
```

集群模式下，整合缓存接入场景，可通过如下方式使用：
```java
   RedisClient redisClient = RedisClientFactory.getInstance(poolName, CacheModeEnum.CLUSTER); 
```

```java
public class FleaRedisClusterClient extends FleaRedisClient {

    private JedisCluster jedisCluster;

    /**
     * <p> Redis集群客户端构造方法 (默认) </p>
     *
     * @since 1.1.0
     */
    private FleaRedisClusterClient() {
        this(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * <p> Redis集群客户端构造方法（指定连接池名） </p>
     *
     * @param poolName 连接池名
     * @since 1.1.0
     */
    private FleaRedisClusterClient(String poolName) {
        super(poolName);
        init();
    }

    /**
     * <p> 初始化Jedis集群实例 </p>
     *
     * @since 1.1.0
     */
    private void init() {
        if (CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME.equals(getPoolName())) {
            jedisCluster = RedisClusterPool.getInstance().getJedisCluster();
        } else {
            jedisCluster = RedisClusterPool.getInstance(getPoolName()).getJedisCluster();
        }

    }

    @Override
    public String set(String key, Object value) {
        if (value instanceof String)
            return jedisCluster.set(key, (String) value);
        else
            return jedisCluster.set(SafeEncoder.encode(key), ObjectUtils.serialize(value));
    }

    // 省略。。。。。。

    /**
     * <p> 内部建造者类 </p>
     */
    public static class Builder {

        private String poolName; // 连接池名

        /**
         * <p> 默认构造器 </p>
         *
         * @since 1.1.0
         */
        public Builder() {
        }

        /**
         * <p> 指定连接池的构造器 </p>
         *
         * @param poolName 连接池名
         * @since 1.1.0
         */
        public Builder(String poolName) {
            this.poolName = poolName;
        }

        /**
         * <p> 构建Redis集群客户端对象 </p>
         *
         * @return Redis集群客户端
         * @since 1.1.0
         */
        public RedisClient build() {
            if (StringUtils.isBlank(poolName)) {
                return new FleaRedisClusterClient();
            } else {
                return new FleaRedisClusterClient(poolName);
            }
        }
    }
}
```
该类的构造函数初始化逻辑，可以看出我们使用了 `RedisClusterPool`， 下面来介绍一下。

## 3.5 定义Redis集群连接池
我们使用 [RedisClusterPool](https://github.com/Huazie/flea-framework/blob/main/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClusterPool.java) 来初始化 **Redis**集群相关配置信息，其中重点是获取**Jedis**集群实例对象 `JedisCluster` ，该类其中一个构造方法如下：

```java
public JedisCluster(Set<HostAndPort> jedisClusterNode, int connectionTimeout, int soTimeout,
          int maxAttempts, String password, String clientName, final GenericObjectPoolConfig poolConfig) {
    super(jedisClusterNode, connectionTimeout, soTimeout, maxAttempts, password, clientName, poolConfig);
}
```

**Redis**集群连接池，用于初始化**Jedis**集群实例。

针对单独缓存接入场景，采用默认连接池初始化的方式； 可参考如下：
```java
    // 初始化默认连接池
    RedisClusterPool.getInstance().initialize();
``` 
针对整合缓存接入场景，采用指定连接池初始化的方式； 可参考如下：
```java
    // 初始化指定连接池
    RedisClusterPool.getInstance(group).initialize(cacheServerList);
```

```java
public class RedisClusterPool {

    private static final ConcurrentMap<String, RedisClusterPool> redisClusterPools = new ConcurrentHashMap<>();

    private static final Object redisClusterPoolLock = new Object();

    private String poolName; // 连接池名

    private JedisCluster jedisCluster; // Jedis集群实例

    private RedisClusterPool(String poolName) {
        this.poolName = poolName;
    }

    /**
     * <p> 获取Redis集群连接池实例 (默认连接池) </p>
     *
     * @return Redis集群连接池实例对象
     * @since 1.1.0
     */
    public static RedisClusterPool getInstance() {
        return getInstance(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * <p> 获取Redis集群连接池实例 (指定连接池名)</p>
     *
     * @param poolName 连接池名
     * @return Redis集群连接池实例
     * @since 1.1.0
     */
    public static RedisClusterPool getInstance(String poolName) {
        if (!redisClusterPools.containsKey(poolName)) {
            synchronized (redisClusterPoolLock) {
                if (!redisClusterPools.containsKey(poolName)) {
                    RedisClusterPool redisClusterPool = new RedisClusterPool(poolName);
                    redisClusterPools.putIfAbsent(poolName, redisClusterPool);
                }
            }
        }
        return redisClusterPools.get(poolName);
    }

    /**
     * <p> 默认初始化 </p>
     *
     * @since 1.1.0
     */
    public void initialize() {
        // 省略。。。。。。
    }

    /**
     * <p> 初始化 (非默认连接池) </p>
     *
     * @param cacheServerList 缓存服务器集
     * @since 1.1.0
     */
    public void initialize(List<CacheServer> cacheServerList) {
        // 省略。。。。。。
    }

    // 省略。。。。。。

    /**
     * <p> 获取Jedis集群实例对象 </p>
     *
     * @return Jedis集群实例对象
     * @since 1.1.0
     */
    public JedisCluster getJedisCluster() {
        if (ObjectUtils.isEmpty(jedisCluster)) {
            throw new FleaCacheConfigException("获取Jedis集群实例对象失败：请先调用initialize初始化");
        }
        return jedisCluster;
    }
}

```

## 3.6 定义Redis集群配置文件
**flea-cache** 读取 [redis.cluster.properties](https://github.com/Huazie/flea-framework/blob/main/flea-config/src/main/resources/flea/cache/redis.cluster.properties)（Redis集群配置文件），用作初始化 **RedisClusterPool** 

```bash
# Redis集群配置
redis.cluster.switch=0

redis.systemName=FleaFrame

redis.cluster.server=127.0.0.1:20011,127.0.0.1:20012,127.0.0.1:20021,127.0.0.1:20022,127.0.0.1:20031,127.0.0.1:20032

redis.cluster.password=huazie123

redis.cluster.connectionTimeout=2000

redis.cluster.soTimeout=2000

# Redis集群客户端连接池配置
redis.pool.maxTotal=100

redis.pool.maxIdle=10

redis.pool.minIdle=0

redis.pool.maxWaitMillis=2000

redis.maxAttempts=5

redis.nullCacheExpiry=10

```

- `redis.cluster.switch` : **Redis**集群配置开关（1：开启 0：关闭），如果不配置也默认开启
- `redis.systemName` : **Redis**缓存所属系统名
- `redis.cluster.server` : **Redis**集群服务节点地址
- `redis.cluster.password` : **Redis**集群服务节点登录密码（集群各节点配置同一个）
- `redis.cluster.connectionTimeout` : **Redis**集群客户端**socket**连接超时时间（单位：ms）
- `redis.cluster.soTimeout` : **Redis**集群客户端**socket**读写超时时间（单位：ms）
- `redis.pool.maxTotal` : **Jedis**连接池最大连接数
- `redis.pool.maxIdle` : **Jedis**连接池最大空闲连接数
- `redis.pool.minIdle` : **Jedis**连接池最小空闲连接数
- `redis.pool.maxWaitMillis` : **Jedis**连接池获取连接时的最大等待时间（单位：ms）
- `redis.maxAttempts` : **Redis**客户端操作最大尝试次数【包含第一次操作】
- `redis.nullCacheExpiry` : 空缓存数据有效期（单位：s）

## 3.7 定义Redis Flea缓存类
[RedisFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/impl/RedisFleaCache.java) 可参考笔者的这篇博文 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，不再赘述。

## 3.8 定义抽象Flea缓存管理类 
[AbstractFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCacheManager.java) 可参考笔者的这篇博文 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/)，不再赘述。

## 3.9 定义Redis集群模式Flea缓存管理类
[RedisClusterFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/manager/RedisClusterFleaCacheManager.java)  继承抽象Flea缓存管理类 `AbstractFleaCacheManager`，用于接入**Flea**框架管理**Redis**缓存。

它的默认构造方法，用于初始化集群模式下默认连接池的**Redis**客户端, 这里需要先初始化**Redis**集群连接池，默认连接池名为【`default`】； 然后通过 `RedisClientFactory` 获取集群模式下默认连接池的Redis客户端 `RedisClient`，可在 **3.10** 查看。

方法 `newCache` 用于创建一个 `RedisFleaCache` 的实例对象，它里面包含了 读、写、删除 和 清空 缓存的基本操作，每一类 **Redis** 缓存数据都对应了一个 `RedisFleaCache` 的实例对象。

```java
public class RedisClusterFleaCacheManager extends AbstractFleaCacheManager {

    private RedisClient redisClient; // Redis客户端

    /**
     * <p> 默认构造方法，初始化集群模式下默认连接池的Redis客户端 </p>
     *
     * @since 1.1.0
     */
    public RedisClusterFleaCacheManager() {
        // 初始化默认连接池
        RedisClusterPool.getInstance().initialize();
        // 获取集群模式下默认连接池的Redis客户端
        redisClient = RedisClientFactory.getInstance(CacheModeEnum.CLUSTER);
    }

    @Override
    protected AbstractFleaCache newCache(String name, int expiry) {
        int nullCacheExpiry = RedisClusterConfig.getConfig().getNullCacheExpiry();
        return new RedisFleaCache(name, expiry, nullCacheExpiry, CacheModeEnum.CLUSTER, redisClient);
    }
}
```

## 3.10 定义Redis客户端工厂类 
[RedisClientFactory](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/RedisClientFactory.java) ，有四种方式获取 **Redis** 客户端：
 - 一是获取分片模式下默认连接池的 **Redis** 客户端，应用在单个缓存接入场景；
 - 二是获取指定模式下默认连接池的 **Redis** 客户端，应用在单个缓存接入场景【**3.9** 采用】；
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

## 3.11 定义Redis客户端策略上下文
[RedisClientStrategyContext](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/common/RedisClientStrategyContext.java)  可参考笔者的这篇博文 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，不再赘述。

## 3.12 定义集群模式Redis客户端策略
 [RedisClusterClientStrategy](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/strategy/RedisClusterClientStrategy.java) 用于新建一个 `Flea Redis` 集群客户端。

```java
/**
 * 集群模式Redis客户端 策略
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.1.0
 */
public class RedisClusterClientStrategy implements IFleaStrategy<RedisClient, String> {

    @Override
    public RedisClient execute(String poolName) throws FleaStrategyException {
        RedisClient originRedisClient;
        // 新建一个Flea Redis集群客户端类实例
        if (CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME.equals(poolName)) {
            originRedisClient = new FleaRedisClusterClient.Builder().build();
        } else {
            originRedisClient = new FleaRedisClusterClient.Builder(poolName).build();
        }
        return originRedisClient;
    }
}
```
好了，到这里我们可以来测试 Redis 集群模式。

## 3.13 Redis集群模式接入自测
单元测试类  [FleaCacheTest](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/test/java/com/huazie/fleaframework/cache/FleaCacheTest.java)

首先，这里需要按照 **Redis集群配置文件** 中的地址部署相应的 **Redis集群** 服务，后续有机会我再出一篇简单的Redis主从集群搭建博文。

```java
@Test
    public void testRedisClusterFleaCache() {
        try {
            // 集群模式下Flea缓存管理类
            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(CacheEnum.RedisCluster.getName());
            AbstractFleaCache cache = manager.getCache("fleamenufavorites");
            LOGGER.debug("Cache={}", cache);
            //## 1.  简单字符串
//            cache.put("menu1", "huazie");
//            cache.put("menu2", null);
//            cache.get("menu1");
//            cache.get("menu2");
//            cache.delete("menu1");
//            cache.delete("menu2");
            cache.clear();
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

## 4.4 定义Redis集群模式Spring缓存管理类
 [RedisClusterSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/manager/RedisClusterSpringCacheManager.java) 继承抽象 **Spring** 缓存管理类 `AbstractSpringCacheManager`，用于接入**Spring**框架管理**Redis**缓存; 基本实现同 **RedisClusterFleaCacheManager**，唯一不同在于 **newCache** 的实现。

它的默认构造方法，用于初始化集群模式下默认连接池的**Redis**客户端, 这里需要先初始化**Redis**连接池，默认连接池名为【`default`】； 然后通过**Redis**客户端工厂类来获取**Redis**客户端。

方法【`newCache`】用于创建一个**Redis Spring**缓存， 而它内部是由**Redis Flea**缓存实现具体的 读、写、删除 和 清空 缓存的基本操作。

```java
public class RedisClusterSpringCacheManager extends AbstractSpringCacheManager {

    private RedisClient redisClient; // Redis客户端

    /**
     * <p> 默认构造方法，初始化集群模式下默认连接池的Redis客户端 </p>
     *
     * @since 1.1.0
     */
    public RedisClusterSpringCacheManager() {
        // 初始化默认连接池
        RedisClusterPool.getInstance().initialize();
        // 获取集群模式下默认连接池的Redis客户端
        redisClient = RedisClientFactory.getInstance(CacheModeEnum.CLUSTER);
    }

    @Override
    protected AbstractSpringCache newCache(String name, int expiry) {
        int nullCacheExpiry = RedisClusterConfig.getConfig().getNullCacheExpiry();
        return new RedisSpringCache(name, expiry, nullCacheExpiry, CacheModeEnum.CLUSTER, redisClient);
    }

}
```

## 4.5 spring 配置

```xml
    <!--
        配置缓存管理 redisClusterSpringCacheManager
        配置缓存时间 configMap (key缓存对象名称 value缓存过期时间)
    -->
    <bean id="redisClusterSpringCacheManager" class="com.huazie.fleaframework.cache.redis.manager.RedisClusterSpringCacheManager">
        <property name="configMap">
            <map>
                <entry key="fleamenufavorites" value="100"/>
            </map>
        </property>
    </bean>

    <!-- 开启缓存 -->
    <cache:annotation-driven cache-manager="redisClusterSpringCacheManager" proxy-target-class="true"/>
    
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
    public void testRedisClusterSpringCache() {
        try {
            // 集群模式下Spring缓存管理类
            AbstractSpringCacheManager manager = (RedisClusterSpringCacheManager) applicationContext.getBean("redisClusterSpringCacheManager");
            AbstractSpringCache cache = manager.getCache("fleamenufavorites");
            LOGGER.debug("Cache = {}", cache);

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
哇哇哇，**Redis** 集群模式接入终于搞定。到目前为止，不论是**Memcached**的接入还是 **Redis分片模式**接入亦或是本篇，都是单一的缓存接入，笔者的 下一篇博文 将介绍如何 [整合Memcached和Redis接入](../../../../../../2019/08/23/flea-framework/flea-cache/flea-cache-corecache/)，以应对日益复杂的业务需求。 敬请期待！！！