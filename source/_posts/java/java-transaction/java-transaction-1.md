---
title: Java事务入门：从基础概念到初步实践
date: 2024-05-29 23:49:30
updated: 2024-05-29 23:49:30
categories:
  - 开发语言-Java
tags:
  - Java
  - JDBC
  - JPA
  - Mybatis
  - Spring Data JPA
  - Spring JDBC
  - 事务管理
---

![](/images/java-logo.png)

# 引言
在 **Java** 语言相关的应用开发中，事务（**Transaction**）是其中一个核心概念，尤其是在涉及数据库操作时。理解并正确使用事务，可以确保应用系统数据的完整性和一致性。本文 **Huazie** 将带您从 **Java** 事务的基础概念出发，通过不同场景的事务管理实操，帮助您快速入门 **Java** 事务。

<!-- more -->

# 1. Java事务基础概念
## 1.1 什么是事务？
事务是数据库操作的基本执行单元，它要么完全地执行，要么完全地不执行。

事务作为一个逻辑工作单位，其关键特性【数据的完整性和一致性】就体现在其 **ACID** 属性上，如下：

- **原子性（Atomicity）**：事务是一个原子操作单元，其对数据的修改要么全都执行，要么全都不执行。
- **一致性（Consistency）**：事务必须使数据库从一个一致性状态变换到另一个一致性状态。
- **隔离性（Isolation）**：事务的隔离性是指一个事务的执行不能被其他事务干扰，即一个事务内部的操作及使用的数据对并发的其他事务是隔离的，并发执行的各个事务之间不能互相干扰。
- **持久性（Durability）**：持久性是指一个事务一旦提交，它对数据库中数据的改变就应该是永久性的。接下来的其他操作或故障不应该对其执行结果有任何影响。
## 1.2 为什么需要事务？
我们先来想一下如下的场景：

**多个数据库操作之间，有部分操作成功，部分操作失败，会咋样？**

显然这个时候，应用系统的数据就会处于不一致的状态。这种不一致可能会带来非常严重的后果，包括但不限于如下：

- **数据错误**：由于某些操作成功而某些操作失败，数据可能处于一个中间状态，既不完全符合操作前的状态，也不符合操作后预期的状态。
- **业务逻辑错误**：当数据出现不一致时，基于这些数据执行的业务逻辑可能会出现错误。例如，在一个订单系统中，如果订单的创建和库存的减少操作没有同时成功，那么可能会导致订单状态与库存状态不匹配，从而影响整个业务流程。
- **系统可靠性降低**：一个不可靠的数据系统可能会导致更多的错误和故障，需要更多的维护和支持，从而增加运营成本。


因此，为了避免这种情况，在多个数据库操作之间，我们就需要一种机制来确保这些操作要么全部成功，要么全部失败。而事务正是为了解决这个问题而存在的。

# 2. Java事务管理
在 **Java** 中，我们可以使用 **JDBC** 、**Spring**、**JPA**、**MyBatis（MyBatis Plus）** 等框架来管理数据库事务。这些框架提供了丰富的 **API** 和工具，使我们能够轻松地管理事务。

## 2.1 JDBC 的事务管理
在 **JDBC** 中，我们可以通过设置 `Connection` 对象的 `autoCommit` 属性来开启或关闭事务。
- 当 `autoCommit` 为 `true` 时，每次执行 **SQL** 语句都会自动提交事务；
- 当 `autoCommit` 为 `false` 时，我们需要手动调用 `commit()` 或 `rollback()` 方法来提交或回滚事务。

我们来看一个简单的 JDBC 事务管理示例：

