---
title: flea-db使用之JPA封装介绍
date: 2019-09-06 16:24:19
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - JPA封装介绍
---

![](/images/jpa-logo.png)


# 引言
**JPA（Java Persistence API）**,即 **Java** 持久层 **API**，它是 **Java** 平台上用于实现对象关系映射 **（Object-Relational Mapping，简称ORM）** 的规范。它定义了 **Java** 对象如何映射到关系型数据库中的表，并提供了一套标准的 **API** 来管理这些映射关系以及数据库中的持久化对象。

为了方便开发人员后续快速接入 和 使用 JPA 操作数据库，本篇 **Huazie** 将向大家介绍笔者 **Flea** 框架下的 **flea-db** 模块封装JPA操作数据库的内容。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 参考
[flea-db使用之封装JPA操作数据库 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-db)

# 2. 依赖

**MySQL** 的 **JDBC** 驱动 [mysql-connector-java-5.1.25.jar](https://mvnrepository.com/artifact/mysql/mysql-connector-java/5.1.25)

```xml
<dependency>
    <groupId>mysql</groupId>
    <artifactId>mysql-connector-java</artifactId>
    <version>5.1.25</version>
</dependency>
```

**JPA** 实现 **EclipseLink** [eclipselink-2.5.0.jar](https://mvnrepository.com/artifact/org.eclipse.persistence/eclipselink/2.5.0)
```xml
<dependency>
    <groupId>org.eclipse.persistence</groupId>
    <artifactId>eclipselink</artifactId>
    <version>2.5.0</version>
</dependency>
```

# 3. 内容讲解

目前示例用的是 **JPA + MySQL** 模式，需要各位本地自行装下 **MySQL** 数据库。

## 3.1 Flea JPA查询对象
[FleaJPAQuery](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQuery.java) 用于实现 **JPA** 标准化方式的数据库查询操作，可以自行组装查询条件。下面对一些关键点进行讲解，且听我细细道来 (这一版并发环境下 可能存在问题，后面我会专门写一篇博文讲解 Flea JPA查询对象的问题，其中引入了对象池的概念 )。

 - 获取FleaJPAQuery实例，并初始化内部成员变量
    - `EntityManager entityManager` ：**JPA** 中用于增删改查的持久化接口
    - `Class sourceClazz` ： 实体类类对象
    - `Class resultClazz` ： 操作结果类类对象
    - `Root root` ： 根SQL表达式对象
    - `CriteriaBuilder criteriaBuilder` ： 标准化生成器
    - `CriteriaQuery criteriaQuery` ： 标准化查询对象
    - `List<Predicate> predicates` ： Where条件集合
    - `List<Order> orders` ： 排序集合
    - `List<Expression> groups` ： 分组集合
    - `getQuery()` ： 获取Flea JPA查询对象。**新版本已废弃**（单例模式，本身没有问题，但是由于获取之后 **Flea JPA** 查询对象还要使用，这在有点并发的环境下就存在问题了；后面我会单独写一篇博文讲解基于对象池的多例模式，既保证并发下各个线程获取的 **Flea JPA** 查询对象之间互不影响，同时也能保证尽可能少的新建 **Flea JPA** 查询对象）
    - `init(EntityManager entityManager, Class sourceClazz, Class resultClazz)` ：获取 `FleaJPAQuery` 实例之后，一定要调用该方法进行初始化
    - `initQueryEntity(Object entity)` ：初始化查询实体，主要用来构建查询条件值，以及分库分表
    
 - 拼接查询条件，添加排序和分组
    - `equal(String attrName, Object value)` ： 等于条件 (单个属性列)
    - `equal(Map<String, Object> paramMap)` ： 等于条件 (多个属性列)
    - `notEqual(String attrName, Object value)` ： 不等于条件 (单个属性列)
    - `notEqual(Map<String, Object> paramMap)` ： 不等于条件 (多个属性列)
    - `isNull(String attrName)` ： `is null` 条件，某属性值为空
    - `isNotNull(String attrName)` ： `is not null` 条件，某属性值为非空
    - `in(String attrName, Collection value)` ： `in` 条件， `attrName` 属性的值在 `value` 集合中
    - `notIn(String attrName, Collection value)` ： `not in` 条件，`attrName` 属性的值不在 `value` 集合中
    - `like(String attrName, String value)` ： `like` 条件， 模糊匹配
    - `le(String attrName, Number value)` ： 小于等于条件
    - `lt(String attrName, Number value)` ： 小于条件
    - `ge(String attrName, Number value)` ： 大于等于条件
    - `gt(String attrName, Number value)` ： 大于条件
    - `between(String attrName, Date startTime, Date endTime)` ： `between and` 条件, 时间区间查询
    - `greaterThan(String attrName, Date value)` ： 大于某个日期值条件
    - `greaterThanOrEqualTo(String attrName, Date value)` ： 大于等于某个日期值条件
    - `lessThan(String attrName, Date value)` ： 小于某个日期值条件
    - `lessThanOrEqualTo(String attrName, Date value)` ： 小于等于某个日期值条件
    - `count()` ： 统计数目，在 `getSingleResult` 调用之前使用
    - `countDistinct()` ： 统计数目(带 `distinct` 参数)，在 `getSingleResult` 调用之前使用
    - `max(String attrName)` ： 设置查询某属性的最大值，在 `getSingleResult` 调用之前使用
    - `min(String attrName)` ： 设置查询某属性的最小值，在 `getSingleResult` 调用之前使用
    - `avg(String attrName)` ： 设置查询某属性的平均值，在 `getSingleResult` 调用之前使用
    - `sum(String attrName)` ： 设置查询某属性的值的总和，在 `getSingleResult` 调用之前使用
    - `sumAsLong(String attrName)` ： 设置查询某属性的值的总和(Long)，在 `getSingleResult` 调用之前使用
    - `sumAsDouble(String attrName)` ： 设置查询某属性的值的总和(Double)，在 `getSingleResult` 调用之前使用
    - `distinct(String attrName)` ： 去重某一列
    - `addOrderby(String attrName, String orderBy)` ： 添加 `order by` 子句
    - `addGroupBy(String attrName)` ： 添加 `group by` 子句

- 获取查询结果（记录行 或 单个结果）
    - `getResultList()` ： 获取查询的记录行结果集合
    - `getResultList(int start, int max)` ： 获取查询的记录行结果集合（设置查询范围）
    - `getSingleResultList()` ： 获取查询的单个属性列结果集合。需要先调用 `distinct`，否则默认返回行记录结果集合
    - `getSingleResultList(int start, int max)` ： 获取查询的单个属性列结果集合（设置查询范围，可用于分页）。需要先调用 `distinct`，否则默认返回行记录结果集合。
    - `getSingleResult()` ： 获取查询的单个结果。需要提前调用 (`count, countDistinct, max, min, avg, sum, sumAsLong, sumAsDouble`)

## 3.2 数据处理的基本接口
[IFleaJPABaseDataHandler](https://github.com/Huazie/flea-frame/blob/dev/flea-frame-db/src/main/java/com/huazie/frame/db/jpa/common/IFleaJPABaseDataHandler.java) 为基本的数据操作接口，其中包含了查询，（批量）添加，（批量）更新，删除等操作。

## 3.3 抽象Flea JPA DAO层接口
[IAbstractFleaJPADAO](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/dao/interfaces/IAbstractFleaJPADAO.java) 实现了基本的查询、（批量）添加、（批量）更新、删除接口
```java
public interface IAbstractFleaJPADAO<T> extends IFleaJPABaseDataHandler<T> {

}
```

## 3.4 抽象Flea JPA DAO层实现
[AbstractFleaJPADAOImpl](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/dao/impl/AbstractFleaJPADAOImpl.java) 中实现上述3中查询、（批量）添加、（批量）更新、删除的接口的具体逻辑。

 - 该类实现上述抽象 **Flea JPA DAO** 层接口，同样有类型T，由子类指定其操作的实体类。
    ```java
    public abstract class AbstractFleaJPADAOImpl<T> implements IAbstractFleaJPADAO<T> 
    ```
 - 无参构造方法，用于获取子类指定的实体类类对象。

    ```java
    /**
     * <p> 获取T类型的Class对象 </p>
     *
     * @since 1.0.0
     */
    public AbstractFleaJPADAOImpl() {
        // 获取泛型类的子类对象的Class对象
        Class<?> clz = getClass();
        // 获取子类对象的泛型父类类型（也就是AbstractDaoImpl<T>）
        ParameterizedType type = (ParameterizedType) clz.getGenericSuperclass();
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("Type={}", type);
        }
        Type[] types = type.getActualTypeArguments();
        clazz = (Class<T>) types[0];
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("ClassName={}", clazz.getName());
        }
    }
    ```
 - 实现接口方法，可参见上述类源码
 - 持久化接口获取，由子类实现（可参考下面的持久化单元 **DAO** 层实现）
    - `getEntityManager()` ：获取实体管理器
    - `getEntityManager(T entity)` ：获取实体管理器【`entity` 实体类对象实例】
    - `getEntityManager(T entity, boolean flag)` ：获取实体管理器【`entity` 实体类对象实例，flag   获取实体管理器标识【`true`：`getFleaNextValue` 获取实体管理器， `false`: 其他场景获取实体管理器】】
    ```java
    protected abstract EntityManager getEntityManager();

    public EntityManager getEntityManager(T entity) throws CommonException {
        return getEntityManager(entity, false);
    }

    private EntityManager getEntityManager(T entity, boolean flag) throws CommonException {
        EntityManager entityManager = getEntityManager();

        // 实体类设置默认库名
        setDefaultLibName(entity);

        // 处理并添加分表信息，如果不存在分表则不处理
        entityManager = LibTableSplitHelper.findTableSplitHandle().handle(entityManager, entity, flag);
        return entityManager;
    }
    ```

 - **Flea JPA** 查询对象获取【这里已经是使用 **Flea JPA** 查询对象池来获取 `FleaJPAQuery`】
 
    ```java
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
## 3.5 定义抽象Flea JPA SV层接口

[IAbstractFleaJPASV](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/service/interfaces/IAbstractFleaJPASV.java) 抽象 **Flea JPA SV** 层接口，继承 `IFleaJPABaseDataHandler` 接口，包含了通用的增删改查接口。
```java
public interface IAbstractFleaJPASV<T> extends IFleaJPABaseDataHandler<T> {
}
```

## 3.6 抽象Flea JPA SV层实现
[AbstractFleaJPASVImpl](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/service/impl/AbstractFleaJPASVImpl.java) 实现上述抽象 **Flea JPA SV** 层接口，相关代码也比较简单，具体接口实现内部调用抽象 **Flea JPA DAO** 层实现。

```java

    @Override
    public T query(long entityId) throws Exception {
        return getDAO().query(entityId);
    }

    // ... 其他接口实现已省略

    protected abstract IAbstractFleaJPADAO<T> getDAO();
```
## 3.7 持久化单元DAO层实现
[FleaAuthDAOImpl](https://github.com/Huazie/flea-framework/blob/dev/flea-auth/src/main/java/com/huazie/fleaframework/auth/base/FleaAuthDAOImpl.java) 与持久化单元一一对应，如果新增一个持久化配置，即需要新增一个持久化单元 **DAO** 层实现，同时 **Spring** 配置中，需要加入对应的持久化单元事务管理者配置。


 - **持久化单元名** ----- `fleaauth`
 - **持久化事务管理者** ----- `fleaauthTransactionManager`
 - **持久化接口对象** ----- `entityManager` （该类由注解定义，由 **Spring** 配置中的 持久化接口工厂 `fleaAuthEntityManagerFactory` 初始化，详细可见下面持久化单元相关配置）

**FleaAuth数据源DAO层父类**

```java
public class FleaAuthDAOImpl<T> extends AbstractFleaJPADAOImpl<T> {

    @PersistenceContext(unitName="fleaauth")
    protected EntityManager entityManager;

    @Override
    @Transactional("fleaAuthTransactionManager")
    public boolean remove(long entityId) throws Exception {
        return super.remove(entityId);
    }
    
    // 其余代码省略。。。

    @Override
    protected EntityManager getEntityManager() {
        return entityManager;
    }

}
```
## 3.8 配置介绍

### 3.8.1 持久化单元配置  

**fleaauth-persistence.xml**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<persistence version="2.0" xmlns="http://java.sun.com/xml/ns/persistence" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://java.sun.com/xml/ns/persistence http://java.sun.com/xml/ns/persistence/persistence_2_0.xsd">

    <persistence-unit name="fleaauth" transaction-type="RESOURCE_LOCAL">
        <!-- provider -->
        <provider>org.eclipse.persistence.jpa.PersistenceProvider</provider>
        <!-- Connection JDBC -->
        <class>具体实体类全名</class>
        <exclude-unlisted-classes>true</exclude-unlisted-classes>

        <properties>
            <property name="javax.persistence.jdbc.driver" value="com.mysql.jdbc.Driver" />
            <property name="javax.persistence.jdbc.url" value="jdbc:mysql://localhost:3306/fleaauth?useUnicode=true&amp;characterEncoding=UTF-8" />
            <property name="javax.persistence.jdbc.user" value="root" />
            <property name="javax.persistence.jdbc.password" value="root" />
        </properties>
    </persistence-unit>
</persistence>
```

### 3.8.2 Spring配置

- `defaultPersistenceManager` ：持久化单元管理器
- `defaultPersistenceProvider` ：持久化提供者
- `defaultLoadTimeWeaver` ：加载时织入器
- `defaultVendorAdapter` ：**JPA** 厂商适配器，对外公开 **EclipseLink** 的持久性提供程序和EntityManager扩展接口
- `defaultJpaDialect` ：**JpaDialect EclipseLink** 持久化服务的实现
- `fleaAuthEntityManagerFactory` ：**JPA** 实体管理器工厂类
- `fleaAuthTransactionManager` ：**JPA** 事务管理器

```xml
<bean id="defaultPersistenceManager" class="org.springframework.orm.jpa.persistenceunit.DefaultPersistenceUnitManager">
    <property name="persistenceXmlLocations">
        <!-- 可以配置多个持久单元 -->
        <list>
            <value>classpath:META-INF/fleaauth-persistence.xml</value>
        </list>
    </property>
</bean>
<bean id="defaultPersistenceProvider" class="org.eclipse.persistence.jpa.PersistenceProvider"/>
<bean id="defaultLoadTimeWeaver" class="org.springframework.instrument.classloading.InstrumentationLoadTimeWeaver"/>
<bean id="defaultVendorAdapter" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaVendorAdapter">
    <!-- 是否在控制台显示sql -->
    <property name="showSql" value="true"/>
</bean>
<bean id="defaultJpaDialect" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaDialect"/>

<!-- 以下部分 与指定持久化单元一一对应 -->
<!-- ################# -->
<!-- FleaAuth TransAction Manager JPA -->
<!-- ################# -->
<bean id="fleaAuthEntityManagerFactory" class="org.springframework.orm.jpa.LocalContainerEntityManagerFactoryBean">
    <property name="persistenceUnitManager" ref="defaultPersistenceManager"/>
    <property name="persistenceUnitName" value="fleaauth"/>
    <property name="persistenceProvider" ref="defaultPersistenceProvider"/>
    <property name="jpaVendorAdapter" ref="defaultVendorAdapter"/>
    <property name="jpaDialect" ref="defaultJpaDialect"/>
    <property name="jpaPropertyMap">
        <map>
            <entry key="eclipselink.weaving" value="false"/>
        </map>
    </property>
</bean>

<bean id="fleaAuthTransactionManager" class="org.springframework.orm.jpa.JpaTransactionManager">
    <property name="entityManagerFactory" ref="fleaAuthEntityManagerFactory"/>
</bean>

<tx:annotation-driven transaction-manager="fleaAuthTransactionManager"/>
```

# 总结
至此，相关 **JPA** 使用已封装完毕，欢迎大家评论区讨论。下一篇博文将介绍 [《JPA接入》](/2019/09/12/flea-framework/flea-db/flea-db-jpa-integration/) ，向大家演示使用 **JPA** 封装代码来操作数据库，敬请期待！！！