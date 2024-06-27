---
title: flea-jersey使用之Flea RESTful接口客户端接入
date: 2019-12-15 10:44:35
updated: 2023-06-28 14:28:31
categories:
  - [开发框架-Flea,flea-jersey]
tags:
  - flea-framework
  - flea-jersey
  - Flea RESTful接口
  - 客户端接入
---

![](/images/flea-logo.png)

# 引言
本篇介绍 **flea-jersey** 模块下的 **flea-jersey-client** 子模块，该模块提供对 **flea-jersey-server** 子模块封装的 `POST`、`PUT`、`DELETE` 和 `GET`资源的调用。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 客户端依赖
项目地址可至GitHub 查看 [flea-jersey-client](https://github.com/Huazie/flea-framework/tree/dev/flea-jersey/flea-jersey-client)
```xml
  <!-- FLEA JERSEY CLIENT-->
    <dependency>
        <groupId>com.huazie.fleaframework</groupId>
        <artifactId>flea-jersey-client</artifactId>
        <version>2.0.0</version>
    </dependency>
```

# 2. 客户端接入步骤
1. 客户端确定待调用的资源服务【参见 [Flea RESTful接口服务端接入](../../../../../../2019/11/29/flea-framework/flea-jersey/flea-jersey-server/)】，并配置资源客户端表；
2. 客户端定义业务入参和业务出参 **POJO** 类；
3. 客户端使用 **FleaJerseyClient** 调用资源服务。

# 3. 具体接入讲解
## 3.1 资源客户端配置
资源客户端【**flea_jersey_res_client**】, 下载鉴权资源服务的资源客户端配置如下：

![](flea_jersey_res_client.png)

|   字段名                |    中文描述        |
|:------------------------|:----------------------|
|client_code          | 客户端编码          |
|resource_url        | 资源地址             |
|resource_code     | 资源编码            |
|service_code        | 服务编码             |
|service_interfaces |资源服务接口类 |
|request_mode      | 请求方式       |
|media_type          | 媒体类型            |
|client_input        |  客户端业务入参  |
|client_output      |  客户端业务出参  |

## 3.2 客户端业务输入和输出参数定义
这里定义的客户端业务入参【**com.huazie.ffs.pojo.upload.input.InputFileUploadInfo**】和 出参【**com.huazie.ffs.pojo.upload.output.OutputFileUploadInfo**】使用的是服务端定义的资源服务入参和出参；
当然这里也可以不一样，原则上只需要客户端业务入参和资源服务入参，客户端业务出参和资源服务出参两边对象转 **JSON** 或 **XML** 的数据内容一致即可。

## 3.3 FleaJerseyClient使用
经过1和2的步骤，客户端接入已经完成一半，下面就可以调用资源服务，可参考如下：
```java
  @Test
    public void testDownloadAuth() {
        try {
            String clientCode = "FLEA_CLIENT_DOWNLOAD_AUTH";

            InputDownloadAuthInfo downloadAuthInfo = new InputDownloadAuthInfo();
            downloadAuthInfo.setFileId("123123123123123123123");

            FleaJerseyClient client = applicationContext.getBean(FleaJerseyClient.class);

            Response<OutputDownloadAuthInfo> response = client.invoke(clientCode, downloadAuthInfo, OutputDownloadAuthInfo.class);

            LOGGER.debug("result = {}", response);
        } catch (Exception e) {
            LOGGER.error("Exception = ", e);
        }
    }
```

# 总结

至此，Flea RESTful接口客户端接入已经完成。上述自测类，可至GitHub查看 [JerseyTest.java](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/test/java/com/huazie/fleaframework/jersey/client/resource/JerseyTest.java)
