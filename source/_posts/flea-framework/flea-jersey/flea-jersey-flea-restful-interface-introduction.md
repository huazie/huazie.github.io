---
title: flea-jersey使用之Flea RESTful接口介绍
date: 2019-11-22 15:00:20
updated: 2023-06-28 22:22:21
categories:
  - [开发框架-Flea,flea-jersey]
tags:
  - flea-framework
  - flea-jersey
  - Flea RESTful接口
---

![](/images/flea-logo.png)

# 引言
相关文档可查看 [Flea RESTful接口规范.docx](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/Flea%20RESTful%E6%8E%A5%E5%8F%A3%E8%A7%84%E8%8C%83.docx) ，点击 **View raw** 即可下载

<!-- more -->

![](flea-jersey-docx.png)

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 总体概述
**Flea RESTful 接口**，即遵守**REST**式风格的接口，基于**Jersey**开发，遵循**JAX-RS**规范。接入 **Flea RESTful 接口** 的应用提供 **RESTful Web Services**（**REST式**的**Web**服务，它是一种遵守**REST式**风格的**Web**服务）。**REST式**的**Web**服务是一种**ROA**（Resource-Oriented Architecture, 面向资源的架构）的应用。其主要特点是方法信息存在于**HTTP**的方法中（比如**GET**、**PUT**、**POST**、**DELETE**），作用域存在于**URI**中。

# 2. 接口定义
## 2.1 接口协议
基于**HTTP**协议，接口完整报文支持**XML**和**JSON**，接口业务报文使用**JSON**。

## 2.2 交互编码
交互内容编码均采用 **UTF-8** 格式

## 2.3 接口地址
服务端地址/自定义部分 （http://ffs.huazie.com/fleafs）
**自定义部分** 可见如下代码中 注解 `ApplicationPath` 内容
```java
@ApplicationPath("/fleafs/*")
public class FleaFSResourceConfig extends FleaResourceConfig {

    public FleaFSResourceConfig() {
        super();
        // 设置 Jersey 过滤器配置文件 路径
        FleaJerseyFilterConfig.setFilePath("flea/jersey/fleafs-jersey-filter.xml");
    }
}
```

## 2.4 资源定义

以上传资源为例，如下贴出上传资源类，其中注解 `Path` 内容会追加到接口地址中来请求（http://ffs.huazie.com/fleafs/upload）。 

```java
@Path("upload")
public class UploadResource extends FleaJerseyPostResource {

}
```

## 2.5 请求报文
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<JERSEY>
    <REQUEST>
        <PUBLIC>
            <SYSTEM_ACCOUNT_ID></SYSTEM_ACCOUNT_ID>
            <ACCOUNT_ID></ACCOUNT_ID>
            <RESOURCE_CODE></RESOURCE_CODE>
            <SERVICE_CODE></SERVICE_CODE>
        </PUBLIC>
        <BUSINESS>
            <INPUT>业务入参JSON报文</INPUT>
        </BUSINESS>
    </REQUEST>
</JERSEY>
```

**公共报文 PUBLIC :**
  - `SYSTEM_ACCOUNT_ID` : 系统账户编号
  - `ACCOUNT_ID` : 账户编号
  - `RESOURCE_CODE` : 资源编码
  - `SERVICE_CODE` : 服务编码

**业务报文 BUSINESS :**
  - `INPUT` : 业务入参JSON报文


## 2.6 响应报文
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<JERSEY>
    <RESPONSE>
        <PUBLIC>
            <RESULT_CODE></RESULT_CODE>
            <RESULT_MESS></RESULT_MESS>
        </PUBLIC>
        <BUSINESS>
            <OUTPUT>业务出参JSON报文</OUTPUT>
        </BUSINESS>
    </RESPONSE>
</JERSEY>
```

**公共报文 PUBLIC :**
  - `RESULT_CODE` : 返回码
  - `RESULT_MESS` : 返回信息

**业务报文 BUSINESS :**
  - `INPUT` : 业务出参JSON报文

# 3.返回码定义

**Flea RESTful**接口资源服务提供方，如果出现异常，应当抛出异常对应的国际码，同时在**Flea Jersey** 国际码和错误码映射表【`flea_jersey_i18n_error_mapping`】中配置异常国际码和错误返回码的映射关系，其中异常国际码由各资源服务提供方自行配置并使用、错误返回码统一按下面的规范定义。

目前，错误返回码包含如下分类：

以下是将你提供的数据转换为 Markdown 格式后的内容：

## 3.1. 成功
| 一码位 | 二码位 | 三码位 | 描述 |
| ---- | ---- | ---- | ---- |
| 0 | 00 | 000 | 成功，RESULT_MESS返回success |

## 3.2. Flea Jersey接口服务端的过滤器错误返回码
| 一码位 | 二码位 | 三码位 | 描述 |
| ---- | ---- | ---- | ---- |
| 1 | 00 | 000 | 请求报文不能为空 |
| 1 | 00 | 001 | 请求公共报文不能为空 |
| 1 | 00 | 002 | 请求业务报文不能为空 |
| 1 | 00 | 003 | 请求公共报文入参【{0}】不能为空 |
| 1 | 00 | 004 | 未能找到指定资源服务配置数据【service_code = {0} ，resource_code = {1}】 |
| 1 | 00 | 005 | 请检查服务端配置【service_code = {0} ，resource_code = {1}】：【{2} = {3}】非法 |
| 1 | 00 | 006 | 资源【{0}】下的服务【{1}】请求异常：配置的出参【{2}】与服务方法【{3}】出参【{4}】类型不一致 |
| 1 | 00 | 007 | 用户【user_id = {0}】不存在或已失效！ |
| 1 | 00 | 008 | 账户【account_id = {0}】不存在或已失效！ |
| 1 | 00 | 009 | 资源【resource_code = {0}】不存在或已失效！ |
| 1 | 00 | 010 | 账户【account_id = {0}】没有权限调用归属于系统【system_account_id = {1}】的资源【{2}】 |
| 1 | 00 | 011 | 当前资源【{0}】不属于指定系统【system_account_id = {1}】，请确认！ |
| 1 | 00 | 012~999 | 保留的过滤器错误返回码 |
| 9 | 99 | 998 | 返回码未配置 |
| 9 | 99 | 999 | 未知异常 (系统异常等，非自定义的异常) |

## 3.3. Flea Jersey接口服务端的业务异常错误返回码
| 一码位 | 二码位 | 三码位 | 描述 |
| ---- | ---- | ---- | ---- |
| 1 | 01~99 | 000~999 | 业务异常错误返回码 |
| 2~8 | 00~99 | 000~999 | 业务异常错误返回码 |
| 9 | 00~99 | 000~997 | 业务异常错误返回码 | 


# 相关文章
<table>
  <tr>
    <td rowspan="5" align="left" > 
      <a href="../../../../../../categories/开发框架-Flea/flea-jersey/">flea-jersey</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="../../../../../../2019/11/29/flea-framework/flea-jersey/flea-jersey-server/">flea-jersey使用之Flea RESTful接口服务端接入</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="../../../../../../2019/12/15/flea-framework/flea-jersey/flea-jersey-client/">flea-jersey使用之Flea RESTful接口客户端接入</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="../../../../../../2019/12/18/flea-framework/flea-jersey/flea-jersey-file-upload/">flea-jersey使用之文件上传接入</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="../../../../../../2019/12/22/flea-framework/flea-jersey/flea-jersey-file-download/">flea-jersey使用之文件下载接入</a> 
    </td>
  </tr>
</table>
