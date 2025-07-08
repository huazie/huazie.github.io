---
title: flea-db使用之SQL模板接入
date: 2019-11-04 10:18:11
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - SQL模板
  - INSERT模板
  - SELECT模板
  - UPDATE模板
  - DELETE模板
---

![](/images/flea-logo.png)

# 引言
本篇将要演示 `SQL` 模板的使用，目前包含 `INSERT` 模板、`SELECT` 模板、`UPDATE` 模板、`DELETE`模板。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 准备工作
为了演示SQL模板接入（参考 [JPA接入](../../../../../../2019/09/12/flea-framework/flea-db/flea-db-jpa-integration/) 中的准备工作），需要如下准备：

 1. MySQL数据库 (客户端可以使用 navicat for mysql)
 2. 新建测试数据库 **fleajpatest**
 3. 新建测试表 **student**

# 2. 使用讲解
## 2.1 SQL模板配置
SQL模板配置包含了SQL模板规则，SQL模板定义，SQL模板参数，SQL关系配置。具体配置可至GitHub，查看 [flea-sql-template.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/flea/db/flea-sql-template.xml)

**SQL** 模板规则，即定义 **SQL** 模板的校验规则，主要包含增删改查的 **4** 类模板。
以 **INSERT SQL** 模板的校验规则配置举例，如下所示【这里属性 `value` 值其实就是 **INSERT SQL** 模板的正则表达式】：

```xml
    <!-- 这边的规则为固定写法，一般不做修改 -->
    <rules>
        <rule id="insert" name="INSERT SQL模板的校验规则配置">
            <property key="sql" value="[ ]*(INSERT)[ ]+(INTO)[ ]+##table##[ ]+\([ ]*##columns##[ ]+\)[ ]+(VALUES)[ ]+\([ ]*##values##[ ]+\)[ ]*" />
        </rule>
    </rules>
```

**SQL** 模板定义，即定义通用的增删改查 **SQL** 模板。
以 **INSERT SQL** 模板定义举例，如下所示：

```xml
    <!-- SQL模板定义配置 -->
    <templates>
        <template id="insert" ruleId="insert" name="INSERT SQL模板" desc="用于原生SQL中INSERT语句的使用">
            <!-- SQL模板数据 -->
            <property key="template" value="INSERT INTO ##table## (##columns## ) VALUES (##values## )" />
            <!-- SQL模板类型 -->
            <property key="type" value="insert"/>
        </template>
    </templates>
```

**SQL** 模板参数，即定义 **SQL** 模板中的参数取值。 
以 **INSERT SQL** 模板参数举例，如下所示：

```xml
    <!-- SQL模板参数配置 -->
    <params>
        <param id="insert" name="SQL模板參數" desc="用于定义SQL模板中的替换参数">
            <!-- 表名 -->
            <property key="table" value="flea_config_data" />
            <!-- 这两个不填，表示表的字段全部使用
            <property key="columns" value="config_id, config_type, config_code, config_name, data1, config_state" />
            <property key="values" value=":configId:, :configType:, :configCode:, :configName:, :data1:, :configState:" />-->
        </param>
    </params>
```

**SQL** 关系配置，用于关联 **SQL** 模板和 **SQL** 模板参数。
以 **INSERT SQL** 关系配置举例，如下所示：

```xml
    <!-- SQL模板和模板参数关联关系配置（简称 SQL关系配置）-->
    <relations>
        <relation id="insert" templateId="insert" paramId="insert" name="SQL关系"/>
    </relations>
```

**relation** 用于定义一条 **SQL** 关系配置：

- `id` : SQL关系编号
- `templateId` : SQL模板编号
- `paramId` : SQL模板参数编号

## 2.2 新增数据
**相关配置可查看 ：** 

```xml
    <param id="insert" name="SQL模板參數" desc="用于定义SQL模板中的替换参数">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- 这两个不填，表示表的字段全部使用-->
        <property key="columns" value="stu_name, stu_age, stu_sex, stu_state" />
        <property key="values" value=":stuName:, :stuAge:, :stuSex:, :stuState:" />
    </param>

    <relation id="insert" templateId="insert" paramId="insert" name="SQL关系"/>
```

**JPA方式接入SQL模板：**
```java
    @Test
    public void testInsertSqlTemplateFromJPA() throws Exception{
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");

        Student student = new Student();
        student.setStuName("王老五");
        student.setStuAge(35);
        student.setStuSex(1);
        student.setStuState(1);

        int ret = studentSV.insert("insert", student);
        LOGGER.debug("result = {}", ret);
    }
```

**运行结果：**
![](jpa-result-add.png)

**新增数据：**
![](jpa-add-data.png)

**JDBC方式接入SQL模板：**
```java
    @Test
    public void testInsertSqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Student student = new Student();
        student.setStuName("钱老六");
        student.setStuAge(30);
        student.setStuSex(1);
        student.setStuState(1);

        int ret = FleaJDBCHelper.insert("insert", student);
        LOGGER.debug("result = {}", ret);
    }
```

**运行结果：**
![](jdbc-result-add.png)

**新增数据：**
![](jdbc-add-data.png)

