---
title: flea-db使用之JPA分库分表实现
date: 2022-07-08 08:45:00
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - JPA分库分表
---

[《开发框架-Flea》](/categories/开发框架-Flea/) [《flea-db》](/categories/开发框架-Flea/flea-db/)

![](/images/flea-logo.png)

# 引言
在开始本篇的讲解之前，我先来说下之前写过的两篇博文【现在已弃用】：
[flea-frame-db使用之基于EntityManager实现JPA分表的数据库操作【旧】](/2019/10/09/flea-framework/flea-db/flea-frame-db-jpa-entitymanager-split-table/)
[flea-frame-db使用之基于FleaJPAQuery实现JPA分表查询【旧】](/2019/10/02/flea-framework/flea-db/flea-frame-db-jpa-fleajpaquery-split-table/)

这两篇都与分表相关，之所以被弃用，是因为在并发场景中这一版的分表存在问题。虽然并发场景有问题，但与之相关的分表配置、分表实现也确实为本篇的分库分表提供了一些基础能力，这些不能被忽视，将会在本篇中一一介绍。

经过重构之后，目前 **flea-db** 模块的结构如下图所示：

![](flea-db.png)

|模块| 描述  |
|:--|:--|
|flea-db-common  |  分库配置、分表配置、SQL模板配置、异常 和 工具类等代码 |
|flea-db-eclipselink|   基于EclipseLink版的JPA实现而定制化的代码 |
|flea-db-jdbc| 基于 JDBC 开发的通用代码 |
|flea-db-jpa| 基于 JPA 开发的通用代码 |

# 1. 名词解释
| 名词          |  解释                                                                  |
|:---------------|:-----------------------------------------------------------|
| 模板库名      |  用作模板的数据库名                                    |
|模板库持久化单元名| 模板库下的持久化单元名，一般和模板库相同|
|模板库事物名      | 模板库下的事物管理器名 ，分库配置中可查看 `<transaction>` 标签|
| 分库名           | 以模板库名为基础，根据分库规则得到的数据库名 |
|分库持久化单元名| 以模板库持久化单元名为基础，根据分库规则得到的持久化单元名，一般和分库名相同 |
|分库事物名    | 以模板库事物名为基础，根据分库规则得到的事物名 |
| 分库转换      | 以模板库名为基础，根据分库规则得到数据库名的过程|
|分库序列键   | 分库规则中`<split>`标签中 seq 的值，组成分库名表达式的一部分；<br/>如果是分库分表，也对应着分表规则中`<split>`标签中 seq 的值 |
|分库序列值   | 分库序列键对应的值，在分库转换中使用   |
| 模板表名        |用作模板的表名|
| 分表名           | 以模板表名为基础，根据分表规则得到的表名 |
| 分表转换      | 以模板表名为基础，根据分表规则得到表名的过程|


# 2. 配置讲解
## 2.1 分库配置
分库配置文件默认路径：[flea/db/flea-lib-split.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/flea/db/flea-lib-split.xml)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<flea-lib-split>
    <libs>
        <!-- 分库配置
            name : 模板库名
            count: 分库总数
            exp  : 分库名表达式 (模板库名)(分库序列键)
        -->
        <lib name="fleaorder" count="2" exp="(FLEA_LIB_NAME)(SEQ)" desc="flea订单库分库规则">
            <!-- 分库事物配置
                name : 模板事物名
                exp  : 分库事物名表达式 (模板事物名)(分库序列键)
            -->
            <transaction name="fleaOrderTransactionManager" exp="(FLEA_TRANSACTION_NAME)(SEQ)"/>
            <splits>
                <!-- 分库转换实现配置
                    key : 分库转换类型关键字【可查看 LibSplitEnum】
                    seq : 分库序列键【】
                    implClass : 分库转换实现类【可自行定义，需实现com.huazie.fleaframework.db.common.lib.split.ILibSplit】
                    注意：
                    （1）key不为空，implClass可不填
                    （2）key为空，implClass必填
                    （3）key 和 implClass 都不为空，implClass需要和分库转换类型枚举中分库转换实现类对应上
                -->
                <split key="DEC_NUM" seq="SEQ"/>
            </splits>
        </lib>
    </libs>

    <!-- 其他模块分库配置文件引入 -->
    <!--<import resource=""/>-->

</flea-lib-split>
```
分库规则相关实现代码，可以移步 GitHub 查看 [FleaSplitUtils##getSplitLib](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-common/src/main/java/com/huazie/fleaframework/db/common/util/FleaSplitUtils.java)

## 2.2 分表配置
分表配置文件默认路径：[flea/db/flea-table-split.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/flea/db/flea-table-split.xml)
分库分表案例中，实体类中 **@Table** 注解定义的表名，我们可以理解为模板表名；实际的分表，根据模板表名和分表规则确定，后面将慢慢讲解。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<flea-table-split>
    <tables>

        <!-- 分表配置
            name : 分表对应的模板表名
            lib  : 分表对应的模板库名
            exp  : 分表名表达式 (FLEA_TABLE_NAME)_(列名大写)_(列名大写)
        -->
        <table name="order" lib="fleaorder" exp="(FLEA_TABLE_NAME)_(ORDER_ID)" desc="Flea订单信息表分表规则">
            <splits>
                <!-- 分表转换实现配置
                    key    : 分表转换类型关键字【可查看 TableSplitEnum】
                    column : 分表属性列字段名
                    seq    : 分库序列键【若不为空，值需对应flea-lib-split.xml中<split seq="SEQ" />】
                    implClass : 分表转换实现类【可自行定义，需实现com.huazie.fleaframework.db.common.table.split.ITableSplit】
                    注意：
                    （1）key不为空，implClass可不填
                    （2）key为空，implClass必填
                    （3）key 和 implClass 都不为空，implClass需要和分表转换类型枚举中分表转换实现类对应上
                -->
                <split key="ONE" column="order_id" seq="SEQ"/>
            </splits>
        </table>

        <table name="order_attr" lib="fleaorder" exp="(FLEA_TABLE_NAME)_(ORDER_ID)" desc="Flea订单属性表分表规则">
            <splits>
                <split key="ONE" column="order_id" seq="SEQ"/>
            </splits>
        </table>

    </tables>

    <!-- 其他模块分表配置文件引入 -->
    <!--<import resource=""/>-->

</flea-table-split>
```

