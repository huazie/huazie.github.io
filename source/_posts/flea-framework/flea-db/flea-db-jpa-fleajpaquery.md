---
title: flea-db使用之基于对象池的FleaJPAQuery
date: 2019-09-18 15:33:05
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
  - [开发语言-Java,Java设计模式]
tags:
  - flea-framework
  - flea-db
  - FleaJPAQuery
  - 对象池
---

![](/images/jpa-logo.png)

# 引言
书接上回《JPA封装介绍》博文，提到 **FleaJPAQuery** 在前一个版本（单例模式）下存在并发的问题，下面首先来分析一下，然后再介绍目前基于对象池的解决方案。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 一、问题分析
1. 上个版本 **FleaJPAQuery** 使用单例模式获取，意味着在同一个服务器中DAO层获取的  **FleaJPAQuery** 有且仅有一个。部分代码如下所示：

    ```java
        private static volatile FleaJPAQuery query;

        private FleaJPAQuery() {
        }
        
        public static FleaJPAQuery getQuery() {
            if (ObjectUtils.isEmpty(query)) {
                synchronized (FleaJPAQuery.class) {
                    if (ObjectUtils.isEmpty(query)) {
                        query = new FleaJPAQuery();
                    }
                }
            }
            return query;
        }

    ```
2. DAO层需要频繁获取 **FleaJPAQuery** ，来实现数据库查询功能，一开始认为单例模式最为适用；但是后来使用发现，**FleaJPAQuery** 获取之后，还需要调用它的公共方法，用于组装查询语句和获取查询结果；在调用期间 **FleaJPAQuery** 实例的成员变量是在不断被修改；最后的查询结果恰好也是依赖这些成员变量去调用(如下**createQuery** 方法)；
 
    ```java
        private TypedQuery createQuery(boolean isSingle) throws DaoException {
            if (ObjectUtils.isEmpty(sourceClazz)) {
                throw new DaoException("ERROR-DB-DAO0000000008");
            }
            if (!isSingle) {
                criteriaQuery.select(root);
            }
            if (CollectionUtils.isNotEmpty(predicates)) {
                criteriaQuery.where(criteriaBuilder.and(predicates.toArray(new Predicate[0])));
            }
            if (CollectionUtils.isNotEmpty(orders)) {
                criteriaQuery.orderBy(orders);
            }
            if (CollectionUtils.isNotEmpty(groups)) {
                criteriaQuery.groupBy(groups);
            }
            return entityManager.createQuery(criteriaQuery);
        }
    ```

3. 在并发下情况下，原始方案DAO层操作获取的 **FleaJPAQuery** 始终是一个；因为存在组装查询语句的过程，不同的数据查询操作之间就会相互影响，导致获取的查询结果不符合预期或者获取查询结果报错；另外组装查询语句的过程也可能直接报错；
4. 基于上面的分析，也就是要每个DAO层操作获取的 **FleaJPAQuery** 之间互不影响，同时又要保证 **FleaJPAQuery** 尽可能少的被新创建；这就需要为**FleaJPAQuery** 维护一个对象池，需要时从对象池中取出来用，用完重置状态后再归还给对象池。


# 二、方案讲解

本文采用 **Apache Commons Pool2** 的对象池化框架，其依赖如下：

```xml
    <!--  -->
    <dependency>
        <groupId>org.apache.commons</groupId>
        <artifactId>commons-pool2</artifactId>
        <version>2.4.3</version>
    </dependency>
```

## 2.1 Flea对象池配置

笔者的 **Flea** 框架下专门为对象池配置定义了一个通用的对象池配置对象 `FleaObjectPoolConfig`，它继承自 `GenericObjectPoolConfig`。 `GenericObjectPoolConfig` 是**Apache Commons Pool2** 库中用于配置 `GenericObjectPool` 对象池行为的类，它提供了丰富的参数以支持各种对象池的定制化需求，具体的参数可参考笔者的[《对象池 GenericObjectPool 配置参数详解》](../../../../../../2019/09/25/java/java-design-patterns/genericobjectpool-config-param/)。


```java
public class FleaObjectPoolConfig extends GenericObjectPoolConfig {
    public FleaObjectPoolConfig() {
        setTestWhileIdle(true);
        setMinEvictableIdleTimeMillis(60000);
        setTimeBetweenEvictionRunsMillis(30000);
        setNumTestsPerEvictionRun(-1);
    }
}
```

