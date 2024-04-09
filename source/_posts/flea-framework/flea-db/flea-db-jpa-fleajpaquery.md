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

[《开发框架-Flea》](/categories/开发框架-Flea/) [《flea-db》](/categories/开发框架-Flea/flea-db/)

![](/images/jpa-logo.png)

# 引言
书接上回JPA封装介绍博文，提到 **FleaJPAQuery** 在前一个版本（单例模式）下存在并发的问题，下面首先来分析一下，然后再介绍目前基于对象池的解决方案。
# 1. 问题分析
1. 上个版本 **FleaJPAQuery** 使用单例模式获取，意味着在同一个服务器中DAO层获取的  **FleaJPAQuery** 有且仅有一个。如下所示：

    ```java
        private static volatile FleaJPAQuery query;

        private FleaJPAQuery() {
        }

        /**
         * <p> 获取Flea JPA查询对象 </p>
         *
         * @return Flea JPA查询对象
         * @since 1.0.0
         */
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
        private EntityManager entityManager; // JPA中用于增删改查的持久化接口
        private Class sourceClazz; // 实体类类对象
        private Class resultClazz; // 操作结果类类对象
        private Root root; // 根SQL表达式对象
        private CriteriaBuilder criteriaBuilder; //标准化生成器
        private CriteriaQuery criteriaQuery; // 标准化查询对象
        private List<Predicate> predicates; // Where条件集合
        private List<Order> orders; // 排序集合
        private List<Expression> groups; // 分组集合
        
        /**
         * <p> 创建查询对象 </p>
         *
         * @return 查询对象
         * @throws DaoException 数据操作层异常类
         */
        private TypedQuery createQuery(boolean isSingle) throws DaoException {
            if (ObjectUtils.isEmpty(sourceClazz)) {
                // 查询非法，实体类类对象为空
                throw new DaoException("ERROR-DB-DAO0000000008");
            }
            if (!isSingle) {
                criteriaQuery.select(root);
            }
            if (CollectionUtils.isNotEmpty(predicates)) {
                // 将所有条件用 and 联合起来
                criteriaQuery.where(criteriaBuilder.and(predicates.toArray(new Predicate[0])));
            }
            if (CollectionUtils.isNotEmpty(orders)) {
                // 将order by 添加到查询语句中
                criteriaQuery.orderBy(orders);
            }
            if (CollectionUtils.isNotEmpty(groups)) {
                // 将group by 添加到查询语句中
                criteriaQuery.groupBy(groups);
            }
            return entityManager.createQuery(criteriaQuery);
        }
    ```

3. 在并发下情况下，DAO层操作获取的 **FleaJPAQuery** 始终是一个；因为存在组装查询语句的过程，不同的数据查询操作之间就会相互影响，导致获取的查询结果不符合预期或者获取查询结果报错；另外组装查询语句的过程也可能直接报错；
 4. 基于上面的分析，也就是要每个DAO层操作获取的 **FleaJPAQuery** 之间互不影响，同时又要保证 **FleaJPAQuery** 尽可能少的被新创建；这就需要为**FleaJPAQuery** 维护一个对象池，需要时从对象池中取出来用，用完重置状态后再归还给对象池。
# 2. 方案讲解

