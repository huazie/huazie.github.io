---
title: flea-frame-db使用之基于FleaJPAQuery实现JPA分表查询【旧】
date: 2019-10-02 09:26:05
updated: 2023-06-28 20:59:54
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - FleaJPAQuery
---

[《开发框架-Flea》](/categories/开发框架-Flea/) [《flea-db》](/categories/开发框架-Flea/flea-db/)

![](/images/jpa-logo.png)

# 引言
本文采用 **EclipseLink** 的 **JPA** 实现，相关 **FleaJPAQuery** 的接入使用请移步我的 [另外几篇博文](/categories/开发框架-Flea/flea-db/)。

首先讨论一下，为了实现 **JPA** 分表查询，我们需要做哪些事情：

 - 分表规则定义（即从主表到分表的转换实现）
 - 分表查询实现（即JPA标准化查询组件根据分表规则查询具体分表）

# 1. JPA标准化查询
在JPA中，实体对应的表由如下注解定义：

```java
@Entity
@Table(name = "flea_login_log")
```
如上可见，实体类实际上只会对应一个表名，单纯从这里是无法实现分表查询。
那么既然这样无法分表，我们选择退而求其次，看看表名是什么时候，被那个对象使用，因为我们可以确认查询最后的表名，一定是使用的注解定义的表名。
下面是调试过程的发现：
**com.huazie.frame.db.jpa.common.FleaJPAQuery**
```java
    /**
     * <p> Flea JPA查询对象池获取之后，一定要调用该方法进行初始化 </p>
     *
     * @param entityManager JPA中用于增删改查的接口
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
    }
```
如下两张图是根SQL表达式对象 Root 的Debug视图，发现存储实际表名的是 **DatabaseTable** 对象。
![](root.png)

![](descriptor.png)

那么既然找到了表名实际相关的地方，下面的重点就是如何在使用的JPA标准化查询的过程中，动态改变查询的表名。

下面给出上述我们需要做的事情的解决方案：

# 2. 分表规则定义
实体类中定义的表名，我们可以理解为主表名；分表名的命名规则首先需要确定一下，定义如下配置：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<tables>

    <!-- 定义分表配置
        name : 分表对应的主表名
        exp  : 分表名表达式 (FLEA_TABLE_NAME)_(列名大写)_(列名大写)
    -->
    <table name="flea_login_log" exp="(FLEA_TABLE_NAME)_(CREATE_DATE)" desc="Flea登录日志表分表规则">
        <splits>
            <!-- 定义分表后缀
                key : 分表类型关键字 (可查看 com.huazie.frame.db.common.table.split.TableSplitEnum )
                column : 分表属性列字段名
                implClass : 分表后缀转换实现类
            -->
            <split key="yyyymm" column="create_date" implClass="com.huazie.frame.db.common.table.split.impl.YYYYMMTableSplitImpl"/>
        </splits>
    </table>

</tables>
```
分表规则相关实现代码，可以移步 GitHub 查看 [TableSplitHelper](https://github.com/Huazie/flea-frame/blob/dev/flea-frame-db/src/main/java/com/huazie/frame/db/common/table/split/TableSplitHelper.java)

# 3. 分表查询实现
在上述分表规则定义中， 我们可以看到分表名表达式**exp**是由 主表名 和 分表字段 组成，分表字段的转换实现规则由**split**定义。
```java
    @Override
    public void handle(CriteriaQuery criteriaQuery, Object entity) throws Exception {

        if (ObjectUtils.isEmpty(criteriaQuery) || ObjectUtils.isEmpty(entity)) {
            return;
        }

        // 获取分表信息（包括主表名 和 分表名 【如果存在分表返回】）
        SplitTable splitTable = EntityUtils.getSplitTable(entity);

        // 存在分表，需要查询指定分表
        if (StringUtils.isNotBlank(splitTable.getSplitTableName())) {
            Set<Root<?>> roots = criteriaQuery.getRoots();
            if (CollectionUtils.isNotEmpty(roots)) {
                // 重新设置 查询的分表表名
                ((EntityTypeImpl<?>) roots.toArray(new Root<?>[0])[0].getModel()).getDescriptor().setTableName(splitTable.getSplitTableName());
            }
        }
    }
```
JPA分表查询相关代码可以 移步 GitHub 查看 [FleaJPAQuery](https://github.com/Huazie/flea-frame/blob/dev/flea-frame-db/src/main/java/com/huazie/frame/db/jpa/common/FleaJPAQuery.java) 和 [EclipseLinkTableSplitHandler](https://github.com/Huazie/flea-frame/blob/dev/flea-frame-db/src/main/java/com/huazie/frame/db/jpa/persistence/impl/EclipseLinkTableSplitHandler.java)； 

# 4. 自测
自测类可以查看 [AuthTest](https://github.com/Huazie/flea-frame/blob/dev/flea-frame-auth/src/test/java/com/huazie/frame/auth/AuthTest.java)。

```java
    @Test
    public void testFleaLoginLog() {
        try {
            FleaLoginLog fleaLoginLog = new FleaLoginLog();
            fleaLoginLog.setLoginIp4("127.0.0");
            fleaLoginLog.setCreateDate(DateUtils.getCurrentTime());

            FleaJPAQueryPool fleaJPAQueryPool = FleaObjectPoolFactory.getFleaObjectPool(FleaJPAQuery.class, FleaJPAQueryPool.class);
            FleaJPAQuery query = fleaJPAQueryPool.getFleaObject();
            LOGGER.debug("FleaJPAQuery: {}", query);
            query.init(em, FleaLoginLog.class, null);
            // 去重查询某一列数据, 模糊查询 para_code
            query.initQueryEntity(fleaLoginLog).distinct("accountId").like("loginIp4");
            List<String> list = query.getSingleResultList();
            LOGGER.debug("List : {}", list);

            FleaJPAQuery query1 = fleaJPAQueryPool.getFleaObject();
            LOGGER.debug("FleaJPAQuery: {}", query1);
            query1.init(em, FleaLoginLog.class, null);
            List<FleaLoginLog> fleaLoginLogList = query1.initQueryEntity(fleaLoginLog).getResultList();
            LOGGER.debug("Resource List : {}", fleaLoginLogList);

        } catch (Exception e) {
            LOGGER.error("Exception:", e);
        }
    }
```


# 更新
这一版本存在并发的问题，目前已经重构，详见笔者后续的 **flea-db使用之JPA分库分表实现**，也可至GitHub查看笔者的 flea-framework中的 [flea-db](https://github.com/Huazie/flea-framework/tree/main/flea-db) 模块。