在 `FleaObjectPoolConfig` 的无参构造方法，初始化了部分默认配置，如下：

- `setTestWhileIdle(true)` ：指明连接是否被空闲连接回收器(如果有)进行检验。如果检测失败，则连接将被从池中去除；`true：是`；
- `setMinEvictableIdleTimeMillis(60000) `：连接在池中保持空闲而不被空闲连接回收器线程（如果有）回收的最小时间值，单位毫秒；
- `setTimeBetweenEvictionRunsMillis(30000)` ：在空闲连接回收器线程运行期间休眠的时间值，以毫秒为单位。 如果设置为非正数，则不运行空闲连接回收器线程；
- `setNumTestsPerEvictionRun(-1)` ：在每次空闲连接回收器线程（如果有）运行时检查的连接数量。如果`numTestsPerEvictionRun>=0`, 则取`numTestsPerEvictionRun` 和池内的链接数的较小值作为每次检测的链接数【`Math.min(numTestsPerEvictionRun, this.idleObjects.size())`】；如果 `numTestsPerEvictionRun<0`，则每次检查的连接数是检查时池内连接的总数除以这个值的绝对值再向上取整的结果【`(int)Math.ceil((double)this.idleObjects.size() / Math.abs((double)numTestsPerEvictionRun))`】。


## 2.2 Flea JPA查询对象池配置

[FleaJPAQueryPoolConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPoolConfig.java) 继承自 **FleaObjectPoolConfig**  。主要初始化以下对象池配置信息：

 1. **maxTotal**  最大连接数
 2. **maxIdle**  最大空闲连接数
 3. **minIdle**  最小空闲连接数
 4. **maxWaitMillis**  获取连接时的最大等待毫秒数

添加 `flea-config.xml` 配置信息如下：

```xml
    <config-items key="flea-jpa-query" desc="Flea JPA查询对象池配置">
        <config-item key="pool.maxTotal" desc="Flea JPA查询对象池最大连接数">100</config-item>
        <config-item key="pool.maxIdle" desc="Flea JPA查询对象池最大空闲连接数">10</config-item>
        <config-item key="pool.minIdle" desc="Flea JPA查询对象池最小空闲连接数">0</config-item>
        <config-item key="pool.maxWaitMillis" desc="Flea JPA查询对象池获取连接时的最大等待毫秒数">2000</config-item>
    </config-items>
```

## 2.3 Flea对象池父类

[FleaObjectPool](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/pool/FleaObjectPool.java) 是 **Flea** 框架下定义的通用对象池抽象实现。 `FleaObjectPool` 内嵌了**GenericObjectPool**作为内部对象池实例，用于存储实际的对象；同时它实现 **java.io.Closeable**，用于处理对象池的关闭。

```java
public abstract class FleaObjectPool<T> implements Closeable {

    private static final Logger LOGGER = LoggerFactory.getLogger(FleaObjectPool.class);
    protected GenericObjectPool<T> fleaObjectPool;

    public FleaObjectPool() {
    }

    public FleaObjectPool(final GenericObjectPoolConfig poolConfig, PooledObjectFactory<T> factory) {
        initFleaObjectPool(poolConfig, factory);
    }

    @Override
    public void close() {
        closeFleaObjectPool();
    }

    public boolean isClosed() {
        return fleaObjectPool.isClosed();
    }

    public void initFleaObjectPool(final GenericObjectPoolConfig poolConfig, PooledObjectFactory<T> factory) {
        if (ObjectUtils.isNotEmpty(fleaObjectPool)) {
            closeFleaObjectPool();
        }
        fleaObjectPool = new GenericObjectPool<T>(factory, poolConfig);
    }

    public T getFleaObject() {
        T object = null;
        try {
            object = fleaObjectPool.borrowObject();
        } catch (Exception e) {
        }
        return object;
    }

    protected void returnObject(final T object) {
        if (ObjectUtils.isEmpty(object)) {
            return;
        }
        try {
            fleaObjectPool.returnObject(object);
        } catch (Exception e) {
        }
    }

    protected void returnFleaObject(final T object) {
        if (ObjectUtils.isNotEmpty(object)) {
            returnObject(object);
        }
    }

    protected void closeFleaObjectPool() {
        try {
            if (ObjectUtils.isNotEmpty(fleaObjectPool)) {
                fleaObjectPool.close();
            }
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error("Could not close the pool, Exception :\n", e);
            }
        }
    }
    
    // 其他方法省略。。。
}
```