```xml
    <!-- 实现对象池化框架 -->
    <dependency>
        <groupId>org.apache.commons</groupId>
        <artifactId>commons-pool2</artifactId>
        <version>2.4.3</version>
    </dependency>
```
## 2.1 Flea对象池配置
```java
/**
 * <p> Flea对象池配置 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class FleaObjectPoolConfig extends GenericObjectPoolConfig {

    /**
     * <p> 无参构造方法，初始化部分默认配置 </p>
     *
     * @since 1.0.0
     */
    public FleaObjectPoolConfig() {
        // 指明连接是否被空闲连接回收器(如果有)进行检验.如果检测失败,则连接将被从池中去除. true：是
        setTestWhileIdle(true);
        // 连接在池中保持空闲而不被空闲连接回收器线程(如果有)回收的最小时间值，单位毫秒
        setMinEvictableIdleTimeMillis(60000);
        // 在空闲连接回收器线程运行期间休眠的时间值,以毫秒为单位. 如果设置为非正数,则不运行空闲连接回收器线程
        setTimeBetweenEvictionRunsMillis(30000);
        // 在每次空闲连接回收器线程(如果有)运行时检查的连接数量
        // 如果 numTestsPerEvictionRun>=0, 则取numTestsPerEvictionRun 和池内的链接数 的较小值 作为每次检测的链接数; Math.min(numTestsPerEvictionRun, this.idleObjects.size())
        // 如果 numTestsPerEvictionRun<0，则每次检查的连接数是检查时池内连接的总数除以这个值的绝对值再向上取整的结果。 (int)Math.ceil((double)this.idleObjects.size() / Math.abs((double)numTestsPerEvictionRun))
        setNumTestsPerEvictionRun(-1);
    }

}
```
## 2.2 Flea JPA查询对象池配置
[FleaJPAQueryPoolConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPoolConfig.java) 继承 **FleaObjectPoolConfig**  。主要初始化以下对象池配置信息：

 1. **maxTotal**  最大连接数
 2. **maxIdle**  最大空闲连接数
 3. **minIdle**  最小空闲连接数
 4. **maxWaitMillis**  获取连接时的最大等待毫秒数

添加flea-config.xml 配置信息如下：

```xml
    <config-items key="flea-jpa-query" desc="Flea JPA查询对象池配置">
        <config-item key="pool.maxTotal" desc="Flea JPA查询对象池最大连接数">100</config-item>
        <config-item key="pool.maxIdle" desc="Flea JPA查询对象池最大空闲连接数">10</config-item>
        <config-item key="pool.minIdle" desc="Flea JPA查询对象池最小空闲连接数">0</config-item>
        <config-item key="pool.maxWaitMillis" desc="Flea JPA查询对象池获取连接时的最大等待毫秒数">2000</config-item>
    </config-items>
```
## 2.3 Flea对象池父类
[FleaObjectPool](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/pool/FleaObjectPool.java) 内嵌了**GenericObjectPool**作为内部对象池实例，用于存储实际的对象；同时它实现 **java.io.Closeable** ，处理对象池的关闭。

