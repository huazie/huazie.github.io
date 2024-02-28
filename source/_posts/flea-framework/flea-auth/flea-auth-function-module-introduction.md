---
title: flea-auth使用之功能子模块介绍
date: 2020-01-20 10:25:42
updated: 2024-02-28 18:07:43
categories:
  - [开发框架-Flea,flea-auth]
tags:
  - flea-framework
  - flea-auth
  - 功能子模块
---

[《开发框架-Flea》](/categories/开发框架-Flea/)

![](/images/flea-logo.png)

# 引言
本篇主要介绍笔者 **授权模块**【flea-auth】下的功能子模块。
# 1. 总览
目前功能子模块包含 **菜单**、**操作**、 **元素** 和 **资源** 四类功能单元；
这些功能都和权限相关联【可参考 [权限子模块](https://blog.csdn.net/u012855229/article/details/103719630) 的 **权限关联表 flea_privilege_rel** 】。
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:------------------------------   |  
|  flea_menu                			|  菜单                            |
|  flea_operation           			|  操作                            |
|  flea_element            			|  元素                           |
|  flea_resource                       |  资源                                             |
|  flea_function_attr        			|  功能扩展属性【模板表】             |
|  flea_function_attr_menu       |  菜单扩展属性【功能类型分表】 |
|  flea_function_attr_operation |  操作扩展属性【功能类型分表】 |
|  flea_function_attr_element   |  元素扩展属性【功能类型分表】 |
|  flea_function_attr_resource  |  资源扩展属性【功能类型分表】 |

# 2. 详述
## 2.1 菜单
授权模块提供的表，可解释为一系列业务逻辑的总和，为完成某种特定功能，而定义的一类功能单元。
|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|menu_id                | 菜单编号                                |
|menu_code           | 菜单编码                                 |
|menu_name          | 菜单名称                                       |
|menu_icon            | 菜单FontAwesome小图标               |
|menu_sort            | 菜单展示顺序(同一个父菜单下)     |
|menu_view            | 菜单对应页面（非叶子菜单的可以为空）     |
|menu_level            | 菜单层级（1：一级 2；二级 3：三级 4：四级）    |
|menu_state            | 菜单状态（0:下线，1: 在用 ）    |
|parent_id               | 父菜单编号    |
|create_date            | 创建日期    |
|done_date             | 修改日期   |
|effective_date        | 生效日期   |
|expiry_date            | 失效日期   |
|remarks                 | 菜单描述   |
## 2.2 操作
授权模块提供的表，可理解为业务逻辑上较为单一的功能单元，如角色新增，权限新增等操作。
|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|operation_id          | 操作编号                                  |
|operation_code      | 操作编码                                    |
|operation_name      | 操作名称                                    |
|operation_desc       | 操作描述             				       |
|operation_state      | 操作状态（0: 删除 1: 正常 ）    |
|create_date            | 创建日期                                  |
|done_date              | 修改日期                                   |
|remarks                  | 备注信息                                    |

## 2.3 元素
授权模块提供的表，目前有页面元素定义，如页面按钮等。
|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|element_id            | 元素编号                              |
|element_code       | 元素编码                                |
|element_name       | 元素名称              				|
|element_desc         | 元素描述    						 |
|element_type         | 元素类型     						|
|element_content    | 元素内容   							 |
|element_state        | 元素状态（0: 删除 1: 正常 ）   |
|create_date            | 创建日期   							 |
|done_date             | 修改日期   							|
|remarks                 | 菜单描述   							|

## 2.4 资源
授权模块提供的表，用于各种资源的授权，目前有Flea Jersey接口资源定义。

|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|resource_id          | 资源编号                                  |
|resource_code      | 资源编码                                    |
|resource_name      | 资源名称                                    |
|resource_desc       | 资源描述             				       |
|resource_state      | 资源状态（0: 删除 1: 正常 ）    |
|create_date            | 创建日期                                  |
|done_date              | 修改日期                                   |
|remarks                  | 备注信息                                    |

## 2.4 功能扩展属性
授权模块提供的表，根据功能类型进行分表，为上述功能单元配置扩展属性。
|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|attr_id       					| 属性编号                              	 |
|function_id             		| 功能编号  |
|function_type             	| 功能类型(菜单、操作、元素)   |
|attr_code       				| 属性码                             		 |
|attr_value         			| 属性值                              			|
|state        					| 属性状态 （0: 删除 1: 正常 ）  |
|create_date            	| 创建日期   								 |
|done_date            	 	| 修改日期   								|
|effective_date            | 生效日期  								 |
|expiry_date               | 失效日期   								|
|remarks                    | 备注信息   								|


如果以后需要新增某种功能的权限控制，只需要如下操作：
- 在功能子模块中新增相关功能表；
- 在权限子模块的 **权限关联表** 中新定义一种 **关联类型** 【**rel_type**】，并绑定上相关授权数据。

