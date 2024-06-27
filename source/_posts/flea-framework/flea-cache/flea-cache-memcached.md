---
title: flea-cache使用之Memcached接入
date: 2019-08-18 22:52:46
updated: 2023-12-20 16:34:29
categories:
  - [开发框架-Flea,flea-cache]
tags:
  - flea-framework
  - flea-cache
  - Memcached
---

![](/images/cache.png)

# 1. 参考
[flea-cache使用之Memcached接入 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-cache)

<!-- more -->

![](flea-cache-memcached.png)

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 2. 依赖
[Memcached-Java-Client-3.0.2.jar](https://mvnrepository.com/artifact/com.whalin/Memcached-Java-Client)
```xml
<!-- Memcached相关 -->
<dependency>
    <groupId>com.whalin</groupId>
    <artifactId>Memcached-Java-Client</artifactId>
    <version>3.0.2</version>
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
[IFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/IFleaCache.java)  自定义缓存接口，包含了一些增删改查的方法
```java
/**
 * <p> 自定义Cache接口类（主要定义了一些增删改查的方法） </p>
 *
 * @author huazie
 */
public interface IFleaCache {

    /**
     * <p> 读缓存 </p>
     *
     * @param key 数据键关键字
     * @return 数据值
     */
    Object get(String key);

    /**
     * <p> 写缓存 </p>
     *
     * @param key   数据键关键字
     * @param value 数据值
     */
    void put(String key, Object value);

    /**
     * <p> 清空所有缓存 </p>
     */
    void clear();

    /**
     * <p> 删除指定数据键关键字对应的缓存 </p>
     *
     * @param key 数据键关键字
     */
    void delete(String key);

    /**
     * <p> 获取 记录当前Cache所有数据键关键字 的Set集合 </p>
     *
     * @return 数据键key的集合
     */
    Set<String> getCacheKey();

    /**
     * <p> 获取缓存所属系统名 </p>
     *
     * @return 缓存所属系统名
     */
    String getSystemName();
}
```
## 3.2 定义抽象Flea缓存类
[AbstractFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCache.java) 抽象缓存类，实现了Flea缓存接口的读、写、删除和清空缓存的基本操作。
```java
/**
 * 抽象Flea Cache类，实现了Flea缓存接口的读、写、删除和清空缓存的基本操作。
 *
 * <p> 它定义本地读、本地写和本地删除的抽象方法，由子类实现具体的读、
 * 写和删除缓存的操作。
 *
 * <p> 在实际调用写缓存操作时，会同时记录当前缓存数据的数据键关键字
 * 【{@code key}】到专门的数据键关键字的缓存中，以Set集合存储。
 *
 * <p> 比如缓存数据主关键字为【{@code name}】，需要存储的数据键关键字为
 * 【{@code key}】，则在实际调用写缓存操作时，会操作两条缓存数据：<br/>
 * 一条是具体的数据缓存，缓存键为“系统名_name_key”，可查看方法
 * 【{@code getNativeKey}】，有效期从配置中获取；<br/>
 * 一条是数据键关键字的缓存，缓存键为“系统名_name”，可查看方法
 * 【{@code getNativeCacheKey}】，默认永久有效。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public abstract class AbstractFleaCache implements IFleaCache {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(AbstractFleaCache.class);

    private final String name;  // 缓存数据主关键字

    private final int expiry;  // 缓存数据有效期（单位：s）

    private final int nullCacheExpiry; // 空缓存数据有效期（单位：s）

    protected CacheEnum cache;  // 缓存实现

    public AbstractFleaCache(String name, int expiry, int nullCacheExpiry) {
        this.name = name;
        this.expiry = expiry;
        this.nullCacheExpiry = nullCacheExpiry;
    }

    @Override
    public Object get(String key) {
        Object value = null;
        try {
            value = getNativeValue(getNativeKey(key));
            if (value instanceof NullCache) {
                value = null;
            }
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error1(new Object() {}, "The action of getting [" + cache.getName() + "] cache occurs exception : ", e);
            }
        }
        return value;
    }

    @Override
    public void put(String key, Object value) {
        try {
            putNativeValue(getNativeKey(key), value, expiry);
            // 将指定Cache的key添加到Set集合，并存于缓存中
            addCacheKey(key);
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error1(new Object() {}, "The action of adding [" + cache.getName() + "] cache occurs exception : ", e);
            }
        }
    }

    @Override
    public void clear() {
        Set<String> keySet = getCacheKey();
        if (CollectionUtils.isNotEmpty(keySet)) {
            for (String key : keySet) {
                deleteNativeValue(getNativeKey(key));
            }
            // 删除 记录当前Cache所有数据键关键字 的缓存
            deleteCacheAllKey();
        }
    }

    @Override
    public void delete(String key) {
        try {
            deleteNativeValue(getNativeKey(key));
            // 从 记录当前Cache所有数据键关键字 的缓存中 删除指定数据键关键字key
            deleteCacheKey(key);
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error1(new Object() {}, "The action of deleting [" + cache.getName() + "] cache occurs exception : ", e);
            }
        }
    }

    /**
     * <p> 将指定数据键关键字{@code key}记录到当前Cache所有数据键关键字的缓存中 </p>
     *
     * @param key 指定Cache的数据键关键字
     * @since 1.0.0
     */
    private void addCacheKey(String key) {
        Set<String> keySet = getCacheKey();
        if (CollectionUtils.isEmpty(keySet)) {
            keySet = new HashSet<>();
        }
        if (!keySet.contains(key)) { // 只有其中不存在，才重新设置
            keySet.add(key);
            putNativeValue(getNativeCacheKey(name), keySet, CommonConstants.NumeralConstants.INT_ZERO);
        }
    }

    /**
     * <p> 从 记录当前Cache所有数据键关键字 的缓存中 删除指定数据键关键字{@code key} </p>
     *
     * @param key 指定Cache的数据键关键字
     * @since 1.0.0
     */
    private void deleteCacheKey(String key) {
        Set<String> keySet = getCacheKey();
        if (CollectionUtils.isNotEmpty(keySet)) {
            // 存在待删除的数据键关键字
            if (keySet.contains(key)) {
                if (CommonConstants.NumeralConstants.INT_ONE == keySet.size()) {
                    deleteCacheAllKey(); // 直接将记录当前Cache所有数据键关键字的缓存从缓存中清空
                } else {
                    // 将数据键关键字从Set集合中删除
                    keySet.remove(key);
                    // 重新覆盖当前Cache所有数据键关键字的缓存信息
                    putNativeValue(getNativeCacheKey(name), keySet, CommonConstants.NumeralConstants.INT_ZERO);
                }
            } else {
                if (LOGGER.isDebugEnabled()) {
                    LOGGER.debug1(new Object() {}, "The CacheKey of [{}] is not exist", key);
                }
            }
        }
    }

    /**
     * <p> 删除 记录当前Cache所有数据键关键字 的缓存 </p>
     *
     * @since 1.0.0
     */
    private void deleteCacheAllKey() {
        try {
            deleteNativeValue(getNativeCacheKey(name));
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error1(new Object() {}, "The action of deleting [" + cache.getName() + "] cache occurs exception : ", e);
            }
        }
    }

    @Override
    @SuppressWarnings(value = "unchecked")
    public Set<String> getCacheKey() {
        Set<String> keySet = null;
        try {
            Object keySetObj = getNativeValue(getNativeCacheKey(name));
            if (ObjectUtils.isNotEmpty(keySetObj) && keySetObj instanceof Set) {
                keySet = (Set<String>) keySetObj;
            }
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error1(new Object() {}, "The action of getting [" + cache.getName() + "] cache occurs exception : ", e);
            }
        }
        return keySet;
    }

    /**
     * <p> 获取缓存值 </p>
     *
     * @param key 缓存数据键关键字
     * @return 缓存值
     * @since 1.0.0
     */
    public abstract Object getNativeValue(String key);

    /**
     * <p> 添加缓存数据 </p>
     *
     * @param key    缓存数据键关键字
     * @param value  缓存值
     * @param expiry 有效期（单位：s）
     * @since 1.0.0
     */
    public abstract Object putNativeValue(String key, Object value, int expiry);

    /**
     * <p> 删除指定缓存数据 </p>
     *
     * @param key 缓存数据键关键字
     * @since 1.0.0
     */
    public abstract Object deleteNativeValue(String key);

    /**
     * <p> 获取实际存储的缓存键【缓存所属系统名 + 缓存名（缓存数据主关键字）+ 缓存数据键（缓存数据关键字）】 </p>
     *
     * @param key 缓存数据键关键字
     * @return 实际存储的缓存键
     * @since 1.0.0
     */
    private String getNativeKey(String key) {
        return StringUtils.strCat(getNativeCacheKey(name), CommonConstants.SymbolConstants.UNDERLINE, key);
    }

    /**
     * <p> 获取缓存主键【包含缓存所属系统名 + 缓存名（缓存数据主关键字）】 </p>
     *
     * @param name 缓存名【缓存数据主关键字】
     * @return 缓存主键【缓存所属系统名 + 缓存名（缓存数据主关键字）】
     * @since 1.0.0
     */
    protected String getNativeCacheKey(String name) {
        return StringUtils.strCat(getSystemName(), CommonConstants.SymbolConstants.UNDERLINE, name);
    }

    // 省略一些get方法
}