分表规则相关实现代码，可以移步 **GitHub** 查看 [FleaSplitUtils##getSplitTable](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-common/src/main/java/com/huazie/fleaframework/db/common/util/FleaSplitUtils.java)

## 2.3 JPA持久化单元配置
**JPA持久化单元**，包含了一组实体类的命名配置 和 数据源配置。实际使用中，一个 **JPA持久化单元** 一般对应一个数据库，其中`<properties>`标签指定具体的数据库配置，包含驱动名、地址、用户和密码；`<class>` 标签指定该数据库下的表对应的实体类。`<exclude-unlisted-classes>` 标签，当设置为 true 时，只有列出的类和 **jars** 将被扫描持久类，否则封闭 **jar** 或目录也将被扫描。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<persistence version="2.0" xmlns="http://java.sun.com/xml/ns/persistence" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://java.sun.com/xml/ns/persistence http://java.sun.com/xml/ns/persistence/persistence_2_0.xsd">

    <persistence-unit name="fleaorder" transaction-type="RESOURCE_LOCAL">
        <!-- provider -->
        <provider>org.eclipse.persistence.jpa.PersistenceProvider</provider>
        <!-- Connection JDBC -->
        <class>com.huazie.fleadbtest.jpa.split.entity.Order</class>
        <class>com.huazie.fleadbtest.jpa.split.entity.OrderAttr</class>
        <class>com.huazie.fleadbtest.jpa.split.entity.OldOrder</class>
        <class>com.huazie.fleadbtest.jpa.split.entity.OldOrderAttr</class>
        <exclude-unlisted-classes>true</exclude-unlisted-classes>

        <properties>
            <property name="javax.persistence.jdbc.driver" value="com.mysql.jdbc.Driver" />
            <property name="javax.persistence.jdbc.url"
                value="jdbc:mysql://localhost:3306/fleaorder?useUnicode=true&amp;characterEncoding=UTF-8" />
            <property name="javax.persistence.jdbc.user" value="root" />
            <property name="javax.persistence.jdbc.password" value="root" />
            <!--<property name="eclipselink.ddl-generation" value="create-tables"/> -->
        </properties>
    </persistence-unit>
</persistence>
```
分库场景，模板库和分库都需要有一个对应的持久化单元配置，详见 [接入演示的持久化单元配置](https://github.com/Huazie/flea-db-test/tree/main/flea-config/src/main/resources/META-INF)。

## 2.4 JPA相关Spring Bean配置
首先是JPA固定的Spring Bean配置，可查看 [fleajpabeans-spring.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/spring/db/jpa/fleajpabeans-spring.xml) ，配置如下：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd">

    <bean id="defaultPersistenceManager"
          class="org.springframework.orm.jpa.persistenceunit.DefaultPersistenceUnitManager">
        <property name="persistenceXmlLocations">
            <!-- 可以配置多个持久单元 -->
            <list>
                <value>classpath:META-INF/fleajpa-persistence.xml</value>
                <value>classpath:META-INF/fleaorder-persistence.xml</value>
                <value>classpath:META-INF/fleaorder1-persistence.xml</value>
                <value>classpath:META-INF/fleaorder2-persistence.xml</value>
            </list>
        </property>
    </bean>

    <bean id="defaultPersistenceProvider" class="org.eclipse.persistence.jpa.PersistenceProvider"/>

    <bean id="defaultVendorAdapter" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaVendorAdapter">
        <property name="showSql" value="true"/>
    </bean>

    <bean id="defaultJpaDialect" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaDialect"/>

</beans>
```
与持久化单元对应的 Bean配置，可查看 [fleaorder-spring.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/spring/db/jpa/fleaorder-spring.xml)，配置 如下：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:tx="http://www.springframework.org/schema/tx"
       xsi:schemaLocation="
       http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd
       http://www.springframework.org/schema/tx http://www.springframework.org/schema/tx/spring-tx.xsd">

    <!-- FleaOrder TransAction Manager JPA -->
    <bean id="fleaOrderEntityManagerFactory"
          class="org.springframework.orm.jpa.LocalContainerEntityManagerFactoryBean">
        <property name="persistenceUnitManager" ref="defaultPersistenceManager"/>
        <property name="persistenceUnitName" value="fleaorder"/>
        <property name="persistenceProvider" ref="defaultPersistenceProvider"/>
        <property name="jpaVendorAdapter" ref="defaultVendorAdapter"/>
        <property name="jpaDialect" ref="defaultJpaDialect"/>
        <property name="jpaPropertyMap">
            <map>
                <entry key="eclipselink.weaving" value="false"/>
                <entry key="eclipselink.logging.thread" value="true"/>
            </map>
        </property>
    </bean>

    <bean id="fleaOrderTransactionManager" class="org.springframework.orm.jpa.JpaTransactionManager">
        <property name="entityManagerFactory" ref="fleaOrderEntityManagerFactory"/>
    </bean>

    <tx:annotation-driven transaction-manager="fleaOrderTransactionManager"/>

    <!-- FleaOrder1 TransAction Manager JPA -->
    <!-- 省略 -->

    <!-- FleaOrder2 TransAction Manager JPA -->
    <!-- 省略 -->

</beans>
```

# 3. 实现讲解
## 3.1 Flea自定义事物切面
Flea自定义事物切面 [FleaTransactionalAspect](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/aspect/FleaTransactionalAspect.java)，拦截由自定义事物注解标记的 **Spring注入** 的方法，
实现在方法调用之前开启事物，调用成功后提交事物，出现异常回滚事务。

Flea自定义事物注解主要标记在两类方法上：
- 一类方法是，**AbstractFleaJPADAOImpl** 的子类的增删改方法；这些方法一般在 某某数据源DAO层实现类 中，注解中需要指定事物名。
- 另一类方法是，除了上一类方法的其他 **Spring注入** 的方法上；需要特别注意的是，自定义事物注解上不仅需要指定事物名、而且还需要指定持久化单元名；

如果存在分库的场景，在调用之前，需要设置当前线程下的分库序列值。

```java
    // 设置当前线程下的分库序列值
    FleaLibUtil.setSplitLibSeqValue("SEQ", "123123123");
    // 调用自定义事物注解标记的方法
```
下面我贴出Flea自定义事物切面的代码，如下：

```java
@Aspect
@Component
public class FleaTransactionalAspect {

    private static final String METHOD_NAME_GET_ENTITY_MANAGER = "getEntityManager";

    @Around("@annotation(com.huazie.fleaframework.db.jpa.transaction.FleaTransactional)")
    public Object invokeWithinTransaction(final ProceedingJoinPoint joinPoint) throws CommonException, FleaException, NoSuchMethodException {
        // 获取当前连接点上的方法
        Method method = FleaAspectUtils.getTargetMethod(joinPoint);
        // 获取当前连接点方法上的自定义Flea事物注解上对应的事物名称
        String transactionName = FleaEntityManager.getTransactionName(method);
        // 获取连接点方法签名上的参数列表
        Object[] args = joinPoint.getArgs();
        // 获取标记Flea事物注解的目标对象
        Object tObj = joinPoint.getTarget();

        // 获取最后一个参数【实体对象】
        FleaEntity fleaEntity = null;
        if (ArrayUtils.isNotEmpty(args)) {
            // 从最后一个参数中获取 Flea实体对象
            fleaEntity = getFleaEntityFromLastParam(args);
        }

        EntityManager entityManager;

        // 标记Flea事物注解的目标对象 为 AbstractFleaJPADAOImpl 的子类
        if (ObjectUtils.isNotEmpty(fleaEntity) && tObj instanceof AbstractFleaJPADAOImpl) {
            // 获取实体管理器
            entityManager = (EntityManager) ReflectUtils.invoke(tObj, METHOD_NAME_GET_ENTITY_MANAGER, fleaEntity, Object.class);
            // 获取分表信息
            SplitTable splitTable = fleaEntity.get(DBConstants.LibTableSplitConstants.SPLIT_TABLE, SplitTable.class);
            // 获取分库信息
            SplitLib splitLib = fleaEntity.get(DBConstants.LibTableSplitConstants.SPLIT_LIB, SplitLib.class);
            if (ObjectUtils.isNotEmpty(splitTable)) {
                splitLib = splitTable.getSplitLib();
            }
            // 分库场景
            if (ObjectUtils.isNotEmpty(splitLib) && splitLib.isExistSplitLib()) {
                transactionName = splitLib.getSplitLibTxName();
            }
        } else {
            // 获取当前连接点方法上的自定义Flea事物注解上对应的持久化单元名
            String unitName = FleaEntityManager.getUnitName(method);
            // 获取分库对象
            SplitLib splitLib = FleaSplitUtils.getSplitLib(unitName, FleaLibUtil.getSplitLibSeqValues());
            // 分库场景
            if (splitLib.isExistSplitLib()) {
                transactionName = splitLib.getSplitLibTxName();
                unitName = splitLib.getSplitLibName();
            }
            entityManager = FleaEntityManager.getEntityManager(unitName, transactionName);
        }

        // 根据事物名，获取配置的事物管理者
        PlatformTransactionManager transactionManager = (PlatformTransactionManager) FleaApplicationContext.getBean(transactionName);
        // 事物名【{0}】非法，请检查！
        ObjectUtils.checkEmpty(transactionManager, DaoException.class, "ERROR-DB-DAO0000000015", transactionName);
        // 新建事物模板对象，用于处理事务生命周期和可能的异常
        FleaTransactionTemplate transactionTemplate = new FleaTransactionTemplate(transactionManager, entityManager);
        return transactionTemplate.execute(new TransactionCallback<Object>() {
            @Override
            public Object doInTransaction(TransactionStatus status) {
                try {
                    return joinPoint.proceed();
                }  catch (Throwable throwable) {
                    ExceptionUtils.throwFleaException(FleaDBException.class, "Proceed with the next advice or target method invocation occurs exception : \n", throwable);
                }
                return null;
            }
        });
    }
}
```
在上述代码中，事物名 和 实体管理器 的获取是重点，因Flea自定义事物注解标记在两类不同的方法上，这两者的获取也不一样。通过事物名可直接从Spring配置中获取定义的事物管理器，事物名对应着spring配置中 `transaction-manager` 对应的属性值，详见 2.4中  [fleaorder-spring.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/spring/db/jpa/fleaorder-spring.xml)   。

最后使用 Flea事物模板，来实现标记 `@FleaTransactional`的方法调用之前开启事物，调用成功后提交事物，出现异常回滚事物。


## 3.2 Flea事物模板
Flea事物模板 [FleaTransactionTemplate](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/transaction/FleaTransactionTemplate.java)，参考 **Spring** 的 **TransactionTemplate**，它是简化程序化事务划分和事务异常处理的模板类。其核心方法是 `execute` , 参数是实现事物回调接口的事务代码。此模板收敛了处理事务生命周期和可能的异常的逻辑，因此事物回调接口的实现和调用代码都不需要显式处理事务。

下面将贴出其核心方法 `execute`，如下：

```java
    @Override
    public <T> T execute(TransactionCallback<T> action) throws TransactionException {
        if (this.transactionManager instanceof CallbackPreferringPlatformTransactionManager) {
            return ((CallbackPreferringPlatformTransactionManager) this.transactionManager).execute(this, action);
        } else {
            // 开启Flea自定义事物
            TransactionStatus status = FleaJPASplitHelper.getHandler().getTransaction(this, transactionManager, entityManager);
            T result;
            try {
                result = action.doInTransaction(status);
            } catch (RuntimeException | Error ex) {
                rollbackOnException(status, ex);
                throw ex;
            } catch (Throwable ex) {
                rollbackOnException(status, ex);
                throw new UndeclaredThrowableException(ex, "TransactionCallback threw undeclared checked exception");
            }
            this.transactionManager.commit(status);
            return result;
        }
    }

    /**
     * 执行回滚，正确处理回滚异常。
     */
    private void rollbackOnException(TransactionStatus status, Throwable ex) throws TransactionException {
        try {
            this.transactionManager.rollback(status);
        } catch (TransactionSystemException ex2) {
            ex2.initApplicationException(ex);
            throw ex2;
        } catch (RuntimeException ex2) {
            throw ex2;
        } catch (Error err) {
            throw err;
        }
    }
```

## 3.3 Flea实体管理器
Flea 实体管理器工具类 [FleaEntityManager](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/persistence/FleaEntityManager.java)，提供了获取持久化上下文交互的实体管理器接口、持久化单元名、事物名、分表信息、各持久化上下文交互接口的静态方法【如： **getFleaNextValue**，**find**，**remove**，**merge**，**persist**，**flush**】。

下面我们来看下整体的代码：

```java
public class FleaEntityManager {

    private static final ConcurrentMap<String, EntityManager> entityManagerMap = new ConcurrentHashMap<>();

    private static final ThreadLocal<Map<Object, Object>> resources = new NamedThreadLocal<>("EntityManager resources");

    private FleaEntityManager() {
    }

    /**
     * 获取指定场景下的实体管理类
     */
    public static EntityManager getEntityManager(String unitName, String transactionName) throws CommonException {

        StringUtils.checkBlank(unitName, DaoException.class, "ERROR-DB-DAO0000000002", "libName");
        StringUtils.checkBlank(transactionName, DaoException.class, "ERROR-DB-DAO0000000002", "transactionName");

        if (!entityManagerMap.containsKey(unitName)) {
            synchronized (entityManagerMap) {
                if (!entityManagerMap.containsKey(unitName)) {
                    // 根据事物名，获取配置的事物管理者
                    JpaTransactionManager manger = (JpaTransactionManager) FleaApplicationContext.getBean(transactionName);
                    // 事物名【{0}】非法，请检查！
                    ObjectUtils.checkEmpty(manger, DaoException.class, "ERROR-DB-DAO0000000015", transactionName);
                    // 获取实体管理者工厂类
                    EntityManagerFactory entityManagerFactory = manger.getEntityManagerFactory();
                    // 创建实体管理者
                    EntityManager entityManager = SharedEntityManagerCreator.createSharedEntityManager(entityManagerFactory);
                    entityManagerMap.put(unitName, entityManager);
                }
            }
        }
        return entityManagerMap.get(unitName);
    }

    /**
     * 从指定类的成员变量上，获取持久化单元名称。在 <b> flea-db </b> 模块中，
     * 该名称一般定义在 {@code AbstractFleaJPADAOImpl} 的子类的成员变量上，由 注解
     * {@code PersistenceContext} 或 注解 {@code FleaPersistenceContext} 进行标识。
     */
    public static String getPersistenceUnitName(Class<?> daoImplClazz) {
        // 省略。。
    }

    /**
     * 从指定类的第一个成员方法上，获取事物名。在 <b> flea-db </b> 模块中，
     * 该名称一般定义在 {@code AbstractFleaJPADAOImpl} 的子类的成员方法上，
     * 由注解 {@code Transactional}或{@code FleaTransactional} 进行标识。
     */
    public static String getTransactionName(Class<?> daoImplClazz) {
        // 省略。。
    }

    /**
     * 从指定类的成员方法上，获取事物名。在 <b> flea-db </b> 模块中，
     * 该名称一般定义在 {@code AbstractFleaJPADAOImpl} 的子类的成员方法上，
     * 由注解 {@code Transactional}或{@code FleaTransactional} 进行标识。
     */
    public static String getTransactionName(Method method) {
        // 省略。。
    }

    /**
     * 从指定类的成员方法上，获取持久化单元名。在 <b> flea-db </b> 模块中，
     * 该名称定义在注解{@code FleaTransactional} 中，用于启动自定的事物。
     */
    public static String getUnitName(Method method) {
        // 省略。。
    }

    /**
     * 根据实体对象，获取实体对应的分表信息
     */
    public static SplitTable getSplitTable(Object entity) throws CommonException {
        // 省略。。
    }

    /**
     * 返回绑定到当前线程的所有资源
     */
    public static Map<Object, Object> getResourceMap() {
        // 省略。。
    }

    /**
     * 检查是否存在绑定到当前线程的给定键的资源。
     */
    public static boolean hasResource(Object key) {
        // 省略。。
    }

    /**
     * 检索绑定到当前线程的给定键的资源。
     */
    public static Object getResource(Object key) {
        // 省略。。
    }

    /**
     * 实际检查给定键绑定的资源的值。
     */
    private static Object doGetResource(Object actualKey) {
        // 省略。。
    }

    /**
     * 将给定键的给定资源绑定到当前线程。
     */
    public static void bindResource(Object key, Object value) throws IllegalStateException {
        // 省略。。
    }

    /**
     * 从当前线程解除给定键的资源绑定。
     */
    public static Object unbindResource(Object key) throws IllegalStateException {
        // 省略。。。
    }

    /**
     * 从当前线程解除给定键的资源绑定。
     */
    public static Object unbindResourceIfPossible(Object key) {
        return doUnbindResource(key);
    }

    /**
     * 实际删除为给定键绑定的资源的值。
     */
    private static Object doUnbindResource(Object actualKey) {
        // 省略。。
    }

    /**
     * 获取下一个主键值
     */
    public static <T> Number getFleaNextValue(EntityManager entityManager, Class<T> entityClass, T entity) {
        return FleaJPASplitHelper.getHandler().getNextValue(entityManager, entityClass, entity);
    }

    /**
     * 根据主键查找表数据
     */
    public static <T> T find(EntityManager entityManager, Object primaryKey, Class<T> entityClass, T entity) {
        return FleaJPASplitHelper.getHandler().find(entityManager, primaryKey, entityClass, entity);
    }

    /**
     * 删除实体类对应的一条数据
     */
    public static <T> boolean remove(EntityManager entityManager, T entity) {
        return FleaJPASplitHelper.getHandler().remove(entityManager, entity);
    }

    /**
     * 将给定实体的状态合并（即更新）到当前持久化上下文中。
     * <p> 注意：调用该方法后，待修改的数据还未更新到数据库中。
     */
    public static <T> T merge(EntityManager entityManager, T entity) {
        return FleaJPASplitHelper.getHandler().merge(entityManager, entity);
    }

    /**
     * 将实体类添加到持久化上下文中，并管理该实体类
     * <p> 注意：调用该方法后，待保存的数据还未添加到数据库中。
     */
    public static <T> void persist(EntityManager entityManager, T entity) {
        FleaJPASplitHelper.getHandler().persist(entityManager, entity);
    }

    /**
     * 将持久化上下文同步到底层数据库。
     */
    public static <T> void flush(EntityManager entityManager, T entity) {
        FleaJPASplitHelper.getHandler().flush(entityManager, entity);
    }
}

```

## 3.4 Flea JPA分库分表处理接口
 [IFleaJPASplitHandler](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/handler/IFleaJPASplitHandler.java)，从上面 3.3中，我们可以看到  Flea实体管理器中的各持久化上下文交互接口的静态方法【如： **getFleaNextValue**，**find**，**remove**，**merge**，**persist**，**flush**】都是调用 [FleaJPASplitHelper.getHandler()](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/FleaJPASplitHelper.java) 的对应方法实现的，也就是 **IFleaJPASplitHandler** 的对应方法。

Flea JPA 分库分表处理者接口，包含分库分表相关的处理接口方法、增删改查的数据操作接口方法。

**下面我们来看看 Flea JPA分库分表处理接口 都有哪些处理方法？**

```java
public interface IFleaJPASplitHandler {

    /**
     * 使用标准化查询时，处理分库分表信息
     */
    void handle(FleaJPAQuery query, Object entity) throws CommonException;

    /**
     * 使用标准化查询时，存在分表场景，具体的JPA查询对象重新设置持久化信息
     */
    void handle(FleaJPAQuery query, TypedQuery typedQuery) throws CommonException;

    /**
     * 使用持久化接口时，处理分库分表信息
     */
    EntityManager handle(EntityManager entityManager, Object entity, boolean flag) throws CommonException;

    /**
     * 分表场景下，取事物管理器中的实体管理器工厂类，并将其作为键，
     * 绑定实体管理器对应的包装类资源到当前线程。以支持JPA的增删改操作。
     */
    TransactionStatus getTransaction(TransactionDefinition definition, PlatformTransactionManager transactionManager, EntityManager entityManager);

    /**
     * 获取下一个主键值
     */
    <T> Number getNextValue(EntityManager entityManager, Class<T> entityClass, T entity);

    /**
     * 根据主键查找表数据
     */
    <T> T find(EntityManager entityManager, Object primaryKey, Class<T> entityClass, T entity);

    /**
     * 删除给定的实体数据
     */
    <T> boolean remove(EntityManager entityManager, T entity);

    /**
     * 将给定实体的状态合并（即更新）到当前持久化上下文中。
     * <p> 注意：调用该方法后，待修改的数据还未更新到数据库中。
     */
    <T> T merge(EntityManager entityManager, T entity);

    /**
     * 将实体类添加到持久化上下文中，并管理该实体类
     * <p> 注意：调用该方法后，待保存的数据还未添加到数据库中。
     */
    <T> void persist(EntityManager entityManager, T entity);

    /**
     * 将持久化上下文同步到底层数据库。
     */
    <T> void flush(EntityManager entityManager, T entity);
}
```

## 3.5 EclipseLink分库分表处理实现
**EclipseLink** 分库分表处理者 [EclipseLinkLibTableSplitHandler](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-eclipselink/src/main/java/com/huazie/fleaframework/db/eclipselink/EclipseLinkLibTableSplitHandler.java)，由自定义的实体管理器实现类处理增删改查等操作。

在讲解 **EclipseLink** 分库分表处理者之前，我们先了解下其父类 [FleaLibTableSplitHandler](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-jpa/src/main/java/com/huazie/fleaframework/db/jpa/handler/impl/FleaLibTableSplitHandler.java)，该类实现了通用的分库分表处理 和 增删改查操作，同时定义了抽象的内部方法由子类实现具体的操作。

下面我们来看一下具体的代码，如下：

```java
public abstract class FleaLibTableSplitHandler implements IFleaJPASplitHandler {

    @Override
    public void handle(FleaJPAQuery query, Object entity) throws CommonException {

        if (ObjectUtils.isEmpty(query) || ObjectUtils.isEmpty(entity) || !(entity instanceof FleaEntity)) {
            return;
        }

        FleaEntity fleaEntity = (FleaEntity) entity;

        // 获取分表信息（包括模板表名 和 分表名 【如果存在分表返回】）
        SplitTable splitTable = FleaEntityManager.getSplitTable(entity);

        SplitLib splitLib;
        // 存在分表，需要查询指定分表
        if (splitTable.isExistSplitTable()) {
            splitLib = splitTable.getSplitLib();
            // 设置分表信息
            fleaEntity.put(DBConstants.LibTableSplitConstants.SPLIT_TABLE, splitTable);
        } else {
            // 获取默认库名，这里的对象池名为持久化单元名【通常对应着库名】
            String libName = query.getPoolName();
            if (ObjectUtils.isEmpty(libName)) {
                return;
            }
            splitLib = FleaSplitUtils.getSplitLib(libName, FleaLibUtil.getSplitLibSeqValues());
        }

        // 分库场景，重新获取对应分库下的实体管理类
        EntityManager splitEntityManager = handleInner(splitLib);

        EntityManager entityManager;
        if (ObjectUtils.isEmpty(splitEntityManager)) {
            entityManager = query.getEntityManager();
        } else {
            entityManager = splitEntityManager;
        }

        // 分表场景 或 分表场景 或 当前线程存在自定义的Flea实体管理器实现, 直接获取
        if (isFleaEntityManagerImpl(entityManager, splitTable, splitLib)) {
            splitEntityManager = getFleaEntityMangerImpl(entityManager);
        }

        if (ObjectUtils.isNotEmpty(splitEntityManager)) {
            // 重新初始化Flea JPA查询对象
            query.init(splitEntityManager, query.getSourceClazz(), query.getResultClazz());
        }
    }

    @Override
    public void handle(FleaJPAQuery query, TypedQuery typedQuery) {
        SplitTable splitTable = getSplitTableFromEntity(query.getEntity());
        // 分表信息为空或不存在分表，默认不处理
        if (!splitTable.isExistSplitTable()) {
            return;
        }
        // 处理类型查询接口的分表信息
        handleInner(query, typedQuery, splitTable);
    }

    @Override
    public EntityManager handle(EntityManager entityManager, Object entity, boolean flag) throws CommonException {

        if (ObjectUtils.isEmpty(entityManager) || ObjectUtils.isEmpty(entity) || !(entity instanceof FleaEntity)) {
            return entityManager;
        }

        // 获取分表信息（包括模板表名 和 分表名 【如果存在分表返回】）
        SplitTable splitTable = FleaEntityManager.getSplitTable(entity);

        FleaEntity fleaEntity = (FleaEntity) entity;

        SplitLib splitLib;
        // 存在分表，则需要操作具体分表
        if (splitTable.isExistSplitTable()) {
            splitLib = splitTable.getSplitLib();
            // 设置分表信息
            fleaEntity.put(DBConstants.LibTableSplitConstants.SPLIT_TABLE, splitTable);
        } else {
            // 获取默认库名
            String libName = fleaEntity.get(DBConstants.LibTableSplitConstants.FLEA_LIB_NAME, String.class);
            if (ObjectUtils.isEmpty(libName)) {
                return entityManager;
            }
            splitLib = FleaSplitUtils.getSplitLib(libName, FleaLibUtil.getSplitLibSeqValues());
            // 设置分库信息
            fleaEntity.put(DBConstants.LibTableSplitConstants.SPLIT_LIB, splitLib);
        }

        // 如果是getFleaNextValue获取实体管理器，并且主键生成器表在模板库中，直接返回实体管理器
        if (flag && splitTable.isGeneratorFlag()) {
            return entityManager;
        }

        // 分库场景，重新获取对应分库下的实体管理类
        EntityManager splitEntityManager = handleInner(splitLib);
        if (ObjectUtils.isNotEmpty(splitEntityManager)) {
            // 分库场景，重新初始化实体管理类
            entityManager = splitEntityManager;
        }
        return entityManager;
    }

    /**
     * 分库场景，重新获取对应分库下的实体管理类
     */
    private EntityManager handleInner(SplitLib splitLib) throws CommonException {
        EntityManager entityManager = null;
        if (ObjectUtils.isNotEmpty(splitLib) && splitLib.isExistSplitLib()) {
            String unitName = splitLib.getSplitLibName();
            String transactionName = splitLib.getSplitLibTxName();
            entityManager = FleaEntityManager.getEntityManager(unitName, transactionName);
        }
        return entityManager;
    }

    @Override
    public TransactionStatus getTransaction(TransactionDefinition definition, PlatformTransactionManager transactionManager, EntityManager entityManager) {
        // JPA事物管理器
        JpaTransactionManager jpaTransactionManager = (JpaTransactionManager) transactionManager;
        Object obj = TransactionSynchronizationManager.getResource(jpaTransactionManager.getEntityManagerFactory());
        if (ObjectUtils.isEmpty(obj)) {
            // 获取Flea实体管理器实现类
            EntityManager fleaEntityManagerImpl = getFleaEntityMangerImpl(entityManager);
            // 新建实体管理器包装类资源，持有Flea实体管理器实现类
            EntityManagerHolder entityManagerHolder = new EntityManagerHolder(fleaEntityManagerImpl);
            // 将实体管理器工厂类的实体管理器包装类资源绑定到当前线程
            TransactionSynchronizationManager.bindResource(jpaTransactionManager.getEntityManagerFactory(), entityManagerHolder);
        }
        // 获取事物状态对象，并开启事物
        return jpaTransactionManager.getTransaction(definition);
    }

    @Override
    public <T> Number getNextValue(EntityManager entityManager, Class<T> entityClass, T entity) {
        SplitTable splitTable = getSplitTableFromEntity(entity);
        return getNextValueInner(entityManager, entityClass, splitTable);
    }

    @Override
    public <T> T find(EntityManager entityManager, Object primaryKey, Class<T> entityClass, T entity) {
        SplitTable splitTable = getSplitTableFromEntity(entity);
        SplitLib splitLib = getSplitLibFromEntity(entity);

        // 分表场景 或 分表场景 或 当前线程存在自定义的Flea实体管理器实现, 直接获取
        T t;
        if (isFleaEntityManagerImpl(entityManager, splitTable, splitLib)) {
            t = findInner(entityManager, primaryKey, entityClass, splitTable);
        } else {
            t =  entityManager.find(entityClass, primaryKey);
        }
        return t;
    }

    @Override
    public <T> boolean remove(EntityManager entityManager, T entity) {
        SplitTable splitTable = getSplitTableFromEntity(entity);
        SplitLib splitLib = getSplitLibFromEntity(entity);

        // 分表场景 或 分表场景 或 当前线程存在自定义的Flea实体管理器实现, 直接获取
        if (isFleaEntityManagerImpl(entityManager, splitTable, splitLib)) {
            // 使用自定义的Flea实体管理器实现，删除实体数据
            removeInner(entityManager, entity);
        } else {
            if (!entityManager.contains(entity)) {
                entity = registerObject(entityManager, entity);
            }
            entityManager.remove(entity);
        }
        return true;
    }

    @Override
    public <T> T merge(EntityManager entityManager, T entity) {
        SplitTable splitTable = getSplitTableFromEntity(entity);
        SplitLib splitLib = getSplitLibFromEntity(entity);

        // 分表场景 或 分表场景 或 当前线程存在自定义的Flea实体管理器实现, 直接获取
        if (isFleaEntityManagerImpl(entityManager, splitTable, splitLib)) {
            // 使用自定义的Flea实体管理器实现，合并实体数据的状态至当前持久化上下文中
            return mergeInner(entityManager, entity);
        } else {
            return entityManager.merge(entity);
        }
    }

    @Override
    public <T> void persist(EntityManager entityManager, T entity) {
        SplitTable splitTable = getSplitTableFromEntity(entity);
        SplitLib splitLib = getSplitLibFromEntity(entity);

        // 分表场景 或 分表场景 或 当前线程存在自定义的Flea实体管理器实现, 直接获取
        if (isFleaEntityManagerImpl(entityManager, splitTable, splitLib)) {
            // 使用自定义的Flea实体管理器实现，向工作单元注册对象
            persistInner(entityManager, entity);
        } else {
            entityManager.persist(entity);
        }
    }

    @Override
    public <T> void flush(EntityManager entityManager, T entity) {
        SplitTable splitTable = getSplitTableFromEntity(entity);
        SplitLib splitLib = getSplitLibFromEntity(entity);

        // 分表场景 或 分表场景 或 当前线程存在自定义的Flea实体管理器实现, 直接获取
        if (isFleaEntityManagerImpl(entityManager, splitTable, splitLib)) {
            // 使用自定义的Flea实体管理器实现，将持久性上下文同步到基础数据库。
            flushInner(entityManager);
        } else {
            entityManager.flush();
        }
    }

    /**
     * 是否使用自定义的Flea实体管理器实现
     */
    private boolean isFleaEntityManagerImpl(EntityManager entityManager, SplitTable splitTable, SplitLib splitLib) {
        return (splitTable.isExistSplitTable() || splitLib.isExistSplitLib() || FleaEntityManager.hasResource(entityManager.getEntityManagerFactory()));
    }

    /**
     * 从实体类对象中获取分表信息
     */
    private SplitTable getSplitTableFromEntity(Object entity) {
        SplitTable splitTable = null;
        if (ObjectUtils.isNotEmpty(entity) && (entity instanceof FleaEntity)) {
            // 获取分表信息
            FleaEntity fleaEntity = (FleaEntity) entity;
            splitTable = fleaEntity.get(DBConstants.LibTableSplitConstants.SPLIT_TABLE, SplitTable.class);
        }
        if (ObjectUtils.isEmpty(splitTable)) {
            splitTable = new SplitTable();
            splitTable.setExistSplitTable(false);
        }
        return splitTable;
    }

    /**
     * 从实体类对象中获取分库信息
     */
    private SplitLib getSplitLibFromEntity(Object entity) {
        SplitLib splitLib = null;
        if (ObjectUtils.isNotEmpty(entity) && (entity instanceof FleaEntity)) {
            // 获取分库信息
            FleaEntity fleaEntity = (FleaEntity) entity;
            splitLib = fleaEntity.get(DBConstants.LibTableSplitConstants.SPLIT_LIB, SplitLib.class);
        }
        if (ObjectUtils.isEmpty(splitLib)) {
            splitLib = new SplitLib();
            splitLib.setExistSplitLib(false);
        }
        return splitLib;
    }

    /**
     * 获取自定义的Flea实体管理器实现
     */
    protected abstract EntityManager getFleaEntityMangerImpl(EntityManager entityManager);

    /**
     * 处理类型查询接口的分表信息
     */
    protected abstract void handleInner(FleaJPAQuery query, TypedQuery typedQuery, SplitTable splitTable);

    /**
     * 自定义的实体管理器实现，获取下一个主键值
     */
    protected abstract <T> Number getNextValueInner(EntityManager entityManager, Class<T> entityClass, SplitTable splitTable);

    /**
     * 使用自定义的实体管理器实现，根据主键查询实体数据
     */
    protected abstract <T> T findInner(EntityManager entityManager, Object primaryKey, Class<T> entityClass, SplitTable splitTable);

    /**
     * 使用自定义的实体管理器实现，删除实体数据
     */
    protected abstract <T> void removeInner(EntityManager entityManager, T entity);

    /**
     * 使用自定义的实体管理器实现，合并实体数据的状态至当前持久化上下文中
     */
    protected abstract <T> T mergeInner(EntityManager entityManager, T entity);

    /**
     * 使用自定义的实体管理器实现，向工作单元注册对象
     */
    protected abstract <T> void persistInner(EntityManager entityManager, T entity);

    /**
     * 使用自定义的实体管理器实现，将持久化上下文同步到底层数据库。
     */
    protected abstract void flushInner(EntityManager entityManager);

    /**
     * 注册实体对象
     */
    protected abstract <T> T registerObject(EntityManager entityManager, T entity);
}
```

好，上面已经基本实现分表处理者的各项接口方法，剩下一些inner方法，需要由特定的JPA实现来定制化，本例中是EclipseLink。

下面我们来看看相关代码，如下：

```java
public class EclipseLinkLibTableSplitHandler extends FleaLibTableSplitHandler {

    @Override
    protected void handleInner(FleaJPAQuery query, TypedQuery typedQuery, SplitTable splitTable) {
        // 获取实体类型
        EntityType entityType = query.getRoot().getModel();
        // 获取实体类对应的持久化信息
        ClassDescriptor classDescriptor = ((EntityTypeImpl) entityType).getDescriptor();
        // 分表场景，这里的entityManager已经重新设置为 FleaEntityManagerImpl
        AbstractSession abstractSession = ((FleaEntityManagerImpl) query.getEntityManager()).getDatabaseSession();
        classDescriptor = ClassDescriptorUtils.getSplitDescriptor(classDescriptor, abstractSession, splitTable);
        // 获取内部DatabaseQuery对象
        ReadAllQuery readAllQuery = typedQuery.unwrap(ReadAllQuery.class);
        // 重新设置实体类的描述符信息
        readAllQuery.setDescriptor(classDescriptor);
        // 重新设置实体类的描述符信息
        readAllQuery.getExpressionBuilder().setQueryClassAndDescriptor(classDescriptor.getJavaClass(), classDescriptor);
    }

    @Override
    protected EntityManager getFleaEntityMangerImpl(EntityManager entityManager) {
        return FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager);
    }

    @Override
    protected <T> Number getNextValueInner(EntityManager entityManager, Class<T> entityClass, SplitTable splitTable) {
        return FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager).getNextValue(entityClass, splitTable);
    }

    @Override
    protected <T> T findInner(EntityManager entityManager, Object primaryKey, Class<T> entityClass, SplitTable splitTable) {
        return FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager).find(entityClass, primaryKey, splitTable);
    }

    @Override
    protected <T> void removeInner(EntityManager entityManager, T entity) {
        FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager).remove(entity);
    }

    @Override
    protected <T> T mergeInner(EntityManager entityManager, T entity) {
        return FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager).merge(entity);
    }

    @Override
    protected <T> void persistInner(EntityManager entityManager, T entity) {
        FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager).persist(entity);
    }

    @Override
    protected void flushInner(EntityManager entityManager) {
        FleaEntityManagerImpl.getFleaEntityManagerImpl(entityManager).flush();
    }

    @Override
    @SuppressWarnings("unchecked")
    protected <T> T registerObject(EntityManager entityManager, T entity) {
        // 如果已经注册过了，直接返回待注册对象
        if (entityManager.contains(entity) || ObjectUtils.isEmpty(entity)) {
            return entity;
        }

        return (T) entityManager.unwrap(UnitOfWork.class).registerObject(entity);
    }
}

