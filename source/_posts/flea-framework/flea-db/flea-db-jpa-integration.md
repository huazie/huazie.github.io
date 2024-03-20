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

[《开发框架-Flea》](/categories/开发框架-Flea/)

![](/images/flea-logo.png)

# 引言

本节内容需要了解JPA封装内容，请参见笔者上篇博文[《JPA封装介绍》](/2019/09/06/flea-framework/flea-db/flea-db-jpa-introduction/)。

# 1. 准备工作
为了演示JPA接入，需要先准备如下：
1. MySQL数据库 (客户端可以使用 navicat for mysql)
2. 新建测试数据库 **fleajpatest**

![](fleajpatest.png)

3. 新建测试表 **student**

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
# 2. 接入讲解
## 2.1 实体类
新建如下实体类Student，对应测试表student

```java
/**
 * <p> 学生表对应的实体类 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
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
## 2.2 持久化单元DAO层实现
上篇博文说到，增加一个持久化单元配置，便需要增加一个持久化单元DAO层实现。针对本次演示新增持久化单元 **fleajpa**，持久化配置文件 **fleajpa-persistence.xml**, Spring配置中新增数据库事务管理者配置，相关内容可参考上一篇博文。下面贴出本次演示的持久化单元DAO层实现代码：

```java
/**
 * <p> FleaJpa数据源DAO层父类 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class FleaJpaDAOImpl<T> extends AbstractFleaJPADAOImpl<T> {
    // 持久化单元，持久化配置文件中定义，spring配置中持久化接口工厂初始化参数
    @PersistenceContext(unitName="fleajpa")
    protected EntityManager entityManager;

    @Override
    // 持久化事务管理者， spring配置文件中定义
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
## 2.3 配置介绍
详细配置信息，可以参考笔者上篇博文，这里不再赘述。
涉及文件 **fleajpa-persistence.xml**  和 **applicationContext.xml**

## 2.4 学生DAO层接口
**IStudentDAO** 继承了抽象Flea JPA DAO层接口，并定义了两个方法，分别获取学生信息列表（分页）和学生总数。
```java
/**
 * <p> 学生DAO层接口 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public interface IStudentDAO extends IAbstractFleaJPADAO<Student> {

    /**
     * <p> 学生信息列表 (分页) </p>
     *
     * @param name      学生姓名，可以模糊查询
     * @param sex       性别
     * @param minAge    最小年龄
     * @param maxAge    最大年龄
     * @param pageNum   查询页
     * @param pageCount 每页总数
     * @return 学生信息列表
     * @throws DaoException 数据操作层异常
     * @since 1.0.0
     */
    List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException;

    /**
     * <p> 学生总数 </p>
     *
     * @param name   学生姓名，可以模糊查询
     * @param sex    性别
     * @param minAge 最小年龄
     * @param maxAge 最大年龄
     * @return 学生总数
     * @throws DaoException 数据操作层异常
     * @since 1.0.0
     */
    int getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException;

}
```
## 2.5 学生DAO层实现
**StudentDAOImpl** 是学生信息的数据操作层实现，继承持久化单元DAO层实现类，并实现了上述学生DAO层接口自定义的两个方法。 具体如何使用 **FleaJPAQuery** 可以参见下面代码 ：
```java
/**
 * <p> 学生DAO层实现类 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
@Repository("studentDAO")
public class StudentDAOImpl extends FleaJpaDAOImpl<Student> implements IStudentDAO {

    @Override
    @SuppressWarnings(value = "unchecked")
    public List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException {
        FleaJPAQuery query = initQuery(name, sex, minAge, maxAge, null);
        List<Student> studentList;
        if (pageNum > 0 && pageCount > 0) {
            // 分页查询
            studentList = query.getResultList((pageNum - 1) * pageCount, pageCount);
        } else {
            // 全量查询
            studentList = query.getResultList();
        }
        return studentList;
    }

    @Override
    @SuppressWarnings(value = "unchecked")
    public long getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException {
        FleaJPAQuery query = initQuery(name, sex, minAge, maxAge, Long.class);
        // 统计数目
        query.countDistinct();
        Object result = query.getSingleResult();
        return Long.parseLong(StringUtils.valueOf(result));
    }

    private FleaJPAQuery initQuery(String name, Integer sex, Integer minAge, Integer maxAge, Class<?> result) throws DaoException{
        FleaJPAQuery query = getQuery(result);
        // 拼接 查询条件
        // 根据姓名 模糊查询, attrName 为 实体类 成员变量名，并非表字段名
        if (StringUtils.isNotEmpty(name)) {
            query.like("stuName", name);
        }
        // 查询性别
        if (ObjectUtils.isNotEmpty(sex)) {
            query.equal("stuSex", sex);
        }
        // 查询年龄范围
        if (ObjectUtils.isNotEmpty(minAge)) {
            // 大于等于
            query.ge("stuAge", minAge);
        }
        if (ObjectUtils.isNotEmpty(maxAge)) {
            // 小于等于
            query.le("stuAge", maxAge);
        }
        return query;
    }
}
```
## 2.6 学生SV层接口
**IStudentSV** 继承抽象Flea JPA SV层接口，并定义两个方法，分别获取学生信息列表（分页）和学生总数。
```java
/**
 * <p> 学生SV层接口定义 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public interface IStudentSV extends IAbstractFleaJPASV<Student> {

    /**
     * <p> 学生信息列表 (分页) </p>
     *
     * @param name      学生姓名，可以模糊查询
     * @param sex       性别
     * @param minAge    最小年龄
     * @param maxAge    最大年龄
     * @param pageNum   查询页
     * @param pageCount 每页总数
     * @return 学生信息列表
     * @throws DaoException 数据操作层异常
     * @since 1.0.0
     */
    List<Student> getStudentList(String name, Integer sex, Integer minAge, Integer maxAge, int pageNum, int pageCount) throws DaoException;

    /**
     * <p> 学生总数 </p>
     *
     * @param name   学生姓名，可以模糊查询
     * @param sex    性别
     * @param minAge 最小年龄
     * @param maxAge 最大年龄
     * @return 学生总数
     * @throws DaoException 数据操作层异常
     * @since 1.0.0
     */
    long getStudentCount(String name, Integer sex, Integer minAge, Integer maxAge) throws DaoException;

}
```
## 2.7 学生SV层实现
**StudentSVImpl** 继承抽象Flea JPA SV层实现类，并实现了上述学生SV层接口的两个自定义方法。具体实现参见如下代码：
```java
/**
 * <p> 学生SV层实现类 </p>
 *  * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
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
    // 可参见AbstractFleaJPASVImpl里的实现
    @Override
    protected IAbstractFleaJPADAO<Student> getDAO() {
        return studentDao;
    }
}
```
## 2.8 JPA接入自测

### 2.8.1 新增学生信息
```java
    private ApplicationContext applicationContext;

    @Before
    public void init() {
        applicationContext = new ClassPathXmlApplicationContext("applicationContext.xml");
        LOGGER.debug("ApplicationContext={}", applicationContext);
    }
    
    @Test
    public void testInsertStudent() {
        try {
            IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");
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
执行结果：
![](result_add.png)
### 2.8.2 更新学生信息
 
```java
     @Test
    public void testStudentUpdate() {
        try {
            IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");
            // 根据主键查询学生信息
            Student student = studentSV.query(3L);
            LOGGER.debug("Before : {}", student);
            student.setStuName("王三麻子");
            student.setStuAge(19);
            // 更新学生信息
            studentSV.update(student);
            // 最后再根据主键查询学生信息
            student = studentSV.query(3L);
            LOGGER.debug("After : {}", student);
        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```
运行结果：
![](result_update.png)
### 2.8.3 删除学生信息
 
```java
    @Test
    public void testStudentDelete() {
        try {
            IStudentSV studentSV = (IStudentSV) applicationContext.getBean("studentSV");
            // 根据主键查询学生信息
            Student student = studentSV.query(3L);
            LOGGER.debug("Before : {}", student);
            // 删除学生信息(里面会先去将学生实体信息查出来，然后再删除)
            studentSV.remove(3L);
            // 最后再根据主键查询学生信息
            student = studentSV.query(3L);
            LOGGER.debug("After : {}", student);
        } catch (Exception e) {
            LOGGER.error("Exception : ", e);
        }
    }
```
运行结果：
![](result_delete.png)

### 2.8.4 查询学生信息（按条件分页查询）
 表里自行再插入些数据，用于测试查询，查询结果因各自表数据而异；
 目前我表中数据如下：
![](https://img-blog.csdnimg.cn/20190912114155255.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTI4NTUyMjk=,size_16,color_FFFFFF,t_70)
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
运行结果：
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
运行结果：
![](result_query_count.png)

# 总结
看到这里，我们的JPA接入工作已经成功完成，本demo工程可以移步到GitHub查看 [flea-jpa-test](https://github.com/Huazie/flea-db-test/tree/main/flea-jpa-test)。

在JPA封装介绍博文中，针对Flea JPA查询对象还存在的一个并发问题，将在后续的博文《flea-db使用之基于对象池的FleaJPAQuery》中介绍。
 