FleaObjectPool 包含两个构造方法：

- **无参构造方法**： 外部可调用 `initFleaObjectPool` 方法初始化

- **带参数构造方法** ：内部直接调用 `initFleaObjectPool` 方法初始化

其他成员方法，如下：

- `initFleaObjectPool(final GenericObjectPoolConfig poolConfig, PooledObjectFactory<T> factory)` ：初始化对象池
    * `poolConfig` ：对象池配置类
    
    * `factory` ： 池化对象工厂类

- `close()` ：调用 `closeFleaObjectPool()` 方法关闭对象池
- `isClosed()` ：是否对象池实例已经关闭【`true` : 对象池实例已关闭】
- `getFleaObject()` ：从对象池中获取一个对象实例【子类对象池可覆盖该方法，可见[FleaJPAQueryPool](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPool.java)】
- `returnObject(final T object)` ：将对象实例归还给对象池
- `returnFleaObject(final T object)` ：将对象实例归还给对象池【子类对象池可覆盖该方法，可见[FleaJPAQueryPool](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPool.java)】
- `closeFleaObjectPool()` ：关闭对象池

## 2.4 Flea JPA查询对象池 

[FleaJPAQueryPool](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPool.java) 继承自 **FleaObjectPool**， 类型 `T` 初始化为 **FleaJPAQuery**；

```java
public class FleaJPAQueryPool extends FleaObjectPool<FleaJPAQuery> {

    private String poolName;

    public FleaJPAQueryPool(GenericObjectPoolConfig poolConfig) {
        super(poolConfig, new FleaJPAQueryFactory());
    }

    public FleaJPAQueryPool(String poolName, GenericObjectPoolConfig poolConfig) {
        this(poolConfig);
        this.poolName = poolName;
    }

    @Override
    public FleaJPAQuery getFleaObject() {
        FleaJPAQuery query = super.getFleaObject();
        query.setFleaObjectPool(this);
        return query;
    }

    @Override
    protected void returnFleaObject(FleaJPAQuery object) {
        if (ObjectUtils.isNotEmpty(object)) {
            object.reset();
            returnObject(object);
        }
    }

    // 部分代码省略

    private static class FleaJPAQueryFactory implements PooledObjectFactory<FleaJPAQuery> {

        @Override
        public PooledObject<FleaJPAQuery> makeObject() throws Exception {
            // 对象池新增对象
            FleaJPAQuery query = new FleaJPAQuery();
            return new DefaultPooledObject<FleaJPAQuery>(query);
        }

        @Override
        public void destroyObject(PooledObject<FleaJPAQuery> p) throws Exception {
            final FleaJPAQuery query = p.getObject();
            if (ObjectUtils.isNotEmpty(query)) {
                query.reset();
            }
        }

        @Override
        public boolean validateObject(PooledObject<FleaJPAQuery> p) {
            return true;
        }

        @Override
        public void activateObject(PooledObject<FleaJPAQuery> p) throws Exception {

        }

        @Override
        public void passivateObject(PooledObject<FleaJPAQuery> p) throws Exception {

        }
    }
}
```


查看上述代码，可以看到在 `getFleaObject` 和 `returnFleaObject` 方法覆写的父类 `FleaObjectPool` 中的方法，操作的就是实际要池化的 `FleaJPAQuery` 对象。

**Flea JPA** 查询对象池化工厂类实现了 `PooledObjectFactory<FleaJPAQuery>` 接口，它用于为对象池创建和管理 `FleaJPAQuery` 对象的生命周期。

以下是 `FleaJPAQueryFactory` 实现的各个方法：

- `makeObject()` ：该方法用于在对象池中创建并返回一个新的 `FleaJPAQuery` 对象。每当对象池需要一个新的 `FleaJPAQuery` 实例时，就会调用此方法。这里创建了一个新的 `FleaJPAQuery` 实例，并将其包装在 `DefaultPooledObject` 中返回。`DefaultPooledObject` 是 **Apache Commons Pool2** 库中的一个类，用于存储和管理池中的对象。
- `destroyObject(PooledObject<FleaJPAQuery> p)` ：-   当一个 `FleaJPAQuery` 对象从对象池中移除并且不再需要时（例如，当对象池被关闭或达到其最大容量时），会调用此方法。
- `validateObject(PooledObject<FleaJPAQuery> p)` ：该方法用于检查传递给它的 `FleaJPAQuery` 对象是否仍然有效或可用。如果对象无效并且应该从池中移除，则为 `false`；否则为 `true`。
- `activateObject(PooledObject<FleaJPAQuery> p)` ：该方法用于在对象从池中借出之前执行特定的操作。
- `passivateObject(PooledObject<FleaJPAQuery> p)` ：该方法用于在对象归还到池中之后执行特定的操作。
    