```
该类实现了IFleaCache接口，同时定义了三个抽象方法 :
```java
    public abstract Object getNativeValue(String key);

    public abstract void putNativeValue(String key, Object value, int expiry);

    public abstract void deleteNativeValue(String key);
```
这三个抽象方法由子类实现具体的读，写，删除缓存的原始操作

## 3.3 定义MemCached Flea缓存类
[MemCachedFleaCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/memcached/impl/MemCachedFleaCache.java) 该类继承 `AbstractFleaCache`，实现`Memcached` 缓存的接入使用；
```java
/**
 * MemCached Flea缓存类，实现了以Flea框架操作MemCached缓存的基本操作方法。
 *
 * <p> 在上述基本操作方法中，实际使用MemCached客户端【{@code} memCachedClient】
 * 读、写和删除MemCached缓存。其中写缓存方法【{@code putNativeValue}】在
 * 添加的数据值为【{@code null}】时，默认添加空缓存数据【{@code NullCache}】
 * 到MemCached中，有效期取初始化参数【{@code nullCacheExpiry}】。
 *
 * <p> 单个缓存接入场景，有效期配置可查看【memcached.properties】中的配置
 * 参数【memcached.nullCacheExpiry】
 *
 * <p> 整合缓存接入场景，有效期配置可查看【flea-cache-config.xml】
 * 中的缓存参数【{@code <cache-param key="fleacore.nullCacheExpiry"
 * desc="空缓存数据有效期（单位：s）">300</cache-param>}】
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class MemCachedFleaCache extends AbstractFleaCache {

    private final MemCachedClient memCachedClient;  // MemCached客户端

    /**
     * <p> 初始化MemCached Flea缓存类 </p>
     *
     * @param name            缓存数据主关键字
     * @param expiry          缓存数据有效期（单位：s）
     * @param nullCacheExpiry 空缓存数据有效期（单位：s）
     * @param memCachedClient MemCached客户端
     * @since 1.0.0
     */
    public MemCachedFleaCache(String name, int expiry, int nullCacheExpiry, MemCachedClient memCachedClient) {
        super(name, expiry, nullCacheExpiry);
        this.memCachedClient = memCachedClient;
        cache = CacheEnum.MemCached;
    }

    @Override
    public Object getNativeValue(String key) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(new Object() {}, "KEY = {}", key);
        }
        return memCachedClient.get(key);
    }

    @Override
    public Object putNativeValue(String key, Object value, int expiry) {
        if (ObjectUtils.isEmpty(value))
            return memCachedClient.set(key, new NullCache(key), new Date(getNullCacheExpiry() * 1000));
        else
            return memCachedClient.set(key, value, new Date(expiry * 1000));
    }

    @Override
    public Object deleteNativeValue(String key) {
        return memCachedClient.delete(key);
    }

    @Override
    public String getSystemName() {
        return MemCachedConfig.getConfig().getSystemName();
    }
}

