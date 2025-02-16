---
title: flea-cache使用之整合Memcached和Redis接入
date: 2019-08-23 09:39:20
updated: 2023-12-28 14:25:16
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - flea-framework
  - flea-cache
  - Memcached
  - Redis
  - 整合缓存接入
---

![](/images/cache.png)

# 一、引言
**Huazie** 的 **flea-framework** 框架下的 **flea-cache**，我们已经介绍了有关 **Memcached 接入** 和 **Redis 接入**；那么本篇我们将要介绍如何 **整合接入 Memcached 和 Redis**。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 二、往期文章

<table>
	<tr>
		<td rowspan="4" align="left"> 
			<a href="/categories/开发框架-Flea/">flea-cache</a> 
		</td>
	</tr>
	<tr>
		<td align="left"> 
			<a href="/2019/08/18/flea-framework/flea-cache/flea-cache-memcached/">flea-cache使用之Memcached接入</a> 
		</td>
	</tr>
	<tr>
		<td align="left"> 
			<a href="/2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/">flea-cache使用之Redis分片模式接入</a> 
		</td>
	</tr>
	<tr>
		<td align="left"> 
			<a href="/2021/11/25/flea-framework/flea-cache/flea-cache-rediscluster/">flea-cache使用之Redis集群模式接入</a> 
		</td>
	</tr>
</table>

# 三、主要内容
## 1. 参考
[flea-cache使用之整合Memcached和Redis接入 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-cache)

![](flea-cache-corecache.png)

## 2. 依赖