## 2.5 Flea对象池构建者接口
[IFleaObjectPoolBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/pool/IFleaObjectPoolBuilder.java) 是 **Flea框架** 下定义的通用的对象池构建者接口，它提供 `build` 方法，用于构建 `Flea` 对象池。

```java
public interface IFleaObjectPoolBuilder {
    FleaObjectPool build(String poolName);
}
```
## 2.6 Flea JPA查询对象池构建者 

[FleaJPAQueryPoolBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPoolBuilder.java) 实现 `Flea` 对象池构建者接口类，用于构建 `Flea JPA` 查询对象池。

```java
public class FleaJPAQueryPoolBuilder implements IFleaObjectPoolBuilder {

    @Override
    public FleaObjectPool build(String poolName) {
        // 获取Flea JPA查询对象池配置
        FleaJPAQueryPoolConfig fleaJPAQueryPoolConfig = FleaJPAQueryPoolConfig.getConfig();
        // 新建 Flea JPA查询对象池
        FleaObjectPool fleaObjectPool = new FleaJPAQueryPool(poolName, fleaJPAQueryPoolConfig);

        return fleaObjectPool;
    }

}
```

## 2.7 Flea对象池工厂类

[FleaObjectPoolFactory](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/pool/FleaObjectPoolFactory.java) 有重载的 `getFleaObjectPool` 方法，分别获取 默认Flea对象池 和 指定对象池名的Flea对象池。

- `getFleaObjectPool(Class<?> objClazz, Class<T> objPoolClazz)` ：获取默认Flea对象池（指定对象Class）
- `getFleaObjectPool(String poolName, Class<?> objClazz, Class<T> objPoolClazz)` ：获取指定对象池名的Flea对象池（指定对象Class）

```java
public class FleaObjectPoolFactory {
    // 存储Flea对象池
    private static final ConcurrentMap<String, FleaObjectPool> fleaObjectPools = new ConcurrentHashMap<String, FleaObjectPool>();

    private static final Object objectPoolLock = new Object();    

    public static <T extends FleaObjectPool> T getFleaObjectPool(Class<?> objClazz, Class<T> objPoolClazz) {
        return getFleaObjectPool(CommonConstants.FleaPoolConstants.DEFAUTL_POOL_NAME, objClazz, objPoolClazz);
    }

    public static <T extends FleaObjectPool> T getFleaObjectPool(String poolName, Class<?> objClazz, Class<T> objPoolClazz) {
        if (StringUtils.isBlank(poolName) || ObjectUtils.isEmpty(objClazz)) {
            return null;
        }
        String poolNameKey = poolName + CommonConstants.SymbolConstants.UNDERLINE + objClazz.getName();
        if (!fleaObjectPools.containsKey(poolNameKey)) {
            synchronized (objectPoolLock) {
                if (!fleaObjectPools.containsKey(poolNameKey)) {
                    fleaObjectPools.put(poolNameKey, build(poolName, objClazz));
                }
            }
        }
        Object objPool = fleaObjectPools.get(poolNameKey);
        if (objPoolClazz.isInstance(objPool)) {
            return objPoolClazz.cast(objPool);
        } else {
            return null;
        }
    }

    private static FleaObjectPool build(String poolName, Class<?> objClazz) {
        String className = objClazz.getSimpleName();
        ConfigItem configItem = FleaConfigManager.getConfigItem(CommonConstants.FleaPoolConstants.FLEA_OBJECT_POOL, className);
        if (ObjectUtils.isEmpty(configItem)) {
            return null;
        }

        String builderImpl = configItem.getValue();
        if (StringUtils.isBlank(builderImpl)) {
            return null;
        }

        FleaObjectPool fleaObjectPool = null;

        IFleaObjectPoolBuilder fleaObjectPoolBuilder = (IFleaObjectPoolBuilder) ReflectUtils.newInstance(builderImpl);
        if (ObjectUtils.isNotEmpty(fleaObjectPoolBuilder)) {
            // 调用指定的类，创建Flea对象池
            fleaObjectPool = fleaObjectPoolBuilder.build(poolName);
        }

        return fleaObjectPool;
    }
}
```