## 2.3 查询数据
**相关配置可查看 ：** 

```xml
    <param id="select" name="SQL模板參數" desc="用于定义SQL模板中的替换参数; 如需查询全部，则设置key=columns的属性值为 *，即可">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- SELECT 显示列 -->
        <property key="columns" value="*" />
        <!-- WHERE 子句 , 出现 xml不能直接识别的需要转义，如 >, < 等-->
        <property key="conditions" value="stu_name LIKE :stuName: AND stu_sex = :stuSex: AND stu_age &gt;= :minAge: AND stu_age &lt;= :maxAge:" />
    </param>

    <relation id="select" templateId="select" paramId="select" name="SQL关系"/>
```
**JPA方式接入SQL模板：**
```java
    @Test
    public void testQuerySqlTemplateFromJPA() throws Exception {
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");

        Student student = new Student();
        student.setStuName("%老%");
        student.setStuSex(1);
        student.put("minAge", 20);
        student.put("maxAge", 40);

        LOGGER.debug("Student List = {}", studentSV.query("select", student));
    }
```

**运行结果：**
![](jpa-result-query.png)

**JDBC方式接入SQL模板：**
```java
    @Test
    public void testQuerySqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Student student = new Student();
        student.setStuName("%老%");
        student.setStuSex(1);
        student.put("minAge", 20);
        student.put("maxAge", 40);

        LOGGER.debug("Student List = {}", FleaJDBCHelper.query("select", student));
    }
```

**运行结果：**
![](jdbc-result-query.png)

## 2.4 更新数据
**相关配置可查看 ：** 

```xml
    <param id="update" name="SQL模板參數" desc="用于定义SQL模板中的替换参数">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- SET 子句 -->
        <property key="sets" value="stu_name = :stuName, stu_age = :stuAge" />
        <!-- WHERE 子句 , 出现 xml不能直接识别的需要转义，如 >, < 等-->
        <property key="conditions" value="stu_name LIKE :sName: AND stu_state = :stuState: AND stu_sex = :stuSex: AND stu_age &gt;= :minAge: AND stu_age &lt;= :maxAge:" />
    </param>
        
    <relation id="update" templateId="update" paramId="update" name="SQL关系"/>
```

**JPA方式接入SQL模板：**

```java
    @Test
    public void testUpdateSqlTemplateFromJPA() throws Exception {
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");

        Student student = new Student();
        student.setStuName("王老五1");
        student.setStuAge(40);
        student.setStuState(1);
        student.setStuSex(1);
        student.put("sName", "%王老五%");
        student.put("minAge", 20);
        student.put("maxAge", 40);

        LOGGER.debug("Result = {}", studentSV.update("update", student));
    }
```

**运行结果：**
![](jpa-result-update.png)
![](jpa-update-data.png)

**JDBC方式接入SQL模板：**
```java
    @Test
    public void testUpdateSqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Student student = new Student();
        student.setStuName("钱老六1");
        student.setStuAge(35);
        student.setStuState(1);
        student.setStuSex(1);
        student.put("sName", "%钱老六%");
        student.put("minAge", 20);
        student.put("maxAge", 40);

        LOGGER.debug("Result = {}", FleaJDBCHelper.update("update", student));
    }
```

**运行结果：**
![](jdbc-result-update.png)
![](jdbc-update-data.png)

## 2.5 删除数据
**相关配置可查看 ：** 

```xml
    <param id="delete" name="SQL模板參數" desc="用于定义SQL模板中的替换参数">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- WHERE 子句 -->
        <property key="conditions" value="stu_name LIKE :stuName: AND stu_state = :stuState: AND stu_sex = :stuSex: AND stu_age &gt;= :minAge: AND stu_age &lt;= :maxAge:" />
    </param>
        
    <relation id="delete" templateId="delete" paramId="delete" name="SQL关系"/>
```

**JPA方式接入SQL模板：**

```java
    @Test
    public void testDeleteSqlTemplateFromJPA() throws Exception {
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");

        Student student = new Student();
        student.setStuName("%王老五%");
        student.setStuState(1);
        student.setStuSex(1);
        student.put("minAge", 20);
        student.put("maxAge", 40);

        LOGGER.debug("Result = {}", studentSV.delete("delete", student));
    }
```

**运行结果：**
![](jpa-result-delete.png)

**JDBC方式接入SQL模板：**
```java
    @Test
    public void testDeleteSqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Student student = new Student();
        student.setStuName("%钱老六%");
        student.setStuState(1);
        student.setStuSex(1);
        student.put("minAge", 20);
        student.put("maxAge", 40);

        LOGGER.debug("Result = {}", FleaJDBCHelper.delete("delete", student));
    }
```

**运行结果：**
![](jdbc-result-delete.png)

## 2.6 分页查询
当前数据库数据如下：
![](student.png)

**相关配置可查看 ：** 