```
到这一步为止，底层的Flea缓存接口和实现已经完成，但目前还不能使用；
## 3.4 定义Memcached连接池
[MemCachedPool](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/memcached/MemCachedPool.java) 用于初始化 `MemCached` 的套接字连接池。
```java
/**
 * Flea MemCached连接池，用于初始化MemCached的套接字连接池。
 *
 * <p> 针对单独缓存接入场景，采用默认连接池初始化的方式；<br/>
 * 可参考如下：
 * <pre>
 *   // 初始化默认连接池
 *   MemCachedPool.getInstance().initialize(); </pre>
 *
 * <p> 针对整合缓存接入场景，采用指定连接池初始化的方式；<br/>
 * 可参考如下：
 * <pre>
 *   // 初始化指定连接池
 *   MemCachedPool.getInstance(group).initialize(cacheServerList); </pre>
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public class MemCachedPool {

    private String poolName; // 连接池名

    private MemCachedConfig memCachedConfig; // MemCached 配置信息

    private SockIOPool sockIOPool; // MemCached SockIOPool

    private MemCachedPool() {
    }

    /**
     * <p> 获取MemCached连接池实例 (默认) </p>
     *
     * @return MemCached连接池实例对象
     * @since 1.0.0
     */
    public static MemCachedPool getInstance() {
        MemCachedPool memCachedPool = new MemCachedPool();
        memCachedPool.memCachedConfig = MemCachedConfig.getConfig();
        memCachedPool.sockIOPool = SockIOPool.getInstance();
        return memCachedPool;
    }

    /**
     * <p> 获取MemCached连接池实例（指定连接池名poolName） </p>
     *
     * @param poolName 连接池名
     * @return MemCached连接池实例对象
     * @since 1.0.0
     */
    public static MemCachedPool getInstance(String poolName) {
        MemCachedPool memCachedPool = new MemCachedPool();
        memCachedPool.poolName = poolName;
        memCachedPool.sockIOPool = SockIOPool.getInstance(poolName);
        return memCachedPool;
    }

    /**
     * <p> 初始化MemCached连接池 </p>
     *
     * @since 1.0.0
     */
    public void initialize() {
        // 省略。。。。。。
    }

    /**
     * <p> 初始化MemCached连接池 </p>
     *
     * @param cacheServerList 缓存服务器集
     * @since 1.0.0
     */
    public void initialize(List<CacheServer> cacheServerList) {
        // 省略。。。。。。
    }

    // 详见 GitHub 链接
}