上述 `build` 方法用于初始化Flea对象池创建，读取如下 `flea-config.xml` 的配置，来获取指定类名的 **Flea** 对象池构建者实现。

以下是Flea JPA查询对象池构建者的配置：

```xml
    <config-items key="flea-object-pool" desc="Flea对象池配置">
        <config-item key="FleaJPAQuery" desc="Flea JPA查询对象池构建者">com.huazie.fleaframework.db.jpa.common.FleaJPAQueryPoolBuilder</config-item>
    </config-items>
```

## 2.8 Flea JPA查询对象改造

[FleaJPAQuery](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQuery.java) 需要实现 **java.io.Closeable** 接口，其他基于对象池的改造如下：

- 添加 **Flea JPA** 查询对象池成员变量，可以使用 `setFleaObjectPool` 方法来设置【这里我们可以查看 `FleaJPAQueryPool` 的 `getFleaObject()` 方法】。
- 实现 **java.io.Closeable** 接口，并在 `close()` 方法中将当前对象归还给对象池。
- `reset()` 方法用于重置 **Flea JPA** 查询对象【这里可以查看 `FleaJPAQueryPool` 的 `returnFleaObject()` 方法】。

```java
    protected FleaJPAQueryPool fleaObjectPool;
    
    public FleaJPAQuery() {
    }

    @Override
    public void close() {
        if (ObjectUtils.isNotEmpty(fleaObjectPool)) {
            fleaObjectPool.returnFleaObject(this);
            fleaObjectPool = null;
        }
    }

    public void setFleaObjectPool(FleaJPAQueryPool fleaObjectPool) {
        this.fleaObjectPool = fleaObjectPool;
    }

    public void reset() {
        entityManager = null;
        sourceClazz = null;
        resultClazz = null;
        root = null;
        criteriaBuilder = null;
        criteriaQuery = null;
        if (CollectionUtils.isNotEmpty(predicates)) {
            predicates.clear();
        }
        orders = null;
        groups = null;
    }
```

## 2.9 抽象Flea JPA DAO层实现改造

这块需要优化 `getQuery` 方法，该方法用于获取指定的 **Flea JPA** 查询对象。


```java
    protected FleaJPAQuery getQuery(Class result) {
        String unitName = FleaEntityManager.getPersistenceUnitName(this.getClass().getSuperclass());
        FleaJPAQueryPool pool;
        if (StringUtils.isBlank(unitName)) {
            pool = FleaObjectPoolFactory.getFleaObjectPool(FleaJPAQuery.class, FleaJPAQueryPool.class);
        } else {
            pool = FleaObjectPoolFactory.getFleaObjectPool(unitName, FleaJPAQuery.class, FleaJPAQueryPool.class);
        }

        if (ObjectUtils.isEmpty(pool)) {
            throw new RuntimeException("Can not get a object pool instance");
        }

        FleaJPAQuery query = pool.getFleaObject();
        query.init(getEntityManager(), entityClass, result);
        return query;
    }
```

上述逻辑，我们来简单总结下：

- 首先，根据当前的 **DAO类** 获取它的父类【即持久化单元 **DAO** 层类】中定义的持久化单元名 `unitName`。
- 接着，如果 `unitName` 不为空，则使用持久化单元名 `unitName` 作为对象池名，来获取 `FleaJPAQueryPool`；否则使用默认对象池名"default"，来获取 `FleaJPAQueryPool`。
- 然后，调用 `FleaJPAQueryPool` 的 `getFleaObject()` 来获取 `FleaJPAQuery`；获取 `FleaJPAQuery` 实例后必须调用 `init` 方法对 **Flea JPA** 查询对象进行初始化。
- 最后，返回 `FleaJPAQuery` 实例对象给调用方，开始组装查询语句，操作数据库。


## 2.10 自测

相关自测类可至 **GitHub** 查看 [FleaJPAQueryTest](https://github.com/Huazie/flea-framework/blob/dev/flea-core/src/test/java/com/huazie/fleaframework/db/jpa/FleaJPAQueryTest.java)。

# 总结
基于对象池的 **FleaJPAQuery** 很好地解决了上一版单例模式引出的并发问题；当你想要尽可能少地创建某个对象，同时又要支持并发环境中使用该对象，不妨试试对象池吧。
