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

[《开发框架-Flea》](/categories/开发框架-Flea/)

![](/images/flea-logo.png)

# 引言
相关文档可参考 [Flea RESTful接口规范.docx](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/Flea%20RESTful%E6%8E%A5%E5%8F%A3%E8%A7%84%E8%8C%83.docx) ，点击 **View raw** 即可下载

<!-- more -->

![](flea-jersey-docx.png)

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 总体概述
**Flea RESTful 接口**，即遵守REST式风格的接口，基于Jersey开发，遵循JAX-RS规范。接入 **Flea RESTful 接口** 的应用提供 **RESTful Web Services**（REST式的Web服务，它是一种遵守REST式风格的Web服务）。REST式的Web服务是一种ROA（Resource-Oriented Architecture, 面向资源的架构）的应用。其主要特点是方法信息存在于HTTP的方法中（比如**GET**、**PUT**、**POST**、**DELETE**），作用域存在于URI中。

# 2. 接口定义
## 2.1 接口协议
基于HTTP协议，业务出入参报文支持 **XML** 和 **JSON**。

## 2.2 交互编码
交互内容编码均采用 **UTF-8** 格式

## 2.3 接口地址
服务端地址/自定义部分 （http://ffs.huazie.com/fleafs）
**自定义部分** 可见如下代码中 注解 **ApplicationPath** 内容
```java
/**
 * <p> FleaFS 资源入口 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
@ApplicationPath("/fleafs/*")
public class FleaFSResourceConfig extends FleaResourceConfig {

    /**
     * <p> 无参构造方法 </p>
     *
     * @since 1.0.0
     */
    public FleaFSResourceConfig() {
        super();
    }
}
```

## 2.4 请求报文
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<JERSEY>
    <REQUEST>
        <PUBLIC>
            <SYSTEM_ACCT_ID></SYSTEM_ACCT_ID>
            <SYSTEM_ACCT_PWD></SYSTEM_ACCT_PWD>
            <ACCT_ID></ACCT_ID>
            <RESOURCE_CODE></RESOURCE_CODE>
            <SERVICE_CODE></SERVICE_CODE>
        </PUBLIC>
        <BUSINESS>
            <INPUT>业务入参JSON报文或XML报文</INPUT>
        </BUSINESS>
    </REQUEST>
</JERSEY>
```
## 2.5 响应报文
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<JERSEY>
    <RESPONSE>
        <PUBLIC>
            <RESULT_CODE></RESULT_CODE>
            <RESULT_MESS></RESULT_MESS>
        </PUBLIC>
        <BUSINESS>
            <OUTPUT>业务出参JSON报文或XML报文</OUTPUT>
        </BUSINESS>
    </RESPONSE>
</JERSEY>
```

# 相关文章
<table>
  <tr>
    <td rowspan="5" align="left" > 
      <a href="/categories/开发框架-Flea/flea-jersey/">flea-jersey</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2019/11/29/flea-framework/flea-jersey/flea-jersey-server/">flea-jersey使用之Flea RESTful接口服务端接入</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2019/12/15/flea-framework/flea-jersey/flea-jersey-client/">flea-jersey使用之Flea RESTful接口客户端接入</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2019/12/18/flea-framework/flea-jersey/flea-jersey-file-upload/">flea-jersey使用之文件上传接入</a> 
    </td>
  </tr>
  <tr>
    <td align="left" > 
      <a href="/2019/12/22/flea-framework/flea-jersey/flea-jersey-file-download/">flea-jersey使用之文件下载接入</a> 
    </td>
  </tr>
</table>