```
上面具体的 **inner** 方法实现，我们可以看到都使用了 **FleaEntityManagerImpl** ，这就是下面将要介绍的 **Flea** 实体管理器 **EclipseLink** 版实现。

## 3.6 Flea实体管理器EclipseLink版实现
**Flea** 实体管理器 **EclipseLink** 版实现 [FleaEntityManagerImpl](https://github.com/Huazie/flea-framework/blob/main/flea-db/flea-db-eclipselink/src/main/java/org/eclipse/persistence/internal/jpa/FleaEntityManagerImpl.java)，继承了 **EclipseLink** 的 **EntityManagerImpl**，它需要有 一个 **EntityManager** 入参来构造。

下面我们来看一下相关代码，如下：

```java
public final class FleaEntityManagerImpl extends EntityManagerImpl {

    /**
     * 获取指定JPA实体管理器工厂类对应的自定义的Flea实体管理器实现
     */
    public static FleaEntityManagerImpl getFleaEntityManagerImpl(EntityManager entityManager) {
        EntityManagerFactory entityManagerFactory = entityManager.getEntityManagerFactory();
        FleaEntityManagerImpl fleaEntityManagerImpl = (FleaEntityManagerImpl) FleaEntityManager.getResource(entityManagerFactory);
        if (ObjectUtils.isEmpty(fleaEntityManagerImpl)) {
            fleaEntityManagerImpl = new FleaEntityManagerImpl(entityManager);
            FleaEntityManager.bindResource(entityManagerFactory, fleaEntityManagerImpl);
        }
        return fleaEntityManagerImpl;
    }

