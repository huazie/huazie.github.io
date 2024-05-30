---
title: flea-auth使用之用户子模块介绍
date: 2020-01-2 11:37:25
updated: 2024-02-28 18:03:09
categories:
  - [开发框架-Flea,flea-auth]
tags:
  - flea-framework
  - flea-auth
  - 用户子模块
---

![](/images/flea-logo.png)

# 引言
本篇主要介绍笔者 **授权模块**【flea-auth】下的用户子模块。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 总览
|  表名                    				   |  中文描述                        |
|:------------------------------------- |:------------------------------   |  
|  flea_account             			|  账户                           |
|  flea_account_attr        			|  账户扩展属性                    |
|  flea_user                				|  用户                            |
|  flea_user_attr           			|  用户扩展属性                     |
|  flea_user_rel       					|  用户关联（角色，角色组）            |
|  flea_user_group          			|  用户组                          			|
|  flea_user_group_rel      		|  用户组关联（角色，角色组）     |
|  flea_real_name_info      		|  实名信息                        |
|  flea_login_log    					|  登录日志【模板表】   |
|  flea_login_log_YYYYMM    	|  登录日志【年月分表】   |

# 2. 详述
## 2.1 账户
授权模块提供的表，这里可细分为 **系统账户** 和 **操作账户** ，如下：

 - **系统账户**，各应用系统在授权模块所注册的账户信息，主要用于各系统之间交互的权限验证；
 - **操作账户**，各应用系统使用者注册的账户信息。

|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|account_id             		| 账户编号                              |
|user_id       				| 用户编号                               |
|account_code         	| 账号                              			|
|account_pwd              | 密码                        				 |
|account_state        		| 账户状态（0：删除，1：正常 ，2：禁用，3：待审核）  |
|create_date            	| 创建日期   								 |
|done_date            	 	| 修改日期   								|
|effective_date            | 生效日期  								 |
|expiry_date               | 失效日期   								|
|remarks                    | 备注信息   								|

## 2.2 账户扩展属性
授权模块提供的账户自定义的属性，用于满足不同应用系统差异化的数据要求；
比如，这里可以自定义账户的类型，用于区分不同账户。
|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|attr_id       					| 属性编号                               |
|account_id             		| 账户编号                              |
|attr_code       				| 属性码                              |
|attr_value         			| 属性值                              			|
|state        					| 属性状态 （0: 删除 1: 正常 ）  |
|create_date            	| 创建日期   								 |
|done_date            	 	| 修改日期   								|
|effective_date            | 生效日期  								 |
|expiry_date               | 失效日期   								|
|remarks                    | 备注信息   								|

## 2.3 用户
授权模块提供的表，与上述账户相对应；
这里可细分 **系统用户** 和 **操作用户** ，如下：

- **系统用户**，各应用系统在授权模块所注册的用户信息；

- **操作用户**，各应用系统使用者注册的用户信息。

|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|user_id       				| 用户编号                               |
|user_name       			| 昵称                               			|
|user_sex       				| 性别（1：男 2：女 3：其他）  |
|user_birthday       		|生日                               				|
|user_address       		| 住址                            			|
|user_email         		| 邮箱                              			|
|user_phone             	 | 手机                      				 |
|group_id              		|用户组编号                      				 |
|user_state        			| 用户状态（0：删除，1：正常 ，2：禁用，3：待审核） |
|create_date            	| 创建日期   										 |
|done_date            	 	| 修改日期   								|
|effective_date            | 生效日期  								 |
|expiry_date               | 失效日期   								|
|remarks                    | 备注信息   								|

## 2.4 用户扩展属性
授权模块提供的用户自定义的属性，用于满足不同应用系统差异化的数据要求。
|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|attr_id       					| 属性编号                              	 |
|user_id             			| 用户编号                             	 |
|attr_code       				| 属性码                             		 |
|attr_value         			| 属性值                              			|
|state        					| 属性状态 （0: 删除 1: 正常 ）  |
|create_date            	| 创建日期   								 |
|done_date            	 	| 修改日期   								|
|effective_date            | 生效日期  								 |
|expiry_date               | 失效日期   								|
|remarks                    | 备注信息   								|

## 2.5 用户关联
授权模块提供的表，目前可关联 **角色**，**角色组**。