```java
/**
 * <p> Flea Object Pool </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public abstract class FleaObjectPool<T> implements Closeable {

    private static final Logger LOGGER = LoggerFactory.getLogger(FleaObjectPool.class);
    protected GenericObjectPool<T> fleaObjectPool; // 内部Flea对象池对象

    /**
     * <p> 外部可调用initFleaObjectPool方法初始化 </p>
     */
    public FleaObjectPool() {
    }

    public FleaObjectPool(final GenericObjectPoolConfig poolConfig, PooledObjectFactory<T> factory) {
        initFleaObjectPool(poolConfig, factory);
    }

    @Override
    public void close() {
        closeFleaObjectPool();
    }

    /**
     * <p> 是否对象池实例已经关闭 </p>
     *
     * @return <code>true</code> : 对象池实例已关闭
     */
    public boolean isClosed() {
        return fleaObjectPool.isClosed();
    }

    /**
     * <p> 初始化对象池 </p>
     *
     * @param poolConfig 对象池配置
     * @param factory    池化对象工厂类
     * @since 1.0.0
     */
    public void initFleaObjectPool(final GenericObjectPoolConfig poolConfig, PooledObjectFactory<T> factory) {
        if (ObjectUtils.isNotEmpty(fleaObjectPool)) {
            closeFleaObjectPool();
        }
        fleaObjectPool = new GenericObjectPool<T>(factory, poolConfig);
    }

    /**
     * <p> 从对象池中获取一个对象实例 </p>
     *
     * @return 池化的对象实例
     * @since 1.0.0
     */
    public T getFleaObject() {
        T object = null;
        try {
            object = fleaObjectPool.borrowObject();
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error("Could not get a object instance from the pool, Exception :\n", e);
            }
        }
        return object;
    }

    /**
     * <p> 将对象实例归还给对象池 </p>
     *
     * @param object 对象实例
     * @since 1.0.0
     */
    protected void returnObject(final T object) {
        if (ObjectUtils.isEmpty(object)) {
            return;
        }
        try {
            fleaObjectPool.returnObject(object);
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error("Could not return the object instance to the pool, Exception :\n", e);
            }
        }
    }

    /**
     * <p> 将对象实例归还给对象池 </p>
     *
     * @param object 对象实例
     * @since 1.0.0
     */
    protected void returnFleaObject(final T object) {
        if (ObjectUtils.isNotEmpty(object)) {
            returnObject(object);
        }
    }

    /**
     * <p> 关闭对象池 </p>
     *
     * @since 1.0.0
     */
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
## 2.4 Flea JPA查询对象池 
[FleaJPAQueryPool](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPool.java) 继承 **FleaObjectPool**， 类型T初始化为 **FleaJPAQuery**；

```java
/**
 * <p> Flea JPA查询对象池 </p>
 * <pre>使用:
 *  // 获取Flea JPA查询对象池实例 （使用默认连接池名"default"即可）
 *  FleaJPAQueryPool pool = FleaObjectPoolFactory.getFleaObjectPool(FleaJPAQuery.class, FleaJPAQueryPool.class);
 *  // 获取Flea JPA查询对象实例
 *  FleaJPAQuery query = pool.getFleaObject();
 * </pre>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class FleaJPAQueryPool extends FleaObjectPool<FleaJPAQuery> {

    private String poolName;

    /**
     * <p> Flea JPA查询对象池构造方法 </p>
     *
     * @param poolConfig 对象池配置
     * @since 1.0.0
     */
    public FleaJPAQueryPool(GenericObjectPoolConfig poolConfig) {
        super(poolConfig, new FleaJPAQueryFactory());
    }

    /**
     * <p> Flea JPA查询对象池构造方法 </p>
     *
     * @param poolConfig 对象池配置
     * @since 1.0.0
     */
    public FleaJPAQueryPool(String poolName, GenericObjectPoolConfig poolConfig) {
        this(poolConfig);
        this.poolName = poolName;
    }

    /**
     * <p> 获取Flea JPA查询对象 </p>
     *
     * @return Flea JPA查询对象
     * @since 1.0.0
     */
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

    /**
     * <p> Flea JPA查询对象池化工厂类 </p>
     *
     * @author huazie
     * @version 1.0.0
     * @since 1.0.0
     */
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
            return false;
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

## 2.5 Flea对象池构建者接口
[IFleaObjectPoolBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/pool/IFleaObjectPoolBuilder.java) 提供 `build` 方法，用于构建 `Flea` 对象池。