```
## 3.5 Memcached配置文件
flea-cache 读取 [memcached.properties](https://github.com/Huazie/flea-framework/blob/dev/flea-config/src/main/resources/flea/cache/memcached.properties)（**Memcached**配置文件），用作初始化 **MemCachedPool**
```bash
# Memcached配置
# Memcached缓存所属系统名
memcached.systemName=FleaFrame

# Memcached服务器地址
memcached.server=127.0.0.1:31113,127.0.0.1:31114

# Memcached服务器权重分配
memcached.weight=1,1

# 初始化时对每个服务器建立的连接数目
memcached.initConn=20

# 每个服务器建立最小的连接数
memcached.minConn=20

# 每个服务器建立最大的连接数
memcached.maxConn=500

# 自查线程周期进行工作，其每次休眠时间（单位：ms）
memcached.maintSleep=60000

# Socket的参数，如果是true在写数据时不缓冲，立即发送出去
memcached.nagle=true

# Socket阻塞读取数据的超时时间（单位：ms）
memcached.socketTO=3000

# Socket连接超时时间（单位：ms）
memcached.socketConnectTO=3000

# 空缓存数据有效期（单位：s）
memcached.nullCacheExpiry=10

# Memcached分布式hash算法
# 0 - native String.hashCode();
# 1 - original compatibility
# 2 - new CRC32 based
# 3 - MD5 Based
memcached.hashingAlg=3