- **用户关联角色** ，记录了实际授予给用户的角色信息；

- **用户关联角色组**，记录了实际授予给用户的角色组中角色信息。

|   字段名                |    中文描述                             |
|:------------------------|:---------------------------------------|
|user_rel_id       		| 用户关联编号                		 |
|user_id             		| 用户编号                               |
|rel_id         			    | 关联编号                             	 |
|rel_type                  | 关联类型                        		 |
|rel_state          		| 关联状态 （0: 删除 1: 正常  ）  |
|create_date            | 创建日期   								 |
|done_date             | 修改日期   								|
|remarks                 | 菜单描述   								|
|rel_ext_a                 | 关联扩展字段A   				|
|rel_ext_b                 | 关联扩展字段B   				|
|rel_ext_c                 | 关联扩展字段C  				|
|rel_ext_x                 | 关联扩展字段X   				|
|rel_ext_y                 | 关联扩展字段Y   				|
|rel_ext_z                 | 关联扩展字段Z   				|

## 2.6 用户组
授权模块提供的表，可以理解为同类型的用户集合；
用户拥有的权限，包含 **自身授权** 和 **其归属的用户组授权**。
|   字段名                				|    中文描述                             |
|:-----------------------------------|:---------------------------------------|
|user_group_id             		| 用户组编号                               |
|user_group_name       		| 用户组名称                                |
|user_group_desc         		| 用户组描述                               |
|user_group_state       		| 用户组状态（0: 删除 1: 正常 ）  |
|create_date            			| 创建日期   								 |
|done_date            				 | 修改日期   								|
|remarks                 				| 菜单描述   								|

## 2.7 用户组关联
授权模块提供的表，目前可关联 **角色**，**角色组**。

- **用户组关联角色** ，记录了实际授予给用户组的角色信息；

- **用户组关联角色组**，记录了实际授予给用户组的角色组中角色信息。

|   字段名                				|    中文描述                             |
|:-----------------------------------|:---------------------------------------|
|user_group_rel_id       		| 用户组关联编号                		 |
|user_group_id             		| 用户组编号                               |
|rel_id         					| 关联编号                             	 |
|rel_type                		| 关联类型                        		 |
|rel_state        				| 关联状态 （0: 删除 1: 正常  ）  |
|create_date            		| 创建日期   								 |
|done_date            			 | 修改日期   								|
|remarks                 		| 菜单描述   								|
|rel_ext_a                 | 关联扩展字段A   				|
|rel_ext_b                 | 关联扩展字段B   				|
|rel_ext_c                 | 关联扩展字段C  				|
|rel_ext_x                 | 关联扩展字段X   				|
|rel_ext_y                 | 关联扩展字段Y   				|
|rel_ext_z                 | 关联扩展字段Z   				|

## 2.8 实名信息
授权模块提供的表，用于记录用户实名认证的信息。
|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|real_name_id       		| 实名编号                             	 |
|cert_type             		| 证件类型（1：身份证）            |
|cert_code       			| 证件号码                             		 |
|cert_name         			| 证件名称                              			|
|cert_address         		| 证件地址                              			|
|real_name_state         | 实名信息状态（0: 删除 1: 正常 ）  |
|create_date            	| 创建日期   								 |
|done_date            	 	| 修改日期   								|
|effective_date            | 生效日期  								 |
|expiry_date               | 失效日期   								|
|remarks                    | 备注信息   								|

## 2.9 登录日志
授权模块提供的表，按年月分表，用于记录操作用户登录和登出系统的日志信息。

|   字段名                    |    中文描述                             |
|:----------------------------|:---------------------------------------|
|login_log_id       			| 登录日志编号                          |
|account_id             		| 账户编号           						 |
|system_account_id     | 系统账户编号                          |
|login_ip4         			| ip4地址                             			|
|login_ip6         			| ip6地址                             			|
|login_area         			| 登录地区                             			|
|login_state         			| 登录状态（1：登录中，2：已退出） |
|login_time         			| 登录时间                            			|
|logout_time        		| 退出时间 									|
|create_date            	| 创建日期   								 |
|done_date            	 	| 修改日期   								|
|remarks                    | 备注信息   								|
|ext1                    		| 扩展字段1   								|
|ext2                   			 | 扩展字段2  								|