```java
    @Test
    public void testJDBCTransaction() throws Exception {
        // flea-config.xml 中配置
        FleaJDBCConfig.init(DBSystemEnum.MySQL.getName(), "fleajpatest");

        Connection conn = null;
        PreparedStatement pstmt1 = null;
        PreparedStatement pstmt2 = null;
        try {
            conn = FleaJDBCConfig.getConfig().getConnection();
            // 关闭自动提交
            conn.setAutoCommit(false);

            // 执行第一条SQL语句
            String sql1 = "UPDATE student SET stu_age = stu_age-10 WHERE stu_name='huazie'";
            pstmt1 = conn.prepareStatement(sql1);
            pstmt1.executeUpdate();

            LOGGER.debug("执行第一条SQL语句");
            // 模拟异常
            throwEx();

            // 执行第二条SQL语句
            String sql2 = "UPDATE student SET stu_age = stu_age+12 WHERE stu_name='huazie'";
            pstmt2 = conn.prepareStatement(sql2);
            pstmt2.executeUpdate();

            LOGGER.debug("执行第二条SQL语句");
            // 提交事务
            conn.commit();

            LOGGER.debug("提交事务");
        } catch (SQLException e) {
            if (null != conn) {
                try {
                    // 回滚事务
                    conn.rollback();
                    LOGGER.debug("回滚事务");
                } catch (SQLException ex) {
                }
            }
        } finally {
            // 关闭资源
            if (null != pstmt2) pstmt2.close();
            if (null != pstmt1) pstmt1.close();
            if (null != conn) conn.close();
        }

    }

    private void throwEx() throws SQLException {
        throw new SQLException("Test Exception");
    }
```

在开始测试之前，先来看看要操作的数据【如下红框所示】：

![](student_data.png)

上述示例中，执行完第一条语句之后，模拟抛出异常，我们来运行一下，可以看到如下结果：

![](result-jdbc.png)

同时检查数据库中发现数据没变，说明第一条 **SQL** 语句也没有执行成功，可见已经被回滚了。

现在，我们将模拟异常的代码去除，然后断点调试下看下两条语句执行完，但还没有提交的情况，如下：

![](result-jdbc-1.png)

从上述截图，可知连接在提交之前，数据还是原来的。

继续执行，连接被提交，我们可以看到 `stu_age` 已经被更改为 **22**，可见示例中的两条 **SQL** 语句已经执行成功。

![](result-jdbc-2.png)

## 2.2 Spring 事务管理

**Spring** 提供了强大的声明式事务管理功能，我们可以将事务管理与业务逻辑分离，在不修改业务代码的情况下，为业务方法添加事务支持。

**Spring** 事务管理的核心接口是 `PlatformTransactionManager`，它是所有事务管理器的抽象。**Spring** 提供了多种事务管理器的实现。我们可以通过配置来选择合适的事务管理器，以便于 **Spring** 结合 **JDBC**、**JPA**、**MyBatis（Mybatis Plus）** 等框架来管理事务。

### 2.2.1 Spring + JDBC 

#### 2.2.1.1 添加 Spring 配置

在 Spring 配置文件中，我们需要配置数据源，如下：

```xml
    <bean id="jdbcDataSource" class="org.springframework.jdbc.datasource.DriverManagerDataSource">
        <property name="driverClassName" value="com.mysql.jdbc.Driver" />
        <property name="url" value="jdbc:mysql://localhost:3306/fleajpatest?useUnicode=true&amp;characterEncoding=UTF-8" />
        <property name="username" value="替换成你的MySQL用户名" />
        <property name="password" value="替换成你的MySQL用户密码" />
    </bean>
```

配置 `JdbcTemplate`，用于使用 `JDBC` 操作数据库，如下：

```xml
    <bean id="jdbcTemplate" class="org.springframework.jdbc.core.JdbcTemplate">
        <property name="dataSource" ref="jdbcDataSource"/>
    </bean>
```

配置事务管理器 `DataSourceTransactionManager`，如下：

```xml
    <bean id="jdbcTransactionManager" class="org.springframework.jdbc.datasource.DataSourceTransactionManager">
        <property name="dataSource" ref="jdbcDataSource" />
    </bean>
```

最后，在 **Spring** 配置文件中，我们需要开启对应的事务管理功能，如下：

```xml
    <tx:annotation-driven transaction-manager="jdbcTransactionManager"/>
```

#### 2.2.1.2 添加业务代码并测试验证

下面让我们添加业务代码，来演示下具体的事务管理，如下：

```java
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class StudentService {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Transactional("jdbcTransactionManager")
    public void service() throws SQLException {

        // 插入一条记录
        String sql = "insert into student(stu_name, stu_age, stu_sex, stu_state) values(?, ?, ?, ?)";
        jdbcTemplate.update(sql, "LGH", "18", 1, 1);

        sql = "update student set stu_state = ? where stu_name = ?";
        jdbcTemplate.update(sql, 0, "LGH");

        throwEx();
    }

    private void throwEx() throws SQLException {
        throw new SQLException("Test Exception");
    }
}
```

需要注意的是，这里抛出的异常还是 **SQLException**，它是 `java.lang.Exception` 的子类。