[Memcached-Java-Client-3.0.2.jar](https://mvnrepository.com/artifact/com.whalin/Memcached-Java-Client)
```xml
<!-- Memcached相关 -->
<dependency>
	<groupId>com.whalin</groupId>
	<artifactId>Memcached-Java-Client</artifactId>
	<version>3.0.2</version>
</dependency>
```
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

## 3. 基础接入
经过上两篇博文的介绍，**Memcached** 和 **Redis** 相信很多笔友都能成功的接入应用系统了。随着业务的复杂度上升，缓存的应用场景不断增多，单独的对接一个缓存系统，已经无法满足业务发展要求。

本文着眼于整合多套缓存接入：一个缓存 **cache** 对应一个缓存数据 **cache-data**，一个缓存数据 **cache-data** 对应一个缓存组 **cache-group**，多个缓存服务器 **cache-server** 关联一个缓存组 **cache-group**，一个缓存组 **cache-group** 对应具体的缓存接入实现（目前支持 **Memcached** 和 **Redis**）。

下面且听我慢慢道来：
### 3.1 Flea缓存配置文件
Flea缓存配置文件 ( [flea-cache-config.xml](https://github.com/Huazie/flea-framework/blob/dev/flea-config/src/main/resources/flea/cache/flea-cache-config.xml))，用来整合 **Memcached** 和 **Redis** 的相关配置，包含了缓存数据，缓存组，缓存服务器，缓存参数以及其他缓存配置项。
```xml
<?xml version="1.0" encoding="UTF-8"?>

<flea-cache-config>

    <!-- 缓存初始化配置项集 -->
    <cache-items key="FleaCacheInit" desc="缓存初始化配置项">
        <cache-item key="systemName" desc="缓存所属系统名">FleaFrame</cache-item>
    </cache-items>

    <!-- Flea缓存建造者配置项集 -->
    <cache-items key="FleaCacheBuilder" desc="Flea缓存建造者实现">
        <cache-item key="MemCached" desc="MemCached的Flea缓存建造者实现">com.huazie.fleaframework.cache.memcached.builder.MemCachedFleaCacheBuilder</cache-item>
        <cache-item key="RedisSharded" desc="Redis分片模式下的Flea缓存建造者实现">com.huazie.fleaframework.cache.redis.builder.RedisShardedFleaCacheBuilder</cache-item>
        <cache-item key="RedisCluster" desc="Redis集群模式下的Flea缓存建造者实现">com.huazie.fleaframework.cache.redis.builder.RedisClusterFleaCacheBuilder</cache-item>
        <cache-item key="RedisSentinel" desc="Redis哨兵模式下的Flea缓存建造者实现">com.huazie.fleaframework.cache.redis.builder.RedisSentinelFleaCacheBuilder</cache-item>
    </cache-items>

    <!-- 缓存参数集 -->
    <cache-params>
        <!-- 通用缓存参数 -->
        <cache-param key="fleacore.nullCacheExpiry" desc="空缓存数据有效期（单位：s）">300</cache-param>
        <!-- Redis 缓存参数 -->
        <cache-param key="redis.switch" desc="Redis分片配置开关（1：开启 0：关闭），如果不配置也默认开启">0</cache-param>
        <cache-param key="redis.connectionTimeout" desc="Redis客户端socket连接超时时间（单位：ms）">2000</cache-param>
        <cache-param key="redis.soTimeout" desc="Redis客户端socket读写超时时间（单位：ms）">2000</cache-param>
        <cache-param key="redis.hashingAlg" desc="Redis分布式hash算法(1:MURMUR_HASH,2:MD5)">1</cache-param>
        <cache-param key="redis.pool.maxTotal" desc="Redis客户端Jedis连接池最大连接数">100</cache-param>
        <cache-param key="redis.pool.maxIdle" desc="Redis客户端Jedis连接池最大空闲连接数">10</cache-param>
        <cache-param key="redis.pool.minIdle" desc="Redis客户端Jedis连接池最小空闲连接数">0</cache-param>
        <cache-param key="redis.pool.maxWaitMillis" desc="Redis客户端Jedis连接池获取连接时的最大等待时间（单位：ms）">2000</cache-param>
        <cache-param key="redis.maxAttempts" desc="Redis客户端操作最大尝试次数【包含第一次操作】">5</cache-param>

        <!-- Redis Cluster 缓存参数-->
        <cache-param key="redis.cluster.switch" desc="Redis集群配置开关（1：开启 0：关闭），如果不配置也默认开启">1</cache-param>
        <cache-param key="redis.cluster.connectionTimeout" desc="Redis集群客户端socket连接超时时间（单位：ms）">2000</cache-param>
        <cache-param key="redis.cluster.soTimeout" desc="Redis集群客户端socket读写超时时间（单位：ms）">2000</cache-param>
        <!-- 可以不用配置，缓存服务器cache-server没有配置，默认使用这里的密码配置 -->
        <!--<cache-param key="redis.cluster.password" desc="Redis集群服务节点登录密码（集群各节点配置同一个）">huazie123</cache-param>-->

        <!-- Redis Sentinel 缓存参数-->
        <cache-param key="redis.sentinel.switch" desc="Redis哨兵配置开关（1：开启 0：关闭），如果不配置也默认开启">1</cache-param>
        <cache-param key="redis.sentinel.connectionTimeout" desc="Redis哨兵客户端socket连接超时时间（单位：ms）">2000</cache-param>
        <cache-param key="redis.sentinel.soTimeout" desc="Redis哨兵客户端socket读写超时时间（单位：ms）">2000</cache-param>
        <!-- 可以不用配置，缓存服务器cache-server没有配置，默认使用这里的密码配置 -->
        <!--<cache-param key="redis.sentinel.password" desc="Redis主从服务器节点登录密码（各节点配置同一个）">huazie123</cache-param>-->

        <!-- MemCached缓存参数 -->
        <cache-param key="memcached.switch" desc="MemCached配置开关（1：开启 0：关闭），如果不配置也默认开启">1</cache-param>
        <cache-param key="memcached.initConn" desc="初始化时对每个服务器建立的连接数目">20</cache-param>
        <cache-param key="memcached.minConn" desc="每个服务器建立最小的连接数">20</cache-param>
        <cache-param key="memcached.maxConn" desc="每个服务器建立最大的连接数">500</cache-param>
        <cache-param key="memcached.maintSleep" desc="自查线程周期进行工作，其每次休眠时间（单位：ms）">60000</cache-param>
        <cache-param key="memcached.nagle" desc="Socket的参数，如果是true在写数据时不缓冲，立即发送出去">true</cache-param>
        <cache-param key="memcached.socketTO" desc="Socket阻塞读取数据的超时时间（单位：ms）">3000</cache-param>
        <cache-param key="memcached.socketConnectTO" desc="Socket连接超时时间（单位：ms）">3000</cache-param>
        <!--
            0 - native String.hashCode();
            1 - original compatibility
            2 - new CRC32 based
            3 - MD5 Based
        -->
        <cache-param key="memcached.hashingAlg" desc="MemCached分布式hash算法">3</cache-param>
    </cache-params>

    <!-- Flea缓存数据集 -->
    <cache-datas>
        <!-- type="缓存数据类型" 对应 flea-cache.xml 中 <cache type="缓存数据类型" />  -->
        <cache-data type="fleaAuth" desc="Flea Auth缓存数据所在组配置">authGroup</cache-data>
        <cache-data type="fleaFrame" desc="Flea Frame配置数据所在组配置">configGroup</cache-data>
        <cache-data type="fleaJersey" desc="Flea Jersey配置数据所在组配置">jerseyGroup</cache-data>
        <cache-data type="fleaDynamic" desc="Flea 动态数据缓存所在组配置">dynamicGroup</cache-data>
    </cache-datas>

    <!-- Flea缓存组集 -->
    <cache-groups>
        <!-- group 对应 <cache-data>group</cache-date>  -->
        <!-- group 的缓存组关联缓存实现 MemCached 对应Flea缓存建造者实现 <cache-item key="MemCached"> -->
        <cache-group group="authGroup" desc="Flea权限数据缓存组">MemCached</cache-group>
        <!-- group 的缓存组关联缓存实现 RedisSharded 对应Flea缓存建造者实现 <cache-item key="RedisSharded"> -->
        <cache-group group="configGroup" desc="Flea配置数据缓存组">RedisSharded</cache-group>
        <!-- group 的缓存组关联缓存实现 RedisCluster 对应Flea缓存建造者实现 <cache-item key="RedisCluster"> -->
        <cache-group group="dynamicGroup" desc="Flea动态数据缓存组">RedisCluster</cache-group>
        <!-- group 的缓存组关联缓存实现 RedisSentinel 对应Flea缓存建造者实现 <cache-item key="RedisSentinel"> -->
        <cache-group group="jerseyGroup" desc="Flea动态数据缓存组">RedisSentinel</cache-group>
    </cache-groups>

    <!-- Flea缓存服务器集 -->
    <cache-servers>
        <cache-server group="authGroup" weight="1" desc="MemCached缓存服务器">127.0.0.1:31113</cache-server>
        <cache-server group="authGroup" weight="1" desc="MemCached缓存服务器">127.0.0.1:31114</cache-server>

        <cache-server group="configGroup" password="huazie123" weight="1" desc="Redis缓存服务器【分片模式】">127.0.0.1:10001</cache-server>
        <cache-server group="configGroup" password="huazie123" weight="1" desc="Redis缓存服务器【分片模式】">127.0.0.1:10002</cache-server>
        <cache-server group="configGroup" password="huazie123" weight="1" desc="Redis缓存服务器【分片模式】">127.0.0.1:10003</cache-server>

        <cache-server group="dynamicGroup" password="huazie123" desc="Redis缓存服务器【集群模式】">127.0.0.1:20011</cache-server>
        <cache-server group="dynamicGroup" password="huazie123" desc="Redis缓存服务器【集群模式】">127.0.0.1:20012</cache-server>
        <cache-server group="dynamicGroup" password="huazie123" desc="Redis缓存服务器【集群模式】">127.0.0.1:20021</cache-server>
        <cache-server group="dynamicGroup" password="huazie123" desc="Redis缓存服务器【集群模式】">127.0.0.1:20022</cache-server>
        <cache-server group="dynamicGroup" password="huazie123" desc="Redis缓存服务器【集群模式】">127.0.0.1:20031</cache-server>
        <cache-server group="dynamicGroup" password="huazie123" desc="Redis缓存服务器【集群模式】">127.0.0.1:20032</cache-server>

        <cache-server group="jerseyGroup" master="mymaster" database="0" desc="Redis哨兵服务器">127.0.0.1:36379</cache-server>
        <cache-server group="jerseyGroup" master="mymaster" database="0" desc="Redis哨兵服务器">127.0.0.1:36380</cache-server>
        <cache-server group="jerseyGroup" master="mymaster" database="0" desc="Redis哨兵服务器">127.0.0.1:36381</cache-server>
    </cache-servers>

</flea-cache-config>
```

### 3.2 Flea缓存定义文件 
Flea缓存定义文件（[flea-cache.xml](https://github.com/Huazie/flea-framework/blob/dev/flea-config/src/main/resources/flea/cache/flea-cache.xml)），用来定义各类缓存，其中 **key** 表示缓存主关键字， **type** 表示一类缓存数据，**expiry** 表示缓存生效时长（单位：秒【0：永久】）。

```xml
<?xml version="1.0" encoding="UTF-8"?>

<flea-cache>

    <caches>
        <!--
            key    : 缓存数据主关键字
            type   : 缓存数据类型，对应 flea-cache-config.xml中 <cache-data type="缓存数据类型">
            expiry : 缓存数据有效期
        -->
        <cache key="fleaconfigdata" type="fleaFrame" expiry="86400" desc="Flea配置数据缓存" />

        <cache key="fleajerseyi18nerrormapping" type="fleaJersey" expiry="86400" desc="Flea Jersey 国际码和错误码映射缓存" />
        <cache key="fleajerseyresservice" type="fleaJersey" expiry="86400" desc="Flea Jersey 资源服务缓存" />
        <cache key="fleajerseyresclient" type="fleaJersey" expiry="86400" desc="Flea Jersey 资源客户端缓存" />
        <cache key="fleajerseyresource" type="fleaJersey" expiry="86400" desc="Flea Jersey 资源缓存" />

        <cache key="fleamenufavorites" type="fleaDynamic" expiry="0" desc="Flea菜单收藏夹数据缓存" />

    </caches>

    <!-- 其他缓存定义配置文件引入 -->
    <!-- flea-auth 授权模块缓存定义配置文件引入 -->
    <import resource="flea/cache/flea-auth-cache.xml"/>

</flea-cache>
```

### 3.3 定义核心Flea缓存类
[CoreFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/core/impl/CoreFleaCache.java) 同样继承抽象Flea缓存 [AbstractFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCache.java)，实现其定义的抽象方法；内部定义成员变量 **fleaCache** 用于指定具体的 **Flea** 缓存实现（这个具体的实现，可参考 [Memcached接入](../../../../../../2019/08/18/flea-framework/flea-cache/flea-cache-memcached/) 和 [Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)），实现的三个方法 **getNativeValue**，**putNativeValue**，**deleteNativeValue** 内部采用具体Flea缓存实现fleaCache相应的方法实现读缓存、写缓存，删缓存；从构造方法可见，**fleaCache** 通过 **FleaCacheFactory.getFleaCache(name)** ，从Flea缓存工厂中获取。

```java
/**
 * 核心Flea缓存类，实现读、写和删除缓存的基本操作方法，用于整合各类缓存的接入。
 *
 * <p> 成员变量【{@code fleaCache}】，在构造方法中根据缓存数据主关键字
 * 指定具体的Flea缓存实现，在 读、写 和 删除缓存的基本操作方法中调用
 * 【{@code fleaCache}】的具体实现。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public class CoreFleaCache extends AbstractFleaCache {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(CoreFleaCache.class);

    private AbstractFleaCache fleaCache; // 指定Flea缓存实现

    /**
     * 初始化核心Flea缓存
     *
     * @param name 缓存数据主关键字
     * @since 1.0.0
     */
    public CoreFleaCache(String name) {
        super(name, CacheConfigUtils.getExpiry(name), CacheConfigUtils.getNullCacheExpiry());
        // 根据缓存数据主关键字name获取指定Flea缓存对象
        fleaCache = FleaCacheFactory.getFleaCache(name);
        // 取指定Flea缓存的缓存类型
        cache = fleaCache.getCache();
    }

    @Override
    public Object getNativeValue(String key) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(new Object() {}, "KEY = {}", key);
        }
        return fleaCache.getNativeValue(key);
    }

    @Override
    public Object putNativeValue(String key, Object value, int expiry) {
        if (LOGGER.isDebugEnabled()) {
            Object obj = new Object() {};
            LOGGER.debug1(obj, "CORE FLEA CACHE, KEY = {}", key);
            LOGGER.debug1(obj, "CORE FLEA CACHE, VALUE = {}", value);
            LOGGER.debug1(obj, "CORE FLEA CACHE, EXPIRY = {}s", expiry);
            LOGGER.debug1(obj, "CORE FLEA CACHE, NULL CACHE EXPIRY = {}s", getNullCacheExpiry());
        }
        return fleaCache.putNativeValue(key, value, expiry);
    }

    @Override
    public Object deleteNativeValue(String key) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(new Object() {}, "KEY = {}", key);
        }
        return fleaCache.deleteNativeValue(key);
    }

    @Override
    public String getSystemName() {
        return CacheConfigUtils.getSystemName();
    }
}
```

### 3.4 定义Flea缓存工厂类 
[FleaCacheFactory](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/common/FleaCacheFactory.java) 根据缓存名（即缓存主关键字）所创建的Flea缓存都存入 **ConcurrentMap<String, AbstractFleaCache>** 中。**newCache** 方法主要是根据缓存名查找相关缓存 **cache**，再找到缓存数据 **cache-data**，接着找到缓存组 **cache-group**，最后根据缓存组所属的缓存系统，查找缓存配置项 **cache-item**，获取对应Flea缓存的建造者实现，可参见 **flea-cache-config.xml** 配置。

```java
/**
 * Flea Cache 工厂类，用于整合各类缓存接入时，创建具体的缓存实现类。
 *
 * <p> 同步集合类 {@code fleaCacheMap}，存储的键为缓存数据主关键字，
 * 对应缓存定义配置文件【flea-cache.xml】中的【{@code <cache
 * key="缓存数据主关键字"></cache>}】；它的值为具体的缓存实现类。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public class FleaCacheFactory {

    private static final ConcurrentMap<String, AbstractFleaCache> fleaCacheMap = new ConcurrentHashMap<>();

    private static final Object fleaCacheMapLock = new Object();    

    private FleaCacheFactory() {
    }

    /**
     * <p> 根据缓存数据主关键字获取指定Flea缓存对象 </p>
     *
     * @param name 缓存数据主关键字（对应 flea-cache.xml {@code <cache key="缓存数据主关键字"></cache>}）
     * @return Flea缓存对象
     * @since 1.0.0
     */
    public static AbstractFleaCache getFleaCache(String name) {
        if (!fleaCacheMap.containsKey(name)) {
            synchronized (fleaCacheMapLock) {
                if (!fleaCacheMap.containsKey(name)) {
                    fleaCacheMap.put(name, newFleaCache(name));
                }
            }
        }
        return fleaCacheMap.get(name);
    }

    /**
     * <p> 根据缓存数据主关键字创建一个Flea缓存对象 </p>
     *
     * @param name 缓存数据主关键字（对应 flea-cache.xml {@code <cache key="缓存数据主关键字"></cache>}）
     * @return Flea缓存对象
     * @since 1.0.0
     */
    private static AbstractFleaCache newFleaCache(String name) {
        // 详见 GitHub 链接
    }
}
```

### 3.5 定义Flea缓存建造者接口类 

[IFleaCacheBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/IFleaCacheBuilder.java) 定义 `build` 方法，用于构建 `Flea` 缓存对象
```java
/**
 * Flea缓存建造者接口类，定义了构建Flea缓存对象的通用接口。
 *
 * <p> 该接口由Flea缓存工厂类【{@code FleaCacheFactory}】使用，
 * 根据不同缓存数据归属缓存组所在的缓存类型，读取缓存配置中的
 * Flea缓存建造者配置项集，通过反射实例化具体的Flea缓存。
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public interface IFleaCacheBuilder {

    /**
     * 构建Flea缓存对象
     *
     * @param name            缓存数据主关键字
     * @param cacheServerList 缓存服务器集
     * @return Flea缓存对象
     * @since 1.0.0
     */
    AbstractFleaCache build(String name, List<CacheServer> cacheServerList);
}
```

### 3.6 定义Memcached Flea缓存建造者
[MemCachedFleaCacheBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/memcached/builder/MemCachedFleaCacheBuilder.java) 实现 **IFleaCacheBuilder**，用于构建基于 **Memcached** 的 **Flea** 缓存，即创建一个 **MemCachedFleaCache**。

```java
/**
 * MemCached Flea缓存建造者实现类，用于整合各类缓存接入时创建MemCached Flea缓存。
 *
 * <p> 缓存定义文件【flea-cache.xml】中，每一个缓存定义配置都对应缓存配置文件
 * 【flea-cache-config.xml】中的一类缓存数据，每类缓存数据都归属一个缓存组，
 * 每个缓存组都映射着具体的缓存实现名，而整合各类缓存接入时，
 * 每个具体的缓存实现名都配置了Flea缓存建造着实现类。
 *
 * <p> 可查看Flea缓存配置文件【flea-cache-config.xml】，获取
 * MemCached Flea缓存建造者配置项【{@code <cache-item key="MemCached">}】
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public class MemCachedFleaCacheBuilder implements IFleaCacheBuilder {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(MemCachedFleaCacheBuilder.class);

    @Override
    public AbstractFleaCache build(String name, List<CacheServer> cacheServerList) {
        if (CollectionUtils.isEmpty(cacheServerList)) {
            throw new FleaCacheConfigException("无法初始化MemCached Flea缓存，缓存服务器列表【cacheServerList】为空");
        }
        // 获取缓存数据有效期（单位：s）
        int expiry = CacheConfigUtils.getExpiry(name);
        // 获取空缓存数据有效期（单位：s）
        int nullCacheExpiry = CacheConfigUtils.getNullCacheExpiry();
        // 获取MemCached服务器所在组名
        String group = cacheServerList.get(0).getGroup();
        // 通过组名来获取 MemCached客户端类
        MemCachedClient memCachedClient = new MemCachedClient(group);
        // 获取MemCachedPool，并初始化连接池
        MemCachedPool.getInstance(group).initialize(cacheServerList);
        // 创建一个MemCached Flea缓存类
        AbstractFleaCache fleaCache = new MemCachedFleaCache(name, expiry, nullCacheExpiry, memCachedClient);

        if (LOGGER.isDebugEnabled()) {
            Object obj = new Object() {};
            LOGGER.debug1(obj, "Pool Name = {}", MemCachedPool.getInstance(group).getPoolName());
            LOGGER.debug1(obj, "Pool = {}", MemCachedPool.getInstance(group).getSockIOPool());
        }

        return fleaCache;
    }
}
```

### 3.7 定义Redis分片模式Flea缓存建造者
[RedisShardedFleaCacheBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/builder/RedisShardedFleaCacheBuilder.java) 实现 **IFleaCacheBuilder**，用于构建基于Redis的Flea缓存，即创建一个 分片模式的 **RedisFleaCache**。

```java
/**
 * Redis分片模式Flea缓存建造者实现类，用于整合各类缓存接入时创建Redis Flea缓存。
 *
 * <p> 缓存定义文件【flea-cache.xml】中，每一个缓存定义配置都对应缓存配置文件
 * 【flea-cache-config.xml】中的一类缓存数据，每类缓存数据都归属一个缓存组，
 * 每个缓存组都映射着具体的缓存实现名，而整合各类缓存接入时，
 * 每个具体的缓存实现名都配置了Flea缓存建造着实现类。
 *
 * <p> 可查看Flea缓存配置文件【flea-cache-config.xml】，
 * 获取Redis Flea缓存建造者配置项【{@code <cache-item key="RedisSharded">}】
 *
 * @author huazie
 * @version 1.1.0
 * @see FleaCacheFactory
 * @since 1.0.0
 */
public class RedisShardedFleaCacheBuilder implements IFleaCacheBuilder {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(RedisShardedFleaCacheBuilder.class);

    @Override
    public AbstractFleaCache build(String name, List<CacheServer> cacheServerList) {
        if (CollectionUtils.isEmpty(cacheServerList)) {
            throw new FleaCacheConfigException("无法初始化分片模式下Redis Flea缓存，缓存服务器列表【cacheServerList】为空");
        }
        // 获取缓存数据有效期（单位：s）
        int expiry = CacheConfigUtils.getExpiry(name);
        // 获取空缓存数据有效期（单位：s）
        int nullCacheExpiry = CacheConfigUtils.getNullCacheExpiry();
        // 获取缓存组名
        String group = cacheServerList.get(0).getGroup();
        // 初始化指定连接池名【group】的Redis分片模式连接池
        RedisShardedPool.getInstance(group).initialize(cacheServerList);
        // 获取分片模式下的指定连接池名【group】的Redis客户端
        RedisClient redisClient = RedisClientFactory.getInstance(group);
        // 创建一个Redis Flea缓存
        AbstractFleaCache fleaCache = new RedisFleaCache(name, expiry, nullCacheExpiry, CacheModeEnum.SHARDED, redisClient);

        if (LOGGER.isDebugEnabled()) {
            Object obj = new Object() {};
            LOGGER.debug1(obj, "Pool Name = {}", RedisShardedPool.getInstance(group).getPoolName());
            LOGGER.debug1(obj, "Pool = {}", RedisShardedPool.getInstance(group).getJedisPool());
        }

        return fleaCache;
    }
}
```

### 3.8 定义Redis集群模式Flea缓存建造者
[RedisClusterFleaCacheBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/builder/RedisClusterFleaCacheBuilder.java) 实现 **IFleaCacheBuilder**，用于构建基于Redis的Flea缓存，即创建一个 集群模式的 **RedisFleaCache**。

```java
/**
 * Redis集群模式Flea缓存建造者实现类，用于整合各类缓存接入时创建Redis Flea缓存。
 *
 * <p> 缓存定义文件【flea-cache.xml】中，每一个缓存定义配置都对应缓存配置文件
 * 【flea-cache-config.xml】中的一类缓存数据，每类缓存数据都归属一个缓存组，
 * 每个缓存组都映射着具体的缓存实现名，而整合各类缓存接入时，
 * 每个具体的缓存实现名都配置了Flea缓存建造着实现类。
 *
 * <p> 可查看Flea缓存配置文件【flea-cache-config.xml】，
 * 获取Redis Flea缓存建造者配置项【{@code <cache-item key="RedisCluster">}】
 *
 * @author huazie
 * @version 1.1.0
 * @see FleaCacheFactory
 * @since 1.1.0
 */
public class RedisClusterFleaCacheBuilder implements IFleaCacheBuilder {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(RedisClusterFleaCacheBuilder.class);

    @Override
    public AbstractFleaCache build(String name, List<CacheServer> cacheServerList) {
        if (CollectionUtils.isEmpty(cacheServerList)) {
            throw new FleaCacheConfigException("无法初始化集群模式下Redis Flea缓存，缓存服务器列表【cacheServerList】为空");

        }
        // 获取缓存数据有效期（单位：s）
        int expiry = CacheConfigUtils.getExpiry(name);
        // 获取空缓存数据有效期（单位：s）
        int nullCacheExpiry = CacheConfigUtils.getNullCacheExpiry();
        // 获取缓存组名
        String group = cacheServerList.get(0).getGroup();
        // 初始化指定连接池名【group】的Redis集群模式连接池
        RedisClusterPool.getInstance(group).initialize(cacheServerList);
        // 获取集群模式下的指定连接池名【group】的Redis客户端类
        RedisClient redisClient = RedisClientFactory.getInstance(group, CacheModeEnum.CLUSTER);
        // 创建一个Redis Flea缓存
        AbstractFleaCache fleaCache = new RedisFleaCache(name, expiry, nullCacheExpiry, CacheModeEnum.CLUSTER, redisClient);

        if (LOGGER.isDebugEnabled()) {
            Object obj = new Object() {};
            LOGGER.debug1(obj, "Pool Name = {}", RedisClusterPool.getInstance(group).getPoolName());
            LOGGER.debug1(obj, "JedisCluster = {}", RedisClusterPool.getInstance(group).getJedisCluster());
        }

        return fleaCache;
    }
}
```

### 3.9 定义Redis哨兵模式Flea缓存建造者
[RedisSentinelFleaCacheBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/redis/builder/RedisSentinelFleaCacheBuilder.java) 实现 **IFleaCacheBuilder**，用于构建基于Redis的Flea缓存，即创建一个哨兵模式的 **RedisFleaCache**。

```java
/**
 * Redis哨兵模式Flea缓存建造者实现类，用于整合各类缓存接入时创建Redis Flea缓存。
 *
 * <p> 缓存定义文件【flea-cache.xml】中，每一个缓存定义配置都对应缓存配置文件
 * 【flea-cache-config.xml】中的一类缓存数据，每类缓存数据都归属一个缓存组，
 * 每个缓存组都映射着具体的缓存实现名，而整合各类缓存接入时，
 * 每个具体的缓存实现名都配置了Flea缓存建造者实现类。
 *
 * <p> 可查看Flea缓存配置文件【flea-cache-config.xml】，
 * 获取Redis Flea缓存建造者配置项【{@code <cache-item key="RedisSentinel">}】
 *
 * @author huazie
 * @version 2.0.0
 * @see FleaCacheFactory
 * @since 2.0.0
 */
public class RedisSentinelFleaCacheBuilder implements IFleaCacheBuilder {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(RedisSentinelFleaCacheBuilder.class);

    @Override
    public AbstractFleaCache build(String name, List<CacheServer> cacheServerList) {
        if (CollectionUtils.isEmpty(cacheServerList)) {
            ExceptionUtils.throwFleaException(FleaCacheConfigException.class, "无法初始化哨兵模式下Redis Flea缓存，缓存服务器列表【cacheServerList】为空");
        }
        // 获取缓存数据有效期（单位：s）
        int expiry = CacheConfigUtils.getExpiry(name);
        // 获取空缓存数据有效期（单位：s）
        int nullCacheExpiry = CacheConfigUtils.getNullCacheExpiry();
        // 获取Redis哨兵配置开关（1：开启 0：关闭），如果不配置也默认开启
        boolean isSwitchOpen = CacheConfigUtils.isSwitchOpen(CacheConstants.RedisConfigConstants.REDIS_SENTINEL_CONFIG_SWITCH);

        AbstractFleaCache fleaCache;
        if (isSwitchOpen) { // 开关启用，按实际缓存处理
            // 获取缓存组名
            String group = cacheServerList.get(0).getGroup();
            // 初始化指定连接池名【group】的Redis哨兵模式连接池
            RedisSentinelPool.getInstance(group).initialize(cacheServerList);
            // 获取哨兵模式下的指定连接池名【group】的Redis客户端
            RedisClient redisClient = RedisClientFactory.getInstance(group, CacheModeEnum.SENTINEL);
            // 创建一个Redis Flea缓存
            fleaCache = new RedisFleaCache(name, expiry, nullCacheExpiry, CacheModeEnum.SENTINEL, redisClient);

            if (LOGGER.isDebugEnabled()) {
                Object obj = new Object() {};
                LOGGER.debug1(obj, "Pool Name = {}", RedisSentinelPool.getInstance(group).getPoolName());
                LOGGER.debug1(obj, "Pool = {}", RedisSentinelPool.getInstance(group).getJedisSentinelPool());
            }
        } else { // 开关关闭，默认返回空缓存实现
            fleaCache = new EmptyFleaCache(name, expiry, nullCacheExpiry);
        }

        return fleaCache;
    }
}
```

### 3.10 定义核心Flea缓存管理类
[CoreFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/core/manager/CoreFleaCacheManager.java) 继承 **AbstractFleaCacheManager** ，实现 **newCache** 方法，用于创建一个核心 **Flea** 缓存。
```java
/**
 * 核心Flea缓存管理类，用于接入Flea框架管理核心Flea缓存。
 *
 * <p> 核心Flea缓存是Flea框架提供出来的整合各类缓存的缓存实现。
 *
 * <p> 方法 {@code newCache}，用于创建一个核心Flea缓存，
 * 它里面包含了读、写、删除和清空缓存的基本操作。
 *
 * @author huazie
 * @version 1.0.0
 * @see CoreFleaCache
 * @since 1.0.0
 */
public class CoreFleaCacheManager extends AbstractFleaCacheManager {

    @Override
    protected AbstractFleaCache newCache(String name, int expiry) {
        return new CoreFleaCache(name);
    }
}
```

### 3.11 整合接入自测 
单元测试类 [FleaCacheTest](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/test/java/com/huazie/fleaframework/cache/FleaCacheTest.java)

首先，这里需要按照 Flea缓存配置文件 (**flea-cache-config.xml**) 中的缓存服务器 **cache-server** 中地址部署相应的 Memcached 和 Redis 服务，可参考笔者的 [这篇博文](../../../../../../2019/08/30/flea-framework/flea-cache/flea-cache-windows-more-services/)。

```java
    @Test
    public void testCoreFleaCache() {
        try {
            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(CacheEnum.FleaCore.getName());
            // 这里根据传入不同的缓存主关键字，可以得到不同的缓存实现
            AbstractFleaCache cache = manager.getCache("fleaconfigdata");
            LOGGER.debug("Cache={}", cache);
            //## 1.  简单字符串
//            cache.put("menu1", "huazie");
//            cache.put("menu2", "helloworld");
            cache.get("menu1");
            cache.get("menu2");
//            cache.delete("menu1");
//            cache.clear();
            cache.getCacheKey();
            LOGGER.debug(cache.getCacheName() + ">>>" + cache.getCacheDesc());
        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
```
经过上面的介绍，核心Flea缓存相关的内容，基本上算是讲解完毕。在不改变现有业务代码的基础上，相关缓存 **cache** 可以通过修改其归属的缓存数据类型 **type**，实现各类缓存数据，多种缓存系统之间的无缝迁移。

## 4. 进阶接入
### 4.1 定义核心Spring缓存类 
[CoreSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/core/impl/CoreSpringCache.java) 继承抽象 **Spring** 缓存 [AbstractSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCache.java)，用于对接 **Spring**；从构造方法可见，该类初始化使用核心Flea缓存类 **CoreFleaCache**。

```java
/**
 * <p> 核心Spring缓存类 </p>
 *
 * @author huazie
 */
public class CoreSpringCache extends AbstractSpringCache {

    /**
     * <p> 带参数构造方法 </p>
     *
     * @param name      缓存主关键字
     * @param fleaCache Flea Cache具体实现
     */
    public CoreSpringCache(String name, IFleaCache fleaCache) {
        super(name, fleaCache);
    }

    /**
     * <p> 带参数构造方法 </p>
     *
     * @param name 缓存主关键字
     */
    public CoreSpringCache(String name) {
        super(name, new CoreFleaCache(name));
    }
}
```

### 4.2 定义核心Spring缓存管理类
[CoreSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/core/manager/CoreSpringCacheManager.java) 继承抽象 **Spring** 缓存管理类 [AbstractSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCacheManager.java)，用于对接 **Spring**；基本实现同核心 **Flea** 缓存管理类 **CoreFleaCacheManager**，唯一不同在于 **newCache** 的实现，这边是 **new** 一个核心 **Spring** 缓存 **CoreSpringCache**。
```java
/**
 * 核心Spring缓存管理类，用于接入Spring框架管理核心Spring缓存。
 *
 * <p> 核心Spring缓存是Flea框架提供出来的整合各类缓存的缓存实现。
 *
 * <p> 方法【{@code newCache}】用于创建一个核心Spring缓存，
 * 而它内部是由核心Flea缓存【{@code CoreFleaCache}】实现具体的
 * 读、写、删除 和 清空 缓存的基本操作。
 *
 * @author huazie
 * @version 1.0.0
 * @see CoreSpringCache
 * @since 1.0.0
 */
public class CoreSpringCacheManager extends AbstractSpringCacheManager {

    @Override
    protected AbstractSpringCache newCache(String name, int expiry) {
        return new CoreSpringCache(name);
    }
}

```

### 4.3 Spring配置

```xml
<!-- 配置核心Flea缓存管理类 RedisSpringCacheManager -->
<bean id="coreSpringCacheManager" class="com.huazie.fleaframework.cache.core.manager.CoreSpringCacheManager" />

<!-- 开启缓存 -->
<cache:annotation-driven cache-manager="coreSpringCacheManager" proxy-target-class="true"/>
```

### 4.4 缓存自测

```java
    private ApplicationContext applicationContext;

    @Before
    public void init() {
        applicationContext = new ClassPathXmlApplicationContext("applicationContext.xml");
        LOGGER.debug("ApplicationContext={}", applicationContext);
    }
    
    @Test
    public void testCoreSpringCache() {
        try {
            AbstractSpringCacheManager manager = (CoreSpringCacheManager) applicationContext.getBean("coreSpringCacheManager");
            LOGGER.debug("CoreSpringCacheManager={}", manager);

            AbstractSpringCache cache = manager.getCache("fleaparadetail");
            LOGGER.debug("Cache={}", cache);

            //## 1.  简单字符串
//			cache.put("menu1", "huazie");
//            cache.get("menu1");
//            cache.get("menu1", String.class);

            //## 2.  简单对象(要是可以序列化的对象)
//			String user = new String("huazie");
//			cache.put("user", user);
//			LOGGER.debug(cache.get("user", String.class));
            cache.clear();

            //## 3.  List塞对象
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

### 4.5 业务逻辑层接入缓存管理
**@Cacheable** 使用，**value** 为缓存名，也作缓存主关键字， **key** 为具体的缓存键

```java
@Cacheable(value = "fleaparadetail", key = "#paraType + '_' + #paraCode")
public FleaParaDetail getParaDetail(String paraType, String paraCode) throws Exception {

    List<FleaParaDetail> fleaParaDetails = fleaParaDetailDao.getParaDetail(paraType, paraCode);
    FleaParaDetail fleaParaDetailValue = null;

    if (CollectionUtils.isNotEmpty(fleaParaDetails)) {
        fleaParaDetailValue = fleaParaDetails.get(0);
    }

    return fleaParaDetailValue;
}
```
# 四、结语
到目前为止，整合 **Memcached** 和 **Redis** 接入的工作已经全部完成，相信各位已经能够接入系统~~~

 

 