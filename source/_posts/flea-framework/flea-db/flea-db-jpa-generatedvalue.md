---
title: JPA主键生成策略介绍
date: 2019-09-03 20:14:12
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - GeneratedValue
  - GenerationType
  - JPA主键生成策略
---

![](/images/jpa-logo.png)

# 引言
接入JPA框架之前，我们有必要了解一下JPA的主键生成策略。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1.  依赖
```xml
<dependency>
    <groupId>org.eclipse.persistence</groupId>
    <artifactId>javax.persistence</artifactId>
    <version>2.1.0</version>
</dependency>
```
# 2. GeneratedValue注解
[GeneratedValue](https://docs.oracle.com/javaee/7/api/javax/persistence/GeneratedValue.html) 是JPA主键生成策略中的一个非常重要的注解。它提供主键值生成策略的规范，可以与 [Id](https://docs.oracle.com/javaee/7/api/javax/persistence/Id.html) 注解一起应用于实体或映射超类的主键属性或字段；它只支持简单的主键，派生的主键不支持使用 。

```java
@Target({METHOD, FIELD})
@Retention(RUNTIME)
public @interface GeneratedValue {
    GenerationType strategy() default AUTO;
    String generator() default "";
}
```
如上代码所示，**GeneratedValue** 注解有 **strategy** 和 **generator** 两个成员变量。

## 2.1 主键生成策略【strategy】
持久化提供程序必须使用主键生成策略来生成被注解的实体的主键。这是一个可选项，默认是 **GenerationType.AUTO**；
**strategy** 的值是主键生成策略枚举类型 [GenerationType](https://docs.oracle.com/javaee/7/api/javax/persistence/GenerationType.html)，包含4个枚举值：【**TABLE**，**SEQUENCE**，**IDENTITY**，**AUTO**】。

## 2.2 主键生成器【generator】
**generator** 指定使用的主键生成器的名称，有 [SequenceGenerator](https://docs.oracle.com/javaee/7/api/javax/persistence/SequenceGenerator.html) 或 [TableGenerator](https://docs.oracle.com/javaee/7/api/javax/persistence/TableGenerator.html) 注解。它为持久化提供程序提供 **ID** 生成器。这也是一个可选项，默认可空。

# 3. GenerationType
[GenerationType](https://docs.oracle.com/javaee/7/api/javax/persistence/GenerationType.html) 定义主键生成策略的类型。包含如下：
## 3.1 GenerationType.TABLE
**TABLE** 指示持久化提供程序必须使用基础数据库表为实体分配主键，以确保唯一性。它的好处是不依赖于具体数据库的实现，代码可移植性高，但由于某些数据库的特性【如主键自增长，序列等等】未能使用到，不推荐优先使用，可作为折中方案。

### 3.1.1 具体用法

```java
    @Id
    @GeneratedValue(strategy = GenerationType.TABLE, generator = "FLEA_LOGIN_LOG_GENERATOR")
    @TableGenerator(
        name = "FLEA_LOGIN_LOG_GENERATOR",
        table = "flea_id_generator",
        catalog = "",
        schema = "",
        pkColumnName = "id_generator_key",
        valueColumnName = "id_generator_value",
        pkColumnValue = "pk_flea_login_log_(CREATE_DATE)",
        initialValue = 0,
        allocationSize = 1,
        uniqueConstraints = {},
        indexes = {}
    )
    @Column(name = "login_log_id", unique = true, nullable = false)
    private Long loginLogId; // 登录日志编号
```

- `name` ：唯一的生成器名称，可以由一个或多个类引用以作为id值的生成器。
- `table` ：【可选】存储生成的ID值的表的名称，默认为持久化提供程序选择的名称。
- `catalog` ：【可选】生成器表所属的数据库目录。
- `schema` ：【可选】生成器表所属的数据库结构。
- `pkColumnName` ：【可选】表中主键列的名称，默认为持久化提供程序选择的名称。
- `valueColumnName` ：【可选】存储最后生成的主键值的列的名称，默认为持久化提供程序选择的名称。
- `pkColumnValue` ：【可选】ID生成器表中的主键值模板，用于将该生成值集与其他可能存储在表中的值区分开；默认为持久化提供程序选择的值，用以存储在生成器表的主键列中。
- `initialValue` ：【可选】用于初始化存储最后生成的值的列的初始值，默认值为 0
- `allocationSize` ：【可选】从ID生成器表中分配ID号时增加的数量, 默认值为 50
- `uniqueConstraints` ：【可选】将在表上放置的其他唯一约束，仅当表生成有效时才使用它们；除了主键约束之外，还应用了这些约束；默认为无其他约束。
- `indexes` ：【可选】表的索引，仅当表生成有效时才使用它们；请注意，不必为主键指定索引，因为主键索引将自动创建。

### 3.1.2 TableGenerator 注解源码

```java
@Target({TYPE, METHOD, FIELD}) 
@Retention(RUNTIME)
public @interface TableGenerator {
    String name();
    String table() default "";
    String catalog() default "";
    String schema() default "";
    String pkColumnName() default "";
    String valueColumnName() default "";
    String pkColumnValue() default "";
    int initialValue() default 0;
    int allocationSize() default 50;
    UniqueConstraint[] uniqueConstraints() default {};
    Index[] indexes() default {};
}
```

`TableGenerator` 定义了一个主键生成器，可以通过名称引用，当在 `GeneratedValue` 注解中指定一个生成器元素时使用。 **表生成器** 可以在实体类或主键字段/属性上指定。生成器名称的作用范围是持久性单元全局的（跨所有生成器类型）。

- `String name()` ：必填项，表示唯一的生成器名称，可以被一个或多个类引用，用于生成id值。
- `String table()` ：可选项，存储生成的id值的表的名称，默认为持久性提供程序选择的名称。
- `String catalog()` ：可选项，表所在的目录名称，默认为默认目录。
- `String schema()` ：可选项，表所在的模式名称，默认为用户默认的模式。
- `String pkColumnName()` ：可选项，表中主键列的名称，默认为提供程序选择的名称。
- `String valueColumnName()` ：可选项，存储最后生成的值的列的名称，默认为提供程序选择的名称。
- `String pkColumnValue()` ：可选项，在生成器表中区分此生成的值集合与可能存储在表中的其他值集合的主键值。默认为提供程序选择的值，以存储在生成器表的主键列中。
- `int initialValue()` ：可选项，用于初始化存储最后生成的值的列的初始值。
- `int allocationSize()` ：可选项，从生成器分配id号码时每次递增的数量。
- `UniqueConstraint[] uniqueConstraints()` ：可选项，要放置在表上的唯一约束条件。仅在表生成器生效时使用。这些约束条件适用于主键约束之外。
- `Index[] indexes()` ：可选项，表的索引。仅在表生成器生效时使用。请注意，对于主键，不必指定索引，因为主键索引将自动创建。


## 3.2 GenerationType.SEQUENCE

**SEQUENCE** 指示持久化提供程序必须使用数据库序列为实体分配主键。该策略只适用于部分支持 **序列** 的数据库系统，比如 **Oracle**。

### 3.2.1 具体用法

```java
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE, generator="PARA_DETAIL_SEQ")
    @SequenceGenerator(
        name="PARA_DETAIL_SEQ", 
        sequenceName="PARA_ID_SEQ",
        catalog = "",
        schema = "",
        initialValue = 0,
        allocationSize = 1
    )
    @Column(name = "para_id", unique = true, nullable = false)
    public Long getParaId() { return paraId; }
```
- `name` ：唯一的生成器名称，可以由一个或多个类引用以作为id值的生成器。
- `sequenceName` ：【可选】主键值对应的数据库序列对象的名称。默认为提供商选择的值。
- `catalog` ：【可选】生成器表所属的数据库目录
- `schema` ：【可选】生成器表所属的数据库结构
- `initialValue` ：【可选】用于初始化存储最后生成的值的列的初始值，默认值为 0
- `allocationSize` ：【可选】从ID生成器表中分配ID号时增加的数量, 默认值为 50


### 3.2.2 SequenceGenerator 注解源码

```java
@Target({TYPE, METHOD, FIELD}) 
@Retention(RUNTIME)
public @interface SequenceGenerator {
    String name();
    String sequenceName() default "";
    String catalog() default "";
    String schema() default "";
    int initialValue() default 1;
    int allocationSize() default 50;
}
```
`SequenceGenerator` 同样定义了一个主键生成器，可以通过名称引用，当在 `GeneratedValue` 注解中指定一个生成器元素时。**序列生成器** 可以在实体类或主键字段或属性上指定。生成器名称的范围是持久单元全局的（跨所有生成器类型）。


- `String name()` ：（必填） 可以被一个或多个类引用的唯一生成器名称，用于主键值的生成器。
- `String sequenceName()` ：（可选）用于获取主键值的数据库序列对象的名称。默认为提供程序选择的值。
- `String catalog()` ：（可选）序列生成器的目录。
- `String schema()` ：（可选）序列生成器的模式。
- `int initialValue()` ：（可选）序列对象开始生成的值。
- `int allocationSize()` ：（可选）从序列分配序列号时要增加的数量。


## 3.3 GenerationType.IDENTITY

**IDENTITY** 指示持久化提供程序必须使用数据库标识列为实体分配主键。该策略只适用于支持 **主键自增长** 的数据库系统，比如 MySQL。

具体用法如下：

```java
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "para_id", unique = true, nullable = false)
    private Long paraId;        // 参数编号
```

## 3.4 GenerationType.AUTO

**AUTO** 指示持久化提供程序应为特定数据库选择适当的策略。 该生成策略可能期望数据库资源存在，或者可能尝试创建一个数据库资源。如果供应商不支持架构生成或无法在运行时创建架构资源，则供应商可能会提供有关如何创建此类资源的文档。

```java
    @Id
    @GeneratedValue
    @Column(name = "para_id", unique = true, nullable = false)
    private Long paraId;        // 参数编号
```

# 4. 各数据库对比

支持对应主键生成策略类型的打 √ ，不支持的打 ×，如下所示：

|            | TABLE | SEQUENCE | IDENTITY | AUTO |
| ---------- | ----- | -------- | -------- | ---- |
| MySQL      | √     | ×        | √        | √    |
| Oracle     | √     | √        | ×        | √    |
| PostgreSQL | √     | √        | √        | √    |
| SQL Server | √     | √        | √        | √    |
| DB2        | √     | √        | √        | √    |

# 总结
本篇我们介绍了 `JPA` 主键生成策略，下一篇基于 `GenerationType.TABLE` 的主键生成器表介绍，敬请期待！！！