继续添加测试代码如下：

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class StudentSpringJDBCTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(StudentSpringJDBCTest.class);

    @Resource(name = "studentService")
    private StudentService studentService;

    @Test
    public void testSpringJDBCTransaction() throws Exception {
        studentService.service();
    }
}
```

现在运行上述测试代码，结果如下：

![](result-spring-jdbc.png)

从上图可知，虽然这里抛出了异常，但是实际数据还是提交成功了。

**那这里是什么原因呢？** 

这就要提到刚才抛出的异常了【即 `SQLException`】，Spring 事务管理能处理的异常一定要是`RuntimeException及其子类` 或者 `Error及其子类`，否则事务无法回滚。

可见如下截图【这块后续有时间展开讲解下】：

![](rollbackOn.png)
![](rollbackOn-1.png)


我们重新来修改一下业务代码，抛出 `RuntimeException` 异常 ：

```java
    @Transactional("jdbcTransactionManager")
    public void service() throws RuntimeException {

        // 省略。。。

        throwEx();
    }

    private void throwEx() throws RuntimeException {
        throw new RuntimeException("Test Exception");
    }
```

继续运行测试代码，执行结果如下：

![](result-spring-jdbc-1.png)

可以看到，除了之前新增的，本次要处理的 **SQL** 并没有执行到数据库。

现在，让我们把抛出异常的代码注释掉，运行结果如下：

![](result-spring-jdbc-2.png)

从上截图中可知，要处理的 **SQL** 都已经成功执行。

### 2.2.2 Spring + JPA

#### 2.2.2.1 JPA相关依赖

**Spring Data JPA 依赖**

```xml
<dependency>
    <groupId>org.springframework.data</groupId>
    <artifactId>spring-data-jpa</artifactId>
    <version>2.5.0</version>
</dependency>
```

**EclipseLink 的 JPA 实现依赖**
```xml
<dependency>
    <groupId>org.eclipse.persistence</groupId>
    <artifactId>eclipselink</artifactId>
    <version>2.5.0</version>
</dependency>
```

#### 2.2.2.2 添加 JPA 配置

在资源文件夹的 `META-INF` 目录下添加 `fleajpa-persistence.xml`，配置持久化单元的内容如下：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<persistence version="2.0" xmlns="http://java.sun.com/xml/ns/persistence" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://java.sun.com/xml/ns/persistence http://java.sun.com/xml/ns/persistence/persistence_2_0.xsd">

    <persistence-unit name="fleajpa" transaction-type="RESOURCE_LOCAL">
        <!-- provider -->
        <provider>org.eclipse.persistence.jpa.PersistenceProvider</provider>
        <!-- Connection JDBC -->
        <class>com.huazie.fleadbtest.jpa.common.entity.Student</class>
        <exclude-unlisted-classes>true</exclude-unlisted-classes>

        <properties>
            <property name="javax.persistence.jdbc.driver" value="com.mysql.jdbc.Driver" />
            <property name="javax.persistence.jdbc.url"
            value="jdbc:mysql://localhost:3306/fleajpatest?useUnicode=true&amp;characterEncoding=UTF-8" />
            <property name="javax.persistence.jdbc.user" value="替换成你的MySQL用户名" />
            <property name="javax.persistence.jdbc.password" value="替换成你的MySQL用户密码" />
        </properties>
    </persistence-unit>

</persistence>

```

#### 2.2.2.3 添加 Spring 配置
在 **Spring** 配置文件中，添加 **JPA** 相关的默认 **Bean**，用来初始化 `LocalContainerEntityManagerFactoryBean`，如下：

```xml
    <bean id="defaultPersistenceManager"
          class="org.springframework.orm.jpa.persistenceunit.DefaultPersistenceUnitManager">
        <property name="persistenceXmlLocations">
            <!-- 可以配置多个持久单元 -->
            <list>
                <value>classpath:META-INF/fleajpa-persistence.xml</value>
            </list>
        </property>
    </bean>

    <bean id="defaultPersistenceProvider" class="org.eclipse.persistence.jpa.PersistenceProvider"/>

    <!--<bean id="defaultLoadTimeWeaver" class="org.springframework.instrument.classloading.InstrumentationLoadTimeWeaver"/>-->

    <bean id="defaultVendorAdapter" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaVendorAdapter">
        <property name="showSql" value="true"/>
    </bean>

    <bean id="defaultJpaDialect" class="org.springframework.orm.jpa.vendor.EclipseLinkJpaDialect"/>

```