```java
/**
 * <p> Flea对象池构建者 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public interface IFleaObjectPoolBuilder {

    /**
     * <p> 构建Flea对象池 </p>
     *
     * @param poolName 对象池名
     * @return Flea对象池实例
     * @since 1.0.0
     */
    FleaObjectPool build(String poolName);

}
```
## 2.6 Flea JPA查询对象池构建者 
[FleaJPAQueryPoolBuilder](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQueryPoolBuilder.java) 实现 `Flea` 对象池构建者接口类，用于构建 `Flea JPA` 查询对象池。
```java
/**
 * <p> Flea JPA查询对象池构建者 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class FleaJPAQueryPoolBuilder implements IFleaObjectPoolBuilder {

    private static final Logger LOGGER = LoggerFactory.getLogger(FleaJPAQueryPoolBuilder.class);

    @Override
    public FleaObjectPool build(String poolName) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaJPAQueryPoolBuilder#build(String) Start");
        }

        // 获取Flea JPA查询对象池配置
        FleaJPAQueryPoolConfig fleaJPAQueryPoolConfig = FleaJPAQueryPoolConfig.getConfig();
        // 新建 Flea JPA查询对象池
        FleaObjectPool fleaObjectPool = new FleaJPAQueryPool(poolName, fleaJPAQueryPoolConfig);

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaJPAQueryPoolBuilder#build(String) FleaJPAQueryPoolConfig = {}", fleaJPAQueryPoolConfig);
            LOGGER.debug("FleaJPAQueryPoolBuilder#build(String) FleaObjectPool = {}", fleaObjectPool);
            LOGGER.debug("FleaJPAQueryPoolBuilder#build(String) End");
        }
        return fleaObjectPool;
    }

}
```
## 2.7 Flea对象池工厂
[FleaObjectPoolFactory](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/pool/FleaObjectPoolFactory.java) 有重载的**getFleaObjectPool** 方法，分别获取 默认Flea对象池 和 指定对象池名的Flea对象池。
```java
/**
 * <p> Flea对象池工厂 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class FleaObjectPoolFactory {

    private static final Logger LOGGER = LoggerFactory.getLogger(FleaObjectPoolFactory.class);
    // 存储Flea对象池
    private static final ConcurrentMap<String, FleaObjectPool> fleaObjectPools = new ConcurrentHashMap<String, FleaObjectPool>();

    /**
     * <p> 默认Flea对象池（指定对象Class） </p>
     *
     * @param objClazz     对象Class
     * @param objPoolClazz 对象池Class
     * @return 默认对象池
     * @since 1.0.0
     */
    public static <T extends FleaObjectPool> T getFleaObjectPool(Class<?> objClazz, Class<T> objPoolClazz) {
        return getFleaObjectPool(CommonConstants.FleaPoolConstants.DEFAUTL_POOL_NAME, objClazz, objPoolClazz);
    }

    /**
     * <p> 指定对象池名的Flea对象池（指定对象Class）</p>
     *
     * @param poolName     对象池名
     * @param objClazz     对象Class
     * @param objPoolClazz 对象池Class
     * @return 指定对象池名的Flea对象池
     * @since 1.0.0
     */
    public static <T extends FleaObjectPool> T getFleaObjectPool(String poolName, Class<?> objClazz, Class<T> objPoolClazz) {
        if (StringUtils.isBlank(poolName) || ObjectUtils.isEmpty(objClazz)) {
            return null;
        }
        String poolNameKey = poolName + CommonConstants.SymbolConstants.UNDERLINE + objClazz.getName();
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaObjectPoolFactory#getFleaObjectPool(String, Class) Pool Name Key = {}", poolNameKey);
        }
        if (!fleaObjectPools.containsKey(poolNameKey)) {
            synchronized (fleaObjectPools) {
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

    /**
     * <p> 初始化Flea对象池创建 </p>
     *
     * @param poolName 对象池名
     * @param objClazz 指定对象Class
     * @since 1.0.0
     */
    private static FleaObjectPool build(String poolName, Class<?> objClazz) {
        String className = objClazz.getSimpleName();
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaObjectPoolFactory#build(String, Class) Object Simple Name = {}", className);
        }
        ConfigItem configItem = FleaConfigManager.getConfigItem(CommonConstants.FleaPoolConstants.FLEA_OBJECT_POOL, className);
        if (ObjectUtils.isEmpty(configItem)) {
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug("FleaObjectPoolFactory#build(String, Class) Can not find the builder by the <config-item key=\"{}\"> from <config-items key=\"{}\"> in [flea-config.xml]",
                        className, CommonConstants.FleaPoolConstants.FLEA_OBJECT_POOL);
            }
            return null;
        }

        String builderImpl = configItem.getValue();
        if (StringUtils.isBlank(builderImpl)) {
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug("FleaObjectPoolFactory#build(String, Class) The builder is empty, found by the <config-item key=\"{}\"> from <config-items key=\"{}\"> in [flea-config.xml]",
                        className, CommonConstants.FleaPoolConstants.FLEA_OBJECT_POOL);
            }
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
添加flea-config.xml 配置信息如下：

```xml
    <config-items key="flea-object-pool" desc="Flea对象池配置">
        <config-item key="FleaJPAQuery" desc="Flea JPA查询对象池构建者">com.huazie.fleaframework.db.jpa.common.FleaJPAQueryPoolBuilder</config-item>
    </config-items>