    /**
     * 通过EntityManagerImpl构建FleaEntityManagerImpl
     */
    private FleaEntityManagerImpl(EntityManager entityManager) {
        super(entityManager.getEntityManagerFactory().unwrap(JpaEntityManagerFactory.class).unwrap(), entityManager.getProperties(), null);
    }

    /**
     * 分表场景下，根据主键查找实体数据
     */
    public <T> T find(Class<T> entityClass, Object primaryKey, SplitTable splitTable) {
        return find(entityClass, primaryKey, null, getQueryHints(entityClass, OperationType.FIND), splitTable);
    }

    /**
     * 分表场景下，根据主键查找实体数据
     */
    public <T> T find(Class<T> entityClass, Object primaryKey, Map<String, Object> properties, SplitTable splitTable) {
        return find(entityClass, primaryKey, null, properties, splitTable);
    }

    /**
     * 分表场景下，根据主键查找实体数据
     */
    public <T> T find(Class<T> entityClass, Object primaryKey, LockModeType lockMode, SplitTable splitTable) {
        return find(entityClass, primaryKey, lockMode, getQueryHints(entityClass, OperationType.FIND), splitTable);
    }

    /**
     * 分表场景下，根据主键查找实体数据
     */
    @SuppressWarnings({"unchecked"})
    public <T> T find(Class<T> entityClass, Object primaryKey, LockModeType lockMode, Map<String, Object> properties, SplitTable splitTable) {
        try {
            verifyOpen();
            if (ObjectUtils.isNotEmpty(lockMode) && !lockMode.equals(LockModeType.NONE)) {
                checkForTransaction(true);
            }
            AbstractSession session = this.databaseSession;
            ClassDescriptor descriptor = session.getDescriptor(entityClass);
            if (descriptor == null || descriptor.isDescriptorTypeAggregate()) {
                throw new IllegalArgumentException(ExceptionLocalization.buildMessage("unknown_bean_class", new Object[]{entityClass}));
            }
            if (!descriptor.shouldBeReadOnly() || !descriptor.isSharedIsolation()) {
                session = (AbstractSession) getActiveSession();
            } else {
                session = (AbstractSession) getReadOnlySession();
            }
            // 确保从当前会话中获取实体类的持久化信息描述符
            if (descriptor.hasTablePerMultitenantPolicy()) {
                descriptor = session.getDescriptor(entityClass);
            }
            // 获取分表对应的实体类的持久化信息描述符
            descriptor = ClassDescriptorUtils.getSplitDescriptor(descriptor, session, splitTable);
            // 复用实体管理器实现类的内部方法
            return (T) findInternal(descriptor, session, primaryKey, lockMode, properties);
        } catch (LockTimeoutException e) {
            throw e;
        } catch (RuntimeException e) {
            setRollbackOnly();
            throw e;
        }
    }