配置 `LocalContainerEntityManagerFactoryBean`，它是 **Spring Data JPA** 中用于创建和管理 `EntityManagerFactory` 的一个核心类，如下所示：

```xml
    <bean id="fleaJpaEntityManagerFactory"
          class="org.springframework.orm.jpa.LocalContainerEntityManagerFactoryBean">
        <property name="persistenceUnitManager" ref="defaultPersistenceManager"/>
        <property name="persistenceUnitName" value="fleajpa"/>
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
```

配置 **JPA** 事务管理器 `JpaTransactionManager`，如下：

```xml
    <bean id="fleaJpaTransactionManager" class="org.springframework.orm.jpa.JpaTransactionManager">
        <property name="entityManagerFactory" ref="fleaJpaEntityManagerFactory"/>
    </bean>
```

启用 **Spring Data JPA** 仓库扫描，如下：

```java
    <jpa:repositories base-package="com.huazie.fleadbtest.jpa.repository"
                      entity-manager-factory-ref="fleaJpaEntityManagerFactory"
                      transaction-manager-ref="fleaJpaTransactionManager"/>
```



最后，在 **Spring** 配置文件中，还需要开启对应的事务管理功能，如下：

```xml
    <tx:annotation-driven transaction-manager="fleaJpaTransactionManager"/>
```
#### 2.2.2.4 添加实体类
新建如下实体类 `Student`，对应测试表 **student**

```java
@Entity
@Table(name = "student")
public class Student implements FleaEntity {

    private static final long serialVersionUID = 1267943552214677159L;

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "STUDENT_SEQ")
    @SequenceGenerator(name = "STUDENT_SEQ")
    @Column(name = "stu_id", unique = true, nullable = false)
    private Long stuId; // 学生编号

    @Column(name = "stu_name", nullable = false)
    private String stuName; // 学生姓名

    @Column(name = "stu_age", nullable = false)
    private Integer stuAge; // 学生年龄

    @Column(name = "stu_sex", nullable = false)
    private Integer stuSex; // 学生性别（1：男 2：女）

    @Column(name = "stu_state", nullable = false)
    private Integer stuState; // 学生状态（0：删除 1：在用）

  // ... 省略get和set方法
    
    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }
}
```
#### 2.2.2.5 添加业务代码并测试验证

首先，添加 `StudentRepository`，用于访问 `Student` 实体相关的数据库操作的接口

```java
public interface StudentRepository extends JpaRepository<Student, Long> {
}
```

接着，添加如下业务代码：

```java
@Service
public class StudentService {

    @Autowired
    private StudentRepository studentRepository;

    @Transactional("fleaJpaTransactionManager")
    public void service() throws RuntimeException {

        Student student = new Student();
        student.setStuName("杜甫");
        student.setStuAge(35);
        student.setStuSex(1);
        student.setStuState(1);

        studentRepository.save(student);

        student.setStuState(0);
        studentRepository.save(student);

        throwEx();
    }

    private void throwEx() throws RuntimeException {
        throw new RuntimeException("Test Exception");
    }
}
```

最后，添加自测类，如下：

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class StudentSpringJPATest {

    @Resource(name = "studentService")
    private StudentService studentService;

    @Test
    public void testSpringJDBCTransaction() throws Exception {
        studentService.service();
    }
}
```

一切准备就绪，我们来运行上述自测类，结果如下：

![](result-spring-jpa.png)

从上述截图可见，业务代码抛出了异常，相关的数据库操作并未执行成功。【**注意：** 这里抛出的异常也一定要是`RuntimeException及其子类` 或者 `Error及其子类`】

将抛出异常的代码注释掉，再来运行一下看看，如下：

![](result-spring-jpa-1.png)

从上图中，我们可以看出业务代码的数据操作已经成功执行到了数据库中。

### 2.2.3 Spring + Mybatis Plus

#### 2.2.3.1 涉及依赖

**alibaba 的数据库连接池 druid 依赖：**

```xml
<dependency>
    <groupId>com.alibaba</groupId>
    <artifactId>druid</artifactId>
    <version>1.2.0</version>
</dependency>
```

**mybatis 依赖：**

```xml
<dependency>
  <groupId>org.mybatis</groupId>
  <artifactId>mybatis</artifactId>
  <version>3.5.9</version>