```
## 3.6 定义抽象Flea缓存管理类
[AbstractFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractFleaCacheManager.java) 用于接入 `Flea` 框架管理缓存。
```java
/**
 * 抽象Flea缓存管理类，用于接入Flea框架管理缓存。
 *
 * <p> 同步集合类【{@code cacheMap}】, 存储的键为缓存数据主关键字，
 * 存储的值为具体的缓存实现类。<br/>
 * 如果是整合各类缓存接入，它的键对应缓存定义配置文件【flea-cache.xml】
 * 中的【{@code <cache key="缓存数据主关键字"></cache>}】；<br/>
 * 如果是单个缓存接入，它的键对应【applicationContext.xml】中
 * 【{@code <entry key="缓存数据主关键字"value="有效期（单位：s）"/>}】；
 *
 * <p> 抽象方法【{@code newCache}】，由子类实现具体的Flea缓存类创建。
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public abstract class AbstractFleaCacheManager {

    private static final ConcurrentMap<String, AbstractFleaCache> cacheMap = new ConcurrentHashMap<>();

    private Map<String, Integer> configMap = new HashMap<>();   // 各缓存的时间Map

    /**
     * <p> 获取所有的Flea缓存 </p>
     *
     * @return 所有的Flea缓存
     * @since 1.0.0
     */
    protected Collection<AbstractFleaCache> loadCaches() {
        return cacheMap.values();
    }

    /**
     * <p> 根据指定缓存名获取缓存对象 </p>
     *
     * @param name 缓存名
     * @return 缓存对象
     * @since 1.0.0
     */
    public AbstractFleaCache getCache(String name) {
        if (!cacheMap.containsKey(name)) {
            synchronized (cacheMap) {
                if (!cacheMap.containsKey(name)) {
                    Integer expiry = configMap.get(name);
                    if (ObjectUtils.isEmpty(expiry)) {
                        expiry = CommonConstants.NumeralConstants.INT_ZERO; // 表示永久
                        configMap.put(name, expiry);
                    }
                    cacheMap.put(name, newCache(name, expiry));
                }
            }
        }
        return cacheMap.get(name);
    }

    /**
     * <p> 新创建一个缓存对象 </p>
     *
     * @param name   缓存名
     * @param expiry 有效期（单位：s  其中0：表示永久）
     * @return 新建的缓存对象
     * @since 1.0.0
     */
    protected abstract AbstractFleaCache newCache(String name, int expiry);

    /**
     * <p> 设置各缓存有效期配置Map </p>
     *
     * @param configMap 有效期配置Map
     * @since 1.0.0
     */
    public void setConfigMap(Map<String, Integer> configMap) {
        this.configMap = configMap;
    }

}

```
## 3.7 定义Memcached Flea缓存管理类 
[MemCachedFleaCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/memcached/manager/MemCachedFleaCacheManager.java) 用于接入 `Flea` 框架管理`MemCached` 缓存
上述 `MemCachedFleaCache` 使用, 需要初始化 `MemCachedPool`
```java
/**
 * MemCached Flea缓存管理类，用于接入Flea框架管理MemCached 缓存。
 *
 * <p> 它的默认构造方法，用于单个缓存接入场景，首先新建MemCached客户端，
 * 然后初始化 MemCached 连接池。
 *
 * <p> 方法 {@code newCache}，用于创建一个MemCached Flea缓存，
 * 它里面包含了 读、写、删除 和 清空 缓存的基本操作。
 *
 * @author huazie
 * @version 1.0.0
 * @see MemCachedFleaCache
 * @since 1.0.0
 */
public class MemCachedFleaCacheManager extends AbstractFleaCacheManager {

    private MemCachedClient memCachedClient;   // MemCached客户端类

    /**
     * 用于新建MemCached客户端，并初始化MemCached连接池。
     *
     * @since 1.0.0
     */
    public MemCachedFleaCacheManager() {
        memCachedClient = new MemCachedClient();
        initPool();
    }

    /**
     * 以传入参数初始化MemCached客户端，并初始化MemCached连接池。
     *
     * @param memCachedClient MemCached客户端
     * @since 1.0.0
     */
    public MemCachedFleaCacheManager(MemCachedClient memCachedClient) {
        this.memCachedClient = memCachedClient;
        initPool();
    }

    /**
     * <p> 初始化MemCached连接池 </p>
     *
     * @since 1.0.0
     */
    private void initPool() {
        MemCachedPool.getInstance().initialize();
    }

