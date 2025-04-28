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
本篇介绍 **flea-jersey**模块下的**flea-jersey-server** 子模块，该模块封装了通用的`POST`、`PUT`、`DELETE` 和 `GET`资源。

<!-- more -->

[![](/images/flea-framework.png)](https://github.com/Huazie/flea-framework)

# 1. 服务端依赖
项目内容可至GitHub 查看 [flea-jersey-server](https://github.com/Huazie/flea-framework/tree/dev/flea-jersey/flea-jersey-server)

相关依赖如下：
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

[FleaResourceConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-server/src/main/java/com/huazie/fleaframework/jersey/server/core/FleaResourceConfig.java) ，即**Flea** 资源配置类，作为 **Jersey** 应用的资源入口，用于配置 **Web** 应用程序。

该类初始化时，从 **Flea Jersey** 资源表中，获取定义的所有资源包名； 并将所有资源包都添加到扫描组件中，以待被递归扫描（包括所有嵌套包）。

```java
public abstract class FleaResourceConfig extends ResourceConfig {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(FleaResourceConfig.class);

    public FleaResourceConfig() {
        init();
    }

    private void init() {
        // 获取所有的资源包名
        List<String> resourcePackages = null;
        try {
            // 获取Web应用上下文对象
            WebApplicationContext webApplicationContext = ContextLoader.getCurrentWebApplicationContext();
            FleaConfigDataSpringBean springBean = webApplicationContext.getBean(FleaConfigDataSpringBean.class);
            resourcePackages = springBean.getResourcePackages();
        } catch (Exception e) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.error1(new Object() {}, "Exception occurs when getting resource packages : \n", e);
            }
        }

        if (CollectionUtils.isNotEmpty(resourcePackages)) {
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug1(new Object() {}, "scan packages : {}", resourcePackages);
            }
            packages(resourcePackages.toArray(new String[0]));
        }
        // 服务端注册MultiPartFeature组件，用于支持 multipart/form-data 媒体类型
        register(MultiPartFeature.class);
    }
}
```

另外每个接入 **Flea Jersey** 的应用，都需创建 **Flea** 资源配置的子类，作为其发布的资源的入口；并在该类上标记注解 `ApplicationPath` ， 其值为该应用对外发布的资源的相对访问路径。

对于 **FleaFS** 应用而言，就是如下的 **FleaFS** 资源入口类：

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

这里我们在入口类的无参构造方法中，可以看到设置 Jersey 过滤器配置文件路径的代码，下面来介绍下这个配置文件：

## 3.2 Jersey 过滤器配置文件

**FleaFS Jersey 过滤器配置文件**，该文件中可定义**FleaFS**应用下的接口处理的**前置、业务服务、后置和异常过滤器链**，并导入了公共的 flea jersey 接口过滤器公共配置文件。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<jersey>
    <filter-chain>
        <!-- 前置过滤器链 -->
        <before>
            <filter clazz="com.huazie.ffs.jersey.filter.FleaFSAuthCheckFilter" order="3" desc="FleaFS业务授权校验过滤器"/>
        </before>
    </filter-chain>
    <!-- flea jersey 接口过滤器公共配置文件引入 -->
    <import resource="flea/jersey/flea-jersey-filter.xml"/>
</jersey>
```

**flea jersey 接口过滤器公共配置文件**，该文件中包含了Flea Jersey默认的**前置、业务服务、后置和异常过滤器链**，I18N国际化映射配置等信息。

```xml
<?xml version="1.0" encoding="UTF-8"?>

<jersey>

    <filter-chain>

        <!-- 前置过滤器链 -->
        <before>
            <filter clazz="com.huazie.fleaframework.jersey.server.filter.impl.DataPreCheckFilter" order="1" desc="数据预校验过滤器"/>
            <filter clazz="com.huazie.fleaframework.jersey.server.filter.impl.AuthCheckFilter" order="2" desc="授权校验过滤器"/>
        </before>

        <!-- 业务服务过滤器链 -->
        <service>
            <filter clazz="com.huazie.fleaframework.jersey.server.filter.impl.InvokeServiceFilter" order="1" desc="服务调用过滤器"/>
        </service>

        <!-- 后置过滤器链 -->
        <after>
            <filter clazz="com.huazie.fleaframework.jersey.server.filter.impl.JerseyLoggerFilter" order="1" desc="Jersey日志记录过滤器"/>
        </after>

        <!-- 异常过滤器链 -->
        <error>
            <filter clazz="com.huazie.fleaframework.jersey.server.filter.impl.ErrorFilter" order="1" desc="异常过滤器"/>
        </error>

    </filter-chain>

    <filter-i18n-error>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000003" errorCode="100000">请求报文不能为空</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000004" errorCode="100001">请求公共报文不能为空</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000005" errorCode="100002">请求业务报文不能为空</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000006" errorCode="100003">请求公共报文入参【{0}】不能为空</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000008" errorCode="100004">未能找到指定资源服务配置数据【service_code = {0} ，resource_code = {1}】</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000009" errorCode="100005">请检查服务端配置【service_code = {0} ，resource_code = {1}】：【{2} = {3}】非法</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-FILTER0000000010" errorCode="100006">资源【{0}】下的服务【{1}】请求异常：配置的出参【{2}】与服务方法【{3}】出参【{4}】类型不一致</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-AUTH-COMMON0000000007" errorCode="100007">用户【user_id = {0}】不存在或已失效！</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-AUTH-COMMON0000000008" errorCode="100008">账户【account_id = {0}】不存在或已失效！</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-AUTH-COMMON0000000021" errorCode="100009">资源【resource_code = {0}】不存在或已失效！</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-AUTH-COMMON0000000022" errorCode="100010">账户【account_id = {0}】没有权限调用归属于系统【system_account_id = {1}】的资源【{2}】</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-AUTH-COMMON0000000023" errorCode="100011">当前资源【{0}】不属于指定系统【system_account_id = {1}】，请确认！</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-COMMON0000000000" errorCode="100012">【{0}】获取失败，请检查！</i18n-error-mapping>
        <i18n-error-mapping i18nCode="ERROR-JERSEY-COMMON0000000001" errorCode="100013">【{0}】不能为空，请检查！</i18n-error-mapping>
    </filter-i18n-error>

</jersey>
```

## 3.3 资源定义
REST服务的核心是对外公布的资源API。如下 [DownloadResource](https://github.com/Huazie/FleaFS/blob/main/fleafs-business/src/main/java/com/huazie/ffs/module/download/web/DownloadResource.java) 下载资源类由注解 `Path` 修饰，其资源路径为 `download`。

```java
@Path("download")
public class DownloadResource extends FleaJerseyFGetResource implements JerseyPostResource {

    /**
     * @see JerseyPostResource#doPostResource(FleaJerseyRequest)
     */
    @Override
    public FleaJerseyResponse doPostResource(FleaJerseyRequest request) {
        return doResource(request);
    }

}
```

下载资源类继承 `FleaJerseyFGetResource`，即 **Flea Jersey文件GET资源**【它只包含了**文件GET资源API**】。同时实现 `JerseyPostResource` 接口，即**Jersey POST 资源接口**【它只包含 **POST资源API**】。

这里主要考虑下载资源仅支持文件的下载【通过 `FleaJerseyFGetResource` 】和下载的鉴权【通过 `JerseyPostResource`】。

有了**下载资源类**，下面就需要我们配置该资源；

资源配置在 **flea_jersey_resource** 表中，需要新增如下配置：

![](flea_jersey_resource.png)

应用服务启动后，**FleaFSResourceConfig** 会扫描所有定义的资源包，即将如上**resource_packages** 字段定义的包都扫描一遍，这样这些包内所有资源类所提供的资源路径将被映射到内存中。详细内容可参考  [FleaResourceConfig](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-server/src/main/java/com/huazie/fleaframework/jersey/server/core/FleaResourceConfig.java)。

## 3.3 资源服务定义
### 3.3.1 资源服务接口

`IFleaDownloadSV` 即下载资源的资源服务接口，如下贴出的是下载授权的资源服务方法。

```java
public interface IFleaDownloadSV {

    /**
     * <p> 下载授权 </p>
     *
     * @param input 下载授权业务入参
     * @return 下载授权业务出参
     * @since 1.0.0
     */
     OutputDownloadAuthInfo downloadAuth(InputDownloadAuthInfo input) throws CommonException;
}
```

### 3.3.2 资源服务实现类

`FleaDownloadSVImpl` 即下载资源的资源服务实现类

```java
@Service
public class FleaDownloadSVImpl implements IFleaDownloadSV {

    @Override
    public OutputDownloadAuthInfo downloadAuth(InputDownloadAuthInfo input) throws CommonException {
        return null;
    }
}
```

### 3.3.3 资源服务配置

有了资源服务接口和实现类，下面就需要进行资源服务配置，该配置在 **flea_jersey_res_service** 表中，如下所示：

![](flea_jersey_res_service.png)

其中 `flea_jersey_res_service` 的表结构如下：

|   字段名         |    中文描述   |
|:-----------------|:-------------|
|service_code      | 服务编码      |
|resource_code     | 资源编码       |
|service_interfaces |资源服务接口类 |
|service_method    | 资源服务方法   |
|service_input        |  资源服务入参  |
|service_output      |  资源服务出参  |

### 3.3.4 资源服务调用

上述资源服务调用逻辑, 可参考 服务调用过滤器 [InvokeServiceFilter](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-server/src/main/java/com/huazie/fleaframework/jersey/server/filter/impl/InvokeServiceFilter.java)。

`InvokeServiceFilter`, 即资源服务调用过滤器实现，它是**Flea Jersey**接口的核心逻辑。

- 首先，它从请求公共报文中获取**资源编码**【`RESOURCE_CODE`】和**服务编码**【`SERVICE_CODE`】；

- 接着，根据它俩获取相关资源服务配置数据，其中包括资源服务接口、方法、出入参【由服务提供方约定】等；

- 再接着，根据服务接口，从**Web**应用上下文中获取**Spring**注入的服务；

- 然后，从请求业务对象中，取请求业务报文**JSON**串，并转换为资源服务方法的入参对象；

- 再然后，通过反射调用对应的资源服务方法，并获取资源服务方法的出参对象；

- 最后，将出参对象转换成业务返回报文**JSON**串，并添加至响应业务对象中返回。

## 3.4 资源服务业务逻辑开发

如下简单演示了资源服务的业务逻辑开发：

```java
    @Override
    public OutputDownloadAuthInfo downloadAuth(InputDownloadAuthInfo input) throws Exception {
        String fileId = input.getFileId();
        if (StringUtils.isBlank(fileId)) {
            // 入参【{0}】不能为空
            throw new ServiceException("ERROR-SERVICE0000000001", "fileId");
        }

        OutputDownloadAuthInfo output = new OutputDownloadAuthInfo();
        // 演示直接塞了一个随机数
        output.setToken(RandomCode.toUUID());

        return output;
    }
```

## 3.5 国际码和错误码映射

代码中出现 `ERROR-SERVICE0000000001` 的异常，需要配置如下国际码和错误码的映射关系： 

（国际码和错误码映射配置表 `flea_jersey_i18n_error_mapping`）

![](flea_jersey_i18n_error_mapping.png)

当然，我们也可以像**3.2**中，在过滤器配置文件中添加国际码和错误码的映射关系配置。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<jersey>
    <!-- 省略其他 -->

    <!-- 一般这里配置框架层面的异常的国际码和错误码的映射关系，业务的还是配置在 flea_jersey_i18n_error_mapping 表里 -->
    <filter-i18n-error>
        <i18n-error-mapping i18nCode="ERROR-SERVICE0000000001" errorCode="110001">入参【{0}】不能为空</i18n-error-mapping>
    </filter-i18n-error>
</jersey>
```

# 总结

至此，**Flea RESTful**接口服务端接入已经完成。下篇 **Huazie** 将介绍 **Flea RESTful** 接口客户端接入，并以此来调用本篇介绍的下载资源服务，敬请期待！