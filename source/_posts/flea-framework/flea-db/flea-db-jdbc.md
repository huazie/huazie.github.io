---
title: flea-db使用之封装JDBC接入
date: 2019-10-16 10:12:20
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - FleaJDBCHelper
  - FleaJDBCConfig
---

![](/images/flea-logo.png)

# 引言
本篇将要演示 FleaJDBCHelper 的使用，该工具类封装了基本的JDBC增删改查的操作，只需简单几步即可实现数据库操作。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 准备工作
为了演示JDBC接入（参考 [JPA接入](../../../../../../2019/09/12/flea-framework/flea-db/flea-db-jpa-integration/) 中的准备工作），需要如下准备：

 1. MySQL数据库 (客户端可以使用 navicat for mysql)
 2. 新建测试数据库 fleajpatest
 3. 新建测试表 student

# 2. 接入讲解
## 2.1 JDBC数据源配置
数据源配置独立出来，定义在 [flea-config.xml](https://github.com/Huazie/flea-jpa-test/blob/main/src/main/resources/flea/flea-config.xml) 中，可添加多个JDBC数据库配置【即 **config-items** 节点】。

```xml
<?xml version="1.0" encoding="UTF-8"?>

<flea-config>

    <!-- 其他配置省略 -->
    <config-items key="mysql-fleajpatest" desc="JDBC数据库配置【key=数据库系统-数据库或数据库用户】">
        <config-item key="driver" desc="mysql数据库驱动名">com.mysql.jdbc.Driver</config-item>
        <config-item key="url" desc="mysql数据库连接地址">jdbc:mysql://localhost:3306/fleajpatest?useUnicode=true&amp;characterEncoding=UTF-8</config-item>
        <config-item key="user" desc="mysql数据库登录用户名">root</config-item>
        <config-item key="password" desc="mysql数据库登录密码">root</config-item>
    </config-items>

</flea-config>
```
## 2.2 定义Flea数据库单元
每个 [FleaDBUnit](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jdbc/src/main/java/com/huazie/fleaframework/db/jdbc/pojo/FleaDBUnit.java) 中， **database** 和 **name** 对应上述 **config-items** 中的 **key**；**driver**、**url**、**user**、**password** 分别对应上述 **config-item** 中的配置。
```java
/**
 * <p> Flea 数据库单元 </p>
 *
 * @author huazie
 */
public class FleaDBUnit {

    private String database;
    private String name;
    private String driver;
    private String url;
    private String user;
    private String password;

    // 省略 set 和 get方法
}
```

- `database` : 数据库管理系统名
- `name` : 数据库名 或 数据库用户名
- `driver` : 数据库驱动名
- `url` : 数据库连接地址
- `user` : 数据库登录用户名
- `password` : 数据库登录密码

## 2.3 定义Flea数据库操作类
[FleaDBOperation](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jdbc/src/main/java/com/huazie/fleaframework/db/jdbc/pojo/FleaDBOperation.java) 封装了JDBC的数据库操作对象，包括数据库连接对象 **Connection**、数据库预编译状态对象 **PreparedStatement** 和 数据库结果集对象 **ResultSet**。该类继承 **Closeable**，实现 **close** 方法，用于每次 **JDBC** 数据库操作后释放资源【这里用到了 [try-with-resource 语法糖](https://docs.oracle.com/javase/tutorial/essential/exceptions/tryResourceClose.html) 】。
```java
/**
 * <p> Flea数据库操作 </p>
 *
 * @author huazie
 */
public class FleaDBOperation implements Closeable {

    private Connection connection;

    private PreparedStatement preparedStatement;

    private ResultSet resultSet;

    @Override
    public void close() {
        FleaJDBCConfig.close(connection, preparedStatement, resultSet);
    }
    
    // 省略 set 和 get方法
}
```
## 2.4 定义Flea JDBC配置类
[FleaJDBCConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jdbc/src/main/java/com/huazie/fleaframework/db/jdbc/config/FleaJDBCConfig.java) 读取数据库的配置信息，该信息存在于 `flea-config.xml` 中。 其中的 **init(String mDatabase, String mName)**  方法用于初始化本次操作的数据库管理系统【mDatabase】、数据库名或数据库用户 【mName】，两者对应数据库配置 **config-items** 中的 **key** 。

```java
public class FleaJDBCConfig {

    private static volatile FleaJDBCConfig config;

    private static final ConcurrentMap<String, FleaDBUnit> fleaDBUnits = new ConcurrentHashMap<>();

    private static final Object fleaDBUnitsLock = new Object();

    private FleaJDBCConfig() {
    }

    /**
     * <p> 读取数据库相关配置信息 </p>
     *
     * @return JDBC配置对象
     */
    public static FleaJDBCConfig getConfig() {

        if (ObjectUtils.isEmpty(config)) {
            synchronized (FleaJDBCConfig.class) {
                if (ObjectUtils.isEmpty(config)) {
                    config = new FleaJDBCConfig();
                }
            }
        }
        return config;
    }

    /**
     * <p> 使用之前先初始化 </p>
     *
     * @param mDatabase 数据库管理系统名称
     * @param mName     数据库名  或  数据库用户
     */
    public static void init(String mDatabase, String mName) {
        FleaFrameManager.getManager().setDBConfigKey(mDatabase, mName);
    }

    /**
     * <p> 建立数据库连接 </p>
     *
     * @return 数据库连接对象
     * @since 1.0.0
     */
    public Connection getConnection() {
        Connection conn = null;
        FleaDBUnit fleaDBUnit;

        String dbConfigKey = FleaFrameManager.getManager().getDBConfigKey();

        if (!fleaDBUnits.containsKey(dbConfigKey)) {
            synchronized (fleaDBUnitsLock) {
                if (!fleaDBUnits.containsKey(dbConfigKey)) {
                    fleaDBUnits.put(dbConfigKey, getFleaDBUnit(dbConfigKey));
                }
            }
        }

        fleaDBUnit = fleaDBUnits.get(dbConfigKey);

        try {
            // 请正确初始化数据库管理系统和数据库（或数据库用户）
            ObjectUtils.checkEmpty(fleaDBUnit, DaoException.class, "ERROR-DB-DAO0000000013");
            Class.forName(fleaDBUnit.getDriver());
            conn = DriverManager.getConnection(fleaDBUnit.getUrl(), fleaDBUnit.getUser(), fleaDBUnit.getPassword());
        } catch (Exception e) {
            LOGGER.error("获取数据库连接异常 ：{}", e.getMessage());
        }
        return conn;
    }

    /**
     * <p> 读取指定配置键的数据库相关配置信息 </p>
     *
     * @param dbConfigKey 数据库配置键
     * @return 数据库配置信息类对象
     */
    private FleaDBUnit getFleaDBUnit(String dbConfigKey) {
        FleaDBUnit fDBUnit = null;
        if (StringUtils.isNotBlank(dbConfigKey)) {
            fDBUnit = new FleaDBUnit();
            String[] dbConfigKeyArr = StringUtils.split(dbConfigKey, CommonConstants.SymbolConstants.HYPHEN);
            if (ArrayUtils.isNotEmpty(dbConfigKeyArr) && CommonConstants.NumeralConstants.INT_TWO == dbConfigKeyArr.length) {
                fDBUnit.setDatabase(dbConfigKeyArr[0]);
                fDBUnit.setName(dbConfigKeyArr[1]);
            }
            fDBUnit.setDriver(FleaConfigManager.getConfigItemValue(dbConfigKey, DBConfigConstants.DB_CONFIG_DRIVER));
            fDBUnit.setUrl(FleaConfigManager.getConfigItemValue(dbConfigKey, DBConfigConstants.DB_CONFIG_URL));
            fDBUnit.setUser(FleaConfigManager.getConfigItemValue(dbConfigKey, DBConfigConstants.DB_CONFIG_USER));
            fDBUnit.setPassword(FleaConfigManager.getConfigItemValue(dbConfigKey, DBConfigConstants.DB_CONFIG_PASSWORD));
        }

        return fDBUnit;
    }

    /**
     * <p> 释放连接Connection </p>
     *
     * @param conn 数据库连接对象
     */
    private static void closeConnection(Connection conn) {
        try {
            if (ObjectUtils.isNotEmpty(conn)) {
                conn.close();
            }
        } catch (SQLException e) {
        }
    }

    /**
     * <p> 释放statement </p>
     *
     * @param statement Statement对象
     */
    private static void closeStatement(Statement statement) {
        try {
            if (ObjectUtils.isNotEmpty(statement)) {
                statement.close();
            }
        } catch (SQLException e) {
        }
    }

    /**
     * <p> 释放ResultSet结果集 </p>
     *
     * @param rs 结果集对象
     */
    private static void closeResultSet(ResultSet rs) {
        try {
            if (ObjectUtils.isNotEmpty(rs)) {
                rs.close();
            }
        } catch (SQLException e) {
        }
    }

    /**
     * <p> 释放资源 </p>
     *
     * @param conn      数据库连接对象
     * @param statement 数据库状态对象
     * @param rs        数据库结果集对象
     */
    public static void close(Connection conn, Statement statement, ResultSet rs) {
        closeResultSet(rs);
        closeStatement(statement);
        closeConnection(conn);
    }
    
    // 省略一些close方法
}
```
上述 **init** 方法中，使用了 [FleaFrameManager.getManager().setDBConfigKey](https://github.com/Huazie/flea-framework/blob/dev/flea-common/src/main/java/com/huazie/fleaframework/common/FleaFrameManager.java)，这里是将 JDBC连接的数据库配置添加到了线程对象中 [ThreadLocal](../../../../../../2021/04/12/java/java-concurrency-learning/java-concurrency-learning4/) 中，代码如下：

```java
    // 当前线程采用JDBC连的数据库前缀配置
    private static ThreadLocal<String> sDBLocal = new ThreadLocal<>(); 
    
    /**
     * <p> 获取当前线程中使用JDBC连接的数据库配置键 </p>
     *
     * @return 当前线程中使用JDBC连接的数据库配置键
     */
    public String getDBConfigKey() {
        return sDBLocal.get();
    }

    /**
     * <p> 设置当前线程中使用JDBC连接的数据库配置键 </p>
     *
     * @param dbSysName 数据库系统名
     * @param dbName    数据库名
     */
    public void setDBConfigKey(String dbSysName, String dbName) {
        if (StringUtils.isNotBlank(dbSysName) && StringUtils.isNotBlank(dbName)) {
            sDBLocal.set(dbSysName.toLowerCase() + CommonConstants.SymbolConstants.HYPHEN + dbName.toLowerCase());
        }
    }
```
## 2.4 定义 Flea JDBC工具类 

在使用 [FleaJDBCHelper](https://github.com/Huazie/flea-framework/blob/dev/flea-db/flea-db-jdbc/src/main/java/com/huazie/fleaframework/db/jdbc/FleaJDBCHelper.java) 之前，一定要先调用一下 **init**。

![](FleaJDBCHelper.png)

# 3. 接入自测 
请查看单元测试类 [StudentJDBCTest](https://github.com/Huazie/flea-jpa-test/blob/master/src/test/java/com/huazie/jpa/StudentJDBCTest.java)
## 3.1 JDBC新增数据
```java
    @Test
    public void testStudentInsert() throws Exception {
        // 初始化数据库配置，用于获取具体操作数据源
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        String sql = "insert into student(stu_name, stu_age, stu_sex, stu_state) values(?, ?, ?, ?)";

        List<Object> paramList = new ArrayList<Object>();
        paramList.add("huazie");
        paramList.add(25);
        paramList.add(1);
        paramList.add(1);

        int ret = FleaJDBCHelper.insert(sql, paramList);

        LOGGER.debug("RESULT = {}", ret);
    }
```
**执行结果：**
![](result_add.png)

## 3.2 JDBC查询数据
这里的查询语句，可以是复杂SQL，返回结果 **List<Map<String, Object>>**
```java
    @Test
    public void testStudentQuery() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        String sql = "select * from student where stu_state = ?";

        List<Object> paramList = new ArrayList<Object>();
        paramList.add(1);

        LOGGER.debug("RESULT LIST = {}", FleaJDBCHelper.query(sql, paramList));
    }
```
**执行结果：**
![](result_query.png)

```java
    @Test
    public void testStudentSingleQuery() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        String sql = "select count(*) from student where stu_state = ?";

        List<Object> paramList = new ArrayList<Object>();
        paramList.add(1);

        LOGGER.debug("COUNT = {}", FleaJDBCHelper.querySingle(sql, paramList));
    }
```
**执行结果：**
![](result_query-1.png)

## 3.3 JDBC更新数据
```java
    @Test
    public void testStudentUpdate() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        String sql = "update student set stu_state = ? where stu_name = ?";

        List<Object> paramList = new ArrayList<Object>();
        paramList.add(2);
        paramList.add("huazie");

        int ret = FleaJDBCHelper.update(sql, paramList);

        LOGGER.debug("RESULT = {}", ret);
    }
```

**执行结果：**
![](result_update.png)

## 3.4 JDBC删除数据
```java
    @Test
    public void testStudentDelete() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        String sql = "delete from student where stu_name = ? and stu_state = ? ";

        List<Object> paramList = new ArrayList<Object>();
        paramList.add("huazie");
        paramList.add(2);

        int ret = FleaJDBCHelper.delete(sql, paramList);

        LOGGER.debug("RESULT = {}", ret);
    }
```

**执行结果：**
![](result_delete.png)