    @Override
    protected AbstractFleaCache newCache(String name, int expiry) {
        int nullCacheExpiry = MemCachedConfig.getConfig().getNullCacheExpiry();
        return new MemCachedFleaCache(name, expiry, nullCacheExpiry, memCachedClient);
    }
}
```
好了，到了这一步，Memcached已接入完成，开始自测

##  3.8 Memcached接入自测 
[FleaCacheTest](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/test/java/com/huazie/fleaframework/cache/FleaCacheTest.java)  单元测试类可点击查看

首先，这里需要按照 `Memcached` 配置文件中的地址部署相应的 `Memcached` 服务，可参考笔者的 [这篇博文](../../../../../../2019/08/30/flea-framework/flea-cache/flea-cache-windows-more-services/)。

下面开始演示我们的 `Memcached` 接入自测：
```java
    @Test
    public void testMemeCachedFleaCache() {
        try {
            AbstractFleaCacheManager manager = FleaCacheManagerFactory.getFleaCacheManager(CacheEnum.MemCached.getName());
            AbstractFleaCache cache = manager.getCache("fleaparadetail");
            LOGGER.debug("Cache={}", cache);
            //## 1.  简单字符串
//            cache.put("menu1", "huazie");
//            cache.get("menu1");
            cache.delete("menu1");
//            cache.getCacheKey();
            LOGGER.debug(cache.getCacheName() + ">>>" + cache.getCacheDesc());
        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
```
# 4. 进阶接入
## 4.1 定义抽象Spring缓存
[AbstractSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCache.java) 与 **AbstractFleaCache** 不同之处在于，实现了 **Spring** 的 **Cache** 接口，用于对接 **Spring**，相关配置后面会介绍一下。
```java
/**
 * 抽象 Spring 缓存类，实现Spring的Cache接口 和 Flea
 * 的IFleaCache接口，由具体的Spring缓存管理类初始化。
 *
 * <p> 它实现了读、写、删除和清空缓存的基本操作方法，
 * 内部由具体Flea缓存实现类【{@code fleaCache}】
 * 调用对应的 读、写、删除 和 清空 缓存的基本操作方法。
 *
 * @author huazie
 * @version 1.1.0
 * @since 1.0.0
 */
public abstract class AbstractSpringCache implements Cache, IFleaCache {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(AbstractSpringCache.class);

    private final String name;  // 缓存主要关键字（用于区分）

    private final IFleaCache fleaCache; // 具体Flea缓存实现

    public AbstractSpringCache(String name, IFleaCache fleaCache) {
        this.name = name;
        this.fleaCache = fleaCache;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public IFleaCache getNativeCache() {
        return fleaCache;
    }

    @Override
    public ValueWrapper get(Object key) {
        if (ObjectUtils.isEmpty(key))
            return null;
        ValueWrapper wrapper = null;
        Object cacheValue = get(key.toString());
        if (ObjectUtils.isNotEmpty(cacheValue)) {
            wrapper = new SimpleValueWrapper(cacheValue);
        }
        return wrapper;
    }

    @Override
    public <T> T get(Object key, Class<T> type) {
        if (ObjectUtils.isEmpty(key) || ObjectUtils.isEmpty(type))
            return null;
        Object cacheValue = get(key.toString());
        if (!type.isInstance(cacheValue)) {
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug1(new Object() {}, "Cached value is not of required type [{}] : {}", type.getName(), cacheValue);
            }
            return null;
        }
        return type.cast(cacheValue);
    }

    @Override
    public <T> T get(Object key, Callable<T> valueLoader) {
        return null;
    }

    @Override
    public Object get(String key) {
        if (StringUtils.isBlank(key))
            return null;
        Object obj = null;
        if (LOGGER.isDebugEnabled()) {
            obj = new Object() {};
            LOGGER.debug1(obj, "KEY = {}", key);
        }
        Object cacheValue = fleaCache.get(key);
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(obj, "VALUE = {}", cacheValue);
        }
        return cacheValue;
    }

    @Override
    public void put(Object key, Object value) {
        if (ObjectUtils.isEmpty(key))
            return;
        put(key.toString(), value);
    }

    @Override
    public ValueWrapper putIfAbsent(Object key, Object value) {
        if (ObjectUtils.isEmpty(key))
            return null;
        ValueWrapper wrapper = null;
        Object cacheValue = get(key.toString());
        if (ObjectUtils.isEmpty(cacheValue)) {
            put(key.toString(), value);
        } else {
            wrapper = new SimpleValueWrapper(cacheValue);
        }
        return wrapper;
    }

    @Override
    public void put(String key, Object value) {
        fleaCache.put(key, value);
    }

    @Override
    public void evict(Object key) {
        if (ObjectUtils.isEmpty(key))
            return;
        delete(key.toString());
    }

    @Override
    public void clear() {
        fleaCache.clear();
    }

    @Override
    public void delete(String key) {
        fleaCache.delete(key);
    }

    @Override
    public Set<String> getCacheKey() {
        return fleaCache.getCacheKey();
    }

    @Override
    public String getSystemName() {
        return fleaCache.getSystemName();
    }
}

```
## 4.2 定义Memcached Spring缓存
[MemCachedSpringCache](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/memcached/impl/MemCachedSpringCache.java) 只定义构造方法，使用 `MemCachedFleaCache` 作为具体缓存实现。
```java
/**
 * MemCached Spring缓存类，继承了抽象Spring缓存父类的
 * 读、写、删除 和 清空 缓存的基本操作方法，由MemCached Spring缓存管理类初始化。
 *
 * <p> 它的构造方法中，必须传入一个具体Flea缓存实现类，这里我们使用
 * MemCached Flea缓存【{@code MemCachedFleaCache}】。
 *
 * @author huazie
 * @version 1.0.0
 * @see MemCachedFleaCache
 * @since 1.0.0
 */
public class MemCachedSpringCache extends AbstractSpringCache {

    /**
     * <p> 带参数的构造方法，初始化MemCached Spring缓存类 </p>
     *
     * @param name      缓存数据主关键字
     * @param fleaCache 具体缓存实现
     * @since 1.0.0
     */
    public MemCachedSpringCache(String name, IFleaCache fleaCache) {
        super(name, fleaCache);
    }

    /**
     * <p> 带参数的构造方法，初始化MemCached Spring缓存类 </p>
     *
     * @param name            缓存数据主关键字
     * @param expiry          缓存数据有效期（单位：s）
     * @param nullCacheExpiry 空缓存数据有效期（单位：s）
     * @param memCachedClient MemCached客户端
     * @since 1.0.0
     */
    public MemCachedSpringCache(String name, int expiry, int nullCacheExpiry, MemCachedClient memCachedClient) {
        this(name, new MemCachedFleaCache(name, expiry, nullCacheExpiry, memCachedClient));
    }
}
```
## 4.3 定义抽象Spring缓存管理类
 [AbstractSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/AbstractSpringCacheManager.java) 继承 `AbstractTransactionSupportingCacheManager`，用于对接 `Spring`。

```java
/**
 * 抽象Spring缓存管理类，用于接入Spring框架管理缓存。
 *
 * <p> 同步集合类【{@code cacheMap}】, 存储的键为缓存数据主关键字，
 * 存储的值为具体的缓存实现类。<br/>
 * 如果是整合各类缓存接入，它的键对应缓存定义配置文件【flea-cache.xml】
 * 中的【{@code <cache key="缓存数据主关键字"></cache>}】；<br/>
 * 如果是单个缓存接入，它的键对应【applicationContext.xml】中
 * 【{@code <entry key="缓存数据主关键字" value="有效期（单位：s）"/>}】；
 *
 * <p> 抽象方法【{@code newCache}】，由子类实现具体的Spring缓存类创建。
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public abstract class AbstractSpringCacheManager extends AbstractTransactionSupportingCacheManager {

    private static final ConcurrentMap<String, AbstractSpringCache> cacheMap = new ConcurrentHashMap<>();

    private Map<String, Integer> configMap = new HashMap<>();   // 各缓存的时间Map

    @Override
    protected Collection<? extends AbstractSpringCache> loadCaches() {
        return cacheMap.values();
    }

    @Override
    public AbstractSpringCache getCache(String name) {
        if (!cacheMap.containsKey(name)) {
            synchronized (cacheMap) {
                if (!cacheMap.containsKey(name)) {
                    Integer expiry = configMap.get(name);
                    if (expiry == null) {
                        expiry = CommonConstants.NumeralConstants.INT_ZERO; // 表示永久
                        configMap.put(name, expiry);
                    }
                    cacheMap.put(name, newCache(name, expiry));
                }
            }
        }
        return cacheMap.get(name);
    }

    /**
     * 由子类实现该方法，新创建一个抽象 Spring 缓存的子类。
     *
     * @param name   缓存名
     * @param expiry 有效期（单位：s  其中0：表示永久）
     * @return 新建的缓存对象
     * @since 1.0.0
     */
    protected abstract AbstractSpringCache newCache(String name, int expiry);

    /**
     * <p> 设置各缓存有效期配置Map </p>
     *
     * @param configMap 有效期配置Map
     * @since 1.0.0
     */
    public void setConfigMap(Map<String, Integer> configMap) {
        this.configMap = configMap;
    }
}
```
## 4.4 定义Memcached Spring缓存管理类
[MemCachedSpringCacheManager](https://github.com/Huazie/flea-framework/blob/dev/flea-cache/src/main/java/com/huazie/fleaframework/cache/memcached/manager/MemCachedSpringCacheManager.java) 基本实现同 `MemCachedFleaCacheManager`，不同在于 `newCache` 方法返回一个  `MemCachedSpringCache` 的对象。

```java
/**
 * MemCached Spring缓存管理类，用于接入Spring框架管理 MemCached 缓存。
 *
 * <p> 它的默认构造方法，用于单个缓存接入场景，首先新建MemCached 客户端，
 * 然后初始化 MemCached 连接池。
 *
 * <p> 方法【{@code newCache}】用于创建一个MemCached Spring缓存，
 * 而它内部是由MemCached Flea缓存实现具体的 读、写、删除 和 清空
 * 缓存的基本操作。
 *
 * @author huazie
 * @version 1.0.0
 * @see MemCachedSpringCache
 * @since 1.0.0
 */
