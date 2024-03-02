---
title: flea-auth使用之角色权限设计初识
date: 2019-11-18 08:57:19
updated: 2023-08-17 16:04:59
categories:
  - [开发框架-Flea,flea-auth]
tags:
  - flea-framework
  - flea-auth
  - RBAC
  - 角色权限设计
  - 用户子模块
  - 角色子模块
  - 权限子模块
  - 功能子模块
---

[《开发框架-Flea》](/categories/开发框架-Flea/)

![](/images/flea-logo.png)

# 引言
本篇将要介绍笔者 `Flea` 框架下的授权模块（[flea-auth](https://github.com/Huazie/flea-framework/tree/dev/flea-auth) ），该模块以 **RBAC** 为基础拓展而来。

# 1. 什么是 RBAC？
**RBAC**（Role-Based Access Control），基于角色的访问控制。其基本思想是，在用户与权限之间新增一个角色的概念。每一种角色关联单个或多个权限。只需要给用户分配适当的角色，该用户就拥有其关联角色下的所有权限。

相较于直接给用户分配权限，带来用户和权限数据的繁杂与冗余；**RBAC** 只需要分配相应的角色给用户即可，而且角色的权限变更相比用户的权限变更所带来的系统开销要小得多。

**RBAC** 模型可以分为：**RBAC0**、**RBAC1**、**RBAC2**、**RBAC3** 四种。其中 **RBAC0** 是最简单的角色权限模型（用户，角色和权限） ，**RBAC1**、**RBAC2**、**RBAC3** 都是以 **RBAC0** 为基础拓展而来。更多 **RBAC** 的了解，可见如下拓展链接：

**拓展：**
* [RBAC模型：基于用户-角色-权限控制的一些思考](http://www.sohu.com/a/245751423_114819) 
* [角色权限设计的100种解法](http://www.woshipm.com/pd/1214616.html)

# 2. flea-auth 都有哪些内容？
**flea-auth** 包含四个子模块：

 - **用户子模块**
 - **角色子模块** 
 - **权限子模块**
 - **功能子模块**

## 2.1 相关表总览
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:------------------------------   |  
|  flea_account             			|  账户                           |
|  flea_account_attr        			|  账户扩展属性                    |
|  flea_real_name_info      		|  实名信息                        |
|  flea_login_log    					|  登录日志【模板表】   |
|  flea_login_log_YYYYMM    	|  登录日志【年月分表】   |
|  flea_user                				|  用户                            |
|  flea_user_attr           			|  用户扩展属性                     |
|  flea_user_rel       					|  用户关联（角色，角色组）            |
|  flea_user_group          			|  用户组                          			|
|  flea_user_group_rel      		|  用户组关联（角色，角色组）     |
|  flea_role                				|  角色                            |
|  flea_role_rel            				|  角色关联（角色， 权限， 权限组）   |
|  flea_role_group          			|  角色组（不参与授权）              |
|  flea_role_group_rel      		|  角色组关联（角色）   |
|  flea_privilege           				|  权限                            					|
|  flea_privilege_rel      			|  权限关联（菜单， 操作， 元素，资源）|
|  flea_privilege_group     		|  权限组                          				|			
|  flea_privilege_group_rel     	|  权限组关联 (权限)                        |
|  flea_menu                			|  菜单                                             |
|  flea_operation           			|  操作                                            |
|  flea_element            			|  元素                                           |
|  flea_resource                       |  资源                                             |
|  flea_function_attr        			|  功能扩展属性【模板表】             |
|  flea_function_attr_menu       |  菜单扩展属性【功能类型分表】 |
|  flea_function_attr_operation |  操作扩展属性【功能类型分表】 |
|  flea_function_attr_element   |  元素扩展属性【功能类型分表】 |
|  flea_function_attr_resource  |  资源扩展属性【功能类型分表】 |

## 2.2 相关表SQL
上述相关表SQL，笔者此处就不再一一列出，可至GitHub查看 [fleaauth.sql](https://github.com/Huazie/flea-framework/tree/dev/flea-auth/fleaauth.sql) , 脚本内容基于MySQL，可作为参考。

## 2.3 用户子模块介绍
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:------------------------------   |  
|  flea_account             			|  账户                           |
|  flea_account_attr        			|  账户扩展属性                    |
|  flea_real_name_info      		|  实名信息                        |
|  flea_login_log    					|  登录日志【模板表】   |
|  flea_login_log_YYYYMM    	|  登录日志【年月分表】   |
|  flea_user                				|  用户                            |
|  flea_user_attr           			|  用户扩展属性                     |
|  flea_user_rel       					|  用户关联（角色，角色组）            |
|  flea_user_group          			|  用户组                          			|
|  flea_user_group_rel      		|  用户组关联（角色，角色组）     |
###  2.3.1 账户
这里可细分为 **系统账户** 和 **操作账户** ：

**系统账户**，各应用系统在授权模块所注册的账户信息，主要用于各系统之间交互的权限验证；
**操作账户**，各应用系统使用者注册的账户信息。

### 2.3.2 账户扩展属性
授权模块提供的账户自定义的属性，用于满足不同应用系统差异化的数据要求；
比如，这里可以自定义账户的类型，用于区分不同账户。

### 2.3.3 用户
与账户相对应，这里可细分 **系统用户** 和 **操作用户** ：

**系统用户**，各应用系统在授权模块所注册的用户信息；
**操作用户**，各应用系统使用者注册的用户信息。

### 2.3.4 用户扩展属性
授权模块提供的用户自定义的属性，用于满足不同应用系统差异化的数据要求。

### 2.3.5 实名信息
授权模块提供的表，用于记录用户实名认证的信息。

### 2.3.6 登录日志
授权模块提供的表，用于记录操作用户登录和登出系统的日志信息。

### 2.3.7 用户关联
授权模块提供的表，目前可关联 **角色**，**角色组**。

**用户关联角色** ，记录了实际授予给用户的角色信息；

**用户关联角色组**，记录了实际授予给用户的角色组中角色信息。

### 2.3.8 用户组
授权模块提供的表，可以理解为同类型的用户集合；用户拥有的权限，包含自身授权和其归属的用户组授权。

### 2.3.9 用户组关联
授权模块提供的表，目前可关联 **角色**，**角色组**。

**用户组关联角色** ，记录了实际授予给用户组的角色信息；

**用户组关联角色组**，记录了实际授予给用户组的角色组中角色信息。

## 2.4 角色子模块介绍
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:------------------------------   |  
|  flea_role                				|  角色                            |
|  flea_role_rel            				|  角色关联（角色， 权限， 权限组）   |
|  flea_role_group          			|  角色组（不参与授权）              |
|  flea_role_group_rel      		|  角色组关联（角色）   |

### 2.4.1 角色
授权模块提供的表，可理解为具备一定权限的一类用户。

### 2.4.2 角色关联
授权模块提供的表，目前可关联 **角色**、**权限**、**权限组**。

**角色关联角色**，引入了如下概念：
 - **角色继承**，关联角色（子角色）可继承被关联角色（父角色）的所有权限；
 - **角色互斥**，关联角色和被关联角色存在权限上的相互制约，在进行用户授权时，两者不能同时授予同一用户；
 - **角色基数约束**，系统中可以拥有这个角色的用户数目限制；

**角色关联权限**，记录了实际给角色绑定的权限信息。

**角色关联权限组**，记录了实际给角色绑定的权限组中的权限信息。
### 2.4.3 角色组
授权模块提供的表，可理解为具备一定权限的一类用户的集合；它本身不参与授权，其下所拥有的权限由其角色成员决定。

### 2.4.4 角色组关联
授权模块提供的表，目前可关联 **角色**。

**角色组关联角色**，引入了如下概念：
 - **组内互斥**，角色组中的角色存在权限上的相互制约，在进行用户授权时，只能选择组内的一个角色授予用户或用户组；

## 2.5 权限子模块介绍
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:-------------------------------------------  |  
|  flea_privilege           				|  权限                            					|
|  flea_privilege_rel      			|  权限关联（菜单， 操作， 元素，资源）|
|  flea_privilege_group     		|  权限组                          				|			
|  flea_privilege_group_rel     	|  权限组关联 (权限)                        |

### 2.5.1 权限
授权模块提供的表，可理解为系统中用户可操作资源的范围和程度。

### 2.5.2 权限关联
授权模块提供的表，目前定义三种功能的关联，分别为 **菜单**、**操作** 、 **元素** 和 **资源**。

### 2.5.3 权限组
授权模块提供的表，可以理解为同类型的权限的集合；权限和权限组之间是多对一的关系，权限表中 **group_id** 记录权限组编号，默认值为-1；权限组关联的功能，即为其下所有权限关联的功能，不单独为权限组关联功能。

### 2.5.4 权限组关联
授权模块提供的表，目前可关联 **权限**。

**权限组关联权限**，引入了如下概念：
 - **组内互斥**，权限组中的权限存在可操作资源的范围和程度上的相互制约，在进行角色授权时，只能选择组内的一个权限授予角色。

## 2.6 功能子模块介绍
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:------------------------------   |  
|  flea_menu                			|  菜单                            |
|  flea_operation           			|  操作                            |
|  flea_element            			|  元素                             |
|  flea_resource                       |  资源                                             |
|  flea_function_attr        			|  功能扩展属性【模板表】             |
|  flea_function_attr_menu       |  菜单扩展属性【功能类型分表】 |
|  flea_function_attr_operation |  操作扩展属性【功能类型分表】 |
|  flea_function_attr_element   |  元素扩展属性【功能类型分表】 |
|  flea_function_attr_resource  |  资源扩展属性【功能类型分表】 |


目前功能子模块包含 **菜单**、**操作** 、**元素** 和 **资源**，这些功能都和权限相关联【可参考 授权模块下的 **权限关联表 flea_privilege_rel** 】，如下：
### 2.6.1 菜单
授权模块提供的表，可解释为一系列业务逻辑的总和，为完成某种特定功能，而定义的一类功能单元。
### 2.6.2 操作
授权模块提供的表，可理解为业务逻辑上较为单一的功能单元，如角色新增，权限新增等。
### 2.6.3 元素
授权模块提供的表，目前有页面元素定义，如页面按钮等。

### 2.6.4 资源
授权模块提供的表，目前有Flea Jersey接口资源定义。
### 2.6.4 功能扩展属性
授权模块提供的表，根据功能类型进行分表，为上述功能单元配置扩展属性。

如果以后需要新增某种功能的权限控制，只需要如下操作：
- 在功能子模块中新增相关功能表；
- 在权限子模块的 **权限关联表** 中新定义一种 **关联类型** 【**rel_type**】，并绑定上相关授权数据。
