---
title: flea-auth使用之权限子模块介绍
date: 2020-01-14 18:21:28
updated: 2024-02-28 18:13:06
categories:
  - [开发框架-Flea,flea-auth]
tags:
  - flea-framework
  - flea-auth
  - 权限子模块
---

![](/images/flea-logo.png)

# 引言
本篇主要介绍笔者 **授权模块**【flea-auth】下的权限子模块。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 总览
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:-------------------------------------------  |  
|  flea_privilege           				|  权限                            					|
|  flea_privilege_rel      			|  权限关联（菜单， 操作， 元素，资源）|
|  flea_privilege_group     		|  权限组                          				|			
|  flea_privilege_group_rel     	|  权限组关联 (权限)                        |

# 2. 详述
## 2.1 权限
授权模块提供的表，可理解为系统中用户可操作资源的范围和程度。
|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|privilege_id             | 权限编号                                |
|privilege_name       | 权限名称                                |
|privilege_desc         | 权限描述                               |
|group_id                  | 权限组编号                           |
|privilege_state        | 权限状态 （0: 删除 1: 正常 ）   |
|create_date            | 创建日期   								 |
|done_date             | 修改日期   								|
|remarks                 | 菜单描述   								|

## 2.2 权限关联
授权模块提供的表，目前定义四种功能的关联，分别为 **菜单**、**操作** 、 **元素** 和 **资源**；
关联类型 【**relat_type**】可以自行定义。

|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|privilege_rel_id       | 权限关联编号                         |
|privilege_id             | 权限编号                                |
|rel_id         			| 关联编号                             	 |
|rel_type                | 关联类型                        		 |
|rel_state        		| 关联状态 （0: 删除 1: 正常 ）   |
|create_date            | 创建日期   								 |
|done_date             | 修改日期   								|
|remarks                 | 菜单描述   								|
|rel_ext_a                 | 关联扩展字段A   				|
|rel_ext_b                 | 关联扩展字段B   				|
|rel_ext_c                 | 关联扩展字段C  				|
|rel_ext_x                 | 关联扩展字段X   				|
|rel_ext_y                 | 关联扩展字段Y   				|
|rel_ext_z                 | 关联扩展字段Z   				|

## 2.3 权限组
授权模块提供的表，可以理解为同类型的权限的集合；
权限和权限组之间是多对一的关系，权限表中 **group_id** 记录权限组编号，默认值为-1；
权限组关联的功能，即为其下所有权限关联的功能，不单独为权限组关联功能。

|   字段名                				|    中文描述                             |
|:-----------------------------------|:---------------------------------------|
|privilege_group_id             | 权限组编号                                |
|privilege_group_name       | 权限组名称                                |
|privilege_group_desc         | 权限组描述                               |
|privilege_group_state        | 权限组状态 （0: 删除 1: 正常 ）   |
|create_date            			| 创建日期   								 |
|done_date             				| 修改日期   								|
|remarks                 				| 菜单描述   								|

## 2.4 权限组关联
授权模块提供的表，目前可关联 **权限**。

|   字段名               				 |    中文描述                             |
|:-----------------------------------|:---------------------------------------|
|privilege_group_rel_id       | 权限组关联编号                       |
|privilege_group_id             | 权限组编号                                |
|rel_id         					| 关联编号                             	 |
|rel_type                		| 关联类型                        		 |
|rel_state        				| 关联状态 （0: 删除 1: 正常 ）  |
|create_date            		| 创建日期   								 |
|done_date             			| 修改日期   								|
|remarks                 			| 菜单描述   								|
|rel_ext_a                 	| 关联扩展字段A   				|
|rel_ext_b                 	| 关联扩展字段B   				|
|rel_ext_c                 	| 关联扩展字段C  				|
|rel_ext_x                 	| 关联扩展字段X   				|
|rel_ext_y                 	| 关联扩展字段Y   				|
|rel_ext_z                 	| 关联扩展字段Z   				|

**权限组关联权限**，引入了如下概念：
 - **组内互斥**，权限组中的权限存在可操作资源的范围和程度上的相互制约，在进行角色授权时，只能选择组内的一个权限授予角色。

