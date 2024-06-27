---
title: flea-db使用之JPA接入
date: 2019-09-12 14:25:48
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - JPA接入
---

![](/images/jpa-logo.png)

# 引言

本节内容需要了解 **JPA** 封装内容，请参见笔者上篇博文[《JPA封装介绍》](../../../../../../2019/09/06/flea-framework/flea-db/flea-db-jpa-introduction/)。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 准备工作

为了演示 **JPA接入**，需要先准备如下：

*   **MySQL** 数据库 (客户端可以使用 **navicat for mysql**)
*   新建测试数据库 **fleajpatest**
![](fleajpatest.png)

*   新建测试表 **student**

    ![](student.png)

    建表语句如下：

    ```sql
    DROP TABLE IF EXISTS `student`;
    CREATE TABLE `student` (
      `stu_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '学生编号',
      `stu_name` varchar(255) NOT NULL COMMENT '学生姓名',
      `stu_age` tinyint(2) NOT NULL COMMENT '学生年龄',
      `stu_sex` tinyint(1) NOT NULL COMMENT '学生性别（1：男 2：女）',
      `stu_state` tinyint(2) NOT NULL COMMENT '学生状态（0：删除 1：在用）',
      PRIMARY KEY (`stu_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
    ```

*   依赖
    **MySQL** 的 **JDBC** 驱动 [mysql-connector-java-5.1.25.jar](https://mvnrepository.com/artifact/mysql/mysql-connector-java/5.1.25)

    ```xml
    <dependency>
        <groupId>mysql</groupId>
        <artifactId>mysql-connector-java</artifactId>
        <version>5.1.25</version>
    </dependency>
    ```

    **FLEA DB ECLIPSELINK**【这里为 **flea-framework** 项目下的包，目前版本 **2.0.0**，暂未发布】

    ```xml
    <dependency>
        <groupId>com.huazie.fleaframework</groupId>
        <artifactId>flea-db-eclipselink</artifactId>
        <version>2.0.0</version>
    </dependency>
    ```

# 2. 接入讲解

## 2.1 实体类

新建如下学生表对应的实体类 `Student`，对应测试表 `student`。

```java
@Entity
@Table(name = "student")
public class Student implements FleaEntity {

    private static final long serialVersionUID = 1267943552214677159L;

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "STUDENT_SEQ")
    @SequenceGenerator(name = "STUDENT_SEQ")
    @Column(name = "stu_id", unique = true, nullable = false)
    private Long stuId;

    @Column(name = "stu_name", nullable = false)
    private String stuName;

    @Column(name = "stu_age", nullable = false)
    private Integer stuAge; 

    @Column(name = "stu_sex", nullable = false)
    private Integer stuSex;

    @Column(name = "stu_state", nullable = false)
    private Integer stuState;

    // ... 省略get和set方法
    
    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }
}
```

- `Long stuId` ： 学生编号【主键】。在笔者的[《JPA主键生成策略介绍》](../../../../../../2019/09/03/flea-framework/flea-db/flea-db-jpa-generatedvalue/) 中，介绍了 `GenerationType.IDENTITY`【适用于支持 **主键自增长** 的数据库系统，比如 **MySQL**】，详细内容可自行查看。
- `String stuName` ： 学生姓名【非空】
- `Integer stuAge` ： 学生年龄【非空】
- `Integer stuSex` ： 学生性别（1：男 2：女）【非空】
- `Integer stuState` ： 学生状态（0：删除 1：在用）【非空】

## 2.2 持久化单元DAO层实现

上篇博文说到，增加一个持久化单元配置，便需要增加一个持久化单元 **DAO** 层实现。针对本次演示新增持久化单元 **fleajpa**，持久化配置文件 **fleajpa-persistence.xml**, **Spring** 配置中新增数据库事务管理者配置，相关内容可参考上一篇博文。下面贴出本次演示的持久化单元 **DAO** 层实现代码：

**FleaJpa数据源DAO层父类**

```java
public class FleaJpaDAOImpl<T> extends AbstractFleaJPADAOImpl<T> {
    @PersistenceContext(unitName="fleajpa")
    protected EntityManager entityManager;

    @Override
    @Transactional("fleaJpaTransactionManager")
    public boolean remove(long entityId) throws Exception {
        return super.remove(entityId);
    }

    // ...其他实现省略

    @Override
    protected EntityManager getEntityManager() {
        return entityManager;
    }

}
```

我们来看上面两个注解，分别是：

- `@PersistenceContext(unitName="fleajpa")` ：持久化上下文注解，其值为持久化单元 `unitName` ，在持久化配置文件中定义，**spring** 配置中 **JPA** 实体管理器工厂初始化该参数。
- `@Transactional("fleaJpaTransactionManager")` ：事务注解，其值为持久化事务管理器， 在 **spring** 配置文件中定义。


## 2.3 配置介绍

详细配置信息，可以参考笔者上篇博文，这里不再赘述。
涉及文件 **fleajpa-persistence.xml**  和 **applicationContext.xml**，文章最后会给出示例工程，可自行查看。

## 2.4 学生DAO层接口

**IStudentDAO** 继承了抽象 **Flea JPA DAO** 层接口，并定义了两个方法，分别获取学生信息列表（分页）和学生总数。

**学生DAO层接口**

```java
public interface IStudentDAO extends IAbstractFleaJPADAO<Student> {
    List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException;

    int getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException;

}
```

- `getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount)` ：获取学生信息列表 (分页)，其中参数如下：
    - `name`   ：  学生姓名，可以模糊查询
    - `sex`   ：   性别
    - `minAge`   ： 最小年龄
    - `maxAge`   ： 最大年龄
    - `pageNum`   ： 查询页
    - `pageCount`   ： 每页总数
- `getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge)` ：获取学生总数
    - `name`   ：  学生姓名，可以模糊查询
    - `sex`   ：   性别
    - `minAge`   ： 最小年龄
    - `maxAge`   ： 最大年龄

## 2.5 学生DAO层实现

**StudentDAOImpl** 是学生信息的数据操作层实现，继承持久化单元DAO层实现类，并实现了上述学生DAO层接口自定义的两个方法。 具体如何使用 **FleaJPAQuery** 可以参见下面代码 ：

**学生DAO层实现类**

```java
@Repository("studentDAO")
public class StudentDAOImpl extends FleaJpaDAOImpl<Student> implements IStudentDAO {

    @Override
    @SuppressWarnings(value = "unchecked")
    public List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException {
        FleaJPAQuery query = initQuery(name, sex, minAge, maxAge, null);
        List<Student> studentList;
        if (pageNum > 0 && pageCount > 0) {
            // 分页查询
            studentList = query.getResultList4Page(pageNum, pageCount);
        } else {
            // 全量查询
            studentList = query.getResultList();
        }
        return studentList;
    }

    @Override
    @SuppressWarnings(value = "unchecked")
    public long getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException {
        Object result = initQuery(name, sex, minAge, maxAge, Long.class).count().getSingleResult();
        return Long.parseLong(StringUtils.valueOf(result));
    }

    private FleaJPAQuery initQuery(String name, Integer sex, Integer minAge, Integer maxAge, Class<?> result) throws DaoException{
        return getQuery(result)
                .like("stuName", name)
                .equal("stuSex", sex)
                .ge("stuAge", minAge)
                .le("stuAge", maxAge);
    }
}
```


上述代码，我们需要重点关注 `initQuery` 私有方法，它实际上用于返回一个已经组装好查询条件的 `FleaJPAQuery` 对象：
- `getQuery(result)` ：在[《flea-db使用之JPA封装介绍》](../../../../../../2019/09/06/flea-framework/flea-db/flea-db-jpa-introduction/) 中的抽象 **Flea JPA DAO** 层实现可以看到，通过 **Flea JPA** 查询对象池来获取 `FleaJPAQuery`。
- `like("stuName", name)` ：根据姓名模糊查询, `attrName` 为 实体类对应的成员变量名，并非表字段名
- `equal("stuSex", sex)` ：查询性别（等于）
- `ge("stuAge", minAge)` ：查询年龄范围 (大于等于)
- `le("stuAge", maxAge)` ：查询年龄范围 (小于等于)


另外，我们看到 **学生DAO层实现类** 的上面还有一个注解：

- `@Repository("studentDAO")` ：在 **Spring** 框架中，它是用来标注数据访问层（**DAO层**）的类。**Spring** 会将 `StudentDAOImpl` 实例化为一个名为 `studentDAO` 的Bean；然后，**在其他Bean** 中通过`@Autowired` 或者 `@Resource` 注解来注入这个 **Bean**。

## 2.6 学生SV层接口

**IStudentSV** 继承抽象Flea JPA SV层接口，并定义两个方法，分别获取学生信息列表（分页）和学生总数。

**学生SV层接口**

```java

public interface IStudentSV extends IAbstractFleaJPASV<Student> {
    List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException;
    
    long getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException;
}
```

## 2.7 学生SV层实现

**StudentSVImpl** 继承抽象Flea JPA SV层实现类，并实现了上述学生SV层接口的两个自定义方法。具体实现参见如下代码：

**学生SV层实现类**

```java
@Service("studentSV")
public class StudentSVImpl extends AbstractFleaJPASVImpl<Student> implements IStudentSV {
    // 注入学生DAO层实现类
    @Autowired
    @Qualifier("studentDAO") 
    private IStudentDAO studentDao; 

    @Override
    public List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException {
        return studentDao.getStudentList(name, sex, minAge, maxAge, pageNum, pageCount);
    }

    @Override
    public long getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException {
        return studentDao.getStudentCount(name, sex, minAge, maxAge);
    }
    
    @Override
    protected IAbstractFleaJPADAO<Student> getDAO() {
        return studentDao;
    }
}
```

上述逻辑不复杂，学生SV层实现类中的两个方法均调用 `IStudentDAO` 接口来查询数据库。

上述代码需要关注三个注解 和 一个方法：

- `@Service("studentSV")` ：在 **Spring** 框架中，它通常用于将一个服务层的类标记为 **Spring** 管理的 **Bean**。`studentSV` 是这个 **Bean** 的名称，可以在其他地方通过这个名称来获取这个 **Bean** 的实例。当 **Spring** 容器启动并扫描到带有 `@Service` 注解的类时，它会创建这个类的实例，并将其注册到Spring应用程序上下文中，使得这个 **Bean** 可以被依赖注入（**DI**）到其他组件中。
- `@Autowired` ：在 **Spring** 框架中，它用于自动装配 **Bean**。
- `@Qualifier("studentDAO")` ：它通常与 `@Autowired` 一起使用，用于指定 `@Autowired` 应该注入的特定 **Bean**。这里就指定了一个名为 `studentDAO` 的 **Bean**。如果你不使用 `@Qualifier` 注解，而是仅仅依赖 `@Autowired` 来注入 **Bean**，那么在存在多个相同类型 **Bean** 的情况下，**Spring** 容器将无法确定应该注入哪一个，从而导致 `NoUniqueBeanDefinitionException` 异常。
- `getDAO()` ：在[《flea-db使用之JPA封装介绍》](../../../../../../2019/09/06/flea-framework/flea-db/flea-db-jpa-introduction/) 中的抽象**Flea JPA SV**层实现，可以看到 `getDAO()` 用于通过的一些增删改查操作，实际的实现需要子类来返回对应的 **DAO层** 实现。


## 2.8 JPA接入自测

首先添加自测类，并注入学生服务类

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class StudentTest {

    private static final Logger LOGGER = FleaLoggerProxy.getProxyInstance(StudentTest.class);

    @Resource(name = "studentSV")
    private IStudentSV studentSV;
    
    // 自测代码详见下面
    
}
```

### 2.8.1 新增学生信息

```java
    @Test
    public void testInsertStudent() {
        try {
            Student student = new Student();
            student.setStuName("张三");
            student.setStuAge(18);
            student.setStuSex(1);
            student.setStuState(1);
            studentSV.save(student);

            student = new Student();
            student.setStuName("李四");
            student.setStuAge(19);
            student.setStuSex(1);
            student.setStuState(1);
            studentSV.save(student);

            student = new Student();
            student.setStuName("王二麻子");
            student.setStuAge(20);
            student.setStuSex(1);
            student.setStuState(1);
            studentSV.save(student);

        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```

**执行结果**：

![](result_add.png)

### 2.8.2 更新学生信息
 
```java
    @Test
    public void testStudentUpdate() {
        try {
            Student student = studentSV.query(3L);
            LOGGER.debug("Before : {}", student);
            student.setStuName("王三麻子");
            student.setStuAge(19);
            studentSV.update(student);
            student = studentSV.query(3L);
            LOGGER.debug("After : {}", student);
        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```
上述演示代码逻辑：
- 首先根据指定主键，调用 `query` 方法查询学生信息，并打印 `Before ：XXX`；
- 然后调用 `update` 方法更新该学生信息；
- 最后再根据指定主键，调用 `query` 方法查询学生信息，并打印` After : XXX`。

**运行结果：**

![](result_update.png)

### 2.8.3 删除学生信息
 
```java
    @Test
    public void testStudentDelete() {
        try {
            Student student = studentSV.query(3L);
            LOGGER.debug("Before : {}", student);
            // 
            studentSV.remove(3L);
            // 最后再根据主键查询学生信息
            student = studentSV.query(3L);
            LOGGER.debug("After : {}", student);
        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```

上述演示代码逻辑：
- 首先根据指定主键，调用 `query` 方法查询学生信息，并打印 `Before ：XXX`；
- 然后调用 remove 方法删除指定主键的学生信息(里面会先去将学生实体信息查出来，然后再删除);
- 最后再根据指定主键，调用 `query` 方法查询学生信息，并打印` After : XXX`.


**运行结果：**

![](result_delete.png)

### 2.8.4 查询学生信息（按条件分页查询）
 表里自行再插入些数据，用于测试查询，查询结果因各自表数据而异；

 目前我表中数据如下：

![](student_record.png)

```java
      @Test
    public void testStudentQueryPage() {
        try {

            IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");
            List<Student> studentList = studentSV.getStudentList("张三", 1, 18, 20, 1, 5);
            LOGGER.debug("Student List = {}", studentList);
        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```

**运行结果：**

![](result_query.png)

### 2.8.5 查询学生总数（按条件查询）
 
```java
    @Test
    public void testStudentQueryCount() {
        try {
            IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");
            long count = studentSV.getStudentCount("张三", 1, 18, 20);
            LOGGER.debug("Student Count = {}", count);
        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```

**运行结果：**

![](result_query_count.png)

# 总结
看到这里，我们的 **JPA接入** 工作已经成功完成，本篇演示示例可以移步到 **GitHub** 查看 [flea-jpa-test](https://github.com/Huazie/flea-db-test/tree/main/flea-jpa-test)。

在 **JPA** 封装介绍博文中，针对 **Flea JPA** 查询对象还存在的一个并发问题，将在后续的博文 《flea-db使用之基于对象池的FleaJPAQuery》 中介绍。