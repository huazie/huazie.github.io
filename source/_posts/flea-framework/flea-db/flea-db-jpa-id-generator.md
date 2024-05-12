---
title: flea-db使用之主键生成器表介绍
date: 2019-09-04 19:26:32
updated: 2024-03-20 16:40:07
categories:
  - [开发框架-Flea,flea-db]
tags:
  - flea-framework
  - flea-db
  - flea-id-generator
  - TableGenerator
  - 分库场景
  - 分表场景
---

[《开发框架-Flea》](/categories/开发框架-Flea/) [《flea-db》](/categories/开发框架-Flea/flea-db/)

![](/images/jpa-logo.png)

# 引言
本篇介绍JPA规范下的主键生成器表，相关主键生成策略可查看 [JPA主键生成策略介绍](/2019/09/03/flea-framework/flea-db/flea-db-jpa-generatedvalue/)。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 主键生成器表
**MySQL** 下的 **flea_id_generator** 表结构：

| 字段 | 名称 |  类型 | 长度 |
|--|--|-- |--|
|id_generator_key  | ID生成器的键【即主键生成策略的键值名称】 |  varchar|  50|
|id_generator_value|ID生成器的值【即主键生成的值】 |  bigint|  20|

相关 SQL 如下：
```sql
CREATE TABLE `flea_id_generator` (
  `id_generator_key` varchar(50) NOT NULL COMMENT 'ID产生器的键【即主键生成策略的键值名称】',
  `id_generator_value` bigint(20) NOT NULL COMMENT 'ID产生器的值【即主键生成的值】',
  UNIQUE KEY `UNIQUE_KEY` (`id_generator_key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `flea_id_generator` VALUES ('pk_order', '0');
```

# 2. 接入介绍
在笔者的 [JPA主键生成策略介绍](/2019/09/03/flea-framework/flea-db/flea-db-jpa-generatedvalue/) 也已经介绍了 **JPA** 规范的 `@TableGenerator` 注解 的相关内容。

## 2.1 通用场景

值得我们注意的是，在上述的主键生成器表中：
- 主键生成器表的字段，对应着 `@TableGenerator` 注解中的 **pkColumnName** 和 **valueColumnName** 两个属性；
- 主键生成器表的表名，对应着 `@TableGenerator`注解中的 **table** 属性。
- `@TableGenerator` 注解中的  **pkColumnValue** 属性，对应主键生成器表中 **id_generator_key** 字段的数据值。

```java
    @Id
    @GeneratedValue(strategy = GenerationType.TABLE, generator = "ORDER_GENERATOR")
    @TableGenerator(
        // 唯一的生成器名称，可以由一个或多个类引用以作为id值的生成器。
        name = "ORDER_GENERATOR",
        // 存储生成的ID值的表的名称
        table = "flea_id_generator",
        // 表中主键列的名称
        pkColumnName = "id_generator_key",
        // 存储最后生成的主键值的列的名称
        valueColumnName = "id_generator_value",
        // ID生成器表中的主键值模板，用于将该生成值集与其他可能存储在表中的值区分开
        pkColumnValue = "pk_order",
        // 从ID生成器表中分配ID号时增加的数量
        allocationSize = 1
    )
    @Column(name = "order_id", unique = true, nullable = false)
    private Long orderId; // 订单编号
```
## 2.2 分表场景
如果存在分表场景，也可以设置分表的主键值模板，如下面的 **pkColumnValue** 属性；这里的 `pk_order_(ORDER_ID)` 中的 `(ORDER_ID)` 需要 分表配置 对应上。

```java
    @Id
    @GeneratedValue(strategy = GenerationType.TABLE, generator = "ORDER_GENERATOR")
    @TableGenerator(
        // 唯一的生成器名称，可以由一个或多个类引用以作为id值的生成器。
        name = "ORDER_GENERATOR",
        // 存储生成的ID值的表的名称
        table = "flea_id_generator",
        // 表中主键列的名称
        pkColumnName = "id_generator_key",
        // 存储最后生成的主键值的列的名称
        valueColumnName = "id_generator_value",
        // ID生成器表中的主键值模板，用于将该生成值集与其他可能存储在表中的值区分开
        pkColumnValue = "pk_order_(ORDER_ID)",
        // 从ID生成器表中分配ID号时增加的数量
        allocationSize = 1
    )
    @Column(name = "order_id", unique = true, nullable = false)
    private Long orderId; // 订单编号
```

flea-db 模块 JPA 的分表配置，参考如下【[flea-table-split.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/flea/db/flea-table-split.xml)】：
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

    </tables>

    <!-- 其他模块分表配置文件引入 -->
    <!--<import resource=""/>-->

</flea-table-split>
```
## 2.3 分库场景
如果存在分库场景，默认主键生成器表在模板库中；当然也可以让主键生成器表存放在每个分库之中，这个时候就需要使用 `@FleaTableGenerator` 注解 ，设置生成器标识 **generatorFlag** 为 **false**，如下所示：

```java
    @Id
    @GeneratedValue(strategy = GenerationType.TABLE, generator = "ORDER_GENERATOR")
    @TableGenerator(
        // 唯一的生成器名称，可以由一个或多个类引用以作为id值的生成器。
        name = "ORDER_GENERATOR",
        // 存储生成的ID值的表的名称
        table = "flea_id_generator",
        // 表中主键列的名称
        pkColumnName = "id_generator_key",
        // 存储最后生成的主键值的列的名称
        valueColumnName = "id_generator_value",
        // ID生成器表中的主键值模板，用于将该生成值集与其他可能存储在表中的值区分开
        pkColumnValue = "pk_order_(ORDER_ID)",
        // 从ID生成器表中分配ID号时增加的数量
        allocationSize = 1
    )
    @FleaTableGenerator(generatorFlag = false)
    @Column(name = "order_id", unique = true, nullable = false)
    private Long orderId; // 订单编号
```

flea-db 模块 JPA 的分库配置，参考如下【[flea-lib-split.xml](https://github.com/Huazie/flea-db-test/blob/main/flea-config/src/main/resources/flea/db/flea-lib-split.xml)】：

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

</flea-table-split>
```

# 3. 接入测试

可以移步 [Huazie](https://github.com/Huazie) 的 GitHub，查看 [《flea-jpa-test》](https://github.com/Huazie/flea-db-test/tree/main/flea-jpa-test) 子项目，该子项目用于 flea-db 模块测试 JPA 相关内容使用。 




