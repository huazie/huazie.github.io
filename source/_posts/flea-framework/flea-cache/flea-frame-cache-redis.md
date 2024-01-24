---
title: flea-frame-cache使用之Redis接入【旧】
date: 2019-08-19 22:50:53
updated: 2023-07-11 11:34:14
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - flea-framework
  - flea-cache
  - Redis
  - Redis分片模式
---

[《开发框架-Flea》](/categories/开发框架-Flea/)

![](/images/cache.png)

# 1. 参考
[flea-frame-cache使用之Redis接入 源代码v1.0.0](https://github.com/Huazie/flea-frame/releases/tag/v1.0.0)

![](flea-frame-cache-redis.png)


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
## 3.1 定义Flea缓存接口 --- IFleaCache
可参考笔者的这篇博文 [Memcached接入](/2019/08/18/flea-cache-memcached/)，不再赘述。

## 3.2 定义抽象Flea缓存类 --- AbstractFleaCache
可参考笔者的这篇博文 [Memcached接入](/2019/08/18/flea-cache-memcached/)，不再赘述。
## 3.3 定义Redis客户端接口类 --- RedisClient
```java
/**
 * <p> Redis客户端对外接口 </p>
 *
 * @author huazie
 */
public interface RedisClient {

    /**
     * <p> 往Redis塞数据 </p>
     *
     * @param key   数据键
     * @param value 数据值
     * @return 状态码 （OK ：成功）
     */
    String set(final String key, final String value);

    /**
     * <p> 往Redis赛数据（用于序列化对象） </p>
     *
     * @param key   数据键
     * @param value 数据值
     * @return 状态码 （OK ：成功）
     */
    String set(final byte[] key, final byte[] value);

    /**
     * <p> 往Redis塞数据 (可以带失效时间) </p>
     * <p> 注意 ： (单位：s)</p>
     *
     * @param key    数据键
     * @param value  数据值
     * @param expiry 失效时间（单位：s）
     * @return 状态码 （OK ：成功）
     */
    String set(final String key, final String value, final int expiry);

    /**
     * <p> 往Redis塞数据 (可以带失效时间，用于序列化对象) </p>
     * <p> 注意 ： (单位：s)</p>
     *
     * @param key    数据键
     * @param value  数据值
     * @param expiry 失效时间（单位：s）
     * @return 状态码 （OK ：成功）
     */
    String set(final byte[] key, final byte[] value, final int expiry);

    /**
     * <p> 往Redis塞数据 (可以带失效时间) </p>
     * <p> 注意：（单位：ms） </p>
     *
     * @param key    数据键
     * @param value  数据值
     * @param expiry 失效时间（单位：ms）
     * @return 状态码 （OK ：成功）
     */
    String set(final String key, final String value, final long expiry);

    /**
     * <p> 往Redis塞数据 (可以带失效时间，用于序列化对象) </p>
     * <p> 注意：（单位：ms） </p>
     *
     * @param key    数据键
     * @param value  数据值
     * @param expiry 失效时间（单位：ms）
     * @return 状态码 （OK ：成功）
     */
    String set(final byte[] key, final byte[] value, final long expiry);

    /**
     * <p> 往Redis塞数据 (带参数) </p>
     *
     * @param key    数据键
     * @param value  数据值
     * @param params 参数
     * @return 状态码 （OK ：成功）
     */
    String set(final String key, final String value, SetParams params);

    /**
     * <p> 往Redis塞数据 (带参数，用于序列化对象) </p>
     *
     * @param key    数据键
     * @param value  数据值
     * @param params 参数
     * @return 状态码 （OK ：成功）
     */
    String set(final byte[] key, final byte[] value, SetParams params);

    /**
     * <p> 从Redis取数据 </p>
     *
     * @param key 数据键
     * @return 数据值
     */
    String get(final String key);

    /**
     * <p> 从Redis取数据（用于获取序列化对象） </p>
     *
     * @param key 数据键
     * @return 数据值
     */
    byte[] get(final byte[] key);

    /**
     * <p> 从Redis中删除数据 </p>
     *
     * @param key 数据键
     * @return 被删除key的数量
     */
    Long del(final String key);

    /**
     * <p> 获取数据所在的Redis服务器ip(主机地址+端口) </p>
     *
     * @param key 数据键
     * @return 当前数据所在的Redis服务器ip
     */
    String getLocation(final String key);

    /**
     * <p> 获取数据所在的Redis服务器ip(主机地址+端口) </p>
     *
     * @param key 数据键(字节数组)
     * @return 当前数据所在的Redis服务器ip
     */
    String getLocation(final byte[] key);

    /**
     * <p> 获取数据所在的Redis服务器主机 </p>
     *
     * @param key 数据键
     * @return 数据所在的Redis服务器主机
     */
    String getHost(final String key);

    /**
     * <p> 获取数据所在的Redis服务器主机 </p>
     *
     * @param key 数据键(字节数组)
     * @return 数据所在的Redis服务器主机
     */
    String getHost(final byte[] key);

    /**
     * <p> 获取数据所在的Redis服务器主机端口 </p>
     *
     * @param key 数据键
     * @return 数据所在的Redis服务器主机端口
     */
    Integer getPort(final String key);

    /**
     * <p> 获取数据所在的Redis服务器主机端口 </p>
     *
     * @param key 数据键(字节数组)
     * @return 数据所在的Redis服务器主机端口
     */
    Integer getPort(final byte[] key);

    /**
     * <p> 获取数据所在的客户端类 </p>
     *
     * @param key 数据键
     * @return 数据所在的客户端类
     */
    Client getClient(final String key);

    /**
     * <p> 获取数据所在的客户端类 </p>
     *
     * @param key 数据键
     * @return 数据所在的客户端类
     */
    Client getClient(final byte[] key);

    /**
     * <p> 获取分布式Redis集群客户端连接池 </p>
     *
     * @return 分布式Redis集群客户端连接池
     */
    ShardedJedisPool getJedisPool();

    /**
     * <p> 设置分布式Redis集群客户端 </p>
     *
     * @param shardedJedis 分布式Redis集群客户端
     */
    void setShardedJedis(ShardedJedis shardedJedis);

    /**
     * <p> 获取分布式Redis集群客户端 </p>
     *
     * @return 分布是Redis集群客户端
     */
    ShardedJedis getShardedJedis();

    /**
     * <p> 获取连接池名 </p>
     *
     * @return 连接池名
     */
    String getPoolName();

    /**
     * <p> 设置连接池名 </p>
     *
     * @param poolName 连接池名
     */
    void setPoolName(String poolName);
}
```
## 3.4 定义Redis客户端实现类 --- FleaRedisClient
该类实现 **RedisClient** 接口， 其中分布式Jedis连接池 **ShardedJedisPool** 用于获取分布式Jedis对象 **ShardedJedis**， **ShardedJedis**可以自行根据初始化的算法，计算当前传入的数据键在某一台初始化的 **Redis** 服务器上，从而操作对数据的添加，查找，删除功能。
```java
/**
 * <p> Flea Redis客户端类 </p>
 *
 * @author huazie
 */
public class FleaRedisClient implements RedisClient {

    private ShardedJedisPool shardedJedisPool; // 分布式Jedis连接池

    private ShardedJedis shardedJedis; // 分布式Jedis对象

    private String poolName; // 连接池名

    /**
     * <p> Redis客户端构造方法 (默认) </p>
     */
    private FleaRedisClient() {
        this(null);
    }

    /**
     * <p> Redis客户端构造方法（指定连接池名）</p>
     *
     * @param poolName 连接池名
     */
    private FleaRedisClient(String poolName) {
        this.poolName = poolName;
        init();
    }

    /**
     * <p> 初始化分布式Jedis连接池 </p>
     */
    private void init() {
        if (StringUtils.isBlank(poolName)) {
            poolName = CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME;
            shardedJedisPool = RedisPool.getInstance().getJedisPool();
        } else {
            shardedJedisPool = RedisPool.getInstance(poolName).getJedisPool();
        }

    }

    @Override
    public String set(final String key, final String value) {
        return shardedJedis.set(key, value);
    }

    @Override
    public String set(byte[] key, byte[] value) {
        return shardedJedis.set(key, value);
    }

    @Override
    public String set(final String key, final String value, final int expiry) {
        return shardedJedis.setex(key, expiry, value);
    }

    @Override
    public String set(byte[] key, byte[] value, int expiry) {
        return shardedJedis.setex(key, expiry, value);
    }

    @Override
    public String set(String key, String value, long expiry) {
        return shardedJedis.psetex(key, expiry, value);
    }

    @Override
    public String set(byte[] key, byte[] value, long expiry) {
        return shardedJedis.psetex(key, expiry, value);
    }

    @Override
    public String set(final String key, final String value, final SetParams params) {
        return shardedJedis.set(key, value, params);
    }

    @Override
    public String set(byte[] key, byte[] value, SetParams params) {
        return shardedJedis.set(key, value, params);
    }

    @Override
    public String get(final String key) {
        return shardedJedis.get(key);
    }

    @Override
    public byte[] get(byte[] key) {
        return shardedJedis.get(key);
    }

    @Override
    public Long del(final String key) {
        return shardedJedis.del(key);
    }

    @Override
    public String getLocation(final String key) {
        return getLocationByKey(key);
    }

    @Override
    public String getLocation(byte[] key) {
        return getLocationByKey(key);
    }

    @Override
    public String getHost(final String key) {
        return getHostByKey(key);
    }

    @Override
    public String getHost(byte[] key) {
        return getHostByKey(key);
    }

    @Override
    public Integer getPort(final String key) {
        return getPortByKey(key);
    }

    @Override
    public Integer getPort(byte[] key) {
        return getPortByKey(key);
    }

    @Override
    public Client getClient(String key) {
        return getClientByKey(key);
    }

    @Override
    public Client getClient(byte[] key) {
        return getClientByKey(key);
    }

    /**
     * <p> 获取数据所在的Redis服务器ip(主机地址+端口) </p>
     *
     * @param key 数据键
     * @return 当前数据所在的Redis服务器ip
     */
    private String getLocationByKey(Object key) {
        StringBuilder location = new StringBuilder();
        Client client = getClientByKey(key);
        if (ObjectUtils.isNotEmpty(client)) {
            location.append(client.getHost()).append(CommonConstants.SymbolConstants.COLON).append(client.getPort());
        }
        return location.toString();
    }

    /**
     * <p> 获取数据所在的Redis服务器主机 </p>
     *
     * @param key 数据键
     * @return 数据所在的Redis服务器主机
     */
    private String getHostByKey(Object key) {
        Client client = getClientByKey(key);
        if (ObjectUtils.isNotEmpty(client)) {
            return client.getHost();
        }
        return null;
    }

    /**
     * <p> 获取数据所在的Redis服务器主机端口 </p>
     *
     * @param key 数据键
     * @return 数据所在的Redis服务器主机端口
     */
    private Integer getPortByKey(Object key) {
        Client client = getClientByKey(key);
        if (ObjectUtils.isNotEmpty(client)) {
            return client.getPort();
        }
        return null;
    }

    /**
     * <p> 获取客户端类 </p>
     *
     * @param key 数据键
     * @return 客户端类
     */
    private Client getClientByKey(Object key) {
        Client client = null;
        if (ObjectUtils.isNotEmpty(key)) {
            if (key instanceof String) {
                client = shardedJedis.getShard(key.toString()).getClient();
            } else if (key instanceof byte[]) {
                client = shardedJedis.getShard((byte[]) key).getClient();
            }
        }
        return client;
    }

    @Override
    public ShardedJedisPool getJedisPool() {
        return shardedJedisPool;
    }

    @Override
    public void setShardedJedis(ShardedJedis shardedJedis) {
        this.shardedJedis = shardedJedis;
    }

    @Override
    public ShardedJedis getShardedJedis() {
        return shardedJedis;
    }

    @Override
    public String getPoolName() {
        return poolName;
    }

    @Override
    public void setPoolName(String poolName) {
        this.poolName = poolName;
        init();
    }

    /**
     * <p> 内部建造者类 </p>
     */
    static class Builder {
       	// 省略。。。。。。
    }
}

```
该类的构造函数初始化逻辑，可以看出我们使用了 **RedisPool**， 下面来介绍一下。
## 3.5 定义Redis连接池 --- RedisPool
**RedisPool** 用于Redis相关配置信息的初始化，其中重点是获取分布式Jedis连接池 **ShardedJedisPool** ，该类其中一个构造方法如下：
```java
/**
 * @param poolConfig 连接池配置信息
 * @param shards Jedis分布式服务器列表
 * @param algo 分布式算法
 */
public ShardedJedisPool(final GenericObjectPoolConfig poolConfig, List<JedisShardInfo> shards,
      Hashing algo) 
```

```java
/**
 * <p>  Flea Redis 连接池 </p>
 *
 * @author huazie
 */
public class RedisPool {

    private static final ConcurrentMap<String, RedisPool> redisPools = new ConcurrentHashMap<>();

    private String poolName; // 连接池名

    private ShardedJedisPool shardedJedisPool; // 分布式Jedis连接池

    private RedisPool(String poolName) {
        this.poolName = poolName;
    }

    /**
     * <p> 获取Redis连接池实例 (指定连接池名)</p>
     *
     * @param poolName 连接池名
     * @return Redis连接池实例对象
     */
    public static RedisPool getInstance(String poolName) {
        if (!redisPools.containsKey(poolName)) {
            synchronized (redisPools) {
                if (!redisPools.containsKey(poolName)) {
                    RedisPool redisPool = new RedisPool(poolName);
                    redisPools.putIfAbsent(poolName, redisPool);
                }
            }
        }
        return redisPools.get(poolName);
    }

    /**
     * <p> 获取Redis连接池实例 (默认) </p>
     *
     * @return Redis连接池实例对象
     */
    public static RedisPool getInstance() {
        return getInstance(CommonConstants.FleaPoolConstants.DEFAULT_POOL_NAME);
    }

    /**
     * <p> 默认初始化 </p>
     */
    void initialize() {
        // ...省略初始化的代码
    }

    /**
     * <p> 初始化 (非默认连接池) </p>
     *
     * @param cacheServerList 缓存服务器集
     * @param cacheParams     缓存参数
     * @since 1.0.0
     */
    void initialize(List<CacheServer> cacheServerList, CacheParams cacheParams) {
        // ...省略初始化的代码
    }

    /**
     * <p> 分布式Redis集群客户端连接池 </p>
     *
     * @return 分布式Redis集群客户端连接池
     */
    public ShardedJedisPool getJedisPool() {
        if (ObjectUtils.isEmpty(shardedJedisPool)) {
            throw new RuntimeException("获取分布式Redis集群客户端连接池失败：请先调用initialize初始化");
        }
        return shardedJedisPool;
    }

    /**
     * <p> 获取当前连接池名 </p>
     *
     * @return 连接池名
     */
    public String getPoolName() {
        return poolName;
    }
}
```
## 3.6 Redis配置文件
**flea-frame-cache** 读取 **redis.properties**（**Redis** 配置文件），用作初始化 **RedisPool**

```bash
# Redis配置
# Redis缓存所属系统名
redis.systemName=FleaFrame
# Redis服务器地址
redis.server=127.0.0.1:10001,127.0.0.1:10002,127.0.0.1:10003

# Redis服务登录密码
redis.password=huazie123,huazie123,huazie123

# Redis服务器权重分配
redis.weight=1,1,1

# Redis客户端socket连接超时时间
redis.connectionTimeout=2000

# Redis客户端socket读写超时时间
redis.soTimeout=2000

# Redis分布式hash算法
# 1 : MURMUR_HASH
# 2 : MD5
redis.hashingAlg=1

# Redis客户端连接池配置
# Redis客户端Jedis连接池最大连接数
redis.pool.maxTotal=100

# Redis客户端Jedis连接池最大空闲连接数
redis.pool.maxIdle=10

# Redis客户端Jedis连接池最小空闲连接数
redis.pool.minIdle=0

# Redis客户端Jedis连接池获取连接时的最大等待毫秒数
redis.pool.maxWaitMillis=2000
```

## 3.7 定义Redis Flea缓存类 --- RedisFleaCache
该类继承抽象Flea缓存类 **AbstractFleaCache** ，其构造方法可见如需要传入Redis客户端 **RedisClient** ，相关使用下面介绍：
```java
/**
 * <p> Redis Flea缓存类 </p>
 *
 * @author huazie
 */
public class RedisFleaCache extends AbstractFleaCache {

    private RedisClient redisClient; // Redis客户端

    /**
     * <p> 带参数的构造方法，初始化Redis Flea缓存类 </p>
     *
     * @param name        缓存主关键字
     * @param expiry      失效时长
     * @param redisClient Redis客户端
     */
    public RedisFleaCache(String name, long expiry, RedisClient redisClient) {
        super(name, expiry);
        this.redisClient = redisClient;
        cache = CacheEnum.Redis;
    }

    @Override
    public Object getNativeValue(String key) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(new Object() {}, "KEY = {}", key);
        }
        // 反序列化
        return ObjectUtils.deserialize(redisClient.get(key.getBytes()));
    }

    @Override
    public void putNativeValue(String key, Object value, long expiry) {
        // 序列化
        if (ObjectUtils.isNotEmpty(value)) {
            byte[] valueBytes = ObjectUtils.serialize(value);
            if (expiry == CommonConstants.NumeralConstants.ZERO) {
                redisClient.set(key.getBytes(), valueBytes);
            } else {
                redisClient.set(key.getBytes(), valueBytes, (int) expiry);
            }
        }

    }

    @Override
    public void deleteNativeValue(String key) {
        redisClient.del(key);
    }

    @Override
    public String getSystemName() {
        return RedisConfig.getConfig().getSystemName();
    }
}
```
## 3.8 定义抽象Flea缓存管理类 --- AbstractFleaCacheManager
可参考笔者的这篇博文 [Memcached接入](/2019/08/18/flea-cache-memcached/)，不再赘述。
## 3.9 定义Redis Flea缓存管理类 --- RedisFleaCacheManager
该类继承抽象Flea缓存管理类 **AbstractFleaCacheManager**，构造方法使用了**Redis** 客户端代理类 **RedisClientProxy** 获取Redis客户端 **RedisClient**，可在 **3.10** 查看。**newCache** 方法返回的是 **RedisFleaCache** 的实例对象，每一类 **Redis** 缓存数据都对应了一个  **RedisFleaCache** 的实例对象。
```java
/**
 * <p> Redis Flea缓存管理类 </p>
 *
 * @author huazie
 */
public class RedisFleaCacheManager extends AbstractFleaCacheManager {

    private RedisClient redisClient;

    /**
     * <p> 默认构造方法，初始化Redis Flea缓存管理类 </p>
     */
    public RedisFleaCacheManager() {
        // 初始化默认连接池
        RedisPool.getInstance().initialize();
        redisClient = RedisClientProxy.getProxyInstance();
    }

    @Override
    protected AbstractFleaCache newCache(String name, long expiry) {
        return new RedisFleaCache(name, expiry, redisClient);
    }
}
```
## 3.10 定义Redis客户端代理类 --- RedisClientProxy
**Redis** 客户端代理类 **RedisClientProxy** 中 **getProxyInstance()** 返回 默认连接池的 **Redis** 客户端，**getProxyInstance(String poolName)** 返回 指定连接池的 **Redis** 客户端。它们返回的都是 **Redis** 客户端接口类 **RedisClient** ，实际代理的是 **Flea Redis** 客户端 **FleaRedisClient**。
```java
/**
 * <p> RedisClient代理类 </p>
 *
 * @author huazie
 */
public class RedisClientProxy extends FleaProxy<RedisClient> {

    private final static ConcurrentMap<String, RedisClient> redisClients = new ConcurrentHashMap<String, RedisClient>();

    /**
     * <p> 获取RedisClient代理类 (默认)</p>
     *
     * @return RedisClient代理类
     * @since 1.0.0
     */
    public static RedisClient getProxyInstance() {
        return getProxyInstance(CacheConstants.FleaCacheConstants.DEFAUTL_POOL_NAME);
    }

    /**
     * <p> 获取RedisClient代理类 (指定连接池名)</p>
     *
     * @param poolName 连接池名
     * @return RedisClient代理类
     */
    public static RedisClient getProxyInstance(String poolName) {
        if (!redisClients.containsKey(poolName)) {
            synchronized (redisClients) {
                if (!redisClients.containsKey(poolName)) {
                    // 新建一个Flea Redis客户端类， 用于被代理
                    RedisClient originRedisClient;
                    if(CacheConstants.FleaCacheConstants.DEFAUTL_POOL_NAME.equals(poolName)) {
                        originRedisClient = new FleaRedisClient();
                    } else {
                        originRedisClient = new FleaRedisClient(poolName);
                    }
                    RedisClient proxyRedisClient = newProxyInstance(originRedisClient.getClass().getClassLoader(), originRedisClient.getClass().getInterfaces(), new RedisClientInvocationHandler(originRedisClient));
                    redisClients.put(poolName, proxyRedisClient);
                }
            }
        }
        return redisClients.get(poolName);
    }
}
```
## 3.11 定义Redis客户端调用处理类 --- RedisClientInvocationHandler
该类在 **RedisClientProxy** 中被调用，用于添加 **Flea Redis** 客户端类相应方法被代理调用前后的自定义操作，包含了代理拦截器 RedisClientProxyInterceptor 和  FleaDebugProxyInterceptor ，异常代理拦截器 FleaErrorProxyInterceptor。

```java
/**
 * <p> Redis客户端调用处理类 </p>
 *
 * @author huazie
 */
public class RedisClientInvocationHandler extends FleaProxyHandler {

    private static List<IFleaProxyInterceptor> proxyInterceptors; // 代理拦截器列表

    private static IFleaExceptionProxyInterceptor exceptionProxyInterceptor; // 异常代理拦截器

    static {
        proxyInterceptors = new ArrayList<>();
        proxyInterceptors.add(new RedisClientProxyInterceptor());
        proxyInterceptors.add(new FleaDebugProxyInterceptor());
        exceptionProxyInterceptor = new FleaErrorProxyInterceptor();
    }

    /**
     * <p> 带参数的构造方法 </p>
     *
     * @param proxyObject 被代理对象实例
     * @since 1.0.0
     */
    public RedisClientInvocationHandler(Object proxyObject) {
        super(proxyObject, proxyInterceptors, exceptionProxyInterceptor);
    }
}
```
## 3.12 定义Redis客户端代理拦截器 --- RedisClientProxyInterceptor
**RedisClientProxyInterceptor** 主要实现代理前的分布式 **Jedis** 对象 **ShardedJedis** 的获取，即方法 **beforeHandle** ；代理后的分布式 **Jedis** 对象 **ShardedJedis** 的关闭，归还相关资源给分布式 **Jedis** 连接池 **ShardedJedisPool**，即方法 **afterHandle**。

```java
/**
 * Redis客户端代理拦截器实现类
 *
 * <p> 方法 {@code beforeHandle} 用于在Redis客户端调用指定方法前，从分布式Jedis连接池中
 * 获取分布式Jedis对象，并将其初始化给Redis客户端类中的分布式Jedis对象。
 *
 * <p> 方法 {@code afterHandle} 用于在Redis客户端调用指定方法后，关闭分布式Jedis对象，
 * 并将它归还给分布式Jedis连接池。
 *
 * @author huazie
 */
public class RedisClientProxyInterceptor implements IFleaProxyInterceptor {

    @Override
    public void beforeHandle(Object proxyObject, Method method, Object[] args) throws Exception {
        RedisClient redisClient = convertProxyObject(proxyObject);
        ShardedJedisPool jedisPool = redisClient.getJedisPool();
        if (ObjectUtils.isNotEmpty(jedisPool)) {
            redisClient.setShardedJedis(jedisPool.getResource());
        }
    }

    @Override
    public void afterHandle(Object proxyObject, Method method, Object[] args, Object result, boolean hasException) throws Exception {
        RedisClient redisClient = convertProxyObject(proxyObject);
        ShardedJedis shardedJedis = redisClient.getShardedJedis();
        if (ObjectUtils.isNotEmpty(shardedJedis)) {
            // 使用后，关闭连接
            shardedJedis.close();
        }
    }

    /**
     * <p> 转换代理对象 </p>
     *
     * @param proxyObject 被代理的对象实例
     * @return Redis客户端
     */
    private RedisClient convertProxyObject(Object proxyObject) throws Exception {
        if (!(proxyObject instanceof RedisClient)) {
            throw new Exception("The proxyObject must implement RedisClient interface");
        }
        return (RedisClient) proxyObject;
    }
}

```
哇，终于Redis接入差不多要完成了，下面一起开始启动单元测试吧
## 3.13 Redis接入自测 --- FleaCacheTest
首先，这里需要按照 **Redis** 配置文件中的地址部署相应的 **Redis** 服务，可参考笔者的 [这篇博文](https://blog.csdn.net/u012855229/article/details/100139652)。
```java
    @Test
    public void testRedisFleaCache() {
        try {
            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(CacheEnum.Redis.getName());
            AbstractFleaCache cache = manager.getCache("fleaparadetail");
            LOGGER.debug("Cache={}", cache);
            //## 1.  简单字符串
//            cache.put("menu1", "huazie");
            cache.get("menu1");
//            cache.delete("menu1");
            cache.getCacheKey();
            LOGGER.debug(cache.getCacheName() + ">>>" + cache.getCacheDesc());
        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
```
# 4. 进阶接入
## 4.1 定义抽象Spring缓存 --- AbstractSpringCache
可参考笔者的这篇博文 [Memcached接入](/2019/08/18/flea-cache-memcached/)，不再赘述。
## 4.2 定义Redis Spring缓存类 --- RedisSpringCache
该类继承抽象 **Spring** 缓存 **AbstractSpringCache**，用于对接 **Spring**； 从构造方法可见，该类初始化还是使用 **RedisFleaCache**。
```java
/**
 * <p> Redis Spring缓存类 </p>
 *
 * @author huazie
 */
public class RedisSpringCache extends AbstractSpringCache {

    /**
     * <p> 带参数的构造方法，初始化Redis Spring缓存类 </p>
     *
     * @param name      缓存主关键字
     * @param fleaCache 具体缓存实现
     */
    public RedisSpringCache(String name, IFleaCache fleaCache) {
        super(name, fleaCache);
    }

    /**
     * <p> 带参数的构造方法，初始化Redis Spring缓存类 </p>
     *
     * @param name        缓存主关键字
     * @param expiry      失效时长
     * @param redisClient Redis客户端
     */
    public RedisSpringCache(String name, long expiry, RedisClient redisClient) {
        this(name, new RedisFleaCache(name, expiry, redisClient));
    }

}
```
## 4.3 定义抽象Spring缓存管理类 --- AbstractSpringCacheManager
可参考笔者的这篇博文 [Memcached接入](/2019/08/18/flea-cache-memcached/)，不再赘述。
## 4.4 定义Redis Spring缓存管理类 --- RedisSpringCacheManager
该类继承抽象 **Spring** 缓存管理类 **AbstractSpringCacheManager**，用于对接**Spring**; 基本实现同 **RedisFleaCacheManager**，唯一不同在于 **newCache** 的实现。
```java
/**
 * <p> Redis Spring缓存管理类 </p>
 *
 * @author huazie
 */
public class RedisSpringCacheManager extends AbstractSpringCacheManager {

    private RedisClient redisClient;

    /**
     * <p> 默认构造方法，初始化Redis Spring缓存管理类 </p>
     */
    public RedisSpringCacheManager() {
        // 初始化默认连接池
        RedisPool.getInstance().initialize();
        redisClient = RedisClientProxy.getProxyInstance();
    }

    @Override
    protected AbstractSpringCache newCache(String name, long expiry) {
        return new RedisSpringCache(name, expiry, redisClient);
    }
}
```
## 4.5 spring 配置
```xml
	<!--
	    配置缓存管理RedisSpringCacheManager
	    配置缓存时间 configMap (key缓存对象名称 value缓存过期时间)
	-->
	<bean id="redisSpringCacheManager" class="com.huazie.frame.cache.redis.RedisSpringCacheManager">
	    <property name="configMap">
	        <map>
	            <entry key="fleaparadetail" value="86400"/>
	        </map>
	    </property>
	</bean>
	
	<!-- 开启缓存 -->
	<cache:annotation-driven cache-manager="redisSpringCacheManager" proxy-target-class="true"/>
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
    public void testRedisSpringCache() {
        try {
            AbstractSpringCacheManager manager = (RedisSpringCacheManager) applicationContext.getBean("redisSpringCacheManager");
            LOGGER.debug("RedisCacheManager={}", manager);

            AbstractSpringCache cache = manager.getCache("fleaparadetail");
            LOGGER.debug("Cache={}", cache);

            //#### 1.  简单字符串
//			cache.put("menu1", "huazie");
//            cache.get("menu1");
//            cache.get("menu1", String.class);

            //#### 2.  简单对象(要是可以序列化的对象)
//			String user = new String("huazie");
//			cache.put("user", user);
//			LOGGER.debug(cache.get("user", String.class));
//            cache.get("FLEA_RES_STATE");
            cache.clear();

            //#### 3.  List塞对象
//			List<String> userList = new ArrayList<String>();
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
好了, **Redis** 的接入工作已经全部完成了。到目前为止，不论是Memcached的接入还是 Redis的接入，都是单一的缓存接入，笔者后续将介绍如何整合Memcached和Redis接入，以应对日益复杂的业务需求。 敬请期待！！！

# 更新
这一版 **Redis接入** 已进行重构，可详见 [flea-cache使用之Redis分片模式接入](/2021/11/18/flea-cache-redissharded/)
