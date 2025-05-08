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
项目内容可至GitHub 查看 [flea-jersey-client](https://github.com/Huazie/flea-framework/tree/dev/flea-jersey/flea-jersey-client)

相关依赖如下：
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
2. 客户端定义业务入参和业务出参 `POJO` 类；
3. 客户端使用 `FleaJerseyClient` 调用资源服务。

# 3. 具体接入讲解
## 3.1 资源客户端配置

添加资源客户端【`flea_jersey_res_client`】配置, 下载鉴权资源服务的资源客户端配置如下：

![](flea_jersey_res_client.png)

其中 `flea_jersey_res_client` 的表结构如下：

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

**下载鉴权业务输入对象**，定义文件编号入参，代码如下：

```java
@Getter
@Setter
@ToString
public class InputDownloadAuthInfo implements Serializable {

    private static final long serialVersionUID = 6849188299874561970L;

    private String fileId; // 文件编号

}
```

**下载鉴权业务输出对象**，定义下载鉴权令牌出参，代码如下：

```java
@Getter
@Setter
@ToString
public class OutputDownloadAuthInfo implements Serializable {

    private static final long serialVersionUID = 5689920399219551237L;

    private String token; // 下载鉴权令牌

}
```

这里定义的客户端业务入参【`InputFileUploadInfo`】和 出参【`OutputFileUploadInfo`】使用的是服务端定义的资源服务入参和出参；当然这里也可以不一样，原则上只需要客户端业务入参和资源服务入参，客户端业务出参和资源服务出参两边对象转 **JSON** 或 **XML** 的数据内容一致即可。

## 3.3 Flea Jersey客户端接入使用

### 3.3.1 FleaJerseyClient

[FleaJerseyClient](https://github.com/huazie/flea-framework/blob/main/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/bean/FleaJerseyClient.java)，即**Flea Jersey** 客户端，对外提供统一的**Jersey**接口客户端调用**API**。

```java
@Component
public class FleaJerseyClient {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(FleaJerseyClient.class);

    private FleaConfigDataSpringBean springBean;

    @Resource
    public void setSpringBean(FleaConfigDataSpringBean springBean) {
        this.springBean = springBean;
    }

    public <T> Response<T> invoke(String clientCode, Object input, Class<T> outputClazz) throws CommonException {

        Object obj = null;
        if (LOGGER.isDebugEnabled()) {
            obj = new Object() {};
            LOGGER.debug1(obj, "Start");
        }

        // 客户端编码不能为空
        StringUtils.checkBlank(clientCode, FleaJerseyClientException.class, "ERROR-JERSEY-CLIENT0000000001");

        // 业务入参不能为空
        ObjectUtils.checkEmpty(input, FleaJerseyClientException.class, "ERROR-JERSEY-CLIENT0000000002");

        // 业务出参类不能为空
        ObjectUtils.checkEmpty(outputClazz, FleaJerseyClientException.class, "ERROR-JERSEY-CLIENT0000000003");

        // 未注入Bean，直接返回null
        if (ObjectUtils.isEmpty(springBean)) {
            if (LOGGER.isErrorEnabled()) {
                LOGGER.debug1(new Object() {}, "未注入配置数据 Spring Bean，请检查");
            }
            return null;
        }

        // 获取Jersey客户端配置
        FleaJerseyResClient resClient = springBean.getResClient(clientCode);
        // 请检查客户端配置【client_code = {0}】：资源服务客户端未配置
        ObjectUtils.checkEmpty(resClient, FleaJerseyClientException.class, "ERROR-JERSEY-CLIENT0000000009", clientCode);

        RequestConfig config = new RequestConfig();
        // 客户端编码
        config.addClientCode(clientCode);
        // 业务入参对象
        config.addInputObj(input);
        // 资源地址
        config.addResourceUrl(resClient.getResourceUrl());
        // 资源编码
        config.addResourceCode(resClient.getResourceCode());
        // 服务编码
        config.addServiceCode(resClient.getServiceCode());
        // 请求方式
        config.addRequestMode(resClient.getRequestMode());
        // 媒体类型
        config.addMediaType(resClient.getMediaType());
        // 业务入参类全名字符串
        config.addClientInput(resClient.getClientInput());
        // 业务出参类全名字符串
        config.addClientOutput(resClient.getClientOutput());

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(obj, "Request Config = {}", config);
        }

        // 传入请求配置，让请求工厂生产一个Flea Jersey请求
        Request request = RequestFactory.getInstance().buildFleaRequest(config);
        Response<T> response = null;
        if (ObjectUtils.isNotEmpty(request)) {
            response = request.doRequest(outputClazz);
        }

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(obj, "Response = {}", response);
            LOGGER.debug1(obj, "End");
        }

        return response;
    }

}
```

它的 `invoke` 方法实现**调用资源服务**的逻辑如下；

- 首先，根据客户端编码获取 **Flea Jersey** 接口客户端配置【`flea_jersey_res_client`】；

- 然后，根据 **Flea Jersey** 接口客户端配置构建通用的请求配置；

- 接着，传入请求配置，让请求工厂生产一个 **Flea Jersey** 请求；

- 最后，执行 **Flea Jersey** 请求。

### 3.3.2 Post 请求

从 **3.1** 中可以看到下载鉴权资源服务的请求方式是 `Post`，也就是说，在上述请求工厂生产 **Flea Jersey** 请求中，它会生产一个 [Post 请求](https://github.com/huazie/flea-framework/blob/main/flea-jersey/flea-jersey-client/src/main/java/com/huazie/fleaframework/jersey/client/request/impl/PostFleaRequest.java) ，如下：

```java
/**
 * POST 请求，对外提供了执行 POST 请求的能力。
 * <p> 注：服务端提供的资源入口方法需包含 POST 注解。
 *
 * @author huazie
 * @version 1.0.0
 * @since 1.0.0
 */
public class PostFleaRequest extends FleaRequest {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(PostFleaRequest.class);

    /**
     * 默认的构造方法
     *
     * @since 1.0.0
     */
    public PostFleaRequest() {
    }

    /**
     * 带请求配置参数的构造方法
     *
     * @param config 请求配置
     * @since 1.0.0
     */
    public PostFleaRequest(RequestConfig config) {
        super(config);
    }

    @Override
    protected void init() {
        modeEnum = RequestModeEnum.POST;
    }

    @Override
    protected FleaJerseyResponse request(WebTarget target, FleaJerseyRequest request) throws CommonException {

        Object obj = null;
        if (LOGGER.isDebugEnabled()) {
            obj = new Object() {};
            LOGGER.debug1(obj, "POST Request, Start");
        }

        Entity<FleaJerseyRequest> entity = Entity.entity(request, toMediaType());

        FleaJerseyResponse response = target.request(toMediaType()).post(entity, FleaJerseyResponse.class);

        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug1(obj, "POST Request, FleaJerseyResponse = {}", response);
            LOGGER.debug1(obj, "POST Request, End");
        }
        return response;
    }
}
```

### 3.3.3 接入自测

经过**3.1**和**3.2**的步骤，客户端接入已经完成一半，下面就可以通过 `FleaJerseyClient` 调用资源服务，可参考如下：

```java
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath:applicationContext.xml"})
public class JerseyTest {

    private static final FleaLogger LOGGER = FleaLoggerProxy.getProxyInstance(JerseyTest.class);

    @Autowired
    private FleaJerseyClient client;

    @Before
    public void init() {
        IFleaUser fleaUser = new FleaUserImpl();
        fleaUser.setAccountId(10000L);
        fleaUser.set("ACCOUNT_CODE", "13218010892");
        FleaSessionManager.setUserInfo(fleaUser);
    }

    @Test
    public void testDownloadAuth() throws CommonException {
        String clientCode = "FLEA_CLIENT_DOWNLOAD_AUTH";

        InputDownloadAuthInfo downloadAuthInfo = new InputDownloadAuthInfo();
        downloadAuthInfo.setFileId("123123123123123123123");

        Response<OutputDownloadAuthInfo> response = client.invoke(clientCode, downloadAuthInfo, OutputDownloadAuthInfo.class);

        LOGGER.debug("result = {}", response);
    }
}
```

# 总结

至此，Flea RESTful接口客户端接入已经完成。上述自测类，可至GitHub查看 [JerseyTest.java](https://github.com/Huazie/flea-framework/blob/dev/flea-jersey/flea-jersey-client/src/test/java/com/huazie/fleaframework/jersey/client/resource/JerseyTest.java)