```

## 2.8 Flea JPA查询对象改造
[FleaJPAQuery](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQuery.java) 实现 **java.io.Closeable** 接口，基于对象池的改造如下：
```java
    protected FleaJPAQueryPool fleaObjectPool; // Flea JPA查询对象池
    
    public FleaJPAQuery() {
    }

    @Override
    public void close() {
        if (ObjectUtils.isNotEmpty(fleaObjectPool)) {
            fleaObjectPool.returnFleaObject(this);
            fleaObjectPool = null;
        }
    }
    
    /**
     * <p> 设置Flea对象池 </p>
     *
     * @param fleaObjectPool Flea JPA查询对象池
     * @since 1.0.0
     */
    public void setFleaObjectPool(FleaJPAQueryPool fleaObjectPool) {
        this.fleaObjectPool = fleaObjectPool;
    }

    /**
     * <p> 重置Flea JPA查询对象 </p>
     *
     * @since 1.0.0
     */
    public void reset() {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaJPAQuery#reset() Start");
            LOGGER.debug("FleaJPAQuery#reset() Before FleaJPAQuery={}", toString());
        }
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
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaJPAQuery#reset() After FleaJPAQuery={}", toString());
            LOGGER.debug("FleaJPAQuery#reset() End");
        }
    }
```
## 2.9 抽象Flea JPA DAO层实现改造

```java
    /**
     * <p> 获取指定的查询对象 </p>
     *
     * @return 自定义Flea JPA查询对象
     * @since 1.0.0
     */
    protected FleaJPAQuery getQuery(Class result) {
        // 获取当前的持久化单元名
        String unitName = FleaEntityManager.getPersistenceUnitName(this.getClass().getSuperclass());
        FleaJPAQueryPool pool;
        if (StringUtils.isBlank(unitName)) {
            // 获取Flea JPA查询对象池 （使用默认对象池名"default"即可）
            pool = FleaObjectPoolFactory.getFleaObjectPool(FleaJPAQuery.class, FleaJPAQueryPool.class);
        } else {
            // 获取Flea JPA查询对象池 （使用持久化单元名unitName作为对象池名）
            pool = FleaObjectPoolFactory.getFleaObjectPool(unitName, FleaJPAQuery.class, FleaJPAQueryPool.class);
        }

        if (ObjectUtils.isEmpty(pool)) {
            throw new RuntimeException("Can not get a object pool instance");
        }

        // 获取Flea JPA查询对象实例
        FleaJPAQuery query = pool.getFleaObject();
        if (LOGGER.isDebugEnabled()) {
            Object obj = new Object() {};
            LOGGER.debug1(obj, "FleaJPAQueryPool = {}", pool);
            LOGGER.debug1(obj, "FleaJPAQuery = {}", query);
        }
        // 获取实例后必须调用该方法,对Flea JPA查询对象进行初始化
        query.init(getEntityManager(), entityClass, result);
        return query;
    }
```
## 2.10 自测
相关自测类可至GitHub查看 [FleaJPAQueryTest](https://github.com/Huazie/flea-framework/blob/dev/flea-core/src/test/java/com/huazie/fleaframework/db/jpa/FleaJPAQueryTest.java)。

# 总结
基于对象池的 **FleaJPAQuery** 很好地解决了上一版单例模式引出的并发问题；当你想要尽可能少地创建某个对象，同时又要支持并发环境中使用该对象，不妨试试对象池吧。