```xml
    <param id="select_1" name="SQL模板參數" desc="用于定义SQL模板中的替换参数; 如需查询全部，则设置key=columns的属性值为 *，即可">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- SELECT 显示列 -->
        <property key="columns" value="*" />
        <!-- WHERE 子句 , 出现 xml不能直接识别的需要转义，如 >, < 等-->
        <property key="conditions" value="stu_name LIKE :stuName: AND stu_sex = :stuSex: AND stu_age &gt;= :minAge: AND stu_age &lt;= :maxAge: ORDER BY stu_id DESC LIMIT :pageStart:, :pageCount:" />
    </param>
    
    <relation id="select_1" templateId="select" paramId="select_1" name="SQL关系"/>
```
**JPA方式接入SQL模板：**
```java
    @Test
    public void testQueryPageSqlTemplateFromJPA() throws Exception {
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");

        Student student = new Student();
        student.setStuName("%张三%");
        student.setStuSex(1);
        student.put("minAge", 18);
        student.put("maxAge", 20);
        int pageNum = 1;   // 第一页
        int pageCount = 5; // 每页5条记录
        student.put("pageStart", (pageNum - 1) * pageCount);
        student.put("pageCount", pageCount);

        List<Student> studentList = studentSV.query("select_1", student);
        LOGGER.debug("Student List = {}", studentList);
        LOGGER.debug("Student Count = {}", studentList.size());
    }
```

运行结果：
![](jpa-result-page.png)

**JDBC方式接入SQL模板：**
```java
     @Test
    public void testQueryPageSqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Student student = new Student();
        student.setStuName("%李四%");
        student.setStuSex(1);
        student.put("minAge", 18);
        student.put("maxAge", 20);
        int pageNum = 1;   // 第一页
        int pageCount = 5; // 每页5条记录
        student.put("pageStart", (pageNum - 1) * pageCount);
        student.put("pageCount", pageCount);

        List<Map<String, Object>> studentList = FleaJDBCHelper.query("select_1", student);
        LOGGER.debug("Student List = {}", studentList);
        LOGGER.debug("Student Count = {}", studentList.size());
    }
```

**运行结果：**
![](jdbc-result-page.png)

## 2.7 单个结果查询--计数
**相关配置可查看 ：** 

```xml
    <param id="select_2" name="SQL模板參數" desc="用于定义SQL模板中的替换参数; 如需查询全部，则设置key=columns的属性值为 *，即可">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- SELECT 显示列 -->
        <property key="columns" value="count(*)" />
        <!-- WHERE 子句 , 出现 xml不能直接识别的需要转义，如 >, < 等-->
        <property key="conditions" value="stu_name LIKE :stuName: AND stu_sex = :stuSex: AND stu_age &gt;= :minAge: AND stu_age &lt;= :maxAge:" />
    </param>
        
    <relation id="select_2" templateId="select" paramId="select_2" name="SQL关系"/>
```
**JPA方式接入SQL模板：**
```java
    @Test
    public void testQueryCountSqlTemplateFromJPA() throws Exception {
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");

        Student student = new Student();
        student.setStuName("%张三%");
        student.setStuSex(1);
        student.put("minAge", 18);
        student.put("maxAge", 20);

        LOGGER.debug("Student Count = {}", studentSV.querySingle("select_2", student));
    }
```

**运行结果：**
![](jpa-result-count.png)

**JDBC方式接入SQL模板：**
```java
    @Test
    public void testQueryCountSqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Student student = new Student();
        student.setStuName("%李四%");
        student.setStuSex(1);
        student.put("minAge", 18);
        student.put("maxAge", 20);

        LOGGER.debug("Student Count = {}", FleaJDBCHelper.querySingle("select_2", student));
    }
```

**运行结果：**
![](jdbc-result-count.png)

## 2.8 单个结果查询--总和
**相关配置可查看 ：** 

```xml
    <param id="select_3" name="SQL模板參數" desc="用于定义SQL模板中的替换参数; 如需查询全部，则设置key=columns的属性值为 *，即可">
        <!-- 表名 -->
        <property key="table" value="student" />
        <!-- SELECT 显示列 -->
        <property key="columns" value="sum(stu_age)" />
        <!-- WHERE 子句 , 出现 xml不能直接识别的需要转义，如 >, < 等-->
        <property key="conditions" value="1=1" />
    </param>
        
    <relation id="select_3" templateId="select" paramId="select_3" name="SQL关系"/>
```
**JPA方式接入SQL模板：**
```java
    @Test
    public void testQuerySumSqlTemplateFromJPA() throws Exception {
        IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");
        LOGGER.debug("Student Age = {}", studentSV.querySingle("select_3", new Student()));
    }
```

**运行结果：**
![](jpa-result-sum.png)

**JDBC方式接入SQL模板：**
```java
    @Test
    public void testQuerySumSqlTemplateFromJDBC() throws Exception {
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");
        LOGGER.debug("Student Age = {}", FleaJDBCHelper.querySingle("select_3", new Student()));
    }
```

**运行结果：**
![](jdbc-result-sum.png)

上述单个结果查询，展示了count和sum，其他avg，max，min等相关内容可以移步 GitHub， 查看 [StudentSqlTemplateTest](https://github.com/Huazie/flea-db-test/blob/main/flea-jpa-test/src/test/java/com/huazie/fleadbtest/StudentSqlTemplateTest.java)