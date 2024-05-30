---
title: flea-jersey使用之Flea RESTful接口服务端接入
date: 2019-11-29 14:20:13
updated: 2023-06-28 14:25:01
categories:
  - [开发框架-Flea,flea-jersey]
tags:
  - flea-framework
  - flea-jersey
  - Flea RESTful接口
  - 服务端接入
---

![](/images/flea-logo.png)

# 引言
本篇介绍 **flea-jersey**模块下的**flea-jersey-server** 子模块，该模块封装了通用的POST、PUT、DELETE 和 GET资源。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 服务端依赖
项目地址可至GitHub 查看 [flea-jersey-server](https://github.com/Huazie/flea-framework/tree/dev/flea-jersey/flea-jersey-server)
```xml
  <!-- FLEA JERSEY SERVER-->
    <dependency>
        <groupId>com.huazie.fleaframework</groupId>
        <artifactId>flea-jersey-server</artifactId>
        <version>2.0.0</version>
    </dependency>
```

# 2. 服务端接入步骤
 1. 服务端自定义资源入口类，继承 [FleaResourceConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-server/src/main/java/com/huazie/fleaframework/jersey/server/core/FleaResourceConfig.java)；
 2. 服务端自定义资源，并配置到资源表中；
 3. 服务端自定义资源服务，并配置到资源服务表中；
 4. 服务端完成资源服务的业务逻辑开发，配置国际码和错误码映射关系。

# 3. 具体接入讲解
## 3.1 资源入口类定义
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
        // 这里加入自定义的配置信息
    }
}
```

## 3.2 资源定义
REST服务的核心是对外公布的资源API。如下 [DownloadResource](https://github.com/Huazie/FleaFS/blob/main/fleafs-business/src/main/java/com/huazie/ffs/module/download/web/DownloadResource.java) 资源类由注解Path修饰，其资源路径为 download。
```java
/**
 * <p> 下载资源类 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
@Path("download")
public class DownloadResource extends Resource implements JerseyGetResource, JerseyPostResource {

    @Override
    public FleaJerseyResponse doGetResource(String requestData) {
        return doResource(requestData);
    }

    @Override
    public FleaJerseyResponse doPostResource(FleaJerseyRequest request) {
        return doResource(request);
    }
}
```
有了资源类，下面就需要配置资源；资源配置在 **flea_jersey_resource** 表中。新增如下配置：

![](flea_jersey_resource.png)

应用服务启动后，**FleaFSResourceConfig** 会扫描所有定义的资源包，即将如上**resource_packages** 字段定义的包都扫描一遍，这样这些包内所有资源类所提供的资源路径将被映射到内存中。详细内容可参考  [FleaResourceConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-server/src/main/java/com/huazie/fleaframework/jersey/server/core/FleaResourceConfig.java)。

## 3.3 资源服务定义
### 3.3.1 资源服务接口

```java
public interface IFleaDownloadSV {

    /**
     * <p> 下载授权 </p>
     *
     * @param input 下载授权业务入参
     * @return 下载授权业务出参
     * @throws Exception
     * @since 1.0.0
     */
    OutputDownloadAuthInfo downloadAuth(InputDownloadAuthInfo input) throws Exception;
}
```
### 3.3.2 资源服务实现类
```java
/**
 * <p> Flea下载服务实现类 </p>
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
@Service
public class FleaDownloadSVImpl implements IFleaDownloadSV {

    @Override
    public OutputDownloadAuthInfo downloadAuth(InputDownloadAuthInfo input) throws Exception {
        return null;
    }
}
```
### 3.3.3 资源服务配置
资源服务配置在 **flea_jersey_res_service** 表中。

![](flea_jersey_res_service.png)

|   字段名                |    中文描述        |
|:------------------------|:----------------------|
|service_code        | 服务编码             |
|resource_code     | 资源编码            |
|service_interfaces |资源服务接口类 |
|service_method    | 资源服务方法   |
|service_input        |  资源服务入参  |
|service_output      |  资源服务出参  |

### 3.3.4 资源服务调用
上述资源服务调用逻辑, 可参考 服务调用过滤器 [InvokeServiceFilter](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-server/src/main/java/com/huazie/fleaframework/jersey/server/filter/impl/InvokeServiceFilter.java)。
## 3.4 资源服务业务逻辑开发
```java
  private static final Logger LOGGER = LoggerFactory.getLogger(FleaDownloadSVImpl.class);

    @Override
    public OutputDownloadAuthInfo downloadAuth(InputDownloadAuthInfo input) throws Exception {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaDownloadSVImpl#downloadAuth(InputDownloadAuthInfo) Start");
        }

        String fileId = input.getFileId();
        if (StringUtils.isBlank(fileId)) {
            // 入参【{0}】不能为空
            throw new ServiceException("ERROR-SERVICE0000000001", "fileId");
        }

        OutputDownloadAuthInfo output = new OutputDownloadAuthInfo();
        // 演示直接塞了一个随机数
        output.setToken(RandomCode.toUUID());

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("FleaDownloadSVImpl#downloadAuth(InputDownloadAuthInfo) End");
        }

        return output;
    }
```

## 3.5 国际码和错误码映射

代码中出现 **ERROR-SERVICE0000000001** 的异常，需要配置如下国际码和错误码的映射关系： （国际码和错误码映射配置表 **flea_jersey_i18n_error_mapping**）

![](flea_jersey_i18n_error_mapping.png)

# 总结

至此，Flea RESTful接口服务端接入已经完成。