    /**
     * 获取指定实体类对应的下一个主键值
     */
    public <T> Number getNextValue(Class<T> entityClass, SplitTable splitTable) {
        // 获取实体类的持久化信息描述符
        AbstractSession session = this.databaseSession;
        ClassDescriptor descriptor = session.getDescriptor(entityClass);
        if (ObjectUtils.isEmpty(descriptor)) {
            throw new IllegalArgumentException(ExceptionLocalization.buildMessage("unknown_bean_class", new Object[]{entityClass}));
        }
        // 获取分表对应的实体类的持久化信息描述符
        descriptor = ClassDescriptorUtils.getSplitDescriptor(descriptor, session, splitTable);

        Number nextValue;
        if (ObjectUtils.isNotEmpty(splitTable) && splitTable.isExistSplitTablePkColumn()) {
            String sequenceName = splitTable.getSplitTablePkColumnValue();
            Sequencing sequencing = session.getSequencing();
            FleaSequencingManager fleaSequencingManager = FleaSequencingManager.getFleaSequencingManager(sequenceName, sequencing, descriptor);
            nextValue = fleaSequencingManager.getNextValue(sequenceName);
        } else {
            nextValue = session.getNextSequenceNumberValue(entityClass);
        }
        return nextValue;
    }

    @Override
    public RepeatableWriteUnitOfWork getActivePersistenceContext(Object txn) {
        // 覆写，详细请查看GitHub
    }
}
```
# 4. 接入讲解
## 4.1 数据库和表
### 4.1.1 模板库
**flea_id_generator** 为主键生成器表，可查看笔者的这篇博文[《flea-db使用之主键生成器表介绍》](/2019/09/04/flea-framework/flea-db/flea-db-jpa-id-generator/)，不再赘述。

![](fleaorder.png)
### 4.1.2 分库1
![](fleaorder1.png)
### 4.1.3 分库2
![](fleaorder2.png)

具体的SQL文件，请参考 [fleaorder.sql](https://github.com/Huazie/flea-db-test/fleaorder.sql)，[fleaorder1.sql](https://github.com/Huazie/flea-db-test/fleaorder1.sql)，[fleaorder2.sql](https://github.com/Huazie/flea-db-test/fleaorder2.sql)

## 4.2 各实体类
|实体表名|    描述       |
|-------------|--------------|
| [OldOrder](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/entity/OldOrder.java) | 旧订单  |
|[OldOrderAttr](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/entity/OldOrderAttr.java) |旧订单属性|
|[Order](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/entity/Order.java)|订单|
|[OrderAttr](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/entity/OrderAttr.java)| 订单属性 |

![](entity.png)

## 4.3 FleaOrder数据源DAO层父类 
[FleaOrderDAOImpl](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/FleaOrderDAOImpl.java) ，该类继承 **AbstractFleaJPADAOImpl**，成员变量 **entityManager** ，由 **PersistenceContext** 注解标记 持久化单元名，这里为模板持久化单元名。

```java
public class FleaOrderDAOImpl<T> extends AbstractFleaJPADAOImpl<T> {