</dependency>
```

**Mybatis Plus 依赖：**

```xml
<dependency>
    <groupId>com.baomidou</groupId>
    <artifactId>mybatis-plus</artifactId>
    <version>3.5.1</version>
</dependency>
```

#### 2.2.3.2 添加 Spring 配置

在 Spring 配置文件中，添加数据源配置，如下：

```xml
  <bean id="dataSource" class="com.alibaba.druid.pool.DruidDataSource" destroy-method="close">
        <property name="driverClassName" value="com.mysql.jdbc.Driver" />
        <property name="url" value="jdbc:mysql://localhost:3306/fleajpatest?useUnicode=true&amp;characterEncoding=UTF-8" />
        <property name="username" value="替换成你的MySQL用户名" />
        <property name="password" value="替换成你的MySQL用户密码" />
    </bean>
```

配置 `SqlSessionFactory`，如下所示：

```xml
  <bean id="sqlSessionFactory" class="com.baomidou.mybatisplus.extension.spring.MybatisSqlSessionFactoryBean">
        <property name="dataSource" ref="dataSource"/>
        <property name="typeAliasesPackage" value="com.huazie.fleadbtest.mybatisplus.entity"/>
    </bean>
```

配置 `MapperScan`，如下所示：

```xml
  <bean class="org.mybatis.spring.mapper.MapperScannerConfigurer">
        <property name="basePackage" value="com.huazie.fleadbtest.mybatisplus.mapper"/>
    </bean>
```

配置事务管理器 `DataSourceTransactionManager`，如下：

```xml
  <bean id="dataSourceTransactionManager" class="org.springframework.jdbc.datasource.DataSourceTransactionManager">
        <property name="dataSource" ref="dataSource"/>
    </bean>
```

最后，在 **Spring** 配置文件中，还需要开启对应的事务管理功能，如下：

```xml
    <tx:annotation-driven transaction-manager="dataSourceTransactionManager"/>
```

#### 2.2.3.3 添加实体类
新建如下实体类 `Student`，对应测试表 **student**

```java
@Data
@TableName("student")
public class Student {

    @TableId(value = "stu_id", type = IdType.AUTO)
    private Long stuId; // 学生编号

    @TableField(value = "stu_name")
    private String stuName; // 学生姓名

    @TableField(value = "stu_age")
    private Integer stuAge; // 学生年龄

    @TableField(value = "stu_sex")
    private Integer stuSex; // 学生性别（1：男 2：女）

    @TableField(value = "stu_state")
    private Integer stuState; // 学生状态（0：删除 1：在用）

}
```

#### 2.2.2.4 添加业务代码并测试验证

添加学生服务层接口类，如下：

```java
public interface IStudentService extends IService<Student> {
    void service() throws RuntimeException;
}
```

添加学生服务层实现类，如下：

```java
@Service("studentServiceImpl")
public class StudentServiceImpl extends ServiceImpl<StudentMapper, Student> implements IStudentService {

    @Override
    @Transactional("dataSourceTransactionManager")
    public void service() throws RuntimeException {
        Student student = new Student();
        student.setStuName("李白");
        student.setStuAge(25);
        student.setStuSex(1);
        student.setStuState(1);

        baseMapper.insert(student);

        student.setStuState(0);

        baseMapper.updateById(student);

        throwEx();
    }

    private void throwEx() throws RuntimeException {
        throw new RuntimeException("Test Exception");
    }

}
```

最后，添加自测类，如下：

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class StudentServiceSpringTest {

    @Resource(name = "studentServiceImpl")
    private IStudentService studentService;

    @Test
    public void testSpringMybatisPlusTransaction() throws RuntimeException {
        studentService.service();
    }
}
```

现在，我们可以运行上述自测类，结果如下：

![](result-spring-mybatisplus.png)

从上述截图可见，自定义异常已经被抛出，并且数据库中也没有执行成功，说明事务已经回滚了。

现在将抛出异常的代码注释掉，再来运行看看，如下：

![](result-spring-mybatisplus-1.png)

从上图可知，相关数据库操作已经成功执行。

# 示例

相关演示示例请查看 **GitHub** 上的 [flea-db-test](https://github.com/huazie/flea-db-test) 测试项目。
# 总结

本文 **Huazie** 从 **Java** 事务的基础概念出发，带大家通过不同的事务管理方式进行实践，进一步加深了对 **Java** 事务的理解。