public class MemCachedSpringCacheManager extends AbstractSpringCacheManager {

    private MemCachedClient memCachedClient;   // MemCached客户端类

    /**
     * 用于新建MemCached客户端，并初始化MemCached连接池。
     *
     * @since 1.0.0
     */
    public MemCachedSpringCacheManager() {
        memCachedClient = new MemCachedClient();
        initPool();
    }

    /**
     * 以传入参数初始化MemCached客户端，并初始化MemCached连接池。
     *
     * @param memCachedClient MemCached客户端
     */
    public MemCachedSpringCacheManager(MemCachedClient memCachedClient) {
        this.memCachedClient = memCachedClient;
        initPool();
    }

    /**
     * <p> 初始化MemCached连接池 </p>
     *
     * @since 1.0.0
     */
    private void initPool() {
        MemCachedPool.getInstance().initialize();
    }

    @Override
    protected AbstractSpringCache newCache(String name, int expiry) {
        int nullCacheExpiry = MemCachedConfig.getConfig().getNullCacheExpiry();
        return new MemCachedSpringCache(name, expiry, nullCacheExpiry, memCachedClient);
    }
}
```
## 4.5 spring 配置
```xml
<!--
   配置缓存管理MemCachedSpringCacheManager
   配置缓存时间 configMap (key缓存对象名称 value缓存过期时间)
