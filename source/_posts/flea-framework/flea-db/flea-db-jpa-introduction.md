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

[《开发框架-Flea》](/categories/开发框架-Flea/)

![](/images/flea-logo.png)


# 1. 参考
[flea-db使用之封装JPA操作数据库 源代码](https://github.com/Huazie/flea-framework/tree/dev/flea-db)

# 2. 依赖
[mysql-connector-java-5.1.25.jar](https://mvnrepository.com/artifact/mysql/mysql-connector-java/5.1.25)

```xml
<!-- 数据库JDBC连接相关 （MySQL的JDBC驱动）-->
<dependency>
    <groupId>mysql</groupId>
    <artifactId>mysql-connector-java</artifactId>
    <version>5.1.25</version>
</dependency>
```
[eclipselink-2.5.0.jar](https://mvnrepository.com/artifact/org.eclipse.persistence/eclipselink/2.5.0)
```xml
<!-- 数据库持久化相关 EclipseLink-->
<dependency>
    <groupId>org.eclipse.persistence</groupId>
    <artifactId>eclipselink</artifactId>
    <version>2.5.0</version>
</dependency>
```
# 3. 内容讲解
目前支持 JPA + MySQL模式，需要各位本地自行装下MySQL数据库。
## 3.1 Flea JPA查询对象
[FleaJPAQuery](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/common/FleaJPAQuery.java) 用于实现 `JPA` 标准化方式的数据库查询操作，可以自行组装查询条件。下面对一些关键点进行讲解，且听我细细道来 (这一版并发环境下 可能存在问题，后面我会专门写一篇博文讲解 Flea JPA查询对象的问题，其中引入了对象池的概念 )。

 - 获取FleaJPAQuery实例，并初始化内部成员变量
	```java
        private static volatile FleaJPAQuery query;
	    
        private EntityManager entityManager; // JPA中用于增删改查的持久化接口
        private Class sourceClazz; // 实体类类对象
        private Class resultClazz; // 操作结果类类对象
        private Root root; // 根SQL表达式对象
        private CriteriaBuilder criteriaBuilder; //标准化生成器
        private CriteriaQuery criteriaQuery; // 标准化查询对象
        private List<Predicate> predicates; // Where条件集合
        private List<Order> orders; // 排序集合
        private List<Expression> groups; // 分组集合

        private FleaJPAQuery() {
        }

        /**
         * <p> 获取Flea JPA查询对象 </p>
         * （单例模式，本身没有问题，但是由于获取之后Flea JPA查询对象还要使用，
         * 这在有点并发的环境下就存在问题了；后面我会单独写一篇博文讲解基于对象池
         * 的多例模式，既保证并发下各个线程获取的Flea JPA查询对象之间互不影响，
         * 同时也能保证尽可能少的新建Flea JPA查询对象） 
         *
         * @return Flea JPA查询对象
         * @since 1.0.0
         */
        public static FleaJPAQuery getQuery() {
            if (ObjectUtils.isEmpty(query)) {
                synchronized (FleaJPAQuery.class) {
                    if (ObjectUtils.isEmpty(query)) {
                        uery = new FleaJPAQuery();
                    }
                }
            }
            return query;
        }

        /**
         * <p> getQuery()之后，一定要调用该方法进行初始化 </p>
         *
         * @param entityManager JPA中用于增删改查的持久化接口
         * @param sourceClazz   实体类类对象
         * @param resultClazz   操作结果类类对象
         * @since 1.0.0
         */
        public void init(EntityManager entityManager, Class sourceClazz, Class resultClazz) {
            this.entityManager = entityManager;
            this.sourceClazz = sourceClazz;
            this.resultClazz = resultClazz;
            // 从持久化接口中获取标准化生成器
            criteriaBuilder = entityManager.getCriteriaBuilder();
            // 通过标准化生成器 获取 标准化查询对象
            if (ObjectUtils.isEmpty(resultClazz)) {
        	    // 行记录查询结果
                criteriaQuery = criteriaBuilder.createQuery(sourceClazz);
            } else {
        	    // 单个查询结果
                criteriaQuery = criteriaBuilder.createQuery(resultClazz);
            }
            // 通过标准化查询对象，获取根SQL表达式对象
            root = criteriaQuery.from(sourceClazz);
            predicates = new ArrayList<Predicate>();
        }
	```
 - 拼接查询条件，添加排序和分组
 	```java
 		// 等于条件 (单个属性列)
 		public void equal(String attrName, Object value) throws DaoException;
 		// 等于条件 (多个属性列)
 		public void equal(Map<String, Object> paramMap) throws DaoException;
 		// 不等于条件 (单个属性列)
 		public void notEqual(String attrName, Object value) throws DaoException;
 		// 等于条件 (多个属性列)
 		public void notEqual(Map<String, Object> paramMap) throws DaoException;
 		// is null 条件，某属性值为空
 		public void isNull(String attrName) throws DaoException;
 		// is not null 条件，某属性值为非空
 		public void isNotNull(String attrName) throws DaoException;
 		// in 条件， attrName属性的值在value集合中
 		public void in(String attrName, Collection value) throws DaoException;
 		//  not in 条件，attrName属性的值不在value集合中
 		public void notIn(String attrName, Collection value) throws DaoException;
 		// like 条件， 模糊匹配
 		public void like(String attrName, String value) throws DaoException;
 		// 小于等于条件
 		public void le(String attrName, Number value) throws DaoException;
 		// 小于条件
 		public void lt(String attrName, Number value) throws DaoException;
 		// 大于等于条件
 		public void ge(String attrName, Number value) throws DaoException;
 		// 大于条件
 		public void gt(String attrName, Number value) throws DaoException;
 		// between and 条件, 时间区间查询
 		public void between(String attrName, Date startTime, Date endTime) throws DaoException;
 		// 大于某个日期值条件
 		public void greaterThan(String attrName, Date value) throws DaoException;
 		// 大于等于某个日期值条件
 		public void greaterThanOrEqualTo(String attrName, Date value) throws DaoException;
 		// 小于某个日期值条件
 		public void lessThan(String attrName, Date value) throws DaoException;
 		// 小于等于某个日期值条件
 		public void lessThanOrEqualTo(String attrName, Date value) throws DaoException;
 		// 统计数目，在getSingleResult调用之前使用
 		public void count();
 		// 统计数目(带distinct参数)，在getSingleResult调用之前使用
 		public void countDistinct();
 		// 设置查询某属性的最大值，在getSingleResult调用之前使用
 		public void max(String attrName) throws DaoException;
 		// 设置查询某属性的最小值，在getSingleResult调用之前使用
 		public void min(String attrName) throws DaoException;
 		// 设置查询某属性的平均值，在getSingleResult调用之前使用
 		public void avg(String attrName) throws DaoException;
 		// 设置查询某属性的值的总和，在getSingleResult调用之前使用
 		public void sum(String attrName) throws DaoException;
 		// 设置查询某属性的值的总和(Long)，在getSingleResult调用之前使用
 		public void sumAsLong(String attrName) throws DaoException;
 		// 设置查询某属性的值的总和(Double)，在getSingleResult调用之前使用
 		public void sumAsDouble(String attrName) throws DaoException;
 		// 去重某一列
 		public void distinct(String attrName) throws DaoException;
 		// 添加order by子句
 		public void addOrderby(String attrName, String orderBy) throws DaoException;
 		// 添加group by子句
 		public void addGroupBy(String attrName) throws DaoException;
 	```

 - 获取查询结果（记录行 或 单个结果）
   ```java
	// 获取查询的记录行结果集合
	public List getResultList() throws DaoException;
	// 获取查询的记录行结果集合（设置查询范围）
	public List getResultList(int start, int max) throws DaoException;
	// 获取查询的单个属性列结果集合
	// 需要先调用 distinct，否则默认返回行记录结果集合
	public List getSingleResultList() throws DaoException;
	// 获取查询的单个属性列结果集合（设置查询范围，可用于分页）
	// 需要先调用 distinct，否则默认返回行记录结果集合
	public List getSingleResultList(int start, int max) throws DaoException;
	// 获取查询的单个结果
	// 需要提前调用 (count, countDistinct, max, min, avg, sum, sumAsLong, sumAsDouble)
	public Object getSingleResult() throws DaoException;
 	```
## 3.2 数据处理的基本接口
[IFleaJPABaseDataHandler](https://github.com/Huazie/flea-frame/blob/dev/flea-frame-db/src/main/java/com/huazie/frame/db/jpa/common/IFleaJPABaseDataHandler.java) 为基本的数据操作接口，其中包含了查询，（批量）添加，（批量）更新，删除等操作。
## 3.3 抽象Flea JPA DAO层接口
[IAbstractFleaJPADAO](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/dao/interfaces/IAbstractFleaJPADAO.java) 实现了基本的查询、（批量）添加、（批量）更新、删除接口
```java
/**
 * <p> 抽象Flea JPA DAO层接口，实现基本的查询、（批量）添加、（批量）更新、删除接口 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public interface IAbstractFleaJPADAO<T> extends IFleaJPABaseDataHandler<T> {

}
```
## 3.4 抽象Flea JPA DAO层实现
[AbstractFleaJPADAOImpl](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/dao/impl/AbstractFleaJPADAOImpl.java) 中实现上述3中查询、（批量）添加、（批量）更新、删除的接口的具体逻辑。

 - 该类实现上述抽象Flea JPA DAO层接口，同样有类型T，由子类指定其操作的实体类。
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
 - 持久化接口获取，由子类实现（可参考下面的持久化单元DAO层实现）
 	
	```java
    /**
     * 获取实体管理器
     *
     * @return 实体管理器
     * @since 1.0.0
     */
    protected abstract EntityManager getEntityManager();

    /**
     * 获取实体管理器
     *
     * @param entity 实体类对象实例
     * @return 实体管理器类
     * @since 1.0.0
     */
    public EntityManager getEntityManager(T entity) throws CommonException {
        return getEntityManager(entity, false);
    }

    /**
     * 获取实体管理器
     *
     * @param entity 实体类对象实例
     * @param flag   获取实体管理器标识【true：getFleaNextValue获取实体管理器， false: 其他场景获取实体管理器】
     * @return 实体管理器类
     * @throws CommonException 通用异常
     * @since 1.2.0
     */
    private EntityManager getEntityManager(T entity, boolean flag) throws CommonException {
        EntityManager entityManager = getEntityManager();

        // 实体类设置默认库名
        setDefaultLibName(entity);

        // 处理并添加分表信息，如果不存在分表则不处理
        entityManager = LibTableSplitHelper.findTableSplitHandle().handle(entityManager, entity, flag);
        return entityManager;
    }
	```

 - Flea JPA查询对象获取
 
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
## 3.5 定义抽象Flea JPA SV层接口

[IAbstractFleaJPASV](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/service/interfaces/IAbstractFleaJPASV.java) 抽象Flea JPA SV层接口，继承 `IFleaJPABaseDataHandler` 接口，包含了通用的增删改查接口。
```java
/**
 * <p> 抽象Flea JPA SV层接口 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public interface IAbstractFleaJPASV<T> extends IFleaJPABaseDataHandler<T> {
}
```
## 3.6 抽象Flea JPA SV层实现
[AbstractFleaJPASVImpl](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/service/impl/AbstractFleaJPASVImpl.java) 实现上述抽象Flea JPA SV层接口，相关代码也比较简单，具体接口实现内部调用抽象Flea JPA DAO层实现。

```java

    @Override
    public T query(long entityId) throws Exception {
        return getDAO().query(entityId);
    }

    // ... 其他接口实现已省略
	
    /**
     * <p> 获取Flea JPA DAO层实现 </p>
     *
     * @return 抽象Flea JPA DAO层实现
     * @since 1.0.0
     */
    protected abstract IAbstractFleaJPADAO<T> getDAO();
```
## 3.7 持久化单元DAO层实现
[FleaAuthDAOImpl](https://github.com/Huazie/flea-framework/blob/dev/flea-auth/src/main/java/com/huazie/fleaframework/auth/base/FleaAuthDAOImpl.java) 与持久化单元一一对应，如果新增一个持久化配置，即需要新增一个持久化单元DAO层实现，同时Spring配置中，需要加入对应的持久化单元事物管理者配置。


 - 持久化单元名 ----- fleaauth
 - 持久化事物管理者 ----- fleaauthTransactionManager
 - 持久化接口对象 ----- entityManager （该类由注解定义，由Spring配置中的 持久化接口工厂 fleaAuthEntityManagerFactory 初始化，详细可见下面持久化单元相关配置）

```java
/**
 * <p> FleaAuth数据源DAO层父类 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
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

 - 持久化单元配置  **fleaauth-persistence.xml**

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

 - Spring配置
 
    ```xml
    <!-- 持久化单元管理器 -->
    <bean id="defaultPersistenceManager" class="org.springframework.orm.jpa.persistenceunit.DefaultPersistenceUnitManager">
        <property name="persistenceXmlLocations">
            <!-- 可以配置多个持久单元 -->
            <list>
                <value>classpath:META-INF/fleaauth-persistence.xml</value>
            </list>
        </property>
    </bean>
    <!-- 持久化提供者 -->
    <bean id="defaultPersistenceProvider" class="org.eclipse.persistence.jpa.PersistenceProvider"/>
    <!-- 加载时织入器 -->
    <bean id="defaultLoadTimeWeaver" class="org.springframework.instrument.classloading.InstrumentationLoadTimeWeaver"/>
    <!-- JPA厂商适配器，对外公开EclipseLink的持久性提供程序和EntityManager扩展接口  -->
    <bean id="defaultVendorAdapter" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaVendorAdapter">
        <!-- 是否在控制台显示sql -->
        <property name="showSql" value="true"/>
    </bean>
    <!-- JpaDialect EclipseLink持久化服务的实现-->
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
至此，相关JPA使用已封装完毕，下一篇博文将介绍 [《JPA接入》](/2019/09/12/flea-framework/flea-db/flea-db-jpa-integration/) ，敬请期待。