    @PersistenceContext(unitName="fleaorder")
    private EntityManager entityManager;

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public Number getFleaNextValue(T entity) throws CommonException {
        return super.getFleaNextValue(entity);
    }

    @Override
    @FleaTransactional(value = "fleaOrderTransactionManager", unitName = "fleaorder")
    public boolean remove(long entityId) throws CommonException {
        return super.remove(entityId);
    }

    @Override
    @FleaTransactional(value = "fleaOrderTransactionManager", unitName = "fleaorder")
    public boolean remove(String entityId) throws CommonException {
        return super.remove(entityId);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public boolean remove(T entity) throws CommonException {
        return super.remove(entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public boolean remove(long entityId, T entity) throws CommonException {
        return super.remove(entityId, entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public boolean remove(String entityId, T entity) throws CommonException {
        return super.remove(entityId, entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public T update(T entity) throws CommonException {
        return super.update(entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public List<T> batchUpdate(List<T> entities) throws CommonException {
        return super.batchUpdate(entities);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public void save(T entity) throws CommonException {
        super.save(entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public void batchSave(List<T> entities) throws CommonException {
        super.batchSave(entities);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public int insert(String relationId, T entity) throws CommonException {
        return super.insert(relationId, entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public int update(String relationId, T entity) throws CommonException {
        return super.update(relationId, entity);
    }

    @Override
    @FleaTransactional("fleaOrderTransactionManager")
    public int delete(String relationId, T entity) throws CommonException {
        return super.delete(relationId, entity);
    }

    @Override
    protected EntityManager getEntityManager() {
        return entityManager;
    }

}
```

## 4.4 各实体的DAO层接口和实现
可至 GitHub 查看如下 [DAO层代码](https://github.com/Huazie/flea-db-test/tree/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/dao)：
![](DAO.png)
## 4.5 各实体的SV层接口和实现
可至 GitHub 查看如下 [SV层代码](https://github.com/Huazie/flea-db-test/tree/main/flea-jpa-test/src/main/java/com/huazie/fleadbtest/jpa/split/service)：

![](SV.png)

# 5. 单元测试
测试之前，先添加主键生成器表中的数据如下：
|id_generator_key| id_generator_value  |
|-----------------------|----------------------------|
|  pk_old_order   |     999999999          |
|  pk_order           |     999999999          |


## 5.1 分库分表测试
 分库分表测试类 [OrderTest](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/test/java/com/huazie/fleadbtest/jpa/split/OrderTest.java)

下面我们来看下，分库分表的新增、查询、更新 和 查询，代码如下：
```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class OrderTest {

    private static final Logger LOGGER = FleaLoggerProxy.getProxyInstance(OrderTest.class);

    @Resource(name = "orderSV")
    private IOrderSV orderSV;

    @Test
    public void testInsertOrder() throws Exception {

        Order order = new Order();
        order.setOrderName("测试订单");
        order.setOrderPrice(0L);
        order.setOrderState(0);
        order.setOrderDate(DateUtils.getCurrentTime());
        // 获取下一个主键值
        Long orderId = (Long) orderSV.getFleaNextValue(null);
        order.setOrderId(orderId);

        orderSV.save(order);
    }

    @Test
    public void testQueryOrder() throws Exception {

        long orderId = 1000000000L;
        Order order = new Order();
        order.setOrderId(orderId);

        order = orderSV.query(orderId, order);

        LOGGER.debug("Order = {}", order);
    }

    @Test
    public void testUpdateOrder() throws Exception {

        long orderId = 1000000000L;
        Order order = new Order();
        order.setOrderId(orderId);

        Set<String> attrNames = new HashSet<>();
        // orderId 为实体类Order中的变量，实际对应表中 order_id 字段
        attrNames.add("orderId");
        List<Order> orderList = orderSV.query(attrNames, order);

        if (CollectionUtils.isNotEmpty(orderList)) {
            order = orderList.get(0);

            LOGGER.debug("Before : {}", order);

            order.setOrderName("修改订单");
            order.setOrderPrice(100L);
            order.setOrderState(1);
            // 更新数据
            orderSV.update(order);
        }

        Order order1 = new Order();
        order1.setOrderId(orderId);

        order1 = orderSV.query(orderId, order1);

        LOGGER.debug("After : {}", order1);
    }

    @Test
    public void testDeleteOrder() throws Exception {
        long orderId = 1000000000L;
        Order order = new Order();
        order.setOrderId(orderId);

        Set<String> attrNames = new HashSet<>();
        attrNames.add("orderId");
        List<Order> orderList = orderSV.query(attrNames, order);

        if (CollectionUtils.isNotEmpty(orderList)) {
            Order order1 = orderList.get(0);
            LOGGER.error("Before : {}", order1);
            // 删除数据
            orderSV.remove(order1);
        }

        Order order2 = orderSV.query(orderId, order);
        LOGGER.error("After : {}", order2);
    }
}
```
## 5.2 分库测试
如果只分库，不分表，需要再执行具体的增删改查之前，设置分库序列值。

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class OldOrderTest {

    private static final Logger LOGGER = FleaLoggerProxy.getProxyInstance(OldOrderTest.class);

    @Resource(name = "oldOrderSV")
    private IOldOrderSV oldOrderSV;

    @Test
    public void testInsertOldOrder() throws Exception {

        OldOrder oldOrder = new OldOrder();
        oldOrder.setOrderName("测试旧订单1");
        oldOrder.setOrderPrice(200L);
        oldOrder.setOrderState(0);
        oldOrder.setOrderDate(DateUtils.getCurrentTime());
        // 获取下一个主键值
        Long orderId = (Long) oldOrderSV.getFleaNextValue(null);
        oldOrder.setOrderId(orderId);

        // 设置分库序列值
        FleaLibUtil.setSplitLibSeqValue("SEQ", orderId);

        oldOrderSV.save(oldOrder);
    }

    @Test
    public void testQueryOldOrder() throws Exception {

        long orderId = 1000000000L;
        OldOrder oldOrder = new OldOrder();
        //oldOrder.setOrderId(orderId);

        // 设置分库序列值
        FleaLibUtil.setSplitLibSeqValue("SEQ", orderId);

        // 分库场景，需要实体类，为了后续从中获取默认库名
        oldOrder = oldOrderSV.query(orderId, oldOrder);

        LOGGER.debug("Order = {}", oldOrder);
    }

    @Test
    public void testUpdateOldOrder() throws Exception {

        long orderId = 1000000000L;

        // 设置分库序列值
        FleaLibUtil.setSplitLibSeqValue("SEQ", orderId);

        OldOrder oldOrder = new OldOrder();
        oldOrder.setOrderId(orderId);

        Set<String> attrNames = new HashSet<>();
        attrNames.add("orderId");
        List<OldOrder> oldOrderList = oldOrderSV.query(attrNames, oldOrder);

        if (CollectionUtils.isNotEmpty(oldOrderList)) {
            oldOrder = oldOrderList.get(0);

            LOGGER.debug("Before : {}", oldOrder);

            oldOrder.setOrderName("修改旧订单1");
            oldOrder.setOrderPrice(200L);
            oldOrder.setOrderState(2);

            oldOrderSV.update(oldOrder);
        }

        OldOrder oldOrder1 = new OldOrder();
        //oldOrder1.setOrderId(orderId);

        oldOrder1 = oldOrderSV.query(orderId, oldOrder1);

        LOGGER.debug("After : {}", oldOrder1);
    }

    @Test
    public void testDeleteOldOrder() throws Exception {

        long orderId = 1000000000L;

        // 设置分库序列值
        FleaLibUtil.setSplitLibSeqValue("SEQ", orderId);

        OldOrder oldOrder = new OldOrder();
        oldOrder.setOrderId(orderId);

        Set<String> attrNames = new HashSet<>();
        attrNames.add("orderId");
        List<OldOrder> orderList = oldOrderSV.query(attrNames, oldOrder);

        if (CollectionUtils.isNotEmpty(orderList)) {
            OldOrder oldOrder1 = orderList.get(0);
            LOGGER.error("Before : {}", oldOrder1);
            oldOrderSV.remove(oldOrder1);
        }

        OldOrder oldOrder2 = new OldOrder();
        oldOrder2 = oldOrderSV.query(orderId, oldOrder2);
        LOGGER.error("After : {}", oldOrder2);
    }

}
```
## 5.3 JPA事物演示
首先我们先看下如何在 **除了数据源DAO层实现类之外** 的方法上使用自定的事物注解 `@FleaTransactional`，
可至 **GitHub** 查看如下代码 ：
![](show-transaction.png)
这里贴出关键使用代码如下：

其中，**value** 的值为 **模板库事物名**，**unitName** 的值为 **模板库持久化单元名**【也为对应 **模板库名**】
```java
    @Override
    @FleaTransactional(value = "fleaOrderTransactionManager", unitName = "fleaorder")
    public void orderTransaction(Long orderId) throws CommonException {
        // 省略。。。
    }
```

**那上面该如何调用呢？**

```java
    @Test
    public void testTransaction() throws Exception {

        long orderId = 1000000000L;

        // 设置分库序列值
        FleaLibUtil.setSplitLibSeqValue("SEQ", orderId);

        fleaOrderModuleSV.orderTransaction(orderId);
    }
```

# 总结
经历了几版的重构，分库分表的实现终于可以和大家见面了，后续我将继续优化和完善 [Flea Framework](https://github.com/Huazie/flea-framework) 这个框架，希望看到的朋友，可以多多支持！！！这篇博文内容有点多，能看到最后属实不容易，再次感谢大家的支持！！！