-->
<bean id="memCachedSpringCacheManager" class="com.huazie.fleaframework.cache.memcached.manager.MemCachedSpringCacheManager">
    <property name="configMap">
        <map>
            <entry key="fleaparadetail" value="86400"/>
        </map>
    </property>
</bean>

 <!-- 开启缓存， 此处定义表示由spring接入来管理缓存访问 -->
 <cache:annotation-driven cache-manager="memCachedSpringCacheManager" proxy-target-class="true"/>
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
    public void testMemCachedSpringCache() {
        try {
            AbstractSpringCacheManager manager = (MemCachedSpringCacheManager) applicationContext.getBean("memCachedSpringCacheManager");
            LOGGER.debug("MemCachedCacheManager={}", manager);

            AbstractSpringCache cache = manager.getCache("fleaparadetail");
            LOGGER.debug("Cache={}", cache);

            Set<String> cacheKey = cache.getCacheKey();
            LOGGER.debug("CacheKey = {}", cacheKey);
            // 缓存清理
//            cache.clear();

            //## 1.  简单字符串
//			cache.put("menu1", "huazie");
//            cache.get("menu1");
//            cache.get("menu1", String.class);

            //## 2.  简单对象(要是可以序列化的对象)
//			String user = new String("huazie");
//			cache.put("user", user);
//			LOGGER.debug(cache.get("user", String.class));

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
## 4.7 业务逻辑层接入缓存管理
**@Cacheable** 使用，**value** 为缓存名，也作缓存主关键字， **key** 为具体的缓存键。

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
# 结语
至此，**Memcached** 的接入工作已经全部完成，下一篇将讲解 [flea-cache使用之Redis分片模式接入](../../../../../../2021/11/18/flea-framework/flea-cache/flea-cache-redissharded/)，敬请期待哦！！！